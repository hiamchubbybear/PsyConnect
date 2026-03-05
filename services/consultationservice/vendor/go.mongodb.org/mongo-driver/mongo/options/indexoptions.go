





package options

import (
	"time"
)



type CreateIndexesOptions struct {
	
	
	
	
	
	
	
	
	
	
	
	
	CommitQuorum interface{}

	
	
	
	
	
	
	MaxTime *time.Duration
}


func CreateIndexes() *CreateIndexesOptions {
	return &CreateIndexesOptions{}
}






func (c *CreateIndexesOptions) SetMaxTime(d time.Duration) *CreateIndexesOptions {
	c.MaxTime = &d
	return c
}


func (c *CreateIndexesOptions) SetCommitQuorumInt(quorum int32) *CreateIndexesOptions {
	c.CommitQuorum = quorum
	return c
}


func (c *CreateIndexesOptions) SetCommitQuorumString(quorum string) *CreateIndexesOptions {
	c.CommitQuorum = quorum
	return c
}


func (c *CreateIndexesOptions) SetCommitQuorumMajority() *CreateIndexesOptions {
	c.CommitQuorum = "majority"
	return c
}


func (c *CreateIndexesOptions) SetCommitQuorumVotingMembers() *CreateIndexesOptions {
	c.CommitQuorum = "votingMembers"
	return c
}






func MergeCreateIndexesOptions(opts ...*CreateIndexesOptions) *CreateIndexesOptions {
	c := CreateIndexes()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.MaxTime != nil {
			c.MaxTime = opt.MaxTime
		}
		if opt.CommitQuorum != nil {
			c.CommitQuorum = opt.CommitQuorum
		}
	}

	return c
}



type DropIndexesOptions struct {
	
	
	
	
	
	
	MaxTime *time.Duration
}


func DropIndexes() *DropIndexesOptions {
	return &DropIndexesOptions{}
}






func (d *DropIndexesOptions) SetMaxTime(duration time.Duration) *DropIndexesOptions {
	d.MaxTime = &duration
	return d
}






func MergeDropIndexesOptions(opts ...*DropIndexesOptions) *DropIndexesOptions {
	c := DropIndexes()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.MaxTime != nil {
			c.MaxTime = opt.MaxTime
		}
	}

	return c
}


type ListIndexesOptions struct {
	
	BatchSize *int32

	
	
	
	
	
	
	MaxTime *time.Duration
}


func ListIndexes() *ListIndexesOptions {
	return &ListIndexesOptions{}
}


func (l *ListIndexesOptions) SetBatchSize(i int32) *ListIndexesOptions {
	l.BatchSize = &i
	return l
}






func (l *ListIndexesOptions) SetMaxTime(d time.Duration) *ListIndexesOptions {
	l.MaxTime = &d
	return l
}






func MergeListIndexesOptions(opts ...*ListIndexesOptions) *ListIndexesOptions {
	c := ListIndexes()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.BatchSize != nil {
			c.BatchSize = opt.BatchSize
		}
		if opt.MaxTime != nil {
			c.MaxTime = opt.MaxTime
		}
	}

	return c
}



type IndexOptions struct {
	
	
	
	
	Background *bool

	
	
	ExpireAfterSeconds *int32

	
	
	Name *string

	
	
	Sparse *bool

	
	
	
	
	StorageEngine interface{}

	
	
	Unique *bool

	
	Version *int32

	
	
	DefaultLanguage *string

	
	
	
	LanguageOverride *string

	
	
	TextVersion *int32

	
	
	
	
	Weights interface{}

	
	
	SphereVersion *int32

	
	
	Bits *int32

	
	
	Max *float64

	
	
	Min *float64

	
	
	
	BucketSize *int32

	
	
	PartialFilterExpression interface{}

	
	
	Collation *Collation

	
	WildcardProjection interface{}

	
	
	Hidden *bool
}


func Index() *IndexOptions {
	return &IndexOptions{}
}




func (i *IndexOptions) SetBackground(background bool) *IndexOptions {
	i.Background = &background
	return i
}


func (i *IndexOptions) SetExpireAfterSeconds(seconds int32) *IndexOptions {
	i.ExpireAfterSeconds = &seconds
	return i
}


func (i *IndexOptions) SetName(name string) *IndexOptions {
	i.Name = &name
	return i
}


func (i *IndexOptions) SetSparse(sparse bool) *IndexOptions {
	i.Sparse = &sparse
	return i
}


func (i *IndexOptions) SetStorageEngine(engine interface{}) *IndexOptions {
	i.StorageEngine = engine
	return i
}


func (i *IndexOptions) SetUnique(unique bool) *IndexOptions {
	i.Unique = &unique
	return i
}


func (i *IndexOptions) SetVersion(version int32) *IndexOptions {
	i.Version = &version
	return i
}


func (i *IndexOptions) SetDefaultLanguage(language string) *IndexOptions {
	i.DefaultLanguage = &language
	return i
}


func (i *IndexOptions) SetLanguageOverride(override string) *IndexOptions {
	i.LanguageOverride = &override
	return i
}


func (i *IndexOptions) SetTextVersion(version int32) *IndexOptions {
	i.TextVersion = &version
	return i
}


func (i *IndexOptions) SetWeights(weights interface{}) *IndexOptions {
	i.Weights = weights
	return i
}


func (i *IndexOptions) SetSphereVersion(version int32) *IndexOptions {
	i.SphereVersion = &version
	return i
}


func (i *IndexOptions) SetBits(bits int32) *IndexOptions {
	i.Bits = &bits
	return i
}


func (i *IndexOptions) SetMax(max float64) *IndexOptions {
	i.Max = &max
	return i
}


func (i *IndexOptions) SetMin(min float64) *IndexOptions {
	i.Min = &min
	return i
}


func (i *IndexOptions) SetBucketSize(bucketSize int32) *IndexOptions {
	i.BucketSize = &bucketSize
	return i
}


func (i *IndexOptions) SetPartialFilterExpression(expression interface{}) *IndexOptions {
	i.PartialFilterExpression = expression
	return i
}


func (i *IndexOptions) SetCollation(collation *Collation) *IndexOptions {
	i.Collation = collation
	return i
}


func (i *IndexOptions) SetWildcardProjection(wildcardProjection interface{}) *IndexOptions {
	i.WildcardProjection = wildcardProjection
	return i
}


func (i *IndexOptions) SetHidden(hidden bool) *IndexOptions {
	i.Hidden = &hidden
	return i
}





func MergeIndexOptions(opts ...*IndexOptions) *IndexOptions {
	i := Index()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Background != nil {
			i.Background = opt.Background
		}
		if opt.ExpireAfterSeconds != nil {
			i.ExpireAfterSeconds = opt.ExpireAfterSeconds
		}
		if opt.Name != nil {
			i.Name = opt.Name
		}
		if opt.Sparse != nil {
			i.Sparse = opt.Sparse
		}
		if opt.StorageEngine != nil {
			i.StorageEngine = opt.StorageEngine
		}
		if opt.Unique != nil {
			i.Unique = opt.Unique
		}
		if opt.Version != nil {
			i.Version = opt.Version
		}
		if opt.DefaultLanguage != nil {
			i.DefaultLanguage = opt.DefaultLanguage
		}
		if opt.LanguageOverride != nil {
			i.LanguageOverride = opt.LanguageOverride
		}
		if opt.TextVersion != nil {
			i.TextVersion = opt.TextVersion
		}
		if opt.Weights != nil {
			i.Weights = opt.Weights
		}
		if opt.SphereVersion != nil {
			i.SphereVersion = opt.SphereVersion
		}
		if opt.Bits != nil {
			i.Bits = opt.Bits
		}
		if opt.Max != nil {
			i.Max = opt.Max
		}
		if opt.Min != nil {
			i.Min = opt.Min
		}
		if opt.BucketSize != nil {
			i.BucketSize = opt.BucketSize
		}
		if opt.PartialFilterExpression != nil {
			i.PartialFilterExpression = opt.PartialFilterExpression
		}
		if opt.Collation != nil {
			i.Collation = opt.Collation
		}
		if opt.WildcardProjection != nil {
			i.WildcardProjection = opt.WildcardProjection
		}
		if opt.Hidden != nil {
			i.Hidden = opt.Hidden
		}
	}

	return i
}
