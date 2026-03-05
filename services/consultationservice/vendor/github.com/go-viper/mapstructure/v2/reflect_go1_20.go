//go:build go1.20

package mapstructure

import "reflect"


func isComparable(v reflect.Value) bool {
	return v.Comparable()
}
