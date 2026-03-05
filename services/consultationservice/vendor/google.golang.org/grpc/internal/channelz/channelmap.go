

package channelz

import (
	"fmt"
	"sort"
	"sync"
	"time"
)


type entry interface {
	
	addChild(id int64, e entry)
	
	deleteChild(id int64)
	
	
	
	triggerDelete()
	
	
	
	deleteSelfIfReady()
	
	getParentID() int64
	Entity
}










type channelMap struct {
	mu               sync.RWMutex
	topLevelChannels map[int64]struct{}
	channels         map[int64]*Channel
	subChannels      map[int64]*SubChannel
	sockets          map[int64]*Socket
	servers          map[int64]*Server
}

func newChannelMap() *channelMap {
	return &channelMap{
		topLevelChannels: make(map[int64]struct{}),
		channels:         make(map[int64]*Channel),
		subChannels:      make(map[int64]*SubChannel),
		sockets:          make(map[int64]*Socket),
		servers:          make(map[int64]*Server),
	}
}

func (c *channelMap) addServer(id int64, s *Server) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s.cm = c
	c.servers[id] = s
}

func (c *channelMap) addChannel(id int64, cn *Channel, isTopChannel bool, pid int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cn.trace.cm = c
	c.channels[id] = cn
	if isTopChannel {
		c.topLevelChannels[id] = struct{}{}
	} else if p := c.channels[pid]; p != nil {
		p.addChild(id, cn)
	} else {
		logger.Infof("channel %d references invalid parent ID %d", id, pid)
	}
}

func (c *channelMap) addSubChannel(id int64, sc *SubChannel, pid int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sc.trace.cm = c
	c.subChannels[id] = sc
	if p := c.channels[pid]; p != nil {
		p.addChild(id, sc)
	} else {
		logger.Infof("subchannel %d references invalid parent ID %d", id, pid)
	}
}

func (c *channelMap) addSocket(s *Socket) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s.cm = c
	c.sockets[s.ID] = s
	if s.Parent == nil {
		logger.Infof("normal socket %d has no parent", s.ID)
	}
	s.Parent.(entry).addChild(s.ID, s)
}






func (c *channelMap) removeEntry(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.findEntry(id).triggerDelete()
}



type tracedChannel interface {
	getChannelTrace() *ChannelTrace
	incrTraceRefCount()
	decrTraceRefCount()
	getRefName() string
}


func (c *channelMap) decrTraceRefCount(id int64) {
	e := c.findEntry(id)
	if v, ok := e.(tracedChannel); ok {
		v.decrTraceRefCount()
		e.deleteSelfIfReady()
	}
}


func (c *channelMap) findEntry(id int64) entry {
	if v, ok := c.channels[id]; ok {
		return v
	}
	if v, ok := c.subChannels[id]; ok {
		return v
	}
	if v, ok := c.servers[id]; ok {
		return v
	}
	if v, ok := c.sockets[id]; ok {
		return v
	}
	return &dummyEntry{idNotFound: id}
}






func (c *channelMap) deleteEntry(id int64) entry {
	if v, ok := c.sockets[id]; ok {
		delete(c.sockets, id)
		return v
	}
	if v, ok := c.subChannels[id]; ok {
		delete(c.subChannels, id)
		return v
	}
	if v, ok := c.channels[id]; ok {
		delete(c.channels, id)
		delete(c.topLevelChannels, id)
		return v
	}
	if v, ok := c.servers[id]; ok {
		delete(c.servers, id)
		return v
	}
	return &dummyEntry{idNotFound: id}
}

func (c *channelMap) traceEvent(id int64, desc *TraceEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()
	child := c.findEntry(id)
	childTC, ok := child.(tracedChannel)
	if !ok {
		return
	}
	childTC.getChannelTrace().append(&traceEvent{Desc: desc.Desc, Severity: desc.Severity, Timestamp: time.Now()})
	if desc.Parent != nil {
		parent := c.findEntry(child.getParentID())
		var chanType RefChannelType
		switch child.(type) {
		case *Channel:
			chanType = RefChannel
		case *SubChannel:
			chanType = RefSubChannel
		}
		if parentTC, ok := parent.(tracedChannel); ok {
			parentTC.getChannelTrace().append(&traceEvent{
				Desc:      desc.Parent.Desc,
				Severity:  desc.Parent.Severity,
				Timestamp: time.Now(),
				RefID:     id,
				RefName:   childTC.getRefName(),
				RefType:   chanType,
			})
			childTC.incrTraceRefCount()
		}
	}
}

type int64Slice []int64

func (s int64Slice) Len() int           { return len(s) }
func (s int64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s int64Slice) Less(i, j int) bool { return s[i] < s[j] }

func copyMap(m map[int64]string) map[int64]string {
	n := make(map[int64]string)
	for k, v := range m {
		n[k] = v
	}
	return n
}

func (c *channelMap) getTopChannels(id int64, maxResults int) ([]*Channel, bool) {
	if maxResults <= 0 {
		maxResults = EntriesPerPage
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	l := int64(len(c.topLevelChannels))
	ids := make([]int64, 0, l)

	for k := range c.topLevelChannels {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool { return ids[i] >= id })
	end := true
	var t []*Channel
	for _, v := range ids[idx:] {
		if len(t) == maxResults {
			end = false
			break
		}
		if cn, ok := c.channels[v]; ok {
			t = append(t, cn)
		}
	}
	return t, end
}

func (c *channelMap) getServers(id int64, maxResults int) ([]*Server, bool) {
	if maxResults <= 0 {
		maxResults = EntriesPerPage
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	ids := make([]int64, 0, len(c.servers))
	for k := range c.servers {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool { return ids[i] >= id })
	end := true
	var s []*Server
	for _, v := range ids[idx:] {
		if len(s) == maxResults {
			end = false
			break
		}
		if svr, ok := c.servers[v]; ok {
			s = append(s, svr)
		}
	}
	return s, end
}

func (c *channelMap) getServerSockets(id int64, startID int64, maxResults int) ([]*Socket, bool) {
	if maxResults <= 0 {
		maxResults = EntriesPerPage
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	svr, ok := c.servers[id]
	if !ok {
		
		return nil, true
	}
	svrskts := svr.sockets
	ids := make([]int64, 0, len(svrskts))
	sks := make([]*Socket, 0, min(len(svrskts), maxResults))
	for k := range svrskts {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	idx := sort.Search(len(ids), func(i int) bool { return ids[i] >= startID })
	end := true
	for _, v := range ids[idx:] {
		if len(sks) == maxResults {
			end = false
			break
		}
		if ns, ok := c.sockets[v]; ok {
			sks = append(sks, ns)
		}
	}
	return sks, end
}

func (c *channelMap) getChannel(id int64) *Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channels[id]
}

func (c *channelMap) getSubChannel(id int64) *SubChannel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.subChannels[id]
}

func (c *channelMap) getSocket(id int64) *Socket {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sockets[id]
}

func (c *channelMap) getServer(id int64) *Server {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.servers[id]
}

type dummyEntry struct {
	
	idNotFound int64
	Entity
}

func (d *dummyEntry) String() string {
	return fmt.Sprintf("non-existent entity #%d", d.idNotFound)
}

func (d *dummyEntry) ID() int64 { return d.idNotFound }

func (d *dummyEntry) addChild(id int64, e entry) {
	
	
	
	
	
	
	
	
	
	logger.Infof("attempt to add child of type %T with id %d to a parent (id=%d) that doesn't currently exist", e, id, d.idNotFound)
}

func (d *dummyEntry) deleteChild(id int64) {
	
	
	logger.Infof("attempt to delete child with id %d from a parent (id=%d) that doesn't currently exist", id, d.idNotFound)
}

func (d *dummyEntry) triggerDelete() {
	logger.Warningf("attempt to delete an entry (id=%d) that doesn't currently exist", d.idNotFound)
}

func (*dummyEntry) deleteSelfIfReady() {
	
}

func (*dummyEntry) getParentID() int64 {
	return 0
}


type Entity interface {
	isEntity()
	fmt.Stringer
	id() int64
}
