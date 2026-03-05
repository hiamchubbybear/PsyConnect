





package bsonoptions





type UIntCodecOptions struct {
	EncodeToMinSize *bool 
}





func UIntCodec() *UIntCodecOptions {
	return &UIntCodecOptions{}
}




func (u *UIntCodecOptions) SetEncodeToMinSize(b bool) *UIntCodecOptions {
	u.EncodeToMinSize = &b
	return u
}





func MergeUIntCodecOptions(opts ...*UIntCodecOptions) *UIntCodecOptions {
	u := UIntCodec()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.EncodeToMinSize != nil {
			u.EncodeToMinSize = opt.EncodeToMinSize
		}
	}

	return u
}
