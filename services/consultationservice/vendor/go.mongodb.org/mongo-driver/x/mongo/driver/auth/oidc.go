





package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)


const MongoDBOIDC = "MONGODB-OIDC"


const EnvironmentProp = "ENVIRONMENT"


const ResourceProp = "TOKEN_RESOURCE"


const AllowedHostsProp = "ALLOWED_HOSTS"


const AzureEnvironmentValue = "azure"


const GCPEnvironmentValue = "gcp"


const TestEnvironmentValue = "test"

const apiVersion = 1
const invalidateSleepTimeout = 100 * time.Millisecond





const machineCallbackTimeout = time.Minute
const humanCallbackTimeout = 5 * time.Minute

var defaultAllowedHosts = []*regexp.Regexp{
	regexp.MustCompile(`^.*[.]mongodb[.]net(:\d+)?$`),
	regexp.MustCompile(`^.*[.]mongodb-qa[.]net(:\d+)?$`),
	regexp.MustCompile(`^.*[.]mongodb-dev[.]net(:\d+)?$`),
	regexp.MustCompile(`^.*[.]mongodbgov[.]net(:\d+)?$`),
	regexp.MustCompile(`^localhost(:\d+)?$`),
	regexp.MustCompile(`^127[.]0[.]0[.]1(:\d+)?$`),
	regexp.MustCompile(`^::1(:\d+)?$`),
}


type OIDCCallback = driver.OIDCCallback


type OIDCArgs = driver.OIDCArgs


type OIDCCredential = driver.OIDCCredential


type IDPInfo = driver.IDPInfo

var _ driver.Authenticator = (*OIDCAuthenticator)(nil)
var _ SpeculativeAuthenticator = (*OIDCAuthenticator)(nil)
var _ SaslClient = (*oidcOneStep)(nil)
var _ SaslClient = (*oidcTwoStep)(nil)




type OIDCAuthenticator struct {
	mu sync.Mutex 

	AuthMechanismProperties map[string]string
	OIDCMachineCallback     OIDCCallback
	OIDCHumanCallback       OIDCCallback

	allowedHosts *[]*regexp.Regexp
	userName     string
	httpClient   *http.Client
	accessToken  string
	refreshToken *string
	idpInfo      *IDPInfo
	tokenGenID   uint64
}



func (oa *OIDCAuthenticator) SetAccessToken(accessToken string) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	oa.accessToken = accessToken
}

func newOIDCAuthenticator(cred *Cred, httpClient *http.Client) (Authenticator, error) {
	if cred.Source != "" && cred.Source != sourceExternal {
		return nil, newAuthError("MONGODB-OIDC source must be empty or $external", nil)
	}
	if cred.Password != "" {
		return nil, fmt.Errorf("password cannot be specified for %q", MongoDBOIDC)
	}
	if cred.Props != nil {
		if env, ok := cred.Props[EnvironmentProp]; ok {
			switch strings.ToLower(env) {
			case AzureEnvironmentValue:
				fallthrough
			case GCPEnvironmentValue:
				if _, ok := cred.Props[ResourceProp]; !ok {
					return nil, fmt.Errorf("%q must be specified for %q %q", ResourceProp, env, EnvironmentProp)
				}
				fallthrough
			case TestEnvironmentValue:
				if cred.OIDCMachineCallback != nil || cred.OIDCHumanCallback != nil {
					return nil, fmt.Errorf("OIDC callbacks are not allowed for %q %q", env, EnvironmentProp)
				}
			}
		}
	}
	oa := &OIDCAuthenticator{
		userName:                cred.Username,
		httpClient:              httpClient,
		AuthMechanismProperties: cred.Props,
		OIDCMachineCallback:     cred.OIDCMachineCallback,
		OIDCHumanCallback:       cred.OIDCHumanCallback,
	}
	err := oa.setAllowedHosts()
	return oa, err
}

func createPatternsForGlobs(hosts []string) ([]*regexp.Regexp, error) {
	var err error
	ret := make([]*regexp.Regexp, len(hosts))
	for i := range hosts {
		hosts[i] = strings.ReplaceAll(hosts[i], ".", "[.]")
		hosts[i] = strings.ReplaceAll(hosts[i], "*", ".*")
		hosts[i] = "^" + hosts[i] + "(:\\d+)?$"
		ret[i], err = regexp.Compile(hosts[i])
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (oa *OIDCAuthenticator) setAllowedHosts() error {
	if oa.AuthMechanismProperties == nil {
		oa.allowedHosts = &defaultAllowedHosts
		return nil
	}

	allowedHosts, ok := oa.AuthMechanismProperties[AllowedHostsProp]
	if !ok {
		oa.allowedHosts = &defaultAllowedHosts
		return nil
	}
	globs := strings.Split(allowedHosts, ",")
	ret, err := createPatternsForGlobs(globs)
	if err != nil {
		return err
	}
	oa.allowedHosts = &ret
	return nil
}

func (oa *OIDCAuthenticator) validateConnectionAddressWithAllowedHosts(conn driver.Connection) error {
	if oa.allowedHosts == nil {
		
		return newAuthError(fmt.Sprintf("%q missing", AllowedHostsProp), nil)
	}
	allowedHosts := *oa.allowedHosts
	if len(allowedHosts) == 0 {
		return newAuthError(fmt.Sprintf("empty %q specified", AllowedHostsProp), nil)
	}
	for _, pattern := range allowedHosts {
		if pattern.MatchString(string(conn.Address())) {
			return nil
		}
	}
	return newAuthError(fmt.Sprintf("address %q not allowed by %q: %v", conn.Address(), AllowedHostsProp, allowedHosts), nil)
}

type oidcOneStep struct {
	userName    string
	accessToken string
}

type oidcTwoStep struct {
	conn driver.Connection
	oa   *OIDCAuthenticator
}

func jwtStepRequest(accessToken string) []byte {
	return bsoncore.NewDocumentBuilder().
		AppendString("jwt", accessToken).
		Build()
}

func principalStepRequest(principal string) []byte {
	doc := bsoncore.NewDocumentBuilder()
	if principal != "" {
		doc.AppendString("n", principal)
	}
	return doc.Build()
}

func (oos *oidcOneStep) Start() (string, []byte, error) {
	return MongoDBOIDC, jwtStepRequest(oos.accessToken), nil
}

func (oos *oidcOneStep) Next(context.Context, []byte) ([]byte, error) {
	return nil, newAuthError("unexpected step in OIDC authentication", nil)
}

func (*oidcOneStep) Completed() bool {
	return true
}

func (ots *oidcTwoStep) Start() (string, []byte, error) {
	return MongoDBOIDC, principalStepRequest(ots.oa.userName), nil
}

func (ots *oidcTwoStep) Next(ctx context.Context, msg []byte) ([]byte, error) {
	var idpInfo IDPInfo
	err := bson.Unmarshal(msg, &idpInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling BSON document: %w", err)
	}

	accessToken, err := ots.oa.getAccessToken(ctx,
		ots.conn,
		&OIDCArgs{
			Version: apiVersion,
			
			IDPInfo: &idpInfo,
			
			RefreshToken: nil,
		},
		
		ots.oa.OIDCHumanCallback)

	return jwtStepRequest(accessToken), err
}

func (*oidcTwoStep) Completed() bool {
	return true
}

func (oa *OIDCAuthenticator) providerCallback() (OIDCCallback, error) {
	env, ok := oa.AuthMechanismProperties[EnvironmentProp]
	if !ok {
		return nil, nil
	}

	switch env {
	case AzureEnvironmentValue:
		resource, ok := oa.AuthMechanismProperties[ResourceProp]
		if !ok {
			return nil, newAuthError(fmt.Sprintf("%q must be specified for Azure OIDC", ResourceProp), nil)
		}
		return getAzureOIDCCallback(oa.userName, resource, oa.httpClient), nil
	case GCPEnvironmentValue:
		resource, ok := oa.AuthMechanismProperties[ResourceProp]
		if !ok {
			return nil, newAuthError(fmt.Sprintf("%q must be specified for GCP OIDC", ResourceProp), nil)
		}
		return getGCPOIDCCallback(resource, oa.httpClient), nil
	}

	return nil, fmt.Errorf("%q %q not supported for MONGODB-OIDC", EnvironmentProp, env)
}


func getAzureOIDCCallback(clientID string, resource string, httpClient *http.Client) OIDCCallback {
	
	
	return func(ctx context.Context, _ *OIDCArgs) (*OIDCCredential, error) {
		resource = url.QueryEscape(resource)
		var uri string
		if clientID != "" {
			uri = fmt.Sprintf("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=%s&client_id=%s", resource, clientID)
		} else {
			uri = fmt.Sprintf("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=%s", resource)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return nil, newAuthError("error creating http request to Azure Identity Provider", err)
		}
		req.Header.Add("Metadata", "true")
		req.Header.Add("Accept", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, newAuthError("error getting access token from Azure Identity Provider", err)
		}
		defer resp.Body.Close()
		var azureResp struct {
			AccessToken string `json:"access_token"`
			ExpiresOn   int64  `json:"expires_on,string"`
		}

		if resp.StatusCode != http.StatusOK {
			return nil, newAuthError(fmt.Sprintf("failed to get a valid response from Azure Identity Provider, http code: %d", resp.StatusCode), nil)
		}
		err = json.NewDecoder(resp.Body).Decode(&azureResp)
		if err != nil {
			return nil, newAuthError("failed parsing result from Azure Identity Provider", err)
		}
		expireTime := time.Unix(azureResp.ExpiresOn, 0)
		return &OIDCCredential{
			AccessToken: azureResp.AccessToken,
			ExpiresAt:   &expireTime,
		}, nil
	}
}


func getGCPOIDCCallback(resource string, httpClient *http.Client) OIDCCallback {
	
	
	return func(ctx context.Context, _ *OIDCArgs) (*OIDCCredential, error) {
		resource = url.QueryEscape(resource)
		uri := fmt.Sprintf("http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=%s", resource)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return nil, newAuthError("error creating http request to GCP Identity Provider", err)
		}
		req.Header.Add("Metadata-Flavor", "Google")
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, newAuthError("error getting access token from GCP Identity Provider", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, newAuthError(fmt.Sprintf("failed to get a valid response from GCP Identity Provider, http code: %d", resp.StatusCode), nil)
		}
		accessToken, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, newAuthError("failed parsing reading response from GCP Identity Provider", err)
		}
		return &OIDCCredential{
			AccessToken: string(accessToken),
			ExpiresAt:   nil,
		}, nil
	}
}

func (oa *OIDCAuthenticator) getAccessToken(
	ctx context.Context,
	conn driver.Connection,
	args *OIDCArgs,
	callback OIDCCallback,
) (string, error) {
	oa.mu.Lock()
	defer oa.mu.Unlock()

	if oa.accessToken != "" {
		return oa.accessToken, nil
	}

	
	if args.RefreshToken != nil {
		cred, err := callback(ctx, args)
		if err == nil && cred != nil {
			oa.accessToken = cred.AccessToken
			oa.tokenGenID++
			conn.SetOIDCTokenGenID(oa.tokenGenID)
			oa.refreshToken = cred.RefreshToken
			return cred.AccessToken, nil
		}
		oa.refreshToken = nil
		args.RefreshToken = nil
	}
	
	cred, err := callback(ctx, args)
	if err != nil {
		return "", err
	}
	
	
	if cred == nil {
		return "", newAuthError("OIDC callback returned nil credential with no specified error", nil)
	}

	oa.accessToken = cred.AccessToken
	oa.tokenGenID++
	conn.SetOIDCTokenGenID(oa.tokenGenID)
	oa.refreshToken = cred.RefreshToken
	
	
	oa.idpInfo = args.IDPInfo

	return cred.AccessToken, nil
}






func (oa *OIDCAuthenticator) invalidateAccessToken(conn driver.Connection) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	tokenGenID := conn.OIDCTokenGenID()
	
	
	
	if tokenGenID == 0 || tokenGenID >= oa.tokenGenID {
		oa.accessToken = ""
		conn.SetOIDCTokenGenID(0)
	}
}



func (oa *OIDCAuthenticator) Reauth(ctx context.Context, cfg *Config) error {
	oa.invalidateAccessToken(cfg.Connection)
	return oa.Auth(ctx, cfg)
}


func (oa *OIDCAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	var err error

	if cfg == nil {
		return newAuthError(fmt.Sprintf("config must be set for %q authentication", MongoDBOIDC), nil)
	}
	conn := cfg.Connection

	oa.mu.Lock()
	cachedAccessToken := oa.accessToken
	cachedRefreshToken := oa.refreshToken
	cachedIDPInfo := oa.idpInfo
	oa.mu.Unlock()

	if cachedAccessToken != "" {
		err = ConductSaslConversation(ctx, cfg, sourceExternal, &oidcOneStep{
			userName:    oa.userName,
			accessToken: cachedAccessToken,
		})
		if err == nil {
			return nil
		}
		
		
		
		oa.invalidateAccessToken(conn)
		time.Sleep(invalidateSleepTimeout)
	}

	if oa.OIDCHumanCallback != nil {
		return oa.doAuthHuman(ctx, cfg, oa.OIDCHumanCallback, cachedIDPInfo, cachedRefreshToken)
	}

	
	var machineCallback OIDCCallback
	if oa.OIDCMachineCallback != nil {
		machineCallback = oa.OIDCMachineCallback
	} else {
		machineCallback, err = oa.providerCallback()
		if err != nil {
			return fmt.Errorf("error getting built-in OIDC provider: %w", err)
		}
	}

	if machineCallback != nil {
		return oa.doAuthMachine(ctx, cfg, machineCallback)
	}
	return newAuthError("no OIDC callback provided", nil)
}

func (oa *OIDCAuthenticator) doAuthHuman(ctx context.Context, cfg *Config, humanCallback OIDCCallback, idpInfo *IDPInfo, refreshToken *string) error {
	
	err := oa.validateConnectionAddressWithAllowedHosts(cfg.Connection)
	if err != nil {
		return err
	}
	subCtx, cancel := context.WithTimeout(ctx, humanCallbackTimeout)
	defer cancel()
	
	if idpInfo != nil {
		accessToken, err := oa.getAccessToken(subCtx,
			cfg.Connection,
			&OIDCArgs{
				Version: apiVersion,
				
				IDPInfo:      idpInfo,
				RefreshToken: refreshToken,
			},
			humanCallback)
		if err != nil {
			return err
		}
		return ConductSaslConversation(
			subCtx,
			cfg,
			sourceExternal,
			&oidcOneStep{accessToken: accessToken},
		)
	}
	
	ots := &oidcTwoStep{
		conn: cfg.Connection,
		oa:   oa,
	}
	return ConductSaslConversation(subCtx, cfg, sourceExternal, ots)
}

func (oa *OIDCAuthenticator) doAuthMachine(ctx context.Context, cfg *Config, machineCallback OIDCCallback) error {
	subCtx, cancel := context.WithTimeout(ctx, machineCallbackTimeout)
	accessToken, err := oa.getAccessToken(subCtx,
		cfg.Connection,
		&OIDCArgs{
			Version: apiVersion,
			
			IDPInfo:      nil,
			RefreshToken: nil,
		},
		machineCallback)
	cancel()
	if err != nil {
		return err
	}
	return ConductSaslConversation(
		ctx,
		cfg,
		sourceExternal,
		&oidcOneStep{accessToken: accessToken},
	)
}


func (oa *OIDCAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	accessToken := oa.accessToken
	if accessToken == "" {
		return nil, nil 
	}

	return newSaslConversation(&oidcOneStep{accessToken: accessToken}, sourceExternal, true), nil
}
