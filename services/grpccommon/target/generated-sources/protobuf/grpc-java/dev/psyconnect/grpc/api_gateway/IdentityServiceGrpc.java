package dev.psyconnect.grpc.api_gateway;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: token_check_identity.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class IdentityServiceGrpc {

  private IdentityServiceGrpc() {}

  public static final String SERVICE_NAME = "identity.IdentityService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest,
      dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> getTokenCheckValidMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "TokenCheckValid",
      requestType = dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest.class,
      responseType = dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest,
      dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> getTokenCheckValidMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest, dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> getTokenCheckValidMethod;
    if ((getTokenCheckValidMethod = IdentityServiceGrpc.getTokenCheckValidMethod) == null) {
      synchronized (IdentityServiceGrpc.class) {
        if ((getTokenCheckValidMethod = IdentityServiceGrpc.getTokenCheckValidMethod) == null) {
          IdentityServiceGrpc.getTokenCheckValidMethod = getTokenCheckValidMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest, dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "TokenCheckValid"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse.getDefaultInstance()))
              .setSchemaDescriptor(new IdentityServiceMethodDescriptorSupplier("TokenCheckValid"))
              .build();
        }
      }
    }
    return getTokenCheckValidMethod;
  }

  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest,
      dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse> getGetUserInfoByProfileIdMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetUserInfoByProfileId",
      requestType = dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest.class,
      responseType = dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest,
      dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse> getGetUserInfoByProfileIdMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest, dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse> getGetUserInfoByProfileIdMethod;
    if ((getGetUserInfoByProfileIdMethod = IdentityServiceGrpc.getGetUserInfoByProfileIdMethod) == null) {
      synchronized (IdentityServiceGrpc.class) {
        if ((getGetUserInfoByProfileIdMethod = IdentityServiceGrpc.getGetUserInfoByProfileIdMethod) == null) {
          IdentityServiceGrpc.getGetUserInfoByProfileIdMethod = getGetUserInfoByProfileIdMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest, dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetUserInfoByProfileId"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new IdentityServiceMethodDescriptorSupplier("GetUserInfoByProfileId"))
              .build();
        }
      }
    }
    return getGetUserInfoByProfileIdMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static IdentityServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<IdentityServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<IdentityServiceStub>() {
        @java.lang.Override
        public IdentityServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new IdentityServiceStub(channel, callOptions);
        }
      };
    return IdentityServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static IdentityServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<IdentityServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<IdentityServiceBlockingStub>() {
        @java.lang.Override
        public IdentityServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new IdentityServiceBlockingStub(channel, callOptions);
        }
      };
    return IdentityServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static IdentityServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<IdentityServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<IdentityServiceFutureStub>() {
        @java.lang.Override
        public IdentityServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new IdentityServiceFutureStub(channel, callOptions);
        }
      };
    return IdentityServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class IdentityServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void tokenCheckValid(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getTokenCheckValidMethod(), responseObserver);
    }

    /**
     */
    public void getUserInfoByProfileId(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetUserInfoByProfileIdMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getTokenCheckValidMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest,
                dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse>(
                  this, METHODID_TOKEN_CHECK_VALID)))
          .addMethod(
            getGetUserInfoByProfileIdMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest,
                dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse>(
                  this, METHODID_GET_USER_INFO_BY_PROFILE_ID)))
          .build();
    }
  }

  /**
   */
  public static final class IdentityServiceStub extends io.grpc.stub.AbstractAsyncStub<IdentityServiceStub> {
    private IdentityServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected IdentityServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new IdentityServiceStub(channel, callOptions);
    }

    /**
     */
    public void tokenCheckValid(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getTokenCheckValidMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getUserInfoByProfileId(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetUserInfoByProfileIdMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class IdentityServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<IdentityServiceBlockingStub> {
    private IdentityServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected IdentityServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new IdentityServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse tokenCheckValid(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getTokenCheckValidMethod(), getCallOptions(), request);
    }

    /**
     */
    public dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse getUserInfoByProfileId(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetUserInfoByProfileIdMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class IdentityServiceFutureStub extends io.grpc.stub.AbstractFutureStub<IdentityServiceFutureStub> {
    private IdentityServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected IdentityServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new IdentityServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> tokenCheckValid(
        dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getTokenCheckValidMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse> getUserInfoByProfileId(
        dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetUserInfoByProfileIdMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_TOKEN_CHECK_VALID = 0;
  private static final int METHODID_GET_USER_INFO_BY_PROFILE_ID = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final IdentityServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(IdentityServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_TOKEN_CHECK_VALID:
          serviceImpl.tokenCheckValid((dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse>) responseObserver);
          break;
        case METHODID_GET_USER_INFO_BY_PROFILE_ID:
          serviceImpl.getUserInfoByProfileId((dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoRequest) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.UserInfoResponse>) responseObserver);
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

  private static abstract class IdentityServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    IdentityServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("IdentityService");
    }
  }

  private static final class IdentityServiceFileDescriptorSupplier
      extends IdentityServiceBaseDescriptorSupplier {
    IdentityServiceFileDescriptorSupplier() {}
  }

  private static final class IdentityServiceMethodDescriptorSupplier
      extends IdentityServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    IdentityServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (IdentityServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new IdentityServiceFileDescriptorSupplier())
              .addMethod(getTokenCheckValidMethod())
              .addMethod(getGetUserInfoByProfileIdMethod())
              .build();
        }
      }
    }
    return result;
  }
}
