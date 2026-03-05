





package gprc_grpc_generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)




const _ = grpc.SupportPackageIsVersion9

const (
	CheckProfileService_CheckProfileExists_FullMethodName   = "/profile.CheckProfileService/CheckProfileExists"
	CheckProfileService_CheckProfileExistsV1_FullMethodName = "/profile.CheckProfileService/CheckProfileExistsV1"
)




type CheckProfileServiceClient interface {
	CheckProfileExists(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error)
	CheckProfileExistsV1(ctx context.Context, in *ProfileRequestV1, opts ...grpc.CallOption) (*ProfileResponseV1, error)
}

type checkProfileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckProfileServiceClient(cc grpc.ClientConnInterface) CheckProfileServiceClient {
	return &checkProfileServiceClient{cc}
}

func (c *checkProfileServiceClient) CheckProfileExists(ctx context.Context, in *ProfileRequest, opts ...grpc.CallOption) (*ProfileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProfileResponse)
	err := c.cc.Invoke(ctx, CheckProfileService_CheckProfileExists_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkProfileServiceClient) CheckProfileExistsV1(ctx context.Context, in *ProfileRequestV1, opts ...grpc.CallOption) (*ProfileResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProfileResponseV1)
	err := c.cc.Invoke(ctx, CheckProfileService_CheckProfileExistsV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}




type CheckProfileServiceServer interface {
	CheckProfileExists(context.Context, *ProfileRequest) (*ProfileResponse, error)
	CheckProfileExistsV1(context.Context, *ProfileRequestV1) (*ProfileResponseV1, error)
	mustEmbedUnimplementedCheckProfileServiceServer()
}






type UnimplementedCheckProfileServiceServer struct{}

func (UnimplementedCheckProfileServiceServer) CheckProfileExists(context.Context, *ProfileRequest) (*ProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckProfileExists not implemented")
}
func (UnimplementedCheckProfileServiceServer) CheckProfileExistsV1(context.Context, *ProfileRequestV1) (*ProfileResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckProfileExistsV1 not implemented")
}
func (UnimplementedCheckProfileServiceServer) mustEmbedUnimplementedCheckProfileServiceServer() {}
func (UnimplementedCheckProfileServiceServer) testEmbeddedByValue()                             {}




type UnsafeCheckProfileServiceServer interface {
	mustEmbedUnimplementedCheckProfileServiceServer()
}

func RegisterCheckProfileServiceServer(s grpc.ServiceRegistrar, srv CheckProfileServiceServer) {
	
	
	
	
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CheckProfileService_ServiceDesc, srv)
}

func _CheckProfileService_CheckProfileExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckProfileServiceServer).CheckProfileExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckProfileService_CheckProfileExists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckProfileServiceServer).CheckProfileExists(ctx, req.(*ProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckProfileService_CheckProfileExistsV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckProfileServiceServer).CheckProfileExistsV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckProfileService_CheckProfileExistsV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckProfileServiceServer).CheckProfileExistsV1(ctx, req.(*ProfileRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}




var CheckProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile.CheckProfileService",
	HandlerType: (*CheckProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckProfileExists",
			Handler:    _CheckProfileService_CheckProfileExists_Handler,
		},
		{
			MethodName: "CheckProfileExistsV1",
			Handler:    _CheckProfileService_CheckProfileExistsV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "check_existed_profile.proto",
}
