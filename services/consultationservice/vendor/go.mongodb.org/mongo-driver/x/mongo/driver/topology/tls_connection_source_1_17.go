





//go:build go1.17
// +build go1.17

package topology

import (
	"context"
	"crypto/tls"
	"net"
)

type tlsConn interface {
	net.Conn

	
	HandshakeContext(ctx context.Context) error
	ConnectionState() tls.ConnectionState
}

var _ tlsConn = (*tls.Conn)(nil)

type tlsConnectionSource interface {
	Client(net.Conn, *tls.Config) tlsConn
}

type tlsConnectionSourceFn func(net.Conn, *tls.Config) tlsConn

var _ tlsConnectionSource = (tlsConnectionSourceFn)(nil)

func (t tlsConnectionSourceFn) Client(nc net.Conn, cfg *tls.Config) tlsConn {
	return t(nc, cfg)
}

var defaultTLSConnectionSource tlsConnectionSourceFn = func(nc net.Conn, cfg *tls.Config) tlsConn {
	return tls.Client(nc, cfg)
}


func clientHandshake(ctx context.Context, client tlsConn) error {
	return client.HandshakeContext(ctx)
}
