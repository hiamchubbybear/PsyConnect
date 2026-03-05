

package grpcutil

import (
	"errors"
	"strings"
)



func ParseMethod(methodName string) (service, method string, _ error) {
	if !strings.HasPrefix(methodName, "/") {
		return "", "", errors.New("invalid method name: should start with /")
	}
	methodName = methodName[1:]

	pos := strings.LastIndex(methodName, "/")
	if pos < 0 {
		return "", "", errors.New("invalid method name: suffix /method is missing")
	}
	return methodName[:pos], methodName[pos+1:], nil
}






const baseContentType = "application/grpc"














func ContentSubtype(contentType string) (string, bool) {
	if contentType == baseContentType {
		return "", true
	}
	if !strings.HasPrefix(contentType, baseContentType) {
		return "", false
	}
	
	switch contentType[len(baseContentType)] {
	case '+', ';':
		
		
		
		return contentType[len(baseContentType)+1:], true
	default:
		return "", false
	}
}




func ContentType(contentSubtype string) string {
	if contentSubtype == "" {
		return baseContentType
	}
	return baseContentType + "+" + contentSubtype
}
