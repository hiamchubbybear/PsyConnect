





package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)



var ErrWrongClient = errors.New("session was not created by this client")

var withTransactionTimeout = 120 * time.Second









type SessionContext interface {
	context.Context
	Session
}

type sessionContext struct {
	context.Context
	Session
}

type sessionKey struct {
}


func NewSessionContext(ctx context.Context, sess Session) SessionContext {
	return &sessionContext{
		Context: context.WithValue(ctx, sessionKey{}, sess),
		Session: sess,
	}
}




func SessionFromContext(ctx context.Context) Session {
	val := ctx.Value(sessionKey{})
	if val == nil {
		return nil
	}

	sess, ok := val.(Session)
	if !ok {
		return nil
	}

	return sess
}











type Session interface {
	
	
	
	StartTransaction(...*options.TransactionOptions) error

	
	
	
	AbortTransaction(context.Context) error

	
	
	
	CommitTransaction(context.Context) error

	
	
	
	
	
	
	
	
	
	
	
	WithTransaction(ctx context.Context, fn func(ctx SessionContext) (interface{}, error),
		opts ...*options.TransactionOptions) (interface{}, error)

	
	EndSession(context.Context)

	
	ClusterTime() bson.Raw

	
	OperationTime() *primitive.Timestamp

	
	Client() *Client

	
	
	ID() bson.Raw

	
	
	AdvanceClusterTime(bson.Raw) error

	
	
	AdvanceOperationTime(*primitive.Timestamp) error

	session()
}





type XSession interface {
	ClientSession() *session.Client
}


type sessionImpl struct {
	clientSession       *session.Client
	client              *Client
	deployment          driver.Deployment
	didCommitAfterStart bool 
}

var _ Session = &sessionImpl{}
var _ XSession = &sessionImpl{}


func (s *sessionImpl) ClientSession() *session.Client {
	return s.clientSession
}


func (s *sessionImpl) ID() bson.Raw {
	return bson.Raw(s.clientSession.SessionID)
}


func (s *sessionImpl) EndSession(ctx context.Context) {
	if s.clientSession.TransactionInProgress() {
		
		_ = s.AbortTransaction(ctx)
	}
	s.clientSession.EndSession()
}


func (s *sessionImpl) WithTransaction(ctx context.Context, fn func(ctx SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (interface{}, error) {
	timeout := time.NewTimer(withTransactionTimeout)
	defer timeout.Stop()
	var err error
	for {
		err = s.StartTransaction(opts...)
		if err != nil {
			return nil, err
		}

		res, err := fn(NewSessionContext(ctx, s))
		if err != nil {
			if s.clientSession.TransactionRunning() {
				
				
				_ = s.AbortTransaction(newBackgroundContext(ctx))
			}

			select {
			case <-timeout.C:
				return nil, err
			default:
			}

			if errorHasLabel(err, driver.TransientTransactionError) {
				continue
			}
			return res, err
		}

		
		
		err = s.clientSession.CheckAbortTransaction()
		if err != nil {
			return res, nil
		}

		
		
		
		
		
		
		
		if ctx.Err() != nil {
			
			
			_ = s.AbortTransaction(newBackgroundContext(ctx))
			return nil, ctx.Err()
		}

	CommitLoop:
		for {
			err = s.CommitTransaction(newBackgroundContext(ctx))
			
			if err == nil {
				return res, nil
			}

			select {
			case <-timeout.C:
				return res, err
			default:
			}

			if cerr, ok := err.(CommandError); ok {
				if cerr.HasErrorLabel(driver.UnknownTransactionCommitResult) && !cerr.IsMaxTimeMSExpiredError() {
					continue
				}
				if cerr.HasErrorLabel(driver.TransientTransactionError) {
					break CommitLoop
				}
			}
			return res, err
		}
	}
}


func (s *sessionImpl) StartTransaction(opts ...*options.TransactionOptions) error {
	err := s.clientSession.CheckStartTransaction()
	if err != nil {
		return err
	}

	s.didCommitAfterStart = false

	topts := options.MergeTransactionOptions(opts...)
	coreOpts := &session.TransactionOptions{
		ReadConcern:    topts.ReadConcern,
		ReadPreference: topts.ReadPreference,
		WriteConcern:   topts.WriteConcern,
		MaxCommitTime:  topts.MaxCommitTime,
	}

	return s.clientSession.StartTransaction(coreOpts)
}


func (s *sessionImpl) AbortTransaction(ctx context.Context) error {
	err := s.clientSession.CheckAbortTransaction()
	if err != nil {
		return err
	}

	
	if s.clientSession.TransactionStarting() || s.didCommitAfterStart {
		return s.clientSession.AbortTransaction()
	}

	selector := makePinnedSelector(s.clientSession, description.WriteSelector())

	s.clientSession.Aborting = true
	_ = operation.NewAbortTransaction().Session(s.clientSession).ClusterClock(s.client.clock).Database("admin").
		Deployment(s.deployment).WriteConcern(s.clientSession.CurrentWc).ServerSelector(selector).
		Retry(driver.RetryOncePerCommand).CommandMonitor(s.client.monitor).
		RecoveryToken(bsoncore.Document(s.clientSession.RecoveryToken)).ServerAPI(s.client.serverAPI).
		Authenticator(s.client.authenticator).Execute(ctx)

	s.clientSession.Aborting = false
	_ = s.clientSession.AbortTransaction()

	return nil
}


func (s *sessionImpl) CommitTransaction(ctx context.Context) error {
	err := s.clientSession.CheckCommitTransaction()
	if err != nil {
		return err
	}

	
	if s.clientSession.TransactionStarting() || s.didCommitAfterStart {
		s.didCommitAfterStart = true
		return s.clientSession.CommitTransaction()
	}

	if s.clientSession.TransactionCommitted() {
		s.clientSession.RetryingCommit = true
	}

	selector := makePinnedSelector(s.clientSession, description.WriteSelector())

	s.clientSession.Committing = true
	op := operation.NewCommitTransaction().
		Session(s.clientSession).ClusterClock(s.client.clock).Database("admin").Deployment(s.deployment).
		WriteConcern(s.clientSession.CurrentWc).ServerSelector(selector).Retry(driver.RetryOncePerCommand).
		CommandMonitor(s.client.monitor).RecoveryToken(bsoncore.Document(s.clientSession.RecoveryToken)).
		ServerAPI(s.client.serverAPI).MaxTime(s.clientSession.CurrentMct).Authenticator(s.client.authenticator)

	err = op.Execute(ctx)
	
	
	if IsTimeout(err) {
		return replaceErrors(err)
	}
	s.clientSession.Committing = false
	commitErr := s.clientSession.CommitTransaction()

	
	s.clientSession.UpdateCommitTransactionWriteConcern()

	if err != nil {
		return replaceErrors(err)
	}
	return commitErr
}


func (s *sessionImpl) ClusterTime() bson.Raw {
	return s.clientSession.ClusterTime
}


func (s *sessionImpl) AdvanceClusterTime(d bson.Raw) error {
	return s.clientSession.AdvanceClusterTime(d)
}


func (s *sessionImpl) OperationTime() *primitive.Timestamp {
	return s.clientSession.OperationTime
}


func (s *sessionImpl) AdvanceOperationTime(ts *primitive.Timestamp) error {
	return s.clientSession.AdvanceOperationTime(ts)
}


func (s *sessionImpl) Client() *Client {
	return s.client
}


func (*sessionImpl) session() {
}



func sessionFromContext(ctx context.Context) *session.Client {
	s := ctx.Value(sessionKey{})
	if ses, ok := s.(*sessionImpl); ses != nil && ok {
		return ses.clientSession
	}

	return nil
}
