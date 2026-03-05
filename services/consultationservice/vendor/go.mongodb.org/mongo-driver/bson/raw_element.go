





package bson

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type RawElement []byte



func (re RawElement) Key() string { return bsoncore.Element(re).Key() }


func (re RawElement) KeyErr() (string, error) { return bsoncore.Element(re).KeyErr() }



func (re RawElement) Value() RawValue { return convertFromCoreValue(bsoncore.Element(re).Value()) }


func (re RawElement) ValueErr() (RawValue, error) {
	val, err := bsoncore.Element(re).ValueErr()
	return convertFromCoreValue(val), err
}


func (re RawElement) Validate() error { return bsoncore.Element(re).Validate() }


func (re RawElement) String() string {
	doc := bsoncore.BuildDocument(nil, re)
	j, err := MarshalExtJSON(Raw(doc), true, false)
	if err != nil {
		return "<malformed>"
	}
	return string(j)
}



func (re RawElement) DebugString() string { return bsoncore.Element(re).DebugString() }
