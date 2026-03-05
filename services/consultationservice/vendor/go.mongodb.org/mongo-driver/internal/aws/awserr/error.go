










package awserr




type Error interface {
	
	error

	
	Code() string

	
	Message() string

	
	OrigErr() error
}






type BatchedErrors interface {
	
	Error

	
	OrigErrs() []error
}





func New(code, message string, origErr error) Error {
	var errs []error
	if origErr != nil {
		errs = append(errs, origErr)
	}
	return newBaseError(code, message, errs)
}



func NewBatchError(code, message string, errs []error) BatchedErrors {
	return newBaseError(code, message, errs)
}
