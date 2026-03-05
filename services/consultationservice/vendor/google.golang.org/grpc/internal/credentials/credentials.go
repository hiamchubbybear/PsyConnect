

package credentials

import (
	"context"
)



type requestInfoKey struct{}


func NewRequestInfoContext(ctx context.Context, ri any) context.Context {
	return context.WithValue(ctx, requestInfoKey{}, ri)
}


func RequestInfoFromContext(ctx context.Context) any {
	return ctx.Value(requestInfoKey{})
}



type clientHandshakeInfoKey struct{}


func ClientHandshakeInfoFromContext(ctx context.Context) any {
	return ctx.Value(clientHandshakeInfoKey{})
}


func NewClientHandshakeInfoContext(ctx context.Context, chi any) context.Context {
	return context.WithValue(ctx, clientHandshakeInfoKey{}, chi)
}
