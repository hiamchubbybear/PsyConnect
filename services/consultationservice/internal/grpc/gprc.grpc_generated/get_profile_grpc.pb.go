





package gprc_grpc_generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)




const _ = grpc.SupportPackageIsVersion9

const (
	ProfileService_GetProfile_FullMethodName      = "/profile.ProfileService/GetProfile"
	ProfileService_GetProfileBatch_FullMethodName = "/profile.ProfileService/GetProfileBatch"
)




type ProfileServiceClient interface {
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error)
	GetProfileBatch(ctx context.Context, in *GetProfileBatchRequest, opts ...grpc.CallOption) (*GetProfileBatchResponse, error)
}

type profileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileServiceClient(cc grpc.ClientConnInterface) ProfileServiceClient {
	return &profileServiceClient{cc}
}

func (c *profileServiceClient) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProfileResponse)
	err := c.cc.Invoke(ctx, ProfileService_GetProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetProfileBatch(ctx context.Context, in *GetProfileBatchRequest, opts ...grpc.CallOption) (*GetProfileBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProfileBatchResponse)
	err := c.cc.Invoke(ctx, ProfileService_GetProfileBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}




type ProfileServiceServer interface {
	GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
	GetProfileBatch(context.Context, *GetProfileBatchRequest) (*GetProfileBatchResponse, error)
	mustEmbedUnimplementedProfileServiceServer()
}






type UnimplementedProfileServiceServer struct{}

func (UnimplementedProfileServiceServer) GetProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}
func (UnimplementedProfileServiceServer) GetProfileBatch(context.Context, *GetProfileBatchRequest) (*GetProfileBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileBatch not implemented")
}
func (UnimplementedProfileServiceServer) mustEmbedUnimplementedProfileServiceServer() {}
func (UnimplementedProfileServiceServer) testEmbeddedByValue()                        {}




type UnsafeProfileServiceServer interface {
	mustEmbedUnimplementedProfileServiceServer()
}

func RegisterProfileServiceServer(s grpc.ServiceRegistrar, srv ProfileServiceServer) {
	
	
	
	
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProfileService_ServiceDesc, srv)
}

func _ProfileService_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetProfile(ctx, req.(*GetProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetProfileBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetProfileBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetProfileBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetProfileBatch(ctx, req.(*GetProfileBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}




var ProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile.ProfileService",
	HandlerType: (*ProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProfile",
			Handler:    _ProfileService_GetProfile_Handler,
		},
		{
			MethodName: "GetProfileBatch",
			Handler:    _ProfileService_GetProfileBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "get_profile.proto",
}
