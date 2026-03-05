





package options

import "time"


type CountOptions struct {
	
	
	
	Collation *Collation

	
	

	
	
	Comment *string

	
	
	
	Hint interface{}

	
	
	Limit *int64

	
	
	
	
	
	
	MaxTime *time.Duration

	
	Skip *int64
}


func Count() *CountOptions {
	return &CountOptions{}
}


func (co *CountOptions) SetCollation(c *Collation) *CountOptions {
	co.Collation = c
	return co
}


func (co *CountOptions) SetComment(c string) *CountOptions {
	co.Comment = &c
	return co
}


func (co *CountOptions) SetHint(h interface{}) *CountOptions {
	co.Hint = h
	return co
}


func (co *CountOptions) SetLimit(i int64) *CountOptions {
	co.Limit = &i
	return co
}






func (co *CountOptions) SetMaxTime(d time.Duration) *CountOptions {
	co.MaxTime = &d
	return co
}


func (co *CountOptions) SetSkip(i int64) *CountOptions {
	co.Skip = &i
	return co
}





func MergeCountOptions(opts ...*CountOptions) *CountOptions {
	countOpts := Count()
	for _, co := range opts {
		if co == nil {
			continue
		}
		if co.Collation != nil {
			countOpts.Collation = co.Collation
		}
		if co.Comment != nil {
			countOpts.Comment = co.Comment
		}
		if co.Hint != nil {
			countOpts.Hint = co.Hint
		}
		if co.Limit != nil {
			countOpts.Limit = co.Limit
		}
		if co.MaxTime != nil {
			countOpts.MaxTime = co.MaxTime
		}
		if co.Skip != nil {
			countOpts.Skip = co.Skip
		}
	}

	return countOpts
}
