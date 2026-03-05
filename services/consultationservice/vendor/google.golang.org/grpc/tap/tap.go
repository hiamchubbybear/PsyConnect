








package tap

import (
	"context"

	"google.golang.org/grpc/metadata"
)


type Info struct {
	
	
	FullMethodName string

	
	Header metadata.MD

	
}

















type ServerInHandle func(ctx context.Context, info *Info) (context.Context, error)
