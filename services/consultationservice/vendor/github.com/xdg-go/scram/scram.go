





package scram

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/xdg-go/stringprep"
)





type HashGeneratorFcn func() hash.Hash



var SHA1 HashGeneratorFcn = func() hash.Hash { return sha1.New() }



var SHA256 HashGeneratorFcn = func() hash.Hash { return sha256.New() }



var SHA512 HashGeneratorFcn = func() hash.Hash { return sha512.New() }





func (f HashGeneratorFcn) NewClient(username, password, authzID string) (*Client, error) {
	var userprep, passprep, authprep string
	var err error

	if userprep, err = stringprep.SASLprep.Prepare(username); err != nil {
		return nil, fmt.Errorf("Error SASLprepping username '%s': %v", username, err)
	}
	if passprep, err = stringprep.SASLprep.Prepare(password); err != nil {
		return nil, fmt.Errorf("Error SASLprepping password '%s': %v", password, err)
	}
	if authprep, err = stringprep.SASLprep.Prepare(authzID); err != nil {
		return nil, fmt.Errorf("Error SASLprepping authzID '%s': %v", authzID, err)
	}

	return newClient(userprep, passprep, authprep, f), nil
}




func (f HashGeneratorFcn) NewClientUnprepped(username, password, authzID string) (*Client, error) {
	return newClient(username, password, authzID, f), nil
}





func (f HashGeneratorFcn) NewServer(cl CredentialLookup) (*Server, error) {
	return newServer(cl, f)
}
