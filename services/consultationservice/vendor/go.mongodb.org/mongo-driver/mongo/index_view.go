





package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)



var ErrInvalidIndexValue = errors.New("invalid index value")


var ErrNonStringIndexName = errors.New("index name must be a string")


var ErrMultipleIndexDrop = errors.New("multiple indexes would be dropped")



type IndexView struct {
	coll *Collection
}


type IndexModel struct {
	
	
	
	Keys interface{}

	
	Options *options.IndexOptions
}

func isNamespaceNotFoundError(err error) bool {
	if de, ok := err.(driver.Error); ok {
		return de.Code == 26
	}
	return false
}







func (iv IndexView) List(ctx context.Context, opts ...*options.ListIndexesOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && iv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(iv.coll.client.sessionPool, iv.coll.client.id)
	}

	err := iv.coll.client.validSession(sess)
	if err != nil {
		closeImplicitSession(sess)
		return nil, err
	}

	selector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(iv.coll.client.localThreshold),
	})
	selector = makeReadPrefSelector(sess, selector, iv.coll.client.localThreshold)

	
	
	op := operation.NewListIndexes().
		Session(sess).CommandMonitor(iv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(iv.coll.client.clock).
		Database(iv.coll.db.name).Collection(iv.coll.name).
		Deployment(iv.coll.client.deployment).ServerAPI(iv.coll.client.serverAPI).
		Timeout(iv.coll.client.timeout).Authenticator(iv.coll.client.authenticator)

	cursorOpts := iv.coll.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(iv.coll.bsonOpts, iv.coll.registry)

	lio := options.MergeListIndexesOptions(opts...)
	if lio.BatchSize != nil {
		op = op.BatchSize(*lio.BatchSize)
		cursorOpts.BatchSize = *lio.BatchSize
	}
	op = op.MaxTime(lio.MaxTime)
	retry := driver.RetryNone
	if iv.coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		
		closeImplicitSession(sess)
		if isNamespaceNotFoundError(err) {
			return newEmptyCursor(), nil
		}

		return nil, replaceErrors(err)
	}

	bc, err := op.Result(cursorOpts)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, iv.coll.bsonOpts, iv.coll.registry, sess)
	return cursor, replaceErrors(err)
}


func (iv IndexView) ListSpecifications(ctx context.Context, opts ...*options.ListIndexesOptions) ([]*IndexSpecification, error) {
	cursor, err := iv.List(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var results []*IndexSpecification
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	ns := iv.coll.db.Name() + "." + iv.coll.Name()
	for _, res := range results {
		
		
		res.Namespace = ns
	}

	return results, nil
}



func (iv IndexView) CreateOne(ctx context.Context, model IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	names, err := iv.CreateMany(ctx, []IndexModel{model}, opts...)
	if err != nil {
		return "", err
	}

	return names[0], nil
}











func (iv IndexView) CreateMany(ctx context.Context, models []IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	names := make([]string, 0, len(models))

	var indexes bsoncore.Document
	aidx, indexes := bsoncore.AppendArrayStart(indexes)

	for i, model := range models {
		if model.Keys == nil {
			return nil, fmt.Errorf("index model keys cannot be nil")
		}

		if isUnorderedMap(model.Keys) {
			return nil, ErrMapForOrderedArgument{"keys"}
		}

		keys, err := marshal(model.Keys, iv.coll.bsonOpts, iv.coll.registry)
		if err != nil {
			return nil, err
		}

		name, err := getOrGenerateIndexName(keys, model)
		if err != nil {
			return nil, err
		}

		names = append(names, name)

		var iidx int32
		iidx, indexes = bsoncore.AppendDocumentElementStart(indexes, strconv.Itoa(i))
		indexes = bsoncore.AppendDocumentElement(indexes, "key", keys)

		if model.Options == nil {
			model.Options = options.Index()
		}
		model.Options.SetName(name)

		optsDoc, err := iv.createOptionsDoc(model.Options)
		if err != nil {
			return nil, err
		}

		indexes = bsoncore.AppendDocument(indexes, optsDoc)

		indexes, err = bsoncore.AppendDocumentEnd(indexes, iidx)
		if err != nil {
			return nil, err
		}
	}

	indexes, err := bsoncore.AppendArrayEnd(indexes, aidx)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)

	if sess == nil && iv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(iv.coll.client.sessionPool, iv.coll.client.id)
		defer sess.EndSession()
	}

	err = iv.coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := iv.coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, iv.coll.writeSelector)

	option := options.MergeCreateIndexesOptions(opts...)

	
	
	
	
	op := operation.NewCreateIndexes(indexes).
		Session(sess).WriteConcern(wc).ClusterClock(iv.coll.client.clock).
		Database(iv.coll.db.name).Collection(iv.coll.name).CommandMonitor(iv.coll.client.monitor).
		Deployment(iv.coll.client.deployment).ServerSelector(selector).ServerAPI(iv.coll.client.serverAPI).
		Timeout(iv.coll.client.timeout).MaxTime(option.MaxTime).Authenticator(iv.coll.client.authenticator)
	if option.CommitQuorum != nil {
		commitQuorum, err := marshalValue(option.CommitQuorum, iv.coll.bsonOpts, iv.coll.registry)
		if err != nil {
			return nil, err
		}

		op.CommitQuorum(commitQuorum)
	}

	err = op.Execute(ctx)
	if err != nil {
		_, err = processWriteError(err)
		return nil, err
	}

	return names, nil
}

func (iv IndexView) createOptionsDoc(opts *options.IndexOptions) (bsoncore.Document, error) {
	optsDoc := bsoncore.Document{}
	if opts.Background != nil {
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "background", *opts.Background)
	}
	if opts.ExpireAfterSeconds != nil {
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "expireAfterSeconds", *opts.ExpireAfterSeconds)
	}
	if opts.Name != nil {
		optsDoc = bsoncore.AppendStringElement(optsDoc, "name", *opts.Name)
	}
	if opts.Sparse != nil {
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "sparse", *opts.Sparse)
	}
	if opts.StorageEngine != nil {
		doc, err := marshal(opts.StorageEngine, iv.coll.bsonOpts, iv.coll.registry)
		if err != nil {
			return nil, err
		}

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "storageEngine", doc)
	}
	if opts.Unique != nil {
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "unique", *opts.Unique)
	}
	if opts.Version != nil {
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "v", *opts.Version)
	}
	if opts.DefaultLanguage != nil {
		optsDoc = bsoncore.AppendStringElement(optsDoc, "default_language", *opts.DefaultLanguage)
	}
	if opts.LanguageOverride != nil {
		optsDoc = bsoncore.AppendStringElement(optsDoc, "language_override", *opts.LanguageOverride)
	}
	if opts.TextVersion != nil {
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "textIndexVersion", *opts.TextVersion)
	}
	if opts.Weights != nil {
		doc, err := marshal(opts.Weights, iv.coll.bsonOpts, iv.coll.registry)
		if err != nil {
			return nil, err
		}

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "weights", doc)
	}
	if opts.SphereVersion != nil {
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "2dsphereIndexVersion", *opts.SphereVersion)
	}
	if opts.Bits != nil {
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "bits", *opts.Bits)
	}
	if opts.Max != nil {
		optsDoc = bsoncore.AppendDoubleElement(optsDoc, "max", *opts.Max)
	}
	if opts.Min != nil {
		optsDoc = bsoncore.AppendDoubleElement(optsDoc, "min", *opts.Min)
	}
	if opts.BucketSize != nil {
		optsDoc = bsoncore.AppendInt32Element(optsDoc, "bucketSize", *opts.BucketSize)
	}
	if opts.PartialFilterExpression != nil {
		doc, err := marshal(opts.PartialFilterExpression, iv.coll.bsonOpts, iv.coll.registry)
		if err != nil {
			return nil, err
		}

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "partialFilterExpression", doc)
	}
	if opts.Collation != nil {
		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "collation", bsoncore.Document(opts.Collation.ToDocument()))
	}
	if opts.WildcardProjection != nil {
		doc, err := marshal(opts.WildcardProjection, iv.coll.bsonOpts, iv.coll.registry)
		if err != nil {
			return nil, err
		}

		optsDoc = bsoncore.AppendDocumentElement(optsDoc, "wildcardProjection", doc)
	}
	if opts.Hidden != nil {
		optsDoc = bsoncore.AppendBooleanElement(optsDoc, "hidden", *opts.Hidden)
	}

	return optsDoc, nil
}

func (iv IndexView) drop(ctx context.Context, index any, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && iv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(iv.coll.client.sessionPool, iv.coll.client.id)
		defer sess.EndSession()
	}

	err := iv.coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := iv.coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, iv.coll.writeSelector)

	dio := options.MergeDropIndexesOptions(opts...)

	
	
	op := operation.NewDropIndexes(index).Session(sess).WriteConcern(wc).CommandMonitor(iv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(iv.coll.client.clock).
		Database(iv.coll.db.name).Collection(iv.coll.name).
		Deployment(iv.coll.client.deployment).ServerAPI(iv.coll.client.serverAPI).
		Timeout(iv.coll.client.timeout).MaxTime(dio.MaxTime).
		Authenticator(iv.coll.client.authenticator)

	err = op.Execute(ctx)
	if err != nil {
		return nil, replaceErrors(err)
	}

	
	ridx, res := bsoncore.AppendDocumentStart(nil)
	res = bsoncore.AppendInt32Element(res, "nIndexesWas", op.Result().NIndexesWas)
	res, _ = bsoncore.AppendDocumentEnd(res, ridx)
	return res, nil
}












func (iv IndexView) DropOne(ctx context.Context, name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	if name == "*" {
		return nil, ErrMultipleIndexDrop
	}

	return iv.drop(ctx, name, opts...)
}






func (iv IndexView) DropOneWithKey(ctx context.Context, keySpecDocument interface{}, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	doc, err := marshal(keySpecDocument, iv.coll.bsonOpts, iv.coll.registry)
	if err != nil {
		return nil, err
	}

	return iv.drop(ctx, doc, opts...)
}









func (iv IndexView) DropAll(ctx context.Context, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	return iv.drop(ctx, "*", opts...)
}

func getOrGenerateIndexName(keySpecDocument bsoncore.Document, model IndexModel) (string, error) {
	if model.Options != nil && model.Options.Name != nil {
		return *model.Options.Name, nil
	}

	name := bytes.NewBufferString("")
	first := true

	elems, err := keySpecDocument.Elements()
	if err != nil {
		return "", err
	}
	for _, elem := range elems {
		if !first {
			_, err := name.WriteRune('_')
			if err != nil {
				return "", err
			}
		}

		_, err := name.WriteString(elem.Key())
		if err != nil {
			return "", err
		}

		_, err = name.WriteRune('_')
		if err != nil {
			return "", err
		}

		var value string

		bsonValue := elem.Value()
		switch bsonValue.Type {
		case bsontype.Int32:
			value = fmt.Sprintf("%d", bsonValue.Int32())
		case bsontype.Int64:
			value = fmt.Sprintf("%d", bsonValue.Int64())
		case bsontype.String:
			value = bsonValue.StringValue()
		default:
			return "", ErrInvalidIndexValue
		}

		_, err = name.WriteString(value)
		if err != nil {
			return "", err
		}

		first = false
	}

	return name.String(), nil
}
