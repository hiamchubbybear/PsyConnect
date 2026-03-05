



















package atomic

//go:generate bin/gen-atomicwrapper -name=String -type=string -wrapped=Value -file=string.go




func (s *String) String() string {
	return s.Load()
}




func (s *String) MarshalText() ([]byte, error) {
	return []byte(s.Load()), nil
}




func (s *String) UnmarshalText(b []byte) error {
	s.Store(string(b))
	return nil
}
