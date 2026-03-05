//go:build !grpcnotrace



package grpc

import (
	"context"

	t "golang.org/x/net/trace"
)

func newTrace(family, title string) traceLog {
	return t.New(family, title)
}

func newTraceContext(ctx context.Context, tr traceLog) context.Context {
	return t.NewContext(ctx, tr)
}

func newTraceEventLog(family, title string) traceEventLog {
	return t.NewEventLog(family, title)
}
