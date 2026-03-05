



package connectivity

import (
	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.Component("core")



type State int

func (s State) String() string {
	switch s {
	case Idle:
		return "IDLE"
	case Connecting:
		return "CONNECTING"
	case Ready:
		return "READY"
	case TransientFailure:
		return "TRANSIENT_FAILURE"
	case Shutdown:
		return "SHUTDOWN"
	default:
		logger.Errorf("unknown connectivity state: %d", s)
		return "INVALID_STATE"
	}
}

const (
	
	Idle State = iota
	
	Connecting
	
	Ready
	
	TransientFailure
	
	Shutdown
)




type ServingMode int

const (
	
	ServingModeStarting ServingMode = iota
	
	
	ServingModeServing
	
	
	
	
	ServingModeNotServing
)

func (s ServingMode) String() string {
	switch s {
	case ServingModeStarting:
		return "STARTING"
	case ServingModeServing:
		return "SERVING"
	case ServingModeNotServing:
		return "NOT_SERVING"
	default:
		logger.Errorf("unknown serving mode: %d", s)
		return "INVALID_MODE"
	}
}
