


package otel 


type ErrorHandler interface {
	
	

	
	
	Handle(error)
	
	
}



type ErrorHandlerFunc func(error)

var _ ErrorHandler = ErrorHandlerFunc(nil)


func (f ErrorHandlerFunc) Handle(err error) {
	f(err)
}
