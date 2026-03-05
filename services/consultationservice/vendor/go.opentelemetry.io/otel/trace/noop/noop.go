











package noop 

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

var (
	

	_ trace.TracerProvider = TracerProvider{}
	_ trace.Tracer         = Tracer{}
	_ trace.Span           = Span{}
)


type TracerProvider struct{ embedded.TracerProvider }


func NewTracerProvider() TracerProvider {
	return TracerProvider{}
}


func (TracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return Tracer{}
}


type Tracer struct{ embedded.Tracer }







func (t Tracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)

	
	
	
	var zeroSC trace.SpanContext
	if sc := span.SpanContext(); !sc.Equal(zeroSC) {
		if !span.IsRecording() {
			
			return ctx, span
		}
		
		span = Span{sc: sc}
	} else {
		
		span = noopSpanInstance
	}
	return trace.ContextWithSpan(ctx, span), span
}

var noopSpanInstance trace.Span = Span{}


type Span struct {
	embedded.Span

	sc trace.SpanContext
}


func (s Span) SpanContext() trace.SpanContext { return s.sc }


func (Span) IsRecording() bool { return false }


func (Span) SetStatus(codes.Code, string) {}


func (Span) SetAttributes(...attribute.KeyValue) {}


func (Span) End(...trace.SpanEndOption) {}


func (Span) RecordError(error, ...trace.EventOption) {}


func (Span) AddEvent(string, ...trace.EventOption) {}


func (Span) AddLink(trace.Link) {}


func (Span) SetName(string) {}


func (Span) TracerProvider() trace.TracerProvider { return TracerProvider{} }
