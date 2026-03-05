//go:build windows && !appengine
// +build windows,!appengine

package isatty

import (
	"errors"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	objectNameInfo uintptr = 1
	fileNameInfo           = 2
	fileTypePipe           = 3
)

var (
	kernel32                         = syscall.NewLazyDLL("kernel32.dll")
	ntdll                            = syscall.NewLazyDLL("ntdll.dll")
	procGetConsoleMode               = kernel32.NewProc("GetConsoleMode")
	procGetFileInformationByHandleEx = kernel32.NewProc("GetFileInformationByHandleEx")
	procGetFileType                  = kernel32.NewProc("GetFileType")
	procNtQueryObject                = ntdll.NewProc("NtQueryObject")
)

func init() {
	
	if procGetFileInformationByHandleEx.Find() != nil {
		procGetFileInformationByHandleEx = nil
	}
}


func IsTerminal(fd uintptr) bool {
	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}




func isCygwinPipeName(name string) bool {
	token := strings.Split(name, "-")
	if len(token) < 5 {
		return false
	}

	if token[0] != `\msys` &&
		token[0] != `\cygwin` &&
		token[0] != `\Device\NamedPipe\msys` &&
		token[0] != `\Device\NamedPipe\cygwin` {
		return false
	}

	if token[1] == "" {
		return false
	}

	if !strings.HasPrefix(token[2], "pty") {
		return false
	}

	if token[3] != `from` && token[3] != `to` {
		return false
	}

	if token[4] != "master" {
		return false
	}

	return true
}






func getFileNameByHandle(fd uintptr) (string, error) {
	if procNtQueryObject == nil {
		return "", errors.New("ntdll.dll: NtQueryObject not supported")
	}

	var buf [4 + syscall.MAX_PATH]uint16
	var result int
	r, _, e := syscall.Syscall6(procNtQueryObject.Addr(), 5,
		fd, objectNameInfo, uintptr(unsafe.Pointer(&buf)), uintptr(2*len(buf)), uintptr(unsafe.Pointer(&result)), 0)
	if r != 0 {
		return "", e
	}
	return string(utf16.Decode(buf[4 : 4+buf[0]/2])), nil
}



func IsCygwinTerminal(fd uintptr) bool {
	if procGetFileInformationByHandleEx == nil {
		name, err := getFileNameByHandle(fd)
		if err != nil {
			return false
		}
		return isCygwinPipeName(name)
	}

	
	ft, _, e := syscall.Syscall(procGetFileType.Addr(), 1, fd, 0, 0)
	if ft != fileTypePipe || e != 0 {
		return false
	}

	var buf [2 + syscall.MAX_PATH]uint16
	r, _, e := syscall.Syscall6(procGetFileInformationByHandleEx.Addr(),
		4, fd, fileNameInfo, uintptr(unsafe.Pointer(&buf)),
		uintptr(len(buf)*2), 0, 0)
	if r == 0 || e != 0 {
		return false
	}

	l := *(*uint32)(unsafe.Pointer(&buf))
	return isCygwinPipeName(string(utf16.Decode(buf[2 : 2+l/2])))
}
