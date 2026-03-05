package magic

import (
	"bytes"
	"encoding/binary"
)

var (
	
	Flac = prefix([]byte("\x66\x4C\x61\x43\x00\x00\x00\x22"))
	
	Midi = prefix([]byte("\x4D\x54\x68\x64"))
	
	Ape = prefix([]byte("\x4D\x41\x43\x20\x96\x0F\x00\x00\x34\x00\x00\x00\x18\x00\x00\x00\x90\xE3"))
	
	MusePack = prefix([]byte("MPCK"))
	
	Au = prefix([]byte("\x2E\x73\x6E\x64"))
	
	Amr = prefix([]byte("\x23\x21\x41\x4D\x52"))
	
	Voc = prefix([]byte("Creative Voice File"))
	
	M3u = prefix([]byte("#EXTM3U"))
	
	AAC = prefix([]byte{0xFF, 0xF1}, []byte{0xFF, 0xF9})
)


func Mp3(raw []byte, limit uint32) bool {
	if len(raw) < 3 {
		return false
	}

	if bytes.HasPrefix(raw, []byte("ID3")) {
		
		
		return true
	}

	
	switch binary.BigEndian.Uint16(raw[:2]) & 0xFFFE {
	case 0xFFFA:
		
		return true
	case 0xFFF2:
		
		return true
	case 0xFFE2:
		
		return true
	}

	return false
}


func Wav(raw []byte, limit uint32) bool {
	return len(raw) > 12 &&
		bytes.Equal(raw[:4], []byte("RIFF")) &&
		bytes.Equal(raw[8:12], []byte{0x57, 0x41, 0x56, 0x45})
}


func Aiff(raw []byte, limit uint32) bool {
	return len(raw) > 12 &&
		bytes.Equal(raw[:4], []byte{0x46, 0x4F, 0x52, 0x4D}) &&
		bytes.Equal(raw[8:12], []byte{0x41, 0x49, 0x46, 0x46})
}


func Qcp(raw []byte, limit uint32) bool {
	return len(raw) > 12 &&
		bytes.Equal(raw[:4], []byte("RIFF")) &&
		bytes.Equal(raw[8:12], []byte("QLCM"))
}
