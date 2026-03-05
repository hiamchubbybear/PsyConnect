


package codec



const maxArrayLen = 1<<((32<<(^uint(0)>>63))-1) - 1




func pruneSignExt(v []byte, pos bool) (n int) {
	if len(v) < 2 {
	} else if pos && v[0] == 0 {
		for ; v[n] == 0 && n+1 < len(v) && (v[n+1]&(1<<7) == 0); n++ {
		}
	} else if !pos && v[0] == 0xff {
		for ; v[n] == 0xff && n+1 < len(v) && (v[n+1]&(1<<7) != 0); n++ {
		}
	}
	return
}

func halfFloatToFloatBits(h uint16) (f uint32) {
	
	
	

	s := uint32(h >> 15)
	m := uint32(h & 0x03ff)
	e := int32((h >> 10) & 0x1f)

	if e == 0 {
		if m == 0 { 
			return s << 31
		}
		
		for (m & 0x0400) == 0 {
			m <<= 1
			e -= 1
		}
		e += 1
		m &= ^uint32(0x0400)
	} else if e == 31 {
		if m == 0 { 
			return (s << 31) | 0x7f800000
		}
		return (s << 31) | 0x7f800000 | (m << 13) 
	}
	e = e + (127 - 15)
	m = m << 13
	return (s << 31) | (uint32(e) << 23) | m
}

func floatToHalfFloatBits(i uint32) (h uint16) {
	
	
	
	
	s := (i >> 16) & 0x8000
	e := int32(((i >> 23) & 0xff) - (127 - 15))
	m := i & 0x7fffff

	var h32 uint32

	if e <= 0 {
		if e < -10 { 
			h32 = s 
		} else {
			m = (m | 0x800000) >> uint32(1-e)
			h32 = s | (m >> 13)
		}
	} else if e == 0xff-(127-15) {
		if m == 0 { 
			h32 = s | 0x7c00
		} else { 
			m >>= 13
			var me uint32
			if m == 0 {
				me = 1
			}
			h32 = s | 0x7c00 | m | me
		}
	} else {
		if e > 30 { 
			h32 = s | 0x7c00
		} else {
			h32 = s | (uint32(e) << 10) | (m >> 13)
		}
	}
	h = uint16(h32)
	return
}





func growCap(oldCap, unit, num uint) (newCap uint) {
	
	
	
	

	
	
	
	
	
	
	
	
	
	

	
	maxCap := num + (oldCap * 3 / 2)
	if unit == 0 || maxCap > maxArrayLen || maxCap < oldCap { 
		return maxArrayLen
	}

	var t1 uint = 1024 
	if unit <= 4 {
		t1 = 8 * 1024
	} else if unit <= 16 {
		t1 = 2 * 1024
	}

	newCap = 2 + num
	if oldCap > 0 {
		if oldCap <= t1 { 
			newCap = num + (oldCap * 2)
		} else { 
			newCap = maxCap
		}
	}

	
	t1 = newCap * unit
	if t2 := t1 % 64; t2 != 0 {
		t1 += 64 - t2
		newCap = t1 / unit
	}

	return
}
