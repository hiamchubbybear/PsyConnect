

package grpclog

import "google.golang.org/grpc/grpclog/internal"




type Logger internal.Logger





func SetLogger(l Logger) {
	internal.LoggerV2Impl = &internal.LoggerWrapper{Logger: l}
}
