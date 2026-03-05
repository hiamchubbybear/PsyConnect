


package baggage 

import (
	"context"

	"go.opentelemetry.io/otel/internal/baggage"
)


func ContextWithBaggage(parent context.Context, b Baggage) context.Context {
	
	return baggage.ContextWithList(parent, b.list)
}


func ContextWithoutBaggage(parent context.Context) context.Context {
	
	return baggage.ContextWithList(parent, nil)
}


func FromContext(ctx context.Context) Baggage {
	
	return Baggage{list: baggage.ListFromContext(ctx)}
}
