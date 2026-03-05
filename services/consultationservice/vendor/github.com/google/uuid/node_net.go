



// +build !js

package uuid

import "net"

var interfaces []net.Interface 







func getHardwareInterface(name string) (string, []byte) {
	if interfaces == nil {
		var err error
		interfaces, err = net.Interfaces()
		if err != nil {
			return "", nil
		}
	}
	for _, ifs := range interfaces {
		if len(ifs.HardwareAddr) >= 6 && (name == "" || name == ifs.Name) {
			return ifs.Name, ifs.HardwareAddr
		}
	}
	return "", nil
}
