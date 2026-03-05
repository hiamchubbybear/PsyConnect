


package internal

import (
	rand "math/rand/v2"
	"time"
)

var (
	
	RandShuffle = rand.Shuffle
	
	
	TimeAfterFunc = func(d time.Duration, f func()) func() {
		timer := time.AfterFunc(d, f)
		return func() { timer.Stop() }
	}
)
