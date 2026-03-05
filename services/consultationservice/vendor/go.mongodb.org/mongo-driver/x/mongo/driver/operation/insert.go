





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type Insert struct {
	authenticator            driver.Authenticator
	bypassDocumentValidation *bool
	comment                  bsoncore.Value
	documents                []bsoncore.Document
	ordered                  *bool
	session                  *session.Client
	clock                    *session.ClusterClock
	collection               string
	monitor                  *event.CommandMonitor
	crypt                    driver.Crypt
	database                 string
	deployment               driver.Deployment
	selector                 description.ServerSelector
	writeConcern             *writeconcern.WriteConcern
	retry                    *driver.RetryMode
	result                   InsertResult
	serverAPI                *driver.ServerAPIOptions
	timeout                  *time.Duration
	bypassEmptyTsReplacement *bool
	logger                   *logger.Logger
}


type InsertResult struct {
	
	N int64
}

func buildInsertResult(response bsoncore.Document) (InsertResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return InsertResult{}, err
	}
	ir := InsertResult{}
	for _, element := range elements {
		if element.Key() == "n" {
			var ok bool
			ir.N, ok = element.Value().AsInt64OK()
			if !ok {
				return ir, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return ir, nil
}


func NewInsert(documents ...bsoncore.Document) *Insert {
	return &Insert{
		documents: documents,
	}
}


func (i *Insert) Result() InsertResult { return i.result }

func (i *Insert) processResponse(info driver.ResponseInfo) error {
	ir, err := buildInsertResult(info.ServerResponse)
	i.result.N += ir.N
	return err
}


func (i *Insert) Execute(ctx context.Context) error {
	if i.deployment == nil {
		return errors.New("the Insert operation must have a Deployment set before Execute can be called")
	}
	batches := &driver.Batches{
		Identifier: "documents",
		Documents:  i.documents,
		Ordered:    i.ordered,
	}

	return driver.Operation{
		CommandFn:         i.command,
		ProcessResponseFn: i.processResponse,
		Batches:           batches,
		RetryMode:         i.retry,
		Type:              driver.Write,
		Client:            i.session,
		Clock:             i.clock,
		CommandMonitor:    i.monitor,
		Crypt:             i.crypt,
		Database:          i.database,
		Deployment:        i.deployment,
		Selector:          i.selector,
		WriteConcern:      i.writeConcern,
		ServerAPI:         i.serverAPI,
		Timeout:           i.timeout,
		Logger:            i.logger,
		Name:              driverutil.InsertOp,
		Authenticator:     i.authenticator,
	}.Execute(ctx)

}

func (i *Insert) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "insert", i.collection)
	if i.bypassDocumentValidation != nil && (desc.WireVersion != nil && desc.WireVersion.Includes(4)) {
		dst = bsoncore.AppendBooleanElement(dst, "bypassDocumentValidation", *i.bypassDocumentValidation)
	}
	if i.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", i.comment)
	}
	if i.ordered != nil {
		dst = bsoncore.AppendBooleanElement(dst, "ordered", *i.ordered)
	}
	if i.bypassEmptyTsReplacement != nil {
		dst = bsoncore.AppendBooleanElement(dst, "bypassEmptyTsReplacement", *i.bypassEmptyTsReplacement)
	}
	return dst, nil
}



func (i *Insert) BypassDocumentValidation(bypassDocumentValidation bool) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.bypassDocumentValidation = &bypassDocumentValidation
	return i
}


func (i *Insert) Comment(comment bsoncore.Value) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.comment = comment
	return i
}



func (i *Insert) Documents(documents ...bsoncore.Document) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.documents = documents
	return i
}



func (i *Insert) Ordered(ordered bool) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.ordered = &ordered
	return i
}


func (i *Insert) Session(session *session.Client) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.session = session
	return i
}


func (i *Insert) ClusterClock(clock *session.ClusterClock) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.clock = clock
	return i
}


func (i *Insert) Collection(collection string) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.collection = collection
	return i
}


func (i *Insert) CommandMonitor(monitor *event.CommandMonitor) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.monitor = monitor
	return i
}


func (i *Insert) Crypt(crypt driver.Crypt) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.crypt = crypt
	return i
}


func (i *Insert) Database(database string) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.database = database
	return i
}


func (i *Insert) Deployment(deployment driver.Deployment) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.deployment = deployment
	return i
}


func (i *Insert) ServerSelector(selector description.ServerSelector) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.selector = selector
	return i
}


func (i *Insert) WriteConcern(writeConcern *writeconcern.WriteConcern) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.writeConcern = writeConcern
	return i
}



func (i *Insert) Retry(retry driver.RetryMode) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.retry = &retry
	return i
}


func (i *Insert) ServerAPI(serverAPI *driver.ServerAPIOptions) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.serverAPI = serverAPI
	return i
}


func (i *Insert) Timeout(timeout *time.Duration) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.timeout = timeout
	return i
}


func (i *Insert) Logger(logger *logger.Logger) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.logger = logger
	return i
}


func (i *Insert) Authenticator(authenticator driver.Authenticator) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.authenticator = authenticator
	return i
}


func (i *Insert) BypassEmptyTsReplacement(bypassEmptyTsReplacement bool) *Insert {
	if i == nil {
		i = new(Insert)
	}

	i.bypassEmptyTsReplacement = &bypassEmptyTsReplacement
	return i
}
