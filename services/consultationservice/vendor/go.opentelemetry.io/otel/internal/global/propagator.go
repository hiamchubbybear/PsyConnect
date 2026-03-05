


package global 

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/propagation"
)




type textMapPropagator struct {
	mtx      sync.Mutex
	once     sync.Once
	delegate propagation.TextMapPropagator
	noop     propagation.TextMapPropagator
}



var _ propagation.TextMapPropagator = (*textMapPropagator)(nil)

func newTextMapPropagator() *textMapPropagator {
	return &textMapPropagator{
		noop: propagation.NewCompositeTextMapPropagator(),
	}
}




func (p *textMapPropagator) SetDelegate(delegate propagation.TextMapPropagator) {
	if delegate == nil {
		return
	}

	p.mtx.Lock()
	p.once.Do(func() { p.delegate = delegate })
	p.mtx.Unlock()
}




func (p *textMapPropagator) effectiveDelegate() propagation.TextMapPropagator {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.delegate != nil {
		return p.delegate
	}
	return p.noop
}


func (p *textMapPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	p.effectiveDelegate().Inject(ctx, carrier)
}


func (p *textMapPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return p.effectiveDelegate().Extract(ctx, carrier)
}


func (p *textMapPropagator) Fields() []string {
	return p.effectiveDelegate().Fields()
}
