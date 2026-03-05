





package credproviders

import (
	"errors"

	"go.mongodb.org/mongo-driver/internal/aws/credentials"
)


const staticProviderName = "StaticProvider"



type StaticProvider struct {
	credentials.Value

	verified bool
	err      error
}

func verify(v credentials.Value) error {
	if !v.HasKeys() {
		return errors.New("failed to retrieve ACCESS_KEY_ID and SECRET_ACCESS_KEY")
	}
	if v.AccessKeyID != "" && v.SecretAccessKey == "" {
		return errors.New("ACCESS_KEY_ID is set, but SECRET_ACCESS_KEY is missing")
	}
	if v.AccessKeyID == "" && v.SecretAccessKey != "" {
		return errors.New("SECRET_ACCESS_KEY is set, but ACCESS_KEY_ID is missing")
	}
	if v.AccessKeyID == "" && v.SecretAccessKey == "" && v.SessionToken != "" {
		return errors.New("AWS_SESSION_TOKEN is set, but ACCESS_KEY_ID and SECRET_ACCESS_KEY are missing")
	}
	return nil

}


func (s *StaticProvider) Retrieve() (credentials.Value, error) {
	if !s.verified {
		s.err = verify(s.Value)
		s.Value.ProviderName = staticProviderName
		s.verified = true
	}
	return s.Value, s.err
}




func (s *StaticProvider) IsExpired() bool {
	return false
}
