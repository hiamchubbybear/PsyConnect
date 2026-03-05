











package unix

import (
	"sync"
	"unsafe"
)


var (
	osreldateOnce sync.Once
	osreldate     uint32
)



const _dragonflyABIChangeVersion = 500705

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
	Data   [12]int8
	Rcf    uint16
	Route  [16]uint16
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
	namlen, ok := direntNamlen(buf)
	if !ok {
		return 0, false
	}
	return (16 + namlen + 1 + 7) &^ 7, true
}

func direntNamlen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(Dirent{}.Namlen), unsafe.Sizeof(Dirent{}.Namlen))
}



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	r, w, err := pipe()
	if err == nil {
		p[0], p[1] = r, w
	}
	return
}



func Pipe2(p []int, flags int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	
	
	r, w, err := pipe2(&pp, flags)
	if err == nil {
		p[0], p[1] = r, w
	}
	return err
}



func pread(fd int, p []byte, offset int64) (n int, err error) {
	return extpread(fd, p, 0, offset)
}



func pwrite(fd int, p []byte, offset int64) (n int, err error) {
	return extpwrite(fd, p, 0, offset)
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
	var _p0 unsafe.Pointer
	var bufsize uintptr
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






func sysctlUname(mib []_C_int, old *byte, oldlen *uintptr) error {
	err := sysctl(mib, old, oldlen, nil, 0)
	if err != nil {
		
		
		
		if err == ENOMEM {
			err = nil
		}
	}
	return err
}

func Uname(uname *Utsname) error {
	mib := []_C_int{CTL_KERN, KERN_OSTYPE}
	n := unsafe.Sizeof(uname.Sysname)
	if err := sysctlUname(mib, &uname.Sysname[0], &n); err != nil {
		return err
	}
	uname.Sysname[unsafe.Sizeof(uname.Sysname)-1] = 0

	mib = []_C_int{CTL_KERN, KERN_HOSTNAME}
	n = unsafe.Sizeof(uname.Nodename)
	if err := sysctlUname(mib, &uname.Nodename[0], &n); err != nil {
		return err
	}
	uname.Nodename[unsafe.Sizeof(uname.Nodename)-1] = 0

	mib = []_C_int{CTL_KERN, KERN_OSRELEASE}
	n = unsafe.Sizeof(uname.Release)
	if err := sysctlUname(mib, &uname.Release[0], &n); err != nil {
		return err
	}
	uname.Release[unsafe.Sizeof(uname.Release)-1] = 0

	mib = []_C_int{CTL_KERN, KERN_VERSION}
	n = unsafe.Sizeof(uname.Version)
	if err := sysctlUname(mib, &uname.Version[0], &n); err != nil {
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
	if err := sysctlUname(mib, &uname.Machine[0], &n); err != nil {
		return err
	}
	uname.Machine[unsafe.Sizeof(uname.Machine)-1] = 0

	return nil
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	return sendfile(outfd, infd, offset, count)
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


































































































