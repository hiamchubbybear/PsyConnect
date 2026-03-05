


















package topology 

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/internal/randutil"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/dns"
)


const (
	topologyDisconnected int64 = iota
	topologyDisconnecting
	topologyConnected
	topologyConnecting
)



var ErrSubscribeAfterClosed = errors.New("cannot subscribe after closeConnection")



var ErrTopologyClosed = errors.New("topology is closed")



var ErrTopologyConnected = errors.New("topology is connected or connecting")



var ErrServerSelectionTimeout = errors.New("server selection timeout")


type MonitorMode uint8


var random = randutil.NewLockedRand()


const (
	AutomaticMode MonitorMode = iota
	SingleMode
)


type Topology struct {
	state int64

	cfg *Config

	desc atomic.Value 

	dnsResolver *dns.Resolver

	done chan struct{}

	pollingRequired   bool
	pollingDone       chan struct{}
	pollingwg         sync.WaitGroup
	rescanSRVInterval time.Duration
	pollHeartbeatTime atomic.Value 

	hosts []string

	updateCallback updateTopologyCallback
	fsm            *fsm

	
	
	
	subscribers         map[uint64]chan description.Topology
	currentSubscriberID uint64
	subscriptionsClosed bool
	subLock             sync.Mutex

	
	
	
	
	serversLock   sync.Mutex
	serversClosed bool
	servers       map[address.Address]*Server

	id primitive.ObjectID
}

var (
	_ driver.Deployment = &Topology{}
	_ driver.Subscriber = &Topology{}
)

type serverSelectionState struct {
	selector    description.ServerSelector
	timeoutChan <-chan time.Time
}

func newServerSelectionState(selector description.ServerSelector, timeoutChan <-chan time.Time) serverSelectionState {
	return serverSelectionState{
		selector:    selector,
		timeoutChan: timeoutChan,
	}
}


func New(cfg *Config) (*Topology, error) {
	if cfg == nil {
		var err error
		cfg, err = NewConfig(options.Client(), nil)
		if err != nil {
			return nil, err
		}
	}

	t := &Topology{
		cfg:               cfg,
		done:              make(chan struct{}),
		pollingDone:       make(chan struct{}),
		rescanSRVInterval: 60 * time.Second,
		fsm:               newFSM(),
		subscribers:       make(map[uint64]chan description.Topology),
		servers:           make(map[address.Address]*Server),
		dnsResolver:       dns.DefaultResolver,
		id:                primitive.NewObjectID(),
	}
	t.desc.Store(description.Topology{})
	t.updateCallback = func(desc description.Server) description.Server {
		return t.apply(context.TODO(), desc)
	}

	if t.cfg.URI != "" {
		connStr, err := connstring.Parse(t.cfg.URI)
		if err != nil {
			return nil, err
		}
		t.pollingRequired = (connStr.Scheme == connstring.SchemeMongoDBSRV) && !t.cfg.LoadBalanced
		t.hosts = connStr.RawHosts
	}

	t.publishTopologyOpeningEvent()

	return t, nil
}

func mustLogTopologyMessage(topo *Topology, level logger.Level) bool {
	return topo.cfg.logger != nil && topo.cfg.logger.LevelComponentEnabled(
		level, logger.ComponentTopology)
}

func logTopologyMessage(topo *Topology, level logger.Level, msg string, keysAndValues ...interface{}) {
	topo.cfg.logger.Print(level,
		logger.ComponentTopology,
		msg,
		logger.SerializeTopology(logger.Topology{
			ID:      topo.id,
			Message: msg,
		}, keysAndValues...)...)
}

func logTopologyThirdPartyUsage(topo *Topology, parsedHosts []string) {
	thirdPartyMessages := [2]string{
		`You appear to be connected to a CosmosDB cluster. For more information regarding feature compatibility and support please visit https://www.mongodb.com/supportability/cosmosdb`,
		`You appear to be connected to a DocumentDB cluster. For more information regarding feature compatibility and support please visit https://www.mongodb.com/supportability/documentdb`,
	}

	thirdPartySuffixes := map[string]int{
		".cosmos.azure.com":            0,
		".docdb.amazonaws.com":         1,
		".docdb-elastic.amazonaws.com": 1,
	}

	hostSet := make([]bool, len(thirdPartyMessages))
	for _, host := range parsedHosts {
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		for suffix, env := range thirdPartySuffixes {
			if !strings.HasSuffix(host, suffix) {
				continue
			}
			if hostSet[env] {
				break
			}
			hostSet[env] = true
			logTopologyMessage(topo, logger.LevelInfo, thirdPartyMessages[env])
		}
	}
}

func mustLogServerSelection(topo *Topology, level logger.Level) bool {
	return topo.cfg.logger != nil && topo.cfg.logger.LevelComponentEnabled(
		level, logger.ComponentServerSelection)
}

func logServerSelection(
	ctx context.Context,
	topo *Topology,
	level logger.Level,
	msg string,
	srvSelector description.ServerSelector,
	keysAndValues ...interface{},
) {
	var srvSelectorString string

	selectorStringer, ok := srvSelector.(fmt.Stringer)
	if ok {
		srvSelectorString = selectorStringer.String()
	}

	operationName, _ := logger.OperationName(ctx)
	operationID, _ := logger.OperationID(ctx)

	topo.cfg.logger.Print(level,
		logger.ComponentServerSelection,
		msg,
		logger.SerializeServerSelection(logger.ServerSelection{
			Selector:            srvSelectorString,
			Operation:           operationName,
			OperationID:         &operationID,
			TopologyDescription: topo.String(),
		}, keysAndValues...)...)
}

func logServerSelectionSucceeded(
	ctx context.Context,
	topo *Topology,
	srvSelector description.ServerSelector,
	server *SelectedServer,
) {
	host, port, err := net.SplitHostPort(server.address.String())
	if err != nil {
		host = server.address.String()
		port = ""
	}

	portInt64, _ := strconv.ParseInt(port, 10, 32)

	logServerSelection(ctx, topo, logger.LevelDebug, logger.ServerSelectionSucceeded, srvSelector,
		logger.KeyServerHost, host,
		logger.KeyServerPort, portInt64)
}

func logServerSelectionFailed(
	ctx context.Context,
	topo *Topology,
	srvSelector description.ServerSelector,
	err error,
) {
	logServerSelection(ctx, topo, logger.LevelDebug, logger.ServerSelectionFailed, srvSelector,
		logger.KeyFailure, err.Error())
}








func logUnexpectedFailure(log *logger.Logger, msg string, callbacks ...func()) {
	r := recover()
	if r == nil {
		return
	}

	defer func() {
		for _, clbk := range callbacks {
			clbk()
		}
	}()

	if log == nil {
		return
	}

	log.Print(logger.LevelInfo, logger.ComponentTopology, fmt.Sprintf("%s: %v", msg, r))
}



func (t *Topology) Connect() error {
	if !atomic.CompareAndSwapInt64(&t.state, topologyDisconnected, topologyConnecting) {
		return ErrTopologyConnected
	}

	t.desc.Store(description.Topology{})
	var err error
	t.serversLock.Lock()

	
	
	if t.cfg.ReplicaSetName != "" {
		t.fsm.SetName = t.cfg.ReplicaSetName
		t.fsm.Kind = description.ReplicaSetNoPrimary
	}

	
	if t.cfg.Mode == SingleMode {
		t.fsm.Kind = description.Single
	}

	for _, a := range t.cfg.SeedList {
		addr := address.Address(a).Canonicalize()
		t.fsm.Servers = append(t.fsm.Servers, description.NewDefaultServer(addr))
	}

	switch {
	case t.cfg.LoadBalanced:
		
		
		
		

		
		t.fsm.Kind = description.LoadBalanced
		t.publishTopologyDescriptionChangedEvent(description.Topology{}, t.fsm.Topology)

		addr := address.Address(t.cfg.SeedList[0]).Canonicalize()
		if err := t.addServer(addr); err != nil {
			t.serversLock.Unlock()
			return err
		}

		
		newServerDesc := t.servers[addr].Description()
		t.publishServerDescriptionChangedEvent(t.fsm.Servers[0], newServerDesc)

		
		oldDesc := t.fsm.Topology
		t.fsm.Servers = []description.Server{newServerDesc}
		t.desc.Store(t.fsm.Topology)
		t.publishTopologyDescriptionChangedEvent(oldDesc, t.fsm.Topology)
	default:
		
		
		
		

		newDesc := description.Topology{
			Kind:                     t.fsm.Kind,
			Servers:                  t.fsm.Servers,
			SessionTimeoutMinutesPtr: t.fsm.SessionTimeoutMinutesPtr,

			
			
			SessionTimeoutMinutes: t.fsm.SessionTimeoutMinutes,
		}
		t.desc.Store(newDesc)
		t.publishTopologyDescriptionChangedEvent(description.Topology{}, t.fsm.Topology)
		for _, a := range t.cfg.SeedList {
			addr := address.Address(a).Canonicalize()
			err = t.addServer(addr)
			if err != nil {
				t.serversLock.Unlock()
				return err
			}
		}
	}

	t.serversLock.Unlock()
	if mustLogTopologyMessage(t, logger.LevelInfo) {
		logTopologyThirdPartyUsage(t, t.hosts)
	}
	if t.pollingRequired {
		
		if len(t.hosts) != 1 {
			return fmt.Errorf("URI with SRV must include one and only one hostname")
		}
		_, _, err = net.SplitHostPort(t.hosts[0])
		if err == nil {
			
			
			return fmt.Errorf("URI with srv must not include a port number")
		}
		go t.pollSRVRecords(t.hosts[0])
		t.pollingwg.Add(1)
	}

	t.subscriptionsClosed = false 

	atomic.StoreInt64(&t.state, topologyConnected)
	return nil
}



func (t *Topology) Disconnect(ctx context.Context) error {
	if !atomic.CompareAndSwapInt64(&t.state, topologyConnected, topologyDisconnecting) {
		return ErrTopologyClosed
	}

	servers := make(map[address.Address]*Server)
	t.serversLock.Lock()
	t.serversClosed = true
	for addr, server := range t.servers {
		servers[addr] = server
	}
	t.serversLock.Unlock()

	for _, server := range servers {
		_ = server.Disconnect(ctx)
		t.publishServerClosedEvent(server.address)
	}

	t.subLock.Lock()
	for id, ch := range t.subscribers {
		close(ch)
		delete(t.subscribers, id)
	}
	t.subscriptionsClosed = true
	t.subLock.Unlock()

	if t.pollingRequired {
		t.pollingDone <- struct{}{}
		t.pollingwg.Wait()
	}

	t.desc.Store(description.Topology{})

	atomic.StoreInt64(&t.state, topologyDisconnected)
	t.publishTopologyClosedEvent()
	return nil
}


func (t *Topology) Description() description.Topology {
	td, ok := t.desc.Load().(description.Topology)
	if !ok {
		td = description.Topology{}
	}
	return td
}


func (t *Topology) Kind() description.TopologyKind { return t.Description().Kind }





func (t *Topology) Subscribe() (*driver.Subscription, error) {
	if atomic.LoadInt64(&t.state) != topologyConnected {
		return nil, errors.New("cannot subscribe to Topology that is not connected")
	}
	ch := make(chan description.Topology, 1)
	td, ok := t.desc.Load().(description.Topology)
	if !ok {
		td = description.Topology{}
	}
	ch <- td

	t.subLock.Lock()
	defer t.subLock.Unlock()
	if t.subscriptionsClosed {
		return nil, ErrSubscribeAfterClosed
	}
	id := t.currentSubscriberID
	t.subscribers[id] = ch
	t.currentSubscriberID++

	return &driver.Subscription{
		Updates: ch,
		ID:      id,
	}, nil
}



func (t *Topology) Unsubscribe(sub *driver.Subscription) error {
	t.subLock.Lock()
	defer t.subLock.Unlock()

	if t.subscriptionsClosed {
		return nil
	}

	ch, ok := t.subscribers[sub.ID]
	if !ok {
		return nil
	}

	close(ch)
	delete(t.subscribers, sub.ID)
	return nil
}



func (t *Topology) RequestImmediateCheck() {
	if atomic.LoadInt64(&t.state) != topologyConnected {
		return
	}
	t.serversLock.Lock()
	for _, server := range t.servers {
		server.RequestImmediateCheck()
	}
	t.serversLock.Unlock()
}




func (t *Topology) SelectServer(ctx context.Context, ss description.ServerSelector) (driver.Server, error) {
	if atomic.LoadInt64(&t.state) != topologyConnected {
		if mustLogServerSelection(t, logger.LevelDebug) {
			logServerSelectionFailed(ctx, t, ss, ErrTopologyClosed)
		}

		return nil, ErrTopologyClosed
	}
	var ssTimeoutCh <-chan time.Time

	if t.cfg.ServerSelectionTimeout > 0 {
		ssTimeout := time.NewTimer(t.cfg.ServerSelectionTimeout)
		ssTimeoutCh = ssTimeout.C
		defer ssTimeout.Stop()
	}

	var doneOnce bool
	var sub *driver.Subscription
	selectionState := newServerSelectionState(ss, ssTimeoutCh)

	
	startTime := time.Now()
	for {
		var suitable []description.Server
		var selectErr error

		if !doneOnce {
			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelection(ctx, t, logger.LevelDebug, logger.ServerSelectionStarted, ss)
			}

			
			
			suitable, selectErr = t.selectServerFromDescription(t.Description(), selectionState)
			doneOnce = true
		} else {
			
			
			if sub == nil {
				var err error
				sub, err = t.Subscribe()
				if err != nil {
					if mustLogServerSelection(t, logger.LevelDebug) {
						logServerSelectionFailed(ctx, t, ss, err)
					}

					return nil, err
				}
				defer func() { _ = t.Unsubscribe(sub) }()
			}

			suitable, selectErr = t.selectServerFromSubscription(ctx, sub.Updates, selectionState)
		}
		if selectErr != nil {
			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelectionFailed(ctx, t, ss, selectErr)
			}

			return nil, selectErr
		}

		if len(suitable) == 0 {
			
			if mustLogServerSelection(t, logger.LevelInfo) {
				elapsed := time.Since(startTime)
				remainingTimeMS := t.cfg.ServerSelectionTimeout - elapsed

				logServerSelection(ctx, t, logger.LevelInfo, logger.ServerSelectionWaiting, ss,
					logger.KeyRemainingTimeMS, remainingTimeMS.Milliseconds())
			}

			continue
		}

		
		
		if len(suitable) == 1 {
			server, err := t.FindServer(suitable[0])
			if err != nil {
				if mustLogServerSelection(t, logger.LevelDebug) {
					logServerSelectionFailed(ctx, t, ss, err)
				}

				return nil, err
			}
			if server == nil {
				continue
			}

			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelectionSucceeded(ctx, t, ss, server)
			}

			return server, nil
		}

		
		
		desc1, desc2 := pick2(suitable)
		server1, err := t.FindServer(desc1)
		if err != nil {
			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelectionFailed(ctx, t, ss, err)
			}

			return nil, err
		}
		server2, err := t.FindServer(desc2)
		if err != nil {
			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelectionFailed(ctx, t, ss, err)
			}

			return nil, err
		}

		
		
		
		
		if server1 == nil || server2 == nil {
			if server1 == nil && server2 == nil {
				continue
			}

			if server1 != nil {
				if mustLogServerSelection(t, logger.LevelDebug) {
					logServerSelectionSucceeded(ctx, t, ss, server1)
				}
				return server1, nil
			}

			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelectionSucceeded(ctx, t, ss, server2)
			}

			return server2, nil
		}

		
		
		
		if server1.OperationCount() < server2.OperationCount() {
			if mustLogServerSelection(t, logger.LevelDebug) {
				logServerSelectionSucceeded(ctx, t, ss, server1)
			}

			return server1, nil
		}

		if mustLogServerSelection(t, logger.LevelDebug) {
			logServerSelectionSucceeded(ctx, t, ss, server2)
		}
		return server2, nil
	}
}





func pick2(ds []description.Server) (description.Server, description.Server) {
	
	idx := random.Intn(len(ds))
	s1 := ds[idx]

	
	
	ds[idx], ds[len(ds)-1] = ds[len(ds)-1], ds[idx]
	ds = ds[:len(ds)-1]

	
	return s1, ds[random.Intn(len(ds))]
}



func (t *Topology) FindServer(selected description.Server) (*SelectedServer, error) {
	if atomic.LoadInt64(&t.state) != topologyConnected {
		return nil, ErrTopologyClosed
	}
	t.serversLock.Lock()
	defer t.serversLock.Unlock()
	server, ok := t.servers[selected.Addr]
	if !ok {
		return nil, nil
	}

	desc := t.Description()
	return &SelectedServer{
		Server: server,
		Kind:   desc.Kind,
	}, nil
}




func (t *Topology) selectServerFromSubscription(ctx context.Context, subscriptionCh <-chan description.Topology,
	selectionState serverSelectionState) ([]description.Server, error) {

	current := t.Description()
	for {
		select {
		case <-ctx.Done():
			return nil, ServerSelectionError{Wrapped: ctx.Err(), Desc: current}
		case <-selectionState.timeoutChan:
			return nil, ServerSelectionError{Wrapped: ErrServerSelectionTimeout, Desc: current}
		case current = <-subscriptionCh:
		}

		suitable, err := t.selectServerFromDescription(current, selectionState)
		if err != nil {
			return nil, err
		}

		if len(suitable) > 0 {
			return suitable, nil
		}
		t.RequestImmediateCheck()
	}
}


func (t *Topology) selectServerFromDescription(desc description.Topology,
	selectionState serverSelectionState) ([]description.Server, error) {

	
	

	if desc.CompatibilityErr != nil {
		return nil, desc.CompatibilityErr
	}

	
	
	
	if desc.Kind == description.LoadBalanced {
		return desc.Servers, nil
	}

	allowedIndexes := make([]int, 0, len(desc.Servers))
	for i, s := range desc.Servers {
		if s.Kind != description.Unknown {
			allowedIndexes = append(allowedIndexes, i)
		}
	}

	allowed := make([]description.Server, len(allowedIndexes))
	for i, idx := range allowedIndexes {
		allowed[i] = desc.Servers[idx]
	}

	suitable, err := selectionState.selector.SelectServer(desc, allowed)
	if err != nil {
		return nil, ServerSelectionError{Wrapped: err, Desc: desc}
	}
	return suitable, nil
}

func (t *Topology) pollSRVRecords(hosts string) {
	defer t.pollingwg.Done()

	serverConfig := newServerConfig(t.cfg.ServerOpts...)
	heartbeatInterval := serverConfig.heartbeatInterval

	pollTicker := time.NewTicker(t.rescanSRVInterval)
	defer pollTicker.Stop()
	t.pollHeartbeatTime.Store(false)
	var doneOnce bool
	defer logUnexpectedFailure(t.cfg.logger, "Encountered unexpected failure polling SRV records", func() {
		if !doneOnce {
			<-t.pollingDone
		}
	})

	for {
		select {
		case <-pollTicker.C:
		case <-t.pollingDone:
			doneOnce = true
			return
		}
		topoKind := t.Description().Kind
		if !(topoKind == description.Unknown || topoKind == description.Sharded) {
			break
		}

		parsedHosts, err := t.dnsResolver.ParseHosts(hosts, t.cfg.SRVServiceName, false)
		
		if err != nil || len(parsedHosts) == 0 {
			if !t.pollHeartbeatTime.Load().(bool) {
				pollTicker.Stop()
				pollTicker = time.NewTicker(heartbeatInterval)
				t.pollHeartbeatTime.Store(true)
			}
			continue
		}
		if t.pollHeartbeatTime.Load().(bool) {
			pollTicker.Stop()
			pollTicker = time.NewTicker(t.rescanSRVInterval)
			t.pollHeartbeatTime.Store(false)
		}

		cont := t.processSRVResults(parsedHosts)
		if !cont {
			break
		}
	}
	<-t.pollingDone
	doneOnce = true
}

func (t *Topology) processSRVResults(parsedHosts []string) bool {
	t.serversLock.Lock()
	defer t.serversLock.Unlock()

	if t.serversClosed {
		return false
	}
	prev := t.fsm.Topology
	diff := diffHostList(t.fsm.Topology, parsedHosts)

	if len(diff.Added) == 0 && len(diff.Removed) == 0 {
		return true
	}

	for _, r := range diff.Removed {
		addr := address.Address(r).Canonicalize()
		s, ok := t.servers[addr]
		if !ok {
			continue
		}
		go func() {
			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()
			_ = s.Disconnect(cancelCtx)
		}()
		delete(t.servers, addr)
		t.fsm.removeServerByAddr(addr)
		t.publishServerClosedEvent(s.address)
	}

	
	
	
	if t.cfg.SRVMaxHosts > 0 && len(t.servers)+len(diff.Added) > t.cfg.SRVMaxHosts {
		random.Shuffle(len(diff.Added), func(i, j int) {
			diff.Added[i], diff.Added[j] = diff.Added[j], diff.Added[i]
		})
	}
	
	for _, a := range diff.Added {
		if t.cfg.SRVMaxHosts > 0 && len(t.servers) >= t.cfg.SRVMaxHosts {
			break
		}
		addr := address.Address(a).Canonicalize()
		_ = t.addServer(addr)
		t.fsm.addServer(addr)
	}

	
	newDesc := description.Topology{
		Kind:                     t.fsm.Kind,
		Servers:                  t.fsm.Servers,
		SessionTimeoutMinutesPtr: t.fsm.SessionTimeoutMinutesPtr,

		
		
		SessionTimeoutMinutes: t.fsm.SessionTimeoutMinutes,
	}
	t.desc.Store(newDesc)

	if !prev.Equal(newDesc) {
		t.publishTopologyDescriptionChangedEvent(prev, newDesc)
	}

	t.subLock.Lock()
	for _, ch := range t.subscribers {
		
		select {
		case <-ch:
		default:
		}
		ch <- newDesc
	}
	t.subLock.Unlock()

	return true
}



func (t *Topology) apply(ctx context.Context, desc description.Server) description.Server {
	t.serversLock.Lock()
	defer t.serversLock.Unlock()

	ind, ok := t.fsm.findServer(desc.Addr)
	if t.serversClosed || !ok {
		return desc
	}

	prev := t.fsm.Topology
	oldDesc := t.fsm.Servers[ind]
	if oldDesc.TopologyVersion.CompareToIncoming(desc.TopologyVersion) > 0 {
		return oldDesc
	}

	var current description.Topology
	current, desc = t.fsm.apply(desc)

	if !oldDesc.Equal(desc) {
		t.publishServerDescriptionChangedEvent(oldDesc, desc)
	}

	diff := diffTopology(prev, current)

	for _, removed := range diff.Removed {
		if s, ok := t.servers[removed.Addr]; ok {
			go func() {
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				_ = s.Disconnect(cancelCtx)
			}()
			delete(t.servers, removed.Addr)
			t.publishServerClosedEvent(s.address)
		}
	}

	for _, added := range diff.Added {
		_ = t.addServer(added.Addr)
	}

	t.desc.Store(current)
	if !prev.Equal(current) {
		t.publishTopologyDescriptionChangedEvent(prev, current)
	}

	t.subLock.Lock()
	for _, ch := range t.subscribers {
		
		select {
		case <-ch:
		default:
		}
		ch <- current
	}
	t.subLock.Unlock()

	return desc
}

func (t *Topology) addServer(addr address.Address) error {
	if _, ok := t.servers[addr]; ok {
		return nil
	}

	svr, err := ConnectServer(addr, t.updateCallback, t.id, t.cfg.ServerOpts...)
	if err != nil {
		return err
	}

	t.servers[addr] = svr

	return nil
}


func (t *Topology) String() string {
	desc := t.Description()

	serversStr := ""
	t.serversLock.Lock()
	defer t.serversLock.Unlock()
	for _, s := range t.servers {
		serversStr += "{ " + s.String() + " }, "
	}
	return fmt.Sprintf("Type: %s, Servers: [%s]", desc.Kind, serversStr)
}


func (t *Topology) publishServerDescriptionChangedEvent(prev description.Server, current description.Server) {
	serverDescriptionChanged := &event.ServerDescriptionChangedEvent{
		Address:             current.Addr,
		TopologyID:          t.id,
		PreviousDescription: prev,
		NewDescription:      current,
	}

	if t.cfg.ServerMonitor != nil && t.cfg.ServerMonitor.ServerDescriptionChanged != nil {
		t.cfg.ServerMonitor.ServerDescriptionChanged(serverDescriptionChanged)
	}
}


func (t *Topology) publishServerClosedEvent(addr address.Address) {
	serverClosed := &event.ServerClosedEvent{
		Address:    addr,
		TopologyID: t.id,
	}

	if t.cfg.ServerMonitor != nil && t.cfg.ServerMonitor.ServerClosed != nil {
		t.cfg.ServerMonitor.ServerClosed(serverClosed)
	}

	if mustLogTopologyMessage(t, logger.LevelDebug) {
		serverHost, serverPort, err := net.SplitHostPort(addr.String())
		if err != nil {
			serverHost = addr.String()
			serverPort = ""
		}

		portInt64, _ := strconv.ParseInt(serverPort, 10, 32)

		logTopologyMessage(t, logger.LevelDebug, logger.TopologyServerClosed,
			logger.KeyServerHost, serverHost,
			logger.KeyServerPort, portInt64)
	}
}


func (t *Topology) publishTopologyDescriptionChangedEvent(prev description.Topology, current description.Topology) {
	topologyDescriptionChanged := &event.TopologyDescriptionChangedEvent{
		TopologyID:          t.id,
		PreviousDescription: prev,
		NewDescription:      current,
	}

	if t.cfg.ServerMonitor != nil && t.cfg.ServerMonitor.TopologyDescriptionChanged != nil {
		t.cfg.ServerMonitor.TopologyDescriptionChanged(topologyDescriptionChanged)
	}

	if mustLogTopologyMessage(t, logger.LevelDebug) {
		logTopologyMessage(t, logger.LevelDebug, logger.TopologyDescriptionChanged,
			logger.KeyPreviousDescription, prev.String(),
			logger.KeyNewDescription, current.String())
	}
}


func (t *Topology) publishTopologyOpeningEvent() {
	topologyOpening := &event.TopologyOpeningEvent{
		TopologyID: t.id,
	}

	if t.cfg.ServerMonitor != nil && t.cfg.ServerMonitor.TopologyOpening != nil {
		t.cfg.ServerMonitor.TopologyOpening(topologyOpening)
	}

	if mustLogTopologyMessage(t, logger.LevelDebug) {
		logTopologyMessage(t, logger.LevelDebug, logger.TopologyOpening)
	}
}


func (t *Topology) publishTopologyClosedEvent() {
	topologyClosed := &event.TopologyClosedEvent{
		TopologyID: t.id,
	}

	if t.cfg.ServerMonitor != nil && t.cfg.ServerMonitor.TopologyClosed != nil {
		t.cfg.ServerMonitor.TopologyClosed(topologyClosed)
	}

	if mustLogTopologyMessage(t, logger.LevelDebug) {
		logTopologyMessage(t, logger.LevelDebug, logger.TopologyClosed)
	}
}
