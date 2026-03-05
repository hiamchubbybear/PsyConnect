





package ocsp

import (
	"crypto"
	"sync"
	"time"

	"golang.org/x/crypto/ocsp"
)

type cacheKey struct {
	HashAlgorithm  crypto.Hash
	IssuerNameHash string
	IssuerKeyHash  string
	SerialNumber   string
}


type Cache interface {
	Update(*ocsp.Request, *ResponseDetails) *ResponseDetails
	Get(request *ocsp.Request) *ResponseDetails
}


type ConcurrentCache struct {
	cache map[cacheKey]*ResponseDetails
	sync.Mutex
}

var _ Cache = (*ConcurrentCache)(nil)


func NewCache() *ConcurrentCache {
	return &ConcurrentCache{
		cache: make(map[cacheKey]*ResponseDetails),
	}
}







func (c *ConcurrentCache) Update(request *ocsp.Request, response *ResponseDetails) *ResponseDetails {
	unknown := response.Status == ocsp.Unknown
	hasUpdateTime := !response.NextUpdate.IsZero()
	canBeCached := !unknown && hasUpdateTime
	key := createCacheKey(request)

	c.Lock()
	defer c.Unlock()

	current, ok := c.cache[key]
	if !ok {
		if canBeCached {
			c.cache[key] = response
		}

		
		
		return response
	}

	
	if unknown {
		return current
	}

	
	
	
	if !hasUpdateTime {
		delete(c.cache, key)
		return response
	}

	
	
	newest := current
	if response.NextUpdate.After(current.NextUpdate) {
		c.cache[key] = response
		newest = response
	}
	return newest
}



func (c *ConcurrentCache) Get(request *ocsp.Request) *ResponseDetails {
	key := createCacheKey(request)

	c.Lock()
	defer c.Unlock()

	response, ok := c.cache[key]
	if !ok {
		return nil
	}

	if time.Now().UTC().Before(response.NextUpdate) {
		return response
	}
	delete(c.cache, key)
	return nil
}

func createCacheKey(request *ocsp.Request) cacheKey {
	return cacheKey{
		HashAlgorithm:  request.HashAlgorithm,
		IssuerNameHash: string(request.IssuerNameHash),
		IssuerKeyHash:  string(request.IssuerKeyHash),
		SerialNumber:   request.SerialNumber.String(),
	}
}
