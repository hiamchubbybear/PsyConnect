





package options


type ReplaceOptions struct {
	
	
	
	
	BypassDocumentValidation *bool

	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	Hint interface{}

	
	
	Upsert *bool

	
	
	
	
	Let interface{}

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func Replace() *ReplaceOptions {
	return &ReplaceOptions{}
}


func (ro *ReplaceOptions) SetBypassDocumentValidation(b bool) *ReplaceOptions {
	ro.BypassDocumentValidation = &b
	return ro
}


func (ro *ReplaceOptions) SetCollation(c *Collation) *ReplaceOptions {
	ro.Collation = c
	return ro
}


func (ro *ReplaceOptions) SetComment(comment interface{}) *ReplaceOptions {
	ro.Comment = comment
	return ro
}


func (ro *ReplaceOptions) SetHint(h interface{}) *ReplaceOptions {
	ro.Hint = h
	return ro
}


func (ro *ReplaceOptions) SetUpsert(b bool) *ReplaceOptions {
	ro.Upsert = &b
	return ro
}


func (ro *ReplaceOptions) SetLet(l interface{}) *ReplaceOptions {
	ro.Let = l
	return ro
}






func MergeReplaceOptions(opts ...*ReplaceOptions) *ReplaceOptions {
	rOpts := Replace()
	for _, ro := range opts {
		if ro == nil {
			continue
		}
		if ro.BypassDocumentValidation != nil {
			rOpts.BypassDocumentValidation = ro.BypassDocumentValidation
		}
		if ro.Collation != nil {
			rOpts.Collation = ro.Collation
		}
		if ro.Comment != nil {
			rOpts.Comment = ro.Comment
		}
		if ro.Hint != nil {
			rOpts.Hint = ro.Hint
		}
		if ro.Upsert != nil {
			rOpts.Upsert = ro.Upsert
		}
		if ro.Let != nil {
			rOpts.Let = ro.Let
		}
		if ro.BypassEmptyTsReplacement != nil {
			rOpts.BypassEmptyTsReplacement = ro.BypassEmptyTsReplacement
		}
	}

	return rOpts
}
