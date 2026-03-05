





package bsonoptions

var defaultOverwriteDuplicatedInlinedFields = true





type StructCodecOptions struct {
	DecodeZeroStruct                 *bool 
	DecodeDeepZeroInline             *bool 
	EncodeOmitDefaultStruct          *bool 
	AllowUnexportedFields            *bool 
	OverwriteDuplicatedInlinedFields *bool 
}





func StructCodec() *StructCodecOptions {
	return &StructCodecOptions{}
}




func (t *StructCodecOptions) SetDecodeZeroStruct(b bool) *StructCodecOptions {
	t.DecodeZeroStruct = &b
	return t
}




func (t *StructCodecOptions) SetDecodeDeepZeroInline(b bool) *StructCodecOptions {
	t.DecodeDeepZeroInline = &b
	return t
}





func (t *StructCodecOptions) SetEncodeOmitDefaultStruct(b bool) *StructCodecOptions {
	t.EncodeOmitDefaultStruct = &b
	return t
}







func (t *StructCodecOptions) SetOverwriteDuplicatedInlinedFields(b bool) *StructCodecOptions {
	t.OverwriteDuplicatedInlinedFields = &b
	return t
}





func (t *StructCodecOptions) SetAllowUnexportedFields(b bool) *StructCodecOptions {
	t.AllowUnexportedFields = &b
	return t
}





func MergeStructCodecOptions(opts ...*StructCodecOptions) *StructCodecOptions {
	s := &StructCodecOptions{
		OverwriteDuplicatedInlinedFields: &defaultOverwriteDuplicatedInlinedFields,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.DecodeZeroStruct != nil {
			s.DecodeZeroStruct = opt.DecodeZeroStruct
		}
		if opt.DecodeDeepZeroInline != nil {
			s.DecodeDeepZeroInline = opt.DecodeDeepZeroInline
		}
		if opt.EncodeOmitDefaultStruct != nil {
			s.EncodeOmitDefaultStruct = opt.EncodeOmitDefaultStruct
		}
		if opt.OverwriteDuplicatedInlinedFields != nil {
			s.OverwriteDuplicatedInlinedFields = opt.OverwriteDuplicatedInlinedFields
		}
		if opt.AllowUnexportedFields != nil {
			s.AllowUnexportedFields = opt.AllowUnexportedFields
		}
	}

	return s
}
