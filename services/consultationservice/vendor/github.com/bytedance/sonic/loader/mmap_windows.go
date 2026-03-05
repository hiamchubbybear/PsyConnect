//go:build windows
// +build windows





package loader

import (
    `syscall`
    `unsafe`
)

const (
    MEM_COMMIT  = 0x00001000
    MEM_RESERVE = 0x00002000
)

var (
    libKernel32                = syscall.NewLazyDLL("KERNEL32.DLL")
    libKernel32_VirtualAlloc   = libKernel32.NewProc("VirtualAlloc")
    libKernel32_VirtualProtect = libKernel32.NewProc("VirtualProtect")
)

func mmap(nb int) uintptr {
    addr, err := winapi_VirtualAlloc(0, nb, MEM_COMMIT|MEM_RESERVE, syscall.PAGE_READWRITE)
    if err != nil {
        panic(err)
    }
    return addr
}

func mprotect(p uintptr, nb int) (oldProtect int) {
    err := winapi_VirtualProtect(p, nb, syscall.PAGE_EXECUTE_READ, &oldProtect)
    if err != nil {
        panic(err)
    }
    return
}



func winapi_VirtualAlloc(lpAddr uintptr, dwSize int, flAllocationType int, flProtect int) (uintptr, error) {
    r1, _, err := libKernel32_VirtualAlloc.Call(
        lpAddr,
        uintptr(dwSize),
        uintptr(flAllocationType),
        uintptr(flProtect),
    )
    if r1 == 0 {
        return 0, err
    }
    return r1, nil
}



func winapi_VirtualProtect(lpAddr uintptr, dwSize int, flNewProtect int, lpflOldProtect *int) error {
    r1, _, err := libKernel32_VirtualProtect.Call(
        lpAddr,
        uintptr(dwSize),
        uintptr(flNewProtect),
        uintptr(unsafe.Pointer(lpflOldProtect)),
    )
    if r1 == 0 {
        return err
    }
    return nil
}
