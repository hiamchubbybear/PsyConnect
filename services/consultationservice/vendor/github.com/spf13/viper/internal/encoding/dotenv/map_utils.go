package dotenv

import (
	"strings"

	"github.com/spf13/cast"
)




func flattenAndMergeMap(shadow, m map[string]any, prefix, delimiter string) map[string]any {
	if shadow != nil && prefix != "" && shadow[prefix] != nil {
		
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]any)
	}

	var m2 map[string]any
	if prefix != "" {
		prefix += delimiter
	}
	for k, val := range m {
		fullKey := prefix + k
		switch val := val.(type) {
		case map[string]any:
			m2 = val
		case map[any]any:
			m2 = cast.ToStringMap(val)
		default:
			
			shadow[strings.ToLower(fullKey)] = val
			continue
		}
		
		shadow = flattenAndMergeMap(shadow, m2, fullKey, delimiter)
	}
	return shadow
}
