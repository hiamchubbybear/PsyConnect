

package transport

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc/internal/grpclog"
	"google.golang.org/grpc/internal/grpcutil"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/status"
)

var updateHeaderTblSize = func(e *hpack.Encoder, v uint32) {
	e.SetMaxDynamicTableSizeLimit(v)
}

type itemNode struct {
	it   any
	next *itemNode
}

type itemList struct {
	head *itemNode
	tail *itemNode
}

func (il *itemList) enqueue(i any) {
	n := &itemNode{it: i}
	if il.tail == nil {
		il.head, il.tail = n, n
		return
	}
	il.tail.next = n
	il.tail = n
}



func (il *itemList) peek() any {
	return il.head.it
}

func (il *itemList) dequeue() any {
	if il.head == nil {
		return nil
	}
	i := il.head.it
	il.head = il.head.next
	if il.head == nil {
		il.tail = nil
	}
	return i
}

func (il *itemList) dequeueAll() *itemNode {
	h := il.head
	il.head, il.tail = nil, nil
	return h
}

func (il *itemList) isEmpty() bool {
	return il.head == nil
}









const maxQueuedTransportResponseFrames = 50

type cbItem interface {
	isTransportResponseFrame() bool
}


type registerStream struct {
	streamID uint32
	wq       *writeQuota
}

func (*registerStream) isTransportResponseFrame() bool { return false }


type headerFrame struct {
	streamID   uint32
	hf         []hpack.HeaderField
	endStream  bool               
	initStream func(uint32) error 
	onWrite    func()
	wq         *writeQuota    
	cleanup    *cleanupStream 
	onOrphaned func(error)    
}

func (h *headerFrame) isTransportResponseFrame() bool {
	return h.cleanup != nil && h.cleanup.rst 
}

type cleanupStream struct {
	streamID uint32
	rst      bool
	rstCode  http2.ErrCode
	onWrite  func()
}

func (c *cleanupStream) isTransportResponseFrame() bool { return c.rst } 

type earlyAbortStream struct {
	httpStatus     uint32
	streamID       uint32
	contentSubtype string
	status         *status.Status
	rst            bool
}

func (*earlyAbortStream) isTransportResponseFrame() bool { return false }

type dataFrame struct {
	streamID  uint32
	endStream bool
	h         []byte
	reader    mem.Reader
	
	
	onEachWrite func()
}

func (*dataFrame) isTransportResponseFrame() bool { return false }

type incomingWindowUpdate struct {
	streamID  uint32
	increment uint32
}

func (*incomingWindowUpdate) isTransportResponseFrame() bool { return false }

type outgoingWindowUpdate struct {
	streamID  uint32
	increment uint32
}

func (*outgoingWindowUpdate) isTransportResponseFrame() bool {
	return false 
}

type incomingSettings struct {
	ss []http2.Setting
}

func (*incomingSettings) isTransportResponseFrame() bool { return true } 

type outgoingSettings struct {
	ss []http2.Setting
}

func (*outgoingSettings) isTransportResponseFrame() bool { return false }

type incomingGoAway struct {
}

func (*incomingGoAway) isTransportResponseFrame() bool { return false }

type goAway struct {
	code      http2.ErrCode
	debugData []byte
	headsUp   bool
	closeConn error 
}

func (*goAway) isTransportResponseFrame() bool { return false }

type ping struct {
	ack  bool
	data [8]byte
}

func (*ping) isTransportResponseFrame() bool { return true }

type outFlowControlSizeRequest struct {
	resp chan uint32
}

func (*outFlowControlSizeRequest) isTransportResponseFrame() bool { return false }





type closeConnection struct{}

func (closeConnection) isTransportResponseFrame() bool { return false }

type outStreamState int

const (
	active outStreamState = iota
	empty
	waitingOnStreamQuota
)

type outStream struct {
	id               uint32
	state            outStreamState
	itl              *itemList
	bytesOutStanding int
	wq               *writeQuota

	next *outStream
	prev *outStream
}

func (s *outStream) deleteSelf() {
	if s.prev != nil {
		s.prev.next = s.next
	}
	if s.next != nil {
		s.next.prev = s.prev
	}
	s.next, s.prev = nil, nil
}

type outStreamList struct {
	
	
	
	
	
	
	
	head *outStream
	tail *outStream
}

func newOutStreamList() *outStreamList {
	head, tail := new(outStream), new(outStream)
	head.next = tail
	tail.prev = head
	return &outStreamList{
		head: head,
		tail: tail,
	}
}

func (l *outStreamList) enqueue(s *outStream) {
	e := l.tail.prev
	e.next = s
	s.prev = e
	s.next = l.tail
	l.tail.prev = s
}


func (l *outStreamList) dequeue() *outStream {
	b := l.head.next
	if b == l.tail {
		return nil
	}
	b.deleteSelf()
	return b
}








type controlBuffer struct {
	wakeupCh chan struct{}   
	done     <-chan struct{} 

	
	
	mu              sync.Mutex
	consumerWaiting bool      
	closed          bool      
	list            *itemList 

	
	
	
	
	
	transportResponseFrames int
	trfChan                 atomic.Pointer[chan struct{}]
}

func newControlBuffer(done <-chan struct{}) *controlBuffer {
	return &controlBuffer{
		wakeupCh: make(chan struct{}, 1),
		list:     &itemList{},
		done:     done,
	}
}




func (c *controlBuffer) throttle() {
	if ch := c.trfChan.Load(); ch != nil {
		select {
		case <-(*ch):
		case <-c.done:
		}
	}
}


func (c *controlBuffer) put(it cbItem) error {
	_, err := c.executeAndPut(nil, it)
	return err
}








func (c *controlBuffer) executeAndPut(f func() bool, it cbItem) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return false, ErrConnClosing
	}
	if f != nil {
		if !f() { 
			return false, nil
		}
	}
	if it == nil {
		return true, nil
	}

	var wakeUp bool
	if c.consumerWaiting {
		wakeUp = true
		c.consumerWaiting = false
	}
	c.list.enqueue(it)
	if it.isTransportResponseFrame() {
		c.transportResponseFrames++
		if c.transportResponseFrames == maxQueuedTransportResponseFrames {
			
			
			ch := make(chan struct{})
			c.trfChan.Store(&ch)
		}
	}
	if wakeUp {
		select {
		case c.wakeupCh <- struct{}{}:
		default:
		}
	}
	return true, nil
}





func (c *controlBuffer) get(block bool) (any, error) {
	for {
		c.mu.Lock()
		frame, err := c.getOnceLocked()
		if frame != nil || err != nil || !block {
			
			
			
			
			c.mu.Unlock()
			return frame, err
		}
		c.consumerWaiting = true
		c.mu.Unlock()

		
		select {
		case <-c.wakeupCh:
		case <-c.done:
			return nil, errors.New("transport closed by client")
		}
	}
}




func (c *controlBuffer) getOnceLocked() (any, error) {
	if c.closed {
		return false, ErrConnClosing
	}
	if c.list.isEmpty() {
		return nil, nil
	}
	h := c.list.dequeue().(cbItem)
	if h.isTransportResponseFrame() {
		if c.transportResponseFrames == maxQueuedTransportResponseFrames {
			
			
			ch := c.trfChan.Swap(nil)
			close(*ch)
		}
		c.transportResponseFrames--
	}
	return h, nil
}




func (c *controlBuffer) finish() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}
	c.closed = true
	
	
	
	for head := c.list.dequeueAll(); head != nil; head = head.next {
		switch v := head.it.(type) {
		case *headerFrame:
			if v.onOrphaned != nil { 
				v.onOrphaned(ErrConnClosing)
			}
		case *dataFrame:
			_ = v.reader.Close()
		}
	}

	
	
	
	ch := c.trfChan.Swap(nil)
	if ch != nil {
		close(*ch)
	}
}

type side int

const (
	clientSide side = iota
	serverSide
)










type loopyWriter struct {
	side      side
	cbuf      *controlBuffer
	sendQuota uint32
	oiws      uint32 
	
	
	
	estdStreams map[uint32]*outStream 
	
	
	
	
	activeStreams *outStreamList
	framer        *framer
	hBuf          *bytes.Buffer  
	hEnc          *hpack.Encoder 
	bdpEst        *bdpEstimator
	draining      bool
	conn          net.Conn
	logger        *grpclog.PrefixLogger
	bufferPool    mem.BufferPool

	
	ssGoAwayHandler func(*goAway) (bool, error)
}

func newLoopyWriter(s side, fr *framer, cbuf *controlBuffer, bdpEst *bdpEstimator, conn net.Conn, logger *grpclog.PrefixLogger, goAwayHandler func(*goAway) (bool, error), bufferPool mem.BufferPool) *loopyWriter {
	var buf bytes.Buffer
	l := &loopyWriter{
		side:            s,
		cbuf:            cbuf,
		sendQuota:       defaultWindowSize,
		oiws:            defaultWindowSize,
		estdStreams:     make(map[uint32]*outStream),
		activeStreams:   newOutStreamList(),
		framer:          fr,
		hBuf:            &buf,
		hEnc:            hpack.NewEncoder(&buf),
		bdpEst:          bdpEst,
		conn:            conn,
		logger:          logger,
		ssGoAwayHandler: goAwayHandler,
		bufferPool:      bufferPool,
	}
	return l
}

const minBatchSize = 1000






















func (l *loopyWriter) run() (err error) {
	defer func() {
		if l.logger.V(logLevel) {
			l.logger.Infof("loopyWriter exiting with error: %v", err)
		}
		if !isIOError(err) {
			l.framer.writer.Flush()
		}
		l.cbuf.finish()
	}()
	for {
		it, err := l.cbuf.get(true)
		if err != nil {
			return err
		}
		if err = l.handle(it); err != nil {
			return err
		}
		if _, err = l.processData(); err != nil {
			return err
		}
		gosched := true
	hasdata:
		for {
			it, err := l.cbuf.get(false)
			if err != nil {
				return err
			}
			if it != nil {
				if err = l.handle(it); err != nil {
					return err
				}
				if _, err = l.processData(); err != nil {
					return err
				}
				continue hasdata
			}
			isEmpty, err := l.processData()
			if err != nil {
				return err
			}
			if !isEmpty {
				continue hasdata
			}
			if gosched {
				gosched = false
				if l.framer.writer.offset < minBatchSize {
					runtime.Gosched()
					continue hasdata
				}
			}
			l.framer.writer.Flush()
			break hasdata
		}
	}
}

func (l *loopyWriter) outgoingWindowUpdateHandler(w *outgoingWindowUpdate) error {
	return l.framer.fr.WriteWindowUpdate(w.streamID, w.increment)
}

func (l *loopyWriter) incomingWindowUpdateHandler(w *incomingWindowUpdate) {
	
	if w.streamID == 0 {
		l.sendQuota += w.increment
		return
	}
	
	if str, ok := l.estdStreams[w.streamID]; ok {
		str.bytesOutStanding -= int(w.increment)
		if strQuota := int(l.oiws) - str.bytesOutStanding; strQuota > 0 && str.state == waitingOnStreamQuota {
			str.state = active
			l.activeStreams.enqueue(str)
			return
		}
	}
}

func (l *loopyWriter) outgoingSettingsHandler(s *outgoingSettings) error {
	return l.framer.fr.WriteSettings(s.ss...)
}

func (l *loopyWriter) incomingSettingsHandler(s *incomingSettings) error {
	l.applySettings(s.ss)
	return l.framer.fr.WriteSettingsAck()
}

func (l *loopyWriter) registerStreamHandler(h *registerStream) {
	str := &outStream{
		id:    h.streamID,
		state: empty,
		itl:   &itemList{},
		wq:    h.wq,
	}
	l.estdStreams[h.streamID] = str
}

func (l *loopyWriter) headerHandler(h *headerFrame) error {
	if l.side == serverSide {
		str, ok := l.estdStreams[h.streamID]
		if !ok {
			if l.logger.V(logLevel) {
				l.logger.Infof("Unrecognized streamID %d in loopyWriter", h.streamID)
			}
			return nil
		}
		
		if !h.endStream {
			return l.writeHeader(h.streamID, h.endStream, h.hf, h.onWrite)
		}
		

		if str.state != empty { 
			
			str.itl.enqueue(h)
			return nil
		}
		if err := l.writeHeader(h.streamID, h.endStream, h.hf, h.onWrite); err != nil {
			return err
		}
		return l.cleanupStreamHandler(h.cleanup)
	}
	
	str := &outStream{
		id:    h.streamID,
		state: empty,
		itl:   &itemList{},
		wq:    h.wq,
	}
	return l.originateStream(str, h)
}

func (l *loopyWriter) originateStream(str *outStream, hdr *headerFrame) error {
	
	
	if l.draining {
		
		hdr.onOrphaned(errStreamDrain)
		return nil
	}
	if err := hdr.initStream(str.id); err != nil {
		return err
	}
	if err := l.writeHeader(str.id, hdr.endStream, hdr.hf, hdr.onWrite); err != nil {
		return err
	}
	l.estdStreams[str.id] = str
	return nil
}

func (l *loopyWriter) writeHeader(streamID uint32, endStream bool, hf []hpack.HeaderField, onWrite func()) error {
	if onWrite != nil {
		onWrite()
	}
	l.hBuf.Reset()
	for _, f := range hf {
		if err := l.hEnc.WriteField(f); err != nil {
			if l.logger.V(logLevel) {
				l.logger.Warningf("Encountered error while encoding headers: %v", err)
			}
		}
	}
	var (
		err               error
		endHeaders, first bool
	)
	first = true
	for !endHeaders {
		size := l.hBuf.Len()
		if size > http2MaxFrameLen {
			size = http2MaxFrameLen
		} else {
			endHeaders = true
		}
		if first {
			first = false
			err = l.framer.fr.WriteHeaders(http2.HeadersFrameParam{
				StreamID:      streamID,
				BlockFragment: l.hBuf.Next(size),
				EndStream:     endStream,
				EndHeaders:    endHeaders,
			})
		} else {
			err = l.framer.fr.WriteContinuation(
				streamID,
				endHeaders,
				l.hBuf.Next(size),
			)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *loopyWriter) preprocessData(df *dataFrame) {
	str, ok := l.estdStreams[df.streamID]
	if !ok {
		return
	}
	
	
	str.itl.enqueue(df)
	if str.state == empty {
		str.state = active
		l.activeStreams.enqueue(str)
	}
}

func (l *loopyWriter) pingHandler(p *ping) error {
	if !p.ack {
		l.bdpEst.timesnap(p.data)
	}
	return l.framer.fr.WritePing(p.ack, p.data)

}

func (l *loopyWriter) outFlowControlSizeRequestHandler(o *outFlowControlSizeRequest) {
	o.resp <- l.sendQuota
}

func (l *loopyWriter) cleanupStreamHandler(c *cleanupStream) error {
	c.onWrite()
	if str, ok := l.estdStreams[c.streamID]; ok {
		
		
		
		delete(l.estdStreams, c.streamID)
		str.deleteSelf()
		for head := str.itl.dequeueAll(); head != nil; head = head.next {
			if df, ok := head.it.(*dataFrame); ok {
				_ = df.reader.Close()
			}
		}
	}
	if c.rst { 
		if err := l.framer.fr.WriteRSTStream(c.streamID, c.rstCode); err != nil {
			return err
		}
	}
	if l.draining && len(l.estdStreams) == 0 {
		
		return errors.New("finished processing active streams while in draining mode")
	}
	return nil
}

func (l *loopyWriter) earlyAbortStreamHandler(eas *earlyAbortStream) error {
	if l.side == clientSide {
		return errors.New("earlyAbortStream not handled on client")
	}
	
	if eas.httpStatus == 0 {
		eas.httpStatus = 200
	}
	headerFields := []hpack.HeaderField{
		{Name: ":status", Value: strconv.Itoa(int(eas.httpStatus))},
		{Name: "content-type", Value: grpcutil.ContentType(eas.contentSubtype)},
		{Name: "grpc-status", Value: strconv.Itoa(int(eas.status.Code()))},
		{Name: "grpc-message", Value: encodeGrpcMessage(eas.status.Message())},
	}

	if err := l.writeHeader(eas.streamID, true, headerFields, nil); err != nil {
		return err
	}
	if eas.rst {
		if err := l.framer.fr.WriteRSTStream(eas.streamID, http2.ErrCodeNo); err != nil {
			return err
		}
	}
	return nil
}

func (l *loopyWriter) incomingGoAwayHandler(*incomingGoAway) error {
	if l.side == clientSide {
		l.draining = true
		if len(l.estdStreams) == 0 {
			
			return errors.New("received GOAWAY with no active streams")
		}
	}
	return nil
}

func (l *loopyWriter) goAwayHandler(g *goAway) error {
	
	if l.ssGoAwayHandler != nil {
		draining, err := l.ssGoAwayHandler(g)
		if err != nil {
			return err
		}
		l.draining = draining
	}
	return nil
}

func (l *loopyWriter) handle(i any) error {
	switch i := i.(type) {
	case *incomingWindowUpdate:
		l.incomingWindowUpdateHandler(i)
	case *outgoingWindowUpdate:
		return l.outgoingWindowUpdateHandler(i)
	case *incomingSettings:
		return l.incomingSettingsHandler(i)
	case *outgoingSettings:
		return l.outgoingSettingsHandler(i)
	case *headerFrame:
		return l.headerHandler(i)
	case *registerStream:
		l.registerStreamHandler(i)
	case *cleanupStream:
		return l.cleanupStreamHandler(i)
	case *earlyAbortStream:
		return l.earlyAbortStreamHandler(i)
	case *incomingGoAway:
		return l.incomingGoAwayHandler(i)
	case *dataFrame:
		l.preprocessData(i)
	case *ping:
		return l.pingHandler(i)
	case *goAway:
		return l.goAwayHandler(i)
	case *outFlowControlSizeRequest:
		l.outFlowControlSizeRequestHandler(i)
	case closeConnection:
		
		
		return ErrConnClosing
	default:
		return fmt.Errorf("transport: unknown control message type %T", i)
	}
	return nil
}

func (l *loopyWriter) applySettings(ss []http2.Setting) {
	for _, s := range ss {
		switch s.ID {
		case http2.SettingInitialWindowSize:
			o := l.oiws
			l.oiws = s.Val
			if o < l.oiws {
				
				for _, stream := range l.estdStreams {
					if stream.state == waitingOnStreamQuota {
						stream.state = active
						l.activeStreams.enqueue(stream)
					}
				}
			}
		case http2.SettingHeaderTableSize:
			updateHeaderTblSize(l.hEnc, s.Val)
		}
	}
}




func (l *loopyWriter) processData() (bool, error) {
	if l.sendQuota == 0 {
		return true, nil
	}
	str := l.activeStreams.dequeue() 
	if str == nil {
		return true, nil
	}
	dataItem := str.itl.peek().(*dataFrame) 
	
	
	
	
	
	

	if len(dataItem.h) == 0 && dataItem.reader.Remaining() == 0 { 
		
		if err := l.framer.fr.WriteData(dataItem.streamID, dataItem.endStream, nil); err != nil {
			return false, err
		}
		str.itl.dequeue() 
		_ = dataItem.reader.Close()
		if str.itl.isEmpty() {
			str.state = empty
		} else if trailer, ok := str.itl.peek().(*headerFrame); ok { 
			if err := l.writeHeader(trailer.streamID, trailer.endStream, trailer.hf, trailer.onWrite); err != nil {
				return false, err
			}
			if err := l.cleanupStreamHandler(trailer.cleanup); err != nil {
				return false, err
			}
		} else {
			l.activeStreams.enqueue(str)
		}
		return false, nil
	}

	
	maxSize := http2MaxFrameLen
	if strQuota := int(l.oiws) - str.bytesOutStanding; strQuota <= 0 { 
		str.state = waitingOnStreamQuota
		return false, nil
	} else if maxSize > strQuota {
		maxSize = strQuota
	}
	if maxSize > int(l.sendQuota) { 
		maxSize = int(l.sendQuota)
	}
	
	hSize := min(maxSize, len(dataItem.h))
	dSize := min(maxSize-hSize, dataItem.reader.Remaining())
	remainingBytes := len(dataItem.h) + dataItem.reader.Remaining() - hSize - dSize
	size := hSize + dSize

	var buf *[]byte

	if hSize != 0 && dSize == 0 {
		buf = &dataItem.h
	} else {
		
		
		
		pool := l.bufferPool
		if pool == nil {
			
			
			pool = mem.DefaultBufferPool()
		}
		buf = pool.Get(size)
		defer pool.Put(buf)

		copy((*buf)[:hSize], dataItem.h)
		_, _ = dataItem.reader.Read((*buf)[hSize:])
	}

	
	str.wq.replenish(size)
	var endStream bool
	
	if dataItem.endStream && remainingBytes == 0 {
		endStream = true
	}
	if dataItem.onEachWrite != nil {
		dataItem.onEachWrite()
	}
	if err := l.framer.fr.WriteData(dataItem.streamID, endStream, (*buf)[:size]); err != nil {
		return false, err
	}
	str.bytesOutStanding += size
	l.sendQuota -= uint32(size)
	dataItem.h = dataItem.h[hSize:]

	if remainingBytes == 0 { 
		_ = dataItem.reader.Close()
		str.itl.dequeue()
	}
	if str.itl.isEmpty() {
		str.state = empty
	} else if trailer, ok := str.itl.peek().(*headerFrame); ok { 
		if err := l.writeHeader(trailer.streamID, trailer.endStream, trailer.hf, trailer.onWrite); err != nil {
			return false, err
		}
		if err := l.cleanupStreamHandler(trailer.cleanup); err != nil {
			return false, err
		}
	} else if int(l.oiws)-str.bytesOutStanding <= 0 { 
		str.state = waitingOnStreamQuota
	} else { 
		l.activeStreams.enqueue(str)
	}
	return false, nil
}
