





//go:build gssapi && !windows && !linux && !darwin
// +build gssapi,!windows,!linux,!darwin

package auth

import (
	"fmt"
	"net/http"
	"runtime"
)


const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(*Cred, *http.Client) (Authenticator, error) {
	return nil, newAuthError(fmt.Sprintf("GSSAPI is not supported on %s", runtime.GOOS), nil)
}
