





package operation

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type Create struct {
	authenticator                driver.Authenticator
	capped                       *bool
	collation                    bsoncore.Document
	changeStreamPreAndPostImages bsoncore.Document
	collectionName               *string
	indexOptionDefaults          bsoncore.Document
	max                          *int64
	pipeline                     bsoncore.Document
	size                         *int64
	storageEngine                bsoncore.Document
	validationAction             *string
	validationLevel              *string
	validator                    bsoncore.Document
	viewOn                       *string
	session                      *session.Client
	clock                        *session.ClusterClock
	monitor                      *event.CommandMonitor
	crypt                        driver.Crypt
	database                     string
	deployment                   driver.Deployment
	selector                     description.ServerSelector
	writeConcern                 *writeconcern.WriteConcern
	serverAPI                    *driver.ServerAPIOptions
	expireAfterSeconds           *int64
	timeSeries                   bsoncore.Document
	encryptedFields              bsoncore.Document
	clusteredIndex               bsoncore.Document
}


func NewCreate(collectionName string) *Create {
	return &Create{
		collectionName: &collectionName,
	}
}

func (c *Create) processResponse(driver.ResponseInfo) error {
	return nil
}


func (c *Create) Execute(ctx context.Context) error {
	if c.deployment == nil {
		return errors.New("the Create operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         c.command,
		ProcessResponseFn: c.processResponse,
		Client:            c.session,
		Clock:             c.clock,
		CommandMonitor:    c.monitor,
		Crypt:             c.crypt,
		Database:          c.database,
		Deployment:        c.deployment,
		Selector:          c.selector,
		WriteConcern:      c.writeConcern,
		ServerAPI:         c.serverAPI,
		Authenticator:     c.authenticator,
	}.Execute(ctx)
}

func (c *Create) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	if c.collectionName != nil {
		dst = bsoncore.AppendStringElement(dst, "create", *c.collectionName)
	}
	if c.capped != nil {
		dst = bsoncore.AppendBooleanElement(dst, "capped", *c.capped)
	}
	if c.changeStreamPreAndPostImages != nil {
		dst = bsoncore.AppendDocumentElement(dst, "changeStreamPreAndPostImages", c.changeStreamPreAndPostImages)
	}
	if c.collation != nil {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", c.collation)
	}
	if c.indexOptionDefaults != nil {
		dst = bsoncore.AppendDocumentElement(dst, "indexOptionDefaults", c.indexOptionDefaults)
	}
	if c.max != nil {
		dst = bsoncore.AppendInt64Element(dst, "max", *c.max)
	}
	if c.pipeline != nil {
		dst = bsoncore.AppendArrayElement(dst, "pipeline", c.pipeline)
	}
	if c.size != nil {
		dst = bsoncore.AppendInt64Element(dst, "size", *c.size)
	}
	if c.storageEngine != nil {
		dst = bsoncore.AppendDocumentElement(dst, "storageEngine", c.storageEngine)
	}
	if c.validationAction != nil {
		dst = bsoncore.AppendStringElement(dst, "validationAction", *c.validationAction)
	}
	if c.validationLevel != nil {
		dst = bsoncore.AppendStringElement(dst, "validationLevel", *c.validationLevel)
	}
	if c.validator != nil {
		dst = bsoncore.AppendDocumentElement(dst, "validator", c.validator)
	}
	if c.viewOn != nil {
		dst = bsoncore.AppendStringElement(dst, "viewOn", *c.viewOn)
	}
	if c.expireAfterSeconds != nil {
		dst = bsoncore.AppendInt64Element(dst, "expireAfterSeconds", *c.expireAfterSeconds)
	}
	if c.timeSeries != nil {
		dst = bsoncore.AppendDocumentElement(dst, "timeseries", c.timeSeries)
	}
	if c.encryptedFields != nil {
		dst = bsoncore.AppendDocumentElement(dst, "encryptedFields", c.encryptedFields)
	}
	if c.clusteredIndex != nil {
		dst = bsoncore.AppendDocumentElement(dst, "clusteredIndex", c.clusteredIndex)
	}
	return dst, nil
}


func (c *Create) Capped(capped bool) *Create {
	if c == nil {
		c = new(Create)
	}

	c.capped = &capped
	return c
}


func (c *Create) Collation(collation bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.collation = collation
	return c
}



func (c *Create) ChangeStreamPreAndPostImages(csppi bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.changeStreamPreAndPostImages = csppi
	return c
}


func (c *Create) CollectionName(collectionName string) *Create {
	if c == nil {
		c = new(Create)
	}

	c.collectionName = &collectionName
	return c
}


func (c *Create) IndexOptionDefaults(indexOptionDefaults bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.indexOptionDefaults = indexOptionDefaults
	return c
}


func (c *Create) Max(max int64) *Create {
	if c == nil {
		c = new(Create)
	}

	c.max = &max
	return c
}


func (c *Create) Pipeline(pipeline bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.pipeline = pipeline
	return c
}


func (c *Create) Size(size int64) *Create {
	if c == nil {
		c = new(Create)
	}

	c.size = &size
	return c
}


func (c *Create) StorageEngine(storageEngine bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.storageEngine = storageEngine
	return c
}


func (c *Create) ValidationAction(validationAction string) *Create {
	if c == nil {
		c = new(Create)
	}

	c.validationAction = &validationAction
	return c
}



func (c *Create) ValidationLevel(validationLevel string) *Create {
	if c == nil {
		c = new(Create)
	}

	c.validationLevel = &validationLevel
	return c
}


func (c *Create) Validator(validator bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.validator = validator
	return c
}


func (c *Create) ViewOn(viewOn string) *Create {
	if c == nil {
		c = new(Create)
	}

	c.viewOn = &viewOn
	return c
}


func (c *Create) Session(session *session.Client) *Create {
	if c == nil {
		c = new(Create)
	}

	c.session = session
	return c
}


func (c *Create) ClusterClock(clock *session.ClusterClock) *Create {
	if c == nil {
		c = new(Create)
	}

	c.clock = clock
	return c
}


func (c *Create) CommandMonitor(monitor *event.CommandMonitor) *Create {
	if c == nil {
		c = new(Create)
	}

	c.monitor = monitor
	return c
}


func (c *Create) Crypt(crypt driver.Crypt) *Create {
	if c == nil {
		c = new(Create)
	}

	c.crypt = crypt
	return c
}


func (c *Create) Database(database string) *Create {
	if c == nil {
		c = new(Create)
	}

	c.database = database
	return c
}


func (c *Create) Deployment(deployment driver.Deployment) *Create {
	if c == nil {
		c = new(Create)
	}

	c.deployment = deployment
	return c
}


func (c *Create) ServerSelector(selector description.ServerSelector) *Create {
	if c == nil {
		c = new(Create)
	}

	c.selector = selector
	return c
}


func (c *Create) WriteConcern(writeConcern *writeconcern.WriteConcern) *Create {
	if c == nil {
		c = new(Create)
	}

	c.writeConcern = writeConcern
	return c
}


func (c *Create) ServerAPI(serverAPI *driver.ServerAPIOptions) *Create {
	if c == nil {
		c = new(Create)
	}

	c.serverAPI = serverAPI
	return c
}


func (c *Create) ExpireAfterSeconds(eas int64) *Create {
	if c == nil {
		c = new(Create)
	}

	c.expireAfterSeconds = &eas
	return c
}


func (c *Create) TimeSeries(timeSeries bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.timeSeries = timeSeries
	return c
}


func (c *Create) EncryptedFields(ef bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.encryptedFields = ef
	return c
}


func (c *Create) ClusteredIndex(ci bsoncore.Document) *Create {
	if c == nil {
		c = new(Create)
	}

	c.clusteredIndex = ci
	return c
}


func (c *Create) Authenticator(authenticator driver.Authenticator) *Create {
	if c == nil {
		c = new(Create)
	}

	c.authenticator = authenticator
	return c
}
