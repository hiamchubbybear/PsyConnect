package kafka

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	metadataAPI "github.com/segmentio/kafka-go/protocol/metadata"
)































































type Writer struct {
	
	
	
	
	
	Addr net.Addr

	
	
	
	
	
	Topic string

	
	
	
	Balancer Balancer

	
	
	
	MaxAttempts int

	
	
	
	
	WriteBackoffMin time.Duration

	
	
	
	
	WriteBackoffMax time.Duration

	
	
	
	
	BatchSize int

	
	
	
	
	BatchBytes int64

	
	
	
	
	BatchTimeout time.Duration

	
	
	
	ReadTimeout time.Duration

	
	
	
	WriteTimeout time.Duration

	
	
	
	
	
	
	
	
	RequiredAcks RequiredAcks

	
	
	
	
	
	
	Async bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Completion func(messages []Message, err error)

	
	Compression Compression

	
	
	Logger Logger

	
	
	ErrorLogger Logger

	
	
	
	Transport RoundTripper

	
	AllowAutoTopicCreation bool

	
	group   sync.WaitGroup
	mutex   sync.Mutex
	closed  bool
	writers map[topicPartition]*partitionWriter

	
	
	
	
	once sync.Once
	*writerStats

	
	
	
	roundRobin RoundRobin

	
	transport *Transport
}






type WriterConfig struct {
	
	
	
	
	
	Brokers []string

	
	
	
	
	
	Topic string

	
	
	
	
	Dialer *Dialer

	
	
	
	Balancer Balancer

	
	
	
	MaxAttempts int

	
	
	
	
	QueueCapacity int

	
	
	
	
	BatchSize int

	
	
	
	
	BatchBytes int

	
	
	
	
	BatchTimeout time.Duration

	
	
	
	ReadTimeout time.Duration

	
	
	
	WriteTimeout time.Duration

	
	
	
	
	RebalanceInterval time.Duration

	
	
	
	IdleConnTimeout time.Duration

	
	
	
	
	RequiredAcks int

	
	
	
	
	Async bool

	
	CompressionCodec

	
	
	Logger Logger

	
	
	ErrorLogger Logger
}

type topicPartition struct {
	topic     string
	partition int32
}


func (config *WriterConfig) Validate() error {
	if len(config.Brokers) == 0 {
		return errors.New("cannot create a kafka writer with an empty list of brokers")
	}
	return nil
}



type WriterStats struct {
	Writes   int64 `metric:"kafka.writer.write.count"     type:"counter"`
	Messages int64 `metric:"kafka.writer.message.count"   type:"counter"`
	Bytes    int64 `metric:"kafka.writer.message.bytes"   type:"counter"`
	Errors   int64 `metric:"kafka.writer.error.count"     type:"counter"`

	BatchTime      DurationStats `metric:"kafka.writer.batch.seconds"`
	BatchQueueTime DurationStats `metric:"kafka.writer.batch.queue.seconds"`
	WriteTime      DurationStats `metric:"kafka.writer.write.seconds"`
	WaitTime       DurationStats `metric:"kafka.writer.wait.seconds"`
	Retries        int64         `metric:"kafka.writer.retries.count" type:"counter"`
	BatchSize      SummaryStats  `metric:"kafka.writer.batch.size"`
	BatchBytes     SummaryStats  `metric:"kafka.writer.batch.bytes"`

	MaxAttempts     int64         `metric:"kafka.writer.attempts.max"  type:"gauge"`
	WriteBackoffMin time.Duration `metric:"kafka.writer.backoff.min"   type:"gauge"`
	WriteBackoffMax time.Duration `metric:"kafka.writer.backoff.max"   type:"gauge"`
	MaxBatchSize    int64         `metric:"kafka.writer.batch.max"     type:"gauge"`
	BatchTimeout    time.Duration `metric:"kafka.writer.batch.timeout" type:"gauge"`
	ReadTimeout     time.Duration `metric:"kafka.writer.read.timeout"  type:"gauge"`
	WriteTimeout    time.Duration `metric:"kafka.writer.write.timeout" type:"gauge"`
	RequiredAcks    int64         `metric:"kafka.writer.acks.required" type:"gauge"`
	Async           bool          `metric:"kafka.writer.async"         type:"gauge"`

	Topic string `tag:"topic"`

	
	
	Dials    int64         `metric:"kafka.writer.dial.count" type:"counter"`
	DialTime DurationStats `metric:"kafka.writer.dial.seconds"`

	
	
	
	
	
	
	Rebalances        int64
	RebalanceInterval time.Duration
	QueueLength       int64
	QueueCapacity     int64
	ClientID          string
}






type writerStats struct {
	dials          counter
	writes         counter
	messages       counter
	bytes          counter
	errors         counter
	dialTime       summary
	batchTime      summary
	batchQueueTime summary
	writeTime      summary
	waitTime       summary
	retries        counter
	batchSize      summary
	batchSizeBytes summary
}






func NewWriter(config WriterConfig) *Writer {
	if err := config.Validate(); err != nil {
		panic(err)
	}

	if config.Dialer == nil {
		config.Dialer = DefaultDialer
	}

	if config.Balancer == nil {
		config.Balancer = &RoundRobin{}
	}

	
	kafkaDialer := DefaultDialer
	if config.Dialer != nil {
		kafkaDialer = config.Dialer
	}

	dialer := (&net.Dialer{
		Timeout:       kafkaDialer.Timeout,
		Deadline:      kafkaDialer.Deadline,
		LocalAddr:     kafkaDialer.LocalAddr,
		DualStack:     kafkaDialer.DualStack,
		FallbackDelay: kafkaDialer.FallbackDelay,
		KeepAlive:     kafkaDialer.KeepAlive,
	})

	var resolver Resolver
	if r, ok := kafkaDialer.Resolver.(*net.Resolver); ok {
		dialer.Resolver = r
	} else {
		resolver = kafkaDialer.Resolver
	}

	stats := new(writerStats)
	
	
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		defer func() {
			stats.dials.observe(1)
			stats.dialTime.observe(int64(time.Since(start)))
		}()
		address, err := lookupHost(ctx, addr, resolver)
		if err != nil {
			return nil, err
		}
		return dialer.DialContext(ctx, network, address)
	}

	idleTimeout := config.IdleConnTimeout
	if idleTimeout == 0 {
		
		
		
		idleTimeout = 9 * time.Minute
	}

	metadataTTL := config.RebalanceInterval
	if metadataTTL == 0 {
		
		metadataTTL = 15 * time.Second
	}

	transport := &Transport{
		Dial:        dial,
		SASL:        kafkaDialer.SASLMechanism,
		TLS:         kafkaDialer.TLS,
		ClientID:    kafkaDialer.ClientID,
		IdleTimeout: idleTimeout,
		MetadataTTL: metadataTTL,
	}

	w := &Writer{
		Addr:         TCP(config.Brokers...),
		Topic:        config.Topic,
		MaxAttempts:  config.MaxAttempts,
		BatchSize:    config.BatchSize,
		Balancer:     config.Balancer,
		BatchBytes:   int64(config.BatchBytes),
		BatchTimeout: config.BatchTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		RequiredAcks: RequiredAcks(config.RequiredAcks),
		Async:        config.Async,
		Logger:       config.Logger,
		ErrorLogger:  config.ErrorLogger,
		Transport:    transport,
		transport:    transport,
		writerStats:  stats,
	}

	if config.RequiredAcks == 0 {
		
		
		w.RequiredAcks = RequireAll
	}

	if config.CompressionCodec != nil {
		w.Compression = Compression(config.CompressionCodec.Code())
	}

	return w
}




func (w *Writer) enter() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.closed {
		return false
	}
	w.group.Add(1)
	return true
}



func (w *Writer) leave() { w.group.Done() }






func (w *Writer) spawn(f func()) {
	w.group.Add(1)
	go func() {
		defer w.group.Done()
		f()
	}()
}





func (w *Writer) Close() error {
	w.mutex.Lock()
	
	
	
	
	w.closed = true

	
	for _, writer := range w.writers {
		writer.close()
	}

	for partition := range w.writers {
		delete(w.writers, partition)
	}

	w.mutex.Unlock()
	w.group.Wait()

	if w.transport != nil {
		w.transport.CloseIdleConnections()
	}

	return nil
}




























func (w *Writer) WriteMessages(ctx context.Context, msgs ...Message) error {
	if w.Addr == nil {
		return errors.New("kafka.(*Writer).WriteMessages: cannot create a kafka writer with a nil address")
	}

	if !w.enter() {
		return io.ErrClosedPipe
	}
	defer w.leave()

	if len(msgs) == 0 {
		return nil
	}

	balancer := w.balancer()
	batchBytes := w.batchBytes()

	for i := range msgs {
		n := int64(msgs[i].totalSize())
		if n > batchBytes {
			
			
			
			
			
			return messageTooLarge(msgs, i)
		}
	}

	
	
	
	
	
	assignments := make(map[topicPartition][]int32)

	for i, msg := range msgs {
		topic, err := w.chooseTopic(msg)
		if err != nil {
			return err
		}

		numPartitions, err := w.partitions(ctx, topic)
		if err != nil {
			return err
		}

		partition := balancer.Balance(msg, loadCachedPartitions(numPartitions)...)

		key := topicPartition{
			topic:     topic,
			partition: int32(partition),
		}

		assignments[key] = append(assignments[key], int32(i))
	}

	batches := w.batchMessages(msgs, assignments)
	if w.Async {
		return nil
	}

	done := ctx.Done()
	hasErrors := false
	for batch := range batches {
		select {
		case <-done:
			return ctx.Err()
		case <-batch.done:
			if batch.err != nil {
				hasErrors = true
			}
		}
	}

	if !hasErrors {
		return nil
	}

	werr := make(WriteErrors, len(msgs))

	for batch, indexes := range batches {
		for _, i := range indexes {
			werr[i] = batch.err
		}
	}
	return werr
}

func (w *Writer) batchMessages(messages []Message, assignments map[topicPartition][]int32) map[*writeBatch][]int32 {
	var batches map[*writeBatch][]int32
	if !w.Async {
		batches = make(map[*writeBatch][]int32, len(assignments))
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.writers == nil {
		w.writers = map[topicPartition]*partitionWriter{}
	}

	for key, indexes := range assignments {
		writer := w.writers[key]
		if writer == nil {
			writer = newPartitionWriter(w, key)
			w.writers[key] = writer
		}
		wbatches := writer.writeMessages(messages, indexes)

		for batch, idxs := range wbatches {
			batches[batch] = idxs
		}
	}

	return batches
}

func (w *Writer) produce(key topicPartition, batch *writeBatch) (*ProduceResponse, error) {
	timeout := w.writeTimeout()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return w.client(timeout).Produce(ctx, &ProduceRequest{
		Partition:    int(key.partition),
		Topic:        key.topic,
		RequiredAcks: w.RequiredAcks,
		Compression:  w.Compression,
		Records: &writerRecords{
			msgs: batch.msgs,
		},
	})
}

func (w *Writer) partitions(ctx context.Context, topic string) (int, error) {
	client := w.client(w.readTimeout())
	
	
	
	
	
	
	r, err := client.transport().RoundTrip(ctx, client.Addr, &metadataAPI.Request{
		TopicNames:             []string{topic},
		AllowAutoTopicCreation: w.AllowAutoTopicCreation,
	})
	if err != nil {
		return 0, err
	}
	for _, t := range r.(*metadataAPI.Response).Topics {
		if t.Name == topic {
			
			if t.ErrorCode != 0 {
				return 0, Error(t.ErrorCode)
			}
			return len(t.Partitions), nil
		}
	}
	return 0, UnknownTopicOrPartition
}

func (w *Writer) client(timeout time.Duration) *Client {
	return &Client{
		Addr:      w.Addr,
		Transport: w.Transport,
		Timeout:   timeout,
	}
}

func (w *Writer) balancer() Balancer {
	if w.Balancer != nil {
		return w.Balancer
	}
	return &w.roundRobin
}

func (w *Writer) maxAttempts() int {
	if w.MaxAttempts > 0 {
		return w.MaxAttempts
	}
	
	
	
	
	return 10
}

func (w *Writer) writeBackoffMin() time.Duration {
	if w.WriteBackoffMin > 0 {
		return w.WriteBackoffMin
	}
	return 100 * time.Millisecond
}

func (w *Writer) writeBackoffMax() time.Duration {
	if w.WriteBackoffMax > 0 {
		return w.WriteBackoffMax
	}
	return 1 * time.Second
}

func (w *Writer) batchSize() int {
	if w.BatchSize > 0 {
		return w.BatchSize
	}
	return 100
}

func (w *Writer) batchBytes() int64 {
	if w.BatchBytes > 0 {
		return w.BatchBytes
	}
	return 1048576
}

func (w *Writer) batchTimeout() time.Duration {
	if w.BatchTimeout > 0 {
		return w.BatchTimeout
	}
	return 1 * time.Second
}

func (w *Writer) readTimeout() time.Duration {
	if w.ReadTimeout > 0 {
		return w.ReadTimeout
	}
	return 10 * time.Second
}

func (w *Writer) writeTimeout() time.Duration {
	if w.WriteTimeout > 0 {
		return w.WriteTimeout
	}
	return 10 * time.Second
}

func (w *Writer) withLogger(do func(Logger)) {
	if w.Logger != nil {
		do(w.Logger)
	}
}

func (w *Writer) withErrorLogger(do func(Logger)) {
	if w.ErrorLogger != nil {
		do(w.ErrorLogger)
	} else {
		w.withLogger(do)
	}
}

func (w *Writer) stats() *writerStats {
	w.once.Do(func() {
		
		
		if w.writerStats == nil {
			w.writerStats = new(writerStats)
		}
	})
	return w.writerStats
}








func (w *Writer) Stats() WriterStats {
	stats := w.stats()
	return WriterStats{
		Dials:           stats.dials.snapshot(),
		Writes:          stats.writes.snapshot(),
		Messages:        stats.messages.snapshot(),
		Bytes:           stats.bytes.snapshot(),
		Errors:          stats.errors.snapshot(),
		DialTime:        stats.dialTime.snapshotDuration(),
		BatchTime:       stats.batchTime.snapshotDuration(),
		BatchQueueTime:  stats.batchQueueTime.snapshotDuration(),
		WriteTime:       stats.writeTime.snapshotDuration(),
		WaitTime:        stats.waitTime.snapshotDuration(),
		Retries:         stats.retries.snapshot(),
		BatchSize:       stats.batchSize.snapshot(),
		BatchBytes:      stats.batchSizeBytes.snapshot(),
		MaxAttempts:     int64(w.maxAttempts()),
		WriteBackoffMin: w.writeBackoffMin(),
		WriteBackoffMax: w.writeBackoffMax(),
		MaxBatchSize:    int64(w.batchSize()),
		BatchTimeout:    w.batchTimeout(),
		ReadTimeout:     w.readTimeout(),
		WriteTimeout:    w.writeTimeout(),
		RequiredAcks:    int64(w.RequiredAcks),
		Async:           w.Async,
		Topic:           w.Topic,
	}
}

func (w *Writer) chooseTopic(msg Message) (string, error) {
	
	
	if w.Topic != "" && msg.Topic != "" {
		return "", errors.New("kafka.(*Writer): Topic must not be specified for both Writer and Message")
	} else if w.Topic == "" && msg.Topic == "" {
		return "", errors.New("kafka.(*Writer): Topic must be specified for Writer or Message")
	}

	
	if msg.Topic != "" {
		return msg.Topic, nil
	}

	return w.Topic, nil
}

type batchQueue struct {
	queue []*writeBatch

	
	
	
	mutex *sync.Mutex
	cond  *sync.Cond

	closed bool
}

func (b *batchQueue) Put(batch *writeBatch) bool {
	b.cond.L.Lock()
	defer b.cond.L.Unlock()
	defer b.cond.Broadcast()

	if b.closed {
		return false
	}
	b.queue = append(b.queue, batch)
	return true
}

func (b *batchQueue) Get() *writeBatch {
	b.cond.L.Lock()
	defer b.cond.L.Unlock()

	for len(b.queue) == 0 && !b.closed {
		b.cond.Wait()
	}

	if len(b.queue) == 0 {
		return nil
	}

	batch := b.queue[0]
	b.queue[0] = nil
	b.queue = b.queue[1:]

	return batch
}

func (b *batchQueue) Close() {
	b.cond.L.Lock()
	defer b.cond.L.Unlock()
	defer b.cond.Broadcast()

	b.closed = true
}

func newBatchQueue(initialSize int) batchQueue {
	bq := batchQueue{
		queue: make([]*writeBatch, 0, initialSize),
		mutex: &sync.Mutex{},
		cond:  &sync.Cond{},
	}

	bq.cond.L = bq.mutex

	return bq
}



type partitionWriter struct {
	meta  topicPartition
	queue batchQueue

	mutex     sync.Mutex
	currBatch *writeBatch

	
	
	w *Writer
}

func newPartitionWriter(w *Writer, key topicPartition) *partitionWriter {
	writer := &partitionWriter{
		meta:  key,
		queue: newBatchQueue(10),
		w:     w,
	}
	w.spawn(writer.writeBatches)
	return writer
}

func (ptw *partitionWriter) writeBatches() {
	for {
		batch := ptw.queue.Get()

		
		
		
		if batch == nil {
			return
		}

		ptw.writeBatch(batch)
	}
}

func (ptw *partitionWriter) writeMessages(msgs []Message, indexes []int32) map[*writeBatch][]int32 {
	ptw.mutex.Lock()
	defer ptw.mutex.Unlock()

	batchSize := ptw.w.batchSize()
	batchBytes := ptw.w.batchBytes()

	var batches map[*writeBatch][]int32
	if !ptw.w.Async {
		batches = make(map[*writeBatch][]int32, 1)
	}

	for _, i := range indexes {
	assignMessage:
		batch := ptw.currBatch
		if batch == nil {
			batch = ptw.newWriteBatch()
			ptw.currBatch = batch
		}
		if !batch.add(msgs[i], batchSize, batchBytes) {
			batch.trigger()
			ptw.queue.Put(batch)
			ptw.currBatch = nil
			goto assignMessage
		}

		if batch.full(batchSize, batchBytes) {
			batch.trigger()
			ptw.queue.Put(batch)
			ptw.currBatch = nil
		}

		if !ptw.w.Async {
			batches[batch] = append(batches[batch], i)
		}
	}
	return batches
}


func (ptw *partitionWriter) newWriteBatch() *writeBatch {
	batch := newWriteBatch(time.Now(), ptw.w.batchTimeout())
	ptw.w.spawn(func() { ptw.awaitBatch(batch) })
	return batch
}




func (ptw *partitionWriter) awaitBatch(batch *writeBatch) {
	select {
	case <-batch.timer.C:
		ptw.mutex.Lock()
		
		
		
		
		
		
		
		if ptw.currBatch == batch {
			ptw.queue.Put(batch)
			ptw.currBatch = nil
		}
		ptw.mutex.Unlock()
	case <-batch.ready:
		
		
		
		batch.timer.Stop()
	}
	stats := ptw.w.stats()
	stats.batchQueueTime.observe(int64(time.Since(batch.time)))
}

func (ptw *partitionWriter) writeBatch(batch *writeBatch) {
	stats := ptw.w.stats()
	stats.batchTime.observe(int64(time.Since(batch.time)))
	stats.batchSize.observe(int64(len(batch.msgs)))
	stats.batchSizeBytes.observe(batch.bytes)

	var res *ProduceResponse
	var err error
	key := ptw.meta
	for attempt, maxAttempts := 0, ptw.w.maxAttempts(); attempt < maxAttempts; attempt++ {
		if attempt != 0 {
			stats.retries.observe(1)
			
			
			
			
			
			
			
			
			
			
			delay := backoff(attempt, ptw.w.writeBackoffMin(), ptw.w.writeBackoffMax())
			ptw.w.withLogger(func(log Logger) {
				log.Printf("backing off %s writing %d messages to %s (partition: %d)", delay, len(batch.msgs), key.topic, key.partition)
			})
			time.Sleep(delay)
		}

		ptw.w.withLogger(func(log Logger) {
			log.Printf("writing %d messages to %s (partition: %d)", len(batch.msgs), key.topic, key.partition)
		})

		start := time.Now()
		res, err = ptw.w.produce(key, batch)

		stats.writes.observe(1)
		stats.messages.observe(int64(len(batch.msgs)))
		stats.bytes.observe(batch.bytes)
		
		
		
		
		
		stats.writeTime.observe(int64(time.Since(start)))

		if res != nil {
			err = res.Error
			stats.waitTime.observe(int64(res.Throttle))
		}

		if err == nil {
			break
		}

		stats.errors.observe(1)

		ptw.w.withErrorLogger(func(log Logger) {
			log.Printf("error writing messages to %s (partition %d, attempt %d): %s", key.topic, key.partition, attempt, err)
		})

		if !isTemporary(err) && !isTransientNetworkError(err) {
			break
		}
	}

	if res != nil {
		for i := range batch.msgs {
			m := &batch.msgs[i]
			m.Topic = key.topic
			m.Partition = int(key.partition)
			m.Offset = res.BaseOffset + int64(i)

			if m.Time.IsZero() {
				m.Time = res.LogAppendTime
			}
		}
	}

	if ptw.w.Completion != nil {
		ptw.w.Completion(batch.msgs, err)
	}

	batch.complete(err)
}

func (ptw *partitionWriter) close() {
	ptw.mutex.Lock()
	defer ptw.mutex.Unlock()

	if ptw.currBatch != nil {
		batch := ptw.currBatch
		ptw.queue.Put(batch)
		ptw.currBatch = nil
		batch.trigger()
	}

	ptw.queue.Close()
}

type writeBatch struct {
	time  time.Time
	msgs  []Message
	size  int
	bytes int64
	ready chan struct{}
	done  chan struct{}
	timer *time.Timer
	err   error 
}

func newWriteBatch(now time.Time, timeout time.Duration) *writeBatch {
	return &writeBatch{
		time:  now,
		ready: make(chan struct{}),
		done:  make(chan struct{}),
		timer: time.NewTimer(timeout),
	}
}

func (b *writeBatch) add(msg Message, maxSize int, maxBytes int64) bool {
	bytes := int64(msg.totalSize())

	if b.size > 0 && (b.bytes+bytes) > maxBytes {
		return false
	}

	if cap(b.msgs) == 0 {
		b.msgs = make([]Message, 0, maxSize)
	}

	b.msgs = append(b.msgs, msg)
	b.size++
	b.bytes += bytes
	return true
}

func (b *writeBatch) full(maxSize int, maxBytes int64) bool {
	return b.size >= maxSize || b.bytes >= maxBytes
}

func (b *writeBatch) trigger() {
	close(b.ready)
}

func (b *writeBatch) complete(err error) {
	b.err = err
	close(b.done)
}

type writerRecords struct {
	msgs   []Message
	index  int
	record Record
	key    bytesReadCloser
	value  bytesReadCloser
}

func (r *writerRecords) ReadRecord() (*Record, error) {
	if r.index >= 0 && r.index < len(r.msgs) {
		m := &r.msgs[r.index]
		r.index++
		r.record = Record{
			Time:    m.Time,
			Headers: m.Headers,
		}
		if m.Key != nil {
			r.key.Reset(m.Key)
			r.record.Key = &r.key
		}
		if m.Value != nil {
			r.value.Reset(m.Value)
			r.record.Value = &r.value
		}
		return &r.record, nil
	}
	return nil, io.EOF
}

type bytesReadCloser struct{ bytes.Reader }

func (*bytesReadCloser) Close() error { return nil }








var partitionsCache atomic.Value

func loadCachedPartitions(numPartitions int) []int {
	partitions, ok := partitionsCache.Load().([]int)
	if ok && len(partitions) >= numPartitions {
		return partitions[:numPartitions]
	}

	const alignment = 128
	n := ((numPartitions / alignment) + 1) * alignment

	partitions = make([]int, n)
	for i := range partitions {
		partitions[i] = i
	}

	partitionsCache.Store(partitions)
	return partitions[:numPartitions]
}
