





package session

import (
	"time"

	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type Server struct {
	SessionID bsoncore.Document
	TxnNumber int64
	LastUsed  time.Time
	Dirty     bool
}



func (ss *Server) expired(topoDesc topologyDescription) bool {
	
	
	if topoDesc.kind == description.LoadBalanced {
		return false
	}

	if topoDesc.timeoutMinutes == nil || *topoDesc.timeoutMinutes <= 0 {
		return true
	}
	timeUnused := time.Since(ss.LastUsed).Minutes()
	return timeUnused > float64(*topoDesc.timeoutMinutes-1)
}



func (ss *Server) updateUseTime() {
	ss.LastUsed = time.Now()
}

func newServerSession() (*Server, error) {
	id, err := uuid.New()
	if err != nil {
		return nil, err
	}

	idx, idDoc := bsoncore.AppendDocumentStart(nil)
	idDoc = bsoncore.AppendBinaryElement(idDoc, "id", UUIDSubtype, id[:])
	idDoc, _ = bsoncore.AppendDocumentEnd(idDoc, idx)

	return &Server{
		SessionID: idDoc,
		LastUsed:  time.Now(),
	}, nil
}


func (ss *Server) IncrementTxnNumber() {
	ss.TxnNumber++
}


func (ss *Server) MarkDirty() {
	ss.Dirty = true
}


const UUIDSubtype byte = 4
