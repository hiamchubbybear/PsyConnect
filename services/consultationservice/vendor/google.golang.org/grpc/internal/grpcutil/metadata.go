

package grpcutil

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type mdExtraKey struct{}


func WithExtraMetadata(ctx context.Context, md metadata.MD) context.Context {
	return context.WithValue(ctx, mdExtraKey{}, md)
}




func ExtraMetadata(ctx context.Context) (md metadata.MD, ok bool) {
	md, ok = ctx.Value(mdExtraKey{}).(metadata.MD)
	return
}
