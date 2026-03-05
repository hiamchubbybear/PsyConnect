package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/describeclientquotas"
)



type DescribeClientQuotasRequest struct {
	
	Addr net.Addr

	
	Components []DescribeClientQuotasRequestComponent

	
	
	Strict bool
}

type DescribeClientQuotasRequestComponent struct {
	
	EntityType string

	
	
	MatchType int8

	
	Match string
}


type DescribeClientQuotasResponse struct {
	
	Throttle time.Duration

	
	
	Error error

	
	Entries []DescribeClientQuotasResponseQuotas
}

type DescribeClientQuotasEntity struct {
	
	EntityType string

	
	EntityName string
}

type DescribeClientQuotasValue struct {
	
	Key string

	
	Value float64
}

type DescribeClientQuotasResponseQuotas struct {
	
	Entities []DescribeClientQuotasEntity

	
	Values []DescribeClientQuotasValue
}



func (c *Client) DescribeClientQuotas(ctx context.Context, req *DescribeClientQuotasRequest) (*DescribeClientQuotasResponse, error) {
	components := make([]describeclientquotas.Component, len(req.Components))

	for componentIdx, component := range req.Components {
		components[componentIdx] = describeclientquotas.Component{
			EntityType: component.EntityType,
			MatchType:  component.MatchType,
			Match:      component.Match,
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &describeclientquotas.Request{
		Components: components,
		Strict:     req.Strict,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).DescribeClientQuotas: %w", err)
	}

	res := m.(*describeclientquotas.Response)
	responseEntries := make([]DescribeClientQuotasResponseQuotas, len(res.Entries))

	for responseEntryIdx, responseEntry := range res.Entries {
		responseEntities := make([]DescribeClientQuotasEntity, len(responseEntry.Entities))
		for responseEntityIdx, responseEntity := range responseEntry.Entities {
			responseEntities[responseEntityIdx] = DescribeClientQuotasEntity{
				EntityType: responseEntity.EntityType,
				EntityName: responseEntity.EntityName,
			}
		}

		responseValues := make([]DescribeClientQuotasValue, len(responseEntry.Values))
		for responseValueIdx, responseValue := range responseEntry.Values {
			responseValues[responseValueIdx] = DescribeClientQuotasValue{
				Key:   responseValue.Key,
				Value: responseValue.Value,
			}
		}
		responseEntries[responseEntryIdx] = DescribeClientQuotasResponseQuotas{
			Entities: responseEntities,
			Values:   responseValues,
		}
	}
	ret := &DescribeClientQuotasResponse{
		Throttle: time.Duration(res.ThrottleTimeMs),
		Entries:  responseEntries,
	}

	return ret, nil
}
