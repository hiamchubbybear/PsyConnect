package magic

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
)

var (
	
	Lnk = prefix([]byte{0x4C, 0x00, 0x00, 0x00, 0x01, 0x14, 0x02, 0x00})
	
	Wasm = prefix([]byte{0x00, 0x61, 0x73, 0x6D})
	
	Exe = prefix([]byte{0x4D, 0x5A})
	
	Elf = prefix([]byte{0x7F, 0x45, 0x4C, 0x46})
	
	Nes = prefix([]byte{0x4E, 0x45, 0x53, 0x1A})
	
	SWF = prefix([]byte("CWS"), []byte("FWS"), []byte("ZWS"))
	
	Torrent = prefix([]byte("d8:announce"))
	
	Par1 = prefix([]byte{0x50, 0x41, 0x52, 0x31})
	
	CBOR = prefix([]byte{0xD9, 0xD9, 0xF7})
)



func classOrMachOFat(in []byte) bool {
	
	
	if len(in) < 8 {
		return false
	}

	return binary.BigEndian.Uint32(in) == macho.MagicFat
}


func Class(raw []byte, limit uint32) bool {
	return classOrMachOFat(raw) && raw[7] > 30
}


func MachO(raw []byte, limit uint32) bool {
	if classOrMachOFat(raw) && raw[7] < 0x14 {
		return true
	}

	if len(raw) < 4 {
		return false
	}

	be := binary.BigEndian.Uint32(raw)
	le := binary.LittleEndian.Uint32(raw)

	return be == macho.Magic32 ||
		le == macho.Magic32 ||
		be == macho.Magic64 ||
		le == macho.Magic64
}



func Dbf(raw []byte, limit uint32) bool {
	if len(raw) < 68 {
		return false
	}

	
	if !(0 < raw[2] && raw[2] < 13 && 0 < raw[3] && raw[3] < 32) {
		return false
	}

	
	if raw[12] != 0x00 || raw[13] != 0x00 || raw[30] != 0x00 || raw[31] != 0x00 {
		return false
	}
	
	
	
	if raw[28] > 0x01 {
		return false
	}

	
	dbfTypes := []byte{
		0x02, 0x03, 0x04, 0x05, 0x30, 0x31, 0x32, 0x42, 0x62, 0x7B, 0x82,
		0x83, 0x87, 0x8A, 0x8B, 0x8E, 0xB3, 0xCB, 0xE5, 0xF5, 0xF4, 0xFB,
	}
	for _, b := range dbfTypes {
		if raw[0] == b {
			return true
		}
	}

	return false
}


func ElfObj(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x01 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x01))
}


func ElfExe(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x02 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x02))
}


func ElfLib(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x03 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x03))
}


func ElfDump(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x04 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x04))
}


func Dcm(raw []byte, limit uint32) bool {
	return len(raw) > 131 &&
		bytes.Equal(raw[128:132], []byte{0x44, 0x49, 0x43, 0x4D})
}


func Marc(raw []byte, limit uint32) bool {
	
	if len(raw) < 24 {
		return false
	}

	
	if !bytes.Equal(raw[20:24], []byte("4500")) {
		return false
	}

	
	for i := 0; i < 5; i++ {
		if raw[i] < '0' || raw[i] > '9' {
			return false
		}
	}

	
	return bytes.Contains(raw[:min(2048, len(raw))], []byte{0x1E})
}
















var Glb = prefix([]byte("\x67\x6C\x54\x46\x02\x00\x00\x00"),
	[]byte("\x67\x6C\x54\x46\x01\x00\x00\x00"))













func TzIf(raw []byte, limit uint32) bool {
	
	if len(raw) < 44 {
		return false
	}

	if !bytes.HasPrefix(raw, []byte("TZif")) {
		return false
	}

	
	if binary.BigEndian.Uint32(raw[36:40]) == 0 {
		return false
	}

	
	return raw[4] == 0x00 || raw[4] == 0x32 || raw[4] == 0x33
}
