





package options


type DeleteOptions struct {
	
	
	
	Collation *Collation

	
	
	Comment interface{}

	
	
	
	
	
	
	Hint interface{}

	
	
	
	
	Let interface{}
}


func Delete() *DeleteOptions {
	return &DeleteOptions{}
}


func (do *DeleteOptions) SetCollation(c *Collation) *DeleteOptions {
	do.Collation = c
	return do
}


func (do *DeleteOptions) SetComment(comment interface{}) *DeleteOptions {
	do.Comment = comment
	return do
}


func (do *DeleteOptions) SetHint(hint interface{}) *DeleteOptions {
	do.Hint = hint
	return do
}


func (do *DeleteOptions) SetLet(let interface{}) *DeleteOptions {
	do.Let = let
	return do
}





func MergeDeleteOptions(opts ...*DeleteOptions) *DeleteOptions {
	dOpts := Delete()
	for _, do := range opts {
		if do == nil {
			continue
		}
		if do.Collation != nil {
			dOpts.Collation = do.Collation
		}
		if do.Comment != nil {
			dOpts.Comment = do.Comment
		}
		if do.Hint != nil {
			dOpts.Hint = do.Hint
		}
		if do.Let != nil {
			dOpts.Let = do.Let
		}
	}

	return dOpts
}
