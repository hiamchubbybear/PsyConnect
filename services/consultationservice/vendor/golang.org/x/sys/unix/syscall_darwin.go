











package unix

import (
	"fmt"
	"syscall"
	"unsafe"
)




func fdopendir(fd int) (dir uintptr, err error) {
	r0, _, e1 := syscall_syscallPtr(libc_fdopendir_trampoline_addr, uintptr(fd), 0, 0)
	dir = uintptr(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

var libc_fdopendir_trampoline_addr uintptr

//go:cgo_import_dynamic libc_fdopendir fdopendir "/usr/lib/libSystem.B.dylib"

func Getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error) {
	
	
	
	
	
	
	
	skip, err := Seek(fd, 0, 1 )
	if err != nil {
		return 0, err
	}

	
	
	
	
	
	
	fd2, err := Openat(fd, ".", O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	d, err := fdopendir(fd2)
	if err != nil {
		Close(fd2)
		return 0, err
	}
	defer closedir(d)

	var cnt int64
	for {
		var entry Dirent
		var entryp *Dirent
		e := readdir_r(d, &entry, &entryp)
		if e != 0 {
			return n, errnoErr(e)
		}
		if entryp == nil {
			break
		}
		if skip > 0 {
			skip--
			cnt++
			continue
		}

		reclen := int(entry.Reclen)
		if reclen > len(buf) {
			
			
			
			
			break
		}

		
		s := unsafe.Slice((*byte)(unsafe.Pointer(&entry)), reclen)
		copy(buf, s)

		buf = buf[reclen:]
		n += reclen
		cnt++
	}
	
	
	_, err = Seek(fd, cnt, 0 )
	if err != nil {
		return n, err
	}

	return n, nil
}


type SockaddrDatalink struct {
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [12]int8
	raw    RawSockaddrDatalink
}


type SockaddrCtl struct {
	ID   uint32
	Unit uint32
	raw  RawSockaddrCtl
}

func (sa *SockaddrCtl) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Sc_len = SizeofSockaddrCtl
	sa.raw.Sc_family = AF_SYSTEM
	sa.raw.Ss_sysaddr = AF_SYS_CONTROL
	sa.raw.Sc_id = sa.ID
	sa.raw.Sc_unit = sa.Unit
	return unsafe.Pointer(&sa.raw), SizeofSockaddrCtl, nil
}





type SockaddrVM struct {
	
	
	
	
	
	CID  uint32
	Port uint32
	raw  RawSockaddrVM
}

func (sa *SockaddrVM) sockaddr() (unsafe.Pointer, _Socklen, error) {
	sa.raw.Len = SizeofSockaddrVM
	sa.raw.Family = AF_VSOCK
	sa.raw.Port = sa.Port
	sa.raw.Cid = sa.CID

	return unsafe.Pointer(&sa.raw), SizeofSockaddrVM, nil
}

func anyToSockaddrGOOS(fd int, rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_SYSTEM:
		pp := (*RawSockaddrCtl)(unsafe.Pointer(rsa))
		if pp.Ss_sysaddr == AF_SYS_CONTROL {
			sa := new(SockaddrCtl)
			sa.ID = pp.Sc_id
			sa.Unit = pp.Sc_unit
			return sa, nil
		}
	case AF_VSOCK:
		pp := (*RawSockaddrVM)(unsafe.Pointer(rsa))
		sa := &SockaddrVM{
			CID:  pp.Cid,
			Port: pp.Port,
		}
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}




const SYS___SYSCTL = SYS_SYSCTL


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
	return readInt(buf, unsafe.Offsetof(Dirent{}.Ino), unsafe.Sizeof(Dirent{}.Ino))
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Reclen), unsafe.Sizeof(Dirent{}.Reclen))
}

func direntNamlen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Namlen), unsafe.Sizeof(Dirent{}.Namlen))
}

func PtraceAttach(pid int) (err error) { return ptrace(PT_ATTACH, pid, 0, 0) }
func PtraceDetach(pid int) (err error) { return ptrace(PT_DETACH, pid, 0, 0) }
func PtraceDenyAttach() (err error)    { return ptrace(PT_DENY_ATTACH, 0, 0, 0) }



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var x [2]int32
	err = pipe(&x)
	if err == nil {
		p[0] = int(x[0])
		p[1] = int(x[1])
	}
	return
}

func Getfsstat(buf []Statfs_t, flags int) (n int, err error) {
	var _p0 unsafe.Pointer
	var bufsize uintptr
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
		bufsize = unsafe.Sizeof(Statfs_t{}) * uintptr(len(buf))
	}
	return getfsstat(_p0, bufsize, flags)
}

func xattrPointer(dest []byte) *byte {
	
	
	
	
	
	var destp *byte
	if len(dest) > 0 {
		destp = &dest[0]
	}
	return destp
}



func Getxattr(path string, attr string, dest []byte) (sz int, err error) {
	return getxattr(path, attr, xattrPointer(dest), len(dest), 0, 0)
}

func Lgetxattr(link string, attr string, dest []byte) (sz int, err error) {
	return getxattr(link, attr, xattrPointer(dest), len(dest), 0, XATTR_NOFOLLOW)
}



func Fgetxattr(fd int, attr string, dest []byte) (sz int, err error) {
	return fgetxattr(fd, attr, xattrPointer(dest), len(dest), 0, 0)
}



func Setxattr(path string, attr string, data []byte, flags int) (err error) {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	return setxattr(path, attr, xattrPointer(data), len(data), 0, flags)
}

func Lsetxattr(link string, attr string, data []byte, flags int) (err error) {
	return setxattr(link, attr, xattrPointer(data), len(data), 0, flags|XATTR_NOFOLLOW)
}



func Fsetxattr(fd int, attr string, data []byte, flags int) (err error) {
	return fsetxattr(fd, attr, xattrPointer(data), len(data), 0, 0)
}



func Removexattr(path string, attr string) (err error) {
	
	
	
	return removexattr(path, attr, 0)
}

func Lremovexattr(link string, attr string) (err error) {
	return removexattr(link, attr, XATTR_NOFOLLOW)
}



func Fremovexattr(fd int, attr string) (err error) {
	return fremovexattr(fd, attr, 0)
}



func Listxattr(path string, dest []byte) (sz int, err error) {
	return listxattr(path, xattrPointer(dest), len(dest), 0)
}

func Llistxattr(link string, dest []byte) (sz int, err error) {
	return listxattr(link, xattrPointer(dest), len(dest), XATTR_NOFOLLOW)
}



func Flistxattr(fd int, dest []byte) (sz int, err error) {
	return flistxattr(fd, xattrPointer(dest), len(dest), 0)
}









func Kill(pid int, signum syscall.Signal) (err error) { return kill(pid, int(signum), 1) }




func IoctlCtlInfo(fd int, ctlInfo *CtlInfo) error {
	return ioctlPtr(fd, CTLIOCGINFO, unsafe.Pointer(ctlInfo))
}


type IfreqMTU struct {
	Name [IFNAMSIZ]byte
	MTU  int32
}



func IoctlGetIfreqMTU(fd int, ifname string) (*IfreqMTU, error) {
	var ifreq IfreqMTU
	copy(ifreq.Name[:], ifname)
	err := ioctlPtr(fd, SIOCGIFMTU, unsafe.Pointer(&ifreq))
	return &ifreq, err
}



func IoctlSetIfreqMTU(fd int, ifreq *IfreqMTU) error {
	return ioctlPtr(fd, SIOCSIFMTU, unsafe.Pointer(ifreq))
}



func RenamexNp(from string, to string, flag uint32) (err error) {
	return renamexNp(from, to, flag)
}



func RenameatxNp(fromfd int, from string, tofd int, to string, flag uint32) (err error) {
	return renameatxNp(fromfd, from, tofd, to, flag)
}



func Uname(uname *Utsname) error {
	mib := []_C_int{CTL_KERN, KERN_OSTYPE}
	n := unsafe.Sizeof(uname.Sysname)
	if err := sysctl(mib, &uname.Sysname[0], &n, nil, 0); err != nil {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_HOSTNAME}
	n = unsafe.Sizeof(uname.Nodename)
	if err := sysctl(mib, &uname.Nodename[0], &n, nil, 0); err != nil {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_OSRELEASE}
	n = unsafe.Sizeof(uname.Release)
	if err := sysctl(mib, &uname.Release[0], &n, nil, 0); err != nil {
		return err
	}

	mib = []_C_int{CTL_KERN, KERN_VERSION}
	n = unsafe.Sizeof(uname.Version)
	if err := sysctl(mib, &uname.Version[0], &n, nil, 0); err != nil {
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
	if err := sysctl(mib, &uname.Machine[0], &n, nil, 0); err != nil {
		return err
	}

	return nil
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	var length = int64(count)
	err = sendfile(infd, outfd, *offset, &length, nil, 0)
	written = int(length)
	return
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

func GetsockoptTCPConnectionInfo(fd, level, opt int) (*TCPConnectionInfo, error) {
	var value TCPConnectionInfo
	vallen := _Socklen(SizeofTCPConnectionInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func SysctlKinfoProc(name string, args ...int) (*KinfoProc, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return nil, err
	}

	var kinfo KinfoProc
	n := uintptr(SizeofKinfoProc)
	if err := sysctl(mib, (*byte)(unsafe.Pointer(&kinfo)), &n, nil, 0); err != nil {
		return nil, err
	}
	if n != SizeofKinfoProc {
		return nil, EIO
	}
	return &kinfo, nil
}

func SysctlKinfoProcSlice(name string, args ...int) ([]KinfoProc, error) {
	mib, err := sysctlmib(name, args...)
	if err != nil {
		return nil, err
	}

	for {
		
		n := uintptr(0)
		if err := sysctl(mib, nil, &n, nil, 0); err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, nil
		}
		if n%SizeofKinfoProc != 0 {
			return nil, fmt.Errorf("sysctl() returned a size of %d, which is not a multiple of %d", n, SizeofKinfoProc)
		}

		
		buf := make([]KinfoProc, n/SizeofKinfoProc)
		if err := sysctl(mib, (*byte)(unsafe.Pointer(&buf[0])), &n, nil, 0); err != nil {
			if err == ENOMEM {
				
				continue
			}
			return nil, err
		}
		if n%SizeofKinfoProc != 0 {
			return nil, fmt.Errorf("sysctl() returned a size of %d, which is not a multiple of %d", n, SizeofKinfoProc)
		}

		
		
		return buf[:n/SizeofKinfoProc], nil
	}
}



func PthreadChdir(path string) (err error) {
	return pthread_chdir_np(path)
}



func PthreadFchdir(fd int) (err error) {
	return pthread_fchdir_np(fd)
}










func Connectx(fd int, srcIf uint32, srcAddr, dstAddr Sockaddr, associd SaeAssocID, flags uint32, iov []Iovec, connid *SaeConnID) (n uintptr, err error) {
	endpoints := SaEndpoints{
		Srcif: srcIf,
	}

	if srcAddr != nil {
		addrp, addrlen, err := srcAddr.sockaddr()
		if err != nil {
			return 0, err
		}
		endpoints.Srcaddr = (*RawSockaddr)(addrp)
		endpoints.Srcaddrlen = uint32(addrlen)
	}

	if dstAddr != nil {
		addrp, addrlen, err := dstAddr.sockaddr()
		if err != nil {
			return 0, err
		}
		endpoints.Dstaddr = (*RawSockaddr)(addrp)
		endpoints.Dstaddrlen = uint32(addrlen)
	}

	err = connectx(fd, &endpoints, associd, flags, iov, &n, connid)
	return
}






































































































