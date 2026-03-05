





package auth

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
)


const PLAIN = "PLAIN"

func newPlainAuthenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	return &PlainAuthenticator{
		Username: cred.Username,
		Password: cred.Password,
	}, nil
}


type PlainAuthenticator struct {
	Username string
	Password string
}


func (a *PlainAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	return ConductSaslConversation(ctx, cfg, sourceExternal, &plainSaslClient{
		username: a.Username,
		password: a.Password,
	})
}


func (a *PlainAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("Plain authentication does not support reauthentication", nil)
}

type plainSaslClient struct {
	username string
	password string
}

var _ SaslClient = (*plainSaslClient)(nil)

func (c *plainSaslClient) Start() (string, []byte, error) {
	b := []byte("\x00" + c.username + "\x00" + c.password)
	return PLAIN, b, nil
}

func (c *plainSaslClient) Next(context.Context, []byte) ([]byte, error) {
	return nil, newAuthError("unexpected server challenge", nil)
}

func (c *plainSaslClient) Completed() bool {
	return true
}
