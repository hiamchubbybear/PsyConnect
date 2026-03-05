





package backoff

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"time"

	grpcbackoff "google.golang.org/grpc/backoff"
)



type Strategy interface {
	
	
	Backoff(retries int) time.Duration
}




var DefaultExponential = Exponential{Config: grpcbackoff.DefaultConfig}



type Exponential struct {
	
	Config grpcbackoff.Config
}



func (bc Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.Config.BaseDelay
	}
	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Config.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	
	
	backoff *= 1 + bc.Config.Jitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}



var ErrResetBackoff = errors.New("reset backoff state")






func RunF(ctx context.Context, f func() error, backoff func(int) time.Duration) {
	attempt := 0
	timer := time.NewTimer(0)
	for ctx.Err() == nil {
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return
		}

		err := f()
		if errors.Is(err, ErrResetBackoff) {
			timer.Reset(0)
			attempt = 0
			continue
		}
		if err != nil {
			return
		}
		timer.Reset(backoff(attempt))
		attempt++
	}
}
