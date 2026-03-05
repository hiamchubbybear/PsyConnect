





package unix

import "unsafe"




func UnixCredentials(ucred *Ucred) []byte {
	b := make([]byte, CmsgSpace(SizeofUcred))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = SOL_SOCKET
	h.Type = SCM_CREDENTIALS
	h.SetLen(CmsgLen(SizeofUcred))
	*(*Ucred)(h.data(0)) = *ucred
	return b
}




func ParseUnixCredentials(m *SocketControlMessage) (*Ucred, error) {
	if m.Header.Level != SOL_SOCKET {
		return nil, EINVAL
	}
	if m.Header.Type != SCM_CREDENTIALS {
		return nil, EINVAL
	}
	ucred := *(*Ucred)(unsafe.Pointer(&m.Data[0]))
	return &ucred, nil
}


func PktInfo4(info *Inet4Pktinfo) []byte {
	b := make([]byte, CmsgSpace(SizeofInet4Pktinfo))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = SOL_IP
	h.Type = IP_PKTINFO
	h.SetLen(CmsgLen(SizeofInet4Pktinfo))
	*(*Inet4Pktinfo)(h.data(0)) = *info
	return b
}


func PktInfo6(info *Inet6Pktinfo) []byte {
	b := make([]byte, CmsgSpace(SizeofInet6Pktinfo))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = SOL_IPV6
	h.Type = IPV6_PKTINFO
	h.SetLen(CmsgLen(SizeofInet6Pktinfo))
	*(*Inet6Pktinfo)(h.data(0)) = *info
	return b
}
