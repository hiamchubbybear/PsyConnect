





package operation

import (
	"context"
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/internal/bsonutil"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/handshake"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/version"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)





const maxClientMetadataSize = 512

const driverName = "mongo-go-driver"


type Hello struct {
	authenticator      driver.Authenticator
	appname            string
	compressors        []string
	saslSupportedMechs string
	d                  driver.Deployment
	clock              *session.ClusterClock
	speculativeAuth    bsoncore.Document
	topologyVersion    *description.TopologyVersion
	maxAwaitTimeMS     *int64
	serverAPI          *driver.ServerAPIOptions
	loadBalanced       bool

	res bsoncore.Document
}

var _ driver.Handshaker = (*Hello)(nil)


func NewHello() *Hello { return &Hello{} }


func (h *Hello) AppName(appname string) *Hello {
	h.appname = appname
	return h
}


func (h *Hello) ClusterClock(clock *session.ClusterClock) *Hello {
	if h == nil {
		h = new(Hello)
	}

	h.clock = clock
	return h
}


func (h *Hello) Compressors(compressors []string) *Hello {
	h.compressors = compressors
	return h
}



func (h *Hello) SASLSupportedMechs(username string) *Hello {
	h.saslSupportedMechs = username
	return h
}


func (h *Hello) Deployment(d driver.Deployment) *Hello {
	h.d = d
	return h
}


func (h *Hello) SpeculativeAuthenticate(doc bsoncore.Document) *Hello {
	h.speculativeAuth = doc
	return h
}


func (h *Hello) TopologyVersion(tv *description.TopologyVersion) *Hello {
	h.topologyVersion = tv
	return h
}


func (h *Hello) MaxAwaitTimeMS(awaitTime int64) *Hello {
	h.maxAwaitTimeMS = &awaitTime
	return h
}


func (h *Hello) ServerAPI(serverAPI *driver.ServerAPIOptions) *Hello {
	h.serverAPI = serverAPI
	return h
}


func (h *Hello) LoadBalanced(lb bool) *Hello {
	h.loadBalanced = lb
	return h
}


func (h *Hello) Result(addr address.Address) description.Server {
	return description.NewServer(addr, bson.Raw(h.res))
}

const dockerEnvPath = "/.dockerenv"

const (
	
	runtimeNameDocker = "docker"

	
	orchestratorNameK8s = "kubernetes"
)







func getFaasEnvName() string {
	envVars := []string{
		driverutil.EnvVarAWSExecutionEnv,
		driverutil.EnvVarAWSLambdaRuntimeAPI,
		driverutil.EnvVarFunctionsWorkerRuntime,
		driverutil.EnvVarKService,
		driverutil.EnvVarFunctionName,
		driverutil.EnvVarVercel,
	}

	
	
	names := make(map[string]struct{})

	for _, envVar := range envVars {
		val := os.Getenv(envVar)
		if val == "" {
			continue
		}

		var name string

		switch envVar {
		case driverutil.EnvVarAWSExecutionEnv:
			if !strings.HasPrefix(val, driverutil.AwsLambdaPrefix) {
				continue
			}

			name = driverutil.EnvNameAWSLambda
		case driverutil.EnvVarAWSLambdaRuntimeAPI:
			name = driverutil.EnvNameAWSLambda
		case driverutil.EnvVarFunctionsWorkerRuntime:
			name = driverutil.EnvNameAzureFunc
		case driverutil.EnvVarKService, driverutil.EnvVarFunctionName:
			name = driverutil.EnvNameGCPFunc
		case driverutil.EnvVarVercel:
			
			delete(names, driverutil.EnvNameAWSLambda)

			name = driverutil.EnvNameVercel
		}

		names[name] = struct{}{}
		if len(names) > 1 {
			
			
			names = nil

			break
		}
	}

	for name := range names {
		return name
	}

	return ""
}

type containerInfo struct {
	runtime      string
	orchestrator string
}




func getContainerEnvInfo() *containerInfo {
	var runtime, orchestrator string
	if _, err := os.Stat(dockerEnvPath); !os.IsNotExist(err) {
		runtime = runtimeNameDocker
	}
	if v := os.Getenv(driverutil.EnvVarK8s); v != "" {
		orchestrator = orchestratorNameK8s
	}
	if runtime != "" || orchestrator != "" {
		return &containerInfo{
			runtime:      runtime,
			orchestrator: orchestrator,
		}
	}
	return nil
}




func appendClientAppName(dst []byte, name string) ([]byte, error) {
	if name == "" {
		return dst, nil
	}

	var idx int32
	idx, dst = bsoncore.AppendDocumentElementStart(dst, "application")

	dst = bsoncore.AppendStringElement(dst, "name", name)

	return bsoncore.AppendDocumentEnd(dst, idx)
}




func appendClientDriver(dst []byte) ([]byte, error) {
	var idx int32
	idx, dst = bsoncore.AppendDocumentElementStart(dst, "driver")

	dst = bsoncore.AppendStringElement(dst, "name", driverName)
	dst = bsoncore.AppendStringElement(dst, "version", version.Driver)

	return bsoncore.AppendDocumentEnd(dst, idx)
}




func appendClientEnv(dst []byte, omitNonName, omitDoc bool) ([]byte, error) {
	if omitDoc {
		return dst, nil
	}

	name := getFaasEnvName()
	container := getContainerEnvInfo()
	
	
	if name == "" && container == nil {
		return dst, nil
	}

	var idx int32

	idx, dst = bsoncore.AppendDocumentElementStart(dst, "env")

	if name != "" {
		dst = bsoncore.AppendStringElement(dst, "name", name)
	}

	addMem := func(envVar string) []byte {
		mem := os.Getenv(envVar)
		if mem == "" {
			return dst
		}

		memInt64, err := strconv.ParseInt(mem, 10, 32)
		if err != nil {
			return dst
		}

		memInt32 := int32(memInt64)

		return bsoncore.AppendInt32Element(dst, "memory_mb", memInt32)
	}

	addRegion := func(envVar string) []byte {
		region := os.Getenv(envVar)
		if region == "" {
			return dst
		}

		return bsoncore.AppendStringElement(dst, "region", region)
	}

	addTimeout := func(envVar string) []byte {
		timeout := os.Getenv(envVar)
		if timeout == "" {
			return dst
		}

		timeoutInt64, err := strconv.ParseInt(timeout, 10, 32)
		if err != nil {
			return dst
		}

		timeoutInt32 := int32(timeoutInt64)
		return bsoncore.AppendInt32Element(dst, "timeout_sec", timeoutInt32)
	}

	if !omitNonName {
		
		switch name {
		case driverutil.EnvNameAWSLambda:
			dst = addMem(driverutil.EnvVarAWSLambdaFunctionMemorySize)
			dst = addRegion(driverutil.EnvVarAWSRegion)
		case driverutil.EnvNameGCPFunc:
			dst = addMem(driverutil.EnvVarFunctionMemoryMB)
			dst = addRegion(driverutil.EnvVarFunctionRegion)
			dst = addTimeout(driverutil.EnvVarFunctionTimeoutSec)
		case driverutil.EnvNameVercel:
			dst = addRegion(driverutil.EnvVarVercelRegion)
		}
	}

	if container != nil {
		var idxCntnr int32
		idxCntnr, dst = bsoncore.AppendDocumentElementStart(dst, "container")
		if container.runtime != "" {
			dst = bsoncore.AppendStringElement(dst, "runtime", container.runtime)
		}
		if container.orchestrator != "" {
			dst = bsoncore.AppendStringElement(dst, "orchestrator", container.orchestrator)
		}
		var err error
		dst, err = bsoncore.AppendDocumentEnd(dst, idxCntnr)
		if err != nil {
			return dst, err
		}
	}

	return bsoncore.AppendDocumentEnd(dst, idx)
}




func appendClientOS(dst []byte, omitNonType bool) ([]byte, error) {
	var idx int32

	idx, dst = bsoncore.AppendDocumentElementStart(dst, "os")

	dst = bsoncore.AppendStringElement(dst, "type", runtime.GOOS)
	if !omitNonType {
		dst = bsoncore.AppendStringElement(dst, "architecture", runtime.GOARCH)
	}

	return bsoncore.AppendDocumentEnd(dst, idx)
}




func appendClientPlatform(dst []byte) []byte {
	return bsoncore.AppendStringElement(dst, "platform", runtime.Version())
}



































func encodeClientMetadata(appname string, maxLen int) ([]byte, error) {
	dst := make([]byte, 0, maxLen)

	omitEnvDoc := false
	omitEnvNonName := false
	omitOSNonType := false
	omitEnvDocument := false
	truncatePlatform := false

retry:
	var idx int32
	idx, dst = bsoncore.AppendDocumentStart(dst)

	var err error
	dst, err = appendClientAppName(dst, appname)
	if err != nil {
		return nil, err
	}

	dst, err = appendClientDriver(dst)
	if err != nil {
		return nil, err
	}

	dst, err = appendClientOS(dst, omitOSNonType)
	if err != nil {
		return nil, err
	}

	if !truncatePlatform {
		dst = appendClientPlatform(dst)
	}

	if !omitEnvDocument {
		dst, err = appendClientEnv(dst, omitEnvNonName, omitEnvDoc)
		if err != nil {
			return nil, err
		}
	}

	dst, err = bsoncore.AppendDocumentEnd(dst, idx)
	if err != nil {
		return nil, err
	}

	if len(dst) > maxLen {
		
		
		
		
		
		
		
		dst = dst[:0]

		if !omitEnvNonName {
			omitEnvNonName = true

			goto retry
		}

		if !omitOSNonType {
			omitOSNonType = true

			goto retry
		}

		if !omitEnvDoc {
			omitEnvDoc = true

			goto retry
		}

		if !truncatePlatform {
			truncatePlatform = true

			goto retry
		}

		
		
		return nil, nil
	}

	return dst, nil
}


func (h *Hello) handshakeCommand(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst, err := h.command(dst, desc)
	if err != nil {
		return dst, err
	}

	if h.saslSupportedMechs != "" {
		dst = bsoncore.AppendStringElement(dst, "saslSupportedMechs", h.saslSupportedMechs)
	}
	if h.speculativeAuth != nil {
		dst = bsoncore.AppendDocumentElement(dst, "speculativeAuthenticate", h.speculativeAuth)
	}
	var idx int32
	idx, dst = bsoncore.AppendArrayElementStart(dst, "compression")
	for i, compressor := range h.compressors {
		dst = bsoncore.AppendStringElement(dst, strconv.Itoa(i), compressor)
	}
	dst, _ = bsoncore.AppendArrayEnd(dst, idx)

	clientMetadata, _ := encodeClientMetadata(h.appname, maxClientMetadataSize)

	
	if len(clientMetadata) > 0 {
		dst = bsoncore.AppendDocumentElement(dst, "client", clientMetadata)
	}

	return dst, nil
}


func (h *Hello) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	
	
	if h.loadBalanced || h.serverAPI != nil || desc.Server.HelloOK {
		dst = bsoncore.AppendInt32Element(dst, "hello", 1)
	} else {
		dst = bsoncore.AppendInt32Element(dst, handshake.LegacyHello, 1)
	}
	dst = bsoncore.AppendBooleanElement(dst, "helloOk", true)

	if tv := h.topologyVersion; tv != nil {
		var tvIdx int32

		tvIdx, dst = bsoncore.AppendDocumentElementStart(dst, "topologyVersion")
		dst = bsoncore.AppendObjectIDElement(dst, "processId", tv.ProcessID)
		dst = bsoncore.AppendInt64Element(dst, "counter", tv.Counter)
		dst, _ = bsoncore.AppendDocumentEnd(dst, tvIdx)
	}
	if h.maxAwaitTimeMS != nil {
		dst = bsoncore.AppendInt64Element(dst, "maxAwaitTimeMS", *h.maxAwaitTimeMS)
	}
	if h.loadBalanced {
		
		
		dst = bsoncore.AppendBooleanElement(dst, "loadBalanced", true)
	}

	return dst, nil
}


func (h *Hello) Execute(ctx context.Context) error {
	if h.d == nil {
		return errors.New("a Hello must have a Deployment set before Execute can be called")
	}

	return h.createOperation().Execute(ctx)
}


func (h *Hello) StreamResponse(ctx context.Context, conn driver.StreamerConnection) error {
	return h.createOperation().ExecuteExhaust(ctx, conn)
}





func isLegacyHandshake(srvAPI *driver.ServerAPIOptions, loadbalanced bool) bool {
	return srvAPI == nil && !loadbalanced
}

func (h *Hello) createOperation() driver.Operation {
	op := driver.Operation{
		Clock:      h.clock,
		CommandFn:  h.command,
		Database:   "admin",
		Deployment: h.d,
		ProcessResponseFn: func(info driver.ResponseInfo) error {
			h.res = info.ServerResponse
			return nil
		},
		ServerAPI: h.serverAPI,
	}

	if isLegacyHandshake(h.serverAPI, h.loadBalanced) {
		op.Legacy = driver.LegacyHandshake
	}

	return op
}



func (h *Hello) GetHandshakeInformation(ctx context.Context, _ address.Address, c driver.Connection) (driver.HandshakeInformation, error) {
	deployment := driver.SingleConnectionDeployment{C: c}

	op := driver.Operation{
		Clock:      h.clock,
		CommandFn:  h.handshakeCommand,
		Deployment: deployment,
		Database:   "admin",
		ProcessResponseFn: func(info driver.ResponseInfo) error {
			h.res = info.ServerResponse
			return nil
		},
		ServerAPI: h.serverAPI,
	}

	if isLegacyHandshake(h.serverAPI, h.loadBalanced) {
		op.Legacy = driver.LegacyHandshake
	}

	if err := op.Execute(ctx); err != nil {
		return driver.HandshakeInformation{}, err
	}

	info := driver.HandshakeInformation{
		Description: h.Result(c.Address()),
	}
	if speculativeAuthenticate, ok := h.res.Lookup("speculativeAuthenticate").DocumentOK(); ok {
		info.SpeculativeAuthenticate = speculativeAuthenticate
	}
	if serverConnectionID, ok := h.res.Lookup("connectionId").AsInt64OK(); ok {
		info.ServerConnectionID = &serverConnectionID
	}

	var err error

	
	
	if saslSupportedMechs, lookupErr := bson.Raw(h.res).LookupErr("saslSupportedMechs"); lookupErr == nil {
		info.SaslSupportedMechs, err = bsonutil.StringSliceFromRawValue("saslSupportedMechs", saslSupportedMechs)
	}
	return info, err
}



func (h *Hello) FinishHandshake(context.Context, driver.Connection) error {
	return nil
}


func (h *Hello) Authenticator(authenticator driver.Authenticator) *Hello {
	if h == nil {
		h = new(Hello)
	}

	h.authenticator = authenticator
	return h
}
