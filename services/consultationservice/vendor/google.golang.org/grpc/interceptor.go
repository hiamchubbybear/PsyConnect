

package grpc

import (
	"context"
)


type UnaryInvoker func(ctx context.Context, method string, req, reply any, cc *ClientConn, opts ...CallOption) error
















type UnaryClientInterceptor func(ctx context.Context, method string, req, reply any, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error


type Streamer func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, opts ...CallOption) (ClientStream, error)
















type StreamClientInterceptor func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)



type UnaryServerInfo struct {
	
	Server any
	
	FullMethod string
}








type UnaryHandler func(ctx context.Context, req any) (any, error)





type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)



type StreamServerInfo struct {
	
	FullMethod string
	
	IsClientStream bool
	
	IsServerStream bool
}





type StreamServerInterceptor func(srv any, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
