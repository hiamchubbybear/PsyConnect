


//go:build !go1.4
// +build !go1.4

package codec

import "errors"








var errCodecSupportedOnlyFromGo14 = errors.New("codec: go 1.3 and below are not supported")

func init() {
	panic(errCodecSupportedOnlyFromGo14)
}
