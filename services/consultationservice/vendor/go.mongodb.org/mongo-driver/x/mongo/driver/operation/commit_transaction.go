





package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type CommitTransaction struct {
	authenticator driver.Authenticator
	maxTime       *time.Duration
	recoveryToken bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	writeConcern  *writeconcern.WriteConcern
	retry         *driver.RetryMode
	serverAPI     *driver.ServerAPIOptions
}


func NewCommitTransaction() *CommitTransaction {
	return &CommitTransaction{}
}

func (ct *CommitTransaction) processResponse(driver.ResponseInfo) error {
	var err error
	return err
}


func (ct *CommitTransaction) Execute(ctx context.Context) error {
	if ct.deployment == nil {
		return errors.New("the CommitTransaction operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         ct.command,
		ProcessResponseFn: ct.processResponse,
		RetryMode:         ct.retry,
		Type:              driver.Write,
		Client:            ct.session,
		Clock:             ct.clock,
		CommandMonitor:    ct.monitor,
		Crypt:             ct.crypt,
		Database:          ct.database,
		Deployment:        ct.deployment,
		MaxTime:           ct.maxTime,
		Selector:          ct.selector,
		WriteConcern:      ct.writeConcern,
		ServerAPI:         ct.serverAPI,
		Name:              driverutil.CommitTransactionOp,
		Authenticator:     ct.authenticator,
	}.Execute(ctx)

}

func (ct *CommitTransaction) command(dst []byte, _ description.SelectedServer) ([]byte, error) {

	dst = bsoncore.AppendInt32Element(dst, "commitTransaction", 1)
	if ct.recoveryToken != nil {
		dst = bsoncore.AppendDocumentElement(dst, "recoveryToken", ct.recoveryToken)
	}
	return dst, nil
}


func (ct *CommitTransaction) MaxTime(maxTime *time.Duration) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.maxTime = maxTime
	return ct
}


func (ct *CommitTransaction) RecoveryToken(recoveryToken bsoncore.Document) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.recoveryToken = recoveryToken
	return ct
}


func (ct *CommitTransaction) Session(session *session.Client) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.session = session
	return ct
}


func (ct *CommitTransaction) ClusterClock(clock *session.ClusterClock) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.clock = clock
	return ct
}


func (ct *CommitTransaction) CommandMonitor(monitor *event.CommandMonitor) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.monitor = monitor
	return ct
}


func (ct *CommitTransaction) Crypt(crypt driver.Crypt) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.crypt = crypt
	return ct
}


func (ct *CommitTransaction) Database(database string) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.database = database
	return ct
}


func (ct *CommitTransaction) Deployment(deployment driver.Deployment) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.deployment = deployment
	return ct
}


func (ct *CommitTransaction) ServerSelector(selector description.ServerSelector) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.selector = selector
	return ct
}


func (ct *CommitTransaction) WriteConcern(writeConcern *writeconcern.WriteConcern) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.writeConcern = writeConcern
	return ct
}



func (ct *CommitTransaction) Retry(retry driver.RetryMode) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.retry = &retry
	return ct
}


func (ct *CommitTransaction) ServerAPI(serverAPI *driver.ServerAPIOptions) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.serverAPI = serverAPI
	return ct
}


func (ct *CommitTransaction) Authenticator(authenticator driver.Authenticator) *CommitTransaction {
	if ct == nil {
		ct = new(CommitTransaction)
	}

	ct.authenticator = authenticator
	return ct
}
