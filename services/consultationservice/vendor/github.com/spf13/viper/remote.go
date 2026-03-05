package viper

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"slices"
)


var SupportedRemoteProviders = []string{"etcd", "etcd3", "consul", "firestore", "nats"}

func resetRemote() {
	SupportedRemoteProviders = []string{"etcd", "etcd3", "consul", "firestore", "nats"}
}

type remoteConfigFactory interface {
	Get(rp RemoteProvider) (io.Reader, error)
	Watch(rp RemoteProvider) (io.Reader, error)
	WatchChannel(rp RemoteProvider) (<-chan *RemoteResponse, chan bool)
}

type RemoteResponse struct {
	Value []byte
	Error error
}


var RemoteConfig remoteConfigFactory



type UnsupportedRemoteProviderError string


func (str UnsupportedRemoteProviderError) Error() string {
	return fmt.Sprintf("Unsupported Remote Provider Type %q", string(str))
}



type RemoteConfigError string


func (rce RemoteConfigError) Error() string {
	return fmt.Sprintf("Remote Configurations Error: %s", string(rce))
}

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}

func (rp defaultRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}





type RemoteProvider interface {
	Provider() string
	Endpoint() string
	Path() string
	SecretKeyring() string
}









func AddRemoteProvider(provider, endpoint, path string) error {
	return v.AddRemoteProvider(provider, endpoint, path)
}

func (v *Viper) AddRemoteProvider(provider, endpoint, path string) error {
	if !slices.Contains(SupportedRemoteProviders, provider) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		v.logger.Info("adding remote provider", "provider", provider, "endpoint", endpoint)

		rp := &defaultRemoteProvider{
			endpoint: endpoint,
			provider: provider,
			path:     path,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}











func AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	return v.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
}

func (v *Viper) AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	if !slices.Contains(SupportedRemoteProviders, provider) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		v.logger.Info("adding remote provider", "provider", provider, "endpoint", endpoint)

		rp := &defaultRemoteProvider{
			endpoint:      endpoint,
			provider:      provider,
			path:          path,
			secretKeyring: secretkeyring,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}

func (v *Viper) providerPathExists(p *defaultRemoteProvider) bool {
	for _, y := range v.remoteProviders {
		if reflect.DeepEqual(y, p) {
			return true
		}
	}
	return false
}



func ReadRemoteConfig() error { return v.ReadRemoteConfig() }

func (v *Viper) ReadRemoteConfig() error {
	return v.getKeyValueConfig()
}

func WatchRemoteConfig() error { return v.WatchRemoteConfig() }
func (v *Viper) WatchRemoteConfig() error {
	return v.watchKeyValueConfig()
}

func (v *Viper) WatchRemoteConfigOnChannel() error {
	return v.watchKeyValueConfigOnChannel()
}


func (v *Viper) getKeyValueConfig() error {
	if RemoteConfig == nil {
		return RemoteConfigError("Enable the remote features by doing a blank import of the viper/remote package: '_ github.com/spf13/viper/remote'")
	}

	if len(v.remoteProviders) == 0 {
		return RemoteConfigError("No Remote Providers")
	}

	for _, rp := range v.remoteProviders {
		val, err := v.getRemoteConfig(rp)
		if err != nil {
			v.logger.Error(fmt.Errorf("get remote config: %w", err).Error())

			continue
		}

		v.kvstore = val

		return nil
	}
	return RemoteConfigError("No Files Found")
}

func (v *Viper) getRemoteConfig(provider RemoteProvider) (map[string]any, error) {
	reader, err := RemoteConfig.Get(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}


func (v *Viper) watchKeyValueConfigOnChannel() error {
	if len(v.remoteProviders) == 0 {
		return RemoteConfigError("No Remote Providers")
	}

	for _, rp := range v.remoteProviders {
		respc, _ := RemoteConfig.WatchChannel(rp)
		
		go func(rc <-chan *RemoteResponse) {
			for {
				b := <-rc
				reader := bytes.NewReader(b.Value)
				v.unmarshalReader(reader, v.kvstore)
			}
		}(respc)
		return nil
	}
	return RemoteConfigError("No Files Found")
}


func (v *Viper) watchKeyValueConfig() error {
	if len(v.remoteProviders) == 0 {
		return RemoteConfigError("No Remote Providers")
	}

	for _, rp := range v.remoteProviders {
		val, err := v.watchRemoteConfig(rp)
		if err != nil {
			v.logger.Error(fmt.Errorf("watch remote config: %w", err).Error())

			continue
		}
		v.kvstore = val
		return nil
	}
	return RemoteConfigError("No Files Found")
}

func (v *Viper) watchRemoteConfig(provider RemoteProvider) (map[string]any, error) {
	reader, err := RemoteConfig.Watch(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}
