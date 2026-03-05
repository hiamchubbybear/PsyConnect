





























package objabi

type RelocType int16

//go:generate stringer -type=RelocType
const (
	R_ADDR RelocType = 1 + iota
	
	
	
	
	
	R_ADDRPOWER
	
	
	R_ADDRARM64
	
	
	R_ADDRMIPS
	
	
	R_ADDROFF
	
	
	
	
	R_WEAKADDROFF
	R_SIZE
	R_CALL
	R_CALLARM
	R_CALLARM64
	R_CALLIND
	R_CALLPOWER
	
	
	R_CALLMIPS
	
	R_CALLRISCV
	R_CONST
	R_PCREL
	
	
	
	
	R_TLS_LE
	
	
	
	
	
	R_TLS_IE
	R_GOTOFF
	R_PLT0
	R_PLT1
	R_PLT2
	R_USEFIELD
	
	
	
	
	R_USETYPE
	
	
	
	
	
	R_METHODOFF
	R_POWER_TOC
	R_GOTPCREL
	
	
	
	R_JMPMIPS

	
	
	R_DWARFSECREF

	
	
	
	R_DWARFFILEREF

	
	
	
	
	
	

	

	
	
	
	R_ARM64_TLS_LE

	
	
	
	R_ARM64_TLS_IE

	
	
	R_ARM64_GOTPCREL

	
	
	R_ARM64_GOT

	
	
	R_ARM64_PCREL

	
	R_ARM64_LDST8

	
	R_ARM64_LDST32

	
	R_ARM64_LDST64

	
	R_ARM64_LDST128

	

	
	
	
	
	R_POWER_TLS_LE

	
	
	
	
	
	R_POWER_TLS_IE

	
	
	
	
	R_POWER_TLS

	
	
	
	
	
	R_ADDRPOWER_DS

	
	
	
	R_ADDRPOWER_GOT

	
	
	
	R_ADDRPOWER_PCREL

	
	
	
	R_ADDRPOWER_TOCREL

	
	
	
	R_ADDRPOWER_TOCREL_DS

	

	
	
	R_RISCV_PCREL_ITYPE

	
	
	R_RISCV_PCREL_STYPE

	
	
	R_PCRELDBL

	
	
	R_ADDRMIPSU
	
	
	R_ADDRMIPSTLS

	
	
	R_ADDRCUOFF

	
	R_WASMIMPORT

	
	
	
	R_XCOFFREF
)






func (r RelocType) IsDirectCall() bool {
	switch r {
	case R_CALL, R_CALLARM, R_CALLARM64, R_CALLMIPS, R_CALLPOWER, R_CALLRISCV:
		return true
	}
	return false
}






func (r RelocType) IsDirectJump() bool {
	switch r {
	case R_JMPMIPS:
		return true
	}
	return false
}



func (r RelocType) IsDirectCallOrJump() bool {
	return r.IsDirectCall() || r.IsDirectJump()
}
