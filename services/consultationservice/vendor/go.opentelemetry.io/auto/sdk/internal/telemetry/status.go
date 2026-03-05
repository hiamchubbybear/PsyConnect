


package telemetry



type StatusCode int32

const (
	
	StatusCodeUnset StatusCode = 0
	
	
	StatusCodeOK StatusCode = 1
	
	StatusCodeError StatusCode = 2
)

var statusCodeStrings = []string{
	"Unset",
	"OK",
	"Error",
}

func (s StatusCode) String() string {
	if s >= 0 && int(s) < len(statusCodeStrings) {
		return statusCodeStrings[s]
	}
	return "<unknown telemetry.StatusCode>"
}



type Status struct {
	
	Message string `json:"message,omitempty"`
	
	Code StatusCode `json:"code,omitempty"`
}
