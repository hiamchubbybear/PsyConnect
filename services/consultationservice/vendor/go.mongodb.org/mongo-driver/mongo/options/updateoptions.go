





package options


type UpdateOptions struct {
	
	
	
	ArrayFilters *ArrayFilters

	
	
	
	
	BypassDocumentValidation *bool

	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	Hint interface{}

	
	
	Upsert *bool

	
	
	
	
	Let interface{}

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func Update() *UpdateOptions {
	return &UpdateOptions{}
}


func (uo *UpdateOptions) SetArrayFilters(af ArrayFilters) *UpdateOptions {
	uo.ArrayFilters = &af
	return uo
}


func (uo *UpdateOptions) SetBypassDocumentValidation(b bool) *UpdateOptions {
	uo.BypassDocumentValidation = &b
	return uo
}


func (uo *UpdateOptions) SetCollation(c *Collation) *UpdateOptions {
	uo.Collation = c
	return uo
}


func (uo *UpdateOptions) SetComment(comment interface{}) *UpdateOptions {
	uo.Comment = comment
	return uo
}


func (uo *UpdateOptions) SetHint(h interface{}) *UpdateOptions {
	uo.Hint = h
	return uo
}


func (uo *UpdateOptions) SetUpsert(b bool) *UpdateOptions {
	uo.Upsert = &b
	return uo
}


func (uo *UpdateOptions) SetLet(l interface{}) *UpdateOptions {
	uo.Let = l
	return uo
}





func MergeUpdateOptions(opts ...*UpdateOptions) *UpdateOptions {
	uOpts := Update()
	for _, uo := range opts {
		if uo == nil {
			continue
		}
		if uo.ArrayFilters != nil {
			uOpts.ArrayFilters = uo.ArrayFilters
		}
		if uo.BypassDocumentValidation != nil {
			uOpts.BypassDocumentValidation = uo.BypassDocumentValidation
		}
		if uo.Collation != nil {
			uOpts.Collation = uo.Collation
		}
		if uo.Comment != nil {
			uOpts.Comment = uo.Comment
		}
		if uo.Hint != nil {
			uOpts.Hint = uo.Hint
		}
		if uo.Upsert != nil {
			uOpts.Upsert = uo.Upsert
		}
		if uo.Let != nil {
			uOpts.Let = uo.Let
		}
		if uo.BypassEmptyTsReplacement != nil {
			uOpts.BypassEmptyTsReplacement = uo.BypassEmptyTsReplacement
		}
	}

	return uOpts
}
