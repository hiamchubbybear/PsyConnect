





package cpu

import (
	"os"
	"strings"
)






var Initialized bool


type CacheLinePad struct{ _ [cacheLineSize]byte }








var X86 struct {
	_                   CacheLinePad
	HasAES              bool 
	HasADX              bool 
	HasAVX              bool 
	HasAVX2             bool 
	HasAVX512           bool 
	HasAVX512F          bool 
	HasAVX512CD         bool 
	HasAVX512ER         bool 
	HasAVX512PF         bool 
	HasAVX512VL         bool 
	HasAVX512BW         bool 
	HasAVX512DQ         bool 
	HasAVX512IFMA       bool 
	HasAVX512VBMI       bool 
	HasAVX5124VNNIW     bool 
	HasAVX5124FMAPS     bool 
	HasAVX512VPOPCNTDQ  bool 
	HasAVX512VPCLMULQDQ bool 
	HasAVX512VNNI       bool 
	HasAVX512GFNI       bool 
	HasAVX512VAES       bool 
	HasAVX512VBMI2      bool 
	HasAVX512BITALG     bool 
	HasAVX512BF16       bool 
	HasAMXTile          bool 
	HasAMXInt8          bool 
	HasAMXBF16          bool 
	HasBMI1             bool 
	HasBMI2             bool 
	HasCX16             bool 
	HasERMS             bool 
	HasFMA              bool 
	HasOSXSAVE          bool 
	HasPCLMULQDQ        bool 
	HasPOPCNT           bool 
	HasRDRAND           bool 
	HasRDSEED           bool 
	HasSSE2             bool 
	HasSSE3             bool 
	HasSSSE3            bool 
	HasSSE41            bool 
	HasSSE42            bool 
	HasAVXIFMA          bool 
	HasAVXVNNI          bool 
	HasAVXVNNIInt8      bool 
	_                   CacheLinePad
}




var ARM64 struct {
	_           CacheLinePad
	HasFP       bool 
	HasASIMD    bool 
	HasEVTSTRM  bool 
	HasAES      bool 
	HasPMULL    bool 
	HasSHA1     bool 
	HasSHA2     bool 
	HasCRC32    bool 
	HasATOMICS  bool 
	HasFPHP     bool 
	HasASIMDHP  bool 
	HasCPUID    bool 
	HasASIMDRDM bool 
	HasJSCVT    bool 
	HasFCMA     bool 
	HasLRCPC    bool 
	HasDCPOP    bool 
	HasSHA3     bool 
	HasSM3      bool 
	HasSM4      bool 
	HasASIMDDP  bool 
	HasSHA512   bool 
	HasSVE      bool 
	HasSVE2     bool 
	HasASIMDFHM bool 
	HasDIT      bool 
	HasI8MM     bool 
	_           CacheLinePad
}





var ARM struct {
	_           CacheLinePad
	HasSWP      bool 
	HasHALF     bool 
	HasTHUMB    bool 
	Has26BIT    bool 
	HasFASTMUL  bool 
	HasFPA      bool 
	HasVFP      bool 
	HasEDSP     bool 
	HasJAVA     bool 
	HasIWMMXT   bool 
	HasCRUNCH   bool 
	HasTHUMBEE  bool 
	HasNEON     bool 
	HasVFPv3    bool 
	HasVFPv3D16 bool 
	HasTLS      bool 
	HasVFPv4    bool 
	HasIDIVA    bool 
	HasIDIVT    bool 
	HasVFPD32   bool 
	HasLPAE     bool 
	HasEVTSTRM  bool 
	HasAES      bool 
	HasPMULL    bool 
	HasSHA1     bool 
	HasSHA2     bool 
	HasCRC32    bool 
	_           CacheLinePad
}




var MIPS64X struct {
	_      CacheLinePad
	HasMSA bool 
	_      CacheLinePad
}








var PPC64 struct {
	_        CacheLinePad
	HasDARN  bool 
	HasSCV   bool 
	IsPOWER8 bool 
	IsPOWER9 bool 
	_        CacheLinePad
}








var S390X struct {
	_         CacheLinePad
	HasZARCH  bool 
	HasSTFLE  bool 
	HasLDISP  bool 
	HasEIMM   bool 
	HasDFP    bool 
	HasETF3EH bool 
	HasMSA    bool 
	HasAES    bool 
	HasAESCBC bool 
	HasAESCTR bool 
	HasAESGCM bool 
	HasGHASH  bool 
	HasSHA1   bool 
	HasSHA256 bool 
	HasSHA512 bool 
	HasSHA3   bool 
	HasVX     bool 
	HasVXE    bool 
	_         CacheLinePad
}









var RISCV64 struct {
	_                 CacheLinePad
	HasFastMisaligned bool 
	HasC              bool 
	HasV              bool 
	HasZba            bool 
	HasZbb            bool 
	HasZbs            bool 
	_                 CacheLinePad
}

func init() {
	archInit()
	initOptions()
	processOptions()
}





var options []option


type option struct {
	Name      string
	Feature   *bool
	Specified bool 
	Enable    bool 
	Required  bool 
}

func processOptions() {
	env := os.Getenv("GODEBUG")
field:
	for env != "" {
		field := ""
		i := strings.IndexByte(env, ',')
		if i < 0 {
			field, env = env, ""
		} else {
			field, env = env[:i], env[i+1:]
		}
		if len(field) < 4 || field[:4] != "cpu." {
			continue
		}
		i = strings.IndexByte(field, '=')
		if i < 0 {
			print("GODEBUG sys/cpu: no value specified for \"", field, "\"\n")
			continue
		}
		key, value := field[4:i], field[i+1:] 

		var enable bool
		switch value {
		case "on":
			enable = true
		case "off":
			enable = false
		default:
			print("GODEBUG sys/cpu: value \"", value, "\" not supported for cpu option \"", key, "\"\n")
			continue field
		}

		if key == "all" {
			for i := range options {
				options[i].Specified = true
				options[i].Enable = enable || options[i].Required
			}
			continue field
		}

		for i := range options {
			if options[i].Name == key {
				options[i].Specified = true
				options[i].Enable = enable
				continue field
			}
		}

		print("GODEBUG sys/cpu: unknown cpu feature \"", key, "\"\n")
	}

	for _, o := range options {
		if !o.Specified {
			continue
		}

		if o.Enable && !*o.Feature {
			print("GODEBUG sys/cpu: can not enable \"", o.Name, "\", missing CPU support\n")
			continue
		}

		if !o.Enable && o.Required {
			print("GODEBUG sys/cpu: can not disable \"", o.Name, "\", required CPU feature\n")
			continue
		}

		*o.Feature = o.Enable
	}
}
