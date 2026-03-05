package kafka

import (
	"context"
	"net"

	"github.com/segmentio/kafka-go/protocol/incrementalalterconfigs"
)

type ConfigOperation int8

const (
	ConfigOperationSet      ConfigOperation = 0
	ConfigOperationDelete   ConfigOperation = 1
	ConfigOperationAppend   ConfigOperation = 2
	ConfigOperationSubtract ConfigOperation = 3
)


type IncrementalAlterConfigsRequest struct {
	
	Addr net.Addr

	
	Resources []IncrementalAlterConfigsRequestResource

	
	
	ValidateOnly bool
}



type IncrementalAlterConfigsRequestResource struct {
	
	ResourceType ResourceType

	
	ResourceName string

	
	Configs []IncrementalAlterConfigsRequestConfig
}



type IncrementalAlterConfigsRequestConfig struct {
	
	Name string

	
	Value string

	
	ConfigOperation ConfigOperation
}


type IncrementalAlterConfigsResponse struct {
	
	Resources []IncrementalAlterConfigsResponseResource
}



type IncrementalAlterConfigsResponseResource struct {
	
	
	Error error

	
	ResourceType ResourceType

	
	ResourceName string
}

func (c *Client) IncrementalAlterConfigs(
	ctx context.Context,
	req *IncrementalAlterConfigsRequest,
) (*IncrementalAlterConfigsResponse, error) {
	apiReq := &incrementalalterconfigs.Request{
		ValidateOnly: req.ValidateOnly,
	}

	for _, res := range req.Resources {
		apiRes := incrementalalterconfigs.RequestResource{
			ResourceType: int8(res.ResourceType),
			ResourceName: res.ResourceName,
		}

		for _, config := range res.Configs {
			apiRes.Configs = append(
				apiRes.Configs,
				incrementalalterconfigs.RequestConfig{
					Name:            config.Name,
					Value:           config.Value,
					ConfigOperation: int8(config.ConfigOperation),
				},
			)
		}

		apiReq.Resources = append(
			apiReq.Resources,
			apiRes,
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

	resp := &IncrementalAlterConfigsResponse{}

	apiResp := protoResp.(*incrementalalterconfigs.Response)
	for _, res := range apiResp.Responses {
		resp.Resources = append(
			resp.Resources,
			IncrementalAlterConfigsResponseResource{
				Error:        makeError(res.ErrorCode, res.ErrorMessage),
				ResourceType: ResourceType(res.ResourceType),
				ResourceName: res.ResourceName,
			},
		)
	}

	return resp, nil
}
