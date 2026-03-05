


//go:build go1.5 && !go1.6
// +build go1.5,!go1.6

package codec

import "os"

var genCheckVendor = os.Getenv("GO15VENDOREXPERIMENT") == "1"
