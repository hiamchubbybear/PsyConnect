





package auth

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

func newDefaultAuthenticator(cred *Cred, httpClient *http.Client) (Authenticator, error) {
	scram, err := newScramSHA256Authenticator(cred, httpClient)
	if err != nil {
		return nil, newAuthError("failed to create internal authenticator", err)
	}
	speculative, ok := scram.(SpeculativeAuthenticator)
	if !ok {
		typeErr := fmt.Errorf("expected SCRAM authenticator to be SpeculativeAuthenticator but got %T", scram)
		return nil, newAuthError("failed to create internal authenticator", typeErr)
	}

	return &DefaultAuthenticator{
		Cred:                     cred,
		speculativeAuthenticator: speculative,
		httpClient:               httpClient,
	}, nil
}



type DefaultAuthenticator struct {
	Cred *Cred

	
	
	speculativeAuthenticator SpeculativeAuthenticator

	httpClient *http.Client
}

var _ SpeculativeAuthenticator = (*DefaultAuthenticator)(nil)


func (a *DefaultAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	return a.speculativeAuthenticator.CreateSpeculativeConversation()
}


func (a *DefaultAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	var actual Authenticator
	var err error

	switch chooseAuthMechanism(cfg) {
	case SCRAMSHA256:
		actual, err = newScramSHA256Authenticator(a.Cred, a.httpClient)
	case SCRAMSHA1:
		actual, err = newScramSHA1Authenticator(a.Cred, a.httpClient)
	default:
		actual, err = newMongoDBCRAuthenticator(a.Cred, a.httpClient)
	}

	if err != nil {
		return newAuthError("error creating authenticator", err)
	}

	return actual.Auth(ctx, cfg)
}


func (a *DefaultAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("DefaultAuthenticator does not support reauthentication", nil)
}




func chooseAuthMechanism(cfg *Config) string {
	if saslSupportedMechs := cfg.HandshakeInfo.SaslSupportedMechs; saslSupportedMechs != nil {
		for _, v := range saslSupportedMechs {
			if v == SCRAMSHA256 {
				return v
			}
		}
	}

	return SCRAMSHA1
}
