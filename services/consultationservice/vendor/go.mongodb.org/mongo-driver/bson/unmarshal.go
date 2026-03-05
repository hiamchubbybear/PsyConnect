





package bson

import (
	"bytes"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)









type Unmarshaler interface {
	UnmarshalBSON([]byte) error
}









type ValueUnmarshaler interface {
	UnmarshalBSONValue(bsontype.Type, []byte) error
}







func Unmarshal(data []byte, val interface{}) error {
	return UnmarshalWithRegistry(DefaultRegistry, data, val)
}














func UnmarshalWithRegistry(r *bsoncodec.Registry, data []byte, val interface{}) error {
	vr := bsonrw.NewBSONDocumentReader(data)
	return unmarshalFromReader(bsoncodec.DecodeContext{Registry: r}, vr, val)
}















func UnmarshalWithContext(dc bsoncodec.DecodeContext, data []byte, val interface{}) error {
	vr := bsonrw.NewBSONDocumentReader(data)
	return unmarshalFromReader(dc, vr, val)
}




func UnmarshalValue(t bsontype.Type, data []byte, val interface{}) error {
	return UnmarshalValueWithRegistry(DefaultRegistry, t, data, val)
}







func UnmarshalValueWithRegistry(r *bsoncodec.Registry, t bsontype.Type, data []byte, val interface{}) error {
	vr := bsonrw.NewBSONValueReader(t, data)
	return unmarshalFromReader(bsoncodec.DecodeContext{Registry: r}, vr, val)
}




func UnmarshalExtJSON(data []byte, canonical bool, val interface{}) error {
	return UnmarshalExtJSONWithRegistry(DefaultRegistry, data, canonical, val)
}


















func UnmarshalExtJSONWithRegistry(r *bsoncodec.Registry, data []byte, canonical bool, val interface{}) error {
	ejvr, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(data), canonical)
	if err != nil {
		return err
	}

	return unmarshalFromReader(bsoncodec.DecodeContext{Registry: r}, ejvr, val)
}



















func UnmarshalExtJSONWithContext(dc bsoncodec.DecodeContext, data []byte, canonical bool, val interface{}) error {
	ejvr, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(data), canonical)
	if err != nil {
		return err
	}

	return unmarshalFromReader(dc, ejvr, val)
}

func unmarshalFromReader(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val interface{}) error {
	dec := decPool.Get().(*Decoder)
	defer decPool.Put(dec)

	err := dec.Reset(vr)
	if err != nil {
		return err
	}
	err = dec.SetContext(dc)
	if err != nil {
		return err
	}

	return dec.Decode(val)
}
