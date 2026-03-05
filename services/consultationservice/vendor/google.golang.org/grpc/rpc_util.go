

package grpc

import (
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)




type Compressor interface {
	
	Do(w io.Writer, p []byte) error
	
	Type() string
}

type gzipCompressor struct {
	pool sync.Pool
}




func NewGZIPCompressor() Compressor {
	c, _ := NewGZIPCompressorWithLevel(gzip.DefaultCompression)
	return c
}







func NewGZIPCompressorWithLevel(level int) (Compressor, error) {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		return nil, fmt.Errorf("grpc: invalid compression level: %d", level)
	}
	return &gzipCompressor{
		pool: sync.Pool{
			New: func() any {
				w, err := gzip.NewWriterLevel(io.Discard, level)
				if err != nil {
					panic(err)
				}
				return w
			},
		},
	}, nil
}

func (c *gzipCompressor) Do(w io.Writer, p []byte) error {
	z := c.pool.Get().(*gzip.Writer)
	defer c.pool.Put(z)
	z.Reset(w)
	if _, err := z.Write(p); err != nil {
		return err
	}
	return z.Close()
}

func (c *gzipCompressor) Type() string {
	return "gzip"
}




type Decompressor interface {
	
	Do(r io.Reader) ([]byte, error)
	
	Type() string
}

type gzipDecompressor struct {
	pool sync.Pool
}




func NewGZIPDecompressor() Decompressor {
	return &gzipDecompressor{}
}

func (d *gzipDecompressor) Do(r io.Reader) ([]byte, error) {
	var z *gzip.Reader
	switch maybeZ := d.pool.Get().(type) {
	case nil:
		newZ, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		z = newZ
	case *gzip.Reader:
		z = maybeZ
		if err := z.Reset(r); err != nil {
			d.pool.Put(z)
			return nil, err
		}
	}

	defer func() {
		z.Close()
		d.pool.Put(z)
	}()
	return io.ReadAll(z)
}

func (d *gzipDecompressor) Type() string {
	return "gzip"
}


type callInfo struct {
	compressorName        string
	failFast              bool
	maxReceiveMessageSize *int
	maxSendMessageSize    *int
	creds                 credentials.PerRPCCredentials
	contentSubtype        string
	codec                 baseCodec
	maxRetryRPCBufferSize int
	onFinish              []func(err error)
}

func defaultCallInfo() *callInfo {
	return &callInfo{
		failFast:              true,
		maxRetryRPCBufferSize: 256 * 1024, 
	}
}



type CallOption interface {
	
	
	before(*callInfo) error

	
	
	after(*callInfo, *csAttempt)
}




type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}





func StaticMethod() CallOption {
	return StaticMethodCallOption{}
}



type StaticMethodCallOption struct {
	EmptyCallOption
}



func Header(md *metadata.MD) CallOption {
	return HeaderCallOption{HeaderAddr: md}
}








type HeaderCallOption struct {
	HeaderAddr *metadata.MD
}

func (o HeaderCallOption) before(*callInfo) error { return nil }
func (o HeaderCallOption) after(_ *callInfo, attempt *csAttempt) {
	*o.HeaderAddr, _ = attempt.transportStream.Header()
}



func Trailer(md *metadata.MD) CallOption {
	return TrailerCallOption{TrailerAddr: md}
}








type TrailerCallOption struct {
	TrailerAddr *metadata.MD
}

func (o TrailerCallOption) before(*callInfo) error { return nil }
func (o TrailerCallOption) after(_ *callInfo, attempt *csAttempt) {
	*o.TrailerAddr = attempt.transportStream.Trailer()
}



func Peer(p *peer.Peer) CallOption {
	return PeerCallOption{PeerAddr: p}
}








type PeerCallOption struct {
	PeerAddr *peer.Peer
}

func (o PeerCallOption) before(*callInfo) error { return nil }
func (o PeerCallOption) after(_ *callInfo, attempt *csAttempt) {
	if x, ok := peer.FromContext(attempt.transportStream.Context()); ok {
		*o.PeerAddr = *x
	}
}








func WaitForReady(waitForReady bool) CallOption {
	return FailFastCallOption{FailFast: !waitForReady}
}




func FailFast(failFast bool) CallOption {
	return FailFastCallOption{FailFast: failFast}
}








type FailFastCallOption struct {
	FailFast bool
}

func (o FailFastCallOption) before(c *callInfo) error {
	c.failFast = o.FailFast
	return nil
}
func (o FailFastCallOption) after(*callInfo, *csAttempt) {}












func OnFinish(onFinish func(err error)) CallOption {
	return OnFinishCallOption{
		OnFinish: onFinish,
	}
}








type OnFinishCallOption struct {
	OnFinish func(error)
}

func (o OnFinishCallOption) before(c *callInfo) error {
	c.onFinish = append(c.onFinish, o.OnFinish)
	return nil
}

func (o OnFinishCallOption) after(*callInfo, *csAttempt) {}




func MaxCallRecvMsgSize(bytes int) CallOption {
	return MaxRecvMsgSizeCallOption{MaxRecvMsgSize: bytes}
}








type MaxRecvMsgSizeCallOption struct {
	MaxRecvMsgSize int
}

func (o MaxRecvMsgSizeCallOption) before(c *callInfo) error {
	c.maxReceiveMessageSize = &o.MaxRecvMsgSize
	return nil
}
func (o MaxRecvMsgSizeCallOption) after(*callInfo, *csAttempt) {}




func MaxCallSendMsgSize(bytes int) CallOption {
	return MaxSendMsgSizeCallOption{MaxSendMsgSize: bytes}
}








type MaxSendMsgSizeCallOption struct {
	MaxSendMsgSize int
}

func (o MaxSendMsgSizeCallOption) before(c *callInfo) error {
	c.maxSendMessageSize = &o.MaxSendMsgSize
	return nil
}
func (o MaxSendMsgSizeCallOption) after(*callInfo, *csAttempt) {}



func PerRPCCredentials(creds credentials.PerRPCCredentials) CallOption {
	return PerRPCCredsCallOption{Creds: creds}
}








type PerRPCCredsCallOption struct {
	Creds credentials.PerRPCCredentials
}

func (o PerRPCCredsCallOption) before(c *callInfo) error {
	c.creds = o.Creds
	return nil
}
func (o PerRPCCredsCallOption) after(*callInfo, *csAttempt) {}









func UseCompressor(name string) CallOption {
	return CompressorCallOption{CompressorType: name}
}







type CompressorCallOption struct {
	CompressorType string
}

func (o CompressorCallOption) before(c *callInfo) error {
	c.compressorName = o.CompressorType
	return nil
}
func (o CompressorCallOption) after(*callInfo, *csAttempt) {}

















func CallContentSubtype(contentSubtype string) CallOption {
	return ContentSubtypeCallOption{ContentSubtype: strings.ToLower(contentSubtype)}
}








type ContentSubtypeCallOption struct {
	ContentSubtype string
}

func (o ContentSubtypeCallOption) before(c *callInfo) error {
	c.contentSubtype = o.ContentSubtype
	return nil
}
func (o ContentSubtypeCallOption) after(*callInfo, *csAttempt) {}



















func ForceCodec(codec encoding.Codec) CallOption {
	return ForceCodecCallOption{Codec: codec}
}








type ForceCodecCallOption struct {
	Codec encoding.Codec
}

func (o ForceCodecCallOption) before(c *callInfo) error {
	c.codec = newCodecV1Bridge(o.Codec)
	return nil
}
func (o ForceCodecCallOption) after(*callInfo, *csAttempt) {}



















func ForceCodecV2(codec encoding.CodecV2) CallOption {
	return ForceCodecV2CallOption{CodecV2: codec}
}








type ForceCodecV2CallOption struct {
	CodecV2 encoding.CodecV2
}

func (o ForceCodecV2CallOption) before(c *callInfo) error {
	c.codec = o.CodecV2
	return nil
}

func (o ForceCodecV2CallOption) after(*callInfo, *csAttempt) {}





func CallCustomCodec(codec Codec) CallOption {
	return CustomCodecCallOption{Codec: codec}
}








type CustomCodecCallOption struct {
	Codec Codec
}

func (o CustomCodecCallOption) before(c *callInfo) error {
	c.codec = newCodecV0Bridge(o.Codec)
	return nil
}
func (o CustomCodecCallOption) after(*callInfo, *csAttempt) {}








func MaxRetryRPCBufferSize(bytes int) CallOption {
	return MaxRetryRPCBufferSizeCallOption{bytes}
}








type MaxRetryRPCBufferSizeCallOption struct {
	MaxRetryRPCBufferSize int
}

func (o MaxRetryRPCBufferSizeCallOption) before(c *callInfo) error {
	c.maxRetryRPCBufferSize = o.MaxRetryRPCBufferSize
	return nil
}
func (o MaxRetryRPCBufferSizeCallOption) after(*callInfo, *csAttempt) {}


type payloadFormat uint8

const (
	compressionNone payloadFormat = 0 
	compressionMade payloadFormat = 1 
)

func (pf payloadFormat) isCompressed() bool {
	return pf == compressionMade
}

type streamReader interface {
	ReadMessageHeader(header []byte) error
	Read(n int) (mem.BufferSlice, error)
}


type parser struct {
	
	
	
	r streamReader

	
	
	header [5]byte

	
	bufferPool mem.BufferPool
}















func (p *parser) recvMsg(maxReceiveMessageSize int) (payloadFormat, mem.BufferSlice, error) {
	err := p.r.ReadMessageHeader(p.header[:])
	if err != nil {
		return 0, nil, err
	}

	pf := payloadFormat(p.header[0])
	length := binary.BigEndian.Uint32(p.header[1:])

	if int64(length) > int64(maxInt) {
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max length allowed on current machine (%d vs. %d)", length, maxInt)
	}
	if int(length) > maxReceiveMessageSize {
		return 0, nil, status.Errorf(codes.ResourceExhausted, "grpc: received message larger than max (%d vs. %d)", length, maxReceiveMessageSize)
	}

	data, err := p.r.Read(int(length))
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, nil, err
	}
	return pf, data, nil
}




func encode(c baseCodec, msg any) (mem.BufferSlice, error) {
	if msg == nil { 
		return nil, nil
	}
	b, err := c.Marshal(msg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
	}
	if bufSize := uint(b.Len()); bufSize > math.MaxUint32 {
		b.Free()
		return nil, status.Errorf(codes.ResourceExhausted, "grpc: message too large (%d bytes)", bufSize)
	}
	return b, nil
}






func compress(in mem.BufferSlice, cp Compressor, compressor encoding.Compressor, pool mem.BufferPool) (mem.BufferSlice, payloadFormat, error) {
	if (compressor == nil && cp == nil) || in.Len() == 0 {
		return nil, compressionNone, nil
	}
	var out mem.BufferSlice
	w := mem.NewWriter(&out, pool)
	wrapErr := func(err error) error {
		out.Free()
		return status.Errorf(codes.Internal, "grpc: error while compressing: %v", err.Error())
	}
	if compressor != nil {
		z, err := compressor.Compress(w)
		if err != nil {
			return nil, 0, wrapErr(err)
		}
		for _, b := range in {
			if _, err := z.Write(b.ReadOnlyData()); err != nil {
				return nil, 0, wrapErr(err)
			}
		}
		if err := z.Close(); err != nil {
			return nil, 0, wrapErr(err)
		}
	} else {
		
		
		
		
		buf := in.MaterializeToBuffer(pool)
		defer buf.Free()
		if err := cp.Do(w, buf.ReadOnlyData()); err != nil {
			return nil, 0, wrapErr(err)
		}
	}
	return out, compressionMade, nil
}

const (
	payloadLen = 1
	sizeLen    = 4
	headerLen  = payloadLen + sizeLen
)



func msgHeader(data, compData mem.BufferSlice, pf payloadFormat) (hdr []byte, payload mem.BufferSlice) {
	hdr = make([]byte, headerLen)
	hdr[0] = byte(pf)

	var length uint32
	if pf.isCompressed() {
		length = uint32(compData.Len())
		payload = compData
	} else {
		length = uint32(data.Len())
		payload = data
	}

	
	binary.BigEndian.PutUint32(hdr[payloadLen:], length)
	return hdr, payload
}

func outPayload(client bool, msg any, dataLength, payloadLength int, t time.Time) *stats.OutPayload {
	return &stats.OutPayload{
		Client:           client,
		Payload:          msg,
		Length:           dataLength,
		WireLength:       payloadLength + headerLen,
		CompressedLength: payloadLength,
		SentTime:         t,
	}
}

func checkRecvPayload(pf payloadFormat, recvCompress string, haveCompressor bool, isServer bool) *status.Status {
	switch pf {
	case compressionNone:
	case compressionMade:
		if recvCompress == "" || recvCompress == encoding.Identity {
			return status.New(codes.Internal, "grpc: compressed flag set with identity or empty encoding")
		}
		if !haveCompressor {
			if isServer {
				return status.Newf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", recvCompress)
			}
			return status.Newf(codes.Internal, "grpc: Decompressor is not installed for grpc-encoding %q", recvCompress)
		}
	default:
		return status.Newf(codes.Internal, "grpc: received unexpected payload format %d", pf)
	}
	return nil
}

type payloadInfo struct {
	compressedLength  int 
	uncompressedBytes mem.BufferSlice
}

func (p *payloadInfo) free() {
	if p != nil && p.uncompressedBytes != nil {
		p.uncompressedBytes.Free()
	}
}







func recvAndDecompress(p *parser, s recvCompressor, dc Decompressor, maxReceiveMessageSize int, payInfo *payloadInfo, compressor encoding.Compressor, isServer bool,
) (out mem.BufferSlice, err error) {
	pf, compressed, err := p.recvMsg(maxReceiveMessageSize)
	if err != nil {
		return nil, err
	}

	compressedLength := compressed.Len()

	if st := checkRecvPayload(pf, s.RecvCompress(), compressor != nil || dc != nil, isServer); st != nil {
		compressed.Free()
		return nil, st.Err()
	}

	if pf.isCompressed() {
		defer compressed.Free()
		
		
		out, err = decompress(compressor, compressed, dc, maxReceiveMessageSize, p.bufferPool)
		if err != nil {
			return nil, err
		}
	} else {
		out = compressed
	}

	if payInfo != nil {
		payInfo.compressedLength = compressedLength
		out.Ref()
		payInfo.uncompressedBytes = out
	}

	return out, nil
}





func decompress(compressor encoding.Compressor, d mem.BufferSlice, dc Decompressor, maxReceiveMessageSize int, pool mem.BufferPool) (mem.BufferSlice, error) {
	if dc != nil {
		uncompressed, err := dc.Do(d.Reader())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "grpc: failed to decompress the received message: %v", err)
		}
		if len(uncompressed) > maxReceiveMessageSize {
			return nil, status.Errorf(codes.ResourceExhausted, "grpc: message after decompression larger than max (%d vs. %d)", len(uncompressed), maxReceiveMessageSize)
		}
		return mem.BufferSlice{mem.SliceBuffer(uncompressed)}, nil
	}
	if compressor != nil {
		dcReader, err := compressor.Decompress(d.Reader())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "grpc: failed to decompress the message: %v", err)
		}

		out, err := mem.ReadAll(io.LimitReader(dcReader, int64(maxReceiveMessageSize)), pool)
		if err != nil {
			out.Free()
			return nil, status.Errorf(codes.Internal, "grpc: failed to read decompressed data: %v", err)
		}

		if out.Len() == maxReceiveMessageSize && !atEOF(dcReader) {
			out.Free()
			return nil, status.Errorf(codes.ResourceExhausted, "grpc: received message after decompression larger than max %d", maxReceiveMessageSize)
		}
		return out, nil
	}
	return nil, status.Errorf(codes.Internal, "grpc: no decompressor available for compressed payload")
}


func atEOF(dcReader io.Reader) bool {
	n, err := dcReader.Read(make([]byte, 1))
	return n == 0 && err == io.EOF
}

type recvCompressor interface {
	RecvCompress() string
}




func recv(p *parser, c baseCodec, s recvCompressor, dc Decompressor, m any, maxReceiveMessageSize int, payInfo *payloadInfo, compressor encoding.Compressor, isServer bool) error {
	data, err := recvAndDecompress(p, s, dc, maxReceiveMessageSize, payInfo, compressor, isServer)
	if err != nil {
		return err
	}

	
	
	defer data.Free()

	if err := c.Unmarshal(data, m); err != nil {
		return status.Errorf(codes.Internal, "grpc: failed to unmarshal the received message: %v", err)
	}

	return nil
}


type rpcInfo struct {
	failfast      bool
	preloaderInfo *compressorInfo
}






type compressorInfo struct {
	codec baseCodec
	cp    Compressor
	comp  encoding.Compressor
}

type rpcInfoContextKey struct{}

func newContextWithRPCInfo(ctx context.Context, failfast bool, codec baseCodec, cp Compressor, comp encoding.Compressor) context.Context {
	return context.WithValue(ctx, rpcInfoContextKey{}, &rpcInfo{
		failfast: failfast,
		preloaderInfo: &compressorInfo{
			codec: codec,
			cp:    cp,
			comp:  comp,
		},
	})
}

func rpcInfoFromContext(ctx context.Context) (s *rpcInfo, ok bool) {
	s, ok = ctx.Value(rpcInfoContextKey{}).(*rpcInfo)
	return
}





func Code(err error) codes.Code {
	return status.Code(err)
}





func ErrorDesc(err error) string {
	return status.Convert(err).Message()
}





func Errorf(c codes.Code, format string, a ...any) error {
	return status.Errorf(c, format, a...)
}

var errContextCanceled = status.Error(codes.Canceled, context.Canceled.Error())
var errContextDeadline = status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error())


func toRPCErr(err error) error {
	switch err {
	case nil, io.EOF:
		return err
	case context.DeadlineExceeded:
		return errContextDeadline
	case context.Canceled:
		return errContextCanceled
	case io.ErrUnexpectedEOF:
		return status.Error(codes.Internal, err.Error())
	}

	switch e := err.(type) {
	case transport.ConnectionError:
		return status.Error(codes.Unavailable, e.Desc)
	case *transport.NewStreamError:
		return toRPCErr(e.Err)
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	return status.Error(codes.Unknown, err.Error())
}


func setCallInfoCodec(c *callInfo) error {
	if c.codec != nil {
		
		
		if c.contentSubtype == "" {
			
			
			
			
			if ec, ok := c.codec.(encoding.CodecV2); ok {
				c.contentSubtype = strings.ToLower(ec.Name())
			}
		}
		return nil
	}

	if c.contentSubtype == "" {
		
		c.codec = getCodec(proto.Name)
		return nil
	}

	
	c.codec = getCodec(c.contentSubtype)
	if c.codec == nil {
		return status.Errorf(codes.Internal, "no codec registered for content-subtype %s", c.contentSubtype)
	}
	return nil
}








const (
	SupportPackageIsVersion3 = true
	SupportPackageIsVersion4 = true
	SupportPackageIsVersion5 = true
	SupportPackageIsVersion6 = true
	SupportPackageIsVersion7 = true
	SupportPackageIsVersion8 = true
	SupportPackageIsVersion9 = true
)

const grpcUA = "grpc-go/" + Version
