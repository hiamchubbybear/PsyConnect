

package grpcsync

import (
	"context"
	"sync"
)




type Subscriber interface {
	
	
	OnMessage(msg any)
}












type PubSub struct {
	cs *CallbackSerializer

	
	mu          sync.Mutex
	msg         any
	subscribers map[Subscriber]bool
}



func NewPubSub(ctx context.Context) *PubSub {
	return &PubSub{
		cs:          NewCallbackSerializer(ctx),
		subscribers: map[Subscriber]bool{},
	}
}









func (ps *PubSub) Subscribe(sub Subscriber) (cancel func()) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subscribers[sub] = true

	if ps.msg != nil {
		msg := ps.msg
		ps.cs.TrySchedule(func(context.Context) {
			ps.mu.Lock()
			defer ps.mu.Unlock()
			if !ps.subscribers[sub] {
				return
			}
			sub.OnMessage(msg)
		})
	}

	return func() {
		ps.mu.Lock()
		defer ps.mu.Unlock()
		delete(ps.subscribers, sub)
	}
}



func (ps *PubSub) Publish(msg any) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.msg = msg
	for sub := range ps.subscribers {
		s := sub
		ps.cs.TrySchedule(func(context.Context) {
			ps.mu.Lock()
			defer ps.mu.Unlock()
			if !ps.subscribers[s] {
				return
			}
			s.OnMessage(msg)
		})
	}
}



func (ps *PubSub) Done() <-chan struct{} {
	return ps.cs.Done()
}
