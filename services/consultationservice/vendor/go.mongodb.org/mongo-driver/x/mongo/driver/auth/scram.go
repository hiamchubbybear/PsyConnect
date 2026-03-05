











package auth

import (
	"context"
	"net/http"

	"github.com/xdg-go/scram"
	"github.com/xdg-go/stringprep"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

const (
	
	SCRAMSHA1 = "SCRAM-SHA-1"

	
	SCRAMSHA256 = "SCRAM-SHA-256"
)

var (
	
	scramStartOptions bsoncore.Document = bsoncore.BuildDocumentFromElements(nil,
		bsoncore.AppendBooleanElement(nil, "skipEmptyExchange", true),
	)
)

func newScramSHA1Authenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	source := cred.Source
	if source == "" {
		source = "admin"
	}
	passdigest := mongoPasswordDigest(cred.Username, cred.Password)
	client, err := scram.SHA1.NewClientUnprepped(cred.Username, passdigest, "")
	if err != nil {
		return nil, newAuthError("error initializing SCRAM-SHA-1 client", err)
	}
	client.WithMinIterations(4096)
	return &ScramAuthenticator{
		mechanism: SCRAMSHA1,
		source:    source,
		client:    client,
	}, nil
}

func newScramSHA256Authenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	source := cred.Source
	if source == "" {
		source = "admin"
	}
	passprep, err := stringprep.SASLprep.Prepare(cred.Password)
	if err != nil {
		return nil, newAuthError("error SASLprepping password", err)
	}
	client, err := scram.SHA256.NewClientUnprepped(cred.Username, passprep, "")
	if err != nil {
		return nil, newAuthError("error initializing SCRAM-SHA-256 client", err)
	}
	client.WithMinIterations(4096)
	return &ScramAuthenticator{
		mechanism: SCRAMSHA256,
		source:    source,
		client:    client,
	}, nil
}


type ScramAuthenticator struct {
	mechanism string
	source    string
	client    *scram.Client
}

var _ SpeculativeAuthenticator = (*ScramAuthenticator)(nil)


func (a *ScramAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	err := ConductSaslConversation(ctx, cfg, a.source, a.createSaslClient())
	if err != nil {
		return newAuthError("sasl conversation error", err)
	}
	return nil
}


func (a *ScramAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("SCRAM does not support reauthentication", nil)
}


func (a *ScramAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	return newSaslConversation(a.createSaslClient(), a.source, true), nil
}

func (a *ScramAuthenticator) createSaslClient() SaslClient {
	return &scramSaslAdapter{
		conversation: a.client.NewConversation(),
		mechanism:    a.mechanism,
	}
}

type scramSaslAdapter struct {
	mechanism    string
	conversation *scram.ClientConversation
}

var _ SaslClient = (*scramSaslAdapter)(nil)
var _ ExtraOptionsSaslClient = (*scramSaslAdapter)(nil)

func (a *scramSaslAdapter) Start() (string, []byte, error) {
	step, err := a.conversation.Step("")
	if err != nil {
		return a.mechanism, nil, err
	}
	return a.mechanism, []byte(step), nil
}

func (a *scramSaslAdapter) Next(_ context.Context, challenge []byte) ([]byte, error) {
	step, err := a.conversation.Step(string(challenge))
	if err != nil {
		return nil, err
	}
	return []byte(step), nil
}

func (a *scramSaslAdapter) Completed() bool {
	return a.conversation.Done()
}

func (*scramSaslAdapter) StartCommandOptions() bsoncore.Document {
	return scramStartOptions
}
