

package vars

import (
    `encoding/json`
    `fmt`
    `reflect`
    `strconv`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

var ERR_too_deep = &json.UnsupportedValueError {
    Str   : "Value nesting too deep",
    Value : reflect.ValueOf("..."),
}

var ERR_nan_or_infinite = &json.UnsupportedValueError {
    Str   : "NaN or ±Infinite",
    Value : reflect.ValueOf("NaN or ±Infinite"),
}

func Error_type(vtype reflect.Type) error {
    return &json.UnsupportedTypeError{Type: vtype}
}

func Error_number(number json.Number) error {
    return &json.UnsupportedValueError {
        Str   : "invalid number literal: " + strconv.Quote(string(number)),
        Value : reflect.ValueOf(number),
    }
}

func Error_unsuppoted(typ *rt.GoType) error {
	return &json.UnsupportedTypeError{Type: typ.Pack() }
}

func Error_marshaler(ret []byte, pos int) error {
    return fmt.Errorf("invalid Marshaler output json syntax at %d: %q", pos, ret)
}

const (
    PanicNilPointerOfNonEmptyString int = 1 + iota
)

func GoPanic(code int, val unsafe.Pointer) {
    switch(code){
    case PanicNilPointerOfNonEmptyString:
        panic(fmt.Sprintf("val: %#v has nil pointer while its length is not zero!\nThis is a nil pointer exception (NPE) problem. There might be a data race issue. It is recommended to execute the tests related to the code with the `-race` compile flag to detect the problem.", (*rt.GoString)(val)))
    default:
        panic("encoder error!")
    }
}
