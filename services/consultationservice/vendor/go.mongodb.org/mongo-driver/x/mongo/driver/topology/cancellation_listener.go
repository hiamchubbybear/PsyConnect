





package topology

import "context"

type cancellationListener interface {
	Listen(context.Context, func())
	StopListening() bool
}
