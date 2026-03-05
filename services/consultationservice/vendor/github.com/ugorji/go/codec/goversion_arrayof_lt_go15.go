


//go:build !go1.5
// +build !go1.5

package codec

import (
	"errors"
	"reflect"
)

const reflectArrayOfSupported = false

var errNoReflectArrayOf = errors.New("codec: reflect.ArrayOf unsupported by this go version")

func reflectArrayOf(count int, elem reflect.Type) reflect.Type {
	panic(errNoReflectArrayOf)
}
