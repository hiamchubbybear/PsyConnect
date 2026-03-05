package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LastOffset  int64 = -1 
	FirstOffset int64 = -2 
)

const (
	
	
	defaultCommitRetries = 3
)

const (
	
	
	
	defaultFetchMinBytes = 1
)

var (
	errOnlyAvailableWithGroup = errors.New("unavailable when GroupID is not set")
	errNotAvailableWithGroup  = errors.New("unavailable when GroupID is set")
)

const (
	
	
	defaultReadBackoffMin = 100 * time.Millisecond
	defaultReadBackoffMax = 1 * time.Second
)














type Reader struct {
	
	config ReaderConfig

	
	msgs chan readerMessage

	
	mutex   sync.Mutex
	join    sync.WaitGroup
	cancel  context.CancelFunc
	stop    context.CancelFunc
	done    chan struct{}
	commits chan commitRequest
	version int64 
	offset  int64
	lag     int64
	closed  bool

	
	
	
	
	
	
	
	
	runError chan error

	
	once  uint32
	stctx context.Context
	
	
	stats *readerStats
}


func (r *Reader) useConsumerGroup() bool { return r.config.GroupID != "" }

func (r *Reader) getTopics() []string {
	if len(r.config.GroupTopics) > 0 {
		return r.config.GroupTopics[:]
	}

	return []string{r.config.Topic}
}



func (r *Reader) useSyncCommits() bool { return r.config.CommitInterval == 0 }

func (r *Reader) unsubscribe() {
	r.cancel()
	r.join.Wait()
	
	
	
	
	
	
	
	
}

func (r *Reader) subscribe(allAssignments map[string][]PartitionAssignment) {
	offsets := make(map[topicPartition]int64)
	for topic, assignments := range allAssignments {
		for _, assignment := range assignments {
			key := topicPartition{
				topic:     topic,
				partition: int32(assignment.ID),
			}
			offsets[key] = assignment.Offset
		}
	}

	r.mutex.Lock()
	r.start(offsets)
	r.mutex.Unlock()

	r.withLogger(func(l Logger) {
		l.Printf("subscribed to topics and partitions: %+v", offsets)
	})
}



func (r *Reader) commitOffsetsWithRetry(gen *Generation, offsetStash offsetStash, retries int) (err error) {
	const (
		backoffDelayMin = 100 * time.Millisecond
		backoffDelayMax = 5 * time.Second
	)

	for attempt := 0; attempt < retries; attempt++ {
		if attempt != 0 {
			if !sleep(r.stctx, backoff(attempt, backoffDelayMin, backoffDelayMax)) {
				return
			}
		}

		if err = gen.CommitOffsets(offsetStash); err == nil {
			return
		}
	}

	return 
}


type offsetStash map[string]map[int]int64


func (o offsetStash) merge(commits []commit) {
	for _, c := range commits {
		offsetsByPartition, ok := o[c.topic]
		if !ok {
			offsetsByPartition = map[int]int64{}
			o[c.topic] = offsetsByPartition
		}

		if offset, ok := offsetsByPartition[c.partition]; !ok || c.offset > offset {
			offsetsByPartition[c.partition] = c.offset
		}
	}
}


func (o offsetStash) reset() {
	for key := range o {
		delete(o, key)
	}
}


func (r *Reader) commitLoopImmediate(ctx context.Context, gen *Generation) {
	offsets := offsetStash{}

	for {
		select {
		case <-ctx.Done():
			
			
			
			
			var errchs []chan<- error
			for hasCommits := true; hasCommits; {
				select {
				case req := <-r.commits:
					offsets.merge(req.commits)
					errchs = append(errchs, req.errch)
				default:
					hasCommits = false
				}
			}
			err := r.commitOffsetsWithRetry(gen, offsets, defaultCommitRetries)
			for _, errch := range errchs {
				
				errch <- err
			}
			return

		case req := <-r.commits:
			offsets.merge(req.commits)
			req.errch <- r.commitOffsetsWithRetry(gen, offsets, defaultCommitRetries)
			offsets.reset()
		}
	}
}



func (r *Reader) commitLoopInterval(ctx context.Context, gen *Generation) {
	ticker := time.NewTicker(r.config.CommitInterval)
	defer ticker.Stop()

	
	
	offsets := offsetStash{}

	commit := func() {
		if err := r.commitOffsetsWithRetry(gen, offsets, defaultCommitRetries); err != nil {
			r.withErrorLogger(func(l Logger) { l.Printf("%v", err) })
		} else {
			offsets.reset()
		}
	}

	for {
		select {
		case <-ctx.Done():
			
			for hasCommits := true; hasCommits; {
				select {
				case req := <-r.commits:
					offsets.merge(req.commits)
				default:
					hasCommits = false
				}
			}
			commit()
			return

		case <-ticker.C:
			commit()

		case req := <-r.commits:
			offsets.merge(req.commits)
		}
	}
}


func (r *Reader) commitLoop(ctx context.Context, gen *Generation) {
	r.withLogger(func(l Logger) {
		l.Printf("started commit for group %s\n", r.config.GroupID)
	})
	defer r.withLogger(func(l Logger) {
		l.Printf("stopped commit for group %s\n", r.config.GroupID)
	})

	if r.useSyncCommits() {
		r.commitLoopImmediate(ctx, gen)
	} else {
		r.commitLoopInterval(ctx, gen)
	}
}





func (r *Reader) run(cg *ConsumerGroup) {
	defer close(r.done)
	defer cg.Close()

	r.withLogger(func(l Logger) {
		l.Printf("entering loop for consumer group, %v\n", r.config.GroupID)
	})

	for {
		
		
		var err error
		var gen *Generation
		for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
			gen, err = cg.Next(r.stctx)
			if err == nil {
				break
			}
			if errors.Is(err, r.stctx.Err()) {
				return
			}
			r.stats.errors.observe(1)
			r.withErrorLogger(func(l Logger) {
				l.Printf("%v", err)
			})
			
		}
		if err != nil {
			
			select {
			case r.runError <- err:
				
				
			default:
				
			}
			continue
		}

		r.stats.rebalances.observe(1)

		r.subscribe(gen.Assignments)

		gen.Start(func(ctx context.Context) {
			r.commitLoop(ctx, gen)
		})
		gen.Start(func(ctx context.Context) {
			
			select {
			case <-ctx.Done():
				
			case <-r.stctx.Done():
				
			}
			r.unsubscribe()
		})
	}
}



type ReaderConfig struct {
	
	Brokers []string

	
	
	GroupID string

	
	
	
	GroupTopics []string

	
	Topic string

	
	
	Partition int

	
	
	Dialer *Dialer

	
	
	QueueCapacity int

	
	
	
	
	
	
	MinBytes int

	
	
	
	
	
	MaxBytes int

	
	
	
	
	MaxWait time.Duration

	
	
	
	ReadBatchTimeout time.Duration

	
	
	ReadLagInterval time.Duration

	
	
	
	
	
	
	
	GroupBalancers []GroupBalancer

	
	
	
	
	
	
	HeartbeatInterval time.Duration

	
	
	
	
	
	
	CommitInterval time.Duration

	
	
	
	
	
	
	
	PartitionWatchInterval time.Duration

	
	
	WatchPartitionChanges bool

	
	
	
	
	
	
	SessionTimeout time.Duration

	
	
	
	
	
	
	
	RebalanceTimeout time.Duration

	
	
	
	
	JoinGroupBackoff time.Duration

	
	
	
	
	
	
	RetentionTime time.Duration

	
	
	
	
	
	
	
	StartOffset int64

	
	
	
	
	ReadBackoffMin time.Duration

	
	
	
	
	ReadBackoffMax time.Duration

	
	
	Logger Logger

	
	
	ErrorLogger Logger

	
	
	
	IsolationLevel IsolationLevel

	
	
	
	MaxAttempts int

	
	
	
	
	OffsetOutOfRangeError bool
}


func (config *ReaderConfig) Validate() error {
	if len(config.Brokers) == 0 {
		return errors.New("cannot create a new kafka reader with an empty list of broker addresses")
	}

	if config.Partition < 0 || config.Partition >= math.MaxInt32 {
		return fmt.Errorf("partition number out of bounds: %d", config.Partition)
	}

	if config.MinBytes < 0 {
		return fmt.Errorf("invalid negative minimum batch size (min = %d)", config.MinBytes)
	}

	if config.MaxBytes < 0 {
		return fmt.Errorf("invalid negative maximum batch size (max = %d)", config.MaxBytes)
	}

	if config.GroupID != "" {
		if config.Partition != 0 {
			return errors.New("either Partition or GroupID may be specified, but not both")
		}

		if len(config.Topic) == 0 && len(config.GroupTopics) == 0 {
			return errors.New("either Topic or GroupTopics must be specified with GroupID")
		}
	} else if len(config.Topic) == 0 {
		return errors.New("cannot create a new kafka reader with an empty topic")
	}

	if config.MinBytes > config.MaxBytes {
		return fmt.Errorf("minimum batch size greater than the maximum (min = %d, max = %d)", config.MinBytes, config.MaxBytes)
	}

	if config.ReadBackoffMax < 0 {
		return fmt.Errorf("ReadBackoffMax out of bounds: %d", config.ReadBackoffMax)
	}

	if config.ReadBackoffMin < 0 {
		return fmt.Errorf("ReadBackoffMin out of bounds: %d", config.ReadBackoffMin)
	}

	return nil
}



type ReaderStats struct {
	Dials      int64 `metric:"kafka.reader.dial.count"      type:"counter"`
	Fetches    int64 `metric:"kafka.reader.fetch.count"     type:"counter"`
	Messages   int64 `metric:"kafka.reader.message.count"   type:"counter"`
	Bytes      int64 `metric:"kafka.reader.message.bytes"   type:"counter"`
	Rebalances int64 `metric:"kafka.reader.rebalance.count" type:"counter"`
	Timeouts   int64 `metric:"kafka.reader.timeout.count"   type:"counter"`
	Errors     int64 `metric:"kafka.reader.error.count"     type:"counter"`

	DialTime   DurationStats `metric:"kafka.reader.dial.seconds"`
	ReadTime   DurationStats `metric:"kafka.reader.read.seconds"`
	WaitTime   DurationStats `metric:"kafka.reader.wait.seconds"`
	FetchSize  SummaryStats  `metric:"kafka.reader.fetch.size"`
	FetchBytes SummaryStats  `metric:"kafka.reader.fetch.bytes"`

	Offset        int64         `metric:"kafka.reader.offset"          type:"gauge"`
	Lag           int64         `metric:"kafka.reader.lag"             type:"gauge"`
	MinBytes      int64         `metric:"kafka.reader.fetch_bytes.min" type:"gauge"`
	MaxBytes      int64         `metric:"kafka.reader.fetch_bytes.max" type:"gauge"`
	MaxWait       time.Duration `metric:"kafka.reader.fetch_wait.max"  type:"gauge"`
	QueueLength   int64         `metric:"kafka.reader.queue.length"    type:"gauge"`
	QueueCapacity int64         `metric:"kafka.reader.queue.capacity"  type:"gauge"`

	ClientID  string `tag:"client_id"`
	Topic     string `tag:"topic"`
	Partition string `tag:"partition"`

	
	
	
	
	DeprecatedFetchesWithTypo int64 `metric:"kafak.reader.fetch.count" type:"counter"`
}


type readerStats struct {
	dials      counter
	fetches    counter
	messages   counter
	bytes      counter
	rebalances counter
	timeouts   counter
	errors     counter
	dialTime   summary
	readTime   summary
	waitTime   summary
	fetchSize  summary
	fetchBytes summary
	offset     gauge
	lag        gauge
	partition  string
}



func NewReader(config ReaderConfig) *Reader {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	if config.GroupID != "" {
		if len(config.GroupBalancers) == 0 {
			config.GroupBalancers = []GroupBalancer{
				RangeGroupBalancer{},
				RoundRobinGroupBalancer{},
			}
		}
	}

	if config.Dialer == nil {
		config.Dialer = DefaultDialer
	}

	if config.MaxBytes == 0 {
		config.MaxBytes = 1e6 
	}

	if config.MinBytes == 0 {
		config.MinBytes = defaultFetchMinBytes
	}

	if config.MaxWait == 0 {
		config.MaxWait = 10 * time.Second
	}

	if config.ReadBatchTimeout == 0 {
		config.ReadBatchTimeout = 10 * time.Second
	}

	if config.ReadLagInterval == 0 {
		config.ReadLagInterval = 1 * time.Minute
	}

	if config.ReadBackoffMin == 0 {
		config.ReadBackoffMin = defaultReadBackoffMin
	}

	if config.ReadBackoffMax == 0 {
		config.ReadBackoffMax = defaultReadBackoffMax
	}

	if config.ReadBackoffMax < config.ReadBackoffMin {
		panic(fmt.Errorf("ReadBackoffMax %d smaller than ReadBackoffMin %d", config.ReadBackoffMax, config.ReadBackoffMin))
	}

	if config.QueueCapacity == 0 {
		config.QueueCapacity = 100
	}

	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}

	
	readerStatsPartition := config.Partition
	if config.GroupID != "" {
		readerStatsPartition = -1
	}

	
	
	version := int64(0)
	if config.GroupID != "" {
		version = 1
	}

	stctx, stop := context.WithCancel(context.Background())
	r := &Reader{
		config:  config,
		msgs:    make(chan readerMessage, config.QueueCapacity),
		cancel:  func() {},
		commits: make(chan commitRequest, config.QueueCapacity),
		stop:    stop,
		offset:  FirstOffset,
		stctx:   stctx,
		stats: &readerStats{
			dialTime:   makeSummary(),
			readTime:   makeSummary(),
			waitTime:   makeSummary(),
			fetchSize:  makeSummary(),
			fetchBytes: makeSummary(),
			
			
			partition: strconv.Itoa(readerStatsPartition),
		},
		version: version,
	}
	if r.useConsumerGroup() {
		r.done = make(chan struct{})
		r.runError = make(chan error)
		cg, err := NewConsumerGroup(ConsumerGroupConfig{
			ID:                     r.config.GroupID,
			Brokers:                r.config.Brokers,
			Dialer:                 r.config.Dialer,
			Topics:                 r.getTopics(),
			GroupBalancers:         r.config.GroupBalancers,
			HeartbeatInterval:      r.config.HeartbeatInterval,
			PartitionWatchInterval: r.config.PartitionWatchInterval,
			WatchPartitionChanges:  r.config.WatchPartitionChanges,
			SessionTimeout:         r.config.SessionTimeout,
			RebalanceTimeout:       r.config.RebalanceTimeout,
			JoinGroupBackoff:       r.config.JoinGroupBackoff,
			RetentionTime:          r.config.RetentionTime,
			StartOffset:            r.config.StartOffset,
			Logger:                 r.config.Logger,
			ErrorLogger:            r.config.ErrorLogger,
		})
		if err != nil {
			panic(err)
		}
		go r.run(cg)
	}

	return r
}


func (r *Reader) Config() ReaderConfig {
	return r.config
}



func (r *Reader) Close() error {
	atomic.StoreUint32(&r.once, 1)

	r.mutex.Lock()
	closed := r.closed
	r.closed = true
	r.mutex.Unlock()

	r.cancel()
	r.stop()
	r.join.Wait()

	if r.done != nil {
		<-r.done
	}

	if !closed {
		close(r.msgs)
	}

	return nil
}













func (r *Reader) ReadMessage(ctx context.Context) (Message, error) {
	m, err := r.FetchMessage(ctx)
	if err != nil {
		return Message{}, fmt.Errorf("fetching message: %w", err)
	}

	if r.useConsumerGroup() {
		if err := r.CommitMessages(ctx, m); err != nil {
			return Message{}, fmt.Errorf("committing message: %w", err)
		}
	}

	return m, nil
}









func (r *Reader) FetchMessage(ctx context.Context) (Message, error) {
	r.activateReadLag()

	for {
		r.mutex.Lock()

		if !r.closed && r.version == 0 {
			r.start(r.getTopicPartitionOffset())
		}

		version := r.version
		r.mutex.Unlock()

		select {
		case <-ctx.Done():
			return Message{}, ctx.Err()

		case err := <-r.runError:
			return Message{}, err

		case m, ok := <-r.msgs:
			if !ok {
				return Message{}, io.EOF
			}

			if m.version >= version {
				r.mutex.Lock()

				switch {
				case m.error != nil:
				case version == r.version:
					r.offset = m.message.Offset + 1
					r.lag = m.watermark - r.offset
				}

				r.mutex.Unlock()

				if errors.Is(m.error, io.EOF) {
					
					
					
					
					m.error = io.ErrUnexpectedEOF
				}

				return m.message, m.error
			}
		}
	}
}













func (r *Reader) CommitMessages(ctx context.Context, msgs ...Message) error {
	if !r.useConsumerGroup() {
		return errOnlyAvailableWithGroup
	}

	var errch <-chan error
	creq := commitRequest{
		commits: makeCommits(msgs...),
	}

	if r.useSyncCommits() {
		ch := make(chan error, 1)
		errch, creq.errch = ch, ch
	}

	select {
	case r.commits <- creq:
	case <-ctx.Done():
		return ctx.Err()
	case <-r.stctx.Done():
		
		
		return io.ErrClosedPipe
	}

	if !r.useSyncCommits() {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errch:
		return err
	}
}












func (r *Reader) ReadLag(ctx context.Context) (lag int64, err error) {
	if r.useConsumerGroup() {
		return 0, errNotAvailableWithGroup
	}

	type offsets struct {
		first int64
		last  int64
	}

	offch := make(chan offsets, 1)
	errch := make(chan error, 1)

	go func() {
		var off offsets
		var err error

		for _, broker := range r.config.Brokers {
			var conn *Conn

			if conn, err = r.config.Dialer.DialLeader(ctx, "tcp", broker, r.config.Topic, r.config.Partition); err != nil {
				continue
			}

			deadline, _ := ctx.Deadline()
			conn.SetDeadline(deadline)

			off.first, off.last, err = conn.ReadOffsets()
			conn.Close()

			if err == nil {
				break
			}
		}

		if err != nil {
			errch <- err
		} else {
			offch <- off
		}
	}()

	select {
	case off := <-offch:
		switch cur := r.Offset(); {
		case cur == FirstOffset:
			lag = off.last - off.first

		case cur == LastOffset:
			lag = 0

		default:
			lag = off.last - cur
		}
	case err = <-errch:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}



func (r *Reader) Offset() int64 {
	if r.useConsumerGroup() {
		return -1
	}

	r.mutex.Lock()
	offset := r.offset
	r.mutex.Unlock()
	r.withLogger(func(log Logger) {
		log.Printf("looking up offset of kafka reader for partition %d of %s: %s", r.config.Partition, r.config.Topic, toHumanOffset(offset))
	})
	return offset
}



func (r *Reader) Lag() int64 {
	if r.useConsumerGroup() {
		return -1
	}

	r.mutex.Lock()
	lag := r.lag
	r.mutex.Unlock()
	return lag
}









func (r *Reader) SetOffset(offset int64) error {
	if r.useConsumerGroup() {
		return errNotAvailableWithGroup
	}

	var err error
	r.mutex.Lock()

	if r.closed {
		err = io.ErrClosedPipe
	} else if offset != r.offset {
		r.withLogger(func(log Logger) {
			log.Printf("setting the offset of the kafka reader for partition %d of %s from %s to %s",
				r.config.Partition, r.config.Topic, toHumanOffset(r.offset), toHumanOffset(offset))
		})
		r.offset = offset

		if r.version != 0 {
			r.start(r.getTopicPartitionOffset())
		}

		r.activateReadLag()
	}

	r.mutex.Unlock()
	return err
}






func (r *Reader) SetOffsetAt(ctx context.Context, t time.Time) error {
	r.mutex.Lock()
	if r.closed {
		r.mutex.Unlock()
		return io.ErrClosedPipe
	}
	r.mutex.Unlock()

	if len(r.config.Brokers) < 1 {
		return errors.New("no brokers in config")
	}
	var conn *Conn
	var err error
	for _, broker := range r.config.Brokers {
		conn, err = r.config.Dialer.DialLeader(ctx, "tcp", broker, r.config.Topic, r.config.Partition)
		if err != nil {
			continue
		}
		deadline, _ := ctx.Deadline()
		conn.SetDeadline(deadline)
		offset, err := conn.ReadOffset(t)
		conn.Close()
		if err != nil {
			return err
		}

		return r.SetOffset(offset)
	}
	return fmt.Errorf("error dialing all brokers, one of the errors: %w", err)
}








func (r *Reader) Stats() ReaderStats {
	stats := ReaderStats{
		Dials:         r.stats.dials.snapshot(),
		Fetches:       r.stats.fetches.snapshot(),
		Messages:      r.stats.messages.snapshot(),
		Bytes:         r.stats.bytes.snapshot(),
		Rebalances:    r.stats.rebalances.snapshot(),
		Timeouts:      r.stats.timeouts.snapshot(),
		Errors:        r.stats.errors.snapshot(),
		DialTime:      r.stats.dialTime.snapshotDuration(),
		ReadTime:      r.stats.readTime.snapshotDuration(),
		WaitTime:      r.stats.waitTime.snapshotDuration(),
		FetchSize:     r.stats.fetchSize.snapshot(),
		FetchBytes:    r.stats.fetchBytes.snapshot(),
		Offset:        r.stats.offset.snapshot(),
		Lag:           r.stats.lag.snapshot(),
		MinBytes:      int64(r.config.MinBytes),
		MaxBytes:      int64(r.config.MaxBytes),
		MaxWait:       r.config.MaxWait,
		QueueLength:   int64(len(r.msgs)),
		QueueCapacity: int64(cap(r.msgs)),
		ClientID:      r.config.Dialer.ClientID,
		Topic:         r.config.Topic,
		Partition:     r.stats.partition,
	}
	
	stats.DeprecatedFetchesWithTypo = stats.Fetches
	return stats
}

func (r *Reader) getTopicPartitionOffset() map[topicPartition]int64 {
	key := topicPartition{topic: r.config.Topic, partition: int32(r.config.Partition)}
	return map[topicPartition]int64{key: r.offset}
}

func (r *Reader) withLogger(do func(Logger)) {
	if r.config.Logger != nil {
		do(r.config.Logger)
	}
}

func (r *Reader) withErrorLogger(do func(Logger)) {
	if r.config.ErrorLogger != nil {
		do(r.config.ErrorLogger)
	} else {
		r.withLogger(do)
	}
}

func (r *Reader) activateReadLag() {
	if r.config.ReadLagInterval > 0 && atomic.CompareAndSwapUint32(&r.once, 0, 1) {
		
		
		if !r.useConsumerGroup() {
			go r.readLag(r.stctx)
		}
	}
}

func (r *Reader) readLag(ctx context.Context) {
	ticker := time.NewTicker(r.config.ReadLagInterval)
	defer ticker.Stop()

	for {
		timeout, cancel := context.WithTimeout(ctx, r.config.ReadLagInterval/2)
		lag, err := r.ReadLag(timeout)
		cancel()

		if err != nil {
			r.stats.errors.observe(1)
			r.withErrorLogger(func(log Logger) {
				log.Printf("kafka reader failed to read lag of partition %d of %s: %s", r.config.Partition, r.config.Topic, err)
			})
		} else {
			r.stats.lag.observe(lag)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (r *Reader) start(offsetsByPartition map[topicPartition]int64) {
	if r.closed {
		
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	r.cancel() 
	r.cancel = cancel
	r.version++

	r.join.Add(len(offsetsByPartition))
	for key, offset := range offsetsByPartition {
		go func(ctx context.Context, key topicPartition, offset int64, join *sync.WaitGroup) {
			defer join.Done()

			(&reader{
				dialer:           r.config.Dialer,
				logger:           r.config.Logger,
				errorLogger:      r.config.ErrorLogger,
				brokers:          r.config.Brokers,
				topic:            key.topic,
				partition:        int(key.partition),
				minBytes:         r.config.MinBytes,
				maxBytes:         r.config.MaxBytes,
				maxWait:          r.config.MaxWait,
				readBatchTimeout: r.config.ReadBatchTimeout,
				backoffDelayMin:  r.config.ReadBackoffMin,
				backoffDelayMax:  r.config.ReadBackoffMax,
				version:          r.version,
				msgs:             r.msgs,
				stats:            r.stats,
				isolationLevel:   r.config.IsolationLevel,
				maxAttempts:      r.config.MaxAttempts,

				
				offsetOutOfRangeError: r.config.OffsetOutOfRangeError,
			}).run(ctx, offset)
		}(ctx, key, offset, &r.join)
	}
}




type reader struct {
	dialer           *Dialer
	logger           Logger
	errorLogger      Logger
	brokers          []string
	topic            string
	partition        int
	minBytes         int
	maxBytes         int
	maxWait          time.Duration
	readBatchTimeout time.Duration
	backoffDelayMin  time.Duration
	backoffDelayMax  time.Duration
	version          int64
	msgs             chan<- readerMessage
	stats            *readerStats
	isolationLevel   IsolationLevel
	maxAttempts      int

	offsetOutOfRangeError bool
}

type readerMessage struct {
	version   int64
	message   Message
	watermark int64
	error     error
}

func (r *reader) run(ctx context.Context, offset int64) {
	
	
	
	
	
	
	
	
	
	for attempt := 0; true; attempt++ {
		if attempt != 0 {
			if !sleep(ctx, backoff(attempt, r.backoffDelayMin, r.backoffDelayMax)) {
				return
			}
		}

		r.withLogger(func(log Logger) {
			log.Printf("initializing kafka reader for partition %d of %s starting at offset %d", r.partition, r.topic, toHumanOffset(offset))
		})

		conn, start, err := r.initialize(ctx, offset)
		if err != nil {
			if errors.Is(err, OffsetOutOfRange) {
				if r.offsetOutOfRangeError {
					r.sendError(ctx, err)
					return
				}

				
				
				
				r.withErrorLogger(func(log Logger) {
					log.Printf("error initializing the kafka reader for partition %d of %s: %s", r.partition, r.topic, err)
				})

				continue
			}

			
			
			
			
			if attempt >= r.maxAttempts {
				r.sendError(ctx, err)
			} else {
				r.stats.errors.observe(1)
				r.withErrorLogger(func(log Logger) {
					log.Printf("error initializing the kafka reader for partition %d of %s: %s", r.partition, r.topic, err)
				})
			}
			continue
		}

		
		
		
		attempt = 0

		
		
		offset = start

		errcount := 0
	readLoop:
		for {
			if !sleep(ctx, backoff(errcount, r.backoffDelayMin, r.backoffDelayMax)) {
				conn.Close()
				return
			}

			offset, err = r.read(ctx, offset, conn)
			switch {
			case err == nil:
				errcount = 0
				continue

			case errors.Is(err, io.EOF):
				
				
				
				
				errcount = 0
				continue

			case errors.Is(err, io.ErrNoProgress):
				
				
				
				
				conn.Close()
				break readLoop

			case errors.Is(err, UnknownTopicOrPartition):
				r.withErrorLogger(func(log Logger) {
					log.Printf("failed to read from current broker %v for partition %d of %s at offset %d: %v", r.brokers, r.partition, r.topic, toHumanOffset(offset), err)
				})

				conn.Close()

				
				
				r.stats.rebalances.observe(1)
				break readLoop

			case errors.Is(err, NotLeaderForPartition):
				r.withErrorLogger(func(log Logger) {
					log.Printf("failed to read from current broker for partition %d of %s at offset %d: %v", r.partition, r.topic, toHumanOffset(offset), err)
				})

				conn.Close()

				
				
				r.stats.rebalances.observe(1)
				break readLoop

			case errors.Is(err, RequestTimedOut):
				
				errcount = 0
				r.withLogger(func(log Logger) {
					log.Printf("no messages received from kafka within the allocated time for partition %d of %s at offset %d: %v", r.partition, r.topic, toHumanOffset(offset), err)
				})
				r.stats.timeouts.observe(1)
				continue

			case errors.Is(err, OffsetOutOfRange):
				first, last, err := r.readOffsets(conn)
				if err != nil {
					r.withErrorLogger(func(log Logger) {
						log.Printf("the kafka reader got an error while attempting to determine whether it was reading before the first offset or after the last offset of partition %d of %s: %s", r.partition, r.topic, err)
					})
					conn.Close()
					break readLoop
				}

				switch {
				case offset < first:
					r.withErrorLogger(func(log Logger) {
						log.Printf("the kafka reader is reading before the first offset for partition %d of %s, skipping from offset %d to %d (%d messages)", r.partition, r.topic, toHumanOffset(offset), first, first-offset)
					})
					offset, errcount = first, 0
					continue 

				case offset < last:
					errcount = 0
					continue 

				default:
					
					r.withErrorLogger(func(log Logger) {
						log.Printf("the kafka reader is reading passed the last offset for partition %d of %s at offset %d", r.partition, r.topic, toHumanOffset(offset))
					})
				}

			case errors.Is(err, context.Canceled):
				
				conn.Close()
				return

			case errors.Is(err, errUnknownCodec):
				
				
				
				r.sendError(ctx, err)
				break readLoop

			default:
				var kafkaError Error
				if errors.As(err, &kafkaError) {
					r.sendError(ctx, err)
				} else {
					r.withErrorLogger(func(log Logger) {
						log.Printf("the kafka reader got an unknown error reading partition %d of %s at offset %d: %s", r.partition, r.topic, toHumanOffset(offset), err)
					})
					r.stats.errors.observe(1)
					conn.Close()
					break readLoop
				}
			}

			errcount++
		}
	}
}

func (r *reader) initialize(ctx context.Context, offset int64) (conn *Conn, start int64, err error) {
	for i := 0; i != len(r.brokers) && conn == nil; i++ {
		broker := r.brokers[i]
		var first, last int64

		t0 := time.Now()
		conn, err = r.dialer.DialLeader(ctx, "tcp", broker, r.topic, r.partition)
		t1 := time.Now()
		r.stats.dials.observe(1)
		r.stats.dialTime.observeDuration(t1.Sub(t0))

		if err != nil {
			continue
		}

		if first, last, err = r.readOffsets(conn); err != nil {
			conn.Close()
			conn = nil
			break
		}

		switch {
		case offset == FirstOffset:
			offset = first

		case offset == LastOffset:
			offset = last

		case offset < first:
			offset = first
		}

		r.withLogger(func(log Logger) {
			log.Printf("the kafka reader for partition %d of %s is seeking to offset %d", r.partition, r.topic, toHumanOffset(offset))
		})

		if start, err = conn.Seek(offset, SeekAbsolute); err != nil {
			conn.Close()
			conn = nil
			break
		}

		conn.SetDeadline(time.Time{})
	}

	return
}

func (r *reader) read(ctx context.Context, offset int64, conn *Conn) (int64, error) {
	r.stats.fetches.observe(1)
	r.stats.offset.observe(offset)

	t0 := time.Now()
	conn.SetReadDeadline(t0.Add(r.maxWait))

	batch := conn.ReadBatchWith(ReadBatchConfig{
		MinBytes:       r.minBytes,
		MaxBytes:       r.maxBytes,
		IsolationLevel: r.isolationLevel,
	})
	highWaterMark := batch.HighWaterMark()

	t1 := time.Now()
	r.stats.waitTime.observeDuration(t1.Sub(t0))

	var msg Message
	var err error
	var size int64
	var bytes int64

	for {
		conn.SetReadDeadline(time.Now().Add(r.readBatchTimeout))

		if msg, err = batch.ReadMessage(); err != nil {
			batch.Close()
			break
		}

		n := int64(len(msg.Key) + len(msg.Value))
		r.stats.messages.observe(1)
		r.stats.bytes.observe(n)

		if err = r.sendMessage(ctx, msg, highWaterMark); err != nil {
			batch.Close()
			break
		}

		offset = msg.Offset + 1
		r.stats.offset.observe(offset)
		r.stats.lag.observe(highWaterMark - offset)

		size++
		bytes += n
	}

	conn.SetReadDeadline(time.Time{})

	t2 := time.Now()
	r.stats.readTime.observeDuration(t2.Sub(t1))
	r.stats.fetchSize.observe(size)
	r.stats.fetchBytes.observe(bytes)
	return offset, err
}

func (r *reader) readOffsets(conn *Conn) (first, last int64, err error) {
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	return conn.ReadOffsets()
}

func (r *reader) sendMessage(ctx context.Context, msg Message, watermark int64) error {
	select {
	case r.msgs <- readerMessage{version: r.version, message: msg, watermark: watermark}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *reader) sendError(ctx context.Context, err error) error {
	select {
	case r.msgs <- readerMessage{version: r.version, error: err}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *reader) withLogger(do func(Logger)) {
	if r.logger != nil {
		do(r.logger)
	}
}

func (r *reader) withErrorLogger(do func(Logger)) {
	if r.errorLogger != nil {
		do(r.errorLogger)
	} else {
		r.withLogger(do)
	}
}



func extractTopics(members []GroupMember) []string {
	visited := map[string]struct{}{}
	var topics []string

	for _, member := range members {
		for _, topic := range member.Topics {
			if _, seen := visited[topic]; seen {
				continue
			}

			topics = append(topics, topic)
			visited[topic] = struct{}{}
		}
	}

	sort.Strings(topics)

	return topics
}

type humanOffset int64

func toHumanOffset(v int64) humanOffset {
	return humanOffset(v)
}

func (offset humanOffset) Format(w fmt.State, _ rune) {
	v := int64(offset)
	switch v {
	case FirstOffset:
		fmt.Fprint(w, "first offset")
	case LastOffset:
		fmt.Fprint(w, "last offset")
	default:
		fmt.Fprint(w, strconv.FormatInt(v, 10))
	}
}
