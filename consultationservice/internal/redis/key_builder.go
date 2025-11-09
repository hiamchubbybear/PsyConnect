package redis

import (
	"fmt"
	"strings"
)

type KeyBuilder struct {
	Service string
}

func NewKeyBuilder(service string) *KeyBuilder {
	return &KeyBuilder{Service: service}
}

func (kb *KeyBuilder) Build(entity string, id string, field ...string) string {

	if len(field) > 0 {
		return fmt.Sprintf("%s:%s:%s:%s",
			kb.Service,
			entity,
			strings.Join(field, ":"),
			id,
		)
	}
	return fmt.Sprintf("%s:%s:%s", kb.Service, entity, id)
}
