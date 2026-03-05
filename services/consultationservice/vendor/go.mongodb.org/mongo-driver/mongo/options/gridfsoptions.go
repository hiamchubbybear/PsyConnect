





package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)


var DefaultName = "fs"


var DefaultChunkSize int32 = 255 * 1024


var DefaultRevision int32 = -1


type BucketOptions struct {
	
	Name *string

	
	ChunkSizeBytes *int32

	
	
	WriteConcern *writeconcern.WriteConcern

	
	
	ReadConcern *readconcern.ReadConcern

	
	
	ReadPreference *readpref.ReadPref
}


func GridFSBucket() *BucketOptions {
	return &BucketOptions{
		Name:           &DefaultName,
		ChunkSizeBytes: &DefaultChunkSize,
	}
}


func (b *BucketOptions) SetName(name string) *BucketOptions {
	b.Name = &name
	return b
}


func (b *BucketOptions) SetChunkSizeBytes(i int32) *BucketOptions {
	b.ChunkSizeBytes = &i
	return b
}


func (b *BucketOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *BucketOptions {
	b.WriteConcern = wc
	return b
}


func (b *BucketOptions) SetReadConcern(rc *readconcern.ReadConcern) *BucketOptions {
	b.ReadConcern = rc
	return b
}


func (b *BucketOptions) SetReadPreference(rp *readpref.ReadPref) *BucketOptions {
	b.ReadPreference = rp
	return b
}





func MergeBucketOptions(opts ...*BucketOptions) *BucketOptions {
	b := GridFSBucket()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Name != nil {
			b.Name = opt.Name
		}
		if opt.ChunkSizeBytes != nil {
			b.ChunkSizeBytes = opt.ChunkSizeBytes
		}
		if opt.WriteConcern != nil {
			b.WriteConcern = opt.WriteConcern
		}
		if opt.ReadConcern != nil {
			b.ReadConcern = opt.ReadConcern
		}
		if opt.ReadPreference != nil {
			b.ReadPreference = opt.ReadPreference
		}
	}

	return b
}


type UploadOptions struct {
	
	ChunkSizeBytes *int32

	
	
	
	Metadata interface{}

	
	Registry *bsoncodec.Registry
}


func GridFSUpload() *UploadOptions {
	return &UploadOptions{Registry: bson.DefaultRegistry}
}


func (u *UploadOptions) SetChunkSizeBytes(i int32) *UploadOptions {
	u.ChunkSizeBytes = &i
	return u
}


func (u *UploadOptions) SetMetadata(doc interface{}) *UploadOptions {
	u.Metadata = doc
	return u
}





func MergeUploadOptions(opts ...*UploadOptions) *UploadOptions {
	u := GridFSUpload()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ChunkSizeBytes != nil {
			u.ChunkSizeBytes = opt.ChunkSizeBytes
		}
		if opt.Metadata != nil {
			u.Metadata = opt.Metadata
		}
		if opt.Registry != nil {
			u.Registry = opt.Registry
		}
	}

	return u
}


type NameOptions struct {
	
	
	
	
	
	
	
	
	
	
	Revision *int32
}


func GridFSName() *NameOptions {
	return &NameOptions{}
}


func (n *NameOptions) SetRevision(r int32) *NameOptions {
	n.Revision = &r
	return n
}





func MergeNameOptions(opts ...*NameOptions) *NameOptions {
	n := GridFSName()
	n.Revision = &DefaultRevision

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Revision != nil {
			n.Revision = opt.Revision
		}
	}

	return n
}


type GridFSFindOptions struct {
	
	
	
	AllowDiskUse *bool

	
	BatchSize *int32

	
	
	
	Limit *int32

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	NoCursorTimeout *bool

	
	Skip *int32

	
	
	Sort interface{}
}


func GridFSFind() *GridFSFindOptions {
	return &GridFSFindOptions{}
}


func (f *GridFSFindOptions) SetAllowDiskUse(b bool) *GridFSFindOptions {
	f.AllowDiskUse = &b
	return f
}


func (f *GridFSFindOptions) SetBatchSize(i int32) *GridFSFindOptions {
	f.BatchSize = &i
	return f
}


func (f *GridFSFindOptions) SetLimit(i int32) *GridFSFindOptions {
	f.Limit = &i
	return f
}






func (f *GridFSFindOptions) SetMaxTime(d time.Duration) *GridFSFindOptions {
	f.MaxTime = &d
	return f
}


func (f *GridFSFindOptions) SetNoCursorTimeout(b bool) *GridFSFindOptions {
	f.NoCursorTimeout = &b
	return f
}


func (f *GridFSFindOptions) SetSkip(i int32) *GridFSFindOptions {
	f.Skip = &i
	return f
}


func (f *GridFSFindOptions) SetSort(sort interface{}) *GridFSFindOptions {
	f.Sort = sort
	return f
}






func MergeGridFSFindOptions(opts ...*GridFSFindOptions) *GridFSFindOptions {
	fo := GridFSFind()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowDiskUse != nil {
			fo.AllowDiskUse = opt.AllowDiskUse
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Limit != nil {
			fo.Limit = opt.Limit
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}
