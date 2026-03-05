





package event 

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
)


type CommandStartedEvent struct {
	Command      bson.Raw
	DatabaseName string
	CommandName  string
	RequestID    int64
	ConnectionID string
	
	
	
	
	
	ServerConnectionID *int32
	
	
	ServerConnectionID64 *int64
	
	
	ServiceID *primitive.ObjectID
}


type CommandFinishedEvent struct {
	
	DurationNanos int64
	Duration      time.Duration
	CommandName   string
	DatabaseName  string
	RequestID     int64
	ConnectionID  string
	
	
	
	
	
	ServerConnectionID *int32
	
	
	ServerConnectionID64 *int64
	
	
	ServiceID *primitive.ObjectID
}


type CommandSucceededEvent struct {
	CommandFinishedEvent
	Reply bson.Raw
}


type CommandFailedEvent struct {
	CommandFinishedEvent
	Failure string
}


type CommandMonitor struct {
	Started   func(context.Context, *CommandStartedEvent)
	Succeeded func(context.Context, *CommandSucceededEvent)
	Failed    func(context.Context, *CommandFailedEvent)
}


const (
	ReasonIdle              = "idle"
	ReasonPoolClosed        = "poolClosed"
	ReasonStale             = "stale"
	ReasonConnectionErrored = "connectionError"
	ReasonTimedOut          = "timeout"
	ReasonError             = "error"
)


const (
	PoolCreated        = "ConnectionPoolCreated"
	PoolReady          = "ConnectionPoolReady"
	PoolCleared        = "ConnectionPoolCleared"
	PoolClosedEvent    = "ConnectionPoolClosed"
	ConnectionCreated  = "ConnectionCreated"
	ConnectionReady    = "ConnectionReady"
	ConnectionClosed   = "ConnectionClosed"
	GetStarted         = "ConnectionCheckOutStarted"
	GetFailed          = "ConnectionCheckOutFailed"
	GetSucceeded       = "ConnectionCheckedOut"
	ConnectionReturned = "ConnectionCheckedIn"
)


type MonitorPoolOptions struct {
	MaxPoolSize        uint64 `json:"maxPoolSize"`
	MinPoolSize        uint64 `json:"minPoolSize"`
	WaitQueueTimeoutMS uint64 `json:"maxIdleTimeMS"`
}


type PoolEvent struct {
	Type         string              `json:"type"`
	Address      string              `json:"address"`
	ConnectionID uint64              `json:"connectionId"`
	PoolOptions  *MonitorPoolOptions `json:"options"`
	Duration     time.Duration       `json:"duration"`
	Reason       string              `json:"reason"`
	
	
	ServiceID    *primitive.ObjectID `json:"serviceId"`
	Interruption bool                `json:"interruptInUseConnections"`
	Error        error               `json:"error"`
}


type PoolMonitor struct {
	Event func(*PoolEvent)
}


type ServerDescriptionChangedEvent struct {
	Address             address.Address
	TopologyID          primitive.ObjectID 
	PreviousDescription description.Server
	NewDescription      description.Server
}


type ServerOpeningEvent struct {
	Address    address.Address
	TopologyID primitive.ObjectID 
}


type ServerClosedEvent struct {
	Address    address.Address
	TopologyID primitive.ObjectID 
}


type TopologyDescriptionChangedEvent struct {
	TopologyID          primitive.ObjectID 
	PreviousDescription description.Topology
	NewDescription      description.Topology
}


type TopologyOpeningEvent struct {
	TopologyID primitive.ObjectID 
}


type TopologyClosedEvent struct {
	TopologyID primitive.ObjectID 
}


type ServerHeartbeatStartedEvent struct {
	ConnectionID string 
	Awaited      bool   
}


type ServerHeartbeatSucceededEvent struct {
	
	DurationNanos int64
	Duration      time.Duration
	Reply         description.Server
	ConnectionID  string 
	Awaited       bool   
}


type ServerHeartbeatFailedEvent struct {
	
	DurationNanos int64
	Duration      time.Duration
	Failure       error
	ConnectionID  string 
	Awaited       bool   
}





type ServerMonitor struct {
	ServerDescriptionChanged func(*ServerDescriptionChangedEvent)
	ServerOpening            func(*ServerOpeningEvent)
	ServerClosed             func(*ServerClosedEvent)
	
	
	TopologyDescriptionChanged func(*TopologyDescriptionChangedEvent)
	TopologyOpening            func(*TopologyOpeningEvent)
	TopologyClosed             func(*TopologyClosedEvent)
	ServerHeartbeatStarted     func(*ServerHeartbeatStartedEvent)
	ServerHeartbeatSucceeded   func(*ServerHeartbeatSucceededEvent)
	ServerHeartbeatFailed      func(*ServerHeartbeatFailedEvent)
}
