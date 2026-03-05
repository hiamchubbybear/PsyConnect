

package channelz

import (
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc/credentials"
)



type SocketMetrics struct {
	
	StreamsStarted atomic.Int64
	
	
	
	StreamsSucceeded atomic.Int64
	
	
	
	StreamsFailed atomic.Int64
	
	MessagesSent     atomic.Int64
	MessagesReceived atomic.Int64
	
	
	KeepAlivesSent atomic.Int64
	
	
	LastLocalStreamCreatedTimestamp atomic.Int64
	
	
	LastRemoteStreamCreatedTimestamp atomic.Int64
	
	LastMessageSentTimestamp atomic.Int64
	
	LastMessageReceivedTimestamp atomic.Int64
}



type EphemeralSocketMetrics struct {
	
	
	
	LocalFlowControlWindow int64
	
	
	
	RemoteFlowControlWindow int64
}


type SocketType string


const (
	SocketTypeNormal = "NormalSocket"
	SocketTypeListen = "ListenSocket"
)




type Socket struct {
	Entity
	SocketType       SocketType
	ID               int64
	Parent           Entity
	cm               *channelMap
	SocketMetrics    SocketMetrics
	EphemeralMetrics func() *EphemeralSocketMetrics

	RefName string
	
	LocalAddr net.Addr
	
	RemoteAddr net.Addr
	
	
	RemoteName string
	
	SocketOptions *SocketOptionData
	
	Security credentials.ChannelzSecurityValue
}



func (ls *Socket) String() string {
	return fmt.Sprintf("%s %s #%d", ls.Parent, ls.SocketType, ls.ID)
}

func (ls *Socket) id() int64 {
	return ls.ID
}

func (ls *Socket) addChild(id int64, e entry) {
	logger.Errorf("cannot add a child (id = %d) of type %T to a listen socket", id, e)
}

func (ls *Socket) deleteChild(id int64) {
	logger.Errorf("cannot delete a child (id = %d) from a listen socket", id)
}

func (ls *Socket) triggerDelete() {
	ls.cm.deleteEntry(ls.ID)
	ls.Parent.(entry).deleteChild(ls.ID)
}

func (ls *Socket) deleteSelfIfReady() {
	logger.Errorf("cannot call deleteSelfIfReady on a listen socket")
}

func (ls *Socket) getParentID() int64 {
	return ls.Parent.id()
}
