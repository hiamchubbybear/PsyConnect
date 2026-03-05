




package stats 

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/metadata"
)


type RPCStats interface {
	isRPCStats()
	
	IsClient() bool
}



type Begin struct {
	
	Client bool
	
	BeginTime time.Time
	
	FailFast bool
	
	IsClientStream bool
	
	IsServerStream bool
	
	
	IsTransparentRetryAttempt bool
}


func (s *Begin) IsClient() bool { return s.Client }

func (s *Begin) isRPCStats() {}



type PickerUpdated struct{}



func (*PickerUpdated) IsClient() bool { return true }

func (*PickerUpdated) isRPCStats() {}


type InPayload struct {
	
	Client bool
	
	
	
	Payload any

	
	
	Length int
	
	
	
	CompressedLength int
	
	
	WireLength int

	
	RecvTime time.Time
}


func (s *InPayload) IsClient() bool { return s.Client }

func (s *InPayload) isRPCStats() {}


type InHeader struct {
	
	Client bool
	
	WireLength int
	
	Compression string
	
	Header metadata.MD

	
	
	FullMethod string
	
	RemoteAddr net.Addr
	
	LocalAddr net.Addr
}


func (s *InHeader) IsClient() bool { return s.Client }

func (s *InHeader) isRPCStats() {}


type InTrailer struct {
	
	Client bool
	
	WireLength int
	
	
	Trailer metadata.MD
}


func (s *InTrailer) IsClient() bool { return s.Client }

func (s *InTrailer) isRPCStats() {}


type OutPayload struct {
	
	Client bool
	
	
	
	Payload any
	
	
	Length int
	
	
	
	CompressedLength int
	
	
	WireLength int
	
	SentTime time.Time
}


func (s *OutPayload) IsClient() bool { return s.Client }

func (s *OutPayload) isRPCStats() {}


type OutHeader struct {
	
	Client bool
	
	Compression string
	
	Header metadata.MD

	
	
	FullMethod string
	
	RemoteAddr net.Addr
	
	LocalAddr net.Addr
}


func (s *OutHeader) IsClient() bool { return s.Client }

func (s *OutHeader) isRPCStats() {}


type OutTrailer struct {
	
	Client bool
	
	
	
	
	WireLength int
	
	
	Trailer metadata.MD
}


func (s *OutTrailer) IsClient() bool { return s.Client }

func (s *OutTrailer) isRPCStats() {}


type End struct {
	
	Client bool
	
	BeginTime time.Time
	
	EndTime time.Time
	
	
	
	Trailer metadata.MD
	
	
	
	Error error
}


func (s *End) IsClient() bool { return s.Client }

func (s *End) isRPCStats() {}


type ConnStats interface {
	isConnStats()
	
	IsClient() bool
}


type ConnBegin struct {
	
	Client bool
}


func (s *ConnBegin) IsClient() bool { return s.Client }

func (s *ConnBegin) isConnStats() {}


type ConnEnd struct {
	
	Client bool
}


func (s *ConnEnd) IsClient() bool { return s.Client }

func (s *ConnEnd) isConnStats() {}






func SetTags(ctx context.Context, b []byte) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "grpc-tags-bin", string(b))
}




func Tags(ctx context.Context) []byte {
	traceValues := metadata.ValueFromIncomingContext(ctx, "grpc-tags-bin")
	if len(traceValues) == 0 {
		return nil
	}
	return []byte(traceValues[len(traceValues)-1])
}






func SetTrace(ctx context.Context, b []byte) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(b))
}




func Trace(ctx context.Context) []byte {
	traceValues := metadata.ValueFromIncomingContext(ctx, "grpc-trace-bin")
	if len(traceValues) == 0 {
		return nil
	}
	return []byte(traceValues[len(traceValues)-1])
}
