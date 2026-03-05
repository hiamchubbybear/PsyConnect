


package trace 

import (
	"context"

	"go.opentelemetry.io/otel/trace/embedded"
)






type Tracer interface {
	
	
	
	embedded.Tracer

	
	
	
	
	
	
	
	
	
	
	
	
	
	Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span)
}
