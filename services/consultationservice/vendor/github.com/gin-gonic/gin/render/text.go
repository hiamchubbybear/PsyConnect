



package render

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/internal/bytesconv"
)


type String struct {
	Format string
	Data   []any
}

var plainContentType = []string{"text/plain; charset=utf-8"}


func (r String) Render(w http.ResponseWriter) error {
	return WriteString(w, r.Format, r.Data)
}


func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType)
}


func WriteString(w http.ResponseWriter, format string, data []any) (err error) {
	writeContentType(w, plainContentType)
	if len(data) > 0 {
		_, err = fmt.Fprintf(w, format, data...)
		return
	}
	_, err = w.Write(bytesconv.StringToBytes(format))
	return
}
