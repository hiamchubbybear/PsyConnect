



package render

import (
	"io"
	"net/http"
	"strconv"
)


type Reader struct {
	ContentType   string
	ContentLength int64
	Reader        io.Reader
	Headers       map[string]string
}


func (r Reader) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	if r.ContentLength >= 0 {
		if r.Headers == nil {
			r.Headers = map[string]string{}
		}
		r.Headers["Content-Length"] = strconv.FormatInt(r.ContentLength, 10)
	}
	r.writeHeaders(w, r.Headers)
	_, err = io.Copy(w, r.Reader)
	return
}


func (r Reader) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, []string{r.ContentType})
}


func (r Reader) writeHeaders(w http.ResponseWriter, headers map[string]string) {
	header := w.Header()
	for k, v := range headers {
		if header.Get(k) == "" {
			header.Set(k, v)
		}
	}
}
