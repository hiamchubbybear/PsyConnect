





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type CreateIndexes struct {
	authenticator driver.Authenticator
	commitQuorum  bsoncore.Value
	indexes       bsoncore.Document
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
	result        CreateIndexesResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}


type CreateIndexesResult struct {
	
	CreatedCollectionAutomatically bool
	
	IndexesAfter int32
	
	IndexesBefore int32
}

func buildCreateIndexesResult(response bsoncore.Document) (CreateIndexesResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return CreateIndexesResult{}, err
	}
	cir := CreateIndexesResult{}
	for _, element := range elements {
		switch element.Key() {
		case "createdCollectionAutomatically":
			var ok bool
			cir.CreatedCollectionAutomatically, ok = element.Value().BooleanOK()
			if !ok {
				return cir, fmt.Errorf("response field 'createdCollectionAutomatically' is type bool, but received BSON type %s", element.Value().Type)
			}
		case "indexesAfter":
			var ok bool
			cir.IndexesAfter, ok = element.Value().AsInt32OK()
			if !ok {
				return cir, fmt.Errorf("response field 'indexesAfter' is type int32, but received BSON type %s", element.Value().Type)
			}
		case "indexesBefore":
			var ok bool
			cir.IndexesBefore, ok = element.Value().AsInt32OK()
			if !ok {
				return cir, fmt.Errorf("response field 'indexesBefore' is type int32, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return cir, nil
}


func NewCreateIndexes(indexes bsoncore.Document) *CreateIndexes {
	return &CreateIndexes{
		indexes: indexes,
	}
}


func (ci *CreateIndexes) Result() CreateIndexesResult { return ci.result }

func (ci *CreateIndexes) processResponse(info driver.ResponseInfo) error {
	var err error
	ci.result, err = buildCreateIndexesResult(info.ServerResponse)
	return err
}


func (ci *CreateIndexes) Execute(ctx context.Context) error {
	if ci.deployment == nil {
		return errors.New("the CreateIndexes operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         ci.command,
		ProcessResponseFn: ci.processResponse,
		Client:            ci.session,
		Clock:             ci.clock,
		CommandMonitor:    ci.monitor,
		Crypt:             ci.crypt,
		Database:          ci.database,
		Deployment:        ci.deployment,
		MaxTime:           ci.maxTime,
		Selector:          ci.selector,
		WriteConcern:      ci.writeConcern,
		ServerAPI:         ci.serverAPI,
		Timeout:           ci.timeout,
		Name:              driverutil.CreateIndexesOp,
		Authenticator:     ci.authenticator,
	}.Execute(ctx)

}

func (ci *CreateIndexes) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "createIndexes", ci.collection)
	if ci.commitQuorum.Type != bsontype.Type(0) {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(9) {
			return nil, errors.New("the 'commitQuorum' command parameter requires a minimum server wire version of 9")
		}
		dst = bsoncore.AppendValueElement(dst, "commitQuorum", ci.commitQuorum)
	}
	if ci.indexes != nil {
		dst = bsoncore.AppendArrayElement(dst, "indexes", ci.indexes)
	}
	return dst, nil
}




func (ci *CreateIndexes) CommitQuorum(commitQuorum bsoncore.Value) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.commitQuorum = commitQuorum
	return ci
}


func (ci *CreateIndexes) Indexes(indexes bsoncore.Document) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.indexes = indexes
	return ci
}


func (ci *CreateIndexes) MaxTime(maxTime *time.Duration) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.maxTime = maxTime
	return ci
}


func (ci *CreateIndexes) Session(session *session.Client) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.session = session
	return ci
}


func (ci *CreateIndexes) ClusterClock(clock *session.ClusterClock) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.clock = clock
	return ci
}


func (ci *CreateIndexes) Collection(collection string) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.collection = collection
	return ci
}


func (ci *CreateIndexes) CommandMonitor(monitor *event.CommandMonitor) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.monitor = monitor
	return ci
}


func (ci *CreateIndexes) Crypt(crypt driver.Crypt) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.crypt = crypt
	return ci
}


func (ci *CreateIndexes) Database(database string) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.database = database
	return ci
}


func (ci *CreateIndexes) Deployment(deployment driver.Deployment) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.deployment = deployment
	return ci
}


func (ci *CreateIndexes) ServerSelector(selector description.ServerSelector) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.selector = selector
	return ci
}


func (ci *CreateIndexes) WriteConcern(writeConcern *writeconcern.WriteConcern) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.writeConcern = writeConcern
	return ci
}


func (ci *CreateIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.serverAPI = serverAPI
	return ci
}


func (ci *CreateIndexes) Timeout(timeout *time.Duration) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.timeout = timeout
	return ci
}


func (ci *CreateIndexes) Authenticator(authenticator driver.Authenticator) *CreateIndexes {
	if ci == nil {
		ci = new(CreateIndexes)
	}

	ci.authenticator = authenticator
	return ci
}
