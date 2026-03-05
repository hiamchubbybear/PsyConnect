package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/deletegroups"
)



type DeleteGroupsRequest struct {
	
	Addr net.Addr

	
	GroupIDs []string
}



type DeleteGroupsResponse struct {
	
	Throttle time.Duration

	
	
	
	
	Errors map[string]error
}



func (c *Client) DeleteGroups(
	ctx context.Context,
	req *DeleteGroupsRequest,
) (*DeleteGroupsResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &deletegroups.Request{
		GroupIDs: req.GroupIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).DeleteGroups: %w", err)
	}

	r := m.(*deletegroups.Response)

	ret := &DeleteGroupsResponse{
		Throttle: makeDuration(r.ThrottleTimeMs),
		Errors:   make(map[string]error, len(r.Responses)),
	}

	for _, t := range r.Responses {
		ret.Errors[t.GroupID] = makeError(t.ErrorCode, "")
	}

	return ret, nil
}
