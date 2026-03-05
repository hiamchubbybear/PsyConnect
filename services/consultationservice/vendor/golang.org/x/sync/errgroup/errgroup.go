








package errgroup

import (
	"context"
	"fmt"
	"sync"
)

type token struct{}






type Group struct {
	cancel func(error)

	wg sync.WaitGroup

	sem chan token

	errOnce sync.Once
	err     error
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done()
}






func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}



func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel(g.err)
	}
	return g.err
}







func (g *Group) Go(f func() error) {
	if g.sem != nil {
		g.sem <- token{}
	}

	g.wg.Add(1)
	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
}





func (g *Group) TryGo(f func() error) bool {
	if g.sem != nil {
		select {
		case g.sem <- token{}:
			
		default:
			return false
		}
	}

	g.wg.Add(1)
	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel(g.err)
				}
			})
		}
	}()
	return true
}









func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}
	if len(g.sem) != 0 {
		panic(fmt.Errorf("errgroup: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}
	g.sem = make(chan token, n)
}
