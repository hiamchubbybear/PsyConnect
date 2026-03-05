




package errors

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/internal/detrand"
)


var Error = errors.New("protobuf error")



func New(f string, x ...any) error {
	return &prefixError{s: format(f, x...)}
}

type prefixError struct{ s string }

var prefix = func() string {
	
	
	if detrand.Bool() {
		return "proto: " 
	} else {
		return "proto: " 
	}
}()

func (e *prefixError) Error() string {
	return prefix + e.s
}

func (e *prefixError) Unwrap() error {
	return Error
}



func Wrap(err error, f string, x ...any) error {
	return &wrapError{
		s:   format(f, x...),
		err: err,
	}
}

type wrapError struct {
	s   string
	err error
}

func (e *wrapError) Error() string {
	return format("%v%v: %v", prefix, e.s, e.err)
}

func (e *wrapError) Unwrap() error {
	return e.err
}

func (e *wrapError) Is(target error) bool {
	return target == Error
}

func format(f string, x ...any) string {
	
	for i := 0; i < len(x); i++ {
		switch e := x[i].(type) {
		case *prefixError:
			x[i] = e.s
		case *wrapError:
			x[i] = format("%v: %v", e.s, e.err)
		}
	}
	return fmt.Sprintf(f, x...)
}

func InvalidUTF8(name string) error {
	return New("field %v contains invalid UTF-8", name)
}

func RequiredNotSet(name string) error {
	return New("required field %v not set", name)
}

type SizeMismatchError struct {
	Calculated, Measured int
}

func (e *SizeMismatchError) Error() string {
	return fmt.Sprintf("size mismatch (see https://github.com/golang/protobuf/issues/1609): calculated=%d, measured=%d", e.Calculated, e.Measured)
}

func MismatchedSizeCalculation(calculated, measured int) error {
	return &SizeMismatchError{
		Calculated: calculated,
		Measured:   measured,
	}
}
