




























package mips

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"fmt"
)

func init() {
	obj.RegisterRegister(obj.RBaseMIPS, REG_LAST+1, rconv)
	obj.RegisterOpcode(obj.ABaseMIPS, Anames)
}

func rconv(r int) string {
	if r == 0 {
		return "NONE"
	}
	if r == REGG {
		
		return "g"
	}
	if REG_R0 <= r && r <= REG_R31 {
		return fmt.Sprintf("R%d", r-REG_R0)
	}
	if REG_F0 <= r && r <= REG_F31 {
		return fmt.Sprintf("F%d", r-REG_F0)
	}
	if REG_M0 <= r && r <= REG_M31 {
		return fmt.Sprintf("M%d", r-REG_M0)
	}
	if REG_FCR0 <= r && r <= REG_FCR31 {
		return fmt.Sprintf("FCR%d", r-REG_FCR0)
	}
	if REG_W0 <= r && r <= REG_W31 {
		return fmt.Sprintf("W%d", r-REG_W0)
	}
	if r == REG_HI {
		return "HI"
	}
	if r == REG_LO {
		return "LO"
	}

	return fmt.Sprintf("Rgok(%d)", r-obj.RBaseMIPS)
}

func DRconv(a int) string {
	s := "C_??"
	if a >= C_NONE && a <= C_NCLASS {
		s = cnames0[a]
	}
	var fp string
	fp += s
	return fp
}
