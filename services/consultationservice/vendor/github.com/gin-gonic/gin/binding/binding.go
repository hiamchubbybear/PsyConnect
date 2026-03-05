



//go:build !nomsgpack

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
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
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
	JSON          BindingBody = jsonBinding{}
	XML           BindingBody = xmlBinding{}
	Form          Binding     = formBinding{}
	Query         Binding     = queryBinding{}
	FormPost      Binding     = formPostBinding{}
	FormMultipart Binding     = formMultipartBinding{}
	ProtoBuf      BindingBody = protobufBinding{}
	MsgPack       BindingBody = msgpackBinding{}
	YAML          BindingBody = yamlBinding{}
	Uri           BindingUri  = uriBinding{}
	Header        Binding     = headerBinding{}
	TOML          BindingBody = tomlBinding{}
)



func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack
	case MIMEYAML, MIMEYAML2:
		return YAML
	case MIMETOML:
		return TOML
	case MIMEMultipartPOSTForm:
		return FormMultipart
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
