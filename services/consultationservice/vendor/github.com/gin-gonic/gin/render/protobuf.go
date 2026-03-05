



package render

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)


type ProtoBuf struct {
	Data any
}

var protobufContentType = []string{"application/x-protobuf"}


func (r ProtoBuf) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := proto.Marshal(r.Data.(proto.Message))
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}


func (r ProtoBuf) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, protobufContentType)
}
