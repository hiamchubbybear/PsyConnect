


//go:build notfastpath || codec.notfastpath
// +build notfastpath codec.notfastpath

package codec

import "reflect"

const fastpathEnabled = false









func fastpathDecodeTypeSwitch(iv interface{}, d *Decoder) bool { return false }
func fastpathEncodeTypeSwitch(iv interface{}, e *Encoder) bool { return false }




func fastpathDecodeSetZeroTypeSwitch(iv interface{}) bool { return false }

type fastpathT struct{}
type fastpathE struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*Encoder, *codecFnInfo, reflect.Value)
	decfn func(*Decoder, *codecFnInfo, reflect.Value)
}
type fastpathA [0]fastpathE

func fastpathAvIndex(rtid uintptr) int { return -1 }

var fastpathAv fastpathA
var fastpathTV fastpathT
