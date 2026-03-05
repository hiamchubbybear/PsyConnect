





package options


type ListCollectionsOptions struct {
	
	NameOnly *bool

	
	BatchSize *int32

	
	
	AuthorizedCollections *bool
}


func ListCollections() *ListCollectionsOptions {
	return &ListCollectionsOptions{}
}


func (lc *ListCollectionsOptions) SetNameOnly(b bool) *ListCollectionsOptions {
	lc.NameOnly = &b
	return lc
}


func (lc *ListCollectionsOptions) SetBatchSize(size int32) *ListCollectionsOptions {
	lc.BatchSize = &size
	return lc
}



func (lc *ListCollectionsOptions) SetAuthorizedCollections(b bool) *ListCollectionsOptions {
	lc.AuthorizedCollections = &b
	return lc
}






func MergeListCollectionsOptions(opts ...*ListCollectionsOptions) *ListCollectionsOptions {
	lc := ListCollections()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.NameOnly != nil {
			lc.NameOnly = opt.NameOnly
		}
		if opt.BatchSize != nil {
			lc.BatchSize = opt.BatchSize
		}
		if opt.AuthorizedCollections != nil {
			lc.AuthorizedCollections = opt.AuthorizedCollections
		}
	}

	return lc
}
