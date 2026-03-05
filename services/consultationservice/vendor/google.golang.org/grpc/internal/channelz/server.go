

package channelz

import (
	"fmt"
	"sync/atomic"
)


type Server struct {
	Entity
	ID      int64
	RefName string

	ServerMetrics ServerMetrics

	closeCalled   bool
	sockets       map[int64]string
	listenSockets map[int64]string
	cm            *channelMap
}


type ServerMetrics struct {
	
	CallsStarted atomic.Int64
	
	CallsSucceeded atomic.Int64
	
	CallsFailed atomic.Int64
	
	LastCallStartedTimestamp atomic.Int64
}


func NewServerMetricsForTesting(started, succeeded, failed, timestamp int64) *ServerMetrics {
	sm := &ServerMetrics{}
	sm.CallsStarted.Store(started)
	sm.CallsSucceeded.Store(succeeded)
	sm.CallsFailed.Store(failed)
	sm.LastCallStartedTimestamp.Store(timestamp)
	return sm
}



func (sm *ServerMetrics) CopyFrom(o *ServerMetrics) {
	sm.CallsStarted.Store(o.CallsStarted.Load())
	sm.CallsSucceeded.Store(o.CallsSucceeded.Load())
	sm.CallsFailed.Store(o.CallsFailed.Load())
	sm.LastCallStartedTimestamp.Store(o.LastCallStartedTimestamp.Load())
}


func (s *Server) ListenSockets() map[int64]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return copyMap(s.listenSockets)
}


func (s *Server) String() string {
	return fmt.Sprintf("Server #%d", s.ID)
}

func (s *Server) id() int64 {
	return s.ID
}

func (s *Server) addChild(id int64, e entry) {
	switch v := e.(type) {
	case *Socket:
		switch v.SocketType {
		case SocketTypeNormal:
			s.sockets[id] = v.RefName
		case SocketTypeListen:
			s.listenSockets[id] = v.RefName
		}
	default:
		logger.Errorf("cannot add a child (id = %d) of type %T to a server", id, e)
	}
}

func (s *Server) deleteChild(id int64) {
	delete(s.sockets, id)
	delete(s.listenSockets, id)
	s.deleteSelfIfReady()
}

func (s *Server) triggerDelete() {
	s.closeCalled = true
	s.deleteSelfIfReady()
}

func (s *Server) deleteSelfIfReady() {
	if !s.closeCalled || len(s.sockets)+len(s.listenSockets) != 0 {
		return
	}
	s.cm.deleteEntry(s.ID)
}

func (s *Server) getParentID() int64 {
	return 0
}
