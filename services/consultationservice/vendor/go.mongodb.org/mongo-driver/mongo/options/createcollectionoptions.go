





package options

import "time"



type DefaultIndexOptions struct {
	
	
	
	StorageEngine interface{}
}


func DefaultIndex() *DefaultIndexOptions {
	return &DefaultIndexOptions{}
}


func (d *DefaultIndexOptions) SetStorageEngine(storageEngine interface{}) *DefaultIndexOptions {
	d.StorageEngine = storageEngine
	return d
}


type TimeSeriesOptions struct {
	
	
	TimeField string

	
	
	
	MetaField *string

	
	
	Granularity *string

	
	
	
	BucketMaxSpan *time.Duration

	
	
	
	
	
	BucketRounding *time.Duration
}


func TimeSeries() *TimeSeriesOptions {
	return &TimeSeriesOptions{}
}


func (tso *TimeSeriesOptions) SetTimeField(timeField string) *TimeSeriesOptions {
	tso.TimeField = timeField
	return tso
}


func (tso *TimeSeriesOptions) SetMetaField(metaField string) *TimeSeriesOptions {
	tso.MetaField = &metaField
	return tso
}


func (tso *TimeSeriesOptions) SetGranularity(granularity string) *TimeSeriesOptions {
	tso.Granularity = &granularity
	return tso
}


func (tso *TimeSeriesOptions) SetBucketMaxSpan(dur time.Duration) *TimeSeriesOptions {
	tso.BucketMaxSpan = &dur

	return tso
}


func (tso *TimeSeriesOptions) SetBucketRounding(dur time.Duration) *TimeSeriesOptions {
	tso.BucketRounding = &dur

	return tso
}


type CreateCollectionOptions struct {
	
	
	Capped *bool

	
	
	Collation *Collation

	
	
	
	
	ChangeStreamPreAndPostImages interface{}

	
	
	DefaultIndexOptions *DefaultIndexOptions

	
	
	
	
	MaxDocuments *int64

	
	SizeInBytes *int64

	
	
	
	StorageEngine interface{}

	
	
	
	ValidationAction *string

	
	
	
	
	ValidationLevel *string

	
	
	
	
	Validator interface{}

	
	
	
	
	
	
	ExpireAfterSeconds *int64

	
	
	
	
	
	
	TimeSeriesOptions *TimeSeriesOptions

	
	
	
	EncryptedFields interface{}

	
	
	
	ClusteredIndex interface{}
}


func CreateCollection() *CreateCollectionOptions {
	return &CreateCollectionOptions{}
}


func (c *CreateCollectionOptions) SetCapped(capped bool) *CreateCollectionOptions {
	c.Capped = &capped
	return c
}


func (c *CreateCollectionOptions) SetCollation(collation *Collation) *CreateCollectionOptions {
	c.Collation = collation
	return c
}


func (c *CreateCollectionOptions) SetChangeStreamPreAndPostImages(csppi interface{}) *CreateCollectionOptions {
	c.ChangeStreamPreAndPostImages = &csppi
	return c
}


func (c *CreateCollectionOptions) SetDefaultIndexOptions(opts *DefaultIndexOptions) *CreateCollectionOptions {
	c.DefaultIndexOptions = opts
	return c
}


func (c *CreateCollectionOptions) SetMaxDocuments(max int64) *CreateCollectionOptions {
	c.MaxDocuments = &max
	return c
}


func (c *CreateCollectionOptions) SetSizeInBytes(size int64) *CreateCollectionOptions {
	c.SizeInBytes = &size
	return c
}


func (c *CreateCollectionOptions) SetStorageEngine(storageEngine interface{}) *CreateCollectionOptions {
	c.StorageEngine = &storageEngine
	return c
}


func (c *CreateCollectionOptions) SetValidationAction(action string) *CreateCollectionOptions {
	c.ValidationAction = &action
	return c
}


func (c *CreateCollectionOptions) SetValidationLevel(level string) *CreateCollectionOptions {
	c.ValidationLevel = &level
	return c
}


func (c *CreateCollectionOptions) SetValidator(validator interface{}) *CreateCollectionOptions {
	c.Validator = validator
	return c
}


func (c *CreateCollectionOptions) SetExpireAfterSeconds(eas int64) *CreateCollectionOptions {
	c.ExpireAfterSeconds = &eas
	return c
}


func (c *CreateCollectionOptions) SetTimeSeriesOptions(timeSeriesOpts *TimeSeriesOptions) *CreateCollectionOptions {
	c.TimeSeriesOptions = timeSeriesOpts
	return c
}


func (c *CreateCollectionOptions) SetEncryptedFields(encryptedFields interface{}) *CreateCollectionOptions {
	c.EncryptedFields = encryptedFields
	return c
}


func (c *CreateCollectionOptions) SetClusteredIndex(clusteredIndex interface{}) *CreateCollectionOptions {
	c.ClusteredIndex = clusteredIndex
	return c
}






func MergeCreateCollectionOptions(opts ...*CreateCollectionOptions) *CreateCollectionOptions {
	cc := CreateCollection()

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.Capped != nil {
			cc.Capped = opt.Capped
		}
		if opt.Collation != nil {
			cc.Collation = opt.Collation
		}
		if opt.ChangeStreamPreAndPostImages != nil {
			cc.ChangeStreamPreAndPostImages = opt.ChangeStreamPreAndPostImages
		}
		if opt.DefaultIndexOptions != nil {
			cc.DefaultIndexOptions = opt.DefaultIndexOptions
		}
		if opt.MaxDocuments != nil {
			cc.MaxDocuments = opt.MaxDocuments
		}
		if opt.SizeInBytes != nil {
			cc.SizeInBytes = opt.SizeInBytes
		}
		if opt.StorageEngine != nil {
			cc.StorageEngine = opt.StorageEngine
		}
		if opt.ValidationAction != nil {
			cc.ValidationAction = opt.ValidationAction
		}
		if opt.ValidationLevel != nil {
			cc.ValidationLevel = opt.ValidationLevel
		}
		if opt.Validator != nil {
			cc.Validator = opt.Validator
		}
		if opt.ExpireAfterSeconds != nil {
			cc.ExpireAfterSeconds = opt.ExpireAfterSeconds
		}
		if opt.TimeSeriesOptions != nil {
			cc.TimeSeriesOptions = opt.TimeSeriesOptions
		}
		if opt.EncryptedFields != nil {
			cc.EncryptedFields = opt.EncryptedFields
		}
		if opt.ClusteredIndex != nil {
			cc.ClusteredIndex = opt.ClusteredIndex
		}
	}

	return cc
}


type CreateViewOptions struct {
	
	
	Collation *Collation
}


func CreateView() *CreateViewOptions {
	return &CreateViewOptions{}
}


func (c *CreateViewOptions) SetCollation(collation *Collation) *CreateViewOptions {
	c.Collation = collation
	return c
}






func MergeCreateViewOptions(opts ...*CreateViewOptions) *CreateViewOptions {
	cv := CreateView()

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.Collation != nil {
			cv.Collation = opt.Collation
		}
	}

	return cv
}
