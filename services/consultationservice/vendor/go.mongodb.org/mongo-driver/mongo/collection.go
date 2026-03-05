





package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
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


type Collection struct {
	client         *Client
	db             *Database
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
}


type aggregateParams struct {
	ctx            context.Context
	pipeline       interface{}
	client         *Client
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	retryRead      bool
	db             string
	col            string
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	readPreference *readpref.ReadPref
	opts           []*options.AggregateOptions
}

func closeImplicitSession(sess *session.Client) {
	if sess != nil && sess.IsImplicit {
		sess.EndSession()
	}
}

func newCollection(db *Database, name string, opts ...*options.CollectionOptions) *Collection {
	collOpt := options.MergeCollectionOptions(opts...)

	rc := db.readConcern
	if collOpt.ReadConcern != nil {
		rc = collOpt.ReadConcern
	}

	wc := db.writeConcern
	if collOpt.WriteConcern != nil {
		wc = collOpt.WriteConcern
	}

	rp := db.readPreference
	if collOpt.ReadPreference != nil {
		rp = collOpt.ReadPreference
	}

	bsonOpts := db.bsonOpts
	if collOpt.BSONOptions != nil {
		bsonOpts = collOpt.BSONOptions
	}

	reg := db.registry
	if collOpt.Registry != nil {
		reg = collOpt.Registry
	}

	readSelector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(rp),
		description.LatencySelector(db.client.localThreshold),
	})

	writeSelector := description.CompositeSelector([]description.ServerSelector{
		description.WriteSelector(),
		description.LatencySelector(db.client.localThreshold),
	})

	coll := &Collection{
		client:         db.client,
		db:             db,
		name:           name,
		readPreference: rp,
		readConcern:    rc,
		writeConcern:   wc,
		readSelector:   readSelector,
		writeSelector:  writeSelector,
		bsonOpts:       bsonOpts,
		registry:       reg,
	}

	return coll
}

func (coll *Collection) copy() *Collection {
	return &Collection{
		client:         coll.client,
		db:             coll.db,
		name:           coll.name,
		readConcern:    coll.readConcern,
		writeConcern:   coll.writeConcern,
		readPreference: coll.readPreference,
		readSelector:   coll.readSelector,
		writeSelector:  coll.writeSelector,
		registry:       coll.registry,
	}
}




func (coll *Collection) Clone(opts ...*options.CollectionOptions) (*Collection, error) {
	copyColl := coll.copy()
	optsColl := options.MergeCollectionOptions(opts...)

	if optsColl.ReadConcern != nil {
		copyColl.readConcern = optsColl.ReadConcern
	}

	if optsColl.WriteConcern != nil {
		copyColl.writeConcern = optsColl.WriteConcern
	}

	if optsColl.ReadPreference != nil {
		copyColl.readPreference = optsColl.ReadPreference
	}

	if optsColl.Registry != nil {
		copyColl.registry = optsColl.Registry
	}

	copyColl.readSelector = description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(copyColl.readPreference),
		description.LatencySelector(copyColl.client.localThreshold),
	})

	return copyColl, nil
}


func (coll *Collection) Name() string {
	return coll.name
}


func (coll *Collection) Database() *Database {
	return coll.db
}








func (coll *Collection) BulkWrite(ctx context.Context, models []WriteModel,
	opts ...*options.BulkWriteOptions) (*BulkWriteResult, error) {

	if len(models) == 0 {
		return nil, ErrEmptySlice
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err := coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	for _, model := range models {
		if model == nil {
			return nil, ErrNilDocument
		}
	}

	bwo := options.MergeBulkWriteOptions(opts...)

	op := bulkWrite{
		comment:                  bwo.Comment,
		ordered:                  bwo.Ordered,
		bypassDocumentValidation: bwo.BypassDocumentValidation,
		models:                   models,
		session:                  sess,
		collection:               coll,
		selector:                 selector,
		writeConcern:             wc,
		let:                      bwo.Let,
		bypassEmptyTsReplacement: bwo.BypassEmptyTsReplacement,
	}

	err = op.execute(ctx)

	return &op.result, replaceErrors(err)
}

func (coll *Collection) insert(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) ([]interface{}, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	result := make([]interface{}, len(documents))
	docs := make([]bsoncore.Document, len(documents))

	for i, doc := range documents {
		bsoncoreDoc, err := marshal(doc, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		bsoncoreDoc, id, err := ensureID(bsoncoreDoc, primitive.NilObjectID, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}

		docs[i] = bsoncoreDoc
		result[i] = id
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err := coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewInsert(docs...).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).Ordered(true).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout).Logger(coll.client.logger).
		Authenticator(coll.client.authenticator)
	imo := options.MergeInsertManyOptions(opts...)
	if imo.BypassDocumentValidation != nil && *imo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*imo.BypassDocumentValidation)
	}
	if imo.Comment != nil {
		comment, err := marshalValue(imo.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op = op.Comment(comment)
	}
	if imo.Ordered != nil {
		op = op.Ordered(*imo.Ordered)
	}
	if imo.BypassEmptyTsReplacement != nil {
		op = op.BypassEmptyTsReplacement(*imo.BypassEmptyTsReplacement)
	}
	retry := driver.RetryNone
	if coll.client.retryWrites {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	var wce driver.WriteCommandError
	if !errors.As(err, &wce) {
		return result, err
	}

	
	for i, we := range wce.WriteErrors {
		
		idIndex := int(we.Index) - i
		
		if imo.Ordered == nil || *imo.Ordered {
			result = result[:idIndex]
			break
		}
		result = append(result[:idIndex], result[idIndex+1:]...)
	}

	return result, err
}










func (coll *Collection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*InsertOneResult, error) {

	ioOpts := options.MergeInsertOneOptions(opts...)
	imOpts := options.InsertMany()

	if ioOpts.BypassDocumentValidation != nil && *ioOpts.BypassDocumentValidation {
		imOpts.SetBypassDocumentValidation(*ioOpts.BypassDocumentValidation)
	}
	if ioOpts.Comment != nil {
		imOpts.SetComment(ioOpts.Comment)
	}
	if ioOpts.BypassEmptyTsReplacement != nil {
		imOpts.BypassEmptyTsReplacement = ioOpts.BypassEmptyTsReplacement
	}
	res, err := coll.insert(ctx, []interface{}{document}, imOpts)

	rr, err := processWriteError(err)
	if rr&rrOne == 0 {
		return nil, err
	}
	return &InsertOneResult{InsertedID: res[0]}, err
}












func (coll *Collection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*InsertManyResult, error) {

	if len(documents) == 0 {
		return nil, ErrEmptySlice
	}

	result, err := coll.insert(ctx, documents, opts...)
	rr, err := processWriteError(err)
	if rr&rrMany == 0 {
		return nil, err
	}

	imResult := &InsertManyResult{InsertedIDs: result}
	var writeException WriteException
	if !errors.As(err, &writeException) {
		return imResult, err
	}

	
	bwErrors := make([]BulkWriteError, 0, len(writeException.WriteErrors))
	for _, we := range writeException.WriteErrors {
		bwErrors = append(bwErrors, BulkWriteError{
			WriteError: we,
			Request:    nil,
		})
	}

	return imResult, BulkWriteException{
		WriteErrors:       bwErrors,
		WriteConcernError: writeException.WriteConcernError,
		Labels:            writeException.Labels,
	}
}

func (coll *Collection) delete(ctx context.Context, filter interface{}, deleteOne bool, expectedRr returnResult,
	opts ...*options.DeleteOptions) (*DeleteResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	var limit int32
	if deleteOne {
		limit = 1
	}
	do := options.MergeDeleteOptions(opts...)
	didx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendDocumentElement(doc, "q", f)
	doc = bsoncore.AppendInt32Element(doc, "limit", limit)
	if do.Collation != nil {
		doc = bsoncore.AppendDocumentElement(doc, "collation", do.Collation.ToDocument())
	}
	if do.Hint != nil {
		if isUnorderedMap(do.Hint) {
			return nil, ErrMapForOrderedArgument{"hint"}
		}
		hint, err := marshalValue(do.Hint, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}

		doc = bsoncore.AppendValueElement(doc, "hint", hint)
	}
	doc, _ = bsoncore.AppendDocumentEnd(doc, didx)

	op := operation.NewDelete(doc).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).Ordered(true).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout).Logger(coll.client.logger).
		Authenticator(coll.client.authenticator)
	if do.Comment != nil {
		comment, err := marshalValue(do.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op = op.Comment(comment)
	}
	if do.Hint != nil {
		op = op.Hint(true)
	}
	if do.Let != nil {
		let, err := marshal(do.Let, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op = op.Let(let)
	}

	
	retryMode := driver.RetryNone
	if deleteOne && coll.client.retryWrites {
		retryMode = driver.RetryOncePerCommand
	}
	op = op.Retry(retryMode)
	rr, err := processWriteError(op.Execute(ctx))
	if rr&expectedRr == 0 {
		return nil, err
	}
	return &DeleteResult{DeletedCount: op.Result().N}, err
}











func (coll *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*DeleteResult, error) {

	return coll.delete(ctx, filter, true, rrOne, opts...)
}











func (coll *Collection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*DeleteResult, error) {

	return coll.delete(ctx, filter, false, rrMany, opts...)
}

func (coll *Collection) updateOrReplace(ctx context.Context, filter bsoncore.Document, update interface{}, multi bool,
	expectedRr returnResult, checkDollarKey bool, opts ...*options.UpdateOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	uo := options.MergeUpdateOptions(opts...)

	
	
	updateDoc, err := createUpdateDoc(
		filter,
		update,
		uo.Hint,
		uo.ArrayFilters,
		uo.Collation,
		uo.Upsert,
		multi,
		checkDollarKey,
		coll.bsonOpts,
		coll.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewUpdate(updateDoc).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).Hint(uo.Hint != nil).
		ArrayFilters(uo.ArrayFilters != nil).Ordered(true).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout).Logger(coll.client.logger).Authenticator(coll.client.authenticator)
	if uo.Let != nil {
		let, err := marshal(uo.Let, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op = op.Let(let)
	}

	if uo.BypassDocumentValidation != nil && *uo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*uo.BypassDocumentValidation)
	}
	if uo.Comment != nil {
		comment, err := marshalValue(uo.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op = op.Comment(comment)
	}
	if uo.BypassEmptyTsReplacement != nil {
		op.BypassEmptyTsReplacement(*uo.BypassEmptyTsReplacement)
	}
	retry := driver.RetryNone
	
	if !multi && coll.client.retryWrites {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)
	err = op.Execute(ctx)

	rr, err := processWriteError(err)
	if rr&expectedRr == 0 {
		return nil, err
	}

	opRes := op.Result()
	res := &UpdateResult{
		MatchedCount:  opRes.N,
		ModifiedCount: opRes.NModified,
		UpsertedCount: int64(len(opRes.Upserted)),
	}
	if len(opRes.Upserted) > 0 {
		res.UpsertedID = opRes.Upserted[0].ID
		res.MatchedCount--
	}

	return res, err
}














func (coll *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*UpdateResult, error) {
	if id == nil {
		return nil, ErrNilValue
	}
	return coll.UpdateOne(ctx, bson.D{{"_id", id}}, update, opts...)
}















func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	return coll.updateOrReplace(ctx, f, update, false, rrOne, true, opts...)
}














func (coll *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	return coll.updateOrReplace(ctx, f, update, true, rrMany, true, opts...)
}














func (coll *Collection) ReplaceOne(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	r, err := marshal(replacement, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	if err := ensureNoDollarKey(r); err != nil {
		return nil, err
	}

	updateOptions := make([]*options.UpdateOptions, 0, len(opts))
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		uOpts := options.Update()
		uOpts.BypassDocumentValidation = opt.BypassDocumentValidation
		uOpts.Collation = opt.Collation
		uOpts.Upsert = opt.Upsert
		uOpts.Hint = opt.Hint
		uOpts.Let = opt.Let
		uOpts.Comment = opt.Comment
		uOpts.BypassEmptyTsReplacement = opt.BypassEmptyTsReplacement
		updateOptions = append(updateOptions, uOpts)
	}

	return coll.updateOrReplace(ctx, f, r, false, rrOne, false, updateOptions...)
}












func (coll *Collection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*Cursor, error) {
	a := aggregateParams{
		ctx:            ctx,
		pipeline:       pipeline,
		client:         coll.client,
		registry:       coll.registry,
		readConcern:    coll.readConcern,
		writeConcern:   coll.writeConcern,
		bsonOpts:       coll.bsonOpts,
		retryRead:      coll.client.retryReads,
		db:             coll.db.name,
		col:            coll.name,
		readSelector:   coll.readSelector,
		writeSelector:  coll.writeSelector,
		readPreference: coll.readPreference,
		opts:           opts,
	}
	return aggregate(a)
}


func aggregate(a aggregateParams) (cur *Cursor, err error) {
	if a.ctx == nil {
		a.ctx = context.Background()
	}

	pipelineArr, hasOutputStage, err := marshalAggregatePipeline(a.pipeline, a.bsonOpts, a.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(a.ctx)
	
	defer func() {
		if err != nil && sess != nil {
			closeImplicitSession(sess)
		}
	}()
	if sess == nil && a.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(a.client.sessionPool, a.client.id)
	}
	if err = a.client.validSession(sess); err != nil {
		return nil, err
	}

	var wc *writeconcern.WriteConcern
	if hasOutputStage {
		wc = a.writeConcern
	}
	rc := a.readConcern
	if sess.TransactionRunning() {
		wc = nil
		rc = nil
	}
	if !writeconcern.AckWrite(wc) {
		closeImplicitSession(sess)
		sess = nil
	}

	selector := makeReadPrefSelector(sess, a.readSelector, a.client.localThreshold)
	if hasOutputStage {
		selector = makeOutputAggregateSelector(sess, a.readPreference, a.client.localThreshold)
	}

	ao := options.MergeAggregateOptions(a.opts...)

	cursorOpts := a.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(a.bsonOpts, a.registry)

	op := operation.NewAggregate(pipelineArr).
		Session(sess).
		WriteConcern(wc).
		ReadConcern(rc).
		ReadPreference(a.readPreference).
		CommandMonitor(a.client.monitor).
		ServerSelector(selector).
		ClusterClock(a.client.clock).
		Database(a.db).
		Collection(a.col).
		Deployment(a.client.deployment).
		Crypt(a.client.cryptFLE).
		ServerAPI(a.client.serverAPI).
		HasOutputStage(hasOutputStage).
		Timeout(a.client.timeout).
		MaxTime(ao.MaxTime).
		Authenticator(a.client.authenticator)

	
	
	
	
	
	
	_, deadlineSet := a.ctx.Deadline()
	op.OmitCSOTMaxTimeMS(deadlineSet)

	if ao.AllowDiskUse != nil {
		op.AllowDiskUse(*ao.AllowDiskUse)
	}
	
	if ao.BatchSize != nil && !(*ao.BatchSize == 0 && hasOutputStage) {
		op.BatchSize(*ao.BatchSize)
		cursorOpts.BatchSize = *ao.BatchSize
	}
	if ao.BypassDocumentValidation != nil && *ao.BypassDocumentValidation {
		op.BypassDocumentValidation(*ao.BypassDocumentValidation)
	}
	if ao.Collation != nil {
		op.Collation(bsoncore.Document(ao.Collation.ToDocument()))
	}
	if ao.MaxAwaitTime != nil {
		cursorOpts.MaxTimeMS = int64(*ao.MaxAwaitTime / time.Millisecond)
	}
	if ao.Comment != nil {
		op.Comment(*ao.Comment)

		commentVal, err := marshalValue(ao.Comment, a.bsonOpts, a.registry)
		if err != nil {
			return nil, err
		}
		cursorOpts.Comment = commentVal
	}
	if ao.Hint != nil {
		if isUnorderedMap(ao.Hint) {
			return nil, ErrMapForOrderedArgument{"hint"}
		}
		hintVal, err := marshalValue(ao.Hint, a.bsonOpts, a.registry)
		if err != nil {
			return nil, err
		}
		op.Hint(hintVal)
	}
	if ao.Let != nil {
		let, err := marshal(ao.Let, a.bsonOpts, a.registry)
		if err != nil {
			return nil, err
		}
		op.Let(let)
	}
	if ao.Custom != nil {
		
		
		customOptions := make(map[string]bsoncore.Value)
		for optionName, optionValue := range ao.Custom {
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(a.registry, optionValue)
			if err != nil {
				return nil, err
			}
			optionValueBSON := bsoncore.Value{Type: bsonType, Data: bsonData}
			customOptions[optionName] = optionValueBSON
		}
		op.CustomOptions(customOptions)
	}

	retry := driver.RetryNone
	if a.retryRead && !hasOutputStage {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(a.ctx)
	if err != nil {
		if wce, ok := err.(driver.WriteCommandError); ok && wce.WriteConcernError != nil {
			return nil, *convertDriverWriteConcernError(wce.WriteConcernError)
		}
		return nil, replaceErrors(err)
	}

	bc, err := op.Result(cursorOpts)
	if err != nil {
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, a.client.bsonOpts, a.registry, sess)
	return cursor, replaceErrors(err)
}









func (coll *Collection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	countOpts := options.MergeCountOptions(opts...)

	pipelineArr, err := countDocumentsAggregatePipeline(filter, coll.bsonOpts, coll.registry, countOpts)
	if err != nil {
		return 0, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}
	if err = coll.client.validSession(sess); err != nil {
		return 0, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewAggregate(pipelineArr).Session(sess).ReadConcern(rc).ReadPreference(coll.readPreference).
		CommandMonitor(coll.client.monitor).ServerSelector(selector).ClusterClock(coll.client.clock).Database(coll.db.name).
		Collection(coll.name).Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout).MaxTime(countOpts.MaxTime).Authenticator(coll.client.authenticator)
	if countOpts.Collation != nil {
		op.Collation(bsoncore.Document(countOpts.Collation.ToDocument()))
	}
	if countOpts.Comment != nil {
		op.Comment(*countOpts.Comment)
	}
	if countOpts.Hint != nil {
		if isUnorderedMap(countOpts.Hint) {
			return 0, ErrMapForOrderedArgument{"hint"}
		}
		hintVal, err := marshalValue(countOpts.Hint, coll.bsonOpts, coll.registry)
		if err != nil {
			return 0, err
		}
		op.Hint(hintVal)
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		return 0, replaceErrors(err)
	}

	batch := op.ResultCursorResponse().FirstBatch
	if batch == nil {
		return 0, errors.New("invalid response from server, no 'firstBatch' field")
	}

	docs, err := batch.Documents()
	if err != nil || len(docs) == 0 {
		return 0, nil
	}

	val, ok := docs[0].Lookup("n").AsInt64OK()
	if !ok {
		return 0, errors.New("invalid response from server, no 'n' field")
	}

	return val, nil
}








func (coll *Collection) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)

	var err error
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return 0, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	co := options.MergeEstimatedDocumentCountOptions(opts...)

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewCount().Session(sess).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).CommandMonitor(coll.client.monitor).
		Deployment(coll.client.deployment).ReadConcern(rc).ReadPreference(coll.readPreference).
		ServerSelector(selector).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout).MaxTime(co.MaxTime).Authenticator(coll.client.authenticator)

	if co.Comment != nil {
		comment, err := marshalValue(co.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return 0, err
		}
		op = op.Comment(comment)
	}

	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op.Retry(retry)

	err = op.Execute(ctx)
	return op.Result().N, replaceErrors(err)
}











func (coll *Collection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]interface{}, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)

	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	option := options.MergeDistinctOptions(opts...)

	op := operation.NewDistinct(fieldName, f).
		Session(sess).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).CommandMonitor(coll.client.monitor).
		Deployment(coll.client.deployment).ReadConcern(rc).ReadPreference(coll.readPreference).
		ServerSelector(selector).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout).MaxTime(option.MaxTime).Authenticator(coll.client.authenticator)

	if option.Collation != nil {
		op.Collation(bsoncore.Document(option.Collation.ToDocument()))
	}
	if option.Comment != nil {
		comment, err := marshalValue(option.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Comment(comment)
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		return nil, replaceErrors(err)
	}

	arr, ok := op.Result().Values.ArrayOK()
	if !ok {
		return nil, fmt.Errorf("response field 'values' is type array, but received BSON type %s", op.Result().Values.Type)
	}

	values, err := arr.Values()
	if err != nil {
		return nil, err
	}

	retArray := make([]interface{}, len(values))

	for i, val := range values {
		raw := bson.RawValue{Type: val.Type, Value: val.Data}
		err = raw.Unmarshal(&retArray[i])
		if err != nil {
			return nil, err
		}
	}

	return retArray, replaceErrors(err)
}









func (coll *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *Cursor, err error) {

	if ctx == nil {
		ctx = context.Background()
	}

	
	
	
	
	
	
	_, deadlineSet := ctx.Deadline()
	return coll.find(ctx, filter, deadlineSet, opts...)
}

func (coll *Collection) find(
	ctx context.Context,
	filter interface{},
	omitCSOTMaxTimeMS bool,
	opts ...*options.FindOptions,
) (cur *Cursor, err error) {

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	
	defer func() {
		if err != nil && sess != nil {
			closeImplicitSession(sess)
		}
	}()
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	fo := options.MergeFindOptions(opts...)

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewFind(f).
		Session(sess).ReadConcern(rc).ReadPreference(coll.readPreference).
		CommandMonitor(coll.client.monitor).ServerSelector(selector).
		ClusterClock(coll.client.clock).Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).ServerAPI(coll.client.serverAPI).
		Timeout(coll.client.timeout).MaxTime(fo.MaxTime).Logger(coll.client.logger).
		OmitCSOTMaxTimeMS(omitCSOTMaxTimeMS).Authenticator(coll.client.authenticator)

	cursorOpts := coll.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(coll.bsonOpts, coll.registry)

	if fo.AllowDiskUse != nil {
		op.AllowDiskUse(*fo.AllowDiskUse)
	}
	if fo.AllowPartialResults != nil {
		op.AllowPartialResults(*fo.AllowPartialResults)
	}
	if fo.BatchSize != nil {
		cursorOpts.BatchSize = *fo.BatchSize
		op.BatchSize(*fo.BatchSize)
	}
	if fo.Collation != nil {
		op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	}
	if fo.Comment != nil {
		op.Comment(*fo.Comment)

		commentVal, err := marshalValue(fo.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		cursorOpts.Comment = commentVal
	}
	if fo.CursorType != nil {
		switch *fo.CursorType {
		case options.Tailable:
			op.Tailable(true)
		case options.TailableAwait:
			op.Tailable(true)
			op.AwaitData(true)
		}
	}
	if fo.Hint != nil {
		if isUnorderedMap(fo.Hint) {
			return nil, ErrMapForOrderedArgument{"hint"}
		}
		hint, err := marshalValue(fo.Hint, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Hint(hint)
	}
	if fo.Let != nil {
		let, err := marshal(fo.Let, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Let(let)
	}
	if fo.Limit != nil {
		limit := *fo.Limit
		if limit < 0 {
			limit = -1 * limit
			op.SingleBatch(true)
		}
		cursorOpts.Limit = int32(limit)
		op.Limit(limit)
	}
	if fo.Max != nil {
		max, err := marshal(fo.Max, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Max(max)
	}
	if fo.MaxAwaitTime != nil {
		cursorOpts.MaxTimeMS = int64(*fo.MaxAwaitTime / time.Millisecond)
	}
	if fo.Min != nil {
		min, err := marshal(fo.Min, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Min(min)
	}
	if fo.NoCursorTimeout != nil {
		op.NoCursorTimeout(*fo.NoCursorTimeout)
	}
	if fo.OplogReplay != nil {
		op.OplogReplay(*fo.OplogReplay)
	}
	if fo.Projection != nil {
		proj, err := marshal(fo.Projection, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Projection(proj)
	}
	if fo.ReturnKey != nil {
		op.ReturnKey(*fo.ReturnKey)
	}
	if fo.ShowRecordID != nil {
		op.ShowRecordID(*fo.ShowRecordID)
	}
	if fo.Skip != nil {
		op.Skip(*fo.Skip)
	}
	if fo.Snapshot != nil {
		op.Snapshot(*fo.Snapshot)
	}
	if fo.Sort != nil {
		if isUnorderedMap(fo.Sort) {
			return nil, ErrMapForOrderedArgument{"sort"}
		}
		sort, err := marshal(fo.Sort, coll.bsonOpts, coll.registry)
		if err != nil {
			return nil, err
		}
		op.Sort(sort)
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	if err = op.Execute(ctx); err != nil {
		return nil, replaceErrors(err)
	}

	bc, err := op.Result(cursorOpts)
	if err != nil {
		return nil, replaceErrors(err)
	}
	return newCursorWithSession(bc, coll.bsonOpts, coll.registry, sess)
}










func (coll *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *SingleResult {

	if ctx == nil {
		ctx = context.Background()
	}

	findOpts := make([]*options.FindOptions, 0, len(opts))
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		findOpts = append(findOpts, &options.FindOptions{
			AllowPartialResults: opt.AllowPartialResults,
			BatchSize:           opt.BatchSize,
			Collation:           opt.Collation,
			Comment:             opt.Comment,
			CursorType:          opt.CursorType,
			Hint:                opt.Hint,
			Max:                 opt.Max,
			MaxAwaitTime:        opt.MaxAwaitTime,
			MaxTime:             opt.MaxTime,
			Min:                 opt.Min,
			NoCursorTimeout:     opt.NoCursorTimeout,
			OplogReplay:         opt.OplogReplay,
			Projection:          opt.Projection,
			ReturnKey:           opt.ReturnKey,
			ShowRecordID:        opt.ShowRecordID,
			Skip:                opt.Skip,
			Snapshot:            opt.Snapshot,
			Sort:                opt.Sort,
		})
	}
	
	
	findOpts = append(findOpts, options.Find().SetLimit(-1))

	cursor, err := coll.find(ctx, filter, false, findOpts...)
	return &SingleResult{
		ctx:      ctx,
		cur:      cursor,
		bsonOpts: coll.bsonOpts,
		reg:      coll.registry,
		err:      replaceErrors(err),
	}
}

func (coll *Collection) findAndModify(ctx context.Context, op *operation.FindAndModify) *SingleResult {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	var err error
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return &SingleResult{err: err}
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	retry := driver.RetryNone
	if coll.client.retryWrites {
		retry = driver.RetryOnce
	}

	op = op.Session(sess).
		WriteConcern(wc).
		CommandMonitor(coll.client.monitor).
		ServerSelector(selector).
		ClusterClock(coll.client.clock).
		Database(coll.db.name).
		Collection(coll.name).
		Deployment(coll.client.deployment).
		Retry(retry).
		Crypt(coll.client.cryptFLE)

	_, err = processWriteError(op.Execute(ctx))
	if err != nil {
		return &SingleResult{err: err}
	}

	return &SingleResult{
		ctx:      ctx,
		rdr:      bson.Raw(op.Result().Value),
		bsonOpts: coll.bsonOpts,
		reg:      coll.registry,
	}
}












func (coll *Collection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *SingleResult {

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return &SingleResult{err: err}
	}
	fod := options.MergeFindOneAndDeleteOptions(opts...)
	op := operation.NewFindAndModify(f).Remove(true).ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout).
		MaxTime(fod.MaxTime).Authenticator(coll.client.authenticator)
	if fod.Collation != nil {
		op = op.Collation(bsoncore.Document(fod.Collation.ToDocument()))
	}
	if fod.Comment != nil {
		comment, err := marshalValue(fod.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Comment(comment)
	}
	if fod.Projection != nil {
		proj, err := marshal(fod.Projection, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Fields(proj)
	}
	if fod.Sort != nil {
		if isUnorderedMap(fod.Sort) {
			return &SingleResult{err: ErrMapForOrderedArgument{"sort"}}
		}
		sort, err := marshal(fod.Sort, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Sort(sort)
	}
	if fod.Hint != nil {
		if isUnorderedMap(fod.Hint) {
			return &SingleResult{err: ErrMapForOrderedArgument{"hint"}}
		}
		hint, err := marshalValue(fod.Hint, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Hint(hint)
	}
	if fod.Let != nil {
		let, err := marshal(fod.Let, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Let(let)
	}

	return coll.findAndModify(ctx, op)
}















func (coll *Collection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *SingleResult {

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return &SingleResult{err: err}
	}
	r, err := marshal(replacement, coll.bsonOpts, coll.registry)
	if err != nil {
		return &SingleResult{err: err}
	}
	if firstElem, err := r.IndexErr(0); err == nil && strings.HasPrefix(firstElem.Key(), "$") {
		return &SingleResult{err: errors.New("replacement document cannot contain keys beginning with '$'")}
	}

	fo := options.MergeFindOneAndReplaceOptions(opts...)
	op := operation.NewFindAndModify(f).Update(bsoncore.Value{Type: bsontype.EmbeddedDocument, Data: r}).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout).MaxTime(fo.MaxTime).Authenticator(coll.client.authenticator)

	if fo.BypassDocumentValidation != nil && *fo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*fo.BypassDocumentValidation)
	}
	if fo.Collation != nil {
		op = op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	}
	if fo.Comment != nil {
		comment, err := marshalValue(fo.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Comment(comment)
	}
	if fo.Projection != nil {
		proj, err := marshal(fo.Projection, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Fields(proj)
	}
	if fo.ReturnDocument != nil {
		op = op.NewDocument(*fo.ReturnDocument == options.After)
	}
	if fo.Sort != nil {
		if isUnorderedMap(fo.Sort) {
			return &SingleResult{err: ErrMapForOrderedArgument{"sort"}}
		}
		sort, err := marshal(fo.Sort, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Sort(sort)
	}
	if fo.Upsert != nil {
		op = op.Upsert(*fo.Upsert)
	}
	if fo.Hint != nil {
		if isUnorderedMap(fo.Hint) {
			return &SingleResult{err: ErrMapForOrderedArgument{"hint"}}
		}
		hint, err := marshalValue(fo.Hint, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Hint(hint)
	}
	if fo.Let != nil {
		let, err := marshal(fo.Let, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Let(let)
	}
	if fo.BypassEmptyTsReplacement != nil {
		op = op.BypassEmptyTsReplacement(*fo.BypassEmptyTsReplacement)
	}

	return coll.findAndModify(ctx, op)
}
















func (coll *Collection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *SingleResult {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := marshal(filter, coll.bsonOpts, coll.registry)
	if err != nil {
		return &SingleResult{err: err}
	}

	fo := options.MergeFindOneAndUpdateOptions(opts...)
	op := operation.NewFindAndModify(f).ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout).
		MaxTime(fo.MaxTime).Authenticator(coll.client.authenticator)

	u, err := marshalUpdateValue(update, coll.bsonOpts, coll.registry, true)
	if err != nil {
		return &SingleResult{err: err}
	}
	op = op.Update(u)

	if fo.ArrayFilters != nil {
		af := fo.ArrayFilters
		reg := coll.registry
		if af.Registry != nil {
			reg = af.Registry
		}
		filtersDoc, err := marshalValue(af.Filters, coll.bsonOpts, reg)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.ArrayFilters(filtersDoc.Data)
	}
	if fo.BypassDocumentValidation != nil && *fo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*fo.BypassDocumentValidation)
	}
	if fo.Collation != nil {
		op = op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	}
	if fo.Comment != nil {
		comment, err := marshalValue(fo.Comment, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Comment(comment)
	}
	if fo.Projection != nil {
		proj, err := marshal(fo.Projection, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Fields(proj)
	}
	if fo.ReturnDocument != nil {
		op = op.NewDocument(*fo.ReturnDocument == options.After)
	}
	if fo.Sort != nil {
		if isUnorderedMap(fo.Sort) {
			return &SingleResult{err: ErrMapForOrderedArgument{"sort"}}
		}
		sort, err := marshal(fo.Sort, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Sort(sort)
	}
	if fo.Upsert != nil {
		op = op.Upsert(*fo.Upsert)
	}
	if fo.Hint != nil {
		if isUnorderedMap(fo.Hint) {
			return &SingleResult{err: ErrMapForOrderedArgument{"hint"}}
		}
		hint, err := marshalValue(fo.Hint, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Hint(hint)
	}
	if fo.Let != nil {
		let, err := marshal(fo.Let, coll.bsonOpts, coll.registry)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Let(let)
	}
	if fo.BypassEmptyTsReplacement != nil {
		op = op.BypassEmptyTsReplacement(*fo.BypassEmptyTsReplacement)
	}

	return coll.findAndModify(ctx, op)
}














func (coll *Collection) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {

	csConfig := changeStreamConfig{
		readConcern:    coll.readConcern,
		readPreference: coll.readPreference,
		client:         coll.client,
		bsonOpts:       coll.bsonOpts,
		registry:       coll.registry,
		streamType:     CollectionStream,
		collectionName: coll.Name(),
		databaseName:   coll.db.Name(),
		crypt:          coll.client.cryptFLE,
	}
	return newChangeStream(ctx, csConfig, pipeline, opts...)
}


func (coll *Collection) Indexes() IndexView {
	return IndexView{coll: coll}
}


func (coll *Collection) SearchIndexes() SearchIndexView {
	c, _ := coll.Clone() 
	c.readConcern = nil
	c.writeConcern = nil
	return SearchIndexView{
		coll: c,
	}
}



func (coll *Collection) Drop(ctx context.Context) error {
	
	
	
	
	ef := coll.db.getEncryptedFieldsFromMap(coll.name)
	if ef == nil && coll.db.client.encryptedFieldsMap != nil {
		var err error
		if ef, err = coll.db.getEncryptedFieldsFromServer(ctx, coll.name); err != nil {
			return err
		}
	}

	if ef != nil {
		return coll.dropEncryptedCollection(ctx, ef)
	}

	return coll.drop(ctx)
}


func (coll *Collection) dropEncryptedCollection(ctx context.Context, ef interface{}) error {
	efBSON, err := marshal(ef, coll.bsonOpts, coll.registry)
	if err != nil {
		return fmt.Errorf("error transforming document: %w", err)
	}

	
	
	escCollection, err := csfle.GetEncryptedStateCollectionName(efBSON, coll.name, csfle.EncryptedStateCollection)
	if err != nil {
		return err
	}
	if err := coll.db.Collection(escCollection).drop(ctx); err != nil {
		return err
	}

	
	ecocCollection, err := csfle.GetEncryptedStateCollectionName(efBSON, coll.name, csfle.EncryptedCompactionCollection)
	if err != nil {
		return err
	}
	if err := coll.db.Collection(ecocCollection).drop(ctx); err != nil {
		return err
	}

	
	return coll.drop(ctx)
}


func (coll *Collection) drop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(coll.client.sessionPool, coll.client.id)
		defer sess.EndSession()
	}

	err := coll.client.validSession(sess)
	if err != nil {
		return err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewDropCollection().
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.cryptFLE).
		ServerAPI(coll.client.serverAPI).Timeout(coll.client.timeout).
		Authenticator(coll.client.authenticator)
	err = op.Execute(ctx)

	
	driverErr, ok := err.(driver.Error)
	if !ok || (ok && !driverErr.NamespaceNotFound()) {
		return replaceErrors(err)
	}
	return nil
}

type pinnedServerSelector struct {
	stringer fmt.Stringer
	fallback description.ServerSelector
	session  *session.Client
}

func (pss pinnedServerSelector) String() string {
	if pss.stringer == nil {
		return ""
	}

	return pss.stringer.String()
}

func (pss pinnedServerSelector) SelectServer(
	t description.Topology,
	svrs []description.Server,
) ([]description.Server, error) {
	if pss.session != nil && pss.session.PinnedServer != nil {
		
		for _, candidate := range svrs {
			if candidate.Addr == pss.session.PinnedServer.Addr {
				return []description.Server{candidate}, nil
			}
		}

		return nil, nil
	}

	return pss.fallback.SelectServer(t, svrs)
}

func makePinnedSelector(sess *session.Client, fallback description.ServerSelector) description.ServerSelector {
	pss := pinnedServerSelector{
		session:  sess,
		fallback: fallback,
	}

	if srvSelectorStringer, ok := fallback.(fmt.Stringer); ok {
		pss.stringer = srvSelectorStringer
	}

	return pss
}

func makeReadPrefSelector(sess *session.Client, selector description.ServerSelector, localThreshold time.Duration) description.ServerSelector {
	if sess != nil && sess.TransactionRunning() {
		selector = description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(sess.CurrentRp),
			description.LatencySelector(localThreshold),
		})
	}

	return makePinnedSelector(sess, selector)
}

func makeOutputAggregateSelector(sess *session.Client, rp *readpref.ReadPref, localThreshold time.Duration) description.ServerSelector {
	if sess != nil && sess.TransactionRunning() {
		
		rp = sess.CurrentRp
	}

	selector := description.CompositeSelector([]description.ServerSelector{
		description.OutputAggregateSelector(rp),
		description.LatencySelector(localThreshold),
	})
	return makePinnedSelector(sess, selector)
}




func isUnorderedMap(val interface{}) bool {
	refValue := reflect.ValueOf(val)
	return refValue.Kind() == reflect.Map && refValue.Len() > 1
}
