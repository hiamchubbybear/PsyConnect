package kafka

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)



var ErrGroupClosed = errors.New("consumer group is closed")



var ErrGenerationEnded = errors.New("consumer group generation has ended")

const (
	
	
	
	
	defaultProtocolType = "consumer"

	
	
	
	
	
	
	defaultHeartbeatInterval = 3 * time.Second

	
	
	defaultSessionTimeout = 30 * time.Second

	
	
	defaultRebalanceTimeout = 30 * time.Second

	
	
	defaultJoinGroupBackoff = 5 * time.Second

	
	
	defaultRetentionTime = -1 * time.Millisecond

	
	
	defaultPartitionWatchTime = 5 * time.Second

	
	
	defaultTimeout = 5 * time.Second
)



type ConsumerGroupConfig struct {
	
	ID string

	
	
	Brokers []string

	
	
	Dialer *Dialer

	
	
	
	Topics []string

	
	
	
	
	
	GroupBalancers []GroupBalancer

	
	
	
	
	HeartbeatInterval time.Duration

	
	
	
	
	
	PartitionWatchInterval time.Duration

	
	
	WatchPartitionChanges bool

	
	
	
	
	SessionTimeout time.Duration

	
	
	
	
	
	RebalanceTimeout time.Duration

	
	
	
	
	JoinGroupBackoff time.Duration

	
	
	
	
	
	
	
	RetentionTime time.Duration

	
	
	
	
	
	StartOffset int64

	
	
	Logger Logger

	
	
	ErrorLogger Logger

	
	
	
	
	
	
	
	
	
	Timeout time.Duration

	
	
	connect func(dialer *Dialer, brokers ...string) (coordinator, error)
}



func (config *ConsumerGroupConfig) Validate() error {

	if len(config.Brokers) == 0 {
		return errors.New("cannot create a consumer group with an empty list of broker addresses")
	}

	if len(config.Topics) == 0 {
		return errors.New("cannot create a consumer group without a topic")
	}

	if config.ID == "" {
		return errors.New("cannot create a consumer group without an ID")
	}

	if config.Dialer == nil {
		config.Dialer = DefaultDialer
	}

	if len(config.GroupBalancers) == 0 {
		config.GroupBalancers = []GroupBalancer{
			RangeGroupBalancer{},
			RoundRobinGroupBalancer{},
		}
	}

	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = defaultHeartbeatInterval
	}

	if config.SessionTimeout == 0 {
		config.SessionTimeout = defaultSessionTimeout
	}

	if config.PartitionWatchInterval == 0 {
		config.PartitionWatchInterval = defaultPartitionWatchTime
	}

	if config.RebalanceTimeout == 0 {
		config.RebalanceTimeout = defaultRebalanceTimeout
	}

	if config.JoinGroupBackoff == 0 {
		config.JoinGroupBackoff = defaultJoinGroupBackoff
	}

	if config.RetentionTime == 0 {
		config.RetentionTime = defaultRetentionTime
	}

	if config.HeartbeatInterval < 0 || (config.HeartbeatInterval/time.Millisecond) >= math.MaxInt32 {
		return fmt.Errorf("HeartbeatInterval out of bounds: %d", config.HeartbeatInterval)
	}

	if config.SessionTimeout < 0 || (config.SessionTimeout/time.Millisecond) >= math.MaxInt32 {
		return fmt.Errorf("SessionTimeout out of bounds: %d", config.SessionTimeout)
	}

	if config.RebalanceTimeout < 0 || (config.RebalanceTimeout/time.Millisecond) >= math.MaxInt32 {
		return fmt.Errorf("RebalanceTimeout out of bounds: %d", config.RebalanceTimeout)
	}

	if config.JoinGroupBackoff < 0 || (config.JoinGroupBackoff/time.Millisecond) >= math.MaxInt32 {
		return fmt.Errorf("JoinGroupBackoff out of bounds: %d", config.JoinGroupBackoff)
	}

	if config.RetentionTime < 0 && config.RetentionTime != defaultRetentionTime {
		return fmt.Errorf("RetentionTime out of bounds: %d", config.RetentionTime)
	}

	if config.PartitionWatchInterval < 0 || (config.PartitionWatchInterval/time.Millisecond) >= math.MaxInt32 {
		return fmt.Errorf("PartitionWachInterval out of bounds %d", config.PartitionWatchInterval)
	}

	if config.StartOffset == 0 {
		config.StartOffset = FirstOffset
	}

	if config.StartOffset != FirstOffset && config.StartOffset != LastOffset {
		return fmt.Errorf("StartOffset is not valid %d", config.StartOffset)
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	if config.connect == nil {
		config.connect = makeConnect(*config)
	}

	return nil
}



type PartitionAssignment struct {
	
	ID int

	
	
	
	
	
	Offset int64
}




type genCtx struct {
	gen *Generation
}

func (c genCtx) Done() <-chan struct{} {
	return c.gen.done
}

func (c genCtx) Err() error {
	select {
	case <-c.gen.done:
		return ErrGenerationEnded
	default:
		return nil
	}
}

func (c genCtx) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c genCtx) Value(interface{}) interface{} {
	return nil
}





type Generation struct {
	
	ID int32

	
	GroupID string

	
	
	MemberID string

	
	
	Assignments map[string][]PartitionAssignment

	conn coordinator

	
	
	
	
	
	
	lock     sync.Mutex
	done     chan struct{}
	closed   bool
	routines int
	joined   chan struct{}

	retentionMillis int64
	log             func(func(Logger))
	logError        func(func(Logger))
}



func (g *Generation) close() {
	g.lock.Lock()
	if !g.closed {
		close(g.done)
		g.closed = true
	}
	
	
	r := g.routines
	g.lock.Unlock()

	
	
	if r > 0 {
		<-g.joined
	}
}














func (g *Generation) Start(fn func(ctx context.Context)) {
	g.lock.Lock()
	defer g.lock.Unlock()

	
	
	
	
	
	
	
	
	if g.closed {
		go fn(genCtx{g})
		return
	}

	
	g.routines++

	go func() {
		fn(genCtx{g})
		g.lock.Lock()
		
		
		
		if !g.closed {
			close(g.done)
			g.closed = true
		}
		g.routines--
		
		
		if g.routines == 0 {
			close(g.joined)
		}
		g.lock.Unlock()
	}()
}




func (g *Generation) CommitOffsets(offsets map[string]map[int]int64) error {
	if len(offsets) == 0 {
		return nil
	}

	topics := make([]offsetCommitRequestV2Topic, 0, len(offsets))
	for topic, partitions := range offsets {
		t := offsetCommitRequestV2Topic{Topic: topic}
		for partition, offset := range partitions {
			t.Partitions = append(t.Partitions, offsetCommitRequestV2Partition{
				Partition: int32(partition),
				Offset:    offset,
			})
		}
		topics = append(topics, t)
	}

	request := offsetCommitRequestV2{
		GroupID:       g.GroupID,
		GenerationID:  g.ID,
		MemberID:      g.MemberID,
		RetentionTime: g.retentionMillis,
		Topics:        topics,
	}

	_, err := g.conn.offsetCommit(request)
	if err == nil {
		
		g.log(func(l Logger) {
			var report []string
			for _, t := range request.Topics {
				report = append(report, fmt.Sprintf("\ttopic: %s", t.Topic))
				for _, p := range t.Partitions {
					report = append(report, fmt.Sprintf("\t\tpartition %d: %d", p.Partition, p.Offset))
				}
			}
			l.Printf("committed offsets for group %s: \n%s", g.GroupID, strings.Join(report, "\n"))
		})
	}

	return err
}




func (g *Generation) heartbeatLoop(interval time.Duration) {
	g.Start(func(ctx context.Context) {
		g.log(func(l Logger) {
			l.Printf("started heartbeat for group, %v [%v]", g.GroupID, interval)
		})
		defer g.log(func(l Logger) {
			l.Printf("stopped heartbeat for group %s\n", g.GroupID)
		})

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, err := g.conn.heartbeat(heartbeatRequestV0{
					GroupID:      g.GroupID,
					GenerationID: g.ID,
					MemberID:     g.MemberID,
				})
				if err != nil {
					return
				}
			}
		}
	})
}







func (g *Generation) partitionWatcher(interval time.Duration, topic string) {
	g.Start(func(ctx context.Context) {
		g.log(func(l Logger) {
			l.Printf("started partition watcher for group, %v, topic %v [%v]", g.GroupID, topic, interval)
		})
		defer g.log(func(l Logger) {
			l.Printf("stopped partition watcher for group, %v, topic %v", g.GroupID, topic)
		})

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		ops, err := g.conn.readPartitions(topic)
		if err != nil {
			g.logError(func(l Logger) {
				l.Printf("Problem getting partitions during startup, %v\n, Returning and setting up nextGeneration", err)
			})
			return
		}
		oParts := len(ops)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ops, err := g.conn.readPartitions(topic)
				switch {
				case err == nil, errors.Is(err, UnknownTopicOrPartition):
					if len(ops) != oParts {
						g.log(func(l Logger) {
							l.Printf("Partition changes found, reblancing group: %v.", g.GroupID)
						})
						return
					}

				default:
					g.logError(func(l Logger) {
						l.Printf("Problem getting partitions while checking for changes, %v", err)
					})
					var kafkaError Error
					if errors.As(err, &kafkaError) {
						continue
					}
					
					
					return
				}
			}
		}
	})
}




type coordinator interface {
	io.Closer
	findCoordinator(findCoordinatorRequestV0) (findCoordinatorResponseV0, error)
	joinGroup(joinGroupRequestV1) (joinGroupResponseV1, error)
	syncGroup(syncGroupRequestV0) (syncGroupResponseV0, error)
	leaveGroup(leaveGroupRequestV0) (leaveGroupResponseV0, error)
	heartbeat(heartbeatRequestV0) (heartbeatResponseV0, error)
	offsetFetch(offsetFetchRequestV1) (offsetFetchResponseV1, error)
	offsetCommit(offsetCommitRequestV2) (offsetCommitResponseV2, error)
	readPartitions(...string) ([]Partition, error)
}







type timeoutCoordinator struct {
	timeout          time.Duration
	sessionTimeout   time.Duration
	rebalanceTimeout time.Duration
	conn             *Conn
}

func (t *timeoutCoordinator) Close() error {
	return t.conn.Close()
}

func (t *timeoutCoordinator) findCoordinator(req findCoordinatorRequestV0) (findCoordinatorResponseV0, error) {
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
		return findCoordinatorResponseV0{}, err
	}
	return t.conn.findCoordinator(req)
}

func (t *timeoutCoordinator) joinGroup(req joinGroupRequestV1) (joinGroupResponseV1, error) {
	
	
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout + t.rebalanceTimeout)); err != nil {
		return joinGroupResponseV1{}, err
	}
	return t.conn.joinGroup(req)
}

func (t *timeoutCoordinator) syncGroup(req syncGroupRequestV0) (syncGroupResponseV0, error) {
	
	
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout + t.sessionTimeout)); err != nil {
		return syncGroupResponseV0{}, err
	}
	return t.conn.syncGroup(req)
}

func (t *timeoutCoordinator) leaveGroup(req leaveGroupRequestV0) (leaveGroupResponseV0, error) {
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
		return leaveGroupResponseV0{}, err
	}
	return t.conn.leaveGroup(req)
}

func (t *timeoutCoordinator) heartbeat(req heartbeatRequestV0) (heartbeatResponseV0, error) {
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
		return heartbeatResponseV0{}, err
	}
	return t.conn.heartbeat(req)
}

func (t *timeoutCoordinator) offsetFetch(req offsetFetchRequestV1) (offsetFetchResponseV1, error) {
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
		return offsetFetchResponseV1{}, err
	}
	return t.conn.offsetFetch(req)
}

func (t *timeoutCoordinator) offsetCommit(req offsetCommitRequestV2) (offsetCommitResponseV2, error) {
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
		return offsetCommitResponseV2{}, err
	}
	return t.conn.offsetCommit(req)
}

func (t *timeoutCoordinator) readPartitions(topics ...string) ([]Partition, error) {
	if err := t.conn.SetDeadline(time.Now().Add(t.timeout)); err != nil {
		return nil, err
	}
	return t.conn.ReadPartitions(topics...)
}





func NewConsumerGroup(config ConsumerGroupConfig) (*ConsumerGroup, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	cg := &ConsumerGroup{
		config: config,
		next:   make(chan *Generation),
		errs:   make(chan error),
		done:   make(chan struct{}),
	}
	cg.wg.Add(1)
	go func() {
		cg.run()
		cg.wg.Done()
	}()
	return cg, nil
}






type ConsumerGroup struct {
	config ConsumerGroupConfig
	next   chan *Generation
	errs   chan error

	closeOnce sync.Once
	wg        sync.WaitGroup
	done      chan struct{}
}




func (cg *ConsumerGroup) Close() error {
	cg.closeOnce.Do(func() {
		close(cg.done)
	})
	cg.wg.Wait()
	return nil
}









func (cg *ConsumerGroup) Next(ctx context.Context) (*Generation, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-cg.done:
		return nil, ErrGroupClosed
	case err := <-cg.errs:
		return nil, err
	case next := <-cg.next:
		return next, nil
	}
}

func (cg *ConsumerGroup) run() {
	
	
	
	
	
	var memberID string
	var err error
	for {
		memberID, err = cg.nextGeneration(memberID)

		
		
		
		var backoff <-chan time.Time

		switch {
		case err == nil:
			
			continue

		case errors.Is(err, ErrGroupClosed):
			
			_ = cg.leaveGroup(memberID)
			return

		case errors.Is(err, RebalanceInProgress):
			
			
			
			
			

		default:
			
			
			
			
			
			_ = cg.leaveGroup(memberID)
			memberID = ""
			backoff = time.After(cg.config.JoinGroupBackoff)
		}
		
		
		select {
		case <-cg.done:
			return
		case cg.errs <- err:
		}
		
		if backoff != nil {
			select {
			case <-cg.done:
				
				return
			case <-backoff:
			}
		}
	}
}

func (cg *ConsumerGroup) nextGeneration(memberID string) (string, error) {
	
	
	
	
	
	
	conn, err := cg.coordinator()
	if err != nil {
		cg.withErrorLogger(func(log Logger) {
			log.Printf("Unable to establish connection to consumer group coordinator for group %s: %v", cg.config.ID, err)
		})
		return memberID, err 
	}
	defer conn.Close()

	var generationID int32
	var groupAssignments GroupMemberAssignments
	var assignments map[string][]int32

	
	
	memberID, generationID, groupAssignments, err = cg.joinGroup(conn, memberID)
	if err != nil {
		cg.withErrorLogger(func(log Logger) {
			log.Printf("Failed to join group %s: %v", cg.config.ID, err)
		})
		return memberID, err
	}
	cg.withLogger(func(log Logger) {
		log.Printf("Joined group %s as member %s in generation %d", cg.config.ID, memberID, generationID)
	})

	
	assignments, err = cg.syncGroup(conn, memberID, generationID, groupAssignments)
	if err != nil {
		cg.withErrorLogger(func(log Logger) {
			log.Printf("Failed to sync group %s: %v", cg.config.ID, err)
		})
		return memberID, err
	}

	
	var offsets map[string]map[int]int64
	offsets, err = cg.fetchOffsets(conn, assignments)
	if err != nil {
		cg.withErrorLogger(func(log Logger) {
			log.Printf("Failed to fetch offsets for group %s: %v", cg.config.ID, err)
		})
		return memberID, err
	}

	
	gen := Generation{
		ID:              generationID,
		GroupID:         cg.config.ID,
		MemberID:        memberID,
		Assignments:     cg.makeAssignments(assignments, offsets),
		conn:            conn,
		done:            make(chan struct{}),
		joined:          make(chan struct{}),
		retentionMillis: int64(cg.config.RetentionTime / time.Millisecond),
		log:             cg.withLogger,
		logError:        cg.withErrorLogger,
	}

	
	
	
	gen.heartbeatLoop(cg.config.HeartbeatInterval)
	if cg.config.WatchPartitionChanges {
		for _, topic := range cg.config.Topics {
			gen.partitionWatcher(cg.config.PartitionWatchInterval, topic)
		}
	}

	
	
	
	
	select {
	case <-cg.done:
		gen.close()
		return memberID, ErrGroupClosed 
	case cg.next <- &gen:
	}

	
	
	select {
	case <-cg.done:
		gen.close()
		return memberID, ErrGroupClosed 
	case <-gen.done:
		
		
		gen.close()
		return memberID, nil
	}
}


func makeConnect(config ConsumerGroupConfig) func(dialer *Dialer, brokers ...string) (coordinator, error) {
	return func(dialer *Dialer, brokers ...string) (coordinator, error) {
		var err error
		for _, broker := range brokers {
			var conn *Conn
			if conn, err = dialer.Dial("tcp", broker); err == nil {
				return &timeoutCoordinator{
					conn:             conn,
					timeout:          config.Timeout,
					sessionTimeout:   config.SessionTimeout,
					rebalanceTimeout: config.RebalanceTimeout,
				}, nil
			}
		}
		return nil, err 
	}
}



func (cg *ConsumerGroup) coordinator() (coordinator, error) {
	
	
	
	
	conn, err := cg.config.connect(cg.config.Dialer, cg.config.Brokers...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	out, err := conn.findCoordinator(findCoordinatorRequestV0{
		CoordinatorKey: cg.config.ID,
	})
	if err == nil && out.ErrorCode != 0 {
		err = Error(out.ErrorCode)
	}
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(out.Coordinator.Host, strconv.Itoa(int(out.Coordinator.Port)))
	return cg.config.connect(cg.config.Dialer, address)
}












func (cg *ConsumerGroup) joinGroup(conn coordinator, memberID string) (string, int32, GroupMemberAssignments, error) {
	request, err := cg.makeJoinGroupRequestV1(memberID)
	if err != nil {
		return "", 0, nil, err
	}

	response, err := conn.joinGroup(request)
	if err == nil && response.ErrorCode != 0 {
		err = Error(response.ErrorCode)
	}
	if err != nil {
		return "", 0, nil, err
	}

	memberID = response.MemberID
	generationID := response.GenerationID

	cg.withLogger(func(l Logger) {
		l.Printf("joined group %s as member %s in generation %d", cg.config.ID, memberID, generationID)
	})

	var assignments GroupMemberAssignments
	if iAmLeader := response.MemberID == response.LeaderID; iAmLeader {
		v, err := cg.assignTopicPartitions(conn, response)
		if err != nil {
			return memberID, 0, nil, err
		}
		assignments = v

		cg.withLogger(func(l Logger) {
			for memberID, assignment := range assignments {
				for topic, partitions := range assignment {
					l.Printf("assigned member/topic/partitions %v/%v/%v", memberID, topic, partitions)
				}
			}
		})
	}

	cg.withLogger(func(l Logger) {
		l.Printf("joinGroup succeeded for response, %v.  generationID=%v, memberID=%v", cg.config.ID, response.GenerationID, response.MemberID)
	})

	return memberID, generationID, assignments, nil
}



func (cg *ConsumerGroup) makeJoinGroupRequestV1(memberID string) (joinGroupRequestV1, error) {
	request := joinGroupRequestV1{
		GroupID:          cg.config.ID,
		MemberID:         memberID,
		SessionTimeout:   int32(cg.config.SessionTimeout / time.Millisecond),
		RebalanceTimeout: int32(cg.config.RebalanceTimeout / time.Millisecond),
		ProtocolType:     defaultProtocolType,
	}

	for _, balancer := range cg.config.GroupBalancers {
		userData, err := balancer.UserData()
		if err != nil {
			return joinGroupRequestV1{}, fmt.Errorf("unable to construct protocol metadata for member, %v: %w", balancer.ProtocolName(), err)
		}
		request.GroupProtocols = append(request.GroupProtocols, joinGroupRequestGroupProtocolV1{
			ProtocolName: balancer.ProtocolName(),
			ProtocolMetadata: groupMetadata{
				Version:  1,
				Topics:   cg.config.Topics,
				UserData: userData,
			}.bytes(),
		})
	}

	return request, nil
}



func (cg *ConsumerGroup) assignTopicPartitions(conn coordinator, group joinGroupResponseV1) (GroupMemberAssignments, error) {
	cg.withLogger(func(l Logger) {
		l.Printf("selected as leader for group, %s\n", cg.config.ID)
	})

	balancer, ok := findGroupBalancer(group.GroupProtocol, cg.config.GroupBalancers)
	if !ok {
		
		
		
		return nil, fmt.Errorf("unable to find selected balancer, %v, for group, %v", group.GroupProtocol, cg.config.ID)
	}

	members, err := cg.makeMemberProtocolMetadata(group.Members)
	if err != nil {
		return nil, err
	}

	topics := extractTopics(members)
	partitions, err := conn.readPartitions(topics...)

	
	
	
	
	if err != nil && !errors.Is(err, UnknownTopicOrPartition) {
		return nil, err
	}

	cg.withLogger(func(l Logger) {
		l.Printf("using '%v' balancer to assign group, %v", group.GroupProtocol, cg.config.ID)
		for _, member := range members {
			l.Printf("found member: %v/%#v", member.ID, member.UserData)
		}
		for _, partition := range partitions {
			l.Printf("found topic/partition: %v/%v", partition.Topic, partition.ID)
		}
	})

	return balancer.AssignGroups(members, partitions), nil
}


func (cg *ConsumerGroup) makeMemberProtocolMetadata(in []joinGroupResponseMemberV1) ([]GroupMember, error) {
	members := make([]GroupMember, 0, len(in))
	for _, item := range in {
		metadata := groupMetadata{}
		reader := bufio.NewReader(bytes.NewReader(item.MemberMetadata))
		if remain, err := (&metadata).readFrom(reader, len(item.MemberMetadata)); err != nil || remain != 0 {
			return nil, fmt.Errorf("unable to read metadata for member, %v: %w", item.MemberID, err)
		}

		members = append(members, GroupMember{
			ID:       item.MemberID,
			Topics:   metadata.Topics,
			UserData: metadata.UserData,
		})
	}
	return members, nil
}











func (cg *ConsumerGroup) syncGroup(conn coordinator, memberID string, generationID int32, memberAssignments GroupMemberAssignments) (map[string][]int32, error) {
	request := cg.makeSyncGroupRequestV0(memberID, generationID, memberAssignments)
	response, err := conn.syncGroup(request)
	if err == nil && response.ErrorCode != 0 {
		err = Error(response.ErrorCode)
	}
	if err != nil {
		return nil, err
	}

	assignments := groupAssignment{}
	reader := bufio.NewReader(bytes.NewReader(response.MemberAssignments))
	if _, err := (&assignments).readFrom(reader, len(response.MemberAssignments)); err != nil {
		return nil, err
	}

	if len(assignments.Topics) == 0 {
		cg.withLogger(func(l Logger) {
			l.Printf("received empty assignments for group, %v as member %s for generation %d", cg.config.ID, memberID, generationID)
		})
	}

	cg.withLogger(func(l Logger) {
		l.Printf("sync group finished for group, %v", cg.config.ID)
	})

	return assignments.Topics, nil
}

func (cg *ConsumerGroup) makeSyncGroupRequestV0(memberID string, generationID int32, memberAssignments GroupMemberAssignments) syncGroupRequestV0 {
	request := syncGroupRequestV0{
		GroupID:      cg.config.ID,
		GenerationID: generationID,
		MemberID:     memberID,
	}

	if memberAssignments != nil {
		request.GroupAssignments = make([]syncGroupRequestGroupAssignmentV0, 0, 1)

		for memberID, topics := range memberAssignments {
			topics32 := make(map[string][]int32)
			for topic, partitions := range topics {
				partitions32 := make([]int32, len(partitions))
				for i := range partitions {
					partitions32[i] = int32(partitions[i])
				}
				topics32[topic] = partitions32
			}
			request.GroupAssignments = append(request.GroupAssignments, syncGroupRequestGroupAssignmentV0{
				MemberID: memberID,
				MemberAssignments: groupAssignment{
					Version: 1,
					Topics:  topics32,
				}.bytes(),
			})
		}

		cg.withLogger(func(logger Logger) {
			logger.Printf("Syncing %d assignments for generation %d as member %s", len(request.GroupAssignments), generationID, memberID)
		})
	}

	return request
}

func (cg *ConsumerGroup) fetchOffsets(conn coordinator, subs map[string][]int32) (map[string]map[int]int64, error) {
	req := offsetFetchRequestV1{
		GroupID: cg.config.ID,
		Topics:  make([]offsetFetchRequestV1Topic, 0, len(cg.config.Topics)),
	}
	for _, topic := range cg.config.Topics {
		req.Topics = append(req.Topics, offsetFetchRequestV1Topic{
			Topic:      topic,
			Partitions: subs[topic],
		})
	}
	offsets, err := conn.offsetFetch(req)
	if err != nil {
		return nil, err
	}

	offsetsByTopic := make(map[string]map[int]int64)
	for _, res := range offsets.Responses {
		offsetsByPartition := map[int]int64{}
		offsetsByTopic[res.Topic] = offsetsByPartition
		for _, pr := range res.PartitionResponses {
			for _, partition := range subs[res.Topic] {
				if partition == pr.Partition {
					offset := pr.Offset
					if offset < 0 {
						offset = cg.config.StartOffset
					}
					offsetsByPartition[int(partition)] = offset
				}
			}
		}
	}

	return offsetsByTopic, nil
}

func (cg *ConsumerGroup) makeAssignments(assignments map[string][]int32, offsets map[string]map[int]int64) map[string][]PartitionAssignment {
	topicAssignments := make(map[string][]PartitionAssignment)
	for _, topic := range cg.config.Topics {
		topicPartitions := assignments[topic]
		topicAssignments[topic] = make([]PartitionAssignment, 0, len(topicPartitions))
		for _, partition := range topicPartitions {
			var offset int64
			partitionOffsets, ok := offsets[topic]
			if ok {
				offset, ok = partitionOffsets[int(partition)]
			}
			if !ok {
				offset = cg.config.StartOffset
			}
			topicAssignments[topic] = append(topicAssignments[topic], PartitionAssignment{
				ID:     int(partition),
				Offset: offset,
			})
		}
	}
	return topicAssignments
}

func (cg *ConsumerGroup) leaveGroup(memberID string) error {
	
	if memberID == "" {
		return nil
	}

	cg.withLogger(func(log Logger) {
		log.Printf("Leaving group %s, member %s", cg.config.ID, memberID)
	})

	
	
	
	
	
	coordinator, err := cg.coordinator()
	if err != nil {
		return err
	}

	_, err = coordinator.leaveGroup(leaveGroupRequestV0{
		GroupID:  cg.config.ID,
		MemberID: memberID,
	})
	if err != nil {
		cg.withErrorLogger(func(log Logger) {
			log.Printf("leave group failed for group, %v, and member, %v: %v", cg.config.ID, memberID, err)
		})
	}

	_ = coordinator.Close()

	return err
}

func (cg *ConsumerGroup) withLogger(do func(Logger)) {
	if cg.config.Logger != nil {
		do(cg.config.Logger)
	}
}

func (cg *ConsumerGroup) withErrorLogger(do func(Logger)) {
	if cg.config.ErrorLogger != nil {
		do(cg.config.ErrorLogger)
	} else {
		cg.withLogger(do)
	}
}
