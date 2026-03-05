





package options


type InsertOneOptions struct {
	
	
	
	
	BypassDocumentValidation *bool

	
	
	Comment interface{}

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func InsertOne() *InsertOneOptions {
	return &InsertOneOptions{}
}


func (ioo *InsertOneOptions) SetBypassDocumentValidation(b bool) *InsertOneOptions {
	ioo.BypassDocumentValidation = &b
	return ioo
}


func (ioo *InsertOneOptions) SetComment(comment interface{}) *InsertOneOptions {
	ioo.Comment = comment
	return ioo
}






func MergeInsertOneOptions(opts ...*InsertOneOptions) *InsertOneOptions {
	ioOpts := InsertOne()
	for _, ioo := range opts {
		if ioo == nil {
			continue
		}
		if ioo.BypassDocumentValidation != nil {
			ioOpts.BypassDocumentValidation = ioo.BypassDocumentValidation
		}
		if ioo.Comment != nil {
			ioOpts.Comment = ioo.Comment
		}
		if ioo.BypassEmptyTsReplacement != nil {
			ioOpts.BypassEmptyTsReplacement = ioo.BypassEmptyTsReplacement
		}
	}

	return ioOpts
}


type InsertManyOptions struct {
	
	
	
	
	BypassDocumentValidation *bool

	
	
	Comment interface{}

	
	Ordered *bool

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func InsertMany() *InsertManyOptions {
	return &InsertManyOptions{
		Ordered: &DefaultOrdered,
	}
}


func (imo *InsertManyOptions) SetBypassDocumentValidation(b bool) *InsertManyOptions {
	imo.BypassDocumentValidation = &b
	return imo
}


func (imo *InsertManyOptions) SetComment(comment interface{}) *InsertManyOptions {
	imo.Comment = comment
	return imo
}


func (imo *InsertManyOptions) SetOrdered(b bool) *InsertManyOptions {
	imo.Ordered = &b
	return imo
}






func MergeInsertManyOptions(opts ...*InsertManyOptions) *InsertManyOptions {
	imOpts := InsertMany()
	for _, imo := range opts {
		if imo == nil {
			continue
		}
		if imo.BypassDocumentValidation != nil {
			imOpts.BypassDocumentValidation = imo.BypassDocumentValidation
		}
		if imo.Comment != nil {
			imOpts.Comment = imo.Comment
		}
		if imo.Ordered != nil {
			imOpts.Ordered = imo.Ordered
		}
		if imo.BypassEmptyTsReplacement != nil {
			imOpts.BypassEmptyTsReplacement = imo.BypassEmptyTsReplacement
		}
	}

	return imOpts
}
