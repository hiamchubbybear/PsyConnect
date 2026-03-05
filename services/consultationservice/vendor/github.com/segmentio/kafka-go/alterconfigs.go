package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/alterconfigs"
)


type AlterConfigsRequest struct {
	
	Addr net.Addr

	
	Resources []AlterConfigRequestResource

	
	
	ValidateOnly bool
}

type AlterConfigRequestResource struct {
	
	ResourceType ResourceType

	
	ResourceName string

	
	Configs []AlterConfigRequestConfig
}

type AlterConfigRequestConfig struct {
	
	Name string

	
	Value string
}


type AlterConfigsResponse struct {
	
	Throttle time.Duration

	
	
	
	
	
	Errors map[AlterConfigsResponseResource]error
}



type AlterConfigsResponseResource struct {
	Type int8
	Name string
}



func (c *Client) AlterConfigs(ctx context.Context, req *AlterConfigsRequest) (*AlterConfigsResponse, error) {
	resources := make([]alterconfigs.RequestResources, len(req.Resources))

	for i, t := range req.Resources {
		configs := make([]alterconfigs.RequestConfig, len(t.Configs))
		for j, v := range t.Configs {
			configs[j] = alterconfigs.RequestConfig{
				Name:  v.Name,
				Value: v.Value,
			}
		}
		resources[i] = alterconfigs.RequestResources{
			ResourceType: int8(t.ResourceType),
			ResourceName: t.ResourceName,
			Configs:      configs,
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &alterconfigs.Request{
		Resources:    resources,
		ValidateOnly: req.ValidateOnly,
	})

	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).AlterConfigs: %w", err)
	}

	res := m.(*alterconfigs.Response)
	ret := &AlterConfigsResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
		Errors:   make(map[AlterConfigsResponseResource]error, len(res.Responses)),
	}

	for _, t := range res.Responses {
		ret.Errors[AlterConfigsResponseResource{
			Type: t.ResourceType,
			Name: t.ResourceName,
		}] = makeError(t.ErrorCode, t.ErrorMessage)
	}

	return ret, nil
}
