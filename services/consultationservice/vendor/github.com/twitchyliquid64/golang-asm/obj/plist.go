



package obj

import (
	"github.com/twitchyliquid64/golang-asm/objabi"
	"fmt"
	"strings"
)

type Plist struct {
	Firstpc *Prog
	Curfn   interface{} 
}



type ProgAlloc func() *Prog

func Flushplist(ctxt *Link, plist *Plist, newprog ProgAlloc, myimportpath string) {
	
	var curtext *LSym
	var etext *Prog
	var text []*LSym

	var plink *Prog
	for p := plist.Firstpc; p != nil; p = plink {
		if ctxt.Debugasm > 0 && ctxt.Debugvlog {
			fmt.Printf("obj: %v\n", p)
		}
		plink = p.Link
		p.Link = nil

		switch p.As {
		case AEND:
			continue

		case ATEXT:
			s := p.From.Sym
			if s == nil {
				
				curtext = nil
				continue
			}
			text = append(text, s)
			etext = p
			curtext = s
			continue

		case AFUNCDATA:
			
			if curtext == nil { 
				continue
			}
			if p.To.Sym.Name == "go_args_stackmap" {
				if p.From.Type != TYPE_CONST || p.From.Offset != objabi.FUNCDATA_ArgsPointerMaps {
					ctxt.Diag("FUNCDATA use of go_args_stackmap(SB) without FUNCDATA_ArgsPointerMaps")
				}
				p.To.Sym = ctxt.LookupDerived(curtext, curtext.Name+".args_stackmap")
			}

		}

		if curtext == nil {
			etext = nil
			continue
		}
		etext.Link = p
		etext = p
	}

	if newprog == nil {
		newprog = ctxt.NewProg
	}

	
	for _, s := range text {
		if !strings.HasPrefix(s.Name, "\"\".") {
			continue
		}
		found := false
		for p := s.Func.Text; p != nil; p = p.Link {
			if p.As == AFUNCDATA && p.From.Type == TYPE_CONST && p.From.Offset == objabi.FUNCDATA_ArgsPointerMaps {
				found = true
				break
			}
		}

		if !found {
			p := Appendp(s.Func.Text, newprog)
			p.As = AFUNCDATA
			p.From.Type = TYPE_CONST
			p.From.Offset = objabi.FUNCDATA_ArgsPointerMaps
			p.To.Type = TYPE_MEM
			p.To.Name = NAME_EXTERN
			p.To.Sym = ctxt.LookupDerived(s, s.Name+".args_stackmap")
		}
	}

	
	for _, s := range text {
		mkfwd(s)
		linkpatch(ctxt, s, newprog)
		ctxt.Arch.Preprocess(ctxt, s, newprog)
		ctxt.Arch.Assemble(ctxt, s, newprog)
		if ctxt.Errors > 0 {
			continue
		}
		linkpcln(ctxt, s)
		if myimportpath != "" {
			ctxt.populateDWARF(plist.Curfn, s, myimportpath)
		}
	}
}

func (ctxt *Link) InitTextSym(s *LSym, flag int) {
	if s == nil {
		
		return
	}
	if s.Func != nil {
		ctxt.Diag("InitTextSym double init for %s", s.Name)
	}
	s.Func = new(FuncInfo)
	if s.OnList() {
		ctxt.Diag("symbol %s listed multiple times", s.Name)
	}
	name := strings.Replace(s.Name, "\"\"", ctxt.Pkgpath, -1)
	s.Func.FuncID = objabi.GetFuncID(name, flag&WRAPPER != 0)
	s.Set(AttrOnList, true)
	s.Set(AttrDuplicateOK, flag&DUPOK != 0)
	s.Set(AttrNoSplit, flag&NOSPLIT != 0)
	s.Set(AttrReflectMethod, flag&REFLECTMETHOD != 0)
	s.Set(AttrWrapper, flag&WRAPPER != 0)
	s.Set(AttrNeedCtxt, flag&NEEDCTXT != 0)
	s.Set(AttrNoFrame, flag&NOFRAME != 0)
	s.Set(AttrTopFrame, flag&TOPFRAME != 0)
	s.Type = objabi.STEXT
	ctxt.Text = append(ctxt.Text, s)

	
	ctxt.dwarfSym(s)
}

func (ctxt *Link) Globl(s *LSym, size int64, flag int) {
	if s.OnList() {
		ctxt.Diag("symbol %s listed multiple times", s.Name)
	}
	s.Set(AttrOnList, true)
	ctxt.Data = append(ctxt.Data, s)
	s.Size = size
	if s.Type == 0 {
		s.Type = objabi.SBSS
	}
	if flag&DUPOK != 0 {
		s.Set(AttrDuplicateOK, true)
	}
	if flag&RODATA != 0 {
		s.Type = objabi.SRODATA
	} else if flag&NOPTR != 0 {
		if s.Type == objabi.SDATA {
			s.Type = objabi.SNOPTRDATA
		} else {
			s.Type = objabi.SNOPTRBSS
		}
	} else if flag&TLSBSS != 0 {
		s.Type = objabi.STLSBSS
	}
	if strings.HasPrefix(s.Name, "\"\"."+StaticNamePref) {
		s.Set(AttrStatic, true)
	}
}




func (ctxt *Link) EmitEntryLiveness(s *LSym, p *Prog, newprog ProgAlloc) *Prog {
	pcdata := ctxt.EmitEntryStackMap(s, p, newprog)
	pcdata = ctxt.EmitEntryRegMap(s, pcdata, newprog)
	return pcdata
}


func (ctxt *Link) EmitEntryStackMap(s *LSym, p *Prog, newprog ProgAlloc) *Prog {
	pcdata := Appendp(p, newprog)
	pcdata.Pos = s.Func.Text.Pos
	pcdata.As = APCDATA
	pcdata.From.Type = TYPE_CONST
	pcdata.From.Offset = objabi.PCDATA_StackMapIndex
	pcdata.To.Type = TYPE_CONST
	pcdata.To.Offset = -1 

	return pcdata
}


func (ctxt *Link) EmitEntryRegMap(s *LSym, p *Prog, newprog ProgAlloc) *Prog {
	pcdata := Appendp(p, newprog)
	pcdata.Pos = s.Func.Text.Pos
	pcdata.As = APCDATA
	pcdata.From.Type = TYPE_CONST
	pcdata.From.Offset = objabi.PCDATA_RegMapIndex
	pcdata.To.Type = TYPE_CONST
	pcdata.To.Offset = -1

	return pcdata
}





func (ctxt *Link) StartUnsafePoint(p *Prog, newprog ProgAlloc) *Prog {
	pcdata := Appendp(p, newprog)
	pcdata.As = APCDATA
	pcdata.From.Type = TYPE_CONST
	pcdata.From.Offset = objabi.PCDATA_RegMapIndex
	pcdata.To.Type = TYPE_CONST
	pcdata.To.Offset = objabi.PCDATA_RegMapUnsafe

	return pcdata
}





func (ctxt *Link) EndUnsafePoint(p *Prog, newprog ProgAlloc, oldval int64) *Prog {
	pcdata := Appendp(p, newprog)
	pcdata.As = APCDATA
	pcdata.From.Type = TYPE_CONST
	pcdata.From.Offset = objabi.PCDATA_RegMapIndex
	pcdata.To.Type = TYPE_CONST
	pcdata.To.Offset = oldval

	return pcdata
}











func MarkUnsafePoints(ctxt *Link, p0 *Prog, newprog ProgAlloc, isUnsafePoint, isRestartable func(*Prog) bool) {
	if isRestartable == nil {
		
		isRestartable = func(*Prog) bool { return false }
	}
	prev := p0
	prevPcdata := int64(-1) 
	prevRestart := int64(0)
	for p := prev.Link; p != nil; p, prev = p.Link, p {
		if p.As == APCDATA && p.From.Offset == objabi.PCDATA_RegMapIndex {
			prevPcdata = p.To.Offset
			continue
		}
		if prevPcdata == objabi.PCDATA_RegMapUnsafe {
			continue 
		}
		if isUnsafePoint(p) {
			q := ctxt.StartUnsafePoint(prev, newprog)
			q.Pc = p.Pc
			q.Link = p
			
			for p.Link != nil && isUnsafePoint(p.Link) {
				p = p.Link
			}
			if p.Link == nil {
				break 
			}
			p = ctxt.EndUnsafePoint(p, newprog, prevPcdata)
			p.Pc = p.Link.Pc
			continue
		}
		if isRestartable(p) {
			val := int64(objabi.PCDATA_Restart1)
			if val == prevRestart {
				val = objabi.PCDATA_Restart2
			}
			prevRestart = val
			q := Appendp(prev, newprog)
			q.As = APCDATA
			q.From.Type = TYPE_CONST
			q.From.Offset = objabi.PCDATA_RegMapIndex
			q.To.Type = TYPE_CONST
			q.To.Offset = val
			q.Pc = p.Pc
			q.Link = p

			if p.Link == nil {
				break 
			}
			if isRestartable(p.Link) {
				
				
				continue
			}
			p = Appendp(p, newprog)
			p.As = APCDATA
			p.From.Type = TYPE_CONST
			p.From.Offset = objabi.PCDATA_RegMapIndex
			p.To.Type = TYPE_CONST
			p.To.Offset = prevPcdata
			p.Pc = p.Link.Pc
		}
	}
}
