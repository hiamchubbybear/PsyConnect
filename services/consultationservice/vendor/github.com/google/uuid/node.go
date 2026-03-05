



package uuid

import (
	"sync"
)

var (
	nodeMu sync.Mutex
	ifname string  
	nodeID [6]byte 
	zeroID [6]byte 
)




func NodeInterface() string {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	return ifname
}







func SetNodeInterface(name string) bool {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	return setNodeInterface(name)
}

func setNodeInterface(name string) bool {
	iname, addr := getHardwareInterface(name) 
	if iname != "" && addr != nil {
		ifname = iname
		copy(nodeID[:], addr)
		return true
	}

	
	
	
	if name == "" {
		ifname = "random"
		randomBits(nodeID[:])
		return true
	}
	return false
}



func NodeID() []byte {
	defer nodeMu.Unlock()
	nodeMu.Lock()
	if nodeID == zeroID {
		setNodeInterface("")
	}
	nid := nodeID
	return nid[:]
}




func SetNodeID(id []byte) bool {
	if len(id) < 6 {
		return false
	}
	defer nodeMu.Unlock()
	nodeMu.Lock()
	copy(nodeID[:], id)
	ifname = "user"
	return true
}



func (uuid UUID) NodeID() []byte {
	var node [6]byte
	copy(node[:], uuid[10:])
	return node[:]
}
