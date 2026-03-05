



package render

import "net/http"


type Render interface {
	
	Render(http.ResponseWriter) error
	
	WriteContentType(w http.ResponseWriter)
}

var (
	_ Render     = (*JSON)(nil)
	_ Render     = (*IndentedJSON)(nil)
	_ Render     = (*SecureJSON)(nil)
	_ Render     = (*JsonpJSON)(nil)
	_ Render     = (*XML)(nil)
	_ Render     = (*String)(nil)
	_ Render     = (*Redirect)(nil)
	_ Render     = (*Data)(nil)
	_ Render     = (*HTML)(nil)
	_ HTMLRender = (*HTMLDebug)(nil)
	_ HTMLRender = (*HTMLProduction)(nil)
	_ Render     = (*YAML)(nil)
	_ Render     = (*Reader)(nil)
	_ Render     = (*AsciiJSON)(nil)
	_ Render     = (*ProtoBuf)(nil)
	_ Render     = (*TOML)(nil)
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
