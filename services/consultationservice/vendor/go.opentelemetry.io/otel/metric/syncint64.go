


package metric 

import (
	"context"

	"go.opentelemetry.io/otel/metric/embedded"
)






type Int64Counter interface {
	
	
	
	embedded.Int64Counter

	
	
	
	
	Add(ctx context.Context, incr int64, options ...AddOption)
}



type Int64CounterConfig struct {
	description string
	unit        string
}



func NewInt64CounterConfig(opts ...Int64CounterOption) Int64CounterConfig {
	var config Int64CounterConfig
	for _, o := range opts {
		config = o.applyInt64Counter(config)
	}
	return config
}


func (c Int64CounterConfig) Description() string {
	return c.description
}


func (c Int64CounterConfig) Unit() string {
	return c.unit
}




type Int64CounterOption interface {
	applyInt64Counter(Int64CounterConfig) Int64CounterConfig
}







type Int64UpDownCounter interface {
	
	
	
	embedded.Int64UpDownCounter

	
	
	
	
	Add(ctx context.Context, incr int64, options ...AddOption)
}



type Int64UpDownCounterConfig struct {
	description string
	unit        string
}



func NewInt64UpDownCounterConfig(opts ...Int64UpDownCounterOption) Int64UpDownCounterConfig {
	var config Int64UpDownCounterConfig
	for _, o := range opts {
		config = o.applyInt64UpDownCounter(config)
	}
	return config
}


func (c Int64UpDownCounterConfig) Description() string {
	return c.description
}


func (c Int64UpDownCounterConfig) Unit() string {
	return c.unit
}




type Int64UpDownCounterOption interface {
	applyInt64UpDownCounter(Int64UpDownCounterConfig) Int64UpDownCounterConfig
}







type Int64Histogram interface {
	
	
	
	embedded.Int64Histogram

	
	
	
	
	Record(ctx context.Context, incr int64, options ...RecordOption)
}



type Int64HistogramConfig struct {
	description              string
	unit                     string
	explicitBucketBoundaries []float64
}



func NewInt64HistogramConfig(opts ...Int64HistogramOption) Int64HistogramConfig {
	var config Int64HistogramConfig
	for _, o := range opts {
		config = o.applyInt64Histogram(config)
	}
	return config
}


func (c Int64HistogramConfig) Description() string {
	return c.description
}


func (c Int64HistogramConfig) Unit() string {
	return c.unit
}


func (c Int64HistogramConfig) ExplicitBucketBoundaries() []float64 {
	return c.explicitBucketBoundaries
}




type Int64HistogramOption interface {
	applyInt64Histogram(Int64HistogramConfig) Int64HistogramConfig
}






type Int64Gauge interface {
	
	
	
	embedded.Int64Gauge

	
	
	
	
	Record(ctx context.Context, value int64, options ...RecordOption)
}



type Int64GaugeConfig struct {
	description string
	unit        string
}



func NewInt64GaugeConfig(opts ...Int64GaugeOption) Int64GaugeConfig {
	var config Int64GaugeConfig
	for _, o := range opts {
		config = o.applyInt64Gauge(config)
	}
	return config
}


func (c Int64GaugeConfig) Description() string {
	return c.description
}


func (c Int64GaugeConfig) Unit() string {
	return c.unit
}




type Int64GaugeOption interface {
	applyInt64Gauge(Int64GaugeConfig) Int64GaugeConfig
}
