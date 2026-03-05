



//go:build (darwin && !ios) || zos

package unix



func SysvShmCtl(id, cmd int, desc *SysvShmDesc) (result int, err error) {
	return shmctl(id, cmd, desc)
}
