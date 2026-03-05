//go:build go1.21
// +build go1.21



package logr

import (
	"context"
	"fmt"
	"log/slog"
)


func FromContext(ctx context.Context) (Logger, error) {
	v := ctx.Value(contextKey{})
	if v == nil {
		return Logger{}, notFoundError{}
	}

	switch v := v.(type) {
	case Logger:
		return v, nil
	case *slog.Logger:
		return FromSlogHandler(v.Handler()), nil
	default:
		
		panic(fmt.Sprintf("unexpected value type for logr context key: %T", v))
	}
}


func FromContextAsSlogLogger(ctx context.Context) *slog.Logger {
	v := ctx.Value(contextKey{})
	if v == nil {
		return nil
	}

	switch v := v.(type) {
	case Logger:
		return slog.New(ToSlogHandler(v))
	case *slog.Logger:
		return v
	default:
		
		panic(fmt.Sprintf("unexpected value type for logr context key: %T", v))
	}
}



func FromContextOrDiscard(ctx context.Context) Logger {
	if logger, err := FromContext(ctx); err == nil {
		return logger
	}
	return Discard()
}



func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}



func NewContextWithSlogLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
