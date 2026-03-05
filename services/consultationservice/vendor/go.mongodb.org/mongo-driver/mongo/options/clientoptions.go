





package options 

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/youmark/pkcs8"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/httputil"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/tag"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const (
	
	
	
	
	
	
	ServerMonitoringModeAuto = connstring.ServerMonitoringModeAuto

	
	
	
	ServerMonitoringModePoll = connstring.ServerMonitoringModePoll

	
	
	
	ServerMonitoringModeStream = connstring.ServerMonitoringModeStream
)





type ContextDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}











































type Credential struct {
	AuthMechanism           string
	AuthMechanismProperties map[string]string
	AuthSource              string
	Username                string
	Password                string
	PasswordSet             bool
	OIDCMachineCallback     OIDCCallback
	OIDCHumanCallback       OIDCCallback
}



type OIDCCallback func(context.Context, *OIDCArgs) (*OIDCCredential, error)


type OIDCArgs struct {
	Version      int
	IDPInfo      *IDPInfo
	RefreshToken *string
}


type OIDCCredential struct {
	AccessToken  string
	ExpiresAt    *time.Time
	RefreshToken *string
}



type IDPInfo struct {
	Issuer        string
	ClientID      string
	RequestScopes []string
}


type BSONOptions struct {
	
	
	UseJSONStructTags bool

	
	
	
	ErrorOnInlineDuplicates bool

	
	
	
	
	IntMinSize bool

	
	
	
	
	
	
	NilMapAsEmpty bool

	
	
	
	
	
	
	NilSliceAsEmpty bool

	
	
	NilByteSliceAsEmpty bool

	
	
	
	OmitZeroStruct bool

	
	
	
	StringifyMapKeysWithFmt bool

	
	
	
	
	AllowTruncatingDoubles bool

	
	
	
	BinaryAsSlice bool

	
	
	
	DefaultDocumentD bool

	
	
	
	DefaultDocumentM bool

	
	
	UseLocalTimeZone bool

	
	
	ZeroMaps bool

	
	
	
	ZeroStructs bool
}



type ClientOptions struct {
	AppName                  *string
	Auth                     *Credential
	AutoEncryptionOptions    *AutoEncryptionOptions
	ConnectTimeout           *time.Duration
	Compressors              []string
	Dialer                   ContextDialer
	Direct                   *bool
	DisableOCSPEndpointCheck *bool
	HeartbeatInterval        *time.Duration
	Hosts                    []string
	HTTPClient               *http.Client
	LoadBalanced             *bool
	LocalThreshold           *time.Duration
	LoggerOptions            *LoggerOptions
	MaxConnIdleTime          *time.Duration
	MaxPoolSize              *uint64
	MinPoolSize              *uint64
	MaxConnecting            *uint64
	PoolMonitor              *event.PoolMonitor
	Monitor                  *event.CommandMonitor
	ServerMonitor            *event.ServerMonitor
	ReadConcern              *readconcern.ReadConcern
	ReadPreference           *readpref.ReadPref
	BSONOptions              *BSONOptions
	Registry                 *bsoncodec.Registry
	ReplicaSet               *string
	RetryReads               *bool
	RetryWrites              *bool
	ServerAPIOptions         *ServerAPIOptions
	ServerMonitoringMode     *string
	ServerSelectionTimeout   *time.Duration
	SRVMaxHosts              *int
	SRVServiceName           *string
	Timeout                  *time.Duration
	TLSConfig                *tls.Config
	WriteConcern             *writeconcern.WriteConcern
	ZlibLevel                *int
	ZstdLevel                *int

	err error
	cs  *connstring.ConnString

	
	
	
	
	AuthenticateToAnything *bool

	
	
	
	
	
	Crypt driver.Crypt

	
	
	
	
	Deployment driver.Deployment

	
	
	
	
	
	SocketTimeout *time.Duration
}


func Client() *ClientOptions {
	return &ClientOptions{
		HTTPClient: httputil.DefaultHTTPClient,
	}
}


func (c *ClientOptions) Validate() error {
	if c.err != nil {
		return c.err
	}
	c.err = c.validate()
	return c.err
}

func (c *ClientOptions) validate() error {
	
	if c.Direct != nil && *c.Direct {
		if len(c.Hosts) > 1 {
			return errors.New("a direct connection cannot be made if multiple hosts are specified")
		}
		if c.cs != nil && c.cs.Scheme == connstring.SchemeMongoDBSRV {
			return errors.New("a direct connection cannot be made if an SRV URI is used")
		}
	}

	if c.MaxPoolSize != nil && c.MinPoolSize != nil && *c.MaxPoolSize != 0 && *c.MinPoolSize > *c.MaxPoolSize {
		return fmt.Errorf("minPoolSize must be less than or equal to maxPoolSize, got minPoolSize=%d maxPoolSize=%d", *c.MinPoolSize, *c.MaxPoolSize)
	}

	
	if c.ServerAPIOptions != nil {
		if err := c.ServerAPIOptions.ServerAPIVersion.Validate(); err != nil {
			return err
		}
	}

	
	if c.LoadBalanced != nil && *c.LoadBalanced {
		if len(c.Hosts) > 1 {
			return connstring.ErrLoadBalancedWithMultipleHosts
		}
		if c.ReplicaSet != nil {
			return connstring.ErrLoadBalancedWithReplicaSet
		}
		if c.Direct != nil && *c.Direct {
			return connstring.ErrLoadBalancedWithDirectConnection
		}
	}

	
	if c.SRVMaxHosts != nil && *c.SRVMaxHosts > 0 {
		if c.ReplicaSet != nil {
			return connstring.ErrSRVMaxHostsWithReplicaSet
		}
		if c.LoadBalanced != nil && *c.LoadBalanced {
			return connstring.ErrSRVMaxHostsWithLoadBalanced
		}
	}

	if mode := c.ServerMonitoringMode; mode != nil && !connstring.IsValidServerMonitoringMode(*mode) {
		return fmt.Errorf("invalid server monitoring mode: %q", *mode)
	}

	
	if c.Auth != nil && c.Auth.AuthMechanism == auth.MongoDBOIDC {
		if c.Auth.Password != "" {
			return fmt.Errorf("password must not be set for the %s auth mechanism", auth.MongoDBOIDC)
		}
		if c.Auth.OIDCMachineCallback != nil && c.Auth.OIDCHumanCallback != nil {
			return fmt.Errorf("cannot set both OIDCMachineCallback and OIDCHumanCallback, only one may be specified")
		}
		if c.Auth.OIDCHumanCallback == nil && c.Auth.AuthMechanismProperties[auth.AllowedHostsProp] != "" {
			return fmt.Errorf("Cannot specify ALLOWED_HOSTS without an OIDCHumanCallback")
		}
		if env, ok := c.Auth.AuthMechanismProperties[auth.EnvironmentProp]; ok {
			switch env {
			case auth.GCPEnvironmentValue, auth.AzureEnvironmentValue:
				if c.Auth.OIDCMachineCallback != nil {
					return fmt.Errorf("OIDCMachineCallback cannot be specified with the %s %q", env, auth.EnvironmentProp)
				}
				if c.Auth.OIDCHumanCallback != nil {
					return fmt.Errorf("OIDCHumanCallback cannot be specified with the %s %q", env, auth.EnvironmentProp)
				}
				if c.Auth.AuthMechanismProperties[auth.ResourceProp] == "" {
					return fmt.Errorf("%q must be set for the %s %q", auth.ResourceProp, env, auth.EnvironmentProp)
				}
			default:
				if c.Auth.AuthMechanismProperties[auth.ResourceProp] != "" {
					return fmt.Errorf("%q must not be set for the %s %q", auth.ResourceProp, env, auth.EnvironmentProp)
				}
			}
		}
	}

	return nil
}



func (c *ClientOptions) GetURI() string {
	if c.cs == nil {
		return ""
	}
	return c.cs.Original
}















func (c *ClientOptions) ApplyURI(uri string) *ClientOptions {
	if c.err != nil {
		return c
	}

	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		c.err = err
		return c
	}
	c.cs = cs

	if cs.AppName != "" {
		c.AppName = &cs.AppName
	}

	
	if cs.HasAuthParameters() {
		c.Auth = &Credential{
			AuthMechanism:           cs.AuthMechanism,
			AuthMechanismProperties: cs.AuthMechanismProperties,
			AuthSource:              cs.AuthSource,
			Username:                cs.Username,
			Password:                cs.Password,
			PasswordSet:             cs.PasswordSet,
		}
	}

	if cs.ConnectSet {
		direct := cs.Connect == connstring.SingleConnect
		c.Direct = &direct
	}

	if cs.DirectConnectionSet {
		c.Direct = &cs.DirectConnection
	}

	if cs.ConnectTimeoutSet {
		c.ConnectTimeout = &cs.ConnectTimeout
	}

	if len(cs.Compressors) > 0 {
		c.Compressors = cs.Compressors
		if stringSliceContains(c.Compressors, "zlib") {
			defaultLevel := wiremessage.DefaultZlibLevel
			c.ZlibLevel = &defaultLevel
		}
		if stringSliceContains(c.Compressors, "zstd") {
			defaultLevel := wiremessage.DefaultZstdLevel
			c.ZstdLevel = &defaultLevel
		}
	}

	if cs.HeartbeatIntervalSet {
		c.HeartbeatInterval = &cs.HeartbeatInterval
	}

	c.Hosts = cs.Hosts

	if cs.LoadBalancedSet {
		c.LoadBalanced = &cs.LoadBalanced
	}

	if cs.LocalThresholdSet {
		c.LocalThreshold = &cs.LocalThreshold
	}

	if cs.MaxConnIdleTimeSet {
		c.MaxConnIdleTime = &cs.MaxConnIdleTime
	}

	if cs.MaxPoolSizeSet {
		c.MaxPoolSize = &cs.MaxPoolSize
	}

	if cs.MinPoolSizeSet {
		c.MinPoolSize = &cs.MinPoolSize
	}

	if cs.MaxConnectingSet {
		c.MaxConnecting = &cs.MaxConnecting
	}

	if cs.ReadConcernLevel != "" {
		c.ReadConcern = readconcern.New(readconcern.Level(cs.ReadConcernLevel))
	}

	if cs.ReadPreference != "" || len(cs.ReadPreferenceTagSets) > 0 || cs.MaxStalenessSet {
		opts := make([]readpref.Option, 0, 1)

		tagSets := tag.NewTagSetsFromMaps(cs.ReadPreferenceTagSets)
		if len(tagSets) > 0 {
			opts = append(opts, readpref.WithTagSets(tagSets...))
		}

		if cs.MaxStaleness != 0 {
			opts = append(opts, readpref.WithMaxStaleness(cs.MaxStaleness))
		}

		mode, err := readpref.ModeFromString(cs.ReadPreference)
		if err != nil {
			c.err = err
			return c
		}

		c.ReadPreference, c.err = readpref.New(mode, opts...)
		if c.err != nil {
			return c
		}
	}

	if cs.RetryWritesSet {
		c.RetryWrites = &cs.RetryWrites
	}

	if cs.RetryReadsSet {
		c.RetryReads = &cs.RetryReads
	}

	if cs.ReplicaSet != "" {
		c.ReplicaSet = &cs.ReplicaSet
	}

	if cs.ServerSelectionTimeoutSet {
		c.ServerSelectionTimeout = &cs.ServerSelectionTimeout
	}

	if cs.SocketTimeoutSet {
		c.SocketTimeout = &cs.SocketTimeout
	}

	if cs.SRVMaxHosts != 0 {
		c.SRVMaxHosts = &cs.SRVMaxHosts
	}

	if cs.SRVServiceName != "" {
		c.SRVServiceName = &cs.SRVServiceName
	}

	if cs.SSL {
		tlsConfig := new(tls.Config)

		if cs.SSLCaFileSet {
			c.err = addCACertFromFile(tlsConfig, cs.SSLCaFile)
			if c.err != nil {
				return c
			}
		}

		if cs.SSLInsecure {
			tlsConfig.InsecureSkipVerify = true
		}

		var x509Subject string
		var keyPasswd string
		if cs.SSLClientCertificateKeyPasswordSet && cs.SSLClientCertificateKeyPassword != nil {
			keyPasswd = cs.SSLClientCertificateKeyPassword()
		}
		if cs.SSLClientCertificateKeyFileSet {
			x509Subject, err = addClientCertFromConcatenatedFile(tlsConfig, cs.SSLClientCertificateKeyFile, keyPasswd)
		} else if cs.SSLCertificateFileSet || cs.SSLPrivateKeyFileSet {
			x509Subject, err = addClientCertFromSeparateFiles(tlsConfig, cs.SSLCertificateFile,
				cs.SSLPrivateKeyFile, keyPasswd)
		}
		if err != nil {
			c.err = err
			return c
		}

		
		if c.Auth != nil && strings.ToLower(c.Auth.AuthMechanism) == "mongodb-x509" &&
			c.Auth.Username == "" {

			
			c.Auth.Username = extractX509UsernameFromSubject(x509Subject)
		}

		c.TLSConfig = tlsConfig
	}

	if cs.JSet || cs.WString != "" || cs.WNumberSet || cs.WTimeoutSet {
		opts := make([]writeconcern.Option, 0, 1)

		if len(cs.WString) > 0 {
			opts = append(opts, writeconcern.WTagSet(cs.WString))
		} else if cs.WNumberSet {
			opts = append(opts, writeconcern.W(cs.WNumber))
		}

		if cs.JSet {
			opts = append(opts, writeconcern.J(cs.J))
		}

		if cs.WTimeoutSet {
			opts = append(opts, writeconcern.WTimeout(cs.WTimeout))
		}

		c.WriteConcern = writeconcern.New(opts...)
	}

	if cs.ZlibLevelSet {
		c.ZlibLevel = &cs.ZlibLevel
	}
	if cs.ZstdLevelSet {
		c.ZstdLevel = &cs.ZstdLevel
	}

	if cs.SSLDisableOCSPEndpointCheckSet {
		c.DisableOCSPEndpointCheck = &cs.SSLDisableOCSPEndpointCheck
	}

	if cs.TimeoutSet {
		c.Timeout = &cs.Timeout
	}

	return c
}




func (c *ClientOptions) SetAppName(s string) *ClientOptions {
	c.AppName = &s
	return c
}




func (c *ClientOptions) SetAuth(auth Credential) *ClientOptions {
	c.Auth = &auth
	return c
}

















func (c *ClientOptions) SetCompressors(comps []string) *ClientOptions {
	c.Compressors = comps

	return c
}




func (c *ClientOptions) SetConnectTimeout(d time.Duration) *ClientOptions {
	c.ConnectTimeout = &d
	return c
}




func (c *ClientOptions) SetDialer(d ContextDialer) *ClientOptions {
	c.Dialer = d
	return c
}
















func (c *ClientOptions) SetDirect(b bool) *ClientOptions {
	c.Direct = &b
	return c
}



func (c *ClientOptions) SetHeartbeatInterval(d time.Duration) *ClientOptions {
	c.HeartbeatInterval = &d
	return c
}






func (c *ClientOptions) SetHosts(s []string) *ClientOptions {
	c.Hosts = s
	return c
}











func (c *ClientOptions) SetLoadBalanced(lb bool) *ClientOptions {
	c.LoadBalanced = &lb
	return c
}





func (c *ClientOptions) SetLocalThreshold(d time.Duration) *ClientOptions {
	c.LocalThreshold = &d
	return c
}



func (c *ClientOptions) SetLoggerOptions(opts *LoggerOptions) *ClientOptions {
	c.LoggerOptions = opts

	return c
}




func (c *ClientOptions) SetMaxConnIdleTime(d time.Duration) *ClientOptions {
	c.MaxConnIdleTime = &d
	return c
}




func (c *ClientOptions) SetMaxPoolSize(u uint64) *ClientOptions {
	c.MaxPoolSize = &u
	return c
}




func (c *ClientOptions) SetMinPoolSize(u uint64) *ClientOptions {
	c.MinPoolSize = &u
	return c
}




func (c *ClientOptions) SetMaxConnecting(u uint64) *ClientOptions {
	c.MaxConnecting = &u
	return c
}



func (c *ClientOptions) SetPoolMonitor(m *event.PoolMonitor) *ClientOptions {
	c.PoolMonitor = m
	return c
}



func (c *ClientOptions) SetMonitor(m *event.CommandMonitor) *ClientOptions {
	c.Monitor = m
	return c
}


func (c *ClientOptions) SetServerMonitor(m *event.ServerMonitor) *ClientOptions {
	c.ServerMonitor = m
	return c
}




func (c *ClientOptions) SetReadConcern(rc *readconcern.ReadConcern) *ClientOptions {
	c.ReadConcern = rc

	return c
}














func (c *ClientOptions) SetReadPreference(rp *readpref.ReadPref) *ClientOptions {
	c.ReadPreference = rp

	return c
}


func (c *ClientOptions) SetBSONOptions(opts *BSONOptions) *ClientOptions {
	c.BSONOptions = opts
	return c
}



func (c *ClientOptions) SetRegistry(registry *bsoncodec.Registry) *ClientOptions {
	c.Registry = registry
	return c
}






func (c *ClientOptions) SetReplicaSet(s string) *ClientOptions {
	c.ReplicaSet = &s
	return c
}












func (c *ClientOptions) SetRetryWrites(b bool) *ClientOptions {
	c.RetryWrites = &b

	return c
}









func (c *ClientOptions) SetRetryReads(b bool) *ClientOptions {
	c.RetryReads = &b
	return c
}




func (c *ClientOptions) SetServerSelectionTimeout(d time.Duration) *ClientOptions {
	c.ServerSelectionTimeout = &d
	return c
}








func (c *ClientOptions) SetSocketTimeout(d time.Duration) *ClientOptions {
	c.SocketTimeout = &d
	return c
}












func (c *ClientOptions) SetTimeout(d time.Duration) *ClientOptions {
	c.Timeout = &d
	return c
}

























func (c *ClientOptions) SetTLSConfig(cfg *tls.Config) *ClientOptions {
	c.TLSConfig = cfg
	return c
}




func (c *ClientOptions) SetHTTPClient(client *http.Client) *ClientOptions {
	c.HTTPClient = client
	return c
}















func (c *ClientOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *ClientOptions {
	c.WriteConcern = wc

	return c
}





func (c *ClientOptions) SetZlibLevel(level int) *ClientOptions {
	c.ZlibLevel = &level

	return c
}




func (c *ClientOptions) SetZstdLevel(level int) *ClientOptions {
	c.ZstdLevel = &level
	return c
}




func (c *ClientOptions) SetAutoEncryptionOptions(opts *AutoEncryptionOptions) *ClientOptions {
	c.AutoEncryptionOptions = opts
	return c
}










func (c *ClientOptions) SetDisableOCSPEndpointCheck(disableCheck bool) *ClientOptions {
	c.DisableOCSPEndpointCheck = &disableCheck
	return c
}




func (c *ClientOptions) SetServerAPIOptions(opts *ServerAPIOptions) *ClientOptions {
	c.ServerAPIOptions = opts
	return c
}





func (c *ClientOptions) SetServerMonitoringMode(mode string) *ClientOptions {
	c.ServerMonitoringMode = &mode

	return c
}




func (c *ClientOptions) SetSRVMaxHosts(srvMaxHosts int) *ClientOptions {
	c.SRVMaxHosts = &srvMaxHosts
	return c
}




func (c *ClientOptions) SetSRVServiceName(srvName string) *ClientOptions {
	c.SRVServiceName = &srvName
	return c
}







func MergeClientOptions(opts ...*ClientOptions) *ClientOptions {
	c := Client()

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.Dialer != nil {
			c.Dialer = opt.Dialer
		}
		if opt.AppName != nil {
			c.AppName = opt.AppName
		}
		if opt.Auth != nil {
			c.Auth = opt.Auth
		}
		if opt.AuthenticateToAnything != nil {
			c.AuthenticateToAnything = opt.AuthenticateToAnything
		}
		if opt.Compressors != nil {
			c.Compressors = opt.Compressors
		}
		if opt.ConnectTimeout != nil {
			c.ConnectTimeout = opt.ConnectTimeout
		}
		if opt.Crypt != nil {
			c.Crypt = opt.Crypt
		}
		if opt.HeartbeatInterval != nil {
			c.HeartbeatInterval = opt.HeartbeatInterval
		}
		if len(opt.Hosts) > 0 {
			c.Hosts = opt.Hosts
		}
		if opt.HTTPClient != nil {
			c.HTTPClient = opt.HTTPClient
		}
		if opt.LoadBalanced != nil {
			c.LoadBalanced = opt.LoadBalanced
		}
		if opt.LocalThreshold != nil {
			c.LocalThreshold = opt.LocalThreshold
		}
		if opt.MaxConnIdleTime != nil {
			c.MaxConnIdleTime = opt.MaxConnIdleTime
		}
		if opt.MaxPoolSize != nil {
			c.MaxPoolSize = opt.MaxPoolSize
		}
		if opt.MinPoolSize != nil {
			c.MinPoolSize = opt.MinPoolSize
		}
		if opt.MaxConnecting != nil {
			c.MaxConnecting = opt.MaxConnecting
		}
		if opt.PoolMonitor != nil {
			c.PoolMonitor = opt.PoolMonitor
		}
		if opt.Monitor != nil {
			c.Monitor = opt.Monitor
		}
		if opt.ServerAPIOptions != nil {
			c.ServerAPIOptions = opt.ServerAPIOptions
		}
		if opt.ServerMonitor != nil {
			c.ServerMonitor = opt.ServerMonitor
		}
		if opt.ReadConcern != nil {
			c.ReadConcern = opt.ReadConcern
		}
		if opt.ReadPreference != nil {
			c.ReadPreference = opt.ReadPreference
		}
		if opt.BSONOptions != nil {
			c.BSONOptions = opt.BSONOptions
		}
		if opt.Registry != nil {
			c.Registry = opt.Registry
		}
		if opt.ReplicaSet != nil {
			c.ReplicaSet = opt.ReplicaSet
		}
		if opt.RetryWrites != nil {
			c.RetryWrites = opt.RetryWrites
		}
		if opt.RetryReads != nil {
			c.RetryReads = opt.RetryReads
		}
		if opt.ServerSelectionTimeout != nil {
			c.ServerSelectionTimeout = opt.ServerSelectionTimeout
		}
		if opt.Direct != nil {
			c.Direct = opt.Direct
		}
		if opt.SocketTimeout != nil {
			c.SocketTimeout = opt.SocketTimeout
		}
		if opt.SRVMaxHosts != nil {
			c.SRVMaxHosts = opt.SRVMaxHosts
		}
		if opt.SRVServiceName != nil {
			c.SRVServiceName = opt.SRVServiceName
		}
		if opt.Timeout != nil {
			c.Timeout = opt.Timeout
		}
		if opt.TLSConfig != nil {
			c.TLSConfig = opt.TLSConfig
		}
		if opt.WriteConcern != nil {
			c.WriteConcern = opt.WriteConcern
		}
		if opt.ZlibLevel != nil {
			c.ZlibLevel = opt.ZlibLevel
		}
		if opt.ZstdLevel != nil {
			c.ZstdLevel = opt.ZstdLevel
		}
		if opt.AutoEncryptionOptions != nil {
			c.AutoEncryptionOptions = opt.AutoEncryptionOptions
		}
		if opt.Deployment != nil {
			c.Deployment = opt.Deployment
		}
		if opt.DisableOCSPEndpointCheck != nil {
			c.DisableOCSPEndpointCheck = opt.DisableOCSPEndpointCheck
		}
		if opt.err != nil {
			c.err = opt.err
		}
		if opt.cs != nil {
			c.cs = opt.cs
		}
		if opt.LoggerOptions != nil {
			c.LoggerOptions = opt.LoggerOptions
		}
		if opt.ServerMonitoringMode != nil {
			c.ServerMonitoringMode = opt.ServerMonitoringMode
		}
	}

	return c
}



func addCACertFromFile(cfg *tls.Config, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if cfg.RootCAs == nil {
		cfg.RootCAs = x509.NewCertPool()
	}
	if !cfg.RootCAs.AppendCertsFromPEM(data) {
		return errors.New("the specified CA file does not contain any valid certificates")
	}

	return nil
}

func addClientCertFromSeparateFiles(cfg *tls.Config, keyFile, certFile, keyPassword string) (string, error) {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", err
	}
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return "", err
	}

	keySize := len(keyData)
	if keySize > 64*1024*1024 {
		return "", errors.New("X.509 key must be less than 64 MiB")
	}
	certSize := len(certData)
	if certSize > 64*1024*1024 {
		return "", errors.New("X.509 certificate must be less than 64 MiB")
	}
	dataSize := keySize + certSize + 1
	if dataSize > math.MaxInt {
		return "", errors.New("size overflow")
	}
	data := make([]byte, 0, dataSize)
	data = append(data, keyData...)
	data = append(data, '\n')
	data = append(data, certData...)
	return addClientCertFromBytes(cfg, data, keyPassword)
}

func addClientCertFromConcatenatedFile(cfg *tls.Config, certKeyFile, keyPassword string) (string, error) {
	data, err := ioutil.ReadFile(certKeyFile)
	if err != nil {
		return "", err
	}

	return addClientCertFromBytes(cfg, data, keyPassword)
}



func addClientCertFromBytes(cfg *tls.Config, data []byte, keyPasswd string) (string, error) {
	var currentBlock *pem.Block
	var certDecodedBlock []byte
	var certBlocks, keyBlocks [][]byte

	remaining := data
	start := 0
	for {
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil {
			break
		}

		if currentBlock.Type == "CERTIFICATE" {
			certBlock := data[start : len(data)-len(remaining)]
			certBlocks = append(certBlocks, certBlock)
			
			
			if certDecodedBlock == nil {
				certDecodedBlock = currentBlock.Bytes
			}
			start += len(certBlock)
		} else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") {
			isEncrypted := x509.IsEncryptedPEMBlock(currentBlock) || strings.Contains(currentBlock.Type, "ENCRYPTED PRIVATE KEY")
			if isEncrypted {
				if keyPasswd == "" {
					return "", fmt.Errorf("no password provided to decrypt private key")
				}

				var keyBytes []byte
				var err error
				
				if x509.IsEncryptedPEMBlock(currentBlock) {
					
					keyBytes, err = x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
					if err != nil {
						return "", err
					}
				} else if strings.Contains(currentBlock.Type, "ENCRYPTED") {
					
					decrypted, err := pkcs8.ParsePKCS8PrivateKey(currentBlock.Bytes, []byte(keyPasswd))
					if err != nil {
						return "", err
					}
					keyBytes, err = x509.MarshalPKCS8PrivateKey(decrypted)
					if err != nil {
						return "", err
					}
				}
				var encoded bytes.Buffer
				err = pem.Encode(&encoded, &pem.Block{Type: currentBlock.Type, Bytes: keyBytes})
				if err != nil {
					return "", fmt.Errorf("error encoding private key as PEM: %w", err)
				}
				keyBlock := encoded.Bytes()
				keyBlocks = append(keyBlocks, keyBlock)
				start = len(data) - len(remaining)
			} else {
				keyBlock := data[start : len(data)-len(remaining)]
				keyBlocks = append(keyBlocks, keyBlock)
				start += len(keyBlock)
			}
		}
	}
	if len(certBlocks) == 0 {
		return "", fmt.Errorf("failed to find CERTIFICATE")
	}
	if len(keyBlocks) == 0 {
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	}

	cert, err := tls.X509KeyPair(bytes.Join(certBlocks, []byte("\n")), bytes.Join(keyBlocks, []byte("\n")))
	if err != nil {
		return "", err
	}

	cfg.Certificates = append(cfg.Certificates, cert)

	
	
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil {
		return "", err
	}

	return crt.Subject.String(), nil
}

func stringSliceContains(source []string, target string) bool {
	for _, str := range source {
		if str == target {
			return true
		}
	}
	return false
}


func extractX509UsernameFromSubject(subject string) string {
	
	pairs := strings.Split(subject, ",")
	for left, right := 0, len(pairs)-1; left < right; left, right = left+1, right-1 {
		pairs[left], pairs[right] = pairs[right], pairs[left]
	}

	return strings.Join(pairs, ",")
}
