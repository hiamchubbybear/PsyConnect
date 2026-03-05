

package internal




type Logger interface {
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Fatalln(args ...any)
	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)
}


type LoggerWrapper struct {
	Logger
}


func (l *LoggerWrapper) Info(args ...any) {
	l.Logger.Print(args...)
}


func (l *LoggerWrapper) Infoln(args ...any) {
	l.Logger.Println(args...)
}


func (l *LoggerWrapper) Infof(format string, args ...any) {
	l.Logger.Printf(format, args...)
}


func (l *LoggerWrapper) Warning(args ...any) {
	l.Logger.Print(args...)
}


func (l *LoggerWrapper) Warningln(args ...any) {
	l.Logger.Println(args...)
}


func (l *LoggerWrapper) Warningf(format string, args ...any) {
	l.Logger.Printf(format, args...)
}


func (l *LoggerWrapper) Error(args ...any) {
	l.Logger.Print(args...)
}


func (l *LoggerWrapper) Errorln(args ...any) {
	l.Logger.Println(args...)
}


func (l *LoggerWrapper) Errorf(format string, args ...any) {
	l.Logger.Printf(format, args...)
}


func (*LoggerWrapper) V(int) bool {
	
	return true
}
