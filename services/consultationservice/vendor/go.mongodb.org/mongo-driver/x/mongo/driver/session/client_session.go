





package session 

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


var ErrSessionEnded = errors.New("ended session was used")


var ErrNoTransactStarted = errors.New("no transaction started")


var ErrTransactInProgress = errors.New("transaction already in progress")


var ErrAbortAfterCommit = errors.New("cannot call abortTransaction after calling commitTransaction")


var ErrAbortTwice = errors.New("cannot call abortTransaction twice")


var ErrCommitAfterAbort = errors.New("cannot call commitTransaction after calling abortTransaction")


var ErrUnackWCUnsupported = errors.New("transactions do not support unacknowledged write concerns")


var ErrSnapshotTransaction = errors.New("transactions are not supported in snapshot sessions")


type TransactionState uint8


const (
	None TransactionState = iota
	Starting
	InProgress
	Committed
	Aborted
)


func (s TransactionState) String() string {
	switch s {
	case None:
		return "none"
	case Starting:
		return "starting"
	case InProgress:
		return "in progress"
	case Committed:
		return "committed"
	case Aborted:
		return "aborted"
	default:
		return "unknown"
	}
}




type LoadBalancedTransactionConnection interface {
	
	WriteWireMessage(context.Context, []byte) error
	ReadWireMessage(ctx context.Context) ([]byte, error)
	Description() description.Server
	Close() error
	ID() string
	ServerConnectionID() *int64
	DriverConnectionID() uint64 
	Address() address.Address
	Stale() bool
	OIDCTokenGenID() uint64
	SetOIDCTokenGenID(uint64)

	
	PinToCursor() error
	PinToTransaction() error
	UnpinFromCursor() error
	UnpinFromTransaction() error
}


type Client struct {
	*Server
	ClientID       uuid.UUID
	ClusterTime    bson.Raw
	Consistent     bool 
	OperationTime  *primitive.Timestamp
	IsImplicit     bool
	Terminated     bool
	RetryingCommit bool
	Committing     bool
	Aborting       bool
	RetryWrite     bool
	RetryRead      bool
	Snapshot       bool

	
	
	CurrentRc  *readconcern.ReadConcern
	CurrentRp  *readpref.ReadPref
	CurrentWc  *writeconcern.WriteConcern
	CurrentMct *time.Duration

	
	transactionRc            *readconcern.ReadConcern
	transactionRp            *readpref.ReadPref
	transactionWc            *writeconcern.WriteConcern
	transactionMaxCommitTime *time.Duration

	pool             *Pool
	TransactionState TransactionState
	PinnedServer     *description.Server
	RecoveryToken    bson.Raw
	PinnedConnection LoadBalancedTransactionConnection
	SnapshotTime     *primitive.Timestamp
}

func getClusterTime(clusterTime bson.Raw) (uint32, uint32) {
	if clusterTime == nil {
		return 0, 0
	}

	clusterTimeVal, err := clusterTime.LookupErr("$clusterTime")
	if err != nil {
		return 0, 0
	}

	timestampVal, err := bson.Raw(clusterTimeVal.Value).LookupErr("clusterTime")
	if err != nil {
		return 0, 0
	}

	return timestampVal.Timestamp()
}


func MaxClusterTime(ct1, ct2 bson.Raw) bson.Raw {
	epoch1, ord1 := getClusterTime(ct1)
	epoch2, ord2 := getClusterTime(ct2)

	switch {
	case epoch1 > epoch2:
		return ct1
	case epoch1 < epoch2:
		return ct2
	case ord1 > ord2:
		return ct1
	case ord1 < ord2:
		return ct2
	}

	return ct1
}


func NewImplicitClientSession(pool *Pool, clientID uuid.UUID) *Client {
	
	
	

	return &Client{
		pool:       pool,
		ClientID:   clientID,
		IsImplicit: true,
	}
}


func NewClientSession(pool *Pool, clientID uuid.UUID, opts ...*ClientOptions) (*Client, error) {
	c := &Client{
		pool:     pool,
		ClientID: clientID,
	}

	mergedOpts := mergeClientOptions(opts...)
	if mergedOpts.DefaultReadPreference != nil {
		c.transactionRp = mergedOpts.DefaultReadPreference
	}
	if mergedOpts.DefaultReadConcern != nil {
		c.transactionRc = mergedOpts.DefaultReadConcern
	}
	if mergedOpts.DefaultWriteConcern != nil {
		c.transactionWc = mergedOpts.DefaultWriteConcern
	}
	if mergedOpts.DefaultMaxCommitTime != nil {
		c.transactionMaxCommitTime = mergedOpts.DefaultMaxCommitTime
	}
	if mergedOpts.Snapshot != nil {
		c.Snapshot = *mergedOpts.Snapshot
	}

	
	
	
	c.Consistent = !c.Snapshot
	if mergedOpts.CausalConsistency != nil {
		c.Consistent = *mergedOpts.CausalConsistency
	}

	if c.Consistent && c.Snapshot {
		return nil, errors.New("causal consistency and snapshot cannot both be set for a session")
	}

	if err := c.SetServer(); err != nil {
		return nil, err
	}

	return c, nil
}


func (c *Client) SetServer() error {
	var err error
	c.Server, err = c.pool.GetSession()
	return err
}


func (c *Client) AdvanceClusterTime(clusterTime bson.Raw) error {
	if c.Terminated {
		return ErrSessionEnded
	}
	c.ClusterTime = MaxClusterTime(c.ClusterTime, clusterTime)
	return nil
}


func (c *Client) AdvanceOperationTime(opTime *primitive.Timestamp) error {
	if c.Terminated {
		return ErrSessionEnded
	}

	if c.OperationTime == nil {
		c.OperationTime = opTime
		return nil
	}

	if opTime.T > c.OperationTime.T {
		c.OperationTime = opTime
	} else if (opTime.T == c.OperationTime.T) && (opTime.I > c.OperationTime.I) {
		c.OperationTime = opTime
	}

	return nil
}




func (c *Client) UpdateUseTime() error {
	if c.Terminated {
		return ErrSessionEnded
	}
	c.updateUseTime()
	return nil
}


func (c *Client) UpdateRecoveryToken(response bson.Raw) {
	if c == nil {
		return
	}

	token, err := response.LookupErr("recoveryToken")
	if err != nil {
		return
	}

	c.RecoveryToken = token.Document()
}


func (c *Client) UpdateSnapshotTime(response bsoncore.Document) {
	if c == nil {
		return
	}

	subDoc := response
	if cur, ok := response.Lookup("cursor").DocumentOK(); ok {
		subDoc = cur
	}

	ssTimeElem, err := subDoc.LookupErr("atClusterTime")
	if err != nil {
		
		return
	}

	t, i := ssTimeElem.Timestamp()
	c.SnapshotTime = &primitive.Timestamp{
		T: t,
		I: i,
	}
}


func (c *Client) ClearPinnedResources() error {
	if c == nil {
		return nil
	}

	c.PinnedServer = nil
	if c.PinnedConnection != nil {
		if err := c.PinnedConnection.UnpinFromTransaction(); err != nil {
			return err
		}
		if err := c.PinnedConnection.Close(); err != nil {
			return err
		}
	}
	c.PinnedConnection = nil
	return nil
}




func (c *Client) unpinConnection() error {
	if c == nil || c.PinnedConnection == nil {
		return nil
	}

	err := c.PinnedConnection.UnpinFromTransaction()
	closeErr := c.PinnedConnection.Close()
	if err == nil && closeErr != nil {
		err = closeErr
	}
	c.PinnedConnection = nil
	return err
}


func (c *Client) EndSession() {
	if c.Terminated {
		return
	}
	c.Terminated = true

	
	
	
	
	_ = c.unpinConnection()
	c.pool.ReturnSession(c.Server)
}


func (c *Client) TransactionInProgress() bool {
	return c.TransactionState == InProgress
}


func (c *Client) TransactionStarting() bool {
	return c.TransactionState == Starting
}



func (c *Client) TransactionRunning() bool {
	return c != nil && (c.TransactionState == Starting || c.TransactionState == InProgress)
}


func (c *Client) TransactionCommitted() bool {
	return c.TransactionState == Committed
}



func (c *Client) CheckStartTransaction() error {
	if c.TransactionState == InProgress || c.TransactionState == Starting {
		return ErrTransactInProgress
	}
	if c.Snapshot {
		return ErrSnapshotTransaction
	}
	return nil
}



func (c *Client) StartTransaction(opts *TransactionOptions) error {
	err := c.CheckStartTransaction()
	if err != nil {
		return err
	}

	c.IncrementTxnNumber()
	c.RetryingCommit = false

	if opts != nil {
		c.CurrentRc = opts.ReadConcern
		c.CurrentRp = opts.ReadPreference
		c.CurrentWc = opts.WriteConcern
		c.CurrentMct = opts.MaxCommitTime
	}

	if c.CurrentRc == nil {
		c.CurrentRc = c.transactionRc
	}

	if c.CurrentRp == nil {
		c.CurrentRp = c.transactionRp
	}

	if c.CurrentWc == nil {
		c.CurrentWc = c.transactionWc
	}

	if c.CurrentMct == nil {
		c.CurrentMct = c.transactionMaxCommitTime
	}

	if !writeconcern.AckWrite(c.CurrentWc) {
		_ = c.clearTransactionOpts()
		return ErrUnackWCUnsupported
	}

	c.TransactionState = Starting
	return c.ClearPinnedResources()
}



func (c *Client) CheckCommitTransaction() error {
	if c.TransactionState == None {
		return ErrNoTransactStarted
	} else if c.TransactionState == Aborted {
		return ErrCommitAfterAbort
	}
	return nil
}



func (c *Client) CommitTransaction() error {
	err := c.CheckCommitTransaction()
	if err != nil {
		return err
	}
	c.TransactionState = Committed
	return nil
}




func (c *Client) UpdateCommitTransactionWriteConcern() {
	wc := c.CurrentWc
	timeout := 10 * time.Second
	if wc != nil && wc.GetWTimeout() != 0 {
		timeout = wc.GetWTimeout()
	}
	c.CurrentWc = wc.WithOptions(writeconcern.WMajority(), writeconcern.WTimeout(timeout))
}



func (c *Client) CheckAbortTransaction() error {
	switch {
	case c.TransactionState == None:
		return ErrNoTransactStarted
	case c.TransactionState == Committed:
		return ErrAbortAfterCommit
	case c.TransactionState == Aborted:
		return ErrAbortTwice
	}
	return nil
}



func (c *Client) AbortTransaction() error {
	err := c.CheckAbortTransaction()
	if err != nil {
		return err
	}
	c.TransactionState = Aborted
	return c.clearTransactionOpts()
}



func (c *Client) StartCommand() error {
	if c == nil {
		return nil
	}

	
	
	if !c.TransactionRunning() && !c.Committing && !c.Aborting {
		return c.ClearPinnedResources()
	}
	return nil
}



func (c *Client) ApplyCommand(desc description.Server) error {
	if c.Committing {
		
		return nil
	}
	if c.TransactionState == Starting {
		c.TransactionState = InProgress
		
		if desc.Kind == description.Mongos {
			c.PinnedServer = &desc
		}
	} else if c.TransactionState == Committed || c.TransactionState == Aborted {
		c.TransactionState = None
		return c.clearTransactionOpts()
	}

	return nil
}

func (c *Client) clearTransactionOpts() error {
	c.RetryingCommit = false
	c.Aborting = false
	c.Committing = false
	c.CurrentWc = nil
	c.CurrentRp = nil
	c.CurrentRc = nil
	c.RecoveryToken = nil

	return c.ClearPinnedResources()
}
