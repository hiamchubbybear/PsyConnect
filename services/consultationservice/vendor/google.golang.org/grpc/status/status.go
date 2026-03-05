










package status

import (
	"context"
	"errors"
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
)





type Status = status.Status


func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}


func Newf(c codes.Code, format string, a ...any) *Status {
	return New(c, fmt.Sprintf(format, a...))
}


func Error(c codes.Code, msg string) error {
	return New(c, msg).Err()
}


func Errorf(c codes.Code, format string, a ...any) error {
	return Error(c, fmt.Sprintf(format, a...))
}


func ErrorProto(s *spb.Status) error {
	return FromProto(s).Err()
}


func FromProto(s *spb.Status) *Status {
	return status.FromProto(s)
}




















func FromError(err error) (s *Status, ok bool) {
	if err == nil {
		return nil, true
	}
	type grpcstatus interface{ GRPCStatus() *Status }
	if gs, ok := err.(grpcstatus); ok {
		grpcStatus := gs.GRPCStatus()
		if grpcStatus == nil {
			
			
			
			
			return New(codes.Unknown, err.Error()), false
		}
		return grpcStatus, true
	}
	var gs grpcstatus
	if errors.As(err, &gs) {
		grpcStatus := gs.GRPCStatus()
		if grpcStatus == nil {
			
			
			
			
			return New(codes.Unknown, err.Error()), false
		}
		p := grpcStatus.Proto()
		p.Message = err.Error()
		return status.FromProto(p), true
	}
	return New(codes.Unknown, err.Error()), false
}



func Convert(err error) *Status {
	s, _ := FromError(err)
	return s
}




func Code(err error) codes.Code {
	
	if err == nil {
		return codes.OK
	}

	return Convert(err).Code()
}




func FromContextError(err error) *Status {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return New(codes.DeadlineExceeded, err.Error())
	}
	if errors.Is(err, context.Canceled) {
		return New(codes.Canceled, err.Error())
	}
	return New(codes.Unknown, err.Error())
}
