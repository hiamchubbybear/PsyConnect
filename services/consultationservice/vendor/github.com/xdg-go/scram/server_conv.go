





package scram

import (
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
)

type serverState int

const (
	serverFirst serverState = iota
	serverFinal
	serverDone
)




type ServerConversation struct {
	nonceGen     NonceGeneratorFcn
	hashGen      HashGeneratorFcn
	credentialCB CredentialLookup
	state        serverState
	credential   StoredCredentials
	valid        bool
	gs2Header    string
	username     string
	authzID      string
	nonce        string
	c1b          string
	s1           string
}





func (sc *ServerConversation) Step(challenge string) (response string, err error) {
	switch sc.state {
	case serverFirst:
		sc.state = serverFinal
		response, err = sc.firstMsg(challenge)
	case serverFinal:
		sc.state = serverDone
		response, err = sc.finalMsg(challenge)
	default:
		response, err = "", errors.New("Conversation already completed")
	}
	return
}


func (sc *ServerConversation) Done() bool {
	return sc.state == serverDone
}



func (sc *ServerConversation) Valid() bool {
	return sc.valid
}



func (sc *ServerConversation) Username() string {
	return sc.username
}




func (sc *ServerConversation) AuthzID() string {
	return sc.authzID
}

func (sc *ServerConversation) firstMsg(c1 string) (string, error) {
	msg, err := parseClientFirst(c1)
	if err != nil {
		sc.state = serverDone
		return "", err
	}

	sc.gs2Header = msg.gs2Header
	sc.username = msg.username
	sc.authzID = msg.authzID

	sc.credential, err = sc.credentialCB(msg.username)
	if err != nil {
		sc.state = serverDone
		return "e=unknown-user", err
	}

	sc.nonce = msg.nonce + sc.nonceGen()
	sc.c1b = msg.c1b
	sc.s1 = fmt.Sprintf("r=%s,s=%s,i=%d",
		sc.nonce,
		base64.StdEncoding.EncodeToString([]byte(sc.credential.Salt)),
		sc.credential.Iters,
	)

	return sc.s1, nil
}



func (sc *ServerConversation) finalMsg(c2 string) (string, error) {
	msg, err := parseClientFinal(c2)
	if err != nil {
		return "", err
	}

	
	
	
	
	if string(msg.cbind) != sc.gs2Header {
		return "e=channel-bindings-dont-match", fmt.Errorf("channel binding received '%s' doesn't match expected '%s'", msg.cbind, sc.gs2Header)
	}

	
	if msg.nonce != sc.nonce {
		return "e=other-error", errors.New("nonce received did not match nonce sent")
	}

	
	authMsg := sc.c1b + "," + sc.s1 + "," + msg.c2wop

	
	clientSignature := computeHMAC(sc.hashGen, sc.credential.StoredKey, []byte(authMsg))
	clientKey := xorBytes([]byte(msg.proof), clientSignature)
	storedKey := computeHash(sc.hashGen, clientKey)

	
	if !hmac.Equal(storedKey, sc.credential.StoredKey) {
		return "e=invalid-proof", errors.New("challenge proof invalid")
	}

	sc.valid = true

	
	serverSignature := computeHMAC(sc.hashGen, sc.credential.ServerKey, []byte(authMsg))
	return "v=" + base64.StdEncoding.EncodeToString(serverSignature), nil
}
