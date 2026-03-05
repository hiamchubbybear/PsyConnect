












package connstring 

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/internal/randutil"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/dns"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const (
	
	
	
	
	
	
	ServerMonitoringModeAuto = "auto"

	
	
	
	ServerMonitoringModePoll = "poll"

	
	
	
	ServerMonitoringModeStream = "stream"
)

var (
	
	
	ErrLoadBalancedWithMultipleHosts = errors.New(
		"loadBalanced cannot be set to true if multiple hosts are specified")

	
	
	ErrLoadBalancedWithReplicaSet = errors.New(
		"loadBalanced cannot be set to true if a replica set name is specified")

	
	
	ErrLoadBalancedWithDirectConnection = errors.New(
		"loadBalanced cannot be set to true if the direct connection option is specified")

	
	
	ErrSRVMaxHostsWithReplicaSet = errors.New(
		"srvMaxHosts cannot be a positive value if a replica set name is specified")

	
	
	ErrSRVMaxHostsWithLoadBalanced = errors.New(
		"srvMaxHosts cannot be a positive value if loadBalanced is set to true")
)


var random = randutil.NewLockedRand()



func ParseAndValidate(s string) (*ConnString, error) {
	connStr, err := Parse(s)
	if err != nil {
		return nil, err
	}
	err = connStr.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating uri: %w", err)
	}
	return connStr, nil
}




func Parse(s string) (*ConnString, error) {
	p := parser{dnsResolver: dns.DefaultResolver}
	connStr, err := p.parse(s)
	if err != nil {
		return nil, fmt.Errorf("error parsing uri: %w", err)
	}
	return connStr, err
}


type ConnString struct {
	Original                           string
	AppName                            string
	AuthMechanism                      string
	AuthMechanismProperties            map[string]string
	AuthMechanismPropertiesSet         bool
	AuthSource                         string
	AuthSourceSet                      bool
	Compressors                        []string
	Connect                            ConnectMode
	ConnectSet                         bool
	DirectConnection                   bool
	DirectConnectionSet                bool
	ConnectTimeout                     time.Duration
	ConnectTimeoutSet                  bool
	Database                           string
	HeartbeatInterval                  time.Duration
	HeartbeatIntervalSet               bool
	Hosts                              []string
	J                                  bool
	JSet                               bool
	LoadBalanced                       bool
	LoadBalancedSet                    bool
	LocalThreshold                     time.Duration
	LocalThresholdSet                  bool
	MaxConnIdleTime                    time.Duration
	MaxConnIdleTimeSet                 bool
	MaxPoolSize                        uint64
	MaxPoolSizeSet                     bool
	MinPoolSize                        uint64
	MinPoolSizeSet                     bool
	MaxConnecting                      uint64
	MaxConnectingSet                   bool
	Password                           string
	PasswordSet                        bool
	RawHosts                           []string
	ReadConcernLevel                   string
	ReadPreference                     string
	ReadPreferenceTagSets              []map[string]string
	RetryWrites                        bool
	RetryWritesSet                     bool
	RetryReads                         bool
	RetryReadsSet                      bool
	MaxStaleness                       time.Duration
	MaxStalenessSet                    bool
	ReplicaSet                         string
	Scheme                             string
	ServerMonitoringMode               string
	ServerSelectionTimeout             time.Duration
	ServerSelectionTimeoutSet          bool
	SocketTimeout                      time.Duration
	SocketTimeoutSet                   bool
	SRVMaxHosts                        int
	SRVServiceName                     string
	SSL                                bool
	SSLSet                             bool
	SSLClientCertificateKeyFile        string
	SSLClientCertificateKeyFileSet     bool
	SSLClientCertificateKeyPassword    func() string
	SSLClientCertificateKeyPasswordSet bool
	SSLCertificateFile                 string
	SSLCertificateFileSet              bool
	SSLPrivateKeyFile                  string
	SSLPrivateKeyFileSet               bool
	SSLInsecure                        bool
	SSLInsecureSet                     bool
	SSLCaFile                          string
	SSLCaFileSet                       bool
	SSLDisableOCSPEndpointCheck        bool
	SSLDisableOCSPEndpointCheckSet     bool
	Timeout                            time.Duration
	TimeoutSet                         bool
	WString                            string
	WNumber                            int
	WNumberSet                         bool
	Username                           string
	UsernameSet                        bool
	ZlibLevel                          int
	ZlibLevelSet                       bool
	ZstdLevel                          int
	ZstdLevelSet                       bool

	WTimeout              time.Duration
	WTimeoutSet           bool
	WTimeoutSetFromOption bool

	Options        map[string][]string
	UnknownOptions map[string][]string
}

func (u *ConnString) String() string {
	return u.Original
}



func (u *ConnString) HasAuthParameters() bool {
	
	
	return u.AuthMechanism != "" || u.AuthMechanismProperties != nil || u.UsernameSet || u.PasswordSet
}


func (u *ConnString) Validate() error {
	var err error

	if err = u.validateAuth(); err != nil {
		return err
	}

	if err = u.validateSSL(); err != nil {
		return err
	}

	
	if u.WNumberSet && u.WNumber == 0 && u.JSet && u.J {
		return writeconcern.ErrInconsistent
	}

	
	if (u.ConnectSet && u.Connect == SingleConnect) ||
		(u.DirectConnectionSet && u.DirectConnection) {
		if len(u.Hosts) > 1 {
			return errors.New("a direct connection cannot be made if multiple hosts are specified")
		}
		if u.Scheme == SchemeMongoDBSRV {
			return errors.New("a direct connection cannot be made if an SRV URI is used")
		}
		if u.LoadBalancedSet && u.LoadBalanced {
			return ErrLoadBalancedWithDirectConnection
		}
	}

	
	if u.LoadBalancedSet && u.LoadBalanced {
		if len(u.Hosts) > 1 {
			return ErrLoadBalancedWithMultipleHosts
		}
		if u.ReplicaSet != "" {
			return ErrLoadBalancedWithReplicaSet
		}
	}

	
	if u.SRVMaxHosts > 0 {
		if u.ReplicaSet != "" {
			return ErrSRVMaxHostsWithReplicaSet
		}
		if u.LoadBalanced {
			return ErrSRVMaxHostsWithLoadBalanced
		}
	}

	
	if u.AuthMechanism == auth.MongoDBOIDC {
		if _, ok := u.AuthMechanismProperties[auth.AllowedHostsProp]; ok {
			return fmt.Errorf(
				"ALLOWED_HOSTS cannot be specified in the URI connection string for the %q auth mechanism, it must be specified through the ClientOptions directly",
				auth.MongoDBOIDC,
			)
		}
	}

	return nil
}

func (u *ConnString) setDefaultAuthParams(dbName string) error {
	
	
	if u.AuthSourceSet && u.AuthSource == "" {
		return errors.New("authSource must be non-empty when supplied in a URI")
	}

	switch strings.ToLower(u.AuthMechanism) {
	case "plain":
		if u.AuthSource == "" {
			u.AuthSource = dbName
			if u.AuthSource == "" {
				u.AuthSource = "$external"
			}
		}
	case "gssapi":
		if u.AuthMechanismProperties == nil {
			u.AuthMechanismProperties = map[string]string{
				"SERVICE_NAME": "mongodb",
			}
		} else if v, ok := u.AuthMechanismProperties["SERVICE_NAME"]; !ok || v == "" {
			u.AuthMechanismProperties["SERVICE_NAME"] = "mongodb"
		}
		fallthrough
	case "mongodb-aws", "mongodb-x509", "mongodb-oidc":
		if u.AuthSource == "" {
			u.AuthSource = "$external"
		} else if u.AuthSource != "$external" {
			return fmt.Errorf("auth source must be $external")
		}
	case "mongodb-cr":
		fallthrough
	case "scram-sha-1":
		fallthrough
	case "scram-sha-256":
		if u.AuthSource == "" {
			u.AuthSource = dbName
			if u.AuthSource == "" {
				u.AuthSource = "admin"
			}
		}
	case "":
		
		if u.AuthSource == "" && (u.AuthMechanismProperties != nil || u.Username != "" || u.PasswordSet) {
			u.AuthSource = dbName
			if u.AuthSource == "" {
				u.AuthSource = "admin"
			}
		}
	default:
		return fmt.Errorf("invalid auth mechanism")
	}
	return nil
}

func (u *ConnString) addOptions(connectionArgPairs []string) error {
	var tlsssl *bool 
	for _, pair := range connectionArgPairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return fmt.Errorf("invalid option")
		}

		key, err := url.QueryUnescape(kv[0])
		if err != nil {
			return fmt.Errorf("invalid option key %q: %w", kv[0], err)
		}

		value, err := url.QueryUnescape(kv[1])
		if err != nil {
			return fmt.Errorf("invalid option value %q: %w", kv[1], err)
		}

		lowerKey := strings.ToLower(key)
		switch lowerKey {
		case "appname":
			u.AppName = value
		case "authmechanism":
			u.AuthMechanism = value
		case "authmechanismproperties":
			u.AuthMechanismProperties = make(map[string]string)
			pairs := strings.Split(value, ",")
			for _, pair := range pairs {
				kv := strings.SplitN(pair, ":", 2)
				if len(kv) != 2 || kv[0] == "" {
					return fmt.Errorf("invalid authMechanism property")
				}
				u.AuthMechanismProperties[kv[0]] = kv[1]
			}
			u.AuthMechanismPropertiesSet = true
		case "authsource":
			u.AuthSource = value
			u.AuthSourceSet = true
		case "compressors":
			compressors := strings.Split(value, ",")
			if len(compressors) < 1 {
				return fmt.Errorf("must have at least 1 compressor")
			}
			u.Compressors = compressors
		case "connect":
			switch strings.ToLower(value) {
			case "automatic":
			case "direct":
				u.Connect = SingleConnect
			default:
				return fmt.Errorf("invalid 'connect' value: %q", value)
			}
			if u.DirectConnectionSet {
				expectedValue := u.Connect == SingleConnect 
				if u.DirectConnection != expectedValue {
					return fmt.Errorf("options connect=%q and directConnection=%v conflict", value, u.DirectConnection)
				}
			}

			u.ConnectSet = true
		case "directconnection":
			switch strings.ToLower(value) {
			case "true":
				u.DirectConnection = true
			case "false":
			default:
				return fmt.Errorf("invalid 'directConnection' value: %q", value)
			}

			if u.ConnectSet {
				expectedValue := AutoConnect
				if u.DirectConnection {
					expectedValue = SingleConnect
				}

				if u.Connect != expectedValue {
					return fmt.Errorf("options connect=%q and directConnection=%q conflict", u.Connect, value)
				}
			}
			u.DirectConnectionSet = true
		case "connecttimeoutms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.ConnectTimeout = time.Duration(n) * time.Millisecond
			u.ConnectTimeoutSet = true
		case "heartbeatintervalms", "heartbeatfrequencyms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.HeartbeatInterval = time.Duration(n) * time.Millisecond
			u.HeartbeatIntervalSet = true
		case "journal":
			switch value {
			case "true":
				u.J = true
			case "false":
				u.J = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			u.JSet = true
		case "loadbalanced":
			switch value {
			case "true":
				u.LoadBalanced = true
			case "false":
				u.LoadBalanced = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			u.LoadBalancedSet = true
		case "localthresholdms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.LocalThreshold = time.Duration(n) * time.Millisecond
			u.LocalThresholdSet = true
		case "maxidletimems":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.MaxConnIdleTime = time.Duration(n) * time.Millisecond
			u.MaxConnIdleTimeSet = true
		case "maxpoolsize":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.MaxPoolSize = uint64(n)
			u.MaxPoolSizeSet = true
		case "minpoolsize":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.MinPoolSize = uint64(n)
			u.MinPoolSizeSet = true
		case "maxconnecting":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.MaxConnecting = uint64(n)
			u.MaxConnectingSet = true
		case "readconcernlevel":
			u.ReadConcernLevel = value
		case "readpreference":
			u.ReadPreference = value
		case "readpreferencetags":
			if value == "" {
				
				
				u.ReadPreferenceTagSets = append(u.ReadPreferenceTagSets, map[string]string{})
				break
			}

			tags := make(map[string]string)
			items := strings.Split(value, ",")
			for _, item := range items {
				parts := strings.Split(item, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid value for %q: %q", key, value)
				}
				tags[parts[0]] = parts[1]
			}
			u.ReadPreferenceTagSets = append(u.ReadPreferenceTagSets, tags)
		case "maxstaleness", "maxstalenessseconds":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.MaxStaleness = time.Duration(n) * time.Second
			u.MaxStalenessSet = true
		case "replicaset":
			u.ReplicaSet = value
		case "retrywrites":
			switch value {
			case "true":
				u.RetryWrites = true
			case "false":
				u.RetryWrites = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			u.RetryWritesSet = true
		case "retryreads":
			switch value {
			case "true":
				u.RetryReads = true
			case "false":
				u.RetryReads = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			u.RetryReadsSet = true
		case "servermonitoringmode":
			if !IsValidServerMonitoringMode(value) {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			u.ServerMonitoringMode = value
		case "serverselectiontimeoutms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.ServerSelectionTimeout = time.Duration(n) * time.Millisecond
			u.ServerSelectionTimeoutSet = true
		case "sockettimeoutms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.SocketTimeout = time.Duration(n) * time.Millisecond
			u.SocketTimeoutSet = true
		case "srvmaxhosts":
			
			if u.Scheme != SchemeMongoDBSRV {
				return fmt.Errorf("cannot specify srvMaxHosts on non-SRV URI")
			}

			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.SRVMaxHosts = n
		case "srvservicename":
			
			if u.Scheme != SchemeMongoDBSRV {
				return fmt.Errorf("cannot specify srvServiceName on non-SRV URI")
			}

			
			
			
			
			if len(value) < 1 || len(value) > 62 {
				return fmt.Errorf("srvServiceName value must be between 1 and 62 characters")
			}
			u.SRVServiceName = value
		case "ssl", "tls":
			switch value {
			case "true":
				u.SSL = true
			case "false":
				u.SSL = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			if tlsssl == nil {
				tlsssl = new(bool)
				*tlsssl = u.SSL
			} else if *tlsssl != u.SSL {
				return errors.New("tls and ssl options, when both specified, must be equivalent")
			}

			u.SSLSet = true
		case "sslclientcertificatekeyfile", "tlscertificatekeyfile":
			u.SSL = true
			u.SSLSet = true
			u.SSLClientCertificateKeyFile = value
			u.SSLClientCertificateKeyFileSet = true
		case "sslclientcertificatekeypassword", "tlscertificatekeyfilepassword":
			u.SSLClientCertificateKeyPassword = func() string { return value }
			u.SSLClientCertificateKeyPasswordSet = true
		case "tlscertificatefile":
			u.SSL = true
			u.SSLSet = true
			u.SSLCertificateFile = value
			u.SSLCertificateFileSet = true
		case "tlsprivatekeyfile":
			u.SSL = true
			u.SSLSet = true
			u.SSLPrivateKeyFile = value
			u.SSLPrivateKeyFileSet = true
		case "sslinsecure", "tlsinsecure":
			switch value {
			case "true":
				u.SSLInsecure = true
			case "false":
				u.SSLInsecure = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			u.SSLInsecureSet = true
		case "sslcertificateauthorityfile", "tlscafile":
			u.SSL = true
			u.SSLSet = true
			u.SSLCaFile = value
			u.SSLCaFileSet = true
		case "timeoutms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.Timeout = time.Duration(n) * time.Millisecond
			u.TimeoutSet = true
		case "tlsdisableocspendpointcheck":
			u.SSL = true
			u.SSLSet = true

			switch value {
			case "true":
				u.SSLDisableOCSPEndpointCheck = true
			case "false":
				u.SSLDisableOCSPEndpointCheck = false
			default:
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.SSLDisableOCSPEndpointCheckSet = true
		case "w":
			if w, err := strconv.Atoi(value); err == nil {
				if w < 0 {
					return fmt.Errorf("invalid value for %q: %q", key, value)
				}

				u.WNumber = w
				u.WNumberSet = true
				u.WString = ""
				break
			}

			u.WString = value
			u.WNumberSet = false

		case "wtimeoutms":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.WTimeout = time.Duration(n) * time.Millisecond
			u.WTimeoutSet = true
		case "wtimeout":
			
			if u.WTimeoutSet {
				break
			}
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}
			u.WTimeout = time.Duration(n) * time.Millisecond
		case "zlibcompressionlevel":
			level, err := strconv.Atoi(value)
			if err != nil || (level < -1 || level > 9) {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			if level == -1 {
				level = wiremessage.DefaultZlibLevel
			}
			u.ZlibLevel = level
			u.ZlibLevelSet = true
		case "zstdcompressionlevel":
			const maxZstdLevel = 22 
			level, err := strconv.Atoi(value)
			if err != nil || (level < -1 || level > maxZstdLevel) {
				return fmt.Errorf("invalid value for %q: %q", key, value)
			}

			if level == -1 {
				level = wiremessage.DefaultZstdLevel
			}
			u.ZstdLevel = level
			u.ZstdLevelSet = true
		default:
			if u.UnknownOptions == nil {
				u.UnknownOptions = make(map[string][]string)
			}
			u.UnknownOptions[lowerKey] = append(u.UnknownOptions[lowerKey], value)
		}

		if u.Options == nil {
			u.Options = make(map[string][]string)
		}
		u.Options[lowerKey] = append(u.Options[lowerKey], value)
	}
	return nil
}

func (u *ConnString) validateAuth() error {
	switch strings.ToLower(u.AuthMechanism) {
	case "mongodb-cr":
		if u.Username == "" {
			return fmt.Errorf("username required for MONGO-CR")
		}
		if u.Password == "" {
			return fmt.Errorf("password required for MONGO-CR")
		}
		if u.AuthMechanismProperties != nil {
			return fmt.Errorf("MONGO-CR cannot have mechanism properties")
		}
	case "mongodb-x509":
		if u.Password != "" {
			return fmt.Errorf("password cannot be specified for MONGO-X509")
		}
		if u.AuthMechanismProperties != nil {
			return fmt.Errorf("MONGO-X509 cannot have mechanism properties")
		}
	case "mongodb-aws":
		if u.Username != "" && u.Password == "" {
			return fmt.Errorf("username without password is invalid for MONGODB-AWS")
		}
		if u.Username == "" && u.Password != "" {
			return fmt.Errorf("password without username is invalid for MONGODB-AWS")
		}
		var token bool
		for k := range u.AuthMechanismProperties {
			if k != "AWS_SESSION_TOKEN" {
				return fmt.Errorf("invalid auth property for MONGODB-AWS")
			}
			token = true
		}
		if token && u.Username == "" && u.Password == "" {
			return fmt.Errorf("token without username and password is invalid for MONGODB-AWS")
		}
	case "gssapi":
		if u.Username == "" {
			return fmt.Errorf("username required for GSSAPI")
		}
		for k := range u.AuthMechanismProperties {
			if k != "SERVICE_NAME" && k != "CANONICALIZE_HOST_NAME" && k != "SERVICE_REALM" && k != "SERVICE_HOST" {
				return fmt.Errorf("invalid auth property for GSSAPI")
			}
		}
	case "plain":
		if u.Username == "" {
			return fmt.Errorf("username required for PLAIN")
		}
		if u.Password == "" {
			return fmt.Errorf("password required for PLAIN")
		}
		if u.AuthMechanismProperties != nil {
			return fmt.Errorf("PLAIN cannot have mechanism properties")
		}
	case "scram-sha-1":
		if u.Username == "" {
			return fmt.Errorf("username required for SCRAM-SHA-1")
		}
		if u.Password == "" {
			return fmt.Errorf("password required for SCRAM-SHA-1")
		}
		if u.AuthMechanismProperties != nil {
			return fmt.Errorf("SCRAM-SHA-1 cannot have mechanism properties")
		}
	case "scram-sha-256":
		if u.Username == "" {
			return fmt.Errorf("username required for SCRAM-SHA-256")
		}
		if u.Password == "" {
			return fmt.Errorf("password required for SCRAM-SHA-256")
		}
		if u.AuthMechanismProperties != nil {
			return fmt.Errorf("SCRAM-SHA-256 cannot have mechanism properties")
		}
	case "mongodb-oidc":
		if u.Password != "" {
			return fmt.Errorf("password cannot be specified for MONGODB-OIDC")
		}
	case "":
		if u.UsernameSet && u.Username == "" {
			return fmt.Errorf("username required if URI contains user info")
		}
	default:
		return fmt.Errorf("invalid auth mechanism")
	}
	return nil
}

func (u *ConnString) validateSSL() error {
	if !u.SSL {
		return nil
	}

	if u.SSLClientCertificateKeyFileSet {
		if u.SSLCertificateFileSet || u.SSLPrivateKeyFileSet {
			return errors.New("the sslClientCertificateKeyFile/tlsCertificateKeyFile URI option cannot be provided " +
				"along with tlsCertificateFile or tlsPrivateKeyFile")
		}
		return nil
	}
	if u.SSLCertificateFileSet && !u.SSLPrivateKeyFileSet {
		return errors.New("the tlsPrivateKeyFile URI option must be provided if the tlsCertificateFile option is specified")
	}
	if u.SSLPrivateKeyFileSet && !u.SSLCertificateFileSet {
		return errors.New("the tlsCertificateFile URI option must be provided if the tlsPrivateKeyFile option is specified")
	}

	if u.SSLInsecureSet && u.SSLDisableOCSPEndpointCheckSet {
		return errors.New("the sslInsecure/tlsInsecure URI option cannot be provided along with " +
			"tlsDisableOCSPEndpointCheck ")
	}
	return nil
}

func sanitizeHost(host string) (string, error) {
	if host == "" {
		return host, nil
	}
	unescaped, err := url.QueryUnescape(host)
	if err != nil {
		return "", fmt.Errorf("invalid host %q: %w", host, err)
	}

	_, port, err := net.SplitHostPort(unescaped)
	
	
	if err != nil {
		if addrError, ok := err.(*net.AddrError); !ok || addrError.Err != "missing port in address" {
			return "", err
		}
	}

	if port != "" {
		d, err := strconv.Atoi(port)
		if err != nil {
			return "", fmt.Errorf("port must be an integer: %w", err)
		}
		if d <= 0 || d >= 65536 {
			return "", fmt.Errorf("port must be in the range [1, 65535]")
		}
	}
	return unescaped, nil
}



type ConnectMode uint8

var _ fmt.Stringer = ConnectMode(0)


const (
	AutoConnect ConnectMode = iota
	SingleConnect
)


func (c ConnectMode) String() string {
	switch c {
	case AutoConnect:
		return "automatic"
	case SingleConnect:
		return "direct"
	default:
		return "unknown"
	}
}


const (
	SchemeMongoDB    = "mongodb"
	SchemeMongoDBSRV = "mongodb+srv"
)

type parser struct {
	dnsResolver *dns.Resolver
}

func (p *parser) parse(original string) (*ConnString, error) {
	connStr := &ConnString{}
	connStr.Original = original
	uri := original

	var err error
	switch {
	case strings.HasPrefix(uri, SchemeMongoDBSRV+"://"):
		connStr.Scheme = SchemeMongoDBSRV
		
		uri = uri[len(SchemeMongoDBSRV)+3:]
	case strings.HasPrefix(uri, SchemeMongoDB+"://"):
		connStr.Scheme = SchemeMongoDB
		
		uri = uri[len(SchemeMongoDB)+3:]
	default:
		return nil, errors.New(`scheme must be "mongodb" or "mongodb+srv"`)
	}

	if idx := strings.Index(uri, "@"); idx != -1 {
		userInfo := uri[:idx]
		uri = uri[idx+1:]

		username := userInfo
		var password string

		if u, p, ok := strings.Cut(userInfo, ":"); ok {
			username = u
			password = p
			connStr.PasswordSet = true
		}

		
		if strings.Contains(username, "/") {
			return nil, fmt.Errorf("unescaped slash in username")
		}
		connStr.Username, err = url.PathUnescape(username)
		if err != nil {
			return nil, fmt.Errorf("invalid username: %w", err)
		}
		connStr.UsernameSet = true

		
		if strings.Contains(password, ":") {
			return nil, fmt.Errorf("unescaped colon in password")
		}
		if strings.Contains(password, "/") {
			return nil, fmt.Errorf("unescaped slash in password")
		}
		connStr.Password, err = url.PathUnescape(password)
		if err != nil {
			return nil, fmt.Errorf("invalid password: %w", err)
		}
	}

	
	hosts := uri
	if idx := strings.IndexAny(uri, "/?@"); idx != -1 {
		if uri[idx] == '@' {
			return nil, fmt.Errorf("unescaped @ sign in user info")
		}
		if uri[idx] == '?' {
			return nil, fmt.Errorf("must have a / before the query ?")
		}
		hosts = uri[:idx]
	}

	for _, host := range strings.Split(hosts, ",") {
		host, err = sanitizeHost(host)
		if err != nil {
			return nil, fmt.Errorf("invalid host %q: %w", host, err)
		}
		if host != "" {
			connStr.RawHosts = append(connStr.RawHosts, host)
		}
	}
	connStr.Hosts = connStr.RawHosts
	uri = uri[len(hosts):]
	extractedDatabase, err := extractDatabaseFromURI(uri)
	if err != nil {
		return nil, err
	}

	uri = extractedDatabase.uri
	connStr.Database = extractedDatabase.db

	
	connectionArgsFromQueryString, err := extractQueryArgsFromURI(uri)
	if err != nil {
		return nil, err
	}

	
	var connectionArgsFromTXT []string
	if connStr.Scheme == SchemeMongoDBSRV && p.dnsResolver != nil {
		connectionArgsFromTXT, err = p.dnsResolver.GetConnectionArgsFromTXT(hosts)
		if err != nil {
			return nil, err
		}

		
		connStr.SSL = true
		connStr.SSLSet = true
	}

	
	connectionArgPairs := make([]string, 0, len(connectionArgsFromTXT)+len(connectionArgsFromQueryString))
	connectionArgPairs = append(connectionArgPairs, connectionArgsFromTXT...)
	connectionArgPairs = append(connectionArgPairs, connectionArgsFromQueryString...)

	err = connStr.addOptions(connectionArgPairs)
	if err != nil {
		return nil, err
	}

	
	if connStr.Scheme == SchemeMongoDBSRV && p.dnsResolver != nil {
		parsedHosts, err := p.dnsResolver.ParseHosts(hosts, connStr.SRVServiceName, true)
		if err != nil {
			return connStr, err
		}

		
		
		if connStr.SRVMaxHosts > 0 && connStr.SRVMaxHosts < len(parsedHosts) {
			random.Shuffle(len(parsedHosts), func(i, j int) {
				parsedHosts[i], parsedHosts[j] = parsedHosts[j], parsedHosts[i]
			})
			parsedHosts = parsedHosts[:connStr.SRVMaxHosts]
		}

		var hosts []string
		for _, host := range parsedHosts {
			host, err = sanitizeHost(host)
			if err != nil {
				return connStr, fmt.Errorf("invalid host %q: %w", host, err)
			}
			if host != "" {
				hosts = append(hosts, host)
			}
		}
		connStr.Hosts = hosts
	}
	if len(connStr.Hosts) == 0 {
		return nil, fmt.Errorf("must have at least 1 host")
	}

	err = connStr.setDefaultAuthParams(extractedDatabase.db)
	if err != nil {
		return nil, err
	}

	
	if connStr.WTimeoutSetFromOption {
		connStr.WTimeoutSet = true
	}

	return connStr, nil
}



func IsValidServerMonitoringMode(mode string) bool {
	return mode == ServerMonitoringModeAuto ||
		mode == ServerMonitoringModeStream ||
		mode == ServerMonitoringModePoll
}

func extractQueryArgsFromURI(uri string) ([]string, error) {
	if len(uri) == 0 {
		return nil, nil
	}

	if uri[0] != '?' {
		return nil, errors.New("must have a ? separator between path and query")
	}

	uri = uri[1:]
	if len(uri) == 0 {
		return nil, nil
	}
	return strings.FieldsFunc(uri, func(r rune) bool { return r == ';' || r == '&' }), nil

}

type extractedDatabase struct {
	uri string
	db  string
}





func extractDatabaseFromURI(uri string) (extractedDatabase, error) {
	if len(uri) == 0 {
		return extractedDatabase{}, nil
	}

	if uri[0] != '/' {
		return extractedDatabase{}, errors.New("must have a / separator between hosts and path")
	}

	uri = uri[1:]
	if len(uri) == 0 {
		return extractedDatabase{}, nil
	}

	database := uri
	if idx := strings.IndexRune(uri, '?'); idx != -1 {
		database = uri[:idx]
	}

	escapedDatabase, err := url.QueryUnescape(database)
	if err != nil {
		return extractedDatabase{}, fmt.Errorf("invalid database %q: %w", database, err)
	}

	uri = uri[len(database):]

	return extractedDatabase{
		uri: uri,
		db:  escapedDatabase,
	}, nil
}
