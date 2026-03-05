


package trace 

import "go.opentelemetry.io/otel/trace/embedded"


















type TracerProvider interface {
	
	
	
	embedded.TracerProvider

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Tracer(name string, options ...TracerOption) Tracer
}
