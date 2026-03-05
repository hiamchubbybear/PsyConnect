



//go:build linux && (mips64 || mips64le)

package cpu


const (
	
	hwcap_MIPS_MSA = 1 << 1
)

func doinit() {
	
	MIPS64X.HasMSA = isSet(hwCap, hwcap_MIPS_MSA)
}

func isSet(hwc uint, value uint) bool {
	return hwc&value != 0
}
