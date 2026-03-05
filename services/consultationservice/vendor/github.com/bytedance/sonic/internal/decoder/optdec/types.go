

package optdec

import (
	"encoding"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

var (
	boolType                = reflect.TypeOf(bool(false))
	byteType                = reflect.TypeOf(byte(0))
	intType                 = reflect.TypeOf(int(0))
	int8Type                = reflect.TypeOf(int8(0))
	int16Type               = reflect.TypeOf(int16(0))
	int32Type               = reflect.TypeOf(int32(0))
	int64Type               = reflect.TypeOf(int64(0))
	uintType                = reflect.TypeOf(uint(0))
	uint8Type               = reflect.TypeOf(uint8(0))
	uint16Type              = reflect.TypeOf(uint16(0))
	uint32Type              = reflect.TypeOf(uint32(0))
	uint64Type              = reflect.TypeOf(uint64(0))
	float32Type             = reflect.TypeOf(float32(0))
	float64Type             = reflect.TypeOf(float64(0))
	stringType              = reflect.TypeOf("")
	bytesType               = reflect.TypeOf([]byte(nil))
	jsonNumberType          = reflect.TypeOf(json.Number(""))
	base64CorruptInputError = reflect.TypeOf(base64.CorruptInputError(0))
	anyType                 = rt.UnpackType(reflect.TypeOf((*interface{})(nil)).Elem())
)

var (
	errorType                   = reflect.TypeOf((*error)(nil)).Elem()
	jsonUnmarshalerType         = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	encodingTextUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func rtype(t reflect.Type) (*rt.GoItab, *rt.GoType) {
	p := (*rt.GoIface)(unsafe.Pointer(&t))
	return p.Itab, (*rt.GoType)(p.Value)
}
