package dev.psyconnect.grpc.consultation_profile;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: check_existed_profile.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class CheckProfileServiceGrpc {

  private CheckProfileServiceGrpc() {}

  public static final String SERVICE_NAME = "profile.CheckProfileService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest,
      dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse> getCheckProfileExistsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CheckProfileExists",
      requestType = dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest.class,
      responseType = dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest,
      dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse> getCheckProfileExistsMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest, dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse> getCheckProfileExistsMethod;
    if ((getCheckProfileExistsMethod = CheckProfileServiceGrpc.getCheckProfileExistsMethod) == null) {
      synchronized (CheckProfileServiceGrpc.class) {
        if ((getCheckProfileExistsMethod = CheckProfileServiceGrpc.getCheckProfileExistsMethod) == null) {
          CheckProfileServiceGrpc.getCheckProfileExistsMethod = getCheckProfileExistsMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest, dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CheckProfileExists"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse.getDefaultInstance()))
              .setSchemaDescriptor(new CheckProfileServiceMethodDescriptorSupplier("CheckProfileExists"))
              .build();
        }
      }
    }
    return getCheckProfileExistsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1,
      dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1> getCheckProfileExistsV1Method;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CheckProfileExistsV1",
      requestType = dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1.class,
      responseType = dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1,
      dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1> getCheckProfileExistsV1Method() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1, dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1> getCheckProfileExistsV1Method;
    if ((getCheckProfileExistsV1Method = CheckProfileServiceGrpc.getCheckProfileExistsV1Method) == null) {
      synchronized (CheckProfileServiceGrpc.class) {
        if ((getCheckProfileExistsV1Method = CheckProfileServiceGrpc.getCheckProfileExistsV1Method) == null) {
          CheckProfileServiceGrpc.getCheckProfileExistsV1Method = getCheckProfileExistsV1Method =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1, dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CheckProfileExistsV1"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1.getDefaultInstance()))
              .setSchemaDescriptor(new CheckProfileServiceMethodDescriptorSupplier("CheckProfileExistsV1"))
              .build();
        }
      }
    }
    return getCheckProfileExistsV1Method;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static CheckProfileServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CheckProfileServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CheckProfileServiceStub>() {
        @java.lang.Override
        public CheckProfileServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CheckProfileServiceStub(channel, callOptions);
        }
      };
    return CheckProfileServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static CheckProfileServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CheckProfileServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CheckProfileServiceBlockingStub>() {
        @java.lang.Override
        public CheckProfileServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CheckProfileServiceBlockingStub(channel, callOptions);
        }
      };
    return CheckProfileServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static CheckProfileServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CheckProfileServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CheckProfileServiceFutureStub>() {
        @java.lang.Override
        public CheckProfileServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CheckProfileServiceFutureStub(channel, callOptions);
        }
      };
    return CheckProfileServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class CheckProfileServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void checkProfileExists(dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCheckProfileExistsMethod(), responseObserver);
    }

    /**
     */
    public void checkProfileExistsV1(dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1 request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCheckProfileExistsV1Method(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCheckProfileExistsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest,
                dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse>(
                  this, METHODID_CHECK_PROFILE_EXISTS)))
          .addMethod(
            getCheckProfileExistsV1Method(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1,
                dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1>(
                  this, METHODID_CHECK_PROFILE_EXISTS_V1)))
          .build();
    }
  }

  /**
   */
  public static final class CheckProfileServiceStub extends io.grpc.stub.AbstractAsyncStub<CheckProfileServiceStub> {
    private CheckProfileServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CheckProfileServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CheckProfileServiceStub(channel, callOptions);
    }

    /**
     */
    public void checkProfileExists(dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCheckProfileExistsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void checkProfileExistsV1(dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1 request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCheckProfileExistsV1Method(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class CheckProfileServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<CheckProfileServiceBlockingStub> {
    private CheckProfileServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CheckProfileServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CheckProfileServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse checkProfileExists(dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCheckProfileExistsMethod(), getCallOptions(), request);
    }

    /**
     */
    public dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1 checkProfileExistsV1(dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1 request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCheckProfileExistsV1Method(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class CheckProfileServiceFutureStub extends io.grpc.stub.AbstractFutureStub<CheckProfileServiceFutureStub> {
    private CheckProfileServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CheckProfileServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CheckProfileServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse> checkProfileExists(
        dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCheckProfileExistsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1> checkProfileExistsV1(
        dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1 request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCheckProfileExistsV1Method(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CHECK_PROFILE_EXISTS = 0;
  private static final int METHODID_CHECK_PROFILE_EXISTS_V1 = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final CheckProfileServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(CheckProfileServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CHECK_PROFILE_EXISTS:
          serviceImpl.checkProfileExists((dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequest) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponse>) responseObserver);
          break;
        case METHODID_CHECK_PROFILE_EXISTS_V1:
          serviceImpl.checkProfileExistsV1((dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileRequestV1) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.ProfileResponseV1>) responseObserver);
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

  private static abstract class CheckProfileServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    CheckProfileServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return dev.psyconnect.grpc.consultation_profile.CheckExistedProfile.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("CheckProfileService");
    }
  }

  private static final class CheckProfileServiceFileDescriptorSupplier
      extends CheckProfileServiceBaseDescriptorSupplier {
    CheckProfileServiceFileDescriptorSupplier() {}
  }

  private static final class CheckProfileServiceMethodDescriptorSupplier
      extends CheckProfileServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    CheckProfileServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (CheckProfileServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new CheckProfileServiceFileDescriptorSupplier())
              .addMethod(getCheckProfileExistsMethod())
              .addMethod(getCheckProfileExistsV1Method())
              .build();
        }
      }
    }
    return result;
  }
}
