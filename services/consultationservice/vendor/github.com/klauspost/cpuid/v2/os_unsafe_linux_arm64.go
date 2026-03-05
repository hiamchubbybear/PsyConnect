

//go:build !nounsafe
// +build !nounsafe

package cpuid

import _ "unsafe" 

//go:linkname hwcap internal/cpu.HWCap
var hwcap uint
