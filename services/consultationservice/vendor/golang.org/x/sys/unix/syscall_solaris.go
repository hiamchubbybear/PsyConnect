











package unix

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)


type syscallFunc uintptr

func rawSysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)


type SockaddrDatalink struct {
	Family uint16
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [244]int8
	raw    RawSockaddrDatalink
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



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	n, err := pipe(&pp)
	if n != 0 {
		return err
	}
	if err == nil {
		p[0] = int(pp[0])
		p[1] = int(pp[1])
	}
	return nil
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



func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(fd, &rsa)
}



func GetsockoptString(fd, level, opt int) (string, error) {
	buf := make([]byte, 256)
	vallen := _Socklen(len(buf))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil {
		return "", err
	}
	return ByteSliceToString(buf[:vallen]), nil
}

const ImplementsGetwd = true



func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	
	_, err = Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	n := clen(buf[:])
	if n < 1 {
		return "", EINVAL
	}
	return string(buf[:n]), nil
}






func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	
	
	if n < 0 || n > 1024 {
		if err != nil {
			return nil, err
		}
		return nil, EINVAL
	} else if n == 0 {
		return nil, nil
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if n == -1 {
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


func ReadDirent(fd int, buf []byte) (n int, err error) {
	
	
	return Getdents(fd, buf, new(uintptr))
}







type WaitStatus uint32

const (
	mask  = 0x7F
	core  = 0x80
	shift = 8

	exited  = 0
	stopped = 0x7F
)

func (w WaitStatus) Exited() bool { return w&mask == exited }

func (w WaitStatus) ExitStatus() int {
	if w&mask != exited {
		return -1
	}
	return int(w >> shift)
}

func (w WaitStatus) Signaled() bool { return w&mask != stopped && w&mask != 0 }

func (w WaitStatus) Signal() syscall.Signal {
	sig := syscall.Signal(w & mask)
	if sig == stopped || sig == 0 {
		return -1
	}
	return sig
}

func (w WaitStatus) CoreDump() bool { return w.Signaled() && w&core != 0 }

func (w WaitStatus) Stopped() bool { return w&mask == stopped && syscall.Signal(w>>shift) != SIGSTOP }

func (w WaitStatus) Continued() bool { return w&mask == stopped && syscall.Signal(w>>shift) == SIGSTOP }

func (w WaitStatus) StopSignal() syscall.Signal {
	if !w.Stopped() {
		return -1
	}
	return syscall.Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int { return -1 }



func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (int, error) {
	var status _C_int
	rpid, err := wait4(int32(pid), &status, options, rusage)
	wpid := int(rpid)
	if wpid == -1 {
		return wpid, err
	}
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return wpid, nil
}



func Gethostname() (name string, err error) {
	var buf [MaxHostNameLen]byte
	n, err := gethostname(buf[:])
	if n != 0 {
		return "", err
	}
	n = clen(buf[:])
	if n < 1 {
		return "", EFAULT
	}
	return string(buf[:n]), nil
}



func Utimes(path string, tv []Timeval) (err error) {
	if tv == nil {
		return utimes(path, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}



func UtimesNano(path string, ts []Timespec) error {
	if ts == nil {
		return utimensat(AT_FDCWD, path, nil, 0)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
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




func FcntlInt(fd uintptr, cmd, arg int) (int, error) {
	valptr, _, errno := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(arg), 0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
	}
	return int(valptr), err
}


func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(lk)), 0, 0, 0)
	if e1 != 0 {
		return e1
	}
	return nil
}



func Futimesat(dirfd int, path string, tv []Timeval) error {
	pathp, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	if tv == nil {
		return futimesat(dirfd, pathp, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return futimesat(dirfd, pathp, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}




func Futimes(fd int, tv []Timeval) error {
	if tv == nil {
		return futimesat(fd, nil, nil)
	}
	if len(tv) != 2 {
		return EINVAL
	}
	return futimesat(fd, nil, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		
		
		
		
		
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		sa.Name = string(unsafe.Slice((*byte)(unsafe.Pointer(&pp.Path[0])), n))
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.Addr = pp.Addr
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		sa.Addr = pp.Addr
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}



func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if nfd == -1 {
		return
	}
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}



func recvmsgRaw(fd int, iov []Iovec, oob []byte, flags int, rsa *RawSockaddrAny) (n, oobn int, recvflags int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var dummy byte
	if len(oob) > 0 {
		
		if emptyIovecs(iov) {
			var iova [1]Iovec
			iova[0].Base = &dummy
			iova[0].SetLen(1)
			iov = iova[:]
		}
		msg.Accrightslen = int32(len(oob))
	}
	if len(iov) > 0 {
		msg.Iov = &iov[0]
		msg.SetIovlen(len(iov))
	}
	if n, err = recvmsg(fd, &msg, flags); n == -1 {
		return
	}
	oobn = int(msg.Accrightslen)
	return
}



func sendmsgN(fd int, iov []Iovec, oob []byte, ptr unsafe.Pointer, salen _Socklen, flags int) (n int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(ptr))
	msg.Namelen = uint32(salen)
	var dummy byte
	var empty bool
	if len(oob) > 0 {
		
		empty = emptyIovecs(iov)
		if empty {
			var iova [1]Iovec
			iova[0].Base = &dummy
			iova[0].SetLen(1)
			iov = iova[:]
		}
		msg.Accrightslen = int32(len(oob))
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



func Acct(path string) (err error) {
	if len(path) == 0 {
		
		return acct(nil)
	}

	pathp, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	return acct(pathp)
}



func Mkdev(major, minor uint32) uint64 {
	return __makedev(NEWDEV, uint(major), uint(minor))
}



func Major(dev uint64) uint32 {
	return uint32(__major(NEWDEV, dev))
}



func Minor(dev uint64) uint32 {
	return uint32(__minor(NEWDEV, dev))
}






func ioctl(fd int, req int, arg uintptr) (err error) {
	_, err = ioctlRet(fd, req, arg)
	return err
}

func ioctlPtr(fd int, req int, arg unsafe.Pointer) (err error) {
	_, err = ioctlPtrRet(fd, req, arg)
	return err
}

func IoctlSetTermio(fd int, req int, value *Termio) error {
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}

func IoctlGetTermio(fd int, req int) (*Termio, error) {
	var value Termio
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return &value, err
}



func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	return sendfile(outfd, infd, offset, count)
}



















































































































type fileObjCookie struct {
	fobj   *fileObj
	cookie interface{}
}


type EventPort struct {
	port  int
	mu    sync.Mutex
	fds   map[uintptr]*fileObjCookie
	paths map[string]*fileObjCookie
	
	
	
	
	
	
	
	
	
	
	
	cookies map[*fileObjCookie]struct{}
}





type PortEvent struct {
	Cookie interface{}
	Events int32
	Fd     uintptr
	Path   string
	Source uint16
	fobj   *fileObj
}



func NewEventPort() (*EventPort, error) {
	port, err := port_create()
	if err != nil {
		return nil, err
	}
	e := &EventPort{
		port:    port,
		fds:     make(map[uintptr]*fileObjCookie),
		paths:   make(map[string]*fileObjCookie),
		cookies: make(map[*fileObjCookie]struct{}),
	}
	return e, nil
}








func (e *EventPort) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	err := Close(e.port)
	if err != nil {
		return err
	}
	e.fds = nil
	e.paths = nil
	e.cookies = nil
	return nil
}


func (e *EventPort) PathIsWatched(path string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, found := e.paths[path]
	return found
}


func (e *EventPort) FdIsWatched(fd uintptr) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, found := e.fds[fd]
	return found
}



func (e *EventPort) AssociatePath(path string, stat os.FileInfo, events int, cookie interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, found := e.paths[path]; found {
		return fmt.Errorf("%v is already associated with this Event Port", path)
	}
	fCookie, err := createFileObjCookie(path, stat, cookie)
	if err != nil {
		return err
	}
	_, err = port_associate(e.port, PORT_SOURCE_FILE, uintptr(unsafe.Pointer(fCookie.fobj)), events, (*byte)(unsafe.Pointer(fCookie)))
	if err != nil {
		return err
	}
	e.paths[path] = fCookie
	e.cookies[fCookie] = struct{}{}
	return nil
}


func (e *EventPort) DissociatePath(path string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	f, ok := e.paths[path]
	if !ok {
		return fmt.Errorf("%v is not associated with this Event Port", path)
	}
	_, err := port_dissociate(e.port, PORT_SOURCE_FILE, uintptr(unsafe.Pointer(f.fobj)))
	
	
	
	if err != nil && err != ENOENT {
		return err
	}
	if err == nil {
		
		fCookie := e.paths[path]
		delete(e.cookies, fCookie)
	}
	delete(e.paths, path)
	return err
}


func (e *EventPort) AssociateFd(fd uintptr, events int, cookie interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, found := e.fds[fd]; found {
		return fmt.Errorf("%v is already associated with this Event Port", fd)
	}
	fCookie, err := createFileObjCookie("", nil, cookie)
	if err != nil {
		return err
	}
	_, err = port_associate(e.port, PORT_SOURCE_FD, fd, events, (*byte)(unsafe.Pointer(fCookie)))
	if err != nil {
		return err
	}
	e.fds[fd] = fCookie
	e.cookies[fCookie] = struct{}{}
	return nil
}


func (e *EventPort) DissociateFd(fd uintptr) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.fds[fd]
	if !ok {
		return fmt.Errorf("%v is not associated with this Event Port", fd)
	}
	_, err := port_dissociate(e.port, PORT_SOURCE_FD, fd)
	if err != nil && err != ENOENT {
		return err
	}
	if err == nil {
		
		fCookie := e.fds[fd]
		delete(e.cookies, fCookie)
	}
	delete(e.fds, fd)
	return err
}

func createFileObjCookie(name string, stat os.FileInfo, cookie interface{}) (*fileObjCookie, error) {
	fCookie := new(fileObjCookie)
	fCookie.cookie = cookie
	if name != "" && stat != nil {
		fCookie.fobj = new(fileObj)
		bs, err := ByteSliceFromString(name)
		if err != nil {
			return nil, err
		}
		fCookie.fobj.Name = (*int8)(unsafe.Pointer(&bs[0]))
		s := stat.Sys().(*syscall.Stat_t)
		fCookie.fobj.Atim.Sec = s.Atim.Sec
		fCookie.fobj.Atim.Nsec = s.Atim.Nsec
		fCookie.fobj.Mtim.Sec = s.Mtim.Sec
		fCookie.fobj.Mtim.Nsec = s.Mtim.Nsec
		fCookie.fobj.Ctim.Sec = s.Ctim.Sec
		fCookie.fobj.Ctim.Nsec = s.Ctim.Nsec
	}
	return fCookie, nil
}


func (e *EventPort) GetOne(t *Timespec) (*PortEvent, error) {
	pe := new(portEvent)
	_, err := port_get(e.port, pe, t)
	if err != nil {
		return nil, err
	}
	p := new(PortEvent)
	e.mu.Lock()
	defer e.mu.Unlock()
	err = e.peIntToExt(pe, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}



func (e *EventPort) peIntToExt(peInt *portEvent, peExt *PortEvent) error {
	if e.cookies == nil {
		return fmt.Errorf("this EventPort is already closed")
	}
	peExt.Events = peInt.Events
	peExt.Source = peInt.Source
	fCookie := (*fileObjCookie)(unsafe.Pointer(peInt.User))
	_, found := e.cookies[fCookie]

	if !found {
		panic("unexpected event port address; may be due to kernel bug; see https://go.dev/issue/54254")
	}
	peExt.Cookie = fCookie.cookie
	delete(e.cookies, fCookie)

	switch peInt.Source {
	case PORT_SOURCE_FD:
		peExt.Fd = uintptr(peInt.Object)
		
		if fobj, ok := e.fds[peExt.Fd]; ok {
			if fobj == fCookie {
				delete(e.fds, peExt.Fd)
			}
		}
	case PORT_SOURCE_FILE:
		peExt.fobj = fCookie.fobj
		peExt.Path = BytePtrToString((*byte)(unsafe.Pointer(peExt.fobj.Name)))
		
		if fobj, ok := e.paths[peExt.Path]; ok {
			if fobj == fCookie {
				delete(e.paths, peExt.Path)
			}
		}
	}
	return nil
}


func (e *EventPort) Pending() (int, error) {
	var n uint32 = 0
	_, err := port_getn(e.port, nil, 0, &n, nil)
	return int(n), err
}





func (e *EventPort) Get(s []PortEvent, min int, timeout *Timespec) (int, error) {
	if min == 0 {
		return 0, fmt.Errorf("need to request at least one event or use Pending() instead")
	}
	if len(s) < min {
		return 0, fmt.Errorf("len(s) (%d) is less than min events requested (%d)", len(s), min)
	}
	got := uint32(min)
	max := uint32(len(s))
	var err error
	ps := make([]portEvent, max)
	_, err = port_getn(e.port, &ps[0], max, &got, timeout)
	
	if err != nil && err != ETIME {
		return 0, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	valid := 0
	for i := 0; i < int(got); i++ {
		err2 := e.peIntToExt(&ps[i], &s[i])
		if err2 != nil {
			if valid == 0 && err == nil {
				
				
				err = err2
			}
			break
		}
		valid = i + 1
	}
	return valid, err
}



func Putmsg(fd int, cl []byte, data []byte, flags int) (err error) {
	var clp, datap *strbuf
	if len(cl) > 0 {
		clp = &strbuf{
			Len: int32(len(cl)),
			Buf: (*int8)(unsafe.Pointer(&cl[0])),
		}
	}
	if len(data) > 0 {
		datap = &strbuf{
			Len: int32(len(data)),
			Buf: (*int8)(unsafe.Pointer(&data[0])),
		}
	}
	return putmsg(fd, clp, datap, flags)
}



func Getmsg(fd int, cl []byte, data []byte) (retCl []byte, retData []byte, flags int, err error) {
	var clp, datap *strbuf
	if len(cl) > 0 {
		clp = &strbuf{
			Maxlen: int32(len(cl)),
			Buf:    (*int8)(unsafe.Pointer(&cl[0])),
		}
	}
	if len(data) > 0 {
		datap = &strbuf{
			Maxlen: int32(len(data)),
			Buf:    (*int8)(unsafe.Pointer(&data[0])),
		}
	}

	if err = getmsg(fd, clp, datap, &flags); err != nil {
		return nil, nil, 0, err
	}

	if len(cl) > 0 {
		retCl = cl[:clp.Len]
	}
	if len(data) > 0 {
		retData = data[:datap.Len]
	}
	return retCl, retData, flags, nil
}

func IoctlSetIntRetInt(fd int, req int, arg int) (int, error) {
	return ioctlRet(fd, req, uintptr(arg))
}

func IoctlSetString(fd int, req int, val string) error {
	bs := make([]byte, len(val)+1)
	copy(bs[:len(bs)-1], val)
	err := ioctlPtr(fd, req, unsafe.Pointer(&bs[0]))
	runtime.KeepAlive(&bs[0])
	return err
}



func (l *Lifreq) SetName(name string) error {
	if len(name) >= len(l.Name) {
		return fmt.Errorf("name cannot be more than %d characters", len(l.Name)-1)
	}
	for i := range name {
		l.Name[i] = int8(name[i])
	}
	return nil
}

func (l *Lifreq) SetLifruInt(d int) {
	*(*int)(unsafe.Pointer(&l.Lifru[0])) = d
}

func (l *Lifreq) GetLifruInt() int {
	return *(*int)(unsafe.Pointer(&l.Lifru[0]))
}

func (l *Lifreq) SetLifruUint(d uint) {
	*(*uint)(unsafe.Pointer(&l.Lifru[0])) = d
}

func (l *Lifreq) GetLifruUint() uint {
	return *(*uint)(unsafe.Pointer(&l.Lifru[0]))
}

func IoctlLifreq(fd int, req int, l *Lifreq) error {
	return ioctlPtr(fd, req, unsafe.Pointer(l))
}



func (s *Strioctl) SetInt(i int) {
	s.Len = int32(unsafe.Sizeof(i))
	s.Dp = (*int8)(unsafe.Pointer(&i))
}

func IoctlSetStrioctlRetInt(fd int, req int, s *Strioctl) (int, error) {
	return ioctlPtrRet(fd, req, unsafe.Pointer(s))
}
















type Ucred struct {
	ucred uintptr
}



func ucredFinalizer(u *Ucred) {
	ucredFree(u.ucred)
}

func GetPeerUcred(fd uintptr) (*Ucred, error) {
	var ucred uintptr
	err := getpeerucred(fd, &ucred)
	if err != nil {
		return nil, err
	}
	result := &Ucred{
		ucred: ucred,
	}
	
	runtime.SetFinalizer(result, ucredFinalizer)
	return result, nil
}

func UcredGet(pid int) (*Ucred, error) {
	ucred, err := ucredGet(pid)
	if err != nil {
		return nil, err
	}
	result := &Ucred{
		ucred: ucred,
	}
	
	runtime.SetFinalizer(result, ucredFinalizer)
	return result, nil
}

func (u *Ucred) Geteuid() int {
	defer runtime.KeepAlive(u)
	return ucredGeteuid(u.ucred)
}

func (u *Ucred) Getruid() int {
	defer runtime.KeepAlive(u)
	return ucredGetruid(u.ucred)
}

func (u *Ucred) Getsuid() int {
	defer runtime.KeepAlive(u)
	return ucredGetsuid(u.ucred)
}

func (u *Ucred) Getegid() int {
	defer runtime.KeepAlive(u)
	return ucredGetegid(u.ucred)
}

func (u *Ucred) Getrgid() int {
	defer runtime.KeepAlive(u)
	return ucredGetrgid(u.ucred)
}

func (u *Ucred) Getsgid() int {
	defer runtime.KeepAlive(u)
	return ucredGetsgid(u.ucred)
}

func (u *Ucred) Getpid() int {
	defer runtime.KeepAlive(u)
	return ucredGetpid(u.ucred)
}
