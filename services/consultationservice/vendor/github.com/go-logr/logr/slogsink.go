//go:build go1.21
// +build go1.21



package logr

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

var (
	_ LogSink          = &slogSink{}
	_ CallDepthLogSink = &slogSink{}
	_ Underlier        = &slogSink{}
)


type Underlier interface {
	
	GetUnderlying() slog.Handler
}

const (
	
	nameKey = "logger"

	
	errKey = "err"
)

type slogSink struct {
	callDepth int
	name      string
	handler   slog.Handler
}

func (l *slogSink) Init(info RuntimeInfo) {
	l.callDepth = info.CallDepth
}

func (l *slogSink) GetUnderlying() slog.Handler {
	return l.handler
}

func (l *slogSink) WithCallDepth(depth int) LogSink {
	newLogger := *l
	newLogger.callDepth += depth
	return &newLogger
}

func (l *slogSink) Enabled(level int) bool {
	return l.handler.Enabled(context.Background(), slog.Level(-level))
}

func (l *slogSink) Info(level int, msg string, kvList ...interface{}) {
	l.log(nil, msg, slog.Level(-level), kvList...)
}

func (l *slogSink) Error(err error, msg string, kvList ...interface{}) {
	l.log(err, msg, slog.LevelError, kvList...)
}

func (l *slogSink) log(err error, msg string, level slog.Level, kvList ...interface{}) {
	var pcs [1]uintptr
	
	runtime.Callers(3+l.callDepth, pcs[:])

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	if l.name != "" {
		record.AddAttrs(slog.String(nameKey, l.name))
	}
	if err != nil {
		record.AddAttrs(slog.Any(errKey, err))
	}
	record.Add(kvList...)
	_ = l.handler.Handle(context.Background(), record)
}

func (l slogSink) WithName(name string) LogSink {
	if l.name != "" {
		l.name += "/"
	}
	l.name += name
	return &l
}

func (l slogSink) WithValues(kvList ...interface{}) LogSink {
	l.handler = l.handler.WithAttrs(kvListToAttrs(kvList...))
	return &l
}

func kvListToAttrs(kvList ...interface{}) []slog.Attr {
	
	record := slog.NewRecord(time.Time{}, 0, "", 0)
	record.Add(kvList...)
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	return attrs
}
