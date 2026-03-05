

package alg

import (
	"encoding"
	"encoding/json"
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/resolver"
	"github.com/bytedance/sonic/internal/rt"
)

func Compact(p *[]byte, v []byte) error {
	buf := vars.NewBuffer()
	err := json.Compact(buf, v)

	
	if err != nil {
		return err
	}

	
	v = buf.Bytes()
	*p = append(*p, v...)

	
	vars.FreeBuffer(buf)
	return nil
}

func EncodeNil(rb *[]byte) error {
	*rb = append(*rb, 'n', 'u', 'l', 'l')
	return nil
}

















func EncodeJsonMarshaler(buf *[]byte, val json.Marshaler, opt uint64) error {
	if ret, err := val.MarshalJSON(); err != nil {
		return err
	} else {
		if opt&(1<<BitCompactMarshaler) != 0 {
			return Compact(buf, ret)
		}
		if opt&(1<<BitNoValidateJSONMarshaler) == 0 {
			if ok, s := Valid(ret); !ok {
				return vars.Error_marshaler(ret, s)
			}
		}
		*buf = append(*buf, ret...)
		return nil
	}
}

func EncodeTextMarshaler(buf *[]byte, val encoding.TextMarshaler, opt uint64) error {
	if ret, err := val.MarshalText(); err != nil {
		return err
	} else {
		if opt&(1<<BitNoQuoteTextMarshaler) != 0 {
			*buf = append(*buf, ret...)
			return nil
		}
		*buf = Quote(*buf, rt.Mem2Str(ret), false)
		return nil
	}
}

func IsZero(val unsafe.Pointer, fv *resolver.FieldMeta) bool {
	rv := reflect.NewAt(fv.Type, val).Elem()
	b1 := fv.IsZero == nil && rv.IsZero()
	b2 := fv.IsZero != nil && fv.IsZero(rv)
	return  b1 || b2
}
