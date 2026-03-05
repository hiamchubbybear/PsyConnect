





package description

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type TopologyVersion struct {
	ProcessID primitive.ObjectID
	Counter   int64
}


func NewTopologyVersion(doc bson.Raw) (*TopologyVersion, error) {
	elements, err := doc.Elements()
	if err != nil {
		return nil, err
	}
	var tv TopologyVersion
	var ok bool
	for _, element := range elements {
		switch element.Key() {
		case "processId":
			tv.ProcessID, ok = element.Value().ObjectIDOK()
			if !ok {
				return nil, fmt.Errorf("expected 'processId' to be a objectID but it's a BSON %s", element.Value().Type)
			}
		case "counter":
			tv.Counter, ok = element.Value().Int64OK()
			if !ok {
				return nil, fmt.Errorf("expected 'counter' to be an int64 but it's a BSON %s", element.Value().Type)
			}
		}
	}
	return &tv, nil
}






func (tv *TopologyVersion) CompareToIncoming(responseTV *TopologyVersion) int {
	if tv == nil || responseTV == nil {
		return -1
	}
	if tv.ProcessID != responseTV.ProcessID {
		return -1
	}
	if tv.Counter == responseTV.Counter {
		return 0
	}
	if tv.Counter < responseTV.Counter {
		return -1
	}
	return 1
}
