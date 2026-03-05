package kafka

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/segmentio/kafka-go/protocol/describegroups"
)


type DescribeGroupsRequest struct {
	
	Addr net.Addr

	
	GroupIDs []string
}


type DescribeGroupsResponse struct {
	
	Groups []DescribeGroupsResponseGroup
}


type DescribeGroupsResponseGroup struct {
	
	
	Error error

	
	GroupID string

	
	GroupState string

	
	Members []DescribeGroupsResponseMember
}


type DescribeGroupsResponseMember struct {
	
	MemberID string

	
	ClientID string

	
	ClientHost string

	
	MemberMetadata DescribeGroupsResponseMemberMetadata

	
	MemberAssignments DescribeGroupsResponseAssignments
}


type DescribeGroupsResponseMemberMetadata struct {
	
	Version int

	
	Topics []string

	
	UserData []byte

	
	
	OwnedPartitions []DescribeGroupsResponseMemberMetadataOwnedPartition
}

type DescribeGroupsResponseMemberMetadataOwnedPartition struct {
	
	Topic string

	
	Partitions []int
}


type DescribeGroupsResponseAssignments struct {
	
	Version int

	
	Topics []GroupMemberTopic

	
	UserData []byte
}



type GroupMemberTopic struct {
	
	Topic string

	
	Partitions []int
}




func (c *Client) DescribeGroups(
	ctx context.Context,
	req *DescribeGroupsRequest,
) (*DescribeGroupsResponse, error) {
	protoResp, err := c.roundTrip(
		ctx,
		req.Addr,
		&describegroups.Request{
			Groups: req.GroupIDs,
		},
	)
	if err != nil {
		return nil, err
	}
	apiResp := protoResp.(*describegroups.Response)
	resp := &DescribeGroupsResponse{}

	for _, apiGroup := range apiResp.Groups {
		group := DescribeGroupsResponseGroup{
			Error:      makeError(apiGroup.ErrorCode, ""),
			GroupID:    apiGroup.GroupID,
			GroupState: apiGroup.GroupState,
		}

		for _, member := range apiGroup.Members {
			decodedMetadata, err := decodeMemberMetadata(member.MemberMetadata)
			if err != nil {
				return nil, err
			}
			decodedAssignments, err := decodeMemberAssignments(member.MemberAssignment)
			if err != nil {
				return nil, err
			}

			group.Members = append(group.Members, DescribeGroupsResponseMember{
				MemberID:          member.MemberID,
				ClientID:          member.ClientID,
				ClientHost:        member.ClientHost,
				MemberAssignments: decodedAssignments,
				MemberMetadata:    decodedMetadata,
			})
		}
		resp.Groups = append(resp.Groups, group)
	}

	return resp, nil
}






func decodeMemberMetadata(rawMetadata []byte) (DescribeGroupsResponseMemberMetadata, error) {
	mm := DescribeGroupsResponseMemberMetadata{}

	if len(rawMetadata) == 0 {
		return mm, nil
	}

	buf := bytes.NewBuffer(rawMetadata)
	bufReader := bufio.NewReader(buf)
	remain := len(rawMetadata)

	var err error
	var version16 int16

	if remain, err = readInt16(bufReader, remain, &version16); err != nil {
		return mm, err
	}
	mm.Version = int(version16)

	if remain, err = readStringArray(bufReader, remain, &mm.Topics); err != nil {
		return mm, err
	}
	if remain, err = readBytes(bufReader, remain, &mm.UserData); err != nil {
		return mm, err
	}

	if mm.Version == 1 && remain > 0 {
		fn := func(r *bufio.Reader, size int) (fnRemain int, fnErr error) {
			op := DescribeGroupsResponseMemberMetadataOwnedPartition{}
			if fnRemain, fnErr = readString(r, size, &op.Topic); fnErr != nil {
				return
			}

			ps := []int32{}
			if fnRemain, fnErr = readInt32Array(r, fnRemain, &ps); fnErr != nil {
				return
			}

			for _, p := range ps {
				op.Partitions = append(op.Partitions, int(p))
			}

			mm.OwnedPartitions = append(mm.OwnedPartitions, op)
			return
		}

		if remain, err = readArrayWith(bufReader, remain, fn); err != nil {
			return mm, err
		}
	}

	if remain != 0 {
		return mm, fmt.Errorf("Got non-zero number of bytes remaining: %d", remain)
	}

	return mm, nil
}






func decodeMemberAssignments(rawAssignments []byte) (DescribeGroupsResponseAssignments, error) {
	ma := DescribeGroupsResponseAssignments{}

	if len(rawAssignments) == 0 {
		return ma, nil
	}

	buf := bytes.NewBuffer(rawAssignments)
	bufReader := bufio.NewReader(buf)
	remain := len(rawAssignments)

	var err error
	var version16 int16

	if remain, err = readInt16(bufReader, remain, &version16); err != nil {
		return ma, err
	}
	ma.Version = int(version16)

	fn := func(r *bufio.Reader, size int) (fnRemain int, fnErr error) {
		item := GroupMemberTopic{}

		if fnRemain, fnErr = readString(r, size, &item.Topic); fnErr != nil {
			return
		}

		partitions := []int32{}

		if fnRemain, fnErr = readInt32Array(r, fnRemain, &partitions); fnErr != nil {
			return
		}
		for _, partition := range partitions {
			item.Partitions = append(item.Partitions, int(partition))
		}

		ma.Topics = append(ma.Topics, item)
		return
	}
	if remain, err = readArrayWith(bufReader, remain, fn); err != nil {
		return ma, err
	}

	if remain, err = readBytes(bufReader, remain, &ma.UserData); err != nil {
		return ma, err
	}

	if remain != 0 {
		return ma, fmt.Errorf("Got non-zero number of bytes remaining: %d", remain)
	}

	return ma, nil
}



func readInt32Array(r *bufio.Reader, sz int, v *[]int32) (remain int, err error) {
	var content []int32
	fn := func(r *bufio.Reader, size int) (fnRemain int, fnErr error) {
		var value int32
		if fnRemain, fnErr = readInt32(r, size, &value); fnErr != nil {
			return
		}
		content = append(content, value)
		return
	}
	if remain, err = readArrayWith(r, sz, fn); err != nil {
		return
	}

	*v = content
	return
}
