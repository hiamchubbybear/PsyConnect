





























package obj





func brloop(p *Prog) *Prog {
	c := 0
	for q := p; q != nil; q = q.To.Target() {
		if q.As != AJMP || q.To.Target() == nil {
			return q
		}
		c++
		if c >= 5000 {
			
			return nil
		}
	}
	panic("unreachable")
}


func checkaddr(ctxt *Link, p *Prog, a *Addr) {
	switch a.Type {
	case TYPE_NONE, TYPE_REGREG2, TYPE_REGLIST:
		return

	case TYPE_BRANCH, TYPE_TEXTSIZE:
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name != 0 {
			break
		}
		return

	case TYPE_MEM:
		return

	case TYPE_CONST:
		
		if a.Name != 0 || a.Sym != nil || a.Reg != 0 {
			ctxt.Diag("argument is TYPE_CONST, should be TYPE_ADDR, in %v", p)
			return
		}

		if a.Reg != 0 || a.Scale != 0 || a.Name != 0 || a.Sym != nil || a.Val != nil {
			break
		}
		return

	case TYPE_FCONST, TYPE_SCONST:
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name != 0 || a.Offset != 0 || a.Sym != nil {
			break
		}
		return

	case TYPE_REG:
		
		
		if a.Scale != 0 || a.Name != 0 || a.Sym != nil {
			break
		}
		return

	case TYPE_ADDR:
		if a.Val != nil {
			break
		}
		if a.Reg == 0 && a.Index == 0 && a.Scale == 0 && a.Name == 0 && a.Sym == nil {
			ctxt.Diag("argument is TYPE_ADDR, should be TYPE_CONST, in %v", p)
		}
		return

	case TYPE_SHIFT, TYPE_REGREG:
		if a.Index != 0 || a.Scale != 0 || a.Name != 0 || a.Sym != nil || a.Val != nil {
			break
		}
		return

	case TYPE_INDIR:
		
		
		if a.Reg != 0 || a.Index != 0 || a.Scale != 0 || a.Name == 0 || a.Offset != 0 || a.Sym == nil || a.Val != nil {
			break
		}
		return
	}

	ctxt.Diag("invalid encoding for argument %v", p)
}

func linkpatch(ctxt *Link, sym *LSym, newprog ProgAlloc) {
	for p := sym.Func.Text; p != nil; p = p.Link {
		checkaddr(ctxt, p, &p.From)
		if p.GetFrom3() != nil {
			checkaddr(ctxt, p, p.GetFrom3())
		}
		checkaddr(ctxt, p, &p.To)

		if ctxt.Arch.Progedit != nil {
			ctxt.Arch.Progedit(ctxt, p, newprog)
		}
		if p.To.Type != TYPE_BRANCH {
			continue
		}
		if p.To.Val != nil {
			continue
		}

		if p.To.Sym != nil {
			continue
		}
		q := sym.Func.Text
		for q != nil && p.To.Offset != q.Pc {
			if q.Forwd != nil && p.To.Offset >= q.Forwd.Pc {
				q = q.Forwd
			} else {
				q = q.Link
			}
		}

		if q == nil {
			name := "<nil>"
			if p.To.Sym != nil {
				name = p.To.Sym.Name
			}
			ctxt.Diag("branch out of range (%#x)\n%v [%s]", uint32(p.To.Offset), p, name)
			p.To.Type = TYPE_NONE
		}

		p.To.SetTarget(q)
	}

	if !ctxt.Flag_optimize {
		return
	}

	
	for p := sym.Func.Text; p != nil; p = p.Link {
		if p.To.Target() == nil {
			continue
		}
		p.To.SetTarget(brloop(p.To.Target()))
		if p.To.Target() != nil && p.To.Type == TYPE_BRANCH {
			p.To.Offset = p.To.Target().Pc
		}
	}
}
