



package render

import (
	"net/http"

	"gopkg.in/yaml.v3"
)


type YAML struct {
	Data any
}

var yamlContentType = []string{"application/yaml; charset=utf-8"}


func (r YAML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	bytes, err := yaml.Marshal(r.Data)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}


func (r YAML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, yamlContentType)
}
