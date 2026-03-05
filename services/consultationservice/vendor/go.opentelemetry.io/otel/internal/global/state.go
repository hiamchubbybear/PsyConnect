


package global 

import (
	"errors"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type (
	errorHandlerHolder struct {
		eh ErrorHandler
	}

	tracerProviderHolder struct {
		tp trace.TracerProvider
	}

	propagatorsHolder struct {
		tm propagation.TextMapPropagator
	}

	meterProviderHolder struct {
		mp metric.MeterProvider
	}
)

var (
	globalErrorHandler  = defaultErrorHandler()
	globalTracer        = defaultTracerValue()
	globalPropagators   = defaultPropagatorsValue()
	globalMeterProvider = defaultMeterProvider()

	delegateErrorHandlerOnce      sync.Once
	delegateTraceOnce             sync.Once
	delegateTextMapPropagatorOnce sync.Once
	delegateMeterOnce             sync.Once
)










func GetErrorHandler() ErrorHandler {
	return globalErrorHandler.Load().(errorHandlerHolder).eh
}







func SetErrorHandler(h ErrorHandler) {
	current := GetErrorHandler()

	if _, cOk := current.(*ErrDelegator); cOk {
		if _, ehOk := h.(*ErrDelegator); ehOk && current == h {
			
			
			Error(
				errors.New("no ErrorHandler delegate configured"),
				"ErrorHandler remains its current value.",
			)
			return
		}
	}

	delegateErrorHandlerOnce.Do(func() {
		if def, ok := current.(*ErrDelegator); ok {
			def.setDelegate(h)
		}
	})
	globalErrorHandler.Store(errorHandlerHolder{eh: h})
}


func TracerProvider() trace.TracerProvider {
	return globalTracer.Load().(tracerProviderHolder).tp
}


func SetTracerProvider(tp trace.TracerProvider) {
	current := TracerProvider()

	if _, cOk := current.(*tracerProvider); cOk {
		if _, tpOk := tp.(*tracerProvider); tpOk && current == tp {
			
			
			Error(
				errors.New("no delegate configured in tracer provider"),
				"Setting tracer provider to its current value. No delegate will be configured",
			)
			return
		}
	}

	delegateTraceOnce.Do(func() {
		if def, ok := current.(*tracerProvider); ok {
			def.setDelegate(tp)
		}
	})
	globalTracer.Store(tracerProviderHolder{tp: tp})
}


func TextMapPropagator() propagation.TextMapPropagator {
	return globalPropagators.Load().(propagatorsHolder).tm
}


func SetTextMapPropagator(p propagation.TextMapPropagator) {
	current := TextMapPropagator()

	if _, cOk := current.(*textMapPropagator); cOk {
		if _, pOk := p.(*textMapPropagator); pOk && current == p {
			
			
			Error(
				errors.New("no delegate configured in text map propagator"),
				"Setting text map propagator to its current value. No delegate will be configured",
			)
			return
		}
	}

	
	
	delegateTextMapPropagatorOnce.Do(func() {
		if def, ok := current.(*textMapPropagator); ok {
			def.SetDelegate(p)
		}
	})
	
	globalPropagators.Store(propagatorsHolder{tm: p})
}


func MeterProvider() metric.MeterProvider {
	return globalMeterProvider.Load().(meterProviderHolder).mp
}


func SetMeterProvider(mp metric.MeterProvider) {
	current := MeterProvider()
	if _, cOk := current.(*meterProvider); cOk {
		if _, mpOk := mp.(*meterProvider); mpOk && current == mp {
			
			
			Error(
				errors.New("no delegate configured in meter provider"),
				"Setting meter provider to its current value. No delegate will be configured",
			)
			return
		}
	}

	delegateMeterOnce.Do(func() {
		if def, ok := current.(*meterProvider); ok {
			def.setDelegate(mp)
		}
	})
	globalMeterProvider.Store(meterProviderHolder{mp: mp})
}

func defaultErrorHandler() *atomic.Value {
	v := &atomic.Value{}
	v.Store(errorHandlerHolder{eh: &ErrDelegator{}})
	return v
}

func defaultTracerValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(tracerProviderHolder{tp: &tracerProvider{}})
	return v
}

func defaultPropagatorsValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(propagatorsHolder{tm: newTextMapPropagator()})
	return v
}

func defaultMeterProvider() *atomic.Value {
	v := &atomic.Value{}
	v.Store(meterProviderHolder{mp: &meterProvider{}})
	return v
}
