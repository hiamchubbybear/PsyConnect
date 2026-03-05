





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type DropIndexes struct {
	authenticator driver.Authenticator
	index         any
	maxTime       *time.Duration
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	result        DropIndexesResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}


type DropIndexesResult struct {
	
	NIndexesWas int32
}

func buildDropIndexesResult(response bsoncore.Document) (DropIndexesResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DropIndexesResult{}, err
	}
	dir := DropIndexesResult{}
	for _, element := range elements {
		if element.Key() == "nIndexesWas" {
			var ok bool
			dir.NIndexesWas, ok = element.Value().AsInt32OK()
			if !ok {
				return dir, fmt.Errorf("response field 'nIndexesWas' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dir, nil
}


func NewDropIndexes(index any) *DropIndexes {
	return &DropIndexes{
		index: index,
	}
}


func (di *DropIndexes) Result() DropIndexesResult { return di.result }

func (di *DropIndexes) processResponse(info driver.ResponseInfo) error {
	var err error
	di.result, err = buildDropIndexesResult(info.ServerResponse)
	return err
}


func (di *DropIndexes) Execute(ctx context.Context) error {
	if di.deployment == nil {
		return errors.New("the DropIndexes operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         di.command,
		ProcessResponseFn: di.processResponse,
		Client:            di.session,
		Clock:             di.clock,
		CommandMonitor:    di.monitor,
		Crypt:             di.crypt,
		Database:          di.database,
		Deployment:        di.deployment,
		MaxTime:           di.maxTime,
		Selector:          di.selector,
		WriteConcern:      di.writeConcern,
		ServerAPI:         di.serverAPI,
		Timeout:           di.timeout,
		Name:              driverutil.DropIndexesOp,
		Authenticator:     di.authenticator,
	}.Execute(ctx)

}

func (di *DropIndexes) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "dropIndexes", di.collection)

	switch t := di.index.(type) {
	case string:
		dst = bsoncore.AppendStringElement(dst, "index", t)
	case bsoncore.Document:
		if di.index != nil {
			dst = bsoncore.AppendDocumentElement(dst, "index", t)
		}
	}

	return dst, nil
}


func (di *DropIndexes) Index(index any) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.index = index
	return di
}


func (di *DropIndexes) MaxTime(maxTime *time.Duration) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.maxTime = maxTime
	return di
}


func (di *DropIndexes) Session(session *session.Client) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.session = session
	return di
}


func (di *DropIndexes) ClusterClock(clock *session.ClusterClock) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.clock = clock
	return di
}


func (di *DropIndexes) Collection(collection string) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.collection = collection
	return di
}


func (di *DropIndexes) CommandMonitor(monitor *event.CommandMonitor) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.monitor = monitor
	return di
}


func (di *DropIndexes) Crypt(crypt driver.Crypt) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.crypt = crypt
	return di
}


func (di *DropIndexes) Database(database string) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.database = database
	return di
}


func (di *DropIndexes) Deployment(deployment driver.Deployment) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.deployment = deployment
	return di
}


func (di *DropIndexes) ServerSelector(selector description.ServerSelector) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.selector = selector
	return di
}


func (di *DropIndexes) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.writeConcern = writeConcern
	return di
}


func (di *DropIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.serverAPI = serverAPI
	return di
}


func (di *DropIndexes) Timeout(timeout *time.Duration) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.timeout = timeout
	return di
}


func (di *DropIndexes) Authenticator(authenticator driver.Authenticator) *DropIndexes {
	if di == nil {
		di = new(DropIndexes)
	}

	di.authenticator = authenticator
	return di
}
