





package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
)



var ErrNoDocuments = errors.New("mongo: no documents in result")




type SingleResult struct {
	ctx      context.Context
	err      error
	cur      *Cursor
	rdr      bson.Raw
	bsonOpts *options.BSONOptions
	reg      *bsoncodec.Registry
}






func NewSingleResultFromDocument(document interface{}, err error, registry *bsoncodec.Registry) *SingleResult {
	if document == nil {
		return &SingleResult{err: ErrNilDocument}
	}
	if registry == nil {
		registry = bson.DefaultRegistry
	}

	cur, createErr := NewCursorFromDocuments([]interface{}{document}, err, registry)
	if createErr != nil {
		return &SingleResult{err: createErr}
	}

	return &SingleResult{
		cur: cur,
		err: err,
		reg: registry,
	}
}







func (sr *SingleResult) Decode(v interface{}) error {
	if sr.err != nil {
		return sr.err
	}
	if sr.reg == nil {
		return bson.ErrNilRegistry
	}

	if sr.err = sr.setRdrContents(); sr.err != nil {
		return sr.err
	}

	dec, err := getDecoder(sr.rdr, sr.bsonOpts, sr.reg)
	if err != nil {
		return fmt.Errorf("error configuring BSON decoder: %w", err)
	}

	return dec.Decode(v)
}





func (sr *SingleResult) Raw() (bson.Raw, error) {
	if sr.err != nil {
		return sr.rdr, sr.err
	}

	if sr.err = sr.setRdrContents(); sr.err != nil {
		return nil, sr.err
	}
	return sr.rdr, nil
}






func (sr *SingleResult) DecodeBytes() (bson.Raw, error) {
	return sr.Raw()
}


func (sr *SingleResult) setRdrContents() error {
	switch {
	case sr.err != nil:
		return sr.err
	case sr.rdr != nil:
		return nil
	case sr.cur != nil:
		defer sr.cur.Close(sr.ctx)

		if !sr.cur.Next(sr.ctx) {
			if err := sr.cur.Err(); err != nil {
				return err
			}

			return ErrNoDocuments
		}
		sr.rdr = sr.cur.Current
		return nil
	}

	return ErrNoDocuments
}





func (sr *SingleResult) Err() error {
	sr.err = sr.setRdrContents()

	return sr.err
}
