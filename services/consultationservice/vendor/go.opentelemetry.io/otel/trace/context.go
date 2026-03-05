


package trace 

import "context"

type traceContextKeyType int

const currentSpanKey traceContextKeyType = iota


func ContextWithSpan(parent context.Context, span Span) context.Context {
	return context.WithValue(parent, currentSpanKey, span)
}





func ContextWithSpanContext(parent context.Context, sc SpanContext) context.Context {
	return ContextWithSpan(parent, nonRecordingSpan{sc: sc})
}





func ContextWithRemoteSpanContext(parent context.Context, rsc SpanContext) context.Context {
	return ContextWithSpanContext(parent, rsc.WithRemote(true))
}





func SpanFromContext(ctx context.Context) Span {
	if ctx == nil {
		return noopSpanInstance
	}
	if span, ok := ctx.Value(currentSpanKey).(Span); ok {
		return span
	}
	return noopSpanInstance
}


func SpanContextFromContext(ctx context.Context) SpanContext {
	return SpanFromContext(ctx).SpanContext()
}
