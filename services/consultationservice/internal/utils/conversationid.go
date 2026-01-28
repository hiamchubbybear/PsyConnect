package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

type Encoder struct {
}

func New() *Encoder {
	return &Encoder{}
}
func (e *Encoder) EncodeConversationId(sender, receiver string) (string, error) {

	if sender == "" || receiver == "" {
		return "", errors.New("invalid UUID ")
	}

	uuids := []string{sender, receiver}
	sort.Strings(uuids)
	hash := sha256.Sum256([]byte(uuids[0] + uuids[1]))
	hexHash := hex.EncodeToString(hash[:])

	return strings.Join([]string{
		hexHash[0:8],
		hexHash[8:12],
		hexHash[12:16],
		hexHash[16:20],
		hexHash[20:32],
	}, "-"), nil

}
