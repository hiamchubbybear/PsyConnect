




package channelz

import (
	"sync/atomic"
	"time"

	"google.golang.org/grpc/internal"
)

var (
	
	
	IDGen IDGenerator

	db = newChannelMap()
	
	EntriesPerPage = 50
	curState       int32
)


func TurnOn() {
	atomic.StoreInt32(&curState, 1)
}

func init() {
	internal.ChannelzTurnOffForTesting = func() {
		atomic.StoreInt32(&curState, 0)
	}
}


func IsOn() bool {
	return atomic.LoadInt32(&curState) == 1
}








func GetTopChannels(id int64, maxResults int) ([]*Channel, bool) {
	return db.getTopChannels(id, maxResults)
}







func GetServers(id int64, maxResults int) ([]*Server, bool) {
	return db.getServers(id, maxResults)
}








func GetServerSockets(id int64, startID int64, maxResults int) ([]*Socket, bool) {
	return db.getServerSockets(id, startID, maxResults)
}


func GetChannel(id int64) *Channel {
	return db.getChannel(id)
}


func GetSubChannel(id int64) *SubChannel {
	return db.getSubChannel(id)
}


func GetSocket(id int64) *Socket {
	return db.getSocket(id)
}


func GetServer(id int64) *Server {
	return db.getServer(id)
}








func RegisterChannel(parent *Channel, target string) *Channel {
	id := IDGen.genID()

	if !IsOn() {
		return &Channel{ID: id}
	}

	isTopChannel := parent == nil

	cn := &Channel{
		ID:          id,
		RefName:     target,
		nestedChans: make(map[int64]string),
		subChans:    make(map[int64]string),
		Parent:      parent,
		trace:       &ChannelTrace{CreationTime: time.Now(), Events: make([]*traceEvent, 0, getMaxTraceEntry())},
	}
	cn.ChannelMetrics.Target.Store(&target)
	db.addChannel(id, cn, isTopChannel, cn.getParentID())
	return cn
}








func RegisterSubChannel(parent *Channel, ref string) *SubChannel {
	id := IDGen.genID()
	sc := &SubChannel{
		ID:      id,
		RefName: ref,
		parent:  parent,
	}

	if !IsOn() {
		return sc
	}

	sc.sockets = make(map[int64]string)
	sc.trace = &ChannelTrace{CreationTime: time.Now(), Events: make([]*traceEvent, 0, getMaxTraceEntry())}
	db.addSubChannel(id, sc, parent.ID)
	return sc
}





func RegisterServer(ref string) *Server {
	id := IDGen.genID()
	if !IsOn() {
		return &Server{ID: id}
	}

	svr := &Server{
		RefName:       ref,
		sockets:       make(map[int64]string),
		listenSockets: make(map[int64]string),
		ID:            id,
	}
	db.addServer(id, svr)
	return svr
}







func RegisterSocket(skt *Socket) *Socket {
	skt.ID = IDGen.genID()
	if IsOn() {
		db.addSocket(skt)
	}
	return skt
}





func RemoveEntry(id int64) {
	if !IsOn() {
		return
	}
	db.removeEntry(id)
}


type IDGenerator struct {
	id int64
}



func (i *IDGenerator) Reset() {
	atomic.StoreInt64(&i.id, 0)
}

func (i *IDGenerator) genID() int64 {
	return atomic.AddInt64(&i.id, 1)
}




type Identifier interface {
	Entity
	channelzIdentifier()
}
