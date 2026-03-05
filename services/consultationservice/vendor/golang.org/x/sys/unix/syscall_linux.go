










package unix

import (
	"encoding/binary"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)



func Access(path string, mode uint32) (err error) {
	return Faccessat(AT_FDCWD, path, mode, 0)
}

func Chmod(path string, mode uint32) (err error) {
	return Fchmodat(AT_FDCWD, path, mode, 0)
}

func Chown(path string, uid int, gid int) (err error) {
	return Fchownat(AT_FDCWD, path, uid, gid, 0)
}

func Creat(path string, mode uint32) (fd int, err error) {
	return Open(path, O_CREAT|O_WRONLY|O_TRUNC, mode)
}

func EpollCreate(size int) (fd int, err error) {
	if size <= 0 {
		return -1, EINVAL
	}
	return EpollCreate1(0)
}




func FanotifyMark(fd int, flags uint, mask uint64, dirFd int, pathname string) (err error) {
	if pathname == "" {
		return fanotifyMark(fd, flags, mask, dirFd, nil)
	}
	p, err := BytePtrFromString(pathname)
	if err != nil {
		return err
	}
	return fanotifyMark(fd, flags, mask, dirFd, p)
}




func Fchmodat(dirfd int, path string, mode uint32, flags int) error {
	
	
	if flags != 0 {
		err := fchmodat2(dirfd, path, mode, flags)
		if err == ENOSYS {
			
			
			if flags&^(AT_SYMLINK_NOFOLLOW|AT_EMPTY_PATH) != 0 {
				return EINVAL
			} else if flags&(AT_SYMLINK_NOFOLLOW|AT_EMPTY_PATH) != 0 {
				return EOPNOTSUPP
			}
		}
		return err
	}
	return fchmodat(dirfd, path, mode)
}

func InotifyInit() (fd int, err error) {
	return InotifyInit1(0)
}
















func Link(oldpath string, newpath string) (err error) {
	return Linkat(AT_FDCWD, oldpath, AT_FDCWD, newpath, 0)
}

func Mkdir(path string, mode uint32) (err error) {
	return Mkdirat(AT_FDCWD, path, mode)
}

func Mknod(path string, mode uint32, dev int) (err error) {
	return Mknodat(AT_FDCWD, path, mode, dev)
}

func Open(path string, mode int, perm uint32) (fd int, err error) {
	return openat(AT_FDCWD, path, mode|O_LARGEFILE, perm)
}



func Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) {
	return openat(dirfd, path, flags|O_LARGEFILE, mode)
}



func Openat2(dirfd int, path string, how *OpenHow) (fd int, err error) {
	return openat2(dirfd, path, how, SizeofOpenHow)
}

func Pipe(p []int) error {
	return Pipe2(p, 0)
}



func Pipe2(p []int, flags int) error {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err := pipe2(&pp, flags)
	if err == nil {
		p[0] = int(pp[0])
		p[1] = int(pp[1])
	}
	return err
}



func Ppoll(fds []PollFd, timeout *Timespec, sigmask *Sigset_t) (n int, err error) {
	if len(fds) == 0 {
		return ppoll(nil, 0, timeout, sigmask)
	}
	return ppoll(&fds[0], len(fds), timeout, sigmask)
}

func Poll(fds []PollFd, timeout int) (n int, err error) {
	var ts *Timespec
	if timeout >= 0 {
		ts = new(Timespec)
		*ts = NsecToTimespec(int64(timeout) * 1e6)
	}
	return Ppoll(fds, ts, nil)
}



func Readlink(path string, buf []byte) (n int, err error) {
	return Readlinkat(AT_FDCWD, path, buf)
}

func Rename(oldpath string, newpath string) (err error) {
	return Renameat(AT_FDCWD, oldpath, AT_FDCWD, newpath)
}

func Rmdir(path string) error {
	return Unlinkat(AT_FDCWD, path, AT_REMOVEDIR)
}



func Symlink(oldpath string, newpath string) (err error) {
	return Symlinkat(oldpath, AT_FDCWD, newpath)
}

func Unlink(path string) error {
	return Unlinkat(AT_FDCWD, path, 0)
}



func Utimes(path string, tv []Timeval) error {
	if tv == nil {
		err := utimensat(AT_FDCWD, path, nil, 0)
		if err != ENOSYS {
			return err
		}
		return utimes(path, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	var ts [2]Timespec
	ts[0] = NsecToTimespec(TimevalToNsec(tv[0]))
	ts[1] = NsecToTimespec(TimevalToNsec(tv[1]))
	err := utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
	if err != ENOSYS {
		return err
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}



func UtimesNano(path string, ts []Timespec) error {
	return UtimesNanoAt(AT_FDCWD, path, ts, 0)
}

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error {
	if ts == nil {
		return utimensat(dirfd, path, nil, flags)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
}

func Futimesat(dirfd int, path string, tv []Timeval) error {
	if tv == nil {
		return futimesat(dirfd, path, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return futimesat(dirfd, path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

func Futimes(fd int, tv []Timeval) (err error) {
	
	
	return Utimes("/proc/self/fd/"+strconv.Itoa(fd), tv)
}

const ImplementsGetwd = true



func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	n, err := Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	
	if n < 1 || n > len(buf) || buf[n-1] != 0 {
		return "", EINVAL
	}
	
	
	
	if buf[0] != '/' {
		return "", ENOENT
	}

	return string(buf[0 : n-1]), nil
}

func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	
	if n < 0 || n > 1<<20 {
		return nil, EINVAL
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil {
		return nil, err
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

func Setgroups(gids []int) (err error) {
	if len(gids) == 0 {
		return setgroups(0, nil)
	}

	a := make([]_Gid_t, len(gids))
	for i, v := range gids {
		a[i] = _Gid_t(v)
	}
	return setgroups(len(a), &a[0])
}

type WaitStatus uint32










const (
	mask    = 0x7F
	core    = 0x80
	exited  = 0x00
	stopped = 0x7F
	shift   = 8
)

func (w WaitStatus) Exited() bool { return w&mask == exited }

func (w WaitStatus) Signaled() bool { return w&mask != stopped && w&mask != exited }

func (w WaitStatus) Stopped() bool { return w&0xFF == stopped }

func (w WaitStatus) Continued() bool { return w == 0xFFFF }

func (w WaitStatus) CoreDump() bool { return w.Signaled() && w&core != 0 }

func (w WaitStatus) ExitStatus() int {
	if !w.Exited() {
		return -1
	}
	return int(w>>shift) & 0xFF
}

func (w WaitStatus) Signal() syscall.Signal {
	if !w.Signaled() {
		return -1
	}
	return syscall.Signal(w & mask)
}

func (w WaitStatus) StopSignal() syscall.Signal {
	if !w.Stopped() {
		return -1
	}
	return syscall.Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int {
	if w.StopSignal() != SIGTRAP {
		return -1
	}
	return int(w>>shift) >> 8
}



func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	var status _C_int
	wpid, err = wait4(pid, &status, options, rusage)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return
}



func Mkfifo(path string, mode uint32) error {
	return Mknod(path, mode|S_IFIFO, 0)
}

func Mkfifoat(dirfd int, path string, mode uint32) error {
	return Mknodat(dirfd, path, mode|S_IFIFO, 0)
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	
	sl := _Socklen(2)
	if n > 0 {
		sl += _Socklen(n) + 1
	}
	if sa.raw.Path[0] == '@' || (sa.raw.Path[0] == 0 && sl > 3) {
		
		sa.raw.Path[0] = 0
		
		sl--
	}

	return unsafe.Pointer(&sa.raw), sl, nil
}


type SockaddrLinklayer struct {
	Protocol uint16
	Ifindex  int
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]byte
	raw      RawSockaddrLinklayer
}

func (sa *SockaddrLinklayer) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_PACKET
	sa.raw.Protocol = sa.Protocol
	sa.raw.Ifindex = int32(sa.Ifindex)
	sa.raw.Hatype = sa.Hatype
	sa.raw.Pkttype = sa.Pkttype
	sa.raw.Halen = sa.Halen
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrLinklayer, nil
}


type SockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
	raw    RawSockaddrNetlink
}

func (sa *SockaddrNetlink) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_NETLINK
	sa.raw.Pad = sa.Pad
	sa.raw.Pid = sa.Pid
	sa.raw.Groups = sa.Groups
	return unsafe.Pointer(&sa.raw), SizeofSockaddrNetlink, nil
}



type SockaddrHCI struct {
	Dev     uint16
	Channel uint16
	raw     RawSockaddrHCI
}

func (sa *SockaddrHCI) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_BLUETOOTH
	sa.raw.Dev = sa.Dev
	sa.raw.Channel = sa.Channel
	return unsafe.Pointer(&sa.raw), SizeofSockaddrHCI, nil
}



type SockaddrL2 struct {
	PSM      uint16
	CID      uint16
	Addr     [6]uint8
	AddrType uint8
	raw      RawSockaddrL2
}

func (sa *SockaddrL2) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_BLUETOOTH
	psm := (*[2]byte)(unsafe.Pointer(&sa.raw.Psm))
	psm[0] = byte(sa.PSM)
	psm[1] = byte(sa.PSM >> 8)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Bdaddr[i] = sa.Addr[len(sa.Addr)-1-i]
	}
	cid := (*[2]byte)(unsafe.Pointer(&sa.raw.Cid))
	cid[0] = byte(sa.CID)
	cid[1] = byte(sa.CID >> 8)
	sa.raw.Bdaddr_type = sa.AddrType
	return unsafe.Pointer(&sa.raw), SizeofSockaddrL2, nil
}
























type SockaddrRFCOMM struct {
	
	Addr [6]uint8

	
	
	Channel uint8

	raw RawSockaddrRFCOMM
}

func (sa *SockaddrRFCOMM) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_BLUETOOTH
	sa.raw.Channel = sa.Channel
	sa.raw.Bdaddr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrRFCOMM, nil
}


















type SockaddrCAN struct {
	Ifindex int
	RxID    uint32
	TxID    uint32
	raw     RawSockaddrCAN
}

func (sa *SockaddrCAN) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_CAN
	sa.raw.Ifindex = int32(sa.Ifindex)
	rx := (*[4]byte)(unsafe.Pointer(&sa.RxID))
	for i := 0; i < 4; i++ {
		sa.raw.Addr[i] = rx[i]
	}
	tx := (*[4]byte)(unsafe.Pointer(&sa.TxID))
	for i := 0; i < 4; i++ {
		sa.raw.Addr[i+4] = tx[i]
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrCAN, nil
}





type SockaddrCANJ1939 struct {
	Ifindex int
	Name    uint64
	PGN     uint32
	Addr    uint8
	raw     RawSockaddrCAN
}

func (sa *SockaddrCANJ1939) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_CAN
	sa.raw.Ifindex = int32(sa.Ifindex)
	n := (*[8]byte)(unsafe.Pointer(&sa.Name))
	for i := 0; i < 8; i++ {
		sa.raw.Addr[i] = n[i]
	}
	p := (*[4]byte)(unsafe.Pointer(&sa.PGN))
	for i := 0; i < 4; i++ {
		sa.raw.Addr[i+8] = p[i]
	}
	sa.raw.Addr[12] = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrCAN, nil
}
































































type SockaddrALG struct {
	Type    string
	Name    string
	Feature uint32
	Mask    uint32
	raw     RawSockaddrALG
}

func (sa *SockaddrALG) sockaddr() (unsafe.Pointer, _Socklen, error) {
	
	if len(sa.Type) > len(sa.raw.Type)-1 {
		return nil, 0, EINVAL
	}
	if len(sa.Name) > len(sa.raw.Name)-1 {
		return nil, 0, EINVAL
	}

	sa.raw.Family = AF_ALG
	sa.raw.Feat = sa.Feature
	sa.raw.Mask = sa.Mask

	copy(sa.raw.Type[:], sa.Type)
	copy(sa.raw.Name[:], sa.Name)

	return unsafe.Pointer(&sa.raw), SizeofSockaddrALG, nil
}





type SockaddrVM struct {
	
	
	
	
	
	CID   uint32
	Port  uint32
	Flags uint8
	raw   RawSockaddrVM
}

func (sa *SockaddrVM) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_VSOCK
	sa.raw.Port = sa.Port
	sa.raw.Cid = sa.CID
	sa.raw.Flags = sa.Flags

	return unsafe.Pointer(&sa.raw), SizeofSockaddrVM, nil
}

type SockaddrXDP struct {
	Flags        uint16
	Ifindex      uint32
	QueueID      uint32
	SharedUmemFD uint32
	raw          RawSockaddrXDP
}

func (sa *SockaddrXDP) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_XDP
	sa.raw.Flags = sa.Flags
	sa.raw.Ifindex = sa.Ifindex
	sa.raw.Queue_id = sa.QueueID
	sa.raw.Shared_umem_fd = sa.SharedUmemFD

	return unsafe.Pointer(&sa.raw), SizeofSockaddrXDP, nil
}








const px_proto_oe = 0

type SockaddrPPPoE struct {
	SID    uint16
	Remote []byte
	Dev    string
	raw    RawSockaddrPPPoX
}

func (sa *SockaddrPPPoE) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if len(sa.Remote) != 6 {
		return nil, 0, EINVAL
	}
	if len(sa.Dev) > IFNAMSIZ-1 {
		return nil, 0, EINVAL
	}

	*(*uint16)(unsafe.Pointer(&sa.raw[0])) = AF_PPPOX
	
	
	
	
	
	
	
	
	binary.BigEndian.PutUint32(sa.raw[2:6], px_proto_oe)
	
	
	binary.BigEndian.PutUint16(sa.raw[6:8], sa.SID)
	copy(sa.raw[8:14], sa.Remote)
	for i := 14; i < 14+IFNAMSIZ; i++ {
		sa.raw[i] = 0
	}
	copy(sa.raw[14:], sa.Dev)
	return unsafe.Pointer(&sa.raw), SizeofSockaddrPPPoX, nil
}



type SockaddrTIPC struct {
	
	
	Scope int

	
	
	
	
	
	
	
	Addr TIPCAddr

	raw RawSockaddrTIPC
}




type TIPCAddr interface {
	tipcAddrtype() uint8
	tipcAddr() [12]byte
}

func (sa *TIPCSocketAddr) tipcAddr() [12]byte {
	var out [12]byte
	copy(out[:], (*(*[unsafe.Sizeof(TIPCSocketAddr{})]byte)(unsafe.Pointer(sa)))[:])
	return out
}

func (sa *TIPCSocketAddr) tipcAddrtype() uint8 { return TIPC_SOCKET_ADDR }

func (sa *TIPCServiceRange) tipcAddr() [12]byte {
	var out [12]byte
	copy(out[:], (*(*[unsafe.Sizeof(TIPCServiceRange{})]byte)(unsafe.Pointer(sa)))[:])
	return out
}

func (sa *TIPCServiceRange) tipcAddrtype() uint8 { return TIPC_SERVICE_RANGE }

func (sa *TIPCServiceName) tipcAddr() [12]byte {
	var out [12]byte
	copy(out[:], (*(*[unsafe.Sizeof(TIPCServiceName{})]byte)(unsafe.Pointer(sa)))[:])
	return out
}

func (sa *TIPCServiceName) tipcAddrtype() uint8 { return TIPC_SERVICE_ADDR }

func (sa *SockaddrTIPC) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Addr == nil {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_TIPC
	sa.raw.Scope = int8(sa.Scope)
	sa.raw.Addrtype = sa.Addr.tipcAddrtype()
	sa.raw.Addr = sa.Addr.tipcAddr()
	return unsafe.Pointer(&sa.raw), SizeofSockaddrTIPC, nil
}


type SockaddrL2TPIP struct {
	Addr   [4]byte
	ConnId uint32
	raw    RawSockaddrL2TPIP
}

func (sa *SockaddrL2TPIP) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_INET
	sa.raw.Conn_id = sa.ConnId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrL2TPIP, nil
}


type SockaddrL2TPIP6 struct {
	Addr   [16]byte
	ZoneId uint32
	ConnId uint32
	raw    RawSockaddrL2TPIP6
}

func (sa *SockaddrL2TPIP6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_INET6
	sa.raw.Conn_id = sa.ConnId
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrL2TPIP6, nil
}


type SockaddrIUCV struct {
	UserID string
	Name   string
	raw    RawSockaddrIUCV
}

func (sa *SockaddrIUCV) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Family = AF_IUCV
	
	
	
	for i := 0; i < 8; i++ {
		sa.raw.Nodeid[i] = ' '
		sa.raw.User_id[i] = ' '
		sa.raw.Name[i] = ' '
	}
	if len(sa.UserID) > 8 || len(sa.Name) > 8 {
		return nil, 0, EINVAL
	}
	for i, b := range []byte(sa.UserID[:]) {
		sa.raw.User_id[i] = int8(b)
	}
	for i, b := range []byte(sa.Name[:]) {
		sa.raw.Name[i] = int8(b)
	}
	return unsafe.Pointer(&sa.raw), SizeofSockaddrIUCV, nil
}

type SockaddrNFC struct {
	DeviceIdx   uint32
	TargetIdx   uint32
	NFCProtocol uint32
	raw         RawSockaddrNFC
}

func (sa *SockaddrNFC) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Sa_family = AF_NFC
	sa.raw.Dev_idx = sa.DeviceIdx
	sa.raw.Target_idx = sa.TargetIdx
	sa.raw.Nfc_protocol = sa.NFCProtocol
	return unsafe.Pointer(&sa.raw), SizeofSockaddrNFC, nil
}

type SockaddrNFCLLCP struct {
	DeviceIdx      uint32
	TargetIdx      uint32
	NFCProtocol    uint32
	DestinationSAP uint8
	SourceSAP      uint8
	ServiceName    string
	raw            RawSockaddrNFCLLCP
}

func (sa *SockaddrNFCLLCP) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Sa_family = AF_NFC
	sa.raw.Dev_idx = sa.DeviceIdx
	sa.raw.Target_idx = sa.TargetIdx
	sa.raw.Nfc_protocol = sa.NFCProtocol
	sa.raw.Dsap = sa.DestinationSAP
	sa.raw.Ssap = sa.SourceSAP
	if len(sa.ServiceName) > len(sa.raw.Service_name) {
		return nil, 0, EINVAL
	}
	copy(sa.raw.Service_name[:], sa.ServiceName)
	sa.raw.SetServiceNameLen(len(sa.ServiceName))
	return unsafe.Pointer(&sa.raw), SizeofSockaddrNFCLLCP, nil
}

var socketProtocol = func(fd int) (int, error) {
	return GetsockoptInt(fd, SOL_SOCKET, SO_PROTOCOL)
}

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_NETLINK:
		pp := (*RawSockaddrNetlink)(unsafe.Pointer(rsa))
		sa := new(SockaddrNetlink)
		sa.Family = pp.Family
		sa.Pad = pp.Pad
		sa.Pid = pp.Pid
		sa.Groups = pp.Groups
		return sa, nil

	case AF_PACKET:
		pp := (*RawSockaddrLinklayer)(unsafe.Pointer(rsa))
		sa := new(SockaddrLinklayer)
		sa.Protocol = pp.Protocol
		sa.Ifindex = int(pp.Ifindex)
		sa.Hatype = pp.Hatype
		sa.Pkttype = pp.Pkttype
		sa.Halen = pp.Halen
		sa.Addr = pp.Addr
		return sa, nil

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		if pp.Path[0] == 0 {
			
			
			
			
			
			pp.Path[0] = '@'
		}

		
		
		
		
		
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		sa.Name = string(unsafe.Slice((*byte)(unsafe.Pointer(&pp.Path[0])), n))
		return sa, nil

	case AF_INET:
		proto, err := socketProtocol(fd)
		if err != nil {
			return nil, err
		}

		switch proto {
		case IPPROTO_L2TP:
			pp := (*RawSockaddrL2TPIP)(unsafe.Pointer(rsa))
			sa := new(SockaddrL2TPIP)
			sa.ConnId = pp.Conn_id
			sa.Addr = pp.Addr
			return sa, nil
		default:
			pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
			sa := new(SockaddrInet4)
			p := (*[2]byte)(unsafe.Pointer(&pp.Port))
			sa.Port = int(p[0])<<8 + int(p[1])
			sa.Addr = pp.Addr
			return sa, nil
		}

	case AF_INET6:
		proto, err := socketProtocol(fd)
		if err != nil {
			return nil, err
		}

		switch proto {
		case IPPROTO_L2TP:
			pp := (*RawSockaddrL2TPIP6)(unsafe.Pointer(rsa))
			sa := new(SockaddrL2TPIP6)
			sa.ConnId = pp.Conn_id
			sa.ZoneId = pp.Scope_id
			sa.Addr = pp.Addr
			return sa, nil
		default:
			pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
			sa := new(SockaddrInet6)
			p := (*[2]byte)(unsafe.Pointer(&pp.Port))
			sa.Port = int(p[0])<<8 + int(p[1])
			sa.ZoneId = pp.Scope_id
			sa.Addr = pp.Addr
			return sa, nil
		}

	case AF_VSOCK:
		pp := (*RawSockaddrVM)(unsafe.Pointer(rsa))
		sa := &SockaddrVM{
			CID:   pp.Cid,
			Port:  pp.Port,
			Flags: pp.Flags,
		}
		return sa, nil
	case AF_BLUETOOTH:
		proto, err := socketProtocol(fd)
		if err != nil {
			return nil, err
		}
		
		switch proto {
		case BTPROTO_L2CAP:
			pp := (*RawSockaddrL2)(unsafe.Pointer(rsa))
			sa := &SockaddrL2{
				PSM:      pp.Psm,
				CID:      pp.Cid,
				Addr:     pp.Bdaddr,
				AddrType: pp.Bdaddr_type,
			}
			return sa, nil
		case BTPROTO_RFCOMM:
			pp := (*RawSockaddrRFCOMM)(unsafe.Pointer(rsa))
			sa := &SockaddrRFCOMM{
				Channel: pp.Channel,
				Addr:    pp.Bdaddr,
			}
			return sa, nil
		}
	case AF_XDP:
		pp := (*RawSockaddrXDP)(unsafe.Pointer(rsa))
		sa := &SockaddrXDP{
			Flags:        pp.Flags,
			Ifindex:      pp.Ifindex,
			QueueID:      pp.Queue_id,
			SharedUmemFD: pp.Shared_umem_fd,
		}
		return sa, nil
	case AF_PPPOX:
		pp := (*RawSockaddrPPPoX)(unsafe.Pointer(rsa))
		if binary.BigEndian.Uint32(pp[2:6]) != px_proto_oe {
			return nil, EINVAL
		}
		sa := &SockaddrPPPoE{
			SID:    binary.BigEndian.Uint16(pp[6:8]),
			Remote: pp[8:14],
		}
		for i := 14; i < 14+IFNAMSIZ; i++ {
			if pp[i] == 0 {
				sa.Dev = string(pp[14:i])
				break
			}
		}
		return sa, nil
	case AF_TIPC:
		pp := (*RawSockaddrTIPC)(unsafe.Pointer(rsa))

		sa := &SockaddrTIPC{
			Scope: int(pp.Scope),
		}

		
		
		switch pp.Addrtype {
		case TIPC_SERVICE_RANGE:
			sa.Addr = (*TIPCServiceRange)(unsafe.Pointer(&pp.Addr))
		case TIPC_SERVICE_ADDR:
			sa.Addr = (*TIPCServiceName)(unsafe.Pointer(&pp.Addr))
		case TIPC_SOCKET_ADDR:
			sa.Addr = (*TIPCSocketAddr)(unsafe.Pointer(&pp.Addr))
		default:
			return nil, EINVAL
		}

		return sa, nil
	case AF_IUCV:
		pp := (*RawSockaddrIUCV)(unsafe.Pointer(rsa))

		var user [8]byte
		var name [8]byte

		for i := 0; i < 8; i++ {
			user[i] = byte(pp.User_id[i])
			name[i] = byte(pp.Name[i])
		}

		sa := &SockaddrIUCV{
			UserID: string(user[:]),
			Name:   string(name[:]),
		}
		return sa, nil

	case AF_CAN:
		proto, err := socketProtocol(fd)
		if err != nil {
			return nil, err
		}

		pp := (*RawSockaddrCAN)(unsafe.Pointer(rsa))

		switch proto {
		case CAN_J1939:
			sa := &SockaddrCANJ1939{
				Ifindex: int(pp.Ifindex),
			}
			name := (*[8]byte)(unsafe.Pointer(&sa.Name))
			for i := 0; i < 8; i++ {
				name[i] = pp.Addr[i]
			}
			pgn := (*[4]byte)(unsafe.Pointer(&sa.PGN))
			for i := 0; i < 4; i++ {
				pgn[i] = pp.Addr[i+8]
			}
			addr := (*[1]byte)(unsafe.Pointer(&sa.Addr))
			addr[0] = pp.Addr[12]
			return sa, nil
		default:
			sa := &SockaddrCAN{
				Ifindex: int(pp.Ifindex),
			}
			rx := (*[4]byte)(unsafe.Pointer(&sa.RxID))
			for i := 0; i < 4; i++ {
				rx[i] = pp.Addr[i]
			}
			tx := (*[4]byte)(unsafe.Pointer(&sa.TxID))
			for i := 0; i < 4; i++ {
				tx[i] = pp.Addr[i+4]
			}
			return sa, nil
		}
	case AF_NFC:
		proto, err := socketProtocol(fd)
		if err != nil {
			return nil, err
		}
		switch proto {
		case NFC_SOCKPROTO_RAW:
			pp := (*RawSockaddrNFC)(unsafe.Pointer(rsa))
			sa := &SockaddrNFC{
				DeviceIdx:   pp.Dev_idx,
				TargetIdx:   pp.Target_idx,
				NFCProtocol: pp.Nfc_protocol,
			}
			return sa, nil
		case NFC_SOCKPROTO_LLCP:
			pp := (*RawSockaddrNFCLLCP)(unsafe.Pointer(rsa))
			if uint64(pp.Service_name_len) > uint64(len(pp.Service_name)) {
				return nil, EINVAL
			}
			sa := &SockaddrNFCLLCP{
				DeviceIdx:      pp.Dev_idx,
				TargetIdx:      pp.Target_idx,
				NFCProtocol:    pp.Nfc_protocol,
				DestinationSAP: pp.Dsap,
				SourceSAP:      pp.Ssap,
				ServiceName:    string(pp.Service_name[:pp.Service_name_len]),
			}
			return sa, nil
		default:
			return nil, EINVAL
		}
	}
	return nil, EAFNOSUPPORT
}

func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, 0)
	if err != nil {
		return
	}
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func Accept4(fd int, flags int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, flags)
	if err != nil {
		return
	}
	if len > SizeofSockaddrAny {
		panic("RawSockaddrAny too small")
	}
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(fd, &rsa)
}

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error) {
	var value IPMreqn
	vallen := _Socklen(SizeofIPMreqn)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptUcred(fd, level, opt int) (*Ucred, error) {
	var value Ucred
	vallen := _Socklen(SizeofUcred)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptTCPInfo(fd, level, opt int) (*TCPInfo, error) {
	var value TCPInfo
	vallen := _Socklen(SizeofTCPInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}







func GetsockoptTCPCCVegasInfo(fd, level, opt int) (*TCPVegasInfo, error) {
	var value [SizeofTCPCCInfo / 4]uint32 
	vallen := _Socklen(SizeofTCPCCInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	out := (*TCPVegasInfo)(unsafe.Pointer(&value[0]))
	return out, err
}







func GetsockoptTCPCCDCTCPInfo(fd, level, opt int) (*TCPDCTCPInfo, error) {
	var value [SizeofTCPCCInfo / 4]uint32 
	vallen := _Socklen(SizeofTCPCCInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	out := (*TCPDCTCPInfo)(unsafe.Pointer(&value[0]))
	return out, err
}







func GetsockoptTCPCCBBRInfo(fd, level, opt int) (*TCPBBRInfo, error) {
	var value [SizeofTCPCCInfo / 4]uint32 
	vallen := _Socklen(SizeofTCPCCInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	out := (*TCPBBRInfo)(unsafe.Pointer(&value[0]))
	return out, err
}



func GetsockoptString(fd, level, opt int) (string, error) {
	buf := make([]byte, 256)
	vallen := _Socklen(len(buf))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil {
		if err == ERANGE {
			buf = make([]byte, vallen)
			err = getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
		}
		if err != nil {
			return "", err
		}
	}
	return ByteSliceToString(buf[:vallen]), nil
}

func GetsockoptTpacketStats(fd, level, opt int) (*TpacketStats, error) {
	var value TpacketStats
	vallen := _Socklen(SizeofTpacketStats)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptTpacketStatsV3(fd, level, opt int) (*TpacketStatsV3, error) {
	var value TpacketStatsV3
	vallen := _Socklen(SizeofTpacketStatsV3)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
}

func SetsockoptPacketMreq(fd, level, opt int, mreq *PacketMreq) error {
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
}



func SetsockoptSockFprog(fd, level, opt int, fprog *SockFprog) error {
	return setsockopt(fd, level, opt, unsafe.Pointer(fprog), unsafe.Sizeof(*fprog))
}

func SetsockoptCanRawFilter(fd, level, opt int, filter []CanFilter) error {
	var p unsafe.Pointer
	if len(filter) > 0 {
		p = unsafe.Pointer(&filter[0])
	}
	return setsockopt(fd, level, opt, p, uintptr(len(filter)*SizeofCanFilter))
}

func SetsockoptTpacketReq(fd, level, opt int, tp *TpacketReq) error {
	return setsockopt(fd, level, opt, unsafe.Pointer(tp), unsafe.Sizeof(*tp))
}

func SetsockoptTpacketReq3(fd, level, opt int, tp *TpacketReq3) error {
	return setsockopt(fd, level, opt, unsafe.Pointer(tp), unsafe.Sizeof(*tp))
}

func SetsockoptTCPRepairOpt(fd, level, opt int, o []TCPRepairOpt) (err error) {
	if len(o) == 0 {
		return EINVAL
	}
	return setsockopt(fd, level, opt, unsafe.Pointer(&o[0]), uintptr(SizeofTCPRepairOpt*len(o)))
}

func SetsockoptTCPMD5Sig(fd, level, opt int, s *TCPMD5Sig) error {
	return setsockopt(fd, level, opt, unsafe.Pointer(s), unsafe.Sizeof(*s))
}

















func KeyctlString(cmd int, id int) (string, error) {
	
	
	
	
	var buffer []byte
	for {
		
		length, err := KeyctlBuffer(cmd, id, buffer, 0)
		if err != nil {
			return "", err
		}

		
		if length <= len(buffer) {
			
			return string(buffer[:length-1]), nil
		}

		
		buffer = make([]byte, length)
	}
}






func KeyctlGetKeyringID(id int, create bool) (ringid int, err error) {
	createInt := 0
	if create {
		createInt = 1
	}
	return KeyctlInt(KEYCTL_GET_KEYRING_ID, id, createInt, 0, 0)
}






func KeyctlSetperm(id int, perm uint32) error {
	_, err := KeyctlInt(KEYCTL_SETPERM, id, int(perm), 0, 0)
	return err
}






func KeyctlJoinSessionKeyring(name string) (ringid int, err error) {
	return keyctlJoin(KEYCTL_JOIN_SESSION_KEYRING, name)
}






func KeyctlSearch(ringid int, keyType, description string, destRingid int) (id int, err error) {
	return keyctlSearch(KEYCTL_SEARCH, ringid, keyType, description, destRingid)
}








func KeyctlInstantiateIOV(id int, payload []Iovec, ringid int) error {
	return keyctlIOV(KEYCTL_INSTANTIATE_IOV, id, payload, ringid)
}












func KeyctlDHCompute(params *KeyctlDHParams, buffer []byte) (size int, err error) {
	return keyctlDH(KEYCTL_DH_COMPUTE, params, buffer)
}



















func KeyctlRestrictKeyring(ringid int, keyType string, restriction string) error {
	if keyType == "" {
		return keyctlRestrictKeyring(KEYCTL_RESTRICT_KEYRING, ringid)
	}
	return keyctlRestrictKeyringByType(KEYCTL_RESTRICT_KEYRING, ringid, keyType, restriction)
}




func recvmsgRaw(fd int, iov []Iovec, oob []byte, flags int, rsa *RawSockaddrAny) (n, oobn int, recvflags int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var dummy byte
	if len(oob) > 0 {
		if emptyIovecs(iov) {
			var sockType int
			sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
			if err != nil {
				return
			}
			
			if sockType != SOCK_DGRAM {
				var iova [1]Iovec
				iova[0].Base = &dummy
				iova[0].SetLen(1)
				iov = iova[:]
			}
		}
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	}
	if len(iov) > 0 {
		msg.Iov = &iov[0]
		msg.SetIovlen(len(iov))
	}
	if n, err = recvmsg(fd, &msg, flags); err != nil {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	return
}

func sendmsgN(fd int, iov []Iovec, oob []byte, ptr unsafe.Pointer, salen _Socklen, flags int) (n int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(ptr)
	msg.Namelen = uint32(salen)
	var dummy byte
	var empty bool
	if len(oob) > 0 {
		empty = emptyIovecs(iov)
		if empty {
			var sockType int
			sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
			if err != nil {
				return 0, err
			}
			
			if sockType != SOCK_DGRAM {
				var iova [1]Iovec
				iova[0].Base = &dummy
				iova[0].SetLen(1)
				iov = iova[:]
			}
		}
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	}
	if len(iov) > 0 {
		msg.Iov = &iov[0]
		msg.SetIovlen(len(iov))
	}
	if n, err = sendmsg(fd, &msg, flags); err != nil {
		return 0, err
	}
	if len(oob) > 0 && empty {
		n = 0
	}
	return n, nil
}


func BindToDevice(fd int, device string) (err error) {
	return SetsockoptString(fd, SOL_SOCKET, SO_BINDTODEVICE, device)
}




func ptracePeek(req int, pid int, addr uintptr, out []byte) (count int, err error) {
	
	

	
	

	var buf [SizeofPtr]byte

	
	
	
	
	
	n := 0
	if addr%SizeofPtr != 0 {
		err = ptracePtr(req, pid, addr-addr%SizeofPtr, unsafe.Pointer(&buf[0]))
		if err != nil {
			return 0, err
		}
		n += copy(out, buf[addr%SizeofPtr:])
		out = out[n:]
	}

	
	for len(out) > 0 {
		
		
		err = ptracePtr(req, pid, addr+uintptr(n), unsafe.Pointer(&buf[0]))
		if err != nil {
			return n, err
		}
		copied := copy(out, buf[0:])
		n += copied
		out = out[copied:]
	}

	return n, nil
}

func PtracePeekText(pid int, addr uintptr, out []byte) (count int, err error) {
	return ptracePeek(PTRACE_PEEKTEXT, pid, addr, out)
}

func PtracePeekData(pid int, addr uintptr, out []byte) (count int, err error) {
	return ptracePeek(PTRACE_PEEKDATA, pid, addr, out)
}

func PtracePeekUser(pid int, addr uintptr, out []byte) (count int, err error) {
	return ptracePeek(PTRACE_PEEKUSR, pid, addr, out)
}

func ptracePoke(pokeReq int, peekReq int, pid int, addr uintptr, data []byte) (count int, err error) {
	
	

	
	n := 0
	if addr%SizeofPtr != 0 {
		var buf [SizeofPtr]byte
		err = ptracePtr(peekReq, pid, addr-addr%SizeofPtr, unsafe.Pointer(&buf[0]))
		if err != nil {
			return 0, err
		}
		n += copy(buf[addr%SizeofPtr:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr-addr%SizeofPtr, word)
		if err != nil {
			return 0, err
		}
		data = data[n:]
	}

	
	for len(data) > SizeofPtr {
		word := *((*uintptr)(unsafe.Pointer(&data[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if err != nil {
			return n, err
		}
		n += SizeofPtr
		data = data[SizeofPtr:]
	}

	
	if len(data) > 0 {
		var buf [SizeofPtr]byte
		err = ptracePtr(peekReq, pid, addr+uintptr(n), unsafe.Pointer(&buf[0]))
		if err != nil {
			return n, err
		}
		copy(buf[0:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if err != nil {
			return n, err
		}
		n += len(data)
	}

	return n, nil
}

func PtracePokeText(pid int, addr uintptr, data []byte) (count int, err error) {
	return ptracePoke(PTRACE_POKETEXT, PTRACE_PEEKTEXT, pid, addr, data)
}

func PtracePokeData(pid int, addr uintptr, data []byte) (count int, err error) {
	return ptracePoke(PTRACE_POKEDATA, PTRACE_PEEKDATA, pid, addr, data)
}

func PtracePokeUser(pid int, addr uintptr, data []byte) (count int, err error) {
	return ptracePoke(PTRACE_POKEUSR, PTRACE_PEEKUSR, pid, addr, data)
}




const elfNT_PRSTATUS = 1

func PtraceGetRegs(pid int, regsout *PtraceRegs) (err error) {
	var iov Iovec
	iov.Base = (*byte)(unsafe.Pointer(regsout))
	iov.SetLen(int(unsafe.Sizeof(*regsout)))
	return ptracePtr(PTRACE_GETREGSET, pid, uintptr(elfNT_PRSTATUS), unsafe.Pointer(&iov))
}

func PtraceSetRegs(pid int, regs *PtraceRegs) (err error) {
	var iov Iovec
	iov.Base = (*byte)(unsafe.Pointer(regs))
	iov.SetLen(int(unsafe.Sizeof(*regs)))
	return ptracePtr(PTRACE_SETREGSET, pid, uintptr(elfNT_PRSTATUS), unsafe.Pointer(&iov))
}

func PtraceSetOptions(pid int, options int) (err error) {
	return ptrace(PTRACE_SETOPTIONS, pid, 0, uintptr(options))
}

func PtraceGetEventMsg(pid int) (msg uint, err error) {
	var data _C_long
	err = ptracePtr(PTRACE_GETEVENTMSG, pid, 0, unsafe.Pointer(&data))
	msg = uint(data)
	return
}

func PtraceCont(pid int, signal int) (err error) {
	return ptrace(PTRACE_CONT, pid, 0, uintptr(signal))
}

func PtraceSyscall(pid int, signal int) (err error) {
	return ptrace(PTRACE_SYSCALL, pid, 0, uintptr(signal))
}

func PtraceSingleStep(pid int) (err error) { return ptrace(PTRACE_SINGLESTEP, pid, 0, 0) }

func PtraceInterrupt(pid int) (err error) { return ptrace(PTRACE_INTERRUPT, pid, 0, 0) }

func PtraceAttach(pid int) (err error) { return ptrace(PTRACE_ATTACH, pid, 0, 0) }

func PtraceSeize(pid int) (err error) { return ptrace(PTRACE_SEIZE, pid, 0, 0) }

func PtraceDetach(pid int) (err error) { return ptrace(PTRACE_DETACH, pid, 0, 0) }



func Reboot(cmd int) (err error) {
	return reboot(LINUX_REBOOT_MAGIC1, LINUX_REBOOT_MAGIC2, cmd, "")
}

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Ino), unsafe.Sizeof(Dirent{}.Ino))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Reclen), unsafe.Sizeof(Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	reclen, ok := direntReclen(buf)
	if !ok {
		return 0, false
	}
	return reclen - uint64(unsafe.Offsetof(Dirent{}.Name)), true
}



func Mount(source string, target string, fstype string, flags uintptr, data string) (err error) {
	
	
	if data == "" {
		return mount(source, target, fstype, flags, nil)
	}
	datap, err := BytePtrFromString(data)
	if err != nil {
		return err
	}
	return mount(source, target, fstype, flags, datap)
}







func MountSetattr(dirfd int, pathname string, flags uint, attr *MountAttr) error {
	return mountSetattr(dirfd, pathname, flags, attr, unsafe.Sizeof(*attr))
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	return sendfile(outfd, infd, offset, count)
}
























func Dup2(oldfd, newfd int) error {
	return Dup3(oldfd, newfd, 0)
}
























func fsconfigCommon(fd int, cmd uint, key string, value *byte, aux int) (err error) {
	var keyp *byte
	if keyp, err = BytePtrFromString(key); err != nil {
		return
	}
	return fsconfig(fd, cmd, keyp, value, aux)
}






func FsconfigSetFlag(fd int, key string) (err error) {
	return fsconfigCommon(fd, FSCONFIG_SET_FLAG, key, nil, 0)
}







func FsconfigSetString(fd int, key string, value string) (err error) {
	var valuep *byte
	if valuep, err = BytePtrFromString(value); err != nil {
		return
	}
	return fsconfigCommon(fd, FSCONFIG_SET_STRING, key, valuep, 0)
}







func FsconfigSetBinary(fd int, key string, value []byte) (err error) {
	if len(value) == 0 {
		return EINVAL
	}
	return fsconfigCommon(fd, FSCONFIG_SET_BINARY, key, &value[0], len(value))
}








func FsconfigSetPath(fd int, key string, path string, atfd int) (err error) {
	var valuep *byte
	if valuep, err = BytePtrFromString(path); err != nil {
		return
	}
	return fsconfigCommon(fd, FSCONFIG_SET_PATH, key, valuep, atfd)
}




func FsconfigSetPathEmpty(fd int, key string, path string, atfd int) (err error) {
	var valuep *byte
	if valuep, err = BytePtrFromString(path); err != nil {
		return
	}
	return fsconfigCommon(fd, FSCONFIG_SET_PATH_EMPTY, key, valuep, atfd)
}







func FsconfigSetFd(fd int, key string, value int) (err error) {
	return fsconfigCommon(fd, FSCONFIG_SET_FD, key, nil, value)
}





func FsconfigCreate(fd int) (err error) {
	return fsconfig(fd, FSCONFIG_CMD_CREATE, nil, nil, 0)
}





func FsconfigReconfigure(fd int) (err error) {
	return fsconfig(fd, FSCONFIG_CMD_RECONFIGURE, nil, nil, 0)
}




func Getpgrp() (pid int) {
	pid, _ = Getpgid(0)
	return
}





func Getrandom(buf []byte, flags int) (n int, err error) {
	vdsoRet, supported := vgetrandom(buf, uint32(flags))
	if supported {
		if vdsoRet < 0 {
			return 0, errnoErr(syscall.Errno(-vdsoRet))
		}
		return vdsoRet, nil
	}
	var p *byte
	if len(buf) > 0 {
		p = &buf[0]
	}
	r, _, e := Syscall(SYS_GETRANDOM, uintptr(unsafe.Pointer(p)), uintptr(len(buf)), uintptr(flags))
	if e != 0 {
		return 0, errnoErr(e)
	}
	return int(r), nil
}





































//go:linkname syscall_prlimit syscall.prlimit
func syscall_prlimit(pid, resource int, newlimit, old *syscall.Rlimit) error

func Prlimit(pid, resource int, newlimit, old *Rlimit) error {
	
	
	return syscall_prlimit(pid, resource, (*syscall.Rlimit)(newlimit), (*syscall.Rlimit)(old))
}




func PrctlRetInt(option int, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr) (int, error) {
	ret, _, err := Syscall6(SYS_PRCTL, uintptr(option), uintptr(arg2), uintptr(arg3), uintptr(arg4), uintptr(arg5), 0)
	if err != 0 {
		return 0, err
	}
	return int(ret), nil
}

func Setuid(uid int) (err error) {
	return syscall.Setuid(uid)
}

func Setgid(gid int) (err error) {
	return syscall.Setgid(gid)
}

func Setreuid(ruid, euid int) (err error) {
	return syscall.Setreuid(ruid, euid)
}

func Setregid(rgid, egid int) (err error) {
	return syscall.Setregid(rgid, egid)
}

func Setresuid(ruid, euid, suid int) (err error) {
	return syscall.Setresuid(ruid, euid, suid)
}

func Setresgid(rgid, egid, sgid int) (err error) {
	return syscall.Setresgid(rgid, egid, sgid)
}




func SetfsgidRetGid(gid int) (int, error) {
	return setfsgid(gid)
}




func SetfsuidRetUid(uid int) (int, error) {
	return setfsuid(uid)
}

func Setfsgid(gid int) error {
	_, err := setfsgid(gid)
	return err
}

func Setfsuid(uid int) error {
	_, err := setfsuid(uid)
	return err
}

func Signalfd(fd int, sigmask *Sigset_t, flags int) (newfd int, err error) {
	return signalfd(fd, sigmask, _C__NSIG/8, flags)
}

































const minIovec = 8


func appendBytes(vecs []Iovec, bs [][]byte) []Iovec {
	for _, b := range bs {
		var v Iovec
		v.SetLen(len(b))
		if len(b) > 0 {
			v.Base = &b[0]
		} else {
			v.Base = (*byte)(unsafe.Pointer(&_zero))
		}
		vecs = append(vecs, v)
	}
	return vecs
}


func offs2lohi(offs int64) (lo, hi uintptr) {
	const longBits = SizeofLong * 8
	return uintptr(offs), uintptr(uint64(offs) >> (longBits - 1) >> 1) 
}

func Readv(fd int, iovs [][]byte) (n int, err error) {
	iovecs := make([]Iovec, 0, minIovec)
	iovecs = appendBytes(iovecs, iovs)
	n, err = readv(fd, iovecs)
	readvRacedetect(iovecs, n, err)
	return n, err
}

func Preadv(fd int, iovs [][]byte, offset int64) (n int, err error) {
	iovecs := make([]Iovec, 0, minIovec)
	iovecs = appendBytes(iovecs, iovs)
	lo, hi := offs2lohi(offset)
	n, err = preadv(fd, iovecs, lo, hi)
	readvRacedetect(iovecs, n, err)
	return n, err
}

func Preadv2(fd int, iovs [][]byte, offset int64, flags int) (n int, err error) {
	iovecs := make([]Iovec, 0, minIovec)
	iovecs = appendBytes(iovecs, iovs)
	lo, hi := offs2lohi(offset)
	n, err = preadv2(fd, iovecs, lo, hi, flags)
	readvRacedetect(iovecs, n, err)
	return n, err
}

func readvRacedetect(iovecs []Iovec, n int, err error) {
	if !raceenabled {
		return
	}
	for i := 0; n > 0 && i < len(iovecs); i++ {
		m := int(iovecs[i].Len)
		if m > n {
			m = n
		}
		n -= m
		if m > 0 {
			raceWriteRange(unsafe.Pointer(iovecs[i].Base), m)
		}
	}
	if err == nil {
		raceAcquire(unsafe.Pointer(&ioSync))
	}
}

func Writev(fd int, iovs [][]byte) (n int, err error) {
	iovecs := make([]Iovec, 0, minIovec)
	iovecs = appendBytes(iovecs, iovs)
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	n, err = writev(fd, iovecs)
	writevRacedetect(iovecs, n)
	return n, err
}

func Pwritev(fd int, iovs [][]byte, offset int64) (n int, err error) {
	iovecs := make([]Iovec, 0, minIovec)
	iovecs = appendBytes(iovecs, iovs)
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	lo, hi := offs2lohi(offset)
	n, err = pwritev(fd, iovecs, lo, hi)
	writevRacedetect(iovecs, n)
	return n, err
}

func Pwritev2(fd int, iovs [][]byte, offset int64, flags int) (n int, err error) {
	iovecs := make([]Iovec, 0, minIovec)
	iovecs = appendBytes(iovecs, iovs)
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	lo, hi := offs2lohi(offset)
	n, err = pwritev2(fd, iovecs, lo, hi, flags)
	writevRacedetect(iovecs, n)
	return n, err
}

func writevRacedetect(iovecs []Iovec, n int) {
	if !raceenabled {
		return
	}
	for i := 0; n > 0 && i < len(iovecs); i++ {
		m := int(iovecs[i].Len)
		if m > n {
			m = n
		}
		n -= m
		if m > 0 {
			raceReadRange(unsafe.Pointer(iovecs[i].Base), m)
		}
	}
}












const (
	mremapFixed     = MREMAP_FIXED
	mremapDontunmap = MREMAP_DONTUNMAP
	mremapMaymove   = MREMAP_MAYMOVE
)



func Vmsplice(fd int, iovs []Iovec, flags int) (int, error) {
	var p unsafe.Pointer
	if len(iovs) > 0 {
		p = unsafe.Pointer(&iovs[0])
	}

	n, _, errno := Syscall6(SYS_VMSPLICE, uintptr(fd), uintptr(p), uintptr(len(iovs)), uintptr(flags), 0, 0)
	if errno != 0 {
		return 0, syscall.Errno(errno)
	}

	return int(n), nil
}

func isGroupMember(gid int) bool {
	groups, err := Getgroups()
	if err != nil {
		return false
	}

	for _, g := range groups {
		if g == gid {
			return true
		}
	}
	return false
}

func isCapDacOverrideSet() bool {
	hdr := CapUserHeader{Version: LINUX_CAPABILITY_VERSION_3}
	data := [2]CapUserData{}
	err := Capget(&hdr, &data[0])

	return err == nil && data[0].Effective&(1<<CAP_DAC_OVERRIDE) != 0
}




func Faccessat(dirfd int, path string, mode uint32, flags int) (err error) {
	if flags == 0 {
		return faccessat(dirfd, path, mode)
	}

	if err := Faccessat2(dirfd, path, mode, flags); err != ENOSYS && err != EPERM {
		return err
	}

	
	
	
	
	

	if flags & ^(AT_SYMLINK_NOFOLLOW|AT_EACCESS) != 0 {
		return EINVAL
	}

	var st Stat_t
	if err := Fstatat(dirfd, path, &st, flags&AT_SYMLINK_NOFOLLOW); err != nil {
		return err
	}

	mode &= 7
	if mode == 0 {
		return nil
	}

	var uid int
	if flags&AT_EACCESS != 0 {
		uid = Geteuid()
		if uid != 0 && isCapDacOverrideSet() {
			
			
			
			uid = 0
		}
	} else {
		uid = Getuid()
	}

	if uid == 0 {
		if mode&1 == 0 {
			
			return nil
		}
		if st.Mode&0111 != 0 {
			
			return nil
		}
		return EACCES
	}

	var fmode uint32
	if uint32(uid) == st.Uid {
		fmode = (st.Mode >> 6) & 7
	} else {
		var gid int
		if flags&AT_EACCESS != 0 {
			gid = Getegid()
		} else {
			gid = Getgid()
		}

		if uint32(gid) == st.Gid || isGroupMember(int(st.Gid)) {
			fmode = (st.Mode >> 3) & 7
		} else {
			fmode = st.Mode & 7
		}
	}

	if fmode&mode == mode {
		return nil
	}

	return EACCES
}









type fileHandle struct {
	Bytes uint32
	Type  int32
}




type FileHandle struct {
	*fileHandle
}


func NewFileHandle(handleType int32, handle []byte) FileHandle {
	const hdrSize = unsafe.Sizeof(fileHandle{})
	buf := make([]byte, hdrSize+uintptr(len(handle)))
	copy(buf[hdrSize:], handle)
	fh := (*fileHandle)(unsafe.Pointer(&buf[0]))
	fh.Type = handleType
	fh.Bytes = uint32(len(handle))
	return FileHandle{fh}
}

func (fh *FileHandle) Size() int   { return int(fh.fileHandle.Bytes) }
func (fh *FileHandle) Type() int32 { return fh.fileHandle.Type }
func (fh *FileHandle) Bytes() []byte {
	n := fh.Size()
	if n == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&fh.fileHandle.Type))+4)), n)
}



func NameToHandleAt(dirfd int, path string, flags int) (handle FileHandle, mountID int, err error) {
	var mid _C_int
	
	
	size := uint32(32 + unsafe.Sizeof(fileHandle{}))
	didResize := false
	for {
		buf := make([]byte, size)
		fh := (*fileHandle)(unsafe.Pointer(&buf[0]))
		fh.Bytes = size - uint32(unsafe.Sizeof(fileHandle{}))
		err = nameToHandleAt(dirfd, path, fh, &mid, flags)
		if err == EOVERFLOW {
			if didResize {
				
				return
			}
			didResize = true
			size = fh.Bytes + uint32(unsafe.Sizeof(fileHandle{}))
			continue
		}
		if err != nil {
			return
		}
		return FileHandle{fh}, int(mid), nil
	}
}



func OpenByHandleAt(mountFD int, handle FileHandle, flags int) (fd int, err error) {
	return openByHandleAt(mountFD, handle.fileHandle, flags)
}



func Klogset(typ int, arg int) (err error) {
	var p unsafe.Pointer
	_, _, errno := Syscall(SYS_SYSLOG, uintptr(typ), uintptr(p), uintptr(arg))
	if errno != 0 {
		return errnoErr(errno)
	}
	return nil
}





type RemoteIovec struct {
	Base uintptr
	Len  int
}

















func MakeItimerval(interval, value time.Duration) Itimerval {
	return Itimerval{
		Interval: NsecToTimeval(interval.Nanoseconds()),
		Value:    NsecToTimeval(value.Nanoseconds()),
	}
}



type ItimerWhich int


const (
	ItimerReal    ItimerWhich = ITIMER_REAL
	ItimerVirtual ItimerWhich = ITIMER_VIRTUAL
	ItimerProf    ItimerWhich = ITIMER_PROF
)



func Getitimer(which ItimerWhich) (Itimerval, error) {
	var it Itimerval
	if err := getitimer(int(which), &it); err != nil {
		return Itimerval{}, err
	}

	return it, nil
}





func Setitimer(which ItimerWhich, it Itimerval) (Itimerval, error) {
	var prev Itimerval
	if err := setitimer(int(which), &it, &prev); err != nil {
		return Itimerval{}, err
	}

	return prev, nil
}



func PthreadSigmask(how int, set, oldset *Sigset_t) error {
	if oldset != nil {
		
		*oldset = Sigset_t{}
	}
	return rtSigprocmask(how, set, oldset, _C__NSIG/8)
}




func Getresuid() (ruid, euid, suid int) {
	var r, e, s _C_int
	getresuid(&r, &e, &s)
	return int(r), int(e), int(s)
}

func Getresgid() (rgid, egid, sgid int) {
	var r, e, s _C_int
	getresgid(&r, &e, &s)
	return int(r), int(e), int(s)
}



func Pselect(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timespec, sigmask *Sigset_t) (n int, err error) {
	
	
	
	var mutableTimeout *Timespec
	if timeout != nil {
		mutableTimeout = new(Timespec)
		*mutableTimeout = *timeout
	}

	
	
	var kernelMask *sigset_argpack
	if sigmask != nil {
		wordBits := 32 << (^uintptr(0) >> 63) 

		
		
		
		sigsetWords := (_C__NSIG - 1 + wordBits - 1) / (wordBits)

		sigsetBytes := uintptr(sigsetWords * (wordBits / 8))
		kernelMask = &sigset_argpack{
			ss:    sigmask,
			ssLen: sigsetBytes,
		}
	}

	return pselect6(nfd, r, w, e, mutableTimeout, kernelMask)
}






func SchedSetAttr(pid int, attr *SchedAttr, flags uint) error {
	if attr == nil {
		return EINVAL
	}
	attr.Size = SizeofSchedAttr
	return schedSetattr(pid, attr, flags)
}



func SchedGetAttr(pid int, flags uint) (*SchedAttr, error) {
	attr := &SchedAttr{}
	if err := schedGetattr(pid, attr, SizeofSchedAttr, flags); err != nil {
		return nil, err
	}
	return attr, nil
}



