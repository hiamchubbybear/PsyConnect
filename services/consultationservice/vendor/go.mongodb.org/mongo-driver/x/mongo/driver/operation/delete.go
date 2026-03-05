





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


type Delete struct {
	authenticator driver.Authenticator
	comment       bsoncore.Value
	deletes       []bsoncore.Document
	ordered       *bool
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
	hint          *bool
	result        DeleteResult
	serverAPI     *driver.ServerAPIOptions
	let           bsoncore.Document
	timeout       *time.Duration
	logger        *logger.Logger
}


type DeleteResult struct {
	
	N int64
}

func buildDeleteResult(response bsoncore.Document) (DeleteResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DeleteResult{}, err
	}
	dr := DeleteResult{}
	for _, element := range elements {
		if element.Key() == "n" {
			var ok bool
			dr.N, ok = element.Value().AsInt64OK()
			if !ok {
				return dr, fmt.Errorf("response field 'n' is type int32 or int64, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dr, nil
}


func NewDelete(deletes ...bsoncore.Document) *Delete {
	return &Delete{
		deletes: deletes,
	}
}


func (d *Delete) Result() DeleteResult { return d.result }

func (d *Delete) processResponse(info driver.ResponseInfo) error {
	dr, err := buildDeleteResult(info.ServerResponse)
	d.result.N += dr.N
	return err
}


func (d *Delete) Execute(ctx context.Context) error {
	if d.deployment == nil {
		return errors.New("the Delete operation must have a Deployment set before Execute can be called")
	}
	batches := &driver.Batches{
		Identifier: "deletes",
		Documents:  d.deletes,
		Ordered:    d.ordered,
	}

	return driver.Operation{
		CommandFn:         d.command,
		ProcessResponseFn: d.processResponse,
		Batches:           batches,
		RetryMode:         d.retry,
		Type:              driver.Write,
		Client:            d.session,
		Clock:             d.clock,
		CommandMonitor:    d.monitor,
		Crypt:             d.crypt,
		Database:          d.database,
		Deployment:        d.deployment,
		Selector:          d.selector,
		WriteConcern:      d.writeConcern,
		ServerAPI:         d.serverAPI,
		Timeout:           d.timeout,
		Logger:            d.logger,
		Name:              driverutil.DeleteOp,
		Authenticator:     d.authenticator,
	}.Execute(ctx)

}

func (d *Delete) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "delete", d.collection)
	if d.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", d.comment)
	}
	if d.ordered != nil {
		dst = bsoncore.AppendBooleanElement(dst, "ordered", *d.ordered)
	}
	if d.hint != nil && *d.hint {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'hint' command parameter requires a minimum server wire version of 5")
		}
		if !d.writeConcern.Acknowledged() {
			return nil, errUnacknowledgedHint
		}
	}
	if d.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", d.let)
	}
	return dst, nil
}




func (d *Delete) Deletes(deletes ...bsoncore.Document) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.deletes = deletes
	return d
}



func (d *Delete) Ordered(ordered bool) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.ordered = &ordered
	return d
}


func (d *Delete) Session(session *session.Client) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.session = session
	return d
}


func (d *Delete) ClusterClock(clock *session.ClusterClock) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.clock = clock
	return d
}


func (d *Delete) Collection(collection string) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.collection = collection
	return d
}


func (d *Delete) Comment(comment bsoncore.Value) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.comment = comment
	return d
}


func (d *Delete) CommandMonitor(monitor *event.CommandMonitor) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.monitor = monitor
	return d
}


func (d *Delete) Crypt(crypt driver.Crypt) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.crypt = crypt
	return d
}


func (d *Delete) Database(database string) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.database = database
	return d
}


func (d *Delete) Deployment(deployment driver.Deployment) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.deployment = deployment
	return d
}


func (d *Delete) ServerSelector(selector description.ServerSelector) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.selector = selector
	return d
}


func (d *Delete) WriteConcern(writeConcern *writeconcern.WriteConcern) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.writeConcern = writeConcern
	return d
}



func (d *Delete) Retry(retry driver.RetryMode) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.retry = &retry
	return d
}




func (d *Delete) Hint(hint bool) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.hint = &hint
	return d
}


func (d *Delete) ServerAPI(serverAPI *driver.ServerAPIOptions) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.serverAPI = serverAPI
	return d
}


func (d *Delete) Let(let bsoncore.Document) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.let = let
	return d
}


func (d *Delete) Timeout(timeout *time.Duration) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.timeout = timeout
	return d
}


func (d *Delete) Logger(logger *logger.Logger) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.logger = logger

	return d
}


func (d *Delete) Authenticator(authenticator driver.Authenticator) *Delete {
	if d == nil {
		d = new(Delete)
	}

	d.authenticator = authenticator
	return d
}
