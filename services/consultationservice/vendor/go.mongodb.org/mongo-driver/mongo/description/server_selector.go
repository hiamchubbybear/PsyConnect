





package description

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/tag"
)




type ServerSelector interface {
	SelectServer(Topology, []Server) ([]Server, error)
}


type ServerSelectorFunc func(Topology, []Server) ([]Server, error)


func (ssf ServerSelectorFunc) SelectServer(t Topology, s []Server) ([]Server, error) {
	return ssf(t, s)
}



type serverSelectorInfo struct {
	Type      string
	Data      string               `json:",omitempty"`
	Selectors []serverSelectorInfo `json:",omitempty"`
}


func (sss serverSelectorInfo) String() string {
	bytes, _ := json.Marshal(sss)

	return string(bytes)
}



type serverSelectorInfoGetter interface {
	info() serverSelectorInfo
}

type compositeSelector struct {
	selectors []ServerSelector
}

func (cs *compositeSelector) info() serverSelectorInfo {
	csInfo := serverSelectorInfo{Type: "compositeSelector"}

	for _, sel := range cs.selectors {
		if getter, ok := sel.(serverSelectorInfoGetter); ok {
			csInfo.Selectors = append(csInfo.Selectors, getter.info())
		}
	}

	return csInfo
}


func (cs *compositeSelector) String() string {
	return cs.info().String()
}











func CompositeSelector(selectors []ServerSelector) ServerSelector {
	return &compositeSelector{selectors: selectors}
}

func (cs *compositeSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	var err error
	for _, sel := range cs.selectors {
		candidates, err = sel.SelectServer(t, candidates)
		if err != nil {
			return nil, err
		}
	}
	return candidates, nil
}

type latencySelector struct {
	latency time.Duration
}


func LatencySelector(latency time.Duration) ServerSelector {
	return &latencySelector{latency: latency}
}

func (latencySelector) info() serverSelectorInfo {
	return serverSelectorInfo{Type: "latencySelector"}
}

func (selector latencySelector) String() string {
	return selector.info().String()
}

func (selector *latencySelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	if selector.latency < 0 {
		return candidates, nil
	}
	if t.Kind == LoadBalanced {
		
		return candidates, nil
	}

	switch len(candidates) {
	case 0, 1:
		return candidates, nil
	default:
		min := time.Duration(math.MaxInt64)
		for _, candidate := range candidates {
			if candidate.AverageRTTSet {
				if candidate.AverageRTT < min {
					min = candidate.AverageRTT
				}
			}
		}

		if min == math.MaxInt64 {
			return candidates, nil
		}

		max := min + selector.latency

		viableIndexes := make([]int, 0, len(candidates))
		for i, candidate := range candidates {
			if candidate.AverageRTTSet {
				if candidate.AverageRTT <= max {
					viableIndexes = append(viableIndexes, i)
				}
			}
		}
		if len(viableIndexes) == len(candidates) {
			return candidates, nil
		}
		result := make([]Server, len(viableIndexes))
		for i, idx := range viableIndexes {
			result[i] = candidates[idx]
		}
		return result, nil
	}
}

type writeServerSelector struct{}


func WriteSelector() ServerSelector {
	return writeServerSelector{}
}

func (writeServerSelector) info() serverSelectorInfo {
	return serverSelectorInfo{Type: "writeSelector"}
}

func (selector writeServerSelector) String() string {
	return selector.info().String()
}

func (writeServerSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	switch t.Kind {
	case Single, LoadBalanced:
		return candidates, nil
	default:
		
		selected := 0
		for _, candidate := range candidates {
			switch candidate.Kind {
			case Mongos, RSPrimary, Standalone:
				selected++
			}
		}

		
		result := make([]Server, 0, selected)
		for _, candidate := range candidates {
			switch candidate.Kind {
			case Mongos, RSPrimary, Standalone:
				result = append(result, candidate)
			}
		}
		return result, nil
	}
}

type readPrefServerSelector struct {
	rp                *readpref.ReadPref
	isOutputAggregate bool
}


func ReadPrefSelector(rp *readpref.ReadPref) ServerSelector {
	return readPrefServerSelector{
		rp:                rp,
		isOutputAggregate: false,
	}
}

func (selector readPrefServerSelector) info() serverSelectorInfo {
	return serverSelectorInfo{
		Type: "readPrefSelector",
		Data: selector.rp.String(),
	}
}

func (selector readPrefServerSelector) String() string {
	return selector.info().String()
}

func (selector readPrefServerSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	if t.Kind == LoadBalanced {
		
		
		
		return candidates, nil
	}

	switch t.Kind {
	case Single:
		return candidates, nil
	case ReplicaSetNoPrimary, ReplicaSetWithPrimary:
		return selectForReplicaSet(selector.rp, selector.isOutputAggregate, t, candidates)
	case Sharded:
		return selectByKind(candidates, Mongos), nil
	}

	return nil, nil
}



func OutputAggregateSelector(rp *readpref.ReadPref) ServerSelector {
	return readPrefServerSelector{
		rp:                rp,
		isOutputAggregate: true,
	}
}

func selectForReplicaSet(rp *readpref.ReadPref, isOutputAggregate bool, t Topology, candidates []Server) ([]Server, error) {
	if err := verifyMaxStaleness(rp, t); err != nil {
		return nil, err
	}

	
	
	if isOutputAggregate {
		for _, s := range candidates {
			if s.WireVersion.Max < 13 {
				return selectByKind(candidates, RSPrimary), nil
			}
		}
	}

	switch rp.Mode() {
	case readpref.PrimaryMode:
		return selectByKind(candidates, RSPrimary), nil
	case readpref.PrimaryPreferredMode:
		selected := selectByKind(candidates, RSPrimary)

		if len(selected) == 0 {
			selected = selectSecondaries(rp, candidates)
			return selectByTagSet(selected, rp.TagSets()), nil
		}

		return selected, nil
	case readpref.SecondaryPreferredMode:
		selected := selectSecondaries(rp, candidates)
		selected = selectByTagSet(selected, rp.TagSets())
		if len(selected) > 0 {
			return selected, nil
		}
		return selectByKind(candidates, RSPrimary), nil
	case readpref.SecondaryMode:
		selected := selectSecondaries(rp, candidates)
		return selectByTagSet(selected, rp.TagSets()), nil
	case readpref.NearestMode:
		selected := selectByKind(candidates, RSPrimary)
		selected = append(selected, selectSecondaries(rp, candidates)...)
		return selectByTagSet(selected, rp.TagSets()), nil
	}

	return nil, fmt.Errorf("unsupported mode: %d", rp.Mode())
}

func selectSecondaries(rp *readpref.ReadPref, candidates []Server) []Server {
	secondaries := selectByKind(candidates, RSSecondary)
	if len(secondaries) == 0 {
		return secondaries
	}
	if maxStaleness, set := rp.MaxStaleness(); set {
		primaries := selectByKind(candidates, RSPrimary)
		if len(primaries) == 0 {
			baseTime := secondaries[0].LastWriteTime
			for i := 1; i < len(secondaries); i++ {
				if secondaries[i].LastWriteTime.After(baseTime) {
					baseTime = secondaries[i].LastWriteTime
				}
			}

			var selected []Server
			for _, secondary := range secondaries {
				estimatedStaleness := baseTime.Sub(secondary.LastWriteTime) + secondary.HeartbeatInterval
				if estimatedStaleness <= maxStaleness {
					selected = append(selected, secondary)
				}
			}

			return selected
		}

		primary := primaries[0]

		var selected []Server
		for _, secondary := range secondaries {
			estimatedStaleness := secondary.LastUpdateTime.Sub(secondary.LastWriteTime) - primary.LastUpdateTime.Sub(primary.LastWriteTime) + secondary.HeartbeatInterval
			if estimatedStaleness <= maxStaleness {
				selected = append(selected, secondary)
			}
		}
		return selected
	}

	return secondaries
}

func selectByTagSet(candidates []Server, tagSets []tag.Set) []Server {
	if len(tagSets) == 0 {
		return candidates
	}

	for _, ts := range tagSets {
		
		
		if len(ts) == 0 {
			return candidates
		}

		var results []Server
		for _, s := range candidates {
			
			if len(s.Tags) > 0 && s.Tags.ContainsAll(ts) {
				results = append(results, s)
			}
		}

		if len(results) > 0 {
			return results
		}
	}

	return []Server{}
}

func selectByKind(candidates []Server, kind ServerKind) []Server {
	
	
	viableIndexes := make([]int, 0, len(candidates))
	for i, s := range candidates {
		if s.Kind == kind {
			viableIndexes = append(viableIndexes, i)
		}
	}
	if len(viableIndexes) == len(candidates) {
		return candidates
	}
	result := make([]Server, len(viableIndexes))
	for i, idx := range viableIndexes {
		result[i] = candidates[idx]
	}
	return result
}

func verifyMaxStaleness(rp *readpref.ReadPref, t Topology) error {
	maxStaleness, set := rp.MaxStaleness()
	if !set {
		return nil
	}

	if maxStaleness < 90*time.Second {
		return fmt.Errorf("max staleness (%s) must be greater than or equal to 90s", maxStaleness)
	}

	if len(t.Servers) < 1 {
		
		return nil
	}

	
	s := t.Servers[0]
	idleWritePeriod := 10 * time.Second

	if maxStaleness < s.HeartbeatInterval+idleWritePeriod {
		return fmt.Errorf(
			"max staleness (%s) must be greater than or equal to the heartbeat interval (%s) plus idle write period (%s)",
			maxStaleness, s.HeartbeatInterval, idleWritePeriod,
		)
	}

	return nil
}
