









package cpuid

import (
	"flag"
	"fmt"
	"math"
	"math/bits"
	"os"
	"runtime"
	"strings"
)





type Vendor int

const (
	VendorUnknown Vendor = iota
	Intel
	AMD
	VIA
	Transmeta
	NSC
	KVM  
	MSVM 
	VMware
	XenHVM
	Bhyve
	Hygon
	SiS
	RDC

	Ampere
	ARM
	Broadcom
	Cavium
	DEC
	Fujitsu
	Infineon
	Motorola
	NVIDIA
	AMCC
	Qualcomm
	Marvell

	QEMU
	QNX
	ACRN
	SRE
	Apple

	lastVendor
)

//go:generate stringer -type=FeatureID,Vendor


type FeatureID int

const (
	
	UNKNOWN = -1

	
	ADX                 FeatureID = iota 
	AESNI                                
	AMD3DNOW                             
	AMD3DNOWEXT                          
	AMXBF16                              
	AMXFP16                              
	AMXINT8                              
	AMXFP8                               
	AMXTILE                              
	AMXTF32                              
	AMXCOMPLEX                           
	APX_F                                
	AVX                                  
	AVX10                                
	AVX10_128                            
	AVX10_256                            
	AVX10_512                            
	AVX2                                 
	AVX512BF16                           
	AVX512BITALG                         
	AVX512BW                             
	AVX512CD                             
	AVX512DQ                             
	AVX512ER                             
	AVX512F                              
	AVX512FP16                           
	AVX512IFMA                           
	AVX512PF                             
	AVX512VBMI                           
	AVX512VBMI2                          
	AVX512VL                             
	AVX512VNNI                           
	AVX512VP2INTERSECT                   
	AVX512VPOPCNTDQ                      
	AVXIFMA                              
	AVXNECONVERT                         
	AVXSLOW                              
	AVXVNNI                              
	AVXVNNIINT8                          
	AVXVNNIINT16                         
	BHI_CTRL                             
	BMI1                                 
	BMI2                                 
	CETIBT                               
	CETSS                                
	CLDEMOTE                             
	CLMUL                                
	CLZERO                               
	CMOV                                 
	CMPCCXADD                            
	CMPSB_SCADBS_SHORT                   
	CMPXCHG8                             
	CPBOOST                              
	CPPC                                 
	CX16                                 
	EFER_LMSLE_UNS                       
	ENQCMD                               
	ERMS                                 
	F16C                                 
	FLUSH_L1D                            
	FMA3                                 
	FMA4                                 
	FP128                                
	FP256                                
	FSRM                                 
	FXSR                                 
	FXSROPT                              
	GFNI                                 
	HLE                                  
	HRESET                               
	HTT                                  
	HWA                                  
	HYBRID_CPU                           
	HYPERVISOR                           
	IA32_ARCH_CAP                        
	IA32_CORE_CAP                        
	IBPB                                 
	IBPB_BRTYPE                          
	IBRS                                 
	IBRS_PREFERRED                       
	IBRS_PROVIDES_SMP                    
	IBS                                  
	IBSBRNTRGT                           
	IBSFETCHSAM                          
	IBSFFV                               
	IBSOPCNT                             
	IBSOPCNTEXT                          
	IBSOPSAM                             
	IBSRDWROPCNT                         
	IBSRIPINVALIDCHK                     
	IBS_FETCH_CTLX                       
	IBS_OPDATA4                          
	IBS_OPFUSE                           
	IBS_PREVENTHOST                      
	IBS_ZEN4                             
	IDPRED_CTRL                          
	INT_WBINVD                           
	INVLPGB                              
	KEYLOCKER                            
	KEYLOCKERW                           
	LAHF                                 
	LAM                                  
	LBRVIRT                              
	LZCNT                                
	MCAOVERFLOW                          
	MCDT_NO                              
	MCOMMIT                              
	MD_CLEAR                             
	MMX                                  
	MMXEXT                               
	MOVBE                                
	MOVDIR64B                            
	MOVDIRI                              
	MOVSB_ZL                             
	MOVU                                 
	MPX                                  
	MSRIRC                               
	MSRLIST                              
	MSR_PAGEFLUSH                        
	NRIPS                                
	NX                                   
	OSXSAVE                              
	PCONFIG                              
	POPCNT                               
	PPIN                                 
	PREFETCHI                            
	PSFD                                 
	RDPRU                                
	RDRAND                               
	RDSEED                               
	RDTSCP                               
	RRSBA_CTRL                           
	RTM                                  
	RTM_ALWAYS_ABORT                     
	SBPB                                 
	SERIALIZE                            
	SEV                                  
	SEV_64BIT                            
	SEV_ALTERNATIVE                      
	SEV_DEBUGSWAP                        
	SEV_ES                               
	SEV_RESTRICTED                       
	SEV_SNP                              
	SGX                                  
	SGXLC                                
	SHA                                  
	SME                                  
	SME_COHERENT                         
	SPEC_CTRL_SSBD                       
	SRBDS_CTRL                           
	SRSO_MSR_FIX                         
	SRSO_NO                              
	SRSO_USER_KERNEL_NO                  
	SSE                                  
	SSE2                                 
	SSE3                                 
	SSE4                                 
	SSE42                                
	SSE4A                                
	SSSE3                                
	STIBP                                
	STIBP_ALWAYSON                       
	STOSB_SHORT                          
	SUCCOR                               
	SVM                                  
	SVMDA                                
	SVMFBASID                            
	SVML                                 
	SVMNP                                
	SVMPF                                
	SVMPFT                               
	SYSCALL                              
	SYSEE                                
	TBM                                  
	TDX_GUEST                            
	TLB_FLUSH_NESTED                     
	TME                                  
	TOPEXT                               
	TSCRATEMSR                           
	TSXLDTRK                             
	VAES                                 
	VMCBCLEAN                            
	VMPL                                 
	VMSA_REGPROT                         
	VMX                                  
	VPCLMULQDQ                           
	VTE                                  
	WAITPKG                              
	WBNOINVD                             
	WRMSRNS                              
	X87                                  
	XGETBV1                              
	XOP                                  
	XSAVE                                
	XSAVEC                               
	XSAVEOPT                             
	XSAVES                               

	
	AESARM   
	ARMCPUID 
	ASIMD    
	ASIMDDP  
	ASIMDHP  
	ASIMDRDM 
	ATOMICS  
	CRC32    
	DCPOP    
	EVTSTRM  
	FCMA     
	FHM      
	FP       
	FPHP     
	GPA      
	JSCVT    
	LRCPC    
	PMULL    
	RNDR     
	TLB      
	TS       
	SHA1     
	SHA2     
	SHA3     
	SHA512   
	SM3      
	SM4      
	SVE      
	
	lastID

	firstID FeatureID = UNKNOWN + 1
)


type CPUInfo struct {
	BrandName              string  
	VendorID               Vendor  
	VendorString           string  
	HypervisorVendorID     Vendor  
	HypervisorVendorString string  
	featureSet             flagSet 
	PhysicalCores          int     
	ThreadsPerCore         int     
	LogicalCores           int     
	Family                 int     
	Model                  int     
	Stepping               int     
	CacheLine              int     
	Hz                     int64   
	BoostFreq              int64   
	Cache                  struct {
		L1I int 
		L1D int 
		L2  int 
		L3  int 
	}
	SGX              SGXSupport
	AMDMemEncryption AMDMemEncryptionSupport
	AVX10Level       uint8

	maxFunc   uint32
	maxExFunc uint32
}

var cpuid func(op uint32) (eax, ebx, ecx, edx uint32)
var cpuidex func(op, op2 uint32) (eax, ebx, ecx, edx uint32)
var xgetbv func(index uint32) (eax, edx uint32)
var rdtscpAsm func() (eax, ebx, ecx, edx uint32)
var darwinHasAVX512 = func() bool { return false }





var CPU CPUInfo

func init() {
	initCPU()
	Detect()
}








func Detect() {
	
	CPU.ThreadsPerCore = 1
	CPU.Cache.L1I = -1
	CPU.Cache.L1D = -1
	CPU.Cache.L2 = -1
	CPU.Cache.L3 = -1
	safe := true
	if detectArmFlag != nil {
		safe = !*detectArmFlag
	}
	addInfo(&CPU, safe)
	if displayFeats != nil && *displayFeats {
		fmt.Println("cpu features:", strings.Join(CPU.FeatureSet(), ","))
		
		os.Exit(1)
	}
	if disableFlag != nil {
		s := strings.Split(*disableFlag, ",")
		for _, feat := range s {
			feat := ParseFeature(strings.TrimSpace(feat))
			if feat != UNKNOWN {
				CPU.featureSet.unset(feat)
			}
		}
	}
}






func DetectARM() {
	addInfo(&CPU, false)
}

var detectArmFlag *bool
var displayFeats *bool
var disableFlag *string






func Flags() {
	disableFlag = flag.String("cpu.disable", "", "disable cpu features; comma separated list")
	displayFeats = flag.Bool("cpu.features", false, "lists cpu features and exits")
	detectArmFlag = flag.Bool("cpu.arm", false, "allow ARM features to be detected; can potentially crash")
}


func (c CPUInfo) Supports(ids ...FeatureID) bool {
	for _, id := range ids {
		if !c.featureSet.inSet(id) {
			return false
		}
	}
	return true
}



func (c *CPUInfo) Has(id FeatureID) bool {
	return c.featureSet.inSet(id)
}


func (c CPUInfo) AnyOf(ids ...FeatureID) bool {
	for _, id := range ids {
		if c.featureSet.inSet(id) {
			return true
		}
	}
	return false
}



type Features *flagSet


func CombineFeatures(ids ...FeatureID) Features {
	var v flagSet
	for _, id := range ids {
		v.set(id)
	}
	return &v
}

func (c *CPUInfo) HasAll(f Features) bool {
	return c.featureSet.hasSetP(f)
}


var oneOfLevel = CombineFeatures(SYSEE, SYSCALL)
var level1Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2)
var level2Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2, CX16, LAHF, POPCNT, SSE3, SSE4, SSE42, SSSE3)
var level3Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2, CX16, LAHF, POPCNT, SSE3, SSE4, SSE42, SSSE3, AVX, AVX2, BMI1, BMI2, F16C, FMA3, LZCNT, MOVBE, OSXSAVE)
var level4Features = CombineFeatures(CMOV, CMPXCHG8, X87, FXSR, MMX, SSE, SSE2, CX16, LAHF, POPCNT, SSE3, SSE4, SSE42, SSSE3, AVX, AVX2, BMI1, BMI2, F16C, FMA3, LZCNT, MOVBE, OSXSAVE, AVX512F, AVX512BW, AVX512CD, AVX512DQ, AVX512VL)




func (c CPUInfo) X64Level() int {
	if !c.featureSet.hasOneOf(oneOfLevel) {
		return 0
	}
	if c.featureSet.hasSetP(level4Features) {
		return 4
	}
	if c.featureSet.hasSetP(level3Features) {
		return 3
	}
	if c.featureSet.hasSetP(level2Features) {
		return 2
	}
	if c.featureSet.hasSetP(level1Features) {
		return 1
	}
	return 0
}


func (c *CPUInfo) Disable(ids ...FeatureID) bool {
	for _, id := range ids {
		c.featureSet.unset(id)
	}
	return true
}



func (c *CPUInfo) Enable(ids ...FeatureID) bool {
	for _, id := range ids {
		c.featureSet.set(id)
	}
	return true
}


func (c CPUInfo) IsVendor(v Vendor) bool {
	return c.VendorID == v
}


func (c CPUInfo) FeatureSet() []string {
	s := make([]string, 0, c.featureSet.nEnabled())
	s = append(s, c.featureSet.Strings()...)
	return s
}




func (c CPUInfo) RTCounter() uint64 {
	if !c.Has(RDTSCP) {
		return 0
	}
	a, _, _, d := rdtscpAsm()
	return uint64(a) | (uint64(d) << 32)
}





func (c CPUInfo) Ia32TscAux() uint32 {
	if !c.Has(RDTSCP) {
		return 0
	}
	_, _, ecx, _ := rdtscpAsm()
	return ecx
}



func (c CPUInfo) SveLengths() (vl, pl uint64) {
	if !c.Has(SVE) {
		return 0, 0
	}
	return getVectorLength()
}





func (c CPUInfo) LogicalCPU() int {
	if c.maxFunc < 1 {
		return -1
	}
	_, ebx, _, _ := cpuid(1)
	return int(ebx >> 24)
}



func (c *CPUInfo) frequencies() {
	c.Hz, c.BoostFreq = 0, 0
	mfi := maxFunctionID()
	if mfi >= 0x15 {
		eax, ebx, ecx, _ := cpuid(0x15)
		if eax != 0 && ebx != 0 && ecx != 0 {
			c.Hz = (int64(ecx) * int64(ebx)) / int64(eax)
		}
	}
	if mfi >= 0x16 {
		a, b, _, _ := cpuid(0x16)
		
		if a&0xffff > 0 {
			c.Hz = int64(a&0xffff) * 1_000_000
		}
		
		if b&0xffff > 0 {
			c.BoostFreq = int64(b&0xffff) * 1_000_000
		}
	}
	if c.Hz > 0 {
		return
	}

	
	
	
	
	
	
	model := c.BrandName
	hz := strings.LastIndex(model, "Hz")
	if hz < 3 {
		return
	}
	var multiplier int64
	switch model[hz-1] {
	case 'M':
		multiplier = 1000 * 1000
	case 'G':
		multiplier = 1000 * 1000 * 1000
	case 'T':
		multiplier = 1000 * 1000 * 1000 * 1000
	}
	if multiplier == 0 {
		return
	}
	freq := int64(0)
	divisor := int64(0)
	decimalShift := int64(1)
	var i int
	for i = hz - 2; i >= 0 && model[i] != ' '; i-- {
		if model[i] >= '0' && model[i] <= '9' {
			freq += int64(model[i]-'0') * decimalShift
			decimalShift *= 10
		} else if model[i] == '.' {
			if divisor != 0 {
				return
			}
			divisor = decimalShift
		} else {
			return
		}
	}
	
	if i < 0 {
		return
	}
	if divisor != 0 {
		c.Hz = (freq * multiplier) / divisor
		return
	}
	c.Hz = freq * multiplier
}



func (c CPUInfo) VM() bool {
	return CPU.featureSet.inSet(HYPERVISOR)
}


type flags uint64


const flagBitsLog2 = 6
const flagBits = 1 << flagBitsLog2
const flagMask = flagBits - 1


type flagSet [(lastID + flagMask) / flagBits]flags

func (s *flagSet) inSet(feat FeatureID) bool {
	return s[feat>>flagBitsLog2]&(1<<(feat&flagMask)) != 0
}

func (s *flagSet) set(feat FeatureID) {
	s[feat>>flagBitsLog2] |= 1 << (feat & flagMask)
}


func (s *flagSet) setIf(cond bool, features ...FeatureID) {
	if cond {
		for _, offset := range features {
			s[offset>>flagBitsLog2] |= 1 << (offset & flagMask)
		}
	}
}

func (s *flagSet) unset(offset FeatureID) {
	bit := flags(1 << (offset & flagMask))
	s[offset>>flagBitsLog2] = s[offset>>flagBitsLog2] & ^bit
}


func (s *flagSet) or(other flagSet) {
	for i, v := range other[:] {
		s[i] |= v
	}
}


func (s *flagSet) hasSet(other flagSet) bool {
	for i, v := range other[:] {
		if s[i]&v != v {
			return false
		}
	}
	return true
}


func (s *flagSet) hasSetP(other *flagSet) bool {
	for i, v := range other[:] {
		if s[i]&v != v {
			return false
		}
	}
	return true
}


func (s *flagSet) hasOneOf(other *flagSet) bool {
	for i, v := range other[:] {
		if s[i]&v != 0 {
			return true
		}
	}
	return false
}


func (s *flagSet) nEnabled() (n int) {
	for _, v := range s[:] {
		n += bits.OnesCount64(uint64(v))
	}
	return n
}

func flagSetWith(feat ...FeatureID) flagSet {
	var res flagSet
	for _, f := range feat {
		res.set(f)
	}
	return res
}



func ParseFeature(s string) FeatureID {
	s = strings.ToUpper(s)
	for i := firstID; i < lastID; i++ {
		if i.String() == s {
			return i
		}
	}
	return UNKNOWN
}


func (s flagSet) Strings() []string {
	if len(s) == 0 {
		return []string{""}
	}
	r := make([]string, 0)
	for i := firstID; i < lastID; i++ {
		if s.inSet(i) {
			r = append(r, i.String())
		}
	}
	return r
}

func maxExtendedFunction() uint32 {
	eax, _, _, _ := cpuid(0x80000000)
	return eax
}

func maxFunctionID() uint32 {
	a, _, _, _ := cpuid(0)
	return a
}

func brandName() string {
	if maxExtendedFunction() >= 0x80000004 {
		v := make([]uint32, 0, 48)
		for i := uint32(0); i < 3; i++ {
			a, b, c, d := cpuid(0x80000002 + i)
			v = append(v, a, b, c, d)
		}
		return strings.Trim(string(valAsString(v...)), " ")
	}
	return "unknown"
}

func threadsPerCore() int {
	mfi := maxFunctionID()
	vend, _ := vendorID()

	if mfi < 0x4 || (vend != Intel && vend != AMD) {
		return 1
	}

	if mfi < 0xb {
		if vend != Intel {
			return 1
		}
		_, b, _, d := cpuid(1)
		if (d & (1 << 28)) != 0 {
			
			v := (b >> 16) & 255
			if v > 1 {
				a4, _, _, _ := cpuid(4)
				
				v2 := (a4 >> 26) + 1
				if v2 > 0 {
					return int(v) / int(v2)
				}
			}
		}
		return 1
	}
	_, b, _, _ := cpuidex(0xb, 0)
	if b&0xffff == 0 {
		if vend == AMD {
			
			
			
			fam, _, _ := familyModel()
			_, _, _, d := cpuid(1)
			if (d&(1<<28)) != 0 && fam >= 23 {
				if maxExtendedFunction() >= 0x8000001e {
					_, b, _, _ := cpuid(0x8000001e)
					return int((b>>8)&0xff) + 1
				}
				return 2
			}
		}
		return 1
	}
	return int(b & 0xffff)
}

func logicalCores() int {
	mfi := maxFunctionID()
	v, _ := vendorID()
	switch v {
	case Intel:
		
		if mfi < 0xb {
			if mfi < 1 {
				return 0
			}
			
			
			
			_, ebx, _, _ := cpuid(1)
			logical := (ebx >> 16) & 0xff
			return int(logical)
		}
		_, b, _, _ := cpuidex(0xb, 1)
		return int(b & 0xffff)
	case AMD, Hygon:
		_, b, _, _ := cpuid(1)
		return int((b >> 16) & 0xff)
	default:
		return 0
	}
}

func familyModel() (family, model, stepping int) {
	if maxFunctionID() < 0x1 {
		return 0, 0, 0
	}
	eax, _, _, _ := cpuid(1)
	
	family = int((eax >> 8) & 0xf)
	extFam := family == 0x6 
	if family == 0xf {
		
		family += int((eax >> 20) & 0xff)
		extFam = true
	}
	
	model = int((eax >> 4) & 0xf)
	if extFam {
		
		model += int((eax >> 12) & 0xf0)
	}
	stepping = int(eax & 0xf)
	return family, model, stepping
}

func physicalCores() int {
	v, _ := vendorID()
	switch v {
	case Intel:
		return logicalCores() / threadsPerCore()
	case AMD, Hygon:
		lc := logicalCores()
		tpc := threadsPerCore()
		if lc > 0 && tpc > 0 {
			return lc / tpc
		}

		
		if maxExtendedFunction() >= 0x80000008 {
			_, _, c, _ := cpuid(0x80000008)
			if c&0xff > 0 {
				return int(c&0xff) + 1
			}
		}
	}
	return 0
}


var vendorMapping = map[string]Vendor{
	"AMDisbetter!": AMD,
	"AuthenticAMD": AMD,
	"CentaurHauls": VIA,
	"GenuineIntel": Intel,
	"TransmetaCPU": Transmeta,
	"GenuineTMx86": Transmeta,
	"Geode by NSC": NSC,
	"VIA VIA VIA ": VIA,
	"KVMKVMKVM":    KVM,
	"Linux KVM Hv": KVM,
	"TCGTCGTCGTCG": QEMU,
	"Microsoft Hv": MSVM,
	"VMwareVMware": VMware,
	"XenVMMXenVMM": XenHVM,
	"bhyve bhyve ": Bhyve,
	"HygonGenuine": Hygon,
	"Vortex86 SoC": SiS,
	"SiS SiS SiS ": SiS,
	"RiseRiseRise": SiS,
	"Genuine  RDC": RDC,
	"QNXQVMBSQG":   QNX,
	"ACRNACRNACRN": ACRN,
	"SRESRESRESRE": SRE,
	"Apple VZ":     Apple,
}

func vendorID() (Vendor, string) {
	_, b, c, d := cpuid(0)
	v := string(valAsString(b, d, c))
	vend, ok := vendorMapping[v]
	if !ok {
		return VendorUnknown, v
	}
	return vend, v
}

func hypervisorVendorID() (Vendor, string) {
	
	_, b, c, d := cpuid(0x40000000)
	v := string(valAsString(b, c, d))
	vend, ok := vendorMapping[v]
	if !ok {
		return VendorUnknown, v
	}
	return vend, v
}

func cacheLine() int {
	if maxFunctionID() < 0x1 {
		return 0
	}

	_, ebx, _, _ := cpuid(1)
	cache := (ebx & 0xff00) >> 5 
	if cache == 0 && maxExtendedFunction() >= 0x80000006 {
		_, _, ecx, _ := cpuid(0x80000006)
		cache = ecx & 0xff 
	}
	
	return int(cache)
}

func (c *CPUInfo) cacheSize() {
	c.Cache.L1D = -1
	c.Cache.L1I = -1
	c.Cache.L2 = -1
	c.Cache.L3 = -1
	vendor, _ := vendorID()
	switch vendor {
	case Intel:
		if maxFunctionID() < 4 {
			return
		}
		c.Cache.L1I, c.Cache.L1D, c.Cache.L2, c.Cache.L3 = 0, 0, 0, 0
		for i := uint32(0); ; i++ {
			eax, ebx, ecx, _ := cpuidex(4, i)
			cacheType := eax & 15
			if cacheType == 0 {
				break
			}
			cacheLevel := (eax >> 5) & 7
			coherency := int(ebx&0xfff) + 1
			partitions := int((ebx>>12)&0x3ff) + 1
			associativity := int((ebx>>22)&0x3ff) + 1
			sets := int(ecx) + 1
			size := associativity * partitions * coherency * sets
			switch cacheLevel {
			case 1:
				if cacheType == 1 {
					
					c.Cache.L1D = size
				} else if cacheType == 2 {
					
					c.Cache.L1I = size
				} else {
					if c.Cache.L1D < 0 {
						c.Cache.L1I = size
					}
					if c.Cache.L1I < 0 {
						c.Cache.L1I = size
					}
				}
			case 2:
				c.Cache.L2 = size
			case 3:
				c.Cache.L3 = size
			}
		}
	case AMD, Hygon:
		
		if maxExtendedFunction() < 0x80000005 {
			return
		}
		_, _, ecx, edx := cpuid(0x80000005)
		c.Cache.L1D = int(((ecx >> 24) & 0xFF) * 1024)
		c.Cache.L1I = int(((edx >> 24) & 0xFF) * 1024)

		if maxExtendedFunction() < 0x80000006 {
			return
		}
		_, _, ecx, _ = cpuid(0x80000006)
		c.Cache.L2 = int(((ecx >> 16) & 0xFFFF) * 1024)

		
		if maxExtendedFunction() < 0x8000001D || !c.Has(TOPEXT) {
			return
		}

		
		
		nSame := 0
		var last uint32
		for i := uint32(0); i < math.MaxUint32; i++ {
			eax, ebx, ecx, _ := cpuidex(0x8000001D, i)

			level := (eax >> 5) & 7
			cacheNumSets := ecx + 1
			cacheLineSize := 1 + (ebx & 2047)
			cachePhysPartitions := 1 + ((ebx >> 12) & 511)
			cacheNumWays := 1 + ((ebx >> 22) & 511)

			typ := eax & 15
			size := int(cacheNumSets * cacheLineSize * cachePhysPartitions * cacheNumWays)
			if typ == 0 {
				return
			}

			
			comb := eax ^ ebx ^ ecx
			if comb == last {
				nSame++
				if nSame == 100 {
					return
				}
			}
			last = comb

			switch level {
			case 1:
				switch typ {
				case 1:
					
					c.Cache.L1D = size
				case 2:
					
					c.Cache.L1I = size
				default:
					if c.Cache.L1D < 0 {
						c.Cache.L1I = size
					}
					if c.Cache.L1I < 0 {
						c.Cache.L1I = size
					}
				}
			case 2:
				c.Cache.L2 = size
			case 3:
				c.Cache.L3 = size
			}
		}
	}
}

type SGXEPCSection struct {
	BaseAddress uint64
	EPCSize     uint64
}

type SGXSupport struct {
	Available           bool
	LaunchControl       bool
	SGX1Supported       bool
	SGX2Supported       bool
	MaxEnclaveSizeNot64 int64
	MaxEnclaveSize64    int64
	EPCSections         []SGXEPCSection
}

func hasSGX(available, lc bool) (rval SGXSupport) {
	rval.Available = available

	if !available {
		return
	}

	rval.LaunchControl = lc

	a, _, _, d := cpuidex(0x12, 0)
	rval.SGX1Supported = a&0x01 != 0
	rval.SGX2Supported = a&0x02 != 0
	rval.MaxEnclaveSizeNot64 = 1 << (d & 0xFF)     
	rval.MaxEnclaveSize64 = 1 << ((d >> 8) & 0xFF) 
	rval.EPCSections = make([]SGXEPCSection, 0)

	for subleaf := uint32(2); subleaf < 2+8; subleaf++ {
		eax, ebx, ecx, edx := cpuidex(0x12, subleaf)
		leafType := eax & 0xf

		if leafType == 0 {
			
			break
		} else if leafType == 1 {
			
			baseAddress := uint64(eax&0xfffff000) + (uint64(ebx&0x000fffff) << 32)
			size := uint64(ecx&0xfffff000) + (uint64(edx&0x000fffff) << 32)

			section := SGXEPCSection{BaseAddress: baseAddress, EPCSize: size}
			rval.EPCSections = append(rval.EPCSections, section)
		}
	}

	return
}

type AMDMemEncryptionSupport struct {
	Available          bool
	CBitPossition      uint32
	NumVMPL            uint32
	PhysAddrReduction  uint32
	NumEntryptedGuests uint32
	MinSevNoEsAsid     uint32
}

func hasAMDMemEncryption(available bool) (rval AMDMemEncryptionSupport) {
	rval.Available = available
	if !available {
		return
	}

	_, b, c, d := cpuidex(0x8000001f, 0)

	rval.CBitPossition = b & 0x3f
	rval.PhysAddrReduction = (b >> 6) & 0x3F
	rval.NumVMPL = (b >> 12) & 0xf
	rval.NumEntryptedGuests = c
	rval.MinSevNoEsAsid = d

	return
}

func support() flagSet {
	var fs flagSet
	mfi := maxFunctionID()
	vend, _ := vendorID()
	if mfi < 0x1 {
		return fs
	}
	family, model, _ := familyModel()

	_, _, c, d := cpuid(1)
	fs.setIf((d&(1<<0)) != 0, X87)
	fs.setIf((d&(1<<8)) != 0, CMPXCHG8)
	fs.setIf((d&(1<<11)) != 0, SYSEE)
	fs.setIf((d&(1<<15)) != 0, CMOV)
	fs.setIf((d&(1<<23)) != 0, MMX)
	fs.setIf((d&(1<<24)) != 0, FXSR)
	fs.setIf((d&(1<<25)) != 0, FXSROPT)
	fs.setIf((d&(1<<25)) != 0, SSE)
	fs.setIf((d&(1<<26)) != 0, SSE2)
	fs.setIf((c&1) != 0, SSE3)
	fs.setIf((c&(1<<5)) != 0, VMX)
	fs.setIf((c&(1<<9)) != 0, SSSE3)
	fs.setIf((c&(1<<19)) != 0, SSE4)
	fs.setIf((c&(1<<20)) != 0, SSE42)
	fs.setIf((c&(1<<25)) != 0, AESNI)
	fs.setIf((c&(1<<1)) != 0, CLMUL)
	fs.setIf(c&(1<<22) != 0, MOVBE)
	fs.setIf(c&(1<<23) != 0, POPCNT)
	fs.setIf(c&(1<<30) != 0, RDRAND)

	
	
	fs.setIf(c&(1<<31) != 0, HYPERVISOR)
	fs.setIf(c&(1<<29) != 0, F16C)
	fs.setIf(c&(1<<13) != 0, CX16)

	if vend == Intel && (d&(1<<28)) != 0 && mfi >= 4 {
		fs.setIf(threadsPerCore() > 1, HTT)
	}
	if vend == AMD && (d&(1<<28)) != 0 && mfi >= 4 {
		fs.setIf(threadsPerCore() > 1, HTT)
	}
	fs.setIf(c&1<<26 != 0, XSAVE)
	fs.setIf(c&1<<27 != 0, OSXSAVE)
	
	const avxCheck = 1<<26 | 1<<27 | 1<<28
	if c&avxCheck == avxCheck {
		
		eax, _ := xgetbv(0)
		if (eax & 0x6) == 0x6 {
			fs.set(AVX)
			switch vend {
			case Intel:
				
				fs.setIf(family == 6 && model < 60, AVXSLOW)
			case AMD:
				
				fs.setIf(family < 23 || (family == 23 && model < 49), AVXSLOW)
			}
		}
	}
	
	
	const fma3Check = 1<<12 | 1<<27
	fs.setIf(c&fma3Check == fma3Check, FMA3)

	
	if mfi >= 7 {
		_, ebx, ecx, edx := cpuidex(7, 0)
		if fs.inSet(AVX) && (ebx&0x00000020) != 0 {
			fs.set(AVX2)
		}
		
		if (ebx & 0x00000008) != 0 {
			fs.set(BMI1)
			fs.setIf((ebx&0x00000100) != 0, BMI2)
		}
		fs.setIf(ebx&(1<<2) != 0, SGX)
		fs.setIf(ebx&(1<<4) != 0, HLE)
		fs.setIf(ebx&(1<<9) != 0, ERMS)
		fs.setIf(ebx&(1<<11) != 0, RTM)
		fs.setIf(ebx&(1<<14) != 0, MPX)
		fs.setIf(ebx&(1<<18) != 0, RDSEED)
		fs.setIf(ebx&(1<<19) != 0, ADX)
		fs.setIf(ebx&(1<<29) != 0, SHA)

		
		fs.setIf(ecx&(1<<5) != 0, WAITPKG)
		fs.setIf(ecx&(1<<7) != 0, CETSS)
		fs.setIf(ecx&(1<<8) != 0, GFNI)
		fs.setIf(ecx&(1<<9) != 0, VAES)
		fs.setIf(ecx&(1<<10) != 0, VPCLMULQDQ)
		fs.setIf(ecx&(1<<13) != 0, TME)
		fs.setIf(ecx&(1<<25) != 0, CLDEMOTE)
		fs.setIf(ecx&(1<<23) != 0, KEYLOCKER)
		fs.setIf(ecx&(1<<27) != 0, MOVDIRI)
		fs.setIf(ecx&(1<<28) != 0, MOVDIR64B)
		fs.setIf(ecx&(1<<29) != 0, ENQCMD)
		fs.setIf(ecx&(1<<30) != 0, SGXLC)

		
		fs.setIf(edx&(1<<4) != 0, FSRM)
		fs.setIf(edx&(1<<9) != 0, SRBDS_CTRL)
		fs.setIf(edx&(1<<10) != 0, MD_CLEAR)
		fs.setIf(edx&(1<<11) != 0, RTM_ALWAYS_ABORT)
		fs.setIf(edx&(1<<14) != 0, SERIALIZE)
		fs.setIf(edx&(1<<15) != 0, HYBRID_CPU)
		fs.setIf(edx&(1<<16) != 0, TSXLDTRK)
		fs.setIf(edx&(1<<18) != 0, PCONFIG)
		fs.setIf(edx&(1<<20) != 0, CETIBT)
		fs.setIf(edx&(1<<26) != 0, IBPB)
		fs.setIf(edx&(1<<27) != 0, STIBP)
		fs.setIf(edx&(1<<28) != 0, FLUSH_L1D)
		fs.setIf(edx&(1<<29) != 0, IA32_ARCH_CAP)
		fs.setIf(edx&(1<<30) != 0, IA32_CORE_CAP)
		fs.setIf(edx&(1<<31) != 0, SPEC_CTRL_SSBD)

		
		eax1, _, _, edx1 := cpuidex(7, 1)
		fs.setIf(fs.inSet(AVX) && eax1&(1<<4) != 0, AVXVNNI)
		fs.setIf(eax1&(1<<7) != 0, CMPCCXADD)
		fs.setIf(eax1&(1<<10) != 0, MOVSB_ZL)
		fs.setIf(eax1&(1<<11) != 0, STOSB_SHORT)
		fs.setIf(eax1&(1<<12) != 0, CMPSB_SCADBS_SHORT)
		fs.setIf(eax1&(1<<22) != 0, HRESET)
		fs.setIf(eax1&(1<<23) != 0, AVXIFMA)
		fs.setIf(eax1&(1<<26) != 0, LAM)

		
		fs.setIf(edx1&(1<<4) != 0, AVXVNNIINT8)
		fs.setIf(edx1&(1<<5) != 0, AVXNECONVERT)
		fs.setIf(edx1&(1<<7) != 0, AMXTF32)
		fs.setIf(edx1&(1<<8) != 0, AMXCOMPLEX)
		fs.setIf(edx1&(1<<10) != 0, AVXVNNIINT16)
		fs.setIf(edx1&(1<<14) != 0, PREFETCHI)
		fs.setIf(edx1&(1<<19) != 0, AVX10)
		fs.setIf(edx1&(1<<21) != 0, APX_F)

		
		if c&((1<<26)|(1<<27)) == (1<<26)|(1<<27) {
			
			eax, _ := xgetbv(0)

			
			
			
			hasAVX512 := (eax>>5)&7 == 7 && (eax>>1)&3 == 3
			if runtime.GOOS == "darwin" {
				hasAVX512 = fs.inSet(AVX) && darwinHasAVX512()
			}
			if hasAVX512 {
				fs.setIf(ebx&(1<<16) != 0, AVX512F)
				fs.setIf(ebx&(1<<17) != 0, AVX512DQ)
				fs.setIf(ebx&(1<<21) != 0, AVX512IFMA)
				fs.setIf(ebx&(1<<26) != 0, AVX512PF)
				fs.setIf(ebx&(1<<27) != 0, AVX512ER)
				fs.setIf(ebx&(1<<28) != 0, AVX512CD)
				fs.setIf(ebx&(1<<30) != 0, AVX512BW)
				fs.setIf(ebx&(1<<31) != 0, AVX512VL)
				
				fs.setIf(ecx&(1<<1) != 0, AVX512VBMI)
				fs.setIf(ecx&(1<<3) != 0, AMXFP8)
				fs.setIf(ecx&(1<<6) != 0, AVX512VBMI2)
				fs.setIf(ecx&(1<<11) != 0, AVX512VNNI)
				fs.setIf(ecx&(1<<12) != 0, AVX512BITALG)
				fs.setIf(ecx&(1<<14) != 0, AVX512VPOPCNTDQ)
				
				fs.setIf(edx&(1<<8) != 0, AVX512VP2INTERSECT)
				fs.setIf(edx&(1<<22) != 0, AMXBF16)
				fs.setIf(edx&(1<<23) != 0, AVX512FP16)
				fs.setIf(edx&(1<<24) != 0, AMXTILE)
				fs.setIf(edx&(1<<25) != 0, AMXINT8)
				
				fs.setIf(eax1&(1<<5) != 0, AVX512BF16)
				fs.setIf(eax1&(1<<19) != 0, WRMSRNS)
				fs.setIf(eax1&(1<<21) != 0, AMXFP16)
				fs.setIf(eax1&(1<<27) != 0, MSRLIST)
			}
		}

		
		_, _, _, edx = cpuidex(7, 2)
		fs.setIf(edx&(1<<0) != 0, PSFD)
		fs.setIf(edx&(1<<1) != 0, IDPRED_CTRL)
		fs.setIf(edx&(1<<2) != 0, RRSBA_CTRL)
		fs.setIf(edx&(1<<4) != 0, BHI_CTRL)
		fs.setIf(edx&(1<<5) != 0, MCDT_NO)

		
		if fs.inSet(KEYLOCKER) && mfi >= 0x19 {
			_, ebx, _, _ := cpuidex(0x19, 0)
			fs.setIf(ebx&5 == 5, KEYLOCKERW) 
		}

		
		if fs.inSet(AVX10) && mfi >= 0x24 {
			_, ebx, _, _ := cpuidex(0x24, 0)
			fs.setIf(ebx&(1<<16) != 0, AVX10_128)
			fs.setIf(ebx&(1<<17) != 0, AVX10_256)
			fs.setIf(ebx&(1<<18) != 0, AVX10_512)
		}
	}

	
	
	
	
	
	
	
	
	
	
	
	
	
	if mfi >= 0xd {
		if fs.inSet(XSAVE) {
			eax, _, _, _ := cpuidex(0xd, 1)
			fs.setIf(eax&(1<<0) != 0, XSAVEOPT)
			fs.setIf(eax&(1<<1) != 0, XSAVEC)
			fs.setIf(eax&(1<<2) != 0, XGETBV1)
			fs.setIf(eax&(1<<3) != 0, XSAVES)
		}
	}
	if maxExtendedFunction() >= 0x80000001 {
		_, _, c, d := cpuid(0x80000001)
		if (c & (1 << 5)) != 0 {
			fs.set(LZCNT)
			fs.set(POPCNT)
		}
		
		fs.setIf((c&(1<<0)) != 0, LAHF)
		fs.setIf((c&(1<<2)) != 0, SVM)
		fs.setIf((c&(1<<6)) != 0, SSE4A)
		fs.setIf((c&(1<<10)) != 0, IBS)
		fs.setIf((c&(1<<22)) != 0, TOPEXT)

		
		fs.setIf(d&(1<<11) != 0, SYSCALL)
		fs.setIf(d&(1<<20) != 0, NX)
		fs.setIf(d&(1<<22) != 0, MMXEXT)
		fs.setIf(d&(1<<23) != 0, MMX)
		fs.setIf(d&(1<<24) != 0, FXSR)
		fs.setIf(d&(1<<25) != 0, FXSROPT)
		fs.setIf(d&(1<<27) != 0, RDTSCP)
		fs.setIf(d&(1<<30) != 0, AMD3DNOWEXT)
		fs.setIf(d&(1<<31) != 0, AMD3DNOW)

		
		if fs.inSet(AVX) {
			fs.setIf((c&(1<<11)) != 0, XOP)
			fs.setIf((c&(1<<16)) != 0, FMA4)
		}

	}
	if maxExtendedFunction() >= 0x80000007 {
		_, b, _, d := cpuid(0x80000007)
		fs.setIf((b&(1<<0)) != 0, MCAOVERFLOW)
		fs.setIf((b&(1<<1)) != 0, SUCCOR)
		fs.setIf((b&(1<<2)) != 0, HWA)
		fs.setIf((d&(1<<9)) != 0, CPBOOST)
	}

	if maxExtendedFunction() >= 0x80000008 {
		_, b, _, _ := cpuid(0x80000008)
		fs.setIf(b&(1<<28) != 0, PSFD)
		fs.setIf(b&(1<<27) != 0, CPPC)
		fs.setIf(b&(1<<24) != 0, SPEC_CTRL_SSBD)
		fs.setIf(b&(1<<23) != 0, PPIN)
		fs.setIf(b&(1<<21) != 0, TLB_FLUSH_NESTED)
		fs.setIf(b&(1<<20) != 0, EFER_LMSLE_UNS)
		fs.setIf(b&(1<<19) != 0, IBRS_PROVIDES_SMP)
		fs.setIf(b&(1<<18) != 0, IBRS_PREFERRED)
		fs.setIf(b&(1<<17) != 0, STIBP_ALWAYSON)
		fs.setIf(b&(1<<15) != 0, STIBP)
		fs.setIf(b&(1<<14) != 0, IBRS)
		fs.setIf((b&(1<<13)) != 0, INT_WBINVD)
		fs.setIf(b&(1<<12) != 0, IBPB)
		fs.setIf((b&(1<<9)) != 0, WBNOINVD)
		fs.setIf((b&(1<<8)) != 0, MCOMMIT)
		fs.setIf((b&(1<<4)) != 0, RDPRU)
		fs.setIf((b&(1<<3)) != 0, INVLPGB)
		fs.setIf((b&(1<<1)) != 0, MSRIRC)
		fs.setIf((b&(1<<0)) != 0, CLZERO)
	}

	if fs.inSet(SVM) && maxExtendedFunction() >= 0x8000000A {
		_, _, _, edx := cpuid(0x8000000A)
		fs.setIf((edx>>0)&1 == 1, SVMNP)
		fs.setIf((edx>>1)&1 == 1, LBRVIRT)
		fs.setIf((edx>>2)&1 == 1, SVML)
		fs.setIf((edx>>3)&1 == 1, NRIPS)
		fs.setIf((edx>>4)&1 == 1, TSCRATEMSR)
		fs.setIf((edx>>5)&1 == 1, VMCBCLEAN)
		fs.setIf((edx>>6)&1 == 1, SVMFBASID)
		fs.setIf((edx>>7)&1 == 1, SVMDA)
		fs.setIf((edx>>10)&1 == 1, SVMPF)
		fs.setIf((edx>>12)&1 == 1, SVMPFT)
	}

	if maxExtendedFunction() >= 0x8000001a {
		eax, _, _, _ := cpuid(0x8000001a)
		fs.setIf((eax>>0)&1 == 1, FP128)
		fs.setIf((eax>>1)&1 == 1, MOVU)
		fs.setIf((eax>>2)&1 == 1, FP256)
	}

	if maxExtendedFunction() >= 0x8000001b && fs.inSet(IBS) {
		eax, _, _, _ := cpuid(0x8000001b)
		fs.setIf((eax>>0)&1 == 1, IBSFFV)
		fs.setIf((eax>>1)&1 == 1, IBSFETCHSAM)
		fs.setIf((eax>>2)&1 == 1, IBSOPSAM)
		fs.setIf((eax>>3)&1 == 1, IBSRDWROPCNT)
		fs.setIf((eax>>4)&1 == 1, IBSOPCNT)
		fs.setIf((eax>>5)&1 == 1, IBSBRNTRGT)
		fs.setIf((eax>>6)&1 == 1, IBSOPCNTEXT)
		fs.setIf((eax>>7)&1 == 1, IBSRIPINVALIDCHK)
		fs.setIf((eax>>8)&1 == 1, IBS_OPFUSE)
		fs.setIf((eax>>9)&1 == 1, IBS_FETCH_CTLX)
		fs.setIf((eax>>10)&1 == 1, IBS_OPDATA4) 
		fs.setIf((eax>>11)&1 == 1, IBS_ZEN4)
	}

	if maxExtendedFunction() >= 0x8000001f && vend == AMD {
		a, _, _, _ := cpuid(0x8000001f)
		fs.setIf((a>>0)&1 == 1, SME)
		fs.setIf((a>>1)&1 == 1, SEV)
		fs.setIf((a>>2)&1 == 1, MSR_PAGEFLUSH)
		fs.setIf((a>>3)&1 == 1, SEV_ES)
		fs.setIf((a>>4)&1 == 1, SEV_SNP)
		fs.setIf((a>>5)&1 == 1, VMPL)
		fs.setIf((a>>10)&1 == 1, SME_COHERENT)
		fs.setIf((a>>11)&1 == 1, SEV_64BIT)
		fs.setIf((a>>12)&1 == 1, SEV_RESTRICTED)
		fs.setIf((a>>13)&1 == 1, SEV_ALTERNATIVE)
		fs.setIf((a>>14)&1 == 1, SEV_DEBUGSWAP)
		fs.setIf((a>>15)&1 == 1, IBS_PREVENTHOST)
		fs.setIf((a>>16)&1 == 1, VTE)
		fs.setIf((a>>24)&1 == 1, VMSA_REGPROT)
	}

	if maxExtendedFunction() >= 0x80000021 && vend == AMD {
		a, _, _, _ := cpuid(0x80000021)
		fs.setIf((a>>31)&1 == 1, SRSO_MSR_FIX)
		fs.setIf((a>>30)&1 == 1, SRSO_USER_KERNEL_NO)
		fs.setIf((a>>29)&1 == 1, SRSO_NO)
		fs.setIf((a>>28)&1 == 1, IBPB_BRTYPE)
		fs.setIf((a>>27)&1 == 1, SBPB)
	}

	if mfi >= 0x20 {
		
		
		
		
		
		
		
		
		
		_, ebx, _, _ := cpuid(0x4000000C)
		fs.setIf(ebx == 0xbe3, TDX_GUEST)
	}

	if mfi >= 0x21 {
		
		_, ebx, ecx, edx := cpuid(0x21)
		identity := string(valAsString(ebx, edx, ecx))
		fs.setIf(identity == "IntelTDX    ", TDX_GUEST)
	}

	return fs
}

func (c *CPUInfo) supportAVX10() uint8 {
	if c.maxFunc >= 0x24 && c.featureSet.inSet(AVX10) {
		_, ebx, _, _ := cpuidex(0x24, 0)
		return uint8(ebx)
	}
	return 0
}

func valAsString(values ...uint32) []byte {
	r := make([]byte, 4*len(values))
	for i, v := range values {
		dst := r[i*4:]
		dst[0] = byte(v & 0xff)
		dst[1] = byte((v >> 8) & 0xff)
		dst[2] = byte((v >> 16) & 0xff)
		dst[3] = byte((v >> 24) & 0xff)
		switch {
		case dst[0] == 0:
			return r[:i*4]
		case dst[1] == 0:
			return r[:i*4+1]
		case dst[2] == 0:
			return r[:i*4+2]
		case dst[3] == 0:
			return r[:i*4+3]
		}
	}
	return r
}
