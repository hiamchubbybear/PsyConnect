





package session

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)


type ClusterClock struct {
	clusterTime bson.Raw
	lock        sync.Mutex
}


func (cc *ClusterClock) GetClusterTime() bson.Raw {
	var ct bson.Raw
	cc.lock.Lock()
	ct = cc.clusterTime
	cc.lock.Unlock()

	return ct
}


func (cc *ClusterClock) AdvanceClusterTime(clusterTime bson.Raw) {
	cc.lock.Lock()
	cc.clusterTime = MaxClusterTime(cc.clusterTime, clusterTime)
	cc.lock.Unlock()
}
