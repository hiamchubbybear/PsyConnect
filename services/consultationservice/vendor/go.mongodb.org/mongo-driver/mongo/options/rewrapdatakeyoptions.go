





package options



type RewrapManyDataKeyOptions struct {
	
	Provider *string

	
	MasterKey interface{}
}


func RewrapManyDataKey() *RewrapManyDataKeyOptions {
	return new(RewrapManyDataKeyOptions)
}


func (rmdko *RewrapManyDataKeyOptions) SetProvider(provider string) *RewrapManyDataKeyOptions {
	rmdko.Provider = &provider
	return rmdko
}


func (rmdko *RewrapManyDataKeyOptions) SetMasterKey(masterKey interface{}) *RewrapManyDataKeyOptions {
	rmdko.MasterKey = masterKey
	return rmdko
}






func MergeRewrapManyDataKeyOptions(opts ...*RewrapManyDataKeyOptions) *RewrapManyDataKeyOptions {
	rmdkOpts := RewrapManyDataKey()
	for _, rmdko := range opts {
		if rmdko == nil {
			continue
		}
		if provider := rmdko.Provider; provider != nil {
			rmdkOpts.Provider = provider
		}
		if masterKey := rmdko.MasterKey; masterKey != nil {
			rmdkOpts.MasterKey = masterKey
		}
	}
	return rmdkOpts
}
