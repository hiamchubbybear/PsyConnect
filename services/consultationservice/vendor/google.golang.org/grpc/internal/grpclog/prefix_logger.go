



package grpclog

import (
	"fmt"

	"google.golang.org/grpc/grpclog"
)




type PrefixLogger struct {
	logger grpclog.DepthLoggerV2
	prefix string
}


func (pl *PrefixLogger) Infof(format string, args ...any) {
	if pl != nil {
		
		format = pl.prefix + format
		pl.logger.InfoDepth(1, fmt.Sprintf(format, args...))
		return
	}
	grpclog.InfoDepth(1, fmt.Sprintf(format, args...))
}


func (pl *PrefixLogger) Warningf(format string, args ...any) {
	if pl != nil {
		format = pl.prefix + format
		pl.logger.WarningDepth(1, fmt.Sprintf(format, args...))
		return
	}
	grpclog.WarningDepth(1, fmt.Sprintf(format, args...))
}


func (pl *PrefixLogger) Errorf(format string, args ...any) {
	if pl != nil {
		format = pl.prefix + format
		pl.logger.ErrorDepth(1, fmt.Sprintf(format, args...))
		return
	}
	grpclog.ErrorDepth(1, fmt.Sprintf(format, args...))
}


func (pl *PrefixLogger) V(l int) bool {
	if pl != nil {
		return pl.logger.V(l)
	}
	return true
}


func NewPrefixLogger(logger grpclog.DepthLoggerV2, prefix string) *PrefixLogger {
	return &PrefixLogger{logger: logger, prefix: prefix}
}
