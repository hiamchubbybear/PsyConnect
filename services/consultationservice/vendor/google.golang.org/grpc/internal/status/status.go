










package status

import (
	"errors"
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/types/known/anypb"
)



type Status struct {
	s *spb.Status
}



func NewWithProto(code codes.Code, message string, statusProto []string) *Status {
	if len(statusProto) != 1 {
		
		return &Status{s: &spb.Status{Code: int32(code), Message: message}}
	}
	st := &spb.Status{}
	if err := proto.Unmarshal([]byte(statusProto[0]), st); err != nil {
		
		return &Status{s: &spb.Status{Code: int32(code), Message: message}}
	}
	if st.Code == int32(code) {
		
		
		return &Status{s: st}
	}
	return &Status{
		s: &spb.Status{
			Code: int32(codes.Internal),
			Message: fmt.Sprintf(
				"grpc-status-details-bin mismatch: grpc-status=%v, grpc-message=%q, grpc-status-details-bin=%+v",
				code, message, st,
			),
		},
	}
}


func New(c codes.Code, msg string) *Status {
	return &Status{s: &spb.Status{Code: int32(c), Message: msg}}
}


func Newf(c codes.Code, format string, a ...any) *Status {
	return New(c, fmt.Sprintf(format, a...))
}


func FromProto(s *spb.Status) *Status {
	return &Status{s: proto.Clone(s).(*spb.Status)}
}


func Err(c codes.Code, msg string) error {
	return New(c, msg).Err()
}


func Errorf(c codes.Code, format string, a ...any) error {
	return Err(c, fmt.Sprintf(format, a...))
}


func (s *Status) Code() codes.Code {
	if s == nil || s.s == nil {
		return codes.OK
	}
	return codes.Code(s.s.Code)
}


func (s *Status) Message() string {
	if s == nil || s.s == nil {
		return ""
	}
	return s.s.Message
}


func (s *Status) Proto() *spb.Status {
	if s == nil {
		return nil
	}
	return proto.Clone(s.s).(*spb.Status)
}


func (s *Status) Err() error {
	if s.Code() == codes.OK {
		return nil
	}
	return &Error{s: s}
}



func (s *Status) WithDetails(details ...protoadapt.MessageV1) (*Status, error) {
	if s.Code() == codes.OK {
		return nil, errors.New("no error details for status with code OK")
	}
	
	p := s.Proto()
	for _, detail := range details {
		m, err := anypb.New(protoadapt.MessageV2Of(detail))
		if err != nil {
			return nil, err
		}
		p.Details = append(p.Details, m)
	}
	return &Status{s: p}, nil
}





func (s *Status) Details() []any {
	if s == nil || s.s == nil {
		return nil
	}
	details := make([]any, 0, len(s.s.Details))
	for _, any := range s.s.Details {
		detail, err := any.UnmarshalNew()
		if err != nil {
			details = append(details, err)
			continue
		}
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		details = append(details, protoadapt.MessageV1Of(detail))
	}
	return details
}

func (s *Status) String() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", s.Code(), s.Message())
}



type Error struct {
	s *Status
}

func (e *Error) Error() string {
	return e.s.String()
}


func (e *Error) GRPCStatus() *Status {
	return e.s
}



func (e *Error) Is(target error) bool {
	tse, ok := target.(*Error)
	if !ok {
		return false
	}
	return proto.Equal(e.s.s, tse.s.s)
}



func IsRestrictedControlPlaneCode(s *Status) bool {
	switch s.Code() {
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.DataLoss:
		return true
	}
	return false
}
