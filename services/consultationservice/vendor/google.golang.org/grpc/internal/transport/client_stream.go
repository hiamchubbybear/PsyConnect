

package transport

import (
	"sync/atomic"

	"golang.org/x/net/http2"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)


type ClientStream struct {
	*Stream 

	ct       *http2Client
	done     chan struct{} 
	doneFunc func()        

	headerChan       chan struct{} 
	headerChanClosed uint32        
	
	
	
	headerValid bool
	header      metadata.MD 
	noHeaders   bool        

	bytesReceived atomic.Bool 
	unprocessed   atomic.Bool 

	status *status.Status 
}


func (s *ClientStream) Read(n int) (mem.BufferSlice, error) {
	b, err := s.Stream.read(n)
	if err == nil {
		s.ct.incrMsgRecv()
	}
	return b, err
}


func (s *ClientStream) Close(err error) {
	var (
		rst     bool
		rstCode http2.ErrCode
	)
	if err != nil {
		rst = true
		rstCode = http2.ErrCodeCancel
	}
	s.ct.closeStream(s, err, rst, rstCode, status.Convert(err), nil, false)
}


func (s *ClientStream) Write(hdr []byte, data mem.BufferSlice, opts *WriteOptions) error {
	return s.ct.write(s, hdr, data, opts)
}


func (s *ClientStream) BytesReceived() bool {
	return s.bytesReceived.Load()
}



func (s *ClientStream) Unprocessed() bool {
	return s.unprocessed.Load()
}

func (s *ClientStream) waitOnHeader() {
	select {
	case <-s.ctx.Done():
		
		
		s.Close(ContextErr(s.ctx.Err()))
		
		
		<-s.headerChan
	case <-s.headerChan:
	}
}



func (s *ClientStream) RecvCompress() string {
	s.waitOnHeader()
	return s.recvCompress
}



func (s *ClientStream) Done() <-chan struct{} {
	return s.done
}





func (s *ClientStream) Header() (metadata.MD, error) {
	s.waitOnHeader()

	if !s.headerValid || s.noHeaders {
		return nil, s.status.Err()
	}

	return s.header.Copy(), nil
}




func (s *ClientStream) TrailersOnly() bool {
	s.waitOnHeader()
	return s.noHeaders
}




func (s *ClientStream) Status() *status.Status {
	return s.status
}
