

package channelz

import (
	"syscall"

	"golang.org/x/sys/unix"
)



type SocketOptionData struct {
	Linger      *unix.Linger
	RecvTimeout *unix.Timeval
	SendTimeout *unix.Timeval
	TCPInfo     *unix.TCPInfo
}



func (s *SocketOptionData) Getsockopt(fd uintptr) {
	if v, err := unix.GetsockoptLinger(int(fd), syscall.SOL_SOCKET, syscall.SO_LINGER); err == nil {
		s.Linger = v
	}
	if v, err := unix.GetsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO); err == nil {
		s.RecvTimeout = v
	}
	if v, err := unix.GetsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDTIMEO); err == nil {
		s.SendTimeout = v
	}
	if v, err := unix.GetsockoptTCPInfo(int(fd), syscall.SOL_TCP, syscall.TCP_INFO); err == nil {
		s.TCPInfo = v
	}
}


func GetSocketOption(socket any) *SocketOptionData {
	c, ok := socket.(syscall.Conn)
	if !ok {
		return nil
	}
	data := &SocketOptionData{}
	if rawConn, err := c.SyscallConn(); err == nil {
		rawConn.Control(data.Getsockopt)
		return data
	}
	return nil
}
