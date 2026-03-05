





package mongo

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	defaultServerSelectionTimeout = 10 * time.Second
	defaultURI                    = "mongodb://localhost:27020"
	defaultPath                   = "mongocryptd"
	serverSelectionTimeoutStr     = "server selection error"
)

var defaultTimeoutArgs = []string{"--idleShutdownTimeoutSecs=60"}
var databaseOpts = options.Database().SetReadConcern(readconcern.New()).SetReadPreference(readpref.Primary())

type mongocryptdClient struct {
	bypassSpawn bool
	client      *Client
	path        string
	spawnArgs   []string
}




func newMongocryptdClient(opts *options.AutoEncryptionOptions) (*mongocryptdClient, error) {
	
	var bypassSpawn bool
	var bypassAutoEncryption bool

	if bypass, ok := opts.ExtraOptions["mongocryptdBypassSpawn"]; ok {
		bypassSpawn = bypass.(bool)
	}
	if opts.BypassAutoEncryption != nil {
		bypassAutoEncryption = *opts.BypassAutoEncryption
	}

	bypassQueryAnalysis := opts.BypassQueryAnalysis != nil && *opts.BypassQueryAnalysis

	mc := &mongocryptdClient{
		
		
		
		
		bypassSpawn: bypassSpawn || bypassAutoEncryption || bypassQueryAnalysis,
	}

	if !mc.bypassSpawn {
		mc.path, mc.spawnArgs = createSpawnArgs(opts.ExtraOptions)
		if err := mc.spawnProcess(); err != nil {
			return nil, err
		}
	}

	
	uri := defaultURI
	if u, ok := opts.ExtraOptions["mongocryptdURI"]; ok {
		uri = u.(string)
	}

	
	client, err := NewClient(options.Client().ApplyURI(uri).SetServerSelectionTimeout(defaultServerSelectionTimeout))
	if err != nil {
		return nil, err
	}
	mc.client = client

	return mc, nil
}


func (mc *mongocryptdClient) markCommand(ctx context.Context, dbName string, cmd bsoncore.Document) (bsoncore.Document, error) {
	
	
	
	ctx = NewSessionContext(ctx, nil)
	db := mc.client.Database(dbName, databaseOpts)

	res, err := db.RunCommand(ctx, cmd).Raw()
	
	if err == nil {
		return bsoncore.Document(res), nil
	}
	
	if mc.bypassSpawn || !strings.Contains(err.Error(), serverSelectionTimeoutStr) {
		return nil, MongocryptdError{Wrapped: err}
	}

	
	if err = mc.spawnProcess(); err != nil {
		return nil, err
	}
	res, err = db.RunCommand(ctx, cmd).Raw()
	if err != nil {
		return nil, MongocryptdError{Wrapped: err}
	}
	return bsoncore.Document(res), nil
}


func (mc *mongocryptdClient) connect(ctx context.Context) error {
	return mc.client.Connect(ctx)
}


func (mc *mongocryptdClient) disconnect(ctx context.Context) error {
	return mc.client.Disconnect(ctx)
}

func (mc *mongocryptdClient) spawnProcess() error {
	
	
	cmd := exec.Command(mc.path, mc.spawnArgs...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}


func createSpawnArgs(opts map[string]interface{}) (string, []string) {
	var spawnArgs []string

	
	path := defaultPath
	if p, ok := opts["mongocryptdPath"]; ok {
		path = p.(string)
	}

	
	if sa, ok := opts["mongocryptdSpawnArgs"]; ok {
		spawnArgs = append(spawnArgs, sa.([]string)...)
	}

	
	var foundTimeout bool
	for _, arg := range spawnArgs {
		
		
		if strings.HasPrefix(arg, "--idleShutdownTimeoutSecs") {
			foundTimeout = true
			break
		}
	}
	if !foundTimeout {
		spawnArgs = append(spawnArgs, defaultTimeoutArgs...)
	}

	return path, spawnArgs
}
