package dev.psyconnect.grpc.api_gateway;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: token_check_identity.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class TokenServiceGrpc {

  private TokenServiceGrpc() {}

  public static final String SERVICE_NAME = "identity.TokenService";

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
    if ((getTokenCheckValidMethod = TokenServiceGrpc.getTokenCheckValidMethod) == null) {
      synchronized (TokenServiceGrpc.class) {
        if ((getTokenCheckValidMethod = TokenServiceGrpc.getTokenCheckValidMethod) == null) {
          TokenServiceGrpc.getTokenCheckValidMethod = getTokenCheckValidMethod =
              io.grpc.MethodDescriptor.<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest, dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "TokenCheckValid"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TokenServiceMethodDescriptorSupplier("TokenCheckValid"))
              .build();
        }
      }
    }
    return getTokenCheckValidMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static TokenServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TokenServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TokenServiceStub>() {
        @java.lang.Override
        public TokenServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TokenServiceStub(channel, callOptions);
        }
      };
    return TokenServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static TokenServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TokenServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TokenServiceBlockingStub>() {
        @java.lang.Override
        public TokenServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TokenServiceBlockingStub(channel, callOptions);
        }
      };
    return TokenServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static TokenServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TokenServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TokenServiceFutureStub>() {
        @java.lang.Override
        public TokenServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TokenServiceFutureStub(channel, callOptions);
        }
      };
    return TokenServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class TokenServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void tokenCheckValid(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getTokenCheckValidMethod(), responseObserver);
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
          .build();
    }
  }

  /**
   */
  public static final class TokenServiceStub extends io.grpc.stub.AbstractAsyncStub<TokenServiceStub> {
    private TokenServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TokenServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TokenServiceStub(channel, callOptions);
    }

    /**
     */
    public void tokenCheckValid(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request,
        io.grpc.stub.StreamObserver<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getTokenCheckValidMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class TokenServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<TokenServiceBlockingStub> {
    private TokenServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TokenServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TokenServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse tokenCheckValid(dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getTokenCheckValidMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class TokenServiceFutureStub extends io.grpc.stub.AbstractFutureStub<TokenServiceFutureStub> {
    private TokenServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TokenServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TokenServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenResponse> tokenCheckValid(
        dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.TokenRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getTokenCheckValidMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_TOKEN_CHECK_VALID = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final TokenServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(TokenServiceImplBase serviceImpl, int methodId) {
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

  private static abstract class TokenServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    TokenServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return dev.psyconnect.grpc.api_gateway.TokenCheckIdentity.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("TokenService");
    }
  }

  private static final class TokenServiceFileDescriptorSupplier
      extends TokenServiceBaseDescriptorSupplier {
    TokenServiceFileDescriptorSupplier() {}
  }

  private static final class TokenServiceMethodDescriptorSupplier
      extends TokenServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    TokenServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (TokenServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new TokenServiceFileDescriptorSupplier())
              .addMethod(getTokenCheckValidMethod())
              .build();
        }
      }
    }
    return result;
  }
}
