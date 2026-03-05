package magic

import (
	"bytes"
	"encoding/binary"
)

var (
	
	SevenZ = prefix([]byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C})
	
	Gzip = prefix([]byte{0x1f, 0x8b})
	
	Fits = prefix([]byte{
		0x53, 0x49, 0x4D, 0x50, 0x4C, 0x45, 0x20, 0x20, 0x3D, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54,
	})
	
	Xar = prefix([]byte{0x78, 0x61, 0x72, 0x21})
	
	Bz2 = prefix([]byte{0x42, 0x5A, 0x68})
	
	Ar = prefix([]byte{0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E})
	
	Deb = offset([]byte{
		0x64, 0x65, 0x62, 0x69, 0x61, 0x6E, 0x2D,
		0x62, 0x69, 0x6E, 0x61, 0x72, 0x79,
	}, 8)
	
	Warc = prefix([]byte("WARC/1.0"), []byte("WARC/1.1"))
	
	Cab = prefix([]byte("MSCF\x00\x00\x00\x00"))
	
	Xz = prefix([]byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00})
	
	Lzip = prefix([]byte{0x4c, 0x5a, 0x49, 0x50})
	
	RPM = prefix([]byte{0xed, 0xab, 0xee, 0xdb}, []byte("drpm"))
	
	Cpio = prefix([]byte("070707"), []byte("070701"), []byte("070702"))
	
	RAR = prefix([]byte("Rar!\x1A\x07\x00"), []byte("Rar!\x1A\x07\x01\x00"))
)


func InstallShieldCab(raw []byte, _ uint32) bool {
	return len(raw) > 7 &&
		bytes.Equal(raw[0:4], []byte("ISc(")) &&
		raw[6] == 0 &&
		(raw[7] == 1 || raw[7] == 2 || raw[7] == 4)
}



func Zstd(raw []byte, limit uint32) bool {
	if len(raw) < 4 {
		return false
	}
	sig := binary.LittleEndian.Uint32(raw)
	
	return (sig >= 0xFD2FB522 && sig <= 0xFD2FB528) ||
		(sig >= 0x184D2A50 && sig <= 0x184D2A5F)
}


func CRX(raw []byte, limit uint32) bool {
	const minHeaderLen = 16
	if len(raw) < minHeaderLen || !bytes.HasPrefix(raw, []byte("Cr24")) {
		return false
	}
	pubkeyLen := binary.LittleEndian.Uint32(raw[8:12])
	sigLen := binary.LittleEndian.Uint32(raw[12:16])
	zipOffset := minHeaderLen + pubkeyLen + sigLen
	if uint32(len(raw)) < zipOffset {
		return false
	}
	return Zip(raw[zipOffset:], limit)
}




func Tar(raw []byte, _ uint32) bool {
	const sizeRecord = 512

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	if len(raw) < sizeRecord {
		return false
	}
	raw = raw[:sizeRecord]

	
	
	if bytes.Contains(raw[:100], []byte("/gpkg-1\x00")) {
		return false
	}

	
	recsum := tarParseOctal(raw[148:156])
	if recsum == -1 {
		return false
	}
	sum1, sum2 := tarChksum(raw)
	return recsum == sum1 || recsum == sum2
}


func tarParseOctal(b []byte) int64 {
	
	
	
	b = bytes.Trim(b, " \x00")

	if len(b) == 0 {
		return -1
	}
	ret := int64(0)
	for _, b := range b {
		if b == 0 {
			break
		}
		if !(b >= '0' && b <= '7') {
			return -1
		}
		ret = (ret << 3) | int64(b-'0')
	}
	return ret
}







func tarChksum(b []byte) (unsigned, signed int64) {
	for i, c := range b {
		if 148 <= i && i < 156 {
			c = ' ' 
		}
		unsigned += int64(c)
		signed += int64(int8(c))
	}
	return unsigned, signed
}
