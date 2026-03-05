





package mongo

import "context"




type backgroundContext struct {
	context.Context
	childValuesCtx context.Context
}



func newBackgroundContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return &backgroundContext{
		Context:        context.Background(),
		childValuesCtx: ctx,
	}
}

func (b *backgroundContext) Value(key interface{}) interface{} {
	return b.childValuesCtx.Value(key)
}
