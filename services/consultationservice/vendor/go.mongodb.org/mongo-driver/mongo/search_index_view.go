





package mongo

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)



type SearchIndexView struct {
	coll *Collection
}


type SearchIndexModel struct {
	
	
	Definition interface{}

	
	Options *options.SearchIndexesOptions
}







func (siv SearchIndexView) List(
	ctx context.Context,
	searchIdxOpts *options.SearchIndexesOptions,
	opts ...*options.ListSearchIndexesOptions,
) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	index := bson.D{}
	if searchIdxOpts != nil && searchIdxOpts.Name != nil {
		index = bson.D{{"name", *searchIdxOpts.Name}}
	}

	aggregateOpts := make([]*options.AggregateOptions, len(opts))
	for i, opt := range opts {
		aggregateOpts[i] = opt.AggregateOpts
	}

	return siv.coll.Aggregate(ctx, Pipeline{{{"$listSearchIndexes", index}}}, aggregateOpts...)
}



func (siv SearchIndexView) CreateOne(
	ctx context.Context,
	model SearchIndexModel,
	opts ...*options.CreateSearchIndexesOptions,
) (string, error) {
	names, err := siv.CreateMany(ctx, []SearchIndexModel{model}, opts...)
	if err != nil {
		return "", err
	}

	return names[0], nil
}








func (siv SearchIndexView) CreateMany(
	ctx context.Context,
	models []SearchIndexModel,
	_ ...*options.CreateSearchIndexesOptions,
) ([]string, error) {
	var indexes bsoncore.Document
	aidx, indexes := bsoncore.AppendArrayStart(indexes)

	for i, model := range models {
		if model.Definition == nil {
			return nil, fmt.Errorf("search index model definition cannot be nil")
		}

		definition, err := marshal(model.Definition, siv.coll.bsonOpts, siv.coll.registry)
		if err != nil {
			return nil, err
		}

		var iidx int32
		iidx, indexes = bsoncore.AppendDocumentElementStart(indexes, strconv.Itoa(i))
		if model.Options != nil && model.Options.Name != nil {
			indexes = bsoncore.AppendStringElement(indexes, "name", *model.Options.Name)
		}
		if model.Options != nil && model.Options.Type != nil {
			indexes = bsoncore.AppendStringElement(indexes, "type", *model.Options.Type)
		}
		indexes = bsoncore.AppendDocumentElement(indexes, "definition", definition)

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

	if sess == nil && siv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(siv.coll.client.sessionPool, siv.coll.client.id)
		defer sess.EndSession()
	}

	err = siv.coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	selector := makePinnedSelector(sess, siv.coll.writeSelector)

	op := operation.NewCreateSearchIndexes(indexes).
		Session(sess).CommandMonitor(siv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(siv.coll.client.clock).
		Collection(siv.coll.name).Database(siv.coll.db.name).
		Deployment(siv.coll.client.deployment).ServerAPI(siv.coll.client.serverAPI).
		Timeout(siv.coll.client.timeout).Authenticator(siv.coll.client.authenticator)

	err = op.Execute(ctx)
	if err != nil {
		_, err = processWriteError(err)
		return nil, err
	}

	indexesCreated := op.Result().IndexesCreated
	names := make([]string, 0, len(indexesCreated))
	for _, index := range indexesCreated {
		names = append(names, index.Name)
	}

	return names, nil
}








func (siv SearchIndexView) DropOne(
	ctx context.Context,
	name string,
	_ ...*options.DropSearchIndexOptions,
) error {
	if name == "*" {
		return ErrMultipleIndexDrop
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && siv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(siv.coll.client.sessionPool, siv.coll.client.id)
		defer sess.EndSession()
	}

	err := siv.coll.client.validSession(sess)
	if err != nil {
		return err
	}

	selector := makePinnedSelector(sess, siv.coll.writeSelector)

	op := operation.NewDropSearchIndex(name).
		Session(sess).CommandMonitor(siv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(siv.coll.client.clock).
		Collection(siv.coll.name).Database(siv.coll.db.name).
		Deployment(siv.coll.client.deployment).ServerAPI(siv.coll.client.serverAPI).
		Timeout(siv.coll.client.timeout).Authenticator(siv.coll.client.authenticator)

	err = op.Execute(ctx)
	if de, ok := err.(driver.Error); ok && de.NamespaceNotFound() {
		return nil
	}
	return err
}









func (siv SearchIndexView) UpdateOne(
	ctx context.Context,
	name string,
	definition interface{},
	_ ...*options.UpdateSearchIndexOptions,
) error {
	if definition == nil {
		return fmt.Errorf("search index definition cannot be nil")
	}

	indexDefinition, err := marshal(definition, siv.coll.bsonOpts, siv.coll.registry)
	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && siv.coll.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(siv.coll.client.sessionPool, siv.coll.client.id)
		defer sess.EndSession()
	}

	err = siv.coll.client.validSession(sess)
	if err != nil {
		return err
	}

	selector := makePinnedSelector(sess, siv.coll.writeSelector)

	op := operation.NewUpdateSearchIndex(name, indexDefinition).
		Session(sess).CommandMonitor(siv.coll.client.monitor).
		ServerSelector(selector).ClusterClock(siv.coll.client.clock).
		Collection(siv.coll.name).Database(siv.coll.db.name).
		Deployment(siv.coll.client.deployment).ServerAPI(siv.coll.client.serverAPI).
		Timeout(siv.coll.client.timeout).Authenticator(siv.coll.client.authenticator)

	return op.Execute(ctx)
}
