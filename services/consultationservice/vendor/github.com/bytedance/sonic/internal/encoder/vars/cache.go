

package vars

import (
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

type Encoder func(
	rb *[]byte,
	vp unsafe.Pointer,
	sb *Stack,
	fv uint64,
) error

func FindOrCompile(vt *rt.GoType, pv bool, compiler func(*rt.GoType, ... interface{}) (interface{}, error)) (interface{}, error) {
	if val := programCache.Get(vt); val != nil {
		return val, nil
	} else if ret, err := programCache.Compute(vt, compiler, pv); err == nil {
		return ret, nil
	} else {
		return nil, err
	}
}

func GetProgram(vt *rt.GoType) (interface{}) {
	return programCache.Get(vt)
}

func ComputeProgram(vt *rt.GoType, compute func(*rt.GoType, ... interface{}) (interface{}, error), pv bool) (interface{}, error) {
	return programCache.Compute(vt, compute, pv)
}