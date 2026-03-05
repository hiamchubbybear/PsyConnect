


package metric 

import (
	"context"

	"go.opentelemetry.io/otel/metric/embedded"
)






type Float64Counter interface {
	
	
	
	embedded.Float64Counter

	
	
	
	
	Add(ctx context.Context, incr float64, options ...AddOption)
}



type Float64CounterConfig struct {
	description string
	unit        string
}



func NewFloat64CounterConfig(opts ...Float64CounterOption) Float64CounterConfig {
	var config Float64CounterConfig
	for _, o := range opts {
		config = o.applyFloat64Counter(config)
	}
	return config
}


func (c Float64CounterConfig) Description() string {
	return c.description
}


func (c Float64CounterConfig) Unit() string {
	return c.unit
}




type Float64CounterOption interface {
	applyFloat64Counter(Float64CounterConfig) Float64CounterConfig
}







type Float64UpDownCounter interface {
	
	
	
	embedded.Float64UpDownCounter

	
	
	
	
	Add(ctx context.Context, incr float64, options ...AddOption)
}



type Float64UpDownCounterConfig struct {
	description string
	unit        string
}



func NewFloat64UpDownCounterConfig(opts ...Float64UpDownCounterOption) Float64UpDownCounterConfig {
	var config Float64UpDownCounterConfig
	for _, o := range opts {
		config = o.applyFloat64UpDownCounter(config)
	}
	return config
}


func (c Float64UpDownCounterConfig) Description() string {
	return c.description
}


func (c Float64UpDownCounterConfig) Unit() string {
	return c.unit
}




type Float64UpDownCounterOption interface {
	applyFloat64UpDownCounter(Float64UpDownCounterConfig) Float64UpDownCounterConfig
}







type Float64Histogram interface {
	
	
	
	embedded.Float64Histogram

	
	
	
	
	Record(ctx context.Context, incr float64, options ...RecordOption)
}



type Float64HistogramConfig struct {
	description              string
	unit                     string
	explicitBucketBoundaries []float64
}



func NewFloat64HistogramConfig(opts ...Float64HistogramOption) Float64HistogramConfig {
	var config Float64HistogramConfig
	for _, o := range opts {
		config = o.applyFloat64Histogram(config)
	}
	return config
}


func (c Float64HistogramConfig) Description() string {
	return c.description
}


func (c Float64HistogramConfig) Unit() string {
	return c.unit
}


func (c Float64HistogramConfig) ExplicitBucketBoundaries() []float64 {
	return c.explicitBucketBoundaries
}




type Float64HistogramOption interface {
	applyFloat64Histogram(Float64HistogramConfig) Float64HistogramConfig
}






type Float64Gauge interface {
	
	
	
	embedded.Float64Gauge

	
	
	
	
	Record(ctx context.Context, value float64, options ...RecordOption)
}



type Float64GaugeConfig struct {
	description string
	unit        string
}



func NewFloat64GaugeConfig(opts ...Float64GaugeOption) Float64GaugeConfig {
	var config Float64GaugeConfig
	for _, o := range opts {
		config = o.applyFloat64Gauge(config)
	}
	return config
}


func (c Float64GaugeConfig) Description() string {
	return c.description
}


func (c Float64GaugeConfig) Unit() string {
	return c.unit
}




type Float64GaugeOption interface {
	applyFloat64Gauge(Float64GaugeConfig) Float64GaugeConfig
}
