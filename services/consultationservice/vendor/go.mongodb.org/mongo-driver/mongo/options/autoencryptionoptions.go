





package options

import (
	"crypto/tls"
	"net/http"

	"go.mongodb.org/mongo-driver/internal/httputil"
)















type AutoEncryptionOptions struct {
	KeyVaultClientOptions *ClientOptions
	KeyVaultNamespace     string
	KmsProviders          map[string]map[string]interface{}
	SchemaMap             map[string]interface{}
	BypassAutoEncryption  *bool
	ExtraOptions          map[string]interface{}
	TLSConfig             map[string]*tls.Config
	HTTPClient            *http.Client
	EncryptedFieldsMap    map[string]interface{}
	BypassQueryAnalysis   *bool
}


func AutoEncryption() *AutoEncryptionOptions {
	return &AutoEncryptionOptions{
		HTTPClient: httputil.DefaultHTTPClient,
	}
}










func (a *AutoEncryptionOptions) SetKeyVaultClientOptions(opts *ClientOptions) *AutoEncryptionOptions {
	a.KeyVaultClientOptions = opts
	return a
}


func (a *AutoEncryptionOptions) SetKeyVaultNamespace(ns string) *AutoEncryptionOptions {
	a.KeyVaultNamespace = ns
	return a
}


func (a *AutoEncryptionOptions) SetKmsProviders(providers map[string]map[string]interface{}) *AutoEncryptionOptions {
	a.KmsProviders = providers
	return a
}








func (a *AutoEncryptionOptions) SetSchemaMap(schemaMap map[string]interface{}) *AutoEncryptionOptions {
	a.SchemaMap = schemaMap
	return a
}









func (a *AutoEncryptionOptions) SetBypassAutoEncryption(bypass bool) *AutoEncryptionOptions {
	a.BypassAutoEncryption = &bypass
	return a
}






























func (a *AutoEncryptionOptions) SetExtraOptions(extraOpts map[string]interface{}) *AutoEncryptionOptions {
	a.ExtraOptions = extraOpts
	return a
}





func (a *AutoEncryptionOptions) SetTLSConfig(tlsOpts map[string]*tls.Config) *AutoEncryptionOptions {
	tlsConfigs := make(map[string]*tls.Config)
	for provider, config := range tlsOpts {
		
		if config.MinVersion == 0 {
			config.MinVersion = tls.VersionTLS12
		}
		tlsConfigs[provider] = config
	}
	a.TLSConfig = tlsConfigs
	return a
}



func (a *AutoEncryptionOptions) SetEncryptedFieldsMap(ef map[string]interface{}) *AutoEncryptionOptions {
	a.EncryptedFieldsMap = ef
	return a
}



func (a *AutoEncryptionOptions) SetBypassQueryAnalysis(bypass bool) *AutoEncryptionOptions {
	a.BypassQueryAnalysis = &bypass
	return a
}





func MergeAutoEncryptionOptions(opts ...*AutoEncryptionOptions) *AutoEncryptionOptions {
	aeo := AutoEncryption()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.KeyVaultClientOptions != nil {
			aeo.KeyVaultClientOptions = opt.KeyVaultClientOptions
		}
		if opt.KeyVaultNamespace != "" {
			aeo.KeyVaultNamespace = opt.KeyVaultNamespace
		}
		if opt.KmsProviders != nil {
			aeo.KmsProviders = opt.KmsProviders
		}
		if opt.SchemaMap != nil {
			aeo.SchemaMap = opt.SchemaMap
		}
		if opt.BypassAutoEncryption != nil {
			aeo.BypassAutoEncryption = opt.BypassAutoEncryption
		}
		if opt.ExtraOptions != nil {
			aeo.ExtraOptions = opt.ExtraOptions
		}
		if opt.TLSConfig != nil {
			aeo.TLSConfig = opt.TLSConfig
		}
		if opt.EncryptedFieldsMap != nil {
			aeo.EncryptedFieldsMap = opt.EncryptedFieldsMap
		}
		if opt.BypassQueryAnalysis != nil {
			aeo.BypassQueryAnalysis = opt.BypassQueryAnalysis
		}
		if opt.HTTPClient != nil {
			aeo.HTTPClient = opt.HTTPClient
		}
	}

	return aeo
}
