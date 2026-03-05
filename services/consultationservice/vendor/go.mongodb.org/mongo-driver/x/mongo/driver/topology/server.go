





package topology

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const minHeartbeatInterval = 500 * time.Millisecond
const wireVersion42 = 8 


const (
	serverDisconnected int64 = iota
	serverDisconnecting
	serverConnected
)

func serverStateString(state int64) string {
	switch state {
	case serverDisconnected:
		return "Disconnected"
	case serverDisconnecting:
		return "Disconnecting"
	case serverConnected:
		return "Connected"
	}

	return ""
}

var (
	
	
	ErrServerClosed = errors.New("server is closed")
	
	
	ErrServerConnected = errors.New("server is connected")

	errCheckCancelled = errors.New("server check cancelled")
	emptyDescription  = description.NewDefaultServer("")
)



type SelectedServer struct {
	*Server

	Kind description.TopologyKind
}


func (ss *SelectedServer) Description() description.SelectedServer {
	sdesc := ss.Server.Description()
	return description.SelectedServer{
		Server: sdesc,
		Kind:   ss.Kind,
	}
}


type Server struct {
	
	
	
	

	state          int64
	operationCount int64

	cfg     *serverConfig
	address address.Address

	
	pool *pool

	
	done          chan struct{}
	checkNow      chan struct{}
	disconnecting chan struct{}
	closewg       sync.WaitGroup

	
	desc                   atomic.Value 
	updateTopologyCallback atomic.Value
	topologyID             primitive.ObjectID

	
	subLock             sync.Mutex
	subscribers         map[uint64]chan description.Server
	currentSubscriberID uint64
	subscriptionsClosed bool

	
	
	
	
	heartbeatLock      sync.Mutex
	conn               *connection
	globalCtx          context.Context
	globalCtxCancel    context.CancelFunc
	heartbeatCtx       context.Context
	heartbeatCtxCancel context.CancelFunc

	processErrorLock sync.Mutex
	rttMonitor       *rttMonitor
	monitorOnce      sync.Once
}




type updateTopologyCallback func(description.Server) description.Server



func ConnectServer(
	addr address.Address,
	updateCallback updateTopologyCallback,
	topologyID primitive.ObjectID,
	opts ...ServerOption,
) (*Server, error) {
	srvr := NewServer(addr, topologyID, opts...)
	err := srvr.Connect(updateCallback)
	if err != nil {
		return nil, err
	}
	return srvr, nil
}



func NewServer(addr address.Address, topologyID primitive.ObjectID, opts ...ServerOption) *Server {
	cfg := newServerConfig(opts...)
	globalCtx, globalCtxCancel := context.WithCancel(context.Background())
	s := &Server{
		state: serverDisconnected,

		cfg:     cfg,
		address: addr,

		done:          make(chan struct{}),
		checkNow:      make(chan struct{}, 1),
		disconnecting: make(chan struct{}),

		topologyID: topologyID,

		subscribers:     make(map[uint64]chan description.Server),
		globalCtx:       globalCtx,
		globalCtxCancel: globalCtxCancel,
	}
	s.desc.Store(description.NewDefaultServer(addr))
	rttCfg := &rttConfig{
		interval:           cfg.heartbeatInterval,
		minRTTWindow:       5 * time.Minute,
		createConnectionFn: s.createConnection,
		createOperationFn:  s.createBaseOperation,
	}
	s.rttMonitor = newRTTMonitor(rttCfg)

	pc := poolConfig{
		Address:          addr,
		MinPoolSize:      cfg.minConns,
		MaxPoolSize:      cfg.maxConns,
		MaxConnecting:    cfg.maxConnecting,
		MaxIdleTime:      cfg.poolMaxIdleTime,
		MaintainInterval: cfg.poolMaintainInterval,
		LoadBalanced:     cfg.loadBalanced,
		PoolMonitor:      cfg.poolMonitor,
		Logger:           cfg.logger,
		handshakeErrFn:   s.ProcessHandshakeError,
	}

	connectionOpts := copyConnectionOpts(cfg.connectionOpts)
	s.pool = newPool(pc, connectionOpts...)
	s.publishServerOpeningEvent(s.address)

	return s
}

func mustLogServerMessage(srv *Server) bool {
	return srv.cfg.logger != nil && srv.cfg.logger.LevelComponentEnabled(
		logger.LevelDebug, logger.ComponentTopology)
}

func logServerMessage(srv *Server, msg string, keysAndValues ...interface{}) {
	serverHost, serverPort, err := net.SplitHostPort(srv.address.String())
	if err != nil {
		serverHost = srv.address.String()
		serverPort = ""
	}

	var driverConnectionID uint64
	var serverConnectionID *int64

	if srv.conn != nil {
		driverConnectionID = srv.conn.driverConnectionID
		serverConnectionID = srv.conn.serverConnectionID
	}

	srv.cfg.logger.Print(logger.LevelDebug,
		logger.ComponentTopology,
		msg,
		logger.SerializeServer(logger.Server{
			DriverConnectionID: driverConnectionID,
			TopologyID:         srv.topologyID,
			Message:            msg,
			ServerConnectionID: serverConnectionID,
			ServerHost:         serverHost,
			ServerPort:         serverPort,
		}, keysAndValues...)...)
}



func (s *Server) Connect(updateCallback updateTopologyCallback) error {
	if !atomic.CompareAndSwapInt64(&s.state, serverDisconnected, serverConnected) {
		return ErrServerConnected
	}

	desc := description.NewDefaultServer(s.address)
	if s.cfg.loadBalanced {
		
		desc.Kind = description.LoadBalancer
	}
	s.desc.Store(desc)
	s.updateTopologyCallback.Store(updateCallback)

	if !s.cfg.monitoringDisabled && !s.cfg.loadBalanced {
		s.closewg.Add(1)
		go s.update()
	}

	
	
	
	
	
	
	
	return s.pool.ready()
}










func (s *Server) Disconnect(ctx context.Context) error {
	if !atomic.CompareAndSwapInt64(&s.state, serverConnected, serverDisconnecting) {
		return ErrServerClosed
	}

	s.updateTopologyCallback.Store((updateTopologyCallback)(nil))

	
	
	
	
	s.globalCtxCancel()
	close(s.done)
	s.cancelCheck()

	s.pool.close(ctx)

	s.closewg.Wait()
	s.rttMonitor.disconnect()
	atomic.StoreInt64(&s.state, serverDisconnected)

	return nil
}


func (s *Server) Connection(ctx context.Context) (driver.Connection, error) {
	if atomic.LoadInt64(&s.state) != serverConnected {
		return nil, ErrServerClosed
	}

	
	
	
	atomic.AddInt64(&s.operationCount, 1)
	conn, err := s.pool.checkOut(ctx)
	if err != nil {
		atomic.AddInt64(&s.operationCount, -1)
		return nil, err
	}

	return &Connection{
		connection: conn,
		cleanupServerFn: func() {
			
			
			
			
			
			
			atomic.AddInt64(&s.operationCount, -1)
		},
	}, nil
}



func (s *Server) ProcessHandshakeError(err error, startingGenerationNumber uint64, serviceID *primitive.ObjectID) {
	
	
	
	if err == nil || s.cfg.loadBalanced && serviceID == nil {
		return
	}
	
	if generation, _ := s.pool.generation.getGeneration(serviceID); startingGenerationNumber < generation {
		return
	}

	
	
	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil {
		return
	}

	
	
	
	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	
	
	
	s.updateDescription(description.NewServerFromError(s.address, wrappedConnErr, nil))
	s.pool.clear(err, serviceID)
	s.cancelCheck()
}


func (s *Server) Description() description.Server {
	return s.desc.Load().(description.Server)
}




func (s *Server) SelectedDescription() description.SelectedServer {
	sdesc := s.Description()
	return description.SelectedServer{
		Server: sdesc,
		Kind:   description.Single,
	}
}




func (s *Server) Subscribe() (*ServerSubscription, error) {
	if atomic.LoadInt64(&s.state) != serverConnected {
		return nil, ErrSubscribeAfterClosed
	}
	ch := make(chan description.Server, 1)
	ch <- s.desc.Load().(description.Server)

	s.subLock.Lock()
	defer s.subLock.Unlock()
	if s.subscriptionsClosed {
		return nil, ErrSubscribeAfterClosed
	}
	id := s.currentSubscriberID
	s.subscribers[id] = ch
	s.currentSubscriberID++

	ss := &ServerSubscription{
		C:  ch,
		s:  s,
		id: id,
	}

	return ss, nil
}



func (s *Server) RequestImmediateCheck() {
	select {
	case s.checkNow <- struct{}{}:
	default:
	}
}




func getWriteConcernErrorForProcessing(err error) (*driver.WriteConcernError, bool) {
	var writeCmdErr driver.WriteCommandError
	if !errors.As(err, &writeCmdErr) {
		return nil, false
	}

	wcerr := writeCmdErr.WriteConcernError
	if wcerr != nil && (wcerr.NodeIsRecovering() || wcerr.NotPrimary()) {
		return wcerr, true
	}
	return nil, false
}


func (s *Server) ProcessError(err error, conn driver.Connection) driver.ProcessErrorResult {
	
	if err == nil {
		return driver.NoChange
	}

	
	
	
	
	
	if conn.Stale() {
		return driver.NoChange
	}

	
	
	
	s.processErrorLock.Lock()
	defer s.processErrorLock.Unlock()

	
	
	
	connDesc := conn.Description()
	wireVersion := connDesc.WireVersion
	serviceID := connDesc.ServiceID

	
	
	serverDesc := s.desc.Load().(description.Server)
	topologyVersion := serverDesc.TopologyVersion

	
	
	
	
	
	
	
	
	
	
	
	if tv := connDesc.TopologyVersion; tv != nil && topologyVersion.CompareToIncoming(tv) < 0 {
		topologyVersion = tv
	}

	
	
	if cerr, ok := err.(driver.Error); ok && (cerr.NodeIsRecovering() || cerr.NotPrimary()) {
		
		if topologyVersion.CompareToIncoming(cerr.TopologyVersion) >= 0 {
			return driver.NoChange
		}

		
		s.updateDescription(description.NewServerFromError(s.address, err, cerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		
		if cerr.NodeIsShuttingDown() || wireVersion == nil || wireVersion.Max < wireVersion42 {
			res = driver.ConnectionPoolCleared
			s.pool.clear(err, serviceID)
		}

		return res
	}
	if wcerr, ok := getWriteConcernErrorForProcessing(err); ok {
		
		if topologyVersion.CompareToIncoming(wcerr.TopologyVersion) >= 0 {
			return driver.NoChange
		}

		
		s.updateDescription(description.NewServerFromError(s.address, err, wcerr.TopologyVersion))
		s.RequestImmediateCheck()

		res := driver.ServerMarkedUnknown
		
		if wcerr.NodeIsShuttingDown() || wireVersion == nil || wireVersion.Max < wireVersion42 {
			res = driver.ConnectionPoolCleared
			s.pool.clear(err, serviceID)
		}
		return res
	}

	wrappedConnErr := unwrapConnectionError(err)
	if wrappedConnErr == nil {
		return driver.NoChange
	}

	
	if netErr, ok := wrappedConnErr.(net.Error); ok && netErr.Timeout() {
		return driver.NoChange
	}
	if errors.Is(wrappedConnErr, context.Canceled) || errors.Is(wrappedConnErr, context.DeadlineExceeded) {
		return driver.NoChange
	}

	
	
	
	s.updateDescription(description.NewServerFromError(s.address, err, nil))
	s.pool.clear(err, serviceID)
	s.cancelCheck()
	return driver.ConnectionPoolCleared
}



func (s *Server) update() {
	defer s.closewg.Done()
	heartbeatTicker := time.NewTicker(s.cfg.heartbeatInterval)
	rateLimiter := time.NewTicker(minHeartbeatInterval)
	defer heartbeatTicker.Stop()
	defer rateLimiter.Stop()
	checkNow := s.checkNow
	done := s.done

	defer logUnexpectedFailure(s.cfg.logger, "Encountered unexpected failure updating server")

	closeServer := func() {
		s.subLock.Lock()
		for id, c := range s.subscribers {
			close(c)
			delete(s.subscribers, id)
		}
		s.subscriptionsClosed = true
		s.subLock.Unlock()

		
		
		if s.conn != nil {
			_ = s.conn.close()
		}
	}

	waitUntilNextCheck := func() {
		
		
		select {
		case <-heartbeatTicker.C:
		case <-checkNow:
		case <-done:
			
			return
		}

		
		select {
		case <-rateLimiter.C:
		case <-done:
			return
		}
	}

	timeoutCnt := 0
	for {
		
		
		select {
		case <-done:
			closeServer()
			return
		default:
		}

		previousDescription := s.Description()

		
		desc, err := s.check()
		if errors.Is(err, errCheckCancelled) {
			if atomic.LoadInt64(&s.state) != serverConnected {
				continue
			}

			
			
			waitUntilNextCheck()
			continue
		}

		if isShortcut := func() bool {
			
			
			
			s.processErrorLock.Lock()
			defer s.processErrorLock.Unlock()

			s.updateDescription(desc)
			
			
			if err := unwrapConnectionError(desc.LastError); err != nil && timeoutCnt < 1 {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					timeoutCnt++
					
					return true
				}
				if err, ok := err.(net.Error); ok && err.Timeout() {
					timeoutCnt++
					
					return true
				}
			}
			if err := desc.LastError; err != nil {
				
				
				
				if timeoutCnt > 0 {
					s.pool.clearAll(err, nil)
				} else {
					s.pool.clear(err, nil)
				}
			}
			
			
			
			timeoutCnt = 0
			return false
		}(); isShortcut {
			continue
		}

		
		
		
		connectionIsStreaming := s.conn != nil && s.conn.getCurrentlyStreaming()
		transitionedFromNetworkError := desc.LastError != nil && unwrapConnectionError(desc.LastError) != nil &&
			previousDescription.Kind != description.Unknown

		if isStreamingEnabled(s) && isStreamable(s) {
			s.monitorOnce.Do(s.rttMonitor.connect)
		}

		if isStreamingEnabled(s) && (isStreamable(s) || connectionIsStreaming) || transitionedFromNetworkError {
			continue
		}

		
		
		waitUntilNextCheck()
	}
}





func (s *Server) updateDescription(desc description.Server) {
	if s.cfg.loadBalanced {
		
		
		return
	}

	defer logUnexpectedFailure(s.cfg.logger, "Encountered unexpected failure updating server description")

	
	
	
	
	
	
	
	
	if desc.Kind != description.Unknown {
		_ = s.pool.ready()
	}

	
	callback, ok := s.updateTopologyCallback.Load().(updateTopologyCallback)
	if ok && callback != nil {
		desc = callback(desc)
	}
	s.desc.Store(desc)

	s.subLock.Lock()
	for _, c := range s.subscribers {
		select {
		
		case <-c:
		default:
		}
		c <- desc
	}
	s.subLock.Unlock()
}



func (s *Server) createConnection() *connection {
	opts := copyConnectionOpts(s.cfg.connectionOpts)
	opts = append(opts,
		WithConnectTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		WithReadTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		WithWriteTimeout(func(time.Duration) time.Duration { return s.cfg.heartbeatTimeout }),
		
		
		WithHandshaker(func(Handshaker) Handshaker {
			return operation.NewHello().AppName(s.cfg.appname).Compressors(s.cfg.compressionOpts).
				ServerAPI(s.cfg.serverAPI)
		}),
		
		WithMonitor(func(*event.CommandMonitor) *event.CommandMonitor { return nil }),
	)

	return newConnection(s.address, opts...)
}

func copyConnectionOpts(opts []ConnectionOption) []ConnectionOption {
	optsCopy := make([]ConnectionOption, len(opts))
	copy(optsCopy, opts)
	return optsCopy
}

func (s *Server) setupHeartbeatConnection() error {
	conn := s.createConnection()

	
	s.heartbeatLock.Lock()
	if s.heartbeatCtxCancel != nil {
		
		s.heartbeatCtxCancel()
	}
	s.heartbeatCtx, s.heartbeatCtxCancel = context.WithCancel(s.globalCtx)
	s.conn = conn
	s.heartbeatLock.Unlock()

	return s.conn.connect(s.heartbeatCtx)
}


func (s *Server) cancelCheck() {
	var conn *connection

	
	s.heartbeatLock.Lock()
	if s.heartbeatCtx != nil {
		s.heartbeatCtxCancel()
	}
	conn = s.conn
	s.heartbeatLock.Unlock()

	if conn == nil {
		return
	}

	
	
	
	conn.closeConnectContext()
	conn.wait()
	_ = conn.close()
}

func (s *Server) checkWasCancelled() bool {
	return s.heartbeatCtx.Err() != nil
}

func (s *Server) createBaseOperation(conn driver.Connection) *operation.Hello {
	return operation.
		NewHello().
		ClusterClock(s.cfg.clock).
		Deployment(driver.SingleConnectionDeployment{C: conn}).
		ServerAPI(s.cfg.serverAPI)
}

func isStreamingEnabled(srv *Server) bool {
	switch srv.cfg.serverMonitoringMode {
	case connstring.ServerMonitoringModeStream:
		return true
	case connstring.ServerMonitoringModePoll:
		return false
	default:
		return driverutil.GetFaasEnvName() == ""
	}
}

func isStreamable(srv *Server) bool {
	return srv.Description().Kind != description.Unknown && srv.Description().TopologyVersion != nil
}

func (s *Server) check() (description.Server, error) {
	var descPtr *description.Server
	var err error
	var duration time.Duration

	start := time.Now()

	
	
	if s.conn == nil || s.conn.closed() || s.checkWasCancelled() {
		connID := "0"
		if s.conn != nil {
			connID = s.conn.ID()
		}
		s.publishServerHeartbeatStartedEvent(connID, false)
		
		err = s.setupHeartbeatConnection()
		duration = time.Since(start)
		connID = "0"
		if s.conn != nil {
			connID = s.conn.ID()
		}
		if err == nil {
			
			s.rttMonitor.addSample(s.conn.helloRTT)
			descPtr = &s.conn.desc
			s.publishServerHeartbeatSucceededEvent(connID, duration, s.conn.desc, false)
		} else {
			err = unwrapConnectionError(err)
			s.publishServerHeartbeatFailedEvent(connID, duration, err, false)
		}
	} else {
		

		
		heartbeatConn := initConnection{s.conn}
		baseOperation := s.createBaseOperation(heartbeatConn)
		previousDescription := s.Description()
		streamable := isStreamingEnabled(s) && isStreamable(s)

		s.publishServerHeartbeatStartedEvent(s.conn.ID(), s.conn.getCurrentlyStreaming() || streamable)

		switch {
		case s.conn.getCurrentlyStreaming():
			
			err = baseOperation.StreamResponse(s.heartbeatCtx, heartbeatConn)
		case streamable:
			
			
			

			
			maxAwaitTimeMS := int64(s.cfg.heartbeatInterval) / 1e6
			
			
			
			socketTimeout := s.cfg.heartbeatTimeout
			if socketTimeout != 0 {
				socketTimeout += s.cfg.heartbeatInterval
			}
			s.conn.setSocketTimeout(socketTimeout)
			baseOperation = baseOperation.TopologyVersion(previousDescription.TopologyVersion).
				MaxAwaitTimeMS(maxAwaitTimeMS)
			s.conn.setCanStream(true)
			err = baseOperation.Execute(s.heartbeatCtx)
		default:
			
			

			s.conn.setSocketTimeout(s.cfg.heartbeatTimeout)
			err = baseOperation.Execute(s.heartbeatCtx)
		}

		duration = time.Since(start)

		
		
		
		if !streamable {
			s.rttMonitor.addSample(duration)
		}

		if err == nil {
			tempDesc := baseOperation.Result(s.address)
			descPtr = &tempDesc
			s.publishServerHeartbeatSucceededEvent(s.conn.ID(), duration, tempDesc, s.conn.getCurrentlyStreaming() || streamable)
		} else {
			
			
			if s.conn != nil {
				_ = s.conn.close()
			}
			s.publishServerHeartbeatFailedEvent(s.conn.ID(), duration, err, s.conn.getCurrentlyStreaming() || streamable)
		}
	}

	if descPtr != nil {
		
		desc := *descPtr
		desc = desc.SetAverageRTT(s.rttMonitor.EWMA())
		desc.HeartbeatInterval = s.cfg.heartbeatInterval
		return desc, nil
	}

	if s.checkWasCancelled() {
		
		
		return emptyDescription, errCheckCancelled
	}

	
	
	topologyVersion := extractTopologyVersion(err)
	s.rttMonitor.reset()
	return description.NewServerFromError(s.address, err, topologyVersion), nil
}

func extractTopologyVersion(err error) *description.TopologyVersion {
	if ce, ok := err.(ConnectionError); ok {
		err = ce.Wrapped
	}

	switch converted := err.(type) {
	case driver.Error:
		return converted.TopologyVersion
	case driver.WriteCommandError:
		if converted.WriteConcernError != nil {
			return converted.WriteConcernError.TopologyVersion
		}
	}

	return nil
}


func (s *Server) RTTMonitor() driver.RTTMonitor {
	return s.rttMonitor
}


func (s *Server) OperationCount() int64 {
	return atomic.LoadInt64(&s.operationCount)
}


func (s *Server) String() string {
	desc := s.Description()
	state := atomic.LoadInt64(&s.state)
	str := fmt.Sprintf("Addr: %s, Type: %s, State: %s",
		s.address, desc.Kind, serverStateString(state))
	if len(desc.Tags) != 0 {
		str += fmt.Sprintf(", Tag sets: %s", desc.Tags)
	}
	if state == serverConnected {
		str += fmt.Sprintf(", Average RTT: %s, Min RTT: %s", desc.AverageRTT, s.RTTMonitor().Min())
	}
	if desc.LastError != nil {
		str += fmt.Sprintf(", Last error: %s", desc.LastError)
	}

	return str
}



type ServerSubscription struct {
	C  <-chan description.Server
	s  *Server
	id uint64
}



func (ss *ServerSubscription) Unsubscribe() error {
	ss.s.subLock.Lock()
	defer ss.s.subLock.Unlock()
	if ss.s.subscriptionsClosed {
		return nil
	}

	ch, ok := ss.s.subscribers[ss.id]
	if !ok {
		return nil
	}

	close(ch)
	delete(ss.s.subscribers, ss.id)

	return nil
}


func (s *Server) publishServerOpeningEvent(addr address.Address) {
	if s == nil {
		return
	}

	serverOpening := &event.ServerOpeningEvent{
		Address:    addr,
		TopologyID: s.topologyID,
	}

	if s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerOpening != nil {
		s.cfg.serverMonitor.ServerOpening(serverOpening)
	}

	if mustLogServerMessage(s) {
		logServerMessage(s, logger.TopologyServerOpening)
	}
}


func (s *Server) publishServerHeartbeatStartedEvent(connectionID string, await bool) {
	serverHeartbeatStarted := &event.ServerHeartbeatStartedEvent{
		ConnectionID: connectionID,
		Awaited:      await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatStarted != nil {
		s.cfg.serverMonitor.ServerHeartbeatStarted(serverHeartbeatStarted)
	}

	if mustLogServerMessage(s) {
		logServerMessage(s, logger.TopologyServerHeartbeatStarted,
			logger.KeyAwaited, await)
	}
}


func (s *Server) publishServerHeartbeatSucceededEvent(connectionID string,
	duration time.Duration,
	desc description.Server,
	await bool,
) {
	serverHeartbeatSucceeded := &event.ServerHeartbeatSucceededEvent{
		DurationNanos: duration.Nanoseconds(),
		Duration:      duration,
		Reply:         desc,
		ConnectionID:  connectionID,
		Awaited:       await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatSucceeded != nil {
		s.cfg.serverMonitor.ServerHeartbeatSucceeded(serverHeartbeatSucceeded)
	}

	if mustLogServerMessage(s) {
		descRaw, _ := bson.Marshal(struct {
			description.Server `bson:",inline"`
			Ok                 int32
		}{
			Server: desc,
			Ok: func() int32 {
				if desc.LastError != nil {
					return 0
				}

				return 1
			}(),
		})

		logServerMessage(s, logger.TopologyServerHeartbeatSucceeded,
			logger.KeyAwaited, await,
			logger.KeyDurationMS, duration.Milliseconds(),
			logger.KeyReply, bson.Raw(descRaw).String())
	}
}


func (s *Server) publishServerHeartbeatFailedEvent(connectionID string,
	duration time.Duration,
	err error,
	await bool,
) {
	serverHeartbeatFailed := &event.ServerHeartbeatFailedEvent{
		DurationNanos: duration.Nanoseconds(),
		Duration:      duration,
		Failure:       err,
		ConnectionID:  connectionID,
		Awaited:       await,
	}

	if s != nil && s.cfg.serverMonitor != nil && s.cfg.serverMonitor.ServerHeartbeatFailed != nil {
		s.cfg.serverMonitor.ServerHeartbeatFailed(serverHeartbeatFailed)
	}

	if mustLogServerMessage(s) {
		logServerMessage(s, logger.TopologyServerHeartbeatFailed,
			logger.KeyAwaited, await,
			logger.KeyDurationMS, duration.Milliseconds(),
			logger.KeyFailure, err.Error())
	}
}


func unwrapConnectionError(err error) error {
	
	

	connErr, ok := err.(ConnectionError)
	if ok {
		return connErr.Wrapped
	}

	driverErr, ok := err.(driver.Error)
	if !ok || !driverErr.NetworkError() {
		return nil
	}

	connErr, ok = driverErr.Wrapped.(ConnectionError)
	if ok {
		return connErr.Wrapped
	}

	return nil
}
