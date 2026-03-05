





package session

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)


type ClientOptions struct {
	CausalConsistency     *bool
	DefaultReadConcern    *readconcern.ReadConcern
	DefaultWriteConcern   *writeconcern.WriteConcern
	DefaultReadPreference *readpref.ReadPref
	DefaultMaxCommitTime  *time.Duration
	Snapshot              *bool
}


type TransactionOptions struct {
	ReadConcern    *readconcern.ReadConcern
	WriteConcern   *writeconcern.WriteConcern
	ReadPreference *readpref.ReadPref
	MaxCommitTime  *time.Duration
}

func mergeClientOptions(opts ...*ClientOptions) *ClientOptions {
	c := &ClientOptions{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.CausalConsistency != nil {
			c.CausalConsistency = opt.CausalConsistency
		}
		if opt.DefaultReadConcern != nil {
			c.DefaultReadConcern = opt.DefaultReadConcern
		}
		if opt.DefaultReadPreference != nil {
			c.DefaultReadPreference = opt.DefaultReadPreference
		}
		if opt.DefaultWriteConcern != nil {
			c.DefaultWriteConcern = opt.DefaultWriteConcern
		}
		if opt.DefaultMaxCommitTime != nil {
			c.DefaultMaxCommitTime = opt.DefaultMaxCommitTime
		}
		if opt.Snapshot != nil {
			c.Snapshot = opt.Snapshot
		}
	}

	return c
}
