package auth






type ReAuthCredentialsListener struct {
	reAuth func(credentials Credentials) error
	onErr  func(err error)
}




func (c *ReAuthCredentialsListener) OnNext(credentials Credentials) {
	if c.reAuth == nil {
		return
	}

	err := c.reAuth(credentials)
	if err != nil {
		c.OnError(err)
	}
}



func (c *ReAuthCredentialsListener) OnError(err error) {
	if c.onErr == nil {
		return
	}

	c.onErr(err)
}



func NewReAuthCredentialsListener(reAuth func(credentials Credentials) error, onErr func(err error)) *ReAuthCredentialsListener {
	return &ReAuthCredentialsListener{
		reAuth: reAuth,
		onErr:  onErr,
	}
}


var _ CredentialsListener = (*ReAuthCredentialsListener)(nil)
