


package otel 

import (
	"github.com/go-logr/logr"

	"go.opentelemetry.io/otel/internal/global"
)


func SetLogger(logger logr.Logger) {
	global.SetLogger(logger)
}
