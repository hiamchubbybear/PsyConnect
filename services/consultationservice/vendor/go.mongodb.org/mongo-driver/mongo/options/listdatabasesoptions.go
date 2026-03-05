





package options


type ListDatabasesOptions struct {
	
	
	NameOnly *bool

	
	
	
	AuthorizedDatabases *bool
}


func ListDatabases() *ListDatabasesOptions {
	return &ListDatabasesOptions{}
}


func (ld *ListDatabasesOptions) SetNameOnly(b bool) *ListDatabasesOptions {
	ld.NameOnly = &b
	return ld
}


func (ld *ListDatabasesOptions) SetAuthorizedDatabases(b bool) *ListDatabasesOptions {
	ld.AuthorizedDatabases = &b
	return ld
}






func MergeListDatabasesOptions(opts ...*ListDatabasesOptions) *ListDatabasesOptions {
	ld := ListDatabases()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.NameOnly != nil {
			ld.NameOnly = opt.NameOnly
		}
		if opt.AuthorizedDatabases != nil {
			ld.AuthorizedDatabases = opt.AuthorizedDatabases
		}
	}

	return ld
}
