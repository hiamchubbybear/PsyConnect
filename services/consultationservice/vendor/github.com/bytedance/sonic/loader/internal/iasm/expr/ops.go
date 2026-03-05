















package expr

import (
	"fmt"
)

func idiv(v int64, d int64) (int64, error) {
	if d != 0 {
		return v / d, nil
	} else {
		return 0, newRuntimeError("division by zero")
	}
}

func imod(v int64, d int64) (int64, error) {
	if d != 0 {
		return v % d, nil
	} else {
		return 0, newRuntimeError("division by zero")
	}
}

func ipow(v int64, e int64) (int64, error) {
	mul := v
	ret := int64(1)

	
	if v < 0 {
		return 0, newRuntimeError(fmt.Sprintf("negative base value: %d", v))
	}

	
	if e < 0 {
		return 0, newRuntimeError(fmt.Sprintf("negative exponent: %d", e))
	}

	
	if (e & 1) != 0 {
		ret *= mul
	}

	
	for e >>= 1; e != 0; e >>= 1 {
		if mul *= mul; (e & 1) != 0 {
			ret *= mul
		}
	}

	
	return ret, nil
}
