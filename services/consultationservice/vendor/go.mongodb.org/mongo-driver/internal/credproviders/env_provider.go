





package credproviders

import (
	"os"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)


const envProviderName = "EnvProvider"


type EnvVar string


func (ev EnvVar) Get() string {
	return os.Getenv(string(ev))
}



type EnvProvider struct {
	AwsAccessKeyIDEnv     EnvVar
	AwsSecretAccessKeyEnv EnvVar
	AwsSessionTokenEnv    EnvVar

	retrieved bool
}


func NewEnvProvider() *EnvProvider {
	return &EnvProvider{
		
		AwsAccessKeyIDEnv: EnvVar("AWS_ACCESS_KEY_ID"),
		
		AwsSecretAccessKeyEnv: EnvVar("AWS_SECRET_ACCESS_KEY"),
		
		AwsSessionTokenEnv: EnvVar("AWS_SESSION_TOKEN"),
	}
}


func (e *EnvProvider) Retrieve() (credentials.Value, error) {
	e.retrieved = false

	v := credentials.Value{
		AccessKeyID:     e.AwsAccessKeyIDEnv.Get(),
		SecretAccessKey: e.AwsSecretAccessKeyEnv.Get(),
		SessionToken:    e.AwsSessionTokenEnv.Get(),
		ProviderName:    envProviderName,
	}
	err := verify(v)
	if err == nil {
		e.retrieved = true
	}

	return v, err
}


func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}
