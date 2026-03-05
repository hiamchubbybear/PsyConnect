


package http2

import "time"


type timer = interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}


type timeTimer struct {
	*time.Timer
}

func (t timeTimer) C() <-chan time.Time { return t.Timer.C }
