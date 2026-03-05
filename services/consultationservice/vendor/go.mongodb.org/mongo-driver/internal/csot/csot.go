





package csot

import (
	"context"
	"time"
)

type timeoutKey struct{}







func MakeTimeoutContext(ctx context.Context, to time.Duration) (context.Context, context.CancelFunc) {
	
	
	cancelFunc := func() {}
	if _, deadlineSet := ctx.Deadline(); to != 0 && !deadlineSet {
		ctx, cancelFunc = context.WithTimeout(ctx, to)
	}

	
	return context.WithValue(ctx, timeoutKey{}, true), cancelFunc
}

func IsTimeoutContext(ctx context.Context) bool {
	return ctx.Value(timeoutKey{}) != nil
}



type ZeroRTTMonitor struct{}


func (zrm *ZeroRTTMonitor) EWMA() time.Duration {
	return 0
}


func (zrm *ZeroRTTMonitor) Min() time.Duration {
	return 0
}


func (zrm *ZeroRTTMonitor) P90() time.Duration {
	return 0
}


func (zrm *ZeroRTTMonitor) Stats() string {
	return ""
}
