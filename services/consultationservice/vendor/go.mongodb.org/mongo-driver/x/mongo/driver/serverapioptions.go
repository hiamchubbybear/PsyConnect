





package driver



const TestServerAPIVersion = "1"



type ServerAPIOptions struct {
	ServerAPIVersion  string
	Strict            *bool
	DeprecationErrors *bool
}


func NewServerAPIOptions(serverAPIVersion string) *ServerAPIOptions {
	return &ServerAPIOptions{ServerAPIVersion: serverAPIVersion}
}


func (s *ServerAPIOptions) SetStrict(strict bool) *ServerAPIOptions {
	s.Strict = &strict
	return s
}


func (s *ServerAPIOptions) SetDeprecationErrors(deprecationErrors bool) *ServerAPIOptions {
	s.DeprecationErrors = &deprecationErrors
	return s
}
