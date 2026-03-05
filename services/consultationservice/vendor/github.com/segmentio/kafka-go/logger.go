package kafka


type Logger interface {
	Printf(string, ...interface{})
}








type LoggerFunc func(string, ...interface{})

func (f LoggerFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }
