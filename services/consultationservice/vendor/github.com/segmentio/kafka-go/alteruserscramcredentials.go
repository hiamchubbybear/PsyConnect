package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/alteruserscramcredentials"
)



type AlterUserScramCredentialsRequest struct {
	
	Addr net.Addr

	
	Deletions []UserScramCredentialsDeletion

	
	Upsertions []UserScramCredentialsUpsertion
}

type ScramMechanism int8

const (
	ScramMechanismUnknown ScramMechanism = iota 
	ScramMechanismSha256                        
	ScramMechanismSha512                        
)

type UserScramCredentialsDeletion struct {
	Name      string
	Mechanism ScramMechanism
}

type UserScramCredentialsUpsertion struct {
	Name           string
	Mechanism      ScramMechanism
	Iterations     int
	Salt           []byte
	SaltedPassword []byte
}



type AlterUserScramCredentialsResponse struct {
	
	Throttle time.Duration

	
	Results []AlterUserScramCredentialsResponseUser
}

type AlterUserScramCredentialsResponseUser struct {
	User  string
	Error error
}



func (c *Client) AlterUserScramCredentials(ctx context.Context, req *AlterUserScramCredentialsRequest) (*AlterUserScramCredentialsResponse, error) {
	deletions := make([]alteruserscramcredentials.RequestUserScramCredentialsDeletion, len(req.Deletions))
	upsertions := make([]alteruserscramcredentials.RequestUserScramCredentialsUpsertion, len(req.Upsertions))

	for deletionIdx, deletion := range req.Deletions {
		deletions[deletionIdx] = alteruserscramcredentials.RequestUserScramCredentialsDeletion{
			Name:      deletion.Name,
			Mechanism: int8(deletion.Mechanism),
		}
	}

	for upsertionIdx, upsertion := range req.Upsertions {
		upsertions[upsertionIdx] = alteruserscramcredentials.RequestUserScramCredentialsUpsertion{
			Name:           upsertion.Name,
			Mechanism:      int8(upsertion.Mechanism),
			Iterations:     int32(upsertion.Iterations),
			Salt:           upsertion.Salt,
			SaltedPassword: upsertion.SaltedPassword,
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &alteruserscramcredentials.Request{
		Deletions:  deletions,
		Upsertions: upsertions,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).AlterUserScramCredentials: %w", err)
	}

	res := m.(*alteruserscramcredentials.Response)
	responseEntries := make([]AlterUserScramCredentialsResponseUser, len(res.Results))

	for responseIdx, responseResult := range res.Results {
		responseEntries[responseIdx] = AlterUserScramCredentialsResponseUser{
			User:  responseResult.User,
			Error: makeError(responseResult.ErrorCode, responseResult.ErrorMessage),
		}
	}
	ret := &AlterUserScramCredentialsResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
		Results:  responseEntries,
	}

	return ret, nil
}
