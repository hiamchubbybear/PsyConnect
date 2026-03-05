

package transport

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


type ServerStream struct {
	*Stream 

	st      internalServerTransport
	ctxDone <-chan struct{}    
	cancel  context.CancelFunc 

	
	
	clientAdvertisedCompressors string
	headerWireLength            int

	
	hdrMu      sync.Mutex
	header     metadata.MD 
	headerSent atomic.Bool 
}


func (s *ServerStream) Read(n int) (mem.BufferSlice, error) {
	b, err := s.Stream.read(n)
	if err == nil {
		s.st.incrMsgRecv()
	}
	return b, err
}


func (s *ServerStream) SendHeader(md metadata.MD) error {
	return s.st.writeHeader(s, md)
}


func (s *ServerStream) Write(hdr []byte, data mem.BufferSlice, opts *WriteOptions) error {
	return s.st.write(s, hdr, data, opts)
}



func (s *ServerStream) WriteStatus(st *status.Status) error {
	return s.st.writeStatus(s, st)
}


func (s *ServerStream) isHeaderSent() bool {
	return s.headerSent.Load()
}



func (s *ServerStream) updateHeaderSent() bool {
	return s.headerSent.Swap(true)
}



func (s *ServerStream) RecvCompress() string {
	return s.recvCompress
}


func (s *ServerStream) SendCompress() string {
	return s.sendCompress
}






func (s *ServerStream) ContentSubtype() string {
	return s.contentSubtype
}


func (s *ServerStream) SetSendCompress(name string) error {
	if s.isHeaderSent() || s.getState() == streamDone {
		return errors.New("transport: set send compressor called after headers sent or stream done")
	}

	s.sendCompress = name
	return nil
}



func (s *ServerStream) SetContext(ctx context.Context) {
	s.ctx = ctx
}



func (s *ServerStream) ClientAdvertisedCompressors() []string {
	values := strings.Split(s.clientAdvertisedCompressors, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}




func (s *ServerStream) Header() (metadata.MD, error) {
	
	
	return s.header.Copy(), nil
}



func (s *ServerStream) HeaderWireLength() int {
	return s.headerWireLength
}



func (s *ServerStream) SetHeader(md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	if s.isHeaderSent() || s.getState() == streamDone {
		return ErrIllegalHeaderWrite
	}
	s.hdrMu.Lock()
	s.header = metadata.Join(s.header, md)
	s.hdrMu.Unlock()
	return nil
}




func (s *ServerStream) SetTrailer(md metadata.MD) error {
	if md.Len() == 0 {
		return nil
	}
	if s.getState() == streamDone {
		return ErrIllegalHeaderWrite
	}
	s.hdrMu.Lock()
	s.trailer = metadata.Join(s.trailer, md)
	s.hdrMu.Unlock()
	return nil
}
