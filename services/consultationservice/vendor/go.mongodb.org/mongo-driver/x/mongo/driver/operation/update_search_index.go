





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


type UpdateSearchIndex struct {
	authenticator driver.Authenticator
	index         string
	definition    bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	result        UpdateSearchIndexResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}


type UpdateSearchIndexResult struct {
	Ok int32
}

func buildUpdateSearchIndexResult(response bsoncore.Document) (UpdateSearchIndexResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return UpdateSearchIndexResult{}, err
	}
	usir := UpdateSearchIndexResult{}
	for _, element := range elements {
		if element.Key() == "ok" {
			var ok bool
			usir.Ok, ok = element.Value().AsInt32OK()
			if !ok {
				return usir, fmt.Errorf("response field 'ok' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return usir, nil
}


func NewUpdateSearchIndex(index string, definition bsoncore.Document) *UpdateSearchIndex {
	return &UpdateSearchIndex{
		index:      index,
		definition: definition,
	}
}


func (usi *UpdateSearchIndex) Result() UpdateSearchIndexResult { return usi.result }

func (usi *UpdateSearchIndex) processResponse(info driver.ResponseInfo) error {
	var err error
	usi.result, err = buildUpdateSearchIndexResult(info.ServerResponse)
	return err
}


func (usi *UpdateSearchIndex) Execute(ctx context.Context) error {
	if usi.deployment == nil {
		return errors.New("the UpdateSearchIndex operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         usi.command,
		ProcessResponseFn: usi.processResponse,
		Client:            usi.session,
		Clock:             usi.clock,
		CommandMonitor:    usi.monitor,
		Crypt:             usi.crypt,
		Database:          usi.database,
		Deployment:        usi.deployment,
		Selector:          usi.selector,
		ServerAPI:         usi.serverAPI,
		Timeout:           usi.timeout,
		Authenticator:     usi.authenticator,
	}.Execute(ctx)

}

func (usi *UpdateSearchIndex) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "updateSearchIndex", usi.collection)
	dst = bsoncore.AppendStringElement(dst, "name", usi.index)
	dst = bsoncore.AppendDocumentElement(dst, "definition", usi.definition)
	return dst, nil
}


func (usi *UpdateSearchIndex) Index(name string) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.index = name
	return usi
}


func (usi *UpdateSearchIndex) Definition(definition bsoncore.Document) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.definition = definition
	return usi
}


func (usi *UpdateSearchIndex) Session(session *session.Client) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.session = session
	return usi
}


func (usi *UpdateSearchIndex) ClusterClock(clock *session.ClusterClock) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.clock = clock
	return usi
}


func (usi *UpdateSearchIndex) Collection(collection string) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.collection = collection
	return usi
}


func (usi *UpdateSearchIndex) CommandMonitor(monitor *event.CommandMonitor) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.monitor = monitor
	return usi
}


func (usi *UpdateSearchIndex) Crypt(crypt driver.Crypt) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.crypt = crypt
	return usi
}


func (usi *UpdateSearchIndex) Database(database string) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.database = database
	return usi
}


func (usi *UpdateSearchIndex) Deployment(deployment driver.Deployment) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.deployment = deployment
	return usi
}


func (usi *UpdateSearchIndex) ServerSelector(selector description.ServerSelector) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.selector = selector
	return usi
}


func (usi *UpdateSearchIndex) ServerAPI(serverAPI *driver.ServerAPIOptions) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.serverAPI = serverAPI
	return usi
}


func (usi *UpdateSearchIndex) Timeout(timeout *time.Duration) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.timeout = timeout
	return usi
}


func (usi *UpdateSearchIndex) Authenticator(authenticator driver.Authenticator) *UpdateSearchIndex {
	if usi == nil {
		usi = new(UpdateSearchIndex)
	}

	usi.authenticator = authenticator
	return usi
}
