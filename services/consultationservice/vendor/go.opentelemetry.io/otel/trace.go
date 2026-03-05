


package otel 

import (
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/trace"
)





func Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return GetTracerProvider().Tracer(name, opts...)
}











func GetTracerProvider() trace.TracerProvider {
	return global.TracerProvider()
}


func SetTracerProvider(tp trace.TracerProvider) {
	global.SetTracerProvider(tp)
}
