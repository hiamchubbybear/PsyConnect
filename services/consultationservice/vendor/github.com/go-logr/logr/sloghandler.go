//go:build go1.21
// +build go1.21



package logr

import (
	"context"
	"log/slog"
)

type slogHandler struct {
	
	sink LogSink
	
	slogSink SlogSink

	
	
	groupPrefix string

	
	
	
	
	
	levelBias slog.Level
}

var _ slog.Handler = &slogHandler{}


const groupSeparator = "."


func (l *slogHandler) GetLevel() slog.Level {
	return l.levelBias
}

func (l *slogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return l.sink != nil && (level >= slog.LevelError || l.sink.Enabled(l.levelFromSlog(level)))
}

func (l *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	if l.slogSink != nil {
		
		if record.Level < slog.LevelError {
			record.Level -= l.levelBias
		}
		return l.slogSink.Handle(ctx, record)
	}

	
	

	kvList := make([]any, 0, 2*record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		kvList = attrToKVs(attr, l.groupPrefix, kvList)
		return true
	})
	if record.Level >= slog.LevelError {
		l.sinkWithCallDepth().Error(nil, record.Message, kvList...)
	} else {
		level := l.levelFromSlog(record.Level)
		l.sinkWithCallDepth().Info(level, record.Message, kvList...)
	}
	return nil
}












func (l *slogHandler) sinkWithCallDepth() LogSink {
	if sink, ok := l.sink.(CallDepthLogSink); ok {
		return sink.WithCallDepth(2)
	}
	return l.sink
}

func (l *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if l.sink == nil || len(attrs) == 0 {
		return l
	}

	clone := *l
	if l.slogSink != nil {
		clone.slogSink = l.slogSink.WithAttrs(attrs)
		clone.sink = clone.slogSink
	} else {
		kvList := make([]any, 0, 2*len(attrs))
		for _, attr := range attrs {
			kvList = attrToKVs(attr, l.groupPrefix, kvList)
		}
		clone.sink = l.sink.WithValues(kvList...)
	}
	return &clone
}

func (l *slogHandler) WithGroup(name string) slog.Handler {
	if l.sink == nil {
		return l
	}
	if name == "" {
		
		return l
	}
	clone := *l
	if l.slogSink != nil {
		clone.slogSink = l.slogSink.WithGroup(name)
		clone.sink = clone.slogSink
	} else {
		clone.groupPrefix = addPrefix(clone.groupPrefix, name)
	}
	return &clone
}



func attrToKVs(attr slog.Attr, groupPrefix string, kvList []any) []any {
	attrVal := attr.Value.Resolve()
	if attrVal.Kind() == slog.KindGroup {
		groupVal := attrVal.Group()
		grpKVs := make([]any, 0, 2*len(groupVal))
		prefix := groupPrefix
		if attr.Key != "" {
			prefix = addPrefix(groupPrefix, attr.Key)
		}
		for _, attr := range groupVal {
			grpKVs = attrToKVs(attr, prefix, grpKVs)
		}
		kvList = append(kvList, grpKVs...)
	} else if attr.Key != "" {
		kvList = append(kvList, addPrefix(groupPrefix, attr.Key), attrVal.Any())
	}

	return kvList
}

func addPrefix(prefix, name string) string {
	if prefix == "" {
		return name
	}
	if name == "" {
		return prefix
	}
	return prefix + groupSeparator + name
}














func (l *slogHandler) levelFromSlog(level slog.Level) int {
	result := -level
	result += l.levelBias 
	if result < 0 {
		result = 0 
	}
	return int(result)
}
