package magic

import (
	"bytes"
	"encoding/binary"
)

var (
	
	Odt = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.text"), 30)
	
	Ott = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.text-template"), 30)
	
	Ods = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.spreadsheet"), 30)
	
	Ots = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.spreadsheet-template"), 30)
	
	Odp = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.presentation"), 30)
	
	Otp = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.presentation-template"), 30)
	
	Odg = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.graphics"), 30)
	
	Otg = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.graphics-template"), 30)
	
	Odf = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.formula"), 30)
	
	Odc = offset([]byte("mimetypeapplication/vnd.oasis.opendocument.chart"), 30)
	
	Epub = offset([]byte("mimetypeapplication/epub+zip"), 30)
	
	Sxc = offset([]byte("mimetypeapplication/vnd.sun.xml.calc"), 30)
)


func Zip(raw []byte, limit uint32) bool {
	return len(raw) > 3 &&
		raw[0] == 0x50 && raw[1] == 0x4B &&
		(raw[2] == 0x3 || raw[2] == 0x5 || raw[2] == 0x7) &&
		(raw[3] == 0x4 || raw[3] == 0x6 || raw[3] == 0x8)
}


func Jar(raw []byte, limit uint32) bool {
	return zipContains(raw, []byte("META-INF/MANIFEST.MF"), false)
}

func zipContains(raw, sig []byte, msoCheck bool) bool {
	b := readBuf(raw)
	pk := []byte("PK\003\004")
	if len(b) < 0x1E {
		return false
	}

	if !b.advance(0x1E) {
		return false
	}
	if bytes.HasPrefix(b, sig) {
		return true
	}

	if msoCheck {
		skipFiles := [][]byte{
			[]byte("[Content_Types].xml"),
			[]byte("_rels/.rels"),
			[]byte("docProps"),
			[]byte("customXml"),
			[]byte("[trash]"),
		}

		hasSkipFile := false
		for _, sf := range skipFiles {
			if bytes.HasPrefix(b, sf) {
				hasSkipFile = true
				break
			}
		}
		if !hasSkipFile {
			return false
		}
	}

	searchOffset := binary.LittleEndian.Uint32(raw[18:]) + 49
	if !b.advance(int(searchOffset)) {
		return false
	}

	nextHeader := bytes.Index(raw[searchOffset:], pk)
	if !b.advance(nextHeader) {
		return false
	}
	if bytes.HasPrefix(b, sig) {
		return true
	}

	for i := 0; i < 4; i++ {
		if !b.advance(0x1A) {
			return false
		}
		nextHeader = bytes.Index(b, pk)
		if nextHeader == -1 {
			return false
		}
		if !b.advance(nextHeader + 0x1E) {
			return false
		}
		if bytes.HasPrefix(b, sig) {
			return true
		}
	}
	return false
}



func APK(raw []byte, _ uint32) bool {
	apkSignatures := [][]byte{
		[]byte("AndroidManifest.xml"),
		[]byte("META-INF/com/android/build/gradle/app-metadata.properties"),
		[]byte("classes.dex"),
		[]byte("resources.arsc"),
		[]byte("res/drawable"),
	}
	for _, sig := range apkSignatures {
		if zipContains(raw, sig, false) {
			return true
		}
	}

	return false
}
