





package options

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/internal/httputil"
)


type ClientEncryptionOptions struct {
	KeyVaultNamespace string
	KmsProviders      map[string]map[string]interface{}
	TLSConfig         map[string]*tls.Config
	HTTPClient        *http.Client
}


func ClientEncryption() *ClientEncryptionOptions {
	return &ClientEncryptionOptions{
		HTTPClient: httputil.DefaultHTTPClient,
	}
}


func (c *ClientEncryptionOptions) SetKeyVaultNamespace(ns string) *ClientEncryptionOptions {
	c.KeyVaultNamespace = ns
	return c
}


func (c *ClientEncryptionOptions) SetKmsProviders(providers map[string]map[string]interface{}) *ClientEncryptionOptions {
	c.KmsProviders = providers
	return c
}





func (c *ClientEncryptionOptions) SetTLSConfig(tlsOpts map[string]*tls.Config) *ClientEncryptionOptions {
	tlsConfigs := make(map[string]*tls.Config)
	for provider, config := range tlsOpts {
		
		if config.MinVersion == 0 {
			config.MinVersion = tls.VersionTLS12
		}
		tlsConfigs[provider] = config
	}
	c.TLSConfig = tlsConfigs
	return c
}
























func BuildTLSConfig(tlsOpts map[string]interface{}) (*tls.Config, error) {
	
	cfg := &tls.Config{MinVersion: tls.VersionTLS12}

	for name := range tlsOpts {
		var err error
		switch name {
		case "tlsCertificateKeyFile", "sslClientCertificateKeyFile":
			clientCertPath, ok := tlsOpts[name].(string)
			if !ok {
				return nil, fmt.Errorf("expected %q value to be of type string, got %T", name, tlsOpts[name])
			}
			
			if keyPwd, found := tlsOpts["tlsCertificateKeyFilePassword"].(string); found {
				_, err = addClientCertFromConcatenatedFile(cfg, clientCertPath, keyPwd)
			} else if keyPwd, found := tlsOpts["sslClientCertificateKeyPassword"].(string); found {
				_, err = addClientCertFromConcatenatedFile(cfg, clientCertPath, keyPwd)
			} else {
				_, err = addClientCertFromConcatenatedFile(cfg, clientCertPath, "")
			}
		case "tlsCertificateKeyFilePassword", "sslClientCertificateKeyPassword":
			continue
		case "tlsCAFile", "sslCertificateAuthorityFile":
			caPath, ok := tlsOpts[name].(string)
			if !ok {
				return nil, fmt.Errorf("expected %q value to be of type string, got %T", name, tlsOpts[name])
			}
			err = addCACertFromFile(cfg, caPath)
		default:
			return nil, fmt.Errorf("unrecognized TLS option %v", name)
		}

		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}





func MergeClientEncryptionOptions(opts ...*ClientEncryptionOptions) *ClientEncryptionOptions {
	ceo := ClientEncryption()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.KeyVaultNamespace != "" {
			ceo.KeyVaultNamespace = opt.KeyVaultNamespace
		}
		if opt.KmsProviders != nil {
			ceo.KmsProviders = opt.KmsProviders
		}
		if opt.TLSConfig != nil {
			ceo.TLSConfig = opt.TLSConfig
		}
		if opt.HTTPClient != nil {
			ceo.HTTPClient = opt.HTTPClient
		}
	}

	return ceo
}
