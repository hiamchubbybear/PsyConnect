package profile;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: matching_creatation.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class CreateMatchingRequestGrpc {

  private CreateMatchingRequestGrpc() {}

  public static final String SERVICE_NAME = "profile.CreateMatchingRequest";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<profile.MatchingCreatation.MatchingRequest,
      profile.MatchingCreatation.MatchingResponse> getCreateMatchingRequestMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreateMatchingRequest",
      requestType = profile.MatchingCreatation.MatchingRequest.class,
      responseType = profile.MatchingCreatation.MatchingResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<profile.MatchingCreatation.MatchingRequest,
      profile.MatchingCreatation.MatchingResponse> getCreateMatchingRequestMethod() {
    io.grpc.MethodDescriptor<profile.MatchingCreatation.MatchingRequest, profile.MatchingCreatation.MatchingResponse> getCreateMatchingRequestMethod;
    if ((getCreateMatchingRequestMethod = CreateMatchingRequestGrpc.getCreateMatchingRequestMethod) == null) {
      synchronized (CreateMatchingRequestGrpc.class) {
        if ((getCreateMatchingRequestMethod = CreateMatchingRequestGrpc.getCreateMatchingRequestMethod) == null) {
          CreateMatchingRequestGrpc.getCreateMatchingRequestMethod = getCreateMatchingRequestMethod =
              io.grpc.MethodDescriptor.<profile.MatchingCreatation.MatchingRequest, profile.MatchingCreatation.MatchingResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CreateMatchingRequest"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  profile.MatchingCreatation.MatchingRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  profile.MatchingCreatation.MatchingResponse.getDefaultInstance()))
              .setSchemaDescriptor(new CreateMatchingRequestMethodDescriptorSupplier("CreateMatchingRequest"))
              .build();
        }
      }
    }
    return getCreateMatchingRequestMethod;
  }

  private static volatile io.grpc.MethodDescriptor<profile.MatchingCreatation.ResponseMatchingRequest,
      profile.MatchingCreatation.ResponseMatchingResponse> getResponseRequestMatchingMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ResponseRequestMatching",
      requestType = profile.MatchingCreatation.ResponseMatchingRequest.class,
      responseType = profile.MatchingCreatation.ResponseMatchingResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<profile.MatchingCreatation.ResponseMatchingRequest,
      profile.MatchingCreatation.ResponseMatchingResponse> getResponseRequestMatchingMethod() {
    io.grpc.MethodDescriptor<profile.MatchingCreatation.ResponseMatchingRequest, profile.MatchingCreatation.ResponseMatchingResponse> getResponseRequestMatchingMethod;
    if ((getResponseRequestMatchingMethod = CreateMatchingRequestGrpc.getResponseRequestMatchingMethod) == null) {
      synchronized (CreateMatchingRequestGrpc.class) {
        if ((getResponseRequestMatchingMethod = CreateMatchingRequestGrpc.getResponseRequestMatchingMethod) == null) {
          CreateMatchingRequestGrpc.getResponseRequestMatchingMethod = getResponseRequestMatchingMethod =
              io.grpc.MethodDescriptor.<profile.MatchingCreatation.ResponseMatchingRequest, profile.MatchingCreatation.ResponseMatchingResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ResponseRequestMatching"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  profile.MatchingCreatation.ResponseMatchingRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  profile.MatchingCreatation.ResponseMatchingResponse.getDefaultInstance()))
              .setSchemaDescriptor(new CreateMatchingRequestMethodDescriptorSupplier("ResponseRequestMatching"))
              .build();
        }
      }
    }
    return getResponseRequestMatchingMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static CreateMatchingRequestStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CreateMatchingRequestStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CreateMatchingRequestStub>() {
        @java.lang.Override
        public CreateMatchingRequestStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CreateMatchingRequestStub(channel, callOptions);
        }
      };
    return CreateMatchingRequestStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static CreateMatchingRequestBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CreateMatchingRequestBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CreateMatchingRequestBlockingStub>() {
        @java.lang.Override
        public CreateMatchingRequestBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CreateMatchingRequestBlockingStub(channel, callOptions);
        }
      };
    return CreateMatchingRequestBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static CreateMatchingRequestFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CreateMatchingRequestFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<CreateMatchingRequestFutureStub>() {
        @java.lang.Override
        public CreateMatchingRequestFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new CreateMatchingRequestFutureStub(channel, callOptions);
        }
      };
    return CreateMatchingRequestFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class CreateMatchingRequestImplBase implements io.grpc.BindableService {

    /**
     */
    public void createMatchingRequest(profile.MatchingCreatation.MatchingRequest request,
        io.grpc.stub.StreamObserver<profile.MatchingCreatation.MatchingResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateMatchingRequestMethod(), responseObserver);
    }

    /**
     */
    public void responseRequestMatching(profile.MatchingCreatation.ResponseMatchingRequest request,
        io.grpc.stub.StreamObserver<profile.MatchingCreatation.ResponseMatchingResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getResponseRequestMatchingMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCreateMatchingRequestMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                profile.MatchingCreatation.MatchingRequest,
                profile.MatchingCreatation.MatchingResponse>(
                  this, METHODID_CREATE_MATCHING_REQUEST)))
          .addMethod(
            getResponseRequestMatchingMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                profile.MatchingCreatation.ResponseMatchingRequest,
                profile.MatchingCreatation.ResponseMatchingResponse>(
                  this, METHODID_RESPONSE_REQUEST_MATCHING)))
          .build();
    }
  }

  /**
   */
  public static final class CreateMatchingRequestStub extends io.grpc.stub.AbstractAsyncStub<CreateMatchingRequestStub> {
    private CreateMatchingRequestStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CreateMatchingRequestStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CreateMatchingRequestStub(channel, callOptions);
    }

    /**
     */
    public void createMatchingRequest(profile.MatchingCreatation.MatchingRequest request,
        io.grpc.stub.StreamObserver<profile.MatchingCreatation.MatchingResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCreateMatchingRequestMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void responseRequestMatching(profile.MatchingCreatation.ResponseMatchingRequest request,
        io.grpc.stub.StreamObserver<profile.MatchingCreatation.ResponseMatchingResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getResponseRequestMatchingMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class CreateMatchingRequestBlockingStub extends io.grpc.stub.AbstractBlockingStub<CreateMatchingRequestBlockingStub> {
    private CreateMatchingRequestBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CreateMatchingRequestBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CreateMatchingRequestBlockingStub(channel, callOptions);
    }

    /**
     */
    public profile.MatchingCreatation.MatchingResponse createMatchingRequest(profile.MatchingCreatation.MatchingRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateMatchingRequestMethod(), getCallOptions(), request);
    }

    /**
     */
    public profile.MatchingCreatation.ResponseMatchingResponse responseRequestMatching(profile.MatchingCreatation.ResponseMatchingRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getResponseRequestMatchingMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class CreateMatchingRequestFutureStub extends io.grpc.stub.AbstractFutureStub<CreateMatchingRequestFutureStub> {
    private CreateMatchingRequestFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CreateMatchingRequestFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CreateMatchingRequestFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<profile.MatchingCreatation.MatchingResponse> createMatchingRequest(
        profile.MatchingCreatation.MatchingRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCreateMatchingRequestMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<profile.MatchingCreatation.ResponseMatchingResponse> responseRequestMatching(
        profile.MatchingCreatation.ResponseMatchingRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getResponseRequestMatchingMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE_MATCHING_REQUEST = 0;
  private static final int METHODID_RESPONSE_REQUEST_MATCHING = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final CreateMatchingRequestImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(CreateMatchingRequestImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE_MATCHING_REQUEST:
          serviceImpl.createMatchingRequest((profile.MatchingCreatation.MatchingRequest) request,
              (io.grpc.stub.StreamObserver<profile.MatchingCreatation.MatchingResponse>) responseObserver);
          break;
        case METHODID_RESPONSE_REQUEST_MATCHING:
          serviceImpl.responseRequestMatching((profile.MatchingCreatation.ResponseMatchingRequest) request,
              (io.grpc.stub.StreamObserver<profile.MatchingCreatation.ResponseMatchingResponse>) responseObserver);
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

  private static abstract class CreateMatchingRequestBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    CreateMatchingRequestBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return profile.MatchingCreatation.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("CreateMatchingRequest");
    }
  }

  private static final class CreateMatchingRequestFileDescriptorSupplier
      extends CreateMatchingRequestBaseDescriptorSupplier {
    CreateMatchingRequestFileDescriptorSupplier() {}
  }

  private static final class CreateMatchingRequestMethodDescriptorSupplier
      extends CreateMatchingRequestBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    CreateMatchingRequestMethodDescriptorSupplier(String methodName) {
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
      synchronized (CreateMatchingRequestGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new CreateMatchingRequestFileDescriptorSupplier())
              .addMethod(getCreateMatchingRequestMethod())
              .addMethod(getResponseRequestMatchingMethod())
              .build();
        }
      }
    }
    return result;
  }
}
