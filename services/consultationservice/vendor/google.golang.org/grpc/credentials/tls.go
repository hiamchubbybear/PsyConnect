

package credentials

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"os"

	"google.golang.org/grpc/grpclog"
	credinternal "google.golang.org/grpc/internal/credentials"
	"google.golang.org/grpc/internal/envconfig"
)

const alpnFailureHelpMessage = "If you upgraded from a grpc-go version earlier than 1.67, your TLS connections may have stopped working due to ALPN enforcement. For more details, see: https://github.com/grpc/grpc-go/issues/434"

var logger = grpclog.Component("credentials")



type TLSInfo struct {
	State tls.ConnectionState
	CommonAuthInfo
	
	SPIFFEID *url.URL
}


func (t TLSInfo) AuthType() string {
	return "tls"
}


func cipherSuiteLookup(cipherSuiteID uint16) string {
	for _, s := range tls.CipherSuites() {
		if s.ID == cipherSuiteID {
			return s.Name
		}
	}
	for _, s := range tls.InsecureCipherSuites() {
		if s.ID == cipherSuiteID {
			return s.Name
		}
	}
	return fmt.Sprintf("unknown ID: %v", cipherSuiteID)
}


func (t TLSInfo) GetSecurityValue() ChannelzSecurityValue {
	v := &TLSChannelzSecurityValue{
		StandardName: cipherSuiteLookup(t.State.CipherSuite),
	}
	
	if len(t.State.PeerCertificates) > 0 {
		v.RemoteCertificate = t.State.PeerCertificates[0].Raw
	}
	return v
}


type tlsCreds struct {
	
	config *tls.Config
}

func (c tlsCreds) Info() ProtocolInfo {
	return ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
		ServerName:       c.config.ServerName,
	}
}

func (c *tlsCreds) ClientHandshake(ctx context.Context, authority string, rawConn net.Conn) (_ net.Conn, _ AuthInfo, err error) {
	
	cfg := credinternal.CloneTLSConfig(c.config)
	if cfg.ServerName == "" {
		serverName, _, err := net.SplitHostPort(authority)
		if err != nil {
			
			serverName = authority
		}
		cfg.ServerName = serverName
	}
	conn := tls.Client(rawConn, cfg)
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- conn.Handshake()
		close(errChannel)
	}()
	select {
	case err := <-errChannel:
		if err != nil {
			conn.Close()
			return nil, nil, err
		}
	case <-ctx.Done():
		conn.Close()
		return nil, nil, ctx.Err()
	}

	
	
	
	
	
	
	
	np := conn.ConnectionState().NegotiatedProtocol
	if np == "" {
		if envconfig.EnforceALPNEnabled {
			conn.Close()
			return nil, nil, fmt.Errorf("credentials: cannot check peer: missing selected ALPN property. %s", alpnFailureHelpMessage)
		}
		logger.Warningf("Allowing TLS connection to server %q with ALPN disabled. TLS connections to servers with ALPN disabled will be disallowed in future grpc-go releases", cfg.ServerName)
	}
	tlsInfo := TLSInfo{
		State: conn.ConnectionState(),
		CommonAuthInfo: CommonAuthInfo{
			SecurityLevel: PrivacyAndIntegrity,
		},
	}
	id := credinternal.SPIFFEIDFromState(conn.ConnectionState())
	if id != nil {
		tlsInfo.SPIFFEID = id
	}
	return credinternal.WrapSyscallConn(rawConn, conn), tlsInfo, nil
}

func (c *tlsCreds) ServerHandshake(rawConn net.Conn) (net.Conn, AuthInfo, error) {
	conn := tls.Server(rawConn, c.config)
	if err := conn.Handshake(); err != nil {
		conn.Close()
		return nil, nil, err
	}
	cs := conn.ConnectionState()
	
	
	
	if cs.NegotiatedProtocol == "" {
		if envconfig.EnforceALPNEnabled {
			conn.Close()
			return nil, nil, fmt.Errorf("credentials: cannot check peer: missing selected ALPN property. %s", alpnFailureHelpMessage)
		} else if logger.V(2) {
			logger.Info("Allowing TLS connection from client with ALPN disabled. TLS connections with ALPN disabled will be disallowed in future grpc-go releases")
		}
	}
	tlsInfo := TLSInfo{
		State: cs,
		CommonAuthInfo: CommonAuthInfo{
			SecurityLevel: PrivacyAndIntegrity,
		},
	}
	id := credinternal.SPIFFEIDFromState(conn.ConnectionState())
	if id != nil {
		tlsInfo.SPIFFEID = id
	}
	return credinternal.WrapSyscallConn(rawConn, conn), tlsInfo, nil
}

func (c *tlsCreds) Clone() TransportCredentials {
	return NewTLS(c.config)
}

func (c *tlsCreds) OverrideServerName(serverNameOverride string) error {
	c.config.ServerName = serverNameOverride
	return nil
}



var tls12ForbiddenCipherSuites = map[uint16]struct{}{
	tls.TLS_RSA_WITH_AES_128_CBC_SHA:         {},
	tls.TLS_RSA_WITH_AES_256_CBC_SHA:         {},
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256:      {},
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384:      {},
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA: {},
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA: {},
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:   {},
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:   {},
}


func NewTLS(c *tls.Config) TransportCredentials {
	config := applyDefaults(c)
	if config.GetConfigForClient != nil {
		oldFn := config.GetConfigForClient
		config.GetConfigForClient = func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
			cfgForClient, err := oldFn(hello)
			if err != nil || cfgForClient == nil {
				return cfgForClient, err
			}
			return applyDefaults(cfgForClient), nil
		}
	}
	return &tlsCreds{config: config}
}

func applyDefaults(c *tls.Config) *tls.Config {
	config := credinternal.CloneTLSConfig(c)
	config.NextProtos = credinternal.AppendH2ToNextProtos(config.NextProtos)
	
	
	
	if config.MinVersion == 0 && (config.MaxVersion == 0 || config.MaxVersion >= tls.VersionTLS12) {
		config.MinVersion = tls.VersionTLS12
	}
	
	
	
	if config.CipherSuites == nil {
		for _, cs := range tls.CipherSuites() {
			if _, ok := tls12ForbiddenCipherSuites[cs.ID]; !ok {
				config.CipherSuites = append(config.CipherSuites, cs.ID)
			}
		}
	}
	return config
}









func NewClientTLSFromCert(cp *x509.CertPool, serverNameOverride string) TransportCredentials {
	return NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp})
}









func NewClientTLSFromFile(certFile, serverNameOverride string) (TransportCredentials, error) {
	b, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}


func NewServerTLSFromCert(cert *tls.Certificate) TransportCredentials {
	return NewTLS(&tls.Config{Certificates: []tls.Certificate{*cert}})
}



func NewServerTLSFromFile(certFile, keyFile string) (TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}}), nil
}








type TLSChannelzSecurityValue struct {
	ChannelzSecurityValue
	StandardName      string
	LocalCertificate  []byte
	RemoteCertificate []byte
}
