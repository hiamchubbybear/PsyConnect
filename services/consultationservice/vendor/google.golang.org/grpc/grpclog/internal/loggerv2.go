

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)


type LoggerV2 interface {
	
	Info(args ...any)
	
	Infoln(args ...any)
	
	Infof(format string, args ...any)
	
	Warning(args ...any)
	
	Warningln(args ...any)
	
	Warningf(format string, args ...any)
	
	Error(args ...any)
	
	Errorln(args ...any)
	
	Errorf(format string, args ...any)
	
	
	
	Fatal(args ...any)
	
	
	
	Fatalln(args ...any)
	
	
	
	Fatalf(format string, args ...any)
	
	V(l int) bool
}









type DepthLoggerV2 interface {
	LoggerV2
	
	InfoDepth(depth int, args ...any)
	
	WarningDepth(depth int, args ...any)
	
	ErrorDepth(depth int, args ...any)
	
	FatalDepth(depth int, args ...any)
}

const (
	
	infoLog int = iota
	
	warningLog
	
	errorLog
	
	fatalLog
)


var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}



var sprintf = fmt.Sprintf



var sprint = fmt.Sprint



var sprintln = fmt.Sprintln



var exit = os.Exit


type loggerT struct {
	m          []*log.Logger
	v          int
	jsonFormat bool
}

func (g *loggerT) output(severity int, s string) {
	sevStr := severityName[severity]
	if !g.jsonFormat {
		g.m[severity].Output(2, sevStr+": "+s)
		return
	}
	
	
	b, _ := json.Marshal(map[string]string{
		"severity": sevStr,
		"message":  s,
	})
	g.m[severity].Output(2, string(b))
}

func (g *loggerT) printf(severity int, format string, args ...any) {
	
	
	
	if lg := g.m[severity]; lg.Writer() == io.Discard {
		return
	}
	g.output(severity, sprintf(format, args...))
}

func (g *loggerT) print(severity int, v ...any) {
	if lg := g.m[severity]; lg.Writer() == io.Discard {
		return
	}
	g.output(severity, sprint(v...))
}

func (g *loggerT) println(severity int, v ...any) {
	if lg := g.m[severity]; lg.Writer() == io.Discard {
		return
	}
	g.output(severity, sprintln(v...))
}

func (g *loggerT) Info(args ...any) {
	g.print(infoLog, args...)
}

func (g *loggerT) Infoln(args ...any) {
	g.println(infoLog, args...)
}

func (g *loggerT) Infof(format string, args ...any) {
	g.printf(infoLog, format, args...)
}

func (g *loggerT) Warning(args ...any) {
	g.print(warningLog, args...)
}

func (g *loggerT) Warningln(args ...any) {
	g.println(warningLog, args...)
}

func (g *loggerT) Warningf(format string, args ...any) {
	g.printf(warningLog, format, args...)
}

func (g *loggerT) Error(args ...any) {
	g.print(errorLog, args...)
}

func (g *loggerT) Errorln(args ...any) {
	g.println(errorLog, args...)
}

func (g *loggerT) Errorf(format string, args ...any) {
	g.printf(errorLog, format, args...)
}

func (g *loggerT) Fatal(args ...any) {
	g.print(fatalLog, args...)
	exit(1)
}

func (g *loggerT) Fatalln(args ...any) {
	g.println(fatalLog, args...)
	exit(1)
}

func (g *loggerT) Fatalf(format string, args ...any) {
	g.printf(fatalLog, format, args...)
	exit(1)
}

func (g *loggerT) V(l int) bool {
	return l <= g.v
}


type LoggerV2Config struct {
	
	Verbosity int
	
	FormatJSON bool
}








func combineLoggers(lower, higher io.Writer) io.Writer {
	if lower == io.Discard {
		return higher
	}
	if higher == io.Discard {
		return lower
	}
	return io.MultiWriter(lower, higher)
}




func NewLoggerV2(infoW, warningW, errorW io.Writer, c LoggerV2Config) LoggerV2 {
	flag := log.LstdFlags
	if c.FormatJSON {
		flag = 0
	}

	warningW = combineLoggers(infoW, warningW)
	errorW = combineLoggers(errorW, warningW)

	fatalW := errorW

	m := []*log.Logger{
		log.New(infoW, "", flag),
		log.New(warningW, "", flag),
		log.New(errorW, "", flag),
		log.New(fatalW, "", flag),
	}
	return &loggerT{m: m, v: c.Verbosity, jsonFormat: c.FormatJSON}
}
