



//go:build (386 || amd64 || amd64p32) && gc

package cpu



func cpuid(eaxArg, ecxArg uint32) (eax, ebx, ecx, edx uint32)



func xgetbv() (eax, edx uint32)
