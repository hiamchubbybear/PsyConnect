


//go:build aix && ppc64 && gccgo

package unix


import "C"
import (
	"syscall"
	"unsafe"
)



func callutimes(_p0 uintptr, times uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.utimes(C.uintptr_t(_p0), C.uintptr_t(times)))
	e1 = syscall.GetErrno()
	return
}



func callutimensat(dirfd int, _p0 uintptr, times uintptr, flag int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.utimensat(C.int(dirfd), C.uintptr_t(_p0), C.uintptr_t(times), C.int(flag)))
	e1 = syscall.GetErrno()
	return
}



func callgetcwd(_p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getcwd(C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callaccept(s int, rsa uintptr, addrlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.accept(C.int(s), C.uintptr_t(rsa), C.uintptr_t(addrlen)))
	e1 = syscall.GetErrno()
	return
}



func callgetdirent(fd int, _p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getdirent(C.int(fd), C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callwait4(pid int, status uintptr, options int, rusage uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.wait4(C.int(pid), C.uintptr_t(status), C.int(options), C.uintptr_t(rusage)))
	e1 = syscall.GetErrno()
	return
}



func callioctl(fd int, req int, arg uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.ioctl(C.int(fd), C.int(req), C.uintptr_t(arg)))
	e1 = syscall.GetErrno()
	return
}



func callioctl_ptr(fd int, req int, arg unsafe.Pointer) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.ioctl(C.int(fd), C.int(req), C.uintptr_t(uintptr(arg))))
	e1 = syscall.GetErrno()
	return
}



func callfcntl(fd uintptr, cmd int, arg uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fcntl(C.uintptr_t(fd), C.int(cmd), C.uintptr_t(arg)))
	e1 = syscall.GetErrno()
	return
}



func callfsync_range(fd int, how int, start int64, length int64) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fsync_range(C.int(fd), C.int(how), C.longlong(start), C.longlong(length)))
	e1 = syscall.GetErrno()
	return
}



func callacct(_p0 uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.acct(C.uintptr_t(_p0)))
	e1 = syscall.GetErrno()
	return
}



func callchdir(_p0 uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.chdir(C.uintptr_t(_p0)))
	e1 = syscall.GetErrno()
	return
}



func callchroot(_p0 uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.chroot(C.uintptr_t(_p0)))
	e1 = syscall.GetErrno()
	return
}



func callclose(fd int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.close(C.int(fd)))
	e1 = syscall.GetErrno()
	return
}



func calldup(oldfd int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.dup(C.int(oldfd)))
	e1 = syscall.GetErrno()
	return
}



func callexit(code int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.exit(C.int(code)))
	e1 = syscall.GetErrno()
	return
}



func callfaccessat(dirfd int, _p0 uintptr, mode uint32, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.faccessat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callfchdir(fd int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fchdir(C.int(fd)))
	e1 = syscall.GetErrno()
	return
}



func callfchmod(fd int, mode uint32) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fchmod(C.int(fd), C.uint(mode)))
	e1 = syscall.GetErrno()
	return
}



func callfchmodat(dirfd int, _p0 uintptr, mode uint32, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fchmodat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callfchownat(dirfd int, _p0 uintptr, uid int, gid int, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fchownat(C.int(dirfd), C.uintptr_t(_p0), C.int(uid), C.int(gid), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callfdatasync(fd int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fdatasync(C.int(fd)))
	e1 = syscall.GetErrno()
	return
}



func callgetpgid(pid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getpgid(C.int(pid)))
	e1 = syscall.GetErrno()
	return
}



func callgetpgrp() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getpgrp())
	e1 = syscall.GetErrno()
	return
}



func callgetpid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getpid())
	e1 = syscall.GetErrno()
	return
}



func callgetppid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getppid())
	e1 = syscall.GetErrno()
	return
}



func callgetpriority(which int, who int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getpriority(C.int(which), C.int(who)))
	e1 = syscall.GetErrno()
	return
}



func callgetrusage(who int, rusage uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getrusage(C.int(who), C.uintptr_t(rusage)))
	e1 = syscall.GetErrno()
	return
}



func callgetsid(pid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getsid(C.int(pid)))
	e1 = syscall.GetErrno()
	return
}



func callkill(pid int, sig int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.kill(C.int(pid), C.int(sig)))
	e1 = syscall.GetErrno()
	return
}



func callsyslog(typ int, _p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.syslog(C.int(typ), C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callmkdir(dirfd int, _p0 uintptr, mode uint32) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mkdir(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode)))
	e1 = syscall.GetErrno()
	return
}



func callmkdirat(dirfd int, _p0 uintptr, mode uint32) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mkdirat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode)))
	e1 = syscall.GetErrno()
	return
}



func callmkfifo(_p0 uintptr, mode uint32) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mkfifo(C.uintptr_t(_p0), C.uint(mode)))
	e1 = syscall.GetErrno()
	return
}



func callmknod(_p0 uintptr, mode uint32, dev int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mknod(C.uintptr_t(_p0), C.uint(mode), C.int(dev)))
	e1 = syscall.GetErrno()
	return
}



func callmknodat(dirfd int, _p0 uintptr, mode uint32, dev int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mknodat(C.int(dirfd), C.uintptr_t(_p0), C.uint(mode), C.int(dev)))
	e1 = syscall.GetErrno()
	return
}



func callnanosleep(time uintptr, leftover uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.nanosleep(C.uintptr_t(time), C.uintptr_t(leftover)))
	e1 = syscall.GetErrno()
	return
}



func callopen64(_p0 uintptr, mode int, perm uint32) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.open64(C.uintptr_t(_p0), C.int(mode), C.uint(perm)))
	e1 = syscall.GetErrno()
	return
}



func callopenat(dirfd int, _p0 uintptr, flags int, mode uint32) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.openat(C.int(dirfd), C.uintptr_t(_p0), C.int(flags), C.uint(mode)))
	e1 = syscall.GetErrno()
	return
}



func callread(fd int, _p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.read(C.int(fd), C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callreadlink(_p0 uintptr, _p1 uintptr, _lenp1 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.readlink(C.uintptr_t(_p0), C.uintptr_t(_p1), C.size_t(_lenp1)))
	e1 = syscall.GetErrno()
	return
}



func callrenameat(olddirfd int, _p0 uintptr, newdirfd int, _p1 uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.renameat(C.int(olddirfd), C.uintptr_t(_p0), C.int(newdirfd), C.uintptr_t(_p1)))
	e1 = syscall.GetErrno()
	return
}



func callsetdomainname(_p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setdomainname(C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callsethostname(_p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.sethostname(C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callsetpgid(pid int, pgid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setpgid(C.int(pid), C.int(pgid)))
	e1 = syscall.GetErrno()
	return
}



func callsetsid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setsid())
	e1 = syscall.GetErrno()
	return
}



func callsettimeofday(tv uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.settimeofday(C.uintptr_t(tv)))
	e1 = syscall.GetErrno()
	return
}



func callsetuid(uid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setuid(C.int(uid)))
	e1 = syscall.GetErrno()
	return
}



func callsetgid(uid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setgid(C.int(uid)))
	e1 = syscall.GetErrno()
	return
}



func callsetpriority(which int, who int, prio int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setpriority(C.int(which), C.int(who), C.int(prio)))
	e1 = syscall.GetErrno()
	return
}



func callstatx(dirfd int, _p0 uintptr, flags int, mask int, stat uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.statx(C.int(dirfd), C.uintptr_t(_p0), C.int(flags), C.int(mask), C.uintptr_t(stat)))
	e1 = syscall.GetErrno()
	return
}



func callsync() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.sync())
	e1 = syscall.GetErrno()
	return
}



func calltimes(tms uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.times(C.uintptr_t(tms)))
	e1 = syscall.GetErrno()
	return
}



func callumask(mask int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.umask(C.int(mask)))
	e1 = syscall.GetErrno()
	return
}



func calluname(buf uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.uname(C.uintptr_t(buf)))
	e1 = syscall.GetErrno()
	return
}



func callunlink(_p0 uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.unlink(C.uintptr_t(_p0)))
	e1 = syscall.GetErrno()
	return
}



func callunlinkat(dirfd int, _p0 uintptr, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.unlinkat(C.int(dirfd), C.uintptr_t(_p0), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callustat(dev int, ubuf uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.ustat(C.int(dev), C.uintptr_t(ubuf)))
	e1 = syscall.GetErrno()
	return
}



func callwrite(fd int, _p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.write(C.int(fd), C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func calldup2(oldfd int, newfd int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.dup2(C.int(oldfd), C.int(newfd)))
	e1 = syscall.GetErrno()
	return
}



func callposix_fadvise64(fd int, offset int64, length int64, advice int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.posix_fadvise64(C.int(fd), C.longlong(offset), C.longlong(length), C.int(advice)))
	e1 = syscall.GetErrno()
	return
}



func callfchown(fd int, uid int, gid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fchown(C.int(fd), C.int(uid), C.int(gid)))
	e1 = syscall.GetErrno()
	return
}



func callfstat(fd int, stat uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fstat(C.int(fd), C.uintptr_t(stat)))
	e1 = syscall.GetErrno()
	return
}



func callfstatat(dirfd int, _p0 uintptr, stat uintptr, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fstatat(C.int(dirfd), C.uintptr_t(_p0), C.uintptr_t(stat), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callfstatfs(fd int, buf uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.fstatfs(C.int(fd), C.uintptr_t(buf)))
	e1 = syscall.GetErrno()
	return
}



func callftruncate(fd int, length int64) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.ftruncate(C.int(fd), C.longlong(length)))
	e1 = syscall.GetErrno()
	return
}



func callgetegid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getegid())
	e1 = syscall.GetErrno()
	return
}



func callgeteuid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.geteuid())
	e1 = syscall.GetErrno()
	return
}



func callgetgid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getgid())
	e1 = syscall.GetErrno()
	return
}



func callgetuid() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getuid())
	e1 = syscall.GetErrno()
	return
}



func calllchown(_p0 uintptr, uid int, gid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.lchown(C.uintptr_t(_p0), C.int(uid), C.int(gid)))
	e1 = syscall.GetErrno()
	return
}



func calllisten(s int, n int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.listen(C.int(s), C.int(n)))
	e1 = syscall.GetErrno()
	return
}



func calllstat(_p0 uintptr, stat uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.lstat(C.uintptr_t(_p0), C.uintptr_t(stat)))
	e1 = syscall.GetErrno()
	return
}



func callpause() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.pause())
	e1 = syscall.GetErrno()
	return
}



func callpread64(fd int, _p0 uintptr, _lenp0 int, offset int64) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.pread64(C.int(fd), C.uintptr_t(_p0), C.size_t(_lenp0), C.longlong(offset)))
	e1 = syscall.GetErrno()
	return
}



func callpwrite64(fd int, _p0 uintptr, _lenp0 int, offset int64) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.pwrite64(C.int(fd), C.uintptr_t(_p0), C.size_t(_lenp0), C.longlong(offset)))
	e1 = syscall.GetErrno()
	return
}



func callselect(nfd int, r uintptr, w uintptr, e uintptr, timeout uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.c_select(C.int(nfd), C.uintptr_t(r), C.uintptr_t(w), C.uintptr_t(e), C.uintptr_t(timeout)))
	e1 = syscall.GetErrno()
	return
}



func callpselect(nfd int, r uintptr, w uintptr, e uintptr, timeout uintptr, sigmask uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.pselect(C.int(nfd), C.uintptr_t(r), C.uintptr_t(w), C.uintptr_t(e), C.uintptr_t(timeout), C.uintptr_t(sigmask)))
	e1 = syscall.GetErrno()
	return
}



func callsetregid(rgid int, egid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setregid(C.int(rgid), C.int(egid)))
	e1 = syscall.GetErrno()
	return
}



func callsetreuid(ruid int, euid int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setreuid(C.int(ruid), C.int(euid)))
	e1 = syscall.GetErrno()
	return
}



func callshutdown(fd int, how int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.shutdown(C.int(fd), C.int(how)))
	e1 = syscall.GetErrno()
	return
}



func callsplice(rfd int, roff uintptr, wfd int, woff uintptr, len int, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.splice(C.int(rfd), C.uintptr_t(roff), C.int(wfd), C.uintptr_t(woff), C.int(len), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callstat(_p0 uintptr, statptr uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.stat(C.uintptr_t(_p0), C.uintptr_t(statptr)))
	e1 = syscall.GetErrno()
	return
}



func callstatfs(_p0 uintptr, buf uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.statfs(C.uintptr_t(_p0), C.uintptr_t(buf)))
	e1 = syscall.GetErrno()
	return
}



func calltruncate(_p0 uintptr, length int64) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.truncate(C.uintptr_t(_p0), C.longlong(length)))
	e1 = syscall.GetErrno()
	return
}



func callbind(s int, addr uintptr, addrlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.bind(C.int(s), C.uintptr_t(addr), C.uintptr_t(addrlen)))
	e1 = syscall.GetErrno()
	return
}



func callconnect(s int, addr uintptr, addrlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.connect(C.int(s), C.uintptr_t(addr), C.uintptr_t(addrlen)))
	e1 = syscall.GetErrno()
	return
}



func callgetgroups(n int, list uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getgroups(C.int(n), C.uintptr_t(list)))
	e1 = syscall.GetErrno()
	return
}



func callsetgroups(n int, list uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setgroups(C.int(n), C.uintptr_t(list)))
	e1 = syscall.GetErrno()
	return
}



func callgetsockopt(s int, level int, name int, val uintptr, vallen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getsockopt(C.int(s), C.int(level), C.int(name), C.uintptr_t(val), C.uintptr_t(vallen)))
	e1 = syscall.GetErrno()
	return
}



func callsetsockopt(s int, level int, name int, val uintptr, vallen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.setsockopt(C.int(s), C.int(level), C.int(name), C.uintptr_t(val), C.uintptr_t(vallen)))
	e1 = syscall.GetErrno()
	return
}



func callsocket(domain int, typ int, proto int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.socket(C.int(domain), C.int(typ), C.int(proto)))
	e1 = syscall.GetErrno()
	return
}



func callsocketpair(domain int, typ int, proto int, fd uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.socketpair(C.int(domain), C.int(typ), C.int(proto), C.uintptr_t(fd)))
	e1 = syscall.GetErrno()
	return
}



func callgetpeername(fd int, rsa uintptr, addrlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getpeername(C.int(fd), C.uintptr_t(rsa), C.uintptr_t(addrlen)))
	e1 = syscall.GetErrno()
	return
}



func callgetsockname(fd int, rsa uintptr, addrlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getsockname(C.int(fd), C.uintptr_t(rsa), C.uintptr_t(addrlen)))
	e1 = syscall.GetErrno()
	return
}



func callrecvfrom(fd int, _p0 uintptr, _lenp0 int, flags int, from uintptr, fromlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.recvfrom(C.int(fd), C.uintptr_t(_p0), C.size_t(_lenp0), C.int(flags), C.uintptr_t(from), C.uintptr_t(fromlen)))
	e1 = syscall.GetErrno()
	return
}



func callsendto(s int, _p0 uintptr, _lenp0 int, flags int, to uintptr, addrlen uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.sendto(C.int(s), C.uintptr_t(_p0), C.size_t(_lenp0), C.int(flags), C.uintptr_t(to), C.uintptr_t(addrlen)))
	e1 = syscall.GetErrno()
	return
}



func callnrecvmsg(s int, msg uintptr, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.nrecvmsg(C.int(s), C.uintptr_t(msg), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callnsendmsg(s int, msg uintptr, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.nsendmsg(C.int(s), C.uintptr_t(msg), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callmunmap(addr uintptr, length uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.munmap(C.uintptr_t(addr), C.uintptr_t(length)))
	e1 = syscall.GetErrno()
	return
}



func callmadvise(_p0 uintptr, _lenp0 int, advice int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.madvise(C.uintptr_t(_p0), C.size_t(_lenp0), C.int(advice)))
	e1 = syscall.GetErrno()
	return
}



func callmprotect(_p0 uintptr, _lenp0 int, prot int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mprotect(C.uintptr_t(_p0), C.size_t(_lenp0), C.int(prot)))
	e1 = syscall.GetErrno()
	return
}



func callmlock(_p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mlock(C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callmlockall(flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mlockall(C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callmsync(_p0 uintptr, _lenp0 int, flags int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.msync(C.uintptr_t(_p0), C.size_t(_lenp0), C.int(flags)))
	e1 = syscall.GetErrno()
	return
}



func callmunlock(_p0 uintptr, _lenp0 int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.munlock(C.uintptr_t(_p0), C.size_t(_lenp0)))
	e1 = syscall.GetErrno()
	return
}



func callmunlockall() (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.munlockall())
	e1 = syscall.GetErrno()
	return
}



func callpipe(p uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.pipe(C.uintptr_t(p)))
	e1 = syscall.GetErrno()
	return
}



func callpoll(fds uintptr, nfds int, timeout int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.poll(C.uintptr_t(fds), C.int(nfds), C.int(timeout)))
	e1 = syscall.GetErrno()
	return
}



func callgettimeofday(tv uintptr, tzp uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.gettimeofday(C.uintptr_t(tv), C.uintptr_t(tzp)))
	e1 = syscall.GetErrno()
	return
}



func calltime(t uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.time(C.uintptr_t(t)))
	e1 = syscall.GetErrno()
	return
}



func callutime(_p0 uintptr, buf uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.utime(C.uintptr_t(_p0), C.uintptr_t(buf)))
	e1 = syscall.GetErrno()
	return
}



func callgetsystemcfg(label int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getsystemcfg(C.int(label)))
	e1 = syscall.GetErrno()
	return
}



func callumount(_p0 uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.umount(C.uintptr_t(_p0)))
	e1 = syscall.GetErrno()
	return
}



func callgetrlimit(resource int, rlim uintptr) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.getrlimit(C.int(resource), C.uintptr_t(rlim)))
	e1 = syscall.GetErrno()
	return
}



func calllseek(fd int, offset int64, whence int) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.lseek(C.int(fd), C.longlong(offset), C.int(whence)))
	e1 = syscall.GetErrno()
	return
}



func callmmap64(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (r1 uintptr, e1 Errno) {
	r1 = uintptr(C.mmap64(C.uintptr_t(addr), C.uintptr_t(length), C.int(prot), C.int(flags), C.int(fd), C.longlong(offset)))
	e1 = syscall.GetErrno()
	return
}
