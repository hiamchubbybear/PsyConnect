package dev.psyconnect.grpc.profile;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: get_profile.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class ProfileServiceGrpc {

  private ProfileServiceGrpc() {}

  public static final String SERVICE_NAME = "profile.ProfileService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.profile.GetProfileRequest,
      dev.psyconnect.grpc.profile.GetProfileResponse> getGetProfileMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetProfile",
      requestType = dev.psyconnect.grpc.profile.GetProfileRequest.class,
      responseType = dev.psyconnect.grpc.profile.GetProfileResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.profile.GetProfileRequest,
      dev.psyconnect.grpc.profile.GetProfileResponse> getGetProfileMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.profile.GetProfileRequest, dev.psyconnect.grpc.profile.GetProfileResponse> getGetProfileMethod;
    if ((getGetProfileMethod = ProfileServiceGrpc.getGetProfileMethod) == null) {
      synchronized (ProfileServiceGrpc.class) {
        if ((getGetProfileMethod = ProfileServiceGrpc.getGetProfileMethod) == null) {
          ProfileServiceGrpc.getGetProfileMethod = getGetProfileMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.profile.GetProfileRequest, dev.psyconnect.grpc.profile.GetProfileResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetProfile"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.profile.GetProfileRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.profile.GetProfileResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ProfileServiceMethodDescriptorSupplier("GetProfile"))
              .build();
        }
      }
    }
    return getGetProfileMethod;
  }

  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.profile.GetProfileBatchRequest,
      dev.psyconnect.grpc.profile.GetProfileBatchResponse> getGetProfileBatchMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetProfileBatch",
      requestType = dev.psyconnect.grpc.profile.GetProfileBatchRequest.class,
      responseType = dev.psyconnect.grpc.profile.GetProfileBatchResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.profile.GetProfileBatchRequest,
      dev.psyconnect.grpc.profile.GetProfileBatchResponse> getGetProfileBatchMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.profile.GetProfileBatchRequest, dev.psyconnect.grpc.profile.GetProfileBatchResponse> getGetProfileBatchMethod;
    if ((getGetProfileBatchMethod = ProfileServiceGrpc.getGetProfileBatchMethod) == null) {
      synchronized (ProfileServiceGrpc.class) {
        if ((getGetProfileBatchMethod = ProfileServiceGrpc.getGetProfileBatchMethod) == null) {
          ProfileServiceGrpc.getGetProfileBatchMethod = getGetProfileBatchMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.profile.GetProfileBatchRequest, dev.psyconnect.grpc.profile.GetProfileBatchResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetProfileBatch"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.profile.GetProfileBatchRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.profile.GetProfileBatchResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ProfileServiceMethodDescriptorSupplier("GetProfileBatch"))
              .build();
        }
      }
    }
    return getGetProfileBatchMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ProfileServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ProfileServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ProfileServiceStub>() {
        @java.lang.Override
        public ProfileServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ProfileServiceStub(channel, callOptions);
        }
      };
    return ProfileServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ProfileServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ProfileServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ProfileServiceBlockingStub>() {
        @java.lang.Override
        public ProfileServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ProfileServiceBlockingStub(channel, callOptions);
        }
      };
    return ProfileServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ProfileServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ProfileServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ProfileServiceFutureStub>() {
        @java.lang.Override
        public ProfileServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ProfileServiceFutureStub(channel, callOptions);
        }
      };
    return ProfileServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class ProfileServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getProfile(dev.psyconnect.grpc.profile.GetProfileRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.profile.GetProfileResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetProfileMethod(), responseObserver);
    }

    /**
     */
    public void getProfileBatch(dev.psyconnect.grpc.profile.GetProfileBatchRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.profile.GetProfileBatchResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetProfileBatchMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetProfileMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.profile.GetProfileRequest,
                dev.psyconnect.grpc.profile.GetProfileResponse>(
                  this, METHODID_GET_PROFILE)))
          .addMethod(
            getGetProfileBatchMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.profile.GetProfileBatchRequest,
                dev.psyconnect.grpc.profile.GetProfileBatchResponse>(
                  this, METHODID_GET_PROFILE_BATCH)))
          .build();
    }
  }

  /**
   */
  public static final class ProfileServiceStub extends io.grpc.stub.AbstractAsyncStub<ProfileServiceStub> {
    private ProfileServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ProfileServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ProfileServiceStub(channel, callOptions);
    }

    /**
     */
    public void getProfile(dev.psyconnect.grpc.profile.GetProfileRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.profile.GetProfileResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetProfileMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getProfileBatch(dev.psyconnect.grpc.profile.GetProfileBatchRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.profile.GetProfileBatchResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetProfileBatchMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class ProfileServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<ProfileServiceBlockingStub> {
    private ProfileServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ProfileServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ProfileServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public dev.psyconnect.grpc.profile.GetProfileResponse getProfile(dev.psyconnect.grpc.profile.GetProfileRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetProfileMethod(), getCallOptions(), request);
    }

    /**
     */
    public dev.psyconnect.grpc.profile.GetProfileBatchResponse getProfileBatch(dev.psyconnect.grpc.profile.GetProfileBatchRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetProfileBatchMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class ProfileServiceFutureStub extends io.grpc.stub.AbstractFutureStub<ProfileServiceFutureStub> {
    private ProfileServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ProfileServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ProfileServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.profile.GetProfileResponse> getProfile(
        dev.psyconnect.grpc.profile.GetProfileRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetProfileMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.profile.GetProfileBatchResponse> getProfileBatch(
        dev.psyconnect.grpc.profile.GetProfileBatchRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetProfileBatchMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_PROFILE = 0;
  private static final int METHODID_GET_PROFILE_BATCH = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ProfileServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ProfileServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_PROFILE:
          serviceImpl.getProfile((dev.psyconnect.grpc.profile.GetProfileRequest) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.profile.GetProfileResponse>) responseObserver);
          break;
        case METHODID_GET_PROFILE_BATCH:
          serviceImpl.getProfileBatch((dev.psyconnect.grpc.profile.GetProfileBatchRequest) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.profile.GetProfileBatchResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class ProfileServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ProfileServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return dev.psyconnect.grpc.profile.GetProfile.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("ProfileService");
    }
  }

  private static final class ProfileServiceFileDescriptorSupplier
      extends ProfileServiceBaseDescriptorSupplier {
    ProfileServiceFileDescriptorSupplier() {}
  }

  private static final class ProfileServiceMethodDescriptorSupplier
      extends ProfileServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ProfileServiceMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (ProfileServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ProfileServiceFileDescriptorSupplier())
              .addMethod(getGetProfileMethod())
              .addMethod(getGetProfileBatchMethod())
              .build();
        }
      }
    }
    return result;
  }
}
