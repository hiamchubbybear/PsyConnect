


//go:build go1.6 && !go1.7
// +build go1.6,!go1.7

package codec

import "os"

var genCheckVendor = os.Getenv("GO15VENDOREXPERIMENT") != "0"
