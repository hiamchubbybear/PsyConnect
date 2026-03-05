package magic

import (
	"bytes"
)

var (
	
	Flv = prefix([]byte("\x46\x4C\x56\x01"))
	
	Asf = prefix([]byte{
		0x30, 0x26, 0xB2, 0x75, 0x8E, 0x66, 0xCF, 0x11,
		0xA6, 0xD9, 0x00, 0xAA, 0x00, 0x62, 0xCE, 0x6C,
	})
	
	Rmvb = prefix([]byte{0x2E, 0x52, 0x4D, 0x46})
)


func WebM(raw []byte, limit uint32) bool {
	return isMatroskaFileTypeMatched(raw, "webm")
}


func Mkv(raw []byte, limit uint32) bool {
	return isMatroskaFileTypeMatched(raw, "matroska")
}






func isMatroskaFileTypeMatched(in []byte, flType string) bool {
	if bytes.HasPrefix(in, []byte("\x1A\x45\xDF\xA3")) {
		return isFileTypeNamePresent(in, flType)
	}
	return false
}





func isFileTypeNamePresent(in []byte, flType string) bool {
	ind, maxInd, lenIn := 0, 4096, len(in)
	if lenIn < maxInd { 
		maxInd = lenIn
	}
	ind = bytes.Index(in[:maxInd], []byte("\x42\x82"))
	if ind > 0 && lenIn > ind+2 {
		ind += 2

		
		
		n := vintWidth(int(in[ind]))
		if lenIn > ind+n {
			return bytes.HasPrefix(in[ind+n:], []byte(flType))
		}
	}
	return false
}


func vintWidth(v int) int {
	mask, max, num := 128, 8, 1
	for num < max && v&mask == 0 {
		mask = mask >> 1
		num++
	}
	return num
}


func Mpeg(raw []byte, limit uint32) bool {
	return len(raw) > 3 && bytes.HasPrefix(raw, []byte{0x00, 0x00, 0x01}) &&
		raw[3] >= 0xB0 && raw[3] <= 0xBF
}


func Avi(raw []byte, limit uint32) bool {
	return len(raw) > 16 &&
		bytes.Equal(raw[:4], []byte("RIFF")) &&
		bytes.Equal(raw[8:16], []byte("AVI LIST"))
}
