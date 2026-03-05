









package tag 

import (
	"bytes"
	"fmt"
)


type Tag struct {
	Name  string
	Value string
}


func (tag Tag) String() string {
	return fmt.Sprintf("%s=%s", tag.Name, tag.Value)
}





func NewTagSetFromMap(m map[string]string) Set {
	var set Set
	for k, v := range m {
		set = append(set, Tag{Name: k, Value: v})
	}

	return set
}





func NewTagSetsFromMaps(maps []map[string]string) []Set {
	sets := make([]Set, 0, len(maps))
	for _, m := range maps {
		sets = append(sets, NewTagSetFromMap(m))
	}
	return sets
}


type Set []Tag


func (ts Set) Contains(name, value string) bool {
	for _, t := range ts {
		if t.Name == name && t.Value == value {
			return true
		}
	}

	return false
}


func (ts Set) ContainsAll(other []Tag) bool {
	for _, ot := range other {
		if !ts.Contains(ot.Name, ot.Value) {
			return false
		}
	}

	return true
}


func (ts Set) String() string {
	var b bytes.Buffer
	for i, tag := range ts {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(tag.String())
	}
	return b.String()
}
