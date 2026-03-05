





package csfle

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	EncryptedCacheCollection      = "ecc"
	EncryptedStateCollection      = "esc"
	EncryptedCompactionCollection = "ecoc"
)


func GetEncryptedStateCollectionName(efBSON bsoncore.Document, dataCollectionName string, stateCollection string) (string, error) {
	fieldName := stateCollection + "Collection"
	val, err := efBSON.LookupErr(fieldName)
	if err != nil {
		if !errors.Is(err, bsoncore.ErrElementNotFound) {
			return "", err
		}
		
		defaultName := "enxcol_." + dataCollectionName + "." + stateCollection
		return defaultName, nil
	}

	stateCollectionName, ok := val.StringValueOK()
	if !ok {
		return "", fmt.Errorf("expected string for '%v', got: %v", fieldName, val.Type)
	}
	return stateCollectionName, nil
}
