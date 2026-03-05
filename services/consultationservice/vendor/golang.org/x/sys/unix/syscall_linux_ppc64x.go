



//go:build linux && (ppc64 || ppc64le)

package unix
























































func setTimespec(sec, nsec int64) Timespec {
	return Timespec{Sec: sec, Nsec: nsec}
}

func setTimeval(sec, usec int64) Timeval {
	return Timeval{Sec: sec, Usec: usec}
}

func (r *PtraceRegs) PC() uint64 { return r.Nip }

func (r *PtraceRegs) SetPC(pc uint64) { r.Nip = pc }

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



func SyncFileRange(fd int, off int64, n int64, flags int) error {
	
	
	return syncFileRange2(fd, flags, off, n)
}



func KexecFileLoad(kernelFd int, initrdFd int, cmdline string, flags int) error {
	cmdlineLen := len(cmdline)
	if cmdlineLen > 0 {
		
		
		
		cmdlineLen++
	}
	return kexecFileLoad(kernelFd, initrdFd, cmdlineLen, cmdline, flags)
}
