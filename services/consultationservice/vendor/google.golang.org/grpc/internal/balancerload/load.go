



package balancerload

import (
	"google.golang.org/grpc/metadata"
)


type Parser interface {
	
	Parse(md metadata.MD) any
}

var parser Parser




func SetParser(lr Parser) {
	parser = lr
}


func Parse(md metadata.MD) any {
	if parser == nil {
		return nil
	}
	return parser.Parse(md)
}
