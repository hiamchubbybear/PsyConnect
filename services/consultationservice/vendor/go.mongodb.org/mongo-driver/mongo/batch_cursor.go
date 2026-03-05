





package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)



type batchCursor interface {
	
	ID() int64

	
	Next(context.Context) bool

	
	
	Batch() *bsoncore.DocumentSequence

	
	Server() driver.Server

	
	Err() error

	
	Close(context.Context) error

	
	
	SetBatchSize(int32)

	
	
	
	
	
	
	SetMaxTime(time.Duration)

	
	
	SetComment(interface{})
}



type changeStreamCursor interface {
	batchCursor
	
	PostBatchResumeToken() bsoncore.Document

	
	KillCursor(context.Context) error
}
