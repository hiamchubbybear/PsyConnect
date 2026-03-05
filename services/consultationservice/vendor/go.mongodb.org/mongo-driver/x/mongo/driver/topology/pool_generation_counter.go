





package topology

import (
	"sync"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


const (
	generationDisconnected int64 = iota
	generationConnected
)



type generationStats struct {
	generation uint64
	numConns   uint64
}




type poolGenerationMap struct {
	
	
	
	state         int64
	generationMap map[primitive.ObjectID]*generationStats

	sync.Mutex
}

func newPoolGenerationMap() *poolGenerationMap {
	pgm := &poolGenerationMap{
		generationMap: make(map[primitive.ObjectID]*generationStats),
	}
	pgm.generationMap[primitive.NilObjectID] = &generationStats{}
	return pgm
}

func (p *poolGenerationMap) connect() {
	atomic.StoreInt64(&p.state, generationConnected)
}

func (p *poolGenerationMap) disconnect() {
	atomic.StoreInt64(&p.state, generationDisconnected)
}



func (p *poolGenerationMap) addConnection(serviceIDPtr *primitive.ObjectID) uint64 {
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	stats, ok := p.generationMap[serviceID]
	if ok {
		
		stats.numConns++
		return stats.generation
	}

	
	stats = &generationStats{
		numConns: 1,
	}
	p.generationMap[serviceID] = stats
	return 0
}

func (p *poolGenerationMap) removeConnection(serviceIDPtr *primitive.ObjectID) {
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	stats, ok := p.generationMap[serviceID]
	if !ok {
		return
	}

	
	
	
	stats.numConns--
	if stats.numConns == 0 {
		delete(p.generationMap, serviceID)
	}
}

func (p *poolGenerationMap) clear(serviceIDPtr *primitive.ObjectID) {
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok {
		stats.generation++
	}
}

func (p *poolGenerationMap) stale(serviceIDPtr *primitive.ObjectID, knownGeneration uint64) bool {
	
	if atomic.LoadInt64(&p.state) == generationDisconnected {
		return true
	}

	if generation, ok := p.getGeneration(serviceIDPtr); ok {
		return knownGeneration < generation
	}
	return false
}

func (p *poolGenerationMap) getGeneration(serviceIDPtr *primitive.ObjectID) (uint64, bool) {
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok {
		return stats.generation, true
	}
	return 0, false
}

func (p *poolGenerationMap) getNumConns(serviceIDPtr *primitive.ObjectID) uint64 {
	serviceID := getServiceID(serviceIDPtr)
	p.Lock()
	defer p.Unlock()

	if stats, ok := p.generationMap[serviceID]; ok {
		return stats.numConns
	}
	return 0
}

func getServiceID(oid *primitive.ObjectID) primitive.ObjectID {
	if oid == nil {
		return primitive.NilObjectID
	}
	return *oid
}
