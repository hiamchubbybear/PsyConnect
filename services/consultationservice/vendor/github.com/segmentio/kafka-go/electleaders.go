package kafka

import (
	"context"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/electleaders"
)


type ElectLeadersRequest struct {
	
	Addr net.Addr

	
	Topic string

	
	Partitions []int

	
	Timeout time.Duration
}


type ElectLeadersResponse struct {
	
	Error error

	
	PartitionResults []ElectLeadersResponsePartitionResult
}


type ElectLeadersResponsePartitionResult struct {
	
	Partition int

	
	
	Error error
}

func (c *Client) ElectLeaders(
	ctx context.Context,
	req *ElectLeadersRequest,
) (*ElectLeadersResponse, error) {
	partitions32 := []int32{}
	for _, partition := range req.Partitions {
		partitions32 = append(partitions32, int32(partition))
	}

	protoResp, err := c.roundTrip(
		ctx,
		req.Addr,
		&electleaders.Request{
			TopicPartitions: []electleaders.RequestTopicPartitions{
				{
					Topic:        req.Topic,
					PartitionIDs: partitions32,
				},
			},
			TimeoutMs: int32(req.Timeout.Milliseconds()),
		},
	)
	if err != nil {
		return nil, err
	}
	apiResp := protoResp.(*electleaders.Response)

	resp := &ElectLeadersResponse{
		Error: makeError(apiResp.ErrorCode, ""),
	}

	for _, topicResult := range apiResp.ReplicaElectionResults {
		for _, partitionResult := range topicResult.PartitionResults {
			resp.PartitionResults = append(
				resp.PartitionResults,
				ElectLeadersResponsePartitionResult{
					Partition: int(partitionResult.PartitionID),
					Error:     makeError(partitionResult.ErrorCode, partitionResult.ErrorMessage),
				},
			)
		}
	}

	return resp, nil
}
