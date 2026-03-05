


package global 

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)






var globalLogger = func() *atomic.Pointer[logr.Logger] {
	l := stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

	p := new(atomic.Pointer[logr.Logger])
	p.Store(&l)
	return p
}()






func SetLogger(l logr.Logger) {
	globalLogger.Store(&l)
}


func GetLogger() logr.Logger {
	return *globalLogger.Load()
}



func Info(msg string, keysAndValues ...interface{}) {
	GetLogger().V(4).Info(msg, keysAndValues...)
}


func Error(err error, msg string, keysAndValues ...interface{}) {
	GetLogger().Error(err, msg, keysAndValues...)
}


func Debug(msg string, keysAndValues ...interface{}) {
	GetLogger().V(8).Info(msg, keysAndValues...)
}



func Warn(msg string, keysAndValues ...interface{}) {
	GetLogger().V(1).Info(msg, keysAndValues...)
}
