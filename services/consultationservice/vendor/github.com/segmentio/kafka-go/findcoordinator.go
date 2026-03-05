package kafka

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/findcoordinator"
)


type CoordinatorKeyType int8

const (
	
	CoordinatorKeyTypeConsumer CoordinatorKeyType = 0

	
	CoordinatorKeyTypeTransaction CoordinatorKeyType = 1
)


type FindCoordinatorRequest struct {
	
	Addr net.Addr

	
	Key string

	
	KeyType CoordinatorKeyType
}


type FindCoordinatorResponseCoordinator struct {
	
	NodeID int

	
	Host string

	
	Port int
}


type FindCoordinatorResponse struct {
	
	Coordinator *FindCoordinatorResponseCoordinator

	
	Throttle time.Duration

	
	
	
	
	Error error
}



func (c *Client) FindCoordinator(ctx context.Context, req *FindCoordinatorRequest) (*FindCoordinatorResponse, error) {

	m, err := c.roundTrip(ctx, req.Addr, &findcoordinator.Request{
		Key:     req.Key,
		KeyType: int8(req.KeyType),
	})

	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).FindCoordinator: %w", err)
	}

	res := m.(*findcoordinator.Response)
	coordinator := &FindCoordinatorResponseCoordinator{
		NodeID: int(res.NodeID),
		Host:   res.Host,
		Port:   int(res.Port),
	}
	ret := &FindCoordinatorResponse{
		Throttle:    makeDuration(res.ThrottleTimeMs),
		Error:       makeError(res.ErrorCode, res.ErrorMessage),
		Coordinator: coordinator,
	}

	return ret, nil
}




type findCoordinatorRequestV0 struct {
	
	
	CoordinatorKey string
}

func (t findCoordinatorRequestV0) size() int32 {
	return sizeofString(t.CoordinatorKey)
}

func (t findCoordinatorRequestV0) writeTo(wb *writeBuffer) {
	wb.writeString(t.CoordinatorKey)
}

type findCoordinatorResponseCoordinatorV0 struct {
	
	NodeID int32

	
	Host string

	
	Port int32
}

func (t findCoordinatorResponseCoordinatorV0) size() int32 {
	return sizeofInt32(t.NodeID) +
		sizeofString(t.Host) +
		sizeofInt32(t.Port)
}

func (t findCoordinatorResponseCoordinatorV0) writeTo(wb *writeBuffer) {
	wb.writeInt32(t.NodeID)
	wb.writeString(t.Host)
	wb.writeInt32(t.Port)
}

func (t *findCoordinatorResponseCoordinatorV0) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt32(r, size, &t.NodeID); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.Host); err != nil {
		return
	}
	if remain, err = readInt32(r, remain, &t.Port); err != nil {
		return
	}
	return
}

type findCoordinatorResponseV0 struct {
	
	ErrorCode int16

	
	Coordinator findCoordinatorResponseCoordinatorV0
}

func (t findCoordinatorResponseV0) size() int32 {
	return sizeofInt16(t.ErrorCode) +
		t.Coordinator.size()
}

func (t findCoordinatorResponseV0) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.ErrorCode)
	t.Coordinator.writeTo(wb)
}

func (t *findCoordinatorResponseV0) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt16(r, size, &t.ErrorCode); err != nil {
		return
	}
	if remain, err = (&t.Coordinator).readFrom(r, remain); err != nil {
		return
	}
	return
}
