


//go:build go1.5
// +build go1.5

package codec

import "reflect"

const reflectArrayOfSupported = true

func reflectArrayOf(count int, elem reflect.Type) reflect.Type {
	return reflect.ArrayOf(count, elem)
}
