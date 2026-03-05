





package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

const sourceExternal = "$external"


type Config = driver.AuthConfig


type AuthenticatorFactory func(*Cred, *http.Client) (Authenticator, error)

var authFactories = make(map[string]AuthenticatorFactory)

func init() {
	RegisterAuthenticatorFactory("", newDefaultAuthenticator)
	RegisterAuthenticatorFactory(SCRAMSHA1, newScramSHA1Authenticator)
	RegisterAuthenticatorFactory(SCRAMSHA256, newScramSHA256Authenticator)
	RegisterAuthenticatorFactory(MONGODBCR, newMongoDBCRAuthenticator)
	RegisterAuthenticatorFactory(PLAIN, newPlainAuthenticator)
	RegisterAuthenticatorFactory(GSSAPI, newGSSAPIAuthenticator)
	RegisterAuthenticatorFactory(MongoDBX509, newMongoDBX509Authenticator)
	RegisterAuthenticatorFactory(MongoDBAWS, newMongoDBAWSAuthenticator)
	RegisterAuthenticatorFactory(MongoDBOIDC, newOIDCAuthenticator)
}


func CreateAuthenticator(name string, cred *Cred, httpClient *http.Client) (Authenticator, error) {
	if f, ok := authFactories[name]; ok {
		return f(cred, httpClient)
	}

	return nil, newAuthError(fmt.Sprintf("unknown authenticator: %s", name), nil)
}


func RegisterAuthenticatorFactory(name string, factory AuthenticatorFactory) {
	authFactories[name] = factory
}




type HandshakeOptions struct {
	AppName               string
	Authenticator         Authenticator
	Compressors           []string
	DBUser                string
	PerformAuthentication func(description.Server) bool
	ClusterClock          *session.ClusterClock
	ServerAPI             *driver.ServerAPIOptions
	LoadBalanced          bool
}

type authHandshaker struct {
	wrapped driver.Handshaker
	options *HandshakeOptions

	handshakeInfo driver.HandshakeInformation
	conversation  SpeculativeConversation
}

var _ driver.Handshaker = (*authHandshaker)(nil)



func (ah *authHandshaker) GetHandshakeInformation(ctx context.Context, addr address.Address, conn driver.Connection) (driver.HandshakeInformation, error) {
	if ah.wrapped != nil {
		return ah.wrapped.GetHandshakeInformation(ctx, addr, conn)
	}

	op := operation.NewHello().
		AppName(ah.options.AppName).
		Compressors(ah.options.Compressors).
		SASLSupportedMechs(ah.options.DBUser).
		ClusterClock(ah.options.ClusterClock).
		ServerAPI(ah.options.ServerAPI).
		LoadBalanced(ah.options.LoadBalanced)

	if ah.options.Authenticator != nil {
		if speculativeAuth, ok := ah.options.Authenticator.(SpeculativeAuthenticator); ok {
			var err error
			ah.conversation, err = speculativeAuth.CreateSpeculativeConversation()
			if err != nil {
				return driver.HandshakeInformation{}, newAuthError("failed to create conversation", err)
			}

			
			
			
			if ah.conversation != nil {
				firstMsg, err := ah.conversation.FirstMessage()
				if err != nil {
					return driver.HandshakeInformation{}, newAuthError("failed to create speculative authentication message", err)
				}

				op = op.SpeculativeAuthenticate(firstMsg)
			}
		}
	}

	var err error
	ah.handshakeInfo, err = op.GetHandshakeInformation(ctx, addr, conn)
	if err != nil {
		return driver.HandshakeInformation{}, newAuthError("handshake failure", err)
	}
	return ah.handshakeInfo, nil
}


func (ah *authHandshaker) FinishHandshake(ctx context.Context, conn driver.Connection) error {
	performAuth := ah.options.PerformAuthentication
	if performAuth == nil {
		performAuth = func(serv description.Server) bool {
			
			return serv.Kind != description.RSArbiter
		}
	}

	desc := conn.Description()
	if performAuth(desc) && ah.options.Authenticator != nil {
		cfg := &Config{
			Description:   desc,
			Connection:    conn,
			ClusterClock:  ah.options.ClusterClock,
			HandshakeInfo: ah.handshakeInfo,
			ServerAPI:     ah.options.ServerAPI,
		}

		if err := ah.authenticate(ctx, cfg); err != nil {
			return newAuthError("auth error", err)
		}
	}

	if ah.wrapped == nil {
		return nil
	}
	return ah.wrapped.FinishHandshake(ctx, conn)
}

func (ah *authHandshaker) authenticate(ctx context.Context, cfg *Config) error {
	
	
	if speculativeResponse := ah.handshakeInfo.SpeculativeAuthenticate; speculativeResponse != nil {
		
		if ah.conversation == nil {
			return errors.New("speculative auth was not attempted but the server included a response")
		}
		return ah.conversation.Finish(ctx, cfg, speculativeResponse)
	}

	
	
	return ah.options.Authenticator.Auth(ctx, cfg)
}


func Handshaker(h driver.Handshaker, options *HandshakeOptions) driver.Handshaker {
	return &authHandshaker{
		wrapped: h,
		options: options,
	}
}


type Authenticator = driver.Authenticator

func newAuthError(msg string, inner error) error {
	return &Error{
		message: msg,
		inner:   inner,
	}
}

func newError(err error, mech string) error {
	return &Error{
		message: fmt.Sprintf("unable to authenticate using mechanism \"%s\"", mech),
		inner:   err,
	}
}


type Error struct {
	message string
	inner   error
}

func (e *Error) Error() string {
	if e.inner == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %s", e.message, e.inner)
}


func (e *Error) Inner() error {
	return e.inner
}


func (e *Error) Unwrap() error {
	return e.inner
}


func (e *Error) Message() string {
	return e.message
}
