






package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/internal/grpclog"
	"google.golang.org/grpc/internal/grpcutil"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)




func NewServerHandlerTransport(w http.ResponseWriter, r *http.Request, stats []stats.Handler, bufferPool mem.BufferPool) (ServerTransport, error) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		msg := fmt.Sprintf("invalid gRPC request method %q", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return nil, errors.New(msg)
	}
	contentType := r.Header.Get("Content-Type")
	
	contentSubtype, validContentType := grpcutil.ContentSubtype(contentType)
	if !validContentType {
		msg := fmt.Sprintf("invalid gRPC request content-type %q", contentType)
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return nil, errors.New(msg)
	}
	if r.ProtoMajor != 2 {
		msg := "gRPC requires HTTP/2"
		http.Error(w, msg, http.StatusHTTPVersionNotSupported)
		return nil, errors.New(msg)
	}
	if _, ok := w.(http.Flusher); !ok {
		msg := "gRPC requires a ResponseWriter supporting http.Flusher"
		http.Error(w, msg, http.StatusInternalServerError)
		return nil, errors.New(msg)
	}

	var localAddr net.Addr
	if la := r.Context().Value(http.LocalAddrContextKey); la != nil {
		localAddr, _ = la.(net.Addr)
	}
	var authInfo credentials.AuthInfo
	if r.TLS != nil {
		authInfo = credentials.TLSInfo{State: *r.TLS, CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity}}
	}
	p := peer.Peer{
		Addr:      strAddr(r.RemoteAddr),
		LocalAddr: localAddr,
		AuthInfo:  authInfo,
	}
	st := &serverHandlerTransport{
		rw:             w,
		req:            r,
		closedCh:       make(chan struct{}),
		writes:         make(chan func()),
		peer:           p,
		contentType:    contentType,
		contentSubtype: contentSubtype,
		stats:          stats,
		bufferPool:     bufferPool,
	}
	st.logger = prefixLoggerForServerHandlerTransport(st)

	if v := r.Header.Get("grpc-timeout"); v != "" {
		to, err := decodeTimeout(v)
		if err != nil {
			msg := fmt.Sprintf("malformed grpc-timeout: %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			return nil, status.Error(codes.Internal, msg)
		}
		st.timeoutSet = true
		st.timeout = to
	}

	metakv := []string{"content-type", contentType}
	if r.Host != "" {
		metakv = append(metakv, ":authority", r.Host)
	}
	for k, vv := range r.Header {
		k = strings.ToLower(k)
		if isReservedHeader(k) && !isWhitelistedHeader(k) {
			continue
		}
		for _, v := range vv {
			v, err := decodeMetadataHeader(k, v)
			if err != nil {
				msg := fmt.Sprintf("malformed binary metadata %q in header %q: %v", v, k, err)
				http.Error(w, msg, http.StatusBadRequest)
				return nil, status.Error(codes.Internal, msg)
			}
			metakv = append(metakv, k, v)
		}
	}
	st.headerMD = metadata.Pairs(metakv...)

	return st, nil
}






type serverHandlerTransport struct {
	rw         http.ResponseWriter
	req        *http.Request
	timeoutSet bool
	timeout    time.Duration

	headerMD metadata.MD

	peer peer.Peer

	closeOnce sync.Once
	closedCh  chan struct{} 

	
	
	
	writes chan func()

	
	
	writeStatusMu sync.Mutex

	
	contentType string
	
	
	contentSubtype string

	stats  []stats.Handler
	logger *grpclog.PrefixLogger

	bufferPool mem.BufferPool
}

func (ht *serverHandlerTransport) Close(err error) {
	ht.closeOnce.Do(func() {
		if ht.logger.V(logLevel) {
			ht.logger.Infof("Closing: %v", err)
		}
		close(ht.closedCh)
	})
}

func (ht *serverHandlerTransport) Peer() *peer.Peer {
	return &peer.Peer{
		Addr:      ht.peer.Addr,
		LocalAddr: ht.peer.LocalAddr,
		AuthInfo:  ht.peer.AuthInfo,
	}
}



type strAddr string

func (a strAddr) Network() string {
	if a != "" {
		
		
		
		
		
		
		
		
		
		return "tcp"
	}
	return ""
}

func (a strAddr) String() string { return string(a) }


func (ht *serverHandlerTransport) do(fn func()) error {
	select {
	case <-ht.closedCh:
		return ErrConnClosing
	case ht.writes <- fn:
		return nil
	}
}

func (ht *serverHandlerTransport) writeStatus(s *ServerStream, st *status.Status) error {
	ht.writeStatusMu.Lock()
	defer ht.writeStatusMu.Unlock()

	headersWritten := s.updateHeaderSent()
	err := ht.do(func() {
		if !headersWritten {
			ht.writePendingHeaders(s)
		}

		
		
		
		ht.rw.(http.Flusher).Flush()

		h := ht.rw.Header()
		h.Set("Grpc-Status", fmt.Sprintf("%d", st.Code()))
		if m := st.Message(); m != "" {
			h.Set("Grpc-Message", encodeGrpcMessage(m))
		}

		s.hdrMu.Lock()
		defer s.hdrMu.Unlock()
		if p := st.Proto(); p != nil && len(p.Details) > 0 {
			delete(s.trailer, grpcStatusDetailsBinHeader)
			stBytes, err := proto.Marshal(p)
			if err != nil {
				
				panic(err)
			}

			h.Set(grpcStatusDetailsBinHeader, encodeBinHeader(stBytes))
		}

		if len(s.trailer) > 0 {
			for k, vv := range s.trailer {
				
				if isReservedHeader(k) {
					continue
				}
				for _, v := range vv {
					
					
					h.Add(http2.TrailerPrefix+k, encodeMetadataHeader(k, v))
				}
			}
		}
	})

	if err == nil { 
		
		
		for _, sh := range ht.stats {
			sh.HandleRPC(s.Context(), &stats.OutTrailer{
				Trailer: s.trailer.Copy(),
			})
		}
	}
	ht.Close(errors.New("finished writing status"))
	return err
}



func (ht *serverHandlerTransport) writePendingHeaders(s *ServerStream) {
	ht.writeCommonHeaders(s)
	ht.writeCustomHeaders(s)
}



func (ht *serverHandlerTransport) writeCommonHeaders(s *ServerStream) {
	h := ht.rw.Header()
	h["Date"] = nil 
	h.Set("Content-Type", ht.contentType)

	
	
	
	
	
	h.Add("Trailer", "Grpc-Status")
	h.Add("Trailer", "Grpc-Message")
	h.Add("Trailer", "Grpc-Status-Details-Bin")

	if s.sendCompress != "" {
		h.Set("Grpc-Encoding", s.sendCompress)
	}
}



func (ht *serverHandlerTransport) writeCustomHeaders(s *ServerStream) {
	h := ht.rw.Header()

	s.hdrMu.Lock()
	for k, vv := range s.header {
		if isReservedHeader(k) {
			continue
		}
		for _, v := range vv {
			h.Add(k, encodeMetadataHeader(k, v))
		}
	}

	s.hdrMu.Unlock()
}

func (ht *serverHandlerTransport) write(s *ServerStream, hdr []byte, data mem.BufferSlice, _ *WriteOptions) error {
	
	
	
	data.Ref()
	headersWritten := s.updateHeaderSent()
	err := ht.do(func() {
		defer data.Free()
		if !headersWritten {
			ht.writePendingHeaders(s)
		}
		ht.rw.Write(hdr)
		for _, b := range data {
			_, _ = ht.rw.Write(b.ReadOnlyData())
		}
		ht.rw.(http.Flusher).Flush()
	})
	if err != nil {
		data.Free()
		return err
	}
	return nil
}

func (ht *serverHandlerTransport) writeHeader(s *ServerStream, md metadata.MD) error {
	if err := s.SetHeader(md); err != nil {
		return err
	}

	headersWritten := s.updateHeaderSent()
	err := ht.do(func() {
		if !headersWritten {
			ht.writePendingHeaders(s)
		}

		ht.rw.WriteHeader(200)
		ht.rw.(http.Flusher).Flush()
	})

	if err == nil {
		for _, sh := range ht.stats {
			
			
			sh.HandleRPC(s.Context(), &stats.OutHeader{
				Header:      md.Copy(),
				Compression: s.sendCompress,
			})
		}
	}
	return err
}

func (ht *serverHandlerTransport) HandleStreams(ctx context.Context, startStream func(*ServerStream)) {
	
	var cancel context.CancelFunc
	if ht.timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, ht.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	
	requestOver := make(chan struct{})
	go func() {
		select {
		case <-requestOver:
		case <-ht.closedCh:
		case <-ht.req.Context().Done():
		}
		cancel()
		ht.Close(errors.New("request is done processing"))
	}()

	ctx = metadata.NewIncomingContext(ctx, ht.headerMD)
	req := ht.req
	s := &ServerStream{
		Stream: &Stream{
			id:             0, 
			ctx:            ctx,
			requestRead:    func(int) {},
			buf:            newRecvBuffer(),
			method:         req.URL.Path,
			recvCompress:   req.Header.Get("grpc-encoding"),
			contentSubtype: ht.contentSubtype,
		},
		cancel:           cancel,
		st:               ht,
		headerWireLength: 0, 
	}
	s.trReader = &transportReader{
		reader:        &recvBufferReader{ctx: s.ctx, ctxDone: s.ctx.Done(), recv: s.buf},
		windowHandler: func(int) {},
	}

	
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)

		for {
			buf := ht.bufferPool.Get(http2MaxFrameLen)
			n, err := req.Body.Read(*buf)
			if n > 0 {
				*buf = (*buf)[:n]
				s.buf.put(recvMsg{buffer: mem.NewBuffer(buf, ht.bufferPool)})
			} else {
				ht.bufferPool.Put(buf)
			}
			if err != nil {
				s.buf.put(recvMsg{err: mapRecvMsgError(err)})
				return
			}
		}
	}()

	
	
	
	
	startStream(s)

	ht.runStream()
	close(requestOver)

	
	req.Body.Close()
	<-readerDone
}

func (ht *serverHandlerTransport) runStream() {
	for {
		select {
		case fn := <-ht.writes:
			fn()
		case <-ht.closedCh:
			return
		}
	}
}

func (ht *serverHandlerTransport) incrMsgRecv() {}

func (ht *serverHandlerTransport) Drain(string) {
	panic("Drain() is not implemented")
}








func mapRecvMsgError(err error) error {
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return err
	}
	if se, ok := err.(http2.StreamError); ok {
		if code, ok := http2ErrConvTab[se.Code]; ok {
			return status.Error(code, se.Error())
		}
	}
	if strings.Contains(err.Error(), "body closed by handler") {
		return status.Error(codes.Canceled, err.Error())
	}
	return connectionErrorf(true, err, "%s", err.Error())
}
