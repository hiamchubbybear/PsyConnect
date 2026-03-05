





package driver

import (
	"errors"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)



var ErrDocumentTooLarge = errors.New("an inserted document is too large")



type Batches struct {
	Identifier string
	Documents  []bsoncore.Document
	Current    []bsoncore.Document
	Ordered    *bool
}



func (b *Batches) Valid() bool { return b != nil && b.Identifier != "" && len(b.Documents) > 0 }



func (b *Batches) ClearBatch() { b.Current = b.Current[:0] }







func (b *Batches) AdvanceBatch(maxCount, targetBatchSize, maxDocSize int) error {
	if len(b.Current) > 0 {
		return nil
	}

	if maxCount <= 0 {
		maxCount = 1
	}

	splitAfter := 0
	size := 0
	for i, doc := range b.Documents {
		if i == maxCount {
			break
		}
		if len(doc) > maxDocSize {
			return ErrDocumentTooLarge
		}
		if size+len(doc) > targetBatchSize {
			break
		}

		size += len(doc)
		splitAfter++
	}

	
	
	if splitAfter == 0 {
		splitAfter = 1
	}

	b.Current, b.Documents = b.Documents[:splitAfter], b.Documents[splitAfter:]
	return nil
}
