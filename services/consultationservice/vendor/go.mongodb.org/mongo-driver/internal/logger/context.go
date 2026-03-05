





package logger

import "context"



type contextKey string

const (
	contextKeyOperation   contextKey = "operation"
	contextKeyOperationID contextKey = "operationID"
)


func WithOperationName(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, contextKeyOperation, operation)
}


func WithOperationID(ctx context.Context, operationID int32) context.Context {
	return context.WithValue(ctx, contextKeyOperationID, operationID)
}


func OperationName(ctx context.Context) (string, bool) {
	operationName := ctx.Value(contextKeyOperation)
	if operationName == nil {
		return "", false
	}

	return operationName.(string), true
}


func OperationID(ctx context.Context) (int32, bool) {
	operationID := ctx.Value(contextKeyOperationID)
	if operationID == nil {
		return 0, false
	}

	return operationID.(int32), true
}
