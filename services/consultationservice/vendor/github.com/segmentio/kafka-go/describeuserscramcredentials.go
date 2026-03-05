package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/describeuserscramcredentials"
)



type DescribeUserScramCredentialsRequest struct {
	
	Addr net.Addr

	
	Users []UserScramCredentialsUser
}

type UserScramCredentialsUser struct {
	Name string
}



type DescribeUserScramCredentialsResponse struct {
	
	Throttle time.Duration

	
	
	
	
	
	Error error

	
	Results []DescribeUserScramCredentialsResponseResult
}

type DescribeUserScramCredentialsResponseResult struct {
	User            string
	CredentialInfos []DescribeUserScramCredentialsCredentialInfo
	Error           error
}

type DescribeUserScramCredentialsCredentialInfo struct {
	Mechanism  ScramMechanism
	Iterations int
}



func (c *Client) DescribeUserScramCredentials(ctx context.Context, req *DescribeUserScramCredentialsRequest) (*DescribeUserScramCredentialsResponse, error) {
	users := make([]describeuserscramcredentials.RequestUser, len(req.Users))

	for userIdx, user := range req.Users {
		users[userIdx] = describeuserscramcredentials.RequestUser{
			Name: user.Name,
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &describeuserscramcredentials.Request{
		Users: users,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).DescribeUserScramCredentials: %w", err)
	}

	res := m.(*describeuserscramcredentials.Response)
	responseResults := make([]DescribeUserScramCredentialsResponseResult, len(res.Results))

	for responseIdx, responseResult := range res.Results {
		credentialInfos := make([]DescribeUserScramCredentialsCredentialInfo, len(responseResult.CredentialInfos))

		for credentialInfoIdx, credentialInfo := range responseResult.CredentialInfos {
			credentialInfos[credentialInfoIdx] = DescribeUserScramCredentialsCredentialInfo{
				Mechanism:  ScramMechanism(credentialInfo.Mechanism),
				Iterations: int(credentialInfo.Iterations),
			}
		}
		responseResults[responseIdx] = DescribeUserScramCredentialsResponseResult{
			User:            responseResult.User,
			CredentialInfos: credentialInfos,
			Error:           makeError(responseResult.ErrorCode, responseResult.ErrorMessage),
		}
	}
	ret := &DescribeUserScramCredentialsResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
		Error:    makeError(res.ErrorCode, res.ErrorMessage),
		Results:  responseResults,
	}

	return ret, nil
}
