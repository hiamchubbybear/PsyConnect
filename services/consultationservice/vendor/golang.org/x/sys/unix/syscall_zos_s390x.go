



//go:build zos && s390x








package unix

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

//go:noescape
func initZosLibVec()

//go:noescape
func GetZosLibVec() uintptr

func init() {
	initZosLibVec()
	r0, _, _ := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS_____GETENV_A<<4, uintptr(unsafe.Pointer(&([]byte("__ZOS_XSYSTRACE\x00"))[0])))
	if r0 != 0 {
		n, _, _ := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS___ATOI_A<<4, r0)
		ZosTraceLevel = int(n)
		r0, _, _ := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS_____GETENV_A<<4, uintptr(unsafe.Pointer(&([]byte("__ZOS_XSYSTRACEFD\x00"))[0])))
		if r0 != 0 {
			fd, _, _ := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS___ATOI_A<<4, r0)
			f := os.NewFile(fd, "zostracefile")
			if f != nil {
				ZosTracefile = f
			}
		}

	}
}

//go:noescape
func CallLeFuncWithErr(funcdesc uintptr, parms ...uintptr) (ret, errno2 uintptr, err Errno)

//go:noescape
func CallLeFuncWithPtrReturn(funcdesc uintptr, parms ...uintptr) (ret, errno2 uintptr, err Errno)






//go:nosplit
func ptrtest(uintptr) uint64



//go:noescape
func safeload(ptr uintptr) (value uintptr, error uintptr)

const (
	entrypointLocationOffset = 8 

	xplinkEyecatcher   = 0x00c300c500c500f1 
	eyecatcherOffset   = 16                 
	ppa1LocationOffset = 8                  

	nameLenOffset = 0x14 
	nameOffset    = 0x16 
)

func getPpaOffset(funcptr uintptr) int64 {
	entrypoint, err := safeload(funcptr + entrypointLocationOffset)
	if err != 0 {
		return -1
	}

	
	val, err := safeload(entrypoint - eyecatcherOffset)
	if err != 0 {
		return -1
	}
	if val != xplinkEyecatcher {
		return -1
	}

	ppaoff, err := safeload(entrypoint - ppa1LocationOffset)
	if err != 0 {
		return -1
	}

	ppaoff >>= 32
	return int64(ppaoff)
}








func funcptrtest(funcptr uintptr, funcName string) uint64 {
	entrypoint, err := safeload(funcptr + entrypointLocationOffset)
	if err != 0 {
		return 1
	}

	ppaoff := getPpaOffset(funcptr)
	if ppaoff == -1 {
		return 1
	}

	
	ppa1 := (entrypoint - eyecatcherOffset) + uintptr(ppaoff)

	nameLen, err := safeload(ppa1 + nameLenOffset)
	if err != 0 {
		return 1
	}

	nameLen >>= 48
	if nameLen > 128 {
		return 1
	}

	
	if funcName == "" {
		return 0
	}

	var funcname [128]byte
	for i := 0; i < int(nameLen); i += 8 {
		v, err := safeload(ppa1 + nameOffset + uintptr(i))
		if err != 0 {
			return 1
		}
		funcname[i] = byte(v >> 56)
		funcname[i+1] = byte(v >> 48)
		funcname[i+2] = byte(v >> 40)
		funcname[i+3] = byte(v >> 32)
		funcname[i+4] = byte(v >> 24)
		funcname[i+5] = byte(v >> 16)
		funcname[i+6] = byte(v >> 8)
		funcname[i+7] = byte(v)
	}

	runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___E2A_L<<4, 
		[]uintptr{uintptr(unsafe.Pointer(&funcname[0])), nameLen})

	name := string(funcname[:nameLen])
	if name != funcName {
		return 1
	}

	return 0
}



func isValidLeFunc(f uintptr) error {
	ret := funcptrtest(f, "")
	if ret != 0 {
		return fmt.Errorf("Bad pointer, not an LE function ")
	}
	return nil
}


func getLeFuncName(f uintptr) (string, error) {
	
	entry := ((*[2]uintptr)(unsafe.Pointer(f)))[1]
	preamp := ((*[4]uint32)(unsafe.Pointer(entry - eyecatcherOffset)))

	offsetPpa1 := preamp[2]
	if offsetPpa1 > 0x0ffff {
		return "", fmt.Errorf("PPA1 offset seems too big 0x%x\n", offsetPpa1)
	}

	ppa1 := uintptr(unsafe.Pointer(preamp)) + uintptr(offsetPpa1)
	res := ptrtest(ppa1)
	if res != 0 {
		return "", fmt.Errorf("PPA1 address not valid")
	}

	size := *(*uint16)(unsafe.Pointer(ppa1 + nameLenOffset))
	if size > 128 {
		return "", fmt.Errorf("Function name seems too long, length=%d\n", size)
	}

	var name [128]byte
	funcname := (*[128]byte)(unsafe.Pointer(ppa1 + nameOffset))
	copy(name[0:size], funcname[0:size])

	runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___E2A_L<<4, 
		[]uintptr{uintptr(unsafe.Pointer(&name[0])), uintptr(size)})

	return string(name[:size]), nil
}


func zosLeVersion() (version, release uint32) {
	p1 := (*(*uintptr)(unsafe.Pointer(uintptr(1208)))) >> 32
	p1 = *(*uintptr)(unsafe.Pointer(uintptr(p1 + 88)))
	p1 = *(*uintptr)(unsafe.Pointer(uintptr(p1 + 8)))
	p1 = *(*uintptr)(unsafe.Pointer(uintptr(p1 + 984)))
	vrm := *(*uint32)(unsafe.Pointer(p1 + 80))
	version = (vrm & 0x00ff0000) >> 16
	release = (vrm & 0x0000ff00) >> 8
	return
}


func ZosStdioFilep(fd int32) uintptr {
	return uintptr(*(*uint64)(unsafe.Pointer(uintptr(*(*uint64)(unsafe.Pointer(uintptr(*(*uint64)(unsafe.Pointer(uintptr(uint64(*(*uint32)(unsafe.Pointer(uintptr(1208)))) + 80))) + uint64((fd+2)<<3))))))))
}

func copyStat(stat *Stat_t, statLE *Stat_LE_t) {
	stat.Dev = uint64(statLE.Dev)
	stat.Ino = uint64(statLE.Ino)
	stat.Nlink = uint64(statLE.Nlink)
	stat.Mode = uint32(statLE.Mode)
	stat.Uid = uint32(statLE.Uid)
	stat.Gid = uint32(statLE.Gid)
	stat.Rdev = uint64(statLE.Rdev)
	stat.Size = statLE.Size
	stat.Atim.Sec = int64(statLE.Atim)
	stat.Atim.Nsec = 0 
	stat.Mtim.Sec = int64(statLE.Mtim)
	stat.Mtim.Nsec = 0 
	stat.Ctim.Sec = int64(statLE.Ctim)
	stat.Ctim.Nsec = 0 
	stat.Blksize = int64(statLE.Blksize)
	stat.Blocks = statLE.Blocks
}

func svcCall(fnptr unsafe.Pointer, argv *unsafe.Pointer, dsa *uint64)
func svcLoad(name *byte) unsafe.Pointer
func svcUnload(name *byte, fnptr unsafe.Pointer) int64

func (d *Dirent) NameString() string {
	if d == nil {
		return ""
	}
	s := string(d.Name[:])
	idx := strings.IndexByte(s, 0)
	if idx == -1 {
		return s
	} else {
		return s[:idx]
	}
}

func DecodeData(dest []byte, sz int, val uint64) {
	for i := 0; i < sz; i++ {
		dest[sz-1-i] = byte((val >> (uint64(i * 8))) & 0xff)
	}
}

func EncodeData(data []byte) uint64 {
	var value uint64
	sz := len(data)
	for i := 0; i < sz; i++ {
		value |= uint64(data[i]) << uint64(((sz - i - 1) * 8))
	}
	return value
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Len = SizeofSockaddrInet4
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Len = SizeofSockaddrInet6
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) || n == 0 {
		return nil, 0, EINVAL
	}
	sa.raw.Len = byte(3 + n) 
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func anyToSockaddr(_ int, rsa *RawSockaddrAny) (Sockaddr, error) {
	
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		
		
		if pp.Len != 0 && pp.Path[0] == 0 {
			
			
			
			
			
			pp.Path[0] = '@'
		}

		
		
		
		
		
		
		
		
		
		
		
		n := 0
		for n < int(pp.Len) && pp.Path[n] != 0 {
			n++
		}
		sa.Name = string(unsafe.Slice((*byte)(unsafe.Pointer(&pp.Path[0])), n))
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}

func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if err != nil {
		return
	}
	
	sa, err = anyToSockaddr(0, &rsa)
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
	
	sa, err = anyToSockaddr(0, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func Ctermid() (tty string, err error) {
	var termdev [1025]byte
	runtime.EnterSyscall()
	r0, err2, err1 := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS___CTERMID_A<<4, uintptr(unsafe.Pointer(&termdev[0])))
	runtime.ExitSyscall()
	if r0 == 0 {
		return "", fmt.Errorf("%s (errno2=0x%x)\n", err1.Error(), err2)
	}
	s := string(termdev[:])
	idx := strings.Index(s, string(rune(0)))
	if idx == -1 {
		tty = s
	} else {
		tty = s[:idx]
	}
	return
}

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = int32(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = int32(length)
}























































func Faccessat2(dirfd int, path string, mode uint32, flags int) (err error) {
	return Faccessat(dirfd, path, mode, flags)
}











func Fstat(fd int, stat *Stat_t) (err error) {
	var statLE Stat_LE_t
	err = fstat(fd, &statLE)
	copyStat(stat, &statLE)
	return
}

func Fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) {
	var statLE Stat_LE_t
	err = fstatat(dirfd, path, &statLE, flags)
	copyStat(stat, &statLE)
	return
}

func impl_Getxattr(path string, attr string, dest []byte) (sz int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(attr)
	if err != nil {
		return
	}
	var _p2 unsafe.Pointer
	if len(dest) > 0 {
		_p2 = unsafe.Pointer(&dest[0])
	} else {
		_p2 = unsafe.Pointer(&_zero)
	}
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___GETXATTR_A<<4, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(dest)))
	sz = int(r0)
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_GetxattrAddr() *(func(path string, attr string, dest []byte) (sz int, err error))

var Getxattr = enter_Getxattr

func enter_Getxattr(path string, attr string, dest []byte) (sz int, err error) {
	funcref := get_GetxattrAddr()
	if validGetxattr() {
		*funcref = impl_Getxattr
	} else {
		*funcref = error_Getxattr
	}
	return (*funcref)(path, attr, dest)
}

func error_Getxattr(path string, attr string, dest []byte) (sz int, err error) {
	return -1, ENOSYS
}

func validGetxattr() bool {
	if funcptrtest(GetZosLibVec()+SYS___GETXATTR_A<<4, "") == 0 {
		if name, err := getLeFuncName(GetZosLibVec() + SYS___GETXATTR_A<<4); err == nil {
			return name == "__getxattr_a"
		}
	}
	return false
}




func impl_Setxattr(path string, attr string, data []byte, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(attr)
	if err != nil {
		return
	}
	var _p2 unsafe.Pointer
	if len(data) > 0 {
		_p2 = unsafe.Pointer(&data[0])
	} else {
		_p2 = unsafe.Pointer(&_zero)
	}
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___SETXATTR_A<<4, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(data)), uintptr(flags))
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_SetxattrAddr() *(func(path string, attr string, data []byte, flags int) (err error))

var Setxattr = enter_Setxattr

func enter_Setxattr(path string, attr string, data []byte, flags int) (err error) {
	funcref := get_SetxattrAddr()
	if validSetxattr() {
		*funcref = impl_Setxattr
	} else {
		*funcref = error_Setxattr
	}
	return (*funcref)(path, attr, data, flags)
}

func error_Setxattr(path string, attr string, data []byte, flags int) (err error) {
	return ENOSYS
}

func validSetxattr() bool {
	if funcptrtest(GetZosLibVec()+SYS___SETXATTR_A<<4, "") == 0 {
		if name, err := getLeFuncName(GetZosLibVec() + SYS___SETXATTR_A<<4); err == nil {
			return name == "__setxattr_a"
		}
	}
	return false
}






















//go:nosplit
func getPipe2Addr() *(func([]int, int) error)

var Pipe2 = pipe2Enter

func pipe2Enter(p []int, flags int) (err error) {
	if funcptrtest(GetZosLibVec()+SYS_PIPE2<<4, "") == 0 {
		*getPipe2Addr() = pipe2Impl
	} else {
		*getPipe2Addr() = pipe2Error
	}
	return (*getPipe2Addr())(p, flags)
}

func pipe2Impl(p []int, flags int) (err error) {
	var pp [2]_C_int
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_PIPE2<<4, uintptr(unsafe.Pointer(&pp[0])), uintptr(flags))
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	} else {
		p[0] = int(pp[0])
		p[1] = int(pp[1])
	}
	return
}
func pipe2Error(p []int, flags int) (err error) {
	return fmt.Errorf("Pipe2 is not available on this system")
}





func Readdir(dir uintptr) (dirent *Dirent, err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___READDIR_A<<4, uintptr(dir))
	runtime.ExitSyscall()
	dirent = (*Dirent)(unsafe.Pointer(r0))
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}















func Ptsname(fd int) (name string, err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS___PTSNAME_A<<4, uintptr(fd))
	runtime.ExitSyscall()
	if r0 == 0 {
		err = errnoErr2(e1, e2)
	} else {
		name = u2s(unsafe.Pointer(r0))
	}
	return
}

func u2s(cstr unsafe.Pointer) string {
	str := (*[1024]uint8)(cstr)
	i := 0
	for str[i] != 0 {
		i++
	}
	return string(str[:i])
}

func Close(fd int) (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_CLOSE<<4, uintptr(fd))
	runtime.ExitSyscall()
	for i := 0; e1 == EAGAIN && i < 10; i++ {
		runtime.EnterSyscall()
		CallLeFuncWithErr(GetZosLibVec()+SYS_USLEEP<<4, uintptr(10))
		runtime.ExitSyscall()
		runtime.EnterSyscall()
		r0, e2, e1 = CallLeFuncWithErr(GetZosLibVec()+SYS_CLOSE<<4, uintptr(fd))
		runtime.ExitSyscall()
	}
	if r0 != 0 {
		err = errnoErr2(e1, e2)
	}
	return
}


func Madvise(b []byte, advice int) (err error) {
	return
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func Munmap(b []byte) (err error) {
	return mapper.Munmap(b)
}

func MmapPtr(fd int, offset int64, addr unsafe.Pointer, length uintptr, prot int, flags int) (ret unsafe.Pointer, err error) {
	xaddr, err := mapper.mmap(uintptr(addr), length, prot, flags, fd, offset)
	return unsafe.Pointer(xaddr), err
}

func MunmapPtr(addr unsafe.Pointer, length uintptr) (err error) {
	return mapper.munmap(uintptr(addr), length)
}






func Getpgrp() (pid int) {
	pid, _ = Getpgid(0)
	return
}







func Getrusage(who int, rusage *Rusage) (err error) {
	var ruz rusage_zos
	err = getrusage(who, &ruz)
	
	rusage.Utime.Sec = ruz.Utime.Sec
	rusage.Utime.Usec = int64(ruz.Utime.Usec)
	rusage.Stime.Sec = ruz.Stime.Sec
	rusage.Stime.Usec = int64(ruz.Stime.Usec)
	return
}












func Lstat(path string, stat *Stat_t) (err error) {
	var statLE Stat_LE_t
	err = lstat(path, &statLE)
	copyStat(stat, &statLE)
	return
}


func isSpecialPath(path []byte) (v bool) {
	var special = [4][8]byte{
		{'V', 'E', 'R', 'S', 'I', 'O', 'N', '/'},
		{'S', 'Y', 'S', 'N', 'A', 'M', 'E', '/'},
		{'S', 'Y', 'S', 'S', 'Y', 'M', 'R', '/'},
		{'S', 'Y', 'S', 'S', 'Y', 'M', 'A', '/'}}

	var i, j int
	for i = 0; i < len(special); i++ {
		for j = 0; j < len(special[i]); j++ {
			if path[j] != special[i][j] {
				break
			}
		}
		if j == len(special[i]) {
			return true
		}
	}
	return false
}

func realpath(srcpath string, abspath []byte) (pathlen int, errno int) {
	var source [1024]byte
	copy(source[:], srcpath)
	source[len(srcpath)] = 0
	ret := runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___REALPATH_A<<4, 
		[]uintptr{uintptr(unsafe.Pointer(&source[0])),
			uintptr(unsafe.Pointer(&abspath[0]))})
	if ret != 0 {
		index := bytes.IndexByte(abspath[:], byte(0))
		if index != -1 {
			return index, 0
		}
	} else {
		errptr := (*int)(unsafe.Pointer(runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___ERRNO<<4, []uintptr{}))) 
		return 0, *errptr
	}
	return 0, 245 
}

func Readlink(path string, buf []byte) (n int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(buf) > 0 {
		_p1 = unsafe.Pointer(&buf[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	n = int(runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___READLINK_A<<4,
		[]uintptr{uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(buf))}))
	runtime.KeepAlive(unsafe.Pointer(_p0))
	if n == -1 {
		value := *(*int32)(unsafe.Pointer(runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___ERRNO<<4, []uintptr{})))
		err = errnoErr(Errno(value))
	} else {
		if buf[0] == '$' {
			if isSpecialPath(buf[1:9]) {
				cnt, err1 := realpath(path, buf)
				if err1 == 0 {
					n = cnt
				}
			}
		}
	}
	return
}

func impl_Readlinkat(dirfd int, path string, buf []byte) (n int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(buf) > 0 {
		_p1 = unsafe.Pointer(&buf[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___READLINKAT_A<<4, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(buf)))
	runtime.ExitSyscall()
	n = int(r0)
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
		return n, err
	} else {
		if buf[0] == '$' {
			if isSpecialPath(buf[1:9]) {
				cnt, err1 := realpath(path, buf)
				if err1 == 0 {
					n = cnt
				}
			}
		}
	}
	return
}

//go:nosplit
func get_ReadlinkatAddr() *(func(dirfd int, path string, buf []byte) (n int, err error))

var Readlinkat = enter_Readlinkat

func enter_Readlinkat(dirfd int, path string, buf []byte) (n int, err error) {
	funcref := get_ReadlinkatAddr()
	if funcptrtest(GetZosLibVec()+SYS___READLINKAT_A<<4, "") == 0 {
		*funcref = impl_Readlinkat
	} else {
		*funcref = error_Readlinkat
	}
	return (*funcref)(dirfd, path, buf)
}

func error_Readlinkat(dirfd int, path string, buf []byte) (n int, err error) {
	n = -1
	err = ENOSYS
	return
}































func Stat(path string, sta *Stat_t) (err error) {
	var statLE Stat_LE_t
	err = stat(path, &statLE)
	copyStat(sta, &statLE)
	return
}














func Open(path string, mode int, perm uint32) (fd int, err error) {
	if mode&O_ACCMODE == 0 {
		mode |= O_RDONLY
	}
	return open(path, mode, perm)
}



func Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) {
	if flags&O_ACCMODE == 0 {
		flags |= O_RDONLY
	}
	return openat(dirfd, path, flags, mode)
}



func Openat2(dirfd int, path string, how *OpenHow) (fd int, err error) {
	if how.Flags&O_ACCMODE == 0 {
		how.Flags |= O_RDONLY
	}
	return openat2(dirfd, path, how, SizeofOpenHow)
}

func ZosFdToPath(dirfd int) (path string, err error) {
	var buffer [1024]byte
	runtime.EnterSyscall()
	ret, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_W_IOCTL<<4, uintptr(dirfd), 17, 1024, uintptr(unsafe.Pointer(&buffer[0])))
	runtime.ExitSyscall()
	if ret == 0 {
		zb := bytes.IndexByte(buffer[:], 0)
		if zb == -1 {
			zb = len(buffer)
		}
		CallLeFuncWithErr(GetZosLibVec()+SYS___E2A_L<<4, uintptr(unsafe.Pointer(&buffer[0])), uintptr(zb))
		return string(buffer[:zb]), nil
	}
	return "", errnoErr2(e1, e2)
}



func Remove(path string) error {
	return remove(path)
}

const ImplementsGetwd = true

func Getcwd(buf []byte) (n int, err error) {
	var p unsafe.Pointer
	if len(buf) > 0 {
		p = unsafe.Pointer(&buf[0])
	} else {
		p = unsafe.Pointer(&_zero)
	}
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS___GETCWD_A<<4, uintptr(p), uintptr(len(buf)))
	runtime.ExitSyscall()
	n = clen(buf) + 1
	if r0 == 0 {
		err = errnoErr2(e1, e2)
	}
	return
}

func Getwd() (wd string, err error) {
	var buf [PathMax]byte
	n, err := Getcwd(buf[0:])
	if err != nil {
		return "", err
	}
	
	if n < 1 || n > len(buf) || buf[n-1] != 0 {
		return "", EINVAL
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

func gettid() uint64

func Gettid() (tid int) {
	return int(gettid())
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

func (w WaitStatus) Signal() Signal {
	if !w.Signaled() {
		return -1
	}
	return Signal(w & mask)
}

func (w WaitStatus) StopSignal() Signal {
	if !w.Stopped() {
		return -1
	}
	return Signal(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int { return -1 }



func Waitid(idType int, id int, info *Siginfo, options int, rusage *Rusage) (err error) {
	return waitid(idType, id, info, options)
}



func impl_Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_WAIT4<<4, uintptr(pid), uintptr(unsafe.Pointer(wstatus)), uintptr(options), uintptr(unsafe.Pointer(rusage)))
	runtime.ExitSyscall()
	wpid = int(r0)
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_Wait4Addr() *(func(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error))

var Wait4 = enter_Wait4

func enter_Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	funcref := get_Wait4Addr()
	if funcptrtest(GetZosLibVec()+SYS_WAIT4<<4, "") == 0 {
		*funcref = impl_Wait4
	} else {
		*funcref = legacyWait4
	}
	return (*funcref)(pid, wstatus, options, rusage)
}

func legacyWait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	
	
	var status _C_int
	wpid, err = waitpid(pid, &status, options)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return
}



func Gettimeofday(tv *Timeval) (err error) {
	var tvz timeval_zos
	err = gettimeofday(&tvz)
	tv.Sec = tvz.Sec
	tv.Usec = int64(tvz.Usec)
	return
}

func Time(t *Time_t) (tt Time_t, err error) {
	var tv Timeval
	err = Gettimeofday(&tv)
	if err != nil {
		return 0, err
	}
	if t != nil {
		*t = Time_t(tv.Sec)
	}
	return Time_t(tv.Sec), nil
}

func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval { 
	return Timeval{Sec: sec, Usec: usec}
}



func Pipe(p []int) (err error) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	err = pipe(&pp)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
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



func validUtimensat() bool {
	if funcptrtest(GetZosLibVec()+SYS___UTIMENSAT_A<<4, "") == 0 {
		if name, err := getLeFuncName(GetZosLibVec() + SYS___UTIMENSAT_A<<4); err == nil {
			return name == "__utimensat_a"
		}
	}
	return false
}



//go:nosplit
func get_UtimesNanoAddr() *(func(path string, ts []Timespec) (err error))

var UtimesNano = enter_UtimesNano

func enter_UtimesNano(path string, ts []Timespec) (err error) {
	funcref := get_UtimesNanoAddr()
	if validUtimensat() {
		*funcref = utimesNanoImpl
	} else {
		*funcref = legacyUtimesNano
	}
	return (*funcref)(path, ts)
}

func utimesNanoImpl(path string, ts []Timespec) (err error) {
	if ts == nil {
		return utimensat(AT_FDCWD, path, nil, 0)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
}

func legacyUtimesNano(path string, ts []Timespec) (err error) {
	if len(ts) != 2 {
		return EINVAL
	}
	
	
	tv := [2]Timeval{
		NsecToTimeval(TimespecToNsec(ts[0])),
		NsecToTimeval(TimespecToNsec(ts[1])),
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}





//go:nosplit
func get_UtimesNanoAtAddr() *(func(dirfd int, path string, ts []Timespec, flags int) (err error))

var UtimesNanoAt = enter_UtimesNanoAt

func enter_UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) (err error) {
	funcref := get_UtimesNanoAtAddr()
	if validUtimensat() {
		*funcref = utimesNanoAtImpl
	} else {
		*funcref = legacyUtimesNanoAt
	}
	return (*funcref)(dirfd, path, ts, flags)
}

func utimesNanoAtImpl(dirfd int, path string, ts []Timespec, flags int) (err error) {
	if ts == nil {
		return utimensat(dirfd, path, nil, flags)
	}
	if len(ts) != 2 {
		return EINVAL
	}
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
}

func legacyUtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) (err error) {
	if path[0] != '/' {
		dirPath, err := ZosFdToPath(dirfd)
		if err != nil {
			return err
		}
		path = dirPath + "/" + path
	}
	if flags == AT_SYMLINK_NOFOLLOW {
		if len(ts) != 2 {
			return EINVAL
		}

		if ts[0].Nsec >= 5e8 {
			ts[0].Sec++
		}
		ts[0].Nsec = 0
		if ts[1].Nsec >= 5e8 {
			ts[1].Sec++
		}
		ts[1].Nsec = 0

		
		
		tv := []Timeval{
			NsecToTimeval(TimespecToNsec(ts[0])),
			NsecToTimeval(TimespecToNsec(ts[1])),
		}
		return Lutimes(path, tv)
	}
	return UtimesNano(path, ts)
}



func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	
	return anyToSockaddr(0, &rsa)
}

const (
	
	nwmHeaderIdentifier    = 0xd5e6d4c8
	nwmFilterIdentifier    = 0xd5e6d4c6
	nwmTCPConnIdentifier   = 0xd5e6d4c3
	nwmRecHeaderIdentifier = 0xd5e6d4d9
	nwmIPStatsIdentifier   = 0xd5e6d4c9d7e2e340
	nwmIPGStatsIdentifier  = 0xd5e6d4c9d7c7e2e3
	nwmTCPStatsIdentifier  = 0xd5e6d4e3c3d7e2e3
	nwmUDPStatsIdentifier  = 0xd5e6d4e4c4d7e2e3
	nwmICMPGStatsEntry     = 0xd5e6d4c9c3d4d7c7
	nwmICMPTStatsEntry     = 0xd5e6d4c9c3d4d7e3

	
	nwmVersion1   = 1
	nwmVersion2   = 2
	nwmCurrentVer = 2

	nwmTCPConnType     = 1
	nwmGlobalStatsType = 14

	
	nwmFilterLclAddrMask = 0x20000000 
	nwmFilterSrcAddrMask = 0x20000000 
	nwmFilterLclPortMask = 0x10000000 
	nwmFilterSrcPortMask = 0x10000000 

	
	nwmTCPStateClosed   = 1
	nwmTCPStateListen   = 2
	nwmTCPStateSynSent  = 3
	nwmTCPStateSynRcvd  = 4
	nwmTCPStateEstab    = 5
	nwmTCPStateFinWait1 = 6
	nwmTCPStateFinWait2 = 7
	nwmTCPStateClosWait = 8
	nwmTCPStateLastAck  = 9
	nwmTCPStateClosing  = 10
	nwmTCPStateTimeWait = 11
	nwmTCPStateDeletTCB = 12

	
	BPF_TCP_CLOSE        = 1
	BPF_TCP_LISTEN       = 2
	BPF_TCP_SYN_SENT     = 3
	BPF_TCP_SYN_RECV     = 4
	BPF_TCP_ESTABLISHED  = 5
	BPF_TCP_FIN_WAIT1    = 6
	BPF_TCP_FIN_WAIT2    = 7
	BPF_TCP_CLOSE_WAIT   = 8
	BPF_TCP_LAST_ACK     = 9
	BPF_TCP_CLOSING      = 10
	BPF_TCP_TIME_WAIT    = 11
	BPF_TCP_NEW_SYN_RECV = -1
	BPF_TCP_MAX_STATES   = -2
)

type nwmTriplet struct {
	offset uint32
	length uint32
	number uint32
}

type nwmQuadruplet struct {
	offset uint32
	length uint32
	number uint32
	match  uint32
}

type nwmHeader struct {
	ident       uint32
	length      uint32
	version     uint16
	nwmType     uint16
	bytesNeeded uint32
	options     uint32
	_           [16]byte
	inputDesc   nwmTriplet
	outputDesc  nwmQuadruplet
}

type nwmFilter struct {
	ident         uint32
	flags         uint32
	resourceName  [8]byte
	resourceId    uint32
	listenerId    uint32
	local         [28]byte 
	remote        [28]byte 
	_             uint16
	_             uint16
	asid          uint16
	_             [2]byte
	tnLuName      [8]byte
	tnMonGrp      uint32
	tnAppl        [8]byte
	applData      [40]byte
	nInterface    [16]byte
	dVipa         [16]byte
	dVipaPfx      uint16
	dVipaPort     uint16
	dVipaFamily   byte
	_             [3]byte
	destXCF       [16]byte
	destXCFPfx    uint16
	destXCFFamily byte
	_             [1]byte
	targIP        [16]byte
	targIPPfx     uint16
	targIPFamily  byte
	_             [1]byte
	_             [20]byte
}

type nwmRecHeader struct {
	ident  uint32
	length uint32
	number byte
	_      [3]byte
}

type nwmTCPStatsEntry struct {
	ident             uint64
	currEstab         uint32
	activeOpened      uint32
	passiveOpened     uint32
	connClosed        uint32
	estabResets       uint32
	attemptFails      uint32
	passiveDrops      uint32
	timeWaitReused    uint32
	inSegs            uint64
	predictAck        uint32
	predictData       uint32
	inDupAck          uint32
	inBadSum          uint32
	inBadLen          uint32
	inShort           uint32
	inDiscOldTime     uint32
	inAllBeforeWin    uint32
	inSomeBeforeWin   uint32
	inAllAfterWin     uint32
	inSomeAfterWin    uint32
	inOutOfOrder      uint32
	inAfterClose      uint32
	inWinProbes       uint32
	inWinUpdates      uint32
	outWinUpdates     uint32
	outSegs           uint64
	outDelayAcks      uint32
	outRsts           uint32
	retransSegs       uint32
	retransTimeouts   uint32
	retransDrops      uint32
	pmtuRetrans       uint32
	pmtuErrors        uint32
	outWinProbes      uint32
	probeDrops        uint32
	keepAliveProbes   uint32
	keepAliveDrops    uint32
	finwait2Drops     uint32
	acceptCount       uint64
	inBulkQSegs       uint64
	inDiscards        uint64
	connFloods        uint32
	connStalls        uint32
	cfgEphemDef       uint16
	ephemInUse        uint16
	ephemHiWater      uint16
	flags             byte
	_                 [1]byte
	ephemExhaust      uint32
	smcRCurrEstabLnks uint32
	smcRLnkActTimeOut uint32
	smcRActLnkOpened  uint32
	smcRPasLnkOpened  uint32
	smcRLnksClosed    uint32
	smcRCurrEstab     uint32
	smcRActiveOpened  uint32
	smcRPassiveOpened uint32
	smcRConnClosed    uint32
	smcRInSegs        uint64
	smcROutSegs       uint64
	smcRInRsts        uint32
	smcROutRsts       uint32
	smcDCurrEstabLnks uint32
	smcDActLnkOpened  uint32
	smcDPasLnkOpened  uint32
	smcDLnksClosed    uint32
	smcDCurrEstab     uint32
	smcDActiveOpened  uint32
	smcDPassiveOpened uint32
	smcDConnClosed    uint32
	smcDInSegs        uint64
	smcDOutSegs       uint64
	smcDInRsts        uint32
	smcDOutRsts       uint32
}

type nwmConnEntry struct {
	ident             uint32
	local             [28]byte 
	remote            [28]byte 
	startTime         [8]byte  
	lastActivity      [8]byte  
	bytesIn           [8]byte  
	bytesOut          [8]byte  
	inSegs            [8]byte  
	outSegs           [8]byte  
	state             uint16
	activeOpen        byte
	flag01            byte
	outBuffered       uint32
	inBuffered        uint32
	maxSndWnd         uint32
	reXmtCount        uint32
	congestionWnd     uint32
	ssThresh          uint32
	roundTripTime     uint32
	roundTripVar      uint32
	sendMSS           uint32
	sndWnd            uint32
	rcvBufSize        uint32
	sndBufSize        uint32
	outOfOrderCount   uint32
	lcl0WindowCount   uint32
	rmt0WindowCount   uint32
	dupacks           uint32
	flag02            byte
	sockOpt6Cont      byte
	asid              uint16
	resourceName      [8]byte
	resourceId        uint32
	subtask           uint32
	sockOpt           byte
	sockOpt6          byte
	clusterConnFlag   byte
	proto             byte
	targetAppl        [8]byte
	luName            [8]byte
	clientUserId      [8]byte
	logMode           [8]byte
	timeStamp         uint32
	timeStampAge      uint32
	serverResourceId  uint32
	intfName          [16]byte
	ttlsStatPol       byte
	ttlsStatConn      byte
	ttlsSSLProt       uint16
	ttlsNegCiph       [2]byte
	ttlsSecType       byte
	ttlsFIPS140Mode   byte
	ttlsUserID        [8]byte
	applData          [40]byte
	inOldestTime      [8]byte 
	outOldestTime     [8]byte 
	tcpTrustedPartner byte
	_                 [3]byte
	bulkDataIntfName  [16]byte
	ttlsNegCiph4      [4]byte
	smcReason         uint32
	lclSMCLinkId      uint32
	rmtSMCLinkId      uint32
	smcStatus         byte
	smcFlags          byte
	_                 [2]byte
	rcvWnd            uint32
	lclSMCBufSz       uint32
	rmtSMCBufSz       uint32
	ttlsSessID        [32]byte
	ttlsSessIDLen     int16
	_                 [1]byte
	smcDStatus        byte
	smcDReason        uint32
}

var svcNameTable [][]byte = [][]byte{
	[]byte("\xc5\xe9\xc2\xd5\xd4\xc9\xc6\xf4"), 
}

const (
	svc_EZBNMIF4 = 0
)

func GetsockoptTCPInfo(fd, level, opt int) (*TCPInfo, error) {
	jobname := []byte("\x5c\x40\x40\x40\x40\x40\x40\x40") 
	responseBuffer := [4096]byte{0}
	var bufferAlet, reasonCode uint32 = 0, 0
	var bufferLen, returnValue, returnCode int32 = 4096, 0, 0

	dsa := [18]uint64{0}
	var argv [7]unsafe.Pointer
	argv[0] = unsafe.Pointer(&jobname[0])
	argv[1] = unsafe.Pointer(&responseBuffer[0])
	argv[2] = unsafe.Pointer(&bufferAlet)
	argv[3] = unsafe.Pointer(&bufferLen)
	argv[4] = unsafe.Pointer(&returnValue)
	argv[5] = unsafe.Pointer(&returnCode)
	argv[6] = unsafe.Pointer(&reasonCode)

	request := (*struct {
		header nwmHeader
		filter nwmFilter
	})(unsafe.Pointer(&responseBuffer[0]))

	EZBNMIF4 := svcLoad(&svcNameTable[svc_EZBNMIF4][0])
	if EZBNMIF4 == nil {
		return nil, errnoErr(EINVAL)
	}

	
	request.header.ident = nwmHeaderIdentifier
	request.header.length = uint32(unsafe.Sizeof(request.header))
	request.header.version = nwmCurrentVer
	request.header.nwmType = nwmGlobalStatsType
	request.header.options = 0x80000000

	svcCall(EZBNMIF4, &argv[0], &dsa[0])

	
	if returnCode != 0 || request.header.outputDesc.offset == 0 {
		return nil, errnoErr(EINVAL)
	}

	
	recHeader := (*nwmRecHeader)(unsafe.Pointer(&responseBuffer[request.header.outputDesc.offset]))
	if recHeader.ident != nwmRecHeaderIdentifier {
		return nil, errnoErr(EINVAL)
	}

	
	var sections []*uint64
	var sectionDesc *nwmTriplet = (*nwmTriplet)(unsafe.Pointer(&responseBuffer[0]))
	for i := uint32(0); i < uint32(recHeader.number); i++ {
		offset := request.header.outputDesc.offset + uint32(unsafe.Sizeof(*recHeader)) + i*uint32(unsafe.Sizeof(*sectionDesc))
		sectionDesc = (*nwmTriplet)(unsafe.Pointer(&responseBuffer[offset]))
		for j := uint32(0); j < sectionDesc.number; j++ {
			offset = request.header.outputDesc.offset + sectionDesc.offset + j*sectionDesc.length
			sections = append(sections, (*uint64)(unsafe.Pointer(&responseBuffer[offset])))
		}
	}

	
	var tcpStats *nwmTCPStatsEntry = nil
	for _, ptr := range sections {
		switch *ptr {
		case nwmTCPStatsIdentifier:
			if tcpStats != nil {
				return nil, errnoErr(EINVAL)
			}
			tcpStats = (*nwmTCPStatsEntry)(unsafe.Pointer(ptr))
		case nwmIPStatsIdentifier:
		case nwmIPGStatsIdentifier:
		case nwmUDPStatsIdentifier:
		case nwmICMPGStatsEntry:
		case nwmICMPTStatsEntry:
		default:
			return nil, errnoErr(EINVAL)
		}
	}
	if tcpStats == nil {
		return nil, errnoErr(EINVAL)
	}

	
	responseBuffer = [4096]byte{0}
	dsa = [18]uint64{0}
	bufferAlet, reasonCode = 0, 0
	bufferLen, returnValue, returnCode = 4096, 0, 0
	nameptr := (*uint32)(unsafe.Pointer(uintptr(0x21c))) 
	nameptr = (*uint32)(unsafe.Pointer(uintptr(*nameptr + 12)))
	argv[0] = unsafe.Pointer(uintptr(*nameptr))

	request.header.ident = nwmHeaderIdentifier
	request.header.length = uint32(unsafe.Sizeof(request.header))
	request.header.version = nwmCurrentVer
	request.header.nwmType = nwmTCPConnType
	request.header.options = 0x80000000

	request.filter.ident = nwmFilterIdentifier

	var localSockaddr RawSockaddrAny
	socklen := _Socklen(SizeofSockaddrAny)
	err := getsockname(fd, &localSockaddr, &socklen)
	if err != nil {
		return nil, errnoErr(EINVAL)
	}
	if localSockaddr.Addr.Family == AF_INET {
		localSockaddr := (*RawSockaddrInet4)(unsafe.Pointer(&localSockaddr.Addr))
		localSockFilter := (*RawSockaddrInet4)(unsafe.Pointer(&request.filter.local[0]))
		localSockFilter.Family = AF_INET
		var i int
		for i = 0; i < 4; i++ {
			if localSockaddr.Addr[i] != 0 {
				break
			}
		}
		if i != 4 {
			request.filter.flags |= nwmFilterLclAddrMask
			for i = 0; i < 4; i++ {
				localSockFilter.Addr[i] = localSockaddr.Addr[i]
			}
		}
		if localSockaddr.Port != 0 {
			request.filter.flags |= nwmFilterLclPortMask
			localSockFilter.Port = localSockaddr.Port
		}
	} else if localSockaddr.Addr.Family == AF_INET6 {
		localSockaddr := (*RawSockaddrInet6)(unsafe.Pointer(&localSockaddr.Addr))
		localSockFilter := (*RawSockaddrInet6)(unsafe.Pointer(&request.filter.local[0]))
		localSockFilter.Family = AF_INET6
		var i int
		for i = 0; i < 16; i++ {
			if localSockaddr.Addr[i] != 0 {
				break
			}
		}
		if i != 16 {
			request.filter.flags |= nwmFilterLclAddrMask
			for i = 0; i < 16; i++ {
				localSockFilter.Addr[i] = localSockaddr.Addr[i]
			}
		}
		if localSockaddr.Port != 0 {
			request.filter.flags |= nwmFilterLclPortMask
			localSockFilter.Port = localSockaddr.Port
		}
	}

	svcCall(EZBNMIF4, &argv[0], &dsa[0])

	
	if returnCode != 0 || request.header.outputDesc.offset == 0 {
		return nil, errnoErr(EINVAL)
	}

	
	conn := (*nwmConnEntry)(unsafe.Pointer(&responseBuffer[request.header.outputDesc.offset]))
	if conn.ident != nwmTCPConnIdentifier {
		return nil, errnoErr(EINVAL)
	}

	
	
	
	
	var tcpinfo TCPInfo
	tcpinfo.State = uint8(conn.state)
	tcpinfo.Ca_state = 0 
	tcpinfo.Retransmits = uint8(tcpStats.retransSegs)
	tcpinfo.Probes = uint8(tcpStats.outWinProbes)
	tcpinfo.Backoff = 0 
	tcpinfo.Options = 0 
	tcpinfo.Rto = tcpStats.retransTimeouts
	tcpinfo.Ato = tcpStats.outDelayAcks
	tcpinfo.Snd_mss = conn.sendMSS
	tcpinfo.Rcv_mss = conn.sendMSS 
	tcpinfo.Unacked = 0            
	tcpinfo.Sacked = 0             
	tcpinfo.Lost = 0               
	tcpinfo.Retrans = conn.reXmtCount
	tcpinfo.Fackets = 0 
	tcpinfo.Last_data_sent = uint32(*(*uint64)(unsafe.Pointer(&conn.lastActivity[0])))
	tcpinfo.Last_ack_sent = uint32(*(*uint64)(unsafe.Pointer(&conn.outOldestTime[0])))
	tcpinfo.Last_data_recv = uint32(*(*uint64)(unsafe.Pointer(&conn.inOldestTime[0])))
	tcpinfo.Last_ack_recv = uint32(*(*uint64)(unsafe.Pointer(&conn.inOldestTime[0])))
	tcpinfo.Pmtu = conn.sendMSS 
	tcpinfo.Rcv_ssthresh = conn.ssThresh
	tcpinfo.Rtt = conn.roundTripTime
	tcpinfo.Rttvar = conn.roundTripVar
	tcpinfo.Snd_ssthresh = conn.ssThresh 
	tcpinfo.Snd_cwnd = conn.congestionWnd
	tcpinfo.Advmss = conn.sendMSS        
	tcpinfo.Reordering = 0               
	tcpinfo.Rcv_rtt = conn.roundTripTime 
	tcpinfo.Rcv_space = conn.sendMSS     
	tcpinfo.Total_retrans = conn.reXmtCount

	svcUnload(&svcNameTable[svc_EZBNMIF4][0], EZBNMIF4)

	return &tcpinfo, nil
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

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) {
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = SizeofSockaddrAny
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	
	if rsa.Addr.Family != AF_UNSPEC {
		
		from, err = anyToSockaddr(0, &rsa)
	}
	return
}

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error) {
	_, err = SendmsgN(fd, p, oob, to, flags)
	return
}

func SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error) {
	var ptr unsafe.Pointer
	var salen _Socklen
	if to != nil {
		var err error
		ptr, salen, err = to.sockaddr()
		if err != nil {
			return 0, err
		}
	}
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(ptr))
	msg.Namelen = int32(salen)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil {
		return 0, err
	}
	if len(oob) > 0 && len(p) == 0 {
		n = 0
	}
	return n, nil
}

func Opendir(name string) (uintptr, error) {
	p, err := BytePtrFromString(name)
	if err != nil {
		return 0, err
	}
	err = nil
	runtime.EnterSyscall()
	dir, e2, e1 := CallLeFuncWithPtrReturn(GetZosLibVec()+SYS___OPENDIR_A<<4, uintptr(unsafe.Pointer(p)))
	runtime.ExitSyscall()
	runtime.KeepAlive(unsafe.Pointer(p))
	if dir == 0 {
		err = errnoErr2(e1, e2)
	}
	return dir, err
}


func clearErrno()

func Closedir(dir uintptr) error {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_CLOSEDIR<<4, dir)
	runtime.ExitSyscall()
	if r0 != 0 {
		return errnoErr2(e1, e2)
	}
	return nil
}

func Seekdir(dir uintptr, pos int) {
	runtime.EnterSyscall()
	CallLeFuncWithErr(GetZosLibVec()+SYS_SEEKDIR<<4, dir, uintptr(pos))
	runtime.ExitSyscall()
}

func Telldir(dir uintptr) (int, error) {
	p, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_TELLDIR<<4, dir)
	pos := int(p)
	if int64(p) == -1 {
		return pos, errnoErr2(e1, e2)
	}
	return pos, nil
}


func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	
	
	var flock [24]byte
	*(*int16)(unsafe.Pointer(&flock[0])) = lk.Type
	*(*int16)(unsafe.Pointer(&flock[2])) = lk.Whence
	*(*int64)(unsafe.Pointer(&flock[4])) = lk.Start
	*(*int64)(unsafe.Pointer(&flock[12])) = lk.Len
	*(*int32)(unsafe.Pointer(&flock[20])) = lk.Pid
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_FCNTL<<4, fd, uintptr(cmd), uintptr(unsafe.Pointer(&flock)))
	runtime.ExitSyscall()
	lk.Type = *(*int16)(unsafe.Pointer(&flock[0]))
	lk.Whence = *(*int16)(unsafe.Pointer(&flock[2]))
	lk.Start = *(*int64)(unsafe.Pointer(&flock[4]))
	lk.Len = *(*int64)(unsafe.Pointer(&flock[12]))
	lk.Pid = *(*int32)(unsafe.Pointer(&flock[20]))
	if r0 == 0 {
		return nil
	}
	return errnoErr2(e1, e2)
}

func impl_Flock(fd int, how int) (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_FLOCK<<4, uintptr(fd), uintptr(how))
	runtime.ExitSyscall()
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_FlockAddr() *(func(fd int, how int) (err error))

var Flock = enter_Flock

func validFlock(fp uintptr) bool {
	if funcptrtest(GetZosLibVec()+SYS_FLOCK<<4, "") == 0 {
		if name, err := getLeFuncName(GetZosLibVec() + SYS_FLOCK<<4); err == nil {
			return name == "flock"
		}
	}
	return false
}

func enter_Flock(fd int, how int) (err error) {
	funcref := get_FlockAddr()
	if validFlock(GetZosLibVec() + SYS_FLOCK<<4) {
		*funcref = impl_Flock
	} else {
		*funcref = legacyFlock
	}
	return (*funcref)(fd, how)
}

func legacyFlock(fd int, how int) error {

	var flock_type int16
	var fcntl_cmd int

	switch how {
	case LOCK_SH | LOCK_NB:
		flock_type = F_RDLCK
		fcntl_cmd = F_SETLK
	case LOCK_EX | LOCK_NB:
		flock_type = F_WRLCK
		fcntl_cmd = F_SETLK
	case LOCK_EX:
		flock_type = F_WRLCK
		fcntl_cmd = F_SETLKW
	case LOCK_UN:
		flock_type = F_UNLCK
		fcntl_cmd = F_SETLKW
	default:
	}

	flock := Flock_t{
		Type:   int16(flock_type),
		Whence: int16(0),
		Start:  int64(0),
		Len:    int64(0),
		Pid:    int32(Getppid()),
	}

	err := FcntlFlock(uintptr(fd), fcntl_cmd, &flock)
	return err
}

func Mlock(b []byte) (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MLOCKALL<<4, _BPX_NONSWAP)
	runtime.ExitSyscall()
	if r0 != 0 {
		err = errnoErr2(e1, e2)
	}
	return
}

func Mlock2(b []byte, flags int) (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MLOCKALL<<4, _BPX_NONSWAP)
	runtime.ExitSyscall()
	if r0 != 0 {
		err = errnoErr2(e1, e2)
	}
	return
}

func Mlockall(flags int) (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MLOCKALL<<4, _BPX_NONSWAP)
	runtime.ExitSyscall()
	if r0 != 0 {
		err = errnoErr2(e1, e2)
	}
	return
}

func Munlock(b []byte) (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MLOCKALL<<4, _BPX_SWAP)
	runtime.ExitSyscall()
	if r0 != 0 {
		err = errnoErr2(e1, e2)
	}
	return
}

func Munlockall() (err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MLOCKALL<<4, _BPX_SWAP)
	runtime.ExitSyscall()
	if r0 != 0 {
		err = errnoErr2(e1, e2)
	}
	return
}

func ClockGettime(clockid int32, ts *Timespec) error {

	var ticks_per_sec uint32 = 100 
	var nsec_per_sec int64 = 1000000000

	if ts == nil {
		return EFAULT
	}
	if clockid == CLOCK_REALTIME || clockid == CLOCK_MONOTONIC {
		var nanotime int64 = runtime.Nanotime1()
		ts.Sec = nanotime / nsec_per_sec
		ts.Nsec = nanotime % nsec_per_sec
	} else if clockid == CLOCK_PROCESS_CPUTIME_ID || clockid == CLOCK_THREAD_CPUTIME_ID {
		var tm Tms
		_, err := Times(&tm)
		if err != nil {
			return EFAULT
		}
		ts.Sec = int64(tm.Utime / ticks_per_sec)
		ts.Nsec = int64(tm.Utime) * nsec_per_sec / int64(ticks_per_sec)
	} else {
		return EINVAL
	}
	return nil
}



//go:nosplit
func get_ChtagAddr() *(func(path string, ccsid uint64, textbit uint64) error)

var Chtag = enter_Chtag

func enter_Chtag(path string, ccsid uint64, textbit uint64) error {
	funcref := get_ChtagAddr()
	if validSetxattr() {
		*funcref = impl_Chtag
	} else {
		*funcref = legacy_Chtag
	}
	return (*funcref)(path, ccsid, textbit)
}

func legacy_Chtag(path string, ccsid uint64, textbit uint64) error {
	tag := ccsid<<16 | textbit<<15
	var tag_buff [8]byte
	DecodeData(tag_buff[:], 8, tag)
	return Setxattr(path, "filetag", tag_buff[:], XATTR_REPLACE)
}

func impl_Chtag(path string, ccsid uint64, textbit uint64) error {
	tag := ccsid<<16 | textbit<<15
	var tag_buff [4]byte
	DecodeData(tag_buff[:], 4, tag)
	return Setxattr(path, "system.filetag", tag_buff[:], XATTR_REPLACE)
}





//go:nosplit
func get_NanosleepAddr() *(func(time *Timespec, leftover *Timespec) error)

var Nanosleep = enter_Nanosleep

func enter_Nanosleep(time *Timespec, leftover *Timespec) error {
	funcref := get_NanosleepAddr()
	if funcptrtest(GetZosLibVec()+SYS_NANOSLEEP<<4, "") == 0 {
		*funcref = impl_Nanosleep
	} else {
		*funcref = legacyNanosleep
	}
	return (*funcref)(time, leftover)
}

func impl_Nanosleep(time *Timespec, leftover *Timespec) error {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_NANOSLEEP<<4, uintptr(unsafe.Pointer(time)), uintptr(unsafe.Pointer(leftover)))
	runtime.ExitSyscall()
	if int64(r0) == -1 {
		return errnoErr2(e1, e2)
	}
	return nil
}

func legacyNanosleep(time *Timespec, leftover *Timespec) error {
	t0 := runtime.Nanotime1()
	var secrem uint32
	var nsecrem uint32
	total := time.Sec*1000000000 + time.Nsec
	elapsed := runtime.Nanotime1() - t0
	var rv int32
	var rc int32
	var err error
	
	for total-elapsed > 1000000000 {
		rv, rc, _ = BpxCondTimedWait(uint32(1), uint32(0), uint32(CW_CONDVAR), &secrem, &nsecrem)
		if rv != 0 && rc != 112 { 
			if leftover != nil && rc == 120 { 
				leftover.Sec = int64(secrem)
				leftover.Nsec = int64(nsecrem)
			}
			err = Errno(rc)
			return err
		}
		elapsed = runtime.Nanotime1() - t0
	}
	
	if total > elapsed {
		rv, rc, _ = BpxCondTimedWait(uint32(0), uint32(total-elapsed), uint32(CW_CONDVAR), &secrem, &nsecrem)
	}
	if leftover != nil && rc == 120 {
		leftover.Sec = int64(secrem)
		leftover.Nsec = int64(nsecrem)
	}
	if rv != 0 && rc != 112 {
		err = Errno(rc)
	}
	return err
}



var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)



var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
)

var ZosTraceLevel int
var ZosTracefile *os.File

var (
	signalNameMapOnce sync.Once
	signalNameMap     map[string]syscall.Signal
)



func errnoErr(e Errno) error {
	switch e {
	case 0:
		return nil
	case EAGAIN:
		return errEAGAIN
	case EINVAL:
		return errEINVAL
	case ENOENT:
		return errENOENT
	}
	return e
}

var reg *regexp.Regexp


func errnoErr2(e Errno, e2 uintptr) error {
	switch e {
	case 0:
		return nil
	case EAGAIN:
		return errEAGAIN
		
	}
	if ZosTraceLevel > 0 {
		var name string
		if reg == nil {
			reg = regexp.MustCompile("(^unix\\.[^/]+$|.*\\/unix\\.[^/]+$)")
		}
		i := 1
		pc, file, line, ok := runtime.Caller(i)
		if ok {
			name = runtime.FuncForPC(pc).Name()
		}
		for ok && reg.MatchString(runtime.FuncForPC(pc).Name()) {
			i += 1
			pc, file, line, ok = runtime.Caller(i)
		}
		if ok {
			if ZosTracefile == nil {
				ZosConsolePrintf("From %s:%d\n", file, line)
				ZosConsolePrintf("%s: %s (errno2=0x%x)\n", name, e.Error(), e2)
			} else {
				fmt.Fprintf(ZosTracefile, "From %s:%d\n", file, line)
				fmt.Fprintf(ZosTracefile, "%s: %s (errno2=0x%x)\n", name, e.Error(), e2)
			}
		} else {
			if ZosTracefile == nil {
				ZosConsolePrintf("%s (errno2=0x%x)\n", e.Error(), e2)
			} else {
				fmt.Fprintf(ZosTracefile, "%s (errno2=0x%x)\n", e.Error(), e2)
			}
		}
	}
	return e
}


func ErrnoName(e Errno) string {
	i := sort.Search(len(errorList), func(i int) bool {
		return errorList[i].num >= e
	})
	if i < len(errorList) && errorList[i].num == e {
		return errorList[i].name
	}
	return ""
}


func SignalName(s syscall.Signal) string {
	i := sort.Search(len(signalList), func(i int) bool {
		return signalList[i].num >= s
	})
	if i < len(signalList) && signalList[i].num == s {
		return signalList[i].name
	}
	return ""
}




func SignalNum(s string) syscall.Signal {
	signalNameMapOnce.Do(func() {
		signalNameMap = make(map[string]syscall.Signal, len(signalList))
		for _, signal := range signalList {
			signalNameMap[signal.name] = signal.num
		}
	})
	return signalNameMap[s]
}


func clen(n []byte) int {
	i := bytes.IndexByte(n, 0)
	if i == -1 {
		i = len(n)
	}
	return i
}



type mmapper struct {
	sync.Mutex
	active map[*byte][]byte 
	mmap   func(addr, length uintptr, prot, flags, fd int, offset int64) (uintptr, error)
	munmap func(addr uintptr, length uintptr) error
}

func (m *mmapper) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	if length <= 0 {
		return nil, EINVAL
	}

	
	flags |= __MAP_64

	
	addr, errno := m.mmap(0, uintptr(length), prot, flags, fd, offset)
	if errno != nil {
		return nil, errno
	}

	
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{addr, length, length}

	
	b := *(*[]byte)(unsafe.Pointer(&sl))

	
	p := &b[cap(b)-1]
	m.Lock()
	defer m.Unlock()
	m.active[p] = b
	return b, nil
}

func (m *mmapper) Munmap(data []byte) (err error) {
	if len(data) == 0 || len(data) != cap(data) {
		return EINVAL
	}

	
	p := &data[cap(data)-1]
	m.Lock()
	defer m.Unlock()
	b := m.active[p]
	if b == nil || &b[0] != &data[0] {
		return EINVAL
	}

	
	if errno := m.munmap(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b))); errno != nil {
		return errno
	}
	delete(m.active, p)
	return nil
}

func Read(fd int, p []byte) (n int, err error) {
	n, err = read(fd, p)
	if raceenabled {
		if n > 0 {
			raceWriteRange(unsafe.Pointer(&p[0]), n)
		}
		if err == nil {
			raceAcquire(unsafe.Pointer(&ioSync))
		}
	}
	return
}

func Write(fd int, p []byte) (n int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	n, err = write(fd, p)
	if raceenabled && n > 0 {
		raceReadRange(unsafe.Pointer(&p[0]), n)
	}
	return
}



var SocketDisableIPv6 bool


type Sockaddr interface {
	sockaddr() (ptr unsafe.Pointer, len _Socklen, err error) 
}


type SockaddrInet4 struct {
	Port int
	Addr [4]byte
	raw  RawSockaddrInet4
}


type SockaddrInet6 struct {
	Port   int
	ZoneId uint32
	Addr   [16]byte
	raw    RawSockaddrInet6
}


type SockaddrUnix struct {
	Name string
	raw  RawSockaddrUnix
}

func Bind(fd int, sa Sockaddr) (err error) {
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return bind(fd, ptr, n)
}

func Connect(fd int, sa Sockaddr) (err error) {
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return connect(fd, ptr, n)
}

func Getpeername(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getpeername(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(fd, &rsa)
}

func GetsockoptByte(fd, level, opt int) (value byte, err error) {
	var n byte
	vallen := _Socklen(1)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return n, err
}

func GetsockoptInt(fd, level, opt int) (value int, err error) {
	var n int32
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return int(n), err
}

func GetsockoptInet4Addr(fd, level, opt int) (value [4]byte, err error) {
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	return value, err
}

func GetsockoptIPMreq(fd, level, opt int) (*IPMreq, error) {
	var value IPMreq
	vallen := _Socklen(SizeofIPMreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptIPv6Mreq(fd, level, opt int) (*IPv6Mreq, error) {
	var value IPv6Mreq
	vallen := _Socklen(SizeofIPv6Mreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptIPv6MTUInfo(fd, level, opt int) (*IPv6MTUInfo, error) {
	var value IPv6MTUInfo
	vallen := _Socklen(SizeofIPv6MTUInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptICMPv6Filter(fd, level, opt int) (*ICMPv6Filter, error) {
	var value ICMPv6Filter
	vallen := _Socklen(SizeofICMPv6Filter)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptLinger(fd, level, opt int) (*Linger, error) {
	var linger Linger
	vallen := _Socklen(SizeofLinger)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&linger), &vallen)
	return &linger, err
}

func GetsockoptTimeval(fd, level, opt int) (*Timeval, error) {
	var tv Timeval
	vallen := _Socklen(unsafe.Sizeof(tv))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&tv), &vallen)
	return &tv, err
}

func GetsockoptUint64(fd, level, opt int) (value uint64, err error) {
	var n uint64
	vallen := _Socklen(8)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return n, err
}

func Recvfrom(fd int, p []byte, flags int) (n int, from Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if n, err = recvfrom(fd, p, flags, &rsa, &len); err != nil {
		return
	}
	if rsa.Addr.Family != AF_UNSPEC {
		from, err = anyToSockaddr(fd, &rsa)
	}
	return
}

func Sendto(fd int, p []byte, flags int, to Sockaddr) (err error) {
	ptr, n, err := to.sockaddr()
	if err != nil {
		return err
	}
	return sendto(fd, p, flags, ptr, n)
}

func SetsockoptByte(fd, level, opt int, value byte) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(&value), 1)
}

func SetsockoptInt(fd, level, opt int, value int) (err error) {
	var n = int32(value)
	return setsockopt(fd, level, opt, unsafe.Pointer(&n), 4)
}

func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(&value[0]), 4)
}

func SetsockoptIPMreq(fd, level, opt int, mreq *IPMreq) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), SizeofIPMreq)
}

func SetsockoptIPv6Mreq(fd, level, opt int, mreq *IPv6Mreq) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), SizeofIPv6Mreq)
}

func SetsockoptICMPv6Filter(fd, level, opt int, filter *ICMPv6Filter) error {
	return setsockopt(fd, level, opt, unsafe.Pointer(filter), SizeofICMPv6Filter)
}

func SetsockoptLinger(fd, level, opt int, l *Linger) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(l), SizeofLinger)
}

func SetsockoptString(fd, level, opt int, s string) (err error) {
	var p unsafe.Pointer
	if len(s) > 0 {
		p = unsafe.Pointer(&[]byte(s)[0])
	}
	return setsockopt(fd, level, opt, p, uintptr(len(s)))
}

func SetsockoptTimeval(fd, level, opt int, tv *Timeval) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(tv), unsafe.Sizeof(*tv))
}

func SetsockoptUint64(fd, level, opt int, value uint64) (err error) {
	return setsockopt(fd, level, opt, unsafe.Pointer(&value), 8)
}

func Socket(domain, typ, proto int) (fd int, err error) {
	if domain == AF_INET6 && SocketDisableIPv6 {
		return -1, EAFNOSUPPORT
	}
	fd, err = socket(domain, typ, proto)
	return
}

func Socketpair(domain, typ, proto int) (fd [2]int, err error) {
	var fdx [2]int32
	err = socketpair(domain, typ, proto, &fdx)
	if err == nil {
		fd[0] = int(fdx[0])
		fd[1] = int(fdx[1])
	}
	return
}

var ioSync int64

func CloseOnExec(fd int) { fcntl(fd, F_SETFD, FD_CLOEXEC) }

func SetNonblock(fd int, nonblocking bool) (err error) {
	flag, err := fcntl(fd, F_GETFL, 0)
	if err != nil {
		return err
	}
	if nonblocking {
		flag |= O_NONBLOCK
	} else {
		flag &= ^O_NONBLOCK
	}
	_, err = fcntl(fd, F_SETFL, flag)
	return err
}






func Exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

func Getag(path string) (ccsid uint16, flag uint16, err error) {
	var val [8]byte
	sz, err := Getxattr(path, "ccsid", val[:])
	if err != nil {
		return
	}
	ccsid = uint16(EncodeData(val[0:sz]))
	sz, err = Getxattr(path, "flags", val[:])
	if err != nil {
		return
	}
	flag = uint16(EncodeData(val[0:sz]) >> 15)
	return
}


func impl_Mount(source string, target string, fstype string, flags uintptr, data string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(source)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(target)
	if err != nil {
		return
	}
	var _p2 *byte
	_p2, err = BytePtrFromString(fstype)
	if err != nil {
		return
	}
	var _p3 *byte
	_p3, err = BytePtrFromString(data)
	if err != nil {
		return
	}
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MOUNT1_A<<4, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(unsafe.Pointer(_p2)), uintptr(flags), uintptr(unsafe.Pointer(_p3)))
	runtime.ExitSyscall()
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_MountAddr() *(func(source string, target string, fstype string, flags uintptr, data string) (err error))

var Mount = enter_Mount

func enter_Mount(source string, target string, fstype string, flags uintptr, data string) (err error) {
	funcref := get_MountAddr()
	if validMount() {
		*funcref = impl_Mount
	} else {
		*funcref = legacyMount
	}
	return (*funcref)(source, target, fstype, flags, data)
}

func legacyMount(source string, target string, fstype string, flags uintptr, data string) (err error) {
	if needspace := 8 - len(fstype); needspace <= 0 {
		fstype = fstype[0:8]
	} else {
		fstype += "        "[0:needspace]
	}
	return mount_LE(target, source, fstype, uint32(flags), int32(len(data)), data)
}

func validMount() bool {
	if funcptrtest(GetZosLibVec()+SYS___MOUNT1_A<<4, "") == 0 {
		if name, err := getLeFuncName(GetZosLibVec() + SYS___MOUNT1_A<<4); err == nil {
			return name == "__mount1_a"
		}
	}
	return false
}




func impl_Unmount(target string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(target)
	if err != nil {
		return
	}
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___UMOUNT2_A<<4, uintptr(unsafe.Pointer(_p0)), uintptr(flags))
	runtime.ExitSyscall()
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_UnmountAddr() *(func(target string, flags int) (err error))

var Unmount = enter_Unmount

func enter_Unmount(target string, flags int) (err error) {
	funcref := get_UnmountAddr()
	if funcptrtest(GetZosLibVec()+SYS___UMOUNT2_A<<4, "") == 0 {
		*funcref = impl_Unmount
	} else {
		*funcref = legacyUnmount
	}
	return (*funcref)(target, flags)
}

func legacyUnmount(name string, mtm int) (err error) {
	
	
	if name[0] != '/' {
		return unmount_LE(name, mtm)
	}
	
	b2s := func(arr []byte) string {
		var str string
		for i := 0; i < len(arr); i++ {
			if arr[i] == 0 {
				str = string(arr[:i])
				break
			}
		}
		return str
	}
	var buffer struct {
		header W_Mnth
		fsinfo [64]W_Mntent
	}
	fs_count, err := W_Getmntent_A((*byte)(unsafe.Pointer(&buffer)), int(unsafe.Sizeof(buffer)))
	if err == nil {
		err = EINVAL
		for i := 0; i < fs_count; i++ {
			if b2s(buffer.fsinfo[i].Mountpoint[:]) == name {
				err = unmount_LE(b2s(buffer.fsinfo[i].Fsname[:]), mtm)
				break
			}
		}
	} else if fs_count == 0 {
		err = EINVAL
	}
	return err
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

func direntLeToDirentUnix(dirent *direntLE, dir uintptr, path string) (Dirent, error) {
	var d Dirent

	d.Ino = uint64(dirent.Ino)
	offset, err := Telldir(dir)
	if err != nil {
		return d, err
	}

	d.Off = int64(offset)
	s := string(bytes.Split(dirent.Name[:], []byte{0})[0])
	copy(d.Name[:], s)

	d.Reclen = uint16(24 + len(d.NameString()))
	var st Stat_t
	path = path + "/" + s
	err = Lstat(path, &st)
	if err != nil {
		return d, err
	}

	d.Type = uint8(st.Mode >> 24)
	return d, err
}

func Getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error) {
	
	
	
	
	
	

	skip, err := Seek(fd, 0, 1 )
	if err != nil {
		return 0, err
	}

	
	path, err := ZosFdToPath(fd)
	if err != nil {
		return 0, err
	}
	d, err := Opendir(path)
	if err != nil {
		return 0, err
	}
	defer Closedir(d)

	var cnt int64
	for {
		var entryLE direntLE
		var entrypLE *direntLE
		e := Readdir_r(d, &entryLE, &entrypLE)
		if e != nil {
			return n, e
		}
		if entrypLE == nil {
			break
		}
		if skip > 0 {
			skip--
			cnt++
			continue
		}

		
		entry, e := direntLeToDirentUnix(&entryLE, d, path)
		if e != nil {
			return n, e
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

func Err2ad() (eadd *int) {
	r0, _, _ := CallLeFuncWithErr(GetZosLibVec() + SYS___ERR2AD<<4)
	eadd = (*int)(unsafe.Pointer(r0))
	return
}

func ZosConsolePrintf(format string, v ...interface{}) (int, error) {
	type __cmsg struct {
		_            uint16
		_            [2]uint8
		__msg_length uint32
		__msg        uintptr
		_            [4]uint8
	}
	msg := fmt.Sprintf(format, v...)
	strptr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&msg)).Data)
	len := (*reflect.StringHeader)(unsafe.Pointer(&msg)).Len
	cmsg := __cmsg{__msg_length: uint32(len), __msg: uintptr(strptr)}
	cmd := uint32(0)
	runtime.EnterSyscall()
	rc, err2, err1 := CallLeFuncWithErr(GetZosLibVec()+SYS_____CONSOLE_A<<4, uintptr(unsafe.Pointer(&cmsg)), 0, uintptr(unsafe.Pointer(&cmd)))
	runtime.ExitSyscall()
	if rc != 0 {
		return 0, fmt.Errorf("%s (errno2=0x%x)\n", err1.Error(), err2)
	}
	return 0, nil
}
func ZosStringToEbcdicBytes(str string, nullterm bool) (ebcdicBytes []byte) {
	if nullterm {
		ebcdicBytes = []byte(str + "\x00")
	} else {
		ebcdicBytes = []byte(str)
	}
	A2e(ebcdicBytes)
	return
}
func ZosEbcdicBytesToString(b []byte, trimRight bool) (str string) {
	res := make([]byte, len(b))
	copy(res, b)
	E2a(res)
	if trimRight {
		str = string(bytes.TrimRight(res, " \x00"))
	} else {
		str = string(res)
	}
	return
}

func fdToPath(dirfd int) (path string, err error) {
	var buffer [1024]byte
	
	ret := runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS_W_IOCTL<<4,
		[]uintptr{uintptr(dirfd), 17, 1024, uintptr(unsafe.Pointer(&buffer[0]))})
	if ret == 0 {
		zb := bytes.IndexByte(buffer[:], 0)
		if zb == -1 {
			zb = len(buffer)
		}
		
		runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___E2A_L<<4,
			[]uintptr{uintptr(unsafe.Pointer(&buffer[0])), uintptr(zb)})
		return string(buffer[:zb]), nil
	}
	
	errno := int(*(*int32)(unsafe.Pointer(runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___ERRNO<<4,
		[]uintptr{}))))
	
	errno2 := int(runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS___ERRNO2<<4,
		[]uintptr{}))
	
	ret = runtime.CallLeFuncByPtr(runtime.XplinkLibvec+SYS_STRERROR_R<<4,
		[]uintptr{uintptr(errno), uintptr(unsafe.Pointer(&buffer[0])), 1024})
	if ret == 0 {
		zb := bytes.IndexByte(buffer[:], 0)
		if zb == -1 {
			zb = len(buffer)
		}
		return "", fmt.Errorf("%s (errno2=0x%x)", buffer[:zb], errno2)
	} else {
		return "", fmt.Errorf("fdToPath errno %d (errno2=0x%x)", errno, errno2)
	}
}

func impl_Mkfifoat(dirfd int, path string, mode uint32) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS___MKFIFOAT_A<<4, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode))
	runtime.ExitSyscall()
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

//go:nosplit
func get_MkfifoatAddr() *(func(dirfd int, path string, mode uint32) (err error))

var Mkfifoat = enter_Mkfifoat

func enter_Mkfifoat(dirfd int, path string, mode uint32) (err error) {
	funcref := get_MkfifoatAddr()
	if funcptrtest(GetZosLibVec()+SYS___MKFIFOAT_A<<4, "") == 0 {
		*funcref = impl_Mkfifoat
	} else {
		*funcref = legacy_Mkfifoat
	}
	return (*funcref)(dirfd, path, mode)
}

func legacy_Mkfifoat(dirfd int, path string, mode uint32) (err error) {
	dirname, err := ZosFdToPath(dirfd)
	if err != nil {
		return err
	}
	return Mkfifo(dirname+"/"+path, mode)
}





func fcntlAsIs(fd uintptr, cmd int, arg uintptr) (val int, err error) {
	runtime.EnterSyscall()
	r0, e2, e1 := CallLeFuncWithErr(GetZosLibVec()+SYS_FCNTL<<4, uintptr(fd), uintptr(cmd), arg)
	runtime.ExitSyscall()
	val = int(r0)
	if int64(r0) == -1 {
		err = errnoErr2(e1, e2)
	}
	return
}

func Fcntl(fd uintptr, cmd int, op interface{}) (ret int, err error) {
	switch op.(type) {
	case *Flock_t:
		err = FcntlFlock(fd, cmd, op.(*Flock_t))
		if err != nil {
			ret = -1
		}
		return
	case int:
		return FcntlInt(fd, cmd, op.(int))
	case *F_cnvrt:
		return fcntlAsIs(fd, cmd, uintptr(unsafe.Pointer(op.(*F_cnvrt))))
	case unsafe.Pointer:
		return fcntlAsIs(fd, cmd, uintptr(op.(unsafe.Pointer)))
	default:
		return -1, EINVAL
	}
	return
}

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	return sendfile(outfd, infd, offset, count)
}

func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) {
	
	originalOffset, err := Seek(infd, 0, SEEK_CUR)
	if err != nil {
		return -1, err
	}
	
	if offset != nil {
		_, err := Seek(infd, *offset, SEEK_SET)
		if err != nil {
			return -1, err
		}
	}

	buf := make([]byte, count)
	readBuf := make([]byte, 0)
	var n int = 0
	for i := 0; i < count; i += n {
		n, err := Read(infd, buf)
		if n == 0 {
			if err != nil {
				return -1, err
			} else { 
				break
			}
		}
		readBuf = append(readBuf, buf...)
		buf = buf[0:0]
	}

	n2, err := Write(outfd, readBuf)
	if err != nil {
		return -1, err
	}

	
	
	if offset != nil {
		*offset = *offset + int64(n)
		
		
		_, err := Seek(infd, originalOffset, SEEK_SET)
		if err != nil {
			return -1, err
		}
	}
	return n2, nil
}
