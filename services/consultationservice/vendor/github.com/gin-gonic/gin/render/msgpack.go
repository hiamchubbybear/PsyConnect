



//go:build !nomsgpack

package render

import (
	"net/http"

	"github.com/ugorji/go/codec"
)



var (
	_ Render = MsgPack{}
)


type MsgPack struct {
	Data any
}

var msgpackContentType = []string{"application/msgpack; charset=utf-8"}


func (r MsgPack) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, msgpackContentType)
}


func (r MsgPack) Render(w http.ResponseWriter) error {
	return WriteMsgPack(w, r.Data)
}


func WriteMsgPack(w http.ResponseWriter, obj any) error {
	writeContentType(w, msgpackContentType)
	var mh codec.MsgpackHandle
	return codec.NewEncoder(w, &mh).Encode(obj)
}
