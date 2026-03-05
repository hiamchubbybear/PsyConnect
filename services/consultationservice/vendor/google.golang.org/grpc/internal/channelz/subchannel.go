

package channelz

import (
	"fmt"
	"sync/atomic"
)


type SubChannel struct {
	Entity
	
	ID int64
	
	RefName       string
	closeCalled   bool
	sockets       map[int64]string
	parent        *Channel
	trace         *ChannelTrace
	traceRefCount int32

	ChannelMetrics ChannelMetrics
}

func (sc *SubChannel) String() string {
	return fmt.Sprintf("%s SubChannel #%d", sc.parent, sc.ID)
}

func (sc *SubChannel) id() int64 {
	return sc.ID
}


func (sc *SubChannel) Sockets() map[int64]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return copyMap(sc.sockets)
}


func (sc *SubChannel) Trace() *ChannelTrace {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return sc.trace.copy()
}

func (sc *SubChannel) addChild(id int64, e entry) {
	if v, ok := e.(*Socket); ok && v.SocketType == SocketTypeNormal {
		sc.sockets[id] = v.RefName
	} else {
		logger.Errorf("cannot add a child (id = %d) of type %T to a subChannel", id, e)
	}
}

func (sc *SubChannel) deleteChild(id int64) {
	delete(sc.sockets, id)
	sc.deleteSelfIfReady()
}

func (sc *SubChannel) triggerDelete() {
	sc.closeCalled = true
	sc.deleteSelfIfReady()
}

func (sc *SubChannel) getParentID() int64 {
	return sc.parent.ID
}








func (sc *SubChannel) deleteSelfFromTree() (deleted bool) {
	if !sc.closeCalled || len(sc.sockets) != 0 {
		return false
	}
	sc.parent.deleteChild(sc.ID)
	return true
}













func (sc *SubChannel) deleteSelfFromMap() (delete bool) {
	return sc.getTraceRefCount() == 0
}







func (sc *SubChannel) deleteSelfIfReady() {
	if !sc.deleteSelfFromTree() {
		return
	}
	if !sc.deleteSelfFromMap() {
		return
	}
	db.deleteEntry(sc.ID)
	sc.trace.clear()
}

func (sc *SubChannel) getChannelTrace() *ChannelTrace {
	return sc.trace
}

func (sc *SubChannel) incrTraceRefCount() {
	atomic.AddInt32(&sc.traceRefCount, 1)
}

func (sc *SubChannel) decrTraceRefCount() {
	atomic.AddInt32(&sc.traceRefCount, -1)
}

func (sc *SubChannel) getTraceRefCount() int {
	i := atomic.LoadInt32(&sc.traceRefCount)
	return int(i)
}

func (sc *SubChannel) getRefName() string {
	return sc.RefName
}
