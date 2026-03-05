



//go:build linux

package unix

import "runtime"



func SysvShmCtl(id, cmd int, desc *SysvShmDesc) (result int, err error) {
	if runtime.GOARCH == "arm" ||
		runtime.GOARCH == "mips64" || runtime.GOARCH == "mips64le" {
		cmd |= ipc_64
	}

	return shmctl(id, cmd, desc)
}
