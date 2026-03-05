





package options

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)


type CollectionOptions struct {
	
	
	ReadConcern *readconcern.ReadConcern

	
	
	WriteConcern *writeconcern.WriteConcern

	
	
	ReadPreference *readpref.ReadPref

	
	
	BSONOptions *BSONOptions

	
	
	Registry *bsoncodec.Registry
}


func Collection() *CollectionOptions {
	return &CollectionOptions{}
}


func (c *CollectionOptions) SetReadConcern(rc *readconcern.ReadConcern) *CollectionOptions {
	c.ReadConcern = rc
	return c
}


func (c *CollectionOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *CollectionOptions {
	c.WriteConcern = wc
	return c
}


func (c *CollectionOptions) SetReadPreference(rp *readpref.ReadPref) *CollectionOptions {
	c.ReadPreference = rp
	return c
}


func (c *CollectionOptions) SetBSONOptions(opts *BSONOptions) *CollectionOptions {
	c.BSONOptions = opts
	return c
}


func (c *CollectionOptions) SetRegistry(r *bsoncodec.Registry) *CollectionOptions {
	c.Registry = r
	return c
}






func MergeCollectionOptions(opts ...*CollectionOptions) *CollectionOptions {
	c := Collection()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			c.ReadConcern = opt.ReadConcern
		}
		if opt.WriteConcern != nil {
			c.WriteConcern = opt.WriteConcern
		}
		if opt.ReadPreference != nil {
			c.ReadPreference = opt.ReadPreference
		}
		if opt.Registry != nil {
			c.Registry = opt.Registry
		}
		if opt.BSONOptions != nil {
			c.BSONOptions = opt.BSONOptions
		}
	}

	return c
}
