

package channelz

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/grpclog"
)

const (
	defaultMaxTraceEntry int32 = 30
)

var maxTraceEntry = defaultMaxTraceEntry



func SetMaxTraceEntry(i int32) {
	atomic.StoreInt32(&maxTraceEntry, i)
}



func ResetMaxTraceEntryToDefault() {
	atomic.StoreInt32(&maxTraceEntry, defaultMaxTraceEntry)
}

func getMaxTraceEntry() int {
	i := atomic.LoadInt32(&maxTraceEntry)
	return int(i)
}


type traceEvent struct {
	
	Desc string
	
	Severity Severity
	
	Timestamp time.Time
	
	
	
	RefID int64
	
	RefName string
	
	RefType RefChannelType
}






type TraceEvent struct {
	Desc     string
	Severity Severity
	Parent   *TraceEvent
}




type ChannelTrace struct {
	cm          *channelMap
	clearCalled bool
	
	CreationTime time.Time
	
	
	EventNum int64
	mu       sync.Mutex
	
	
	Events []*traceEvent
}

func (c *ChannelTrace) copy() *ChannelTrace {
	return &ChannelTrace{
		CreationTime: c.CreationTime,
		EventNum:     c.EventNum,
		Events:       append(([]*traceEvent)(nil), c.Events...),
	}
}

func (c *ChannelTrace) append(e *traceEvent) {
	c.mu.Lock()
	if len(c.Events) == getMaxTraceEntry() {
		del := c.Events[0]
		c.Events = c.Events[1:]
		if del.RefID != 0 {
			
			go func() {
				
				c.cm.mu.Lock()
				c.cm.decrTraceRefCount(del.RefID)
				c.cm.mu.Unlock()
			}()
		}
	}
	e.Timestamp = time.Now()
	c.Events = append(c.Events, e)
	c.EventNum++
	c.mu.Unlock()
}

func (c *ChannelTrace) clear() {
	if c.clearCalled {
		return
	}
	c.clearCalled = true
	c.mu.Lock()
	for _, e := range c.Events {
		if e.RefID != 0 {
			
			c.cm.decrTraceRefCount(e.RefID)
		}
	}
	c.mu.Unlock()
}




type Severity int

const (
	
	CtUnknown Severity = iota
	
	CtInfo
	
	CtWarning
	
	CtError
)


type RefChannelType int

const (
	
	RefUnknown RefChannelType = iota
	
	RefChannel
	
	RefSubChannel
	
	RefServer
	
	RefListenSocket
	
	RefNormalSocket
)

var refChannelTypeToString = map[RefChannelType]string{
	RefUnknown:      "Unknown",
	RefChannel:      "Channel",
	RefSubChannel:   "SubChannel",
	RefServer:       "Server",
	RefListenSocket: "ListenSocket",
	RefNormalSocket: "NormalSocket",
}


func (r RefChannelType) String() string {
	return refChannelTypeToString[r]
}





func AddTraceEvent(l grpclog.DepthLoggerV2, e Entity, depth int, desc *TraceEvent) {
	
	d := fmt.Sprintf("[%s]%s", e, desc.Desc)
	switch desc.Severity {
	case CtUnknown, CtInfo:
		l.InfoDepth(depth+1, d)
	case CtWarning:
		l.WarningDepth(depth+1, d)
	case CtError:
		l.ErrorDepth(depth+1, d)
	}

	if getMaxTraceEntry() == 0 {
		return
	}
	if IsOn() {
		db.traceEvent(e.id(), desc)
	}
}
