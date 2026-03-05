





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type DropSearchIndex struct {
	authenticator driver.Authenticator
	index         string
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	result        DropSearchIndexResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}


type DropSearchIndexResult struct {
	Ok int32
}

func buildDropSearchIndexResult(response bsoncore.Document) (DropSearchIndexResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DropSearchIndexResult{}, err
	}
	dsir := DropSearchIndexResult{}
	for _, element := range elements {
		if element.Key() == "ok" {
			var ok bool
			dsir.Ok, ok = element.Value().AsInt32OK()
			if !ok {
				return dsir, fmt.Errorf("response field 'ok' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dsir, nil
}


func NewDropSearchIndex(index string) *DropSearchIndex {
	return &DropSearchIndex{
		index: index,
	}
}


func (dsi *DropSearchIndex) Result() DropSearchIndexResult { return dsi.result }

func (dsi *DropSearchIndex) processResponse(info driver.ResponseInfo) error {
	var err error
	dsi.result, err = buildDropSearchIndexResult(info.ServerResponse)
	return err
}


func (dsi *DropSearchIndex) Execute(ctx context.Context) error {
	if dsi.deployment == nil {
		return errors.New("the DropSearchIndex operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         dsi.command,
		ProcessResponseFn: dsi.processResponse,
		Client:            dsi.session,
		Clock:             dsi.clock,
		CommandMonitor:    dsi.monitor,
		Crypt:             dsi.crypt,
		Database:          dsi.database,
		Deployment:        dsi.deployment,
		Selector:          dsi.selector,
		ServerAPI:         dsi.serverAPI,
		Timeout:           dsi.timeout,
		Authenticator:     dsi.authenticator,
	}.Execute(ctx)

}

func (dsi *DropSearchIndex) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "dropSearchIndex", dsi.collection)
	dst = bsoncore.AppendStringElement(dst, "name", dsi.index)
	return dst, nil
}


func (dsi *DropSearchIndex) Index(index string) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.index = index
	return dsi
}


func (dsi *DropSearchIndex) Session(session *session.Client) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.session = session
	return dsi
}


func (dsi *DropSearchIndex) ClusterClock(clock *session.ClusterClock) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.clock = clock
	return dsi
}


func (dsi *DropSearchIndex) Collection(collection string) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.collection = collection
	return dsi
}


func (dsi *DropSearchIndex) CommandMonitor(monitor *event.CommandMonitor) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.monitor = monitor
	return dsi
}


func (dsi *DropSearchIndex) Crypt(crypt driver.Crypt) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.crypt = crypt
	return dsi
}


func (dsi *DropSearchIndex) Database(database string) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.database = database
	return dsi
}


func (dsi *DropSearchIndex) Deployment(deployment driver.Deployment) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.deployment = deployment
	return dsi
}


func (dsi *DropSearchIndex) ServerSelector(selector description.ServerSelector) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.selector = selector
	return dsi
}


func (dsi *DropSearchIndex) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.serverAPI = serverAPI
	return dsi
}


func (dsi *DropSearchIndex) Timeout(timeout *time.Duration) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.timeout = timeout
	return dsi
}


func (dsi *DropSearchIndex) Authenticator(authenticator driver.Authenticator) *DropSearchIndex {
	if dsi == nil {
		dsi = new(DropSearchIndex)
	}

	dsi.authenticator = authenticator
	return dsi
}
