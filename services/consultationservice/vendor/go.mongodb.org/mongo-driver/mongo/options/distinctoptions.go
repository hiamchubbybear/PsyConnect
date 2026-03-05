





package options

import "time"


type DistinctOptions struct {
	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	MaxTime *time.Duration
}


func Distinct() *DistinctOptions {
	return &DistinctOptions{}
}


func (do *DistinctOptions) SetCollation(c *Collation) *DistinctOptions {
	do.Collation = c
	return do
}


func (do *DistinctOptions) SetComment(comment interface{}) *DistinctOptions {
	do.Comment = comment
	return do
}






func (do *DistinctOptions) SetMaxTime(d time.Duration) *DistinctOptions {
	do.MaxTime = &d
	return do
}






func MergeDistinctOptions(opts ...*DistinctOptions) *DistinctOptions {
	distinctOpts := Distinct()
	for _, do := range opts {
		if do == nil {
			continue
		}
		if do.Collation != nil {
			distinctOpts.Collation = do.Collation
		}
		if do.Comment != nil {
			distinctOpts.Comment = do.Comment
		}
		if do.MaxTime != nil {
			distinctOpts.MaxTime = do.MaxTime
		}
	}

	return distinctOpts
}
