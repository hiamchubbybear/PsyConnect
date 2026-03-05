





package bsonoptions





type ByteSliceCodecOptions struct {
	EncodeNilAsEmpty *bool 
}





func ByteSliceCodec() *ByteSliceCodecOptions {
	return &ByteSliceCodecOptions{}
}




func (bs *ByteSliceCodecOptions) SetEncodeNilAsEmpty(b bool) *ByteSliceCodecOptions {
	bs.EncodeNilAsEmpty = &b
	return bs
}





func MergeByteSliceCodecOptions(opts ...*ByteSliceCodecOptions) *ByteSliceCodecOptions {
	bs := ByteSliceCodec()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.EncodeNilAsEmpty != nil {
			bs.EncodeNilAsEmpty = opt.EncodeNilAsEmpty
		}
	}

	return bs
}
