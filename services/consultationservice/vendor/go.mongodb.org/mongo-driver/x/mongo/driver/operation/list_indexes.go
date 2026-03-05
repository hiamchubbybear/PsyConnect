





package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type ListIndexes struct {
	authenticator driver.Authenticator
	batchSize     *int32
	maxTime       *time.Duration
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	retry         *driver.RetryMode
	crypt         driver.Crypt
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration

	result driver.CursorResponse
}


func NewListIndexes() *ListIndexes {
	return &ListIndexes{}
}


func (li *ListIndexes) Result(opts driver.CursorOptions) (*driver.BatchCursor, error) {

	clientSession := li.session

	clock := li.clock
	opts.ServerAPI = li.serverAPI
	return driver.NewBatchCursor(li.result, clientSession, clock, opts)
}

func (li *ListIndexes) processResponse(info driver.ResponseInfo) error {
	var err error

	li.result, err = driver.NewCursorResponse(info)
	return err

}


func (li *ListIndexes) Execute(ctx context.Context) error {
	if li.deployment == nil {
		return errors.New("the ListIndexes operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         li.command,
		ProcessResponseFn: li.processResponse,

		Client:         li.session,
		Clock:          li.clock,
		CommandMonitor: li.monitor,
		Database:       li.database,
		Deployment:     li.deployment,
		MaxTime:        li.maxTime,
		Selector:       li.selector,
		Crypt:          li.crypt,
		Legacy:         driver.LegacyListIndexes,
		RetryMode:      li.retry,
		Type:           driver.Read,
		ServerAPI:      li.serverAPI,
		Timeout:        li.timeout,
		Name:           driverutil.ListIndexesOp,
		Authenticator:  li.authenticator,
	}.Execute(ctx)

}

func (li *ListIndexes) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "listIndexes", li.collection)
	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)

	if li.batchSize != nil {

		cursorDoc = bsoncore.AppendInt32Element(cursorDoc, "batchSize", *li.batchSize)
	}
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)
	dst = bsoncore.AppendDocumentElement(dst, "cursor", cursorDoc)

	return dst, nil
}


func (li *ListIndexes) BatchSize(batchSize int32) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.batchSize = &batchSize
	return li
}


func (li *ListIndexes) MaxTime(maxTime *time.Duration) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.maxTime = maxTime
	return li
}


func (li *ListIndexes) Session(session *session.Client) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.session = session
	return li
}


func (li *ListIndexes) ClusterClock(clock *session.ClusterClock) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.clock = clock
	return li
}


func (li *ListIndexes) Collection(collection string) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.collection = collection
	return li
}


func (li *ListIndexes) CommandMonitor(monitor *event.CommandMonitor) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.monitor = monitor
	return li
}


func (li *ListIndexes) Database(database string) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.database = database
	return li
}


func (li *ListIndexes) Deployment(deployment driver.Deployment) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.deployment = deployment
	return li
}


func (li *ListIndexes) ServerSelector(selector description.ServerSelector) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.selector = selector
	return li
}



func (li *ListIndexes) Retry(retry driver.RetryMode) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.retry = &retry
	return li
}


func (li *ListIndexes) Crypt(crypt driver.Crypt) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.crypt = crypt
	return li
}


func (li *ListIndexes) ServerAPI(serverAPI *driver.ServerAPIOptions) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.serverAPI = serverAPI
	return li
}


func (li *ListIndexes) Timeout(timeout *time.Duration) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.timeout = timeout
	return li
}


func (li *ListIndexes) Authenticator(authenticator driver.Authenticator) *ListIndexes {
	if li == nil {
		li = new(ListIndexes)
	}

	li.authenticator = authenticator
	return li
}
