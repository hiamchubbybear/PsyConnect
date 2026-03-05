





package bson

import (
	"bytes"
	"encoding/json"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultDstCap = 256

var bvwPool = bsonrw.NewBSONValueWriterPool()
var extjPool = bsonrw.NewExtJSONValueWriterPool()







type Marshaler interface {
	MarshalBSON() ([]byte, error)
}








type ValueMarshaler interface {
	MarshalBSONValue() (bsontype.Type, []byte, error)
}







func Marshal(val interface{}) ([]byte, error) {
	return MarshalWithRegistry(DefaultRegistry, val)
}



















func MarshalAppend(dst []byte, val interface{}) ([]byte, error) {
	return MarshalAppendWithRegistry(DefaultRegistry, dst, val)
}


















func MarshalWithRegistry(r *bsoncodec.Registry, val interface{}) ([]byte, error) {
	dst := make([]byte, 0)
	return MarshalAppendWithRegistry(r, dst, val)
}



















func MarshalWithContext(ec bsoncodec.EncodeContext, val interface{}) ([]byte, error) {
	dst := make([]byte, 0)
	return MarshalAppendWithContext(ec, dst, val)
}




















func MarshalAppendWithRegistry(r *bsoncodec.Registry, dst []byte, val interface{}) ([]byte, error) {
	return MarshalAppendWithContext(bsoncodec.EncodeContext{Registry: r}, dst, val)
}


var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}





















func MarshalAppendWithContext(ec bsoncodec.EncodeContext, dst []byte, val interface{}) ([]byte, error) {
	sw := bufPool.Get().(*bytes.Buffer)
	defer func() {
		
		
		
		
		
		
		
		
		
		
		
		
		if sw.Cap() < 16*1024*1024 && sw.Cap()/2 < sw.Len() {
			bufPool.Put(sw)
		}
	}()

	sw.Reset()
	vw := bvwPool.Get(sw)
	defer bvwPool.Put(vw)

	enc := encPool.Get().(*Encoder)
	defer encPool.Put(enc)

	err := enc.Reset(vw)
	if err != nil {
		return nil, err
	}
	err = enc.SetContext(ec)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(val)
	if err != nil {
		return nil, err
	}

	return append(dst, sw.Bytes()...), nil
}





func MarshalValue(val interface{}) (bsontype.Type, []byte, error) {
	return MarshalValueWithRegistry(DefaultRegistry, val)
}






func MarshalValueAppend(dst []byte, val interface{}) (bsontype.Type, []byte, error) {
	return MarshalValueAppendWithRegistry(DefaultRegistry, dst, val)
}





func MarshalValueWithRegistry(r *bsoncodec.Registry, val interface{}) (bsontype.Type, []byte, error) {
	dst := make([]byte, 0)
	return MarshalValueAppendWithRegistry(r, dst, val)
}





func MarshalValueWithContext(ec bsoncodec.EncodeContext, val interface{}) (bsontype.Type, []byte, error) {
	dst := make([]byte, 0)
	return MarshalValueAppendWithContext(ec, dst, val)
}






func MarshalValueAppendWithRegistry(r *bsoncodec.Registry, dst []byte, val interface{}) (bsontype.Type, []byte, error) {
	return MarshalValueAppendWithContext(bsoncodec.EncodeContext{Registry: r}, dst, val)
}






func MarshalValueAppendWithContext(ec bsoncodec.EncodeContext, dst []byte, val interface{}) (bsontype.Type, []byte, error) {
	
	sw := new(bsonrw.SliceWriter)
	*sw = dst
	vwFlusher := bvwPool.GetAtModeElement(sw)

	
	enc := encPool.Get().(*Encoder)
	defer encPool.Put(enc)
	if err := enc.Reset(vwFlusher); err != nil {
		return 0, nil, err
	}
	if err := enc.SetContext(ec); err != nil {
		return 0, nil, err
	}
	if err := enc.Encode(val); err != nil {
		return 0, nil, err
	}

	
	
	
	if err := vwFlusher.Flush(); err != nil {
		return 0, nil, err
	}
	buffer := *sw
	return bsontype.Type(buffer[0]), buffer[2:], nil
}


func MarshalExtJSON(val interface{}, canonical, escapeHTML bool) ([]byte, error) {
	return MarshalExtJSONWithRegistry(DefaultRegistry, val, canonical, escapeHTML)
}



















func MarshalExtJSONAppend(dst []byte, val interface{}, canonical, escapeHTML bool) ([]byte, error) {
	return MarshalExtJSONAppendWithRegistry(DefaultRegistry, dst, val, canonical, escapeHTML)
}

















func MarshalExtJSONWithRegistry(r *bsoncodec.Registry, val interface{}, canonical, escapeHTML bool) ([]byte, error) {
	dst := make([]byte, 0, defaultDstCap)
	return MarshalExtJSONAppendWithContext(bsoncodec.EncodeContext{Registry: r}, dst, val, canonical, escapeHTML)
}


















func MarshalExtJSONWithContext(ec bsoncodec.EncodeContext, val interface{}, canonical, escapeHTML bool) ([]byte, error) {
	dst := make([]byte, 0, defaultDstCap)
	return MarshalExtJSONAppendWithContext(ec, dst, val, canonical, escapeHTML)
}




















func MarshalExtJSONAppendWithRegistry(r *bsoncodec.Registry, dst []byte, val interface{}, canonical, escapeHTML bool) ([]byte, error) {
	return MarshalExtJSONAppendWithContext(bsoncodec.EncodeContext{Registry: r}, dst, val, canonical, escapeHTML)
}





















func MarshalExtJSONAppendWithContext(ec bsoncodec.EncodeContext, dst []byte, val interface{}, canonical, escapeHTML bool) ([]byte, error) {
	sw := new(bsonrw.SliceWriter)
	*sw = dst
	ejvw := extjPool.Get(sw, canonical, escapeHTML)
	defer extjPool.Put(ejvw)

	enc := encPool.Get().(*Encoder)
	defer encPool.Put(enc)

	err := enc.Reset(ejvw)
	if err != nil {
		return nil, err
	}
	err = enc.SetContext(ec)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(val)
	if err != nil {
		return nil, err
	}

	return *sw, nil
}


func IndentExtJSON(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return json.Indent(dst, src, prefix, indent)
}



func MarshalExtJSONIndent(val interface{}, canonical, escapeHTML bool, prefix, indent string) ([]byte, error) {
	marshaled, err := MarshalExtJSON(val, canonical, escapeHTML)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = IndentExtJSON(&buf, marshaled, prefix, indent)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
