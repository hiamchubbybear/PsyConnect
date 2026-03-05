





package mongo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)




type Cursor struct {
	
	
	Current bson.Raw

	bc            batchCursor
	batch         *bsoncore.DocumentSequence
	batchLength   int
	bsonOpts      *options.BSONOptions
	registry      *bsoncodec.Registry
	clientSession *session.Client

	err error
}

func newCursor(
	bc batchCursor,
	bsonOpts *options.BSONOptions,
	registry *bsoncodec.Registry,
) (*Cursor, error) {
	return newCursorWithSession(bc, bsonOpts, registry, nil)
}

func newCursorWithSession(
	bc batchCursor,
	bsonOpts *options.BSONOptions,
	registry *bsoncodec.Registry,
	clientSession *session.Client,
) (*Cursor, error) {
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	if bc == nil {
		return nil, errors.New("batch cursor must not be nil")
	}
	c := &Cursor{
		bc:            bc,
		bsonOpts:      bsonOpts,
		registry:      registry,
		clientSession: clientSession,
	}
	if bc.ID() == 0 {
		c.closeImplicitSession()
	}

	
	
	c.batchLength = c.bc.Batch().DocumentCount()
	return c, nil
}

func newEmptyCursor() *Cursor {
	return &Cursor{bc: driver.NewEmptyBatchCursor()}
}





func NewCursorFromDocuments(documents []interface{}, err error, registry *bsoncodec.Registry) (*Cursor, error) {
	if registry == nil {
		registry = bson.DefaultRegistry
	}

	
	var docsBytes []byte
	for _, doc := range documents {
		switch t := doc.(type) {
		case nil:
			return nil, ErrNilDocument
		case []byte:
			
			doc = bson.Raw(t)
		}
		var marshalErr error
		docsBytes, marshalErr = bson.MarshalAppendWithRegistry(registry, docsBytes, doc)
		if marshalErr != nil {
			return nil, marshalErr
		}
	}

	c := &Cursor{
		bc:       driver.NewBatchCursorFromDocuments(docsBytes),
		registry: registry,
		err:      err,
	}

	
	
	c.batch = c.bc.Batch()
	c.batchLength = c.bc.Batch().DocumentCount()
	return c, nil
}


func (c *Cursor) ID() int64 { return c.bc.ID() }








func (c *Cursor) Next(ctx context.Context) bool {
	return c.next(ctx, false)
}













func (c *Cursor) TryNext(ctx context.Context) bool {
	return c.next(ctx, true)
}

func (c *Cursor) next(ctx context.Context, nonBlocking bool) bool {
	
	if c.err != nil {
		return false
	}

	if ctx == nil {
		ctx = context.Background()
	}
	doc, err := c.batch.Next()
	switch {
	case err == nil:
		
		c.batchLength--
		c.Current = bson.Raw(doc)
		return true
	case errors.Is(err, io.EOF): 
	default:
		c.err = err
		return false
	}

	
	
	for {
		
		if !c.bc.Next(ctx) {
			
			c.err = replaceErrors(c.bc.Err())
			if c.err != nil {
				return false
			}
			
			if c.bc.ID() == 0 {
				c.closeImplicitSession()
				return false
			}
			
			
			if nonBlocking {
				return false
			}
			continue
		}

		
		if c.bc.ID() == 0 {
			c.closeImplicitSession()
		}

		
		c.batch = c.bc.Batch()
		c.batchLength = c.batch.DocumentCount()
		doc, err = c.batch.Next()
		switch {
		case err == nil:
			c.batchLength--
			c.Current = bson.Raw(doc)
			return true
		case errors.Is(err, io.EOF): 
		default:
			c.err = err
			return false
		}
	}
}

func getDecoder(
	data []byte,
	opts *options.BSONOptions,
	reg *bsoncodec.Registry,
) (*bson.Decoder, error) {
	dec, err := bson.NewDecoder(bsonrw.NewBSONDocumentReader(data))
	if err != nil {
		return nil, err
	}

	if opts != nil {
		if opts.AllowTruncatingDoubles {
			dec.AllowTruncatingDoubles()
		}
		if opts.BinaryAsSlice {
			dec.BinaryAsSlice()
		}
		if opts.DefaultDocumentD {
			dec.DefaultDocumentD()
		}
		if opts.DefaultDocumentM {
			dec.DefaultDocumentM()
		}
		if opts.UseJSONStructTags {
			dec.UseJSONStructTags()
		}
		if opts.UseLocalTimeZone {
			dec.UseLocalTimeZone()
		}
		if opts.ZeroMaps {
			dec.ZeroMaps()
		}
		if opts.ZeroStructs {
			dec.ZeroStructs()
		}
	}

	if reg != nil {
		
		if err := dec.SetRegistry(reg); err != nil {
			return nil, err
		}
	}

	return dec, nil
}



func (c *Cursor) Decode(val interface{}) error {
	dec, err := getDecoder(c.Current, c.bsonOpts, c.registry)
	if err != nil {
		return fmt.Errorf("error configuring BSON decoder: %w", err)
	}

	return dec.Decode(val)
}


func (c *Cursor) Err() error { return c.err }



func (c *Cursor) Close(ctx context.Context) error {
	defer c.closeImplicitSession()
	return replaceErrors(c.bc.Close(ctx))
}







func (c *Cursor) All(ctx context.Context, results interface{}) error {
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}

	elementType := sliceVal.Type().Elem()
	var index int
	var err error

	
	
	
	defer c.Close(context.Background())

	batch := c.batch 
	for {
		sliceVal, index, err = c.addFromBatch(sliceVal, elementType, batch, index)
		if err != nil {
			return err
		}

		if !c.bc.Next(ctx) {
			break
		}

		batch = c.bc.Batch()
	}

	if err = replaceErrors(c.bc.Err()); err != nil {
		return err
	}

	resultsVal.Elem().Set(sliceVal.Slice(0, index))
	return nil
}



func (c *Cursor) RemainingBatchLength() int {
	return c.batchLength
}



func (c *Cursor) addFromBatch(sliceVal reflect.Value, elemType reflect.Type, batch *bsoncore.DocumentSequence,
	index int) (reflect.Value, int, error) {

	docs, err := batch.Documents()
	if err != nil {
		return sliceVal, index, err
	}

	for _, doc := range docs {
		if sliceVal.Len() == index {
			
			newElem := reflect.New(elemType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		}

		currElem := sliceVal.Index(index).Addr().Interface()
		dec, err := getDecoder(doc, c.bsonOpts, c.registry)
		if err != nil {
			return sliceVal, index, fmt.Errorf("error configuring BSON decoder: %w", err)
		}
		err = dec.Decode(currElem)
		if err != nil {
			return sliceVal, index, err
		}

		index++
	}

	return sliceVal, index, nil
}

func (c *Cursor) closeImplicitSession() {
	if c.clientSession != nil && c.clientSession.IsImplicit {
		c.clientSession.EndSession()
	}
}





func (c *Cursor) SetBatchSize(batchSize int32) {
	c.bc.SetBatchSize(batchSize)
}







func (c *Cursor) SetMaxTime(dur time.Duration) {
	c.bc.SetMaxTime(dur)
}



func (c *Cursor) SetComment(comment interface{}) {
	c.bc.SetComment(comment)
}






func BatchCursorFromCursor(c *Cursor) *driver.BatchCursor {
	bc, _ := c.bc.(*driver.BatchCursor)
	return bc
}
