package dev.psyconnect.identity_service.grpc.server;

import java.util.Optional;
import java.util.UUID;

import dev.psyconnect.grpc.api_gateway.IdentityServiceGrpc;
import dev.psyconnect.grpc.api_gateway.TokenCheckIdentity;
import dev.psyconnect.identity_service.model.Account;
import dev.psyconnect.identity_service.repository.BlackListTokenRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;

@GrpcService
public class IdentityGatewayServer extends IdentityServiceGrpc.IdentityServiceImplBase {

    private final BlackListTokenRepository blackListTokenRepository;
    private final UserAccountRepository userAccountRepository;

    public IdentityGatewayServer(
            BlackListTokenRepository blackListTokenRepository, UserAccountRepository userAccountRepository) {
        super();
        this.blackListTokenRepository = blackListTokenRepository;
        this.userAccountRepository = userAccountRepository;
    }

    @Override
    public void tokenCheckValid(
            TokenCheckIdentity.TokenRequest request,
            StreamObserver<TokenCheckIdentity.TokenResponse> responseObserver) {

        Boolean valid = !blackListTokenRepository.existsByToken(request.getToken());
        TokenCheckIdentity.TokenResponse response =
                TokenCheckIdentity.TokenResponse.newBuilder().setValid(valid).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void getUserInfoByProfileId(
            TokenCheckIdentity.UserInfoRequest request,
            StreamObserver<TokenCheckIdentity.UserInfoResponse> responseObserver) {
        try {
            Optional<Account> account = userAccountRepository.findByProfileId(UUID.fromString(request.getProfileId()));
            if (account.isPresent()) {
                responseObserver.onNext(TokenCheckIdentity.UserInfoResponse.newBuilder()
                        .setSuccess(true)
                        .setEmail(account.get().getEmail())
                        .setUsername(account.get().getUsername())
                        .build());
            } else {
                responseObserver.onNext(TokenCheckIdentity.UserInfoResponse.newBuilder()
                        .setSuccess(false)
                        .build());
            }
        } catch (Exception e) {
            responseObserver.onNext(TokenCheckIdentity.UserInfoResponse.newBuilder()
                    .setSuccess(false)
                    .build());
        } finally {
            responseObserver.onCompleted();
        }
    }
}
