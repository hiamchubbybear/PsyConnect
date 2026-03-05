





package options


var DefaultOrdered = true


type BulkWriteOptions struct {
	
	
	
	
	BypassDocumentValidation *bool

	
	
	Comment interface{}

	
	Ordered *bool

	
	
	
	
	Let interface{}

	
	
	
	
	BypassEmptyTsReplacement *bool
}


func BulkWrite() *BulkWriteOptions {
	return &BulkWriteOptions{
		Ordered: &DefaultOrdered,
	}
}


func (b *BulkWriteOptions) SetComment(comment interface{}) *BulkWriteOptions {
	b.Comment = comment
	return b
}


func (b *BulkWriteOptions) SetOrdered(ordered bool) *BulkWriteOptions {
	b.Ordered = &ordered
	return b
}


func (b *BulkWriteOptions) SetBypassDocumentValidation(bypass bool) *BulkWriteOptions {
	b.BypassDocumentValidation = &bypass
	return b
}





func (b *BulkWriteOptions) SetLet(let interface{}) *BulkWriteOptions {
	b.Let = &let
	return b
}






func MergeBulkWriteOptions(opts ...*BulkWriteOptions) *BulkWriteOptions {
	b := BulkWrite()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Comment != nil {
			b.Comment = opt.Comment
		}
		if opt.Ordered != nil {
			b.Ordered = opt.Ordered
		}
		if opt.BypassDocumentValidation != nil {
			b.BypassDocumentValidation = opt.BypassDocumentValidation
		}
		if opt.Let != nil {
			b.Let = opt.Let
		}
		if opt.BypassEmptyTsReplacement != nil {
			b.BypassEmptyTsReplacement = opt.BypassEmptyTsReplacement
		}
	}

	return b
}
