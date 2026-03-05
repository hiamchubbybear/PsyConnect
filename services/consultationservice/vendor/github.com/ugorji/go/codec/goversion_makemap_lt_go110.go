


//go:build !go1.10
// +build !go1.10

package codec

import "reflect"

func makeMapReflect(t reflect.Type, size int) reflect.Value {
	return reflect.MakeMap(t)
}
