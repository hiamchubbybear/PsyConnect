





package topology

import "go.mongodb.org/mongo-driver/mongo/description"


type hostlistDiff struct {
	Added   []string
	Removed []string
}


func diffHostList(t description.Topology, hostlist []string) hostlistDiff {
	var diff hostlistDiff

	oldServers := make(map[string]bool)
	for _, s := range t.Servers {
		oldServers[s.Addr.String()] = true
	}

	for _, addr := range hostlist {
		if oldServers[addr] {
			delete(oldServers, addr)
		} else {
			diff.Added = append(diff.Added, addr)
		}
	}

	for addr := range oldServers {
		diff.Removed = append(diff.Removed, addr)
	}

	return diff
}


type topologyDiff struct {
	Added   []description.Server
	Removed []description.Server
}


func diffTopology(old, new description.Topology) topologyDiff {
	var diff topologyDiff

	oldServers := make(map[string]bool)
	for _, s := range old.Servers {
		oldServers[s.Addr.String()] = true
	}

	for _, s := range new.Servers {
		addr := s.Addr.String()
		if oldServers[addr] {
			delete(oldServers, addr)
		} else {
			diff.Added = append(diff.Added, s)
		}
	}

	for _, s := range old.Servers {
		addr := s.Addr.String()
		if oldServers[addr] {
			diff.Removed = append(diff.Removed, s)
		}
	}

	return diff
}
