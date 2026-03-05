



package s2

import (
	"bytes"
	"encoding/binary"
	"sync"
)

const (
	
	MinDictSize = 16

	
	MaxDictSize = 65536

	
	MaxDictSrcOffset = 65535
)


type Dict struct {
	dict   []byte
	repeat int 

	fast, better, best sync.Once
	fastTable          *[1 << 14]uint16

	betterTableShort *[1 << 14]uint16
	betterTableLong  *[1 << 17]uint16

	bestTableShort *[1 << 16]uint32
	bestTableLong  *[1 << 19]uint32
}



func NewDict(dict []byte) *Dict {
	if len(dict) == 0 {
		return nil
	}
	var d Dict
	
	r, n := binary.Uvarint(dict)
	if n <= 0 {
		return nil
	}
	dict = dict[n:]
	d.dict = dict
	if cap(d.dict) < len(d.dict)+16 {
		d.dict = append(make([]byte, 0, len(d.dict)+16), d.dict...)
	}
	if len(dict) < MinDictSize || len(dict) > MaxDictSize {
		return nil
	}
	d.repeat = int(r)
	if d.repeat > len(dict) {
		return nil
	}
	return &d
}



func (d *Dict) Bytes() []byte {
	dst := make([]byte, binary.MaxVarintLen16+len(d.dict))
	return append(dst[:binary.PutUvarint(dst, uint64(d.repeat))], d.dict...)
}









func MakeDict(data []byte, searchStart []byte) *Dict {
	if len(data) == 0 {
		return nil
	}
	if len(data) > MaxDictSize {
		data = data[len(data)-MaxDictSize:]
	}
	var d Dict
	dict := data
	d.dict = dict
	if cap(d.dict) < len(d.dict)+16 {
		d.dict = append(make([]byte, 0, len(d.dict)+16), d.dict...)
	}
	if len(dict) < MinDictSize {
		return nil
	}

	
	for s := len(searchStart); s > 4; s-- {
		if idx := bytes.LastIndex(data, searchStart[:s]); idx >= 0 && idx <= len(data)-8 {
			d.repeat = idx
			break
		}
	}

	return &d
}













func (d *Dict) Encode(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if cap(dst) < n {
		dst = make([]byte, n)
	} else {
		dst = dst[:n]
	}

	
	dstP := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:dstP]
	}
	if len(src) < minNonLiteralBlockSize {
		dstP += emitLiteral(dst[dstP:], src)
		return dst[:dstP]
	}
	n := encodeBlockDictGo(dst[dstP:], src, d)
	if n > 0 {
		dstP += n
		return dst[:dstP]
	}
	
	dstP += emitLiteral(dst[dstP:], src)
	return dst[:dstP]
}
















func (d *Dict) EncodeBetter(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if len(dst) < n {
		dst = make([]byte, n)
	}

	
	dstP := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:dstP]
	}
	if len(src) < minNonLiteralBlockSize {
		dstP += emitLiteral(dst[dstP:], src)
		return dst[:dstP]
	}
	n := encodeBlockBetterDict(dst[dstP:], src, d)
	if n > 0 {
		dstP += n
		return dst[:dstP]
	}
	
	dstP += emitLiteral(dst[dstP:], src)
	return dst[:dstP]
}
















func (d *Dict) EncodeBest(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if len(dst) < n {
		dst = make([]byte, n)
	}

	
	dstP := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:dstP]
	}
	if len(src) < minNonLiteralBlockSize {
		dstP += emitLiteral(dst[dstP:], src)
		return dst[:dstP]
	}
	n := encodeBlockBest(dst[dstP:], src, d)
	if n > 0 {
		dstP += n
		return dst[:dstP]
	}
	
	dstP += emitLiteral(dst[dstP:], src)
	return dst[:dstP]
}






func (d *Dict) Decode(dst, src []byte) ([]byte, error) {
	dLen, s, err := decodedLen(src)
	if err != nil {
		return nil, err
	}
	if dLen <= cap(dst) {
		dst = dst[:dLen]
	} else {
		dst = make([]byte, dLen)
	}
	if s2DecodeDict(dst, src[s:], d) != 0 {
		return nil, ErrCorrupt
	}
	return dst, nil
}

func (d *Dict) initFast() {
	d.fast.Do(func() {
		const (
			tableBits    = 14
			maxTableSize = 1 << tableBits
		)

		var table [maxTableSize]uint16
		
		for i := 0; i < len(d.dict)-8-2; i += 3 {
			x0 := load64(d.dict, i)
			h0 := hash6(x0, tableBits)
			h1 := hash6(x0>>8, tableBits)
			h2 := hash6(x0>>16, tableBits)
			table[h0] = uint16(i)
			table[h1] = uint16(i + 1)
			table[h2] = uint16(i + 2)
		}
		d.fastTable = &table
	})
}

func (d *Dict) initBetter() {
	d.better.Do(func() {
		const (
			
			lTableBits    = 17
			maxLTableSize = 1 << lTableBits

			
			sTableBits    = 14
			maxSTableSize = 1 << sTableBits
		)

		var lTable [maxLTableSize]uint16
		var sTable [maxSTableSize]uint16

		
		for i := 0; i < len(d.dict)-8; i++ {
			cv := load64(d.dict, i)
			lTable[hash7(cv, lTableBits)] = uint16(i)
			sTable[hash4(cv, sTableBits)] = uint16(i)
		}
		d.betterTableShort = &sTable
		d.betterTableLong = &lTable
	})
}

func (d *Dict) initBest() {
	d.best.Do(func() {
		const (
			
			lTableBits    = 19
			maxLTableSize = 1 << lTableBits

			
			sTableBits    = 16
			maxSTableSize = 1 << sTableBits
		)

		var lTable [maxLTableSize]uint32
		var sTable [maxSTableSize]uint32

		
		for i := 0; i < len(d.dict)-8; i++ {
			cv := load64(d.dict, i)
			hashL := hash8(cv, lTableBits)
			hashS := hash4(cv, sTableBits)
			candidateL := lTable[hashL]
			candidateS := sTable[hashS]
			lTable[hashL] = uint32(i) | candidateL<<16
			sTable[hashS] = uint32(i) | candidateS<<16
		}
		d.bestTableShort = &sTable
		d.bestTableLong = &lTable
	})
}
