





























package obj

import (
	"bufio"
	"github.com/twitchyliquid64/golang-asm/dwarf"
	"github.com/twitchyliquid64/golang-asm/goobj"
	"github.com/twitchyliquid64/golang-asm/objabi"
	"github.com/twitchyliquid64/golang-asm/src"
	"github.com/twitchyliquid64/golang-asm/sys"
	"fmt"
	"sync"
)















































































































































type Addr struct {
	Reg    int16
	Index  int16
	Scale  int16 
	Type   AddrType
	Name   AddrName
	Class  int8
	Offset int64
	Sym    *LSym

	
	
	
	
	
	Val interface{}
}

type AddrName int8

const (
	NAME_NONE AddrName = iota
	NAME_EXTERN
	NAME_STATIC
	NAME_AUTO
	NAME_PARAM
	
	
	NAME_GOTREF
	
	NAME_TOCREF
)

//go:generate stringer -type AddrType

type AddrType uint8

const (
	TYPE_NONE AddrType = iota
	TYPE_BRANCH
	TYPE_TEXTSIZE
	TYPE_MEM
	TYPE_CONST
	TYPE_FCONST
	TYPE_SCONST
	TYPE_REG
	TYPE_ADDR
	TYPE_SHIFT
	TYPE_REGREG
	TYPE_REGREG2
	TYPE_INDIR
	TYPE_REGLIST
)

func (a *Addr) Target() *Prog {
	if a.Type == TYPE_BRANCH && a.Val != nil {
		return a.Val.(*Prog)
	}
	return nil
}
func (a *Addr) SetTarget(t *Prog) {
	if a.Type != TYPE_BRANCH {
		panic("setting branch target when type is not TYPE_BRANCH")
	}
	a.Val = t
}
































type Prog struct {
	Ctxt     *Link    
	Link     *Prog    
	From     Addr     
	RestArgs []Addr   
	To       Addr     
	Pool     *Prog    
	Forwd    *Prog    
	Rel      *Prog    
	Pc       int64    
	Pos      src.XPos 
	Spadj    int32    
	As       As       
	Reg      int16    
	RegTo2   int16    
	Mark     uint16   
	Optab    uint16   
	Scond    uint8    
	Back     uint8    
	Ft       uint8    
	Tt       uint8    
	Isize    uint8    
}





func (p *Prog) From3Type() AddrType {
	if p.RestArgs == nil {
		return TYPE_NONE
	}
	return p.RestArgs[0].Type
}










func (p *Prog) GetFrom3() *Addr {
	if p.RestArgs == nil {
		return nil
	}
	return &p.RestArgs[0]
}





func (p *Prog) SetFrom3(a Addr) {
	p.RestArgs = []Addr{a}
}






type As int16


const (
	AXXX As = iota
	ACALL
	ADUFFCOPY
	ADUFFZERO
	AEND
	AFUNCDATA
	AJMP
	ANOP
	APCALIGN
	APCDATA
	ARET
	AGETCALLERPC
	ATEXT
	AUNDEF
	A_ARCHSPECIFIC
)








const (
	ABase386 = (1 + iota) << 11
	ABaseARM
	ABaseAMD64
	ABasePPC64
	ABaseARM64
	ABaseMIPS
	ABaseRISCV
	ABaseS390X
	ABaseWasm

	AllowedOpCodes = 1 << 11            
	AMask          = AllowedOpCodes - 1 
)



type LSym struct {
	Name string
	Type objabi.SymKind
	Attribute

	RefIdx int 
	Size   int64
	Gotype *LSym
	P      []byte
	R      []Reloc

	Func *FuncInfo

	Pkg    string
	PkgIdx int32
	SymIdx int32 
}


type FuncInfo struct {
	Args     int32
	Locals   int32
	Align    int32
	FuncID   objabi.FuncID
	Text     *Prog
	Autot    map[*LSym]struct{}
	Pcln     Pcln
	InlMarks []InlMark

	dwarfInfoSym       *LSym
	dwarfLocSym        *LSym
	dwarfRangesSym     *LSym
	dwarfAbsFnSym      *LSym
	dwarfDebugLinesSym *LSym

	GCArgs             *LSym
	GCLocals           *LSym
	GCRegs             *LSym 
	StackObjects       *LSym
	OpenCodedDeferInfo *LSym

	FuncInfoSym *LSym
}

type InlMark struct {
	
	
	
	
	p  *Prog
	id int32
}





func (fi *FuncInfo) AddInlMark(p *Prog, id int32) {
	fi.InlMarks = append(fi.InlMarks, InlMark{p: p, id: id})
}



func (fi *FuncInfo) RecordAutoType(gotype *LSym) {
	if fi.Autot == nil {
		fi.Autot = make(map[*LSym]struct{})
	}
	fi.Autot[gotype] = struct{}{}
}

//go:generate stringer -type ABI


type ABI uint8

const (
	
	
	
	
	
	ABI0 ABI = iota

	
	
	
	
	ABIInternal

	ABICount
)


type Attribute uint32

const (
	AttrDuplicateOK Attribute = 1 << iota
	AttrCFunc
	AttrNoSplit
	AttrLeaf
	AttrWrapper
	AttrNeedCtxt
	AttrNoFrame
	AttrOnList
	AttrStatic

	
	AttrMakeTypelink

	
	
	
	
	
	
	AttrReflectMethod

	
	
	
	
	
	
	AttrLocal

	
	
	AttrWasInlined

	
	
	AttrTopFrame

	
	
	AttrIndexed

	
	
	
	
	AttrUsedInIface

	
	AttrContentAddressable

	
	
	
	
	
	attrABIBase
)

func (a Attribute) DuplicateOK() bool        { return a&AttrDuplicateOK != 0 }
func (a Attribute) MakeTypelink() bool       { return a&AttrMakeTypelink != 0 }
func (a Attribute) CFunc() bool              { return a&AttrCFunc != 0 }
func (a Attribute) NoSplit() bool            { return a&AttrNoSplit != 0 }
func (a Attribute) Leaf() bool               { return a&AttrLeaf != 0 }
func (a Attribute) OnList() bool             { return a&AttrOnList != 0 }
func (a Attribute) ReflectMethod() bool      { return a&AttrReflectMethod != 0 }
func (a Attribute) Local() bool              { return a&AttrLocal != 0 }
func (a Attribute) Wrapper() bool            { return a&AttrWrapper != 0 }
func (a Attribute) NeedCtxt() bool           { return a&AttrNeedCtxt != 0 }
func (a Attribute) NoFrame() bool            { return a&AttrNoFrame != 0 }
func (a Attribute) Static() bool             { return a&AttrStatic != 0 }
func (a Attribute) WasInlined() bool         { return a&AttrWasInlined != 0 }
func (a Attribute) TopFrame() bool           { return a&AttrTopFrame != 0 }
func (a Attribute) Indexed() bool            { return a&AttrIndexed != 0 }
func (a Attribute) UsedInIface() bool        { return a&AttrUsedInIface != 0 }
func (a Attribute) ContentAddressable() bool { return a&AttrContentAddressable != 0 }

func (a *Attribute) Set(flag Attribute, value bool) {
	if value {
		*a |= flag
	} else {
		*a &^= flag
	}
}

func (a Attribute) ABI() ABI { return ABI(a / attrABIBase) }
func (a *Attribute) SetABI(abi ABI) {
	const mask = 1 
	*a = (*a &^ (mask * attrABIBase)) | Attribute(abi)*attrABIBase
}

var textAttrStrings = [...]struct {
	bit Attribute
	s   string
}{
	{bit: AttrDuplicateOK, s: "DUPOK"},
	{bit: AttrMakeTypelink, s: ""},
	{bit: AttrCFunc, s: "CFUNC"},
	{bit: AttrNoSplit, s: "NOSPLIT"},
	{bit: AttrLeaf, s: "LEAF"},
	{bit: AttrOnList, s: ""},
	{bit: AttrReflectMethod, s: "REFLECTMETHOD"},
	{bit: AttrLocal, s: "LOCAL"},
	{bit: AttrWrapper, s: "WRAPPER"},
	{bit: AttrNeedCtxt, s: "NEEDCTXT"},
	{bit: AttrNoFrame, s: "NOFRAME"},
	{bit: AttrStatic, s: "STATIC"},
	{bit: AttrWasInlined, s: ""},
	{bit: AttrTopFrame, s: "TOPFRAME"},
	{bit: AttrIndexed, s: ""},
	{bit: AttrContentAddressable, s: ""},
}


func (a Attribute) TextAttrString() string {
	var s string
	for _, x := range textAttrStrings {
		if a&x.bit != 0 {
			if x.s != "" {
				s += x.s + "|"
			}
			a &^= x.bit
		}
	}
	switch a.ABI() {
	case ABI0:
	case ABIInternal:
		s += "ABIInternal|"
		a.SetABI(0) 
	}
	if a != 0 {
		s += fmt.Sprintf("UnknownAttribute(%d)|", a)
	}
	
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

func (s *LSym) String() string {
	return s.Name
}


func (s *LSym) CanBeAnSSASym() {
}

type Pcln struct {
	Pcsp        Pcdata
	Pcfile      Pcdata
	Pcline      Pcdata
	Pcinline    Pcdata
	Pcdata      []Pcdata
	Funcdata    []*LSym
	Funcdataoff []int64
	UsedFiles   map[goobj.CUFileIndex]struct{} 
	InlTree     InlTree                        
}

type Reloc struct {
	Off  int32
	Siz  uint8
	Type objabi.RelocType
	Add  int64
	Sym  *LSym
}

type Auto struct {
	Asym    *LSym
	Aoffset int32
	Name    AddrName
	Gotype  *LSym
}

type Pcdata struct {
	P []byte
}



type Link struct {
	Headtype           objabi.HeadType
	Arch               *LinkArch
	Debugasm           int
	Debugvlog          bool
	Debugpcln          string
	Flag_shared        bool
	Flag_dynlink       bool
	Flag_linkshared    bool
	Flag_optimize      bool
	Flag_locationlists bool
	Retpoline          bool 
	Bso                *bufio.Writer
	Pathname           string
	Pkgpath            string           
	hashmu             sync.Mutex       
	hash               map[string]*LSym 
	funchash           map[string]*LSym 
	statichash         map[string]*LSym 
	PosTable           src.PosTable
	InlTree            InlTree 
	DwFixups           *DwarfFixupTable
	Imports            []goobj.ImportedPkg
	DiagFunc           func(string, ...interface{})
	DiagFlush          func()
	DebugInfo          func(fn *LSym, info *LSym, curfn interface{}) ([]dwarf.Scope, dwarf.InlCalls) 
	GenAbstractFunc    func(fn *LSym)
	Errors             int

	InParallel    bool 
	UseBASEntries bool 
	IsAsm         bool 

	
	Text []*LSym
	Data []*LSym

	
	
	
	
	
	
	
	
	ABIAliases []*LSym

	
	
	
	
	constSyms []*LSym

	
	
	pkgIdx map[string]int32

	defs         []*LSym 
	hashed64defs []*LSym 
	hasheddefs   []*LSym 
	nonpkgdefs   []*LSym 
	nonpkgrefs   []*LSym 

	Fingerprint goobj.FingerprintType 
}

func (ctxt *Link) Diag(format string, args ...interface{}) {
	ctxt.Errors++
	ctxt.DiagFunc(format, args...)
}

func (ctxt *Link) Logf(format string, args ...interface{}) {
	fmt.Fprintf(ctxt.Bso, format, args...)
	ctxt.Bso.Flush()
}





func (ctxt *Link) FixedFrameSize() int64 {
	switch ctxt.Arch.Family {
	case sys.AMD64, sys.I386, sys.Wasm:
		return 0
	case sys.PPC64:
		
		
		return int64(4 * ctxt.Arch.PtrSize)
	default:
		return int64(ctxt.Arch.PtrSize)
	}
}


type LinkArch struct {
	*sys.Arch
	Init           func(*Link)
	Preprocess     func(*Link, *LSym, ProgAlloc)
	Assemble       func(*Link, *LSym, ProgAlloc)
	Progedit       func(*Link, *Prog, ProgAlloc)
	UnaryDst       map[As]bool 
	DWARFRegisters map[int16]int16
}
