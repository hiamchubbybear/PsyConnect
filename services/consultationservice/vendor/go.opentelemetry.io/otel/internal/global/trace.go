


package global 



import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/auto/sdk"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)





type tracerProvider struct {
	embedded.TracerProvider

	mtx      sync.Mutex
	tracers  map[il]*tracer
	delegate trace.TracerProvider
}



var _ trace.TracerProvider = &tracerProvider{}








func (p *tracerProvider) setDelegate(provider trace.TracerProvider) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.delegate = provider

	if len(p.tracers) == 0 {
		return
	}

	for _, t := range p.tracers {
		t.setDelegate(provider)
	}

	p.tracers = nil
}


func (p *tracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.delegate != nil {
		return p.delegate.Tracer(name, opts...)
	}

	

	c := trace.NewTracerConfig(opts...)
	key := il{
		name:    name,
		version: c.InstrumentationVersion(),
		schema:  c.SchemaURL(),
		attrs:   c.InstrumentationAttributes(),
	}

	if p.tracers == nil {
		p.tracers = make(map[il]*tracer)
	}

	if val, ok := p.tracers[key]; ok {
		return val
	}

	t := &tracer{name: name, opts: opts, provider: p}
	p.tracers[key] = t
	return t
}

type il struct {
	name    string
	version string
	schema  string
	attrs   attribute.Set
}





type tracer struct {
	embedded.Tracer

	name     string
	opts     []trace.TracerOption
	provider *tracerProvider

	delegate atomic.Value
}


var _ trace.Tracer = &tracer{}







func (t *tracer) setDelegate(provider trace.TracerProvider) {
	t.delegate.Store(provider.Tracer(t.name, t.opts...))
}



func (t *tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	delegate := t.delegate.Load()
	if delegate != nil {
		return delegate.(trace.Tracer).Start(ctx, name, opts...)
	}

	return t.newSpan(ctx, autoInstEnabled, name, opts)
}








var autoInstEnabled = new(bool)

func (t *tracer) newSpan(ctx context.Context, autoSpan *bool, name string, opts []trace.SpanStartOption) (context.Context, trace.Span) {
	
	
	
	
	

	if *autoSpan {
		tracer := sdk.TracerProvider().Tracer(t.name, t.opts...)
		return tracer.Start(ctx, name, opts...)
	}

	s := nonRecordingSpan{sc: trace.SpanContextFromContext(ctx), tracer: t}
	ctx = trace.ContextWithSpan(ctx, s)
	return ctx, s
}




type nonRecordingSpan struct {
	embedded.Span

	sc     trace.SpanContext
	tracer *tracer
}

var _ trace.Span = nonRecordingSpan{}


func (s nonRecordingSpan) SpanContext() trace.SpanContext { return s.sc }


func (nonRecordingSpan) IsRecording() bool { return false }


func (nonRecordingSpan) SetStatus(codes.Code, string) {}


func (nonRecordingSpan) SetError(bool) {}


func (nonRecordingSpan) SetAttributes(...attribute.KeyValue) {}


func (nonRecordingSpan) End(...trace.SpanEndOption) {}


func (nonRecordingSpan) RecordError(error, ...trace.EventOption) {}


func (nonRecordingSpan) AddEvent(string, ...trace.EventOption) {}


func (nonRecordingSpan) AddLink(trace.Link) {}


func (nonRecordingSpan) SetName(string) {}

func (s nonRecordingSpan) TracerProvider() trace.TracerProvider { return s.tracer.provider }
