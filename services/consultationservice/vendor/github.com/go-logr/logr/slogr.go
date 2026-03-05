//go:build go1.21
// +build go1.21



package logr

import (
	"context"
	"log/slog"
)





func FromSlogHandler(handler slog.Handler) Logger {
	if handler, ok := handler.(*slogHandler); ok {
		if handler.sink == nil {
			return Discard()
		}
		return New(handler.sink).V(int(handler.levelBias))
	}
	return New(&slogSink{handler: handler})
}



















func ToSlogHandler(logger Logger) slog.Handler {
	if sink, ok := logger.GetSink().(*slogSink); ok && logger.GetV() == 0 {
		return sink.handler
	}

	handler := &slogHandler{sink: logger.GetSink(), levelBias: slog.Level(logger.GetV())}
	if slogSink, ok := handler.sink.(SlogSink); ok {
		handler.slogSink = slogSink
	}
	return handler
}
























type SlogSink interface {
	LogSink

	Handle(ctx context.Context, record slog.Record) error
	WithAttrs(attrs []slog.Attr) SlogSink
	WithGroup(name string) SlogSink
}
