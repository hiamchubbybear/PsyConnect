//go:build !go1.21
// +build !go1.21



package logr

import (
	"context"
)


func FromContext(ctx context.Context) (Logger, error) {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v, nil
	}

	return Logger{}, notFoundError{}
}



func FromContextOrDiscard(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}

	return Discard()
}



func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
