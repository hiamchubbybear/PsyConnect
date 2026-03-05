





package bsonutil

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)




func StringSliceFromRawValue(name string, val bson.RawValue) ([]string, error) {
	arr, ok := val.ArrayOK()
	if !ok {
		return nil, fmt.Errorf("expected '%s' to be an array but it's a BSON %s", name, val.Type)
	}

	arrayValues, err := arr.Values()
	if err != nil {
		return nil, err
	}

	strs := make([]string, 0, len(arrayValues))
	for _, arrayVal := range arrayValues {
		str, ok := arrayVal.StringValueOK()
		if !ok {
			return nil, fmt.Errorf("expected '%s' to be an array of strings, but found a BSON %s", name, arrayVal.Type)
		}
		strs = append(strs, str)
	}
	return strs, nil
}


func RawToDocuments(doc bson.Raw) []bson.Raw {
	values, err := doc.Values()
	if err != nil {
		panic(fmt.Sprintf("error converting BSON document to values: %v", err))
	}

	out := make([]bson.Raw, len(values))
	for i := range values {
		out[i] = values[i].Document()
	}

	return out
}


func RawToInterfaces(docs ...bson.Raw) []interface{} {
	out := make([]interface{}, len(docs))
	for i := range docs {
		out[i] = docs[i]
	}
	return out
}
