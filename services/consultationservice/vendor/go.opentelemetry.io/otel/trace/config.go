


package trace 

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)


type TracerConfig struct {
	instrumentationVersion string
	
	schemaURL string
	attrs     attribute.Set
}


func (t *TracerConfig) InstrumentationVersion() string {
	return t.instrumentationVersion
}



func (t *TracerConfig) InstrumentationAttributes() attribute.Set {
	return t.attrs
}


func (t *TracerConfig) SchemaURL() string {
	return t.schemaURL
}


func NewTracerConfig(options ...TracerOption) TracerConfig {
	var config TracerConfig
	for _, option := range options {
		config = option.apply(config)
	}
	return config
}


type TracerOption interface {
	apply(TracerConfig) TracerConfig
}

type tracerOptionFunc func(TracerConfig) TracerConfig

func (fn tracerOptionFunc) apply(cfg TracerConfig) TracerConfig {
	return fn(cfg)
}


type SpanConfig struct {
	attributes []attribute.KeyValue
	timestamp  time.Time
	links      []Link
	newRoot    bool
	spanKind   SpanKind
	stackTrace bool
}


func (cfg *SpanConfig) Attributes() []attribute.KeyValue {
	return cfg.attributes
}


func (cfg *SpanConfig) Timestamp() time.Time {
	return cfg.timestamp
}


func (cfg *SpanConfig) StackTrace() bool {
	return cfg.stackTrace
}


func (cfg *SpanConfig) Links() []Link {
	return cfg.links
}




func (cfg *SpanConfig) NewRoot() bool {
	return cfg.newRoot
}


func (cfg *SpanConfig) SpanKind() SpanKind {
	return cfg.spanKind
}





func NewSpanStartConfig(options ...SpanStartOption) SpanConfig {
	var c SpanConfig
	for _, option := range options {
		c = option.applySpanStart(c)
	}
	return c
}





func NewSpanEndConfig(options ...SpanEndOption) SpanConfig {
	var c SpanConfig
	for _, option := range options {
		c = option.applySpanEnd(c)
	}
	return c
}



type SpanStartOption interface {
	applySpanStart(SpanConfig) SpanConfig
}

type spanOptionFunc func(SpanConfig) SpanConfig

func (fn spanOptionFunc) applySpanStart(cfg SpanConfig) SpanConfig {
	return fn(cfg)
}



type SpanEndOption interface {
	applySpanEnd(SpanConfig) SpanConfig
}


type EventConfig struct {
	attributes []attribute.KeyValue
	timestamp  time.Time
	stackTrace bool
}


func (cfg *EventConfig) Attributes() []attribute.KeyValue {
	return cfg.attributes
}


func (cfg *EventConfig) Timestamp() time.Time {
	return cfg.timestamp
}


func (cfg *EventConfig) StackTrace() bool {
	return cfg.stackTrace
}





func NewEventConfig(options ...EventOption) EventConfig {
	var c EventConfig
	for _, option := range options {
		c = option.applyEvent(c)
	}
	if c.timestamp.IsZero() {
		c.timestamp = time.Now()
	}
	return c
}


type EventOption interface {
	applyEvent(EventConfig) EventConfig
}


type SpanOption interface {
	SpanStartOption
	SpanEndOption
}


type SpanStartEventOption interface {
	SpanStartOption
	EventOption
}


type SpanEndEventOption interface {
	SpanEndOption
	EventOption
}

type attributeOption []attribute.KeyValue

func (o attributeOption) applySpan(c SpanConfig) SpanConfig {
	c.attributes = append(c.attributes, []attribute.KeyValue(o)...)
	return c
}
func (o attributeOption) applySpanStart(c SpanConfig) SpanConfig { return o.applySpan(c) }
func (o attributeOption) applyEvent(c EventConfig) EventConfig {
	c.attributes = append(c.attributes, []attribute.KeyValue(o)...)
	return c
}

var _ SpanStartEventOption = attributeOption{}










func WithAttributes(attributes ...attribute.KeyValue) SpanStartEventOption {
	return attributeOption(attributes)
}


type SpanEventOption interface {
	SpanOption
	EventOption
}

type timestampOption time.Time

func (o timestampOption) applySpan(c SpanConfig) SpanConfig {
	c.timestamp = time.Time(o)
	return c
}
func (o timestampOption) applySpanStart(c SpanConfig) SpanConfig { return o.applySpan(c) }
func (o timestampOption) applySpanEnd(c SpanConfig) SpanConfig   { return o.applySpan(c) }
func (o timestampOption) applyEvent(c EventConfig) EventConfig {
	c.timestamp = time.Time(o)
	return c
}

var _ SpanEventOption = timestampOption{}



func WithTimestamp(t time.Time) SpanEventOption {
	return timestampOption(t)
}

type stackTraceOption bool

func (o stackTraceOption) applyEvent(c EventConfig) EventConfig {
	c.stackTrace = bool(o)
	return c
}

func (o stackTraceOption) applySpan(c SpanConfig) SpanConfig {
	c.stackTrace = bool(o)
	return c
}
func (o stackTraceOption) applySpanEnd(c SpanConfig) SpanConfig { return o.applySpan(c) }


func WithStackTrace(b bool) SpanEndEventOption {
	return stackTraceOption(b)
}



func WithLinks(links ...Link) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.links = append(cfg.links, links...)
		return cfg
	})
}




func WithNewRoot() SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.newRoot = true
		return cfg
	})
}


func WithSpanKind(kind SpanKind) SpanStartOption {
	return spanOptionFunc(func(cfg SpanConfig) SpanConfig {
		cfg.spanKind = kind
		return cfg
	})
}


func WithInstrumentationVersion(version string) TracerOption {
	return tracerOptionFunc(func(cfg TracerConfig) TracerConfig {
		cfg.instrumentationVersion = version
		return cfg
	})
}




func WithInstrumentationAttributes(attr ...attribute.KeyValue) TracerOption {
	return tracerOptionFunc(func(config TracerConfig) TracerConfig {
		config.attrs = attribute.NewSet(attr...)
		return config
	})
}


func WithSchemaURL(schemaURL string) TracerOption {
	return tracerOptionFunc(func(cfg TracerConfig) TracerConfig {
		cfg.schemaURL = schemaURL
		return cfg
	})
}
