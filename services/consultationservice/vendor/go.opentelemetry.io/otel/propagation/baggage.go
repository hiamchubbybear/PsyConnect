


package propagation 

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

const baggageHeader = "baggage"





type Baggage struct{}

var _ TextMapPropagator = Baggage{}


func (b Baggage) Inject(ctx context.Context, carrier TextMapCarrier) {
	bStr := baggage.FromContext(ctx).String()
	if bStr != "" {
		carrier.Set(baggageHeader, bStr)
	}
}


func (b Baggage) Extract(parent context.Context, carrier TextMapCarrier) context.Context {
	bStr := carrier.Get(baggageHeader)
	if bStr == "" {
		return parent
	}

	bag, err := baggage.Parse(bStr)
	if err != nil {
		return parent
	}
	return baggage.ContextWithBaggage(parent, bag)
}


func (b Baggage) Fields() []string {
	return []string{baggageHeader}
}
