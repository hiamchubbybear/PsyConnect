



//go:build (linux && 386) || (linux && arm) || (linux && mips) || (linux && mipsle) || (linux && ppc)

package unix

func init() {
	
	
	fcntl64Syscall = SYS_FCNTL64
}
