









package awserr

import (
	"fmt"
)





func SprintError(code, message, extra string, origErr error) string {
	msg := fmt.Sprintf("%s: %s", code, message)
	if extra != "" {
		msg = fmt.Sprintf("%s\n\t%s", msg, extra)
	}
	if origErr != nil {
		msg = fmt.Sprintf("%s\ncaused by: %s", msg, origErr.Error())
	}
	return msg
}






type baseError struct {
	
	code string

	
	message string

	
	
	errs []error
}











func newBaseError(code, message string, origErrs []error) *baseError {
	b := &baseError{
		code:    code,
		message: message,
		errs:    origErrs,
	}

	return b
}






func (b baseError) Error() string {
	size := len(b.errs)
	if size > 0 {
		return SprintError(b.code, b.message, "", errorList(b.errs))
	}

	return SprintError(b.code, b.message, "", nil)
}



func (b baseError) String() string {
	return b.Error()
}


func (b baseError) Code() string {
	return b.code
}


func (b baseError) Message() string {
	return b.message
}




func (b baseError) OrigErr() error {
	switch len(b.errs) {
	case 0:
		return nil
	case 1:
		return b.errs[0]
	default:
		if err, ok := b.errs[0].(Error); ok {
			return NewBatchError(err.Code(), err.Message(), b.errs[1:])
		}
		return NewBatchError("BatchedErrors",
			"multiple errors occurred", b.errs)
	}
}



func (b baseError) OrigErrs() []error {
	return b.errs
}


type errorList []error




func (e errorList) Error() string {
	msg := ""
	
	if size := len(e); size > 0 {
		for i := 0; i < size; i++ {
			msg += e[i].Error()
			
			
			
			if i+1 < size {
				msg += "\n"
			}
		}
	}
	return msg
}
