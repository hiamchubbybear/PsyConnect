//go:build unix



package internal

import (
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)







func NetDialerWithTCPKeepalive() *net.Dialer {
	return &net.Dialer{
		
		
		
		KeepAlive: time.Duration(-1),
		
		
		
		
		
		
		Control: func(_, _ string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1)
			})
		},
	}
}
