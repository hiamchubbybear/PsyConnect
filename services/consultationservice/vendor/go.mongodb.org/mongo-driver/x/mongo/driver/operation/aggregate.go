





package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type Aggregate struct {
	authenticator            driver.Authenticator
	allowDiskUse             *bool
	batchSize                *int32
	bypassDocumentValidation *bool
	collation                bsoncore.Document
	comment                  *string
	hint                     bsoncore.Value
	maxTime                  *time.Duration
	pipeline                 bsoncore.Document
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	database                 string
	deployment               driver.Deployment
	readConcern              *readconcern.ReadConcern
	readPreference           *readpref.ReadPref
	retry                    *driver.RetryMode
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	crypt                    driver.Crypt
	serverAPI                *driver.ServerAPIOptions
	let                      bsoncore.Document
	hasOutputStage           bool
	customOptions            map[string]bsoncore.Value
	timeout                  *time.Duration
	omitCSOTMaxTimeMS        bool

	result driver.CursorResponse
}


func NewAggregate(pipeline bsoncore.Document) *Aggregate {
	return &Aggregate{
		pipeline: pipeline,
	}
}


func (a *Aggregate) Result(opts driver.CursorOptions) (*driver.BatchCursor, error) {

	clientSession := a.session

	clock := a.clock
	opts.ServerAPI = a.serverAPI
	return driver.NewBatchCursor(a.result, clientSession, clock, opts)
}



func (a *Aggregate) ResultCursorResponse() driver.CursorResponse {
	return a.result
}

func (a *Aggregate) processResponse(info driver.ResponseInfo) error {
	var err error

	a.result, err = driver.NewCursorResponse(info)
	return err

}


func (a *Aggregate) Execute(ctx context.Context) error {
	if a.deployment == nil {
		return errors.New("the Aggregate operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         a.command,
		ProcessResponseFn: a.processResponse,

		Client:                         a.session,
		Clock:                          a.clock,
		CommandMonitor:                 a.monitor,
		Database:                       a.database,
		Deployment:                     a.deployment,
		ReadConcern:                    a.readConcern,
		ReadPreference:                 a.readPreference,
		Type:                           driver.Read,
		RetryMode:                      a.retry,
		Selector:                       a.selector,
		WriteConcern:                   a.writeConcern,
		Crypt:                          a.crypt,
		MinimumWriteConcernWireVersion: 5,
		ServerAPI:                      a.serverAPI,
		IsOutputAggregate:              a.hasOutputStage,
		MaxTime:                        a.maxTime,
		Timeout:                        a.timeout,
		Name:                           driverutil.AggregateOp,
		OmitCSOTMaxTimeMS:              a.omitCSOTMaxTimeMS,
		Authenticator:                  a.authenticator,
	}.Execute(ctx)

}

func (a *Aggregate) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	header := bsoncore.Value{Type: bsontype.String, Data: bsoncore.AppendString(nil, a.collection)}
	if a.collection == "" {
		header = bsoncore.Value{Type: bsontype.Int32, Data: []byte{0x01, 0x00, 0x00, 0x00}}
	}
	dst = bsoncore.AppendValueElement(dst, "aggregate", header)

	cursorIdx, cursorDoc := bsoncore.AppendDocumentStart(nil)
	if a.allowDiskUse != nil {

		dst = bsoncore.AppendBooleanElement(dst, "allowDiskUse", *a.allowDiskUse)
	}
	if a.batchSize != nil {
		cursorDoc = bsoncore.AppendInt32Element(cursorDoc, "batchSize", *a.batchSize)
	}
	if a.bypassDocumentValidation != nil {

		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *a.bypassDocumentValidation)
	}
	if a.collation != nil {

		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", a.collation)
	}
	if a.comment != nil {

		dst = bsoncore.AppendStringElement(dst, "comment", *a.comment)
	}
	if a.hint.Type != bsontype.Type(0) {

		dst = bsoncore.AppendValueElement(dst, "hint", a.hint)
	}
	if a.pipeline != nil {

		dst = bsoncore.AppendArrayElement(dst, "pipeline", a.pipeline)
	}
	if a.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", a.let)
	}
	for optionName, optionValue := range a.customOptions {
		dst = bsoncore.AppendValueElement(dst, optionName, optionValue)
	}
	cursorDoc, _ = bsoncore.AppendDocumentEnd(cursorDoc, cursorIdx)
	dst = bsoncore.AppendDocumentElement(dst, "cursor", cursorDoc)

	return dst, nil
}


func (a *Aggregate) AllowDiskUse(allowDiskUse bool) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.allowDiskUse = &allowDiskUse
	return a
}


func (a *Aggregate) BatchSize(batchSize int32) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.batchSize = &batchSize
	return a
}


func (a *Aggregate) BypassDocumentValidation(bypassDocumentValidation bool) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.bypassDocumentValidation = &bypassDocumentValidation
	return a
}


func (a *Aggregate) Collation(collation bsoncore.Document) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.collation = collation
	return a
}


func (a *Aggregate) Comment(comment string) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.comment = &comment
	return a
}


func (a *Aggregate) Hint(hint bsoncore.Value) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.hint = hint
	return a
}


func (a *Aggregate) MaxTime(maxTime *time.Duration) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.maxTime = maxTime
	return a
}


func (a *Aggregate) Pipeline(pipeline bsoncore.Document) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.pipeline = pipeline
	return a
}


func (a *Aggregate) Session(session *session.Client) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.session = session
	return a
}


func (a *Aggregate) ClusterClock(clock *session.ClusterClock) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.clock = clock
	return a
}


func (a *Aggregate) Collection(collection string) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.collection = collection
	return a
}


func (a *Aggregate) CommandMonitor(monitor *event.CommandMonitor) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.monitor = monitor
	return a
}


func (a *Aggregate) Database(database string) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.database = database
	return a
}


func (a *Aggregate) Deployment(deployment driver.Deployment) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.deployment = deployment
	return a
}


func (a *Aggregate) ReadConcern(readConcern *readconcern.ReadConcern) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.readConcern = readConcern
	return a
}


func (a *Aggregate) ReadPreference(readPreference *readpref.ReadPref) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.readPreference = readPreference
	return a
}


func (a *Aggregate) ServerSelector(selector description.ServerSelector) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.selector = selector
	return a
}


func (a *Aggregate) WriteConcern(writeConcern *writeconcern.WriteConcern) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.writeConcern = writeConcern
	return a
}




func (a *Aggregate) Retry(retry driver.RetryMode) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.retry = &retry
	return a
}


func (a *Aggregate) Crypt(crypt driver.Crypt) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.crypt = crypt
	return a
}


func (a *Aggregate) ServerAPI(serverAPI *driver.ServerAPIOptions) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.serverAPI = serverAPI
	return a
}


func (a *Aggregate) Let(let bsoncore.Document) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.let = let
	return a
}



func (a *Aggregate) HasOutputStage(hos bool) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.hasOutputStage = hos
	return a
}


func (a *Aggregate) CustomOptions(co map[string]bsoncore.Value) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.customOptions = co
	return a
}


func (a *Aggregate) Timeout(timeout *time.Duration) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.timeout = timeout
	return a
}




func (a *Aggregate) OmitCSOTMaxTimeMS(omit bool) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.omitCSOTMaxTimeMS = omit
	return a
}


func (a *Aggregate) Authenticator(authenticator driver.Authenticator) *Aggregate {
	if a == nil {
		a = new(Aggregate)
	}

	a.authenticator = authenticator
	return a
}
