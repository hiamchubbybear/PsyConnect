





package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)


type BulkWriteResult struct {
	
	InsertedCount int64

	
	MatchedCount int64

	
	ModifiedCount int64

	
	DeletedCount int64

	
	UpsertedCount int64

	
	UpsertedIDs map[int64]interface{}
}


type InsertOneResult struct {
	
	InsertedID interface{}
}


type InsertManyResult struct {
	
	InsertedIDs []interface{}
}




type DeleteResult struct {
	DeletedCount int64 `bson:"n"` 
}



type RewrapManyDataKeyResult struct {
	*BulkWriteResult
}


type ListDatabasesResult struct {
	
	Databases []DatabaseSpecification

	
	
	TotalSize int64
}

func newListDatabasesResultFromOperation(res operation.ListDatabasesResult) ListDatabasesResult {
	var ldr ListDatabasesResult
	ldr.Databases = make([]DatabaseSpecification, 0, len(res.Databases))
	for _, spec := range res.Databases {
		ldr.Databases = append(
			ldr.Databases,
			DatabaseSpecification{Name: spec.Name, SizeOnDisk: spec.SizeOnDisk, Empty: spec.Empty},
		)
	}
	ldr.TotalSize = res.TotalSize
	return ldr
}


type DatabaseSpecification struct {
	Name       string 
	SizeOnDisk int64  
	Empty      bool   
}


type UpdateResult struct {
	MatchedCount  int64       
	ModifiedCount int64       
	UpsertedCount int64       
	UpsertedID    interface{} 
}





func (result *UpdateResult) UnmarshalBSON(b []byte) error {
	
	elems, err := bson.Raw(b).Elements()
	if err != nil {
		return err
	}

	for _, elem := range elems {
		switch elem.Key() {
		case "n":
			switch elem.Value().Type {
			case bson.TypeInt32:
				result.MatchedCount = int64(elem.Value().Int32())
			case bson.TypeInt64:
				result.MatchedCount = elem.Value().Int64()
			default:
				return fmt.Errorf("Received invalid type for n, should be Int32 or Int64, received %s", elem.Value().Type)
			}
		case "nModified":
			switch elem.Value().Type {
			case bson.TypeInt32:
				result.ModifiedCount = int64(elem.Value().Int32())
			case bson.TypeInt64:
				result.ModifiedCount = elem.Value().Int64()
			default:
				return fmt.Errorf("Received invalid type for nModified, should be Int32 or Int64, received %s", elem.Value().Type)
			}
		case "upserted":
			switch elem.Value().Type {
			case bson.TypeArray:
				e, err := elem.Value().Array().IndexErr(0)
				if err != nil {
					break
				}
				if e.Value().Type != bson.TypeEmbeddedDocument {
					break
				}
				var d struct {
					ID interface{} `bson:"_id"`
				}
				err = bson.Unmarshal(e.Value().Document(), &d)
				if err != nil {
					return err
				}
				result.UpsertedID = d.ID
			default:
				return fmt.Errorf("Received invalid type for upserted, should be Array, received %s", elem.Value().Type)
			}
		}
	}

	return nil
}



type IndexSpecification struct {
	
	Name string

	
	Namespace string

	
	KeysDocument bson.Raw

	
	Version int32

	
	
	ExpireAfterSeconds *int32

	
	
	Sparse *bool

	
	
	Unique *bool

	
	Clustered *bool
}

var _ bson.Unmarshaler = (*IndexSpecification)(nil)

type unmarshalIndexSpecification struct {
	Name               string   `bson:"name"`
	Namespace          string   `bson:"ns"`
	KeysDocument       bson.Raw `bson:"key"`
	Version            int32    `bson:"v"`
	ExpireAfterSeconds *int32   `bson:"expireAfterSeconds"`
	Sparse             *bool    `bson:"sparse"`
	Unique             *bool    `bson:"unique"`
	Clustered          *bool    `bson:"clustered"`
}




func (i *IndexSpecification) UnmarshalBSON(data []byte) error {
	var temp unmarshalIndexSpecification
	if err := bson.Unmarshal(data, &temp); err != nil {
		return err
	}

	i.Name = temp.Name
	i.Namespace = temp.Namespace
	i.KeysDocument = temp.KeysDocument
	i.Version = temp.Version
	i.ExpireAfterSeconds = temp.ExpireAfterSeconds
	i.Sparse = temp.Sparse
	i.Unique = temp.Unique
	i.Clustered = temp.Clustered
	return nil
}



type CollectionSpecification struct {
	
	Name string

	
	Type string

	
	ReadOnly bool

	
	
	UUID *primitive.Binary

	
	Options bson.Raw

	
	
	IDIndex *IndexSpecification
}

var _ bson.Unmarshaler = (*CollectionSpecification)(nil)



type unmarshalCollectionSpecification struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
	Info *struct {
		ReadOnly bool              `bson:"readOnly"`
		UUID     *primitive.Binary `bson:"uuid"`
	} `bson:"info"`
	Options bson.Raw            `bson:"options"`
	IDIndex *IndexSpecification `bson:"idIndex"`
}





func (cs *CollectionSpecification) UnmarshalBSON(data []byte) error {
	var temp unmarshalCollectionSpecification
	if err := bson.Unmarshal(data, &temp); err != nil {
		return err
	}

	cs.Name = temp.Name
	cs.Type = temp.Type
	if cs.Type == "" {
		
		
		cs.Type = "collection"
	}
	if temp.Info != nil {
		cs.ReadOnly = temp.Info.ReadOnly
		cs.UUID = temp.Info.UUID
	}
	cs.Options = temp.Options
	cs.IDIndex = temp.IDIndex
	return nil
}
