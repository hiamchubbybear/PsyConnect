




package metadata

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
)

type mdKeyType string

const mdKey = mdKeyType("grpc.internal.address.metadata")

type mdValue metadata.MD

func (m mdValue) Equal(o any) bool {
	om, ok := o.(mdValue)
	if !ok {
		return false
	}
	if len(m) != len(om) {
		return false
	}
	for k, v := range m {
		ov := om[k]
		if len(ov) != len(v) {
			return false
		}
		for i, ve := range v {
			if ov[i] != ve {
				return false
			}
		}
	}
	return true
}


func Get(addr resolver.Address) metadata.MD {
	attrs := addr.Attributes
	if attrs == nil {
		return nil
	}
	md, _ := attrs.Value(mdKey).(mdValue)
	return metadata.MD(md)
}





func Set(addr resolver.Address, md metadata.MD) resolver.Address {
	addr.Attributes = addr.Attributes.WithValue(mdKey, mdValue(md))
	return addr
}


func Validate(md metadata.MD) error {
	for k, vals := range md {
		if err := ValidatePair(k, vals...); err != nil {
			return err
		}
	}
	return nil
}


func hasNotPrintable(msg string) bool {
	
	for i := 0; i < len(msg); i++ {
		if msg[i] < 0x20 || msg[i] > 0x7E {
			return true
		}
	}
	return false
}







func ValidatePair(key string, vals ...string) error {
	
	if key == "" {
		return fmt.Errorf("there is an empty key in the header")
	}
	
	if key[0] == ':' {
		return nil
	}
	
	for i := 0; i < len(key); i++ {
		r := key[i]
		if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '.' && r != '-' && r != '_' {
			return fmt.Errorf("header key %q contains illegal characters not in [0-9a-z-_.]", key)
		}
	}
	if strings.HasSuffix(key, "-bin") {
		return nil
	}
	
	for _, val := range vals {
		if hasNotPrintable(val) {
			return fmt.Errorf("header key %q contains value with non-printable ASCII characters", key)
		}
	}
	return nil
}
