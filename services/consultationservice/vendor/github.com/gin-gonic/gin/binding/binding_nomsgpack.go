



//go:build nomsgpack

package binding

import "net/http"


const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEYAML              = "application/x-yaml"
	MIMEYAML2             = "application/yaml"
	MIMETOML              = "application/toml"
)




type Binding interface {
	Name() string
	Bind(*http.Request, any) error
}



type BindingBody interface {
	Binding
	BindBody([]byte, any) error
}



type BindingUri interface {
	Name() string
	BindUri(map[string][]string, any) error
}





type StructValidator interface {
	
	
	
	
	
	ValidateStruct(any) error

	
	
	Engine() any
}




var Validator StructValidator = &defaultValidator{}



var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Header        = headerBinding{}
	TOML          = tomlBinding{}
)



func Default(method, contentType string) Binding {
	if method == "GET" {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEYAML, MIMEYAML2:
		return YAML
	case MIMEMultipartPOSTForm:
		return FormMultipart
	case MIMETOML:
		return TOML
	default: 
		return Form
	}
}

func validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
