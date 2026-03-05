package magic

import "bytes"

var (
	
	
	Pdf = prefix(
		
		[]byte("%PDF-"),
		
		[]byte("\012%PDF-"),
		
		[]byte("\xef\xbb\xbf%PDF-"),
	)
	
	Fdf = prefix([]byte("%FDF"))
	
	Mobi = offset([]byte("BOOKMOBI"), 60)
	
	Lit = prefix([]byte("ITOLITLS"))
)


func DjVu(raw []byte, limit uint32) bool {
	if len(raw) < 12 {
		return false
	}
	if !bytes.HasPrefix(raw, []byte{0x41, 0x54, 0x26, 0x54, 0x46, 0x4F, 0x52, 0x4D}) {
		return false
	}
	return bytes.HasPrefix(raw[12:], []byte("DJVM")) ||
		bytes.HasPrefix(raw[12:], []byte("DJVU")) ||
		bytes.HasPrefix(raw[12:], []byte("DJVI")) ||
		bytes.HasPrefix(raw[12:], []byte("THUM"))
}


func P7s(raw []byte, limit uint32) bool {
	
	if bytes.HasPrefix(raw, []byte("-----BEGIN PKCS7")) {
		return true
	}
	
	if len(raw) < 20 {
		return false
	}
	
	startHeader := [][]byte{{0x30, 0x80}, {0x30, 0x81}, {0x30, 0x82}, {0x30, 0x83}, {0x30, 0x84}}
	signedDataMatch := []byte{0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x07}
	
	for i, match := range startHeader {
		
		if bytes.HasPrefix(raw, match) {
			if bytes.HasPrefix(raw[i+2:], signedDataMatch) {
				return true
			}
		}
	}

	return false
}
