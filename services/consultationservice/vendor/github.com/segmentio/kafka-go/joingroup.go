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
	"github.com/segmentio/kafka-go/protocol/joingroup"
)


type JoinGroupRequest struct {
	
	Addr net.Addr

	
	GroupID string

	
	
	SessionTimeout time.Duration

	
	RebalanceTimeout time.Duration

	
	MemberID string

	
	GroupInstanceID string

	
	ProtocolType string

	
	Protocols []GroupProtocol
}


type GroupProtocol struct {
	
	Name string

	
	Metadata GroupProtocolSubscription
}

type GroupProtocolSubscription struct {
	
	Topics []string

	
	UserData []byte

	
	OwnedPartitions map[string][]int
}


type JoinGroupResponse struct {
	
	
	
	
	Error error

	
	Throttle time.Duration

	
	GenerationID int

	
	ProtocolName string

	
	ProtocolType string

	
	LeaderID string

	
	MemberID string

	
	Members []JoinGroupResponseMember
}


type JoinGroupResponseMember struct {
	
	ID string

	
	GroupInstanceID string

	
	Metadata GroupProtocolSubscription
}


func (c *Client) JoinGroup(ctx context.Context, req *JoinGroupRequest) (*JoinGroupResponse, error) {
	joinGroup := joingroup.Request{
		GroupID:            req.GroupID,
		SessionTimeoutMS:   int32(req.SessionTimeout.Milliseconds()),
		RebalanceTimeoutMS: int32(req.RebalanceTimeout.Milliseconds()),
		MemberID:           req.MemberID,
		GroupInstanceID:    req.GroupInstanceID,
		ProtocolType:       req.ProtocolType,
		Protocols:          make([]joingroup.RequestProtocol, 0, len(req.Protocols)),
	}

	for _, proto := range req.Protocols {
		protoMeta := consumer.Subscription{
			Version:         consumer.MaxVersionSupported,
			Topics:          proto.Metadata.Topics,
			UserData:        proto.Metadata.UserData,
			OwnedPartitions: make([]consumer.TopicPartition, 0, len(proto.Metadata.OwnedPartitions)),
		}
		for topic, partitions := range proto.Metadata.OwnedPartitions {
			tp := consumer.TopicPartition{
				Topic:      topic,
				Partitions: make([]int32, 0, len(partitions)),
			}
			for _, partition := range partitions {
				tp.Partitions = append(tp.Partitions, int32(partition))
			}
			protoMeta.OwnedPartitions = append(protoMeta.OwnedPartitions, tp)
		}

		metaBytes, err := protocol.Marshal(consumer.MaxVersionSupported, protoMeta)
		if err != nil {
			return nil, fmt.Errorf("kafka.(*Client).JoinGroup: %w", err)
		}

		joinGroup.Protocols = append(joinGroup.Protocols, joingroup.RequestProtocol{
			Name:     proto.Name,
			Metadata: metaBytes,
		})
	}

	m, err := c.roundTrip(ctx, req.Addr, &joinGroup)
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).JoinGroup: %w", err)
	}

	r := m.(*joingroup.Response)

	res := &JoinGroupResponse{
		Error:        makeError(r.ErrorCode, ""),
		Throttle:     makeDuration(r.ThrottleTimeMS),
		GenerationID: int(r.GenerationID),
		ProtocolName: r.ProtocolName,
		ProtocolType: r.ProtocolType,
		LeaderID:     r.LeaderID,
		MemberID:     r.MemberID,
		Members:      make([]JoinGroupResponseMember, 0, len(r.Members)),
	}

	for _, member := range r.Members {
		var meta consumer.Subscription
		err = protocol.Unmarshal(member.Metadata, consumer.MaxVersionSupported, &meta)
		if err != nil {
			return nil, fmt.Errorf("kafka.(*Client).JoinGroup: %w", err)
		}
		subscription := GroupProtocolSubscription{
			Topics:          meta.Topics,
			UserData:        meta.UserData,
			OwnedPartitions: make(map[string][]int, len(meta.OwnedPartitions)),
		}
		for _, owned := range meta.OwnedPartitions {
			subscription.OwnedPartitions[owned.Topic] = make([]int, 0, len(owned.Partitions))
			for _, partition := range owned.Partitions {
				subscription.OwnedPartitions[owned.Topic] = append(subscription.OwnedPartitions[owned.Topic], int(partition))
			}
		}
		res.Members = append(res.Members, JoinGroupResponseMember{
			ID:              member.MemberID,
			GroupInstanceID: member.GroupInstanceID,
			Metadata:        subscription,
		})
	}

	return res, nil
}

type groupMetadata struct {
	Version  int16
	Topics   []string
	UserData []byte
}

func (t groupMetadata) size() int32 {
	return sizeofInt16(t.Version) +
		sizeofStringArray(t.Topics) +
		sizeofBytes(t.UserData)
}

func (t groupMetadata) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.Version)
	wb.writeStringArray(t.Topics)
	wb.writeBytes(t.UserData)
}

func (t groupMetadata) bytes() []byte {
	buf := bytes.NewBuffer(nil)
	t.writeTo(&writeBuffer{w: buf})
	return buf.Bytes()
}

func (t *groupMetadata) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt16(r, size, &t.Version); err != nil {
		return
	}
	if remain, err = readStringArray(r, remain, &t.Topics); err != nil {
		return
	}
	if remain, err = readBytes(r, remain, &t.UserData); err != nil {
		return
	}
	return
}

type joinGroupRequestGroupProtocolV1 struct {
	ProtocolName     string
	ProtocolMetadata []byte
}

func (t joinGroupRequestGroupProtocolV1) size() int32 {
	return sizeofString(t.ProtocolName) +
		sizeofBytes(t.ProtocolMetadata)
}

func (t joinGroupRequestGroupProtocolV1) writeTo(wb *writeBuffer) {
	wb.writeString(t.ProtocolName)
	wb.writeBytes(t.ProtocolMetadata)
}

type joinGroupRequestV1 struct {
	
	GroupID string

	
	
	SessionTimeout int32

	
	
	RebalanceTimeout int32

	
	
	MemberID string

	
	ProtocolType string

	
	GroupProtocols []joinGroupRequestGroupProtocolV1
}

func (t joinGroupRequestV1) size() int32 {
	return sizeofString(t.GroupID) +
		sizeofInt32(t.SessionTimeout) +
		sizeofInt32(t.RebalanceTimeout) +
		sizeofString(t.MemberID) +
		sizeofString(t.ProtocolType) +
		sizeofArray(len(t.GroupProtocols), func(i int) int32 { return t.GroupProtocols[i].size() })
}

func (t joinGroupRequestV1) writeTo(wb *writeBuffer) {
	wb.writeString(t.GroupID)
	wb.writeInt32(t.SessionTimeout)
	wb.writeInt32(t.RebalanceTimeout)
	wb.writeString(t.MemberID)
	wb.writeString(t.ProtocolType)
	wb.writeArray(len(t.GroupProtocols), func(i int) { t.GroupProtocols[i].writeTo(wb) })
}

type joinGroupResponseMemberV1 struct {
	
	MemberID       string
	MemberMetadata []byte
}

func (t joinGroupResponseMemberV1) size() int32 {
	return sizeofString(t.MemberID) +
		sizeofBytes(t.MemberMetadata)
}

func (t joinGroupResponseMemberV1) writeTo(wb *writeBuffer) {
	wb.writeString(t.MemberID)
	wb.writeBytes(t.MemberMetadata)
}

func (t *joinGroupResponseMemberV1) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readString(r, size, &t.MemberID); err != nil {
		return
	}
	if remain, err = readBytes(r, remain, &t.MemberMetadata); err != nil {
		return
	}
	return
}

type joinGroupResponseV1 struct {
	
	ErrorCode int16

	
	GenerationID int32

	
	GroupProtocol string

	
	LeaderID string

	
	MemberID string
	Members  []joinGroupResponseMemberV1
}

func (t joinGroupResponseV1) size() int32 {
	return sizeofInt16(t.ErrorCode) +
		sizeofInt32(t.GenerationID) +
		sizeofString(t.GroupProtocol) +
		sizeofString(t.LeaderID) +
		sizeofString(t.MemberID) +
		sizeofArray(len(t.MemberID), func(i int) int32 { return t.Members[i].size() })
}

func (t joinGroupResponseV1) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.ErrorCode)
	wb.writeInt32(t.GenerationID)
	wb.writeString(t.GroupProtocol)
	wb.writeString(t.LeaderID)
	wb.writeString(t.MemberID)
	wb.writeArray(len(t.Members), func(i int) { t.Members[i].writeTo(wb) })
}

func (t *joinGroupResponseV1) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt16(r, size, &t.ErrorCode); err != nil {
		return
	}
	if remain, err = readInt32(r, remain, &t.GenerationID); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.GroupProtocol); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.LeaderID); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.MemberID); err != nil {
		return
	}

	fn := func(r *bufio.Reader, size int) (fnRemain int, fnErr error) {
		var item joinGroupResponseMemberV1
		if fnRemain, fnErr = (&item).readFrom(r, size); fnErr != nil {
			return
		}
		t.Members = append(t.Members, item)
		return
	}
	if remain, err = readArrayWith(r, remain, fn); err != nil {
		return
	}

	return
}
