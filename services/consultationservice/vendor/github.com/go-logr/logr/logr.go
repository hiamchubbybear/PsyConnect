
































































































































































































package logr




func New(sink LogSink) Logger {
	logger := Logger{}
	logger.setSink(sink)
	if sink != nil {
		sink.Init(runtimeInfo)
	}
	return logger
}




func (l *Logger) setSink(sink LogSink) {
	l.sink = sink
}


func (l Logger) GetSink() LogSink {
	return l.sink
}


func (l Logger) WithSink(sink LogSink) Logger {
	l.setSink(sink)
	return l
}










type Logger struct {
	sink  LogSink
	level int
}



func (l Logger) Enabled() bool {
	
	
	
	
	
	
	return l.sink != nil && l.sink.Enabled(l.level)
}







func (l Logger) Info(msg string, keysAndValues ...any) {
	if l.sink == nil {
		return
	}
	if l.sink.Enabled(l.level) { 
		if withHelper, ok := l.sink.(CallStackHelperLogSink); ok {
			withHelper.GetCallStackHelper()()
		}
		l.sink.Info(l.level, msg, keysAndValues...)
	}
}











func (l Logger) Error(err error, msg string, keysAndValues ...any) {
	if l.sink == nil {
		return
	}
	if withHelper, ok := l.sink.(CallStackHelperLogSink); ok {
		withHelper.GetCallStackHelper()()
	}
	l.sink.Error(err, msg, keysAndValues...)
}





func (l Logger) V(level int) Logger {
	if l.sink == nil {
		return l
	}
	if level < 0 {
		level = 0
	}
	l.level += level
	return l
}



func (l Logger) GetV() int {
	
	return l.level
}



func (l Logger) WithValues(keysAndValues ...any) Logger {
	if l.sink == nil {
		return l
	}
	l.setSink(l.sink.WithValues(keysAndValues...))
	return l
}






func (l Logger) WithName(name string) Logger {
	if l.sink == nil {
		return l
	}
	l.setSink(l.sink.WithName(name))
	return l
}
















func (l Logger) WithCallDepth(depth int) Logger {
	if l.sink == nil {
		return l
	}
	if withCallDepth, ok := l.sink.(CallDepthLogSink); ok {
		l.setSink(withCallDepth.WithCallDepth(depth))
	}
	return l
}















func (l Logger) WithCallStackHelper() (func(), Logger) {
	if l.sink == nil {
		return func() {}, l
	}
	var helper func()
	if withCallDepth, ok := l.sink.(CallDepthLogSink); ok {
		l.setSink(withCallDepth.WithCallDepth(1))
	}
	if withHelper, ok := l.sink.(CallStackHelperLogSink); ok {
		helper = withHelper.GetCallStackHelper()
	} else {
		helper = func() {}
	}
	return helper, l
}


func (l Logger) IsZero() bool {
	return l.sink == nil
}



type RuntimeInfo struct {
	
	
	
	
	CallDepth int
}


var runtimeInfo = RuntimeInfo{
	CallDepth: 1,
}



type LogSink interface {
	
	
	Init(info RuntimeInfo)

	
	
	
	Enabled(level int) bool

	
	
	
	
	Info(level int, msg string, keysAndValues ...any)

	
	
	Error(err error, msg string, keysAndValues ...any)

	
	
	WithValues(keysAndValues ...any) LogSink

	
	
	WithName(name string) LogSink
}











type CallDepthLogSink interface {
	
	
	
	
	
	
	
	
	
	
	
	WithCallDepth(depth int) LogSink
}


















type CallStackHelperLogSink interface {
	
	
	
	GetCallStackHelper() func()
}





type Marshaler interface {
	
	
	
	
	
	
	
	
	
	
	MarshalLog() any
}
