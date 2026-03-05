

//go:build linux

package unix

import (
	"syscall"
	"unsafe"
)



func FanotifyInit(flags uint, event_f_flags uint) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_FANOTIFY_INIT, uintptr(flags), uintptr(event_f_flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func fchmodat(dirfd int, path string, mode uint32) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_FCHMODAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func fchmodat2(dirfd int, path string, mode uint32, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_FCHMODAT2, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ioctl(fd int, req uint, arg uintptr) (err error) {
	_, _, e1 := Syscall(SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(arg))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ioctlPtr(fd int, req uint, arg unsafe.Pointer) (err error) {
	_, _, e1 := Syscall(SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(arg))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(oldpath)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(newpath)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_LINKAT, uintptr(olddirfd), uintptr(unsafe.Pointer(_p0)), uintptr(newdirfd), uintptr(unsafe.Pointer(_p1)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags), uintptr(mode), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func openat2(dirfd int, path string, open_how *OpenHow, size int) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_OPENAT2, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(open_how)), uintptr(size), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func pipe2(p *[2]_C_int, flags int) (err error) {
	_, _, e1 := RawSyscall(SYS_PIPE2, uintptr(unsafe.Pointer(p)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ppoll(fds *PollFd, nfds int, timeout *Timespec, sigmask *Sigset_t) (n int, err error) {
	r0, _, e1 := Syscall6(SYS_PPOLL, uintptr(unsafe.Pointer(fds)), uintptr(nfds), uintptr(unsafe.Pointer(timeout)), uintptr(unsafe.Pointer(sigmask)), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Readlinkat(dirfd int, path string, buf []byte) (n int, err error) {
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
	r0, _, e1 := Syscall6(SYS_READLINKAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(buf)), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Symlinkat(oldpath string, newdirfd int, newpath string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(oldpath)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(newpath)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_SYMLINKAT, uintptr(unsafe.Pointer(_p0)), uintptr(newdirfd), uintptr(unsafe.Pointer(_p1)))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Unlinkat(dirfd int, path string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_UNLINKAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func utimensat(dirfd int, path string, times *[2]Timespec, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_UTIMENSAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(times)), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Getcwd(buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_GETCWD, uintptr(_p0), uintptr(len(buf)), 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, err error) {
	r0, _, e1 := Syscall6(SYS_WAIT4, uintptr(pid), uintptr(unsafe.Pointer(wstatus)), uintptr(options), uintptr(unsafe.Pointer(rusage)), 0, 0)
	wpid = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Waitid(idType int, id int, info *Siginfo, options int, rusage *Rusage) (err error) {
	_, _, e1 := Syscall6(SYS_WAITID, uintptr(idType), uintptr(id), uintptr(unsafe.Pointer(info)), uintptr(options), uintptr(unsafe.Pointer(rusage)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func KeyctlInt(cmd int, arg2 int, arg3 int, arg4 int, arg5 int) (ret int, err error) {
	r0, _, e1 := Syscall6(SYS_KEYCTL, uintptr(cmd), uintptr(arg2), uintptr(arg3), uintptr(arg4), uintptr(arg5), 0)
	ret = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func KeyctlBuffer(cmd int, arg2 int, buf []byte, arg5 int) (ret int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_KEYCTL, uintptr(cmd), uintptr(arg2), uintptr(_p0), uintptr(len(buf)), uintptr(arg5), 0)
	ret = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func keyctlJoin(cmd int, arg2 string) (ret int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(arg2)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_KEYCTL, uintptr(cmd), uintptr(unsafe.Pointer(_p0)), 0)
	ret = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func keyctlSearch(cmd int, arg2 int, arg3 string, arg4 string, arg5 int) (ret int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(arg3)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(arg4)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_KEYCTL, uintptr(cmd), uintptr(arg2), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(arg5), 0)
	ret = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func keyctlIOV(cmd int, arg2 int, payload []Iovec, arg5 int) (err error) {
	var _p0 unsafe.Pointer
	if len(payload) > 0 {
		_p0 = unsafe.Pointer(&payload[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall6(SYS_KEYCTL, uintptr(cmd), uintptr(arg2), uintptr(_p0), uintptr(len(payload)), uintptr(arg5), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func keyctlDH(cmd int, arg2 *KeyctlDHParams, buf []byte) (ret int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_KEYCTL, uintptr(cmd), uintptr(unsafe.Pointer(arg2)), uintptr(_p0), uintptr(len(buf)), 0, 0)
	ret = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func keyctlRestrictKeyringByType(cmd int, arg2 int, keyType string, restriction string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(keyType)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(restriction)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_KEYCTL, uintptr(cmd), uintptr(arg2), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func keyctlRestrictKeyring(cmd int, arg2 int) (err error) {
	_, _, e1 := Syscall(SYS_KEYCTL, uintptr(cmd), uintptr(arg2), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ptrace(request int, pid int, addr uintptr, data uintptr) (err error) {
	_, _, e1 := Syscall6(SYS_PTRACE, uintptr(request), uintptr(pid), uintptr(addr), uintptr(data), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ptracePtr(request int, pid int, addr uintptr, data unsafe.Pointer) (err error) {
	_, _, e1 := Syscall6(SYS_PTRACE, uintptr(request), uintptr(pid), uintptr(addr), uintptr(data), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func reboot(magic1 uint, magic2 uint, cmd int, arg string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(arg)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_REBOOT, uintptr(magic1), uintptr(magic2), uintptr(cmd), uintptr(unsafe.Pointer(_p0)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func mount(source string, target string, fstype string, flags uintptr, data *byte) (err error) {
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
	_, _, e1 := Syscall6(SYS_MOUNT, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(unsafe.Pointer(_p2)), uintptr(flags), uintptr(unsafe.Pointer(data)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func mountSetattr(dirfd int, pathname string, flags uint, attr *MountAttr, size uintptr) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(pathname)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_MOUNT_SETATTR, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags), uintptr(unsafe.Pointer(attr)), uintptr(size), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Acct(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_ACCT, uintptr(unsafe.Pointer(_p0)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func AddKey(keyType string, description string, payload []byte, ringid int) (id int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(keyType)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(description)
	if err != nil {
		return
	}
	var _p2 unsafe.Pointer
	if len(payload) > 0 {
		_p2 = unsafe.Pointer(&payload[0])
	} else {
		_p2 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_ADD_KEY, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(payload)), uintptr(ringid), 0)
	id = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Adjtimex(buf *Timex) (state int, err error) {
	r0, _, e1 := Syscall(SYS_ADJTIMEX, uintptr(unsafe.Pointer(buf)), 0, 0)
	state = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Capget(hdr *CapUserHeader, data *CapUserData) (err error) {
	_, _, e1 := RawSyscall(SYS_CAPGET, uintptr(unsafe.Pointer(hdr)), uintptr(unsafe.Pointer(data)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Capset(hdr *CapUserHeader, data *CapUserData) (err error) {
	_, _, e1 := RawSyscall(SYS_CAPSET, uintptr(unsafe.Pointer(hdr)), uintptr(unsafe.Pointer(data)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Chdir(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_CHDIR, uintptr(unsafe.Pointer(_p0)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Chroot(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_CHROOT, uintptr(unsafe.Pointer(_p0)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ClockAdjtime(clockid int32, buf *Timex) (state int, err error) {
	r0, _, e1 := Syscall(SYS_CLOCK_ADJTIME, uintptr(clockid), uintptr(unsafe.Pointer(buf)), 0)
	state = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ClockGetres(clockid int32, res *Timespec) (err error) {
	_, _, e1 := Syscall(SYS_CLOCK_GETRES, uintptr(clockid), uintptr(unsafe.Pointer(res)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ClockGettime(clockid int32, time *Timespec) (err error) {
	_, _, e1 := Syscall(SYS_CLOCK_GETTIME, uintptr(clockid), uintptr(unsafe.Pointer(time)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ClockSettime(clockid int32, time *Timespec) (err error) {
	_, _, e1 := Syscall(SYS_CLOCK_SETTIME, uintptr(clockid), uintptr(unsafe.Pointer(time)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ClockNanosleep(clockid int32, flags int, request *Timespec, remain *Timespec) (err error) {
	_, _, e1 := Syscall6(SYS_CLOCK_NANOSLEEP, uintptr(clockid), uintptr(flags), uintptr(unsafe.Pointer(request)), uintptr(unsafe.Pointer(remain)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Close(fd int) (err error) {
	_, _, e1 := Syscall(SYS_CLOSE, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func CloseRange(first uint, last uint, flags uint) (err error) {
	_, _, e1 := Syscall(SYS_CLOSE_RANGE, uintptr(first), uintptr(last), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func CopyFileRange(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error) {
	r0, _, e1 := Syscall6(SYS_COPY_FILE_RANGE, uintptr(rfd), uintptr(unsafe.Pointer(roff)), uintptr(wfd), uintptr(unsafe.Pointer(woff)), uintptr(len), uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func DeleteModule(name string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(name)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_DELETE_MODULE, uintptr(unsafe.Pointer(_p0)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Dup(oldfd int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_DUP, uintptr(oldfd), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Dup3(oldfd int, newfd int, flags int) (err error) {
	_, _, e1 := Syscall(SYS_DUP3, uintptr(oldfd), uintptr(newfd), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func EpollCreate1(flag int) (fd int, err error) {
	r0, _, e1 := RawSyscall(SYS_EPOLL_CREATE1, uintptr(flag), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func EpollCtl(epfd int, op int, fd int, event *EpollEvent) (err error) {
	_, _, e1 := RawSyscall6(SYS_EPOLL_CTL, uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Eventfd(initval uint, flags int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_EVENTFD2, uintptr(initval), uintptr(flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Exit(code int) {
	SyscallNoError(SYS_EXIT_GROUP, uintptr(code), 0, 0)
	return
}



func Fchdir(fd int) (err error) {
	_, _, e1 := Syscall(SYS_FCHDIR, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fchmod(fd int, mode uint32) (err error) {
	_, _, e1 := Syscall(SYS_FCHMOD, uintptr(fd), uintptr(mode), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_FCHOWNAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(uid), uintptr(gid), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fdatasync(fd int) (err error) {
	_, _, e1 := Syscall(SYS_FDATASYNC, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fgetxattr(fd int, attr string, dest []byte) (sz int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(attr)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(dest) > 0 {
		_p1 = unsafe.Pointer(&dest[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_FGETXATTR, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(dest)), 0, 0)
	sz = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func FinitModule(fd int, params string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(params)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_FINIT_MODULE, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Flistxattr(fd int, dest []byte) (sz int, err error) {
	var _p0 unsafe.Pointer
	if len(dest) > 0 {
		_p0 = unsafe.Pointer(&dest[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FLISTXATTR, uintptr(fd), uintptr(_p0), uintptr(len(dest)))
	sz = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Flock(fd int, how int) (err error) {
	_, _, e1 := Syscall(SYS_FLOCK, uintptr(fd), uintptr(how), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fremovexattr(fd int, attr string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(attr)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_FREMOVEXATTR, uintptr(fd), uintptr(unsafe.Pointer(_p0)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fsetxattr(fd int, attr string, dest []byte, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(attr)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(dest) > 0 {
		_p1 = unsafe.Pointer(&dest[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall6(SYS_FSETXATTR, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(dest)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fsync(fd int) (err error) {
	_, _, e1 := Syscall(SYS_FSYNC, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fsmount(fd int, flags int, mountAttrs int) (fsfd int, err error) {
	r0, _, e1 := Syscall(SYS_FSMOUNT, uintptr(fd), uintptr(flags), uintptr(mountAttrs))
	fsfd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fsopen(fsName string, flags int) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(fsName)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_FSOPEN, uintptr(unsafe.Pointer(_p0)), uintptr(flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Fspick(dirfd int, pathName string, flags int) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(pathName)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_FSPICK, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags))
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func fsconfig(fd int, cmd uint, key *byte, value *byte, aux int) (err error) {
	_, _, e1 := Syscall6(SYS_FSCONFIG, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(value)), uintptr(aux), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Getdents(fd int, buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_GETDENTS64, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Getpgid(pid int) (pgid int, err error) {
	r0, _, e1 := RawSyscall(SYS_GETPGID, uintptr(pid), 0, 0)
	pgid = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Getpid() (pid int) {
	r0, _ := RawSyscallNoError(SYS_GETPID, 0, 0, 0)
	pid = int(r0)
	return
}



func Getppid() (ppid int) {
	r0, _ := RawSyscallNoError(SYS_GETPPID, 0, 0, 0)
	ppid = int(r0)
	return
}



func Getpriority(which int, who int) (prio int, err error) {
	r0, _, e1 := Syscall(SYS_GETPRIORITY, uintptr(which), uintptr(who), 0)
	prio = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Getrusage(who int, rusage *Rusage) (err error) {
	_, _, e1 := RawSyscall(SYS_GETRUSAGE, uintptr(who), uintptr(unsafe.Pointer(rusage)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Getsid(pid int) (sid int, err error) {
	r0, _, e1 := RawSyscall(SYS_GETSID, uintptr(pid), 0, 0)
	sid = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Gettid() (tid int) {
	r0, _ := RawSyscallNoError(SYS_GETTID, 0, 0, 0)
	tid = int(r0)
	return
}



func Getxattr(path string, attr string, dest []byte) (sz int, err error) {
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
	r0, _, e1 := Syscall6(SYS_GETXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(dest)), 0, 0)
	sz = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func InitModule(moduleImage []byte, params string) (err error) {
	var _p0 unsafe.Pointer
	if len(moduleImage) > 0 {
		_p0 = unsafe.Pointer(&moduleImage[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(params)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_INIT_MODULE, uintptr(_p0), uintptr(len(moduleImage)), uintptr(unsafe.Pointer(_p1)))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func InotifyAddWatch(fd int, pathname string, mask uint32) (watchdesc int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(pathname)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_INOTIFY_ADD_WATCH, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(mask))
	watchdesc = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func InotifyInit1(flags int) (fd int, err error) {
	r0, _, e1 := RawSyscall(SYS_INOTIFY_INIT1, uintptr(flags), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func InotifyRmWatch(fd int, watchdesc uint32) (success int, err error) {
	r0, _, e1 := RawSyscall(SYS_INOTIFY_RM_WATCH, uintptr(fd), uintptr(watchdesc), 0)
	success = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Kill(pid int, sig syscall.Signal) (err error) {
	_, _, e1 := RawSyscall(SYS_KILL, uintptr(pid), uintptr(sig), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Klogctl(typ int, buf []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_SYSLOG, uintptr(typ), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Lgetxattr(path string, attr string, dest []byte) (sz int, err error) {
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
	r0, _, e1 := Syscall6(SYS_LGETXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(dest)), 0, 0)
	sz = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Listxattr(path string, dest []byte) (sz int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(dest) > 0 {
		_p1 = unsafe.Pointer(&dest[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_LISTXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(dest)))
	sz = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Llistxattr(path string, dest []byte) (sz int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(dest) > 0 {
		_p1 = unsafe.Pointer(&dest[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_LLISTXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(dest)))
	sz = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Lremovexattr(path string, attr string) (err error) {
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
	_, _, e1 := Syscall(SYS_LREMOVEXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Lsetxattr(path string, attr string, data []byte, flags int) (err error) {
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
	_, _, e1 := Syscall6(SYS_LSETXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(data)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func MemfdCreate(name string, flags int) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(name)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_MEMFD_CREATE, uintptr(unsafe.Pointer(_p0)), uintptr(flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Mkdirat(dirfd int, path string, mode uint32) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_MKDIRAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Mknodat(dirfd int, path string, mode uint32, dev int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_MKNODAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode), uintptr(dev), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func MoveMount(fromDirfd int, fromPathName string, toDirfd int, toPathName string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(fromPathName)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(toPathName)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_MOVE_MOUNT, uintptr(fromDirfd), uintptr(unsafe.Pointer(_p0)), uintptr(toDirfd), uintptr(unsafe.Pointer(_p1)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Nanosleep(time *Timespec, leftover *Timespec) (err error) {
	_, _, e1 := Syscall(SYS_NANOSLEEP, uintptr(unsafe.Pointer(time)), uintptr(unsafe.Pointer(leftover)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func OpenTree(dfd int, fileName string, flags uint) (r int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(fileName)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_OPEN_TREE, uintptr(dfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags))
	r = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func PerfEventOpen(attr *PerfEventAttr, pid int, cpu int, groupFd int, flags int) (fd int, err error) {
	r0, _, e1 := Syscall6(SYS_PERF_EVENT_OPEN, uintptr(unsafe.Pointer(attr)), uintptr(pid), uintptr(cpu), uintptr(groupFd), uintptr(flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func PivotRoot(newroot string, putold string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(newroot)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(putold)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_PIVOT_ROOT, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Prctl(option int, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr) (err error) {
	_, _, e1 := Syscall6(SYS_PRCTL, uintptr(option), uintptr(arg2), uintptr(arg3), uintptr(arg4), uintptr(arg5), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func pselect6(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timespec, sigmask *sigset_argpack) (n int, err error) {
	r0, _, e1 := Syscall6(SYS_PSELECT6, uintptr(nfd), uintptr(unsafe.Pointer(r)), uintptr(unsafe.Pointer(w)), uintptr(unsafe.Pointer(e)), uintptr(unsafe.Pointer(timeout)), uintptr(unsafe.Pointer(sigmask)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func read(fd int, p []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_READ, uintptr(fd), uintptr(_p0), uintptr(len(p)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Removexattr(path string, attr string) (err error) {
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
	_, _, e1 := Syscall(SYS_REMOVEXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Renameat2(olddirfd int, oldpath string, newdirfd int, newpath string, flags uint) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(oldpath)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(newpath)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_RENAMEAT2, uintptr(olddirfd), uintptr(unsafe.Pointer(_p0)), uintptr(newdirfd), uintptr(unsafe.Pointer(_p1)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func RequestKey(keyType string, description string, callback string, destRingid int) (id int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(keyType)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(description)
	if err != nil {
		return
	}
	var _p2 *byte
	_p2, err = BytePtrFromString(callback)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_REQUEST_KEY, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(unsafe.Pointer(_p2)), uintptr(destRingid), 0, 0)
	id = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Setdomainname(p []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_SETDOMAINNAME, uintptr(_p0), uintptr(len(p)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Sethostname(p []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_SETHOSTNAME, uintptr(_p0), uintptr(len(p)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Setpgid(pid int, pgid int) (err error) {
	_, _, e1 := RawSyscall(SYS_SETPGID, uintptr(pid), uintptr(pgid), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Setsid() (pid int, err error) {
	r0, _, e1 := RawSyscall(SYS_SETSID, 0, 0, 0)
	pid = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Settimeofday(tv *Timeval) (err error) {
	_, _, e1 := RawSyscall(SYS_SETTIMEOFDAY, uintptr(unsafe.Pointer(tv)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Setns(fd int, nstype int) (err error) {
	_, _, e1 := Syscall(SYS_SETNS, uintptr(fd), uintptr(nstype), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Setpriority(which int, who int, prio int) (err error) {
	_, _, e1 := Syscall(SYS_SETPRIORITY, uintptr(which), uintptr(who), uintptr(prio))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Setxattr(path string, attr string, data []byte, flags int) (err error) {
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
	_, _, e1 := Syscall6(SYS_SETXATTR, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(_p2), uintptr(len(data)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func signalfd(fd int, sigmask *Sigset_t, maskSize uintptr, flags int) (newfd int, err error) {
	r0, _, e1 := Syscall6(SYS_SIGNALFD4, uintptr(fd), uintptr(unsafe.Pointer(sigmask)), uintptr(maskSize), uintptr(flags), 0, 0)
	newfd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Statx(dirfd int, path string, flags int, mask int, stat *Statx_t) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_STATX, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags), uintptr(mask), uintptr(unsafe.Pointer(stat)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Sync() {
	SyscallNoError(SYS_SYNC, 0, 0, 0)
	return
}



func Syncfs(fd int) (err error) {
	_, _, e1 := Syscall(SYS_SYNCFS, uintptr(fd), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Sysinfo(info *Sysinfo_t) (err error) {
	_, _, e1 := RawSyscall(SYS_SYSINFO, uintptr(unsafe.Pointer(info)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func TimerfdCreate(clockid int, flags int) (fd int, err error) {
	r0, _, e1 := RawSyscall(SYS_TIMERFD_CREATE, uintptr(clockid), uintptr(flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func TimerfdGettime(fd int, currValue *ItimerSpec) (err error) {
	_, _, e1 := RawSyscall(SYS_TIMERFD_GETTIME, uintptr(fd), uintptr(unsafe.Pointer(currValue)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func TimerfdSettime(fd int, flags int, newValue *ItimerSpec, oldValue *ItimerSpec) (err error) {
	_, _, e1 := RawSyscall6(SYS_TIMERFD_SETTIME, uintptr(fd), uintptr(flags), uintptr(unsafe.Pointer(newValue)), uintptr(unsafe.Pointer(oldValue)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Tgkill(tgid int, tid int, sig syscall.Signal) (err error) {
	_, _, e1 := RawSyscall(SYS_TGKILL, uintptr(tgid), uintptr(tid), uintptr(sig))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Times(tms *Tms) (ticks uintptr, err error) {
	r0, _, e1 := RawSyscall(SYS_TIMES, uintptr(unsafe.Pointer(tms)), 0, 0)
	ticks = uintptr(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Umask(mask int) (oldmask int) {
	r0, _ := RawSyscallNoError(SYS_UMASK, uintptr(mask), 0, 0)
	oldmask = int(r0)
	return
}



func Uname(buf *Utsname) (err error) {
	_, _, e1 := RawSyscall(SYS_UNAME, uintptr(unsafe.Pointer(buf)), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Unmount(target string, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(target)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_UMOUNT2, uintptr(unsafe.Pointer(_p0)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Unshare(flags int) (err error) {
	_, _, e1 := Syscall(SYS_UNSHARE, uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func write(fd int, p []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_WRITE, uintptr(fd), uintptr(_p0), uintptr(len(p)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func exitThread(code int) (err error) {
	_, _, e1 := Syscall(SYS_EXIT, uintptr(code), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func readv(fd int, iovs []Iovec) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(iovs) > 0 {
		_p0 = unsafe.Pointer(&iovs[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_READV, uintptr(fd), uintptr(_p0), uintptr(len(iovs)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func writev(fd int, iovs []Iovec) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(iovs) > 0 {
		_p0 = unsafe.Pointer(&iovs[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_WRITEV, uintptr(fd), uintptr(_p0), uintptr(len(iovs)))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func preadv(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(iovs) > 0 {
		_p0 = unsafe.Pointer(&iovs[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PREADV, uintptr(fd), uintptr(_p0), uintptr(len(iovs)), uintptr(offs_l), uintptr(offs_h), 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func pwritev(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(iovs) > 0 {
		_p0 = unsafe.Pointer(&iovs[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PWRITEV, uintptr(fd), uintptr(_p0), uintptr(len(iovs)), uintptr(offs_l), uintptr(offs_h), 0)
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func preadv2(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr, flags int) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(iovs) > 0 {
		_p0 = unsafe.Pointer(&iovs[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PREADV2, uintptr(fd), uintptr(_p0), uintptr(len(iovs)), uintptr(offs_l), uintptr(offs_h), uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func pwritev2(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr, flags int) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(iovs) > 0 {
		_p0 = unsafe.Pointer(&iovs[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PWRITEV2, uintptr(fd), uintptr(_p0), uintptr(len(iovs)), uintptr(offs_l), uintptr(offs_h), uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func munmap(addr uintptr, length uintptr) (err error) {
	_, _, e1 := Syscall(SYS_MUNMAP, uintptr(addr), uintptr(length), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func mremap(oldaddr uintptr, oldlength uintptr, newlength uintptr, flags int, newaddr uintptr) (xaddr uintptr, err error) {
	r0, _, e1 := Syscall6(SYS_MREMAP, uintptr(oldaddr), uintptr(oldlength), uintptr(newlength), uintptr(flags), uintptr(newaddr), 0)
	xaddr = uintptr(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Madvise(b []byte, advice int) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_MADVISE, uintptr(_p0), uintptr(len(b)), uintptr(advice))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Mprotect(b []byte, prot int) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_MPROTECT, uintptr(_p0), uintptr(len(b)), uintptr(prot))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Mlock(b []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_MLOCK, uintptr(_p0), uintptr(len(b)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Mlockall(flags int) (err error) {
	_, _, e1 := Syscall(SYS_MLOCKALL, uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Msync(b []byte, flags int) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_MSYNC, uintptr(_p0), uintptr(len(b)), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Munlock(b []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_MUNLOCK, uintptr(_p0), uintptr(len(b)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Munlockall() (err error) {
	_, _, e1 := Syscall(SYS_MUNLOCKALL, 0, 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func faccessat(dirfd int, path string, mode uint32) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall(SYS_FACCESSAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Faccessat2(dirfd int, path string, mode uint32, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_FACCESSAT2, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(mode), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func nameToHandleAt(dirFD int, pathname string, fh *fileHandle, mountID *_C_int, flags int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(pathname)
	if err != nil {
		return
	}
	_, _, e1 := Syscall6(SYS_NAME_TO_HANDLE_AT, uintptr(dirFD), uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(fh)), uintptr(unsafe.Pointer(mountID)), uintptr(flags), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func openByHandleAt(mountFD int, fh *fileHandle, flags int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_OPEN_BY_HANDLE_AT, uintptr(mountFD), uintptr(unsafe.Pointer(fh)), uintptr(flags))
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ProcessVMReadv(pid int, localIov []Iovec, remoteIov []RemoteIovec, flags uint) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(localIov) > 0 {
		_p0 = unsafe.Pointer(&localIov[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	var _p1 unsafe.Pointer
	if len(remoteIov) > 0 {
		_p1 = unsafe.Pointer(&remoteIov[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PROCESS_VM_READV, uintptr(pid), uintptr(_p0), uintptr(len(localIov)), uintptr(_p1), uintptr(len(remoteIov)), uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func ProcessVMWritev(pid int, localIov []Iovec, remoteIov []RemoteIovec, flags uint) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(localIov) > 0 {
		_p0 = unsafe.Pointer(&localIov[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	var _p1 unsafe.Pointer
	if len(remoteIov) > 0 {
		_p1 = unsafe.Pointer(&remoteIov[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PROCESS_VM_WRITEV, uintptr(pid), uintptr(_p0), uintptr(len(localIov)), uintptr(_p1), uintptr(len(remoteIov)), uintptr(flags))
	n = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func PidfdOpen(pid int, flags int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_PIDFD_OPEN, uintptr(pid), uintptr(flags), 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func PidfdGetfd(pidfd int, targetfd int, flags int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_PIDFD_GETFD, uintptr(pidfd), uintptr(targetfd), uintptr(flags))
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func PidfdSendSignal(pidfd int, sig Signal, info *Siginfo, flags int) (err error) {
	_, _, e1 := Syscall6(SYS_PIDFD_SEND_SIGNAL, uintptr(pidfd), uintptr(sig), uintptr(unsafe.Pointer(info)), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func shmat(id int, addr uintptr, flag int) (ret uintptr, err error) {
	r0, _, e1 := Syscall(SYS_SHMAT, uintptr(id), uintptr(addr), uintptr(flag))
	ret = uintptr(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func shmctl(id int, cmd int, buf *SysvShmDesc) (result int, err error) {
	r0, _, e1 := Syscall(SYS_SHMCTL, uintptr(id), uintptr(cmd), uintptr(unsafe.Pointer(buf)))
	result = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func shmdt(addr uintptr) (err error) {
	_, _, e1 := Syscall(SYS_SHMDT, uintptr(addr), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func shmget(key int, size int, flag int) (id int, err error) {
	r0, _, e1 := Syscall(SYS_SHMGET, uintptr(key), uintptr(size), uintptr(flag))
	id = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func getitimer(which int, currValue *Itimerval) (err error) {
	_, _, e1 := Syscall(SYS_GETITIMER, uintptr(which), uintptr(unsafe.Pointer(currValue)), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func setitimer(which int, newValue *Itimerval, oldValue *Itimerval) (err error) {
	_, _, e1 := Syscall(SYS_SETITIMER, uintptr(which), uintptr(unsafe.Pointer(newValue)), uintptr(unsafe.Pointer(oldValue)))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func rtSigprocmask(how int, set *Sigset_t, oldset *Sigset_t, sigsetsize uintptr) (err error) {
	_, _, e1 := RawSyscall6(SYS_RT_SIGPROCMASK, uintptr(how), uintptr(unsafe.Pointer(set)), uintptr(unsafe.Pointer(oldset)), uintptr(sigsetsize), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func getresuid(ruid *_C_int, euid *_C_int, suid *_C_int) {
	RawSyscallNoError(SYS_GETRESUID, uintptr(unsafe.Pointer(ruid)), uintptr(unsafe.Pointer(euid)), uintptr(unsafe.Pointer(suid)))
	return
}



func getresgid(rgid *_C_int, egid *_C_int, sgid *_C_int) {
	RawSyscallNoError(SYS_GETRESGID, uintptr(unsafe.Pointer(rgid)), uintptr(unsafe.Pointer(egid)), uintptr(unsafe.Pointer(sgid)))
	return
}



func schedSetattr(pid int, attr *SchedAttr, flags uint) (err error) {
	_, _, e1 := Syscall(SYS_SCHED_SETATTR, uintptr(pid), uintptr(unsafe.Pointer(attr)), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func schedGetattr(pid int, attr *SchedAttr, size uint, flags uint) (err error) {
	_, _, e1 := Syscall6(SYS_SCHED_GETATTR, uintptr(pid), uintptr(unsafe.Pointer(attr)), uintptr(size), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Cachestat(fd uint, crange *CachestatRange, cstat *Cachestat_t, flags uint) (err error) {
	_, _, e1 := Syscall6(SYS_CACHESTAT, uintptr(fd), uintptr(unsafe.Pointer(crange)), uintptr(unsafe.Pointer(cstat)), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}



func Mseal(b []byte, flags uint) (err error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := Syscall(SYS_MSEAL, uintptr(_p0), uintptr(len(b)), uintptr(flags))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
