



package proxyattributes

import (
	"net/url"

	"google.golang.org/grpc/resolver"
)

type keyType string

const proxyOptionsKey = keyType("grpc.resolver.delegatingresolver.proxyOptions")



type Options struct {
	User        *url.Userinfo
	ConnectAddr string
}


func Set(addr resolver.Address, opts Options) resolver.Address {
	addr.Attributes = addr.Attributes.WithValue(proxyOptionsKey, opts)
	return addr
}




func Get(addr resolver.Address) (Options, bool) {
	if a := addr.Attributes.Value(proxyOptionsKey); a != nil {
		return a.(Options), true
	}
	return Options{}, false
}
