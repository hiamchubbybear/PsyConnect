



package wasm

import (
	"bytes"
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/objabi"
	"github.com/twitchyliquid64/golang-asm/sys"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

var Register = map[string]int16{
	"SP":    REG_SP,
	"CTXT":  REG_CTXT,
	"g":     REG_g,
	"RET0":  REG_RET0,
	"RET1":  REG_RET1,
	"RET2":  REG_RET2,
	"RET3":  REG_RET3,
	"PAUSE": REG_PAUSE,

	"R0":  REG_R0,
	"R1":  REG_R1,
	"R2":  REG_R2,
	"R3":  REG_R3,
	"R4":  REG_R4,
	"R5":  REG_R5,
	"R6":  REG_R6,
	"R7":  REG_R7,
	"R8":  REG_R8,
	"R9":  REG_R9,
	"R10": REG_R10,
	"R11": REG_R11,
	"R12": REG_R12,
	"R13": REG_R13,
	"R14": REG_R14,
	"R15": REG_R15,

	"F0":  REG_F0,
	"F1":  REG_F1,
	"F2":  REG_F2,
	"F3":  REG_F3,
	"F4":  REG_F4,
	"F5":  REG_F5,
	"F6":  REG_F6,
	"F7":  REG_F7,
	"F8":  REG_F8,
	"F9":  REG_F9,
	"F10": REG_F10,
	"F11": REG_F11,
	"F12": REG_F12,
	"F13": REG_F13,
	"F14": REG_F14,
	"F15": REG_F15,

	"F16": REG_F16,
	"F17": REG_F17,
	"F18": REG_F18,
	"F19": REG_F19,
	"F20": REG_F20,
	"F21": REG_F21,
	"F22": REG_F22,
	"F23": REG_F23,
	"F24": REG_F24,
	"F25": REG_F25,
	"F26": REG_F26,
	"F27": REG_F27,
	"F28": REG_F28,
	"F29": REG_F29,
	"F30": REG_F30,
	"F31": REG_F31,

	"PC_B": REG_PC_B,
}

var registerNames []string

func init() {
	obj.RegisterRegister(MINREG, MAXREG, rconv)
	obj.RegisterOpcode(obj.ABaseWasm, Anames)

	registerNames = make([]string, MAXREG-MINREG)
	for name, reg := range Register {
		registerNames[reg-MINREG] = name
	}
}

func rconv(r int) string {
	return registerNames[r-MINREG]
}

var unaryDst = map[obj.As]bool{
	ASet:          true,
	ATee:          true,
	ACall:         true,
	ACallIndirect: true,
	ACallImport:   true,
	ABr:           true,
	ABrIf:         true,
	ABrTable:      true,
	AI32Store:     true,
	AI64Store:     true,
	AF32Store:     true,
	AF64Store:     true,
	AI32Store8:    true,
	AI32Store16:   true,
	AI64Store8:    true,
	AI64Store16:   true,
	AI64Store32:   true,
	ACALLNORESUME: true,
}

var Linkwasm = obj.LinkArch{
	Arch:       sys.ArchWasm,
	Init:       instinit,
	Preprocess: preprocess,
	Assemble:   assemble,
	UnaryDst:   unaryDst,
}

var (
	morestack       *obj.LSym
	morestackNoCtxt *obj.LSym
	gcWriteBarrier  *obj.LSym
	sigpanic        *obj.LSym
	sigpanic0       *obj.LSym
	deferreturn     *obj.LSym
	jmpdefer        *obj.LSym
)

const (
	
	WasmImport = 1 << 0
)

func instinit(ctxt *obj.Link) {
	morestack = ctxt.Lookup("runtime.morestack")
	morestackNoCtxt = ctxt.Lookup("runtime.morestack_noctxt")
	gcWriteBarrier = ctxt.Lookup("runtime.gcWriteBarrier")
	sigpanic = ctxt.LookupABI("runtime.sigpanic", obj.ABIInternal)
	sigpanic0 = ctxt.LookupABI("runtime.sigpanic", 0) 
	deferreturn = ctxt.LookupABI("runtime.deferreturn", obj.ABIInternal)
	
	
	
	
	jmpdefer = ctxt.LookupABI(`"".jmpdefer`, obj.ABIInternal)
}

func preprocess(ctxt *obj.Link, s *obj.LSym, newprog obj.ProgAlloc) {
	appendp := func(p *obj.Prog, as obj.As, args ...obj.Addr) *obj.Prog {
		if p.As != obj.ANOP {
			p2 := obj.Appendp(p, newprog)
			p2.Pc = p.Pc
			p = p2
		}
		p.As = as
		switch len(args) {
		case 0:
			p.From = obj.Addr{}
			p.To = obj.Addr{}
		case 1:
			if unaryDst[as] {
				p.From = obj.Addr{}
				p.To = args[0]
			} else {
				p.From = args[0]
				p.To = obj.Addr{}
			}
		case 2:
			p.From = args[0]
			p.To = args[1]
		default:
			panic("bad args")
		}
		return p
	}

	framesize := s.Func.Text.To.Offset
	if framesize < 0 {
		panic("bad framesize")
	}
	s.Func.Args = s.Func.Text.To.Val.(int32)
	s.Func.Locals = int32(framesize)

	if s.Func.Text.From.Sym.Wrapper() {
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		

		gpanic := obj.Addr{
			Type:   obj.TYPE_MEM,
			Reg:    REGG,
			Offset: 4 * 8, 
		}

		panicargp := obj.Addr{
			Type:   obj.TYPE_MEM,
			Reg:    REG_R0,
			Offset: 0, 
		}

		p := s.Func.Text
		p = appendp(p, AMOVD, gpanic, regAddr(REG_R0))

		p = appendp(p, AGet, regAddr(REG_R0))
		p = appendp(p, AI64Eqz)
		p = appendp(p, ANot)
		p = appendp(p, AIf)

		p = appendp(p, AGet, regAddr(REG_SP))
		p = appendp(p, AI64ExtendI32U)
		p = appendp(p, AI64Const, constAddr(framesize+8))
		p = appendp(p, AI64Add)
		p = appendp(p, AI64Load, panicargp)

		p = appendp(p, AI64Eq)
		p = appendp(p, AIf)
		p = appendp(p, AMOVD, regAddr(REG_SP), panicargp)
		p = appendp(p, AEnd)

		p = appendp(p, AEnd)
	}

	if framesize > 0 {
		p := s.Func.Text
		p = appendp(p, AGet, regAddr(REG_SP))
		p = appendp(p, AI32Const, constAddr(framesize))
		p = appendp(p, AI32Sub)
		p = appendp(p, ASet, regAddr(REG_SP))
		p.Spadj = int32(framesize)
	}

	
	
	numResumePoints := 0
	explicitBlockDepth := 0
	pc := int64(0) 
	var tableIdxs []uint64
	tablePC := int64(0)
	base := ctxt.PosTable.Pos(s.Func.Text.Pos).Base()
	for p := s.Func.Text; p != nil; p = p.Link {
		prevBase := base
		base = ctxt.PosTable.Pos(p.Pos).Base()
		switch p.As {
		case ABlock, ALoop, AIf:
			explicitBlockDepth++

		case AEnd:
			if explicitBlockDepth == 0 {
				panic("End without block")
			}
			explicitBlockDepth--

		case ARESUMEPOINT:
			if explicitBlockDepth != 0 {
				panic("RESUME can only be used on toplevel")
			}
			p.As = AEnd
			for tablePC <= pc {
				tableIdxs = append(tableIdxs, uint64(numResumePoints))
				tablePC++
			}
			numResumePoints++
			pc++

		case obj.ACALL:
			if explicitBlockDepth != 0 {
				panic("CALL can only be used on toplevel, try CALLNORESUME instead")
			}
			appendp(p, ARESUMEPOINT)
		}

		p.Pc = pc

		
		
		
		
		if p.As == ACALLNORESUME || p.As == obj.ANOP || p.As == ANop || p.Spadj != 0 || base != prevBase {
			pc++
			if p.To.Sym == sigpanic {
				
				
				
				
				pc++
			}
		}
	}
	tableIdxs = append(tableIdxs, uint64(numResumePoints))
	s.Size = pc + 1

	if !s.Func.Text.From.Sym.NoSplit() {
		p := s.Func.Text

		if framesize <= objabi.StackSmall {
			
			
			
			
			
			

			p = appendp(p, AGet, regAddr(REG_SP))
			p = appendp(p, AGet, regAddr(REGG))
			p = appendp(p, AI32WrapI64)
			p = appendp(p, AI32Load, constAddr(2*int64(ctxt.Arch.PtrSize))) 
			p = appendp(p, AI32LeU)
		} else {
			
			
			
			
			
			
			
			
			

			p = appendp(p, AGet, regAddr(REG_SP))
			p = appendp(p, AGet, regAddr(REGG))
			p = appendp(p, AI32WrapI64)
			p = appendp(p, AI32Load, constAddr(2*int64(ctxt.Arch.PtrSize))) 
			p = appendp(p, AI32Const, constAddr(int64(framesize)-objabi.StackSmall))
			p = appendp(p, AI32Add)
			p = appendp(p, AI32LeU)
		}
		

		p = appendp(p, AIf)
		p = appendp(p, obj.ACALL, constAddr(0))
		if s.Func.Text.From.Sym.NeedCtxt() {
			p.To = obj.Addr{Type: obj.TYPE_MEM, Name: obj.NAME_EXTERN, Sym: morestack}
		} else {
			p.To = obj.Addr{Type: obj.TYPE_MEM, Name: obj.NAME_EXTERN, Sym: morestackNoCtxt}
		}
		p = appendp(p, AEnd)
	}

	
	
	var entryPointLoopBranches []*obj.Prog
	var unwindExitBranches []*obj.Prog
	currentDepth := 0
	for p := s.Func.Text; p != nil; p = p.Link {
		switch p.As {
		case ABlock, ALoop, AIf:
			currentDepth++
		case AEnd:
			currentDepth--
		}

		switch p.As {
		case obj.AJMP:
			jmp := *p
			p.As = obj.ANOP

			if jmp.To.Type == obj.TYPE_BRANCH {
				
				p = appendp(p, AI32Const, constAddr(jmp.To.Val.(*obj.Prog).Pc))
				p = appendp(p, ASet, regAddr(REG_PC_B)) 
				p = appendp(p, ABr)                     
				entryPointLoopBranches = append(entryPointLoopBranches, p)
				break
			}

			
			switch jmp.To.Type {
			case obj.TYPE_MEM:
				if !notUsePC_B[jmp.To.Sym.Name] {
					
					p = appendp(p, AI32Const, constAddr(0))
				}
				p = appendp(p, ACall, jmp.To)

			case obj.TYPE_NONE:
				
				p = appendp(p, AI32WrapI64)
				p = appendp(p, AI32Const, constAddr(16)) 
				p = appendp(p, AI32ShrU)

				
				
				
				
				p = appendp(p, ASet, regAddr(REG_PC_B))
				p = appendp(p, AI32Const, constAddr(0))
				p = appendp(p, AGet, regAddr(REG_PC_B))

				p = appendp(p, ACallIndirect)

			default:
				panic("bad target for JMP")
			}

			p = appendp(p, AReturn)

		case obj.ACALL, ACALLNORESUME:
			call := *p
			p.As = obj.ANOP

			pcAfterCall := call.Link.Pc
			if call.To.Sym == sigpanic {
				pcAfterCall-- 
			}

			
			
			if call.To.Sym == deferreturn {
				p = appendp(p, ALoop)
			}

			
			p = appendp(p, AGet, regAddr(REG_SP))
			p = appendp(p, AI32Const, constAddr(8))
			p = appendp(p, AI32Sub)
			p = appendp(p, ASet, regAddr(REG_SP))

			
			p = appendp(p, AGet, regAddr(REG_SP))
			p = appendp(p, AI64Const, obj.Addr{
				Type:   obj.TYPE_ADDR,
				Name:   obj.NAME_EXTERN,
				Sym:    s,           
				Offset: pcAfterCall, 
			})
			p = appendp(p, AI64Store, constAddr(0))

			
			switch call.To.Type {
			case obj.TYPE_MEM:
				if !notUsePC_B[call.To.Sym.Name] {
					
					p = appendp(p, AI32Const, constAddr(0))
				}
				p = appendp(p, ACall, call.To)

			case obj.TYPE_NONE:
				
				p = appendp(p, AI32WrapI64)
				p = appendp(p, AI32Const, constAddr(16)) 
				p = appendp(p, AI32ShrU)

				
				
				
				
				p = appendp(p, ASet, regAddr(REG_PC_B))
				p = appendp(p, AI32Const, constAddr(0))
				p = appendp(p, AGet, regAddr(REG_PC_B))

				p = appendp(p, ACallIndirect)

			default:
				panic("bad target for CALL")
			}

			
			if call.To.Sym == gcWriteBarrier {
				break
			}

			
			
			
			
			if call.To.Sym == jmpdefer {
				p = appendp(p, AReturn)
				break
			}

			
			if call.As == ACALLNORESUME && call.To.Sym != sigpanic && call.To.Sym != sigpanic0 { 
				
				p = appendp(p, AIf)
				p = appendp(p, obj.AUNDEF)
				p = appendp(p, AEnd)
			} else {
				
				p = appendp(p, ABrIf)
				unwindExitBranches = append(unwindExitBranches, p)
			}

			
			if call.To.Sym == deferreturn {
				
				p = appendp(p, AGet, regAddr(REG_SP))
				p = appendp(p, AI32Const, constAddr(8))
				p = appendp(p, AI32Sub)
				p = appendp(p, AI32Load16U, constAddr(0))
				p = appendp(p, ATee, regAddr(REG_PC_B))

				p = appendp(p, AI32Const, constAddr(call.Pc))
				p = appendp(p, AI32Eq)
				p = appendp(p, ABrIf, constAddr(0))
				p = appendp(p, AEnd) 
			}

		case obj.ARET, ARETUNWIND:
			ret := *p
			p.As = obj.ANOP

			if framesize > 0 {
				
				p = appendp(p, AGet, regAddr(REG_SP))
				p = appendp(p, AI32Const, constAddr(framesize))
				p = appendp(p, AI32Add)
				p = appendp(p, ASet, regAddr(REG_SP))
				
				
			}

			if ret.To.Type == obj.TYPE_MEM {
				
				p = appendp(p, AI32Const, constAddr(0))

				
				p = appendp(p, ACall, ret.To)
				p = appendp(p, AReturn)
				break
			}

			
			p = appendp(p, AGet, regAddr(REG_SP))
			p = appendp(p, AI32Const, constAddr(8))
			p = appendp(p, AI32Add)
			p = appendp(p, ASet, regAddr(REG_SP))

			if ret.As == ARETUNWIND {
				
				p = appendp(p, AI32Const, constAddr(1))
				p = appendp(p, AReturn)
				break
			}

			
			p = appendp(p, AI32Const, constAddr(0))
			p = appendp(p, AReturn)
		}
	}

	for p := s.Func.Text; p != nil; p = p.Link {
		switch p.From.Name {
		case obj.NAME_AUTO:
			p.From.Offset += int64(framesize)
		case obj.NAME_PARAM:
			p.From.Reg = REG_SP
			p.From.Offset += int64(framesize) + 8 
		}

		switch p.To.Name {
		case obj.NAME_AUTO:
			p.To.Offset += int64(framesize)
		case obj.NAME_PARAM:
			p.To.Reg = REG_SP
			p.To.Offset += int64(framesize) + 8 
		}

		switch p.As {
		case AGet:
			if p.From.Type == obj.TYPE_ADDR {
				get := *p
				p.As = obj.ANOP

				switch get.From.Name {
				case obj.NAME_EXTERN:
					p = appendp(p, AI64Const, get.From)
				case obj.NAME_AUTO, obj.NAME_PARAM:
					p = appendp(p, AGet, regAddr(get.From.Reg))
					if get.From.Reg == REG_SP {
						p = appendp(p, AI64ExtendI32U)
					}
					if get.From.Offset != 0 {
						p = appendp(p, AI64Const, constAddr(get.From.Offset))
						p = appendp(p, AI64Add)
					}
				default:
					panic("bad Get: invalid name")
				}
			}

		case AI32Load, AI64Load, AF32Load, AF64Load, AI32Load8S, AI32Load8U, AI32Load16S, AI32Load16U, AI64Load8S, AI64Load8U, AI64Load16S, AI64Load16U, AI64Load32S, AI64Load32U:
			if p.From.Type == obj.TYPE_MEM {
				as := p.As
				from := p.From

				p.As = AGet
				p.From = regAddr(from.Reg)

				if from.Reg != REG_SP {
					p = appendp(p, AI32WrapI64)
				}

				p = appendp(p, as, constAddr(from.Offset))
			}

		case AMOVB, AMOVH, AMOVW, AMOVD:
			mov := *p
			p.As = obj.ANOP

			var loadAs obj.As
			var storeAs obj.As
			switch mov.As {
			case AMOVB:
				loadAs = AI64Load8U
				storeAs = AI64Store8
			case AMOVH:
				loadAs = AI64Load16U
				storeAs = AI64Store16
			case AMOVW:
				loadAs = AI64Load32U
				storeAs = AI64Store32
			case AMOVD:
				loadAs = AI64Load
				storeAs = AI64Store
			}

			appendValue := func() {
				switch mov.From.Type {
				case obj.TYPE_CONST:
					p = appendp(p, AI64Const, constAddr(mov.From.Offset))

				case obj.TYPE_ADDR:
					switch mov.From.Name {
					case obj.NAME_NONE, obj.NAME_PARAM, obj.NAME_AUTO:
						p = appendp(p, AGet, regAddr(mov.From.Reg))
						if mov.From.Reg == REG_SP {
							p = appendp(p, AI64ExtendI32U)
						}
						p = appendp(p, AI64Const, constAddr(mov.From.Offset))
						p = appendp(p, AI64Add)
					case obj.NAME_EXTERN:
						p = appendp(p, AI64Const, mov.From)
					default:
						panic("bad name for MOV")
					}

				case obj.TYPE_REG:
					p = appendp(p, AGet, mov.From)
					if mov.From.Reg == REG_SP {
						p = appendp(p, AI64ExtendI32U)
					}

				case obj.TYPE_MEM:
					p = appendp(p, AGet, regAddr(mov.From.Reg))
					if mov.From.Reg != REG_SP {
						p = appendp(p, AI32WrapI64)
					}
					p = appendp(p, loadAs, constAddr(mov.From.Offset))

				default:
					panic("bad MOV type")
				}
			}

			switch mov.To.Type {
			case obj.TYPE_REG:
				appendValue()
				if mov.To.Reg == REG_SP {
					p = appendp(p, AI32WrapI64)
				}
				p = appendp(p, ASet, mov.To)

			case obj.TYPE_MEM:
				switch mov.To.Name {
				case obj.NAME_NONE, obj.NAME_PARAM:
					p = appendp(p, AGet, regAddr(mov.To.Reg))
					if mov.To.Reg != REG_SP {
						p = appendp(p, AI32WrapI64)
					}
				case obj.NAME_EXTERN:
					p = appendp(p, AI32Const, obj.Addr{Type: obj.TYPE_ADDR, Name: obj.NAME_EXTERN, Sym: mov.To.Sym})
				default:
					panic("bad MOV name")
				}
				appendValue()
				p = appendp(p, storeAs, constAddr(mov.To.Offset))

			default:
				panic("bad MOV type")
			}

		case ACallImport:
			p.As = obj.ANOP
			p = appendp(p, AGet, regAddr(REG_SP))
			p = appendp(p, ACall, obj.Addr{Type: obj.TYPE_MEM, Name: obj.NAME_EXTERN, Sym: s})
			p.Mark = WasmImport
		}
	}

	{
		p := s.Func.Text
		if len(unwindExitBranches) > 0 {
			p = appendp(p, ABlock) 
			for _, b := range unwindExitBranches {
				b.To = obj.Addr{Type: obj.TYPE_BRANCH, Val: p}
			}
		}
		if len(entryPointLoopBranches) > 0 {
			p = appendp(p, ALoop) 
			for _, b := range entryPointLoopBranches {
				b.To = obj.Addr{Type: obj.TYPE_BRANCH, Val: p}
			}
		}
		if numResumePoints > 0 {
			
			for i := 0; i < numResumePoints+1; i++ {
				p = appendp(p, ABlock)
			}
			p = appendp(p, AGet, regAddr(REG_PC_B)) 
			p = appendp(p, ABrTable, obj.Addr{Val: tableIdxs})
			p = appendp(p, AEnd) 
		}
		for p.Link != nil {
			p = p.Link 
		}
		if len(entryPointLoopBranches) > 0 {
			p = appendp(p, AEnd) 
		}
		p = appendp(p, obj.AUNDEF)
		if len(unwindExitBranches) > 0 {
			p = appendp(p, AEnd) 
			p = appendp(p, AI32Const, constAddr(1))
		}
	}

	currentDepth = 0
	blockDepths := make(map[*obj.Prog]int)
	for p := s.Func.Text; p != nil; p = p.Link {
		switch p.As {
		case ABlock, ALoop, AIf:
			currentDepth++
			blockDepths[p] = currentDepth
		case AEnd:
			currentDepth--
		}

		switch p.As {
		case ABr, ABrIf:
			if p.To.Type == obj.TYPE_BRANCH {
				blockDepth, ok := blockDepths[p.To.Val.(*obj.Prog)]
				if !ok {
					panic("label not at block")
				}
				p.To = constAddr(int64(currentDepth - blockDepth))
			}
		}
	}
}

func constAddr(value int64) obj.Addr {
	return obj.Addr{Type: obj.TYPE_CONST, Offset: value}
}

func regAddr(reg int16) obj.Addr {
	return obj.Addr{Type: obj.TYPE_REG, Reg: reg}
}



var notUsePC_B = map[string]bool{
	"_rt0_wasm_js":           true,
	"wasm_export_run":        true,
	"wasm_export_resume":     true,
	"wasm_export_getsp":      true,
	"wasm_pc_f_loop":         true,
	"runtime.wasmMove":       true,
	"runtime.wasmZero":       true,
	"runtime.wasmDiv":        true,
	"runtime.wasmTruncS":     true,
	"runtime.wasmTruncU":     true,
	"runtime.gcWriteBarrier": true,
	"cmpbody":                true,
	"memeqbody":              true,
	"memcmp":                 true,
	"memchr":                 true,
}

func assemble(ctxt *obj.Link, s *obj.LSym, newprog obj.ProgAlloc) {
	type regVar struct {
		global bool
		index  uint64
	}

	type varDecl struct {
		count uint64
		typ   valueType
	}

	hasLocalSP := false
	regVars := [MAXREG - MINREG]*regVar{
		REG_SP - MINREG:    {true, 0},
		REG_CTXT - MINREG:  {true, 1},
		REG_g - MINREG:     {true, 2},
		REG_RET0 - MINREG:  {true, 3},
		REG_RET1 - MINREG:  {true, 4},
		REG_RET2 - MINREG:  {true, 5},
		REG_RET3 - MINREG:  {true, 6},
		REG_PAUSE - MINREG: {true, 7},
	}
	var varDecls []*varDecl
	useAssemblyRegMap := func() {
		for i := int16(0); i < 16; i++ {
			regVars[REG_R0+i-MINREG] = &regVar{false, uint64(i)}
		}
	}

	
	
	switch s.Name {
	case "_rt0_wasm_js", "wasm_export_run", "wasm_export_resume", "wasm_export_getsp", "wasm_pc_f_loop",
		"runtime.wasmMove", "runtime.wasmZero", "runtime.wasmDiv", "runtime.wasmTruncS", "runtime.wasmTruncU", "memeqbody":
		varDecls = []*varDecl{}
		useAssemblyRegMap()
	case "memchr", "memcmp":
		varDecls = []*varDecl{{count: 2, typ: i32}}
		useAssemblyRegMap()
	case "cmpbody":
		varDecls = []*varDecl{{count: 2, typ: i64}}
		useAssemblyRegMap()
	case "runtime.gcWriteBarrier":
		varDecls = []*varDecl{{count: 4, typ: i64}}
		useAssemblyRegMap()
	default:
		
		regVars[REG_PC_B-MINREG] = &regVar{false, 0}
		hasLocalSP = true

		var regUsed [MAXREG - MINREG]bool
		for p := s.Func.Text; p != nil; p = p.Link {
			if p.From.Reg != 0 {
				regUsed[p.From.Reg-MINREG] = true
			}
			if p.To.Reg != 0 {
				regUsed[p.To.Reg-MINREG] = true
			}
		}

		regs := []int16{REG_SP}
		for reg := int16(REG_R0); reg <= REG_F31; reg++ {
			if regUsed[reg-MINREG] {
				regs = append(regs, reg)
			}
		}

		var lastDecl *varDecl
		for i, reg := range regs {
			t := regType(reg)
			if lastDecl == nil || lastDecl.typ != t {
				lastDecl = &varDecl{
					count: 0,
					typ:   t,
				}
				varDecls = append(varDecls, lastDecl)
			}
			lastDecl.count++
			if reg != REG_SP {
				regVars[reg-MINREG] = &regVar{false, 1 + uint64(i)}
			}
		}
	}

	w := new(bytes.Buffer)

	writeUleb128(w, uint64(len(varDecls)))
	for _, decl := range varDecls {
		writeUleb128(w, decl.count)
		w.WriteByte(byte(decl.typ))
	}

	if hasLocalSP {
		
		updateLocalSP(w)
	}

	for p := s.Func.Text; p != nil; p = p.Link {
		switch p.As {
		case AGet:
			if p.From.Type != obj.TYPE_REG {
				panic("bad Get: argument is not a register")
			}
			reg := p.From.Reg
			v := regVars[reg-MINREG]
			if v == nil {
				panic("bad Get: invalid register")
			}
			if reg == REG_SP && hasLocalSP {
				writeOpcode(w, ALocalGet)
				writeUleb128(w, 1) 
				continue
			}
			if v.global {
				writeOpcode(w, AGlobalGet)
			} else {
				writeOpcode(w, ALocalGet)
			}
			writeUleb128(w, v.index)
			continue

		case ASet:
			if p.To.Type != obj.TYPE_REG {
				panic("bad Set: argument is not a register")
			}
			reg := p.To.Reg
			v := regVars[reg-MINREG]
			if v == nil {
				panic("bad Set: invalid register")
			}
			if reg == REG_SP && hasLocalSP {
				writeOpcode(w, ALocalTee)
				writeUleb128(w, 1) 
			}
			if v.global {
				writeOpcode(w, AGlobalSet)
			} else {
				if p.Link.As == AGet && p.Link.From.Reg == reg {
					writeOpcode(w, ALocalTee)
					p = p.Link
				} else {
					writeOpcode(w, ALocalSet)
				}
			}
			writeUleb128(w, v.index)
			continue

		case ATee:
			if p.To.Type != obj.TYPE_REG {
				panic("bad Tee: argument is not a register")
			}
			reg := p.To.Reg
			v := regVars[reg-MINREG]
			if v == nil {
				panic("bad Tee: invalid register")
			}
			writeOpcode(w, ALocalTee)
			writeUleb128(w, v.index)
			continue

		case ANot:
			writeOpcode(w, AI32Eqz)
			continue

		case obj.AUNDEF:
			writeOpcode(w, AUnreachable)
			continue

		case obj.ANOP, obj.ATEXT, obj.AFUNCDATA, obj.APCDATA:
			
			continue
		}

		writeOpcode(w, p.As)

		switch p.As {
		case ABlock, ALoop, AIf:
			if p.From.Offset != 0 {
				
				w.WriteByte(0x80 - byte(p.From.Offset))
				continue
			}
			w.WriteByte(0x40)

		case ABr, ABrIf:
			if p.To.Type != obj.TYPE_CONST {
				panic("bad Br/BrIf")
			}
			writeUleb128(w, uint64(p.To.Offset))

		case ABrTable:
			idxs := p.To.Val.([]uint64)
			writeUleb128(w, uint64(len(idxs)-1))
			for _, idx := range idxs {
				writeUleb128(w, idx)
			}

		case ACall:
			switch p.To.Type {
			case obj.TYPE_CONST:
				writeUleb128(w, uint64(p.To.Offset))

			case obj.TYPE_MEM:
				if p.To.Name != obj.NAME_EXTERN && p.To.Name != obj.NAME_STATIC {
					fmt.Println(p.To)
					panic("bad name for Call")
				}
				r := obj.Addrel(s)
				r.Off = int32(w.Len())
				r.Type = objabi.R_CALL
				if p.Mark&WasmImport != 0 {
					r.Type = objabi.R_WASMIMPORT
				}
				r.Sym = p.To.Sym
				if hasLocalSP {
					
					updateLocalSP(w)
				}

			default:
				panic("bad type for Call")
			}

		case ACallIndirect:
			writeUleb128(w, uint64(p.To.Offset))
			w.WriteByte(0x00) 
			if hasLocalSP {
				
				updateLocalSP(w)
			}

		case AI32Const, AI64Const:
			if p.From.Name == obj.NAME_EXTERN {
				r := obj.Addrel(s)
				r.Off = int32(w.Len())
				r.Type = objabi.R_ADDR
				r.Sym = p.From.Sym
				r.Add = p.From.Offset
				break
			}
			writeSleb128(w, p.From.Offset)

		case AF32Const:
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, math.Float32bits(float32(p.From.Val.(float64))))
			w.Write(b)

		case AF64Const:
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, math.Float64bits(p.From.Val.(float64)))
			w.Write(b)

		case AI32Load, AI64Load, AF32Load, AF64Load, AI32Load8S, AI32Load8U, AI32Load16S, AI32Load16U, AI64Load8S, AI64Load8U, AI64Load16S, AI64Load16U, AI64Load32S, AI64Load32U:
			if p.From.Offset < 0 {
				panic("negative offset for *Load")
			}
			if p.From.Type != obj.TYPE_CONST {
				panic("bad type for *Load")
			}
			if p.From.Offset > math.MaxUint32 {
				ctxt.Diag("bad offset in %v", p)
			}
			writeUleb128(w, align(p.As))
			writeUleb128(w, uint64(p.From.Offset))

		case AI32Store, AI64Store, AF32Store, AF64Store, AI32Store8, AI32Store16, AI64Store8, AI64Store16, AI64Store32:
			if p.To.Offset < 0 {
				panic("negative offset")
			}
			if p.From.Offset > math.MaxUint32 {
				ctxt.Diag("bad offset in %v", p)
			}
			writeUleb128(w, align(p.As))
			writeUleb128(w, uint64(p.To.Offset))

		case ACurrentMemory, AGrowMemory:
			w.WriteByte(0x00)

		}
	}

	w.WriteByte(0x0b) 

	s.P = w.Bytes()
}

func updateLocalSP(w *bytes.Buffer) {
	writeOpcode(w, AGlobalGet)
	writeUleb128(w, 0) 
	writeOpcode(w, ALocalSet)
	writeUleb128(w, 1) 
}

func writeOpcode(w *bytes.Buffer, as obj.As) {
	switch {
	case as < AUnreachable:
		panic(fmt.Sprintf("unexpected assembler op: %s", as))
	case as < AEnd:
		w.WriteByte(byte(as - AUnreachable + 0x00))
	case as < ADrop:
		w.WriteByte(byte(as - AEnd + 0x0B))
	case as < ALocalGet:
		w.WriteByte(byte(as - ADrop + 0x1A))
	case as < AI32Load:
		w.WriteByte(byte(as - ALocalGet + 0x20))
	case as < AI32TruncSatF32S:
		w.WriteByte(byte(as - AI32Load + 0x28))
	case as < ALast:
		w.WriteByte(0xFC)
		w.WriteByte(byte(as - AI32TruncSatF32S + 0x00))
	default:
		panic(fmt.Sprintf("unexpected assembler op: %s", as))
	}
}

type valueType byte

const (
	i32 valueType = 0x7F
	i64 valueType = 0x7E
	f32 valueType = 0x7D
	f64 valueType = 0x7C
)

func regType(reg int16) valueType {
	switch {
	case reg == REG_SP:
		return i32
	case reg >= REG_R0 && reg <= REG_R15:
		return i64
	case reg >= REG_F0 && reg <= REG_F15:
		return f32
	case reg >= REG_F16 && reg <= REG_F31:
		return f64
	default:
		panic("invalid register")
	}
}

func align(as obj.As) uint64 {
	switch as {
	case AI32Load8S, AI32Load8U, AI64Load8S, AI64Load8U, AI32Store8, AI64Store8:
		return 0
	case AI32Load16S, AI32Load16U, AI64Load16S, AI64Load16U, AI32Store16, AI64Store16:
		return 1
	case AI32Load, AF32Load, AI64Load32S, AI64Load32U, AI32Store, AF32Store, AI64Store32:
		return 2
	case AI64Load, AF64Load, AI64Store, AF64Store:
		return 3
	default:
		panic("align: bad op")
	}
}

func writeUleb128(w io.ByteWriter, v uint64) {
	if v < 128 {
		w.WriteByte(uint8(v))
		return
	}
	more := true
	for more {
		c := uint8(v & 0x7f)
		v >>= 7
		more = v != 0
		if more {
			c |= 0x80
		}
		w.WriteByte(c)
	}
}

func writeSleb128(w io.ByteWriter, v int64) {
	more := true
	for more {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		more = !((v == 0 && s == 0) || (v == -1 && s != 0))
		if more {
			c |= 0x80
		}
		w.WriteByte(c)
	}
}
