package kafka

import (
	"sort"
)


type GroupMember struct {
	
	ID string

	
	Topics []string

	
	
	UserData []byte
}


type GroupMemberAssignments map[string]map[string][]int


type GroupBalancer interface {
	
	ProtocolName() string

	
	
	
	
	
	
	UserData() ([]byte, error)

	
	
	AssignGroups(members []GroupMember, partitions []Partition) GroupMemberAssignments
}












type RangeGroupBalancer struct{}

func (r RangeGroupBalancer) ProtocolName() string {
	return "range"
}

func (r RangeGroupBalancer) UserData() ([]byte, error) {
	return nil, nil
}

func (r RangeGroupBalancer) AssignGroups(members []GroupMember, topicPartitions []Partition) GroupMemberAssignments {
	groupAssignments := GroupMemberAssignments{}
	membersByTopic := findMembersByTopic(members)

	for topic, members := range membersByTopic {
		partitions := findPartitions(topic, topicPartitions)
		partitionCount := len(partitions)
		memberCount := len(members)

		for memberIndex, member := range members {
			assignmentsByTopic, ok := groupAssignments[member.ID]
			if !ok {
				assignmentsByTopic = map[string][]int{}
				groupAssignments[member.ID] = assignmentsByTopic
			}

			minIndex := memberIndex * partitionCount / memberCount
			maxIndex := (memberIndex + 1) * partitionCount / memberCount

			for partitionIndex, partition := range partitions {
				if partitionIndex >= minIndex && partitionIndex < maxIndex {
					assignmentsByTopic[topic] = append(assignmentsByTopic[topic], partition)
				}
			}
		}
	}

	return groupAssignments
}












type RoundRobinGroupBalancer struct{}

func (r RoundRobinGroupBalancer) ProtocolName() string {
	return "roundrobin"
}

func (r RoundRobinGroupBalancer) UserData() ([]byte, error) {
	return nil, nil
}

func (r RoundRobinGroupBalancer) AssignGroups(members []GroupMember, topicPartitions []Partition) GroupMemberAssignments {
	groupAssignments := GroupMemberAssignments{}
	membersByTopic := findMembersByTopic(members)
	for topic, members := range membersByTopic {
		partitionIDs := findPartitions(topic, topicPartitions)
		memberCount := len(members)

		for memberIndex, member := range members {
			assignmentsByTopic, ok := groupAssignments[member.ID]
			if !ok {
				assignmentsByTopic = map[string][]int{}
				groupAssignments[member.ID] = assignmentsByTopic
			}

			for partitionIndex, partition := range partitionIDs {
				if (partitionIndex % memberCount) == memberIndex {
					assignmentsByTopic[topic] = append(assignmentsByTopic[topic], partition)
				}
			}
		}
	}

	return groupAssignments
}
















type RackAffinityGroupBalancer struct {
	
	
	
	Rack string
}

func (r RackAffinityGroupBalancer) ProtocolName() string {
	return "rack-affinity"
}

func (r RackAffinityGroupBalancer) AssignGroups(members []GroupMember, partitions []Partition) GroupMemberAssignments {
	membersByTopic := make(map[string][]GroupMember)
	for _, m := range members {
		for _, t := range m.Topics {
			membersByTopic[t] = append(membersByTopic[t], m)
		}
	}

	partitionsByTopic := make(map[string][]Partition)
	for _, p := range partitions {
		partitionsByTopic[p.Topic] = append(partitionsByTopic[p.Topic], p)
	}

	assignments := GroupMemberAssignments{}
	for topic := range membersByTopic {
		topicAssignments := r.assignTopic(membersByTopic[topic], partitionsByTopic[topic])
		for member, parts := range topicAssignments {
			memberAssignments, ok := assignments[member]
			if !ok {
				memberAssignments = make(map[string][]int)
				assignments[member] = memberAssignments
			}
			memberAssignments[topic] = parts
		}
	}
	return assignments
}

func (r RackAffinityGroupBalancer) UserData() ([]byte, error) {
	return []byte(r.Rack), nil
}

func (r *RackAffinityGroupBalancer) assignTopic(members []GroupMember, partitions []Partition) map[string][]int {
	zonedPartitions := make(map[string][]int)
	for _, part := range partitions {
		zone := part.Leader.Rack
		zonedPartitions[zone] = append(zonedPartitions[zone], part.ID)
	}

	zonedConsumers := make(map[string][]string)
	for _, member := range members {
		zone := string(member.UserData)
		zonedConsumers[zone] = append(zonedConsumers[zone], member.ID)
	}

	targetPerMember := len(partitions) / len(members)
	remainder := len(partitions) % len(members)
	assignments := make(map[string][]int)

	
	
	
	for zone, parts := range zonedPartitions {
		consumers := zonedConsumers[zone]
		if len(consumers) == 0 {
			continue
		}

		
		
		partsPerMember := len(parts) / len(consumers)
		if partsPerMember > targetPerMember {
			partsPerMember = targetPerMember
		}

		for _, consumer := range consumers {
			assignments[consumer] = append(assignments[consumer], parts[:partsPerMember]...)
			parts = parts[partsPerMember:]
		}

		
		
		
		leftover := len(parts)
		if partsPerMember == targetPerMember {
			if leftover > remainder {
				leftover = remainder
			}
			if leftover > len(consumers) {
				leftover = len(consumers)
			}
			remainder -= leftover
		}

		
		
		
		for i := 0; i < leftover; i++ {
			assignments[consumers[i]] = append(assignments[consumers[i]], parts[i])
		}
		parts = parts[leftover:]

		if len(parts) == 0 {
			delete(zonedPartitions, zone)
		} else {
			zonedPartitions[zone] = parts
		}
	}

	
	var remaining []int
	for _, partitions := range zonedPartitions {
		remaining = append(remaining, partitions...)
	}

	for _, member := range members {
		assigned := assignments[member.ID]
		delta := targetPerMember - len(assigned)
		
		
		
		if delta >= 0 && remainder > 0 {
			delta++
			remainder--
		}
		if delta > 0 {
			assignments[member.ID] = append(assigned, remaining[:delta]...)
			remaining = remaining[delta:]
		}
	}

	return assignments
}



func findPartitions(topic string, partitions []Partition) []int {
	var ids []int
	for _, partition := range partitions {
		if partition.Topic == topic {
			ids = append(ids, partition.ID)
		}
	}
	return ids
}


func findMembersByTopic(members []GroupMember) map[string][]GroupMember {
	membersByTopic := map[string][]GroupMember{}
	for _, member := range members {
		for _, topic := range member.Topics {
			membersByTopic[topic] = append(membersByTopic[topic], member)
		}
	}

	
	
	
	
	
	
	
	
	
	
	
	
	for _, members := range membersByTopic {
		sort.Slice(members, func(i, j int) bool {
			return members[i].ID < members[j].ID
		})
	}

	return membersByTopic
}



func findGroupBalancer(protocolName string, balancers []GroupBalancer) (GroupBalancer, bool) {
	for _, balancer := range balancers {
		if balancer.ProtocolName() == protocolName {
			return balancer, true
		}
	}
	return nil, false
}
