//go:build !unix && !windows



package internal

import (
	"net"
)


func NetDialerWithTCPKeepalive() *net.Dialer {
	return &net.Dialer{}
}
