


package otel 

import (
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/propagation"
)



func GetTextMapPropagator() propagation.TextMapPropagator {
	return global.TextMapPropagator()
}


func SetTextMapPropagator(propagator propagation.TextMapPropagator) {
	global.SetTextMapPropagator(propagator)
}
