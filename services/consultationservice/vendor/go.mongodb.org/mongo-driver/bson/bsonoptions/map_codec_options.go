





package bsonoptions





type MapCodecOptions struct {
	DecodeZerosMap   *bool 
	EncodeNilAsEmpty *bool 
	
	
	
	
	
	EncodeKeysWithStringer *bool
}





func MapCodec() *MapCodecOptions {
	return &MapCodecOptions{}
}




func (t *MapCodecOptions) SetDecodeZerosMap(b bool) *MapCodecOptions {
	t.DecodeZerosMap = &b
	return t
}




func (t *MapCodecOptions) SetEncodeNilAsEmpty(b bool) *MapCodecOptions {
	t.EncodeNilAsEmpty = &b
	return t
}








func (t *MapCodecOptions) SetEncodeKeysWithStringer(b bool) *MapCodecOptions {
	t.EncodeKeysWithStringer = &b
	return t
}





func MergeMapCodecOptions(opts ...*MapCodecOptions) *MapCodecOptions {
	s := MapCodec()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.DecodeZerosMap != nil {
			s.DecodeZerosMap = opt.DecodeZerosMap
		}
		if opt.EncodeNilAsEmpty != nil {
			s.EncodeNilAsEmpty = opt.EncodeNilAsEmpty
		}
		if opt.EncodeKeysWithStringer != nil {
			s.EncodeKeysWithStringer = opt.EncodeKeysWithStringer
		}
	}

	return s
}
