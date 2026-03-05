

















package x86_64














func (self *Program) ADDQ(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("ADDQ", 2, Operands{v0, v1})
	
	if isImm32(v0) && v1 == RAX {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48)
			m.emit(0x05)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isImm8Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0x83)
			m.emit(0xc0 | lcode(v[1]))
			m.imm1(toImmAny(v[0]))
		})
	}
	
	if isImm32Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0x81)
			m.emit(0xc0 | lcode(v[1]))
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[0])<<2 | hcode(v[1]))
			m.emit(0x01)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x03)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isM64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x03)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isImm8Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0x83)
			m.mrsd(0, addr(v[1]), 1)
			m.imm1(toImmAny(v[0]))
		})
	}
	
	if isImm32Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0x81)
			m.mrsd(0, addr(v[1]), 1)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[0]), addr(v[1]))
			m.emit(0x01)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for ADDQ")
	}
	return p
}








func (self *Program) CALLQ(v0 interface{}) *Instruction {
	p := self.alloc("CALLQ", 1, Operands{v0})
	
	if isReg64(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(0, v[0], false)
			m.emit(0xff)
			m.emit(0xd0 | lcode(v[0]))
		})
	}
	
	if isM64(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(0, addr(v[0]), false)
			m.emit(0xff)
			m.mrsd(2, addr(v[0]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for CALLQ")
	}
	return p
}














func (self *Program) CMPQ(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("CMPQ", 2, Operands{v0, v1})
	
	if isImm32(v0) && v1 == RAX {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48)
			m.emit(0x3d)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isImm8Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0x83)
			m.emit(0xf8 | lcode(v[1]))
			m.imm1(toImmAny(v[0]))
		})
	}
	
	if isImm32Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0x81)
			m.emit(0xf8 | lcode(v[1]))
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[0])<<2 | hcode(v[1]))
			m.emit(0x39)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x3b)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isM64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x3b)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isImm8Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0x83)
			m.mrsd(7, addr(v[1]), 1)
			m.imm1(toImmAny(v[0]))
		})
	}
	
	if isImm32Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0x81)
			m.mrsd(7, addr(v[1]), 1)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[0]), addr(v[1]))
			m.emit(0x39)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for CMPQ")
	}
	return p
}








func (self *Program) JBE(v0 interface{}) *Instruction {
	p := self.alloc("JBE", 1, Operands{v0})
	p.branch = _B_conditional
	
	if isRel8(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x76)
			m.imm1(relv(v[0]))
		})
	}
	
	if isRel32(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x0f)
			m.emit(0x86)
			m.imm4(relv(v[0]))
		})
	}
	
	if isLabel(v0) {
		p.add(_F_rel1, func(m *_Encoding, v []interface{}) {
			m.emit(0x76)
			m.imm1(relv(v[0]))
		})
		p.add(_F_rel4, func(m *_Encoding, v []interface{}) {
			m.emit(0x0f)
			m.emit(0x86)
			m.imm4(relv(v[0]))
		})
	}
	if p.len == 0 {
		panic("invalid operands for JBE")
	}
	return p
}








func (self *Program) JMP(v0 interface{}) *Instruction {
	p := self.alloc("JMP", 1, Operands{v0})
	p.branch = _B_unconditional
	
	if isRel8(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xeb)
			m.imm1(relv(v[0]))
		})
	}
	
	if isRel32(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xe9)
			m.imm4(relv(v[0]))
		})
	}
	
	if isLabel(v0) {
		p.add(_F_rel1, func(m *_Encoding, v []interface{}) {
			m.emit(0xeb)
			m.imm1(relv(v[0]))
		})
		p.add(_F_rel4, func(m *_Encoding, v []interface{}) {
			m.emit(0xe9)
			m.imm4(relv(v[0]))
		})
	}
	if p.len == 0 {
		panic("invalid operands for JMP")
	}
	return p
}








func (self *Program) JMPQ(v0 interface{}) *Instruction {
	p := self.alloc("JMPQ", 1, Operands{v0})
	
	if isReg64(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(0, v[0], false)
			m.emit(0xff)
			m.emit(0xe0 | lcode(v[0]))
		})
	}
	
	if isM64(v0) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(0, addr(v[0]), false)
			m.emit(0xff)
			m.mrsd(4, addr(v[0]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for JMPQ")
	}
	return p
}







func (self *Program) LEAQ(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("LEAQ", 2, Operands{v0, v1})
	
	if isM(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x8d)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for LEAQ")
	}
	return p
}






















func (self *Program) MOVQ(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("MOVQ", 2, Operands{v0, v1})
	
	if isImm32Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0xc7)
			m.emit(0xc0 | lcode(v[1]))
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isImm64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0xb8 | lcode(v[1]))
			m.imm8(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[0])<<2 | hcode(v[1]))
			m.emit(0x89)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x8b)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isM64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x8b)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isImm32Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0xc7)
			m.mrsd(0, addr(v[1]), 1)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[0]), addr(v[1]))
			m.emit(0x89)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	
	if isMM(v0) && isReg64(v1) {
		self.require(ISA_MMX)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[0])<<2 | hcode(v[1]))
			m.emit(0x0f)
			m.emit(0x7e)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
	}
	
	if isReg64(v0) && isMM(v1) {
		self.require(ISA_MMX)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x0f)
			m.emit(0x6e)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isMM(v0) && isMM(v1) {
		self.require(ISA_MMX)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(hcode(v[1]), v[0], false)
			m.emit(0x0f)
			m.emit(0x6f)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(hcode(v[0]), v[1], false)
			m.emit(0x0f)
			m.emit(0x7f)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
	}
	
	if isM64(v0) && isMM(v1) {
		self.require(ISA_MMX)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(hcode(v[1]), addr(v[0]), false)
			m.emit(0x0f)
			m.emit(0x6f)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x0f)
			m.emit(0x6e)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isMM(v0) && isM64(v1) {
		self.require(ISA_MMX)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(hcode(v[0]), addr(v[1]), false)
			m.emit(0x0f)
			m.emit(0x7f)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[0]), addr(v[1]))
			m.emit(0x0f)
			m.emit(0x7e)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	
	if isXMM(v0) && isReg64(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x66)
			m.emit(0x48 | hcode(v[0])<<2 | hcode(v[1]))
			m.emit(0x0f)
			m.emit(0x7e)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
	}
	
	if isReg64(v0) && isXMM(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x66)
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x0f)
			m.emit(0x6e)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isXMM(v0) && isXMM(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf3)
			m.rexo(hcode(v[1]), v[0], false)
			m.emit(0x0f)
			m.emit(0x7e)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x66)
			m.rexo(hcode(v[0]), v[1], false)
			m.emit(0x0f)
			m.emit(0xd6)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
	}
	
	if isM64(v0) && isXMM(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf3)
			m.rexo(hcode(v[1]), addr(v[0]), false)
			m.emit(0x0f)
			m.emit(0x7e)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x66)
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x0f)
			m.emit(0x6e)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isXMM(v0) && isM64(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x66)
			m.rexo(hcode(v[0]), addr(v[1]), false)
			m.emit(0x0f)
			m.emit(0xd6)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x66)
			m.rexm(1, hcode(v[0]), addr(v[1]))
			m.emit(0x0f)
			m.emit(0x7e)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for MOVQ")
	}
	return p
}









func (self *Program) MOVSD(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("MOVSD", 2, Operands{v0, v1})
	
	if isXMM(v0) && isXMM(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf2)
			m.rexo(hcode(v[1]), v[0], false)
			m.emit(0x0f)
			m.emit(0x10)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf2)
			m.rexo(hcode(v[0]), v[1], false)
			m.emit(0x0f)
			m.emit(0x11)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
	}
	
	if isM64(v0) && isXMM(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf2)
			m.rexo(hcode(v[1]), addr(v[0]), false)
			m.emit(0x0f)
			m.emit(0x10)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isXMM(v0) && isM64(v1) {
		self.require(ISA_SSE2)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf2)
			m.rexo(hcode(v[0]), addr(v[1]), false)
			m.emit(0x0f)
			m.emit(0x11)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for MOVSD")
	}
	return p
}








func (self *Program) MOVSLQ(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("MOVSLQ", 2, Operands{v0, v1})
	
	if isReg32(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x63)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isM32(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x63)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for MOVSLQ")
	}
	return p
}









func (self *Program) MOVSS(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("MOVSS", 2, Operands{v0, v1})
	
	if isXMM(v0) && isXMM(v1) {
		self.require(ISA_SSE)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf3)
			m.rexo(hcode(v[1]), v[0], false)
			m.emit(0x0f)
			m.emit(0x10)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf3)
			m.rexo(hcode(v[0]), v[1], false)
			m.emit(0x0f)
			m.emit(0x11)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
	}
	
	if isM32(v0) && isXMM(v1) {
		self.require(ISA_SSE)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf3)
			m.rexo(hcode(v[1]), addr(v[0]), false)
			m.emit(0x0f)
			m.emit(0x10)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isXMM(v0) && isM32(v1) {
		self.require(ISA_SSE)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xf3)
			m.rexo(hcode(v[0]), addr(v[1]), false)
			m.emit(0x0f)
			m.emit(0x11)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for MOVSS")
	}
	return p
}








func (self *Program) RET(vv ...interface{}) *Instruction {
	var p *Instruction
	switch len(vv) {
	case 0:
		p = self.alloc("RET", 0, Operands{})
	case 1:
		p = self.alloc("RET", 1, Operands{vv[0]})
	default:
		panic("instruction RET takes 0 or 1 operands")
	}
	
	if len(vv) == 0 {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xc3)
		})
	}
	
	if len(vv) == 1 && isImm16(vv[0]) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xc2)
			m.imm2(toImmAny(v[0]))
		})
	}
	if p.len == 0 {
		panic("invalid operands for RET")
	}
	return p
}














func (self *Program) SUBQ(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("SUBQ", 2, Operands{v0, v1})
	
	if isImm32(v0) && v1 == RAX {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48)
			m.emit(0x2d)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isImm8Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0x83)
			m.emit(0xe8 | lcode(v[1]))
			m.imm1(toImmAny(v[0]))
		})
	}
	
	if isImm32Ext(v0, 8) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1]))
			m.emit(0x81)
			m.emit(0xe8 | lcode(v[1]))
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[0])<<2 | hcode(v[1]))
			m.emit(0x29)
			m.emit(0xc0 | lcode(v[0])<<3 | lcode(v[1]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0x48 | hcode(v[1])<<2 | hcode(v[0]))
			m.emit(0x2b)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isM64(v0) && isReg64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[1]), addr(v[0]))
			m.emit(0x2b)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	
	if isImm8Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0x83)
			m.mrsd(5, addr(v[1]), 1)
			m.imm1(toImmAny(v[0]))
		})
	}
	
	if isImm32Ext(v0, 8) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, 0, addr(v[1]))
			m.emit(0x81)
			m.mrsd(5, addr(v[1]), 1)
			m.imm4(toImmAny(v[0]))
		})
	}
	
	if isReg64(v0) && isM64(v1) {
		p.domain = DomainGeneric
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexm(1, hcode(v[0]), addr(v[1]))
			m.emit(0x29)
			m.mrsd(lcode(v[0]), addr(v[1]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for SUBQ")
	}
	return p
}












func (self *Program) VPERMIL2PD(v0 interface{}, v1 interface{}, v2 interface{}, v3 interface{}, v4 interface{}) *Instruction {
	p := self.alloc("VPERMIL2PD", 5, Operands{v0, v1, v2, v3, v4})
	
	if isImm4(v0) && isXMM(v1) && isXMM(v2) && isXMM(v3) && isXMM(v4) {
		self.require(ISA_XOP)
		p.domain = DomainAMDSpecific
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xc4)
			m.emit(0xe3 ^ (hcode(v[4]) << 7) ^ (hcode(v[2]) << 5))
			m.emit(0x79 ^ (hlcode(v[3]) << 3))
			m.emit(0x49)
			m.emit(0xc0 | lcode(v[4])<<3 | lcode(v[2]))
			m.emit((hlcode(v[1]) << 4) | imml(v[0]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xc4)
			m.emit(0xe3 ^ (hcode(v[4]) << 7) ^ (hcode(v[1]) << 5))
			m.emit(0xf9 ^ (hlcode(v[3]) << 3))
			m.emit(0x49)
			m.emit(0xc0 | lcode(v[4])<<3 | lcode(v[1]))
			m.emit((hlcode(v[2]) << 4) | imml(v[0]))
		})
	}
	
	if isImm4(v0) && isM128(v1) && isXMM(v2) && isXMM(v3) && isXMM(v4) {
		self.require(ISA_XOP)
		p.domain = DomainAMDSpecific
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.vex3(0xc4, 0b11, 0x81, hcode(v[4]), addr(v[1]), hlcode(v[3]))
			m.emit(0x49)
			m.mrsd(lcode(v[4]), addr(v[1]), 1)
			m.emit((hlcode(v[2]) << 4) | imml(v[0]))
		})
	}
	
	if isImm4(v0) && isXMM(v1) && isM128(v2) && isXMM(v3) && isXMM(v4) {
		self.require(ISA_XOP)
		p.domain = DomainAMDSpecific
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.vex3(0xc4, 0b11, 0x01, hcode(v[4]), addr(v[2]), hlcode(v[3]))
			m.emit(0x49)
			m.mrsd(lcode(v[4]), addr(v[2]), 1)
			m.emit((hlcode(v[1]) << 4) | imml(v[0]))
		})
	}
	
	if isImm4(v0) && isYMM(v1) && isYMM(v2) && isYMM(v3) && isYMM(v4) {
		self.require(ISA_XOP)
		p.domain = DomainAMDSpecific
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xc4)
			m.emit(0xe3 ^ (hcode(v[4]) << 7) ^ (hcode(v[2]) << 5))
			m.emit(0x7d ^ (hlcode(v[3]) << 3))
			m.emit(0x49)
			m.emit(0xc0 | lcode(v[4])<<3 | lcode(v[2]))
			m.emit((hlcode(v[1]) << 4) | imml(v[0]))
		})
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.emit(0xc4)
			m.emit(0xe3 ^ (hcode(v[4]) << 7) ^ (hcode(v[1]) << 5))
			m.emit(0xfd ^ (hlcode(v[3]) << 3))
			m.emit(0x49)
			m.emit(0xc0 | lcode(v[4])<<3 | lcode(v[1]))
			m.emit((hlcode(v[2]) << 4) | imml(v[0]))
		})
	}
	
	if isImm4(v0) && isM256(v1) && isYMM(v2) && isYMM(v3) && isYMM(v4) {
		self.require(ISA_XOP)
		p.domain = DomainAMDSpecific
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.vex3(0xc4, 0b11, 0x85, hcode(v[4]), addr(v[1]), hlcode(v[3]))
			m.emit(0x49)
			m.mrsd(lcode(v[4]), addr(v[1]), 1)
			m.emit((hlcode(v[2]) << 4) | imml(v[0]))
		})
	}
	
	if isImm4(v0) && isYMM(v1) && isM256(v2) && isYMM(v3) && isYMM(v4) {
		self.require(ISA_XOP)
		p.domain = DomainAMDSpecific
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.vex3(0xc4, 0b11, 0x05, hcode(v[4]), addr(v[2]), hlcode(v[3]))
			m.emit(0x49)
			m.mrsd(lcode(v[4]), addr(v[2]), 1)
			m.emit((hlcode(v[1]) << 4) | imml(v[0]))
		})
	}
	if p.len == 0 {
		panic("invalid operands for VPERMIL2PD")
	}
	return p
}








func (self *Program) XORPS(v0 interface{}, v1 interface{}) *Instruction {
	p := self.alloc("XORPS", 2, Operands{v0, v1})
	
	if isXMM(v0) && isXMM(v1) {
		self.require(ISA_SSE)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(hcode(v[1]), v[0], false)
			m.emit(0x0f)
			m.emit(0x57)
			m.emit(0xc0 | lcode(v[1])<<3 | lcode(v[0]))
		})
	}
	
	if isM128(v0) && isXMM(v1) {
		self.require(ISA_SSE)
		p.domain = DomainMMXSSE
		p.add(0, func(m *_Encoding, v []interface{}) {
			m.rexo(hcode(v[1]), addr(v[0]), false)
			m.emit(0x0f)
			m.emit(0x57)
			m.mrsd(lcode(v[1]), addr(v[0]), 1)
		})
	}
	if p.len == 0 {
		panic("invalid operands for XORPS")
	}
	return p
}
