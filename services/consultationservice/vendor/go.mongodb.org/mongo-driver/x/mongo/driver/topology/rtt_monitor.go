





package topology

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

const (
	rttAlphaValue = 0.2
	minSamples    = 10
	maxSamples    = 500
)

type rttConfig struct {
	
	
	interval time.Duration

	
	
	timeout time.Duration

	minRTTWindow       time.Duration
	createConnectionFn func() *connection
	createOperationFn  func(driver.Connection) *operation.Hello
}

type rttMonitor struct {
	mu sync.RWMutex 

	
	
	
	connMu        sync.Mutex
	samples       []time.Duration
	offset        int
	minRTT        time.Duration
	rtt90         time.Duration
	averageRTT    time.Duration
	averageRTTSet bool

	closeWg  sync.WaitGroup
	cfg      *rttConfig
	ctx      context.Context
	cancelFn context.CancelFunc
}

var _ driver.RTTMonitor = &rttMonitor{}

func newRTTMonitor(cfg *rttConfig) *rttMonitor {
	if cfg.interval <= 0 {
		panic("RTT monitor interval must be greater than 0")
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	
	numSamples := int(math.Max(minSamples, math.Min(maxSamples, float64((cfg.minRTTWindow)/cfg.interval))))

	return &rttMonitor{
		samples:  make([]time.Duration, numSamples),
		cfg:      cfg,
		ctx:      ctx,
		cancelFn: cancel,
	}
}

func (r *rttMonitor) connect() {
	r.connMu.Lock()
	defer r.connMu.Unlock()

	r.closeWg.Add(1)

	go func() {
		defer r.closeWg.Done()

		r.start()
	}()
}

func (r *rttMonitor) disconnect() {
	r.connMu.Lock()
	defer r.connMu.Unlock()

	r.cancelFn()

	
	r.closeWg.Wait()
}

func (r *rttMonitor) start() {
	var conn *connection
	defer func() {
		if conn != nil {
			
			
			
			
			conn.closeConnectContext()
			conn.wait()
			_ = conn.close()
		}
	}()

	ticker := time.NewTicker(r.cfg.interval)
	defer ticker.Stop()

	for {
		conn := r.cfg.createConnectionFn()
		err := conn.connect(r.ctx)

		
		
		
		if err == nil {
			r.addSample(conn.helloRTT)
			r.runHellos(conn)
		}

		
		
		_ = conn.close()

		
		
		select {
		case <-ticker.C:
		case <-r.ctx.Done():
			return
		}
	}
}



func (r *rttMonitor) runHellos(conn *connection) {
	ticker := time.NewTicker(r.cfg.interval)
	defer ticker.Stop()

	for {
		
		
		select {
		case <-ticker.C:
		case <-r.ctx.Done():
			return
		}

		
		
		
		
		
		
		
		timeout := r.cfg.timeout
		if timeout <= 0 {
			timeout = conn.config.connectTimeout
		}
		ctx, cancel := context.WithTimeout(r.ctx, timeout)

		start := time.Now()
		err := r.cfg.createOperationFn(initConnection{conn}).Execute(ctx)
		cancel()
		if err != nil {
			return
		}
		
		
		
		r.addSample(time.Since(start))
	}
}



func (r *rttMonitor) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.samples {
		r.samples[i] = 0
	}
	r.offset = 0
	r.minRTT = 0
	r.rtt90 = 0
	r.averageRTT = 0
	r.averageRTTSet = false
}

func (r *rttMonitor) addSample(rtt time.Duration) {
	
	
	r.mu.Lock()
	defer r.mu.Unlock()

	r.samples[r.offset] = rtt
	r.offset = (r.offset + 1) % len(r.samples)
	
	
	
	r.minRTT = min(r.samples, minSamples)
	r.rtt90 = percentile(90.0, r.samples, minSamples)

	if !r.averageRTTSet {
		r.averageRTT = rtt
		r.averageRTTSet = true
		return
	}

	r.averageRTT = time.Duration(rttAlphaValue*float64(rtt) + (1-rttAlphaValue)*float64(r.averageRTT))
}




func min(samples []time.Duration, minSamples int) time.Duration {
	count := 0
	min := time.Duration(math.MaxInt64)
	for _, d := range samples {
		if d > 0 {
			count++
		}
		if d > 0 && d < min {
			min = d
		}
	}
	if count == 0 || count < minSamples {
		return 0
	}
	return min
}




func percentile(perc float64, samples []time.Duration, minSamples int) time.Duration {
	
	floatSamples := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if sample > 0 {
			floatSamples = append(floatSamples, float64(sample))
		}
	}
	if len(floatSamples) == 0 || len(floatSamples) < minSamples {
		return 0
	}

	p, err := stats.Percentile(floatSamples, perc)
	if err != nil {
		panic(fmt.Errorf("x/mongo/driver/topology: error calculating %f percentile RTT: %w for samples:\n%v", perc, err, floatSamples))
	}
	return time.Duration(p)
}


func (r *rttMonitor) EWMA() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.averageRTT
}


func (r *rttMonitor) Min() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.minRTT
}


func (r *rttMonitor) P90() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.rtt90
}


func (r *rttMonitor) Stats() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	
	var sum float64
	floatSamples := make([]float64, 0, len(r.samples))
	for _, sample := range r.samples {
		if sample > 0 {
			floatSamples = append(floatSamples, float64(sample))
			sum += float64(sample)
		}
	}

	var avg, stdDev float64
	if len(floatSamples) > 0 {
		avg = sum / float64(len(floatSamples))

		var err error
		stdDev, err = stats.StandardDeviation(floatSamples)
		if err != nil {
			panic(fmt.Errorf("x/mongo/driver/topology: error calculating standard deviation RTT: %w for samples:\n%v", err, floatSamples))
		}
	}

	return fmt.Sprintf(
		"network round-trip time stats: avg: %v, min: %v, 90th pct: %v, stddev: %v",
		time.Duration(avg),
		r.minRTT,
		r.rtt90,
		time.Duration(stdDev))
}
