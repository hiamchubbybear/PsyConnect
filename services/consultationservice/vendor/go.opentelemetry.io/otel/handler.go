


package otel 

import (
	"go.opentelemetry.io/otel/internal/global"
)


var _ ErrorHandler = (*global.ErrDelegator)(nil)










func GetErrorHandler() ErrorHandler { return global.GetErrorHandler() }







func SetErrorHandler(h ErrorHandler) { global.SetErrorHandler(h) }


func Handle(err error) { global.GetErrorHandler().Handle(err) }
