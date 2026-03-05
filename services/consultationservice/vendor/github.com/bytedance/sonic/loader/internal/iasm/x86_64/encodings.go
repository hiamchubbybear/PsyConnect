















package x86_64

import (
	"encoding/binary"
	"math"
)



func imml(v interface{}) byte {
	return byte(toImmAny(v) & 0x0f)
}

func relv(v interface{}) int64 {
	switch r := v.(type) {
	case *Label:
		return 0
	case RelativeOffset:
		return int64(r)
	default:
		panic("invalid relative offset")
	}
}

func addr(v interface{}) interface{} {
	switch a := v.(*MemoryOperand).Addr; a.Type {
	case Memory:
		return a.Memory
	case Offset:
		return a.Offset
	case Reference:
		return a.Reference
	default:
		panic("invalid memory operand type")
	}
}

func bcode(v interface{}) byte {
	if m, ok := v.(*MemoryOperand); !ok {
		panic("v is not a memory operand")
	} else if m.Broadcast == 0 {
		return 0
	} else {
		return 1
	}
}

func vcode(v interface{}) byte {
	switch r := v.(type) {
	case XMMRegister:
		return byte(r)
	case YMMRegister:
		return byte(r)
	case ZMMRegister:
		return byte(r)
	case MaskedRegister:
		return vcode(r.Reg)
	default:
		panic("v is not a vector register")
	}
}

func kcode(v interface{}) byte {
	switch r := v.(type) {
	case KRegister:
		return byte(r)
	case XMMRegister:
		return 0
	case YMMRegister:
		return 0
	case ZMMRegister:
		return 0
	case RegisterMask:
		return byte(r.K)
	case MaskedRegister:
		return byte(r.Mask.K)
	case *MemoryOperand:
		return toKcodeMem(r)
	default:
		panic("v is not a maskable operand")
	}
}

func zcode(v interface{}) byte {
	switch r := v.(type) {
	case KRegister:
		return 0
	case XMMRegister:
		return 0
	case YMMRegister:
		return 0
	case ZMMRegister:
		return 0
	case RegisterMask:
		return toZcodeRegM(r)
	case MaskedRegister:
		return toZcodeRegM(r.Mask)
	case *MemoryOperand:
		return toZcodeMem(r)
	default:
		panic("v is not a maskable operand")
	}
}

func lcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r & 0x07)
	case Register16:
		return byte(r & 0x07)
	case Register32:
		return byte(r & 0x07)
	case Register64:
		return byte(r & 0x07)
	case KRegister:
		return byte(r & 0x07)
	case MMRegister:
		return byte(r & 0x07)
	case XMMRegister:
		return byte(r & 0x07)
	case YMMRegister:
		return byte(r & 0x07)
	case ZMMRegister:
		return byte(r & 0x07)
	case MaskedRegister:
		return lcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func hcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r>>3) & 1
	case Register16:
		return byte(r>>3) & 1
	case Register32:
		return byte(r>>3) & 1
	case Register64:
		return byte(r>>3) & 1
	case KRegister:
		return byte(r>>3) & 1
	case MMRegister:
		return byte(r>>3) & 1
	case XMMRegister:
		return byte(r>>3) & 1
	case YMMRegister:
		return byte(r>>3) & 1
	case ZMMRegister:
		return byte(r>>3) & 1
	case MaskedRegister:
		return hcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func ecode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r>>4) & 1
	case Register16:
		return byte(r>>4) & 1
	case Register32:
		return byte(r>>4) & 1
	case Register64:
		return byte(r>>4) & 1
	case KRegister:
		return byte(r>>4) & 1
	case MMRegister:
		return byte(r>>4) & 1
	case XMMRegister:
		return byte(r>>4) & 1
	case YMMRegister:
		return byte(r>>4) & 1
	case ZMMRegister:
		return byte(r>>4) & 1
	case MaskedRegister:
		return ecode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func hlcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return toHLcodeReg8(r)
	case Register16:
		return byte(r & 0x0f)
	case Register32:
		return byte(r & 0x0f)
	case Register64:
		return byte(r & 0x0f)
	case KRegister:
		return byte(r & 0x0f)
	case MMRegister:
		return byte(r & 0x0f)
	case XMMRegister:
		return byte(r & 0x0f)
	case YMMRegister:
		return byte(r & 0x0f)
	case ZMMRegister:
		return byte(r & 0x0f)
	case MaskedRegister:
		return hlcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func ehcode(v interface{}) byte {
	switch r := v.(type) {
	case Register8:
		return byte(r>>3) & 0x03
	case Register16:
		return byte(r>>3) & 0x03
	case Register32:
		return byte(r>>3) & 0x03
	case Register64:
		return byte(r>>3) & 0x03
	case KRegister:
		return byte(r>>3) & 0x03
	case MMRegister:
		return byte(r>>3) & 0x03
	case XMMRegister:
		return byte(r>>3) & 0x03
	case YMMRegister:
		return byte(r>>3) & 0x03
	case ZMMRegister:
		return byte(r>>3) & 0x03
	case MaskedRegister:
		return ehcode(r.Reg)
	default:
		panic("v is not a register")
	}
}

func toImmAny(v interface{}) int64 {
	if x, ok := asInt64(v); ok {
		return x
	} else {
		panic("value is not an integer")
	}
}

func toHcodeOpt(v interface{}) byte {
	if v == nil {
		return 0
	} else {
		return hcode(v)
	}
}

func toEcodeVMM(v interface{}, x byte) byte {
	switch r := v.(type) {
	case XMMRegister:
		return ecode(r)
	case YMMRegister:
		return ecode(r)
	case ZMMRegister:
		return ecode(r)
	default:
		return x
	}
}

func toKcodeMem(v *MemoryOperand) byte {
	if !v.Masked {
		return 0
	} else {
		return byte(v.Mask.K)
	}
}

func toZcodeMem(v *MemoryOperand) byte {
	if !v.Masked || v.Mask.Z {
		return 0
	} else {
		return 1
	}
}

func toZcodeRegM(v RegisterMask) byte {
	if v.Z {
		return 1
	} else {
		return 0
	}
}

func toHLcodeReg8(v Register8) byte {
	switch v {
	case AH:
		fallthrough
	case BH:
		fallthrough
	case CH:
		fallthrough
	case DH:
		panic("ah/bh/ch/dh registers never use 4-bit encoding")
	default:
		return byte(v & 0x0f)
	}
}



const (
	_N_inst = 16
)

const (
	_F_rel1 = 1 << iota
	_F_rel4
)

type _Encoding struct {
	len     int
	flags   int
	bytes   [_N_inst]byte
	encoder func(m *_Encoding, v []interface{})
}


func (self *_Encoding) buf(n int) []byte {
	if i := self.len; i+n > _N_inst {
		panic("instruction too long")
	} else {
		return self.bytes[i:]
	}
}


func (self *_Encoding) emit(v byte) {
	self.buf(1)[0] = v
	self.len++
}


func (self *_Encoding) imm1(v int64) {
	self.emit(byte(v))
}


func (self *_Encoding) imm2(v int64) {
	binary.LittleEndian.PutUint16(self.buf(2), uint16(v))
	self.len += 2
}


func (self *_Encoding) imm4(v int64) {
	binary.LittleEndian.PutUint32(self.buf(4), uint32(v))
	self.len += 4
}


func (self *_Encoding) imm8(v int64) {
	binary.LittleEndian.PutUint64(self.buf(8), uint64(v))
	self.len += 8
}





































func (self *_Encoding) vex2(lpp byte, r byte, rm interface{}, vvvv byte) {
	var b byte
	var x byte

	
	if r > 1 {
		panic("VEX.R must be a 1-bit mask")
	}

	
	if lpp&^0b111 != 0 {
		panic("VEX.Lpp must be a 3-bit mask")
	}

	
	if vvvv&^0b1111 != 0 {
		panic("VEX.vvvv must be a 4-bit mask")
	}

	
	if rm != nil {
		switch v := rm.(type) {
		case *Label:
			break
		case Register:
			b = hcode(v)
		case MemoryAddress:
			b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
		case RelativeOffset:
			break
		default:
			panic("rm is expected to be a register or a memory address")
		}
	}

	
	if x == 0 && b == 0 {
		self.emit(0xc5)
		self.emit(0xf8 ^ (r << 7) ^ (vvvv << 3) ^ lpp)
	} else {
		self.emit(0xc4)
		self.emit(0xe1 ^ (r << 7) ^ (x << 6) ^ (b << 5))
		self.emit(0x78 ^ (vvvv << 3) ^ lpp)
	}
}





















func (self *_Encoding) vex3(esc byte, mmmmm byte, wlpp byte, r byte, rm interface{}, vvvv byte) {
	var b byte
	var x byte

	
	if r > 1 {
		panic("VEX.R must be a 1-bit mask")
	}

	
	if vvvv&^0b1111 != 0 {
		panic("VEX.vvvv must be a 4-bit mask")
	}

	
	if esc != 0xc4 && esc != 0x8f {
		panic("escape must be a 3-byte VEX (0xc4) or XOP (0x8f) prefix")
	}

	
	if wlpp&^0b10000111 != 0 {
		panic("VEX.W____Lpp is expected to have no bits set except 0, 1, 2 and 7")
	}

	
	if mmmmm&^0b11111 != 0 {
		panic("VEX.m-mmmm is expected to be a 5-bit mask")
	}

	
	switch v := rm.(type) {
	case *Label:
		break
	case MemoryAddress:
		b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
	case RelativeOffset:
		break
	default:
		panic("rm is expected to be a register or a memory address")
	}

	
	self.emit(esc)
	self.emit(0xe0 ^ (r << 7) ^ (x << 6) ^ (b << 5) ^ mmmmm)
	self.emit(0x78 ^ (vvvv << 3) ^ wlpp)
}


func (self *_Encoding) evex(mm byte, w1pp byte, ll byte, rr byte, rm interface{}, vvvvv byte, aaa byte, zz byte, bb byte) {
	var b byte
	var x byte

	
	if bb > 1 {
		panic("EVEX.b must be a 1-bit mask")
	}

	
	if zz > 1 {
		panic("EVEX.z must be a 1-bit mask")
	}

	
	if mm&^0b11 != 0 {
		panic("EVEX.mm must be a 2-bit mask")
	}

	
	if ll&^0b11 != 0 {
		panic("EVEX.L'L must be a 2-bit mask")
	}

	
	if rr&^0b11 != 0 {
		panic("EVEX.R'R must be a 2-bit mask")
	}

	
	if aaa&^0b111 != 0 {
		panic("EVEX.aaa must be a 3-bit mask")
	}

	
	if vvvvv&^0b11111 != 0 {
		panic("EVEX.v'vvvv must be a 5-bit mask")
	}

	
	if w1pp&^0b10000011 != 0b100 {
		panic("EVEX.W____1pp is expected to have no bits set except 0, 1, 2, and 7")
	}

	
	r1, r0 := rr>>1, rr&1
	v1, v0 := vvvvv>>4, vvvvv&0b1111

	
	if rm != nil {
		switch m := rm.(type) {
		case *Label:
			break
		case Register:
			b, x = hcode(m), ecode(m)
		case MemoryAddress:
			b, x, v1 = toHcodeOpt(m.Base), toHcodeOpt(m.Index), toEcodeVMM(m.Index, v1)
		case RelativeOffset:
			break
		default:
			panic("rm is expected to be a register or a memory address")
		}
	}

	
	p0 := (r0 << 7) | (x << 6) | (b << 5) | (r1 << 4) | mm
	p1 := (v0 << 3) | w1pp
	p2 := (zz << 7) | (ll << 5) | (b << 4) | (v1 << 3) | aaa

	
	self.emit(0x62)
	self.emit(p0 ^ 0xf0)
	self.emit(p1 ^ 0x78)
	self.emit(p2 ^ 0x08)
}


func (self *_Encoding) rexm(w byte, r byte, rm interface{}) {
	var b byte
	var x byte

	
	if r != 0 && r != 1 {
		panic("REX.R must be 0 or 1")
	}

	
	if w != 0 && w != 1 {
		panic("REX.W must be 0 or 1")
	}

	
	switch v := rm.(type) {
	case *Label:
		break
	case MemoryAddress:
		b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
	case RelativeOffset:
		break
	default:
		panic("rm is expected to be a register or a memory address")
	}

	
	self.emit(0x40 | (w << 3) | (r << 2) | (x << 1) | b)
}


func (self *_Encoding) rexo(r byte, rm interface{}, force bool) {
	var b byte
	var x byte

	
	if r != 0 && r != 1 {
		panic("REX.R must be 0 or 1")
	}

	
	switch v := rm.(type) {
	case *Label:
		break
	case Register:
		b = hcode(v)
	case MemoryAddress:
		b, x = toHcodeOpt(v.Base), toHcodeOpt(v.Index)
	case RelativeOffset:
		break
	default:
		panic("rm is expected to be a register or a memory address")
	}

	
	if force || r != 0 || x != 0 || b != 0 {
		self.emit(0x40 | (r << 2) | (x << 1) | b)
	}
}














func (self *_Encoding) mrsd(reg byte, rm interface{}, disp8v int32) {
	var ok bool
	var mm MemoryAddress
	var ro RelativeOffset

	
	if reg > 7 {
		panic("invalid register bits")
	}

	
	switch disp8v {
	case 1:
		break
	case 2:
		break
	case 4:
		break
	case 8:
		break
	case 16:
		break
	case 32:
		break
	case 64:
		break
	default:
		panic("invalid displacement size")
	}

	
	if _, ok = rm.(*Label); ok {
		self.emit(0x05 | (reg << 3))
		self.imm4(0)
		return
	}

	
	if ro, ok = rm.(RelativeOffset); ok {
		self.emit(0x05 | (reg << 3))
		self.imm4(int64(ro))
		return
	}

	
	if mm, ok = rm.(MemoryAddress); !ok {
		panic("rm must be a memory address")
	}

	
	if mm.Base == nil && mm.Index == nil {
		self.emit(0x04 | (reg << 3))
		self.emit(0x25)
		self.imm4(int64(mm.Displacement))
		return
	}

	
	if mm.Index == nil && lcode(mm.Base) != 0b100 {
		cc := lcode(mm.Base)
		dv := mm.Displacement

		
		if dv == 0 && mm.Base != RBP && mm.Base != R13 {
			if cc == 0b101 {
				panic("rbp/r13 is not encodable as a base register (interpreted as disp32 address)")
			} else {
				self.emit((reg << 3) | cc)
				return
			}
		}

		
		if dq := dv / disp8v; dq >= math.MinInt8 && dq <= math.MaxInt8 && dv%disp8v == 0 {
			self.emit(0x40 | (reg << 3) | cc)
			self.imm1(int64(dq))
			return
		}

		
		self.emit(0x80 | (reg << 3) | cc)
		self.imm4(int64(mm.Displacement))
		return
	}

	
	if mm.Index == RSP {
		panic("rsp is not encodable as an index register (interpreted as no index)")
	}

	
	var scale byte
	var index byte = 0x04

	
	if mm.Scale != 0 {
		switch mm.Scale {
		case 1:
			scale = 0
		case 2:
			scale = 1
		case 4:
			scale = 2
		case 8:
			scale = 3
		default:
			panic("invalid scale value")
		}
	}

	
	if mm.Index != nil {
		index = lcode(mm.Index)
	}

	
	if mm.Base == nil {
		self.emit((reg << 3) | 0b100)
		self.emit((scale << 6) | (index << 3) | 0b101)
		self.imm4(int64(mm.Displacement))
		return
	}

	
	cc := lcode(mm.Base)
	dv := mm.Displacement

	
	if dv == 0 && cc != 0b101 {
		self.emit((reg << 3) | 0b100)
		self.emit((scale << 6) | (index << 3) | cc)
		return
	}

	
	if dq := dv / disp8v; dq >= math.MinInt8 && dq <= math.MaxInt8 && dv%disp8v == 0 {
		self.emit(0x44 | (reg << 3))
		self.emit((scale << 6) | (index << 3) | cc)
		self.imm1(int64(dq))
		return
	}

	
	self.emit(0x84 | (reg << 3))
	self.emit((scale << 6) | (index << 3) | cc)
	self.imm4(int64(mm.Displacement))
}


func (self *_Encoding) encode(v []interface{}) int {
	self.len = 0
	self.encoder(self, v)
	return self.len
}
