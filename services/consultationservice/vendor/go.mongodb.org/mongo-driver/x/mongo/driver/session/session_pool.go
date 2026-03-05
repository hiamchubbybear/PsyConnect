





package session

import (
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type Node struct {
	*Server
	next *Node
	prev *Node
}



type topologyDescription struct {
	kind           description.TopologyKind
	timeoutMinutes *int64
}


type Pool struct {
	
	checkedOut int64

	descChan       <-chan description.Topology
	head           *Node
	tail           *Node
	latestTopology topologyDescription
	mutex          sync.Mutex 
}

func (p *Pool) createServerSession() (*Server, error) {
	s, err := newServerSession()
	if err != nil {
		return nil, err
	}

	atomic.AddInt64(&p.checkedOut, 1)
	return s, nil
}


func NewPool(descChan <-chan description.Topology) *Pool {
	p := &Pool{
		descChan: descChan,
	}

	return p
}


func (p *Pool) updateTimeout() {
	select {
	case newDesc := <-p.descChan:
		p.latestTopology = topologyDescription{
			kind:           newDesc.Kind,
			timeoutMinutes: newDesc.SessionTimeoutMinutesPtr,
		}
	default:
		
	}
}


func (p *Pool) GetSession() (*Server, error) {
	p.mutex.Lock() 
	defer p.mutex.Unlock()

	
	if p.head == nil && p.tail == nil {
		return p.createServerSession()
	}

	p.updateTimeout()
	for p.head != nil {
		
		if p.head.expired(p.latestTopology) {
			p.head = p.head.next
			continue
		}

		
		session := p.head.Server
		if p.head.next != nil {
			p.head.next.prev = nil
		}
		if p.tail == p.head {
			p.tail = nil
			p.head = nil
		} else {
			p.head = p.head.next
		}

		atomic.AddInt64(&p.checkedOut, 1)
		return session, nil
	}

	
	p.tail = nil 
	return p.createServerSession()
}


func (p *Pool) ReturnSession(ss *Server) {
	if ss == nil {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	atomic.AddInt64(&p.checkedOut, -1)
	p.updateTimeout()
	
	
	for p.tail != nil && p.tail.expired(p.latestTopology) {
		if p.tail.prev != nil {
			p.tail.prev.next = nil
		}
		p.tail = p.tail.prev
	}

	
	if ss.expired(p.latestTopology) {
		return
	}

	
	if ss.Dirty {
		return
	}

	newNode := &Node{
		Server: ss,
		next:   nil,
		prev:   nil,
	}

	
	if p.tail == nil {
		p.head = newNode
		p.tail = newNode
		return
	}

	
	newNode.next = p.head
	p.head.prev = newNode
	p.head = newNode
}


func (p *Pool) IDSlice() []bsoncore.Document {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var ids []bsoncore.Document
	for node := p.head; node != nil; node = node.next {
		ids = append(ids, node.SessionID)
	}

	return ids
}


func (p *Pool) String() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	s := ""
	for head := p.head; head != nil; head = head.next {
		s += head.SessionID.String() + "\n"
	}

	return s
}


func (p *Pool) CheckedOut() int64 {
	return atomic.LoadInt64(&p.checkedOut)
}
