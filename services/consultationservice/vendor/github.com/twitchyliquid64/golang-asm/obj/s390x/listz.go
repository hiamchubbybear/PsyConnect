




























package s390x

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"fmt"
)

func init() {
	obj.RegisterRegister(obj.RBaseS390X, REG_R0+1024, rconv)
	obj.RegisterOpcode(obj.ABaseS390X, Anames)
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
	if REG_AR0 <= r && r <= REG_AR15 {
		return fmt.Sprintf("AR%d", r-REG_AR0)
	}
	if REG_V0 <= r && r <= REG_V31 {
		return fmt.Sprintf("V%d", r-REG_V0)
	}
	return fmt.Sprintf("Rgok(%d)", r-obj.RBaseS390X)
}

func DRconv(a int) string {
	s := "C_??"
	if a >= C_NONE && a <= C_NCLASS {
		s = cnamesz[a]
	}
	var fp string
	fp += s
	return fp
}
