

package grpcsync

import (
	"context"

	"google.golang.org/grpc/internal/buffer"
)







type CallbackSerializer struct {
	
	
	
	done chan struct{}

	callbacks *buffer.Unbounded
}






func NewCallbackSerializer(ctx context.Context) *CallbackSerializer {
	cs := &CallbackSerializer{
		done:      make(chan struct{}),
		callbacks: buffer.NewUnbounded(),
	}
	go cs.run(ctx)
	return cs
}








func (cs *CallbackSerializer) TrySchedule(f func(ctx context.Context)) {
	cs.callbacks.Put(f)
}








func (cs *CallbackSerializer) ScheduleOr(f func(ctx context.Context), onFailure func()) {
	if cs.callbacks.Put(f) != nil {
		onFailure()
	}
}

func (cs *CallbackSerializer) run(ctx context.Context) {
	defer close(cs.done)

	
	
	
	
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			
			
		case cb := <-cs.callbacks.Get():
			cs.callbacks.Load()
			cb.(func(context.Context))(ctx)
		}
	}

	
	cs.callbacks.Close()

	
	for cb := range cs.callbacks.Get() {
		cs.callbacks.Load()
		cb.(func(context.Context))(ctx)
	}
}



func (cs *CallbackSerializer) Done() <-chan struct{} {
	return cs.done
}
