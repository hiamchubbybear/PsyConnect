package kafka

import (
	"fmt"
	"strings"
)


type ResourceType int8

const (
	ResourceTypeUnknown ResourceType = 0
	ResourceTypeAny     ResourceType = 1
	ResourceTypeTopic   ResourceType = 2
	ResourceTypeGroup   ResourceType = 3
	
	ResourceTypeBroker          ResourceType = 4
	ResourceTypeCluster         ResourceType = 4
	ResourceTypeTransactionalID ResourceType = 5
	ResourceTypeDelegationToken ResourceType = 6
)

func (rt ResourceType) String() string {
	mapping := map[ResourceType]string{
		ResourceTypeUnknown: "Unknown",
		ResourceTypeAny:     "Any",
		ResourceTypeTopic:   "Topic",
		ResourceTypeGroup:   "Group",
		
		
		ResourceTypeCluster:         "Cluster",
		ResourceTypeTransactionalID: "Transactionalid",
		ResourceTypeDelegationToken: "Delegationtoken",
	}
	s, ok := mapping[rt]
	if !ok {
		s = mapping[ResourceTypeUnknown]
	}
	return s
}

func (rt ResourceType) MarshalText() ([]byte, error) {
	return []byte(rt.String()), nil
}

func (rt *ResourceType) UnmarshalText(text []byte) error {
	normalized := strings.ToLower(string(text))
	mapping := map[string]ResourceType{
		"unknown":         ResourceTypeUnknown,
		"any":             ResourceTypeAny,
		"topic":           ResourceTypeTopic,
		"group":           ResourceTypeGroup,
		"broker":          ResourceTypeBroker,
		"cluster":         ResourceTypeCluster,
		"transactionalid": ResourceTypeTransactionalID,
		"delegationtoken": ResourceTypeDelegationToken,
	}
	parsed, ok := mapping[normalized]
	if !ok {
		*rt = ResourceTypeUnknown
		return fmt.Errorf("cannot parse %s as a ResourceType", normalized)
	}
	*rt = parsed
	return nil
}


type PatternType int8

const (
	
	
	PatternTypeUnknown PatternType = 0
	
	PatternTypeAny PatternType = 1
	
	PatternTypeMatch PatternType = 2
	
	
	
	PatternTypeLiteral PatternType = 3
	
	
	
	PatternTypePrefixed PatternType = 4
)

func (pt PatternType) String() string {
	mapping := map[PatternType]string{
		PatternTypeUnknown:  "Unknown",
		PatternTypeAny:      "Any",
		PatternTypeMatch:    "Match",
		PatternTypeLiteral:  "Literal",
		PatternTypePrefixed: "Prefixed",
	}
	s, ok := mapping[pt]
	if !ok {
		s = mapping[PatternTypeUnknown]
	}
	return s
}

func (pt PatternType) MarshalText() ([]byte, error) {
	return []byte(pt.String()), nil
}

func (pt *PatternType) UnmarshalText(text []byte) error {
	normalized := strings.ToLower(string(text))
	mapping := map[string]PatternType{
		"unknown":  PatternTypeUnknown,
		"any":      PatternTypeAny,
		"match":    PatternTypeMatch,
		"literal":  PatternTypeLiteral,
		"prefixed": PatternTypePrefixed,
	}
	parsed, ok := mapping[normalized]
	if !ok {
		*pt = PatternTypeUnknown
		return fmt.Errorf("cannot parse %s as a PatternType", normalized)
	}
	*pt = parsed
	return nil
}
