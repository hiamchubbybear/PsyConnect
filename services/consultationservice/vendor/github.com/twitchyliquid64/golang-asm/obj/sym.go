






























package obj

import (
	"github.com/twitchyliquid64/golang-asm/goobj"
	"github.com/twitchyliquid64/golang-asm/objabi"
	"fmt"
	"log"
	"math"
	"sort"
)

func Linknew(arch *LinkArch) *Link {
	ctxt := new(Link)
	ctxt.hash = make(map[string]*LSym)
	ctxt.funchash = make(map[string]*LSym)
	ctxt.statichash = make(map[string]*LSym)
	ctxt.Arch = arch
	ctxt.Pathname = objabi.WorkingDir()

	if err := ctxt.Headtype.Set(objabi.GOOS); err != nil {
		log.Fatalf("unknown goos %s", objabi.GOOS)
	}

	ctxt.Flag_optimize = true
	return ctxt
}



func (ctxt *Link) LookupDerived(s *LSym, name string) *LSym {
	if s.Static() {
		return ctxt.LookupStatic(name)
	}
	return ctxt.Lookup(name)
}



func (ctxt *Link) LookupStatic(name string) *LSym {
	s := ctxt.statichash[name]
	if s == nil {
		s = &LSym{Name: name, Attribute: AttrStatic}
		ctxt.statichash[name] = s
	}
	return s
}



func (ctxt *Link) LookupABI(name string, abi ABI) *LSym {
	return ctxt.LookupABIInit(name, abi, nil)
}




func (ctxt *Link) LookupABIInit(name string, abi ABI, init func(s *LSym)) *LSym {
	var hash map[string]*LSym
	switch abi {
	case ABI0:
		hash = ctxt.hash
	case ABIInternal:
		hash = ctxt.funchash
	default:
		panic("unknown ABI")
	}

	ctxt.hashmu.Lock()
	s := hash[name]
	if s == nil {
		s = &LSym{Name: name}
		s.SetABI(abi)
		hash[name] = s
		if init != nil {
			init(s)
		}
	}
	ctxt.hashmu.Unlock()
	return s
}



func (ctxt *Link) Lookup(name string) *LSym {
	return ctxt.LookupInit(name, nil)
}




func (ctxt *Link) LookupInit(name string, init func(s *LSym)) *LSym {
	ctxt.hashmu.Lock()
	s := ctxt.hash[name]
	if s == nil {
		s = &LSym{Name: name}
		ctxt.hash[name] = s
		if init != nil {
			init(s)
		}
	}
	ctxt.hashmu.Unlock()
	return s
}

func (ctxt *Link) Float32Sym(f float32) *LSym {
	i := math.Float32bits(f)
	name := fmt.Sprintf("$f32.%08x", i)
	return ctxt.LookupInit(name, func(s *LSym) {
		s.Size = 4
		s.WriteFloat32(ctxt, 0, f)
		s.Type = objabi.SRODATA
		s.Set(AttrLocal, true)
		s.Set(AttrContentAddressable, true)
		ctxt.constSyms = append(ctxt.constSyms, s)
	})
}

func (ctxt *Link) Float64Sym(f float64) *LSym {
	i := math.Float64bits(f)
	name := fmt.Sprintf("$f64.%016x", i)
	return ctxt.LookupInit(name, func(s *LSym) {
		s.Size = 8
		s.WriteFloat64(ctxt, 0, f)
		s.Type = objabi.SRODATA
		s.Set(AttrLocal, true)
		s.Set(AttrContentAddressable, true)
		ctxt.constSyms = append(ctxt.constSyms, s)
	})
}

func (ctxt *Link) Int64Sym(i int64) *LSym {
	name := fmt.Sprintf("$i64.%016x", uint64(i))
	return ctxt.LookupInit(name, func(s *LSym) {
		s.Size = 8
		s.WriteInt(ctxt, 0, 8, i)
		s.Type = objabi.SRODATA
		s.Set(AttrLocal, true)
		s.Set(AttrContentAddressable, true)
		ctxt.constSyms = append(ctxt.constSyms, s)
	})
}




func (ctxt *Link) NumberSyms() {
	if ctxt.Headtype == objabi.Haix {
		
		
		
		
		sort.Slice(ctxt.Data, func(i, j int) bool {
			return ctxt.Data[i].Name < ctxt.Data[j].Name
		})
	}

	
	
	sort.Slice(ctxt.constSyms, func(i, j int) bool {
		return ctxt.constSyms[i].Name < ctxt.constSyms[j].Name
	})
	ctxt.Data = append(ctxt.Data, ctxt.constSyms...)
	ctxt.constSyms = nil

	ctxt.pkgIdx = make(map[string]int32)
	ctxt.defs = []*LSym{}
	ctxt.hashed64defs = []*LSym{}
	ctxt.hasheddefs = []*LSym{}
	ctxt.nonpkgdefs = []*LSym{}

	var idx, hashedidx, hashed64idx, nonpkgidx int32
	ctxt.traverseSyms(traverseDefs, func(s *LSym) {
		
		
		if s.ContentAddressable() && (ctxt.Pkgpath != "" || len(s.R) == 0) {
			if len(s.P) <= 8 && len(s.R) == 0 { 
				s.PkgIdx = goobj.PkgIdxHashed64
				s.SymIdx = hashed64idx
				if hashed64idx != int32(len(ctxt.hashed64defs)) {
					panic("bad index")
				}
				ctxt.hashed64defs = append(ctxt.hashed64defs, s)
				hashed64idx++
			} else {
				s.PkgIdx = goobj.PkgIdxHashed
				s.SymIdx = hashedidx
				if hashedidx != int32(len(ctxt.hasheddefs)) {
					panic("bad index")
				}
				ctxt.hasheddefs = append(ctxt.hasheddefs, s)
				hashedidx++
			}
		} else if isNonPkgSym(ctxt, s) {
			s.PkgIdx = goobj.PkgIdxNone
			s.SymIdx = nonpkgidx
			if nonpkgidx != int32(len(ctxt.nonpkgdefs)) {
				panic("bad index")
			}
			ctxt.nonpkgdefs = append(ctxt.nonpkgdefs, s)
			nonpkgidx++
		} else {
			s.PkgIdx = goobj.PkgIdxSelf
			s.SymIdx = idx
			if idx != int32(len(ctxt.defs)) {
				panic("bad index")
			}
			ctxt.defs = append(ctxt.defs, s)
			idx++
		}
		s.Set(AttrIndexed, true)
	})

	ipkg := int32(1) 
	nonpkgdef := nonpkgidx
	ctxt.traverseSyms(traverseRefs|traverseAux, func(rs *LSym) {
		if rs.PkgIdx != goobj.PkgIdxInvalid {
			return
		}
		if !ctxt.Flag_linkshared {
			
			
			
			if i := goobj.BuiltinIdx(rs.Name, int(rs.ABI())); i != -1 {
				rs.PkgIdx = goobj.PkgIdxBuiltin
				rs.SymIdx = int32(i)
				rs.Set(AttrIndexed, true)
				return
			}
		}
		pkg := rs.Pkg
		if rs.ContentAddressable() {
			
			panic("hashed refs unsupported for now")
		}
		if pkg == "" || pkg == "\"\"" || pkg == "_" || !rs.Indexed() {
			rs.PkgIdx = goobj.PkgIdxNone
			rs.SymIdx = nonpkgidx
			rs.Set(AttrIndexed, true)
			if nonpkgidx != nonpkgdef+int32(len(ctxt.nonpkgrefs)) {
				panic("bad index")
			}
			ctxt.nonpkgrefs = append(ctxt.nonpkgrefs, rs)
			nonpkgidx++
			return
		}
		if k, ok := ctxt.pkgIdx[pkg]; ok {
			rs.PkgIdx = k
			return
		}
		rs.PkgIdx = ipkg
		ctxt.pkgIdx[pkg] = ipkg
		ipkg++
	})
}



func isNonPkgSym(ctxt *Link, s *LSym) bool {
	if ctxt.IsAsm && !s.Static() {
		
		
		return true
	}
	if ctxt.Flag_linkshared {
		
		
		return true
	}
	if s.Pkg == "_" {
		
		
		return true
	}
	if s.DuplicateOK() {
		
		return true
	}
	return false
}




const StaticNamePref = ".stmp_"

type traverseFlag uint32

const (
	traverseDefs traverseFlag = 1 << iota
	traverseRefs
	traverseAux

	traverseAll = traverseDefs | traverseRefs | traverseAux
)


func (ctxt *Link) traverseSyms(flag traverseFlag, fn func(*LSym)) {
	lists := [][]*LSym{ctxt.Text, ctxt.Data, ctxt.ABIAliases}
	for _, list := range lists {
		for _, s := range list {
			if flag&traverseDefs != 0 {
				fn(s)
			}
			if flag&traverseRefs != 0 {
				for _, r := range s.R {
					if r.Sym != nil {
						fn(r.Sym)
					}
				}
			}
			if flag&traverseAux != 0 {
				if s.Gotype != nil {
					fn(s.Gotype)
				}
				if s.Type == objabi.STEXT {
					f := func(parent *LSym, aux *LSym) {
						fn(aux)
					}
					ctxt.traverseFuncAux(flag, s, f)
				}
			}
		}
	}
}

func (ctxt *Link) traverseFuncAux(flag traverseFlag, fsym *LSym, fn func(parent *LSym, aux *LSym)) {
	pc := &fsym.Func.Pcln
	if flag&traverseAux == 0 {
		
		
		panic("should not be here")
	}
	for _, d := range pc.Funcdata {
		if d != nil {
			fn(fsym, d)
		}
	}
	files := ctxt.PosTable.FileTable()
	usedFiles := make([]goobj.CUFileIndex, 0, len(pc.UsedFiles))
	for f := range pc.UsedFiles {
		usedFiles = append(usedFiles, f)
	}
	sort.Slice(usedFiles, func(i, j int) bool { return usedFiles[i] < usedFiles[j] })
	for _, f := range usedFiles {
		if filesym := ctxt.Lookup(files[f]); filesym != nil {
			fn(fsym, filesym)
		}
	}
	for _, call := range pc.InlTree.nodes {
		if call.Func != nil {
			fn(fsym, call.Func)
		}
		f, _ := linkgetlineFromPos(ctxt, call.Pos)
		if filesym := ctxt.Lookup(f); filesym != nil {
			fn(fsym, filesym)
		}
	}
	dwsyms := []*LSym{fsym.Func.dwarfRangesSym, fsym.Func.dwarfLocSym, fsym.Func.dwarfDebugLinesSym, fsym.Func.dwarfInfoSym}
	for _, dws := range dwsyms {
		if dws == nil || dws.Size == 0 {
			continue
		}
		fn(fsym, dws)
		if flag&traverseRefs != 0 {
			for _, r := range dws.R {
				if r.Sym != nil {
					fn(dws, r.Sym)
				}
			}
		}
	}
}


func (ctxt *Link) traverseAuxSyms(flag traverseFlag, fn func(parent *LSym, aux *LSym)) {
	lists := [][]*LSym{ctxt.Text, ctxt.Data, ctxt.ABIAliases}
	for _, list := range lists {
		for _, s := range list {
			if s.Gotype != nil {
				if flag&traverseDefs != 0 {
					fn(s, s.Gotype)
				}
			}
			if s.Type != objabi.STEXT {
				continue
			}
			ctxt.traverseFuncAux(flag, s, fn)
		}
	}
}
