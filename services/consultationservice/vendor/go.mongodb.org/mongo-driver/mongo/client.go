





package mongo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/httputil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	mcopts "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

const (
	defaultLocalThreshold = 15 * time.Millisecond
	defaultMaxPoolSize    = 100
)

var (
	
	keyVaultCollOpts = options.Collection().SetReadConcern(readconcern.Majority()).
				SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	endSessionsBatchSize = 10000
)






type Client struct {
	id             uuid.UUID
	deployment     driver.Deployment
	localThreshold time.Duration
	retryWrites    bool
	retryReads     bool
	clock          *session.ClusterClock
	readPreference *readpref.ReadPref
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
	monitor        *event.CommandMonitor
	serverAPI      *driver.ServerAPIOptions
	serverMonitor  *event.ServerMonitor
	sessionPool    *session.Pool
	timeout        *time.Duration
	httpClient     *http.Client
	logger         *logger.Logger

	
	keyVaultClientFLE  *Client
	keyVaultCollFLE    *Collection
	mongocryptdFLE     *mongocryptdClient
	cryptFLE           driver.Crypt
	metadataClientFLE  *Client
	internalClientFLE  *Client
	encryptedFieldsMap map[string]interface{}
	authenticator      driver.Authenticator
}























func Connect(ctx context.Context, opts ...*options.ClientOptions) (*Client, error) {
	c, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}
















func NewClient(opts ...*options.ClientOptions) (*Client, error) {
	clientOpt := options.MergeClientOptions(opts...)

	id, err := uuid.New()
	if err != nil {
		return nil, err
	}
	client := &Client{id: id}

	
	client.clock = new(session.ClusterClock)

	
	client.localThreshold = defaultLocalThreshold
	if clientOpt.LocalThreshold != nil {
		client.localThreshold = *clientOpt.LocalThreshold
	}
	
	if clientOpt.Monitor != nil {
		client.monitor = clientOpt.Monitor
	}
	
	if clientOpt.ServerMonitor != nil {
		client.serverMonitor = clientOpt.ServerMonitor
	}
	
	client.readConcern = readconcern.New()
	if clientOpt.ReadConcern != nil {
		client.readConcern = clientOpt.ReadConcern
	}
	
	client.readPreference = readpref.Primary()
	if clientOpt.ReadPreference != nil {
		client.readPreference = clientOpt.ReadPreference
	}
	
	if clientOpt.BSONOptions != nil {
		client.bsonOpts = clientOpt.BSONOptions
	}
	
	client.registry = bson.DefaultRegistry
	if clientOpt.Registry != nil {
		client.registry = clientOpt.Registry
	}
	
	client.retryWrites = true 
	if clientOpt.RetryWrites != nil {
		client.retryWrites = *clientOpt.RetryWrites
	}
	client.retryReads = true
	if clientOpt.RetryReads != nil {
		client.retryReads = *clientOpt.RetryReads
	}
	
	client.timeout = clientOpt.Timeout
	client.httpClient = clientOpt.HTTPClient
	
	if clientOpt.WriteConcern != nil {
		client.writeConcern = clientOpt.WriteConcern
	}
	
	if clientOpt.AutoEncryptionOptions != nil {
		if err := client.configureAutoEncryption(clientOpt); err != nil {
			return nil, err
		}
	} else {
		client.cryptFLE = clientOpt.Crypt
	}

	
	if clientOpt.Deployment != nil {
		client.deployment = clientOpt.Deployment
	}

	
	if clientOpt.MaxPoolSize == nil {
		clientOpt.SetMaxPoolSize(defaultMaxPoolSize)
	}

	if clientOpt.Auth != nil {
		client.authenticator, err = auth.CreateAuthenticator(
			clientOpt.Auth.AuthMechanism,
			topology.ConvertCreds(clientOpt.Auth),
			clientOpt.HTTPClient,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator: %w", err)
		}
	}

	cfg, err := topology.NewConfigWithAuthenticator(clientOpt, client.clock, client.authenticator)
	if err != nil {
		return nil, err
	}

	client.serverAPI = topology.ServerAPIFromServerOptions(cfg.ServerOpts)

	if client.deployment == nil {
		client.deployment, err = topology.New(cfg)
		if err != nil {
			return nil, replaceErrors(err)
		}
	}

	
	client.logger, err = newLogger(clientOpt.LoggerOptions)
	if err != nil {
		return nil, fmt.Errorf("invalid logger options: %w", err)
	}

	return client, nil
}








func (c *Client) Connect(ctx context.Context) error {
	if connector, ok := c.deployment.(driver.Connector); ok {
		err := connector.Connect()
		if err != nil {
			return replaceErrors(err)
		}
	}

	if c.mongocryptdFLE != nil {
		if err := c.mongocryptdFLE.connect(ctx); err != nil {
			return err
		}
	}

	if c.internalClientFLE != nil {
		if err := c.internalClientFLE.Connect(ctx); err != nil {
			return err
		}
	}

	if c.keyVaultClientFLE != nil && c.keyVaultClientFLE != c.internalClientFLE && c.keyVaultClientFLE != c {
		if err := c.keyVaultClientFLE.Connect(ctx); err != nil {
			return err
		}
	}

	if c.metadataClientFLE != nil && c.metadataClientFLE != c.internalClientFLE && c.metadataClientFLE != c {
		if err := c.metadataClientFLE.Connect(ctx); err != nil {
			return err
		}
	}

	var updateChan <-chan description.Topology
	if subscriber, ok := c.deployment.(driver.Subscriber); ok {
		sub, err := subscriber.Subscribe()
		if err != nil {
			return replaceErrors(err)
		}
		updateChan = sub.Updates
	}
	c.sessionPool = session.NewPool(updateChan)
	return nil
}









func (c *Client) Disconnect(ctx context.Context) error {
	if c.logger != nil {
		defer c.logger.Close()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if c.httpClient == httputil.DefaultHTTPClient {
		defer httputil.CloseIdleHTTPConnections(c.httpClient)
	}

	c.endSessions(ctx)
	if c.mongocryptdFLE != nil {
		if err := c.mongocryptdFLE.disconnect(ctx); err != nil {
			return err
		}
	}

	if c.internalClientFLE != nil {
		if err := c.internalClientFLE.Disconnect(ctx); err != nil {
			return err
		}
	}

	if c.keyVaultClientFLE != nil && c.keyVaultClientFLE != c.internalClientFLE && c.keyVaultClientFLE != c {
		if err := c.keyVaultClientFLE.Disconnect(ctx); err != nil {
			return err
		}
	}
	if c.metadataClientFLE != nil && c.metadataClientFLE != c.internalClientFLE && c.metadataClientFLE != c {
		if err := c.metadataClientFLE.Disconnect(ctx); err != nil {
			return err
		}
	}
	if c.cryptFLE != nil {
		c.cryptFLE.Close()
	}

	if disconnector, ok := c.deployment.(driver.Disconnector); ok {
		return replaceErrors(disconnector.Disconnect(ctx))
	}

	return nil
}












func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if rp == nil {
		rp = c.readPreference
	}

	db := c.Database("admin")
	res := db.RunCommand(ctx, bson.D{
		{"ping", 1},
	}, options.RunCmd().SetReadPreference(rp))

	return replaceErrors(res.Err())
}











func (c *Client) StartSession(opts ...*options.SessionOptions) (Session, error) {
	if c.sessionPool == nil {
		return nil, ErrClientDisconnected
	}

	sopts := options.MergeSessionOptions(opts...)
	coreOpts := &session.ClientOptions{
		DefaultReadConcern:    c.readConcern,
		DefaultReadPreference: c.readPreference,
		DefaultWriteConcern:   c.writeConcern,
	}
	if sopts.CausalConsistency != nil {
		coreOpts.CausalConsistency = sopts.CausalConsistency
	}
	if sopts.DefaultReadConcern != nil {
		coreOpts.DefaultReadConcern = sopts.DefaultReadConcern
	}
	if sopts.DefaultWriteConcern != nil {
		coreOpts.DefaultWriteConcern = sopts.DefaultWriteConcern
	}
	if sopts.DefaultReadPreference != nil {
		coreOpts.DefaultReadPreference = sopts.DefaultReadPreference
	}
	if sopts.DefaultMaxCommitTime != nil {
		coreOpts.DefaultMaxCommitTime = sopts.DefaultMaxCommitTime
	}
	if sopts.Snapshot != nil {
		coreOpts.Snapshot = sopts.Snapshot
	}

	sess, err := session.NewClientSession(c.sessionPool, c.id, coreOpts)
	if err != nil {
		return nil, replaceErrors(err)
	}

	
	sess.RetryWrite = false
	sess.RetryRead = c.retryReads

	return &sessionImpl{
		clientSession: sess,
		client:        c,
		deployment:    c.deployment,
	}, nil
}

func (c *Client) endSessions(ctx context.Context) {
	if c.sessionPool == nil {
		return
	}

	sessionIDs := c.sessionPool.IDSlice()
	op := operation.NewEndSessions(nil).ClusterClock(c.clock).Deployment(c.deployment).
		ServerSelector(description.ReadPrefSelector(readpref.PrimaryPreferred())).CommandMonitor(c.monitor).
		Database("admin").Crypt(c.cryptFLE).ServerAPI(c.serverAPI)

	totalNumIDs := len(sessionIDs)
	var currentBatch []bsoncore.Document
	for i := 0; i < totalNumIDs; i++ {
		currentBatch = append(currentBatch, sessionIDs[i])

		
		if ((i+1)%endSessionsBatchSize) == 0 || i == totalNumIDs-1 {
			
			_, marshalVal, err := bson.MarshalValue(currentBatch)
			if err == nil {
				_ = op.SessionIDs(marshalVal).Execute(ctx)
			}

			currentBatch = currentBatch[:0]
		}
	}
}

func (c *Client) configureAutoEncryption(clientOpts *options.ClientOptions) error {
	c.encryptedFieldsMap = clientOpts.AutoEncryptionOptions.EncryptedFieldsMap
	if err := c.configureKeyVaultClientFLE(clientOpts); err != nil {
		return err
	}
	if err := c.configureMetadataClientFLE(clientOpts); err != nil {
		return err
	}

	mc, err := c.newMongoCrypt(clientOpts.AutoEncryptionOptions)
	if err != nil {
		return err
	}

	
	if mc.CryptSharedLibVersionString() == "" {
		mongocryptdFLE, err := newMongocryptdClient(clientOpts.AutoEncryptionOptions)
		if err != nil {
			return err
		}
		c.mongocryptdFLE = mongocryptdFLE
	}

	c.configureCryptFLE(mc, clientOpts.AutoEncryptionOptions)
	return nil
}

func (c *Client) getOrCreateInternalClient(clientOpts *options.ClientOptions) (*Client, error) {
	if c.internalClientFLE != nil {
		return c.internalClientFLE, nil
	}

	internalClientOpts := options.MergeClientOptions(clientOpts)
	internalClientOpts.AutoEncryptionOptions = nil
	internalClientOpts.SetMinPoolSize(0)
	var err error
	c.internalClientFLE, err = NewClient(internalClientOpts)
	return c.internalClientFLE, err
}

func (c *Client) configureKeyVaultClientFLE(clientOpts *options.ClientOptions) error {
	
	var err error
	aeOpts := clientOpts.AutoEncryptionOptions
	switch {
	case aeOpts.KeyVaultClientOptions != nil:
		c.keyVaultClientFLE, err = NewClient(aeOpts.KeyVaultClientOptions)
	case clientOpts.MaxPoolSize != nil && *clientOpts.MaxPoolSize == 0:
		c.keyVaultClientFLE = c
	default:
		c.keyVaultClientFLE, err = c.getOrCreateInternalClient(clientOpts)
	}

	if err != nil {
		return err
	}

	dbName, collName := splitNamespace(aeOpts.KeyVaultNamespace)
	c.keyVaultCollFLE = c.keyVaultClientFLE.Database(dbName).Collection(collName, keyVaultCollOpts)
	return nil
}

func (c *Client) configureMetadataClientFLE(clientOpts *options.ClientOptions) error {
	
	aeOpts := clientOpts.AutoEncryptionOptions
	if aeOpts.BypassAutoEncryption != nil && *aeOpts.BypassAutoEncryption {
		
		return nil
	}
	if clientOpts.MaxPoolSize != nil && *clientOpts.MaxPoolSize == 0 {
		c.metadataClientFLE = c
		return nil
	}

	var err error
	c.metadataClientFLE, err = c.getOrCreateInternalClient(clientOpts)
	return err
}

func (c *Client) newMongoCrypt(opts *options.AutoEncryptionOptions) (*mongocrypt.MongoCrypt, error) {
	
	cryptSchemaMap := make(map[string]bsoncore.Document)
	for k, v := range opts.SchemaMap {
		schema, err := marshal(v, c.bsonOpts, c.registry)
		if err != nil {
			return nil, err
		}
		cryptSchemaMap[k] = schema
	}

	
	cryptEncryptedFieldsMap := make(map[string]bsoncore.Document)
	for k, v := range opts.EncryptedFieldsMap {
		encryptedFields, err := marshal(v, c.bsonOpts, c.registry)
		if err != nil {
			return nil, err
		}
		cryptEncryptedFieldsMap[k] = encryptedFields
	}

	kmsProviders, err := marshal(opts.KmsProviders, c.bsonOpts, c.registry)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS providers document: %w", err)
	}

	
	
	cryptSharedLibPath := ""
	if val, ok := opts.ExtraOptions["cryptSharedLibPath"]; ok {
		str, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf(
				`expected AutoEncryption extra option "cryptSharedLibPath" to be a string, but is a %T`, val)
		}
		cryptSharedLibPath = str
	}

	
	
	
	cryptSharedLibDisabled := false
	if v, ok := opts.ExtraOptions["__cryptSharedLibDisabledForTestOnly"]; ok {
		cryptSharedLibDisabled = v.(bool)
	}

	bypassAutoEncryption := opts.BypassAutoEncryption != nil && *opts.BypassAutoEncryption
	bypassQueryAnalysis := opts.BypassQueryAnalysis != nil && *opts.BypassQueryAnalysis

	mc, err := mongocrypt.NewMongoCrypt(mcopts.MongoCrypt().
		SetKmsProviders(kmsProviders).
		SetLocalSchemaMap(cryptSchemaMap).
		SetBypassQueryAnalysis(bypassQueryAnalysis).
		SetEncryptedFieldsMap(cryptEncryptedFieldsMap).
		SetCryptSharedLibDisabled(cryptSharedLibDisabled || bypassAutoEncryption).
		SetCryptSharedLibOverridePath(cryptSharedLibPath).
		SetHTTPClient(opts.HTTPClient))
	if err != nil {
		return nil, err
	}

	var cryptSharedLibRequired bool
	if val, ok := opts.ExtraOptions["cryptSharedLibRequired"]; ok {
		b, ok := val.(bool)
		if !ok {
			return nil, fmt.Errorf(
				`expected AutoEncryption extra option "cryptSharedLibRequired" to be a bool, but is a %T`, val)
		}
		cryptSharedLibRequired = b
	}

	
	
	
	if cryptSharedLibRequired && mc.CryptSharedLibVersionString() == "" {
		return nil, errors.New(
			`AutoEncryption extra option "cryptSharedLibRequired" is true, but we failed to load the crypt_shared library`)
	}

	return mc, nil
}


func (c *Client) configureCryptFLE(mc *mongocrypt.MongoCrypt, opts *options.AutoEncryptionOptions) {
	bypass := opts.BypassAutoEncryption != nil && *opts.BypassAutoEncryption
	kr := keyRetriever{coll: c.keyVaultCollFLE}
	var cir collInfoRetriever
	
	
	
	if !bypass {
		cir = collInfoRetriever{client: c.metadataClientFLE}
	}

	c.cryptFLE = driver.NewCrypt(&driver.CryptOptions{
		MongoCrypt:           mc,
		CollInfoFn:           cir.cryptCollInfo,
		KeyFn:                kr.cryptKeys,
		MarkFn:               c.mongocryptdFLE.markCommand,
		TLSConfig:            opts.TLSConfig,
		BypassAutoEncryption: bypass,
	})
}


func (c *Client) validSession(sess *session.Client) error {
	if sess != nil && sess.ClientID != c.id {
		return ErrWrongClient
	}
	return nil
}


func (c *Client) Database(name string, opts ...*options.DatabaseOptions) *Database {
	return newDatabase(c, name, opts...)
}










func (c *Client) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (ListDatabasesResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)

	err := c.validSession(sess)
	if err != nil {
		return ListDatabasesResult{}, err
	}
	if sess == nil && c.sessionPool != nil {
		sess = session.NewImplicitClientSession(c.sessionPool, c.id)
		defer sess.EndSession()
	}

	err = c.validSession(sess)
	if err != nil {
		return ListDatabasesResult{}, err
	}

	filterDoc, err := marshal(filter, c.bsonOpts, c.registry)
	if err != nil {
		return ListDatabasesResult{}, err
	}

	selector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(c.localThreshold),
	})
	selector = makeReadPrefSelector(sess, selector, c.localThreshold)

	ldo := options.MergeListDatabasesOptions(opts...)
	op := operation.NewListDatabases(filterDoc).
		Session(sess).ReadPreference(c.readPreference).CommandMonitor(c.monitor).
		ServerSelector(selector).ClusterClock(c.clock).Database("admin").Deployment(c.deployment).Crypt(c.cryptFLE).
		ServerAPI(c.serverAPI).Timeout(c.timeout).Authenticator(c.authenticator)

	if ldo.NameOnly != nil {
		op = op.NameOnly(*ldo.NameOnly)
	}
	if ldo.AuthorizedDatabases != nil {
		op = op.AuthorizedDatabases(*ldo.AuthorizedDatabases)
	}

	retry := driver.RetryNone
	if c.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		return ListDatabasesResult{}, replaceErrors(err)
	}

	return newListDatabasesResultFromOperation(op.Result()), nil
}












func (c *Client) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	opts = append(opts, options.ListDatabases().SetNameOnly(true))

	res, err := c.ListDatabases(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, spec := range res.Databases {
		names = append(names, spec.Name)
	}

	return names, nil
}











func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error {
	return fn(NewSessionContext(ctx, sess))
}












func (c *Client) UseSession(ctx context.Context, fn func(SessionContext) error) error {
	return c.UseSessionWithOptions(ctx, options.Session(), fn)
}





func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(SessionContext) error) error {
	defaultSess, err := c.StartSession(opts)
	if err != nil {
		return err
	}

	defer defaultSess.EndSession(ctx)
	return fn(NewSessionContext(ctx, defaultSess))
}














func (c *Client) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {
	if c.sessionPool == nil {
		return nil, ErrClientDisconnected
	}

	csConfig := changeStreamConfig{
		readConcern:    c.readConcern,
		readPreference: c.readPreference,
		client:         c,
		bsonOpts:       c.bsonOpts,
		registry:       c.registry,
		streamType:     ClientStream,
		crypt:          c.cryptFLE,
	}

	return newChangeStream(ctx, csConfig, pipeline, opts...)
}



func (c *Client) NumberSessionsInProgress() int {
	
	
	
	return int(c.sessionPool.CheckedOut())
}


func (c *Client) Timeout() *time.Duration {
	return c.timeout
}

func (c *Client) createBaseCursorOptions() driver.CursorOptions {
	return driver.CursorOptions{
		CommandMonitor: c.monitor,
		Crypt:          c.cryptFLE,
		ServerAPI:      c.serverAPI,
	}
}



func newLogger(opts *options.LoggerOptions) (*logger.Logger, error) {
	
	if opts == nil {
		opts = options.Logger()
	}

	
	
	if (len(opts.ComponentLevels) == 0) &&
		!logger.EnvHasComponentVariables() {

		return nil, nil
	}

	
	componentLevels := make(map[logger.Component]logger.Level)
	for component, level := range opts.ComponentLevels {
		componentLevels[logger.Component(component)] = logger.Level(level)
	}

	return logger.New(opts.Sink, opts.MaxDocumentLength, componentLevels)
}
