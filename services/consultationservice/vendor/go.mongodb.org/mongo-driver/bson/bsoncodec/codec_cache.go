





package bsoncodec

import (
	"reflect"
	"sync"
	"sync/atomic"
)



func init() {
	if s := reflect.Kind(len(kindEncoderCache{}.entries)).String(); s != "kind27" {
		panic("The capacity of kindEncoderCache is too small.\n" +
			"This is due to a new type being added to reflect.Kind.")
	}
}


var _ = (kindEncoderCache{}).entries[reflect.UnsafePointer]
var _ = (kindDecoderCache{}).entries[reflect.UnsafePointer]

type typeEncoderCache struct {
	cache sync.Map 
}

func (c *typeEncoderCache) Store(rt reflect.Type, enc ValueEncoder) {
	c.cache.Store(rt, enc)
}

func (c *typeEncoderCache) Load(rt reflect.Type) (ValueEncoder, bool) {
	if v, _ := c.cache.Load(rt); v != nil {
		return v.(ValueEncoder), true
	}
	return nil, false
}

func (c *typeEncoderCache) LoadOrStore(rt reflect.Type, enc ValueEncoder) ValueEncoder {
	if v, loaded := c.cache.LoadOrStore(rt, enc); loaded {
		enc = v.(ValueEncoder)
	}
	return enc
}

func (c *typeEncoderCache) Clone() *typeEncoderCache {
	cc := new(typeEncoderCache)
	c.cache.Range(func(k, v interface{}) bool {
		if k != nil && v != nil {
			cc.cache.Store(k, v)
		}
		return true
	})
	return cc
}

type typeDecoderCache struct {
	cache sync.Map 
}

func (c *typeDecoderCache) Store(rt reflect.Type, dec ValueDecoder) {
	c.cache.Store(rt, dec)
}

func (c *typeDecoderCache) Load(rt reflect.Type) (ValueDecoder, bool) {
	if v, _ := c.cache.Load(rt); v != nil {
		return v.(ValueDecoder), true
	}
	return nil, false
}

func (c *typeDecoderCache) LoadOrStore(rt reflect.Type, dec ValueDecoder) ValueDecoder {
	if v, loaded := c.cache.LoadOrStore(rt, dec); loaded {
		dec = v.(ValueDecoder)
	}
	return dec
}

func (c *typeDecoderCache) Clone() *typeDecoderCache {
	cc := new(typeDecoderCache)
	c.cache.Range(func(k, v interface{}) bool {
		if k != nil && v != nil {
			cc.cache.Store(k, v)
		}
		return true
	})
	return cc
}





type kindEncoderCacheEntry struct {
	enc ValueEncoder
}

type kindEncoderCache struct {
	entries [reflect.UnsafePointer + 1]atomic.Value 
}

func (c *kindEncoderCache) Store(rt reflect.Kind, enc ValueEncoder) {
	if enc != nil && rt < reflect.Kind(len(c.entries)) {
		c.entries[rt].Store(&kindEncoderCacheEntry{enc: enc})
	}
}

func (c *kindEncoderCache) Load(rt reflect.Kind) (ValueEncoder, bool) {
	if rt < reflect.Kind(len(c.entries)) {
		if ent, ok := c.entries[rt].Load().(*kindEncoderCacheEntry); ok {
			return ent.enc, ent.enc != nil
		}
	}
	return nil, false
}

func (c *kindEncoderCache) Clone() *kindEncoderCache {
	cc := new(kindEncoderCache)
	for i, v := range c.entries {
		if val := v.Load(); val != nil {
			cc.entries[i].Store(val)
		}
	}
	return cc
}





type kindDecoderCacheEntry struct {
	dec ValueDecoder
}

type kindDecoderCache struct {
	entries [reflect.UnsafePointer + 1]atomic.Value 
}

func (c *kindDecoderCache) Store(rt reflect.Kind, dec ValueDecoder) {
	if rt < reflect.Kind(len(c.entries)) {
		c.entries[rt].Store(&kindDecoderCacheEntry{dec: dec})
	}
}

func (c *kindDecoderCache) Load(rt reflect.Kind) (ValueDecoder, bool) {
	if rt < reflect.Kind(len(c.entries)) {
		if ent, ok := c.entries[rt].Load().(*kindDecoderCacheEntry); ok {
			return ent.dec, ent.dec != nil
		}
	}
	return nil, false
}

func (c *kindDecoderCache) Clone() *kindDecoderCache {
	cc := new(kindDecoderCache)
	for i, v := range c.entries {
		if val := v.Load(); val != nil {
			cc.entries[i].Store(val)
		}
	}
	return cc
}
