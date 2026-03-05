





package options


type SearchIndexesOptions struct {
	Name *string
	Type *string
}


func SearchIndexes() *SearchIndexesOptions {
	return &SearchIndexesOptions{}
}


func (sio *SearchIndexesOptions) SetName(name string) *SearchIndexesOptions {
	sio.Name = &name
	return sio
}


func (sio *SearchIndexesOptions) SetType(typ string) *SearchIndexesOptions {
	sio.Type = &typ
	return sio
}



type CreateSearchIndexesOptions struct {
}


type ListSearchIndexesOptions struct {
	AggregateOpts *AggregateOptions
}


type DropSearchIndexOptions struct {
}


type UpdateSearchIndexOptions struct {
}
