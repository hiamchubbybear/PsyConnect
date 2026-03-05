





package options

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)


type TransactionOptions struct {
	
	
	ReadConcern *readconcern.ReadConcern

	
	
	ReadPreference *readpref.ReadPref

	
	
	WriteConcern *writeconcern.WriteConcern

	
	

	
	
	
	
	
	
	
	MaxCommitTime *time.Duration
}


func Transaction() *TransactionOptions {
	return &TransactionOptions{}
}


func (t *TransactionOptions) SetReadConcern(rc *readconcern.ReadConcern) *TransactionOptions {
	t.ReadConcern = rc
	return t
}


func (t *TransactionOptions) SetReadPreference(rp *readpref.ReadPref) *TransactionOptions {
	t.ReadPreference = rp
	return t
}


func (t *TransactionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *TransactionOptions {
	t.WriteConcern = wc
	return t
}






func (t *TransactionOptions) SetMaxCommitTime(mct *time.Duration) *TransactionOptions {
	t.MaxCommitTime = mct
	return t
}






func MergeTransactionOptions(opts ...*TransactionOptions) *TransactionOptions {
	t := Transaction()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			t.ReadConcern = opt.ReadConcern
		}
		if opt.ReadPreference != nil {
			t.ReadPreference = opt.ReadPreference
		}
		if opt.WriteConcern != nil {
			t.WriteConcern = opt.WriteConcern
		}
		if opt.MaxCommitTime != nil {
			t.MaxCommitTime = opt.MaxCommitTime
		}
	}

	return t
}
