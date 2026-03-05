



//go:build loong64 && linux

package unix

import "unsafe"















func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) {
	var ts *Timespec
	if timeout != nil {
		ts = &Timespec{Sec: timeout.Sec, Nsec: timeout.Usec * 1000}
	}
	return pselect6(nfd, r, w, e, ts, nil)
}







func timespecFromStatxTimestamp(x StatxTimestamp) Timespec {
	return Timespec{
		Sec:  x.Sec,
		Nsec: int64(x.Nsec),
	}
}

func Fstatat(fd int, path string, stat *Stat_t, flags int) error {
	var r Statx_t
	
	if err := Statx(fd, path, AT_NO_AUTOMOUNT|flags, STATX_BASIC_STATS, &r); err != nil {
		return err
	}

	stat.Dev = Mkdev(r.Dev_major, r.Dev_minor)
	stat.Ino = r.Ino
	stat.Mode = uint32(r.Mode)
	stat.Nlink = r.Nlink
	stat.Uid = r.Uid
	stat.Gid = r.Gid
	stat.Rdev = Mkdev(r.Rdev_major, r.Rdev_minor)
	
	
	stat.Size = int64(r.Size)
	stat.Blksize = int32(r.Blksize)
	stat.Blocks = int64(r.Blocks)
	stat.Atim = timespecFromStatxTimestamp(r.Atime)
	stat.Mtim = timespecFromStatxTimestamp(r.Mtime)
	stat.Ctim = timespecFromStatxTimestamp(r.Ctime)

	return nil
}

func Fstat(fd int, stat *Stat_t) (err error) {
	return Fstatat(fd, "", stat, AT_EMPTY_PATH)
}

func Stat(path string, stat *Stat_t) (err error) {
	return Fstatat(AT_FDCWD, path, stat, 0)
}

func Lchown(path string, uid int, gid int) (err error) {
	return Fchownat(AT_FDCWD, path, uid, gid, AT_SYMLINK_NOFOLLOW)
}

func Lstat(path string, stat *Stat_t) (err error) {
	return Fstatat(AT_FDCWD, path, stat, AT_SYMLINK_NOFOLLOW)
}





func Ustat(dev int, ubuf *Ustat_t) (err error) {
	return ENOSYS
}




















func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: sec, Usec: usec}
}

func Getrlimit(resource int, rlim *Rlimit) (err error) {
	err = Prlimit(0, resource, nil, rlim)
	return
}

func futimesat(dirfd int, path string, tv *[2]Timeval) (err error) {
	if tv == nil {
		return utimensat(dirfd, path, nil, 0)
	}

	ts := []Timespec{
		NsecToTimespec(TimevalToNsec(tv[0])),
		NsecToTimespec(TimevalToNsec(tv[1])),
	}
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
}

func Time(t *Time_t) (Time_t, error) {
	var tv Timeval
	err := Gettimeofday(&tv)
	if err != nil {
		return 0, err
	}
	if t != nil {
		*t = Time_t(tv.Sec)
	}
	return Time_t(tv.Sec), nil
}

func Utime(path string, buf *Utimbuf) error {
	tv := []Timeval{
		{Sec: buf.Actime},
		{Sec: buf.Modtime},
	}
	return Utimes(path, tv)
}

func utimes(path string, tv *[2]Timeval) (err error) {
	if tv == nil {
		return utimensat(AT_FDCWD, path, nil, 0)
	}

	ts := []Timespec{
		NsecToTimespec(TimevalToNsec(tv[0])),
		NsecToTimespec(TimevalToNsec(tv[1])),
	}
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
}

func (r *PtraceRegs) PC() uint64 { return r.Era }

func (r *PtraceRegs) SetPC(era uint64) { r.Era = era }

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint64(length)
}

func (msghdr *Msghdr) SetIovlen(length int) {
	msghdr.Iovlen = uint64(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint64(length)
}

func (rsa *RawSockaddrNFCLLCP) SetServiceNameLen(length int) {
	rsa.Service_name_len = uint64(length)
}

func Pause() error {
	_, err := ppoll(nil, 0, nil, nil)
	return err
}

func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error) {
	return Renameat2(olddirfd, oldpath, newdirfd, newpath, 0)
}



func KexecFileLoad(kernelFd int, initrdFd int, cmdline string, flags int) error {
	cmdlineLen := len(cmdline)
	if cmdlineLen > 0 {
		
		
		
		cmdlineLen++
	}
	return kexecFileLoad(kernelFd, initrdFd, cmdlineLen, cmdline, flags)
}

const SYS_FSTATAT = SYS_NEWFSTATAT
