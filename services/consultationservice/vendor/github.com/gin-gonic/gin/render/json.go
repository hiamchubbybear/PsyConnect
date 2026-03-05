



package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin/internal/bytesconv"
	"github.com/gin-gonic/gin/internal/json"
)


type JSON struct {
	Data any
}


type IndentedJSON struct {
	Data any
}


type SecureJSON struct {
	Prefix string
	Data   any
}


type JsonpJSON struct {
	Callback string
	Data     any
}


type AsciiJSON struct {
	Data any
}


type PureJSON struct {
	Data any
}

var (
	jsonContentType      = []string{"application/json; charset=utf-8"}
	jsonpContentType     = []string{"application/javascript; charset=utf-8"}
	jsonASCIIContentType = []string{"application/json"}
)


func (r JSON) Render(w http.ResponseWriter) error {
	return WriteJSON(w, r.Data)
}


func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}


func WriteJSON(w http.ResponseWriter, obj any) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}


func (r IndentedJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}


func (r IndentedJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}


func (r SecureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	
	if bytes.HasPrefix(jsonBytes, bytesconv.StringToBytes("[")) && bytes.HasSuffix(jsonBytes,
		bytesconv.StringToBytes("]")) {
		if _, err = w.Write(bytesconv.StringToBytes(r.Prefix)); err != nil {
			return err
		}
	}
	_, err = w.Write(jsonBytes)
	return err
}


func (r SecureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}


func (r JsonpJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		_, err = w.Write(ret)
		return err
	}

	callback := template.JSEscapeString(r.Callback)
	if _, err = w.Write(bytesconv.StringToBytes(callback)); err != nil {
		return err
	}

	if _, err = w.Write(bytesconv.StringToBytes("(")); err != nil {
		return err
	}

	if _, err = w.Write(ret); err != nil {
		return err
	}

	if _, err = w.Write(bytesconv.StringToBytes(");")); err != nil {
		return err
	}

	return nil
}


func (r JsonpJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}


func (r AsciiJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, r := range bytesconv.BytesToString(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	_, err = w.Write(buffer.Bytes())
	return err
}


func (r AsciiJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonASCIIContentType)
}


func (r PureJSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(r.Data)
}


func (r PureJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
