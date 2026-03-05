package kafka

import (
	"context"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/listpartitionreassignments"
)


type ListPartitionReassignmentsRequest struct {
	
	Addr net.Addr

	
	Topics map[string]ListPartitionReassignmentsRequestTopic

	
	Timeout time.Duration
}



type ListPartitionReassignmentsRequestTopic struct {
	
	PartitionIndexes []int
}


type ListPartitionReassignmentsResponse struct {
	
	
	Error error

	
	Topics map[string]ListPartitionReassignmentsResponseTopic
}



type ListPartitionReassignmentsResponseTopic struct {
	
	Partitions []ListPartitionReassignmentsResponsePartition
}



type ListPartitionReassignmentsResponsePartition struct {
	
	PartitionIndex int

	
	Replicas []int

	
	AddingReplicas []int

	
	RemovingReplicas []int
}

func (c *Client) ListPartitionReassignments(
	ctx context.Context,
	req *ListPartitionReassignmentsRequest,
) (*ListPartitionReassignmentsResponse, error) {
	apiReq := &listpartitionreassignments.Request{
		TimeoutMs: int32(req.Timeout.Milliseconds()),
	}

	for topicName, topicReq := range req.Topics {
		apiReq.Topics = append(
			apiReq.Topics,
			listpartitionreassignments.RequestTopic{
				Name:             topicName,
				PartitionIndexes: intToInt32Array(topicReq.PartitionIndexes),
			},
		)
	}

	protoResp, err := c.roundTrip(
		ctx,
		req.Addr,
		apiReq,
	)
	if err != nil {
		return nil, err
	}
	apiResp := protoResp.(*listpartitionreassignments.Response)

	resp := &ListPartitionReassignmentsResponse{
		Error:  makeError(apiResp.ErrorCode, apiResp.ErrorMessage),
		Topics: make(map[string]ListPartitionReassignmentsResponseTopic),
	}

	for _, topicResult := range apiResp.Topics {
		respTopic := ListPartitionReassignmentsResponseTopic{}
		for _, partitionResult := range topicResult.Partitions {
			respTopic.Partitions = append(
				respTopic.Partitions,
				ListPartitionReassignmentsResponsePartition{
					PartitionIndex:   int(partitionResult.PartitionIndex),
					Replicas:         int32ToIntArray(partitionResult.Replicas),
					AddingReplicas:   int32ToIntArray(partitionResult.AddingReplicas),
					RemovingReplicas: int32ToIntArray(partitionResult.RemovingReplicas),
				},
			)
		}
		resp.Topics[topicResult.Name] = respTopic
	}

	return resp, nil
}

func intToInt32Array(arr []int) []int32 {
	if arr == nil {
		return nil
	}
	res := make([]int32, len(arr))
	for i := range arr {
		res[i] = int32(arr[i])
	}
	return res
}

func int32ToIntArray(arr []int32) []int {
	if arr == nil {
		return nil
	}
	res := make([]int, len(arr))
	for i := range arr {
		res[i] = int(arr[i])
	}
	return res
}
