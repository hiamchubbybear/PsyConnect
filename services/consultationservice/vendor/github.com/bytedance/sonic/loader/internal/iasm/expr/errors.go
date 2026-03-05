















package expr

import (
	"fmt"
)


type SyntaxError struct {
	Pos    int
	Reason string
}

func newSyntaxError(pos int, reason string) *SyntaxError {
	return &SyntaxError{
		Pos:    pos,
		Reason: reason,
	}
}

func (self *SyntaxError) Error() string {
	return fmt.Sprintf("Syntax error at position %d: %s", self.Pos, self.Reason)
}


type RuntimeError struct {
	Reason string
}

func newRuntimeError(reason string) *RuntimeError {
	return &RuntimeError{
		Reason: reason,
	}
}

func (self *RuntimeError) Error() string {
	return "Runtime error: " + self.Reason
}
