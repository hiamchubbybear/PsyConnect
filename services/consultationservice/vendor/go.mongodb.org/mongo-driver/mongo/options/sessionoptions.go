





package options

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)


var DefaultCausalConsistency = true


type SessionOptions struct {
	
	
	
	CausalConsistency *bool

	
	
	DefaultReadConcern *readconcern.ReadConcern

	
	
	DefaultReadPreference *readpref.ReadPref

	
	
	DefaultWriteConcern *writeconcern.WriteConcern

	
	
	
	
	
	
	DefaultMaxCommitTime *time.Duration

	
	
	
	Snapshot *bool
}


func Session() *SessionOptions {
	return &SessionOptions{}
}


func (s *SessionOptions) SetCausalConsistency(b bool) *SessionOptions {
	s.CausalConsistency = &b
	return s
}


func (s *SessionOptions) SetDefaultReadConcern(rc *readconcern.ReadConcern) *SessionOptions {
	s.DefaultReadConcern = rc
	return s
}


func (s *SessionOptions) SetDefaultReadPreference(rp *readpref.ReadPref) *SessionOptions {
	s.DefaultReadPreference = rp
	return s
}


func (s *SessionOptions) SetDefaultWriteConcern(wc *writeconcern.WriteConcern) *SessionOptions {
	s.DefaultWriteConcern = wc
	return s
}







func (s *SessionOptions) SetDefaultMaxCommitTime(mct *time.Duration) *SessionOptions {
	s.DefaultMaxCommitTime = mct
	return s
}


func (s *SessionOptions) SetSnapshot(b bool) *SessionOptions {
	s.Snapshot = &b
	return s
}






func MergeSessionOptions(opts ...*SessionOptions) *SessionOptions {
	s := Session()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.CausalConsistency != nil {
			s.CausalConsistency = opt.CausalConsistency
		}
		if opt.DefaultReadConcern != nil {
			s.DefaultReadConcern = opt.DefaultReadConcern
		}
		if opt.DefaultReadPreference != nil {
			s.DefaultReadPreference = opt.DefaultReadPreference
		}
		if opt.DefaultWriteConcern != nil {
			s.DefaultWriteConcern = opt.DefaultWriteConcern
		}
		if opt.DefaultMaxCommitTime != nil {
			s.DefaultMaxCommitTime = opt.DefaultMaxCommitTime
		}
		if opt.Snapshot != nil {
			s.Snapshot = opt.Snapshot
		}
	}
	if s.CausalConsistency == nil && (s.Snapshot == nil || !*s.Snapshot) {
		s.CausalConsistency = &DefaultCausalConsistency
	}

	return s
}
