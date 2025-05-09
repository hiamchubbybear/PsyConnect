package utils

import (
	"errors"
	"strconv"
)

type Converter struct{}

func (c Converter) StringToInt64(rawString string) (int64, error) {
	num, err := strconv.ParseInt(rawString, 10, 64)
	if err != nil {
		return 0, errors.New("Failed to convert from string to int6")
	}
	return num, nil
}
