





package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

type changeStreamDeployment struct {
	topologyKind description.TopologyKind
	server       driver.Server
	conn         driver.Connection
}

var _ driver.Deployment = (*changeStreamDeployment)(nil)
var _ driver.Server = (*changeStreamDeployment)(nil)
var _ driver.ErrorProcessor = (*changeStreamDeployment)(nil)

func (c *changeStreamDeployment) SelectServer(context.Context, description.ServerSelector) (driver.Server, error) {
	return c, nil
}

func (c *changeStreamDeployment) Kind() description.TopologyKind {
	return c.topologyKind
}

func (c *changeStreamDeployment) Connection(context.Context) (driver.Connection, error) {
	return c.conn, nil
}

func (c *changeStreamDeployment) RTTMonitor() driver.RTTMonitor {
	return c.server.RTTMonitor()
}

func (c *changeStreamDeployment) ProcessError(err error, conn driver.Connection) driver.ProcessErrorResult {
	ep, ok := c.server.(driver.ErrorProcessor)
	if !ok {
		return driver.NoChange
	}

	return ep.ProcessError(err, conn)
}
