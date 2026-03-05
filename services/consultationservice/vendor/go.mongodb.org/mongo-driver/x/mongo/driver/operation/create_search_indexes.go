





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type CreateSearchIndexes struct {
	authenticator driver.Authenticator
	indexes       bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	result        CreateSearchIndexesResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}


type CreateSearchIndexResult struct {
	Name string
}


type CreateSearchIndexesResult struct {
	IndexesCreated []CreateSearchIndexResult
}

func buildCreateSearchIndexesResult(response bsoncore.Document) (CreateSearchIndexesResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return CreateSearchIndexesResult{}, err
	}
	csir := CreateSearchIndexesResult{}
	for _, element := range elements {
		switch element.Key() {
		case "indexesCreated":
			arr, ok := element.Value().ArrayOK()
			if !ok {
				return csir, fmt.Errorf("response field 'indexesCreated' is type array, but received BSON type %s", element.Value().Type)
			}

			var values []bsoncore.Value
			values, err = arr.Values()
			if err != nil {
				break
			}

			for _, val := range values {
				valDoc, ok := val.DocumentOK()
				if !ok {
					return csir, fmt.Errorf("indexesCreated value is type document, but received BSON type %s", val.Type)
				}
				var indexesCreated CreateSearchIndexResult
				if err = bson.Unmarshal(valDoc, &indexesCreated); err != nil {
					return csir, err
				}
				csir.IndexesCreated = append(csir.IndexesCreated, indexesCreated)
			}
		}
	}
	return csir, nil
}


func NewCreateSearchIndexes(indexes bsoncore.Document) *CreateSearchIndexes {
	return &CreateSearchIndexes{
		indexes: indexes,
	}
}


func (csi *CreateSearchIndexes) Result() CreateSearchIndexesResult { return csi.result }

func (csi *CreateSearchIndexes) processResponse(info driver.ResponseInfo) error {
	var err error
	csi.result, err = buildCreateSearchIndexesResult(info.ServerResponse)
	return err
}


func (csi *CreateSearchIndexes) Execute(ctx context.Context) error {
	if csi.deployment == nil {
		return errors.New("the CreateSearchIndexes operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         csi.command,
		ProcessResponseFn: csi.processResponse,
		Client:            csi.session,
		Clock:             csi.clock,
		CommandMonitor:    csi.monitor,
		Crypt:             csi.crypt,
		Database:          csi.database,
		Deployment:        csi.deployment,
		Selector:          csi.selector,
		ServerAPI:         csi.serverAPI,
		Timeout:           csi.timeout,
		Authenticator:     csi.authenticator,
	}.Execute(ctx)

}

func (csi *CreateSearchIndexes) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "createSearchIndexes", csi.collection)
	if csi.indexes != nil {
		dst = bsoncore.AppendArrayElement(dst, "indexes", csi.indexes)
	}
	return dst, nil
}


func (csi *CreateSearchIndexes) Indexes(indexes bsoncore.Document) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.indexes = indexes
	return csi
}


func (csi *CreateSearchIndexes) Session(session *session.Client) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.session = session
	return csi
}


func (csi *CreateSearchIndexes) ClusterClock(clock *session.ClusterClock) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.clock = clock
	return csi
}


func (csi *CreateSearchIndexes) Collection(collection string) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.collection = collection
	return csi
}


func (csi *CreateSearchIndexes) CommandMonitor(monitor *event.CommandMonitor) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.monitor = monitor
	return csi
}


func (csi *CreateSearchIndexes) Crypt(crypt driver.Crypt) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.crypt = crypt
	return csi
}


func (csi *CreateSearchIndexes) Database(database string) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.database = database
	return csi
}


func (csi *CreateSearchIndexes) Deployment(deployment driver.Deployment) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.deployment = deployment
	return csi
}


func (csi *CreateSearchIndexes) ServerSelector(selector description.ServerSelector) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.selector = selector
	return csi
}


func (csi *CreateSearchIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.serverAPI = serverAPI
	return csi
}


func (csi *CreateSearchIndexes) Timeout(timeout *time.Duration) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.timeout = timeout
	return csi
}


func (csi *CreateSearchIndexes) Authenticator(authenticator driver.Authenticator) *CreateSearchIndexes {
	if csi == nil {
		csi = new(CreateSearchIndexes)
	}

	csi.authenticator = authenticator
	return csi
}
