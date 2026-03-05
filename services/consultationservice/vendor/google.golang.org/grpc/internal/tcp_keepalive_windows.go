//go:build windows



package internal

import (
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)







func NetDialerWithTCPKeepalive() *net.Dialer {
	return &net.Dialer{
		
		
		
		KeepAlive: time.Duration(-1),
		
		
		
		
		
		
		Control: func(_, _ string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_KEEPALIVE, 1)
			})
		},
	}
}
