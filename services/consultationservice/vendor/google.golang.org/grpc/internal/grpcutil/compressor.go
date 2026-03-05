

package grpcutil

import (
	"strings"
)


var RegisteredCompressorNames []string


func IsCompressorNameRegistered(name string) bool {
	for _, compressor := range RegisteredCompressorNames {
		if compressor == name {
			return true
		}
	}
	return false
}



func RegisteredCompressors() string {
	return strings.Join(RegisteredCompressorNames, ",")
}
