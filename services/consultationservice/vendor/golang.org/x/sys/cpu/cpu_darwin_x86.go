



//go:build darwin && amd64 && gc

package cpu



















func darwinSupportsAVX512() bool {
	return darwinSysctlEnabled([]byte("hw.optional.avx512f\x00")) && darwinKernelVersionCheck(21, 3, 0)
}


func darwinKernelVersionCheck(major, minor, patch int) bool {
	var release [256]byte
	err := darwinOSRelease(&release)
	if err != nil {
		return false
	}

	var mmp [3]int
	c := 0
Loop:
	for _, b := range release[:] {
		switch {
		case b >= '0' && b <= '9':
			mmp[c] = 10*mmp[c] + int(b-'0')
		case b == '.':
			c++
			if c > 2 {
				return false
			}
		case b == 0:
			break Loop
		default:
			return false
		}
	}
	if c != 2 {
		return false
	}
	return mmp[0] > major || mmp[0] == major && (mmp[1] > minor || mmp[1] == minor && mmp[2] >= patch)
}
