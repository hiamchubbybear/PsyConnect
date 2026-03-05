//go:build (!amd64 && !arm && !arm64) || appengine || !gc || noasm
// +build !amd64,!arm,!arm64 appengine !gc noasm

package lz4block

import (
	"encoding/binary"
)

func decodeBlock(dst, src, dict []byte) (ret int) {
	
	dst = dst[:len(dst):len(dst)]
	src = src[:len(src):len(src)]

	const hasError = -2

	if len(src) == 0 {
		return hasError
	}

	defer func() {
		if recover() != nil {
			ret = hasError
		}
	}()

	var si, di uint
	for si < uint(len(src)) {
		
		b := uint(src[si])
		si++

		
		if lLen := b >> 4; lLen > 0 {
			switch {
			case lLen < 0xF && si+16 < uint(len(src)):
				
				
				
				
				copy(dst[di:], src[si:si+16])
				si += lLen
				di += lLen
				if mLen := b & 0xF; mLen < 0xF {
					
					
					
					mLen += 4
					if offset := u16(src[si:]); mLen <= offset && offset < di {
						i := di - offset
						end := i + 18
						copy(dst[di:], dst[i:end])
						si += 2
						di += mLen
						continue
					}
				}
			case lLen == 0xF:
				for {
					x := uint(src[si])
					if lLen += x; int(lLen) < 0 {
						return hasError
					}
					si++
					if x != 0xFF {
						break
					}
				}
				fallthrough
			default:
				copy(dst[di:di+lLen], src[si:si+lLen])
				si += lLen
				di += lLen
			}
		}

		mLen := b & 0xF
		if si == uint(len(src)) && mLen == 0 {
			break
		} else if si >= uint(len(src)) {
			return hasError
		}

		offset := u16(src[si:])
		if offset == 0 {
			return hasError
		}
		si += 2

		
		mLen += minMatch
		if mLen == minMatch+0xF {
			for {
				x := uint(src[si])
				if mLen += x; int(mLen) < 0 {
					return hasError
				}
				si++
				if x != 0xFF {
					break
				}
			}
		}

		
		if di < offset {
			
			
			fromDict := dict[uint(len(dict))+di-offset:]
			n := uint(copy(dst[di:di+mLen], fromDict))
			di += n
			if mLen -= n; mLen == 0 {
				continue
			}
			
			
			
		}

		expanded := dst[di-offset:]
		if mLen > offset {
			
			bytesToCopy := offset * (mLen / offset)
			for n := offset; n <= bytesToCopy+offset; n *= 2 {
				copy(expanded[n:], expanded[:n])
			}
			di += bytesToCopy
			mLen -= bytesToCopy
		}
		di += uint(copy(dst[di:di+mLen], expanded[:mLen]))
	}

	return int(di)
}

func u16(p []byte) uint { return uint(binary.LittleEndian.Uint16(p)) }
