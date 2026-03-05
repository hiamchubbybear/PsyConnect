package magic

import (
	"bytes"
)

var (
	
	
	
	AVIF = ftyp([]byte("avif"), []byte("avis"))
	
	ThreeGP = ftyp(
		[]byte("3gp1"), []byte("3gp2"), []byte("3gp3"), []byte("3gp4"),
		[]byte("3gp5"), []byte("3gp6"), []byte("3gp7"), []byte("3gs7"),
		[]byte("3ge6"), []byte("3ge7"), []byte("3gg6"),
	)
	
	ThreeG2 = ftyp(
		[]byte("3g24"), []byte("3g25"), []byte("3g26"), []byte("3g2a"),
		[]byte("3g2b"), []byte("3g2c"), []byte("KDDI"),
	)
	
	AMp4 = ftyp(
		
		[]byte("F4A "), []byte("F4B "),
		
		[]byte("M4B "), []byte("M4P "),
		
		[]byte("MSNV"),
		
		[]byte("NDAS"),
	)
	
	Mqv = ftyp([]byte("mqt "))
	
	M4a = ftyp([]byte("M4A "))
	
	M4v = ftyp([]byte("M4V "), []byte("M4VH"), []byte("M4VP"))
	
	Heic = ftyp([]byte("heic"), []byte("heix"))
	
	HeicSequence = ftyp([]byte("hevc"), []byte("hevx"))
	
	Heif = ftyp([]byte("mif1"), []byte("heim"), []byte("heis"), []byte("avic"))
	
	HeifSequence = ftyp([]byte("msf1"), []byte("hevm"), []byte("hevs"), []byte("avcs"))
	
	Mj2 = ftyp([]byte("mj2s"), []byte("mjp2"), []byte("MFSM"), []byte("MGSV"))
	
	
	
	Dvb = ftyp(
		[]byte("dby1"), []byte("dsms"), []byte("dts1"), []byte("dts2"),
		[]byte("dts3"), []byte("dxo "), []byte("dmb1"), []byte("dmpf"),
		[]byte("drc1"), []byte("dv1a"), []byte("dv1b"), []byte("dv2a"),
		[]byte("dv2b"), []byte("dv3a"), []byte("dv3b"), []byte("dvr1"),
		[]byte("dvt1"), []byte("emsg"))
	
)





func QuickTime(raw []byte, _ uint32) bool {
	if len(raw) < 12 {
		return false
	}
	
	
	
	
	if bytes.Equal(raw[4:12], []byte("ftypqt  ")) ||
		bytes.Equal(raw[4:12], []byte("ftypmoov")) {
		return raw[0] == 0x00
	}
	basicAtomTypes := [][]byte{
		[]byte("moov\x00"),
		[]byte("mdat\x00"),
		[]byte("free\x00"),
		[]byte("skip\x00"),
		[]byte("pnot\x00"),
	}
	for _, a := range basicAtomTypes {
		if bytes.Equal(raw[4:9], a) {
			return true
		}
	}
	return bytes.Equal(raw[:8], []byte("\x00\x00\x00\x08wide"))
}





func Mp4(raw []byte, _ uint32) bool {
	if len(raw) < 12 {
		return false
	}
	
	
	
	
	if raw[0] != 0 {
		return false
	}
	return bytes.Equal(raw[4:8], []byte("ftyp"))
}
