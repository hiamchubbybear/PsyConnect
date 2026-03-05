package kafka

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	heartbeatAPI "github.com/segmentio/kafka-go/protocol/heartbeat"
)


type HeartbeatRequest struct {
	
	Addr net.Addr

	
	GroupID string

	
	GenerationID int32

	
	MemberID string

	
	GroupInstanceID string
}


type HeartbeatResponse struct {
	
	Error error

	
	
	
	
	Throttle time.Duration
}

type heartbeatRequestV0 struct {
	
	GroupID string

	
	GenerationID int32

	
	MemberID string
}


func (c *Client) Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &heartbeatAPI.Request{
		GroupID:         req.GroupID,
		GenerationID:    req.GenerationID,
		MemberID:        req.MemberID,
		GroupInstanceID: req.GroupInstanceID,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).Heartbeat: %w", err)
	}

	res := m.(*heartbeatAPI.Response)

	ret := &HeartbeatResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
	}

	if res.ErrorCode != 0 {
		ret.Error = Error(res.ErrorCode)
	}

	return ret, nil
}

func (t heartbeatRequestV0) size() int32 {
	return sizeofString(t.GroupID) +
		sizeofInt32(t.GenerationID) +
		sizeofString(t.MemberID)
}

func (t heartbeatRequestV0) writeTo(wb *writeBuffer) {
	wb.writeString(t.GroupID)
	wb.writeInt32(t.GenerationID)
	wb.writeString(t.MemberID)
}

type heartbeatResponseV0 struct {
	
	ErrorCode int16
}

func (t heartbeatResponseV0) size() int32 {
	return sizeofInt16(t.ErrorCode)
}

func (t heartbeatResponseV0) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.ErrorCode)
}

func (t *heartbeatResponseV0) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	if remain, err = readInt16(r, sz, &t.ErrorCode); err != nil {
		return
	}
	return
}
