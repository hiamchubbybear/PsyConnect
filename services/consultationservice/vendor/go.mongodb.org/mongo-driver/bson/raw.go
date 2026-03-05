





package bson

import (
	"errors"
	"io"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


var ErrNilReader = errors.New("nil reader")





type Raw []byte



func ReadDocument(r io.Reader) (Raw, error) {
	doc, err := bsoncore.NewDocumentFromReader(r)
	return Raw(doc), err
}





func NewFromIOReader(r io.Reader) (Raw, error) {
	return ReadDocument(r)
}



func (r Raw) Validate() (err error) { return bsoncore.Document(r).Validate() }





func (r Raw) Lookup(key ...string) RawValue {
	return convertFromCoreValue(bsoncore.Document(r).Lookup(key...))
}



func (r Raw) LookupErr(key ...string) (RawValue, error) {
	val, err := bsoncore.Document(r).LookupErr(key...)
	return convertFromCoreValue(val), err
}




func (r Raw) Elements() ([]RawElement, error) {
	doc := bsoncore.Document(r)
	if len(doc) == 0 {
		return nil, nil
	}
	elems, err := doc.Elements()
	if err != nil {
		return nil, err
	}
	relems := make([]RawElement, 0, len(elems))
	for _, elem := range elems {
		relems = append(relems, RawElement(elem))
	}
	return relems, nil
}




func (r Raw) Values() ([]RawValue, error) {
	vals, err := bsoncore.Document(r).Values()
	rvals := make([]RawValue, 0, len(vals))
	for _, val := range vals {
		rvals = append(rvals, convertFromCoreValue(val))
	}
	return rvals, err
}



func (r Raw) Index(index uint) RawElement { return RawElement(bsoncore.Document(r).Index(index)) }


func (r Raw) IndexErr(index uint) (RawElement, error) {
	elem, err := bsoncore.Document(r).IndexErr(index)
	return RawElement(elem), err
}


func (r Raw) String() string { return bsoncore.Document(r).String() }
