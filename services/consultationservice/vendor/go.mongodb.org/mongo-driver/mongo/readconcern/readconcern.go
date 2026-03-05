









package readconcern 

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)






type ReadConcern struct {
	Level string
}






type Option func(concern *ReadConcern)






func Level(level string) Option {
	return func(concern *ReadConcern) {
		concern.Level = level
	}
}






func Local() *ReadConcern {
	return New(Level("local"))
}






func Majority() *ReadConcern {
	return New(Level("majority"))
}






func Linearizable() *ReadConcern {
	return New(Level("linearizable"))
}






func Available() *ReadConcern {
	return New(Level("available"))
}






func Snapshot() *ReadConcern {
	return New(Level("snapshot"))
}






func New(options ...Option) *ReadConcern {
	concern := &ReadConcern{}

	for _, option := range options {
		option(concern)
	}

	return concern
}




func (rc *ReadConcern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if rc == nil {
		return 0, nil, errors.New("cannot marshal nil ReadConcern")
	}

	var elems []byte

	if len(rc.Level) > 0 {
		elems = bsoncore.AppendStringElement(elems, "level", rc.Level)
	}

	return bsontype.EmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
}




func (rc *ReadConcern) GetLevel() string {
	return rc.Level
}
