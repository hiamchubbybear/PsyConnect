












package driver 

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)




type AuthConfig struct {
	Description   description.Server
	Connection    Connection
	ClusterClock  *session.ClusterClock
	HandshakeInfo HandshakeInformation
	ServerAPI     *ServerAPIOptions
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
	Issuer        string   `bson:"issuer"`
	ClientID      string   `bson:"clientId"`
	RequestScopes []string `bson:"requestScopes"`
}





type Authenticator interface {
	
	Auth(context.Context, *AuthConfig) error
	Reauth(context.Context, *AuthConfig) error
}


type Cred struct {
	Source              string
	Username            string
	Password            string
	PasswordSet         bool
	Props               map[string]string
	OIDCMachineCallback OIDCCallback
	OIDCHumanCallback   OIDCCallback
}


type Deployment interface {
	SelectServer(context.Context, description.ServerSelector) (Server, error)
	Kind() description.TopologyKind
}


type Connector interface {
	Connect() error
}


type Disconnector interface {
	Disconnect(context.Context) error
}



type Subscription struct {
	Updates <-chan description.Topology
	ID      uint64
}



type Subscriber interface {
	Subscribe() (*Subscription, error)
	Unsubscribe(*Subscription) error
}



type Server interface {
	Connection(context.Context) (Connection, error)

	
	RTTMonitor() RTTMonitor
}


type Connection interface {
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
}


type RTTMonitor interface {
	
	EWMA() time.Duration

	
	Min() time.Duration

	
	P90() time.Duration

	
	Stats() string
}

var _ RTTMonitor = &csot.ZeroRTTMonitor{}







type PinnedConnection interface {
	Connection
	PinToCursor() error
	PinToTransaction() error
	UnpinFromCursor() error
	UnpinFromTransaction() error
}




var _ PinnedConnection = (session.LoadBalancedTransactionConnection)(nil)


type LocalAddresser interface {
	LocalAddress() address.Address
}


type Expirable interface {
	Expire() error
	Alive() bool
}









type StreamerConnection interface {
	Connection
	SetStreaming(bool)
	CurrentlyStreaming() bool
	SupportsStreaming() bool
}




type Compressor interface {
	CompressWireMessage(src, dst []byte) ([]byte, error)
}




type ProcessErrorResult int

const (
	
	NoChange ProcessErrorResult = iota
	
	ServerMarkedUnknown
	
	
	ConnectionPoolCleared
)




type ErrorProcessor interface {
	ProcessError(err error, conn Connection) ProcessErrorResult
}






type HandshakeInformation struct {
	Description             description.Server
	SpeculativeAuthenticate bsoncore.Document
	ServerConnectionID      *int64
	SaslSupportedMechs      []string
}




type Handshaker interface {
	GetHandshakeInformation(context.Context, address.Address, Connection) (HandshakeInformation, error)
	FinishHandshake(context.Context, Connection) error
}


type SingleServerDeployment struct{ Server }

var _ Deployment = SingleServerDeployment{}



func (ssd SingleServerDeployment) SelectServer(context.Context, description.ServerSelector) (Server, error) {
	return ssd.Server, nil
}


func (SingleServerDeployment) Kind() description.TopologyKind { return description.Single }




type SingleConnectionDeployment struct{ C Connection }

var _ Deployment = SingleConnectionDeployment{}
var _ Server = SingleConnectionDeployment{}




func (scd SingleConnectionDeployment) SelectServer(context.Context, description.ServerSelector) (Server, error) {
	return scd, nil
}


func (SingleConnectionDeployment) Kind() description.TopologyKind { return description.Single }


func (scd SingleConnectionDeployment) Connection(context.Context) (Connection, error) {
	return scd.C, nil
}


func (scd SingleConnectionDeployment) RTTMonitor() RTTMonitor {
	return &csot.ZeroRTTMonitor{}
}








type Type uint


const (
	_ Type = iota
	Write
	Read
)


type RetryMode uint




const (
	
	RetryNone RetryMode = iota
	
	RetryOnce
	
	
	
	RetryOncePerCommand
	
	
	RetryContext
)


func (rm RetryMode) Enabled() bool {
	return rm == RetryOnce || rm == RetryOncePerCommand || rm == RetryContext
}
