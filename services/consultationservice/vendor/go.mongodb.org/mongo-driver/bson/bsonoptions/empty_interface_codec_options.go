





package bsonoptions





type EmptyInterfaceCodecOptions struct {
	DecodeBinaryAsSlice *bool 
}





func EmptyInterfaceCodec() *EmptyInterfaceCodecOptions {
	return &EmptyInterfaceCodecOptions{}
}




func (e *EmptyInterfaceCodecOptions) SetDecodeBinaryAsSlice(b bool) *EmptyInterfaceCodecOptions {
	e.DecodeBinaryAsSlice = &b
	return e
}





func MergeEmptyInterfaceCodecOptions(opts ...*EmptyInterfaceCodecOptions) *EmptyInterfaceCodecOptions {
	e := EmptyInterfaceCodec()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.DecodeBinaryAsSlice != nil {
			e.DecodeBinaryAsSlice = opt.DecodeBinaryAsSlice
		}
	}

	return e
}
