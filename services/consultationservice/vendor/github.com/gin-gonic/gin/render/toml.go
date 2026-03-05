



package render

import (
	"net/http"

	"github.com/pelletier/go-toml/v2"
)


type TOML struct {
	Data any
}

var TOMLContentType = []string{"application/toml; charset=utf-8"}


func (r TOML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := toml.Marshal(r.Data)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}


func (r TOML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, TOMLContentType)
}
