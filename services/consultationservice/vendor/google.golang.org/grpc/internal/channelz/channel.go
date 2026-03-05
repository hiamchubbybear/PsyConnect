

package channelz

import (
	"fmt"
	"sync/atomic"

	"google.golang.org/grpc/connectivity"
)



type Channel struct {
	Entity
	
	ID int64
	
	RefName string

	closeCalled bool
	nestedChans map[int64]string
	subChans    map[int64]string
	Parent      *Channel
	trace       *ChannelTrace
	
	
	traceRefCount int32

	
	
	ChannelMetrics ChannelMetrics
}



func (c *Channel) channelzIdentifier() {}



func (c *Channel) String() string {
	if c.Parent == nil {
		return fmt.Sprintf("Channel #%d", c.ID)
	}
	return fmt.Sprintf("%s Channel #%d", c.Parent, c.ID)
}

func (c *Channel) id() int64 {
	return c.ID
}



func (c *Channel) SubChans() map[int64]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return copyMap(c.subChans)
}



func (c *Channel) NestedChans() map[int64]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return copyMap(c.nestedChans)
}


func (c *Channel) Trace() *ChannelTrace {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return c.trace.copy()
}



type ChannelMetrics struct {
	
	State atomic.Pointer[connectivity.State]
	
	Target atomic.Pointer[string]
	
	CallsStarted atomic.Int64
	
	CallsSucceeded atomic.Int64
	
	CallsFailed atomic.Int64
	
	LastCallStartedTimestamp atomic.Int64
}


func (c *ChannelMetrics) CopyFrom(o *ChannelMetrics) {
	c.State.Store(o.State.Load())
	c.Target.Store(o.Target.Load())
	c.CallsStarted.Store(o.CallsStarted.Load())
	c.CallsSucceeded.Store(o.CallsSucceeded.Load())
	c.CallsFailed.Store(o.CallsFailed.Load())
	c.LastCallStartedTimestamp.Store(o.LastCallStartedTimestamp.Load())
}



func (c *ChannelMetrics) Equal(o any) bool {
	oc, ok := o.(*ChannelMetrics)
	if !ok {
		return false
	}
	if (c.State.Load() == nil) != (oc.State.Load() == nil) {
		return false
	}
	if c.State.Load() != nil && *c.State.Load() != *oc.State.Load() {
		return false
	}
	if (c.Target.Load() == nil) != (oc.Target.Load() == nil) {
		return false
	}
	if c.Target.Load() != nil && *c.Target.Load() != *oc.Target.Load() {
		return false
	}
	return c.CallsStarted.Load() == oc.CallsStarted.Load() &&
		c.CallsFailed.Load() == oc.CallsFailed.Load() &&
		c.CallsSucceeded.Load() == oc.CallsSucceeded.Load() &&
		c.LastCallStartedTimestamp.Load() == oc.LastCallStartedTimestamp.Load()
}

func strFromPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}



func (c *ChannelMetrics) String() string {
	return fmt.Sprintf("State: %v, Target: %s, CallsStarted: %v, CallsSucceeded: %v, CallsFailed: %v, LastCallStartedTimestamp: %v",
		c.State.Load(), strFromPointer(c.Target.Load()), c.CallsStarted.Load(), c.CallsSucceeded.Load(), c.CallsFailed.Load(), c.LastCallStartedTimestamp.Load(),
	)
}



func NewChannelMetricForTesting(state connectivity.State, target string, started, succeeded, failed, timestamp int64) *ChannelMetrics {
	c := &ChannelMetrics{}
	c.State.Store(&state)
	c.Target.Store(&target)
	c.CallsStarted.Store(started)
	c.CallsSucceeded.Store(succeeded)
	c.CallsFailed.Store(failed)
	c.LastCallStartedTimestamp.Store(timestamp)
	return c
}

func (c *Channel) addChild(id int64, e entry) {
	switch v := e.(type) {
	case *SubChannel:
		c.subChans[id] = v.RefName
	case *Channel:
		c.nestedChans[id] = v.RefName
	default:
		logger.Errorf("cannot add a child (id = %d) of type %T to a channel", id, e)
	}
}

func (c *Channel) deleteChild(id int64) {
	delete(c.subChans, id)
	delete(c.nestedChans, id)
	c.deleteSelfIfReady()
}

func (c *Channel) triggerDelete() {
	c.closeCalled = true
	c.deleteSelfIfReady()
}

func (c *Channel) getParentID() int64 {
	if c.Parent == nil {
		return -1
	}
	return c.Parent.ID
}








func (c *Channel) deleteSelfFromTree() (deleted bool) {
	if !c.closeCalled || len(c.subChans)+len(c.nestedChans) != 0 {
		return false
	}
	
	if c.Parent != nil {
		c.Parent.deleteChild(c.ID)
	}
	return true
}













func (c *Channel) deleteSelfFromMap() (delete bool) {
	return c.getTraceRefCount() == 0
}







func (c *Channel) deleteSelfIfReady() {
	if !c.deleteSelfFromTree() {
		return
	}
	if !c.deleteSelfFromMap() {
		return
	}
	db.deleteEntry(c.ID)
	c.trace.clear()
}

func (c *Channel) getChannelTrace() *ChannelTrace {
	return c.trace
}

func (c *Channel) incrTraceRefCount() {
	atomic.AddInt32(&c.traceRefCount, 1)
}

func (c *Channel) decrTraceRefCount() {
	atomic.AddInt32(&c.traceRefCount, -1)
}

func (c *Channel) getTraceRefCount() int {
	i := atomic.LoadInt32(&c.traceRefCount)
	return int(i)
}

func (c *Channel) getRefName() string {
	return c.RefName
}
