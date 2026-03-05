









package credentials

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/internal/aws/awserr"
	"golang.org/x/sync/singleflight"
)





type Value struct {
	
	AccessKeyID string

	
	SecretAccessKey string

	
	SessionToken string

	
	ProviderName string
}



func (v Value) HasKeys() bool {
	return len(v.AccessKeyID) != 0 && len(v.SecretAccessKey) != 0
}







type Provider interface {
	
	
	Retrieve() (Value, error)

	
	
	IsExpired() bool
}


type ProviderWithContext interface {
	Provider

	RetrieveWithContext(context.Context) (Value, error)
}















type Credentials struct {
	sf singleflight.Group

	m        sync.RWMutex
	creds    Value
	provider Provider
}


func NewCredentials(provider Provider) *Credentials {
	c := &Credentials{
		provider: provider,
	}
	return c
}











func (c *Credentials) GetWithContext(ctx context.Context) (Value, error) {
	
	select {
	case curCreds, ok := <-c.asyncIsExpired():
		
		
		if ok {
			return curCreds, nil
		}
	case <-ctx.Done():
		return Value{}, awserr.New("RequestCanceled",
			"request context canceled", ctx.Err())
	}

	
	
	
	resCh := c.sf.DoChan("", func() (interface{}, error) {
		return c.singleRetrieve(&suppressedContext{ctx})
	})
	select {
	case res := <-resCh:
		return res.Val.(Value), res.Err
	case <-ctx.Done():
		return Value{}, awserr.New("RequestCanceled",
			"request context canceled", ctx.Err())
	}
}

func (c *Credentials) singleRetrieve(ctx context.Context) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if curCreds := c.creds; !c.isExpiredLocked(curCreds) {
		return curCreds, nil
	}

	var creds Value
	var err error
	if p, ok := c.provider.(ProviderWithContext); ok {
		creds, err = p.RetrieveWithContext(ctx)
	} else {
		creds, err = c.provider.Retrieve()
	}
	if err == nil {
		c.creds = creds
	}

	return creds, err
}



func (c *Credentials) asyncIsExpired() <-chan Value {
	ch := make(chan Value, 1)
	go func() {
		c.m.RLock()
		defer c.m.RUnlock()

		if curCreds := c.creds; !c.isExpiredLocked(curCreds) {
			ch <- curCreds
		}

		close(ch)
	}()

	return ch
}


func (c *Credentials) isExpiredLocked(creds interface{}) bool {
	return creds == nil || creds.(Value) == Value{} || c.provider.IsExpired()
}

type suppressedContext struct {
	context.Context
}

func (s *suppressedContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (s *suppressedContext) Done() <-chan struct{} {
	return nil
}

func (s *suppressedContext) Err() error {
	return nil
}
