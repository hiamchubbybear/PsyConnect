







package backoff

import "time"


type Config struct {
	
	BaseDelay time.Duration
	
	
	Multiplier float64
	
	Jitter float64
	
	MaxDelay time.Duration
}






var DefaultConfig = Config{
	BaseDelay:  1.0 * time.Second,
	Multiplier: 1.6,
	Jitter:     0.2,
	MaxDelay:   120 * time.Second,
}
