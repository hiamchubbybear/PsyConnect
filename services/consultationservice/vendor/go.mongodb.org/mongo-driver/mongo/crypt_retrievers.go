





package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type keyRetriever struct {
	coll *Collection
}

func (kr *keyRetriever) cryptKeys(ctx context.Context, filter bsoncore.Document) ([]bsoncore.Document, error) {
	
	
	ctx = NewSessionContext(ctx, nil)
	cursor, err := kr.coll.Find(ctx, filter)
	if err != nil {
		return nil, EncryptionKeyVaultError{Wrapped: err}
	}
	defer cursor.Close(ctx)

	var results []bsoncore.Document
	for cursor.Next(ctx) {
		cur := make([]byte, len(cursor.Current))
		copy(cur, cursor.Current)
		results = append(results, cur)
	}
	if err = cursor.Err(); err != nil {
		return nil, EncryptionKeyVaultError{Wrapped: err}
	}

	return results, nil
}


type collInfoRetriever struct {
	client *Client
}

func (cir *collInfoRetriever) cryptCollInfo(ctx context.Context, db string, filter bsoncore.Document) (bsoncore.Document, error) {
	
	
	ctx = NewSessionContext(ctx, nil)
	cursor, err := cir.client.Database(db).ListCollections(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, cursor.Err()
	}

	res := make([]byte, len(cursor.Current))
	copy(res, cursor.Current)
	return res, nil
}
