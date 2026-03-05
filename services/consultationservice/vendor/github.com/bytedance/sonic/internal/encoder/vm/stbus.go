

package vm

import (
	"unsafe"
	_ "unsafe"

	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/ir"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
)

func EncodeTypedPointer(buf *[]byte, vt *rt.GoType, vp *unsafe.Pointer, sb *vars.Stack, fv uint64) error {
	if vt == nil {
		return alg.EncodeNil(buf)
	} else if pp, err := vars.FindOrCompile(vt, (fv&(1<<alg.BitPointerValue)) != 0, compiler); err != nil {
		return err
	} else if vt.Indirect() {
		return Execute(buf, *vp, sb, fv, pp.(*ir.Program))
	} else {
		return Execute(buf, unsafe.Pointer(vp), sb, fv, pp.(*ir.Program))
	}
}

var compiler func(*rt.GoType, ... interface{}) (interface{}, error)

func SetCompiler(c func(*rt.GoType, ... interface{}) (interface{}, error)) {
	compiler = c
}
