





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type Count struct {
	authenticator  driver.Authenticator
	maxTime        *time.Duration
	query          bsoncore.Document
	session        *session.Client
	clock          *session.ClusterClock
	collection     string
	comment        bsoncore.Value
	monitor        *event.CommandMonitor
	crypt          driver.Crypt
	database       string
	deployment     driver.Deployment
	readConcern    *readconcern.ReadConcern
	readPreference *readpref.ReadPref
	selector       description.ServerSelector
	retry          *driver.RetryMode
	result         CountResult
	serverAPI      *driver.ServerAPIOptions
	timeout        *time.Duration
}


type CountResult struct {
	
	N int64
}

func buildCountResult(response bsoncore.Document) (CountResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return CountResult{}, err
	}
	cr := CountResult{}
	for _, element := range elements {
		switch element.Key() {
		case "n": 
			var ok bool
			cr.N, ok = element.Value().AsInt64OK()
			if !ok {
				return cr, fmt.Errorf("response field 'n' is type int64, but received BSON type %s",
					element.Value().Type)
			}
		case "cursor": 
			firstBatch, err := element.Value().Document().LookupErr("firstBatch")
			if err != nil {
				return cr, err
			}

			
			val := firstBatch.Array().Index(0)
			count, err := val.Document().LookupErr("n")
			if err != nil {
				return cr, err
			}

			
			var ok bool
			cr.N, ok = count.AsInt64OK()
			if !ok {
				return cr, fmt.Errorf("response field 'n' is type int64, but received BSON type %s",
					element.Value().Type)
			}
		}
	}
	return cr, nil
}


func NewCount() *Count {
	return &Count{}
}


func (c *Count) Result() CountResult { return c.result }

func (c *Count) processResponse(info driver.ResponseInfo) error {
	var err error
	c.result, err = buildCountResult(info.ServerResponse)
	return err
}


func (c *Count) Execute(ctx context.Context) error {
	if c.deployment == nil {
		return errors.New("the Count operation must have a Deployment set before Execute can be called")
	}

	err := driver.Operation{
		CommandFn:         c.command,
		ProcessResponseFn: c.processResponse,
		RetryMode:         c.retry,
		Type:              driver.Read,
		Client:            c.session,
		Clock:             c.clock,
		CommandMonitor:    c.monitor,
		Crypt:             c.crypt,
		Database:          c.database,
		Deployment:        c.deployment,
		MaxTime:           c.maxTime,
		ReadConcern:       c.readConcern,
		ReadPreference:    c.readPreference,
		Selector:          c.selector,
		ServerAPI:         c.serverAPI,
		Timeout:           c.timeout,
		Name:              driverutil.CountOp,
		Authenticator:     c.authenticator,
	}.Execute(ctx)

	
	if err != nil {
		dErr, ok := err.(driver.Error)
		if ok && dErr.Code == 26 {
			err = nil
		}
	}
	return err
}

func (c *Count) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "count", c.collection)
	if c.query != nil {
		dst = bsoncore.AppendDocumentElement(dst, "query", c.query)
	}
	if c.comment.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "comment", c.comment)
	}
	return dst, nil
}


func (c *Count) MaxTime(maxTime *time.Duration) *Count {
	if c == nil {
		c = new(Count)
	}

	c.maxTime = maxTime
	return c
}


func (c *Count) Query(query bsoncore.Document) *Count {
	if c == nil {
		c = new(Count)
	}

	c.query = query
	return c
}


func (c *Count) Session(session *session.Client) *Count {
	if c == nil {
		c = new(Count)
	}

	c.session = session
	return c
}


func (c *Count) ClusterClock(clock *session.ClusterClock) *Count {
	if c == nil {
		c = new(Count)
	}

	c.clock = clock
	return c
}


func (c *Count) Collection(collection string) *Count {
	if c == nil {
		c = new(Count)
	}

	c.collection = collection
	return c
}


func (c *Count) Comment(comment bsoncore.Value) *Count {
	if c == nil {
		c = new(Count)
	}

	c.comment = comment
	return c
}


func (c *Count) CommandMonitor(monitor *event.CommandMonitor) *Count {
	if c == nil {
		c = new(Count)
	}

	c.monitor = monitor
	return c
}


func (c *Count) Crypt(crypt driver.Crypt) *Count {
	if c == nil {
		c = new(Count)
	}

	c.crypt = crypt
	return c
}


func (c *Count) Database(database string) *Count {
	if c == nil {
		c = new(Count)
	}

	c.database = database
	return c
}


func (c *Count) Deployment(deployment driver.Deployment) *Count {
	if c == nil {
		c = new(Count)
	}

	c.deployment = deployment
	return c
}


func (c *Count) ReadConcern(readConcern *readconcern.ReadConcern) *Count {
	if c == nil {
		c = new(Count)
	}

	c.readConcern = readConcern
	return c
}


func (c *Count) ReadPreference(readPreference *readpref.ReadPref) *Count {
	if c == nil {
		c = new(Count)
	}

	c.readPreference = readPreference
	return c
}


func (c *Count) ServerSelector(selector description.ServerSelector) *Count {
	if c == nil {
		c = new(Count)
	}

	c.selector = selector
	return c
}



func (c *Count) Retry(retry driver.RetryMode) *Count {
	if c == nil {
		c = new(Count)
	}

	c.retry = &retry
	return c
}


func (c *Count) ServerAPI(serverAPI *driver.ServerAPIOptions) *Count {
	if c == nil {
		c = new(Count)
	}

	c.serverAPI = serverAPI
	return c
}


func (c *Count) Timeout(timeout *time.Duration) *Count {
	if c == nil {
		c = new(Count)
	}

	c.timeout = timeout
	return c
}


func (c *Count) Authenticator(authenticator driver.Authenticator) *Count {
	if c == nil {
		c = new(Count)
	}

	c.authenticator = authenticator
	return c
}
