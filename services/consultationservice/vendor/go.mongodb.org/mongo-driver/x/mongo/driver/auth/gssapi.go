





//go:build gssapi && (windows || linux || darwin)
// +build gssapi
// +build windows linux darwin

package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/internal/gssapi"
)


const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	if cred.Source != "" && cred.Source != sourceExternal {
		return nil, newAuthError("GSSAPI source must be empty or $external", nil)
	}

	return &GSSAPIAuthenticator{
		Username:    cred.Username,
		Password:    cred.Password,
		PasswordSet: cred.PasswordSet,
		Props:       cred.Props,
	}, nil
}


type GSSAPIAuthenticator struct {
	Username    string
	Password    string
	PasswordSet bool
	Props       map[string]string
}


func (a *GSSAPIAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	target := cfg.Description.Addr.String()
	hostname, _, err := net.SplitHostPort(target)
	if err != nil {
		return newAuthError(fmt.Sprintf("invalid endpoint (%s) specified: %s", target, err), nil)
	}

	client, err := gssapi.New(hostname, a.Username, a.Password, a.PasswordSet, a.Props)

	if err != nil {
		return newAuthError("error creating gssapi", err)
	}
	return ConductSaslConversation(ctx, cfg, sourceExternal, client)
}


func (a *GSSAPIAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("GSSAPI does not support reauthentication", nil)
}
