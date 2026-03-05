//go:build !linux
// +build !linux





package syscall

import (
	"net"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"
)

var once sync.Once
var logger = grpclog.Component("core")

func log() {
	once.Do(func() {
		logger.Info("CPU time info is unavailable on non-linux environments.")
	})
}



func GetCPUTime() int64 {
	log()
	return 0
}


type Rusage struct{}


func GetRusage() *Rusage {
	log()
	return nil
}



func CPUTimeDiff(*Rusage, *Rusage) (float64, float64) {
	log()
	return 0, 0
}


func SetTCPUserTimeout(net.Conn, time.Duration) error {
	log()
	return nil
}



func GetTCPUserTimeout(net.Conn) (int, error) {
	log()
	return -1, nil
}
