

package x86

import (
	"unsafe"
	_ "unsafe"

	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/loader"
	_ "github.com/cloudwego/base64x"
)

var compiler func(*rt.GoType, ... interface{}) (interface{}, error)

func SetCompiler(c func(*rt.GoType, ... interface{}) (interface{}, error)) {
	compiler = c
}

func ptoenc(p loader.Function) vars.Encoder {
    return *(*vars.Encoder)(unsafe.Pointer(&p))
}

func EncodeTypedPointer(buf *[]byte, vt *rt.GoType, vp *unsafe.Pointer, sb *vars.Stack, fv uint64) error {
	if vt == nil {
		return alg.EncodeNil(buf)
	} else if fn, err := vars.FindOrCompile(vt, (fv&(1<<alg.BitPointerValue)) != 0, compiler); err != nil {
		return err
	} else if vt.Indirect() {
		return	fn.(vars.Encoder)(buf, *vp, sb, fv)
	} else {
		return fn.(vars.Encoder)(buf, unsafe.Pointer(vp), sb, fv)
	}
}

