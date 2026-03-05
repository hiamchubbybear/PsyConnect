





package topology

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/description"
)


type ConnectionError struct {
	ConnectionID string
	Wrapped      error

	
	
	init    bool
	message string
}


func (e ConnectionError) Error() string {
	message := e.message
	if e.init {
		fullMsg := "error occurred during connection handshake"
		if message != "" {
			fullMsg = fmt.Sprintf("%s: %s", fullMsg, message)
		}
		message = fullMsg
	}
	if e.Wrapped != nil && message != "" {
		return fmt.Sprintf("connection(%s) %s: %s", e.ConnectionID, message, e.Wrapped.Error())
	}
	if e.Wrapped != nil {
		return fmt.Sprintf("connection(%s) %s", e.ConnectionID, e.Wrapped.Error())
	}
	return fmt.Sprintf("connection(%s) %s", e.ConnectionID, message)
}


func (e ConnectionError) Unwrap() error {
	return e.Wrapped
}


type ServerSelectionError struct {
	Desc    description.Topology
	Wrapped error
}


func (e ServerSelectionError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("server selection error: %s, current topology: { %s }", e.Wrapped.Error(), e.Desc.String())
	}
	return fmt.Sprintf("server selection error: current topology: { %s }", e.Desc.String())
}


func (e ServerSelectionError) Unwrap() error {
	return e.Wrapped
}


type WaitQueueTimeoutError struct {
	Wrapped              error
	pinnedConnections    *pinnedConnections
	maxPoolSize          uint64
	totalConnections     int
	availableConnections int
	waitDuration         time.Duration
}

type pinnedConnections struct {
	cursorConnections      uint64
	transactionConnections uint64
}


func (w WaitQueueTimeoutError) Error() string {
	errorMsg := "timed out while checking out a connection from connection pool"
	switch {
	case w.Wrapped == nil:
	case errors.Is(w.Wrapped, context.Canceled):
		errorMsg = fmt.Sprintf(
			"%s: %s",
			"canceled while checking out a connection from connection pool",
			w.Wrapped.Error(),
		)
	default:
		errorMsg = fmt.Sprintf(
			"%s: %s",
			errorMsg,
			w.Wrapped.Error(),
		)
	}

	msg := fmt.Sprintf("%s; total connections: %d, maxPoolSize: %d, ", errorMsg, w.totalConnections, w.maxPoolSize)
	if pinnedConnections := w.pinnedConnections; pinnedConnections != nil {
		openConnectionCount := uint64(w.totalConnections) -
			pinnedConnections.cursorConnections -
			pinnedConnections.transactionConnections
		msg += fmt.Sprintf("connections in use by cursors: %d, connections in use by transactions: %d, connections in use by other operations: %d, ",
			pinnedConnections.cursorConnections,
			pinnedConnections.transactionConnections,
			openConnectionCount,
		)
	}
	msg += fmt.Sprintf("idle connections: %d, wait duration: %s", w.availableConnections, w.waitDuration.String())
	return msg
}


func (w WaitQueueTimeoutError) Unwrap() error {
	return w.Wrapped
}
