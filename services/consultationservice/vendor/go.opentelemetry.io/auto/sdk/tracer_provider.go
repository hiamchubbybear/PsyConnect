


package sdk

import (
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)








func TracerProvider() trace.TracerProvider { return tracerProviderInstance }

var tracerProviderInstance = new(tracerProvider)

type tracerProvider struct{ noop.TracerProvider }

var _ trace.TracerProvider = tracerProvider{}

func (p tracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	cfg := trace.NewTracerConfig(opts...)
	return tracer{
		name:      name,
		version:   cfg.InstrumentationVersion(),
		schemaURL: cfg.SchemaURL(),
	}
}
