


package sdk

import (
	"log/slog"
	"os"
	"strconv"
)


var maxSpan = newSpanLimits()

type spanLimits struct {
	
	
	
	
	
	
	Attrs int
	
	
	
	
	
	
	AttrValueLen int
	
	
	
	
	Events int
	
	
	
	
	EventAttrs int
	
	
	
	
	Links int
	
	
	
	
	LinkAttrs int
}

func newSpanLimits() spanLimits {
	return spanLimits{
		Attrs: firstEnv(
			128,
			"OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT",
			"OTEL_ATTRIBUTE_COUNT_LIMIT",
		),
		AttrValueLen: firstEnv(
			-1, 
			"OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT",
			"OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT",
		),
		Events:     firstEnv(128, "OTEL_SPAN_EVENT_COUNT_LIMIT"),
		EventAttrs: firstEnv(128, "OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT"),
		Links:      firstEnv(128, "OTEL_SPAN_LINK_COUNT_LIMIT"),
		LinkAttrs:  firstEnv(128, "OTEL_LINK_ATTRIBUTE_COUNT_LIMIT"),
	}
}




func firstEnv(defaultVal int, keys ...string) int {
	for _, key := range keys {
		strV := os.Getenv(key)
		if strV == "" {
			continue
		}

		v, err := strconv.Atoi(strV)
		if err == nil {
			return v
		}
		slog.Warn(
			"invalid limit environment variable",
			"error", err,
			"key", key,
			"value", strV,
		)
	}

	return defaultVal
}
