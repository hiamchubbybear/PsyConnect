





//go:build !gssapi
// +build !gssapi

package auth

import "net/http"


const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(*Cred, *http.Client) (Authenticator, error) {
	return nil, newAuthError("GSSAPI support not enabled during build (-tags gssapi)", nil)
}
