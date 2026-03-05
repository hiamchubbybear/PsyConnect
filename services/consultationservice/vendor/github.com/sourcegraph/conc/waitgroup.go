package conc

import (
	"sync"

	"github.com/sourcegraph/conc/panics"
)


func NewWaitGroup() *WaitGroup {
	return &WaitGroup{}
}









type WaitGroup struct {
	wg sync.WaitGroup
	pc panics.Catcher
}


func (h *WaitGroup) Go(f func()) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.pc.Try(f)
	}()
}



func (h *WaitGroup) Wait() {
	h.wg.Wait()

	
	h.pc.Repanic()
}



func (h *WaitGroup) WaitAndRecover() *panics.Recovered {
	h.wg.Wait()

	
	return h.pc.Recovered()
}
