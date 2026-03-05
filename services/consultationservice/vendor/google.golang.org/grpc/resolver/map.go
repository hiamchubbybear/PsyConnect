

package resolver

type addressMapEntry struct {
	addr  Address
	value any
}





type AddressMap struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	m map[Address]addressMapEntryList
}

func toMapKey(addr *Address) Address {
	return Address{Addr: addr.Addr, ServerName: addr.ServerName}
}

type addressMapEntryList []*addressMapEntry


func NewAddressMap() *AddressMap {
	return &AddressMap{m: make(map[Address]addressMapEntryList)}
}



func (l addressMapEntryList) find(addr Address) int {
	for i, entry := range l {
		
		
		if entry.addr.Attributes.Equal(addr.Attributes) {
			return i
		}
	}
	return -1
}


func (a *AddressMap) Get(addr Address) (value any, ok bool) {
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	if entry := entryList.find(addr); entry != -1 {
		return entryList[entry].value, true
	}
	return nil, false
}


func (a *AddressMap) Set(addr Address, value any) {
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	if entry := entryList.find(addr); entry != -1 {
		entryList[entry].value = value
		return
	}
	a.m[addrKey] = append(entryList, &addressMapEntry{addr: addr, value: value})
}


func (a *AddressMap) Delete(addr Address) {
	addrKey := toMapKey(&addr)
	entryList := a.m[addrKey]
	entry := entryList.find(addr)
	if entry == -1 {
		return
	}
	if len(entryList) == 1 {
		entryList = nil
	} else {
		copy(entryList[entry:], entryList[entry+1:])
		entryList = entryList[:len(entryList)-1]
	}
	a.m[addrKey] = entryList
}


func (a *AddressMap) Len() int {
	ret := 0
	for _, entryList := range a.m {
		ret += len(entryList)
	}
	return ret
}


func (a *AddressMap) Keys() []Address {
	ret := make([]Address, 0, a.Len())
	for _, entryList := range a.m {
		for _, entry := range entryList {
			ret = append(ret, entry.addr)
		}
	}
	return ret
}


func (a *AddressMap) Values() []any {
	ret := make([]any, 0, a.Len())
	for _, entryList := range a.m {
		for _, entry := range entryList {
			ret = append(ret, entry.value)
		}
	}
	return ret
}

type endpointNode struct {
	addrs map[string]struct{}
}



func (en *endpointNode) Equal(en2 *endpointNode) bool {
	if len(en.addrs) != len(en2.addrs) {
		return false
	}
	for addr := range en.addrs {
		if _, ok := en2.addrs[addr]; !ok {
			return false
		}
	}
	return true
}

func toEndpointNode(endpoint Endpoint) endpointNode {
	en := make(map[string]struct{})
	for _, addr := range endpoint.Addresses {
		en[addr.Addr] = struct{}{}
	}
	return endpointNode{
		addrs: en,
	}
}





type EndpointMap struct {
	endpoints map[*endpointNode]any
}


func NewEndpointMap() *EndpointMap {
	return &EndpointMap{
		endpoints: make(map[*endpointNode]any),
	}
}


func (em *EndpointMap) Get(e Endpoint) (value any, ok bool) {
	en := toEndpointNode(e)
	if endpoint := em.find(en); endpoint != nil {
		return em.endpoints[endpoint], true
	}
	return nil, false
}


func (em *EndpointMap) Set(e Endpoint, value any) {
	en := toEndpointNode(e)
	if endpoint := em.find(en); endpoint != nil {
		em.endpoints[endpoint] = value
		return
	}
	em.endpoints[&en] = value
}


func (em *EndpointMap) Len() int {
	return len(em.endpoints)
}






func (em *EndpointMap) Keys() []Endpoint {
	ret := make([]Endpoint, 0, len(em.endpoints))
	for en := range em.endpoints {
		var endpoint Endpoint
		for addr := range en.addrs {
			endpoint.Addresses = append(endpoint.Addresses, Address{Addr: addr})
		}
		ret = append(ret, endpoint)
	}
	return ret
}


func (em *EndpointMap) Values() []any {
	ret := make([]any, 0, len(em.endpoints))
	for _, val := range em.endpoints {
		ret = append(ret, val)
	}
	return ret
}




func (em EndpointMap) find(e endpointNode) *endpointNode {
	for endpoint := range em.endpoints {
		if e.Equal(endpoint) {
			return endpoint
		}
	}
	return nil
}


func (em *EndpointMap) Delete(e Endpoint) {
	en := toEndpointNode(e)
	if entry := em.find(en); entry != nil {
		delete(em.endpoints, entry)
	}
}
