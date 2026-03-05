

package auth




type StreamingCredentialsProvider interface {
	
	
	
	
	Subscribe(listener CredentialsListener) (Credentials, UnsubscribeFunc, error)
}



type UnsubscribeFunc func() error





type CredentialsListener interface {
	OnNext(credentials Credentials)
	OnError(err error)
}



type Credentials interface {
	
	BasicAuth() (username string, password string)
	
	
	
	RawCredentials() string
}

type basicAuth struct {
	username string
	password string
}


func (b *basicAuth) RawCredentials() string {
	return b.username + ":" + b.password
}


func (b *basicAuth) BasicAuth() (username string, password string) {
	return b.username, b.password
}


func NewBasicCredentials(username, password string) Credentials {
	return &basicAuth{
		username: username,
		password: password,
	}
}
