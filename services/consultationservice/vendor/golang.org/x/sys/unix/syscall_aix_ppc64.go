



//go:build aix && ppc64

package unix






func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: int64(sec), Usec: int32(usec)}
}

func (iov *Iovec) SetLen(length int) {
	iov.Len = uint64(length)
}

func (msghdr *Msghdr) SetControllen(length int) {
	msghdr.Controllen = uint32(length)
}

func (msghdr *Msghdr) SetIovlen(length int) {
	msghdr.Iovlen = int32(length)
}

func (cmsg *Cmsghdr) SetLen(length int) {
	cmsg.Len = uint32(length)
}






func fixStatTimFields(stat *Stat_t) {
	stat.Atim.Nsec >>= 32
	stat.Mtim.Nsec >>= 32
	stat.Ctim.Nsec >>= 32
}

func Fstat(fd int, stat *Stat_t) error {
	err := fstat(fd, stat)
	if err != nil {
		return err
	}
	fixStatTimFields(stat)
	return nil
}

func Fstatat(dirfd int, path string, stat *Stat_t, flags int) error {
	err := fstatat(dirfd, path, stat, flags)
	if err != nil {
		return err
	}
	fixStatTimFields(stat)
	return nil
}

func Lstat(path string, stat *Stat_t) error {
	err := lstat(path, stat)
	if err != nil {
		return err
	}
	fixStatTimFields(stat)
	return nil
}

func Stat(path string, statptr *Stat_t) error {
	err := stat(path, statptr)
	if err != nil {
		return err
	}
	fixStatTimFields(statptr)
	return nil
}
