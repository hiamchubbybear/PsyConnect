package dev.psyconnect.grpc;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: profile_creation_request.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class ProfileCreationServiceGrpc {

  private ProfileCreationServiceGrpc() {}

  public static final String SERVICE_NAME = "dev.psyconnect.grpc.ProfileCreationService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.ProfileCreationRequest,
      dev.psyconnect.grpc.ProfileCreationResponse> getCreateUserMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "createUser",
      requestType = dev.psyconnect.grpc.ProfileCreationRequest.class,
      responseType = dev.psyconnect.grpc.ProfileCreationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.ProfileCreationRequest,
      dev.psyconnect.grpc.ProfileCreationResponse> getCreateUserMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.ProfileCreationRequest, dev.psyconnect.grpc.ProfileCreationResponse> getCreateUserMethod;
    if ((getCreateUserMethod = ProfileCreationServiceGrpc.getCreateUserMethod) == null) {
      synchronized (ProfileCreationServiceGrpc.class) {
        if ((getCreateUserMethod = ProfileCreationServiceGrpc.getCreateUserMethod) == null) {
          ProfileCreationServiceGrpc.getCreateUserMethod = getCreateUserMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.ProfileCreationRequest, dev.psyconnect.grpc.ProfileCreationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "createUser"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.ProfileCreationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.ProfileCreationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ProfileCreationServiceMethodDescriptorSupplier("createUser"))
              .build();
        }
      }
    }
    return getCreateUserMethod;
  }

  private static volatile io.grpc.MethodDescriptor<dev.psyconnect.grpc.Hello,
      dev.psyconnect.grpc.HelloResponse> getHellowordMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "helloword",
      requestType = dev.psyconnect.grpc.Hello.class,
      responseType = dev.psyconnect.grpc.HelloResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<dev.psyconnect.grpc.Hello,
      dev.psyconnect.grpc.HelloResponse> getHellowordMethod() {
    io.grpc.MethodDescriptor<dev.psyconnect.grpc.Hello, dev.psyconnect.grpc.HelloResponse> getHellowordMethod;
    if ((getHellowordMethod = ProfileCreationServiceGrpc.getHellowordMethod) == null) {
      synchronized (ProfileCreationServiceGrpc.class) {
        if ((getHellowordMethod = ProfileCreationServiceGrpc.getHellowordMethod) == null) {
          ProfileCreationServiceGrpc.getHellowordMethod = getHellowordMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.Hello, dev.psyconnect.grpc.HelloResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "helloword"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.Hello.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.HelloResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ProfileCreationServiceMethodDescriptorSupplier("helloword"))
              .build();
        }
      }
    }
    return getHellowordMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ProfileCreationServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ProfileCreationServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ProfileCreationServiceStub>() {
        @java.lang.Override
        public ProfileCreationServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ProfileCreationServiceStub(channel, callOptions);
        }
      };
    return ProfileCreationServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ProfileCreationServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ProfileCreationServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ProfileCreationServiceBlockingStub>() {
        @java.lang.Override
        public ProfileCreationServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ProfileCreationServiceBlockingStub(channel, callOptions);
        }
      };
    return ProfileCreationServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ProfileCreationServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ProfileCreationServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ProfileCreationServiceFutureStub>() {
        @java.lang.Override
        public ProfileCreationServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ProfileCreationServiceFutureStub(channel, callOptions);
        }
      };
    return ProfileCreationServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class ProfileCreationServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void createUser(dev.psyconnect.grpc.ProfileCreationRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.ProfileCreationResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateUserMethod(), responseObserver);
    }

    /**
     */
    public void helloword(dev.psyconnect.grpc.Hello request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.HelloResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getHellowordMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCreateUserMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.ProfileCreationRequest,
                dev.psyconnect.grpc.ProfileCreationResponse>(
                  this, METHODID_CREATE_USER)))
          .addMethod(
            getHellowordMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                dev.psyconnect.grpc.Hello,
                dev.psyconnect.grpc.HelloResponse>(
                  this, METHODID_HELLOWORD)))
          .build();
    }
  }

  /**
   */
  public static final class ProfileCreationServiceStub extends io.grpc.stub.AbstractAsyncStub<ProfileCreationServiceStub> {
    private ProfileCreationServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ProfileCreationServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ProfileCreationServiceStub(channel, callOptions);
    }

    /**
     */
    public void createUser(dev.psyconnect.grpc.ProfileCreationRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.ProfileCreationResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCreateUserMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void helloword(dev.psyconnect.grpc.Hello request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.HelloResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getHellowordMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class ProfileCreationServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<ProfileCreationServiceBlockingStub> {
    private ProfileCreationServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ProfileCreationServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ProfileCreationServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public dev.psyconnect.grpc.ProfileCreationResponse createUser(dev.psyconnect.grpc.ProfileCreationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateUserMethod(), getCallOptions(), request);
    }

    /**
     */
    public dev.psyconnect.grpc.HelloResponse helloword(dev.psyconnect.grpc.Hello request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getHellowordMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class ProfileCreationServiceFutureStub extends io.grpc.stub.AbstractFutureStub<ProfileCreationServiceFutureStub> {
    private ProfileCreationServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ProfileCreationServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ProfileCreationServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.ProfileCreationResponse> createUser(
        dev.psyconnect.grpc.ProfileCreationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCreateUserMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.HelloResponse> helloword(
        dev.psyconnect.grpc.Hello request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getHellowordMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE_USER = 0;
  private static final int METHODID_HELLOWORD = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ProfileCreationServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ProfileCreationServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE_USER:
          serviceImpl.createUser((dev.psyconnect.grpc.ProfileCreationRequest) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.ProfileCreationResponse>) responseObserver);
          break;
        case METHODID_HELLOWORD:
          serviceImpl.helloword((dev.psyconnect.grpc.Hello) request,
              (io.grpc.stub.StreamObserver<dev.psyconnect.grpc.HelloResponse>) responseObserver);
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

  private static abstract class ProfileCreationServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ProfileCreationServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return dev.psyconnect.grpc.ProfileCreationRequestOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("ProfileCreationService");
    }
  }

  private static final class ProfileCreationServiceFileDescriptorSupplier
      extends ProfileCreationServiceBaseDescriptorSupplier {
    ProfileCreationServiceFileDescriptorSupplier() {}
  }

  private static final class ProfileCreationServiceMethodDescriptorSupplier
      extends ProfileCreationServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ProfileCreationServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (ProfileCreationServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ProfileCreationServiceFileDescriptorSupplier())
              .addMethod(getCreateUserMethod())
              .addMethod(getHellowordMethod())
              .build();
        }
      }
    }
    return result;
  }
}
