





package description

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)


type Topology struct {
	Servers []Server
	SetName string
	Kind    TopologyKind
	
	SessionTimeoutMinutes    uint32
	SessionTimeoutMinutesPtr *int64
	CompatibilityErr         error
}


func (t Topology) String() string {
	var serversStr string
	for _, s := range t.Servers {
		serversStr += "{ " + s.String() + " }, "
	}
	return fmt.Sprintf("Type: %s, Servers: [%s]", t.Kind, serversStr)
}


func (t Topology) Equal(other Topology) bool {
	if t.Kind != other.Kind {
		return false
	}

	topoServers := make(map[string]Server)
	for _, s := range t.Servers {
		topoServers[s.Addr.String()] = s
	}

	otherServers := make(map[string]Server)
	for _, s := range other.Servers {
		otherServers[s.Addr.String()] = s
	}

	if len(topoServers) != len(otherServers) {
		return false
	}

	for _, server := range topoServers {
		otherServer := otherServers[server.Addr.String()]

		if !server.Equal(otherServer) {
			return false
		}
	}

	return true
}








func (t Topology) HasReadableServer(mode readpref.Mode) bool {
	switch t.Kind {
	case Single, Sharded:
		return hasAvailableServer(t.Servers, 0)
	case ReplicaSetWithPrimary:
		return hasAvailableServer(t.Servers, mode)
	case ReplicaSetNoPrimary, ReplicaSet:
		if mode == readpref.PrimaryMode {
			return false
		}
		
		if !mode.IsValid() {
			return false
		}

		return hasAvailableServer(t.Servers, mode)
	}
	return false
}







func (t Topology) HasWritableServer() bool {
	return t.HasReadableServer(readpref.PrimaryMode)
}


func hasAvailableServer(servers []Server, mode readpref.Mode) bool {
	switch mode {
	case readpref.PrimaryMode:
		for _, s := range servers {
			if s.Kind == RSPrimary {
				return true
			}
		}
		return false
	case readpref.PrimaryPreferredMode, readpref.SecondaryPreferredMode, readpref.NearestMode:
		for _, s := range servers {
			if s.Kind == RSPrimary || s.Kind == RSSecondary {
				return true
			}
		}
		return false
	case readpref.SecondaryMode:
		for _, s := range servers {
			if s.Kind == RSSecondary {
				return true
			}
		}
		return false
	}

	
	for _, s := range servers {
		switch s.Kind {
		case Standalone,
			RSMember,
			RSPrimary,
			RSSecondary,
			RSArbiter,
			RSGhost,
			Mongos:
			return true
		}
	}

	return false
}
