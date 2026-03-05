





package operation

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)


type EndSessions struct {
	authenticator driver.Authenticator
	sessionIDs    bsoncore.Document
	session       *session.Client
	clock         *session.ClusterClock
	monitor       *event.CommandMonitor
	crypt         driver.Crypt
	database      string
	deployment    driver.Deployment
	selector      description.ServerSelector
	serverAPI     *driver.ServerAPIOptions
}


func NewEndSessions(sessionIDs bsoncore.Document) *EndSessions {
	return &EndSessions{
		sessionIDs: sessionIDs,
	}
}

func (es *EndSessions) processResponse(driver.ResponseInfo) error {
	var err error
	return err
}


func (es *EndSessions) Execute(ctx context.Context) error {
	if es.deployment == nil {
		return errors.New("the EndSessions operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         es.command,
		ProcessResponseFn: es.processResponse,
		Client:            es.session,
		Clock:             es.clock,
		CommandMonitor:    es.monitor,
		Crypt:             es.crypt,
		Database:          es.database,
		Deployment:        es.deployment,
		Selector:          es.selector,
		ServerAPI:         es.serverAPI,
		Name:              driverutil.EndSessionsOp,
		Authenticator:     es.authenticator,
	}.Execute(ctx)

}

func (es *EndSessions) command(dst []byte, _ description.SelectedServer) ([]byte, error) {
	if es.sessionIDs != nil {
		dst = bsoncore.AppendArrayElement(dst, "endSessions", es.sessionIDs)
	}
	return dst, nil
}


func (es *EndSessions) SessionIDs(sessionIDs bsoncore.Document) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.sessionIDs = sessionIDs
	return es
}


func (es *EndSessions) Session(session *session.Client) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.session = session
	return es
}


func (es *EndSessions) ClusterClock(clock *session.ClusterClock) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.clock = clock
	return es
}


func (es *EndSessions) CommandMonitor(monitor *event.CommandMonitor) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.monitor = monitor
	return es
}


func (es *EndSessions) Crypt(crypt driver.Crypt) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.crypt = crypt
	return es
}


func (es *EndSessions) Database(database string) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.database = database
	return es
}


func (es *EndSessions) Deployment(deployment driver.Deployment) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.deployment = deployment
	return es
}


func (es *EndSessions) ServerSelector(selector description.ServerSelector) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.selector = selector
	return es
}


func (es *EndSessions) ServerAPI(serverAPI *driver.ServerAPIOptions) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.serverAPI = serverAPI
	return es
}


func (es *EndSessions) Authenticator(authenticator driver.Authenticator) *EndSessions {
	if es == nil {
		es = new(EndSessions)
	}

	es.authenticator = authenticator
	return es
}
