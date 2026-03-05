//go:build !windows
// +build !windows



package loader

import (
    `syscall`
)

const (
    _AP = syscall.MAP_ANON  | syscall.MAP_PRIVATE
    _RX = syscall.PROT_READ | syscall.PROT_EXEC
    _RW = syscall.PROT_READ | syscall.PROT_WRITE
)


func mmap(nb int) uintptr {
    if m, _, e := syscall.RawSyscall6(syscall.SYS_MMAP, 0, uintptr(nb), _RW, _AP, 0, 0); e != 0 {
        panic(e)
    } else {
        return m
    }
}

func mprotect(p uintptr, nb int) {
    if _, _, err := syscall.RawSyscall(syscall.SYS_MPROTECT, p, uintptr(nb), _RX); err != 0 {
        panic(err)
    }
}
