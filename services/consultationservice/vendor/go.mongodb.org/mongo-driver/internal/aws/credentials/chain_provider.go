









package credentials

import (
	"go.mongodb.org/mongo-driver/internal/aws/awserr"
)














type ChainProvider struct {
	Providers []Provider
	curr      Provider
}



func NewChainCredentials(providers []Provider) *Credentials {
	return NewCredentials(&ChainProvider{
		Providers: append([]Provider{}, providers...),
	})
}






func (c *ChainProvider) Retrieve() (Value, error) {
	var errs = make([]error, 0, len(c.Providers))
	for _, p := range c.Providers {
		creds, err := p.Retrieve()
		if err == nil {
			c.curr = p
			return creds, nil
		}
		errs = append(errs, err)
	}
	c.curr = nil

	var err = awserr.NewBatchError("NoCredentialProviders", "no valid providers in chain", errs)
	return Value{}, err
}



func (c *ChainProvider) IsExpired() bool {
	if c.curr != nil {
		return c.curr.IsExpired()
	}

	return true
}
