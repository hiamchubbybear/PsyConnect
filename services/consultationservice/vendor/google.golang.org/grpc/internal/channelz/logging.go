

package channelz

import (
	"fmt"

	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.Component("channelz")


func Info(l grpclog.DepthLoggerV2, e Entity, args ...any) {
	AddTraceEvent(l, e, 1, &TraceEvent{
		Desc:     fmt.Sprint(args...),
		Severity: CtInfo,
	})
}


func Infof(l grpclog.DepthLoggerV2, e Entity, format string, args ...any) {
	AddTraceEvent(l, e, 1, &TraceEvent{
		Desc:     fmt.Sprintf(format, args...),
		Severity: CtInfo,
	})
}


func Warning(l grpclog.DepthLoggerV2, e Entity, args ...any) {
	AddTraceEvent(l, e, 1, &TraceEvent{
		Desc:     fmt.Sprint(args...),
		Severity: CtWarning,
	})
}


func Warningf(l grpclog.DepthLoggerV2, e Entity, format string, args ...any) {
	AddTraceEvent(l, e, 1, &TraceEvent{
		Desc:     fmt.Sprintf(format, args...),
		Severity: CtWarning,
	})
}


func Error(l grpclog.DepthLoggerV2, e Entity, args ...any) {
	AddTraceEvent(l, e, 1, &TraceEvent{
		Desc:     fmt.Sprint(args...),
		Severity: CtError,
	})
}


func Errorf(l grpclog.DepthLoggerV2, e Entity, format string, args ...any) {
	AddTraceEvent(l, e, 1, &TraceEvent{
		Desc:     fmt.Sprintf(format, args...),
		Severity: CtError,
	})
}
