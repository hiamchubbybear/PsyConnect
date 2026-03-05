





package options

import (
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)



type Collation struct {
	Locale          string `bson:",omitempty"` 
	CaseLevel       bool   `bson:",omitempty"` 
	CaseFirst       string `bson:",omitempty"` 
	Strength        int    `bson:",omitempty"` 
	NumericOrdering bool   `bson:",omitempty"` 
	Alternate       string `bson:",omitempty"` 
	MaxVariable     string `bson:",omitempty"` 
	Normalization   bool   `bson:",omitempty"` 
	Backwards       bool   `bson:",omitempty"` 
}




func (co *Collation) ToDocument() bson.Raw {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	if co.Locale != "" {
		doc = bsoncore.AppendStringElement(doc, "locale", co.Locale)
	}
	if co.CaseLevel {
		doc = bsoncore.AppendBooleanElement(doc, "caseLevel", true)
	}
	if co.CaseFirst != "" {
		doc = bsoncore.AppendStringElement(doc, "caseFirst", co.CaseFirst)
	}
	if co.Strength != 0 {
		doc = bsoncore.AppendInt32Element(doc, "strength", int32(co.Strength))
	}
	if co.NumericOrdering {
		doc = bsoncore.AppendBooleanElement(doc, "numericOrdering", true)
	}
	if co.Alternate != "" {
		doc = bsoncore.AppendStringElement(doc, "alternate", co.Alternate)
	}
	if co.MaxVariable != "" {
		doc = bsoncore.AppendStringElement(doc, "maxVariable", co.MaxVariable)
	}
	if co.Normalization {
		doc = bsoncore.AppendBooleanElement(doc, "normalization", true)
	}
	if co.Backwards {
		doc = bsoncore.AppendBooleanElement(doc, "backwards", true)
	}
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc
}



type CursorType int8

const (
	
	NonTailable CursorType = iota
	
	Tailable
	
	
	TailableAwait
)



type ReturnDocument int8

const (
	
	Before ReturnDocument = iota
	
	After
)


type FullDocument string

const (
	
	Default FullDocument = "default"
	
	Off FullDocument = "off"
	
	Required FullDocument = "required"
	
	
	UpdateLookup FullDocument = "updateLookup"
	
	
	WhenAvailable FullDocument = "whenAvailable"
)







type ArrayFilters struct {
	
	
	
	Registry *bsoncodec.Registry

	Filters []interface{} 
}




func (af *ArrayFilters) ToArray() ([]bson.Raw, error) {
	registry := af.Registry
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	filters := make([]bson.Raw, 0, len(af.Filters))
	for _, f := range af.Filters {
		filter, err := bson.MarshalWithRegistry(registry, f)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}





func (af *ArrayFilters) ToArrayDocument() (bson.Raw, error) {
	registry := af.Registry
	if registry == nil {
		registry = bson.DefaultRegistry
	}

	idx, arr := bsoncore.AppendArrayStart(nil)
	for i, f := range af.Filters {
		filter, err := bson.MarshalWithRegistry(registry, f)
		if err != nil {
			return nil, err
		}

		arr = bsoncore.AppendDocumentElement(arr, strconv.Itoa(i), filter)
	}
	arr, _ = bsoncore.AppendArrayEnd(arr, idx)
	return arr, nil
}





type MarshalError struct {
	Value interface{}
	Err   error
}




func (me MarshalError) Error() string {
	return fmt.Sprintf("cannot transform type %s to a bson.Raw", reflect.TypeOf(me.Value))
}
