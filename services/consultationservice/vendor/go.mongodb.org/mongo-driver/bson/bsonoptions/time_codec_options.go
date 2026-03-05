





package bsonoptions





type TimeCodecOptions struct {
	UseLocalTimeZone *bool 
}





func TimeCodec() *TimeCodecOptions {
	return &TimeCodecOptions{}
}




func (t *TimeCodecOptions) SetUseLocalTimeZone(b bool) *TimeCodecOptions {
	t.UseLocalTimeZone = &b
	return t
}





func MergeTimeCodecOptions(opts ...*TimeCodecOptions) *TimeCodecOptions {
	t := TimeCodec()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.UseLocalTimeZone != nil {
			t.UseLocalTimeZone = opt.UseLocalTimeZone
		}
	}

	return t
}
