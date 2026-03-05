





package gprc_grpc_generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)




const _ = grpc.SupportPackageIsVersion9

const (
	CreateMatchingRequest_CreateMatchingRequest_FullMethodName   = "/profile.CreateMatchingRequest/CreateMatchingRequest"
	CreateMatchingRequest_ResponseRequestMatching_FullMethodName = "/profile.CreateMatchingRequest/ResponseRequestMatching"
)




type CreateMatchingRequestClient interface {
	CreateMatchingRequest(ctx context.Context, in *MatchingRequest, opts ...grpc.CallOption) (*MatchingResponse, error)
	ResponseRequestMatching(ctx context.Context, in *ResponseMatchingRequest, opts ...grpc.CallOption) (*ResponseMatchingResponse, error)
}

type createMatchingRequestClient struct {
	cc grpc.ClientConnInterface
}

func NewCreateMatchingRequestClient(cc grpc.ClientConnInterface) CreateMatchingRequestClient {
	return &createMatchingRequestClient{cc}
}

func (c *createMatchingRequestClient) CreateMatchingRequest(ctx context.Context, in *MatchingRequest, opts ...grpc.CallOption) (*MatchingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MatchingResponse)
	err := c.cc.Invoke(ctx, CreateMatchingRequest_CreateMatchingRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *createMatchingRequestClient) ResponseRequestMatching(ctx context.Context, in *ResponseMatchingRequest, opts ...grpc.CallOption) (*ResponseMatchingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseMatchingResponse)
	err := c.cc.Invoke(ctx, CreateMatchingRequest_ResponseRequestMatching_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}




type CreateMatchingRequestServer interface {
	CreateMatchingRequest(context.Context, *MatchingRequest) (*MatchingResponse, error)
	ResponseRequestMatching(context.Context, *ResponseMatchingRequest) (*ResponseMatchingResponse, error)
	mustEmbedUnimplementedCreateMatchingRequestServer()
}






type UnimplementedCreateMatchingRequestServer struct{}

func (UnimplementedCreateMatchingRequestServer) CreateMatchingRequest(context.Context, *MatchingRequest) (*MatchingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMatchingRequest not implemented")
}
func (UnimplementedCreateMatchingRequestServer) ResponseRequestMatching(context.Context, *ResponseMatchingRequest) (*ResponseMatchingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResponseRequestMatching not implemented")
}
func (UnimplementedCreateMatchingRequestServer) mustEmbedUnimplementedCreateMatchingRequestServer() {}
func (UnimplementedCreateMatchingRequestServer) testEmbeddedByValue()                               {}




type UnsafeCreateMatchingRequestServer interface {
	mustEmbedUnimplementedCreateMatchingRequestServer()
}

func RegisterCreateMatchingRequestServer(s grpc.ServiceRegistrar, srv CreateMatchingRequestServer) {
	
	
	
	
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CreateMatchingRequest_ServiceDesc, srv)
}

func _CreateMatchingRequest_CreateMatchingRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreateMatchingRequestServer).CreateMatchingRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreateMatchingRequest_CreateMatchingRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreateMatchingRequestServer).CreateMatchingRequest(ctx, req.(*MatchingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CreateMatchingRequest_ResponseRequestMatching_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseMatchingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreateMatchingRequestServer).ResponseRequestMatching(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreateMatchingRequest_ResponseRequestMatching_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreateMatchingRequestServer).ResponseRequestMatching(ctx, req.(*ResponseMatchingRequest))
	}
	return interceptor(ctx, in, info, handler)
}




var CreateMatchingRequest_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile.CreateMatchingRequest",
	HandlerType: (*CreateMatchingRequestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMatchingRequest",
			Handler:    _CreateMatchingRequest_CreateMatchingRequest_Handler,
		},
		{
			MethodName: "ResponseRequestMatching",
			Handler:    _CreateMatchingRequest_ResponseRequestMatching_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "matching.proto",
}
