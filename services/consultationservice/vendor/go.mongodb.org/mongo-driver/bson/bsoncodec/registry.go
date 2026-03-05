





package bsoncodec

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)




var ErrNilType = errors.New("cannot perform a decoder lookup on <nil>")




var ErrNotPointer = errors.New("non-pointer provided to LookupDecoder")




type ErrNoEncoder struct {
	Type reflect.Type
}

func (ene ErrNoEncoder) Error() string {
	if ene.Type == nil {
		return "no encoder found for <nil>"
	}
	return "no encoder found for " + ene.Type.String()
}




type ErrNoDecoder struct {
	Type reflect.Type
}

func (end ErrNoDecoder) Error() string {
	return "no decoder found for " + end.Type.String()
}




type ErrNoTypeMapEntry struct {
	Type bsontype.Type
}

func (entme ErrNoTypeMapEntry) Error() string {
	return "no type map entry found for " + entme.Type.String()
}




var ErrNotInterface = errors.New("The provided type is not an interface")





type RegistryBuilder struct {
	registry *Registry
}




func NewRegistryBuilder() *RegistryBuilder {
	return &RegistryBuilder{
		registry: NewRegistry(),
	}
}




func (rb *RegistryBuilder) RegisterCodec(t reflect.Type, codec ValueCodec) *RegistryBuilder {
	rb.RegisterTypeEncoder(t, codec)
	rb.RegisterTypeDecoder(t, codec)
	return rb
}










func (rb *RegistryBuilder) RegisterTypeEncoder(t reflect.Type, enc ValueEncoder) *RegistryBuilder {
	rb.registry.RegisterTypeEncoder(t, enc)
	return rb
}






func (rb *RegistryBuilder) RegisterHookEncoder(t reflect.Type, enc ValueEncoder) *RegistryBuilder {
	rb.registry.RegisterInterfaceEncoder(t, enc)
	return rb
}










func (rb *RegistryBuilder) RegisterTypeDecoder(t reflect.Type, dec ValueDecoder) *RegistryBuilder {
	rb.registry.RegisterTypeDecoder(t, dec)
	return rb
}






func (rb *RegistryBuilder) RegisterHookDecoder(t reflect.Type, dec ValueDecoder) *RegistryBuilder {
	rb.registry.RegisterInterfaceDecoder(t, dec)
	return rb
}




func (rb *RegistryBuilder) RegisterEncoder(t reflect.Type, enc ValueEncoder) *RegistryBuilder {
	if t == tEmpty {
		rb.registry.RegisterTypeEncoder(t, enc)
		return rb
	}
	switch t.Kind() {
	case reflect.Interface:
		rb.registry.RegisterInterfaceEncoder(t, enc)
	default:
		rb.registry.RegisterTypeEncoder(t, enc)
	}
	return rb
}




func (rb *RegistryBuilder) RegisterDecoder(t reflect.Type, dec ValueDecoder) *RegistryBuilder {
	if t == nil {
		rb.registry.RegisterTypeDecoder(t, dec)
		return rb
	}
	if t == tEmpty {
		rb.registry.RegisterTypeDecoder(t, dec)
		return rb
	}
	switch t.Kind() {
	case reflect.Interface:
		rb.registry.RegisterInterfaceDecoder(t, dec)
	default:
		rb.registry.RegisterTypeDecoder(t, dec)
	}
	return rb
}





func (rb *RegistryBuilder) RegisterDefaultEncoder(kind reflect.Kind, enc ValueEncoder) *RegistryBuilder {
	rb.registry.RegisterKindEncoder(kind, enc)
	return rb
}





func (rb *RegistryBuilder) RegisterDefaultDecoder(kind reflect.Kind, dec ValueDecoder) *RegistryBuilder {
	rb.registry.RegisterKindDecoder(kind, dec)
	return rb
}












func (rb *RegistryBuilder) RegisterTypeMapEntry(bt bsontype.Type, rt reflect.Type) *RegistryBuilder {
	rb.registry.RegisterTypeMapEntry(bt, rt)
	return rb
}




func (rb *RegistryBuilder) Build() *Registry {
	r := &Registry{
		interfaceEncoders: append([]interfaceValueEncoder(nil), rb.registry.interfaceEncoders...),
		interfaceDecoders: append([]interfaceValueDecoder(nil), rb.registry.interfaceDecoders...),
		typeEncoders:      rb.registry.typeEncoders.Clone(),
		typeDecoders:      rb.registry.typeDecoders.Clone(),
		kindEncoders:      rb.registry.kindEncoders.Clone(),
		kindDecoders:      rb.registry.kindDecoders.Clone(),
	}
	rb.registry.typeMap.Range(func(k, v interface{}) bool {
		if k != nil && v != nil {
			r.typeMap.Store(k, v)
		}
		return true
	})
	return r
}



type Registry struct {
	interfaceEncoders []interfaceValueEncoder
	interfaceDecoders []interfaceValueDecoder
	typeEncoders      *typeEncoderCache
	typeDecoders      *typeDecoderCache
	kindEncoders      *kindEncoderCache
	kindDecoders      *kindDecoderCache
	typeMap           sync.Map 
}


func NewRegistry() *Registry {
	return &Registry{
		typeEncoders: new(typeEncoderCache),
		typeDecoders: new(typeDecoderCache),
		kindEncoders: new(kindEncoderCache),
		kindDecoders: new(kindDecoderCache),
	}
}











func (r *Registry) RegisterTypeEncoder(valueType reflect.Type, enc ValueEncoder) {
	r.typeEncoders.Store(valueType, enc)
}











func (r *Registry) RegisterTypeDecoder(valueType reflect.Type, dec ValueDecoder) {
	r.typeDecoders.Store(valueType, dec)
}













func (r *Registry) RegisterKindEncoder(kind reflect.Kind, enc ValueEncoder) {
	r.kindEncoders.Store(kind, enc)
}













func (r *Registry) RegisterKindDecoder(kind reflect.Kind, dec ValueDecoder) {
	r.kindDecoders.Store(kind, dec)
}







func (r *Registry) RegisterInterfaceEncoder(iface reflect.Type, enc ValueEncoder) {
	if iface.Kind() != reflect.Interface {
		panicStr := fmt.Errorf("RegisterInterfaceEncoder expects a type with kind reflect.Interface, "+
			"got type %s with kind %s", iface, iface.Kind())
		panic(panicStr)
	}

	for idx, encoder := range r.interfaceEncoders {
		if encoder.i == iface {
			r.interfaceEncoders[idx].ve = enc
			return
		}
	}

	r.interfaceEncoders = append(r.interfaceEncoders, interfaceValueEncoder{i: iface, ve: enc})
}







func (r *Registry) RegisterInterfaceDecoder(iface reflect.Type, dec ValueDecoder) {
	if iface.Kind() != reflect.Interface {
		panicStr := fmt.Errorf("RegisterInterfaceDecoder expects a type with kind reflect.Interface, "+
			"got type %s with kind %s", iface, iface.Kind())
		panic(panicStr)
	}

	for idx, decoder := range r.interfaceDecoders {
		if decoder.i == iface {
			r.interfaceDecoders[idx].vd = dec
			return
		}
	}

	r.interfaceDecoders = append(r.interfaceDecoders, interfaceValueDecoder{i: iface, vd: dec})
}










func (r *Registry) RegisterTypeMapEntry(bt bsontype.Type, rt reflect.Type) {
	r.typeMap.Store(bt, rt)
}














func (r *Registry) LookupEncoder(valueType reflect.Type) (ValueEncoder, error) {
	if valueType == nil {
		return nil, ErrNoEncoder{Type: valueType}
	}
	enc, found := r.lookupTypeEncoder(valueType)
	if found {
		if enc == nil {
			return nil, ErrNoEncoder{Type: valueType}
		}
		return enc, nil
	}

	enc, found = r.lookupInterfaceEncoder(valueType, true)
	if found {
		return r.typeEncoders.LoadOrStore(valueType, enc), nil
	}

	if v, ok := r.kindEncoders.Load(valueType.Kind()); ok {
		return r.storeTypeEncoder(valueType, v), nil
	}
	return nil, ErrNoEncoder{Type: valueType}
}

func (r *Registry) storeTypeEncoder(rt reflect.Type, enc ValueEncoder) ValueEncoder {
	return r.typeEncoders.LoadOrStore(rt, enc)
}

func (r *Registry) lookupTypeEncoder(rt reflect.Type) (ValueEncoder, bool) {
	return r.typeEncoders.Load(rt)
}

func (r *Registry) lookupInterfaceEncoder(valueType reflect.Type, allowAddr bool) (ValueEncoder, bool) {
	if valueType == nil {
		return nil, false
	}
	for _, ienc := range r.interfaceEncoders {
		if valueType.Implements(ienc.i) {
			return ienc.ve, true
		}
		if allowAddr && valueType.Kind() != reflect.Ptr && reflect.PtrTo(valueType).Implements(ienc.i) {
			
			
			defaultEnc, found := r.lookupInterfaceEncoder(valueType, false)
			if !found {
				defaultEnc, _ = r.kindEncoders.Load(valueType.Kind())
			}
			return newCondAddrEncoder(ienc.ve, defaultEnc), true
		}
	}
	return nil, false
}














func (r *Registry) LookupDecoder(valueType reflect.Type) (ValueDecoder, error) {
	if valueType == nil {
		return nil, ErrNilType
	}
	dec, found := r.lookupTypeDecoder(valueType)
	if found {
		if dec == nil {
			return nil, ErrNoDecoder{Type: valueType}
		}
		return dec, nil
	}

	dec, found = r.lookupInterfaceDecoder(valueType, true)
	if found {
		return r.storeTypeDecoder(valueType, dec), nil
	}

	if v, ok := r.kindDecoders.Load(valueType.Kind()); ok {
		return r.storeTypeDecoder(valueType, v), nil
	}
	return nil, ErrNoDecoder{Type: valueType}
}

func (r *Registry) lookupTypeDecoder(valueType reflect.Type) (ValueDecoder, bool) {
	return r.typeDecoders.Load(valueType)
}

func (r *Registry) storeTypeDecoder(typ reflect.Type, dec ValueDecoder) ValueDecoder {
	return r.typeDecoders.LoadOrStore(typ, dec)
}

func (r *Registry) lookupInterfaceDecoder(valueType reflect.Type, allowAddr bool) (ValueDecoder, bool) {
	for _, idec := range r.interfaceDecoders {
		if valueType.Implements(idec.i) {
			return idec.vd, true
		}
		if allowAddr && valueType.Kind() != reflect.Ptr && reflect.PtrTo(valueType).Implements(idec.i) {
			
			
			defaultDec, found := r.lookupInterfaceDecoder(valueType, false)
			if !found {
				defaultDec, _ = r.kindDecoders.Load(valueType.Kind())
			}
			return newCondAddrDecoder(idec.vd, defaultDec), true
		}
	}
	return nil, false
}





func (r *Registry) LookupTypeMapEntry(bt bsontype.Type) (reflect.Type, error) {
	v, ok := r.typeMap.Load(bt)
	if v == nil || !ok {
		return nil, ErrNoTypeMapEntry{Type: bt}
	}
	return v.(reflect.Type), nil
}

type interfaceValueEncoder struct {
	i  reflect.Type
	ve ValueEncoder
}

type interfaceValueDecoder struct {
	i  reflect.Type
	vd ValueDecoder
}
