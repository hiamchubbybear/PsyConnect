





package httputil

import (
	"net/http"
)


var DefaultHTTPClient = &http.Client{
	Transport: http.DefaultTransport.(*http.Transport).Clone(),
}






func CloseIdleHTTPConnections(client *http.Client) {
	type closeIdler interface {
		CloseIdleConnections()
	}
	if tr, ok := client.Transport.(closeIdler); ok {
		tr.CloseIdleConnections()
	}
}
