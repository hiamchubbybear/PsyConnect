//go:build !appengine && !noasm && gc
// +build !appengine,!noasm,gc

package s2

const hasAmd64Asm = true









func encodeBlock(dst, src []byte) (d int) {
	const (
		
		limit12B = 16 << 10
		
		limit10B = 4 << 10
		
		limit8B = 512
	)

	if len(src) >= 4<<20 {
		return encodeBlockAsm(dst, src)
	}
	if len(src) >= limit12B {
		return encodeBlockAsm4MB(dst, src)
	}
	if len(src) >= limit10B {
		return encodeBlockAsm12B(dst, src)
	}
	if len(src) >= limit8B {
		return encodeBlockAsm10B(dst, src)
	}
	if len(src) < minNonLiteralBlockSize {
		return 0
	}
	return encodeBlockAsm8B(dst, src)
}









func encodeBlockBetter(dst, src []byte) (d int) {
	const (
		
		limit12B = 16 << 10
		
		limit10B = 4 << 10
		
		limit8B = 512
	)

	if len(src) > 4<<20 {
		return encodeBetterBlockAsm(dst, src)
	}
	if len(src) >= limit12B {
		return encodeBetterBlockAsm4MB(dst, src)
	}
	if len(src) >= limit10B {
		return encodeBetterBlockAsm12B(dst, src)
	}
	if len(src) >= limit8B {
		return encodeBetterBlockAsm10B(dst, src)
	}
	if len(src) < minNonLiteralBlockSize {
		return 0
	}
	return encodeBetterBlockAsm8B(dst, src)
}









func encodeBlockSnappy(dst, src []byte) (d int) {
	const (
		
		limit12B = 16 << 10
		
		limit10B = 4 << 10
		
		limit8B = 512
	)
	if len(src) >= 64<<10 {
		return encodeSnappyBlockAsm(dst, src)
	}
	if len(src) >= limit12B {
		return encodeSnappyBlockAsm64K(dst, src)
	}
	if len(src) >= limit10B {
		return encodeSnappyBlockAsm12B(dst, src)
	}
	if len(src) >= limit8B {
		return encodeSnappyBlockAsm10B(dst, src)
	}
	if len(src) < minNonLiteralBlockSize {
		return 0
	}
	return encodeSnappyBlockAsm8B(dst, src)
}









func encodeBlockBetterSnappy(dst, src []byte) (d int) {
	const (
		
		limit12B = 16 << 10
		
		limit10B = 4 << 10
		
		limit8B = 512
	)
	if len(src) >= 64<<10 {
		return encodeSnappyBetterBlockAsm(dst, src)
	}
	if len(src) >= limit12B {
		return encodeSnappyBetterBlockAsm64K(dst, src)
	}
	if len(src) >= limit10B {
		return encodeSnappyBetterBlockAsm12B(dst, src)
	}
	if len(src) >= limit8B {
		return encodeSnappyBetterBlockAsm10B(dst, src)
	}
	if len(src) < minNonLiteralBlockSize {
		return 0
	}
	return encodeSnappyBetterBlockAsm8B(dst, src)
}
