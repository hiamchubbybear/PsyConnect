





package bsonoptions





type SliceCodecOptions struct {
	EncodeNilAsEmpty *bool 
}





func SliceCodec() *SliceCodecOptions {
	return &SliceCodecOptions{}
}




func (s *SliceCodecOptions) SetEncodeNilAsEmpty(b bool) *SliceCodecOptions {
	s.EncodeNilAsEmpty = &b
	return s
}





func MergeSliceCodecOptions(opts ...*SliceCodecOptions) *SliceCodecOptions {
	s := SliceCodec()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.EncodeNilAsEmpty != nil {
			s.EncodeNilAsEmpty = opt.EncodeNilAsEmpty
		}
	}

	return s
}
