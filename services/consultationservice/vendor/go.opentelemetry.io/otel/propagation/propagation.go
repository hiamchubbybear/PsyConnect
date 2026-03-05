


package propagation 

import (
	"context"
	"net/http"
)


type TextMapCarrier interface {
	
	

	
	Get(key string) string
	
	

	
	Set(key string, value string)
	
	

	
	Keys() []string
	
	
}



type MapCarrier map[string]string


var _ TextMapCarrier = MapCarrier{}


func (c MapCarrier) Get(key string) string {
	return c[key]
}


func (c MapCarrier) Set(key, value string) {
	c[key] = value
}


func (c MapCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}


type HeaderCarrier http.Header


func (hc HeaderCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}


func (hc HeaderCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}


func (hc HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}



type TextMapPropagator interface {
	
	

	
	Inject(ctx context.Context, carrier TextMapCarrier)
	
	

	
	Extract(ctx context.Context, carrier TextMapCarrier) context.Context
	
	

	
	Fields() []string
	
	
}

type compositeTextMapPropagator []TextMapPropagator

func (p compositeTextMapPropagator) Inject(ctx context.Context, carrier TextMapCarrier) {
	for _, i := range p {
		i.Inject(ctx, carrier)
	}
}

func (p compositeTextMapPropagator) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
	for _, i := range p {
		ctx = i.Extract(ctx, carrier)
	}
	return ctx
}

func (p compositeTextMapPropagator) Fields() []string {
	unique := make(map[string]struct{})
	for _, i := range p {
		for _, k := range i.Fields() {
			unique[k] = struct{}{}
		}
	}

	fields := make([]string, 0, len(unique))
	for k := range unique {
		fields = append(fields, k)
	}
	return fields
}









func NewCompositeTextMapPropagator(p ...TextMapPropagator) TextMapPropagator {
	return compositeTextMapPropagator(p)
}
