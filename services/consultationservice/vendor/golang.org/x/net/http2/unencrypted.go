



package http2

import (
	"crypto/tls"
	"errors"
	"net"
)

const nextProtoUnencryptedHTTP2 = "unencrypted_http2"










func unencryptedNetConnFromTLSConn(tc *tls.Conn) (net.Conn, error) {
	conner, ok := tc.NetConn().(interface {
		UnencryptedNetConn() net.Conn
	})
	if !ok {
		return nil, errors.New("http2: TLS conn unexpectedly found in unencrypted handoff")
	}
	return conner.UnencryptedNetConn(), nil
}
