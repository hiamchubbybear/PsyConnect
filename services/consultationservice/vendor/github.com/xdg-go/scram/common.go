





package scram

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"strings"
)





type NonceGeneratorFcn func() string



type derivedKeys struct {
	ClientKey []byte
	StoredKey []byte
	ServerKey []byte
}





type KeyFactors struct {
	Salt  string
	Iters int
}









type StoredCredentials struct {
	KeyFactors
	StoredKey []byte
	ServerKey []byte
}







type CredentialLookup func(string) (StoredCredentials, error)

func defaultNonceGenerator() string {
	raw := make([]byte, 24)
	nonce := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	rand.Read(raw)
	base64.StdEncoding.Encode(nonce, raw)
	return string(nonce)
}

func encodeName(s string) string {
	return strings.Replace(strings.Replace(s, "=", "=3D", -1), ",", "=2C", -1)
}

func decodeName(s string) (string, error) {
	
	return strings.Replace(strings.Replace(s, "=2C", ",", -1), "=3D", "=", -1), nil
}

func computeHash(hg HashGeneratorFcn, b []byte) []byte {
	h := hg()
	h.Write(b)
	return h.Sum(nil)
}

func computeHMAC(hg HashGeneratorFcn, key, data []byte) []byte {
	mac := hmac.New(hg, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func xorBytes(a, b []byte) []byte {
	
	xor := make([]byte, len(a))
	for i := range a {
		xor[i] = a[i] ^ b[i]
	}
	return xor
}
