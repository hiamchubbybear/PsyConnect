





package uuid

import (
	"encoding/hex"
	"io"

	"go.mongodb.org/mongo-driver/internal/randutil"
)


type UUID [16]byte



type source struct {
	random io.Reader
}


func (s *source) new() (UUID, error) {
	var uuid UUID
	_, err := io.ReadFull(s.random, uuid[:])
	if err != nil {
		return UUID{}, err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 
	uuid[8] = (uuid[8] & 0x3f) | 0x80 
	return uuid, nil
}



func newSource() *source {
	return &source{
		random: randutil.NewLockedRand(),
	}
}


var globalSource = newSource()





func New() (UUID, error) {
	return globalSource.new()
}

func (uuid UUID) String() string {
	var str [36]byte
	hex.Encode(str[:], uuid[:4])
	str[8] = '-'
	hex.Encode(str[9:13], uuid[4:6])
	str[13] = '-'
	hex.Encode(str[14:18], uuid[6:8])
	str[18] = '-'
	hex.Encode(str[19:23], uuid[8:10])
	str[23] = '-'
	hex.Encode(str[24:], uuid[10:])
	return string(str[:])
}
