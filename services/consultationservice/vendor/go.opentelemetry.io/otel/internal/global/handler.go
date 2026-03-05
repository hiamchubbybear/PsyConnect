


package global 

import (
	"log"
	"sync/atomic"
)


type ErrorHandler interface {
	
	
	Handle(error)
}

type ErrDelegator struct {
	delegate atomic.Pointer[ErrorHandler]
}


var _ ErrorHandler = (*ErrDelegator)(nil)

func (d *ErrDelegator) Handle(err error) {
	if eh := d.delegate.Load(); eh != nil {
		(*eh).Handle(err)
		return
	}
	log.Print(err)
}


func (d *ErrDelegator) setDelegate(eh ErrorHandler) {
	d.delegate.Store(&eh)
}
