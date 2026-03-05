



package unix

import "unsafe"


func FcntlInt(fd uintptr, cmd, arg int) (int, error) {
	return fcntl(int(fd), cmd, arg)
}


func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	_, err := fcntl(int(fd), cmd, int(uintptr(unsafe.Pointer(lk))))
	return err
}


func FcntlFstore(fd uintptr, cmd int, fstore *Fstore_t) error {
	_, err := fcntl(int(fd), cmd, int(uintptr(unsafe.Pointer(fstore))))
	return err
}
