


package stats

import "context"


type Labels struct {
	
	TelemetryLabels map[string]string
}

type labelsKey struct{}


func GetLabels(ctx context.Context) *Labels {
	labels, _ := ctx.Value(labelsKey{}).(*Labels)
	return labels
}


func SetLabels(ctx context.Context, labels *Labels) context.Context {
	
	return context.WithValue(ctx, labelsKey{}, labels)
}
