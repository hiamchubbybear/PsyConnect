


package metric 

import "go.opentelemetry.io/otel/attribute"



type Observable interface {
	observable()
}


type InstrumentOption interface {
	Int64CounterOption
	Int64UpDownCounterOption
	Int64HistogramOption
	Int64GaugeOption
	Int64ObservableCounterOption
	Int64ObservableUpDownCounterOption
	Int64ObservableGaugeOption

	Float64CounterOption
	Float64UpDownCounterOption
	Float64HistogramOption
	Float64GaugeOption
	Float64ObservableCounterOption
	Float64ObservableUpDownCounterOption
	Float64ObservableGaugeOption
}


type HistogramOption interface {
	Int64HistogramOption
	Float64HistogramOption
}

type descOpt string

func (o descOpt) applyFloat64Counter(c Float64CounterConfig) Float64CounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64UpDownCounter(c Float64UpDownCounterConfig) Float64UpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64Histogram(c Float64HistogramConfig) Float64HistogramConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64Gauge(c Float64GaugeConfig) Float64GaugeConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableCounter(c Float64ObservableCounterConfig) Float64ObservableCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableUpDownCounter(c Float64ObservableUpDownCounterConfig) Float64ObservableUpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyFloat64ObservableGauge(c Float64ObservableGaugeConfig) Float64ObservableGaugeConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64Counter(c Int64CounterConfig) Int64CounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64UpDownCounter(c Int64UpDownCounterConfig) Int64UpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64Histogram(c Int64HistogramConfig) Int64HistogramConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64Gauge(c Int64GaugeConfig) Int64GaugeConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableCounter(c Int64ObservableCounterConfig) Int64ObservableCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableUpDownCounter(c Int64ObservableUpDownCounterConfig) Int64ObservableUpDownCounterConfig {
	c.description = string(o)
	return c
}

func (o descOpt) applyInt64ObservableGauge(c Int64ObservableGaugeConfig) Int64ObservableGaugeConfig {
	c.description = string(o)
	return c
}


func WithDescription(desc string) InstrumentOption { return descOpt(desc) }

type unitOpt string

func (o unitOpt) applyFloat64Counter(c Float64CounterConfig) Float64CounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64UpDownCounter(c Float64UpDownCounterConfig) Float64UpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64Histogram(c Float64HistogramConfig) Float64HistogramConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64Gauge(c Float64GaugeConfig) Float64GaugeConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableCounter(c Float64ObservableCounterConfig) Float64ObservableCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableUpDownCounter(c Float64ObservableUpDownCounterConfig) Float64ObservableUpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyFloat64ObservableGauge(c Float64ObservableGaugeConfig) Float64ObservableGaugeConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64Counter(c Int64CounterConfig) Int64CounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64UpDownCounter(c Int64UpDownCounterConfig) Int64UpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64Histogram(c Int64HistogramConfig) Int64HistogramConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64Gauge(c Int64GaugeConfig) Int64GaugeConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableCounter(c Int64ObservableCounterConfig) Int64ObservableCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableUpDownCounter(c Int64ObservableUpDownCounterConfig) Int64ObservableUpDownCounterConfig {
	c.unit = string(o)
	return c
}

func (o unitOpt) applyInt64ObservableGauge(c Int64ObservableGaugeConfig) Int64ObservableGaugeConfig {
	c.unit = string(o)
	return c
}




func WithUnit(u string) InstrumentOption { return unitOpt(u) }




func WithExplicitBucketBoundaries(bounds ...float64) HistogramOption { return bucketOpt(bounds) }

type bucketOpt []float64

func (o bucketOpt) applyFloat64Histogram(c Float64HistogramConfig) Float64HistogramConfig {
	c.explicitBucketBoundaries = o
	return c
}

func (o bucketOpt) applyInt64Histogram(c Int64HistogramConfig) Int64HistogramConfig {
	c.explicitBucketBoundaries = o
	return c
}



type AddOption interface {
	applyAdd(AddConfig) AddConfig
}


type AddConfig struct {
	attrs attribute.Set
}


func NewAddConfig(opts []AddOption) AddConfig {
	config := AddConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyAdd(config)
	}
	return config
}


func (c AddConfig) Attributes() attribute.Set {
	return c.attrs
}



type RecordOption interface {
	applyRecord(RecordConfig) RecordConfig
}


type RecordConfig struct {
	attrs attribute.Set
}


func NewRecordConfig(opts []RecordOption) RecordConfig {
	config := RecordConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyRecord(config)
	}
	return config
}


func (c RecordConfig) Attributes() attribute.Set {
	return c.attrs
}



type ObserveOption interface {
	applyObserve(ObserveConfig) ObserveConfig
}


type ObserveConfig struct {
	attrs attribute.Set
}


func NewObserveConfig(opts []ObserveOption) ObserveConfig {
	config := ObserveConfig{attrs: *attribute.EmptySet()}
	for _, o := range opts {
		config = o.applyObserve(config)
	}
	return config
}


func (c ObserveConfig) Attributes() attribute.Set {
	return c.attrs
}


type MeasurementOption interface {
	AddOption
	RecordOption
	ObserveOption
}

type attrOpt struct {
	set attribute.Set
}



func mergeSets(a, b attribute.Set) attribute.Set {
	
	iter := attribute.NewMergeIterator(&b, &a)
	merged := make([]attribute.KeyValue, 0, a.Len()+b.Len())
	for iter.Next() {
		merged = append(merged, iter.Attribute())
	}
	return attribute.NewSet(merged...)
}

func (o attrOpt) applyAdd(c AddConfig) AddConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyRecord(c RecordConfig) RecordConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}

func (o attrOpt) applyObserve(c ObserveConfig) ObserveConfig {
	switch {
	case o.set.Len() == 0:
	case c.attrs.Len() == 0:
		c.attrs = o.set
	default:
		c.attrs = mergeSets(c.attrs, o.set)
	}
	return c
}







func WithAttributeSet(attributes attribute.Set) MeasurementOption {
	return attrOpt{set: attributes}
}
















func WithAttributes(attributes ...attribute.KeyValue) MeasurementOption {
	cp := make([]attribute.KeyValue, len(attributes))
	copy(cp, attributes)
	return attrOpt{set: attribute.NewSet(cp...)}
}
