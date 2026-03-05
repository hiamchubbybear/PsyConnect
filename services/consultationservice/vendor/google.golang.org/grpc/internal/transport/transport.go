




package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

const logLevel = 2



type recvMsg struct {
	buffer mem.Buffer
	
	
	
	err error
}







type recvBuffer struct {
	c       chan recvMsg
	mu      sync.Mutex
	backlog []recvMsg
	err     error
}

func newRecvBuffer() *recvBuffer {
	b := &recvBuffer{
		c: make(chan recvMsg, 1),
	}
	return b
}

func (b *recvBuffer) put(r recvMsg) {
	b.mu.Lock()
	if b.err != nil {
		
		
		r.buffer.Free()
		b.mu.Unlock()
		
		
		return
	}
	b.err = r.err
	if len(b.backlog) == 0 {
		select {
		case b.c <- r:
			b.mu.Unlock()
			return
		default:
		}
	}
	b.backlog = append(b.backlog, r)
	b.mu.Unlock()
}

func (b *recvBuffer) load() {
	b.mu.Lock()
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = recvMsg{}
			b.backlog = b.backlog[1:]
		default:
		}
	}
	b.mu.Unlock()
}





func (b *recvBuffer) get() <-chan recvMsg {
	return b.c
}



type recvBufferReader struct {
	closeStream func(error) 
	ctx         context.Context
	ctxDone     <-chan struct{} 
	recv        *recvBuffer
	last        mem.Buffer 
	err         error
}

func (r *recvBufferReader) ReadMessageHeader(header []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.last != nil {
		n, r.last = mem.ReadUnsafe(header, r.last)
		return n, nil
	}
	if r.closeStream != nil {
		n, r.err = r.readMessageHeaderClient(header)
	} else {
		n, r.err = r.readMessageHeader(header)
	}
	return n, r.err
}





func (r *recvBufferReader) Read(n int) (buf mem.Buffer, err error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.last != nil {
		buf = r.last
		if r.last.Len() > n {
			buf, r.last = mem.SplitUnsafe(buf, n)
		} else {
			r.last = nil
		}
		return buf, nil
	}
	if r.closeStream != nil {
		buf, r.err = r.readClient(n)
	} else {
		buf, r.err = r.read(n)
	}
	return buf, r.err
}

func (r *recvBufferReader) readMessageHeader(header []byte) (n int, err error) {
	select {
	case <-r.ctxDone:
		return 0, ContextErr(r.ctx.Err())
	case m := <-r.recv.get():
		return r.readMessageHeaderAdditional(m, header)
	}
}

func (r *recvBufferReader) read(n int) (buf mem.Buffer, err error) {
	select {
	case <-r.ctxDone:
		return nil, ContextErr(r.ctx.Err())
	case m := <-r.recv.get():
		return r.readAdditional(m, n)
	}
}

func (r *recvBufferReader) readMessageHeaderClient(header []byte) (n int, err error) {
	
	
	
	select {
	case <-r.ctxDone:
		
		
		
		
		
		
		
		
		
		
		
		
		
		r.closeStream(ContextErr(r.ctx.Err()))
		m := <-r.recv.get()
		return r.readMessageHeaderAdditional(m, header)
	case m := <-r.recv.get():
		return r.readMessageHeaderAdditional(m, header)
	}
}

func (r *recvBufferReader) readClient(n int) (buf mem.Buffer, err error) {
	
	
	
	select {
	case <-r.ctxDone:
		
		
		
		
		
		
		
		
		
		
		
		
		
		r.closeStream(ContextErr(r.ctx.Err()))
		m := <-r.recv.get()
		return r.readAdditional(m, n)
	case m := <-r.recv.get():
		return r.readAdditional(m, n)
	}
}

func (r *recvBufferReader) readMessageHeaderAdditional(m recvMsg, header []byte) (n int, err error) {
	r.recv.load()
	if m.err != nil {
		if m.buffer != nil {
			m.buffer.Free()
		}
		return 0, m.err
	}

	n, r.last = mem.ReadUnsafe(header, m.buffer)

	return n, nil
}

func (r *recvBufferReader) readAdditional(m recvMsg, n int) (b mem.Buffer, err error) {
	r.recv.load()
	if m.err != nil {
		if m.buffer != nil {
			m.buffer.Free()
		}
		return nil, m.err
	}

	if m.buffer.Len() > n {
		m.buffer, r.last = mem.SplitUnsafe(m.buffer, n)
	}

	return m.buffer, nil
}

type streamState uint32

const (
	streamActive    streamState = iota
	streamWriteDone             
	streamReadDone              
	streamDone                  
)


type Stream struct {
	id           uint32
	ctx          context.Context 
	method       string          
	recvCompress string
	sendCompress string
	buf          *recvBuffer
	trReader     *transportReader
	fc           *inFlow
	wq           *writeQuota

	
	
	requestRead func(int)

	state streamState

	
	
	contentSubtype string

	trailer metadata.MD 
}

func (s *Stream) swapState(st streamState) streamState {
	return streamState(atomic.SwapUint32((*uint32)(&s.state), uint32(st)))
}

func (s *Stream) compareAndSwapState(oldState, newState streamState) bool {
	return atomic.CompareAndSwapUint32((*uint32)(&s.state), uint32(oldState), uint32(newState))
}

func (s *Stream) getState() streamState {
	return streamState(atomic.LoadUint32((*uint32)(&s.state)))
}





func (s *Stream) Trailer() metadata.MD {
	return s.trailer.Copy()
}


func (s *Stream) Context() context.Context {
	return s.ctx
}


func (s *Stream) Method() string {
	return s.method
}

func (s *Stream) write(m recvMsg) {
	s.buf.put(m)
}









func (s *Stream) ReadMessageHeader(header []byte) (err error) {
	
	if er := s.trReader.er; er != nil {
		return er
	}
	s.requestRead(len(header))
	for len(header) != 0 {
		n, err := s.trReader.ReadMessageHeader(header)
		header = header[n:]
		if len(header) == 0 {
			err = nil
		}
		if err != nil {
			if n > 0 && err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return err
		}
	}
	return nil
}


func (s *Stream) read(n int) (data mem.BufferSlice, err error) {
	
	if er := s.trReader.er; er != nil {
		return nil, er
	}
	s.requestRead(n)
	for n != 0 {
		buf, err := s.trReader.Read(n)
		var bufLen int
		if buf != nil {
			bufLen = buf.Len()
		}
		n -= bufLen
		if n == 0 {
			err = nil
		}
		if err != nil {
			if bufLen > 0 && err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			data.Free()
			return nil, err
		}
		data = append(data, buf)
	}
	return data, nil
}





type transportReader struct {
	reader *recvBufferReader
	
	
	windowHandler func(int)
	er            error
}

func (t *transportReader) ReadMessageHeader(header []byte) (int, error) {
	n, err := t.reader.ReadMessageHeader(header)
	if err != nil {
		t.er = err
		return 0, err
	}
	t.windowHandler(n)
	return n, nil
}

func (t *transportReader) Read(n int) (mem.Buffer, error) {
	buf, err := t.reader.Read(n)
	if err != nil {
		t.er = err
		return buf, err
	}
	t.windowHandler(buf.Len())
	return buf, nil
}



func (s *Stream) GoString() string {
	return fmt.Sprintf("<stream: %p, %v>", s, s.method)
}


type transportState int

const (
	reachable transportState = iota
	closing
	draining
)


type ServerConfig struct {
	MaxStreams            uint32
	ConnectionTimeout     time.Duration
	Credentials           credentials.TransportCredentials
	InTapHandle           tap.ServerInHandle
	StatsHandlers         []stats.Handler
	KeepaliveParams       keepalive.ServerParameters
	KeepalivePolicy       keepalive.EnforcementPolicy
	InitialWindowSize     int32
	InitialConnWindowSize int32
	WriteBufferSize       int
	ReadBufferSize        int
	SharedWriteBuffer     bool
	ChannelzParent        *channelz.Server
	MaxHeaderListSize     *uint32
	HeaderTableSize       *uint32
	BufferPool            mem.BufferPool
}


type ConnectOptions struct {
	
	UserAgent string
	
	Dialer func(context.Context, string) (net.Conn, error)
	
	FailOnNonTempDialError bool
	
	PerRPCCredentials []credentials.PerRPCCredentials
	
	
	TransportCredentials credentials.TransportCredentials
	
	
	CredsBundle credentials.Bundle
	
	KeepaliveParams keepalive.ClientParameters
	
	StatsHandlers []stats.Handler
	
	InitialWindowSize int32
	
	InitialConnWindowSize int32
	
	WriteBufferSize int
	
	ReadBufferSize int
	
	SharedWriteBuffer bool
	
	ChannelzParent *channelz.SubChannel
	
	MaxHeaderListSize *uint32
	
	BufferPool mem.BufferPool
}



type WriteOptions struct {
	
	
	Last bool
}


type CallHdr struct {
	
	Host string

	
	Method string

	
	
	SendCompress string

	
	Creds credentials.PerRPCCredentials

	
	
	
	
	
	
	ContentSubtype string

	PreviousAttempts int 

	DoneFunc func() 
}



type ClientTransport interface {
	
	
	
	Close(err error)

	
	
	
	
	
	GracefulClose()

	
	NewStream(ctx context.Context, callHdr *CallHdr) (*ClientStream, error)

	
	
	
	
	
	Error() <-chan struct{}

	
	
	
	GoAway() <-chan struct{}

	
	
	GetGoAwayReason() (GoAwayReason, string)

	
	RemoteAddr() net.Addr
}






type ServerTransport interface {
	
	HandleStreams(context.Context, func(*ServerStream))

	
	
	
	Close(err error)

	
	Peer() *peer.Peer

	
	Drain(debugData string)
}

type internalServerTransport interface {
	ServerTransport
	writeHeader(s *ServerStream, md metadata.MD) error
	write(s *ServerStream, hdr []byte, data mem.BufferSlice, opts *WriteOptions) error
	writeStatus(s *ServerStream, st *status.Status) error
	incrMsgRecv()
}


func connectionErrorf(temp bool, e error, format string, a ...any) ConnectionError {
	return ConnectionError{
		Desc: fmt.Sprintf(format, a...),
		temp: temp,
		err:  e,
	}
}



type ConnectionError struct {
	Desc string
	temp bool
	err  error
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("connection error: desc = %q", e.Desc)
}


func (e ConnectionError) Temporary() bool {
	return e.temp
}


func (e ConnectionError) Origin() error {
	
	
	if e.err == nil {
		return e
	}
	return e.err
}



func (e ConnectionError) Unwrap() error {
	return e.err
}

var (
	
	ErrConnClosing = connectionErrorf(true, nil, "transport is closing")
	
	
	
	errStreamDrain = status.Error(codes.Unavailable, "the connection is draining")
	
	
	errStreamDone = errors.New("the stream is done")
	
	
	statusGoAway = status.New(codes.Unavailable, "the stream is rejected because server is draining the connection")
)


type GoAwayReason uint8

const (
	
	GoAwayInvalid GoAwayReason = 0
	
	GoAwayNoReason GoAwayReason = 1
	
	
	
	GoAwayTooManyPings GoAwayReason = 2
)


func ContextErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	}
	return status.Errorf(codes.Internal, "Unexpected error from context packet: %v", err)
}
