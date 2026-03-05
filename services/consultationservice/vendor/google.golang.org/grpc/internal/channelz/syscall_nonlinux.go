//go:build !linux



package channelz

import (
	"sync"
)

var once sync.Once




type SocketOptionData struct {
}




func (s *SocketOptionData) Getsockopt(uintptr) {
	once.Do(func() {
		logger.Warning("Channelz: socket options are not supported on non-linux environments")
	})
}


func GetSocketOption(any) *SocketOptionData {
	return nil
}
