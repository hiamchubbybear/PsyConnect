package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/describeconfigs"
)


type DescribeConfigsRequest struct {
	
	Addr net.Addr

	
	Resources []DescribeConfigRequestResource

	
	IncludeSynonyms bool

	
	IncludeDocumentation bool
}

type DescribeConfigRequestResource struct {
	
	ResourceType ResourceType

	
	ResourceName string

	
	ConfigNames []string
}


type DescribeConfigsResponse struct {
	
	Throttle time.Duration

	
	Resources []DescribeConfigResponseResource
}


type DescribeConfigResponseResource struct {
	
	ResourceType int8

	
	ResourceName string

	
	Error error

	
	ConfigEntries []DescribeConfigResponseConfigEntry
}


type DescribeConfigResponseConfigEntry struct {
	ConfigName  string
	ConfigValue string
	ReadOnly    bool

	
	IsDefault bool

	
	ConfigSource int8

	IsSensitive bool

	
	ConfigSynonyms []DescribeConfigResponseConfigSynonym

	
	ConfigType int8

	
	ConfigDocumentation string
}


type DescribeConfigResponseConfigSynonym struct {
	
	ConfigName string

	
	ConfigValue string

	
	ConfigSource int8
}



func (c *Client) DescribeConfigs(ctx context.Context, req *DescribeConfigsRequest) (*DescribeConfigsResponse, error) {
	resources := make([]describeconfigs.RequestResource, len(req.Resources))

	for i, t := range req.Resources {
		resources[i] = describeconfigs.RequestResource{
			ResourceType: int8(t.ResourceType),
			ResourceName: t.ResourceName,
			ConfigNames:  t.ConfigNames,
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &describeconfigs.Request{
		Resources:            resources,
		IncludeSynonyms:      req.IncludeSynonyms,
		IncludeDocumentation: req.IncludeDocumentation,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).DescribeConfigs: %w", err)
	}

	res := m.(*describeconfigs.Response)
	ret := &DescribeConfigsResponse{
		Throttle:  makeDuration(res.ThrottleTimeMs),
		Resources: make([]DescribeConfigResponseResource, len(res.Resources)),
	}

	for i, t := range res.Resources {

		configEntries := make([]DescribeConfigResponseConfigEntry, len(t.ConfigEntries))
		for j, v := range t.ConfigEntries {

			configSynonyms := make([]DescribeConfigResponseConfigSynonym, len(v.ConfigSynonyms))
			for k, cs := range v.ConfigSynonyms {
				configSynonyms[k] = DescribeConfigResponseConfigSynonym{
					ConfigName:   cs.ConfigName,
					ConfigValue:  cs.ConfigValue,
					ConfigSource: cs.ConfigSource,
				}
			}

			configEntries[j] = DescribeConfigResponseConfigEntry{
				ConfigName:          v.ConfigName,
				ConfigValue:         v.ConfigValue,
				ReadOnly:            v.ReadOnly,
				ConfigSource:        v.ConfigSource,
				IsDefault:           v.IsDefault,
				IsSensitive:         v.IsSensitive,
				ConfigSynonyms:      configSynonyms,
				ConfigType:          v.ConfigType,
				ConfigDocumentation: v.ConfigDocumentation,
			}
		}

		ret.Resources[i] = DescribeConfigResponseResource{
			ResourceType:  t.ResourceType,
			ResourceName:  t.ResourceName,
			Error:         makeError(t.ErrorCode, t.ErrorMessage),
			ConfigEntries: configEntries,
		}
	}

	return ret, nil
}
