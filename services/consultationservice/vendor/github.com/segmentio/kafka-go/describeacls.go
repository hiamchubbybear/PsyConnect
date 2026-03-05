package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/describeacls"
)



type DescribeACLsRequest struct {
	
	Addr net.Addr

	
	Filter ACLFilter
}

type ACLFilter struct {
	ResourceTypeFilter ResourceType
	ResourceNameFilter string
	
	ResourcePatternTypeFilter PatternType
	PrincipalFilter           string
	HostFilter                string
	Operation                 ACLOperationType
	PermissionType            ACLPermissionType
}



type DescribeACLsResponse struct {
	
	Throttle time.Duration

	
	
	Error error

	
	Resources []ACLResource
}

type ACLResource struct {
	ResourceType ResourceType
	ResourceName string
	PatternType  PatternType
	ACLs         []ACLDescription
}

type ACLDescription struct {
	Principal      string
	Host           string
	Operation      ACLOperationType
	PermissionType ACLPermissionType
}

func (c *Client) DescribeACLs(ctx context.Context, req *DescribeACLsRequest) (*DescribeACLsResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &describeacls.Request{
		Filter: describeacls.ACLFilter{
			ResourceTypeFilter:        int8(req.Filter.ResourceTypeFilter),
			ResourceNameFilter:        req.Filter.ResourceNameFilter,
			ResourcePatternTypeFilter: int8(req.Filter.ResourcePatternTypeFilter),
			PrincipalFilter:           req.Filter.PrincipalFilter,
			HostFilter:                req.Filter.HostFilter,
			Operation:                 int8(req.Filter.Operation),
			PermissionType:            int8(req.Filter.PermissionType),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).DescribeACLs: %w", err)
	}

	res := m.(*describeacls.Response)
	resources := make([]ACLResource, len(res.Resources))

	for resourceIdx, respResource := range res.Resources {
		descriptions := make([]ACLDescription, len(respResource.ACLs))

		for descriptionIdx, respDescription := range respResource.ACLs {
			descriptions[descriptionIdx] = ACLDescription{
				Principal:      respDescription.Principal,
				Host:           respDescription.Host,
				Operation:      ACLOperationType(respDescription.Operation),
				PermissionType: ACLPermissionType(respDescription.PermissionType),
			}
		}

		resources[resourceIdx] = ACLResource{
			ResourceType: ResourceType(respResource.ResourceType),
			ResourceName: respResource.ResourceName,
			PatternType:  PatternType(respResource.PatternType),
			ACLs:         descriptions,
		}
	}

	ret := &DescribeACLsResponse{
		Throttle:  makeDuration(res.ThrottleTimeMs),
		Error:     makeError(res.ErrorCode, res.ErrorMessage),
		Resources: resources,
	}

	return ret, nil
}
