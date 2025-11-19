package dev.psyconnect.identity_service.grpc.server;

import dev.psyconnect.grpc.api_gateway.TokenCheckIdentity;
import dev.psyconnect.grpc.api_gateway.TokenServiceGrpc;
import dev.psyconnect.identity_service.repository.BlackListTokenRepository;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;

@GrpcService
public class TokenGatewayServer extends TokenServiceGrpc.TokenServiceImplBase {

    private final BlackListTokenRepository blackListTokenRepository;

    public TokenGatewayServer(BlackListTokenRepository blackListTokenRepository) {
        super();
        this.blackListTokenRepository = blackListTokenRepository;
    }

    @Override
    public void tokenCheckValid(
            TokenCheckIdentity.TokenRequest request,
            StreamObserver<TokenCheckIdentity.TokenResponse> responseObserver) {
        // Check token
        Boolean valid = !blackListTokenRepository.existsByToken(request.getToken());
        TokenCheckIdentity.TokenResponse response =
                TokenCheckIdentity.TokenResponse.newBuilder().setValid(valid).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
