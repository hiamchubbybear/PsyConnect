





package logger

import (
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CommandFailed                    = "Command failed"
	CommandStarted                   = "Command started"
	CommandSucceeded                 = "Command succeeded"
	ConnectionPoolCreated            = "Connection pool created"
	ConnectionPoolReady              = "Connection pool ready"
	ConnectionPoolCleared            = "Connection pool cleared"
	ConnectionPoolClosed             = "Connection pool closed"
	ConnectionCreated                = "Connection created"
	ConnectionReady                  = "Connection ready"
	ConnectionClosed                 = "Connection closed"
	ConnectionCheckoutStarted        = "Connection checkout started"
	ConnectionCheckoutFailed         = "Connection checkout failed"
	ConnectionCheckedOut             = "Connection checked out"
	ConnectionCheckedIn              = "Connection checked in"
	ServerSelectionFailed            = "Server selection failed"
	ServerSelectionStarted           = "Server selection started"
	ServerSelectionSucceeded         = "Server selection succeeded"
	ServerSelectionWaiting           = "Waiting for suitable server to become available"
	TopologyClosed                   = "Stopped topology monitoring"
	TopologyDescriptionChanged       = "Topology description changed"
	TopologyOpening                  = "Starting topology monitoring"
	TopologyServerClosed             = "Stopped server monitoring"
	TopologyServerHeartbeatFailed    = "Server heartbeat failed"
	TopologyServerHeartbeatStarted   = "Server heartbeat started"
	TopologyServerHeartbeatSucceeded = "Server heartbeat succeeded"
	TopologyServerOpening            = "Starting server monitoring"
)

const (
	KeyAwaited             = "awaited"
	KeyCommand             = "command"
	KeyCommandName         = "commandName"
	KeyDatabaseName        = "databaseName"
	KeyDriverConnectionID  = "driverConnectionId"
	KeyDurationMS          = "durationMS"
	KeyError               = "error"
	KeyFailure             = "failure"
	KeyMaxConnecting       = "maxConnecting"
	KeyMaxIdleTimeMS       = "maxIdleTimeMS"
	KeyMaxPoolSize         = "maxPoolSize"
	KeyMessage             = "message"
	KeyMinPoolSize         = "minPoolSize"
	KeyNewDescription      = "newDescription"
	KeyOperation           = "operation"
	KeyOperationID         = "operationId"
	KeyPreviousDescription = "previousDescription"
	KeyRemainingTimeMS     = "remainingTimeMS"
	KeyReason              = "reason"
	KeyReply               = "reply"
	KeyRequestID           = "requestId"
	KeySelector            = "selector"
	KeyServerConnectionID  = "serverConnectionId"
	KeyServerHost          = "serverHost"
	KeyServerPort          = "serverPort"
	KeyServiceID           = "serviceId"
	KeyTimestamp           = "timestamp"
	KeyTopologyDescription = "topologyDescription"
	KeyTopologyID          = "topologyId"
)


type KeyValues []interface{}


func (kvs *KeyValues) Add(key string, value interface{}) {
	*kvs = append(*kvs, key, value)
}

const (
	ReasonConnClosedStale              = "Connection became stale because the pool was cleared"
	ReasonConnClosedIdle               = "Connection has been available but unused for longer than the configured max idle time"
	ReasonConnClosedError              = "An error occurred while using the connection"
	ReasonConnClosedPoolClosed         = "Connection pool was closed"
	ReasonConnCheckoutFailedTimout     = "Wait queue timeout elapsed without a connection becoming available"
	ReasonConnCheckoutFailedError      = "An error occurred while trying to establish a new connection"
	ReasonConnCheckoutFailedPoolClosed = "Connection pool was closed"
)



type Component int

const (
	
	ComponentAll Component = iota

	
	ComponentCommand

	
	ComponentTopology

	
	ComponentServerSelection

	
	ComponentConnection
)

const (
	mongoDBLogAllEnvVar             = "MONGODB_LOG_ALL"
	mongoDBLogCommandEnvVar         = "MONGODB_LOG_COMMAND"
	mongoDBLogTopologyEnvVar        = "MONGODB_LOG_TOPOLOGY"
	mongoDBLogServerSelectionEnvVar = "MONGODB_LOG_SERVER_SELECTION"
	mongoDBLogConnectionEnvVar      = "MONGODB_LOG_CONNECTION"
)

var componentEnvVarMap = map[string]Component{
	mongoDBLogAllEnvVar:             ComponentAll,
	mongoDBLogCommandEnvVar:         ComponentCommand,
	mongoDBLogTopologyEnvVar:        ComponentTopology,
	mongoDBLogServerSelectionEnvVar: ComponentServerSelection,
	mongoDBLogConnectionEnvVar:      ComponentConnection,
}



func EnvHasComponentVariables() bool {
	for envVar := range componentEnvVarMap {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}



type Command struct {
	
	DriverConnectionID uint64              
	Name               string              
	DatabaseName       string              
	Message            string              
	OperationID        int32               
	RequestID          int64               
	ServerConnectionID *int64              
	ServerHost         string              
	ServerPort         string              
	ServiceID          *primitive.ObjectID 
}




func SerializeCommand(cmd Command, extraKeysAndValues ...interface{}) KeyValues {
	
	keysAndValues := KeyValues{
		KeyCommandName, cmd.Name,
		KeyDatabaseName, cmd.DatabaseName,
		KeyDriverConnectionID, cmd.DriverConnectionID,
		KeyMessage, cmd.Message,
		KeyOperationID, cmd.OperationID,
		KeyRequestID, cmd.RequestID,
		KeyServerHost, cmd.ServerHost,
	}

	
	for i := 0; i < len(extraKeysAndValues); i += 2 {
		keysAndValues.Add(extraKeysAndValues[i].(string), extraKeysAndValues[i+1])
	}

	port, err := strconv.ParseInt(cmd.ServerPort, 10, 32)
	if err == nil {
		keysAndValues.Add(KeyServerPort, port)
	}

	
	if cmd.ServerConnectionID != nil {
		keysAndValues.Add(KeyServerConnectionID, *cmd.ServerConnectionID)
	}

	
	if cmd.ServiceID != nil {
		keysAndValues.Add(KeyServiceID, cmd.ServiceID.Hex())
	}

	return keysAndValues
}


type Connection struct {
	Message    string 
	ServerHost string 
	ServerPort string 
}



func SerializeConnection(conn Connection, extraKeysAndValues ...interface{}) KeyValues {
	
	keysAndValues := KeyValues{
		KeyMessage, conn.Message,
		KeyServerHost, conn.ServerHost,
	}

	
	for i := 0; i < len(extraKeysAndValues); i += 2 {
		keysAndValues.Add(extraKeysAndValues[i].(string), extraKeysAndValues[i+1])
	}

	port, err := strconv.ParseInt(conn.ServerPort, 10, 32)
	if err == nil {
		keysAndValues.Add(KeyServerPort, port)
	}

	return keysAndValues
}


type Server struct {
	DriverConnectionID uint64             
	TopologyID         primitive.ObjectID 
	Message            string             
	ServerConnectionID *int64             
	ServerHost         string             
	ServerPort         string             
}



func SerializeServer(srv Server, extraKV ...interface{}) KeyValues {
	
	keysAndValues := KeyValues{
		KeyDriverConnectionID, srv.DriverConnectionID,
		KeyMessage, srv.Message,
		KeyServerHost, srv.ServerHost,
		KeyTopologyID, srv.TopologyID.Hex(),
	}

	if connID := srv.ServerConnectionID; connID != nil {
		keysAndValues.Add(KeyServerConnectionID, *connID)
	}

	port, err := strconv.ParseInt(srv.ServerPort, 10, 32)
	if err == nil {
		keysAndValues.Add(KeyServerPort, port)
	}

	
	for i := 0; i < len(extraKV); i += 2 {
		keysAndValues.Add(extraKV[i].(string), extraKV[i+1])
	}

	return keysAndValues
}



type ServerSelection struct {
	Selector            string
	OperationID         *int32
	Operation           string
	TopologyDescription string
}



func SerializeServerSelection(srvSelection ServerSelection, extraKV ...interface{}) KeyValues {
	keysAndValues := KeyValues{
		KeySelector, srvSelection.Selector,
		KeyOperation, srvSelection.Operation,
		KeyTopologyDescription, srvSelection.TopologyDescription,
	}

	if srvSelection.OperationID != nil {
		keysAndValues.Add(KeyOperationID, *srvSelection.OperationID)
	}

	
	for i := 0; i < len(extraKV); i += 2 {
		keysAndValues.Add(extraKV[i].(string), extraKV[i+1])
	}

	return keysAndValues
}


type Topology struct {
	ID      primitive.ObjectID 
	Message string             
}



func SerializeTopology(topo Topology, extraKV ...interface{}) KeyValues {
	keysAndValues := KeyValues{
		KeyTopologyID, topo.ID.Hex(),
	}

	
	for i := 0; i < len(extraKV); i += 2 {
		keysAndValues.Add(extraKV[i].(string), extraKV[i+1])
	}

	return keysAndValues
}
