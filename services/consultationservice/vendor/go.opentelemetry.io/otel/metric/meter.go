


package metric 

import (
	"context"

	"go.opentelemetry.io/otel/metric/embedded"
)







type MeterProvider interface {
	
	
	
	embedded.MeterProvider

	
	
	
	
	
	
	
	
	
	Meter(name string, opts ...MeterOption) Meter
}






type Meter interface {
	
	
	
	embedded.Meter

	
	
	
	
	
	
	
	Int64Counter(name string, options ...Int64CounterOption) (Int64Counter, error)

	
	
	
	
	
	
	
	
	Int64UpDownCounter(name string, options ...Int64UpDownCounterOption) (Int64UpDownCounter, error)

	
	
	
	
	
	
	
	
	Int64Histogram(name string, options ...Int64HistogramOption) (Int64Histogram, error)

	
	
	
	
	
	
	
	Int64Gauge(name string, options ...Int64GaugeOption) (Int64Gauge, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	Int64ObservableCounter(name string, options ...Int64ObservableCounterOption) (Int64ObservableCounter, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	Int64ObservableUpDownCounter(name string, options ...Int64ObservableUpDownCounterOption) (Int64ObservableUpDownCounter, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	Int64ObservableGauge(name string, options ...Int64ObservableGaugeOption) (Int64ObservableGauge, error)

	
	
	
	
	
	
	
	
	Float64Counter(name string, options ...Float64CounterOption) (Float64Counter, error)

	
	
	
	
	
	
	
	
	Float64UpDownCounter(name string, options ...Float64UpDownCounterOption) (Float64UpDownCounter, error)

	
	
	
	
	
	
	
	
	Float64Histogram(name string, options ...Float64HistogramOption) (Float64Histogram, error)

	
	
	
	
	
	
	
	Float64Gauge(name string, options ...Float64GaugeOption) (Float64Gauge, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	Float64ObservableCounter(name string, options ...Float64ObservableCounterOption) (Float64ObservableCounter, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	Float64ObservableUpDownCounter(name string, options ...Float64ObservableUpDownCounterOption) (Float64ObservableUpDownCounter, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	Float64ObservableGauge(name string, options ...Float64ObservableGaugeOption) (Float64ObservableGauge, error)

	
	
	
	
	
	
	
	
	
	
	
	
	
	RegisterCallback(f Callback, instruments ...Observable) (Registration, error)
}













type Callback func(context.Context, Observer) error






type Observer interface {
	
	
	
	embedded.Observer

	
	ObserveFloat64(obsrv Float64Observable, value float64, opts ...ObserveOption)

	
	ObserveInt64(obsrv Int64Observable, value int64, opts ...ObserveOption)
}







type Registration interface {
	
	
	
	embedded.Registration

	
	
	
	Unregister() error
}
