

package transport

import (
	"fmt"

	"google.golang.org/grpc/grpclog"
	internalgrpclog "google.golang.org/grpc/internal/grpclog"
)

var logger = grpclog.Component("transport")

func prefixLoggerForServerTransport(p *http2Server) *internalgrpclog.PrefixLogger {
	return internalgrpclog.NewPrefixLogger(logger, fmt.Sprintf("[server-transport %p] ", p))
}

func prefixLoggerForServerHandlerTransport(p *serverHandlerTransport) *internalgrpclog.PrefixLogger {
	return internalgrpclog.NewPrefixLogger(logger, fmt.Sprintf("[server-handler-transport %p] ", p))
}

func prefixLoggerForClientTransport(p *http2Client) *internalgrpclog.PrefixLogger {
	return internalgrpclog.NewPrefixLogger(logger, fmt.Sprintf("[client-transport %p] ", p))
}
