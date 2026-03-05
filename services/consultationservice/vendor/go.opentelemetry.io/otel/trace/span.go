


package trace 

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace/embedded"
)









type Span interface {
	
	
	
	embedded.Span

	
	
	
	
	End(options ...SpanEndOption)

	
	AddEvent(name string, options ...EventOption)

	
	
	
	
	AddLink(link Link)

	
	
	IsRecording() bool

	
	
	
	
	RecordError(err error, options ...EventOption)

	
	
	SpanContext() SpanContext

	
	
	
	
	SetStatus(code codes.Code, description string)

	
	SetName(name string)

	
	
	
	SetAttributes(kv ...attribute.KeyValue)

	
	
	TracerProvider() TracerProvider
}
















type Link struct {
	
	SpanContext SpanContext

	
	Attributes []attribute.KeyValue
}



func LinkFromContext(ctx context.Context, attrs ...attribute.KeyValue) Link {
	return Link{
		SpanContext: SpanContextFromContext(ctx),
		Attributes:  attrs,
	}
}


type SpanKind int






const (
	
	
	
	SpanKindUnspecified SpanKind = 0
	
	
	SpanKindInternal SpanKind = 1
	
	
	SpanKindServer SpanKind = 2
	
	
	SpanKindClient SpanKind = 3
	
	
	
	
	
	
	
	SpanKindProducer SpanKind = 4
	
	
	
	
	SpanKindConsumer SpanKind = 5
)



func ValidateSpanKind(spanKind SpanKind) SpanKind {
	switch spanKind {
	case SpanKindInternal,
		SpanKindServer,
		SpanKindClient,
		SpanKindProducer,
		SpanKindConsumer:
		
		return spanKind
	default:
		return SpanKindInternal
	}
}


func (sk SpanKind) String() string {
	switch sk {
	case SpanKindInternal:
		return "internal"
	case SpanKindServer:
		return "server"
	case SpanKindClient:
		return "client"
	case SpanKindProducer:
		return "producer"
	case SpanKindConsumer:
		return "consumer"
	default:
		return "unspecified"
	}
}
