


package trace 

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace/embedded"
)







func NewNoopTracerProvider() TracerProvider {
	return noopTracerProvider{}
}

type noopTracerProvider struct{ embedded.TracerProvider }

var _ TracerProvider = noopTracerProvider{}


func (p noopTracerProvider) Tracer(string, ...TracerOption) Tracer {
	return noopTracer{}
}


type noopTracer struct{ embedded.Tracer }

var _ Tracer = noopTracer{}



func (t noopTracer) Start(ctx context.Context, name string, _ ...SpanStartOption) (context.Context, Span) {
	span := SpanFromContext(ctx)
	if _, ok := span.(nonRecordingSpan); !ok {
		
		span = noopSpanInstance
	}
	return ContextWithSpan(ctx, span), span
}


type noopSpan struct{ embedded.Span }

var noopSpanInstance Span = noopSpan{}


func (noopSpan) SpanContext() SpanContext { return SpanContext{} }


func (noopSpan) IsRecording() bool { return false }


func (noopSpan) SetStatus(codes.Code, string) {}


func (noopSpan) SetError(bool) {}


func (noopSpan) SetAttributes(...attribute.KeyValue) {}


func (noopSpan) End(...SpanEndOption) {}


func (noopSpan) RecordError(error, ...EventOption) {}


func (noopSpan) AddEvent(string, ...EventOption) {}


func (noopSpan) AddLink(Link) {}


func (noopSpan) SetName(string) {}


func (noopSpan) TracerProvider() TracerProvider { return noopTracerProvider{} }
