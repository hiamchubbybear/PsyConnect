





package options

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)


type DatabaseOptions struct {
	
	
	ReadConcern *readconcern.ReadConcern

	
	
	WriteConcern *writeconcern.WriteConcern

	
	
	ReadPreference *readpref.ReadPref

	
	
	BSONOptions *BSONOptions

	
	
	Registry *bsoncodec.Registry
}


func Database() *DatabaseOptions {
	return &DatabaseOptions{}
}


func (d *DatabaseOptions) SetReadConcern(rc *readconcern.ReadConcern) *DatabaseOptions {
	d.ReadConcern = rc
	return d
}


func (d *DatabaseOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *DatabaseOptions {
	d.WriteConcern = wc
	return d
}


func (d *DatabaseOptions) SetReadPreference(rp *readpref.ReadPref) *DatabaseOptions {
	d.ReadPreference = rp
	return d
}


func (d *DatabaseOptions) SetBSONOptions(opts *BSONOptions) *DatabaseOptions {
	d.BSONOptions = opts
	return d
}


func (d *DatabaseOptions) SetRegistry(r *bsoncodec.Registry) *DatabaseOptions {
	d.Registry = r
	return d
}






func MergeDatabaseOptions(opts ...*DatabaseOptions) *DatabaseOptions {
	d := Database()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadConcern != nil {
			d.ReadConcern = opt.ReadConcern
		}
		if opt.WriteConcern != nil {
			d.WriteConcern = opt.WriteConcern
		}
		if opt.ReadPreference != nil {
			d.ReadPreference = opt.ReadPreference
		}
		if opt.Registry != nil {
			d.Registry = opt.Registry
		}
		if opt.BSONOptions != nil {
			d.BSONOptions = opt.BSONOptions
		}
	}

	return d
}
