





package scram

import (
	"sync"

	"github.com/xdg-go/pbkdf2"
)















type Client struct {
	sync.RWMutex
	username string
	password string
	authzID  string
	minIters int
	nonceGen NonceGeneratorFcn
	hashGen  HashGeneratorFcn
	cache    map[KeyFactors]derivedKeys
}

func newClient(username, password, authzID string, fcn HashGeneratorFcn) *Client {
	return &Client{
		username: username,
		password: password,
		authzID:  authzID,
		minIters: 4096,
		nonceGen: defaultNonceGenerator,
		hashGen:  fcn,
		cache:    make(map[KeyFactors]derivedKeys),
	}
}


func (c *Client) WithMinIterations(n int) *Client {
	c.Lock()
	defer c.Unlock()
	c.minIters = n
	return c
}




func (c *Client) WithNonceGenerator(ng NonceGeneratorFcn) *Client {
	c.Lock()
	defer c.Unlock()
	c.nonceGen = ng
	return c
}




func (c *Client) NewConversation() *ClientConversation {
	c.RLock()
	defer c.RUnlock()
	return &ClientConversation{
		client:   c,
		nonceGen: c.nonceGen,
		hashGen:  c.hashGen,
		minIters: c.minIters,
	}
}

func (c *Client) getDerivedKeys(kf KeyFactors) derivedKeys {
	dk, ok := c.getCache(kf)
	if !ok {
		dk = c.computeKeys(kf)
		c.setCache(kf, dk)
	}
	return dk
}





func (c *Client) GetStoredCredentials(kf KeyFactors) StoredCredentials {
	dk := c.getDerivedKeys(kf)
	return StoredCredentials{
		KeyFactors: kf,
		StoredKey:  dk.StoredKey,
		ServerKey:  dk.ServerKey,
	}
}

func (c *Client) computeKeys(kf KeyFactors) derivedKeys {
	h := c.hashGen()
	saltedPassword := pbkdf2.Key([]byte(c.password), []byte(kf.Salt), kf.Iters, h.Size(), c.hashGen)
	clientKey := computeHMAC(c.hashGen, saltedPassword, []byte("Client Key"))

	return derivedKeys{
		ClientKey: clientKey,
		StoredKey: computeHash(c.hashGen, clientKey),
		ServerKey: computeHMAC(c.hashGen, saltedPassword, []byte("Server Key")),
	}
}

func (c *Client) getCache(kf KeyFactors) (derivedKeys, bool) {
	c.RLock()
	defer c.RUnlock()
	dk, ok := c.cache[kf]
	return dk, ok
}

func (c *Client) setCache(kf KeyFactors, dk derivedKeys) {
	c.Lock()
	defer c.Unlock()
	c.cache[kf] = dk
	return
}
