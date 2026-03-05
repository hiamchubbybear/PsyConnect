



//go:build (darwin && !ios) || linux || zos

package unix

import "unsafe"



func SysvShmAttach(id int, addr uintptr, flag int) ([]byte, error) {
	addr, errno := shmat(id, addr, flag)
	if errno != nil {
		return nil, errno
	}

	
	var info SysvShmDesc

	_, err := SysvShmCtl(id, IPC_STAT, &info)
	if err != nil {
		

		
		shmdt(addr)
		return nil, err
	}

	
	b := unsafe.Slice((*byte)(unsafe.Pointer(addr)), int(info.Segsz))
	return b, nil
}




func SysvShmDetach(data []byte) error {
	if len(data) == 0 {
		return EINVAL
	}

	return shmdt(uintptr(unsafe.Pointer(&data[0])))
}



func SysvShmGet(key, size, flag int) (id int, err error) {
	return shmget(key, size, flag)
}
