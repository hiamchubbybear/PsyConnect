






package grpclog

import (
	"os"

	"google.golang.org/grpc/grpclog/internal"
)

func init() {
	SetLoggerV2(newLoggerV2())
}


func V(l int) bool {
	return internal.LoggerV2Impl.V(l)
}


func Info(args ...any) {
	internal.LoggerV2Impl.Info(args...)
}


func Infof(format string, args ...any) {
	internal.LoggerV2Impl.Infof(format, args...)
}


func Infoln(args ...any) {
	internal.LoggerV2Impl.Infoln(args...)
}


func Warning(args ...any) {
	internal.LoggerV2Impl.Warning(args...)
}


func Warningf(format string, args ...any) {
	internal.LoggerV2Impl.Warningf(format, args...)
}


func Warningln(args ...any) {
	internal.LoggerV2Impl.Warningln(args...)
}


func Error(args ...any) {
	internal.LoggerV2Impl.Error(args...)
}


func Errorf(format string, args ...any) {
	internal.LoggerV2Impl.Errorf(format, args...)
}


func Errorln(args ...any) {
	internal.LoggerV2Impl.Errorln(args...)
}



func Fatal(args ...any) {
	internal.LoggerV2Impl.Fatal(args...)
	
	os.Exit(1)
}



func Fatalf(format string, args ...any) {
	internal.LoggerV2Impl.Fatalf(format, args...)
	
	os.Exit(1)
}



func Fatalln(args ...any) {
	internal.LoggerV2Impl.Fatalln(args...)
	
	os.Exit(1)
}




func Print(args ...any) {
	internal.LoggerV2Impl.Info(args...)
}




func Printf(format string, args ...any) {
	internal.LoggerV2Impl.Infof(format, args...)
}




func Println(args ...any) {
	internal.LoggerV2Impl.Infoln(args...)
}







func InfoDepth(depth int, args ...any) {
	if internal.DepthLoggerV2Impl != nil {
		internal.DepthLoggerV2Impl.InfoDepth(depth, args...)
	} else {
		internal.LoggerV2Impl.Infoln(args...)
	}
}







func WarningDepth(depth int, args ...any) {
	if internal.DepthLoggerV2Impl != nil {
		internal.DepthLoggerV2Impl.WarningDepth(depth, args...)
	} else {
		internal.LoggerV2Impl.Warningln(args...)
	}
}







func ErrorDepth(depth int, args ...any) {
	if internal.DepthLoggerV2Impl != nil {
		internal.DepthLoggerV2Impl.ErrorDepth(depth, args...)
	} else {
		internal.LoggerV2Impl.Errorln(args...)
	}
}







func FatalDepth(depth int, args ...any) {
	if internal.DepthLoggerV2Impl != nil {
		internal.DepthLoggerV2Impl.FatalDepth(depth, args...)
	} else {
		internal.LoggerV2Impl.Fatalln(args...)
	}
	os.Exit(1)
}
