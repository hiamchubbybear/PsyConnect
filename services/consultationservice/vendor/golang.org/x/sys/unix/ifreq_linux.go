



//go:build linux

package unix

import (
	"unsafe"
)













type Ifreq struct{ raw ifreq }




func NewIfreq(name string) (*Ifreq, error) {
	
	if len(name) >= IFNAMSIZ {
		return nil, EINVAL
	}

	var ifr ifreq
	copy(ifr.Ifrn[:], name)

	return &Ifreq{raw: ifr}, nil
}




func (ifr *Ifreq) Name() string {
	return ByteSliceToString(ifr.raw.Ifrn[:])
}








func (ifr *Ifreq) Inet4Addr() ([]byte, error) {
	raw := *(*RawSockaddrInet4)(unsafe.Pointer(&ifr.raw.Ifru[:SizeofSockaddrInet4][0]))
	if raw.Family != AF_INET {
		
		return nil, EINVAL
	}

	return raw.Addr[:], nil
}




func (ifr *Ifreq) SetInet4Addr(v []byte) error {
	if len(v) != 4 {
		return EINVAL
	}

	var addr [4]byte
	copy(addr[:], v)

	ifr.clear()
	*(*RawSockaddrInet4)(
		unsafe.Pointer(&ifr.raw.Ifru[:SizeofSockaddrInet4][0]),
	) = RawSockaddrInet4{
		
		Family: AF_INET,
		Addr:   addr,
	}

	return nil
}


func (ifr *Ifreq) Uint16() uint16 {
	return *(*uint16)(unsafe.Pointer(&ifr.raw.Ifru[:2][0]))
}


func (ifr *Ifreq) SetUint16(v uint16) {
	ifr.clear()
	*(*uint16)(unsafe.Pointer(&ifr.raw.Ifru[:2][0])) = v
}


func (ifr *Ifreq) Uint32() uint32 {
	return *(*uint32)(unsafe.Pointer(&ifr.raw.Ifru[:4][0]))
}


func (ifr *Ifreq) SetUint32(v uint32) {
	ifr.clear()
	*(*uint32)(unsafe.Pointer(&ifr.raw.Ifru[:4][0])) = v
}



func (ifr *Ifreq) clear() {
	for i := range ifr.raw.Ifru {
		ifr.raw.Ifru[i] = 0
	}
}






type ifreqData struct {
	name [IFNAMSIZ]byte
	
	
	
	data unsafe.Pointer
	
	_ [len(ifreq{}.Ifru) - SizeofPtr]byte
}



func (ifr Ifreq) withData(p unsafe.Pointer) ifreqData {
	return ifreqData{
		name: ifr.raw.Ifrn,
		data: p,
	}
}
