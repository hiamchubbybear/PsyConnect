package sasl

import "context"

type ctxKey struct{}







type Mechanism interface {
	
	
	
	Name() string

	
	
	
	
	
	
	
	
	
	Start(ctx context.Context) (sess StateMachine, ir []byte, err error)
}










type StateMachine interface {
	
	
	
	
	Next(ctx context.Context, challenge []byte) (done bool, response []byte, err error)
}


type Metadata struct {
	
	
	Host string
	Port int
}


func WithMetadata(ctx context.Context, m *Metadata) context.Context {
	return context.WithValue(ctx, ctxKey{}, m)
}


func MetadataFromContext(ctx context.Context) *Metadata {
	m, _ := ctx.Value(ctxKey{}).(*Metadata)
	return m
}
