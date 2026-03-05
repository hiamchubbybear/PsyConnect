//go:build grpcnotrace



package grpc





import (
	"context"
	"fmt"
)

type notrace struct{}

func (notrace) LazyLog(x fmt.Stringer, sensitive bool) {}
func (notrace) LazyPrintf(format string, a ...any)     {}
func (notrace) SetError()                              {}
func (notrace) SetRecycler(f func(any))                {}
func (notrace) SetTraceInfo(traceID, spanID uint64)    {}
func (notrace) SetMaxEvents(m int)                     {}
func (notrace) Finish()                                {}

func newTrace(family, title string) traceLog {
	return notrace{}
}

func newTraceContext(ctx context.Context, tr traceLog) context.Context {
	return ctx
}

func newTraceEventLog(family, title string) traceEventLog {
	return nil
}
