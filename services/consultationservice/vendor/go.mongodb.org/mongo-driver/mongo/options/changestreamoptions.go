





package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type ChangeStreamOptions struct {
	
	BatchSize *int32

	
	
	
	Collation *Collation

	
	
	Comment *string

	
	
	FullDocument *FullDocument

	
	
	FullDocumentBeforeChange *FullDocument

	
	MaxAwaitTime *time.Duration

	
	
	
	ResumeAfter interface{}

	
	
	
	ShowExpandedEvents *bool

	
	
	
	StartAtOperationTime *primitive.Timestamp

	
	
	
	
	
	StartAfter interface{}

	
	
	
	Custom bson.M

	
	
	
	CustomPipeline bson.M
}


func ChangeStream() *ChangeStreamOptions {
	cso := &ChangeStreamOptions{}
	return cso
}


func (cso *ChangeStreamOptions) SetBatchSize(i int32) *ChangeStreamOptions {
	cso.BatchSize = &i
	return cso
}


func (cso *ChangeStreamOptions) SetCollation(c Collation) *ChangeStreamOptions {
	cso.Collation = &c
	return cso
}


func (cso *ChangeStreamOptions) SetComment(comment string) *ChangeStreamOptions {
	cso.Comment = &comment
	return cso
}


func (cso *ChangeStreamOptions) SetFullDocument(fd FullDocument) *ChangeStreamOptions {
	cso.FullDocument = &fd
	return cso
}


func (cso *ChangeStreamOptions) SetFullDocumentBeforeChange(fdbc FullDocument) *ChangeStreamOptions {
	cso.FullDocumentBeforeChange = &fdbc
	return cso
}


func (cso *ChangeStreamOptions) SetMaxAwaitTime(d time.Duration) *ChangeStreamOptions {
	cso.MaxAwaitTime = &d
	return cso
}


func (cso *ChangeStreamOptions) SetResumeAfter(rt interface{}) *ChangeStreamOptions {
	cso.ResumeAfter = rt
	return cso
}


func (cso *ChangeStreamOptions) SetShowExpandedEvents(see bool) *ChangeStreamOptions {
	cso.ShowExpandedEvents = &see
	return cso
}


func (cso *ChangeStreamOptions) SetStartAtOperationTime(t *primitive.Timestamp) *ChangeStreamOptions {
	cso.StartAtOperationTime = t
	return cso
}


func (cso *ChangeStreamOptions) SetStartAfter(sa interface{}) *ChangeStreamOptions {
	cso.StartAfter = sa
	return cso
}





func (cso *ChangeStreamOptions) SetCustom(c bson.M) *ChangeStreamOptions {
	cso.Custom = c
	return cso
}




func (cso *ChangeStreamOptions) SetCustomPipeline(cp bson.M) *ChangeStreamOptions {
	cso.CustomPipeline = cp
	return cso
}






func MergeChangeStreamOptions(opts ...*ChangeStreamOptions) *ChangeStreamOptions {
	csOpts := ChangeStream()
	for _, cso := range opts {
		if cso == nil {
			continue
		}
		if cso.BatchSize != nil {
			csOpts.BatchSize = cso.BatchSize
		}
		if cso.Collation != nil {
			csOpts.Collation = cso.Collation
		}
		if cso.Comment != nil {
			csOpts.Comment = cso.Comment
		}
		if cso.FullDocument != nil {
			csOpts.FullDocument = cso.FullDocument
		}
		if cso.FullDocumentBeforeChange != nil {
			csOpts.FullDocumentBeforeChange = cso.FullDocumentBeforeChange
		}
		if cso.MaxAwaitTime != nil {
			csOpts.MaxAwaitTime = cso.MaxAwaitTime
		}
		if cso.ResumeAfter != nil {
			csOpts.ResumeAfter = cso.ResumeAfter
		}
		if cso.ShowExpandedEvents != nil {
			csOpts.ShowExpandedEvents = cso.ShowExpandedEvents
		}
		if cso.StartAtOperationTime != nil {
			csOpts.StartAtOperationTime = cso.StartAtOperationTime
		}
		if cso.StartAfter != nil {
			csOpts.StartAfter = cso.StartAfter
		}
		if cso.Custom != nil {
			csOpts.Custom = cso.Custom
		}
		if cso.CustomPipeline != nil {
			csOpts.CustomPipeline = cso.CustomPipeline
		}
	}

	return csOpts
}
