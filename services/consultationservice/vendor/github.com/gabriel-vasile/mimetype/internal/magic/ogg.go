package magic

import (
	"bytes"
)




func Ogg(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte("\x4F\x67\x67\x53\x00"))
}


func OggAudio(raw []byte, limit uint32) bool {
	return len(raw) >= 37 && (bytes.HasPrefix(raw[28:], []byte("\x7fFLAC")) ||
		bytes.HasPrefix(raw[28:], []byte("\x01vorbis")) ||
		bytes.HasPrefix(raw[28:], []byte("OpusHead")) ||
		bytes.HasPrefix(raw[28:], []byte("Speex\x20\x20\x20")))
}


func OggVideo(raw []byte, limit uint32) bool {
	return len(raw) >= 37 && (bytes.HasPrefix(raw[28:], []byte("\x80theora")) ||
		bytes.HasPrefix(raw[28:], []byte("fishead\x00")) ||
		bytes.HasPrefix(raw[28:], []byte("\x01video\x00\x00\x00"))) 
}
