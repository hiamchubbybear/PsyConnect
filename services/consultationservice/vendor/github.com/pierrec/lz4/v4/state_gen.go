

package lz4

import "strconv"

func _() {
	
	
	var x [1]struct{}
	_ = x[noState-0]
	_ = x[errorState-1]
	_ = x[newState-2]
	_ = x[readState-3]
	_ = x[writeState-4]
	_ = x[closedState-5]
}

const _aState_name = "noStateerrorStatenewStatereadStatewriteStateclosedState"

var _aState_index = [...]uint8{0, 7, 17, 25, 34, 44, 55}

func (i aState) String() string {
	if i >= aState(len(_aState_index)-1) {
		return "aState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _aState_name[_aState_index[i]:_aState_index[i+1]]
}
