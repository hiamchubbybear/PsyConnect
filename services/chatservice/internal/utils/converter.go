package utils

import (
	"errors"
	"strconv"
)

type Converter struct{}

func (c *Converter) StringToInt64(rawString string) (int64, error) {
	num, err := strconv.ParseInt(rawString, 10, 64)
	if err != nil {
		return 0, errors.New("failed to convert from string to int64´")
	}
	return num, nil
}

func (c *Converter) StringToInt(rawString string) (int, error) {
	num, err := strconv.ParseInt(rawString, 10, 0)
	if err != nil {
		return 0, errors.New("failed to convert from string to int")
	}
	return int(num), nil
}
