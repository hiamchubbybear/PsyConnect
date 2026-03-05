





package options

import "time"


type EstimatedDocumentCountOptions struct {
	
	
	Comment interface{}

	
	
	
	
	
	
	MaxTime *time.Duration
}


func EstimatedDocumentCount() *EstimatedDocumentCountOptions {
	return &EstimatedDocumentCountOptions{}
}


func (eco *EstimatedDocumentCountOptions) SetComment(comment interface{}) *EstimatedDocumentCountOptions {
	eco.Comment = comment
	return eco
}






func (eco *EstimatedDocumentCountOptions) SetMaxTime(d time.Duration) *EstimatedDocumentCountOptions {
	eco.MaxTime = &d
	return eco
}






func MergeEstimatedDocumentCountOptions(opts ...*EstimatedDocumentCountOptions) *EstimatedDocumentCountOptions {
	e := EstimatedDocumentCount()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Comment != nil {
			e.Comment = opt.Comment
		}
		if opt.MaxTime != nil {
			e.MaxTime = opt.MaxTime
		}
	}

	return e
}
