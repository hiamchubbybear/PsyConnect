



package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/internal/bytesconv"
)


const AuthUserKey = "user"


const AuthProxyUserKey = "proxy_user"


type Accounts map[string]string

type authPair struct {
	value string
	user  string
}

type authPairs []authPair

func (a authPairs) searchCredential(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if subtle.ConstantTimeCompare(bytesconv.StringToBytes(pair.value), bytesconv.StringToBytes(authValue)) == 1 {
			return pair.user, true
		}
	}
	return "", false
}





func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		
		user, found := pairs.searchCredential(c.requestHeader("Authorization"))
		if !found {
			
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		
		
		c.Set(AuthUserKey, user)
	}
}



func BasicAuth(accounts Accounts) HandlerFunc {
	return BasicAuthForRealm(accounts, "")
}

func processAccounts(accounts Accounts) authPairs {
	length := len(accounts)
	assert1(length > 0, "Empty list of authorized credentials")
	pairs := make(authPairs, 0, length)
	for user, password := range accounts {
		assert1(user != "", "User can not be empty")
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			value: value,
			user:  user,
		})
	}
	return pairs
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(bytesconv.StringToBytes(base))
}



func BasicAuthForProxy(accounts Accounts, realm string) HandlerFunc {
	if realm == "" {
		realm = "Proxy Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *Context) {
		proxyUser, found := pairs.searchCredential(c.requestHeader("Proxy-Authorization"))
		if !found {
			
			c.Header("Proxy-Authenticate", realm)
			c.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		}
		
		
		c.Set(AuthProxyUserKey, proxyUser)
	}
}
