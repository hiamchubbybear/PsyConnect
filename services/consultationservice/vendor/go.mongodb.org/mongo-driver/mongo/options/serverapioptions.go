





package options

import (
	"fmt"
)










type ServerAPIOptions struct {
	ServerAPIVersion  ServerAPIVersion
	Strict            *bool
	DeprecationErrors *bool
}


func ServerAPI(serverAPIVersion ServerAPIVersion) *ServerAPIOptions {
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


type ServerAPIVersion string

const (
	
	ServerAPIVersion1 ServerAPIVersion = "1"
)


func (sav ServerAPIVersion) Validate() error {
	if sav == ServerAPIVersion1 {
		return nil
	}
	return fmt.Errorf("api version %q not supported; this driver version only supports API version \"1\"", sav)
}
