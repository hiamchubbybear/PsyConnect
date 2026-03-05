





package ocsp

import "net/http"


type VerifyOptions struct {
	Cache                   Cache
	DisableEndpointChecking bool
	HTTPClient              *http.Client
}
