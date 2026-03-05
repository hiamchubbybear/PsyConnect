package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/deleteacls"
)



type DeleteACLsRequest struct {
	
	Addr net.Addr

	
	Filters []DeleteACLsFilter
}

type DeleteACLsFilter struct {
	ResourceTypeFilter        ResourceType
	ResourceNameFilter        string
	ResourcePatternTypeFilter PatternType
	PrincipalFilter           string
	HostFilter                string
	Operation                 ACLOperationType
	PermissionType            ACLPermissionType
}



type DeleteACLsResponse struct {
	
	Throttle time.Duration

	
	Results []DeleteACLsResult
}

type DeleteACLsResult struct {
	Error        error
	MatchingACLs []DeleteACLsMatchingACLs
}

type DeleteACLsMatchingACLs struct {
	Error               error
	ResourceType        ResourceType
	ResourceName        string
	ResourcePatternType PatternType
	Principal           string
	Host                string
	Operation           ACLOperationType
	PermissionType      ACLPermissionType
}



func (c *Client) DeleteACLs(ctx context.Context, req *DeleteACLsRequest) (*DeleteACLsResponse, error) {
	filters := make([]deleteacls.RequestFilter, 0, len(req.Filters))

	for _, filter := range req.Filters {
		filters = append(filters, deleteacls.RequestFilter{
			ResourceTypeFilter:        int8(filter.ResourceTypeFilter),
			ResourceNameFilter:        filter.ResourceNameFilter,
			ResourcePatternTypeFilter: int8(filter.ResourcePatternTypeFilter),
			PrincipalFilter:           filter.PrincipalFilter,
			HostFilter:                filter.HostFilter,
			Operation:                 int8(filter.Operation),
			PermissionType:            int8(filter.PermissionType),
		})
	}

	m, err := c.roundTrip(ctx, req.Addr, &deleteacls.Request{
		Filters: filters,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).DeleteACLs: %w", err)
	}

	res := m.(*deleteacls.Response)

	results := make([]DeleteACLsResult, 0, len(res.FilterResults))

	for _, result := range res.FilterResults {
		matchingACLs := make([]DeleteACLsMatchingACLs, 0, len(result.MatchingACLs))

		for _, matchingACL := range result.MatchingACLs {
			matchingACLs = append(matchingACLs, DeleteACLsMatchingACLs{
				Error:               makeError(matchingACL.ErrorCode, matchingACL.ErrorMessage),
				ResourceType:        ResourceType(matchingACL.ResourceType),
				ResourceName:        matchingACL.ResourceName,
				ResourcePatternType: PatternType(matchingACL.ResourcePatternType),
				Principal:           matchingACL.Principal,
				Host:                matchingACL.Host,
				Operation:           ACLOperationType(matchingACL.Operation),
				PermissionType:      ACLPermissionType(matchingACL.PermissionType),
			})
		}

		results = append(results, DeleteACLsResult{
			Error:        makeError(result.ErrorCode, result.ErrorMessage),
			MatchingACLs: matchingACLs,
		})
	}

	ret := &DeleteACLsResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
		Results:  results,
	}

	return ret, nil
}
