package concurrent

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
	"reflect"
)


var HandlePanic = func(recovered interface{}, funcName string) {
	ErrorLogger.Println(fmt.Sprintf("%s panic: %v", funcName, recovered))
	ErrorLogger.Println(string(debug.Stack()))
}



type UnboundedExecutor struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	activeGoroutinesMutex *sync.Mutex
	activeGoroutines      map[string]int
	HandlePanic           func(recovered interface{}, funcName string)
}





var GlobalUnboundedExecutor = NewUnboundedExecutor()




func NewUnboundedExecutor() *UnboundedExecutor {
	ctx, cancel := context.WithCancel(context.TODO())
	return &UnboundedExecutor{
		ctx:                   ctx,
		cancel:                cancel,
		activeGoroutinesMutex: &sync.Mutex{},
		activeGoroutines:      map[string]int{},
	}
}



func (executor *UnboundedExecutor) Go(handler func(ctx context.Context)) {
	pc := reflect.ValueOf(handler).Pointer()
	f := runtime.FuncForPC(pc)
	funcName := f.Name()
	file, line := f.FileLine(pc)
	executor.activeGoroutinesMutex.Lock()
	defer executor.activeGoroutinesMutex.Unlock()
	startFrom := fmt.Sprintf("%s:%d", file, line)
	executor.activeGoroutines[startFrom] += 1
	go func() {
		defer func() {
			recovered := recover()
			
			
			if recovered != nil {
				if executor.HandlePanic == nil {
					HandlePanic(recovered, funcName)
				} else {
					executor.HandlePanic(recovered, funcName)
				}
			}
			executor.activeGoroutinesMutex.Lock()
			executor.activeGoroutines[startFrom] -= 1
			executor.activeGoroutinesMutex.Unlock()
		}()
		handler(executor.ctx)
	}()
}


func (executor *UnboundedExecutor) Stop() {
	executor.cancel()
}



func (executor *UnboundedExecutor) StopAndWaitForever() {
	executor.StopAndWait(context.Background())
}



func (executor *UnboundedExecutor) StopAndWait(ctx context.Context) {
	executor.cancel()
	for {
		oneHundredMilliseconds := time.NewTimer(time.Millisecond * 100)
		select {
		case <-oneHundredMilliseconds.C:
			if executor.checkNoActiveGoroutines() {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (executor *UnboundedExecutor) checkNoActiveGoroutines() bool {
	executor.activeGoroutinesMutex.Lock()
	defer executor.activeGoroutinesMutex.Unlock()
	for startFrom, count := range executor.activeGoroutines {
		if count > 0 {
			InfoLogger.Println("UnboundedExecutor is still waiting goroutines to quit",
				"startFrom", startFrom,
				"count", count)
			return false
		}
	}
	return true
}
