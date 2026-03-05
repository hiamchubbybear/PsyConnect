





package operation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type DropCollection struct {
	authenticator driver.Authenticator
	session       *session.Client
	clock         *session.ClusterClock
	collection    string
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	result        DropCollectionResult
	serverAPI     *driver.ServerAPIOptions
	timeout       *time.Duration
}


type DropCollectionResult struct {
	
	NIndexesWas int32
	
	Ns string
}

func buildDropCollectionResult(response bsoncore.Document) (DropCollectionResult, error) {
	elements, err := response.Elements()
	if err != nil {
		return DropCollectionResult{}, err
	}
	dcr := DropCollectionResult{}
	for _, element := range elements {
		switch element.Key() {
		case "nIndexesWas":
			var ok bool
			dcr.NIndexesWas, ok = element.Value().AsInt32OK()
			if !ok {
				return dcr, fmt.Errorf("response field 'nIndexesWas' is type int32, but received BSON type %s", element.Value().Type)
			}
		case "ns":
			var ok bool
			dcr.Ns, ok = element.Value().StringValueOK()
			if !ok {
				return dcr, fmt.Errorf("response field 'ns' is type string, but received BSON type %s", element.Value().Type)
			}
		}
	}
	return dcr, nil
}


func NewDropCollection() *DropCollection {
	return &DropCollection{}
}


func (dc *DropCollection) Result() DropCollectionResult { return dc.result }

func (dc *DropCollection) processResponse(info driver.ResponseInfo) error {
	var err error
	dc.result, err = buildDropCollectionResult(info.ServerResponse)
	return err
}


func (dc *DropCollection) Execute(ctx context.Context) error {
	if dc.deployment == nil {
		return errors.New("the DropCollection operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         dc.command,
		ProcessResponseFn: dc.processResponse,
		Client:            dc.session,
		Clock:             dc.clock,
		CommandMonitor:    dc.monitor,
		Crypt:             dc.crypt,
		Database:          dc.database,
		Deployment:        dc.deployment,
		Selector:          dc.selector,
		WriteConcern:      dc.writeConcern,
		ServerAPI:         dc.serverAPI,
		Timeout:           dc.timeout,
		Name:              driverutil.DropOp,
		Authenticator:     dc.authenticator,
	}.Execute(ctx)

}

func (dc *DropCollection) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "drop", dc.collection)
	return dst, nil
}


func (dc *DropCollection) Session(session *session.Client) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.session = session
	return dc
}


func (dc *DropCollection) ClusterClock(clock *session.ClusterClock) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.clock = clock
	return dc
}


func (dc *DropCollection) Collection(collection string) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.collection = collection
	return dc
}


func (dc *DropCollection) CommandMonitor(monitor *event.CommandMonitor) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.monitor = monitor
	return dc
}


func (dc *DropCollection) Crypt(crypt driver.Crypt) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.crypt = crypt
	return dc
}


func (dc *DropCollection) Database(database string) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.database = database
	return dc
}


func (dc *DropCollection) Deployment(deployment driver.Deployment) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.deployment = deployment
	return dc
}


func (dc *DropCollection) ServerSelector(selector description.ServerSelector) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.selector = selector
	return dc
}


func (dc *DropCollection) WriteConcern(writeConcern *writeconcern.WriteConcern) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.writeConcern = writeConcern
	return dc
}


func (dc *DropCollection) ServerAPI(serverAPI *driver.ServerAPIOptions) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.serverAPI = serverAPI
	return dc
}


func (dc *DropCollection) Timeout(timeout *time.Duration) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.timeout = timeout
	return dc
}


func (dc *DropCollection) Authenticator(authenticator driver.Authenticator) *DropCollection {
	if dc == nil {
		dc = new(DropCollection)
	}

	dc.authenticator = authenticator
	return dc
}
