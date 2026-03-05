





























package objabi


type SymKind uint8





//go:generate stringer -type=SymKind
const (
	
	Sxxx SymKind = iota
	
	STEXT
	
	SRODATA
	
	SNOPTRDATA
	
	SDATA
	
	SBSS
	
	SNOPTRBSS
	
	STLSBSS
	
	SDWARFCUINFO
	SDWARFCONST
	SDWARFFCN
	SDWARFABSFCN
	SDWARFTYPE
	SDWARFVAR
	SDWARFRANGE
	SDWARFLOC
	SDWARFLINES
	
	
	
	
	
	
	SABIALIAS
	
	SLIBFUZZER_EXTRA_COUNTER
	

)
