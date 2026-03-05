



//go:build !go1.24

package http2

import "net/http"




func fillNetHTTPServerConfig(conf *http2Config, srv *http.Server) {}

func fillNetHTTPTransportConfig(conf *http2Config, tr *http.Transport) {}
