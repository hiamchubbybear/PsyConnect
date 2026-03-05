

package grpclog

import (
	"io"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/grpclog/internal"
)


type LoggerV2 internal.LoggerV2



func SetLoggerV2(l LoggerV2) {
	if _, ok := l.(*componentData); ok {
		panic("cannot use component logger as grpclog logger")
	}
	internal.LoggerV2Impl = l
	internal.DepthLoggerV2Impl, _ = l.(internal.DepthLoggerV2)
}






func NewLoggerV2(infoW, warningW, errorW io.Writer) LoggerV2 {
	return internal.NewLoggerV2(infoW, warningW, errorW, internal.LoggerV2Config{})
}



func NewLoggerV2WithVerbosity(infoW, warningW, errorW io.Writer, v int) LoggerV2 {
	return internal.NewLoggerV2(infoW, warningW, errorW, internal.LoggerV2Config{Verbosity: v})
}



func newLoggerV2() LoggerV2 {
	errorW := io.Discard
	warningW := io.Discard
	infoW := io.Discard

	logLevel := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch logLevel {
	case "", "ERROR", "error": 
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	vLevel := os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	}

	jsonFormat := strings.EqualFold(os.Getenv("GRPC_GO_LOG_FORMATTER"), "json")

	return internal.NewLoggerV2(infoW, warningW, errorW, internal.LoggerV2Config{
		Verbosity:  v,
		FormatJSON: jsonFormat,
	})
}









type DepthLoggerV2 internal.DepthLoggerV2
