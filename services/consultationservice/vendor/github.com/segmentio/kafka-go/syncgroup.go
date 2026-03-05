package kafka

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol"
	"github.com/segmentio/kafka-go/protocol/consumer"
	"github.com/segmentio/kafka-go/protocol/syncgroup"
)


type SyncGroupRequest struct {
	
	Addr net.Addr

	
	GroupID string

	
	GenerationID int

	
	MemberID string

	
	GroupInstanceID string

	
	ProtocolType string

	
	ProtocolName string

	
	Assignments []SyncGroupRequestAssignment
}


type SyncGroupRequestAssignment struct {
	
	MemberID string

	
	Assignment GroupProtocolAssignment
}


type SyncGroupResponse struct {
	
	
	
	
	Error error

	
	Throttle time.Duration

	
	ProtocolType string

	
	ProtocolName string

	
	Assignment GroupProtocolAssignment
}


type GroupProtocolAssignment struct {
	
	AssignedPartitions map[string][]int

	
	UserData []byte
}


func (c *Client) SyncGroup(ctx context.Context, req *SyncGroupRequest) (*SyncGroupResponse, error) {
	syncGroup := syncgroup.Request{
		GroupID:         req.GroupID,
		GenerationID:    int32(req.GenerationID),
		MemberID:        req.MemberID,
		GroupInstanceID: req.GroupInstanceID,
		ProtocolType:    req.ProtocolType,
		ProtocolName:    req.ProtocolName,
		Assignments:     make([]syncgroup.RequestAssignment, 0, len(req.Assignments)),
	}

	for _, assignment := range req.Assignments {
		assign := consumer.Assignment{
			Version:            consumer.MaxVersionSupported,
			AssignedPartitions: make([]consumer.TopicPartition, 0, len(assignment.Assignment.AssignedPartitions)),
			UserData:           assignment.Assignment.UserData,
		}

		for topic, partitions := range assignment.Assignment.AssignedPartitions {
			tp := consumer.TopicPartition{
				Topic:      topic,
				Partitions: make([]int32, 0, len(partitions)),
			}
			for _, partition := range partitions {
				tp.Partitions = append(tp.Partitions, int32(partition))
			}
			assign.AssignedPartitions = append(assign.AssignedPartitions, tp)
		}

		assignBytes, err := protocol.Marshal(consumer.MaxVersionSupported, assign)
		if err != nil {
			return nil, fmt.Errorf("kafka.(*Client).SyncGroup: %w", err)
		}

		syncGroup.Assignments = append(syncGroup.Assignments, syncgroup.RequestAssignment{
			MemberID:   assignment.MemberID,
			Assignment: assignBytes,
		})
	}

	m, err := c.roundTrip(ctx, req.Addr, &syncGroup)
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).SyncGroup: %w", err)
	}

	r := m.(*syncgroup.Response)

	var assignment consumer.Assignment
	err = protocol.Unmarshal(r.Assignments, consumer.MaxVersionSupported, &assignment)
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).SyncGroup: %w", err)
	}

	res := &SyncGroupResponse{
		Throttle:     makeDuration(r.ThrottleTimeMS),
		Error:        makeError(r.ErrorCode, ""),
		ProtocolType: r.ProtocolType,
		ProtocolName: r.ProtocolName,
		Assignment: GroupProtocolAssignment{
			AssignedPartitions: make(map[string][]int, len(assignment.AssignedPartitions)),
			UserData:           assignment.UserData,
		},
	}
	partitions := map[string][]int{}
	for _, topicPartition := range assignment.AssignedPartitions {
		for _, partition := range topicPartition.Partitions {
			partitions[topicPartition.Topic] = append(partitions[topicPartition.Topic], int(partition))
		}
	}
	res.Assignment.AssignedPartitions = partitions

	return res, nil
}

type groupAssignment struct {
	Version  int16
	Topics   map[string][]int32
	UserData []byte
}

func (t groupAssignment) size() int32 {
	sz := sizeofInt16(t.Version) + sizeofInt16(int16(len(t.Topics)))

	for topic, partitions := range t.Topics {
		sz += sizeofString(topic) + sizeofInt32Array(partitions)
	}

	return sz + sizeofBytes(t.UserData)
}

func (t groupAssignment) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.Version)
	wb.writeInt32(int32(len(t.Topics)))

	for topic, partitions := range t.Topics {
		wb.writeString(topic)
		wb.writeInt32Array(partitions)
	}

	wb.writeBytes(t.UserData)
}

func (t *groupAssignment) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	
	
	
	if size == 0 {
		t.Topics = map[string][]int32{}
		return 0, nil
	}

	if remain, err = readInt16(r, size, &t.Version); err != nil {
		return
	}
	if remain, err = readMapStringInt32(r, remain, &t.Topics); err != nil {
		return
	}
	if remain, err = readBytes(r, remain, &t.UserData); err != nil {
		return
	}

	return
}

func (t groupAssignment) bytes() []byte {
	buf := bytes.NewBuffer(nil)
	t.writeTo(&writeBuffer{w: buf})
	return buf.Bytes()
}

type syncGroupRequestGroupAssignmentV0 struct {
	
	MemberID string

	
	
	
	MemberAssignments []byte
}

func (t syncGroupRequestGroupAssignmentV0) size() int32 {
	return sizeofString(t.MemberID) +
		sizeofBytes(t.MemberAssignments)
}

func (t syncGroupRequestGroupAssignmentV0) writeTo(wb *writeBuffer) {
	wb.writeString(t.MemberID)
	wb.writeBytes(t.MemberAssignments)
}

type syncGroupRequestV0 struct {
	
	GroupID string

	
	GenerationID int32

	
	MemberID string

	GroupAssignments []syncGroupRequestGroupAssignmentV0
}

func (t syncGroupRequestV0) size() int32 {
	return sizeofString(t.GroupID) +
		sizeofInt32(t.GenerationID) +
		sizeofString(t.MemberID) +
		sizeofArray(len(t.GroupAssignments), func(i int) int32 { return t.GroupAssignments[i].size() })
}

func (t syncGroupRequestV0) writeTo(wb *writeBuffer) {
	wb.writeString(t.GroupID)
	wb.writeInt32(t.GenerationID)
	wb.writeString(t.MemberID)
	wb.writeArray(len(t.GroupAssignments), func(i int) { t.GroupAssignments[i].writeTo(wb) })
}

type syncGroupResponseV0 struct {
	
	ErrorCode int16

	
	
	
	MemberAssignments []byte
}

func (t syncGroupResponseV0) size() int32 {
	return sizeofInt16(t.ErrorCode) +
		sizeofBytes(t.MemberAssignments)
}

func (t syncGroupResponseV0) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.ErrorCode)
	wb.writeBytes(t.MemberAssignments)
}

func (t *syncGroupResponseV0) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	if remain, err = readInt16(r, sz, &t.ErrorCode); err != nil {
		return
	}
	if remain, err = readBytes(r, remain, &t.MemberAssignments); err != nil {
		return
	}
	return
}
