





package singleflight 

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
)



var errGoexit = errors.New("runtime.Goexit was called")



type panicError struct {
	value interface{}
	stack []byte
}


func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func (p *panicError) Unwrap() error {
	err, ok := p.value.(error)
	if !ok {
		return nil
	}

	return err
}

func newPanicError(v interface{}) error {
	stack := debug.Stack()

	
	
	
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}


type call struct {
	wg sync.WaitGroup

	
	
	val interface{}
	err error

	
	
	
	dups  int
	chans []chan<- Result
}



type Group struct {
	mu sync.Mutex       
	m  map[string]*call 
}



type Result struct {
	Val    interface{}
	Err    error
	Shared bool
}






func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}





func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)

	return ch
}


func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

	
	
	defer func() {
		
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			
			
			if len(c.chans) > 0 {
				go panic(e)
				select {} 
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			
		} else {
			
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				
				
				
				
				
				
				
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}




func (g *Group) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}
