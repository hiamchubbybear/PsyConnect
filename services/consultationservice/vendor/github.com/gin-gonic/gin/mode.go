



package gin

import (
	"flag"
	"io"
	"os"

	"github.com/gin-gonic/gin/binding"
)


const EnvGinMode = "GIN_MODE"

const (
	
	DebugMode = "debug"
	
	ReleaseMode = "release"
	
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)









var DefaultWriter io.Writer = os.Stdout


var DefaultErrorWriter io.Writer = os.Stderr

var (
	ginMode  = debugCode
	modeName = DebugMode
)

func init() {
	mode := os.Getenv(EnvGinMode)
	SetMode(mode)
}


func SetMode(value string) {
	if value == "" {
		if flag.Lookup("test.v") != nil {
			value = TestMode
		} else {
			value = DebugMode
		}
	}

	switch value {
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("gin mode unknown: " + value + " (available mode: debug release test)")
	}

	modeName = value
}


func DisableBindValidation() {
	binding.Validator = nil
}



func EnableJsonDecoderUseNumber() {
	binding.EnableDecoderUseNumber = true
}



func EnableJsonDecoderDisallowUnknownFields() {
	binding.EnableDecoderDisallowUnknownFields = true
}


func Mode() string {
	return modeName
}
