












package compact 

import (
	"sort"
	"strings"

	"golang.org/x/text/internal/language"
)


type ID uint16

func getCoreIndex(t language.Tag) (id ID, ok bool) {
	cci, ok := language.GetCompactCore(t)
	if !ok {
		return 0, false
	}
	i := sort.Search(len(coreTags), func(i int) bool {
		return cci <= coreTags[i]
	})
	if i == len(coreTags) || coreTags[i] != cci {
		return 0, false
	}
	return ID(i), true
}


func (id ID) Parent() ID {
	return parents[id]
}


func (id ID) Tag() language.Tag {
	if int(id) >= len(coreTags) {
		return specialTags[int(id)-len(coreTags)]
	}
	return coreTags[id].Tag()
}

var specialTags []language.Tag

func init() {
	tags := strings.Split(specialTagsStr, " ")
	specialTags = make([]language.Tag, len(tags))
	for i, t := range tags {
		specialTags[i] = language.MustParse(t)
	}
}
