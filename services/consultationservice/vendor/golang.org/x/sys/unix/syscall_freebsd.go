











package unix

import (
	"errors"
	"sync"
	"unsafe"
)


var (
	osreldateOnce sync.Once
	osreldate     uint32
)

func supportsABI(ver uint32) bool {
	osreldateOnce.Do(func() { osreldate, _ = SysctlUint32("kern.osreldate") })
	return osreldate >= ver
}


type SockaddrDatalink struct {
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [46]int8
	raw    RawSockaddrDatalink
}

func anyToSockaddrGOOS(fd int, rsa *RawSockaddrAny) (Sockaddr, error) {
	return nil, EAFNOSUPPORT
}


func nametomib(name string) (mib []_C_int, err error) {
	const siz = unsafe.Sizeof(mib[0])

	
	
	
	
	
	
	
	var buf [CTL_MAXNAME + 2]_C_int
	n := uintptr(CTL_MAXNAME) * siz

	p := (*byte)(unsafe.Pointer(&buf[0]))
	bytes, err := ByteSliceFromString(name)
	if err != nil {
		return nil, err
	}

	
	
	if err = sysctl([]_C_int{0, 3}, p, &n, &bytes[0], uintptr(len(name))); err != nil {
		return nil, err
	}
	return buf[0 : n/siz], nil
}

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Fileno), unsafe.Sizeof(Dirent{}.Fileno))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Reclen), unsafe.Sizeof(Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Namlen), unsafe.Sizeof(Dirent{}.Namlen))
}

func Pipe(p []int) (err error) {
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

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error) {
	var value IPMreqn
	vallen := _Socklen(SizeofIPMreqn)
	errno := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, errno
}

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
}



func GetsockoptXucred(fd, level, opt int) (*Xucred, error) {
	x := new(Xucred)
	vallen := _Socklen(SizeofXucred)
	err := getsockopt(fd, level, opt, unsafe.Pointer(x), &vallen)
	return x, err
}

func Accept4(fd, flags int) (nfd int, sa Sockaddr, err error) {
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



func Getfsstat(buf []Statfs_t, flags int) (n int, err error) {
	var (
		_p0     unsafe.Pointer
		bufsize uintptr
	)
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
		bufsize = unsafe.Sizeof(Statfs_t{}) * uintptr(len(buf))
	}
	r0, _, e1 := Syscall(SYS_GETFSSTAT, uintptr(_p0), bufsize, uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}






func Uname(uname *Utsname) error {
	mib := []_C_int{CTL_KERN, KERN_OSTYPE}
	n := unsafe.Sizeof(uname.Sysname)
	
	if err := sysctl(mib, &uname.Sysname[0], &n, nil, 0); err != nil && !errors.Is(err, ENOMEM) {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_HOSTNAME}
	n = unsafe.Sizeof(uname.Nodename)
	if err := sysctl(mib, &uname.Nodename[0], &n, nil, 0); err != nil && !errors.Is(err, ENOMEM) {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_OSRELEASE}
	n = unsafe.Sizeof(uname.Release)
	if err := sysctl(mib, &uname.Release[0], &n, nil, 0); err != nil && !errors.Is(err, ENOMEM) {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_VERSION}
	n = unsafe.Sizeof(uname.Version)
	if err := sysctl(mib, &uname.Version[0], &n, nil, 0); err != nil && !errors.Is(err, ENOMEM) {
		return err
	}

	
	
	for i, b := range uname.Version {
		if b == '\n' || b == '\t' {
			if i == len(uname.Version)-1 {
				uname.Version[i] = 0
			} else {
				uname.Version[i] = ' '
			}
		}
	}

	mib = []_C_int{CTL_HW, HW_MACHINE}
	n = unsafe.Sizeof(uname.Machine)
	if err := sysctl(mib, &uname.Machine[0], &n, nil, 0); err != nil && !errors.Is(err, ENOMEM) {
		return err
	}

	return nil
}

func Stat(path string, st *Stat_t) (err error) {
	return Fstatat(AT_FDCWD, path, st, 0)
}

func Lstat(path string, st *Stat_t) (err error) {
	return Fstatat(AT_FDCWD, path, st, AT_SYMLINK_NOFOLLOW)
}

func Getdents(fd int, buf []byte) (n int, err error) {
	return Getdirentries(fd, buf, nil)
}

func Getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error) {
	if basep == nil || unsafe.Sizeof(*basep) == 8 {
		return getdirentries(fd, buf, (*uint64)(unsafe.Pointer(basep)))
	}
	
	
	var base uint64 = uint64(*basep)
	n, err = getdirentries(fd, buf, &base)
	*basep = uintptr(base)
	if base>>32 != 0 {
		
		
		
		err = EIO
	}
	return
}

func Mknod(path string, mode uint32, dev uint64) (err error) {
	return Mknodat(AT_FDCWD, path, mode, dev)
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	return sendfile(outfd, infd, offset, count)
}




func PtraceAttach(pid int) (err error) {
	return ptrace(PT_ATTACH, pid, 0, 0)
}

func PtraceCont(pid int, signal int) (err error) {
	return ptrace(PT_CONTINUE, pid, 1, signal)
}

func PtraceDetach(pid int) (err error) {
	return ptrace(PT_DETACH, pid, 1, 0)
}

func PtraceGetFpRegs(pid int, fpregsout *FpReg) (err error) {
	return ptracePtr(PT_GETFPREGS, pid, unsafe.Pointer(fpregsout), 0)
}

func PtraceGetRegs(pid int, regsout *Reg) (err error) {
	return ptracePtr(PT_GETREGS, pid, unsafe.Pointer(regsout), 0)
}

func PtraceIO(req int, pid int, offs uintptr, out []byte, countin int) (count int, err error) {
	ioDesc := PtraceIoDesc{
		Op:   int32(req),
		Offs: offs,
	}
	if countin > 0 {
		_ = out[:countin] 
		ioDesc.Addr = &out[0]
	} else if out != nil {
		ioDesc.Addr = (*byte)(unsafe.Pointer(&_zero))
	}
	ioDesc.SetLen(countin)

	err = ptracePtr(PT_IO, pid, unsafe.Pointer(&ioDesc), 0)
	return int(ioDesc.Len), err
}

func PtraceLwpEvents(pid int, enable int) (err error) {
	return ptrace(PT_LWP_EVENTS, pid, 0, enable)
}

func PtraceLwpInfo(pid int, info *PtraceLwpInfoStruct) (err error) {
	return ptracePtr(PT_LWPINFO, pid, unsafe.Pointer(info), int(unsafe.Sizeof(*info)))
}

func PtracePeekData(pid int, addr uintptr, out []byte) (count int, err error) {
	return PtraceIO(PIOD_READ_D, pid, addr, out, SizeofLong)
}

func PtracePeekText(pid int, addr uintptr, out []byte) (count int, err error) {
	return PtraceIO(PIOD_READ_I, pid, addr, out, SizeofLong)
}

func PtracePokeData(pid int, addr uintptr, data []byte) (count int, err error) {
	return PtraceIO(PIOD_WRITE_D, pid, addr, data, SizeofLong)
}

func PtracePokeText(pid int, addr uintptr, data []byte) (count int, err error) {
	return PtraceIO(PIOD_WRITE_I, pid, addr, data, SizeofLong)
}

func PtraceSetRegs(pid int, regs *Reg) (err error) {
	return ptracePtr(PT_SETREGS, pid, unsafe.Pointer(regs), 0)
}

func PtraceSingleStep(pid int) (err error) {
	return ptrace(PT_STEP, pid, 1, 0)
}

func Dup3(oldfd, newfd, flags int) error {
	if oldfd == newfd || flags&^O_CLOEXEC != 0 {
		return EINVAL
	}
	how := F_DUP2FD
	if flags&O_CLOEXEC != 0 {
		how = F_DUP2FD_CLOEXEC
	}
	_, err := fcntl(oldfd, how, newfd)
	return err
}

















































































































