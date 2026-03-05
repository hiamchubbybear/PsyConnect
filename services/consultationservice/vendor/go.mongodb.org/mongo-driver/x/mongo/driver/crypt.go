





package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

const (
	defaultKmsPort    = 443
	defaultKmsTimeout = 10 * time.Second
)


type CollectionInfoFn func(ctx context.Context, db string, filter bsoncore.Document) (bsoncore.Document, error)


type KeyRetrieverFn func(ctx context.Context, filter bsoncore.Document) ([]bsoncore.Document, error)


type MarkCommandFn func(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error)


type CryptOptions struct {
	MongoCrypt           *mongocrypt.MongoCrypt
	CollInfoFn           CollectionInfoFn
	KeyFn                KeyRetrieverFn
	MarkFn               MarkCommandFn
	TLSConfig            map[string]*tls.Config
	BypassAutoEncryption bool
	BypassQueryAnalysis  bool
}






type Crypt interface {
	
	Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error)
	
	Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error)
	
	CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error)
	
	EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error)
	
	EncryptExplicitExpression(ctx context.Context, val bsoncore.Document, opts *options.ExplicitEncryptionOptions) (bsoncore.Document, error)
	
	DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error)
	
	Close()
	
	BypassAutoEncryption() bool
	
	
	RewrapDataKey(ctx context.Context, filter []byte, opts *options.RewrapManyDataKeyOptions) ([]bsoncore.Document, error)
}



type crypt struct {
	mongoCrypt *mongocrypt.MongoCrypt
	collInfoFn CollectionInfoFn
	keyFn      KeyRetrieverFn
	markFn     MarkCommandFn
	tlsConfig  map[string]*tls.Config

	bypassAutoEncryption bool
}


func NewCrypt(opts *CryptOptions) Crypt {
	c := &crypt{
		mongoCrypt:           opts.MongoCrypt,
		collInfoFn:           opts.CollInfoFn,
		keyFn:                opts.KeyFn,
		markFn:               opts.MarkFn,
		tlsConfig:            opts.TLSConfig,
		bypassAutoEncryption: opts.BypassAutoEncryption,
	}
	return c
}


func (c *crypt) Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error) {
	if c.bypassAutoEncryption {
		return cmd, nil
	}

	cryptCtx, err := c.mongoCrypt.CreateEncryptionContext(db, cmd)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, db)
}


func (c *crypt) Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error) {
	cryptCtx, err := c.mongoCrypt.CreateDecryptionContext(cmdResponse)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
}


func (c *crypt) CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error) {
	cryptCtx, err := c.mongoCrypt.CreateDataKeyContext(kmsProvider, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
}



func (c *crypt) RewrapDataKey(ctx context.Context, filter []byte,
	opts *options.RewrapManyDataKeyOptions) ([]bsoncore.Document, error) {

	cryptCtx, err := c.mongoCrypt.RewrapDataKeyContext(filter, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	rewrappedBSON, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return nil, err
	}
	if rewrappedBSON == nil {
		return nil, nil
	}

	
	
	rewrappedDocumentBytes, err := rewrappedBSON.LookupErr("v")
	if err != nil {
		return nil, err
	}

	
	rewrappedDocsArray, ok := rewrappedDocumentBytes.ArrayOK()
	if !ok {
		return nil, fmt.Errorf("expected results from mongocrypt_ctx_rewrap_many_datakey_init to be an array")
	}

	rewrappedDocumentValues, err := rewrappedDocsArray.Values()
	if err != nil {
		return nil, err
	}

	rewrappedDocuments := []bsoncore.Document{}
	for _, rewrappedDocumentValue := range rewrappedDocumentValues {
		if rewrappedDocumentValue.Type != bsontype.EmbeddedDocument {
			
			
			return nil, fmt.Errorf("expected value of type %q, got: %q",
				bsontype.EmbeddedDocument.String(),
				rewrappedDocumentValue.Type.String())
		}
		rewrappedDocuments = append(rewrappedDocuments, rewrappedDocumentValue.Document())
	}
	return rewrappedDocuments, nil
}


func (c *crypt) EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendValueElement(doc, "v", val)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitEncryptionContext(doc, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return 0, nil, err
	}

	sub, data := res.Lookup("v").Binary()
	return sub, data, nil
}


func (c *crypt) EncryptExplicitExpression(ctx context.Context, expr bsoncore.Document, opts *options.ExplicitEncryptionOptions) (bsoncore.Document, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendDocumentElement(doc, "v", expr)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitEncryptionExpressionContext(doc, opts)
	if err != nil {
		return nil, err
	}
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return nil, err
	}

	encryptedExpr := res.Lookup("v").Document()
	return encryptedExpr, nil
}


func (c *crypt) DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error) {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendBinaryElement(doc, "v", subtype, data)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitDecryptionContext(doc)
	if err != nil {
		return bsoncore.Value{}, err
	}
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil {
		return bsoncore.Value{}, err
	}

	return res.Lookup("v"), nil
}


func (c *crypt) Close() {
	c.mongoCrypt.Close()
}

func (c *crypt) BypassAutoEncryption() bool {
	return c.bypassAutoEncryption
}

func (c *crypt) executeStateMachine(ctx context.Context, cryptCtx *mongocrypt.Context, db string) (bsoncore.Document, error) {
	var err error
	for {
		state := cryptCtx.State()
		switch state {
		case mongocrypt.NeedMongoCollInfo:
			err = c.collectionInfo(ctx, cryptCtx, db)
		case mongocrypt.NeedMongoMarkings:
			err = c.markCommand(ctx, cryptCtx, db)
		case mongocrypt.NeedMongoKeys:
			err = c.retrieveKeys(ctx, cryptCtx)
		case mongocrypt.NeedKms:
			err = c.decryptKeys(cryptCtx)
		case mongocrypt.Ready:
			return cryptCtx.Finish()
		case mongocrypt.Done:
			return nil, nil
		case mongocrypt.NeedKmsCredentials:
			err = c.provideKmsProviders(ctx, cryptCtx)
		default:
			return nil, fmt.Errorf("invalid Crypt state: %v", state)
		}
		if err != nil {
			return nil, err
		}
	}
}

func (c *crypt) collectionInfo(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error {
	op, err := cryptCtx.NextOperation()
	if err != nil {
		return err
	}

	collInfo, err := c.collInfoFn(ctx, db, op)
	if err != nil {
		return err
	}
	if collInfo != nil {
		if err = cryptCtx.AddOperationResult(collInfo); err != nil {
			return err
		}
	}

	return cryptCtx.CompleteOperation()
}

func (c *crypt) markCommand(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error {
	op, err := cryptCtx.NextOperation()
	if err != nil {
		return err
	}

	markedCmd, err := c.markFn(ctx, db, op)
	if err != nil {
		return err
	}
	if err = cryptCtx.AddOperationResult(markedCmd); err != nil {
		return err
	}

	return cryptCtx.CompleteOperation()
}

func (c *crypt) retrieveKeys(ctx context.Context, cryptCtx *mongocrypt.Context) error {
	op, err := cryptCtx.NextOperation()
	if err != nil {
		return err
	}

	keys, err := c.keyFn(ctx, op)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err = cryptCtx.AddOperationResult(key); err != nil {
			return err
		}
	}

	return cryptCtx.CompleteOperation()
}

func (c *crypt) decryptKeys(cryptCtx *mongocrypt.Context) error {
	for {
		kmsCtx := cryptCtx.NextKmsContext()
		if kmsCtx == nil {
			break
		}

		if err := c.decryptKey(kmsCtx); err != nil {
			return err
		}
	}

	return cryptCtx.FinishKmsContexts()
}

func (c *crypt) decryptKey(kmsCtx *mongocrypt.KmsContext) error {
	host, err := kmsCtx.HostName()
	if err != nil {
		return err
	}
	msg, err := kmsCtx.Message()
	if err != nil {
		return err
	}

	
	addr := host
	if idx := strings.IndexByte(host, ':'); idx == -1 {
		addr = fmt.Sprintf("%s:%d", host, defaultKmsPort)
	}

	kmsProvider := kmsCtx.KMSProvider()
	tlsCfg := c.tlsConfig[kmsProvider]
	if tlsCfg == nil {
		tlsCfg = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	if err = conn.SetWriteDeadline(time.Now().Add(defaultKmsTimeout)); err != nil {
		return err
	}
	if _, err = conn.Write(msg); err != nil {
		return err
	}

	for {
		bytesNeeded := kmsCtx.BytesNeeded()
		if bytesNeeded == 0 {
			return nil
		}

		res := make([]byte, bytesNeeded)
		bytesRead, err := conn.Read(res)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if err = kmsCtx.FeedResponse(res[:bytesRead]); err != nil {
			return err
		}
	}
}

func (c *crypt) provideKmsProviders(ctx context.Context, cryptCtx *mongocrypt.Context) error {
	kmsProviders, err := c.mongoCrypt.GetKmsProviders(ctx)
	if err != nil {
		return err
	}
	return cryptCtx.ProvideKmsProviders(kmsProviders)
}
