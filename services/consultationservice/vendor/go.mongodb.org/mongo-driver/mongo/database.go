





package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/internal/csfle"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var (
	defaultRunCmdOpts = []*options.RunCmdOptions{options.RunCmd().SetReadPreference(readpref.Primary())}
)


type Database struct {
	client         *Client
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
}

func newDatabase(client *Client, name string, opts ...*options.DatabaseOptions) *Database {
	dbOpt := options.MergeDatabaseOptions(opts...)

	rc := client.readConcern
	if dbOpt.ReadConcern != nil {
		rc = dbOpt.ReadConcern
	}

	rp := client.readPreference
	if dbOpt.ReadPreference != nil {
		rp = dbOpt.ReadPreference
	}

	wc := client.writeConcern
	if dbOpt.WriteConcern != nil {
		wc = dbOpt.WriteConcern
	}

	bsonOpts := client.bsonOpts
	if dbOpt.BSONOptions != nil {
		bsonOpts = dbOpt.BSONOptions
	}

	reg := client.registry
	if dbOpt.Registry != nil {
		reg = dbOpt.Registry
	}

	db := &Database{
		client:         client,
		name:           name,
		readPreference: rp,
		readConcern:    rc,
		writeConcern:   wc,
		bsonOpts:       bsonOpts,
		registry:       reg,
	}

	db.readSelector = description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(db.readPreference),
		description.LatencySelector(db.client.localThreshold),
	})

	db.writeSelector = description.CompositeSelector([]description.ServerSelector{
		description.WriteSelector(),
		description.LatencySelector(db.client.localThreshold),
	})

	return db
}


func (db *Database) Client() *Client {
	return db.client
}


func (db *Database) Name() string {
	return db.name
}


func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return newCollection(db, name, opts...)
}













func (db *Database) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*Cursor, error) {
	a := aggregateParams{
		ctx:            ctx,
		pipeline:       pipeline,
		client:         db.client,
		registry:       db.registry,
		readConcern:    db.readConcern,
		writeConcern:   db.writeConcern,
		retryRead:      db.client.retryReads,
		db:             db.name,
		readSelector:   db.readSelector,
		writeSelector:  db.writeSelector,
		readPreference: db.readPreference,
		opts:           opts,
	}
	return aggregate(a)
}

func (db *Database) processRunCommand(ctx context.Context, cmd interface{},
	cursorCommand bool, opts ...*options.RunCmdOptions) (*operation.Command, *session.Client, error) {
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
	}

	err := db.client.validSession(sess)
	if err != nil {
		return nil, sess, err
	}

	ro := options.MergeRunCmdOptions(append(defaultRunCmdOpts, opts...)...)
	if sess != nil && sess.TransactionRunning() && ro.ReadPreference != nil && ro.ReadPreference.Mode() != readpref.PrimaryMode {
		return nil, sess, errors.New("read preference in a transaction must be primary")
	}

	if isUnorderedMap(cmd) {
		return nil, sess, ErrMapForOrderedArgument{"cmd"}
	}

	runCmdDoc, err := marshal(cmd, db.bsonOpts, db.registry)
	if err != nil {
		return nil, sess, err
	}
	readSelect := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(ro.ReadPreference),
		description.LatencySelector(db.client.localThreshold),
	})
	if sess != nil && sess.PinnedServer != nil {
		readSelect = makePinnedSelector(sess, readSelect)
	}

	var op *operation.Command
	switch cursorCommand {
	case true:
		cursorOpts := db.client.createBaseCursorOptions()

		cursorOpts.MarshalValueEncoderFn = newEncoderFn(db.bsonOpts, db.registry)

		op = operation.NewCursorCommand(runCmdDoc, cursorOpts)
	default:
		op = operation.NewCommand(runCmdDoc)
	}

	return op.Session(sess).CommandMonitor(db.client.monitor).
		ServerSelector(readSelect).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).
		Crypt(db.client.cryptFLE).ReadPreference(ro.ReadPreference).ServerAPI(db.client.serverAPI).
		Timeout(db.client.timeout).Logger(db.client.logger).Authenticator(db.client.authenticator), sess, nil
}



















func (db *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *SingleResult {
	if ctx == nil {
		ctx = context.Background()
	}

	op, sess, err := db.processRunCommand(ctx, runCommand, false, opts...)
	defer closeImplicitSession(sess)
	if err != nil {
		return &SingleResult{err: err}
	}

	err = op.Execute(ctx)
	
	_, convErr := processWriteError(err)
	return &SingleResult{
		ctx:      ctx,
		err:      convErr,
		rdr:      bson.Raw(op.Result()),
		bsonOpts: db.bsonOpts,
		reg:      db.registry,
	}
}















func (db *Database) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	op, sess, err := db.processRunCommand(ctx, runCommand, true, opts...)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	if err = op.Execute(ctx); err != nil {
		closeImplicitSession(sess)
		if errors.Is(err, driver.ErrNoCursor) {
			return nil, errors.New(
				"database response does not contain a cursor; try using RunCommand instead")
		}
		return nil, replaceErrors(err)
	}

	bc, err := op.ResultCursor()
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, db.bsonOpts, db.registry, sess)
	return cursor, replaceErrors(err)
}



func (db *Database) Drop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
		defer sess.EndSession()
	}

	err := db.client.validSession(sess)
	if err != nil {
		return err
	}

	wc := db.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, db.writeSelector)

	op := operation.NewDropDatabase().
		Session(sess).WriteConcern(wc).CommandMonitor(db.client.monitor).
		ServerSelector(selector).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.cryptFLE).
		ServerAPI(db.client.serverAPI).Authenticator(db.client.authenticator)

	err = op.Execute(ctx)

	driverErr, ok := err.(driver.Error)
	if err != nil && (!ok || !driverErr.NamespaceNotFound()) {
		return replaceErrors(err)
	}
	return nil
}















func (db *Database) ListCollectionSpecifications(ctx context.Context, filter interface{},
	opts ...*options.ListCollectionsOptions) ([]*CollectionSpecification, error) {

	cursor, err := db.ListCollections(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var specs []*CollectionSpecification
	err = cursor.All(ctx, &specs)
	if err != nil {
		return nil, err
	}

	for _, spec := range specs {
		
		
		if spec.IDIndex != nil && spec.IDIndex.Namespace == "" {
			spec.IDIndex.Namespace = db.name + "." + spec.Name
		}
	}
	return specs, nil
}














func (db *Database) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	filterDoc, err := marshal(filter, db.bsonOpts, db.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
	}

	err = db.client.validSession(sess)
	if err != nil {
		closeImplicitSession(sess)
		return nil, err
	}

	selector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(db.client.localThreshold),
	})
	selector = makeReadPrefSelector(sess, selector, db.client.localThreshold)

	lco := options.MergeListCollectionsOptions(opts...)
	op := operation.NewListCollections(filterDoc).
		Session(sess).ReadPreference(db.readPreference).CommandMonitor(db.client.monitor).
		ServerSelector(selector).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.cryptFLE).
		ServerAPI(db.client.serverAPI).Timeout(db.client.timeout).Authenticator(db.client.authenticator)

	cursorOpts := db.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(db.bsonOpts, db.registry)

	if lco.NameOnly != nil {
		op = op.NameOnly(*lco.NameOnly)
	}
	if lco.BatchSize != nil {
		cursorOpts.BatchSize = *lco.BatchSize
		op = op.BatchSize(*lco.BatchSize)
	}
	if lco.AuthorizedCollections != nil {
		op = op.AuthorizedCollections(*lco.AuthorizedCollections)
	}

	retry := driver.RetryNone
	if db.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	bc, err := op.Result(cursorOpts)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, db.bsonOpts, db.registry, sess)
	return cursor, replaceErrors(err)
}















func (db *Database) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	opts = append(opts, options.ListCollections().SetNameOnly(true))

	res, err := db.ListCollections(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)

	names := make([]string, 0)
	for res.Next(ctx) {
		elem, err := res.Current.LookupErr("name")
		if err != nil {
			return nil, err
		}

		if elem.Type != bson.TypeString {
			return nil, fmt.Errorf("incorrect type for 'name'. got %v. want %v", elem.Type, bson.TypeString)
		}

		elemName := elem.StringValue()
		names = append(names, elemName)
	}

	res.Close(ctx)
	return names, nil
}


func (db *Database) ReadConcern() *readconcern.ReadConcern {
	return db.readConcern
}


func (db *Database) ReadPreference() *readpref.ReadPref {
	return db.readPreference
}


func (db *Database) WriteConcern() *writeconcern.WriteConcern {
	return db.writeConcern
}














func (db *Database) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {

	csConfig := changeStreamConfig{
		readConcern:    db.readConcern,
		readPreference: db.readPreference,
		client:         db.client,
		registry:       db.registry,
		streamType:     DatabaseStream,
		databaseName:   db.Name(),
		crypt:          db.client.cryptFLE,
	}
	return newChangeStream(ctx, csConfig, pipeline, opts...)
}









func (db *Database) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	cco := options.MergeCreateCollectionOptions(opts...)
	
	
	ef := cco.EncryptedFields
	
	if ef == nil {
		ef = db.getEncryptedFieldsFromMap(name)
	}
	if ef != nil {
		return db.createCollectionWithEncryptedFields(ctx, name, ef, opts...)
	}

	return db.createCollection(ctx, name, opts...)
}



func (db *Database) getEncryptedFieldsFromServer(ctx context.Context, collectionName string) (interface{}, error) {
	
	collSpecs, err := db.ListCollectionSpecifications(ctx, bson.D{{"name", collectionName}})
	if err != nil {
		return nil, err
	}
	if len(collSpecs) == 0 {
		return nil, nil
	}
	if len(collSpecs) > 1 {
		return nil, fmt.Errorf("expected 1 or 0 results from listCollections, got %v", len(collSpecs))
	}
	collSpec := collSpecs[0]
	rawValue, err := collSpec.Options.LookupErr("encryptedFields")
	if errors.Is(err, bsoncore.ErrElementNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	encryptedFields, ok := rawValue.DocumentOK()
	if !ok {
		return nil, fmt.Errorf("expected encryptedFields of %v to be document, got %v", collectionName, rawValue.Type)
	}

	return encryptedFields, nil
}



func (db *Database) getEncryptedFieldsFromMap(collectionName string) interface{} {
	
	efMap := db.client.encryptedFieldsMap
	if efMap == nil {
		return nil
	}

	namespace := db.name + "." + collectionName

	ef, ok := efMap[namespace]
	if ok {
		return ef
	}
	return nil
}


func (db *Database) createCollectionWithEncryptedFields(ctx context.Context, name string, ef interface{}, opts ...*options.CreateCollectionOptions) error {
	efBSON, err := marshal(ef, db.bsonOpts, db.registry)
	if err != nil {
		return fmt.Errorf("error transforming document: %w", err)
	}

	
	
	
	{
		const QEv2WireVersion = 21
		server, err := db.client.deployment.SelectServer(ctx, description.WriteSelector())
		if err != nil {
			return fmt.Errorf("error selecting server to check maxWireVersion: %w", err)
		}
		conn, err := server.Connection(ctx)
		if err != nil {
			return fmt.Errorf("error getting connection to check maxWireVersion: %w", err)
		}
		defer conn.Close()
		wireVersionRange := conn.Description().WireVersion
		if wireVersionRange.Max < QEv2WireVersion {
			return fmt.Errorf("Driver support of Queryable Encryption is incompatible with server. Upgrade server to use Queryable Encryption. Got maxWireVersion %v but need maxWireVersion >= %v", wireVersionRange.Max, QEv2WireVersion)
		}
	}

	

	stateCollectionOpts := options.CreateCollection().
		SetClusteredIndex(bson.D{{"key", bson.D{{"_id", 1}}}, {"unique", true}})
	
	escCollection, err := csfle.GetEncryptedStateCollectionName(efBSON, name, csfle.EncryptedStateCollection)
	if err != nil {
		return err
	}

	if err := db.createCollection(ctx, escCollection, stateCollectionOpts); err != nil {
		return err
	}

	
	ecocCollection, err := csfle.GetEncryptedStateCollectionName(efBSON, name, csfle.EncryptedCompactionCollection)
	if err != nil {
		return err
	}

	if err := db.createCollection(ctx, ecocCollection, stateCollectionOpts); err != nil {
		return err
	}

	
	op, err := db.createCollectionOperation(name, opts...)
	if err != nil {
		return err
	}

	op.EncryptedFields(efBSON)
	if err := db.executeCreateOperation(ctx, op); err != nil {
		return err
	}

	
	if _, err := db.Collection(name).Indexes().CreateOne(ctx, IndexModel{Keys: bson.D{{"__safeContent__", 1}}}); err != nil {
		return fmt.Errorf("error creating safeContent index: %w", err)
	}

	return nil
}


func (db *Database) createCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	op, err := db.createCollectionOperation(name, opts...)
	if err != nil {
		return err
	}
	return db.executeCreateOperation(ctx, op)
}

func (db *Database) createCollectionOperation(name string, opts ...*options.CreateCollectionOptions) (*operation.Create, error) {
	cco := options.MergeCreateCollectionOptions(opts...)
	op := operation.NewCreate(name).ServerAPI(db.client.serverAPI).Authenticator(db.client.authenticator)

	if cco.Capped != nil {
		op.Capped(*cco.Capped)
	}
	if cco.Collation != nil {
		op.Collation(bsoncore.Document(cco.Collation.ToDocument()))
	}
	if cco.ChangeStreamPreAndPostImages != nil {
		csppi, err := marshal(cco.ChangeStreamPreAndPostImages, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.ChangeStreamPreAndPostImages(csppi)
	}
	if cco.DefaultIndexOptions != nil {
		idx, doc := bsoncore.AppendDocumentStart(nil)
		if cco.DefaultIndexOptions.StorageEngine != nil {
			storageEngine, err := marshal(cco.DefaultIndexOptions.StorageEngine, db.bsonOpts, db.registry)
			if err != nil {
				return nil, err
			}

			doc = bsoncore.AppendDocumentElement(doc, "storageEngine", storageEngine)
		}
		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		if err != nil {
			return nil, err
		}

		op.IndexOptionDefaults(doc)
	}
	if cco.MaxDocuments != nil {
		op.Max(*cco.MaxDocuments)
	}
	if cco.SizeInBytes != nil {
		op.Size(*cco.SizeInBytes)
	}
	if cco.StorageEngine != nil {
		storageEngine, err := marshal(cco.StorageEngine, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.StorageEngine(storageEngine)
	}
	if cco.ValidationAction != nil {
		op.ValidationAction(*cco.ValidationAction)
	}
	if cco.ValidationLevel != nil {
		op.ValidationLevel(*cco.ValidationLevel)
	}
	if cco.Validator != nil {
		validator, err := marshal(cco.Validator, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.Validator(validator)
	}
	if cco.ExpireAfterSeconds != nil {
		op.ExpireAfterSeconds(*cco.ExpireAfterSeconds)
	}
	if cco.TimeSeriesOptions != nil {
		idx, doc := bsoncore.AppendDocumentStart(nil)
		doc = bsoncore.AppendStringElement(doc, "timeField", cco.TimeSeriesOptions.TimeField)

		if cco.TimeSeriesOptions.MetaField != nil {
			doc = bsoncore.AppendStringElement(doc, "metaField", *cco.TimeSeriesOptions.MetaField)
		}
		if cco.TimeSeriesOptions.Granularity != nil {
			doc = bsoncore.AppendStringElement(doc, "granularity", *cco.TimeSeriesOptions.Granularity)
		}

		if cco.TimeSeriesOptions.BucketMaxSpan != nil {
			bmss := int64(*cco.TimeSeriesOptions.BucketMaxSpan / time.Second)

			doc = bsoncore.AppendInt64Element(doc, "bucketMaxSpanSeconds", bmss)
		}

		if cco.TimeSeriesOptions.BucketRounding != nil {
			brs := int64(*cco.TimeSeriesOptions.BucketRounding / time.Second)

			doc = bsoncore.AppendInt64Element(doc, "bucketRoundingSeconds", brs)
		}

		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		if err != nil {
			return nil, err
		}

		op.TimeSeries(doc)
	}
	if cco.ClusteredIndex != nil {
		clusteredIndex, err := marshal(cco.ClusteredIndex, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.ClusteredIndex(clusteredIndex)
	}

	return op, nil
}














func (db *Database) CreateView(ctx context.Context, viewName, viewOn string, pipeline interface{},
	opts ...*options.CreateViewOptions) error {

	pipelineArray, _, err := marshalAggregatePipeline(pipeline, db.bsonOpts, db.registry)
	if err != nil {
		return err
	}

	op := operation.NewCreate(viewName).
		ViewOn(viewOn).
		Pipeline(pipelineArray).
		ServerAPI(db.client.serverAPI).
		Authenticator(db.client.authenticator)
	cvo := options.MergeCreateViewOptions(opts...)
	if cvo.Collation != nil {
		op.Collation(bsoncore.Document(cvo.Collation.ToDocument()))
	}

	return db.executeCreateOperation(ctx, op)
}

func (db *Database) executeCreateOperation(ctx context.Context, op *operation.Create) error {
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
		defer sess.EndSession()
	}

	err := db.client.validSession(sess)
	if err != nil {
		return err
	}

	wc := db.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, db.writeSelector)
	op = op.Session(sess).
		WriteConcern(wc).
		CommandMonitor(db.client.monitor).
		ServerSelector(selector).
		ClusterClock(db.client.clock).
		Database(db.name).
		Deployment(db.client.deployment).
		Crypt(db.client.cryptFLE)

	return replaceErrors(op.Execute(ctx))
}
