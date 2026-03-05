package lz4

import (
	"errors"
	"fmt"
	"io"

	"github.com/pierrec/lz4/v4/internal/lz4errors"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=aState -output state_gen.go

const (
	noState     aState = iota 
	errorState                
	newState                  
	readState                 
	writeState                
	closedState               
)

type (
	aState uint8
	_State struct {
		states []aState
		state  aState
		err    error
	}
)

func (s *_State) init(states []aState) {
	s.states = states
	s.state = states[0]
}

func (s *_State) reset() {
	s.state = s.states[0]
	s.err = nil
}



func (s *_State) next(err error) bool {
	if err != nil {
		s.err = fmt.Errorf("%s: %w", s.state, err)
		s.state = errorState
		return true
	}
	s.state = s.states[s.state]
	return false
}


func (s *_State) nextd(errp *error) bool {
	return errp != nil && s.next(*errp)
}


func (s *_State) check(errp *error) {
	if s.state == errorState || errp == nil {
		return
	}
	if err := *errp; err != nil {
		s.err = fmt.Errorf("%w[%s]", err, s.state)
		if !errors.Is(err, io.EOF) {
			s.state = errorState
		}
	}
}

func (s *_State) fail() error {
	s.state = errorState
	s.err = fmt.Errorf("%w[%s]", lz4errors.ErrInternalUnhandledState, s.state)
	return s.err
}
