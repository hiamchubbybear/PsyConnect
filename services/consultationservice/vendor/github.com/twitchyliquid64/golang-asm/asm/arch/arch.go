




package arch

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/arm"
	"github.com/twitchyliquid64/golang-asm/obj/arm64"
	"github.com/twitchyliquid64/golang-asm/obj/mips"
	"github.com/twitchyliquid64/golang-asm/obj/ppc64"
	"github.com/twitchyliquid64/golang-asm/obj/riscv"
	"github.com/twitchyliquid64/golang-asm/obj/s390x"
	"github.com/twitchyliquid64/golang-asm/obj/wasm"
	"github.com/twitchyliquid64/golang-asm/obj/x86"
	"fmt"
	"strings"
)


const (
	RFP = -(iota + 1)
	RSB
	RSP
	RPC
)


type Arch struct {
	*obj.LinkArch
	
	Instructions map[string]obj.As
	
	Register map[string]int16
	
	RegisterPrefix map[string]bool
	
	RegisterNumber func(string, int16) (int16, bool)
	
	IsJump func(word string) bool
}



func nilRegisterNumber(name string, n int16) (int16, bool) {
	return 0, false
}



func Set(GOARCH string) *Arch {
	switch GOARCH {
	case "386":
		return archX86(&x86.Link386)
	case "amd64":
		return archX86(&x86.Linkamd64)
	case "arm":
		return archArm()
	case "arm64":
		return archArm64()
	case "mips":
		return archMips(&mips.Linkmips)
	case "mipsle":
		return archMips(&mips.Linkmipsle)
	case "mips64":
		return archMips64(&mips.Linkmips64)
	case "mips64le":
		return archMips64(&mips.Linkmips64le)
	case "ppc64":
		return archPPC64(&ppc64.Linkppc64)
	case "ppc64le":
		return archPPC64(&ppc64.Linkppc64le)
	case "riscv64":
		return archRISCV64()
	case "s390x":
		return archS390x()
	case "wasm":
		return archWasm()
	}
	return nil
}

func jumpX86(word string) bool {
	return word[0] == 'J' || word == "CALL" || strings.HasPrefix(word, "LOOP") || word == "XBEGIN"
}

func jumpRISCV(word string) bool {
	switch word {
	case "BEQ", "BEQZ", "BGE", "BGEU", "BGEZ", "BGT", "BGTU", "BGTZ", "BLE", "BLEU", "BLEZ",
		"BLT", "BLTU", "BLTZ", "BNE", "BNEZ", "CALL", "JAL", "JALR", "JMP":
		return true
	}
	return false
}

func jumpWasm(word string) bool {
	return word == "JMP" || word == "CALL" || word == "Call" || word == "Br" || word == "BrIf"
}

func archX86(linkArch *obj.LinkArch) *Arch {
	register := make(map[string]int16)
	
	for i, s := range x86.Register {
		register[s] = int16(i + x86.REG_AL)
	}
	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range x86.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseAMD64
		}
	}
	
	instructions["JA"] = x86.AJHI   
	instructions["JAE"] = x86.AJCC  
	instructions["JB"] = x86.AJCS   
	instructions["JBE"] = x86.AJLS  
	instructions["JC"] = x86.AJCS   
	instructions["JCC"] = x86.AJCC  
	instructions["JCS"] = x86.AJCS  
	instructions["JE"] = x86.AJEQ   
	instructions["JEQ"] = x86.AJEQ  
	instructions["JG"] = x86.AJGT   
	instructions["JGE"] = x86.AJGE  
	instructions["JGT"] = x86.AJGT  
	instructions["JHI"] = x86.AJHI  
	instructions["JHS"] = x86.AJCC  
	instructions["JL"] = x86.AJLT   
	instructions["JLE"] = x86.AJLE  
	instructions["JLO"] = x86.AJCS  
	instructions["JLS"] = x86.AJLS  
	instructions["JLT"] = x86.AJLT  
	instructions["JMI"] = x86.AJMI  
	instructions["JNA"] = x86.AJLS  
	instructions["JNAE"] = x86.AJCS 
	instructions["JNB"] = x86.AJCC  
	instructions["JNBE"] = x86.AJHI 
	instructions["JNC"] = x86.AJCC  
	instructions["JNE"] = x86.AJNE  
	instructions["JNG"] = x86.AJLE  
	instructions["JNGE"] = x86.AJLT 
	instructions["JNL"] = x86.AJGE  
	instructions["JNLE"] = x86.AJGT 
	instructions["JNO"] = x86.AJOC  
	instructions["JNP"] = x86.AJPC  
	instructions["JNS"] = x86.AJPL  
	instructions["JNZ"] = x86.AJNE  
	instructions["JO"] = x86.AJOS   
	instructions["JOC"] = x86.AJOC  
	instructions["JOS"] = x86.AJOS  
	instructions["JP"] = x86.AJPS   
	instructions["JPC"] = x86.AJPC  
	instructions["JPE"] = x86.AJPS  
	instructions["JPL"] = x86.AJPL  
	instructions["JPO"] = x86.AJPC  
	instructions["JPS"] = x86.AJPS  
	instructions["JS"] = x86.AJMI   
	instructions["JZ"] = x86.AJEQ   
	instructions["MASKMOVDQU"] = x86.AMASKMOVOU
	instructions["MOVD"] = x86.AMOVQ
	instructions["MOVDQ2Q"] = x86.AMOVQ
	instructions["MOVNTDQ"] = x86.AMOVNTO
	instructions["MOVOA"] = x86.AMOVO
	instructions["PSLLDQ"] = x86.APSLLO
	instructions["PSRLDQ"] = x86.APSRLO
	instructions["PADDD"] = x86.APADDL

	return &Arch{
		LinkArch:       linkArch,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: nil,
		RegisterNumber: nilRegisterNumber,
		IsJump:         jumpX86,
	}
}

func archArm() *Arch {
	register := make(map[string]int16)
	
	
	for i := arm.REG_R0; i < arm.REG_SPSR; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	
	delete(register, "R10")
	register["g"] = arm.REG_R10
	for i := 0; i < 16; i++ {
		register[fmt.Sprintf("C%d", i)] = int16(i)
	}

	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	register["SP"] = RSP
	registerPrefix := map[string]bool{
		"F": true,
		"R": true,
	}

	
	register["MB_SY"] = arm.REG_MB_SY
	register["MB_ST"] = arm.REG_MB_ST
	register["MB_ISH"] = arm.REG_MB_ISH
	register["MB_ISHST"] = arm.REG_MB_ISHST
	register["MB_NSH"] = arm.REG_MB_NSH
	register["MB_NSHST"] = arm.REG_MB_NSHST
	register["MB_OSH"] = arm.REG_MB_OSH
	register["MB_OSHST"] = arm.REG_MB_OSHST

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range arm.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseARM
		}
	}
	
	instructions["B"] = obj.AJMP
	instructions["BL"] = obj.ACALL
	
	
	
	instructions["MCR"] = aMCR

	return &Arch{
		LinkArch:       &arm.Linkarm,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: armRegisterNumber,
		IsJump:         jumpArm,
	}
}

func archArm64() *Arch {
	register := make(map[string]int16)
	
	
	register[obj.Rconv(arm64.REGSP)] = int16(arm64.REGSP)
	for i := arm64.REG_R0; i <= arm64.REG_R31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	
	register["R18_PLATFORM"] = register["R18"]
	delete(register, "R18")
	for i := arm64.REG_F0; i <= arm64.REG_F31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := arm64.REG_V0; i <= arm64.REG_V31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}

	
	for i := 0; i < len(arm64.SystemReg); i++ {
		register[arm64.SystemReg[i].Name] = arm64.SystemReg[i].Reg
	}

	register["LR"] = arm64.REGLINK
	register["DAIFSet"] = arm64.REG_DAIFSet
	register["DAIFClr"] = arm64.REG_DAIFClr
	register["PLDL1KEEP"] = arm64.REG_PLDL1KEEP
	register["PLDL1STRM"] = arm64.REG_PLDL1STRM
	register["PLDL2KEEP"] = arm64.REG_PLDL2KEEP
	register["PLDL2STRM"] = arm64.REG_PLDL2STRM
	register["PLDL3KEEP"] = arm64.REG_PLDL3KEEP
	register["PLDL3STRM"] = arm64.REG_PLDL3STRM
	register["PLIL1KEEP"] = arm64.REG_PLIL1KEEP
	register["PLIL1STRM"] = arm64.REG_PLIL1STRM
	register["PLIL2KEEP"] = arm64.REG_PLIL2KEEP
	register["PLIL2STRM"] = arm64.REG_PLIL2STRM
	register["PLIL3KEEP"] = arm64.REG_PLIL3KEEP
	register["PLIL3STRM"] = arm64.REG_PLIL3STRM
	register["PSTL1KEEP"] = arm64.REG_PSTL1KEEP
	register["PSTL1STRM"] = arm64.REG_PSTL1STRM
	register["PSTL2KEEP"] = arm64.REG_PSTL2KEEP
	register["PSTL2STRM"] = arm64.REG_PSTL2STRM
	register["PSTL3KEEP"] = arm64.REG_PSTL3KEEP
	register["PSTL3STRM"] = arm64.REG_PSTL3STRM

	
	register["EQ"] = arm64.COND_EQ
	register["NE"] = arm64.COND_NE
	register["HS"] = arm64.COND_HS
	register["CS"] = arm64.COND_HS
	register["LO"] = arm64.COND_LO
	register["CC"] = arm64.COND_LO
	register["MI"] = arm64.COND_MI
	register["PL"] = arm64.COND_PL
	register["VS"] = arm64.COND_VS
	register["VC"] = arm64.COND_VC
	register["HI"] = arm64.COND_HI
	register["LS"] = arm64.COND_LS
	register["GE"] = arm64.COND_GE
	register["LT"] = arm64.COND_LT
	register["GT"] = arm64.COND_GT
	register["LE"] = arm64.COND_LE
	register["AL"] = arm64.COND_AL
	register["NV"] = arm64.COND_NV
	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	register["SP"] = RSP
	
	delete(register, "R28")
	register["g"] = arm64.REG_R28
	registerPrefix := map[string]bool{
		"F": true,
		"R": true,
		"V": true,
	}

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range arm64.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseARM64
		}
	}
	
	instructions["B"] = arm64.AB
	instructions["BL"] = arm64.ABL

	return &Arch{
		LinkArch:       &arm64.Linkarm64,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: arm64RegisterNumber,
		IsJump:         jumpArm64,
	}

}

func archPPC64(linkArch *obj.LinkArch) *Arch {
	register := make(map[string]int16)
	
	
	for i := ppc64.REG_R0; i <= ppc64.REG_R31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := ppc64.REG_F0; i <= ppc64.REG_F31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := ppc64.REG_V0; i <= ppc64.REG_V31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := ppc64.REG_VS0; i <= ppc64.REG_VS63; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := ppc64.REG_CR0; i <= ppc64.REG_CR7; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := ppc64.REG_MSR; i <= ppc64.REG_CR; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	register["CR"] = ppc64.REG_CR
	register["XER"] = ppc64.REG_XER
	register["LR"] = ppc64.REG_LR
	register["CTR"] = ppc64.REG_CTR
	register["FPSCR"] = ppc64.REG_FPSCR
	register["MSR"] = ppc64.REG_MSR
	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	
	delete(register, "R30")
	register["g"] = ppc64.REG_R30
	registerPrefix := map[string]bool{
		"CR":  true,
		"F":   true,
		"R":   true,
		"SPR": true,
	}

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range ppc64.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABasePPC64
		}
	}
	
	instructions["BR"] = ppc64.ABR
	instructions["BL"] = ppc64.ABL

	return &Arch{
		LinkArch:       linkArch,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: ppc64RegisterNumber,
		IsJump:         jumpPPC64,
	}
}

func archMips(linkArch *obj.LinkArch) *Arch {
	register := make(map[string]int16)
	
	
	for i := mips.REG_R0; i <= mips.REG_R31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}

	for i := mips.REG_F0; i <= mips.REG_F31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := mips.REG_M0; i <= mips.REG_M31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := mips.REG_FCR0; i <= mips.REG_FCR31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	register["HI"] = mips.REG_HI
	register["LO"] = mips.REG_LO
	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	
	delete(register, "R30")
	register["g"] = mips.REG_R30

	registerPrefix := map[string]bool{
		"F":   true,
		"FCR": true,
		"M":   true,
		"R":   true,
	}

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range mips.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseMIPS
		}
	}
	
	instructions["JAL"] = mips.AJAL

	return &Arch{
		LinkArch:       linkArch,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: mipsRegisterNumber,
		IsJump:         jumpMIPS,
	}
}

func archMips64(linkArch *obj.LinkArch) *Arch {
	register := make(map[string]int16)
	
	
	for i := mips.REG_R0; i <= mips.REG_R31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := mips.REG_F0; i <= mips.REG_F31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := mips.REG_M0; i <= mips.REG_M31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := mips.REG_FCR0; i <= mips.REG_FCR31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := mips.REG_W0; i <= mips.REG_W31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	register["HI"] = mips.REG_HI
	register["LO"] = mips.REG_LO
	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	
	delete(register, "R30")
	register["g"] = mips.REG_R30
	
	delete(register, "R28")
	register["RSB"] = mips.REG_R28
	registerPrefix := map[string]bool{
		"F":   true,
		"FCR": true,
		"M":   true,
		"R":   true,
		"W":   true,
	}

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range mips.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseMIPS
		}
	}
	
	instructions["JAL"] = mips.AJAL

	return &Arch{
		LinkArch:       linkArch,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: mipsRegisterNumber,
		IsJump:         jumpMIPS,
	}
}

func archRISCV64() *Arch {
	register := make(map[string]int16)

	
	for i := riscv.REG_X0; i <= riscv.REG_X31; i++ {
		name := fmt.Sprintf("X%d", i-riscv.REG_X0)
		register[name] = int16(i)
	}
	for i := riscv.REG_F0; i <= riscv.REG_F31; i++ {
		name := fmt.Sprintf("F%d", i-riscv.REG_F0)
		register[name] = int16(i)
	}

	
	register["ZERO"] = riscv.REG_ZERO
	register["RA"] = riscv.REG_RA
	register["SP"] = riscv.REG_SP
	register["GP"] = riscv.REG_GP
	register["TP"] = riscv.REG_TP
	register["T0"] = riscv.REG_T0
	register["T1"] = riscv.REG_T1
	register["T2"] = riscv.REG_T2
	register["S0"] = riscv.REG_S0
	register["S1"] = riscv.REG_S1
	register["A0"] = riscv.REG_A0
	register["A1"] = riscv.REG_A1
	register["A2"] = riscv.REG_A2
	register["A3"] = riscv.REG_A3
	register["A4"] = riscv.REG_A4
	register["A5"] = riscv.REG_A5
	register["A6"] = riscv.REG_A6
	register["A7"] = riscv.REG_A7
	register["S2"] = riscv.REG_S2
	register["S3"] = riscv.REG_S3
	register["S4"] = riscv.REG_S4
	register["S5"] = riscv.REG_S5
	register["S6"] = riscv.REG_S6
	register["S7"] = riscv.REG_S7
	register["S8"] = riscv.REG_S8
	register["S9"] = riscv.REG_S9
	register["S10"] = riscv.REG_S10
	register["S11"] = riscv.REG_S11
	register["T3"] = riscv.REG_T3
	register["T4"] = riscv.REG_T4
	register["T5"] = riscv.REG_T5
	register["T6"] = riscv.REG_T6

	
	register["g"] = riscv.REG_G
	register["CTXT"] = riscv.REG_CTXT
	register["TMP"] = riscv.REG_TMP

	
	register["FT0"] = riscv.REG_FT0
	register["FT1"] = riscv.REG_FT1
	register["FT2"] = riscv.REG_FT2
	register["FT3"] = riscv.REG_FT3
	register["FT4"] = riscv.REG_FT4
	register["FT5"] = riscv.REG_FT5
	register["FT6"] = riscv.REG_FT6
	register["FT7"] = riscv.REG_FT7
	register["FS0"] = riscv.REG_FS0
	register["FS1"] = riscv.REG_FS1
	register["FA0"] = riscv.REG_FA0
	register["FA1"] = riscv.REG_FA1
	register["FA2"] = riscv.REG_FA2
	register["FA3"] = riscv.REG_FA3
	register["FA4"] = riscv.REG_FA4
	register["FA5"] = riscv.REG_FA5
	register["FA6"] = riscv.REG_FA6
	register["FA7"] = riscv.REG_FA7
	register["FS2"] = riscv.REG_FS2
	register["FS3"] = riscv.REG_FS3
	register["FS4"] = riscv.REG_FS4
	register["FS5"] = riscv.REG_FS5
	register["FS6"] = riscv.REG_FS6
	register["FS7"] = riscv.REG_FS7
	register["FS8"] = riscv.REG_FS8
	register["FS9"] = riscv.REG_FS9
	register["FS10"] = riscv.REG_FS10
	register["FS11"] = riscv.REG_FS11
	register["FT8"] = riscv.REG_FT8
	register["FT9"] = riscv.REG_FT9
	register["FT10"] = riscv.REG_FT10
	register["FT11"] = riscv.REG_FT11

	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range riscv.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseRISCV
		}
	}

	return &Arch{
		LinkArch:       &riscv.LinkRISCV64,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: nil,
		RegisterNumber: nilRegisterNumber,
		IsJump:         jumpRISCV,
	}
}

func archS390x() *Arch {
	register := make(map[string]int16)
	
	
	for i := s390x.REG_R0; i <= s390x.REG_R15; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := s390x.REG_F0; i <= s390x.REG_F15; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := s390x.REG_V0; i <= s390x.REG_V31; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	for i := s390x.REG_AR0; i <= s390x.REG_AR15; i++ {
		register[obj.Rconv(i)] = int16(i)
	}
	register["LR"] = s390x.REG_LR
	
	register["SB"] = RSB
	register["FP"] = RFP
	register["PC"] = RPC
	
	delete(register, "R13")
	register["g"] = s390x.REG_R13
	registerPrefix := map[string]bool{
		"AR": true,
		"F":  true,
		"R":  true,
	}

	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range s390x.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseS390X
		}
	}
	
	instructions["BR"] = s390x.ABR
	instructions["BL"] = s390x.ABL

	return &Arch{
		LinkArch:       &s390x.Links390x,
		Instructions:   instructions,
		Register:       register,
		RegisterPrefix: registerPrefix,
		RegisterNumber: s390xRegisterNumber,
		IsJump:         jumpS390x,
	}
}

func archWasm() *Arch {
	instructions := make(map[string]obj.As)
	for i, s := range obj.Anames {
		instructions[s] = obj.As(i)
	}
	for i, s := range wasm.Anames {
		if obj.As(i) >= obj.A_ARCHSPECIFIC {
			instructions[s] = obj.As(i) + obj.ABaseWasm
		}
	}

	return &Arch{
		LinkArch:       &wasm.Linkwasm,
		Instructions:   instructions,
		Register:       wasm.Register,
		RegisterPrefix: nil,
		RegisterNumber: nilRegisterNumber,
		IsJump:         jumpWasm,
	}
}
