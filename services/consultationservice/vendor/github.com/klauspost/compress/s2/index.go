



package s2

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

const (
	S2IndexHeader   = "s2idx\x00"
	S2IndexTrailer  = "\x00xdi2s"
	maxIndexEntries = 1 << 16
)


type Index struct {
	TotalUncompressed int64 
	TotalCompressed   int64 
	info              []struct {
		compressedOffset   int64
		uncompressedOffset int64
	}
	estBlockUncomp int64
}

func (i *Index) reset(maxBlock int) {
	i.estBlockUncomp = int64(maxBlock)
	i.TotalCompressed = -1
	i.TotalUncompressed = -1
	if len(i.info) > 0 {
		i.info = i.info[:0]
	}
}


func (i *Index) allocInfos(n int) {
	if n > maxIndexEntries {
		panic("n > maxIndexEntries")
	}
	i.info = make([]struct {
		compressedOffset   int64
		uncompressedOffset int64
	}, 0, n)
}



func (i *Index) add(compressedOffset, uncompressedOffset int64) error {
	if i == nil {
		return nil
	}
	lastIdx := len(i.info) - 1
	if lastIdx >= 0 {
		latest := i.info[lastIdx]
		if latest.uncompressedOffset == uncompressedOffset {
			
			
			latest.compressedOffset = compressedOffset
			i.info[lastIdx] = latest
			return nil
		}
		if latest.uncompressedOffset > uncompressedOffset {
			return fmt.Errorf("internal error: Earlier uncompressed received (%d > %d)", latest.uncompressedOffset, uncompressedOffset)
		}
		if latest.compressedOffset > compressedOffset {
			return fmt.Errorf("internal error: Earlier compressed received (%d > %d)", latest.uncompressedOffset, uncompressedOffset)
		}
	}
	i.info = append(i.info, struct {
		compressedOffset   int64
		uncompressedOffset int64
	}{compressedOffset: compressedOffset, uncompressedOffset: uncompressedOffset})
	return nil
}









func (i *Index) Find(offset int64) (compressedOff, uncompressedOff int64, err error) {
	if i.TotalUncompressed < 0 {
		return 0, 0, ErrCorrupt
	}
	if offset < 0 {
		offset = i.TotalUncompressed + offset
		if offset < 0 {
			return 0, 0, io.ErrUnexpectedEOF
		}
	}
	if offset > i.TotalUncompressed {
		return 0, 0, io.ErrUnexpectedEOF
	}
	if len(i.info) > 200 {
		n := sort.Search(len(i.info), func(n int) bool {
			return i.info[n].uncompressedOffset > offset
		})
		if n == 0 {
			n = 1
		}
		return i.info[n-1].compressedOffset, i.info[n-1].uncompressedOffset, nil
	}
	for _, info := range i.info {
		if info.uncompressedOffset > offset {
			break
		}
		compressedOff = info.compressedOffset
		uncompressedOff = info.uncompressedOffset
	}
	return compressedOff, uncompressedOff, nil
}


func (i *Index) reduce() {
	if len(i.info) < maxIndexEntries && i.estBlockUncomp >= 1<<20 {
		return
	}

	
	removeN := (len(i.info) + 1) / maxIndexEntries
	src := i.info
	j := 0

	
	for i.estBlockUncomp*(int64(removeN)+1) < 1<<20 && len(i.info)/(removeN+1) > 1000 {
		removeN++
	}
	for idx := 0; idx < len(src); idx++ {
		i.info[j] = src[idx]
		j++
		idx += removeN
	}
	i.info = i.info[:j]
	
	i.estBlockUncomp += i.estBlockUncomp * int64(removeN)
}

func (i *Index) appendTo(b []byte, uncompTotal, compTotal int64) []byte {
	i.reduce()
	var tmp [binary.MaxVarintLen64]byte

	initSize := len(b)
	
	b = append(b, ChunkTypeIndex, 0, 0, 0)
	b = append(b, []byte(S2IndexHeader)...)
	
	n := binary.PutVarint(tmp[:], uncompTotal)
	b = append(b, tmp[:n]...)
	
	n = binary.PutVarint(tmp[:], compTotal)
	b = append(b, tmp[:n]...)
	
	n = binary.PutVarint(tmp[:], i.estBlockUncomp)
	b = append(b, tmp[:n]...)
	
	n = binary.PutVarint(tmp[:], int64(len(i.info)))
	b = append(b, tmp[:n]...)

	
	var hasUncompressed byte
	for idx, info := range i.info {
		if idx == 0 {
			if info.uncompressedOffset != 0 {
				hasUncompressed = 1
				break
			}
			continue
		}
		if info.uncompressedOffset != i.info[idx-1].uncompressedOffset+i.estBlockUncomp {
			hasUncompressed = 1
			break
		}
	}
	b = append(b, hasUncompressed)

	
	if hasUncompressed == 1 {
		for idx, info := range i.info {
			uOff := info.uncompressedOffset
			if idx > 0 {
				prev := i.info[idx-1]
				uOff -= prev.uncompressedOffset + (i.estBlockUncomp)
			}
			n = binary.PutVarint(tmp[:], uOff)
			b = append(b, tmp[:n]...)
		}
	}

	
	cPredict := i.estBlockUncomp / 2

	for idx, info := range i.info {
		cOff := info.compressedOffset
		if idx > 0 {
			prev := i.info[idx-1]
			cOff -= prev.compressedOffset + cPredict
			
			cPredict += cOff / 2
		}
		n = binary.PutVarint(tmp[:], cOff)
		b = append(b, tmp[:n]...)
	}

	
	
	binary.LittleEndian.PutUint32(tmp[:], uint32(len(b)-initSize+4+len(S2IndexTrailer)))
	b = append(b, tmp[:4]...)
	
	b = append(b, []byte(S2IndexTrailer)...)

	
	chunkLen := len(b) - initSize - skippableFrameHeader
	b[initSize+1] = uint8(chunkLen >> 0)
	b[initSize+2] = uint8(chunkLen >> 8)
	b[initSize+3] = uint8(chunkLen >> 16)
	
	return b
}



func (i *Index) Load(b []byte) ([]byte, error) {
	if len(b) <= 4+len(S2IndexHeader)+len(S2IndexTrailer) {
		return b, io.ErrUnexpectedEOF
	}
	if b[0] != ChunkTypeIndex {
		return b, ErrCorrupt
	}
	chunkLen := int(b[1]) | int(b[2])<<8 | int(b[3])<<16
	b = b[4:]

	
	if len(b) < chunkLen {
		return b, io.ErrUnexpectedEOF
	}
	if !bytes.Equal(b[:len(S2IndexHeader)], []byte(S2IndexHeader)) {
		return b, ErrUnsupported
	}
	b = b[len(S2IndexHeader):]

	
	if v, n := binary.Varint(b); n <= 0 || v < 0 {
		return b, ErrCorrupt
	} else {
		i.TotalUncompressed = v
		b = b[n:]
	}

	
	if v, n := binary.Varint(b); n <= 0 {
		return b, ErrCorrupt
	} else {
		i.TotalCompressed = v
		b = b[n:]
	}

	
	if v, n := binary.Varint(b); n <= 0 {
		return b, ErrCorrupt
	} else {
		if v < 0 {
			return b, ErrCorrupt
		}
		i.estBlockUncomp = v
		b = b[n:]
	}

	var entries int
	if v, n := binary.Varint(b); n <= 0 {
		return b, ErrCorrupt
	} else {
		if v < 0 || v > maxIndexEntries {
			return b, ErrCorrupt
		}
		entries = int(v)
		b = b[n:]
	}
	if cap(i.info) < entries {
		i.allocInfos(entries)
	}
	i.info = i.info[:entries]

	if len(b) < 1 {
		return b, io.ErrUnexpectedEOF
	}
	hasUncompressed := b[0]
	b = b[1:]
	if hasUncompressed&1 != hasUncompressed {
		return b, ErrCorrupt
	}

	
	for idx := range i.info {
		var uOff int64
		if hasUncompressed != 0 {
			
			if v, n := binary.Varint(b); n <= 0 {
				return b, ErrCorrupt
			} else {
				uOff = v
				b = b[n:]
			}
		}

		if idx > 0 {
			prev := i.info[idx-1].uncompressedOffset
			uOff += prev + (i.estBlockUncomp)
			if uOff <= prev {
				return b, ErrCorrupt
			}
		}
		if uOff < 0 {
			return b, ErrCorrupt
		}
		i.info[idx].uncompressedOffset = uOff
	}

	
	cPredict := i.estBlockUncomp / 2

	
	for idx := range i.info {
		var cOff int64
		if v, n := binary.Varint(b); n <= 0 {
			return b, ErrCorrupt
		} else {
			cOff = v
			b = b[n:]
		}

		if idx > 0 {
			
			cPredictNew := cPredict + cOff/2

			prev := i.info[idx-1].compressedOffset
			cOff += prev + cPredict
			if cOff <= prev {
				return b, ErrCorrupt
			}
			cPredict = cPredictNew
		}
		if cOff < 0 {
			return b, ErrCorrupt
		}
		i.info[idx].compressedOffset = cOff
	}
	if len(b) < 4+len(S2IndexTrailer) {
		return b, io.ErrUnexpectedEOF
	}
	
	b = b[4:]

	
	if !bytes.Equal(b[:len(S2IndexTrailer)], []byte(S2IndexTrailer)) {
		return b, ErrCorrupt
	}
	return b[len(S2IndexTrailer):], nil
}






func (i *Index) LoadStream(rs io.ReadSeeker) error {
	
	_, err := rs.Seek(-10, io.SeekEnd)
	if err != nil {
		return err
	}
	var tmp [10]byte
	_, err = io.ReadFull(rs, tmp[:])
	if err != nil {
		return err
	}
	
	if !bytes.Equal(tmp[4:4+len(S2IndexTrailer)], []byte(S2IndexTrailer)) {
		return ErrUnsupported
	}
	sz := binary.LittleEndian.Uint32(tmp[:4])
	if sz > maxChunkSize+skippableFrameHeader {
		return ErrCorrupt
	}
	_, err = rs.Seek(-int64(sz), io.SeekEnd)
	if err != nil {
		return err
	}

	
	buf := make([]byte, sz)
	_, err = io.ReadFull(rs, buf)
	if err != nil {
		return err
	}
	_, err = i.Load(buf)
	return err
}






func IndexStream(r io.Reader) ([]byte, error) {
	var i Index
	var buf [maxChunkSize]byte
	var readHeader bool
	for {
		_, err := io.ReadFull(r, buf[:4])
		if err != nil {
			if err == io.EOF {
				return i.appendTo(nil, i.TotalUncompressed, i.TotalCompressed), nil
			}
			return nil, err
		}
		
		startChunk := i.TotalCompressed
		i.TotalCompressed += 4

		chunkType := buf[0]
		if !readHeader {
			if chunkType != chunkTypeStreamIdentifier {
				return nil, ErrCorrupt
			}
			readHeader = true
		}
		chunkLen := int(buf[1]) | int(buf[2])<<8 | int(buf[3])<<16
		if chunkLen < checksumSize {
			return nil, ErrCorrupt
		}

		i.TotalCompressed += int64(chunkLen)
		_, err = io.ReadFull(r, buf[:chunkLen])
		if err != nil {
			return nil, io.ErrUnexpectedEOF
		}
		
		
		switch chunkType {
		case chunkTypeCompressedData:
			
			
			dLen, err := DecodedLen(buf[checksumSize:])
			if err != nil {
				return nil, err
			}
			if dLen > maxBlockSize {
				return nil, ErrCorrupt
			}
			if i.estBlockUncomp == 0 {
				
				i.estBlockUncomp = int64(dLen)
			}
			err = i.add(startChunk, i.TotalUncompressed)
			if err != nil {
				return nil, err
			}
			i.TotalUncompressed += int64(dLen)
			continue
		case chunkTypeUncompressedData:
			n2 := chunkLen - checksumSize
			if n2 > maxBlockSize {
				return nil, ErrCorrupt
			}
			if i.estBlockUncomp == 0 {
				
				i.estBlockUncomp = int64(n2)
			}
			err = i.add(startChunk, i.TotalUncompressed)
			if err != nil {
				return nil, err
			}
			i.TotalUncompressed += int64(n2)
			continue
		case chunkTypeStreamIdentifier:
			
			if chunkLen != len(magicBody) {
				return nil, ErrCorrupt
			}

			if string(buf[:len(magicBody)]) != magicBody {
				if string(buf[:len(magicBody)]) != magicBodySnappy {
					return nil, ErrCorrupt
				}
			}

			continue
		}

		if chunkType <= 0x7f {
			
			return nil, ErrUnsupported
		}
		if chunkLen > maxChunkSize {
			return nil, ErrUnsupported
		}
		
		
	}
}


func (i *Index) JSON() []byte {
	x := struct {
		TotalUncompressed int64 `json:"total_uncompressed"` 
		TotalCompressed   int64 `json:"total_compressed"`   
		Offsets           []struct {
			CompressedOffset   int64 `json:"compressed"`
			UncompressedOffset int64 `json:"uncompressed"`
		} `json:"offsets"`
		EstBlockUncomp int64 `json:"est_block_uncompressed"`
	}{
		TotalUncompressed: i.TotalUncompressed,
		TotalCompressed:   i.TotalCompressed,
		EstBlockUncomp:    i.estBlockUncomp,
	}
	for _, v := range i.info {
		x.Offsets = append(x.Offsets, struct {
			CompressedOffset   int64 `json:"compressed"`
			UncompressedOffset int64 `json:"uncompressed"`
		}{CompressedOffset: v.compressedOffset, UncompressedOffset: v.uncompressedOffset})
	}
	b, _ := json.MarshalIndent(x, "", "  ")
	return b
}







func RemoveIndexHeaders(b []byte) []byte {
	const save = 4 + len(S2IndexHeader) + len(S2IndexTrailer) + 4
	if len(b) <= save {
		return nil
	}
	if b[0] != ChunkTypeIndex {
		return nil
	}
	chunkLen := int(b[1]) | int(b[2])<<8 | int(b[3])<<16
	b = b[4:]

	
	if len(b) < chunkLen {
		return nil
	}
	b = b[:chunkLen]

	if !bytes.Equal(b[:len(S2IndexHeader)], []byte(S2IndexHeader)) {
		return nil
	}
	b = b[len(S2IndexHeader):]
	if !bytes.HasSuffix(b, []byte(S2IndexTrailer)) {
		return nil
	}
	b = bytes.TrimSuffix(b, []byte(S2IndexTrailer))

	if len(b) < 4 {
		return nil
	}
	return b[:len(b)-4]
}




func RestoreIndexHeaders(in []byte) []byte {
	if len(in) == 0 {
		return in
	}
	b := make([]byte, 0, 4+len(S2IndexHeader)+len(in)+len(S2IndexTrailer)+4)
	b = append(b, ChunkTypeIndex, 0, 0, 0)
	b = append(b, []byte(S2IndexHeader)...)
	b = append(b, in...)

	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(len(b)+4+len(S2IndexTrailer)))
	b = append(b, tmp[:4]...)
	
	b = append(b, []byte(S2IndexTrailer)...)

	chunkLen := len(b) - skippableFrameHeader
	b[1] = uint8(chunkLen >> 0)
	b[2] = uint8(chunkLen >> 8)
	b[3] = uint8(chunkLen >> 16)
	return b
}
