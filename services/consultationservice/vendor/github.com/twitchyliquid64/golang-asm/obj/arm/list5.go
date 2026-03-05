





























package arm

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"fmt"
)

func init() {
	obj.RegisterRegister(obj.RBaseARM, MAXREG, rconv)
	obj.RegisterOpcode(obj.ABaseARM, Anames)
	obj.RegisterRegisterList(obj.RegListARMLo, obj.RegListARMHi, rlconv)
	obj.RegisterOpSuffix("arm", obj.CConvARM)
}

func rconv(r int) string {
	if r == 0 {
		return "NONE"
	}
	if r == REGG {
		
		return "g"
	}
	if REG_R0 <= r && r <= REG_R15 {
		return fmt.Sprintf("R%d", r-REG_R0)
	}
	if REG_F0 <= r && r <= REG_F15 {
		return fmt.Sprintf("F%d", r-REG_F0)
	}

	switch r {
	case REG_FPSR:
		return "FPSR"

	case REG_FPCR:
		return "FPCR"

	case REG_CPSR:
		return "CPSR"

	case REG_SPSR:
		return "SPSR"

	case REG_MB_SY:
		return "MB_SY"
	case REG_MB_ST:
		return "MB_ST"
	case REG_MB_ISH:
		return "MB_ISH"
	case REG_MB_ISHST:
		return "MB_ISHST"
	case REG_MB_NSH:
		return "MB_NSH"
	case REG_MB_NSHST:
		return "MB_NSHST"
	case REG_MB_OSH:
		return "MB_OSH"
	case REG_MB_OSHST:
		return "MB_OSHST"
	}

	return fmt.Sprintf("Rgok(%d)", r-obj.RBaseARM)
}

func DRconv(a int) string {
	s := "C_??"
	if a >= C_NONE && a <= C_NCLASS {
		s = cnames5[a]
	}
	var fp string
	fp += s
	return fp
}

func rlconv(list int64) string {
	str := ""
	for i := 0; i < 16; i++ {
		if list&(1<<uint(i)) != 0 {
			if str == "" {
				str += "["
			} else {
				str += ","
			}
			
			if i == REGG-REG_R0 {
				str += "g"
			} else {
				str += fmt.Sprintf("R%d", i)
			}
		}
	}

	str += "]"
	return str
}
