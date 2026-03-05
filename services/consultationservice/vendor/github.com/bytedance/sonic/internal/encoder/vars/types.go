

package vars

import (
    `encoding`
    `encoding/json`
    `reflect`
)

var (
    ByteType                 = reflect.TypeOf(byte(0))
    JsonNumberType           = reflect.TypeOf(json.Number(""))
    JsonUnsupportedValueType = reflect.TypeOf(new(json.UnsupportedValueError))
)

var (
    ErrorType                 = reflect.TypeOf((*error)(nil)).Elem()
    JsonMarshalerType         = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
    EncodingTextMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func IsSimpleByte(vt reflect.Type) bool {
    if vt.Kind() != ByteType.Kind() {
        return false
    } else {
        return !isEitherMarshaler(vt) && !isEitherMarshaler(reflect.PtrTo(vt))
    }
}

func isEitherMarshaler(vt reflect.Type) bool {
    return vt.Implements(JsonMarshalerType) || vt.Implements(EncodingTextMarshalerType)
}
