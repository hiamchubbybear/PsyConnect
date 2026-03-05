

package cpu

import (
    `fmt`
    `os`

    `github.com/klauspost/cpuid/v2`
)

var (
    HasAVX2 = cpuid.CPU.Has(cpuid.AVX2)
    HasSSE = cpuid.CPU.Has(cpuid.SSE)
)

func init() {
    switch v := os.Getenv("SONIC_MODE"); v {
        case ""       : break
        case "auto"   : break
        case "noavx"  : HasAVX2 = false
        
        case "noavx2" : HasAVX2 = false
        default       : panic(fmt.Sprintf("invalid mode: '%s', should be one of 'auto', 'noavx', 'noavx2'", v))
    }
}
