






package address 

import (
	"net"
	"strings"
)

const defaultPort = "27017"


type Address string



func (a Address) Network() string {
	if strings.HasSuffix(string(a), "sock") {
		return "unix"
	}
	return "tcp"
}



func (a Address) String() string {
	
	s := strings.ToLower(string(a))
	if len(s) == 0 {
		return ""
	}
	if a.Network() != "unix" {
		_, _, err := net.SplitHostPort(s)
		if err != nil && strings.Contains(err.Error(), "missing port in address") {
			s += ":" + defaultPort
		}
	}

	return s
}


func (a Address) Canonicalize() Address {
	return Address(a.String())
}
