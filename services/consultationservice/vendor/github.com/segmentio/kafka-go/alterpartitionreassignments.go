package kafka

import (
	"context"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/alterpartitionreassignments"
)


type AlterPartitionReassignmentsRequest struct {
	
	Addr net.Addr

	
	
	Topic string

	
	Assignments []AlterPartitionReassignmentsRequestAssignment

	
	Timeout time.Duration
}



type AlterPartitionReassignmentsRequestAssignment struct {
	
	Topic string

	
	PartitionID int

	
	BrokerIDs []int
}


type AlterPartitionReassignmentsResponse struct {
	
	
	Error error

	
	PartitionResults []AlterPartitionReassignmentsResponsePartitionResult
}



type AlterPartitionReassignmentsResponsePartitionResult struct {
	
	Topic string

	
	PartitionID int

	
	
	Error error
}

func (c *Client) AlterPartitionReassignments(
	ctx context.Context,
	req *AlterPartitionReassignmentsRequest,
) (*AlterPartitionReassignmentsResponse, error) {
	apiTopicMap := make(map[string]*alterpartitionreassignments.RequestTopic)

	for _, assignment := range req.Assignments {
		topic := assignment.Topic
		if topic == "" {
			topic = req.Topic
		}

		apiTopic := apiTopicMap[topic]
		if apiTopic == nil {
			apiTopic = &alterpartitionreassignments.RequestTopic{
				Name: topic,
			}
			apiTopicMap[topic] = apiTopic
		}

		replicas := []int32{}
		for _, brokerID := range assignment.BrokerIDs {
			replicas = append(replicas, int32(brokerID))
		}

		apiTopic.Partitions = append(
			apiTopic.Partitions,
			alterpartitionreassignments.RequestPartition{
				PartitionIndex: int32(assignment.PartitionID),
				Replicas:       replicas,
			},
		)
	}

	apiReq := &alterpartitionreassignments.Request{
		TimeoutMs: int32(req.Timeout.Milliseconds()),
	}

	for _, apiTopic := range apiTopicMap {
		apiReq.Topics = append(apiReq.Topics, *apiTopic)
	}

	protoResp, err := c.roundTrip(
		ctx,
		req.Addr,
		apiReq,
	)
	if err != nil {
		return nil, err
	}
	apiResp := protoResp.(*alterpartitionreassignments.Response)

	resp := &AlterPartitionReassignmentsResponse{
		Error: makeError(apiResp.ErrorCode, apiResp.ErrorMessage),
	}

	for _, topicResult := range apiResp.Results {
		for _, partitionResult := range topicResult.Partitions {
			resp.PartitionResults = append(
				resp.PartitionResults,
				AlterPartitionReassignmentsResponsePartitionResult{
					Topic:       topicResult.Name,
					PartitionID: int(partitionResult.PartitionIndex),
					Error:       makeError(partitionResult.ErrorCode, partitionResult.ErrorMessage),
				},
			)
		}
	}

	return resp, nil
}
