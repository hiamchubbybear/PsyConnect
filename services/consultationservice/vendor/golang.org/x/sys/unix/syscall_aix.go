



//go:build aix








package unix

import "unsafe"



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



func Utimes(path string, tv []Timeval) error {
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}



func UtimesNano(path string, ts []Timespec) error {
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
	if n > len(sa.raw.Path) {
		return nil, 0, EINVAL
	}
	if n == len(sa.raw.Path) && name[0] != '@' {
		return nil, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = uint8(name[i])
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



const ImplementsGetwd = true

func Getwd() (ret string, err error) {
	for len := uint64(4096); ; len *= 2 {
		b := make([]byte, len)
		err := getcwd(b)
		if err == nil {
			i := 0
			for b[i] != 0 {
				i++
			}
			return string(b[0:i]), nil
		}
		if err != ERANGE {
			return "", err
		}
	}
}

func Getcwd(buf []byte) (n int, err error) {
	err = getcwd(buf)
	if err == nil {
		i := 0
		for buf[i] != 0 {
			i++
		}
		n = i + 1
	}
	return
}

func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	
	if n < 0 || n > 1000 {
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
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	}
	if len(iov) > 0 {
		msg.Iov = &iov[0]
		msg.SetIovlen(len(iov))
	}
	if n, err = recvmsg(fd, &msg, flags); n == -1 {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
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
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
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

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)

		
		
		n := SizeofSockaddrUnix - 3 
		for i := 0; i < n; i++ {
			if pp.Path[i] == 0 {
				n = i
				break
			}
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

func Gettimeofday(tv *Timeval) (err error) {
	err = gettimeofday(tv, nil)
	return
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	return sendfile(outfd, infd, offset, count)
}


func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	return -1, ENOSYS
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



func Getdents(fd int, buf []byte) (n int, err error) {
	return getdirent(fd, buf)
}



func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	var status _C_int
	var r Pid_t
	err = ERESTART
	
	
	for err == ERESTART {
		r, err = wait4(Pid_t(pid), &status, options, rusage)
	}
	wpid = int(r)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return
}



type WaitStatus uint32

func (w WaitStatus) Stopped() bool { return w&0x40 != 0 }
func (w WaitStatus) StopSignal() Signal {
	if !w.Stopped() {
		return -1
	}
	return Signal(w>>8) & 0xFF
}

func (w WaitStatus) Exited() bool { return w&0xFF == 0 }
func (w WaitStatus) ExitStatus() int {
	if !w.Exited() {
		return -1
	}
	return int((w >> 8) & 0xFF)
}

func (w WaitStatus) Signaled() bool { return w&0x40 == 0 && w&0xFF != 0 }
func (w WaitStatus) Signal() Signal {
	if !w.Signaled() {
		return -1
	}
	return Signal(w>>16) & 0xFF
}

func (w WaitStatus) Continued() bool { return w&0x01000000 != 0 }

func (w WaitStatus) CoreDump() bool { return w&0x80 == 0x80 }

func (w WaitStatus) TrapCause() int { return -1 }


















func Fsync(fd int) error {
	return fsyncRange(fd, O_SYNC, 0, 0)
}

















































































































func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err = pipe(&pp)
	if err == nil {
		p[0] = int(pp[0])
		p[1] = int(pp[1])
	}
	return
}



func Poll(fds []PollFd, timeout int) (n int, err error) {
	if len(fds) == 0 {
		return poll(nil, 0, timeout)
	}
	return poll(&fds[0], len(fds), timeout)
}









func Unmount(target string, flags int) (err error) {
	if flags != 0 {
		
		return ENOSYS
	}
	return umount(target)
}
