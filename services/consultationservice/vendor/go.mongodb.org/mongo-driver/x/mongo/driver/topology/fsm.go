





package topology

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/ptrutil"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
)

var (
	
	MinSupportedMongoDBVersion = "3.6"

	
	SupportedWireVersions = description.NewVersionRange(6, 25)
)

type fsm struct {
	description.Topology
	maxElectionID    primitive.ObjectID
	maxSetVersion    uint32
	compatible       atomic.Value
	compatibilityErr error
}

func newFSM() *fsm {
	f := fsm{}
	f.compatible.Store(true)
	return &f
}















func selectFSMSessionTimeout(f *fsm, s description.Server) *int64 {
	oldMinutes := f.SessionTimeoutMinutesPtr
	comp := ptrutil.CompareInt64(oldMinutes, s.SessionTimeoutMinutesPtr)

	
	
	
	
	
	
	
	if s.DataBearing() && (comp == 1 || comp == 2) {
		return s.SessionTimeoutMinutesPtr
	}

	
	
	
	if oldMinutes != nil {
		return oldMinutes
	}

	timeout := s.SessionTimeoutMinutesPtr
	for _, server := range f.Servers {
		
		
		if !server.DataBearing() {
			continue
		}

		srvTimeout := server.SessionTimeoutMinutesPtr
		comp := ptrutil.CompareInt64(timeout, srvTimeout)

		if comp <= 0 { 
			continue
		}

		timeout = server.SessionTimeoutMinutesPtr
	}

	return timeout
}







func (f *fsm) apply(s description.Server) (description.Topology, description.Server) {
	newServers := make([]description.Server, len(f.Servers))
	copy(newServers, f.Servers)

	
	
	serverTimeoutMinutes := selectFSMSessionTimeout(f, s)

	f.Topology = description.Topology{
		Kind:    f.Kind,
		Servers: newServers,
		SetName: f.SetName,
	}

	f.Topology.SessionTimeoutMinutesPtr = serverTimeoutMinutes

	if serverTimeoutMinutes != nil {
		f.SessionTimeoutMinutes = uint32(*serverTimeoutMinutes)
	}

	if _, ok := f.findServer(s.Addr); !ok {
		return f.Topology, s
	}

	updatedDesc := s
	switch f.Kind {
	case description.Unknown:
		updatedDesc = f.applyToUnknown(s)
	case description.Sharded:
		updatedDesc = f.applyToSharded(s)
	case description.ReplicaSetNoPrimary:
		updatedDesc = f.applyToReplicaSetNoPrimary(s)
	case description.ReplicaSetWithPrimary:
		updatedDesc = f.applyToReplicaSetWithPrimary(s)
	case description.Single:
		updatedDesc = f.applyToSingle(s)
	}

	for _, server := range f.Servers {
		if server.WireVersion != nil {
			if server.WireVersion.Max < SupportedWireVersions.Min {
				f.compatible.Store(false)
				f.compatibilityErr = fmt.Errorf(
					"server at %s reports wire version %d, but this version of the Go driver requires "+
						"at least %d (MongoDB %s)",
					server.Addr.String(),
					server.WireVersion.Max,
					SupportedWireVersions.Min,
					MinSupportedMongoDBVersion,
				)
				f.Topology.CompatibilityErr = f.compatibilityErr
				return f.Topology, s
			}

			if server.WireVersion.Min > SupportedWireVersions.Max {
				f.compatible.Store(false)
				f.compatibilityErr = fmt.Errorf(
					"server at %s requires wire version %d, but this version of the Go driver only supports up to %d",
					server.Addr.String(),
					server.WireVersion.Min,
					SupportedWireVersions.Max,
				)
				f.Topology.CompatibilityErr = f.compatibilityErr
				return f.Topology, s
			}
		}
	}

	f.compatible.Store(true)
	f.compatibilityErr = nil

	return f.Topology, updatedDesc
}

func (f *fsm) applyToReplicaSetNoPrimary(s description.Server) description.Server {
	switch s.Kind {
	case description.Standalone, description.Mongos:
		f.removeServerByAddr(s.Addr)
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.updateRSWithoutPrimary(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
	}

	return s
}

func (f *fsm) applyToReplicaSetWithPrimary(s description.Server) description.Server {
	switch s.Kind {
	case description.Standalone, description.Mongos:
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.updateRSWithPrimaryFromMember(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
		f.checkIfHasPrimary()
	}

	return s
}

func (f *fsm) applyToSharded(s description.Server) description.Server {
	switch s.Kind {
	case description.Mongos, description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		f.removeServerByAddr(s.Addr)
	}

	return s
}

func (f *fsm) applyToSingle(s description.Server) description.Server {
	switch s.Kind {
	case description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.Mongos:
		if f.SetName != "" {
			f.removeServerByAddr(s.Addr)
			return s
		}

		f.replaceServer(s)
	case description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		
		
		
		
		
		
		if f.SetName != "" && f.SetName != s.SetName {
			s = description.Server{
				Addr: s.Addr,
				Kind: description.Unknown,
			}
		}

		f.replaceServer(s)
	}

	return s
}

func (f *fsm) applyToUnknown(s description.Server) description.Server {
	switch s.Kind {
	case description.Mongos:
		f.setKind(description.Sharded)
		f.replaceServer(s)
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.setKind(description.ReplicaSetNoPrimary)
		f.updateRSWithoutPrimary(s)
	case description.Standalone:
		f.updateUnknownWithStandalone(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
	}

	return s
}

func (f *fsm) checkIfHasPrimary() {
	if _, ok := f.findPrimary(); ok {
		f.setKind(description.ReplicaSetWithPrimary)
	} else {
		f.setKind(description.ReplicaSetNoPrimary)
	}
}


func hasStalePrimary(fsm fsm, srv description.Server) bool {
	
	compRes := bytes.Compare(srv.ElectionID[:], fsm.maxElectionID[:])

	if wireVersion := srv.WireVersion; wireVersion != nil && wireVersion.Max >= 17 {
		
		
		
		
		return compRes == -1 || (compRes != 1 && srv.SetVersion < fsm.maxSetVersion)
	}

	
	
	
	return compRes == -1 || fsm.maxSetVersion > srv.SetVersion
}




func transferEVTuple(srv description.Server, fsm *fsm) bool {
	stalePrimary := hasStalePrimary(*fsm, srv)

	if wireVersion := srv.WireVersion; wireVersion != nil && wireVersion.Max >= 17 {
		if stalePrimary {
			fsm.checkIfHasPrimary()
			return false
		}

		fsm.maxElectionID = srv.ElectionID
		fsm.maxSetVersion = srv.SetVersion

		return true
	}

	if srv.SetVersion != 0 && !srv.ElectionID.IsZero() {
		if stalePrimary {
			fsm.replaceServer(description.Server{
				Addr: srv.Addr,
				LastError: fmt.Errorf(
					"was a primary, but its set version or election id is stale"),
			})

			fsm.checkIfHasPrimary()

			return false
		}

		fsm.maxElectionID = srv.ElectionID
	}

	if srv.SetVersion > fsm.maxSetVersion {
		fsm.maxSetVersion = srv.SetVersion
	}

	return true
}

func (f *fsm) updateRSFromPrimary(srv description.Server) {
	if f.SetName == "" {
		f.SetName = srv.SetName
	} else if f.SetName != srv.SetName {
		f.removeServerByAddr(srv.Addr)
		f.checkIfHasPrimary()

		return
	}

	if ok := transferEVTuple(srv, f); !ok {
		return
	}

	if j, ok := f.findPrimary(); ok {
		f.setServer(j, description.Server{
			Addr:      f.Servers[j].Addr,
			LastError: fmt.Errorf("was a primary, but a new primary was discovered"),
		})
	}

	f.replaceServer(srv)

	for j := len(f.Servers) - 1; j >= 0; j-- {
		found := false
		for _, member := range srv.Members {
			if member == f.Servers[j].Addr {
				found = true
				break
			}
		}

		if !found {
			f.removeServer(j)
		}
	}

	for _, member := range srv.Members {
		if _, ok := f.findServer(member); !ok {
			f.addServer(member)
		}
	}

	f.checkIfHasPrimary()
}

func (f *fsm) updateRSWithPrimaryFromMember(s description.Server) {
	if f.SetName != s.SetName {
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	}

	if s.Addr != s.CanonicalAddr {
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	}

	f.replaceServer(s)

	if _, ok := f.findPrimary(); !ok {
		f.setKind(description.ReplicaSetNoPrimary)
	}
}

func (f *fsm) updateRSWithoutPrimary(s description.Server) {
	if f.SetName == "" {
		f.SetName = s.SetName
	} else if f.SetName != s.SetName {
		f.removeServerByAddr(s.Addr)
		return
	}

	for _, member := range s.Members {
		if _, ok := f.findServer(member); !ok {
			f.addServer(member)
		}
	}

	if s.Addr != s.CanonicalAddr {
		f.removeServerByAddr(s.Addr)
		return
	}

	f.replaceServer(s)
}

func (f *fsm) updateUnknownWithStandalone(s description.Server) {
	if len(f.Servers) > 1 {
		f.removeServerByAddr(s.Addr)
		return
	}

	f.setKind(description.Single)
	f.replaceServer(s)
}

func (f *fsm) addServer(addr address.Address) {
	f.Servers = append(f.Servers, description.Server{
		Addr: addr.Canonicalize(),
	})
}

func (f *fsm) findPrimary() (int, bool) {
	for i, s := range f.Servers {
		if s.Kind == description.RSPrimary {
			return i, true
		}
	}

	return 0, false
}

func (f *fsm) findServer(addr address.Address) (int, bool) {
	canon := addr.Canonicalize()
	for i, s := range f.Servers {
		if canon == s.Addr {
			return i, true
		}
	}

	return 0, false
}

func (f *fsm) removeServer(i int) {
	f.Servers = append(f.Servers[:i], f.Servers[i+1:]...)
}

func (f *fsm) removeServerByAddr(addr address.Address) {
	if i, ok := f.findServer(addr); ok {
		f.removeServer(i)
	}
}

func (f *fsm) replaceServer(s description.Server) {
	if i, ok := f.findServer(s.Addr); ok {
		f.setServer(i, s)
	}
}

func (f *fsm) setServer(i int, s description.Server) {
	f.Servers[i] = s
}

func (f *fsm) setKind(k description.TopologyKind) {
	f.Kind = k
}
