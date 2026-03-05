





package bsoncodec

import (
	"reflect"
	"strings"
)




type StructTagParser interface {
	ParseStructTags(reflect.StructField) (StructTags, error)
}





type StructTagParserFunc func(reflect.StructField) (StructTags, error)


func (stpf StructTagParserFunc) ParseStructTags(sf reflect.StructField) (StructTags, error) {
	return stpf(sf)
}



























type StructTags struct {
	Name      string
	OmitEmpty bool
	MinSize   bool
	Truncate  bool
	Inline    bool
	Skip      bool
}




























var DefaultStructTagParser StructTagParserFunc = func(sf reflect.StructField) (StructTags, error) {
	key := strings.ToLower(sf.Name)
	tag, ok := sf.Tag.Lookup("bson")
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		tag = string(sf.Tag)
	}
	return parseTags(key, tag)
}

func parseTags(key string, tag string) (StructTags, error) {
	var st StructTags
	if tag == "-" {
		st.Skip = true
		return st, nil
	}

	for idx, str := range strings.Split(tag, ",") {
		if idx == 0 && str != "" {
			key = str
		}
		switch str {
		case "omitempty":
			st.OmitEmpty = true
		case "minsize":
			st.MinSize = true
		case "truncate":
			st.Truncate = true
		case "inline":
			st.Inline = true
		}
	}

	st.Name = key

	return st, nil
}







var JSONFallbackStructTagParser StructTagParserFunc = func(sf reflect.StructField) (StructTags, error) {
	key := strings.ToLower(sf.Name)
	tag, ok := sf.Tag.Lookup("bson")
	if !ok {
		tag, ok = sf.Tag.Lookup("json")
	}
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		tag = string(sf.Tag)
	}

	return parseTags(key, tag)
}
