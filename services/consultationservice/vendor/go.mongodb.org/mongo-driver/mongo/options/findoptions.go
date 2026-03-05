





package options

import (
	"time"
)


type FindOptions struct {
	
	
	
	
	AllowDiskUse *bool

	
	
	AllowPartialResults *bool

	
	BatchSize *int32

	
	
	
	Collation *Collation

	
	
	Comment *string

	
	
	CursorType *CursorType

	
	
	
	Hint interface{}

	
	
	
	Limit *int64

	
	
	Max interface{}

	
	
	
	MaxAwaitTime *time.Duration

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	Min interface{}

	
	
	NoCursorTimeout *bool

	
	
	
	
	OplogReplay *bool

	
	
	Projection interface{}

	
	
	ReturnKey *bool

	
	
	ShowRecordID *bool

	
	Skip *int64

	
	
	
	
	Snapshot *bool

	
	
	Sort interface{}

	
	
	
	
	Let interface{}
}


func Find() *FindOptions {
	return &FindOptions{}
}


func (f *FindOptions) SetAllowDiskUse(b bool) *FindOptions {
	f.AllowDiskUse = &b
	return f
}


func (f *FindOptions) SetAllowPartialResults(b bool) *FindOptions {
	f.AllowPartialResults = &b
	return f
}


func (f *FindOptions) SetBatchSize(i int32) *FindOptions {
	f.BatchSize = &i
	return f
}


func (f *FindOptions) SetCollation(collation *Collation) *FindOptions {
	f.Collation = collation
	return f
}


func (f *FindOptions) SetComment(comment string) *FindOptions {
	f.Comment = &comment
	return f
}


func (f *FindOptions) SetCursorType(ct CursorType) *FindOptions {
	f.CursorType = &ct
	return f
}


func (f *FindOptions) SetHint(hint interface{}) *FindOptions {
	f.Hint = hint
	return f
}


func (f *FindOptions) SetLet(let interface{}) *FindOptions {
	f.Let = let
	return f
}


func (f *FindOptions) SetLimit(i int64) *FindOptions {
	f.Limit = &i
	return f
}


func (f *FindOptions) SetMax(max interface{}) *FindOptions {
	f.Max = max
	return f
}


func (f *FindOptions) SetMaxAwaitTime(d time.Duration) *FindOptions {
	f.MaxAwaitTime = &d
	return f
}






func (f *FindOptions) SetMaxTime(d time.Duration) *FindOptions {
	f.MaxTime = &d
	return f
}


func (f *FindOptions) SetMin(min interface{}) *FindOptions {
	f.Min = min
	return f
}


func (f *FindOptions) SetNoCursorTimeout(b bool) *FindOptions {
	f.NoCursorTimeout = &b
	return f
}




func (f *FindOptions) SetOplogReplay(b bool) *FindOptions {
	f.OplogReplay = &b
	return f
}


func (f *FindOptions) SetProjection(projection interface{}) *FindOptions {
	f.Projection = projection
	return f
}


func (f *FindOptions) SetReturnKey(b bool) *FindOptions {
	f.ReturnKey = &b
	return f
}


func (f *FindOptions) SetShowRecordID(b bool) *FindOptions {
	f.ShowRecordID = &b
	return f
}


func (f *FindOptions) SetSkip(i int64) *FindOptions {
	f.Skip = &i
	return f
}




func (f *FindOptions) SetSnapshot(b bool) *FindOptions {
	f.Snapshot = &b
	return f
}


func (f *FindOptions) SetSort(sort interface{}) *FindOptions {
	f.Sort = sort
	return f
}





func MergeFindOptions(opts ...*FindOptions) *FindOptions {
	fo := Find()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowDiskUse != nil {
			fo.AllowDiskUse = opt.AllowDiskUse
		}
		if opt.AllowPartialResults != nil {
			fo.AllowPartialResults = opt.AllowPartialResults
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.CursorType != nil {
			fo.CursorType = opt.CursorType
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Let != nil {
			fo.Let = opt.Let
		}
		if opt.Limit != nil {
			fo.Limit = opt.Limit
		}
		if opt.Max != nil {
			fo.Max = opt.Max
		}
		if opt.MaxAwaitTime != nil {
			fo.MaxAwaitTime = opt.MaxAwaitTime
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Min != nil {
			fo.Min = opt.Min
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.OplogReplay != nil {
			fo.OplogReplay = opt.OplogReplay
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnKey != nil {
			fo.ReturnKey = opt.ReturnKey
		}
		if opt.ShowRecordID != nil {
			fo.ShowRecordID = opt.ShowRecordID
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Snapshot != nil {
			fo.Snapshot = opt.Snapshot
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}


type FindOneOptions struct {
	
	
	AllowPartialResults *bool

	
	
	
	BatchSize *int32

	
	
	
	Collation *Collation

	
	
	Comment *string

	
	
	
	
	CursorType *CursorType

	
	
	
	Hint interface{}

	
	
	Max interface{}

	
	
	
	
	
	MaxAwaitTime *time.Duration

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	Min interface{}

	
	
	
	
	NoCursorTimeout *bool

	
	
	
	
	OplogReplay *bool

	
	
	Projection interface{}

	
	
	ReturnKey *bool

	
	
	ShowRecordID *bool

	
	Skip *int64

	
	
	
	
	Snapshot *bool

	
	
	Sort interface{}
}


func FindOne() *FindOneOptions {
	return &FindOneOptions{}
}


func (f *FindOneOptions) SetAllowPartialResults(b bool) *FindOneOptions {
	f.AllowPartialResults = &b
	return f
}




func (f *FindOneOptions) SetBatchSize(i int32) *FindOneOptions {
	f.BatchSize = &i
	return f
}


func (f *FindOneOptions) SetCollation(collation *Collation) *FindOneOptions {
	f.Collation = collation
	return f
}


func (f *FindOneOptions) SetComment(comment string) *FindOneOptions {
	f.Comment = &comment
	return f
}




func (f *FindOneOptions) SetCursorType(ct CursorType) *FindOneOptions {
	f.CursorType = &ct
	return f
}


func (f *FindOneOptions) SetHint(hint interface{}) *FindOneOptions {
	f.Hint = hint
	return f
}


func (f *FindOneOptions) SetMax(max interface{}) *FindOneOptions {
	f.Max = max
	return f
}




func (f *FindOneOptions) SetMaxAwaitTime(d time.Duration) *FindOneOptions {
	f.MaxAwaitTime = &d
	return f
}






func (f *FindOneOptions) SetMaxTime(d time.Duration) *FindOneOptions {
	f.MaxTime = &d
	return f
}


func (f *FindOneOptions) SetMin(min interface{}) *FindOneOptions {
	f.Min = min
	return f
}




func (f *FindOneOptions) SetNoCursorTimeout(b bool) *FindOneOptions {
	f.NoCursorTimeout = &b
	return f
}





func (f *FindOneOptions) SetOplogReplay(b bool) *FindOneOptions {
	f.OplogReplay = &b
	return f
}


func (f *FindOneOptions) SetProjection(projection interface{}) *FindOneOptions {
	f.Projection = projection
	return f
}


func (f *FindOneOptions) SetReturnKey(b bool) *FindOneOptions {
	f.ReturnKey = &b
	return f
}


func (f *FindOneOptions) SetShowRecordID(b bool) *FindOneOptions {
	f.ShowRecordID = &b
	return f
}


func (f *FindOneOptions) SetSkip(i int64) *FindOneOptions {
	f.Skip = &i
	return f
}




func (f *FindOneOptions) SetSnapshot(b bool) *FindOneOptions {
	f.Snapshot = &b
	return f
}


func (f *FindOneOptions) SetSort(sort interface{}) *FindOneOptions {
	f.Sort = sort
	return f
}






func MergeFindOneOptions(opts ...*FindOneOptions) *FindOneOptions {
	fo := FindOne()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowPartialResults != nil {
			fo.AllowPartialResults = opt.AllowPartialResults
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.CursorType != nil {
			fo.CursorType = opt.CursorType
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Max != nil {
			fo.Max = opt.Max
		}
		if opt.MaxAwaitTime != nil {
			fo.MaxAwaitTime = opt.MaxAwaitTime
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Min != nil {
			fo.Min = opt.Min
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.OplogReplay != nil {
			fo.OplogReplay = opt.OplogReplay
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnKey != nil {
			fo.ReturnKey = opt.ReturnKey
		}
		if opt.ShowRecordID != nil {
			fo.ShowRecordID = opt.ShowRecordID
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Snapshot != nil {
			fo.Snapshot = opt.Snapshot
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}


type FindOneAndReplaceOptions struct {
	
	
	
	
	BypassDocumentValidation *bool

	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	Projection interface{}

	
	
	ReturnDocument *ReturnDocument

	
	
	
	Sort interface{}

	
	
	Upsert *bool

	
	
	
	
	
	
	Hint interface{}

	
	
	
	
	Let interface{}

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func FindOneAndReplace() *FindOneAndReplaceOptions {
	return &FindOneAndReplaceOptions{}
}


func (f *FindOneAndReplaceOptions) SetBypassDocumentValidation(b bool) *FindOneAndReplaceOptions {
	f.BypassDocumentValidation = &b
	return f
}


func (f *FindOneAndReplaceOptions) SetCollation(collation *Collation) *FindOneAndReplaceOptions {
	f.Collation = collation
	return f
}


func (f *FindOneAndReplaceOptions) SetComment(comment interface{}) *FindOneAndReplaceOptions {
	f.Comment = comment
	return f
}






func (f *FindOneAndReplaceOptions) SetMaxTime(d time.Duration) *FindOneAndReplaceOptions {
	f.MaxTime = &d
	return f
}


func (f *FindOneAndReplaceOptions) SetProjection(projection interface{}) *FindOneAndReplaceOptions {
	f.Projection = projection
	return f
}


func (f *FindOneAndReplaceOptions) SetReturnDocument(rd ReturnDocument) *FindOneAndReplaceOptions {
	f.ReturnDocument = &rd
	return f
}


func (f *FindOneAndReplaceOptions) SetSort(sort interface{}) *FindOneAndReplaceOptions {
	f.Sort = sort
	return f
}


func (f *FindOneAndReplaceOptions) SetUpsert(b bool) *FindOneAndReplaceOptions {
	f.Upsert = &b
	return f
}


func (f *FindOneAndReplaceOptions) SetHint(hint interface{}) *FindOneAndReplaceOptions {
	f.Hint = hint
	return f
}


func (f *FindOneAndReplaceOptions) SetLet(let interface{}) *FindOneAndReplaceOptions {
	f.Let = let
	return f
}






func MergeFindOneAndReplaceOptions(opts ...*FindOneAndReplaceOptions) *FindOneAndReplaceOptions {
	fo := FindOneAndReplace()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.BypassDocumentValidation != nil {
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnDocument != nil {
			fo.ReturnDocument = opt.ReturnDocument
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
		if opt.Upsert != nil {
			fo.Upsert = opt.Upsert
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Let != nil {
			fo.Let = opt.Let
		}
		if opt.BypassEmptyTsReplacement != nil {
			fo.BypassEmptyTsReplacement = opt.BypassEmptyTsReplacement
		}
	}

	return fo
}


type FindOneAndUpdateOptions struct {
	
	
	
	ArrayFilters *ArrayFilters

	
	
	
	
	BypassDocumentValidation *bool

	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	Projection interface{}

	
	
	ReturnDocument *ReturnDocument

	
	
	
	Sort interface{}

	
	
	Upsert *bool

	
	
	
	
	
	
	Hint interface{}

	
	
	
	
	Let interface{}

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func FindOneAndUpdate() *FindOneAndUpdateOptions {
	return &FindOneAndUpdateOptions{}
}


func (f *FindOneAndUpdateOptions) SetBypassDocumentValidation(b bool) *FindOneAndUpdateOptions {
	f.BypassDocumentValidation = &b
	return f
}


func (f *FindOneAndUpdateOptions) SetArrayFilters(filters ArrayFilters) *FindOneAndUpdateOptions {
	f.ArrayFilters = &filters
	return f
}


func (f *FindOneAndUpdateOptions) SetCollation(collation *Collation) *FindOneAndUpdateOptions {
	f.Collation = collation
	return f
}


func (f *FindOneAndUpdateOptions) SetComment(comment interface{}) *FindOneAndUpdateOptions {
	f.Comment = comment
	return f
}






func (f *FindOneAndUpdateOptions) SetMaxTime(d time.Duration) *FindOneAndUpdateOptions {
	f.MaxTime = &d
	return f
}


func (f *FindOneAndUpdateOptions) SetProjection(projection interface{}) *FindOneAndUpdateOptions {
	f.Projection = projection
	return f
}


func (f *FindOneAndUpdateOptions) SetReturnDocument(rd ReturnDocument) *FindOneAndUpdateOptions {
	f.ReturnDocument = &rd
	return f
}


func (f *FindOneAndUpdateOptions) SetSort(sort interface{}) *FindOneAndUpdateOptions {
	f.Sort = sort
	return f
}


func (f *FindOneAndUpdateOptions) SetUpsert(b bool) *FindOneAndUpdateOptions {
	f.Upsert = &b
	return f
}


func (f *FindOneAndUpdateOptions) SetHint(hint interface{}) *FindOneAndUpdateOptions {
	f.Hint = hint
	return f
}


func (f *FindOneAndUpdateOptions) SetLet(let interface{}) *FindOneAndUpdateOptions {
	f.Let = let
	return f
}






func MergeFindOneAndUpdateOptions(opts ...*FindOneAndUpdateOptions) *FindOneAndUpdateOptions {
	fo := FindOneAndUpdate()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ArrayFilters != nil {
			fo.ArrayFilters = opt.ArrayFilters
		}
		if opt.BypassDocumentValidation != nil {
			fo.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.ReturnDocument != nil {
			fo.ReturnDocument = opt.ReturnDocument
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
		if opt.Upsert != nil {
			fo.Upsert = opt.Upsert
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Let != nil {
			fo.Let = opt.Let
		}
		if opt.BypassEmptyTsReplacement != nil {
			fo.BypassEmptyTsReplacement = opt.BypassEmptyTsReplacement
		}
	}

	return fo
}


type FindOneAndDeleteOptions struct {
	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	MaxTime *time.Duration

	
	
	Projection interface{}

	
	
	
	Sort interface{}

	
	
	
	
	
	
	Hint interface{}

	
	
	
	
	Let interface{}
}


func FindOneAndDelete() *FindOneAndDeleteOptions {
	return &FindOneAndDeleteOptions{}
}


func (f *FindOneAndDeleteOptions) SetCollation(collation *Collation) *FindOneAndDeleteOptions {
	f.Collation = collation
	return f
}


func (f *FindOneAndDeleteOptions) SetComment(comment interface{}) *FindOneAndDeleteOptions {
	f.Comment = comment
	return f
}






func (f *FindOneAndDeleteOptions) SetMaxTime(d time.Duration) *FindOneAndDeleteOptions {
	f.MaxTime = &d
	return f
}


func (f *FindOneAndDeleteOptions) SetProjection(projection interface{}) *FindOneAndDeleteOptions {
	f.Projection = projection
	return f
}


func (f *FindOneAndDeleteOptions) SetSort(sort interface{}) *FindOneAndDeleteOptions {
	f.Sort = sort
	return f
}


func (f *FindOneAndDeleteOptions) SetHint(hint interface{}) *FindOneAndDeleteOptions {
	f.Hint = hint
	return f
}


func (f *FindOneAndDeleteOptions) SetLet(let interface{}) *FindOneAndDeleteOptions {
	f.Let = let
	return f
}






func MergeFindOneAndDeleteOptions(opts ...*FindOneAndDeleteOptions) *FindOneAndDeleteOptions {
	fo := FindOneAndDelete()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Collation != nil {
			fo.Collation = opt.Collation
		}
		if opt.Comment != nil {
			fo.Comment = opt.Comment
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.Projection != nil {
			fo.Projection = opt.Projection
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
		if opt.Hint != nil {
			fo.Hint = opt.Hint
		}
		if opt.Let != nil {
			fo.Let = opt.Let
		}
	}

	return fo
}
