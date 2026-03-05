

package flate

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"strings"
)





func (f *decompressor) huffmanBytesBuffer() {
	const (
		stateInit = iota 
		stateDict
	)
	fr := f.r.(*bytes.Buffer)

	
	
	
	fnb, fb, dict := f.nb, f.b, &f.dict

	switch f.stepState {
	case stateInit:
		goto readLiteral
	case stateDict:
		goto copyHistory
	}

readLiteral:
	
	{
		var v int
		{
			
			
			
			
			
			n := uint(f.hl.maxRead)
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hl.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hl.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hl.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					v = int(chunk >> huffmanValueShift)
					break
				}
			}
		}

		var length int
		switch {
		case v < 256:
			dict.writeByte(byte(v))
			if dict.availWrite() == 0 {
				f.toRead = dict.readFlush()
				f.step = (*decompressor).huffmanBytesBuffer
				f.stepState = stateInit
				f.b, f.nb = fb, fnb
				return
			}
			goto readLiteral
		case v == 256:
			f.b, f.nb = fb, fnb
			f.finishBlock()
			return
		
		case v < 265:
			length = v - (257 - 3)
		case v < maxNumLit:
			val := decCodeToLen[(v - 257)]
			length = int(val.length) + 3
			n := uint(val.extra)
			for fnb < n {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits n>0:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			length += int(fb & bitMask32[n])
			fb >>= n & regSizeMaskUint32
			fnb -= n
		default:
			if debugDecode {
				fmt.Println(v, ">= maxNumLit")
			}
			f.err = CorruptInputError(f.roffset)
			f.b, f.nb = fb, fnb
			return
		}

		var dist uint32
		if f.hd == nil {
			for fnb < 5 {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<5:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			dist = uint32(bits.Reverse8(uint8(fb & 0x1F << 3)))
			fb >>= 5
			fnb -= 5
		} else {
			
			
			
			
			n := uint(f.hd.maxRead)
			
			
			
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hd.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hd.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hd.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					dist = uint32(chunk >> huffmanValueShift)
					break
				}
			}
		}

		switch {
		case dist < 4:
			dist++
		case dist < maxNumDist:
			nb := uint(dist-2) >> 1
			
			extra := (dist & 1) << (nb & regSizeMaskUint32)
			for fnb < nb {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<nb:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			extra |= fb & bitMask32[nb]
			fb >>= nb & regSizeMaskUint32
			fnb -= nb
			dist = 1<<((nb+1)&regSizeMaskUint32) + 1 + extra
			
		default:
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist too big:", dist, maxNumDist)
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		
		if dist > uint32(dict.histSize()) {
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist > dict.histSize():", dist, dict.histSize())
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		f.copyLen, f.copyDist = length, int(dist)
		goto copyHistory
	}

copyHistory:
	
	{
		cnt := dict.tryWriteCopy(f.copyDist, f.copyLen)
		if cnt == 0 {
			cnt = dict.writeCopy(f.copyDist, f.copyLen)
		}
		f.copyLen -= cnt

		if dict.availWrite() == 0 || f.copyLen > 0 {
			f.toRead = dict.readFlush()
			f.step = (*decompressor).huffmanBytesBuffer 
			f.stepState = stateDict
			f.b, f.nb = fb, fnb
			return
		}
		goto readLiteral
	}
	
}





func (f *decompressor) huffmanBytesReader() {
	const (
		stateInit = iota 
		stateDict
	)
	fr := f.r.(*bytes.Reader)

	
	
	
	fnb, fb, dict := f.nb, f.b, &f.dict

	switch f.stepState {
	case stateInit:
		goto readLiteral
	case stateDict:
		goto copyHistory
	}

readLiteral:
	
	{
		var v int
		{
			
			
			
			
			
			n := uint(f.hl.maxRead)
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hl.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hl.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hl.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					v = int(chunk >> huffmanValueShift)
					break
				}
			}
		}

		var length int
		switch {
		case v < 256:
			dict.writeByte(byte(v))
			if dict.availWrite() == 0 {
				f.toRead = dict.readFlush()
				f.step = (*decompressor).huffmanBytesReader
				f.stepState = stateInit
				f.b, f.nb = fb, fnb
				return
			}
			goto readLiteral
		case v == 256:
			f.b, f.nb = fb, fnb
			f.finishBlock()
			return
		
		case v < 265:
			length = v - (257 - 3)
		case v < maxNumLit:
			val := decCodeToLen[(v - 257)]
			length = int(val.length) + 3
			n := uint(val.extra)
			for fnb < n {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits n>0:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			length += int(fb & bitMask32[n])
			fb >>= n & regSizeMaskUint32
			fnb -= n
		default:
			if debugDecode {
				fmt.Println(v, ">= maxNumLit")
			}
			f.err = CorruptInputError(f.roffset)
			f.b, f.nb = fb, fnb
			return
		}

		var dist uint32
		if f.hd == nil {
			for fnb < 5 {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<5:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			dist = uint32(bits.Reverse8(uint8(fb & 0x1F << 3)))
			fb >>= 5
			fnb -= 5
		} else {
			
			
			
			
			n := uint(f.hd.maxRead)
			
			
			
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hd.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hd.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hd.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					dist = uint32(chunk >> huffmanValueShift)
					break
				}
			}
		}

		switch {
		case dist < 4:
			dist++
		case dist < maxNumDist:
			nb := uint(dist-2) >> 1
			
			extra := (dist & 1) << (nb & regSizeMaskUint32)
			for fnb < nb {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<nb:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			extra |= fb & bitMask32[nb]
			fb >>= nb & regSizeMaskUint32
			fnb -= nb
			dist = 1<<((nb+1)&regSizeMaskUint32) + 1 + extra
			
		default:
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist too big:", dist, maxNumDist)
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		
		if dist > uint32(dict.histSize()) {
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist > dict.histSize():", dist, dict.histSize())
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		f.copyLen, f.copyDist = length, int(dist)
		goto copyHistory
	}

copyHistory:
	
	{
		cnt := dict.tryWriteCopy(f.copyDist, f.copyLen)
		if cnt == 0 {
			cnt = dict.writeCopy(f.copyDist, f.copyLen)
		}
		f.copyLen -= cnt

		if dict.availWrite() == 0 || f.copyLen > 0 {
			f.toRead = dict.readFlush()
			f.step = (*decompressor).huffmanBytesReader 
			f.stepState = stateDict
			f.b, f.nb = fb, fnb
			return
		}
		goto readLiteral
	}
	
}





func (f *decompressor) huffmanBufioReader() {
	const (
		stateInit = iota 
		stateDict
	)
	fr := f.r.(*bufio.Reader)

	
	
	
	fnb, fb, dict := f.nb, f.b, &f.dict

	switch f.stepState {
	case stateInit:
		goto readLiteral
	case stateDict:
		goto copyHistory
	}

readLiteral:
	
	{
		var v int
		{
			
			
			
			
			
			n := uint(f.hl.maxRead)
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hl.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hl.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hl.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					v = int(chunk >> huffmanValueShift)
					break
				}
			}
		}

		var length int
		switch {
		case v < 256:
			dict.writeByte(byte(v))
			if dict.availWrite() == 0 {
				f.toRead = dict.readFlush()
				f.step = (*decompressor).huffmanBufioReader
				f.stepState = stateInit
				f.b, f.nb = fb, fnb
				return
			}
			goto readLiteral
		case v == 256:
			f.b, f.nb = fb, fnb
			f.finishBlock()
			return
		
		case v < 265:
			length = v - (257 - 3)
		case v < maxNumLit:
			val := decCodeToLen[(v - 257)]
			length = int(val.length) + 3
			n := uint(val.extra)
			for fnb < n {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits n>0:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			length += int(fb & bitMask32[n])
			fb >>= n & regSizeMaskUint32
			fnb -= n
		default:
			if debugDecode {
				fmt.Println(v, ">= maxNumLit")
			}
			f.err = CorruptInputError(f.roffset)
			f.b, f.nb = fb, fnb
			return
		}

		var dist uint32
		if f.hd == nil {
			for fnb < 5 {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<5:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			dist = uint32(bits.Reverse8(uint8(fb & 0x1F << 3)))
			fb >>= 5
			fnb -= 5
		} else {
			
			
			
			
			n := uint(f.hd.maxRead)
			
			
			
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hd.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hd.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hd.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					dist = uint32(chunk >> huffmanValueShift)
					break
				}
			}
		}

		switch {
		case dist < 4:
			dist++
		case dist < maxNumDist:
			nb := uint(dist-2) >> 1
			
			extra := (dist & 1) << (nb & regSizeMaskUint32)
			for fnb < nb {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<nb:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			extra |= fb & bitMask32[nb]
			fb >>= nb & regSizeMaskUint32
			fnb -= nb
			dist = 1<<((nb+1)&regSizeMaskUint32) + 1 + extra
			
		default:
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist too big:", dist, maxNumDist)
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		
		if dist > uint32(dict.histSize()) {
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist > dict.histSize():", dist, dict.histSize())
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		f.copyLen, f.copyDist = length, int(dist)
		goto copyHistory
	}

copyHistory:
	
	{
		cnt := dict.tryWriteCopy(f.copyDist, f.copyLen)
		if cnt == 0 {
			cnt = dict.writeCopy(f.copyDist, f.copyLen)
		}
		f.copyLen -= cnt

		if dict.availWrite() == 0 || f.copyLen > 0 {
			f.toRead = dict.readFlush()
			f.step = (*decompressor).huffmanBufioReader 
			f.stepState = stateDict
			f.b, f.nb = fb, fnb
			return
		}
		goto readLiteral
	}
	
}





func (f *decompressor) huffmanStringsReader() {
	const (
		stateInit = iota 
		stateDict
	)
	fr := f.r.(*strings.Reader)

	
	
	
	fnb, fb, dict := f.nb, f.b, &f.dict

	switch f.stepState {
	case stateInit:
		goto readLiteral
	case stateDict:
		goto copyHistory
	}

readLiteral:
	
	{
		var v int
		{
			
			
			
			
			
			n := uint(f.hl.maxRead)
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hl.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hl.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hl.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					v = int(chunk >> huffmanValueShift)
					break
				}
			}
		}

		var length int
		switch {
		case v < 256:
			dict.writeByte(byte(v))
			if dict.availWrite() == 0 {
				f.toRead = dict.readFlush()
				f.step = (*decompressor).huffmanStringsReader
				f.stepState = stateInit
				f.b, f.nb = fb, fnb
				return
			}
			goto readLiteral
		case v == 256:
			f.b, f.nb = fb, fnb
			f.finishBlock()
			return
		
		case v < 265:
			length = v - (257 - 3)
		case v < maxNumLit:
			val := decCodeToLen[(v - 257)]
			length = int(val.length) + 3
			n := uint(val.extra)
			for fnb < n {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits n>0:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			length += int(fb & bitMask32[n])
			fb >>= n & regSizeMaskUint32
			fnb -= n
		default:
			if debugDecode {
				fmt.Println(v, ">= maxNumLit")
			}
			f.err = CorruptInputError(f.roffset)
			f.b, f.nb = fb, fnb
			return
		}

		var dist uint32
		if f.hd == nil {
			for fnb < 5 {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<5:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			dist = uint32(bits.Reverse8(uint8(fb & 0x1F << 3)))
			fb >>= 5
			fnb -= 5
		} else {
			
			
			
			
			n := uint(f.hd.maxRead)
			
			
			
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hd.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hd.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hd.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					dist = uint32(chunk >> huffmanValueShift)
					break
				}
			}
		}

		switch {
		case dist < 4:
			dist++
		case dist < maxNumDist:
			nb := uint(dist-2) >> 1
			
			extra := (dist & 1) << (nb & regSizeMaskUint32)
			for fnb < nb {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<nb:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			extra |= fb & bitMask32[nb]
			fb >>= nb & regSizeMaskUint32
			fnb -= nb
			dist = 1<<((nb+1)&regSizeMaskUint32) + 1 + extra
			
		default:
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist too big:", dist, maxNumDist)
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		
		if dist > uint32(dict.histSize()) {
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist > dict.histSize():", dist, dict.histSize())
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		f.copyLen, f.copyDist = length, int(dist)
		goto copyHistory
	}

copyHistory:
	
	{
		cnt := dict.tryWriteCopy(f.copyDist, f.copyLen)
		if cnt == 0 {
			cnt = dict.writeCopy(f.copyDist, f.copyLen)
		}
		f.copyLen -= cnt

		if dict.availWrite() == 0 || f.copyLen > 0 {
			f.toRead = dict.readFlush()
			f.step = (*decompressor).huffmanStringsReader 
			f.stepState = stateDict
			f.b, f.nb = fb, fnb
			return
		}
		goto readLiteral
	}
	
}





func (f *decompressor) huffmanGenericReader() {
	const (
		stateInit = iota 
		stateDict
	)
	fr := f.r.(Reader)

	
	
	
	fnb, fb, dict := f.nb, f.b, &f.dict

	switch f.stepState {
	case stateInit:
		goto readLiteral
	case stateDict:
		goto copyHistory
	}

readLiteral:
	
	{
		var v int
		{
			
			
			
			
			
			n := uint(f.hl.maxRead)
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hl.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hl.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hl.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					v = int(chunk >> huffmanValueShift)
					break
				}
			}
		}

		var length int
		switch {
		case v < 256:
			dict.writeByte(byte(v))
			if dict.availWrite() == 0 {
				f.toRead = dict.readFlush()
				f.step = (*decompressor).huffmanGenericReader
				f.stepState = stateInit
				f.b, f.nb = fb, fnb
				return
			}
			goto readLiteral
		case v == 256:
			f.b, f.nb = fb, fnb
			f.finishBlock()
			return
		
		case v < 265:
			length = v - (257 - 3)
		case v < maxNumLit:
			val := decCodeToLen[(v - 257)]
			length = int(val.length) + 3
			n := uint(val.extra)
			for fnb < n {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits n>0:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			length += int(fb & bitMask32[n])
			fb >>= n & regSizeMaskUint32
			fnb -= n
		default:
			if debugDecode {
				fmt.Println(v, ">= maxNumLit")
			}
			f.err = CorruptInputError(f.roffset)
			f.b, f.nb = fb, fnb
			return
		}

		var dist uint32
		if f.hd == nil {
			for fnb < 5 {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<5:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			dist = uint32(bits.Reverse8(uint8(fb & 0x1F << 3)))
			fb >>= 5
			fnb -= 5
		} else {
			
			
			
			
			n := uint(f.hd.maxRead)
			
			
			
			for {
				for fnb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b, f.nb = fb, fnb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					fb |= uint32(c) << (fnb & regSizeMaskUint32)
					fnb += 8
				}
				chunk := f.hd.chunks[fb&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hd.links[chunk>>huffmanValueShift][(fb>>huffmanChunkBits)&f.hd.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= fnb {
					if n == 0 {
						f.b, f.nb = fb, fnb
						if debugDecode {
							fmt.Println("huffsym: n==0")
						}
						f.err = CorruptInputError(f.roffset)
						return
					}
					fb = fb >> (n & regSizeMaskUint32)
					fnb = fnb - n
					dist = uint32(chunk >> huffmanValueShift)
					break
				}
			}
		}

		switch {
		case dist < 4:
			dist++
		case dist < maxNumDist:
			nb := uint(dist-2) >> 1
			
			extra := (dist & 1) << (nb & regSizeMaskUint32)
			for fnb < nb {
				c, err := fr.ReadByte()
				if err != nil {
					f.b, f.nb = fb, fnb
					if debugDecode {
						fmt.Println("morebits f.nb<nb:", err)
					}
					f.err = err
					return
				}
				f.roffset++
				fb |= uint32(c) << (fnb & regSizeMaskUint32)
				fnb += 8
			}
			extra |= fb & bitMask32[nb]
			fb >>= nb & regSizeMaskUint32
			fnb -= nb
			dist = 1<<((nb+1)&regSizeMaskUint32) + 1 + extra
			
		default:
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist too big:", dist, maxNumDist)
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		
		if dist > uint32(dict.histSize()) {
			f.b, f.nb = fb, fnb
			if debugDecode {
				fmt.Println("dist > dict.histSize():", dist, dict.histSize())
			}
			f.err = CorruptInputError(f.roffset)
			return
		}

		f.copyLen, f.copyDist = length, int(dist)
		goto copyHistory
	}

copyHistory:
	
	{
		cnt := dict.tryWriteCopy(f.copyDist, f.copyLen)
		if cnt == 0 {
			cnt = dict.writeCopy(f.copyDist, f.copyLen)
		}
		f.copyLen -= cnt

		if dict.availWrite() == 0 || f.copyLen > 0 {
			f.toRead = dict.readFlush()
			f.step = (*decompressor).huffmanGenericReader 
			f.stepState = stateDict
			f.b, f.nb = fb, fnb
			return
		}
		goto readLiteral
	}
	
}

func (f *decompressor) huffmanBlockDecoder() func() {
	switch f.r.(type) {
	case *bytes.Buffer:
		return f.huffmanBytesBuffer
	case *bytes.Reader:
		return f.huffmanBytesReader
	case *bufio.Reader:
		return f.huffmanBufioReader
	case *strings.Reader:
		return f.huffmanStringsReader
	case Reader:
		return f.huffmanGenericReader
	default:
		return f.huffmanGenericReader
	}
}
