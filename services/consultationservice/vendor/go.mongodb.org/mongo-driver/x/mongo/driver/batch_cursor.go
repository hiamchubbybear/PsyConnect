





package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/codecutil"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)



var ErrNoCursor = errors.New("database response does not contain a cursor")



type BatchCursor struct {
	clientSession        *session.Client
	clock                *session.ClusterClock
	comment              interface{}
	encoderFn            codecutil.EncoderFn
	database             string
	collection           string
	id                   int64
	err                  error
	server               Server
	serverDescription    description.Server
	errorProcessor       ErrorProcessor 
	connection           PinnedConnection
	batchSize            int32
	maxTimeMS            int64
	currentBatch         *bsoncore.DocumentSequence
	firstBatch           bool
	cmdMonitor           *event.CommandMonitor
	postBatchResumeToken bsoncore.Document
	crypt                Crypt
	serverAPI            *ServerAPIOptions

	
	limit       int32
	numReturned int32 
}



type CursorResponse struct {
	Server               Server
	ErrorProcessor       ErrorProcessor 
	Connection           PinnedConnection
	Desc                 description.Server
	FirstBatch           *bsoncore.DocumentSequence
	Database             string
	Collection           string
	ID                   int64
	postBatchResumeToken bsoncore.Document
}






func NewCursorResponse(info ResponseInfo) (CursorResponse, error) {
	response := info.ServerResponse
	cur, err := response.LookupErr("cursor")
	if errors.Is(err, bsoncore.ErrElementNotFound) {
		return CursorResponse{}, ErrNoCursor
	}
	if err != nil {
		return CursorResponse{}, fmt.Errorf("error getting cursor from database response: %w", err)
	}
	curDoc, ok := cur.DocumentOK()
	if !ok {
		return CursorResponse{}, fmt.Errorf("cursor should be an embedded document but is BSON type %s", cur.Type)
	}
	elems, err := curDoc.Elements()
	if err != nil {
		return CursorResponse{}, fmt.Errorf("error getting elements from cursor: %w", err)
	}
	curresp := CursorResponse{Server: info.Server, Desc: info.ConnectionDescription}

	for _, elem := range elems {
		switch elem.Key() {
		case "firstBatch":
			arr, ok := elem.Value().ArrayOK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("firstBatch should be an array but is a BSON %s", elem.Value().Type)
			}
			curresp.FirstBatch = &bsoncore.DocumentSequence{Style: bsoncore.ArrayStyle, Data: arr}
		case "ns":
			ns, ok := elem.Value().StringValueOK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("ns should be a string but is a BSON %s", elem.Value().Type)
			}
			database, collection, ok := strings.Cut(ns, ".")
			if !ok {
				return CursorResponse{}, errors.New("ns field must contain a valid namespace, but is missing '.'")
			}
			curresp.Database = database
			curresp.Collection = collection
		case "id":
			curresp.ID, ok = elem.Value().Int64OK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("id should be an int64 but it is a BSON %s", elem.Value().Type)
			}
		case "postBatchResumeToken":
			curresp.postBatchResumeToken, ok = elem.Value().DocumentOK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("post batch resume token should be a document but it is a BSON %s", elem.Value().Type)
			}
		}
	}

	
	
	if curresp.Desc.LoadBalanced() && curresp.ID != 0 {
		
		ep, ok := curresp.Server.(ErrorProcessor)
		if !ok {
			return CursorResponse{}, fmt.Errorf("expected Server used to establish a cursor to implement ErrorProcessor, but got %T", curresp.Server)
		}
		curresp.ErrorProcessor = ep

		refConn, ok := info.Connection.(PinnedConnection)
		if !ok {
			return CursorResponse{}, fmt.Errorf("expected Connection used to establish a cursor to implement PinnedConnection, but got %T", info.Connection)
		}
		if err := refConn.PinToCursor(); err != nil {
			return CursorResponse{}, fmt.Errorf("error incrementing connection reference count when creating a cursor: %w", err)
		}
		curresp.Connection = refConn
	}

	return curresp, nil
}


type CursorOptions struct {
	BatchSize             int32
	Comment               bsoncore.Value
	MaxTimeMS             int64
	Limit                 int32
	CommandMonitor        *event.CommandMonitor
	Crypt                 Crypt
	ServerAPI             *ServerAPIOptions
	MarshalValueEncoderFn func(io.Writer) (*bson.Encoder, error)
}


func NewBatchCursor(cr CursorResponse, clientSession *session.Client, clock *session.ClusterClock, opts CursorOptions) (*BatchCursor, error) {
	ds := cr.FirstBatch
	bc := &BatchCursor{
		clientSession:        clientSession,
		clock:                clock,
		comment:              opts.Comment,
		database:             cr.Database,
		collection:           cr.Collection,
		id:                   cr.ID,
		server:               cr.Server,
		connection:           cr.Connection,
		errorProcessor:       cr.ErrorProcessor,
		batchSize:            opts.BatchSize,
		maxTimeMS:            opts.MaxTimeMS,
		cmdMonitor:           opts.CommandMonitor,
		firstBatch:           true,
		postBatchResumeToken: cr.postBatchResumeToken,
		crypt:                opts.Crypt,
		serverAPI:            opts.ServerAPI,
		serverDescription:    cr.Desc,
		encoderFn:            opts.MarshalValueEncoderFn,
	}

	if ds != nil {
		bc.numReturned = int32(ds.DocumentCount())
	}
	if cr.Desc.WireVersion == nil {
		bc.limit = opts.Limit

		
		if bc.limit != 0 && bc.limit < bc.numReturned {
			for i := int32(0); i < bc.limit; i++ {
				_, err := ds.Next()
				if err != nil {
					return nil, err
				}
			}
			ds.Data = ds.Data[:ds.Pos]
			ds.ResetIterator()
		}
	}

	bc.currentBatch = ds
	return bc, nil
}


func NewEmptyBatchCursor() *BatchCursor {
	return &BatchCursor{currentBatch: new(bsoncore.DocumentSequence)}
}



func NewBatchCursorFromDocuments(documents []byte) *BatchCursor {
	return &BatchCursor{
		currentBatch: &bsoncore.DocumentSequence{
			Data:  documents,
			Style: bsoncore.SequenceStyle,
		},
		
		
		id:     0,
		server: nil,
	}
}


func (bc *BatchCursor) ID() int64 {
	return bc.id
}






func (bc *BatchCursor) Next(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	if bc.firstBatch {
		bc.firstBatch = false
		return !bc.currentBatch.Empty()
	}

	if bc.id == 0 || bc.server == nil {
		return false
	}

	bc.getMore(ctx)

	return !bc.currentBatch.Empty()
}



func (bc *BatchCursor) Batch() *bsoncore.DocumentSequence { return bc.currentBatch }


func (bc *BatchCursor) Err() error { return bc.err }


func (bc *BatchCursor) Close(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	err := bc.KillCursor(ctx)
	bc.id = 0
	bc.currentBatch.Data = nil
	bc.currentBatch.Style = 0
	bc.currentBatch.ResetIterator()

	connErr := bc.unpinConnection()
	if err == nil {
		err = connErr
	}
	return err
}

func (bc *BatchCursor) unpinConnection() error {
	if bc.connection == nil {
		return nil
	}

	err := bc.connection.UnpinFromCursor()
	closeErr := bc.connection.Close()
	if err == nil && closeErr != nil {
		err = closeErr
	}
	bc.connection = nil
	return err
}


func (bc *BatchCursor) Server() Server {
	return bc.server
}

func (bc *BatchCursor) clearBatch() {
	bc.currentBatch.Data = bc.currentBatch.Data[:0]
}


func (bc *BatchCursor) KillCursor(ctx context.Context) error {
	if bc.server == nil || bc.id == 0 {
		return nil
	}

	return Operation{
		CommandFn: func(dst []byte, _ description.SelectedServer) ([]byte, error) {
			dst = bsoncore.AppendStringElement(dst, "killCursors", bc.collection)
			dst = bsoncore.BuildArrayElement(dst, "cursors", bsoncore.Value{Type: bsontype.Int64, Data: bsoncore.AppendInt64(nil, bc.id)})
			return dst, nil
		},
		Database:       bc.database,
		Deployment:     bc.getOperationDeployment(),
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyKillCursors,
		CommandMonitor: bc.cmdMonitor,
		ServerAPI:      bc.serverAPI,

		
		
		
		
		omitReadPreference: true,
	}.Execute(ctx)
}





func calcGetMoreBatchSize(bc BatchCursor) (int32, bool) {
	gmBatchSize := bc.batchSize

	
	if bc.limit != 0 && bc.numReturned+bc.batchSize >= bc.limit {
		gmBatchSize = bc.limit - bc.numReturned
		if gmBatchSize <= 0 {
			return gmBatchSize, false
		}
	}

	return gmBatchSize, true
}

func (bc *BatchCursor) getMore(ctx context.Context) {
	bc.clearBatch()
	if bc.id == 0 {
		return
	}

	numToReturn, ok := calcGetMoreBatchSize(*bc)
	if !ok {
		if err := bc.Close(ctx); err != nil {
			bc.err = err
		}

		return
	}

	bc.err = Operation{
		CommandFn: func(dst []byte, _ description.SelectedServer) ([]byte, error) {
			dst = bsoncore.AppendInt64Element(dst, "getMore", bc.id)
			dst = bsoncore.AppendStringElement(dst, "collection", bc.collection)
			if numToReturn > 0 {
				dst = bsoncore.AppendInt32Element(dst, "batchSize", numToReturn)
			}
			if bc.maxTimeMS > 0 {
				dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", bc.maxTimeMS)
			}

			comment, err := codecutil.MarshalValue(bc.comment, bc.encoderFn)
			if err != nil {
				return nil, fmt.Errorf("error marshaling comment as a BSON value: %w", err)
			}

			
			if comment.Type != bsontype.Type(0) && bc.serverDescription.WireVersion.Max >= 9 {
				dst = bsoncore.AppendValueElement(dst, "comment", comment)
			}

			return dst, nil
		},
		Database:   bc.database,
		Deployment: bc.getOperationDeployment(),
		ProcessResponseFn: func(info ResponseInfo) error {
			response := info.ServerResponse
			id, ok := response.Lookup("cursor", "id").Int64OK()
			if !ok {
				return fmt.Errorf("cursor.id should be an int64 but is a BSON %s", response.Lookup("cursor", "id").Type)
			}
			bc.id = id

			batch, ok := response.Lookup("cursor", "nextBatch").ArrayOK()
			if !ok {
				return fmt.Errorf("cursor.nextBatch should be an array but is a BSON %s", response.Lookup("cursor", "nextBatch").Type)
			}
			bc.currentBatch.Style = bsoncore.ArrayStyle
			bc.currentBatch.Data = batch
			bc.currentBatch.ResetIterator()
			bc.numReturned += int32(bc.currentBatch.DocumentCount()) 

			pbrt, err := response.LookupErr("cursor", "postBatchResumeToken")
			if err != nil {
				
				return nil
			}

			pbrtDoc, ok := pbrt.DocumentOK()
			if !ok {
				bc.err = fmt.Errorf("expected BSON type for post batch resume token to be EmbeddedDocument but got %s", pbrt.Type)
				return nil
			}

			bc.postBatchResumeToken = pbrtDoc

			return nil
		},
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyGetMore,
		CommandMonitor: bc.cmdMonitor,
		Crypt:          bc.crypt,
		ServerAPI:      bc.serverAPI,

		
		
		
		
		omitReadPreference: true,
	}.Execute(ctx)

	
	if bc.id == 0 {
		err := bc.unpinConnection()
		if err != nil && bc.err == nil {
			bc.err = err
		}
	}

	
	
	
	if driverErr, ok := bc.err.(Error); ok && driverErr.NetworkError() && bc.connection != nil {
		bc.id = 0
	}

	
	if bc.limit != 0 && bc.numReturned >= bc.limit {
		
		err := bc.KillCursor(ctx)
		if err != nil && bc.err == nil {
			bc.err = err
		}
	}
}


func (bc *BatchCursor) PostBatchResumeToken() bsoncore.Document {
	return bc.postBatchResumeToken
}


func (bc *BatchCursor) SetBatchSize(size int32) {
	bc.batchSize = size
}







func (bc *BatchCursor) SetMaxTime(dur time.Duration) {
	bc.maxTimeMS = int64(dur / time.Millisecond)
}


func (bc *BatchCursor) SetComment(comment interface{}) {
	bc.comment = comment
}

func (bc *BatchCursor) getOperationDeployment() Deployment {
	if bc.connection != nil {
		return &loadBalancedCursorDeployment{
			errorProcessor: bc.errorProcessor,
			conn:           bc.connection,
		}
	}
	return SingleServerDeployment{bc.server}
}




type loadBalancedCursorDeployment struct {
	errorProcessor ErrorProcessor
	conn           PinnedConnection
}

var _ Deployment = (*loadBalancedCursorDeployment)(nil)
var _ Server = (*loadBalancedCursorDeployment)(nil)
var _ ErrorProcessor = (*loadBalancedCursorDeployment)(nil)

func (lbcd *loadBalancedCursorDeployment) SelectServer(_ context.Context, _ description.ServerSelector) (Server, error) {
	return lbcd, nil
}

func (lbcd *loadBalancedCursorDeployment) Kind() description.TopologyKind {
	return description.LoadBalanced
}

func (lbcd *loadBalancedCursorDeployment) Connection(_ context.Context) (Connection, error) {
	return lbcd.conn, nil
}


func (lbcd *loadBalancedCursorDeployment) RTTMonitor() RTTMonitor {
	return &csot.ZeroRTTMonitor{}
}

func (lbcd *loadBalancedCursorDeployment) ProcessError(err error, conn Connection) ProcessErrorResult {
	return lbcd.errorProcessor.ProcessError(err, conn)
}
