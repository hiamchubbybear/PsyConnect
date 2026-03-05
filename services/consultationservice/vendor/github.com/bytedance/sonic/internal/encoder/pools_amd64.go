

package encoder

import (
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/encoder/x86"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/option"
)


func ForceUseJit() {
	x86.SetCompiler(makeEncoderX86)
	pretouchType = pretouchTypeX86
	encodeTypedPointer = x86.EncodeTypedPointer
	vars.UseVM = false
}

func init() {
	if vars.UseVM {
		ForceUseVM()
	} else {
		ForceUseJit()
	}
}

var _KeepAlive struct {
	rb    *[]byte
	vp    unsafe.Pointer
	sb    *vars.Stack
	fv    uint64
	err   error
	frame [x86.FP_offs]byte
}

func makeEncoderX86(vt *rt.GoType, ex ...interface{}) (interface{}, error) {
	pp, err := NewCompiler().Compile(vt.Pack(), ex[0].(bool))
	if err != nil {
		return nil, err
	}
	as := x86.NewAssembler(pp)
	as.Name = vt.String()
	return as.Load(), nil
}

func pretouchTypeX86(_vt reflect.Type, opts option.CompileOptions, v uint8) (map[reflect.Type]uint8, error) {
	
	compiler := NewCompiler().apply(opts)

	
	vt := rt.UnpackType(_vt)
	if val := vars.GetProgram(vt); val != nil {
		return nil, nil
	} else if _, err := vars.ComputeProgram(vt, makeEncoderX86, v == 1); err == nil {
		return compiler.rec, nil
	} else {
		return nil, err
	}
}

