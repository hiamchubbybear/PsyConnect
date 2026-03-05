//go:build go1.21
// +build go1.21



package funcr

import (
	"context"
	"log/slog"

	"github.com/go-logr/logr"
)

var _ logr.SlogSink = &fnlogger{}

const extraSlogSinkDepth = 3 

func (l fnlogger) Handle(_ context.Context, record slog.Record) error {
	kvList := make([]any, 0, 2*record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		kvList = attrToKVs(attr, kvList)
		return true
	})

	if record.Level >= slog.LevelError {
		l.WithCallDepth(extraSlogSinkDepth).Error(nil, record.Message, kvList...)
	} else {
		level := l.levelFromSlog(record.Level)
		l.WithCallDepth(extraSlogSinkDepth).Info(level, record.Message, kvList...)
	}
	return nil
}

func (l fnlogger) WithAttrs(attrs []slog.Attr) logr.SlogSink {
	kvList := make([]any, 0, 2*len(attrs))
	for _, attr := range attrs {
		kvList = attrToKVs(attr, kvList)
	}
	l.AddValues(kvList)
	return &l
}

func (l fnlogger) WithGroup(name string) logr.SlogSink {
	l.startGroup(name)
	return &l
}



func attrToKVs(attr slog.Attr, kvList []any) []any {
	attrVal := attr.Value.Resolve()
	if attrVal.Kind() == slog.KindGroup {
		groupVal := attrVal.Group()
		grpKVs := make([]any, 0, 2*len(groupVal))
		for _, attr := range groupVal {
			grpKVs = attrToKVs(attr, grpKVs)
		}
		if attr.Key == "" {
			
			kvList = append(kvList, grpKVs...)
		} else {
			kvList = append(kvList, attr.Key, PseudoStruct(grpKVs))
		}
	} else if attr.Key != "" {
		kvList = append(kvList, attr.Key, attrVal.Any())
	}

	return kvList
}














func (l fnlogger) levelFromSlog(level slog.Level) int {
	result := -level
	if result < 0 {
		result = 0 
	}
	return int(result)
}
