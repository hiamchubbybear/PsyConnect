





package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/handshake"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const defaultLocalThreshold = 15 * time.Millisecond

var (
	
	ErrNoDocCommandResponse = errors.New("command returned no documents")
	
	ErrMultiDocCommandResponse = errors.New("command returned multiple documents")
	
	ErrReplyDocumentMismatch = errors.New("number of documents returned does not match numberReturned field")
	
	ErrNonPrimaryReadPref = errors.New("read preference in a transaction must be primary")
	
	errDatabaseNameEmpty = errors.New("database name cannot be empty")
)

const (
	
	cryptMaxBsonObjectSize uint32 = 2097152
	
	cryptMinWireVersion int32 = 8
	
	readSnapshotMinWireVersion int32 = 13
)


type RetryablePoolError interface {
	Retryable() bool
}


type labeledError interface {
	error
	HasErrorLabel(string) bool
}



type InvalidOperationError struct{ MissingField string }

func (err InvalidOperationError) Error() string {
	return "the " + err.MissingField + " field must be set on Operation"
}



type opReply struct {
	responseFlags wiremessage.ReplyFlag
	cursorID      int64
	startingFrom  int32
	numReturned   int32
	documents     []bsoncore.Document
	err           error
}


type startedInformation struct {
	cmd                      bsoncore.Document
	requestID                int32
	cmdName                  string
	documentSequenceIncluded bool
	connID                   string
	driverConnectionID       uint64 
	serverConnID             *int64
	redacted                 bool
	serviceID                *primitive.ObjectID
	serverAddress            address.Address
}


type finishedInformation struct {
	cmdName            string
	requestID          int32
	response           bsoncore.Document
	cmdErr             error
	connID             string
	driverConnectionID uint64 
	serverConnID       *int64
	redacted           bool
	serviceID          *primitive.ObjectID
	serverAddress      address.Address
	duration           time.Duration
}




func convertInt64PtrToInt32Ptr(i64 *int64) *int32 {
	if i64 == nil {
		return nil
	}

	if *i64 > math.MaxInt32 || *i64 < math.MinInt32 {
		return nil
	}

	i32 := int32(*i64)
	return &i32
}







func (info finishedInformation) success() bool {
	if _, ok := info.cmdErr.(WriteCommandError); ok {
		return true
	}

	return info.cmdErr == nil
}


type ResponseInfo struct {
	ServerResponse        bsoncore.Document
	Server                Server
	Connection            Connection
	ConnectionDescription description.Server
	CurrentIndex          int
}

func redactStartedInformationCmd(op Operation, info startedInformation) bson.Raw {
	var cmdCopy bson.Raw

	
	
	
	if !info.redacted {
		cmdCopy = make([]byte, len(info.cmd))
		copy(cmdCopy, info.cmd)

		if info.documentSequenceIncluded {
			
			cmdCopy = cmdCopy[:len(info.cmd)-1]
			cmdCopy = op.addBatchArray(cmdCopy)

			
			cmdCopy, _ = bsoncore.AppendDocumentEnd(cmdCopy, 0)
		}
	}

	return cmdCopy
}

func redactFinishedInformationResponse(info finishedInformation) bson.Raw {
	if !info.redacted {
		return bson.Raw(info.response)
	}

	return bson.Raw{}
}











type Operation struct {
	
	
	
	
	CommandFn func(dst []byte, desc description.SelectedServer) ([]byte, error)

	
	Database string

	
	
	
	
	Deployment Deployment

	
	
	
	ProcessResponseFn func(ResponseInfo) error

	
	
	
	Selector description.ServerSelector

	
	
	ReadPreference *readpref.ReadPref

	
	
	
	ReadConcern *readconcern.ReadConcern

	
	
	MinimumReadConcernWireVersion int32

	
	
	
	WriteConcern *writeconcern.WriteConcern

	
	
	MinimumWriteConcernWireVersion int32

	
	
	
	
	
	
	Client *session.Client

	
	
	
	Clock *session.ClusterClock

	
	
	
	
	
	RetryMode *RetryMode

	
	
	
	Type Type

	
	
	
	
	Batches *Batches

	
	
	
	Legacy LegacyOperationKind

	
	
	CommandMonitor *event.CommandMonitor

	
	Crypt Crypt

	
	ServerAPI *ServerAPIOptions

	
	
	IsOutputAggregate bool

	
	MaxTime *time.Duration

	
	
	Timeout *time.Duration

	Logger *logger.Logger

	
	
	Name string

	
	
	
	OmitCSOTMaxTimeMS bool

	
	
	Authenticator Authenticator

	
	
	
	
	omitReadPreference bool
}


func (op Operation) shouldEncrypt() bool {
	return op.Crypt != nil && !op.Crypt.BypassAutoEncryption()
}








func filterDeprioritizedServers(candidates, deprioritized []description.Server) []description.Server {
	if len(deprioritized) == 0 {
		return candidates
	}

	dpaSet := make(map[address.Address]*description.Server)
	for i, srv := range deprioritized {
		dpaSet[srv.Addr] = &deprioritized[i]
	}

	allowed := []description.Server{}

	
	
	for _, candidate := range candidates {
		if srv, ok := dpaSet[candidate.Addr]; !ok || !srv.Equal(candidate) {
			allowed = append(allowed, candidate)
		}
	}

	
	
	
	if len(allowed) == 0 {
		return candidates
	}

	return allowed
}




type opServerSelector struct {
	selector             description.ServerSelector
	deprioritizedServers []description.Server
}



func (oss *opServerSelector) SelectServer(
	topo description.Topology,
	candidates []description.Server,
) ([]description.Server, error) {
	selectedServers, err := oss.selector.SelectServer(topo, candidates)
	if err != nil {
		return nil, err
	}

	filteredServers := filterDeprioritizedServers(selectedServers, oss.deprioritizedServers)

	return filteredServers, nil
}


func (op Operation) selectServer(
	ctx context.Context,
	requestID int32,
	deprioritized []description.Server,
) (Server, error) {
	if err := op.Validate(); err != nil {
		return nil, err
	}

	selector := op.Selector
	if selector == nil {
		rp := op.ReadPreference
		if rp == nil {
			rp = readpref.Primary()
		}
		selector = description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(rp),
			description.LatencySelector(defaultLocalThreshold),
		})
	}

	oss := &opServerSelector{
		selector:             selector,
		deprioritizedServers: deprioritized,
	}

	ctx = logger.WithOperationName(ctx, op.Name)
	ctx = logger.WithOperationID(ctx, requestID)

	return op.Deployment.SelectServer(ctx, oss)
}


func (op Operation) getServerAndConnection(
	ctx context.Context,
	requestID int32,
	deprioritized []description.Server,
) (Server, Connection, error) {
	server, err := op.selectServer(ctx, requestID, deprioritized)
	if err != nil {
		if op.Client != nil &&
			!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning() {
			err = Error{
				Message: err.Error(),
				Labels:  []string{TransientTransactionError},
				Wrapped: err,
			}
		}
		return nil, nil, err
	}

	
	
	if op.Client != nil && op.Client.PinnedConnection != nil {
		return server, op.Client.PinnedConnection, nil
	}

	
	conn, err := server.Connection(ctx)
	if err != nil {
		return nil, nil, err
	}

	
	if conn.Description().LoadBalanced() && op.Client != nil && op.Client.TransactionStarting() {
		pinnedConn, ok := conn.(PinnedConnection)
		if !ok {
			
			_ = conn.Close()
			return nil, nil, fmt.Errorf("expected Connection used to start a transaction to be a PinnedConnection, but got %T", conn)
		}
		if err := pinnedConn.PinToTransaction(); err != nil {
			
			_ = conn.Close()
			return nil, nil, fmt.Errorf("error incrementing connection reference count when starting a transaction: %w", err)
		}
		op.Client.PinnedConnection = pinnedConn
	}

	return server, conn, nil
}


func (op Operation) Validate() error {
	if op.CommandFn == nil {
		return InvalidOperationError{MissingField: "CommandFn"}
	}
	if op.Deployment == nil {
		return InvalidOperationError{MissingField: "Deployment"}
	}
	if op.Database == "" {
		return errDatabaseNameEmpty
	}
	if op.Client != nil && !writeconcern.AckWrite(op.WriteConcern) {
		return errors.New("session provided for an unacknowledged write")
	}
	return nil
}

var memoryPool = sync.Pool{
	New: func() interface{} {
		
		b := make([]byte, 1024)
		
		return &b
	},
}


func (op Operation) Execute(ctx context.Context) error {
	err := op.Validate()
	if err != nil {
		return err
	}

	
	
	if op.Timeout != nil && !csot.IsTimeoutContext(ctx) {
		newCtx, cancelFunc := csot.MakeTimeoutContext(ctx, *op.Timeout)
		
		ctx = newCtx
		
		defer cancelFunc()
	}

	if op.Client != nil {
		if err := op.Client.StartCommand(); err != nil {
			return err
		}
	}

	var retries int
	if op.RetryMode != nil {
		switch op.Type {
		case Write:
			if op.Client == nil {
				break
			}
			switch *op.RetryMode {
			case RetryOnce, RetryOncePerCommand:
				retries = 1
			case RetryContext:
				retries = -1
			}
		case Read:
			switch *op.RetryMode {
			case RetryOnce, RetryOncePerCommand:
				retries = 1
			case RetryContext:
				retries = -1
			}
		}
	}
	
	
	retryEnabled := op.RetryMode != nil && op.RetryMode.Enabled()
	if csot.IsTimeoutContext(ctx) && retryEnabled {
		retries = -1
	}

	var srvr Server
	var conn Connection
	var res bsoncore.Document
	var operationErr WriteCommandError
	var prevErr error
	var prevIndefiniteErr error
	batching := op.Batches.Valid()
	retrySupported := false
	first := true
	currIndex := 0

	
	
	
	var deprioritizedServers []description.Server

	
	
	resetForRetry := func(err error) {
		retries--
		prevErr = err

		
		
		if err, ok := err.(labeledError); ok {
			
			
			
			
			if prevIndefiniteErr == nil {
				prevIndefiniteErr = err
			}

			
			
			if !err.HasErrorLabel(NoWritesPerformed) && err.HasErrorLabel(RetryableWriteError) {
				prevIndefiniteErr = err
			}
		}

		
		
		if conn != nil {
			
			
			if desc := conn.Description; desc != nil && op.Deployment.Kind() == description.Sharded {
				deprioritizedServers = []description.Server{conn.Description()}
			}

			conn.Close()
		}

		
		srvr = nil
		conn = nil
	}

	wm := memoryPool.Get().(*[]byte)
	defer func() {
		
		
		
		
		
		
		
		
		if c := cap(*wm); c < 16*1024*1024 && c/2 < len(*wm) {
			memoryPool.Put(wm)
		}
	}()
	for {
		
		
		
		if errors.Is(prevErr, context.Canceled) || errors.Is(prevErr, context.DeadlineExceeded) {
			return prevErr
		}

		requestID := wiremessage.NextRequestID()

		
		if srvr == nil || conn == nil {
			srvr, conn, err = op.getServerAndConnection(ctx, requestID, deprioritizedServers)
			if err != nil {
				
				
				
				if rerr, ok := err.(RetryablePoolError); ok && rerr.Retryable() && retries != 0 {
					resetForRetry(err)
					continue
				}

				
				
				if prevErr != nil {
					return prevErr
				}
				return err
			}
			defer conn.Close()

			
			
			
			if op.Client != nil && op.Client.Server == nil && op.Client.IsImplicit {
				if op.Client.Terminated {
					return fmt.Errorf("unexpected nil session for a terminated implicit session")
				}
				if err := op.Client.SetServer(); err != nil {
					return err
				}
			}
		}

		
		if first {
			
			
			
			
			
			
			
			retrySupported = op.retryable(conn.Description())

			
			
			
			
			
			
			if retrySupported && op.RetryMode != nil && op.Type == Write && op.Client != nil {
				op.Client.RetryWrite = false
				if op.RetryMode.Enabled() {
					op.Client.RetryWrite = true
					if !op.Client.Committing && !op.Client.Aborting {
						op.Client.IncrementTxnNumber()
					}
				}
			}

			first = false
		}

		maxTimeMS, err := op.calculateMaxTimeMS(ctx, srvr.RTTMonitor())
		if err != nil {
			return err
		}

		
		
		if conn.Description().IsCryptd {
			maxTimeMS = 0
		}

		desc := description.SelectedServer{Server: conn.Description(), Kind: op.Deployment.Kind()}

		if batching {
			targetBatchSize := desc.MaxDocumentSize
			maxDocSize := desc.MaxDocumentSize
			if op.shouldEncrypt() {
				
				
				
				
				targetBatchSize = cryptMaxBsonObjectSize
			}

			err = op.Batches.AdvanceBatch(int(desc.MaxBatchCount), int(targetBatchSize), int(maxDocSize))
			if err != nil {
				
				return err
			}
		}

		var startedInfo startedInformation
		*wm, startedInfo, err = op.createWireMessage(ctx, maxTimeMS, (*wm)[:0], desc, conn, requestID)

		if err != nil {
			return err
		}

		
		startedInfo.connID = conn.ID()
		startedInfo.driverConnectionID = conn.DriverConnectionID()
		startedInfo.cmdName = op.getCommandName(startedInfo.cmd)

		
		
		
		
		if startedInfo.cmdName != op.Name {
			op.Name = startedInfo.cmdName
		}

		startedInfo.redacted = op.redactCommand(startedInfo.cmdName, startedInfo.cmd)
		startedInfo.serviceID = conn.Description().ServiceID
		startedInfo.serverConnID = conn.ServerConnectionID()
		startedInfo.serverAddress = conn.Description().Addr

		op.publishStartedEvent(ctx, startedInfo)

		
		moreToCome := wiremessage.IsMsgMoreToCome(*wm)

		
		if compressor, ok := conn.(Compressor); ok && op.canCompress(startedInfo.cmdName) {
			b := memoryPool.Get().(*[]byte)
			*b, err = compressor.CompressWireMessage(*wm, (*b)[:0])
			memoryPool.Put(wm)
			wm = b
			if err != nil {
				return err
			}
		}

		finishedInfo := finishedInformation{
			cmdName:            startedInfo.cmdName,
			driverConnectionID: startedInfo.driverConnectionID,
			requestID:          startedInfo.requestID,
			connID:             startedInfo.connID,
			serverConnID:       startedInfo.serverConnID,
			redacted:           startedInfo.redacted,
			serviceID:          startedInfo.serviceID,
			serverAddress:      desc.Server.Addr,
		}

		startedTime := time.Now()

		
		
		
		if ctx.Err() != nil {
			err = ctx.Err()
		} else if deadline, ok := ctx.Deadline(); ok {
			if csot.IsTimeoutContext(ctx) && time.Now().Add(srvr.RTTMonitor().P90()).After(deadline) {
				err = fmt.Errorf(
					"remaining time %v until context deadline is less than 90th percentile network round-trip time: %w\n%v",
					time.Until(deadline),
					ErrDeadlineWouldBeExceeded,
					srvr.RTTMonitor().Stats())
			} else if time.Now().Add(srvr.RTTMonitor().Min()).After(deadline) {
				err = context.DeadlineExceeded
			}
		}

		if err == nil {
			
			
			roundTrip := op.roundTrip
			if moreToCome {
				roundTrip = op.moreToComeRoundTrip
			}
			res, err = roundTrip(ctx, conn, *wm)

			if ep, ok := srvr.(ErrorProcessor); ok {
				_ = ep.ProcessError(err, conn)
			}
		}

		finishedInfo.response = res
		finishedInfo.cmdErr = err
		finishedInfo.duration = time.Since(startedTime)

		op.publishFinishedEvent(ctx, finishedInfo)

		
		
		var prevIndefiniteErrIsSet bool

		
		
	checkError:
		var perr error
		switch tt := err.(type) {
		case WriteCommandError:
			if e := err.(WriteCommandError); retrySupported && op.Type == Write && e.UnsupportedStorageEngine() {
				return ErrUnsupportedStorageEngine
			}

			connDesc := conn.Description()
			retryableErr := tt.Retryable(connDesc.WireVersion)
			preRetryWriteLabelVersion := connDesc.WireVersion != nil && connDesc.WireVersion.Max < 9
			inTransaction := op.Client != nil &&
				!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning()
			
			
			if retryableErr && preRetryWriteLabelVersion && retryEnabled && !inTransaction {
				tt.Labels = append(tt.Labels, RetryableWriteError)
			}

			
			
			
			if retrySupported && retryableErr && retries != 0 {
				if op.Client != nil && op.Client.Committing {
					
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				}
				resetForRetry(tt)
				continue
			}

			
			
			
			if tt.HasErrorLabel(NoWritesPerformed) && !prevIndefiniteErrIsSet {
				err = prevIndefiniteErr
				prevIndefiniteErrIsSet = true

				goto checkError
			}

			
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				_ = op.ProcessResponseFn(info)
			}

			if batching && len(tt.WriteErrors) > 0 && currIndex > 0 {
				for i := range tt.WriteErrors {
					tt.WriteErrors[i].Index += int64(currIndex)
				}
			}

			
			
			if batching && (op.Batches.Ordered == nil || *op.Batches.Ordered) && len(tt.WriteErrors) > 0 {
				return tt
			}
			if op.Client != nil && op.Client.Committing && tt.WriteConcernError != nil {
				
				err := Error{
					Name:    tt.WriteConcernError.Name,
					Code:    int32(tt.WriteConcernError.Code),
					Message: tt.WriteConcernError.Message,
					Labels:  tt.Labels,
					Raw:     tt.Raw,
				}
				
				
				if err.Code != unknownReplWriteConcernCode && err.Code != unsatisfiableWriteConcernCode {
					err.Labels = append(err.Labels, UnknownTransactionCommitResult)
				}
				if retryableErr && retryEnabled {
					err.Labels = append(err.Labels, RetryableWriteError)
				}
				return err
			}
			operationErr.WriteConcernError = tt.WriteConcernError
			operationErr.WriteErrors = append(operationErr.WriteErrors, tt.WriteErrors...)
			operationErr.Labels = tt.Labels
			operationErr.Raw = tt.Raw
		case Error:
			
			
			if tt.Code == 391 {
				if op.Authenticator != nil {
					cfg := AuthConfig{
						Description:  conn.Description(),
						Connection:   conn,
						ClusterClock: op.Clock,
						ServerAPI:    op.ServerAPI,
					}
					if err := op.Authenticator.Reauth(ctx, &cfg); err != nil {
						return fmt.Errorf("error reauthenticating: %w", err)
					}
					if op.Client != nil && op.Client.Committing {
						
						op.Client.UpdateCommitTransactionWriteConcern()
						op.WriteConcern = op.Client.CurrentWc
					}
					resetForRetry(tt)
					continue
				}
			}
			if tt.HasErrorLabel(TransientTransactionError) || tt.HasErrorLabel(UnknownTransactionCommitResult) {
				if err := op.Client.ClearPinnedResources(); err != nil {
					return err
				}
			}

			if e := err.(Error); retrySupported && op.Type == Write && e.UnsupportedStorageEngine() {
				return ErrUnsupportedStorageEngine
			}

			connDesc := conn.Description()
			var retryableErr bool
			if op.Type == Write {
				retryableErr = tt.RetryableWrite(connDesc.WireVersion)
				preRetryWriteLabelVersion := connDesc.WireVersion != nil && connDesc.WireVersion.Max < 9
				inTransaction := op.Client != nil &&
					!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning()
				
				
				if retryEnabled && !inTransaction &&
					(tt.HasErrorLabel(NetworkError) || (retryableErr && preRetryWriteLabelVersion)) {
					tt.Labels = append(tt.Labels, RetryableWriteError)
				}
			} else {
				retryableErr = tt.RetryableRead()
			}

			
			
			
			if retrySupported && retryableErr && retries != 0 {
				if op.Client != nil && op.Client.Committing {
					
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				}
				resetForRetry(tt)
				continue
			}

			
			
			
			if tt.HasErrorLabel(NoWritesPerformed) && !prevIndefiniteErrIsSet {
				err = prevIndefiniteErr
				prevIndefiniteErrIsSet = true

				goto checkError
			}

			
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				_ = op.ProcessResponseFn(info)
			}

			if op.Client != nil && op.Client.Committing && (retryableErr || tt.Code == 50) {
				
				tt.Labels = append(tt.Labels, UnknownTransactionCommitResult)
			}
			return tt
		case nil:
			if moreToCome {
				return ErrUnacknowledgedWrite
			}
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				perr = op.ProcessResponseFn(info)
			}
			if perr != nil {
				return perr
			}
		default:
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				_ = op.ProcessResponseFn(info)
			}
			return err
		}

		
		
		
		if batching && len(op.Batches.Documents) > 0 {
			
			
			
			
			if retrySupported && op.Client != nil && op.RetryMode != nil {
				if op.RetryMode.Enabled() {
					op.Client.IncrementTxnNumber()
				}
				
				
				if *op.RetryMode == RetryOncePerCommand && !csot.IsTimeoutContext(ctx) {
					retries = 1
				}
			}
			currIndex += len(op.Batches.Current)
			op.Batches.ClearBatch()
			continue
		}
		break
	}
	if len(operationErr.WriteErrors) > 0 || operationErr.WriteConcernError != nil {
		return operationErr
	}
	return nil
}



func (op Operation) retryable(desc description.Server) bool {
	switch op.Type {
	case Write:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) {
			return true
		}
		if retryWritesSupported(desc) &&
			op.Client != nil && !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) &&
			writeconcern.AckWrite(op.WriteConcern) {
			return true
		}
	case Read:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) {
			return true
		}
		if op.Client == nil || !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) {
			return true
		}
	}
	return false
}



func (op Operation) roundTrip(ctx context.Context, conn Connection, wm []byte) ([]byte, error) {
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil {
		return nil, op.networkError(err)
	}
	return op.readWireMessage(ctx, conn)
}

func (op Operation) readWireMessage(ctx context.Context, conn Connection) (result []byte, err error) {
	wm, err := conn.ReadWireMessage(ctx)
	if err != nil {
		return nil, op.networkError(err)
	}

	
	
	if streamer, ok := conn.(StreamerConnection); ok {
		streamer.SetStreaming(wiremessage.IsMsgMoreToCome(wm))
	}

	length, _, _, opcode, rem, ok := wiremessage.ReadHeader(wm)
	if !ok || len(wm) < int(length) {
		return nil, errors.New("malformed wire message: insufficient bytes")
	}
	if opcode == wiremessage.OpCompressed {
		rawsize := length - 16 
		
		opcode, rem, err = op.decompressWireMessage(rem[:rawsize])
		if err != nil {
			return nil, err
		}
	}

	
	res, err := op.decodeResult(ctx, opcode, rem)
	
	
	op.updateClusterTimes(res)
	op.updateOperationTime(res)
	op.Client.UpdateRecoveryToken(bson.Raw(res))

	
	if op.Name == driverutil.FindOp || op.Name == driverutil.AggregateOp || op.Name == driverutil.DistinctOp {
		op.Client.UpdateSnapshotTime(res)
	}

	if err != nil {
		return res, err
	}

	
	if op.Crypt != nil {
		res, err = op.Crypt.Decrypt(ctx, res)
	}
	return res, err
}




func (op Operation) networkError(err error) error {
	if err == nil {
		return nil
	}

	labels := []string{NetworkError}
	if op.Client != nil {
		op.Client.MarkDirty()
	}
	if op.Client != nil && op.Client.TransactionRunning() && !op.Client.Committing {
		labels = append(labels, TransientTransactionError)
	}
	if op.Client != nil && op.Client.Committing {
		labels = append(labels, UnknownTransactionCommitResult)
	}
	return Error{Message: err.Error(), Labels: labels, Wrapped: err}
}



func (op *Operation) moreToComeRoundTrip(ctx context.Context, conn Connection, wm []byte) (result []byte, err error) {
	err = conn.WriteWireMessage(ctx, wm)
	if err != nil {
		if op.Client != nil {
			op.Client.MarkDirty()
		}
		err = Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}, Wrapped: err}
	}
	return bsoncore.BuildDocument(nil, bsoncore.AppendInt32Element(nil, "ok", 1)), err
}


func (Operation) decompressWireMessage(wm []byte) (wiremessage.OpCode, []byte, error) {
	
	opcode, rem, ok := wiremessage.ReadCompressedOriginalOpCode(wm)
	if !ok {
		return 0, nil, errors.New("malformed OP_COMPRESSED: missing original opcode")
	}
	uncompressedSize, rem, ok := wiremessage.ReadCompressedUncompressedSize(rem)
	if !ok {
		return 0, nil, errors.New("malformed OP_COMPRESSED: missing uncompressed size")
	}
	
	compressorID, rem, ok := wiremessage.ReadCompressedCompressorID(rem)
	if !ok {
		return 0, nil, errors.New("malformed OP_COMPRESSED: missing compressor ID")
	}

	opts := CompressionOpts{
		Compressor:       compressorID,
		UncompressedSize: uncompressedSize,
	}
	uncompressed, err := DecompressPayload(rem, opts)
	if err != nil {
		return 0, nil, err
	}

	return opcode, uncompressed, nil
}

func (op Operation) addBatchArray(dst []byte) []byte {
	aidx, dst := bsoncore.AppendArrayElementStart(dst, op.Batches.Identifier)
	for i, doc := range op.Batches.Current {
		dst = bsoncore.AppendDocumentElement(dst, strconv.Itoa(i), doc)
	}
	dst, _ = bsoncore.AppendArrayEnd(dst, aidx)
	return dst
}

func (op Operation) createLegacyHandshakeWireMessage(
	maxTimeMS uint64,
	dst []byte,
	desc description.SelectedServer,
) ([]byte, startedInformation, error) {
	var info startedInformation
	flags := op.secondaryOK(desc)
	var wmindex int32
	info.requestID = wiremessage.NextRequestID()
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, flags)

	dollarCmd := [...]byte{'.', '$', 'c', 'm', 'd'}

	
	dst = append(dst, op.Database...)
	dst = append(dst, dollarCmd[:]...)
	dst = append(dst, 0x00)
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, -1)

	wrapper := int32(-1)
	rp, err := op.createReadPref(desc, true)
	if err != nil {
		return dst, info, err
	}
	if len(rp) > 0 {
		wrapper, dst = bsoncore.AppendDocumentStart(dst)
		dst = bsoncore.AppendHeader(dst, bsontype.EmbeddedDocument, "$query")
	}
	idx, dst := bsoncore.AppendDocumentStart(dst)
	dst, err = op.CommandFn(dst, desc)
	if err != nil {
		return dst, info, err
	}

	if op.Batches != nil && len(op.Batches.Current) > 0 {
		dst = op.addBatchArray(dst)
	}

	dst, err = op.addReadConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst, err = op.addWriteConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst, err = op.addSession(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst = op.addClusterTime(dst, desc)
	dst = op.addServerAPI(dst)
	
	
	if maxTimeMS > 0 {
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", int64(maxTimeMS))
	}

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)
	
	info.cmd = dst[idx:]

	if len(rp) > 0 {
		var err error
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
		dst, err = bsoncore.AppendDocumentEnd(dst, wrapper)
		if err != nil {
			return dst, info, err
		}
	}

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), info, nil
}

func (op Operation) createMsgWireMessage(
	ctx context.Context,
	maxTimeMS uint64,
	dst []byte,
	desc description.SelectedServer,
	conn Connection,
	requestID int32,
) ([]byte, startedInformation, error) {
	var info startedInformation
	var flags wiremessage.MsgFlag
	var wmindex int32
	
	
	if op.WriteConcern != nil && !writeconcern.AckWrite(op.WriteConcern) && (op.Batches == nil || len(op.Batches.Documents) == 0) {
		flags = wiremessage.MoreToCome
	}
	
	
	if streamer, ok := conn.(StreamerConnection); ok && streamer.SupportsStreaming() {
		flags |= wiremessage.ExhaustAllowed
	}

	info.requestID = requestID
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpMsg)
	dst = wiremessage.AppendMsgFlags(dst, flags)
	
	dst = wiremessage.AppendMsgSectionType(dst, wiremessage.SingleDocument)

	idx, dst := bsoncore.AppendDocumentStart(dst)

	dst, err := op.addCommandFields(ctx, dst, desc)
	if err != nil {
		return dst, info, err
	}
	dst, err = op.addReadConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}
	dst, err = op.addWriteConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}
	dst, err = op.addSession(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst = op.addClusterTime(dst, desc)
	dst = op.addServerAPI(dst)
	
	
	if maxTimeMS > 0 {
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", int64(maxTimeMS))
	}

	dst = bsoncore.AppendStringElement(dst, "$db", op.Database)
	rp, err := op.createReadPref(desc, false)
	if err != nil {
		return dst, info, err
	}
	if len(rp) > 0 {
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
	}

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)
	
	info.cmd = dst[idx:]

	
	
	if !op.shouldEncrypt() && op.Batches != nil && len(op.Batches.Current) > 0 {
		info.documentSequenceIncluded = true
		dst = wiremessage.AppendMsgSectionType(dst, wiremessage.DocumentSequence)
		idx, dst = bsoncore.ReserveLength(dst)

		dst = append(dst, op.Batches.Identifier...)
		dst = append(dst, 0x00)

		for _, doc := range op.Batches.Current {
			dst = append(dst, doc...)
		}

		dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	}

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), info, nil
}



func isLegacyHandshake(op Operation, desc description.SelectedServer) bool {
	isInitialHandshake := desc.WireVersion == nil || desc.WireVersion.Max == 0

	return op.Legacy == LegacyHandshake && isInitialHandshake
}

func (op Operation) createWireMessage(
	ctx context.Context,
	maxTimeMS uint64,
	dst []byte,
	desc description.SelectedServer,
	conn Connection,
	requestID int32,
) ([]byte, startedInformation, error) {
	if isLegacyHandshake(op, desc) {
		return op.createLegacyHandshakeWireMessage(maxTimeMS, dst, desc)
	}

	return op.createMsgWireMessage(ctx, maxTimeMS, dst, desc, conn, requestID)
}



func (op Operation) addCommandFields(ctx context.Context, dst []byte, desc description.SelectedServer) ([]byte, error) {
	if !op.shouldEncrypt() {
		return op.CommandFn(dst, desc)
	}

	if desc.WireVersion.Max < cryptMinWireVersion {
		return dst, errors.New("auto-encryption requires a MongoDB version of 4.2")
	}

	
	cidx, cmdDst := bsoncore.AppendDocumentStart(nil)
	var err error
	cmdDst, err = op.CommandFn(cmdDst, desc)
	if err != nil {
		return dst, err
	}
	
	if op.Batches != nil && len(op.Batches.Current) > 0 {
		cmdDst = op.addBatchArray(cmdDst)
	}
	cmdDst, _ = bsoncore.AppendDocumentEnd(cmdDst, cidx)

	
	encrypted, err := op.Crypt.Encrypt(ctx, op.Database, cmdDst)
	if err != nil {
		return dst, err
	}
	
	dst = append(dst, encrypted[4:len(encrypted)-1]...)
	return dst, nil
}


func (op Operation) addServerAPI(dst []byte) []byte {
	sa := op.ServerAPI
	if sa == nil {
		return dst
	}

	dst = bsoncore.AppendStringElement(dst, "apiVersion", sa.ServerAPIVersion)
	if sa.Strict != nil {
		dst = bsoncore.AppendBooleanElement(dst, "apiStrict", *sa.Strict)
	}
	if sa.DeprecationErrors != nil {
		dst = bsoncore.AppendBooleanElement(dst, "apiDeprecationErrors", *sa.DeprecationErrors)
	}
	return dst
}

func (op Operation) addReadConcern(dst []byte, desc description.SelectedServer) ([]byte, error) {
	if op.MinimumReadConcernWireVersion > 0 && (desc.WireVersion == nil || !desc.WireVersion.Includes(op.MinimumReadConcernWireVersion)) {
		return dst, nil
	}
	rc := op.ReadConcern
	client := op.Client
	
	if client != nil && client.TransactionStarting() && client.CurrentRc != nil {
		rc = client.CurrentRc
	}

	
	if rc == nil && client != nil && client.TransactionStarting() && client.Consistent && client.OperationTime != nil {
		rc = readconcern.New()
	}

	if client != nil && client.Snapshot {
		if desc.WireVersion.Max < readSnapshotMinWireVersion {
			return dst, errors.New("snapshot reads require MongoDB 5.0 or later")
		}
		rc = readconcern.Snapshot()
	}

	if rc == nil {
		return dst, nil
	}

	_, data, err := rc.MarshalBSONValue() 
	if err != nil {
		return dst, err
	}

	if sessionsSupported(desc.WireVersion) && client != nil {
		if client.Consistent && client.OperationTime != nil {
			data = data[:len(data)-1] 
			data = bsoncore.AppendTimestampElement(data, "afterClusterTime", client.OperationTime.T, client.OperationTime.I)
			data, _ = bsoncore.AppendDocumentEnd(data, 0)
		}
		if client.Snapshot && client.SnapshotTime != nil {
			data = data[:len(data)-1] 
			data = bsoncore.AppendTimestampElement(data, "atClusterTime", client.SnapshotTime.T, client.SnapshotTime.I)
			data, _ = bsoncore.AppendDocumentEnd(data, 0)
		}
	}

	if len(data) == bsoncore.EmptyDocumentLength {
		return dst, nil
	}
	return bsoncore.AppendDocumentElement(dst, "readConcern", data), nil
}

func (op Operation) addWriteConcern(dst []byte, desc description.SelectedServer) ([]byte, error) {
	if op.MinimumWriteConcernWireVersion > 0 && (desc.WireVersion == nil || !desc.WireVersion.Includes(op.MinimumWriteConcernWireVersion)) {
		return dst, nil
	}
	wc := op.WriteConcern
	if wc == nil {
		return dst, nil
	}

	t, data, err := wc.MarshalBSONValue()
	if errors.Is(err, writeconcern.ErrEmptyWriteConcern) {
		return dst, nil
	}
	if err != nil {
		return dst, err
	}

	return append(bsoncore.AppendHeader(dst, t, "writeConcern"), data...), nil
}

func (op Operation) addSession(dst []byte, desc description.SelectedServer) ([]byte, error) {
	client := op.Client

	
	
	if client != nil && !client.IsImplicit && desc.SessionTimeoutMinutesPtr == nil {
		return nil, fmt.Errorf("current topology does not support sessions")
	}

	if client == nil || !sessionsSupported(desc.WireVersion) || desc.SessionTimeoutMinutesPtr == nil {
		return dst, nil
	}
	if err := client.UpdateUseTime(); err != nil {
		return dst, err
	}
	dst = bsoncore.AppendDocumentElement(dst, "lsid", client.SessionID)

	var addedTxnNumber bool
	if op.Type == Write && client.RetryWrite {
		addedTxnNumber = true
		dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
	}
	if client.TransactionRunning() || client.RetryingCommit {
		if !addedTxnNumber {
			dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
		}
		if client.TransactionStarting() {
			dst = bsoncore.AppendBooleanElement(dst, "startTransaction", true)
		}
		dst = bsoncore.AppendBooleanElement(dst, "autocommit", false)
	}

	return dst, client.ApplyCommand(desc.Server)
}

func (op Operation) addClusterTime(dst []byte, desc description.SelectedServer) []byte {
	client, clock := op.Client, op.Clock
	if (clock == nil && client == nil) || !sessionsSupported(desc.WireVersion) {
		return dst
	}
	clusterTime := clock.GetClusterTime()
	if client != nil {
		clusterTime = session.MaxClusterTime(clusterTime, client.ClusterTime)
	}
	if clusterTime == nil {
		return dst
	}
	val, err := clusterTime.LookupErr("$clusterTime")
	if err != nil {
		return dst
	}
	return append(bsoncore.AppendHeader(dst, val.Type, "$clusterTime"), val.Value...)
	
}






func (op Operation) calculateMaxTimeMS(ctx context.Context, mon RTTMonitor) (uint64, error) {
	
	
	
	
	
	
	
	
	
	
	if csot.IsTimeoutContext(ctx) && !op.OmitCSOTMaxTimeMS {
		if deadline, ok := ctx.Deadline(); ok {
			remainingTimeout := time.Until(deadline)
			rtt90 := mon.P90()
			maxTime := remainingTimeout - rtt90

			
			
			maxTimeMS := int64((maxTime + (time.Millisecond - 1)) / time.Millisecond)
			if maxTimeMS <= 0 {
				return 0, fmt.Errorf(
					"negative maxTimeMS: remaining time %v until context deadline is less than 90th percentile network round-trip time (%v): %w",
					remainingTimeout,
					mon.Stats(),
					ErrDeadlineWouldBeExceeded)
			}

			
			
			
			
			
			if maxTimeMS > math.MaxInt32 {
				return 0, nil
			}

			return uint64(maxTimeMS), nil
		}
	} else if op.MaxTime != nil {
		
		
		if *op.MaxTime < 0 {
			return 0, ErrNegativeMaxTime
		}
		
		
		return uint64((*op.MaxTime + (time.Millisecond - 1)) / time.Millisecond), nil
	}
	return 0, nil
}




func (op Operation) updateClusterTimes(response bsoncore.Document) {
	
	value, err := response.LookupErr("$clusterTime")
	if err != nil {
		
		return
	}
	clusterTime := bsoncore.BuildDocumentFromElements(nil, bsoncore.AppendValueElement(nil, "$clusterTime", value))

	sess, clock := op.Client, op.Clock

	if sess != nil {
		_ = sess.AdvanceClusterTime(bson.Raw(clusterTime))
	}

	if clock != nil {
		clock.AdvanceClusterTime(bson.Raw(clusterTime))
	}
}




func (op Operation) updateOperationTime(response bsoncore.Document) {
	sess := op.Client
	if sess == nil {
		return
	}

	opTimeElem, err := response.LookupErr("operationTime")
	if err != nil {
		
		return
	}

	t, i := opTimeElem.Timestamp()
	_ = sess.AdvanceOperationTime(&primitive.Timestamp{
		T: t,
		I: i,
	})
}

func (op Operation) getReadPrefBasedOnTransaction() (*readpref.ReadPref, error) {
	if op.Client != nil && op.Client.TransactionRunning() {
		
		rp := op.Client.CurrentRp
		
		
		if rp != nil && !op.Client.TransactionStarting() && rp.Mode() != readpref.PrimaryMode {
			return nil, ErrNonPrimaryReadPref
		}
		return rp, nil
	}
	return op.ReadPreference, nil
}




func (op Operation) createReadPref(desc description.SelectedServer, isOpQuery bool) (bsoncore.Document, error) {
	if op.omitReadPreference {
		return nil, nil
	}

	
	
	if desc.Server.Kind == description.Standalone || (isOpQuery && desc.Server.Kind != description.Mongos) ||
		op.Type == Write || (op.IsOutputAggregate && desc.Server.WireVersion.Max < 13) {
		
		
		
		
		
		
		return nil, nil
	}

	idx, doc := bsoncore.AppendDocumentStart(nil)
	rp, err := op.getReadPrefBasedOnTransaction()
	if err != nil {
		return nil, err
	}

	if rp == nil {
		if desc.Kind == description.Single && desc.Server.Kind != description.Mongos {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		}
		return nil, nil
	}

	switch rp.Mode() {
	case readpref.PrimaryMode:
		if desc.Server.Kind == description.Mongos {
			return nil, nil
		}
		if desc.Kind == description.Single {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		}

		
		
		
		
		
		
		return nil, nil
	case readpref.PrimaryPreferredMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
	case readpref.SecondaryPreferredMode:
		_, ok := rp.MaxStaleness()
		if desc.Server.Kind == description.Mongos && isOpQuery && !ok && len(rp.TagSets()) == 0 && rp.HedgeEnabled() == nil {
			return nil, nil
		}
		doc = bsoncore.AppendStringElement(doc, "mode", "secondaryPreferred")
	case readpref.SecondaryMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "secondary")
	case readpref.NearestMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "nearest")
	}

	sets := make([]bsoncore.Document, 0, len(rp.TagSets()))
	for _, ts := range rp.TagSets() {
		i, set := bsoncore.AppendDocumentStart(nil)
		for _, t := range ts {
			set = bsoncore.AppendStringElement(set, t.Name, t.Value)
		}
		set, _ = bsoncore.AppendDocumentEnd(set, i)
		sets = append(sets, set)
	}
	if len(sets) > 0 {
		var aidx int32
		aidx, doc = bsoncore.AppendArrayElementStart(doc, "tags")
		for i, set := range sets {
			doc = bsoncore.AppendDocumentElement(doc, strconv.Itoa(i), set)
		}
		doc, _ = bsoncore.AppendArrayEnd(doc, aidx)
	}

	if d, ok := rp.MaxStaleness(); ok {
		doc = bsoncore.AppendInt32Element(doc, "maxStalenessSeconds", int32(d.Seconds()))
	}

	if hedgeEnabled := rp.HedgeEnabled(); hedgeEnabled != nil {
		var hedgeIdx int32
		hedgeIdx, doc = bsoncore.AppendDocumentElementStart(doc, "hedge")
		doc = bsoncore.AppendBooleanElement(doc, "enabled", *hedgeEnabled)
		doc, err = bsoncore.AppendDocumentEnd(doc, hedgeIdx)
		if err != nil {
			return nil, fmt.Errorf("error creating hedge document: %w", err)
		}
	}

	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc, nil
}

func (op Operation) secondaryOK(desc description.SelectedServer) wiremessage.QueryFlag {
	if desc.Kind == description.Single && desc.Server.Kind != description.Mongos {
		return wiremessage.SecondaryOK
	}

	if rp := op.ReadPreference; rp != nil && rp.Mode() != readpref.PrimaryMode {
		return wiremessage.SecondaryOK
	}

	return 0
}

func (Operation) canCompress(cmd string) bool {
	if cmd == handshake.LegacyHello || cmd == "hello" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "authenticate" ||
		cmd == "createUser" || cmd == "updateUser" || cmd == "copydbSaslStart" || cmd == "copydbgetnonce" || cmd == "copydb" {
		return false
	}
	return true
}




func (Operation) decodeOpReply(wm []byte) opReply {
	var reply opReply
	var ok bool

	reply.responseFlags, wm, ok = wiremessage.ReadReplyFlags(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing flags")
		return reply
	}
	reply.cursorID, wm, ok = wiremessage.ReadReplyCursorID(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing cursorID")
		return reply
	}
	reply.startingFrom, wm, ok = wiremessage.ReadReplyStartingFrom(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing startingFrom")
		return reply
	}
	reply.numReturned, wm, ok = wiremessage.ReadReplyNumberReturned(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing numberReturned")
		return reply
	}
	reply.documents, _, ok = wiremessage.ReadReplyDocuments(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: could not read documents from reply")
	}

	if reply.responseFlags&wiremessage.QueryFailure == wiremessage.QueryFailure {
		reply.err = QueryFailureError{
			Message:  "command failure",
			Response: reply.documents[0],
		}
		return reply
	}
	if reply.responseFlags&wiremessage.CursorNotFound == wiremessage.CursorNotFound {
		reply.err = ErrCursorNotFound
		return reply
	}
	if reply.numReturned != int32(len(reply.documents)) {
		reply.err = ErrReplyDocumentMismatch
		return reply
	}

	return reply
}

func (op Operation) decodeResult(ctx context.Context, opcode wiremessage.OpCode, wm []byte) (bsoncore.Document, error) {
	switch opcode {
	case wiremessage.OpReply:
		reply := op.decodeOpReply(wm)
		if reply.err != nil {
			return nil, reply.err
		}
		if reply.numReturned == 0 {
			return nil, ErrNoDocCommandResponse
		}
		if reply.numReturned > 1 {
			return nil, ErrMultiDocCommandResponse
		}
		rdr := reply.documents[0]
		if err := rdr.Validate(); err != nil {
			return nil, NewCommandResponseError("malformed OP_REPLY: invalid document", err)
		}

		return rdr, ExtractErrorFromServerResponse(ctx, rdr)
	case wiremessage.OpMsg:
		_, wm, ok := wiremessage.ReadMsgFlags(wm)
		if !ok {
			return nil, errors.New("malformed wire message: missing OP_MSG flags")
		}

		var res bsoncore.Document
		for len(wm) > 0 {
			var stype wiremessage.SectionType
			stype, wm, ok = wiremessage.ReadMsgSectionType(wm)
			if !ok {
				return nil, errors.New("malformed wire message: insuffienct bytes to read section type")
			}

			switch stype {
			case wiremessage.SingleDocument:
				res, wm, ok = wiremessage.ReadMsgSectionSingleDocument(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read single document")
				}
			case wiremessage.DocumentSequence:
				_, _, wm, ok = wiremessage.ReadMsgSectionDocumentSequence(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read document sequence")
				}
			default:
				return nil, fmt.Errorf("malformed wire message: unknown section type %v", stype)
			}
		}

		err := res.Validate()
		if err != nil {
			return nil, NewCommandResponseError("malformed OP_MSG: invalid document", err)
		}

		return res, ExtractErrorFromServerResponse(ctx, res)
	default:
		return nil, fmt.Errorf("cannot decode result from %s", opcode)
	}
}


func (op Operation) getCommandName(doc []byte) string {
	
	idx := bytes.IndexByte(doc[5:], 0x00) 
	return string(doc[5 : idx+5])
}

func (op *Operation) redactCommand(cmd string, doc bsoncore.Document) bool {
	if cmd == "authenticate" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "createUser" ||
		cmd == "updateUser" || cmd == "copydbgetnonce" || cmd == "copydbsaslstart" || cmd == "copydb" {

		return true
	}
	if strings.ToLower(cmd) != handshake.LegacyHelloLowercase && cmd != "hello" {
		return false
	}

	
	_, err := doc.LookupErr("speculativeAuthenticate")
	return err == nil
}


func (op Operation) canLogCommandMessage() bool {
	return op.Logger != nil && op.Logger.LevelComponentEnabled(logger.LevelDebug, logger.ComponentCommand)
}

func (op Operation) canPublishStartedEvent() bool {
	return op.CommandMonitor != nil && op.CommandMonitor.Started != nil
}




func (op Operation) publishStartedEvent(ctx context.Context, info startedInformation) {
	
	if op.canLogCommandMessage() {
		host, port, _ := net.SplitHostPort(info.serverAddress.String())

		redactedCmd := redactStartedInformationCmd(op, info).String()
		formattedCmd := logger.FormatMessage(redactedCmd, op.Logger.MaxDocumentLength)

		op.Logger.Print(logger.LevelDebug,
			logger.ComponentCommand,
			logger.CommandStarted,
			logger.SerializeCommand(logger.Command{
				DriverConnectionID: info.driverConnectionID,
				Message:            logger.CommandStarted,
				Name:               info.cmdName,
				DatabaseName:       op.Database,
				RequestID:          int64(info.requestID),
				ServerConnectionID: info.serverConnID,
				ServerHost:         host,
				ServerPort:         port,
				ServiceID:          info.serviceID,
			},
				logger.KeyCommand, formattedCmd)...)

	}

	if op.canPublishStartedEvent() {
		started := &event.CommandStartedEvent{
			Command:              redactStartedInformationCmd(op, info),
			DatabaseName:         op.Database,
			CommandName:          info.cmdName,
			RequestID:            int64(info.requestID),
			ConnectionID:         info.connID,
			ServerConnectionID:   convertInt64PtrToInt32Ptr(info.serverConnID),
			ServerConnectionID64: info.serverConnID,
			ServiceID:            info.serviceID,
		}
		op.CommandMonitor.Started(ctx, started)
	}
}




func (op Operation) canPublishFinishedEvent(info finishedInformation) bool {
	success := info.success()

	return op.CommandMonitor != nil &&
		(!success || op.CommandMonitor.Succeeded != nil) &&
		(success || op.CommandMonitor.Failed != nil)
}



func (op Operation) publishFinishedEvent(ctx context.Context, info finishedInformation) {
	if op.canLogCommandMessage() && info.success() {
		host, port, _ := net.SplitHostPort(info.serverAddress.String())

		redactedReply := redactFinishedInformationResponse(info).String()
		formattedReply := logger.FormatMessage(redactedReply, op.Logger.MaxDocumentLength)

		op.Logger.Print(logger.LevelDebug,
			logger.ComponentCommand,
			logger.CommandSucceeded,
			logger.SerializeCommand(logger.Command{
				DriverConnectionID: info.driverConnectionID,
				Message:            logger.CommandSucceeded,
				Name:               info.cmdName,
				DatabaseName:       op.Database,
				RequestID:          int64(info.requestID),
				ServerConnectionID: info.serverConnID,
				ServerHost:         host,
				ServerPort:         port,
				ServiceID:          info.serviceID,
			},
				logger.KeyDurationMS, info.duration.Milliseconds(),
				logger.KeyReply, formattedReply)...)
	}

	if op.canLogCommandMessage() && !info.success() {
		host, port, _ := net.SplitHostPort(info.serverAddress.String())

		formattedReply := logger.FormatMessage(info.cmdErr.Error(), op.Logger.MaxDocumentLength)

		op.Logger.Print(logger.LevelDebug,
			logger.ComponentCommand,
			logger.CommandFailed,
			logger.SerializeCommand(logger.Command{
				DriverConnectionID: info.driverConnectionID,
				Message:            logger.CommandFailed,
				Name:               info.cmdName,
				DatabaseName:       op.Database,
				RequestID:          int64(info.requestID),
				ServerConnectionID: info.serverConnID,
				ServerHost:         host,
				ServerPort:         port,
				ServiceID:          info.serviceID,
			},
				logger.KeyDurationMS, info.duration.Milliseconds(),
				logger.KeyFailure, formattedReply)...)
	}

	
	if !op.canPublishFinishedEvent(info) {
		return
	}

	finished := event.CommandFinishedEvent{
		CommandName:          info.cmdName,
		DatabaseName:         op.Database,
		RequestID:            int64(info.requestID),
		ConnectionID:         info.connID,
		Duration:             info.duration,
		DurationNanos:        info.duration.Nanoseconds(),
		ServerConnectionID:   convertInt64PtrToInt32Ptr(info.serverConnID),
		ServerConnectionID64: info.serverConnID,
		ServiceID:            info.serviceID,
	}

	if info.success() {
		successEvent := &event.CommandSucceededEvent{
			Reply:                redactFinishedInformationResponse(info),
			CommandFinishedEvent: finished,
		}
		op.CommandMonitor.Succeeded(ctx, successEvent)

		return
	}

	failedEvent := &event.CommandFailedEvent{
		Failure:              info.cmdErr.Error(),
		CommandFinishedEvent: finished,
	}
	op.CommandMonitor.Failed(ctx, failedEvent)
}


func sessionsSupported(wireVersion *description.VersionRange) bool {
	return wireVersion != nil
}


func retryWritesSupported(s description.Server) bool {
	return s.SessionTimeoutMinutesPtr != nil && s.Kind != description.Standalone
}
