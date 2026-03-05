





package bsonoptions

var defaultDecodeOIDAsHex = true





type StringCodecOptions struct {
	DecodeObjectIDAsHex *bool 
}





func StringCodec() *StringCodecOptions {
	return &StringCodecOptions{}
}





func (t *StringCodecOptions) SetDecodeObjectIDAsHex(b bool) *StringCodecOptions {
	t.DecodeObjectIDAsHex = &b
	return t
}





func MergeStringCodecOptions(opts ...*StringCodecOptions) *StringCodecOptions {
	s := &StringCodecOptions{&defaultDecodeOIDAsHex}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.DecodeObjectIDAsHex != nil {
			s.DecodeObjectIDAsHex = opt.DecodeObjectIDAsHex
		}
	}

	return s
}
