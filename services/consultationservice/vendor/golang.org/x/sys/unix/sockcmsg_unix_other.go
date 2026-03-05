



//go:build aix || darwin || freebsd || linux || netbsd || openbsd || solaris || zos

package unix

import (
	"runtime"
)


func cmsgAlignOf(salen int) int {
	salign := SizeofPtr

	
	
	switch runtime.GOOS {
	case "aix":
		
		salign = 1
	case "darwin", "ios", "illumos", "solaris":
		
		
		
		if SizeofPtr == 8 {
			salign = 4
		}
	case "netbsd", "openbsd":
		
		if runtime.GOARCH == "arm" {
			salign = 8
		}
		
		if runtime.GOOS == "netbsd" && runtime.GOARCH == "arm64" {
			salign = 16
		}
	case "zos":
		
		
		salign = SizeofInt
	}

	return (salen + salign - 1) & ^(salign - 1)
}
