





package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var (
	
	ErrMissingResumeToken = errors.New("cannot provide resume functionality when the resume token is missing")
	
	ErrNilCursor = errors.New("cursor is nil")

	minResumableLabelWireVersion int32 = 9 
	networkErrorLabel                  = "NetworkError"
	resumableErrorLabel                = "ResumableChangeStreamError"
	errorCursorNotFound          int32 = 43 

	
	resumableChangeStreamErrors = map[int32]struct{}{
		6:     {}, 
		7:     {}, 
		89:    {}, 
		91:    {}, 
		189:   {}, 
		262:   {}, 
		9001:  {}, 
		10107: {}, 
		11600: {}, 
		11602: {}, 
		13435: {}, 
		13436: {}, 
		63:    {}, 
		150:   {}, 
		13388: {}, 
		234:   {}, 
		133:   {}, 
	}
)





type ChangeStream struct {
	
	
	Current bson.Raw

	aggregate       *operation.Aggregate
	pipelineSlice   []bsoncore.Document
	pipelineOptions map[string]bsoncore.Value
	cursor          changeStreamCursor
	cursorOptions   driver.CursorOptions
	batch           []bsoncore.Document
	resumeToken     bson.Raw
	err             error
	sess            *session.Client
	client          *Client
	bsonOpts        *options.BSONOptions
	registry        *bsoncodec.Registry
	streamType      StreamType
	options         *options.ChangeStreamOptions
	selector        description.ServerSelector
	operationTime   *primitive.Timestamp
	wireVersion     *description.VersionRange
}

type changeStreamConfig struct {
	readConcern    *readconcern.ReadConcern
	readPreference *readpref.ReadPref
	client         *Client
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
	streamType     StreamType
	collectionName string
	databaseName   string
	crypt          driver.Crypt
}

func newChangeStream(ctx context.Context, config changeStreamConfig, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	cursorOpts := config.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(config.bsonOpts, config.registry)

	cs := &ChangeStream{
		client:     config.client,
		bsonOpts:   config.bsonOpts,
		registry:   config.registry,
		streamType: config.streamType,
		options:    options.MergeChangeStreamOptions(opts...),
		selector: description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(config.readPreference),
			description.LatencySelector(config.client.localThreshold),
		}),
		cursorOptions: cursorOpts,
	}

	cs.sess = sessionFromContext(ctx)
	if cs.sess == nil && cs.client.sessionPool != nil {
		cs.sess = session.NewImplicitClientSession(cs.client.sessionPool, cs.client.id)
	}
	if cs.err = cs.client.validSession(cs.sess); cs.err != nil {
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	}

	cs.aggregate = operation.NewAggregate(nil).
		ReadPreference(config.readPreference).ReadConcern(config.readConcern).
		Deployment(cs.client.deployment).ClusterClock(cs.client.clock).
		CommandMonitor(cs.client.monitor).Session(cs.sess).ServerSelector(cs.selector).Retry(driver.RetryNone).
		ServerAPI(cs.client.serverAPI).Crypt(config.crypt).Timeout(cs.client.timeout).
		Authenticator(cs.client.authenticator)

	if cs.options.Collation != nil {
		cs.aggregate.Collation(bsoncore.Document(cs.options.Collation.ToDocument()))
	}
	if comment := cs.options.Comment; comment != nil {
		cs.aggregate.Comment(*comment)

		commentVal, err := marshalValue(comment, cs.bsonOpts, cs.registry)
		if err != nil {
			return nil, err
		}
		cs.cursorOptions.Comment = commentVal
	}
	if cs.options.BatchSize != nil {
		cs.aggregate.BatchSize(*cs.options.BatchSize)
		cs.cursorOptions.BatchSize = *cs.options.BatchSize
	}
	if cs.options.MaxAwaitTime != nil {
		cs.cursorOptions.MaxTimeMS = int64(*cs.options.MaxAwaitTime / time.Millisecond)
	}
	if cs.options.Custom != nil {
		
		
		customOptions := make(map[string]bsoncore.Value)
		for optionName, optionValue := range cs.options.Custom {
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(cs.registry, optionValue)
			if err != nil {
				cs.err = err
				closeImplicitSession(cs.sess)
				return nil, cs.Err()
			}
			optionValueBSON := bsoncore.Value{Type: bsonType, Data: bsonData}
			customOptions[optionName] = optionValueBSON
		}
		cs.aggregate.CustomOptions(customOptions)
	}
	if cs.options.CustomPipeline != nil {
		
		
		cs.pipelineOptions = make(map[string]bsoncore.Value)
		for optionName, optionValue := range cs.options.CustomPipeline {
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(cs.registry, optionValue)
			if err != nil {
				cs.err = err
				closeImplicitSession(cs.sess)
				return nil, cs.Err()
			}
			optionValueBSON := bsoncore.Value{Type: bsonType, Data: bsonData}
			cs.pipelineOptions[optionName] = optionValueBSON
		}
	}

	switch cs.streamType {
	case ClientStream:
		cs.aggregate.Database("admin")
	case DatabaseStream:
		cs.aggregate.Database(config.databaseName)
	case CollectionStream:
		cs.aggregate.Collection(config.collectionName).Database(config.databaseName)
	default:
		closeImplicitSession(cs.sess)
		return nil, fmt.Errorf("must supply a valid StreamType in config, instead of %v", cs.streamType)
	}

	
	
	resumeToken := cs.options.StartAfter
	if resumeToken == nil {
		resumeToken = cs.options.ResumeAfter
	}
	var marshaledToken bson.Raw
	if resumeToken != nil {
		if marshaledToken, cs.err = bson.Marshal(resumeToken); cs.err != nil {
			closeImplicitSession(cs.sess)
			return nil, cs.Err()
		}
	}
	cs.resumeToken = marshaledToken

	if cs.err = cs.buildPipelineSlice(pipeline); cs.err != nil {
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	}
	var pipelineArr bsoncore.Document
	pipelineArr, cs.err = cs.pipelineToBSON()
	cs.aggregate.Pipeline(pipelineArr)

	if cs.err = cs.executeOperation(ctx, false); cs.err != nil {
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	}

	return cs, cs.Err()
}

func (cs *ChangeStream) createOperationDeployment(server driver.Server, connection driver.Connection) driver.Deployment {
	return &changeStreamDeployment{
		topologyKind: cs.client.deployment.Kind(),
		server:       server,
		conn:         connection,
	}
}

func (cs *ChangeStream) executeOperation(ctx context.Context, resuming bool) error {
	var server driver.Server
	var conn driver.Connection

	if server, cs.err = cs.client.deployment.SelectServer(ctx, cs.selector); cs.err != nil {
		return cs.Err()
	}
	if conn, cs.err = server.Connection(ctx); cs.err != nil {
		return cs.Err()
	}
	defer conn.Close()
	cs.wireVersion = conn.Description().WireVersion

	cs.aggregate.Deployment(cs.createOperationDeployment(server, conn))

	if resuming {
		cs.replaceOptions(cs.wireVersion)

		csOptDoc, err := cs.createPipelineOptionsDoc()
		if err != nil {
			return err
		}
		pipIdx, pipDoc := bsoncore.AppendDocumentStart(nil)
		pipDoc = bsoncore.AppendDocumentElement(pipDoc, "$changeStream", csOptDoc)
		if pipDoc, cs.err = bsoncore.AppendDocumentEnd(pipDoc, pipIdx); cs.err != nil {
			return cs.Err()
		}
		cs.pipelineSlice[0] = pipDoc

		var plArr bsoncore.Document
		if plArr, cs.err = cs.pipelineToBSON(); cs.err != nil {
			return cs.Err()
		}
		cs.aggregate.Pipeline(plArr)
	}

	
	
	
	if cs.client.timeout != nil && !csot.IsTimeoutContext(ctx) {
		newCtx, cancelFunc := csot.MakeTimeoutContext(ctx, *cs.client.timeout)
		
		ctx = newCtx
		
		defer cancelFunc()
	}

	
	
	var retries int
	if cs.client.retryReads {
		retries = 1
	}
	if csot.IsTimeoutContext(ctx) {
		retries = -1
	}

	var err error
AggregateExecuteLoop:
	for {
		err = cs.aggregate.Execute(ctx)
		
		if err == nil || retries == 0 {
			break AggregateExecuteLoop
		}

		switch tt := err.(type) {
		case driver.Error:
			
			if !tt.RetryableRead() {
				break AggregateExecuteLoop
			}

			
			
			retries--
			server, err = cs.client.deployment.SelectServer(ctx, cs.selector)
			if err != nil {
				break AggregateExecuteLoop
			}

			conn.Close()
			conn, err = server.Connection(ctx)
			if err != nil {
				break AggregateExecuteLoop
			}
			defer conn.Close()

			
			cs.wireVersion = conn.Description().WireVersion

			
			cs.aggregate.Deployment(cs.createOperationDeployment(server, conn))
		default:
			
			break AggregateExecuteLoop
		}
	}
	if err != nil {
		cs.err = replaceErrors(err)
		return cs.err
	}

	cr := cs.aggregate.ResultCursorResponse()
	cr.Server = server

	cs.cursor, cs.err = driver.NewBatchCursor(cr, cs.sess, cs.client.clock, cs.cursorOptions)
	if cs.err = replaceErrors(cs.err); cs.err != nil {
		return cs.Err()
	}

	cs.updatePbrtFromCommand()
	if cs.options.StartAtOperationTime == nil && cs.options.ResumeAfter == nil &&
		cs.options.StartAfter == nil && cs.wireVersion.Max >= 7 &&
		cs.emptyBatch() && cs.resumeToken == nil {
		cs.operationTime = cs.sess.OperationTime
	}

	return cs.Err()
}


func (cs *ChangeStream) updatePbrtFromCommand() {
	
	if pbrt := cs.cursor.PostBatchResumeToken(); cs.emptyBatch() && pbrt != nil {
		cs.resumeToken = bson.Raw(pbrt)
	}
}

func (cs *ChangeStream) storeResumeToken() error {
	
	
	var tokenDoc bson.Raw
	if len(cs.batch) == 0 {
		if pbrt := cs.cursor.PostBatchResumeToken(); pbrt != nil {
			tokenDoc = bson.Raw(pbrt)
		}
	}

	if tokenDoc == nil {
		var ok bool
		tokenDoc, ok = cs.Current.Lookup("_id").DocumentOK()
		if !ok {
			_ = cs.Close(context.Background())
			return ErrMissingResumeToken
		}
	}

	cs.resumeToken = tokenDoc
	return nil
}

func (cs *ChangeStream) buildPipelineSlice(pipeline interface{}) error {
	val := reflect.ValueOf(pipeline)
	if !val.IsValid() || !(val.Kind() == reflect.Slice) {
		cs.err = errors.New("can only marshal slices and arrays into aggregation pipelines, but got invalid")
		return cs.err
	}

	cs.pipelineSlice = make([]bsoncore.Document, 0, val.Len()+1)

	csIdx, csDoc := bsoncore.AppendDocumentStart(nil)

	csDocTemp, err := cs.createPipelineOptionsDoc()
	if err != nil {
		return err
	}
	csDoc = bsoncore.AppendDocumentElement(csDoc, "$changeStream", csDocTemp)
	csDoc, cs.err = bsoncore.AppendDocumentEnd(csDoc, csIdx)
	if cs.err != nil {
		return cs.err
	}
	cs.pipelineSlice = append(cs.pipelineSlice, csDoc)

	for i := 0; i < val.Len(); i++ {
		var elem []byte
		elem, cs.err = marshal(val.Index(i).Interface(), cs.bsonOpts, cs.registry)
		if cs.err != nil {
			return cs.err
		}

		cs.pipelineSlice = append(cs.pipelineSlice, elem)
	}

	return cs.err
}

func (cs *ChangeStream) createPipelineOptionsDoc() (bsoncore.Document, error) {
	plDocIdx, plDoc := bsoncore.AppendDocumentStart(nil)

	if cs.streamType == ClientStream {
		plDoc = bsoncore.AppendBooleanElement(plDoc, "allChangesForCluster", true)
	}

	if cs.options.FullDocument != nil && *cs.options.FullDocument != options.Default {
		plDoc = bsoncore.AppendStringElement(plDoc, "fullDocument", string(*cs.options.FullDocument))
	}

	if cs.options.FullDocumentBeforeChange != nil {
		plDoc = bsoncore.AppendStringElement(plDoc, "fullDocumentBeforeChange", string(*cs.options.FullDocumentBeforeChange))
	}

	if cs.options.ResumeAfter != nil {
		var raDoc bsoncore.Document
		raDoc, cs.err = marshal(cs.options.ResumeAfter, cs.bsonOpts, cs.registry)
		if cs.err != nil {
			return nil, cs.err
		}

		plDoc = bsoncore.AppendDocumentElement(plDoc, "resumeAfter", raDoc)
	}

	if cs.options.ShowExpandedEvents != nil {
		plDoc = bsoncore.AppendBooleanElement(plDoc, "showExpandedEvents", *cs.options.ShowExpandedEvents)
	}

	if cs.options.StartAfter != nil {
		var saDoc bsoncore.Document
		saDoc, cs.err = marshal(cs.options.StartAfter, cs.bsonOpts, cs.registry)
		if cs.err != nil {
			return nil, cs.err
		}

		plDoc = bsoncore.AppendDocumentElement(plDoc, "startAfter", saDoc)
	}

	if cs.options.StartAtOperationTime != nil {
		plDoc = bsoncore.AppendTimestampElement(plDoc, "startAtOperationTime", cs.options.StartAtOperationTime.T, cs.options.StartAtOperationTime.I)
	}

	
	for optionName, optionValue := range cs.pipelineOptions {
		plDoc = bsoncore.AppendValueElement(plDoc, optionName, optionValue)
	}

	if plDoc, cs.err = bsoncore.AppendDocumentEnd(plDoc, plDocIdx); cs.err != nil {
		return nil, cs.err
	}

	return plDoc, nil
}

func (cs *ChangeStream) pipelineToBSON() (bsoncore.Document, error) {
	pipelineDocIdx, pipelineArr := bsoncore.AppendArrayStart(nil)
	for i, doc := range cs.pipelineSlice {
		pipelineArr = bsoncore.AppendDocumentElement(pipelineArr, strconv.Itoa(i), doc)
	}
	if pipelineArr, cs.err = bsoncore.AppendArrayEnd(pipelineArr, pipelineDocIdx); cs.err != nil {
		return nil, cs.err
	}
	return pipelineArr, cs.err
}

func (cs *ChangeStream) replaceOptions(wireVersion *description.VersionRange) {
	
	if cs.resumeToken != nil {
		cs.options.SetResumeAfter(cs.resumeToken)
		cs.options.SetStartAfter(nil)
		cs.options.SetStartAtOperationTime(nil)
		return
	}

	
	
	if (cs.sess.OperationTime != nil || cs.options.StartAtOperationTime != nil) && wireVersion.Max >= 7 {
		opTime := cs.options.StartAtOperationTime
		if cs.operationTime != nil {
			opTime = cs.sess.OperationTime
		}

		cs.options.SetStartAtOperationTime(opTime)
		cs.options.SetResumeAfter(nil)
		cs.options.SetStartAfter(nil)
		return
	}

	
	cs.options.SetResumeAfter(nil)
	cs.options.SetStartAfter(nil)
	cs.options.SetStartAtOperationTime(nil)
}


func (cs *ChangeStream) ID() int64 {
	if cs.cursor == nil {
		return 0
	}
	return cs.cursor.ID()
}



func (cs *ChangeStream) RemainingBatchLength() int {
	return len(cs.batch)
}




func (cs *ChangeStream) SetBatchSize(size int32) {
	
	
	cs.cursorOptions.BatchSize = size
	cs.cursor.SetBatchSize(size)
}



func (cs *ChangeStream) Decode(val interface{}) error {
	if cs.cursor == nil {
		return ErrNilCursor
	}

	dec, err := getDecoder(cs.Current, cs.bsonOpts, cs.registry)
	if err != nil {
		return fmt.Errorf("error configuring BSON decoder: %w", err)
	}
	return dec.Decode(val)
}


func (cs *ChangeStream) Err() error {
	if cs.err != nil {
		return replaceErrors(cs.err)
	}
	if cs.cursor == nil {
		return nil
	}

	return replaceErrors(cs.cursor.Err())
}



func (cs *ChangeStream) Close(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	defer closeImplicitSession(cs.sess)

	if cs.cursor == nil {
		return nil 
	}

	cs.err = replaceErrors(cs.cursor.Close(ctx))
	cs.cursor = nil
	return cs.Err()
}



func (cs *ChangeStream) ResumeToken() bson.Raw {
	return cs.resumeToken
}








func (cs *ChangeStream) Next(ctx context.Context) bool {
	return cs.next(ctx, false)
}












func (cs *ChangeStream) TryNext(ctx context.Context) bool {
	return cs.next(ctx, true)
}

func (cs *ChangeStream) next(ctx context.Context, nonBlocking bool) bool {
	
	if cs.err != nil {
		return false
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if len(cs.batch) == 0 {
		cs.loopNext(ctx, nonBlocking)
		if cs.err != nil {
			cs.err = replaceErrors(cs.err)
			return false
		}
		if len(cs.batch) == 0 {
			return false
		}
	}

	
	cs.Current = bson.Raw(cs.batch[0])
	cs.batch = cs.batch[1:]
	if cs.err = cs.storeResumeToken(); cs.err != nil {
		return false
	}
	return true
}

func (cs *ChangeStream) loopNext(ctx context.Context, nonBlocking bool) {
	for {
		if cs.cursor == nil {
			return
		}

		if cs.cursor.Next(ctx) {
			
			cs.batch, cs.err = cs.cursor.Batch().Documents()
			return
		}

		cs.err = replaceErrors(cs.cursor.Err())
		if cs.err == nil {
			
			if cs.ID() == 0 {
				return
			}

			
			
			cs.updatePbrtFromCommand()
			if nonBlocking {
				
				return
			}
			continue 
		}

		if !cs.isResumableError() {
			return
		}

		
		_ = cs.cursor.Close(ctx)
		if cs.err = cs.executeOperation(ctx, true); cs.err != nil {
			return
		}
	}
}

func (cs *ChangeStream) isResumableError() bool {
	var commandErr CommandError
	if !errors.As(cs.err, &commandErr) || commandErr.HasErrorLabel(networkErrorLabel) {
		
		return true
	}

	if commandErr.Code == errorCursorNotFound {
		return true
	}

	
	if cs.wireVersion != nil && cs.wireVersion.Includes(minResumableLabelWireVersion) {
		return commandErr.HasErrorLabel(resumableErrorLabel)
	}

	
	_, resumable := resumableChangeStreamErrors[commandErr.Code]
	return resumable
}


func (cs *ChangeStream) emptyBatch() bool {
	return cs.cursor.Batch().Empty()
}


type StreamType uint8



const (
	CollectionStream StreamType = iota
	DatabaseStream
	ClientStream
)
