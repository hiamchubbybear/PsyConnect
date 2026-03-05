



package networktype

import (
	"google.golang.org/grpc/resolver"
)


type keyType string

const key = keyType("grpc.internal.transport.networktype")


func Set(address resolver.Address, networkType string) resolver.Address {
	address.Attributes = address.Attributes.WithValue(key, networkType)
	return address
}



func Get(address resolver.Address) (string, bool) {
	v := address.Attributes.Value(key)
	if v == nil {
		return "", false
	}
	return v.(string), true
}
