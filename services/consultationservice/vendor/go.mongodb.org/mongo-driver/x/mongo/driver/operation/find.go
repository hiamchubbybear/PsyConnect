





package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type Find struct {
	authenticator       driver.Authenticator
	allowDiskUse        *bool
	allowPartialResults *bool
	awaitData           *bool
	batchSize           *int32
	collation           bsoncore.Document
	comment             *string
	filter              bsoncore.Document
	hint                bsoncore.Value
	let                 bsoncore.Document
	limit               *int64
	max                 bsoncore.Document
	maxTime             *time.Duration
	min                 bsoncore.Document
	noCursorTimeout     *bool
	oplogReplay         *bool
	projection          bsoncore.Document
	returnKey           *bool
	showRecordID        *bool
	singleBatch         *bool
	skip                *int64
	snapshot            *bool
	sort                bsoncore.Document
	tailable            *bool
	session             *session.Client
	clock               *session.ClusterClock
	collection          string
	monitor             *event.CommandMonitor
	crypt               driver.Crypt
	database            string
	deployment          driver.Deployment
	readConcern         *readconcern.ReadConcern
	readPreference      *readpref.ReadPref
	selector            description.ServerSelector
	retry               *driver.RetryMode
	result              driver.CursorResponse
	serverAPI           *driver.ServerAPIOptions
	timeout             *time.Duration
	omitCSOTMaxTimeMS   bool
	logger              *logger.Logger
}


func NewFind(filter bsoncore.Document) *Find {
	return &Find{
		filter: filter,
	}
}


func (f *Find) Result(opts driver.CursorOptions) (*driver.BatchCursor, error) {
	opts.ServerAPI = f.serverAPI
	return driver.NewBatchCursor(f.result, f.session, f.clock, opts)
}

func (f *Find) processResponse(info driver.ResponseInfo) error {
	var err error
	f.result, err = driver.NewCursorResponse(info)
	return err
}


func (f *Find) Execute(ctx context.Context) error {
	if f.deployment == nil {
		return errors.New("the Find operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         f.command,
		ProcessResponseFn: f.processResponse,
		RetryMode:         f.retry,
		Type:              driver.Read,
		Client:            f.session,
		Clock:             f.clock,
		CommandMonitor:    f.monitor,
		Crypt:             f.crypt,
		Database:          f.database,
		Deployment:        f.deployment,
		MaxTime:           f.maxTime,
		ReadConcern:       f.readConcern,
		ReadPreference:    f.readPreference,
		Selector:          f.selector,
		Legacy:            driver.LegacyFind,
		ServerAPI:         f.serverAPI,
		Timeout:           f.timeout,
		Logger:            f.logger,
		Name:              driverutil.FindOp,
		OmitCSOTMaxTimeMS: f.omitCSOTMaxTimeMS,
		Authenticator:     f.authenticator,
	}.Execute(ctx)

}

func (f *Find) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "find", f.collection)
	if f.allowDiskUse != nil {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(4) {
			return nil, errors.New("the 'allowDiskUse' command parameter requires a minimum server wire version of 4")
		}
		dst = bsoncore.AppendBooleanElement(dst, "allowDiskUse", *f.allowDiskUse)
	}
	if f.allowPartialResults != nil {
		dst = bsoncore.AppendBooleanElement(dst, "allowPartialResults", *f.allowPartialResults)
	}
	if f.awaitData != nil {
		dst = bsoncore.AppendBooleanElement(dst, "awaitData", *f.awaitData)
	}
	if f.batchSize != nil {
		dst = bsoncore.AppendInt32Element(dst, "batchSize", *f.batchSize)
	}
	if f.collation != nil {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", f.collation)
	}
	if f.comment != nil {
		dst = bsoncore.AppendStringElement(dst, "comment", *f.comment)
	}
	if f.filter != nil {
		dst = bsoncore.AppendDocumentElement(dst, "filter", f.filter)
	}
	if f.hint.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "hint", f.hint)
	}
	if f.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", f.let)
	}
	if f.limit != nil {
		dst = bsoncore.AppendInt64Element(dst, "limit", *f.limit)
	}
	if f.max != nil {
		dst = bsoncore.AppendDocumentElement(dst, "max", f.max)
	}
	if f.min != nil {
		dst = bsoncore.AppendDocumentElement(dst, "min", f.min)
	}
	if f.noCursorTimeout != nil {
		dst = bsoncore.AppendBooleanElement(dst, "noCursorTimeout", *f.noCursorTimeout)
	}
	if f.oplogReplay != nil {
		dst = bsoncore.AppendBooleanElement(dst, "oplogReplay", *f.oplogReplay)
	}
	if f.projection != nil {
		dst = bsoncore.AppendDocumentElement(dst, "projection", f.projection)
	}
	if f.returnKey != nil {
		dst = bsoncore.AppendBooleanElement(dst, "returnKey", *f.returnKey)
	}
	if f.showRecordID != nil {
		dst = bsoncore.AppendBooleanElement(dst, "showRecordId", *f.showRecordID)
	}
	if f.singleBatch != nil {
		dst = bsoncore.AppendBooleanElement(dst, "singleBatch", *f.singleBatch)
	}
	if f.skip != nil {
		dst = bsoncore.AppendInt64Element(dst, "skip", *f.skip)
	}
	if f.snapshot != nil {
		dst = bsoncore.AppendBooleanElement(dst, "snapshot", *f.snapshot)
	}
	if f.sort != nil {
		dst = bsoncore.AppendDocumentElement(dst, "sort", f.sort)
	}
	if f.tailable != nil {
		dst = bsoncore.AppendBooleanElement(dst, "tailable", *f.tailable)
	}
	return dst, nil
}


func (f *Find) AllowDiskUse(allowDiskUse bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.allowDiskUse = &allowDiskUse
	return f
}


func (f *Find) AllowPartialResults(allowPartialResults bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.allowPartialResults = &allowPartialResults
	return f
}


func (f *Find) AwaitData(awaitData bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.awaitData = &awaitData
	return f
}


func (f *Find) BatchSize(batchSize int32) *Find {
	if f == nil {
		f = new(Find)
	}

	f.batchSize = &batchSize
	return f
}


func (f *Find) Collation(collation bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.collation = collation
	return f
}


func (f *Find) Comment(comment string) *Find {
	if f == nil {
		f = new(Find)
	}

	f.comment = &comment
	return f
}


func (f *Find) Filter(filter bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.filter = filter
	return f
}


func (f *Find) Hint(hint bsoncore.Value) *Find {
	if f == nil {
		f = new(Find)
	}

	f.hint = hint
	return f
}


func (f *Find) Let(let bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.let = let
	return f
}


func (f *Find) Limit(limit int64) *Find {
	if f == nil {
		f = new(Find)
	}

	f.limit = &limit
	return f
}


func (f *Find) Max(max bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.max = max
	return f
}


func (f *Find) MaxTime(maxTime *time.Duration) *Find {
	if f == nil {
		f = new(Find)
	}

	f.maxTime = maxTime
	return f
}


func (f *Find) Min(min bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.min = min
	return f
}


func (f *Find) NoCursorTimeout(noCursorTimeout bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.noCursorTimeout = &noCursorTimeout
	return f
}


func (f *Find) OplogReplay(oplogReplay bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.oplogReplay = &oplogReplay
	return f
}


func (f *Find) Projection(projection bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.projection = projection
	return f
}


func (f *Find) ReturnKey(returnKey bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.returnKey = &returnKey
	return f
}


func (f *Find) ShowRecordID(showRecordID bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.showRecordID = &showRecordID
	return f
}


func (f *Find) SingleBatch(singleBatch bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.singleBatch = &singleBatch
	return f
}


func (f *Find) Skip(skip int64) *Find {
	if f == nil {
		f = new(Find)
	}

	f.skip = &skip
	return f
}


func (f *Find) Snapshot(snapshot bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.snapshot = &snapshot
	return f
}


func (f *Find) Sort(sort bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.sort = sort
	return f
}


func (f *Find) Tailable(tailable bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.tailable = &tailable
	return f
}


func (f *Find) Session(session *session.Client) *Find {
	if f == nil {
		f = new(Find)
	}

	f.session = session
	return f
}


func (f *Find) ClusterClock(clock *session.ClusterClock) *Find {
	if f == nil {
		f = new(Find)
	}

	f.clock = clock
	return f
}


func (f *Find) Collection(collection string) *Find {
	if f == nil {
		f = new(Find)
	}

	f.collection = collection
	return f
}


func (f *Find) CommandMonitor(monitor *event.CommandMonitor) *Find {
	if f == nil {
		f = new(Find)
	}

	f.monitor = monitor
	return f
}


func (f *Find) Crypt(crypt driver.Crypt) *Find {
	if f == nil {
		f = new(Find)
	}

	f.crypt = crypt
	return f
}


func (f *Find) Database(database string) *Find {
	if f == nil {
		f = new(Find)
	}

	f.database = database
	return f
}


func (f *Find) Deployment(deployment driver.Deployment) *Find {
	if f == nil {
		f = new(Find)
	}

	f.deployment = deployment
	return f
}


func (f *Find) ReadConcern(readConcern *readconcern.ReadConcern) *Find {
	if f == nil {
		f = new(Find)
	}

	f.readConcern = readConcern
	return f
}


func (f *Find) ReadPreference(readPreference *readpref.ReadPref) *Find {
	if f == nil {
		f = new(Find)
	}

	f.readPreference = readPreference
	return f
}


func (f *Find) ServerSelector(selector description.ServerSelector) *Find {
	if f == nil {
		f = new(Find)
	}

	f.selector = selector
	return f
}



func (f *Find) Retry(retry driver.RetryMode) *Find {
	if f == nil {
		f = new(Find)
	}

	f.retry = &retry
	return f
}


func (f *Find) ServerAPI(serverAPI *driver.ServerAPIOptions) *Find {
	if f == nil {
		f = new(Find)
	}

	f.serverAPI = serverAPI
	return f
}


func (f *Find) Timeout(timeout *time.Duration) *Find {
	if f == nil {
		f = new(Find)
	}

	f.timeout = timeout
	return f
}




func (f *Find) OmitCSOTMaxTimeMS(omit bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.omitCSOTMaxTimeMS = omit
	return f
}


func (f *Find) Logger(logger *logger.Logger) *Find {
	if f == nil {
		f = new(Find)
	}

	f.logger = logger
	return f
}


func (f *Find) Authenticator(authenticator driver.Authenticator) *Find {
	if f == nil {
		f = new(Find)
	}

	f.authenticator = authenticator
	return f
}
