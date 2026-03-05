





package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)


type AggregateOptions struct {
	
	
	AllowDiskUse *bool

	
	BatchSize *int32

	
	
	
	
	BypassDocumentValidation *bool

	
	
	
	Collation *Collation

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	MaxAwaitTime *time.Duration

	
	
	Comment *string

	
	
	
	Hint interface{}

	
	
	
	
	Let interface{}

	
	
	
	Custom bson.M
}


func Aggregate() *AggregateOptions {
	return &AggregateOptions{}
}


func (ao *AggregateOptions) SetAllowDiskUse(b bool) *AggregateOptions {
	ao.AllowDiskUse = &b
	return ao
}


func (ao *AggregateOptions) SetBatchSize(i int32) *AggregateOptions {
	ao.BatchSize = &i
	return ao
}


func (ao *AggregateOptions) SetBypassDocumentValidation(b bool) *AggregateOptions {
	ao.BypassDocumentValidation = &b
	return ao
}


func (ao *AggregateOptions) SetCollation(c *Collation) *AggregateOptions {
	ao.Collation = c
	return ao
}






func (ao *AggregateOptions) SetMaxTime(d time.Duration) *AggregateOptions {
	ao.MaxTime = &d
	return ao
}


func (ao *AggregateOptions) SetMaxAwaitTime(d time.Duration) *AggregateOptions {
	ao.MaxAwaitTime = &d
	return ao
}


func (ao *AggregateOptions) SetComment(s string) *AggregateOptions {
	ao.Comment = &s
	return ao
}


func (ao *AggregateOptions) SetHint(h interface{}) *AggregateOptions {
	ao.Hint = h
	return ao
}


func (ao *AggregateOptions) SetLet(let interface{}) *AggregateOptions {
	ao.Let = let
	return ao
}





func (ao *AggregateOptions) SetCustom(c bson.M) *AggregateOptions {
	ao.Custom = c
	return ao
}






func MergeAggregateOptions(opts ...*AggregateOptions) *AggregateOptions {
	aggOpts := Aggregate()
	for _, ao := range opts {
		if ao == nil {
			continue
		}
		if ao.AllowDiskUse != nil {
			aggOpts.AllowDiskUse = ao.AllowDiskUse
		}
		if ao.BatchSize != nil {
			aggOpts.BatchSize = ao.BatchSize
		}
		if ao.BypassDocumentValidation != nil {
			aggOpts.BypassDocumentValidation = ao.BypassDocumentValidation
		}
		if ao.Collation != nil {
			aggOpts.Collation = ao.Collation
		}
		if ao.MaxTime != nil {
			aggOpts.MaxTime = ao.MaxTime
		}
		if ao.MaxAwaitTime != nil {
			aggOpts.MaxAwaitTime = ao.MaxAwaitTime
		}
		if ao.Comment != nil {
			aggOpts.Comment = ao.Comment
		}
		if ao.Hint != nil {
			aggOpts.Hint = ao.Hint
		}
		if ao.Let != nil {
			aggOpts.Let = ao.Let
		}
		if ao.Custom != nil {
			aggOpts.Custom = ao.Custom
		}
	}

	return aggOpts
}
