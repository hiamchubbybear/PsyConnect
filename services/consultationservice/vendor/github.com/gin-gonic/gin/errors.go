



package gin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/internal/json"
)


type ErrorType uint64

const (
	
	ErrorTypeBind ErrorType = 1 << 63
	
	ErrorTypeRender ErrorType = 1 << 62
	
	ErrorTypePrivate ErrorType = 1 << 0
	
	ErrorTypePublic ErrorType = 1 << 1
	
	ErrorTypeAny ErrorType = 1<<64 - 1
	
	ErrorTypeNu = 2
)


type Error struct {
	Err  error
	Type ErrorType
	Meta any
}

type errorMsgs []*Error

var _ error = (*Error)(nil)


func (msg *Error) SetType(flags ErrorType) *Error {
	msg.Type = flags
	return msg
}


func (msg *Error) SetMeta(data any) *Error {
	msg.Meta = data
	return msg
}


func (msg *Error) JSON() any {
	jsonData := H{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				jsonData[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			jsonData["meta"] = msg.Meta
		}
	}
	if _, ok := jsonData["error"]; !ok {
		jsonData["error"] = msg.Error()
	}
	return jsonData
}


func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}


func (msg Error) Error() string {
	return msg.Err.Error()
}


func (msg *Error) IsType(flags ErrorType) bool {
	return (msg.Type & flags) > 0
}


func (msg *Error) Unwrap() error {
	return msg.Err
}



func (a errorMsgs) ByType(typ ErrorType) errorMsgs {
	if len(a) == 0 {
		return nil
	}
	if typ == ErrorTypeAny {
		return a
	}
	var result errorMsgs
	for _, msg := range a {
		if msg.IsType(typ) {
			result = append(result, msg)
		}
	}
	return result
}



func (a errorMsgs) Last() *Error {
	if length := len(a); length > 0 {
		return a[length-1]
	}
	return nil
}








func (a errorMsgs) Errors() []string {
	if len(a) == 0 {
		return nil
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

func (a errorMsgs) JSON() any {
	switch length := len(a); length {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		jsonData := make([]any, length)
		for i, err := range a {
			jsonData[i] = err.JSON()
		}
		return jsonData
	}
}


func (a errorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
