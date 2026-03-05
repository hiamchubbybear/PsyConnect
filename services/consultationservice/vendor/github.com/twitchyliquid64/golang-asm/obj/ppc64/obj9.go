




























package ppc64

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/objabi"
	"github.com/twitchyliquid64/golang-asm/sys"
)

func progedit(ctxt *obj.Link, p *obj.Prog, newprog obj.ProgAlloc) {
	p.From.Class = 0
	p.To.Class = 0

	c := ctxt9{ctxt: ctxt, newprog: newprog}

	
	switch p.As {
	case ABR,
		ABL,
		obj.ARET,
		obj.ADUFFZERO,
		obj.ADUFFCOPY:
		if p.To.Sym != nil {
			p.To.Type = obj.TYPE_BRANCH
		}
	}

	
	switch p.As {
	case AFMOVS:
		if p.From.Type == obj.TYPE_FCONST {
			f32 := float32(p.From.Val.(float64))
			p.From.Type = obj.TYPE_MEM
			p.From.Sym = ctxt.Float32Sym(f32)
			p.From.Name = obj.NAME_EXTERN
			p.From.Offset = 0
		}

	case AFMOVD:
		if p.From.Type == obj.TYPE_FCONST {
			f64 := p.From.Val.(float64)
			
			if f64 != 0 {
				p.From.Type = obj.TYPE_MEM
				p.From.Sym = ctxt.Float64Sym(f64)
				p.From.Name = obj.NAME_EXTERN
				p.From.Offset = 0
			}
		}

		
	case AMOVD:
		if p.From.Type == obj.TYPE_CONST && p.From.Name == obj.NAME_NONE && p.From.Reg == 0 && int64(int32(p.From.Offset)) != p.From.Offset {
			p.From.Type = obj.TYPE_MEM
			p.From.Sym = ctxt.Int64Sym(p.From.Offset)
			p.From.Name = obj.NAME_EXTERN
			p.From.Offset = 0
		}
	}

	
	switch p.As {
	case ASUBC:
		if p.From.Type == obj.TYPE_CONST {
			p.From.Offset = -p.From.Offset
			p.As = AADDC
		}

	case ASUBCCC:
		if p.From.Type == obj.TYPE_CONST {
			p.From.Offset = -p.From.Offset
			p.As = AADDCCC
		}

	case ASUB:
		if p.From.Type == obj.TYPE_CONST {
			p.From.Offset = -p.From.Offset
			p.As = AADD
		}
	}
	if c.ctxt.Headtype == objabi.Haix {
		c.rewriteToUseTOC(p)
	} else if c.ctxt.Flag_dynlink {
		c.rewriteToUseGot(p)
	}
}



func (c *ctxt9) rewriteToUseTOC(p *obj.Prog) {
	if p.As == obj.ATEXT || p.As == obj.AFUNCDATA || p.As == obj.ACALL || p.As == obj.ARET || p.As == obj.AJMP {
		return
	}

	if p.As == obj.ADUFFCOPY || p.As == obj.ADUFFZERO {
		
		
		if !c.ctxt.Flag_dynlink {
			return
		}
		
		
		
		
		
		
		var sym *obj.LSym
		if p.As == obj.ADUFFZERO {
			sym = c.ctxt.Lookup("runtime.duffzero")
		} else {
			sym = c.ctxt.Lookup("runtime.duffcopy")
		}
		
		symtoc := c.ctxt.LookupInit("TOC."+sym.Name, func(s *obj.LSym) {
			s.Type = objabi.SDATA
			s.Set(obj.AttrDuplicateOK, true)
			s.Set(obj.AttrStatic, true)
			c.ctxt.Data = append(c.ctxt.Data, s)
			s.WriteAddr(c.ctxt, 0, 8, sym, 0)
		})

		offset := p.To.Offset
		p.As = AMOVD
		p.From.Type = obj.TYPE_MEM
		p.From.Name = obj.NAME_TOCREF
		p.From.Sym = symtoc
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R12
		p.To.Name = obj.NAME_NONE
		p.To.Offset = 0
		p.To.Sym = nil
		p1 := obj.Appendp(p, c.newprog)
		p1.As = AADD
		p1.From.Type = obj.TYPE_CONST
		p1.From.Offset = offset
		p1.To.Type = obj.TYPE_REG
		p1.To.Reg = REG_R12
		p2 := obj.Appendp(p1, c.newprog)
		p2.As = AMOVD
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = REG_R12
		p2.To.Type = obj.TYPE_REG
		p2.To.Reg = REG_LR
		p3 := obj.Appendp(p2, c.newprog)
		p3.As = obj.ACALL
		p3.To.Type = obj.TYPE_REG
		p3.To.Reg = REG_LR
	}

	var source *obj.Addr
	if p.From.Name == obj.NAME_EXTERN || p.From.Name == obj.NAME_STATIC {
		if p.From.Type == obj.TYPE_ADDR {
			if p.As == ADWORD {
				
				return
			}
			if p.As != AMOVD {
				c.ctxt.Diag("do not know how to handle TYPE_ADDR in %v", p)
				return
			}
			if p.To.Type != obj.TYPE_REG {
				c.ctxt.Diag("do not know how to handle LEAQ-type insn to non-register in %v", p)
				return
			}
		} else if p.From.Type != obj.TYPE_MEM {
			c.ctxt.Diag("do not know how to handle %v without TYPE_MEM", p)
			return
		}
		source = &p.From

	} else if p.To.Name == obj.NAME_EXTERN || p.To.Name == obj.NAME_STATIC {
		if p.To.Type != obj.TYPE_MEM {
			c.ctxt.Diag("do not know how to handle %v without TYPE_MEM", p)
			return
		}
		if source != nil {
			c.ctxt.Diag("cannot handle symbols on both sides in %v", p)
			return
		}
		source = &p.To
	} else {
		return

	}

	if source.Sym == nil {
		c.ctxt.Diag("do not know how to handle nil symbol in %v", p)
		return
	}

	if source.Sym.Type == objabi.STLSBSS {
		return
	}

	
	symtoc := c.ctxt.LookupInit("TOC."+source.Sym.Name, func(s *obj.LSym) {
		s.Type = objabi.SDATA
		s.Set(obj.AttrDuplicateOK, true)
		s.Set(obj.AttrStatic, true)
		c.ctxt.Data = append(c.ctxt.Data, s)
		s.WriteAddr(c.ctxt, 0, 8, source.Sym, 0)
	})

	if source.Type == obj.TYPE_ADDR {
		
		
		p.From.Type = obj.TYPE_MEM
		p.From.Sym = symtoc
		p.From.Name = obj.NAME_TOCREF

		if p.From.Offset != 0 {
			q := obj.Appendp(p, c.newprog)
			q.As = AADD
			q.From.Type = obj.TYPE_CONST
			q.From.Offset = p.From.Offset
			p.From.Offset = 0
			q.To = p.To
		}
		return

	}

	
	
	

	q := obj.Appendp(p, c.newprog)
	q.As = AMOVD
	q.From.Type = obj.TYPE_MEM
	q.From.Sym = symtoc
	q.From.Name = obj.NAME_TOCREF
	q.To.Type = obj.TYPE_REG
	q.To.Reg = REGTMP

	q = obj.Appendp(q, c.newprog)
	q.As = p.As
	q.From = p.From
	q.To = p.To
	if p.From.Name != obj.NAME_NONE {
		q.From.Type = obj.TYPE_MEM
		q.From.Reg = REGTMP
		q.From.Name = obj.NAME_NONE
		q.From.Sym = nil
	} else if p.To.Name != obj.NAME_NONE {
		q.To.Type = obj.TYPE_MEM
		q.To.Reg = REGTMP
		q.To.Name = obj.NAME_NONE
		q.To.Sym = nil
	} else {
		c.ctxt.Diag("unreachable case in rewriteToUseTOC with %v", p)
	}

	obj.Nopout(p)
}


func (c *ctxt9) rewriteToUseGot(p *obj.Prog) {
	if p.As == obj.ADUFFCOPY || p.As == obj.ADUFFZERO {
		
		
		
		
		
		
		var sym *obj.LSym
		if p.As == obj.ADUFFZERO {
			sym = c.ctxt.Lookup("runtime.duffzero")
		} else {
			sym = c.ctxt.Lookup("runtime.duffcopy")
		}
		offset := p.To.Offset
		p.As = AMOVD
		p.From.Type = obj.TYPE_MEM
		p.From.Name = obj.NAME_GOTREF
		p.From.Sym = sym
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R12
		p.To.Name = obj.NAME_NONE
		p.To.Offset = 0
		p.To.Sym = nil
		p1 := obj.Appendp(p, c.newprog)
		p1.As = AADD
		p1.From.Type = obj.TYPE_CONST
		p1.From.Offset = offset
		p1.To.Type = obj.TYPE_REG
		p1.To.Reg = REG_R12
		p2 := obj.Appendp(p1, c.newprog)
		p2.As = AMOVD
		p2.From.Type = obj.TYPE_REG
		p2.From.Reg = REG_R12
		p2.To.Type = obj.TYPE_REG
		p2.To.Reg = REG_LR
		p3 := obj.Appendp(p2, c.newprog)
		p3.As = obj.ACALL
		p3.To.Type = obj.TYPE_REG
		p3.To.Reg = REG_LR
	}

	
	
	
	if p.From.Type == obj.TYPE_ADDR && p.From.Name == obj.NAME_EXTERN && !p.From.Sym.Local() {
		
		
		if p.As != AMOVD {
			c.ctxt.Diag("do not know how to handle TYPE_ADDR in %v with -dynlink", p)
		}
		if p.To.Type != obj.TYPE_REG {
			c.ctxt.Diag("do not know how to handle LEAQ-type insn to non-register in %v with -dynlink", p)
		}
		p.From.Type = obj.TYPE_MEM
		p.From.Name = obj.NAME_GOTREF
		if p.From.Offset != 0 {
			q := obj.Appendp(p, c.newprog)
			q.As = AADD
			q.From.Type = obj.TYPE_CONST
			q.From.Offset = p.From.Offset
			q.To = p.To
			p.From.Offset = 0
		}
	}
	if p.GetFrom3() != nil && p.GetFrom3().Name == obj.NAME_EXTERN {
		c.ctxt.Diag("don't know how to handle %v with -dynlink", p)
	}
	var source *obj.Addr
	
	
	
	if p.From.Name == obj.NAME_EXTERN && !p.From.Sym.Local() {
		if p.To.Name == obj.NAME_EXTERN && !p.To.Sym.Local() {
			c.ctxt.Diag("cannot handle NAME_EXTERN on both sides in %v with -dynlink", p)
		}
		source = &p.From
	} else if p.To.Name == obj.NAME_EXTERN && !p.To.Sym.Local() {
		source = &p.To
	} else {
		return
	}
	if p.As == obj.ATEXT || p.As == obj.AFUNCDATA || p.As == obj.ACALL || p.As == obj.ARET || p.As == obj.AJMP {
		return
	}
	if source.Sym.Type == objabi.STLSBSS {
		return
	}
	if source.Type != obj.TYPE_MEM {
		c.ctxt.Diag("don't know how to handle %v with -dynlink", p)
	}
	p1 := obj.Appendp(p, c.newprog)
	p2 := obj.Appendp(p1, c.newprog)

	p1.As = AMOVD
	p1.From.Type = obj.TYPE_MEM
	p1.From.Sym = source.Sym
	p1.From.Name = obj.NAME_GOTREF
	p1.To.Type = obj.TYPE_REG
	p1.To.Reg = REGTMP

	p2.As = p.As
	p2.From = p.From
	p2.To = p.To
	if p.From.Name == obj.NAME_EXTERN {
		p2.From.Reg = REGTMP
		p2.From.Name = obj.NAME_NONE
		p2.From.Sym = nil
	} else if p.To.Name == obj.NAME_EXTERN {
		p2.To.Reg = REGTMP
		p2.To.Name = obj.NAME_NONE
		p2.To.Sym = nil
	} else {
		return
	}
	obj.Nopout(p)
}

func preprocess(ctxt *obj.Link, cursym *obj.LSym, newprog obj.ProgAlloc) {
	
	if cursym.Func.Text == nil || cursym.Func.Text.Link == nil {
		return
	}

	c := ctxt9{ctxt: ctxt, cursym: cursym, newprog: newprog}

	p := c.cursym.Func.Text
	textstksiz := p.To.Offset
	if textstksiz == -8 {
		
		p.From.Sym.Set(obj.AttrNoFrame, true)
		textstksiz = 0
	}
	if textstksiz%8 != 0 {
		c.ctxt.Diag("frame size %d not a multiple of 8", textstksiz)
	}
	if p.From.Sym.NoFrame() {
		if textstksiz != 0 {
			c.ctxt.Diag("NOFRAME functions must have a frame size of 0, not %d", textstksiz)
		}
	}

	c.cursym.Func.Args = p.To.Val.(int32)
	c.cursym.Func.Locals = int32(textstksiz)

	

	var q *obj.Prog
	var q1 *obj.Prog
	for p := c.cursym.Func.Text; p != nil; p = p.Link {
		switch p.As {
		
		case obj.ATEXT:
			q = p

			p.Mark |= LABEL | LEAF | SYNC
			if p.Link != nil {
				p.Link.Mark |= LABEL
			}

		case ANOR:
			q = p
			if p.To.Type == obj.TYPE_REG {
				if p.To.Reg == REGZERO {
					p.Mark |= LABEL | SYNC
				}
			}

		case ALWAR,
			ALBAR,
			ASTBCCC,
			ASTWCCC,
			AECIWX,
			AECOWX,
			AEIEIO,
			AICBI,
			AISYNC,
			ATLBIE,
			ATLBIEL,
			ASLBIA,
			ASLBIE,
			ASLBMFEE,
			ASLBMFEV,
			ASLBMTE,
			ADCBF,
			ADCBI,
			ADCBST,
			ADCBT,
			ADCBTST,
			ADCBZ,
			ASYNC,
			ATLBSYNC,
			APTESYNC,
			ALWSYNC,
			ATW,
			AWORD,
			ARFI,
			ARFCI,
			ARFID,
			AHRFID:
			q = p
			p.Mark |= LABEL | SYNC
			continue

		case AMOVW, AMOVWZ, AMOVD:
			q = p
			if p.From.Reg >= REG_SPECIAL || p.To.Reg >= REG_SPECIAL {
				p.Mark |= LABEL | SYNC
			}
			continue

		case AFABS,
			AFABSCC,
			AFADD,
			AFADDCC,
			AFCTIW,
			AFCTIWCC,
			AFCTIWZ,
			AFCTIWZCC,
			AFDIV,
			AFDIVCC,
			AFMADD,
			AFMADDCC,
			AFMOVD,
			AFMOVDU,
			
			AFMOVS,
			AFMOVSU,

			
			AFMSUB,
			AFMSUBCC,
			AFMUL,
			AFMULCC,
			AFNABS,
			AFNABSCC,
			AFNEG,
			AFNEGCC,
			AFNMADD,
			AFNMADDCC,
			AFNMSUB,
			AFNMSUBCC,
			AFRSP,
			AFRSPCC,
			AFSUB,
			AFSUBCC:
			q = p

			p.Mark |= FLOAT
			continue

		case ABL,
			ABCL,
			obj.ADUFFZERO,
			obj.ADUFFCOPY:
			c.cursym.Func.Text.Mark &^= LEAF
			fallthrough

		case ABC,
			ABEQ,
			ABGE,
			ABGT,
			ABLE,
			ABLT,
			ABNE,
			ABR,
			ABVC,
			ABVS:
			p.Mark |= BRANCH
			q = p
			q1 = p.To.Target()
			if q1 != nil {
				

				if q1.Mark&LEAF == 0 {
					q1.Mark |= LABEL
				}
			} else {
				p.Mark |= LABEL
			}
			q1 = p.Link
			if q1 != nil {
				q1.Mark |= LABEL
			}
			continue

		case AFCMPO, AFCMPU:
			q = p
			p.Mark |= FCMP | FLOAT
			continue

		case obj.ARET:
			q = p
			if p.Link != nil {
				p.Link.Mark |= LABEL
			}
			continue

		case obj.ANOP:
			
			
			continue

		default:
			q = p
			continue
		}
	}

	autosize := int32(0)
	var p1 *obj.Prog
	var p2 *obj.Prog
	for p := c.cursym.Func.Text; p != nil; p = p.Link {
		o := p.As
		switch o {
		case obj.ATEXT:
			autosize = int32(textstksiz)

			if p.Mark&LEAF != 0 && autosize == 0 {
				
				p.From.Sym.Set(obj.AttrNoFrame, true)
			}

			if !p.From.Sym.NoFrame() {
				
				
				autosize += int32(c.ctxt.FixedFrameSize())
			}

			if p.Mark&LEAF != 0 && autosize < objabi.StackSmall {
				
				
				p.From.Sym.Set(obj.AttrNoSplit, true)
			}

			p.To.Offset = int64(autosize)

			q = p

			if c.ctxt.Flag_shared && c.cursym.Name != "runtime.duffzero" && c.cursym.Name != "runtime.duffcopy" {
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				q = obj.Appendp(q, c.newprog)
				q.As = AWORD
				q.Pos = p.Pos
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = 0x3c4c0000
				q = obj.Appendp(q, c.newprog)
				q.As = AWORD
				q.Pos = p.Pos
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = 0x38420000
				rel := obj.Addrel(c.cursym)
				rel.Off = 0
				rel.Siz = 8
				rel.Sym = c.ctxt.Lookup(".TOC.")
				rel.Type = objabi.R_ADDRPOWER_PCREL
			}

			if !c.cursym.Func.Text.From.Sym.NoSplit() {
				q = c.stacksplit(q, autosize) 
			}

			
			
			
			if autosize != 0 && c.cursym.Name != "runtime.racecallbackthunk" {
				
				
				
				if autosize >= -BIG && autosize <= BIG {
					
					q = obj.Appendp(q, c.newprog)
					q.As = AMOVD
					q.Pos = p.Pos
					q.From.Type = obj.TYPE_REG
					q.From.Reg = REG_LR
					q.To.Type = obj.TYPE_REG
					q.To.Reg = REGTMP

					q = obj.Appendp(q, c.newprog)
					q.As = AMOVDU
					q.Pos = p.Pos
					q.From.Type = obj.TYPE_REG
					q.From.Reg = REGTMP
					q.To.Type = obj.TYPE_MEM
					q.To.Offset = int64(-autosize)
					q.To.Reg = REGSP
					q.Spadj = autosize
				} else {
					
					
					
					
					
					
					q = obj.Appendp(q, c.newprog)
					q.As = AMOVD
					q.Pos = p.Pos
					q.From.Type = obj.TYPE_REG
					q.From.Reg = REG_LR
					q.To.Type = obj.TYPE_REG
					q.To.Reg = REG_R29 

					q = c.ctxt.StartUnsafePoint(q, c.newprog)

					q = obj.Appendp(q, c.newprog)
					q.As = AMOVD
					q.Pos = p.Pos
					q.From.Type = obj.TYPE_REG
					q.From.Reg = REG_R29
					q.To.Type = obj.TYPE_MEM
					q.To.Offset = int64(-autosize)
					q.To.Reg = REGSP

					q = obj.Appendp(q, c.newprog)
					q.As = AADD
					q.Pos = p.Pos
					q.From.Type = obj.TYPE_CONST
					q.From.Offset = int64(-autosize)
					q.To.Type = obj.TYPE_REG
					q.To.Reg = REGSP
					q.Spadj = +autosize

					q = c.ctxt.EndUnsafePoint(q, c.newprog, -1)

				}
			} else if c.cursym.Func.Text.Mark&LEAF == 0 {
				
				
				
				c.cursym.Func.Text.Mark |= LEAF
			}

			if c.cursym.Func.Text.Mark&LEAF != 0 {
				c.cursym.Set(obj.AttrLeaf, true)
				break
			}

			if c.ctxt.Flag_shared {
				q = obj.Appendp(q, c.newprog)
				q.As = AMOVD
				q.Pos = p.Pos
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_R2
				q.To.Type = obj.TYPE_MEM
				q.To.Reg = REGSP
				q.To.Offset = 24
			}

			if c.cursym.Func.Text.From.Sym.Wrapper() {
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				
				

				q = obj.Appendp(q, c.newprog)

				q.As = AMOVD
				q.From.Type = obj.TYPE_MEM
				q.From.Reg = REGG
				q.From.Offset = 4 * int64(c.ctxt.Arch.PtrSize) 
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_R3

				q = obj.Appendp(q, c.newprog)
				q.As = ACMP
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_R0
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_R3

				q = obj.Appendp(q, c.newprog)
				q.As = ABEQ
				q.To.Type = obj.TYPE_BRANCH
				p1 = q

				q = obj.Appendp(q, c.newprog)
				q.As = AMOVD
				q.From.Type = obj.TYPE_MEM
				q.From.Reg = REG_R3
				q.From.Offset = 0 
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_R4

				q = obj.Appendp(q, c.newprog)
				q.As = AADD
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = int64(autosize) + c.ctxt.FixedFrameSize()
				q.Reg = REGSP
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_R5

				q = obj.Appendp(q, c.newprog)
				q.As = ACMP
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_R4
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_R5

				q = obj.Appendp(q, c.newprog)
				q.As = ABNE
				q.To.Type = obj.TYPE_BRANCH
				p2 = q

				q = obj.Appendp(q, c.newprog)
				q.As = AADD
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = c.ctxt.FixedFrameSize()
				q.Reg = REGSP
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_R6

				q = obj.Appendp(q, c.newprog)
				q.As = AMOVD
				q.From.Type = obj.TYPE_REG
				q.From.Reg = REG_R6
				q.To.Type = obj.TYPE_MEM
				q.To.Reg = REG_R3
				q.To.Offset = 0 

				q = obj.Appendp(q, c.newprog)

				q.As = obj.ANOP
				p1.To.SetTarget(q)
				p2.To.SetTarget(q)
			}

		case obj.ARET:
			if p.From.Type == obj.TYPE_CONST {
				c.ctxt.Diag("using BECOME (%v) is not supported!", p)
				break
			}

			retTarget := p.To.Sym

			if c.cursym.Func.Text.Mark&LEAF != 0 {
				if autosize == 0 || c.cursym.Name == "runtime.racecallbackthunk" {
					p.As = ABR
					p.From = obj.Addr{}
					if retTarget == nil {
						p.To.Type = obj.TYPE_REG
						p.To.Reg = REG_LR
					} else {
						p.To.Type = obj.TYPE_BRANCH
						p.To.Sym = retTarget
					}
					p.Mark |= BRANCH
					break
				}

				p.As = AADD
				p.From.Type = obj.TYPE_CONST
				p.From.Offset = int64(autosize)
				p.To.Type = obj.TYPE_REG
				p.To.Reg = REGSP
				p.Spadj = -autosize

				q = c.newprog()
				q.As = ABR
				q.Pos = p.Pos
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REG_LR
				q.Mark |= BRANCH
				q.Spadj = +autosize

				q.Link = p.Link
				p.Link = q
				break
			}

			p.As = AMOVD
			p.From.Type = obj.TYPE_MEM
			p.From.Offset = 0
			p.From.Reg = REGSP
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REGTMP

			q = c.newprog()
			q.As = AMOVD
			q.Pos = p.Pos
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REGTMP
			q.To.Type = obj.TYPE_REG
			q.To.Reg = REG_LR

			q.Link = p.Link
			p.Link = q
			p = q

			if false {
				
				q = c.newprog()

				q.As = AMOVD
				q.Pos = p.Pos
				q.From.Type = obj.TYPE_MEM
				q.From.Offset = 0
				q.From.Reg = REGTMP
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REGTMP

				q.Link = p.Link
				p.Link = q
				p = q
			}
			prev := p
			if autosize != 0 && c.cursym.Name != "runtime.racecallbackthunk" {
				q = c.newprog()
				q.As = AADD
				q.Pos = p.Pos
				q.From.Type = obj.TYPE_CONST
				q.From.Offset = int64(autosize)
				q.To.Type = obj.TYPE_REG
				q.To.Reg = REGSP
				q.Spadj = -autosize

				q.Link = p.Link
				prev.Link = q
				prev = q
			}

			q1 = c.newprog()
			q1.As = ABR
			q1.Pos = p.Pos
			if retTarget == nil {
				q1.To.Type = obj.TYPE_REG
				q1.To.Reg = REG_LR
			} else {
				q1.To.Type = obj.TYPE_BRANCH
				q1.To.Sym = retTarget
			}
			q1.Mark |= BRANCH
			q1.Spadj = +autosize

			q1.Link = q.Link
			prev.Link = q1
		case AADD:
			if p.To.Type == obj.TYPE_REG && p.To.Reg == REGSP && p.From.Type == obj.TYPE_CONST {
				p.Spadj = int32(-p.From.Offset)
			}
		case AMOVDU:
			if p.To.Type == obj.TYPE_MEM && p.To.Reg == REGSP {
				p.Spadj = int32(-p.To.Offset)
			}
			if p.From.Type == obj.TYPE_MEM && p.From.Reg == REGSP {
				p.Spadj = int32(-p.From.Offset)
			}
		case obj.AGETCALLERPC:
			if cursym.Leaf() {
				
				p.As = AMOVD
				p.From.Type = obj.TYPE_REG
				p.From.Reg = REG_LR
			} else {
				
				p.As = AMOVD
				p.From.Type = obj.TYPE_MEM
				p.From.Reg = REGSP
			}
		}
	}
}


func (c *ctxt9) stacksplit(p *obj.Prog, framesize int32) *obj.Prog {
	p0 := p 

	
	p = obj.Appendp(p, c.newprog)

	p.As = AMOVD
	p.From.Type = obj.TYPE_MEM
	p.From.Reg = REGG
	p.From.Offset = 2 * int64(c.ctxt.Arch.PtrSize) 
	if c.cursym.CFunc() {
		p.From.Offset = 3 * int64(c.ctxt.Arch.PtrSize) 
	}
	p.To.Type = obj.TYPE_REG
	p.To.Reg = REG_R3

	
	
	
	
	p = c.ctxt.StartUnsafePoint(p, c.newprog)

	var q *obj.Prog
	if framesize <= objabi.StackSmall {
		
		
		p = obj.Appendp(p, c.newprog)

		p.As = ACMPU
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_R3
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REGSP
	} else if framesize <= objabi.StackBig {
		
		
		
		p = obj.Appendp(p, c.newprog)

		p.As = AADD
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = -(int64(framesize) - objabi.StackSmall)
		p.Reg = REGSP
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R4

		p = obj.Appendp(p, c.newprog)
		p.As = ACMPU
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_R3
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R4
	} else {
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		p = obj.Appendp(p, c.newprog)

		p.As = ACMP
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_R3
		p.To.Type = obj.TYPE_CONST
		p.To.Offset = objabi.StackPreempt

		p = obj.Appendp(p, c.newprog)
		q = p
		p.As = ABEQ
		p.To.Type = obj.TYPE_BRANCH

		p = obj.Appendp(p, c.newprog)
		p.As = AADD
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(objabi.StackGuard)
		p.Reg = REGSP
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R4

		p = obj.Appendp(p, c.newprog)
		p.As = ASUB
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_R3
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R4

		p = obj.Appendp(p, c.newprog)
		p.As = AMOVD
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = int64(framesize) + int64(objabi.StackGuard) - objabi.StackSmall
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REGTMP

		p = obj.Appendp(p, c.newprog)
		p.As = ACMPU
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REGTMP
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R4
	}

	
	p = obj.Appendp(p, c.newprog)
	q1 := p

	p.As = ABLT
	p.To.Type = obj.TYPE_BRANCH

	
	p = obj.Appendp(p, c.newprog)

	p.As = AMOVD
	p.From.Type = obj.TYPE_REG
	p.From.Reg = REG_LR
	p.To.Type = obj.TYPE_REG
	p.To.Reg = REG_R5
	if q != nil {
		q.To.SetTarget(p)
	}

	p = c.ctxt.EmitEntryStackMap(c.cursym, p, c.newprog)

	var morestacksym *obj.LSym
	if c.cursym.CFunc() {
		morestacksym = c.ctxt.Lookup("runtime.morestackc")
	} else if !c.cursym.Func.Text.From.Sym.NeedCtxt() {
		morestacksym = c.ctxt.Lookup("runtime.morestack_noctxt")
	} else {
		morestacksym = c.ctxt.Lookup("runtime.morestack")
	}

	if c.ctxt.Flag_shared {
		
		
		
		
		
		

		
		p = obj.Appendp(p, c.newprog)
		p.As = AMOVD
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_R2
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = REGSP
		p.To.Offset = 8
	}

	if c.ctxt.Flag_dynlink {
		
		
		
		
		
		
		
		
		
		
		
		
		

		
		p = obj.Appendp(p, c.newprog)
		p.As = AMOVD
		p.From.Type = obj.TYPE_MEM
		p.From.Sym = morestacksym
		p.From.Name = obj.NAME_GOTREF
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R12

		
		p = obj.Appendp(p, c.newprog)
		p.As = AMOVD
		p.From.Type = obj.TYPE_REG
		p.From.Reg = REG_R12
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_LR

		
		p = obj.Appendp(p, c.newprog)
		p.As = obj.ACALL
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_LR
	} else {
		
		p = obj.Appendp(p, c.newprog)

		p.As = ABL
		p.To.Type = obj.TYPE_BRANCH
		p.To.Sym = morestacksym
	}

	if c.ctxt.Flag_shared {
		
		p = obj.Appendp(p, c.newprog)
		p.As = AMOVD
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = REGSP
		p.From.Offset = 8
		p.To.Type = obj.TYPE_REG
		p.To.Reg = REG_R2
	}

	p = c.ctxt.EndUnsafePoint(p, c.newprog, -1)

	
	p = obj.Appendp(p, c.newprog)
	p.As = ABR
	p.To.Type = obj.TYPE_BRANCH
	p.To.SetTarget(p0.Link)

	
	p = obj.Appendp(p, c.newprog)

	p.As = obj.ANOP 
	q1.To.SetTarget(p)

	return p
}

var Linkppc64 = obj.LinkArch{
	Arch:           sys.ArchPPC64,
	Init:           buildop,
	Preprocess:     preprocess,
	Assemble:       span9,
	Progedit:       progedit,
	DWARFRegisters: PPC64DWARFRegisters,
}

var Linkppc64le = obj.LinkArch{
	Arch:           sys.ArchPPC64LE,
	Init:           buildop,
	Preprocess:     preprocess,
	Assemble:       span9,
	Progedit:       progedit,
	DWARFRegisters: PPC64DWARFRegisters,
}
