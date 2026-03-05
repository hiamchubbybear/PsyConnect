




























package ppc64

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/objabi"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"sort"
)




type ctxt9 struct {
	ctxt       *obj.Link
	newprog    obj.ProgAlloc
	cursym     *obj.LSym
	autosize   int32
	instoffset int64
	pc         int64
}



const (
	funcAlign     = 16
	funcAlignMask = funcAlign - 1
)

const (
	r0iszero = 1
)

type Optab struct {
	as    obj.As 
	a1    uint8
	a2    uint8
	a3    uint8
	a4    uint8
	type_ int8 
	size  int8
	param int16
}









var optab = []Optab{
	{obj.ATEXT, C_LEXT, C_NONE, C_NONE, C_TEXTSIZE, 0, 0, 0},
	{obj.ATEXT, C_LEXT, C_NONE, C_LCON, C_TEXTSIZE, 0, 0, 0},
	{obj.ATEXT, C_ADDR, C_NONE, C_NONE, C_TEXTSIZE, 0, 0, 0},
	{obj.ATEXT, C_ADDR, C_NONE, C_LCON, C_TEXTSIZE, 0, 0, 0},
	
	{AMOVD, C_REG, C_NONE, C_NONE, C_REG, 1, 4, 0},
	{AMOVB, C_REG, C_NONE, C_NONE, C_REG, 12, 4, 0},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_REG, 13, 4, 0},
	{AMOVW, C_REG, C_NONE, C_NONE, C_REG, 12, 4, 0},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_REG, 13, 4, 0},
	{AADD, C_REG, C_REG, C_NONE, C_REG, 2, 4, 0},
	{AADD, C_REG, C_NONE, C_NONE, C_REG, 2, 4, 0},
	{AADD, C_SCON, C_REG, C_NONE, C_REG, 4, 4, 0},
	{AADD, C_SCON, C_NONE, C_NONE, C_REG, 4, 4, 0},
	{AADD, C_ADDCON, C_REG, C_NONE, C_REG, 4, 4, 0},
	{AADD, C_ADDCON, C_NONE, C_NONE, C_REG, 4, 4, 0},
	{AADD, C_UCON, C_REG, C_NONE, C_REG, 20, 4, 0},
	{AADD, C_UCON, C_NONE, C_NONE, C_REG, 20, 4, 0},
	{AADD, C_ANDCON, C_REG, C_NONE, C_REG, 22, 8, 0},
	{AADD, C_ANDCON, C_NONE, C_NONE, C_REG, 22, 8, 0},
	{AADD, C_LCON, C_REG, C_NONE, C_REG, 22, 12, 0},
	{AADD, C_LCON, C_NONE, C_NONE, C_REG, 22, 12, 0},
	{AADDIS, C_ADDCON, C_REG, C_NONE, C_REG, 20, 4, 0},
	{AADDIS, C_ADDCON, C_NONE, C_NONE, C_REG, 20, 4, 0},
	{AADDC, C_REG, C_REG, C_NONE, C_REG, 2, 4, 0},
	{AADDC, C_REG, C_NONE, C_NONE, C_REG, 2, 4, 0},
	{AADDC, C_ADDCON, C_REG, C_NONE, C_REG, 4, 4, 0},
	{AADDC, C_ADDCON, C_NONE, C_NONE, C_REG, 4, 4, 0},
	{AADDC, C_LCON, C_REG, C_NONE, C_REG, 22, 12, 0},
	{AADDC, C_LCON, C_NONE, C_NONE, C_REG, 22, 12, 0},
	{AAND, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0}, 
	{AAND, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{AANDCC, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0},
	{AANDCC, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{AANDCC, C_ANDCON, C_NONE, C_NONE, C_REG, 58, 4, 0},
	{AANDCC, C_ANDCON, C_REG, C_NONE, C_REG, 58, 4, 0},
	{AANDCC, C_UCON, C_NONE, C_NONE, C_REG, 59, 4, 0},
	{AANDCC, C_UCON, C_REG, C_NONE, C_REG, 59, 4, 0},
	{AANDCC, C_ADDCON, C_NONE, C_NONE, C_REG, 23, 8, 0},
	{AANDCC, C_ADDCON, C_REG, C_NONE, C_REG, 23, 8, 0},
	{AANDCC, C_LCON, C_NONE, C_NONE, C_REG, 23, 12, 0},
	{AANDCC, C_LCON, C_REG, C_NONE, C_REG, 23, 12, 0},
	{AANDISCC, C_ANDCON, C_NONE, C_NONE, C_REG, 59, 4, 0},
	{AANDISCC, C_ANDCON, C_REG, C_NONE, C_REG, 59, 4, 0},
	{AMULLW, C_REG, C_REG, C_NONE, C_REG, 2, 4, 0},
	{AMULLW, C_REG, C_NONE, C_NONE, C_REG, 2, 4, 0},
	{AMULLW, C_ADDCON, C_REG, C_NONE, C_REG, 4, 4, 0},
	{AMULLW, C_ADDCON, C_NONE, C_NONE, C_REG, 4, 4, 0},
	{AMULLW, C_ANDCON, C_REG, C_NONE, C_REG, 4, 4, 0},
	{AMULLW, C_ANDCON, C_NONE, C_NONE, C_REG, 4, 4, 0},
	{AMULLW, C_LCON, C_REG, C_NONE, C_REG, 22, 12, 0},
	{AMULLW, C_LCON, C_NONE, C_NONE, C_REG, 22, 12, 0},
	{ASUBC, C_REG, C_REG, C_NONE, C_REG, 10, 4, 0},
	{ASUBC, C_REG, C_NONE, C_NONE, C_REG, 10, 4, 0},
	{ASUBC, C_REG, C_NONE, C_ADDCON, C_REG, 27, 4, 0},
	{ASUBC, C_REG, C_NONE, C_LCON, C_REG, 28, 12, 0},
	{AOR, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0}, 
	{AOR, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{AOR, C_ANDCON, C_NONE, C_NONE, C_REG, 58, 4, 0},
	{AOR, C_ANDCON, C_REG, C_NONE, C_REG, 58, 4, 0},
	{AOR, C_UCON, C_NONE, C_NONE, C_REG, 59, 4, 0},
	{AOR, C_UCON, C_REG, C_NONE, C_REG, 59, 4, 0},
	{AOR, C_ADDCON, C_NONE, C_NONE, C_REG, 23, 8, 0},
	{AOR, C_ADDCON, C_REG, C_NONE, C_REG, 23, 8, 0},
	{AOR, C_LCON, C_NONE, C_NONE, C_REG, 23, 12, 0},
	{AOR, C_LCON, C_REG, C_NONE, C_REG, 23, 12, 0},
	{AORIS, C_ANDCON, C_NONE, C_NONE, C_REG, 59, 4, 0},
	{AORIS, C_ANDCON, C_REG, C_NONE, C_REG, 59, 4, 0},
	{ADIVW, C_REG, C_REG, C_NONE, C_REG, 2, 4, 0}, 
	{ADIVW, C_REG, C_NONE, C_NONE, C_REG, 2, 4, 0},
	{ASUB, C_REG, C_REG, C_NONE, C_REG, 10, 4, 0}, 
	{ASUB, C_REG, C_NONE, C_NONE, C_REG, 10, 4, 0},
	{ASLW, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{ASLW, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0},
	{ASLD, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{ASLD, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0},
	{ASLD, C_SCON, C_REG, C_NONE, C_REG, 25, 4, 0},
	{ASLD, C_SCON, C_NONE, C_NONE, C_REG, 25, 4, 0},
	{ASLW, C_SCON, C_REG, C_NONE, C_REG, 57, 4, 0},
	{ASLW, C_SCON, C_NONE, C_NONE, C_REG, 57, 4, 0},
	{ASRAW, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{ASRAW, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0},
	{ASRAW, C_SCON, C_REG, C_NONE, C_REG, 56, 4, 0},
	{ASRAW, C_SCON, C_NONE, C_NONE, C_REG, 56, 4, 0},
	{ASRAD, C_REG, C_NONE, C_NONE, C_REG, 6, 4, 0},
	{ASRAD, C_REG, C_REG, C_NONE, C_REG, 6, 4, 0},
	{ASRAD, C_SCON, C_REG, C_NONE, C_REG, 56, 4, 0},
	{ASRAD, C_SCON, C_NONE, C_NONE, C_REG, 56, 4, 0},
	{ARLWMI, C_SCON, C_REG, C_LCON, C_REG, 62, 4, 0},
	{ARLWMI, C_REG, C_REG, C_LCON, C_REG, 63, 4, 0},
	{ARLDMI, C_SCON, C_REG, C_LCON, C_REG, 30, 4, 0},
	{ARLDC, C_SCON, C_REG, C_LCON, C_REG, 29, 4, 0},
	{ARLDCL, C_SCON, C_REG, C_LCON, C_REG, 29, 4, 0},
	{ARLDCL, C_REG, C_REG, C_LCON, C_REG, 14, 4, 0},
	{ARLDICL, C_REG, C_REG, C_LCON, C_REG, 14, 4, 0},
	{ARLDICL, C_SCON, C_REG, C_LCON, C_REG, 14, 4, 0},
	{ARLDCL, C_REG, C_NONE, C_LCON, C_REG, 14, 4, 0},
	{AFADD, C_FREG, C_NONE, C_NONE, C_FREG, 2, 4, 0},
	{AFADD, C_FREG, C_FREG, C_NONE, C_FREG, 2, 4, 0},
	{AFABS, C_FREG, C_NONE, C_NONE, C_FREG, 33, 4, 0},
	{AFABS, C_NONE, C_NONE, C_NONE, C_FREG, 33, 4, 0},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_FREG, 33, 4, 0},
	{AFMADD, C_FREG, C_FREG, C_FREG, C_FREG, 34, 4, 0},
	{AFMUL, C_FREG, C_NONE, C_NONE, C_FREG, 32, 4, 0},
	{AFMUL, C_FREG, C_FREG, C_NONE, C_FREG, 32, 4, 0},

	
	{AMOVD, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVW, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVWZ, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVBZ, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVBZU, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVB, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVBU, C_REG, C_REG, C_NONE, C_ZOREG, 7, 4, REGZERO},
	{AMOVD, C_REG, C_NONE, C_NONE, C_SEXT, 7, 4, REGSB},
	{AMOVW, C_REG, C_NONE, C_NONE, C_SEXT, 7, 4, REGSB},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_SEXT, 7, 4, REGSB},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_SEXT, 7, 4, REGSB},
	{AMOVB, C_REG, C_NONE, C_NONE, C_SEXT, 7, 4, REGSB},
	{AMOVD, C_REG, C_NONE, C_NONE, C_SAUTO, 7, 4, REGSP},
	{AMOVW, C_REG, C_NONE, C_NONE, C_SAUTO, 7, 4, REGSP},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_SAUTO, 7, 4, REGSP},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_SAUTO, 7, 4, REGSP},
	{AMOVB, C_REG, C_NONE, C_NONE, C_SAUTO, 7, 4, REGSP},
	{AMOVD, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AMOVW, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AMOVBZU, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AMOVB, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AMOVBU, C_REG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},

	
	{AMOVD, C_ZOREG, C_REG, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVW, C_ZOREG, C_REG, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVWZ, C_ZOREG, C_REG, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVBZ, C_ZOREG, C_REG, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVBZU, C_ZOREG, C_REG, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVB, C_ZOREG, C_REG, C_NONE, C_REG, 9, 8, REGZERO},
	{AMOVBU, C_ZOREG, C_REG, C_NONE, C_REG, 9, 8, REGZERO},
	{AMOVD, C_SEXT, C_NONE, C_NONE, C_REG, 8, 4, REGSB},
	{AMOVW, C_SEXT, C_NONE, C_NONE, C_REG, 8, 4, REGSB},
	{AMOVWZ, C_SEXT, C_NONE, C_NONE, C_REG, 8, 4, REGSB},
	{AMOVBZ, C_SEXT, C_NONE, C_NONE, C_REG, 8, 4, REGSB},
	{AMOVB, C_SEXT, C_NONE, C_NONE, C_REG, 9, 8, REGSB},
	{AMOVD, C_SAUTO, C_NONE, C_NONE, C_REG, 8, 4, REGSP},
	{AMOVW, C_SAUTO, C_NONE, C_NONE, C_REG, 8, 4, REGSP},
	{AMOVWZ, C_SAUTO, C_NONE, C_NONE, C_REG, 8, 4, REGSP},
	{AMOVBZ, C_SAUTO, C_NONE, C_NONE, C_REG, 8, 4, REGSP},
	{AMOVB, C_SAUTO, C_NONE, C_NONE, C_REG, 9, 8, REGSP},
	{AMOVD, C_SOREG, C_NONE, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVW, C_SOREG, C_NONE, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVWZ, C_SOREG, C_NONE, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVBZ, C_SOREG, C_NONE, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVBZU, C_SOREG, C_NONE, C_NONE, C_REG, 8, 4, REGZERO},
	{AMOVB, C_SOREG, C_NONE, C_NONE, C_REG, 9, 8, REGZERO},
	{AMOVBU, C_SOREG, C_NONE, C_NONE, C_REG, 9, 8, REGZERO},

	
	{AMOVD, C_REG, C_NONE, C_NONE, C_LEXT, 35, 8, REGSB},
	{AMOVW, C_REG, C_NONE, C_NONE, C_LEXT, 35, 8, REGSB},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_LEXT, 35, 8, REGSB},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_LEXT, 35, 8, REGSB},
	{AMOVB, C_REG, C_NONE, C_NONE, C_LEXT, 35, 8, REGSB},
	{AMOVD, C_REG, C_NONE, C_NONE, C_LAUTO, 35, 8, REGSP},
	{AMOVW, C_REG, C_NONE, C_NONE, C_LAUTO, 35, 8, REGSP},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_LAUTO, 35, 8, REGSP},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_LAUTO, 35, 8, REGSP},
	{AMOVB, C_REG, C_NONE, C_NONE, C_LAUTO, 35, 8, REGSP},
	{AMOVD, C_REG, C_NONE, C_NONE, C_LOREG, 35, 8, REGZERO},
	{AMOVW, C_REG, C_NONE, C_NONE, C_LOREG, 35, 8, REGZERO},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_LOREG, 35, 8, REGZERO},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_LOREG, 35, 8, REGZERO},
	{AMOVB, C_REG, C_NONE, C_NONE, C_LOREG, 35, 8, REGZERO},
	{AMOVD, C_REG, C_NONE, C_NONE, C_ADDR, 74, 8, 0},
	{AMOVW, C_REG, C_NONE, C_NONE, C_ADDR, 74, 8, 0},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_ADDR, 74, 8, 0},
	{AMOVBZ, C_REG, C_NONE, C_NONE, C_ADDR, 74, 8, 0},
	{AMOVB, C_REG, C_NONE, C_NONE, C_ADDR, 74, 8, 0},

	
	{AMOVD, C_LEXT, C_NONE, C_NONE, C_REG, 36, 8, REGSB},
	{AMOVW, C_LEXT, C_NONE, C_NONE, C_REG, 36, 8, REGSB},
	{AMOVWZ, C_LEXT, C_NONE, C_NONE, C_REG, 36, 8, REGSB},
	{AMOVBZ, C_LEXT, C_NONE, C_NONE, C_REG, 36, 8, REGSB},
	{AMOVB, C_LEXT, C_NONE, C_NONE, C_REG, 37, 12, REGSB},
	{AMOVD, C_LAUTO, C_NONE, C_NONE, C_REG, 36, 8, REGSP},
	{AMOVW, C_LAUTO, C_NONE, C_NONE, C_REG, 36, 8, REGSP},
	{AMOVWZ, C_LAUTO, C_NONE, C_NONE, C_REG, 36, 8, REGSP},
	{AMOVBZ, C_LAUTO, C_NONE, C_NONE, C_REG, 36, 8, REGSP},
	{AMOVB, C_LAUTO, C_NONE, C_NONE, C_REG, 37, 12, REGSP},
	{AMOVD, C_LOREG, C_NONE, C_NONE, C_REG, 36, 8, REGZERO},
	{AMOVW, C_LOREG, C_NONE, C_NONE, C_REG, 36, 8, REGZERO},
	{AMOVWZ, C_LOREG, C_NONE, C_NONE, C_REG, 36, 8, REGZERO},
	{AMOVBZ, C_LOREG, C_NONE, C_NONE, C_REG, 36, 8, REGZERO},
	{AMOVB, C_LOREG, C_NONE, C_NONE, C_REG, 37, 12, REGZERO},
	{AMOVD, C_ADDR, C_NONE, C_NONE, C_REG, 75, 8, 0},
	{AMOVW, C_ADDR, C_NONE, C_NONE, C_REG, 75, 8, 0},
	{AMOVWZ, C_ADDR, C_NONE, C_NONE, C_REG, 75, 8, 0},
	{AMOVBZ, C_ADDR, C_NONE, C_NONE, C_REG, 75, 8, 0},
	{AMOVB, C_ADDR, C_NONE, C_NONE, C_REG, 76, 12, 0},

	{AMOVD, C_TLS_LE, C_NONE, C_NONE, C_REG, 79, 4, 0},
	{AMOVD, C_TLS_IE, C_NONE, C_NONE, C_REG, 80, 8, 0},

	{AMOVD, C_GOTADDR, C_NONE, C_NONE, C_REG, 81, 8, 0},
	{AMOVD, C_TOCADDR, C_NONE, C_NONE, C_REG, 95, 8, 0},

	
	{AMOVD, C_SECON, C_NONE, C_NONE, C_REG, 3, 4, REGSB},
	{AMOVD, C_SACON, C_NONE, C_NONE, C_REG, 3, 4, REGSP},
	{AMOVD, C_LECON, C_NONE, C_NONE, C_REG, 26, 8, REGSB},
	{AMOVD, C_LACON, C_NONE, C_NONE, C_REG, 26, 8, REGSP},
	{AMOVD, C_ADDCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVD, C_ANDCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVW, C_SECON, C_NONE, C_NONE, C_REG, 3, 4, REGSB}, 
	{AMOVW, C_SACON, C_NONE, C_NONE, C_REG, 3, 4, REGSP},
	{AMOVW, C_LECON, C_NONE, C_NONE, C_REG, 26, 8, REGSB},
	{AMOVW, C_LACON, C_NONE, C_NONE, C_REG, 26, 8, REGSP},
	{AMOVW, C_ADDCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVW, C_ANDCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVWZ, C_SECON, C_NONE, C_NONE, C_REG, 3, 4, REGSB}, 
	{AMOVWZ, C_SACON, C_NONE, C_NONE, C_REG, 3, 4, REGSP},
	{AMOVWZ, C_LECON, C_NONE, C_NONE, C_REG, 26, 8, REGSB},
	{AMOVWZ, C_LACON, C_NONE, C_NONE, C_REG, 26, 8, REGSP},
	{AMOVWZ, C_ADDCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVWZ, C_ANDCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},

	
	{AMOVD, C_UCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVD, C_LCON, C_NONE, C_NONE, C_REG, 19, 8, 0},
	{AMOVW, C_UCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVW, C_LCON, C_NONE, C_NONE, C_REG, 19, 8, 0},
	{AMOVWZ, C_UCON, C_NONE, C_NONE, C_REG, 3, 4, REGZERO},
	{AMOVWZ, C_LCON, C_NONE, C_NONE, C_REG, 19, 8, 0},
	{AMOVHBR, C_ZOREG, C_REG, C_NONE, C_REG, 45, 4, 0},
	{AMOVHBR, C_ZOREG, C_NONE, C_NONE, C_REG, 45, 4, 0},
	{AMOVHBR, C_REG, C_REG, C_NONE, C_ZOREG, 44, 4, 0},
	{AMOVHBR, C_REG, C_NONE, C_NONE, C_ZOREG, 44, 4, 0},
	{ASYSCALL, C_NONE, C_NONE, C_NONE, C_NONE, 5, 4, 0},
	{ASYSCALL, C_REG, C_NONE, C_NONE, C_NONE, 77, 12, 0},
	{ASYSCALL, C_SCON, C_NONE, C_NONE, C_NONE, 77, 12, 0},
	{ABEQ, C_NONE, C_NONE, C_NONE, C_SBRA, 16, 4, 0},
	{ABEQ, C_CREG, C_NONE, C_NONE, C_SBRA, 16, 4, 0},
	{ABR, C_NONE, C_NONE, C_NONE, C_LBRA, 11, 4, 0},
	{ABR, C_NONE, C_NONE, C_NONE, C_LBRAPIC, 11, 8, 0},
	{ABC, C_SCON, C_REG, C_NONE, C_SBRA, 16, 4, 0},
	{ABC, C_SCON, C_REG, C_NONE, C_LBRA, 17, 4, 0},
	{ABR, C_NONE, C_NONE, C_NONE, C_LR, 18, 4, 0},
	{ABR, C_NONE, C_NONE, C_NONE, C_CTR, 18, 4, 0},
	{ABR, C_REG, C_NONE, C_NONE, C_CTR, 18, 4, 0},
	{ABR, C_NONE, C_NONE, C_NONE, C_ZOREG, 15, 8, 0},
	{ABC, C_NONE, C_REG, C_NONE, C_LR, 18, 4, 0},
	{ABC, C_NONE, C_REG, C_NONE, C_CTR, 18, 4, 0},
	{ABC, C_SCON, C_REG, C_NONE, C_LR, 18, 4, 0},
	{ABC, C_SCON, C_REG, C_NONE, C_CTR, 18, 4, 0},
	{ABC, C_NONE, C_NONE, C_NONE, C_ZOREG, 15, 8, 0},
	{AFMOVD, C_SEXT, C_NONE, C_NONE, C_FREG, 8, 4, REGSB},
	{AFMOVD, C_SAUTO, C_NONE, C_NONE, C_FREG, 8, 4, REGSP},
	{AFMOVD, C_SOREG, C_NONE, C_NONE, C_FREG, 8, 4, REGZERO},
	{AFMOVD, C_LEXT, C_NONE, C_NONE, C_FREG, 36, 8, REGSB},
	{AFMOVD, C_LAUTO, C_NONE, C_NONE, C_FREG, 36, 8, REGSP},
	{AFMOVD, C_LOREG, C_NONE, C_NONE, C_FREG, 36, 8, REGZERO},
	{AFMOVD, C_ZCON, C_NONE, C_NONE, C_FREG, 24, 4, 0},
	{AFMOVD, C_ADDCON, C_NONE, C_NONE, C_FREG, 24, 8, 0},
	{AFMOVD, C_ADDR, C_NONE, C_NONE, C_FREG, 75, 8, 0},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_SEXT, 7, 4, REGSB},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_SAUTO, 7, 4, REGSP},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_SOREG, 7, 4, REGZERO},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_LEXT, 35, 8, REGSB},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_LAUTO, 35, 8, REGSP},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_LOREG, 35, 8, REGZERO},
	{AFMOVD, C_FREG, C_NONE, C_NONE, C_ADDR, 74, 8, 0},
	{AFMOVSX, C_ZOREG, C_REG, C_NONE, C_FREG, 45, 4, 0},
	{AFMOVSX, C_ZOREG, C_NONE, C_NONE, C_FREG, 45, 4, 0},
	{AFMOVSX, C_FREG, C_REG, C_NONE, C_ZOREG, 44, 4, 0},
	{AFMOVSX, C_FREG, C_NONE, C_NONE, C_ZOREG, 44, 4, 0},
	{AFMOVSZ, C_ZOREG, C_REG, C_NONE, C_FREG, 45, 4, 0},
	{AFMOVSZ, C_ZOREG, C_NONE, C_NONE, C_FREG, 45, 4, 0},
	{ASYNC, C_NONE, C_NONE, C_NONE, C_NONE, 46, 4, 0},
	{AWORD, C_LCON, C_NONE, C_NONE, C_NONE, 40, 4, 0},
	{ADWORD, C_LCON, C_NONE, C_NONE, C_NONE, 31, 8, 0},
	{ADWORD, C_DCON, C_NONE, C_NONE, C_NONE, 31, 8, 0},
	{AADDME, C_REG, C_NONE, C_NONE, C_REG, 47, 4, 0},
	{AEXTSB, C_REG, C_NONE, C_NONE, C_REG, 48, 4, 0},
	{AEXTSB, C_NONE, C_NONE, C_NONE, C_REG, 48, 4, 0},
	{AISEL, C_LCON, C_REG, C_REG, C_REG, 84, 4, 0},
	{AISEL, C_ZCON, C_REG, C_REG, C_REG, 84, 4, 0},
	{ANEG, C_REG, C_NONE, C_NONE, C_REG, 47, 4, 0},
	{ANEG, C_NONE, C_NONE, C_NONE, C_REG, 47, 4, 0},
	{AREM, C_REG, C_NONE, C_NONE, C_REG, 50, 12, 0},
	{AREM, C_REG, C_REG, C_NONE, C_REG, 50, 12, 0},
	{AREMU, C_REG, C_NONE, C_NONE, C_REG, 50, 16, 0},
	{AREMU, C_REG, C_REG, C_NONE, C_REG, 50, 16, 0},
	{AREMD, C_REG, C_NONE, C_NONE, C_REG, 51, 12, 0},
	{AREMD, C_REG, C_REG, C_NONE, C_REG, 51, 12, 0},
	{AMTFSB0, C_SCON, C_NONE, C_NONE, C_NONE, 52, 4, 0},
	{AMOVFL, C_FPSCR, C_NONE, C_NONE, C_FREG, 53, 4, 0},
	{AMOVFL, C_FREG, C_NONE, C_NONE, C_FPSCR, 64, 4, 0},
	{AMOVFL, C_FREG, C_NONE, C_LCON, C_FPSCR, 64, 4, 0},
	{AMOVFL, C_LCON, C_NONE, C_NONE, C_FPSCR, 65, 4, 0},
	{AMOVD, C_MSR, C_NONE, C_NONE, C_REG, 54, 4, 0},  
	{AMOVD, C_REG, C_NONE, C_NONE, C_MSR, 54, 4, 0},  
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_MSR, 54, 4, 0}, 

	
	{APOPCNTD, C_REG, C_NONE, C_NONE, C_REG, 93, 4, 0}, 
	{ACMPB, C_REG, C_REG, C_NONE, C_REG, 92, 4, 0},     
	{ACMPEQB, C_REG, C_REG, C_NONE, C_CREG, 92, 4, 0},  
	{ACMPEQB, C_REG, C_NONE, C_NONE, C_REG, 70, 4, 0},
	{AFTDIV, C_FREG, C_FREG, C_NONE, C_SCON, 92, 4, 0},  
	{AFTSQRT, C_FREG, C_NONE, C_NONE, C_SCON, 93, 4, 0}, 
	{ACOPY, C_REG, C_NONE, C_NONE, C_REG, 92, 4, 0},     
	{ADARN, C_SCON, C_NONE, C_NONE, C_REG, 92, 4, 0},    
	{ALDMX, C_SOREG, C_NONE, C_NONE, C_REG, 45, 4, 0},   
	{AMADDHD, C_REG, C_REG, C_REG, C_REG, 83, 4, 0},     
	{AADDEX, C_REG, C_REG, C_SCON, C_REG, 94, 4, 0},     
	{ACRAND, C_CREG, C_NONE, C_NONE, C_CREG, 2, 4, 0},   

	

	
	{ALV, C_SOREG, C_NONE, C_NONE, C_VREG, 45, 4, 0}, 

	
	{ASTV, C_VREG, C_NONE, C_NONE, C_SOREG, 44, 4, 0}, 

	
	{AVAND, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVOR, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},  

	
	{AVADDUM, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVADDCU, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVADDUS, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVADDSS, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVADDE, C_VREG, C_VREG, C_VREG, C_VREG, 83, 4, 0},  

	
	{AVSUBUM, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSUBCU, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSUBUS, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSUBSS, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSUBE, C_VREG, C_VREG, C_VREG, C_VREG, 83, 4, 0},  

	
	{AVMULESB, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 9},  
	{AVPMSUM, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},   
	{AVMSUMUDM, C_VREG, C_VREG, C_VREG, C_VREG, 83, 4, 0}, 

	
	{AVR, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 

	
	{AVS, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},     
	{AVSA, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},    
	{AVSOI, C_ANDCON, C_VREG, C_VREG, C_VREG, 83, 4, 0}, 

	
	{AVCLZ, C_VREG, C_NONE, C_NONE, C_VREG, 85, 4, 0},    
	{AVPOPCNT, C_VREG, C_NONE, C_NONE, C_VREG, 85, 4, 0}, 

	
	{AVCMPEQ, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},   
	{AVCMPGT, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},   
	{AVCMPNEZB, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 

	
	{AVMRGOW, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 

	
	{AVPERM, C_VREG, C_VREG, C_VREG, C_VREG, 83, 4, 0}, 

	
	{AVBPERMQ, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 

	
	{AVSEL, C_VREG, C_VREG, C_VREG, C_VREG, 83, 4, 0}, 

	
	{AVSPLTB, C_SCON, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSPLTB, C_ADDCON, C_VREG, C_NONE, C_VREG, 82, 4, 0},
	{AVSPLTISB, C_SCON, C_NONE, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSPLTISB, C_ADDCON, C_NONE, C_NONE, C_VREG, 82, 4, 0},

	
	{AVCIPH, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0},  
	{AVNCIPH, C_VREG, C_VREG, C_NONE, C_VREG, 82, 4, 0}, 
	{AVSBOX, C_VREG, C_NONE, C_NONE, C_VREG, 82, 4, 0},  

	
	{AVSHASIGMA, C_ANDCON, C_VREG, C_ANDCON, C_VREG, 82, 4, 0}, 

	
	{ALXVD2X, C_SOREG, C_NONE, C_NONE, C_VSREG, 87, 4, 0}, 
	{ALXV, C_SOREG, C_NONE, C_NONE, C_VSREG, 96, 4, 0},    
	{ALXVL, C_REG, C_REG, C_NONE, C_VSREG, 98, 4, 0},      

	
	{ASTXVD2X, C_VSREG, C_NONE, C_NONE, C_SOREG, 86, 4, 0}, 
	{ASTXV, C_VSREG, C_NONE, C_NONE, C_SOREG, 97, 4, 0},    
	{ASTXVL, C_VSREG, C_REG, C_NONE, C_REG, 99, 4, 0},      

	
	{ALXSDX, C_SOREG, C_NONE, C_NONE, C_VSREG, 87, 4, 0}, 

	
	{ASTXSDX, C_VSREG, C_NONE, C_NONE, C_SOREG, 86, 4, 0}, 

	
	{ALXSIWAX, C_SOREG, C_NONE, C_NONE, C_VSREG, 87, 4, 0}, 

	
	{ASTXSIWX, C_VSREG, C_NONE, C_NONE, C_SOREG, 86, 4, 0}, 

	
	{AMFVSRD, C_VSREG, C_NONE, C_NONE, C_REG, 88, 4, 0}, 
	{AMFVSRD, C_FREG, C_NONE, C_NONE, C_REG, 88, 4, 0},
	{AMFVSRD, C_VREG, C_NONE, C_NONE, C_REG, 88, 4, 0},

	
	{AMTVSRD, C_REG, C_NONE, C_NONE, C_VSREG, 88, 4, 0}, 
	{AMTVSRD, C_REG, C_REG, C_NONE, C_VSREG, 88, 4, 0},
	{AMTVSRD, C_REG, C_NONE, C_NONE, C_FREG, 88, 4, 0},
	{AMTVSRD, C_REG, C_NONE, C_NONE, C_VREG, 88, 4, 0},

	
	{AXXLAND, C_VSREG, C_VSREG, C_NONE, C_VSREG, 90, 4, 0}, 
	{AXXLOR, C_VSREG, C_VSREG, C_NONE, C_VSREG, 90, 4, 0},  

	
	{AXXSEL, C_VSREG, C_VSREG, C_VSREG, C_VSREG, 91, 4, 0}, 

	
	{AXXMRGHW, C_VSREG, C_VSREG, C_NONE, C_VSREG, 90, 4, 0}, 

	
	{AXXSPLTW, C_VSREG, C_NONE, C_SCON, C_VSREG, 89, 4, 0},  
	{AXXSPLTIB, C_SCON, C_NONE, C_NONE, C_VSREG, 100, 4, 0}, 

	
	{AXXPERM, C_VSREG, C_VSREG, C_NONE, C_VSREG, 90, 4, 0}, 

	
	{AXXSLDWI, C_VSREG, C_VSREG, C_SCON, C_VSREG, 90, 4, 0}, 

	
	{AXXBRQ, C_VSREG, C_NONE, C_NONE, C_VSREG, 101, 4, 0}, 

	
	{AXSCVDPSP, C_VSREG, C_NONE, C_NONE, C_VSREG, 89, 4, 0}, 

	
	{AXVCVDPSP, C_VSREG, C_NONE, C_NONE, C_VSREG, 89, 4, 0}, 

	
	{AXSCVDPSXDS, C_VSREG, C_NONE, C_NONE, C_VSREG, 89, 4, 0}, 

	
	{AXSCVSXDDP, C_VSREG, C_NONE, C_NONE, C_VSREG, 89, 4, 0}, 

	
	{AXVCVDPSXDS, C_VSREG, C_NONE, C_NONE, C_VSREG, 89, 4, 0}, 

	
	{AXVCVSXDDP, C_VSREG, C_NONE, C_NONE, C_VSREG, 89, 4, 0}, 

	
	{AMOVD, C_REG, C_NONE, C_NONE, C_SPR, 66, 4, 0},
	{AMOVD, C_REG, C_NONE, C_NONE, C_LR, 66, 4, 0},
	{AMOVD, C_REG, C_NONE, C_NONE, C_CTR, 66, 4, 0},
	{AMOVD, C_REG, C_NONE, C_NONE, C_XER, 66, 4, 0},
	{AMOVD, C_SPR, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVD, C_LR, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVD, C_CTR, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVD, C_XER, C_NONE, C_NONE, C_REG, 66, 4, 0},

	
	{AMOVW, C_REG, C_NONE, C_NONE, C_SPR, 66, 4, 0},
	{AMOVW, C_REG, C_NONE, C_NONE, C_CTR, 66, 4, 0},
	{AMOVW, C_REG, C_NONE, C_NONE, C_XER, 66, 4, 0},
	{AMOVW, C_SPR, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVW, C_XER, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_SPR, 66, 4, 0},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_CTR, 66, 4, 0},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_XER, 66, 4, 0},
	{AMOVWZ, C_SPR, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVWZ, C_XER, C_NONE, C_NONE, C_REG, 66, 4, 0},
	{AMOVFL, C_FPSCR, C_NONE, C_NONE, C_CREG, 73, 4, 0},
	{AMOVFL, C_CREG, C_NONE, C_NONE, C_CREG, 67, 4, 0},
	{AMOVW, C_CREG, C_NONE, C_NONE, C_REG, 68, 4, 0},
	{AMOVWZ, C_CREG, C_NONE, C_NONE, C_REG, 68, 4, 0},
	{AMOVFL, C_REG, C_NONE, C_NONE, C_LCON, 69, 4, 0},
	{AMOVFL, C_REG, C_NONE, C_NONE, C_CREG, 69, 4, 0},
	{AMOVW, C_REG, C_NONE, C_NONE, C_CREG, 69, 4, 0},
	{AMOVWZ, C_REG, C_NONE, C_NONE, C_CREG, 69, 4, 0},
	{ACMP, C_REG, C_NONE, C_NONE, C_REG, 70, 4, 0},
	{ACMP, C_REG, C_REG, C_NONE, C_REG, 70, 4, 0},
	{ACMP, C_REG, C_NONE, C_NONE, C_ADDCON, 71, 4, 0},
	{ACMP, C_REG, C_REG, C_NONE, C_ADDCON, 71, 4, 0},
	{ACMPU, C_REG, C_NONE, C_NONE, C_REG, 70, 4, 0},
	{ACMPU, C_REG, C_REG, C_NONE, C_REG, 70, 4, 0},
	{ACMPU, C_REG, C_NONE, C_NONE, C_ANDCON, 71, 4, 0},
	{ACMPU, C_REG, C_REG, C_NONE, C_ANDCON, 71, 4, 0},
	{AFCMPO, C_FREG, C_NONE, C_NONE, C_FREG, 70, 4, 0},
	{AFCMPO, C_FREG, C_REG, C_NONE, C_FREG, 70, 4, 0},
	{ATW, C_LCON, C_REG, C_NONE, C_REG, 60, 4, 0},
	{ATW, C_LCON, C_REG, C_NONE, C_ADDCON, 61, 4, 0},
	{ADCBF, C_ZOREG, C_NONE, C_NONE, C_NONE, 43, 4, 0},
	{ADCBF, C_SOREG, C_NONE, C_NONE, C_NONE, 43, 4, 0},
	{ADCBF, C_ZOREG, C_REG, C_NONE, C_SCON, 43, 4, 0},
	{ADCBF, C_SOREG, C_NONE, C_NONE, C_SCON, 43, 4, 0},
	{AECOWX, C_REG, C_REG, C_NONE, C_ZOREG, 44, 4, 0},
	{AECIWX, C_ZOREG, C_REG, C_NONE, C_REG, 45, 4, 0},
	{AECOWX, C_REG, C_NONE, C_NONE, C_ZOREG, 44, 4, 0},
	{AECIWX, C_ZOREG, C_NONE, C_NONE, C_REG, 45, 4, 0},
	{ALDAR, C_ZOREG, C_NONE, C_NONE, C_REG, 45, 4, 0},
	{ALDAR, C_ZOREG, C_NONE, C_ANDCON, C_REG, 45, 4, 0},
	{AEIEIO, C_NONE, C_NONE, C_NONE, C_NONE, 46, 4, 0},
	{ATLBIE, C_REG, C_NONE, C_NONE, C_NONE, 49, 4, 0},
	{ATLBIE, C_SCON, C_NONE, C_NONE, C_REG, 49, 4, 0},
	{ASLBMFEE, C_REG, C_NONE, C_NONE, C_REG, 55, 4, 0},
	{ASLBMTE, C_REG, C_NONE, C_NONE, C_REG, 55, 4, 0},
	{ASTSW, C_REG, C_NONE, C_NONE, C_ZOREG, 44, 4, 0},
	{ASTSW, C_REG, C_NONE, C_LCON, C_ZOREG, 41, 4, 0},
	{ALSW, C_ZOREG, C_NONE, C_NONE, C_REG, 45, 4, 0},
	{ALSW, C_ZOREG, C_NONE, C_LCON, C_REG, 42, 4, 0},
	{obj.AUNDEF, C_NONE, C_NONE, C_NONE, C_NONE, 78, 4, 0},
	{obj.APCDATA, C_LCON, C_NONE, C_NONE, C_LCON, 0, 0, 0},
	{obj.AFUNCDATA, C_SCON, C_NONE, C_NONE, C_ADDR, 0, 0, 0},
	{obj.ANOP, C_NONE, C_NONE, C_NONE, C_NONE, 0, 0, 0},
	{obj.ANOP, C_LCON, C_NONE, C_NONE, C_NONE, 0, 0, 0}, 
	{obj.ANOP, C_REG, C_NONE, C_NONE, C_NONE, 0, 0, 0},  
	{obj.ANOP, C_FREG, C_NONE, C_NONE, C_NONE, 0, 0, 0},
	{obj.ADUFFZERO, C_NONE, C_NONE, C_NONE, C_LBRA, 11, 4, 0}, 
	{obj.ADUFFCOPY, C_NONE, C_NONE, C_NONE, C_LBRA, 11, 4, 0}, 
	{obj.APCALIGN, C_LCON, C_NONE, C_NONE, C_NONE, 0, 0, 0},   

	{obj.AXXX, C_NONE, C_NONE, C_NONE, C_NONE, 0, 4, 0},
}

var oprange [ALAST & obj.AMask][]Optab

var xcmp [C_NCLASS][C_NCLASS]bool


func addpad(pc, a int64, ctxt *obj.Link, cursym *obj.LSym) int {
	
	
	switch a {
	case 8:
		if pc&7 != 0 {
			return 4
		}
	case 16:
		
		
		switch pc & 15 {
		case 4, 12:
			return 4
		case 8:
			return 8
		}
	case 32:
		
		
		switch pc & 31 {
		case 4, 20:
			return 12
		case 8, 24:
			return 8
		case 12, 28:
			return 4
		}
		
		
		
		
		
		if ctxt.Headtype != objabi.Haix && cursym.Func.Align < 32 {
			cursym.Func.Align = 32
		}
	default:
		ctxt.Diag("Unexpected alignment: %d for PCALIGN directive\n", a)
	}
	return 0
}

func span9(ctxt *obj.Link, cursym *obj.LSym, newprog obj.ProgAlloc) {
	p := cursym.Func.Text
	if p == nil || p.Link == nil { 
		return
	}

	if oprange[AANDN&obj.AMask] == nil {
		ctxt.Diag("ppc64 ops not initialized, call ppc64.buildop first")
	}

	c := ctxt9{ctxt: ctxt, newprog: newprog, cursym: cursym, autosize: int32(p.To.Offset)}

	pc := int64(0)
	p.Pc = pc

	var m int
	var o *Optab
	for p = p.Link; p != nil; p = p.Link {
		p.Pc = pc
		o = c.oplook(p)
		m = int(o.size)
		if m == 0 {
			if p.As == obj.APCALIGN {
				a := c.vregoff(&p.From)
				m = addpad(pc, a, ctxt, cursym)
			} else {
				if p.As != obj.ANOP && p.As != obj.AFUNCDATA && p.As != obj.APCDATA {
					ctxt.Diag("zero-width instruction\n%v", p)
				}
				continue
			}
		}
		pc += int64(m)
	}

	c.cursym.Size = pc

	
	bflag := 1

	var otxt int64
	var q *obj.Prog
	for bflag != 0 {
		bflag = 0
		pc = 0
		for p = c.cursym.Func.Text.Link; p != nil; p = p.Link {
			p.Pc = pc
			o = c.oplook(p)

			
			if (o.type_ == 16 || o.type_ == 17) && p.To.Target() != nil {
				otxt = p.To.Target().Pc - pc
				if otxt < -(1<<15)+10 || otxt >= (1<<15)-10 {
					q = c.newprog()
					q.Link = p.Link
					p.Link = q
					q.As = ABR
					q.To.Type = obj.TYPE_BRANCH
					q.To.SetTarget(p.To.Target())
					p.To.SetTarget(q)
					q = c.newprog()
					q.Link = p.Link
					p.Link = q
					q.As = ABR
					q.To.Type = obj.TYPE_BRANCH
					q.To.SetTarget(q.Link.Link)

					
					
					bflag = 1
				}
			}

			m = int(o.size)
			if m == 0 {
				if p.As == obj.APCALIGN {
					a := c.vregoff(&p.From)
					m = addpad(pc, a, ctxt, cursym)
				} else {
					if p.As != obj.ANOP && p.As != obj.AFUNCDATA && p.As != obj.APCDATA {
						ctxt.Diag("zero-width instruction\n%v", p)
					}
					continue
				}
			}

			pc += int64(m)
		}

		c.cursym.Size = pc
	}

	if r := pc & funcAlignMask; r != 0 {
		pc += funcAlign - r
	}

	c.cursym.Size = pc

	

	c.cursym.Grow(c.cursym.Size)

	bp := c.cursym.P
	var i int32
	var out [6]uint32
	for p := c.cursym.Func.Text.Link; p != nil; p = p.Link {
		c.pc = p.Pc
		o = c.oplook(p)
		if int(o.size) > 4*len(out) {
			log.Fatalf("out array in span9 is too small, need at least %d for %v", o.size/4, p)
		}
		
		if o.type_ == 0 && p.As == obj.APCALIGN {
			pad := LOP_RRR(OP_OR, REGZERO, REGZERO, REGZERO)
			aln := c.vregoff(&p.From)
			v := addpad(p.Pc, aln, c.ctxt, c.cursym)
			if v > 0 {
				
				for i = 0; i < int32(v/4); i++ {
					c.ctxt.Arch.ByteOrder.PutUint32(bp, pad)
					bp = bp[4:]
				}
			}
		} else {
			c.asmout(p, o, out[:])
			for i = 0; i < int32(o.size/4); i++ {
				c.ctxt.Arch.ByteOrder.PutUint32(bp, out[i])
				bp = bp[4:]
			}
		}
	}
}

func isint32(v int64) bool {
	return int64(int32(v)) == v
}

func isuint32(v uint64) bool {
	return uint64(uint32(v)) == v
}

func (c *ctxt9) aclass(a *obj.Addr) int {
	switch a.Type {
	case obj.TYPE_NONE:
		return C_NONE

	case obj.TYPE_REG:
		if REG_R0 <= a.Reg && a.Reg <= REG_R31 {
			return C_REG
		}
		if REG_F0 <= a.Reg && a.Reg <= REG_F31 {
			return C_FREG
		}
		if REG_V0 <= a.Reg && a.Reg <= REG_V31 {
			return C_VREG
		}
		if REG_VS0 <= a.Reg && a.Reg <= REG_VS63 {
			return C_VSREG
		}
		if REG_CR0 <= a.Reg && a.Reg <= REG_CR7 || a.Reg == REG_CR {
			return C_CREG
		}
		if REG_SPR0 <= a.Reg && a.Reg <= REG_SPR0+1023 {
			switch a.Reg {
			case REG_LR:
				return C_LR

			case REG_XER:
				return C_XER

			case REG_CTR:
				return C_CTR
			}

			return C_SPR
		}

		if REG_DCR0 <= a.Reg && a.Reg <= REG_DCR0+1023 {
			return C_SPR
		}
		if a.Reg == REG_FPSCR {
			return C_FPSCR
		}
		if a.Reg == REG_MSR {
			return C_MSR
		}
		return C_GOK

	case obj.TYPE_MEM:
		switch a.Name {
		case obj.NAME_EXTERN,
			obj.NAME_STATIC:
			if a.Sym == nil {
				break
			}
			c.instoffset = a.Offset
			if a.Sym != nil { 
				if a.Sym.Type == objabi.STLSBSS {
					if c.ctxt.Flag_shared {
						return C_TLS_IE
					} else {
						return C_TLS_LE
					}
				}
				return C_ADDR
			}
			return C_LEXT

		case obj.NAME_GOTREF:
			return C_GOTADDR

		case obj.NAME_TOCREF:
			return C_TOCADDR

		case obj.NAME_AUTO:
			c.instoffset = int64(c.autosize) + a.Offset
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SAUTO
			}
			return C_LAUTO

		case obj.NAME_PARAM:
			c.instoffset = int64(c.autosize) + a.Offset + c.ctxt.FixedFrameSize()
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SAUTO
			}
			return C_LAUTO

		case obj.NAME_NONE:
			c.instoffset = a.Offset
			if c.instoffset == 0 {
				return C_ZOREG
			}
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SOREG
			}
			return C_LOREG
		}

		return C_GOK

	case obj.TYPE_TEXTSIZE:
		return C_TEXTSIZE

	case obj.TYPE_FCONST:
		
		
		f64 := a.Val.(float64)
		if f64 == 0 {
			if math.Signbit(f64) {
				return C_ADDCON
			}
			return C_ZCON
		}
		log.Fatalf("Unexpected nonzero FCONST operand %v", a)

	case obj.TYPE_CONST,
		obj.TYPE_ADDR:
		switch a.Name {
		case obj.NAME_NONE:
			c.instoffset = a.Offset
			if a.Reg != 0 {
				if -BIG <= c.instoffset && c.instoffset <= BIG {
					return C_SACON
				}
				if isint32(c.instoffset) {
					return C_LACON
				}
				return C_DACON
			}

		case obj.NAME_EXTERN,
			obj.NAME_STATIC:
			s := a.Sym
			if s == nil {
				return C_GOK
			}

			c.instoffset = a.Offset

			
			return C_LCON

		case obj.NAME_AUTO:
			c.instoffset = int64(c.autosize) + a.Offset
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SACON
			}
			return C_LACON

		case obj.NAME_PARAM:
			c.instoffset = int64(c.autosize) + a.Offset + c.ctxt.FixedFrameSize()
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SACON
			}
			return C_LACON

		default:
			return C_GOK
		}

		if c.instoffset >= 0 {
			if c.instoffset == 0 {
				return C_ZCON
			}
			if c.instoffset <= 0x7fff {
				return C_SCON
			}
			if c.instoffset <= 0xffff {
				return C_ANDCON
			}
			if c.instoffset&0xffff == 0 && isuint32(uint64(c.instoffset)) { 
				return C_UCON
			}
			if isint32(c.instoffset) || isuint32(uint64(c.instoffset)) {
				return C_LCON
			}
			return C_DCON
		}

		if c.instoffset >= -0x8000 {
			return C_ADDCON
		}
		if c.instoffset&0xffff == 0 && isint32(c.instoffset) {
			return C_UCON
		}
		if isint32(c.instoffset) {
			return C_LCON
		}
		return C_DCON

	case obj.TYPE_BRANCH:
		if a.Sym != nil && c.ctxt.Flag_dynlink {
			return C_LBRAPIC
		}
		return C_SBRA
	}

	return C_GOK
}

func prasm(p *obj.Prog) {
	fmt.Printf("%v\n", p)
}

func (c *ctxt9) oplook(p *obj.Prog) *Optab {
	a1 := int(p.Optab)
	if a1 != 0 {
		return &optab[a1-1]
	}
	a1 = int(p.From.Class)
	if a1 == 0 {
		a1 = c.aclass(&p.From) + 1
		p.From.Class = int8(a1)
	}

	a1--
	a3 := C_NONE + 1
	if p.GetFrom3() != nil {
		a3 = int(p.GetFrom3().Class)
		if a3 == 0 {
			a3 = c.aclass(p.GetFrom3()) + 1
			p.GetFrom3().Class = int8(a3)
		}
	}

	a3--
	a4 := int(p.To.Class)
	if a4 == 0 {
		a4 = c.aclass(&p.To) + 1
		p.To.Class = int8(a4)
	}

	a4--
	a2 := C_NONE
	if p.Reg != 0 {
		if REG_R0 <= p.Reg && p.Reg <= REG_R31 {
			a2 = C_REG
		} else if REG_V0 <= p.Reg && p.Reg <= REG_V31 {
			a2 = C_VREG
		} else if REG_VS0 <= p.Reg && p.Reg <= REG_VS63 {
			a2 = C_VSREG
		} else if REG_F0 <= p.Reg && p.Reg <= REG_F31 {
			a2 = C_FREG
		}
	}

	
	ops := oprange[p.As&obj.AMask]
	c1 := &xcmp[a1]
	c3 := &xcmp[a3]
	c4 := &xcmp[a4]
	for i := range ops {
		op := &ops[i]
		if int(op.a2) == a2 && c1[op.a1] && c3[op.a3] && c4[op.a4] {
			p.Optab = uint16(cap(optab) - cap(ops) + i + 1)
			return op
		}
	}

	c.ctxt.Diag("illegal combination %v %v %v %v %v", p.As, DRconv(a1), DRconv(a2), DRconv(a3), DRconv(a4))
	prasm(p)
	if ops == nil {
		ops = optab
	}
	return &ops[0]
}

func cmp(a int, b int) bool {
	if a == b {
		return true
	}
	switch a {
	case C_LCON:
		if b == C_ZCON || b == C_SCON || b == C_UCON || b == C_ADDCON || b == C_ANDCON {
			return true
		}

	case C_ADDCON:
		if b == C_ZCON || b == C_SCON {
			return true
		}

	case C_ANDCON:
		if b == C_ZCON || b == C_SCON {
			return true
		}

	case C_SPR:
		if b == C_LR || b == C_XER || b == C_CTR {
			return true
		}

	case C_UCON:
		if b == C_ZCON {
			return true
		}

	case C_SCON:
		if b == C_ZCON {
			return true
		}

	case C_LACON:
		if b == C_SACON {
			return true
		}

	case C_LBRA:
		if b == C_SBRA {
			return true
		}

	case C_LEXT:
		if b == C_SEXT {
			return true
		}

	case C_LAUTO:
		if b == C_SAUTO {
			return true
		}

	case C_REG:
		if b == C_ZCON {
			return r0iszero != 0 
		}

	case C_LOREG:
		if b == C_ZOREG || b == C_SOREG {
			return true
		}

	case C_SOREG:
		if b == C_ZOREG {
			return true
		}

	case C_ANY:
		return true
	}

	return false
}

type ocmp []Optab

func (x ocmp) Len() int {
	return len(x)
}

func (x ocmp) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}




func (x ocmp) Less(i, j int) bool {
	p1 := &x[i]
	p2 := &x[j]
	n := int(p1.as) - int(p2.as)
	
	if n != 0 {
		return n < 0
	}
	
	
	n = int(p1.size) - int(p2.size)
	if n != 0 {
		return n < 0
	}
	
	
	n = int(p1.a1) - int(p2.a1)
	if n != 0 {
		return n < 0
	}
	n = int(p1.a2) - int(p2.a2)
	if n != 0 {
		return n < 0
	}
	n = int(p1.a3) - int(p2.a3)
	if n != 0 {
		return n < 0
	}
	n = int(p1.a4) - int(p2.a4)
	if n != 0 {
		return n < 0
	}
	return false
}




func opset(a, b0 obj.As) {
	oprange[a&obj.AMask] = oprange[b0]
}


func buildop(ctxt *obj.Link) {
	if oprange[AANDN&obj.AMask] != nil {
		
		
		
		return
	}

	var n int

	for i := 0; i < C_NCLASS; i++ {
		for n = 0; n < C_NCLASS; n++ {
			if cmp(n, i) {
				xcmp[i][n] = true
			}
		}
	}
	for n = 0; optab[n].as != obj.AXXX; n++ {
	}
	sort.Sort(ocmp(optab[:n]))
	for i := 0; i < n; i++ {
		r := optab[i].as
		r0 := r & obj.AMask
		start := i
		for optab[i].as == r {
			i++
		}
		oprange[r0] = optab[start:i]
		i--

		switch r {
		default:
			ctxt.Diag("unknown op in build: %v", r)
			log.Fatalf("instruction missing from switch in asm9.go:buildop: %v", r)

		case ADCBF: 
			opset(ADCBI, r0)

			opset(ADCBST, r0)
			opset(ADCBT, r0)
			opset(ADCBTST, r0)
			opset(ADCBZ, r0)
			opset(AICBI, r0)

		case AECOWX: 
			opset(ASTWCCC, r0)
			opset(ASTHCCC, r0)
			opset(ASTBCCC, r0)
			opset(ASTDCCC, r0)

		case AREM: 
			opset(AREM, r0)

		case AREMU:
			opset(AREMU, r0)

		case AREMD:
			opset(AREMDU, r0)

		case ADIVW: 
			opset(AMULHW, r0)

			opset(AMULHWCC, r0)
			opset(AMULHWU, r0)
			opset(AMULHWUCC, r0)
			opset(AMULLWCC, r0)
			opset(AMULLWVCC, r0)
			opset(AMULLWV, r0)
			opset(ADIVWCC, r0)
			opset(ADIVWV, r0)
			opset(ADIVWVCC, r0)
			opset(ADIVWU, r0)
			opset(ADIVWUCC, r0)
			opset(ADIVWUV, r0)
			opset(ADIVWUVCC, r0)
			opset(AMODUD, r0)
			opset(AMODUW, r0)
			opset(AMODSD, r0)
			opset(AMODSW, r0)
			opset(AADDCC, r0)
			opset(AADDCV, r0)
			opset(AADDCVCC, r0)
			opset(AADDV, r0)
			opset(AADDVCC, r0)
			opset(AADDE, r0)
			opset(AADDECC, r0)
			opset(AADDEV, r0)
			opset(AADDEVCC, r0)
			opset(AMULHD, r0)
			opset(AMULHDCC, r0)
			opset(AMULHDU, r0)
			opset(AMULHDUCC, r0)
			opset(AMULLD, r0)
			opset(AMULLDCC, r0)
			opset(AMULLDVCC, r0)
			opset(AMULLDV, r0)
			opset(ADIVD, r0)
			opset(ADIVDCC, r0)
			opset(ADIVDE, r0)
			opset(ADIVDEU, r0)
			opset(ADIVDECC, r0)
			opset(ADIVDEUCC, r0)
			opset(ADIVDVCC, r0)
			opset(ADIVDV, r0)
			opset(ADIVDU, r0)
			opset(ADIVDUV, r0)
			opset(ADIVDUVCC, r0)
			opset(ADIVDUCC, r0)

		case ACRAND:
			opset(ACRANDN, r0)
			opset(ACREQV, r0)
			opset(ACRNAND, r0)
			opset(ACRNOR, r0)
			opset(ACROR, r0)
			opset(ACRORN, r0)
			opset(ACRXOR, r0)

		case APOPCNTD: 
			opset(APOPCNTW, r0)
			opset(APOPCNTB, r0)
			opset(ACNTTZW, r0)
			opset(ACNTTZWCC, r0)
			opset(ACNTTZD, r0)
			opset(ACNTTZDCC, r0)

		case ACOPY: 
			opset(APASTECC, r0)

		case AMADDHD: 
			opset(AMADDHDU, r0)
			opset(AMADDLD, r0)

		case AMOVBZ: 
			opset(AMOVH, r0)
			opset(AMOVHZ, r0)

		case AMOVBZU: 
			opset(AMOVHU, r0)

			opset(AMOVHZU, r0)
			opset(AMOVWU, r0)
			opset(AMOVWZU, r0)
			opset(AMOVDU, r0)
			opset(AMOVMW, r0)

		case ALV: 
			opset(ALVEBX, r0)
			opset(ALVEHX, r0)
			opset(ALVEWX, r0)
			opset(ALVX, r0)
			opset(ALVXL, r0)
			opset(ALVSL, r0)
			opset(ALVSR, r0)

		case ASTV: 
			opset(ASTVEBX, r0)
			opset(ASTVEHX, r0)
			opset(ASTVEWX, r0)
			opset(ASTVX, r0)
			opset(ASTVXL, r0)

		case AVAND: 
			opset(AVAND, r0)
			opset(AVANDC, r0)
			opset(AVNAND, r0)

		case AVMRGOW: 
			opset(AVMRGEW, r0)

		case AVOR: 
			opset(AVOR, r0)
			opset(AVORC, r0)
			opset(AVXOR, r0)
			opset(AVNOR, r0)
			opset(AVEQV, r0)

		case AVADDUM: 
			opset(AVADDUBM, r0)
			opset(AVADDUHM, r0)
			opset(AVADDUWM, r0)
			opset(AVADDUDM, r0)
			opset(AVADDUQM, r0)

		case AVADDCU: 
			opset(AVADDCUQ, r0)
			opset(AVADDCUW, r0)

		case AVADDUS: 
			opset(AVADDUBS, r0)
			opset(AVADDUHS, r0)
			opset(AVADDUWS, r0)

		case AVADDSS: 
			opset(AVADDSBS, r0)
			opset(AVADDSHS, r0)
			opset(AVADDSWS, r0)

		case AVADDE: 
			opset(AVADDEUQM, r0)
			opset(AVADDECUQ, r0)

		case AVSUBUM: 
			opset(AVSUBUBM, r0)
			opset(AVSUBUHM, r0)
			opset(AVSUBUWM, r0)
			opset(AVSUBUDM, r0)
			opset(AVSUBUQM, r0)

		case AVSUBCU: 
			opset(AVSUBCUQ, r0)
			opset(AVSUBCUW, r0)

		case AVSUBUS: 
			opset(AVSUBUBS, r0)
			opset(AVSUBUHS, r0)
			opset(AVSUBUWS, r0)

		case AVSUBSS: 
			opset(AVSUBSBS, r0)
			opset(AVSUBSHS, r0)
			opset(AVSUBSWS, r0)

		case AVSUBE: 
			opset(AVSUBEUQM, r0)
			opset(AVSUBECUQ, r0)

		case AVMULESB: 
			opset(AVMULOSB, r0)
			opset(AVMULEUB, r0)
			opset(AVMULOUB, r0)
			opset(AVMULESH, r0)
			opset(AVMULOSH, r0)
			opset(AVMULEUH, r0)
			opset(AVMULOUH, r0)
			opset(AVMULESW, r0)
			opset(AVMULOSW, r0)
			opset(AVMULEUW, r0)
			opset(AVMULOUW, r0)
			opset(AVMULUWM, r0)
		case AVPMSUM: 
			opset(AVPMSUMB, r0)
			opset(AVPMSUMH, r0)
			opset(AVPMSUMW, r0)
			opset(AVPMSUMD, r0)

		case AVR: 
			opset(AVRLB, r0)
			opset(AVRLH, r0)
			opset(AVRLW, r0)
			opset(AVRLD, r0)

		case AVS: 
			opset(AVSLB, r0)
			opset(AVSLH, r0)
			opset(AVSLW, r0)
			opset(AVSL, r0)
			opset(AVSLO, r0)
			opset(AVSRB, r0)
			opset(AVSRH, r0)
			opset(AVSRW, r0)
			opset(AVSR, r0)
			opset(AVSRO, r0)
			opset(AVSLD, r0)
			opset(AVSRD, r0)

		case AVSA: 
			opset(AVSRAB, r0)
			opset(AVSRAH, r0)
			opset(AVSRAW, r0)
			opset(AVSRAD, r0)

		case AVSOI: 
			opset(AVSLDOI, r0)

		case AVCLZ: 
			opset(AVCLZB, r0)
			opset(AVCLZH, r0)
			opset(AVCLZW, r0)
			opset(AVCLZD, r0)

		case AVPOPCNT: 
			opset(AVPOPCNTB, r0)
			opset(AVPOPCNTH, r0)
			opset(AVPOPCNTW, r0)
			opset(AVPOPCNTD, r0)

		case AVCMPEQ: 
			opset(AVCMPEQUB, r0)
			opset(AVCMPEQUBCC, r0)
			opset(AVCMPEQUH, r0)
			opset(AVCMPEQUHCC, r0)
			opset(AVCMPEQUW, r0)
			opset(AVCMPEQUWCC, r0)
			opset(AVCMPEQUD, r0)
			opset(AVCMPEQUDCC, r0)

		case AVCMPGT: 
			opset(AVCMPGTUB, r0)
			opset(AVCMPGTUBCC, r0)
			opset(AVCMPGTUH, r0)
			opset(AVCMPGTUHCC, r0)
			opset(AVCMPGTUW, r0)
			opset(AVCMPGTUWCC, r0)
			opset(AVCMPGTUD, r0)
			opset(AVCMPGTUDCC, r0)
			opset(AVCMPGTSB, r0)
			opset(AVCMPGTSBCC, r0)
			opset(AVCMPGTSH, r0)
			opset(AVCMPGTSHCC, r0)
			opset(AVCMPGTSW, r0)
			opset(AVCMPGTSWCC, r0)
			opset(AVCMPGTSD, r0)
			opset(AVCMPGTSDCC, r0)

		case AVCMPNEZB: 
			opset(AVCMPNEZBCC, r0)
			opset(AVCMPNEB, r0)
			opset(AVCMPNEBCC, r0)
			opset(AVCMPNEH, r0)
			opset(AVCMPNEHCC, r0)
			opset(AVCMPNEW, r0)
			opset(AVCMPNEWCC, r0)

		case AVPERM: 
			opset(AVPERMXOR, r0)
			opset(AVPERMR, r0)

		case AVBPERMQ: 
			opset(AVBPERMD, r0)

		case AVSEL: 
			opset(AVSEL, r0)

		case AVSPLTB: 
			opset(AVSPLTH, r0)
			opset(AVSPLTW, r0)

		case AVSPLTISB: 
			opset(AVSPLTISH, r0)
			opset(AVSPLTISW, r0)

		case AVCIPH: 
			opset(AVCIPHER, r0)
			opset(AVCIPHERLAST, r0)

		case AVNCIPH: 
			opset(AVNCIPHER, r0)
			opset(AVNCIPHERLAST, r0)

		case AVSBOX: 
			opset(AVSBOX, r0)

		case AVSHASIGMA: 
			opset(AVSHASIGMAW, r0)
			opset(AVSHASIGMAD, r0)

		case ALXVD2X: 
			opset(ALXVDSX, r0)
			opset(ALXVW4X, r0)
			opset(ALXVH8X, r0)
			opset(ALXVB16X, r0)

		case ALXV: 
			opset(ALXV, r0)

		case ALXVL: 
			opset(ALXVLL, r0)
			opset(ALXVX, r0)

		case ASTXVD2X: 
			opset(ASTXVW4X, r0)
			opset(ASTXVH8X, r0)
			opset(ASTXVB16X, r0)

		case ASTXV: 
			opset(ASTXV, r0)

		case ASTXVL: 
			opset(ASTXVLL, r0)
			opset(ASTXVX, r0)

		case ALXSDX: 
			opset(ALXSDX, r0)

		case ASTXSDX: 
			opset(ASTXSDX, r0)

		case ALXSIWAX: 
			opset(ALXSIWZX, r0)

		case ASTXSIWX: 
			opset(ASTXSIWX, r0)

		case AMFVSRD: 
			opset(AMFFPRD, r0)
			opset(AMFVRD, r0)
			opset(AMFVSRWZ, r0)
			opset(AMFVSRLD, r0)

		case AMTVSRD: 
			opset(AMTFPRD, r0)
			opset(AMTVRD, r0)
			opset(AMTVSRWA, r0)
			opset(AMTVSRWZ, r0)
			opset(AMTVSRDD, r0)
			opset(AMTVSRWS, r0)

		case AXXLAND: 
			opset(AXXLANDC, r0)
			opset(AXXLEQV, r0)
			opset(AXXLNAND, r0)

		case AXXLOR: 
			opset(AXXLORC, r0)
			opset(AXXLNOR, r0)
			opset(AXXLORQ, r0)
			opset(AXXLXOR, r0)

		case AXXSEL: 
			opset(AXXSEL, r0)

		case AXXMRGHW: 
			opset(AXXMRGLW, r0)

		case AXXSPLTW: 
			opset(AXXSPLTW, r0)

		case AXXSPLTIB: 
			opset(AXXSPLTIB, r0)

		case AXXPERM: 
			opset(AXXPERM, r0)

		case AXXSLDWI: 
			opset(AXXPERMDI, r0)
			opset(AXXSLDWI, r0)

		case AXXBRQ: 
			opset(AXXBRD, r0)
			opset(AXXBRW, r0)
			opset(AXXBRH, r0)

		case AXSCVDPSP: 
			opset(AXSCVSPDP, r0)
			opset(AXSCVDPSPN, r0)
			opset(AXSCVSPDPN, r0)

		case AXVCVDPSP: 
			opset(AXVCVSPDP, r0)

		case AXSCVDPSXDS: 
			opset(AXSCVDPSXWS, r0)
			opset(AXSCVDPUXDS, r0)
			opset(AXSCVDPUXWS, r0)

		case AXSCVSXDDP: 
			opset(AXSCVUXDDP, r0)
			opset(AXSCVSXDSP, r0)
			opset(AXSCVUXDSP, r0)

		case AXVCVDPSXDS: 
			opset(AXVCVDPSXDS, r0)
			opset(AXVCVDPSXWS, r0)
			opset(AXVCVDPUXDS, r0)
			opset(AXVCVDPUXWS, r0)
			opset(AXVCVSPSXDS, r0)
			opset(AXVCVSPSXWS, r0)
			opset(AXVCVSPUXDS, r0)
			opset(AXVCVSPUXWS, r0)

		case AXVCVSXDDP: 
			opset(AXVCVSXWDP, r0)
			opset(AXVCVUXDDP, r0)
			opset(AXVCVUXWDP, r0)
			opset(AXVCVSXDSP, r0)
			opset(AXVCVSXWSP, r0)
			opset(AXVCVUXDSP, r0)
			opset(AXVCVUXWSP, r0)

		case AAND: 
			opset(AANDN, r0)
			opset(AANDNCC, r0)
			opset(AEQV, r0)
			opset(AEQVCC, r0)
			opset(ANAND, r0)
			opset(ANANDCC, r0)
			opset(ANOR, r0)
			opset(ANORCC, r0)
			opset(AORCC, r0)
			opset(AORN, r0)
			opset(AORNCC, r0)
			opset(AXORCC, r0)

		case AADDME: 
			opset(AADDMECC, r0)

			opset(AADDMEV, r0)
			opset(AADDMEVCC, r0)
			opset(AADDZE, r0)
			opset(AADDZECC, r0)
			opset(AADDZEV, r0)
			opset(AADDZEVCC, r0)
			opset(ASUBME, r0)
			opset(ASUBMECC, r0)
			opset(ASUBMEV, r0)
			opset(ASUBMEVCC, r0)
			opset(ASUBZE, r0)
			opset(ASUBZECC, r0)
			opset(ASUBZEV, r0)
			opset(ASUBZEVCC, r0)

		case AADDC:
			opset(AADDCCC, r0)

		case ABEQ:
			opset(ABGE, r0)
			opset(ABGT, r0)
			opset(ABLE, r0)
			opset(ABLT, r0)
			opset(ABNE, r0)
			opset(ABVC, r0)
			opset(ABVS, r0)

		case ABR:
			opset(ABL, r0)

		case ABC:
			opset(ABCL, r0)

		case AEXTSB: 
			opset(AEXTSBCC, r0)

			opset(AEXTSH, r0)
			opset(AEXTSHCC, r0)
			opset(ACNTLZW, r0)
			opset(ACNTLZWCC, r0)
			opset(ACNTLZD, r0)
			opset(AEXTSW, r0)
			opset(AEXTSWCC, r0)
			opset(ACNTLZDCC, r0)

		case AFABS: 
			opset(AFABSCC, r0)

			opset(AFNABS, r0)
			opset(AFNABSCC, r0)
			opset(AFNEG, r0)
			opset(AFNEGCC, r0)
			opset(AFRSP, r0)
			opset(AFRSPCC, r0)
			opset(AFCTIW, r0)
			opset(AFCTIWCC, r0)
			opset(AFCTIWZ, r0)
			opset(AFCTIWZCC, r0)
			opset(AFCTID, r0)
			opset(AFCTIDCC, r0)
			opset(AFCTIDZ, r0)
			opset(AFCTIDZCC, r0)
			opset(AFCFID, r0)
			opset(AFCFIDCC, r0)
			opset(AFCFIDU, r0)
			opset(AFCFIDUCC, r0)
			opset(AFCFIDS, r0)
			opset(AFCFIDSCC, r0)
			opset(AFRES, r0)
			opset(AFRESCC, r0)
			opset(AFRIM, r0)
			opset(AFRIMCC, r0)
			opset(AFRIP, r0)
			opset(AFRIPCC, r0)
			opset(AFRIZ, r0)
			opset(AFRIZCC, r0)
			opset(AFRIN, r0)
			opset(AFRINCC, r0)
			opset(AFRSQRTE, r0)
			opset(AFRSQRTECC, r0)
			opset(AFSQRT, r0)
			opset(AFSQRTCC, r0)
			opset(AFSQRTS, r0)
			opset(AFSQRTSCC, r0)

		case AFADD:
			opset(AFADDS, r0)
			opset(AFADDCC, r0)
			opset(AFADDSCC, r0)
			opset(AFCPSGN, r0)
			opset(AFCPSGNCC, r0)
			opset(AFDIV, r0)
			opset(AFDIVS, r0)
			opset(AFDIVCC, r0)
			opset(AFDIVSCC, r0)
			opset(AFSUB, r0)
			opset(AFSUBS, r0)
			opset(AFSUBCC, r0)
			opset(AFSUBSCC, r0)

		case AFMADD:
			opset(AFMADDCC, r0)
			opset(AFMADDS, r0)
			opset(AFMADDSCC, r0)
			opset(AFMSUB, r0)
			opset(AFMSUBCC, r0)
			opset(AFMSUBS, r0)
			opset(AFMSUBSCC, r0)
			opset(AFNMADD, r0)
			opset(AFNMADDCC, r0)
			opset(AFNMADDS, r0)
			opset(AFNMADDSCC, r0)
			opset(AFNMSUB, r0)
			opset(AFNMSUBCC, r0)
			opset(AFNMSUBS, r0)
			opset(AFNMSUBSCC, r0)
			opset(AFSEL, r0)
			opset(AFSELCC, r0)

		case AFMUL:
			opset(AFMULS, r0)
			opset(AFMULCC, r0)
			opset(AFMULSCC, r0)

		case AFCMPO:
			opset(AFCMPU, r0)

		case AISEL:
			opset(AISEL, r0)

		case AMTFSB0:
			opset(AMTFSB0CC, r0)
			opset(AMTFSB1, r0)
			opset(AMTFSB1CC, r0)

		case ANEG: 
			opset(ANEGCC, r0)

			opset(ANEGV, r0)
			opset(ANEGVCC, r0)

		case AOR: 
			opset(AXOR, r0)

		case AORIS: 
			opset(AXORIS, r0)

		case ASLW:
			opset(ASLWCC, r0)
			opset(ASRW, r0)
			opset(ASRWCC, r0)
			opset(AROTLW, r0)

		case ASLD:
			opset(ASLDCC, r0)
			opset(ASRD, r0)
			opset(ASRDCC, r0)
			opset(AROTL, r0)

		case ASRAW: 
			opset(ASRAWCC, r0)

		case ASRAD: 
			opset(ASRADCC, r0)

		case ASUB: 
			opset(ASUB, r0)

			opset(ASUBCC, r0)
			opset(ASUBV, r0)
			opset(ASUBVCC, r0)
			opset(ASUBCCC, r0)
			opset(ASUBCV, r0)
			opset(ASUBCVCC, r0)
			opset(ASUBE, r0)
			opset(ASUBECC, r0)
			opset(ASUBEV, r0)
			opset(ASUBEVCC, r0)

		case ASYNC:
			opset(AISYNC, r0)
			opset(ALWSYNC, r0)
			opset(APTESYNC, r0)
			opset(ATLBSYNC, r0)

		case ARLWMI:
			opset(ARLWMICC, r0)
			opset(ARLWNM, r0)
			opset(ARLWNMCC, r0)
			opset(ACLRLSLWI, r0)

		case ARLDMI:
			opset(ARLDMICC, r0)
			opset(ARLDIMI, r0)
			opset(ARLDIMICC, r0)

		case ARLDC:
			opset(ARLDCCC, r0)

		case ARLDCL:
			opset(ARLDCR, r0)
			opset(ARLDCLCC, r0)
			opset(ARLDCRCC, r0)

		case ARLDICL:
			opset(ARLDICLCC, r0)
			opset(ARLDICR, r0)
			opset(ARLDICRCC, r0)
			opset(ARLDIC, r0)
			opset(ARLDICCC, r0)
			opset(ACLRLSLDI, r0)

		case AFMOVD:
			opset(AFMOVDCC, r0)
			opset(AFMOVDU, r0)
			opset(AFMOVS, r0)
			opset(AFMOVSU, r0)

		case ALDAR:
			opset(ALBAR, r0)
			opset(ALHAR, r0)
			opset(ALWAR, r0)

		case ASYSCALL: 
			opset(ARFI, r0)

			opset(ARFCI, r0)
			opset(ARFID, r0)
			opset(AHRFID, r0)

		case AMOVHBR:
			opset(AMOVWBR, r0)
			opset(AMOVDBR, r0)

		case ASLBMFEE:
			opset(ASLBMFEV, r0)

		case ATW:
			opset(ATD, r0)

		case ATLBIE:
			opset(ASLBIE, r0)
			opset(ATLBIEL, r0)

		case AEIEIO:
			opset(ASLBIA, r0)

		case ACMP:
			opset(ACMPW, r0)

		case ACMPU:
			opset(ACMPWU, r0)

		case ACMPB:
			opset(ACMPB, r0)

		case AFTDIV:
			opset(AFTDIV, r0)

		case AFTSQRT:
			opset(AFTSQRT, r0)

		case AADD,
			AADDIS,
			AANDCC, 
			AANDISCC,
			AFMOVSX,
			AFMOVSZ,
			ALSW,
			AMOVW,
			
			AMOVWZ, 
			AMOVD,  
			AMOVB,  
			AMOVBU, 
			AMOVFL,
			AMULLW,
			
			ASUBC, 
			ASTSW,
			ASLBMTE,
			AWORD,
			ADWORD,
			ADARN,
			ALDMX,
			AVMSUMUDM,
			AADDEX,
			ACMPEQB,
			AECIWX,
			obj.ANOP,
			obj.ATEXT,
			obj.AUNDEF,
			obj.AFUNCDATA,
			obj.APCALIGN,
			obj.APCDATA,
			obj.ADUFFZERO,
			obj.ADUFFCOPY:
			break
		}
	}
}

func OPVXX1(o uint32, xo uint32, oe uint32) uint32 {
	return o<<26 | xo<<1 | oe<<11
}

func OPVXX2(o uint32, xo uint32, oe uint32) uint32 {
	return o<<26 | xo<<2 | oe<<11
}

func OPVXX2VA(o uint32, xo uint32, oe uint32) uint32 {
	return o<<26 | xo<<2 | oe<<16
}

func OPVXX3(o uint32, xo uint32, oe uint32) uint32 {
	return o<<26 | xo<<3 | oe<<11
}

func OPVXX4(o uint32, xo uint32, oe uint32) uint32 {
	return o<<26 | xo<<4 | oe<<11
}

func OPDQ(o uint32, xo uint32, oe uint32) uint32 {
	return o<<26 | xo | oe<<4
}

func OPVX(o uint32, xo uint32, oe uint32, rc uint32) uint32 {
	return o<<26 | xo | oe<<11 | rc&1
}

func OPVC(o uint32, xo uint32, oe uint32, rc uint32) uint32 {
	return o<<26 | xo | oe<<11 | (rc&1)<<10
}

func OPVCC(o uint32, xo uint32, oe uint32, rc uint32) uint32 {
	return o<<26 | xo<<1 | oe<<10 | rc&1
}

func OPCC(o uint32, xo uint32, rc uint32) uint32 {
	return OPVCC(o, xo, 0, rc)
}


func AOP_RRR(op uint32, d uint32, a uint32, b uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | (b&31)<<11
}


func AOP_RR(op uint32, d uint32, a uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<11
}


func AOP_RRRR(op uint32, d uint32, a uint32, b uint32, c uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | (b&31)<<11 | (c&31)<<6
}

func AOP_IRR(op uint32, d uint32, a uint32, simm uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | simm&0xFFFF
}


func AOP_VIRR(op uint32, d uint32, a uint32, simm uint32) uint32 {
	return op | (d&31)<<21 | (simm&0xFFFF)<<16 | (a&31)<<11
}


func AOP_IIRR(op uint32, d uint32, a uint32, sbit uint32, simm uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | (sbit&1)<<15 | (simm&0xF)<<11
}


func AOP_IRRR(op uint32, d uint32, a uint32, b uint32, simm uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | (b&31)<<11 | (simm&0xF)<<6
}


func AOP_IR(op uint32, d uint32, simm uint32) uint32 {
	return op | (d&31)<<21 | (simm&31)<<16
}


func AOP_XX1(op uint32, d uint32, a uint32, b uint32) uint32 {
	
	
	r := d - REG_VS0
	return op | (r&31)<<21 | (a&31)<<16 | (b&31)<<11 | (r&32)>>5
}


func AOP_XX2(op uint32, d uint32, a uint32, b uint32) uint32 {
	xt := d - REG_VS0
	xb := b - REG_VS0
	return op | (xt&31)<<21 | (a&3)<<16 | (xb&31)<<11 | (xb&32)>>4 | (xt&32)>>5
}


func AOP_XX3(op uint32, d uint32, a uint32, b uint32) uint32 {
	xt := d - REG_VS0
	xa := a - REG_VS0
	xb := b - REG_VS0
	return op | (xt&31)<<21 | (xa&31)<<16 | (xb&31)<<11 | (xa&32)>>3 | (xb&32)>>4 | (xt&32)>>5
}


func AOP_XX3I(op uint32, d uint32, a uint32, b uint32, c uint32) uint32 {
	xt := d - REG_VS0
	xa := a - REG_VS0
	xb := b - REG_VS0
	return op | (xt&31)<<21 | (xa&31)<<16 | (xb&31)<<11 | (c&3)<<8 | (xa&32)>>3 | (xb&32)>>4 | (xt&32)>>5
}


func AOP_XX4(op uint32, d uint32, a uint32, b uint32, c uint32) uint32 {
	xt := d - REG_VS0
	xa := a - REG_VS0
	xb := b - REG_VS0
	xc := c - REG_VS0
	return op | (xt&31)<<21 | (xa&31)<<16 | (xb&31)<<11 | (xc&31)<<6 | (xc&32)>>2 | (xa&32)>>3 | (xb&32)>>4 | (xt&32)>>5
}


func AOP_DQ(op uint32, d uint32, a uint32, b uint32) uint32 {
	
	
	r := d - REG_VS0
	
	
	
	
	
	
	dq := b >> 4
	return op | (r&31)<<21 | (a&31)<<16 | (dq&4095)<<4 | (r&32)>>2
}


func AOP_Z23I(op uint32, d uint32, a uint32, b uint32, c uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | (b&31)<<11 | (c&3)<<7
}


func AOP_RRRI(op uint32, d uint32, a uint32, b uint32, c uint32) uint32 {
	return op | (d&31)<<21 | (a&31)<<16 | (b&31)<<11 | (c & 1)
}

func LOP_RRR(op uint32, a uint32, s uint32, b uint32) uint32 {
	return op | (s&31)<<21 | (a&31)<<16 | (b&31)<<11
}

func LOP_IRR(op uint32, a uint32, s uint32, uimm uint32) uint32 {
	return op | (s&31)<<21 | (a&31)<<16 | uimm&0xFFFF
}

func OP_BR(op uint32, li uint32, aa uint32) uint32 {
	return op | li&0x03FFFFFC | aa<<1
}

func OP_BC(op uint32, bo uint32, bi uint32, bd uint32, aa uint32) uint32 {
	return op | (bo&0x1F)<<21 | (bi&0x1F)<<16 | bd&0xFFFC | aa<<1
}

func OP_BCR(op uint32, bo uint32, bi uint32) uint32 {
	return op | (bo&0x1F)<<21 | (bi&0x1F)<<16
}

func OP_RLW(op uint32, a uint32, s uint32, sh uint32, mb uint32, me uint32) uint32 {
	return op | (s&31)<<21 | (a&31)<<16 | (sh&31)<<11 | (mb&31)<<6 | (me&31)<<1
}

func AOP_RLDIC(op uint32, a uint32, s uint32, sh uint32, m uint32) uint32 {
	return op | (s&31)<<21 | (a&31)<<16 | (sh&31)<<11 | ((sh&32)>>5)<<1 | (m&31)<<6 | ((m&32)>>5)<<5
}

func AOP_ISEL(op uint32, t uint32, a uint32, b uint32, bc uint32) uint32 {
	return op | (t&31)<<21 | (a&31)<<16 | (b&31)<<11 | (bc&0x1F)<<6
}

const (
	
	OP_ADD    = 31<<26 | 266<<1 | 0<<10 | 0
	OP_ADDI   = 14<<26 | 0<<1 | 0<<10 | 0
	OP_ADDIS  = 15<<26 | 0<<1 | 0<<10 | 0
	OP_ANDI   = 28<<26 | 0<<1 | 0<<10 | 0
	OP_EXTSB  = 31<<26 | 954<<1 | 0<<10 | 0
	OP_EXTSH  = 31<<26 | 922<<1 | 0<<10 | 0
	OP_EXTSW  = 31<<26 | 986<<1 | 0<<10 | 0
	OP_ISEL   = 31<<26 | 15<<1 | 0<<10 | 0
	OP_MCRF   = 19<<26 | 0<<1 | 0<<10 | 0
	OP_MCRFS  = 63<<26 | 64<<1 | 0<<10 | 0
	OP_MCRXR  = 31<<26 | 512<<1 | 0<<10 | 0
	OP_MFCR   = 31<<26 | 19<<1 | 0<<10 | 0
	OP_MFFS   = 63<<26 | 583<<1 | 0<<10 | 0
	OP_MFMSR  = 31<<26 | 83<<1 | 0<<10 | 0
	OP_MFSPR  = 31<<26 | 339<<1 | 0<<10 | 0
	OP_MFSR   = 31<<26 | 595<<1 | 0<<10 | 0
	OP_MFSRIN = 31<<26 | 659<<1 | 0<<10 | 0
	OP_MTCRF  = 31<<26 | 144<<1 | 0<<10 | 0
	OP_MTFSF  = 63<<26 | 711<<1 | 0<<10 | 0
	OP_MTFSFI = 63<<26 | 134<<1 | 0<<10 | 0
	OP_MTMSR  = 31<<26 | 146<<1 | 0<<10 | 0
	OP_MTMSRD = 31<<26 | 178<<1 | 0<<10 | 0
	OP_MTSPR  = 31<<26 | 467<<1 | 0<<10 | 0
	OP_MTSR   = 31<<26 | 210<<1 | 0<<10 | 0
	OP_MTSRIN = 31<<26 | 242<<1 | 0<<10 | 0
	OP_MULLW  = 31<<26 | 235<<1 | 0<<10 | 0
	OP_MULLD  = 31<<26 | 233<<1 | 0<<10 | 0
	OP_OR     = 31<<26 | 444<<1 | 0<<10 | 0
	OP_ORI    = 24<<26 | 0<<1 | 0<<10 | 0
	OP_ORIS   = 25<<26 | 0<<1 | 0<<10 | 0
	OP_RLWINM = 21<<26 | 0<<1 | 0<<10 | 0
	OP_RLWNM  = 23<<26 | 0<<1 | 0<<10 | 0
	OP_SUBF   = 31<<26 | 40<<1 | 0<<10 | 0
	OP_RLDIC  = 30<<26 | 4<<1 | 0<<10 | 0
	OP_RLDICR = 30<<26 | 2<<1 | 0<<10 | 0
	OP_RLDICL = 30<<26 | 0<<1 | 0<<10 | 0
	OP_RLDCL  = 30<<26 | 8<<1 | 0<<10 | 0
)

func oclass(a *obj.Addr) int {
	return int(a.Class) - 1
}

const (
	D_FORM = iota
	DS_FORM
)








func (c *ctxt9) opform(insn uint32) int {
	switch insn {
	default:
		c.ctxt.Diag("bad insn in loadform: %x", insn)
	case OPVCC(58, 0, 0, 0), 
		OPVCC(58, 0, 0, 1),        
		OPVCC(58, 0, 0, 0) | 1<<1, 
		OPVCC(62, 0, 0, 0),        
		OPVCC(62, 0, 0, 1):        
		return DS_FORM
	case OP_ADDI, 
		OPVCC(32, 0, 0, 0), 
		OPVCC(33, 0, 0, 0), 
		OPVCC(34, 0, 0, 0), 
		OPVCC(35, 0, 0, 0), 
		OPVCC(40, 0, 0, 0), 
		OPVCC(41, 0, 0, 0), 
		OPVCC(42, 0, 0, 0), 
		OPVCC(43, 0, 0, 0), 
		OPVCC(46, 0, 0, 0), 
		OPVCC(48, 0, 0, 0), 
		OPVCC(49, 0, 0, 0), 
		OPVCC(50, 0, 0, 0), 
		OPVCC(51, 0, 0, 0), 
		OPVCC(36, 0, 0, 0), 
		OPVCC(37, 0, 0, 0), 
		OPVCC(38, 0, 0, 0), 
		OPVCC(39, 0, 0, 0), 
		OPVCC(44, 0, 0, 0), 
		OPVCC(45, 0, 0, 0), 
		OPVCC(47, 0, 0, 0), 
		OPVCC(52, 0, 0, 0), 
		OPVCC(53, 0, 0, 0), 
		OPVCC(54, 0, 0, 0), 
		OPVCC(55, 0, 0, 0): 
		return D_FORM
	}
	return 0
}



func (c *ctxt9) symbolAccess(s *obj.LSym, d int64, reg int16, op uint32) (o1, o2 uint32) {
	if c.ctxt.Headtype == objabi.Haix {
		
		c.ctxt.Diag("symbolAccess called for %s", s.Name)
	}
	var base uint32
	form := c.opform(op)
	if c.ctxt.Flag_shared {
		base = REG_R2
	} else {
		base = REG_R0
	}
	o1 = AOP_IRR(OP_ADDIS, REGTMP, base, 0)
	o2 = AOP_IRR(op, uint32(reg), REGTMP, 0)
	rel := obj.Addrel(c.cursym)
	rel.Off = int32(c.pc)
	rel.Siz = 8
	rel.Sym = s
	rel.Add = d
	if c.ctxt.Flag_shared {
		switch form {
		case D_FORM:
			rel.Type = objabi.R_ADDRPOWER_TOCREL
		case DS_FORM:
			rel.Type = objabi.R_ADDRPOWER_TOCREL_DS
		}

	} else {
		switch form {
		case D_FORM:
			rel.Type = objabi.R_ADDRPOWER
		case DS_FORM:
			rel.Type = objabi.R_ADDRPOWER_DS
		}
	}
	return
}


func getmask(m []byte, v uint32) bool {
	m[1] = 0
	m[0] = m[1]
	if v != ^uint32(0) && v&(1<<31) != 0 && v&1 != 0 { 
		if getmask(m, ^v) {
			i := int(m[0])
			m[0] = m[1] + 1
			m[1] = byte(i - 1)
			return true
		}

		return false
	}

	for i := 0; i < 32; i++ {
		if v&(1<<uint(31-i)) != 0 {
			m[0] = byte(i)
			for {
				m[1] = byte(i)
				i++
				if i >= 32 || v&(1<<uint(31-i)) == 0 {
					break
				}
			}

			for ; i < 32; i++ {
				if v&(1<<uint(31-i)) != 0 {
					return false
				}
			}
			return true
		}
	}

	return false
}

func (c *ctxt9) maskgen(p *obj.Prog, m []byte, v uint32) {
	if !getmask(m, v) {
		c.ctxt.Diag("cannot generate mask #%x\n%v", v, p)
	}
}


func getmask64(m []byte, v uint64) bool {
	m[1] = 0
	m[0] = m[1]
	for i := 0; i < 64; i++ {
		if v&(uint64(1)<<uint(63-i)) != 0 {
			m[0] = byte(i)
			for {
				m[1] = byte(i)
				i++
				if i >= 64 || v&(uint64(1)<<uint(63-i)) == 0 {
					break
				}
			}

			for ; i < 64; i++ {
				if v&(uint64(1)<<uint(63-i)) != 0 {
					return false
				}
			}
			return true
		}
	}

	return false
}

func (c *ctxt9) maskgen64(p *obj.Prog, m []byte, v uint64) {
	if !getmask64(m, v) {
		c.ctxt.Diag("cannot generate mask #%x\n%v", v, p)
	}
}

func loadu32(r int, d int64) uint32 {
	v := int32(d >> 16)
	if isuint32(uint64(d)) {
		return LOP_IRR(OP_ORIS, uint32(r), REGZERO, uint32(v))
	}
	return AOP_IRR(OP_ADDIS, uint32(r), REGZERO, uint32(v))
}

func high16adjusted(d int32) uint16 {
	if d&0x8000 != 0 {
		return uint16((d >> 16) + 1)
	}
	return uint16(d >> 16)
}

func (c *ctxt9) asmout(p *obj.Prog, o *Optab, out []uint32) {
	o1 := uint32(0)
	o2 := uint32(0)
	o3 := uint32(0)
	o4 := uint32(0)
	o5 := uint32(0)

	
	switch o.type_ {
	default:
		c.ctxt.Diag("unknown type %d", o.type_)
		prasm(p)

	case 0: 
		break

	case 1: 
		if p.To.Reg == REGZERO && p.From.Type == obj.TYPE_CONST {
			v := c.regoff(&p.From)
			if r0iszero != 0  && v != 0 {
				
				c.ctxt.Diag("literal operation on R0\n%v", p)
			}

			o1 = LOP_IRR(OP_ADDI, REGZERO, REGZERO, uint32(v))
			break
		}

		o1 = LOP_RRR(OP_OR, uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.From.Reg))

	case 2: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), uint32(p.From.Reg))

	case 3: 
		d := c.vregoff(&p.From)

		v := int32(d)
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		if r0iszero != 0  && p.To.Reg == 0 && (r != 0 || v != 0) {
			c.ctxt.Diag("literal operation on R0\n%v", p)
		}
		a := OP_ADDI
		if o.a1 == C_UCON {
			if d&0xffff != 0 {
				log.Fatalf("invalid handling of %v", p)
			}
			
			
			v >>= 16
			if r == REGZERO && isuint32(uint64(d)) {
				o1 = LOP_IRR(OP_ORIS, uint32(p.To.Reg), REGZERO, uint32(v))
				break
			}

			a = OP_ADDIS
		} else if int64(int16(d)) != d {
			
			if o.a1 == C_ANDCON {
				
				if r == 0 || r == REGZERO {
					o1 = LOP_IRR(uint32(OP_ORI), uint32(p.To.Reg), uint32(0), uint32(v))
					break
				}
				
			} else if o.a1 != C_ADDCON {
				log.Fatalf("invalid handling of %v", p)
			}
		}

		o1 = AOP_IRR(uint32(a), uint32(p.To.Reg), uint32(r), uint32(v))

	case 4: 
		v := c.regoff(&p.From)

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		if r0iszero != 0  && p.To.Reg == 0 {
			c.ctxt.Diag("literal operation on R0\n%v", p)
		}
		if int32(int16(v)) != v {
			log.Fatalf("mishandled instruction %v", p)
		}
		o1 = AOP_IRR(c.opirr(p.As), uint32(p.To.Reg), uint32(r), uint32(v))

	case 5: 
		o1 = c.oprrr(p.As)

	case 6: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		
		switch p.As {
		case AROTL:
			o1 = AOP_RLDIC(OP_RLDCL, uint32(p.To.Reg), uint32(r), uint32(p.From.Reg), uint32(0))
		case AROTLW:
			o1 = OP_RLW(OP_RLWNM, uint32(p.To.Reg), uint32(r), uint32(p.From.Reg), 0, 31)
		default:
			o1 = LOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), uint32(p.From.Reg))
		}

	case 7: 
		r := int(p.To.Reg)

		if r == 0 {
			r = int(o.param)
		}
		v := c.regoff(&p.To)
		if p.To.Type == obj.TYPE_MEM && p.To.Index != 0 {
			if v != 0 {
				c.ctxt.Diag("illegal indexed instruction\n%v", p)
			}
			if c.ctxt.Flag_shared && r == REG_R13 {
				rel := obj.Addrel(c.cursym)
				rel.Off = int32(c.pc)
				rel.Siz = 4
				
				
				
				
				
				rel.Sym = c.ctxt.Lookup("runtime.tls_g")
				rel.Type = objabi.R_POWER_TLS
			}
			o1 = AOP_RRR(c.opstorex(p.As), uint32(p.From.Reg), uint32(p.To.Index), uint32(r))
		} else {
			if int32(int16(v)) != v {
				log.Fatalf("mishandled instruction %v", p)
			}
			
			inst := c.opstore(p.As)
			if c.opform(inst) == DS_FORM && v&0x3 != 0 {
				log.Fatalf("invalid offset for DS form load/store %v", p)
			}
			o1 = AOP_IRR(inst, uint32(p.From.Reg), uint32(r), uint32(v))
		}

	case 8: 
		r := int(p.From.Reg)

		if r == 0 {
			r = int(o.param)
		}
		v := c.regoff(&p.From)
		if p.From.Type == obj.TYPE_MEM && p.From.Index != 0 {
			if v != 0 {
				c.ctxt.Diag("illegal indexed instruction\n%v", p)
			}
			if c.ctxt.Flag_shared && r == REG_R13 {
				rel := obj.Addrel(c.cursym)
				rel.Off = int32(c.pc)
				rel.Siz = 4
				rel.Sym = c.ctxt.Lookup("runtime.tls_g")
				rel.Type = objabi.R_POWER_TLS
			}
			o1 = AOP_RRR(c.oploadx(p.As), uint32(p.To.Reg), uint32(p.From.Index), uint32(r))
		} else {
			if int32(int16(v)) != v {
				log.Fatalf("mishandled instruction %v", p)
			}
			
			inst := c.opload(p.As)
			if c.opform(inst) == DS_FORM && v&0x3 != 0 {
				log.Fatalf("invalid offset for DS form load/store %v", p)
			}
			o1 = AOP_IRR(inst, uint32(p.To.Reg), uint32(r), uint32(v))
		}

	case 9: 
		r := int(p.From.Reg)

		if r == 0 {
			r = int(o.param)
		}
		v := c.regoff(&p.From)
		if p.From.Type == obj.TYPE_MEM && p.From.Index != 0 {
			if v != 0 {
				c.ctxt.Diag("illegal indexed instruction\n%v", p)
			}
			o1 = AOP_RRR(c.oploadx(p.As), uint32(p.To.Reg), uint32(p.From.Index), uint32(r))
		} else {
			o1 = AOP_IRR(c.opload(p.As), uint32(p.To.Reg), uint32(r), uint32(v))
		}
		o2 = LOP_RRR(OP_EXTSB, uint32(p.To.Reg), uint32(p.To.Reg), 0)

	case 10: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(r))

	case 11: 
		v := int32(0)

		if p.To.Target() != nil {
			v = int32(p.To.Target().Pc - p.Pc)
			if v&03 != 0 {
				c.ctxt.Diag("odd branch target address\n%v", p)
				v &^= 03
			}

			if v < -(1<<25) || v >= 1<<24 {
				c.ctxt.Diag("branch too far\n%v", p)
			}
		}

		o1 = OP_BR(c.opirr(p.As), uint32(v), 0)
		if p.To.Sym != nil {
			rel := obj.Addrel(c.cursym)
			rel.Off = int32(c.pc)
			rel.Siz = 4
			rel.Sym = p.To.Sym
			v += int32(p.To.Offset)
			if v&03 != 0 {
				c.ctxt.Diag("odd branch target address\n%v", p)
				v &^= 03
			}

			rel.Add = int64(v)
			rel.Type = objabi.R_CALLPOWER
		}
		o2 = 0x60000000 

	case 12: 
		if p.To.Reg == REGZERO && p.From.Type == obj.TYPE_CONST {
			v := c.regoff(&p.From)
			if r0iszero != 0  && v != 0 {
				c.ctxt.Diag("literal operation on R0\n%v", p)
			}

			o1 = LOP_IRR(OP_ADDI, REGZERO, REGZERO, uint32(v))
			break
		}

		if p.As == AMOVW {
			o1 = LOP_RRR(OP_EXTSW, uint32(p.To.Reg), uint32(p.From.Reg), 0)
		} else {
			o1 = LOP_RRR(OP_EXTSB, uint32(p.To.Reg), uint32(p.From.Reg), 0)
		}

	case 13: 
		if p.As == AMOVBZ {
			o1 = OP_RLW(OP_RLWINM, uint32(p.To.Reg), uint32(p.From.Reg), 0, 24, 31)
		} else if p.As == AMOVH {
			o1 = LOP_RRR(OP_EXTSH, uint32(p.To.Reg), uint32(p.From.Reg), 0)
		} else if p.As == AMOVHZ {
			o1 = OP_RLW(OP_RLWINM, uint32(p.To.Reg), uint32(p.From.Reg), 0, 16, 31)
		} else if p.As == AMOVWZ {
			o1 = OP_RLW(OP_RLDIC, uint32(p.To.Reg), uint32(p.From.Reg), 0, 0, 0) | 1<<5 
		} else {
			c.ctxt.Diag("internal: bad mov[bhw]z\n%v", p)
		}

	case 14: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		d := c.vregoff(p.GetFrom3())
		var a int
		switch p.As {

		
		
		
		case ARLDCL, ARLDCLCC:
			var mask [2]uint8
			c.maskgen64(p, mask[:], uint64(d))

			a = int(mask[0]) 
			if mask[1] != 63 {
				c.ctxt.Diag("invalid mask for rotate: %x (end != bit 63)\n%v", uint64(d), p)
			}
			o1 = LOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), uint32(p.From.Reg))
			o1 |= (uint32(a) & 31) << 6
			if a&0x20 != 0 {
				o1 |= 1 << 5 
			}

		case ARLDCR, ARLDCRCC:
			var mask [2]uint8
			c.maskgen64(p, mask[:], uint64(d))

			a = int(mask[1]) 
			if mask[0] != 0 {
				c.ctxt.Diag("invalid mask for rotate: %x %x (start != 0)\n%v", uint64(d), mask[0], p)
			}
			o1 = LOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), uint32(p.From.Reg))
			o1 |= (uint32(a) & 31) << 6
			if a&0x20 != 0 {
				o1 |= 1 << 5 
			}

		
		case ARLDICR, ARLDICRCC:
			me := int(d)
			sh := c.regoff(&p.From)
			if me < 0 || me > 63 || sh > 63 {
				c.ctxt.Diag("Invalid me or sh for RLDICR: %x %x\n%v", int(d), sh)
			}
			o1 = AOP_RLDIC(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), uint32(sh), uint32(me))

		case ARLDICL, ARLDICLCC, ARLDIC, ARLDICCC:
			mb := int(d)
			sh := c.regoff(&p.From)
			if mb < 0 || mb > 63 || sh > 63 {
				c.ctxt.Diag("Invalid mb or sh for RLDIC, RLDICL: %x %x\n%v", mb, sh)
			}
			o1 = AOP_RLDIC(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), uint32(sh), uint32(mb))

		case ACLRLSLDI:
			
			
			
			
			b := int(d)
			n := c.regoff(&p.From)
			if n > int32(b) || b > 63 {
				c.ctxt.Diag("Invalid n or b for CLRLSLDI: %x %x\n%v", n, b)
			}
			o1 = AOP_RLDIC(OP_RLDIC, uint32(p.To.Reg), uint32(r), uint32(n), uint32(b)-uint32(n))

		default:
			c.ctxt.Diag("unexpected op in rldc case\n%v", p)
			a = 0
		}

	case 17, 
		16: 
		a := 0

		r := int(p.Reg)

		if p.From.Type == obj.TYPE_CONST {
			a = int(c.regoff(&p.From))
		} else if p.From.Type == obj.TYPE_REG {
			if r != 0 {
				c.ctxt.Diag("unexpected register setting for branch with CR: %d\n", r)
			}
			
			switch p.From.Reg {
			case REG_CR0:
				r = BI_CR0
			case REG_CR1:
				r = BI_CR1
			case REG_CR2:
				r = BI_CR2
			case REG_CR3:
				r = BI_CR3
			case REG_CR4:
				r = BI_CR4
			case REG_CR5:
				r = BI_CR5
			case REG_CR6:
				r = BI_CR6
			case REG_CR7:
				r = BI_CR7
			default:
				c.ctxt.Diag("unrecognized register: expecting CR\n")
			}
		}
		v := int32(0)
		if p.To.Target() != nil {
			v = int32(p.To.Target().Pc - p.Pc)
		}
		if v&03 != 0 {
			c.ctxt.Diag("odd branch target address\n%v", p)
			v &^= 03
		}

		if v < -(1<<16) || v >= 1<<15 {
			c.ctxt.Diag("branch too far\n%v", p)
		}
		o1 = OP_BC(c.opirr(p.As), uint32(a), uint32(r), uint32(v), 0)

	case 15: 
		var v int32
		if p.As == ABC || p.As == ABCL {
			v = c.regoff(&p.To) & 31
		} else {
			v = 20 
		}
		o1 = AOP_RRR(OP_MTSPR, uint32(p.To.Reg), 0, 0) | (REG_LR&0x1f)<<16 | ((REG_LR>>5)&0x1f)<<11
		o2 = OPVCC(19, 16, 0, 0)
		if p.As == ABL || p.As == ABCL {
			o2 |= 1
		}
		o2 = OP_BCR(o2, uint32(v), uint32(p.To.Index))

	case 18: 
		var v int32
		if p.As == ABC || p.As == ABCL {
			v = c.regoff(&p.From) & 31
		} else {
			v = 20 
		}
		r := int(p.Reg)
		if r == 0 {
			r = 0
		}
		switch oclass(&p.To) {
		case C_CTR:
			o1 = OPVCC(19, 528, 0, 0)

		case C_LR:
			o1 = OPVCC(19, 16, 0, 0)

		default:
			c.ctxt.Diag("bad optab entry (18): %d\n%v", p.To.Class, p)
			v = 0
		}

		if p.As == ABL || p.As == ABCL {
			o1 |= 1
		}
		o1 = OP_BCR(o1, uint32(v), uint32(r))

	case 19: 
		d := c.vregoff(&p.From)

		if p.From.Sym == nil {
			o1 = loadu32(int(p.To.Reg), d)
			o2 = LOP_IRR(OP_ORI, uint32(p.To.Reg), uint32(p.To.Reg), uint32(int32(d)))
		} else {
			o1, o2 = c.symbolAccess(p.From.Sym, d, p.To.Reg, OP_ADDI)
		}

	case 20: 
		v := c.regoff(&p.From)

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		if p.As == AADD && (r0iszero == 0  && p.Reg == 0 || r0iszero != 0  && p.To.Reg == 0) {
			c.ctxt.Diag("literal operation on R0\n%v", p)
		}
		if p.As == AADDIS {
			o1 = AOP_IRR(c.opirr(p.As), uint32(p.To.Reg), uint32(r), uint32(v))
		} else {
			o1 = AOP_IRR(c.opirr(AADDIS), uint32(p.To.Reg), uint32(r), uint32(v)>>16)
		}

	case 22: 
		if p.To.Reg == REGTMP || p.Reg == REGTMP {
			c.ctxt.Diag("can't synthesize large constant\n%v", p)
		}
		d := c.vregoff(&p.From)
		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		if p.From.Sym != nil {
			c.ctxt.Diag("%v is not supported", p)
		}
		
		
		if o.size == 8 {
			o1 = LOP_IRR(OP_ORI, REGTMP, REGZERO, uint32(int32(d)))
			o2 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), REGTMP, uint32(r))
		} else {
			o1 = loadu32(REGTMP, d)
			o2 = LOP_IRR(OP_ORI, REGTMP, REGTMP, uint32(int32(d)))
			o3 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), REGTMP, uint32(r))
		}

	case 23: 
		if p.To.Reg == REGTMP || p.Reg == REGTMP {
			c.ctxt.Diag("can't synthesize large constant\n%v", p)
		}
		d := c.vregoff(&p.From)
		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}

		
		
		if o.size == 8 {
			o1 = LOP_IRR(OP_ADDI, REGZERO, REGTMP, uint32(int32(d)))
			o2 = LOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), REGTMP, uint32(r))
		} else {
			o1 = loadu32(REGTMP, d)
			o2 = LOP_IRR(OP_ORI, REGTMP, REGTMP, uint32(int32(d)))
			o3 = LOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), REGTMP, uint32(r))
		}
		if p.From.Sym != nil {
			c.ctxt.Diag("%v is not supported", p)
		}

	case 24: 
		o1 = AOP_XX3I(c.oprrr(AXXLXOR), uint32(p.To.Reg), uint32(p.To.Reg), uint32(p.To.Reg), uint32(0))
		
		if o.size == 8 {
			o2 = AOP_RRR(c.oprrr(AFNEG), uint32(p.To.Reg), 0, uint32(p.To.Reg))
		}

	case 25:
		
		v := c.regoff(&p.From)

		if v < 0 {
			v = 0
		} else if v > 63 {
			v = 63
		}
		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		var a int
		op := uint32(0)
		switch p.As {
		case ASLD, ASLDCC:
			a = int(63 - v)
			op = OP_RLDICR

		case ASRD, ASRDCC:
			a = int(v)
			v = 64 - v
			op = OP_RLDICL
		case AROTL:
			a = int(0)
			op = OP_RLDICL
		default:
			c.ctxt.Diag("unexpected op in sldi case\n%v", p)
			a = 0
			o1 = 0
		}

		o1 = AOP_RLDIC(op, uint32(p.To.Reg), uint32(r), uint32(v), uint32(a))
		if p.As == ASLDCC || p.As == ASRDCC {
			o1 |= 1 
		}

	case 26: 
		if p.To.Reg == REGTMP {
			c.ctxt.Diag("can't synthesize large constant\n%v", p)
		}
		v := c.regoff(&p.From)
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = AOP_IRR(OP_ADDIS, REGTMP, uint32(r), uint32(high16adjusted(v)))
		o2 = AOP_IRR(OP_ADDI, uint32(p.To.Reg), REGTMP, uint32(v))

	case 27: 
		v := c.regoff(p.GetFrom3())

		r := int(p.From.Reg)
		o1 = AOP_IRR(c.opirr(p.As), uint32(p.To.Reg), uint32(r), uint32(v))

	case 28: 
		if p.To.Reg == REGTMP || p.From.Reg == REGTMP {
			c.ctxt.Diag("can't synthesize large constant\n%v", p)
		}
		v := c.regoff(p.GetFrom3())
		o1 = AOP_IRR(OP_ADDIS, REGTMP, REGZERO, uint32(v)>>16)
		o2 = LOP_IRR(OP_ORI, REGTMP, REGTMP, uint32(v))
		o3 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), REGTMP)
		if p.From.Sym != nil {
			c.ctxt.Diag("%v is not supported", p)
		}

	case 29: 
		v := c.regoff(&p.From)

		d := c.vregoff(p.GetFrom3())
		var mask [2]uint8
		c.maskgen64(p, mask[:], uint64(d))
		var a int
		switch p.As {
		case ARLDC, ARLDCCC:
			a = int(mask[0]) 
			if int32(mask[1]) != (63 - v) {
				c.ctxt.Diag("invalid mask for shift: %x %x (shift %d)\n%v", uint64(d), mask[1], v, p)
			}

		case ARLDCL, ARLDCLCC:
			a = int(mask[0]) 
			if mask[1] != 63 {
				c.ctxt.Diag("invalid mask for shift: %x %s (shift %d)\n%v", uint64(d), mask[1], v, p)
			}

		case ARLDCR, ARLDCRCC:
			a = int(mask[1]) 
			if mask[0] != 0 {
				c.ctxt.Diag("invalid mask for shift: %x %x (shift %d)\n%v", uint64(d), mask[0], v, p)
			}

		default:
			c.ctxt.Diag("unexpected op in rldic case\n%v", p)
			a = 0
		}

		o1 = AOP_RRR(c.opirr(p.As), uint32(p.Reg), uint32(p.To.Reg), (uint32(v) & 0x1F))
		o1 |= (uint32(a) & 31) << 6
		if v&0x20 != 0 {
			o1 |= 1 << 1
		}
		if a&0x20 != 0 {
			o1 |= 1 << 5 
		}

	case 30: 
		v := c.regoff(&p.From)

		d := c.vregoff(p.GetFrom3())

		
		
		switch p.As {
		case ARLDMI, ARLDMICC:
			var mask [2]uint8
			c.maskgen64(p, mask[:], uint64(d))
			if int32(mask[1]) != (63 - v) {
				c.ctxt.Diag("invalid mask for shift: %x %x (shift %d)\n%v", uint64(d), mask[1], v, p)
			}
			o1 = AOP_RRR(c.opirr(p.As), uint32(p.Reg), uint32(p.To.Reg), (uint32(v) & 0x1F))
			o1 |= (uint32(mask[0]) & 31) << 6
			if v&0x20 != 0 {
				o1 |= 1 << 1
			}
			if mask[0]&0x20 != 0 {
				o1 |= 1 << 5 
			}

		
		case ARLDIMI, ARLDIMICC:
			o1 = AOP_RRR(c.opirr(p.As), uint32(p.Reg), uint32(p.To.Reg), (uint32(v) & 0x1F))
			o1 |= (uint32(d) & 31) << 6
			if d&0x20 != 0 {
				o1 |= 1 << 5
			}
			if v&0x20 != 0 {
				o1 |= 1 << 1
			}
		}

	case 31: 
		d := c.vregoff(&p.From)

		if c.ctxt.Arch.ByteOrder == binary.BigEndian {
			o1 = uint32(d >> 32)
			o2 = uint32(d)
		} else {
			o1 = uint32(d)
			o2 = uint32(d >> 32)
		}

		if p.From.Sym != nil {
			rel := obj.Addrel(c.cursym)
			rel.Off = int32(c.pc)
			rel.Siz = 8
			rel.Sym = p.From.Sym
			rel.Add = p.From.Offset
			rel.Type = objabi.R_ADDR
			o2 = 0
			o1 = o2
		}

	case 32: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), 0) | (uint32(p.From.Reg)&31)<<6

	case 33: 
		r := int(p.From.Reg)

		if oclass(&p.From) == C_NONE {
			r = int(p.To.Reg)
		}
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), 0, uint32(r))

	case 34: 
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg)) | (uint32(p.GetFrom3().Reg)&31)<<6

	case 35: 
		v := c.regoff(&p.To)

		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		
		inst := c.opstore(p.As)
		if c.opform(inst) == DS_FORM && v&0x3 != 0 {
			log.Fatalf("invalid offset for DS form load/store %v", p)
		}
		o1 = AOP_IRR(OP_ADDIS, REGTMP, uint32(r), uint32(high16adjusted(v)))
		o2 = AOP_IRR(inst, uint32(p.From.Reg), REGTMP, uint32(v))

	case 36: 
		v := c.regoff(&p.From)

		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = AOP_IRR(OP_ADDIS, REGTMP, uint32(r), uint32(high16adjusted(v)))
		o2 = AOP_IRR(c.opload(p.As), uint32(p.To.Reg), REGTMP, uint32(v))

	case 37: 
		v := c.regoff(&p.From)

		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = AOP_IRR(OP_ADDIS, REGTMP, uint32(r), uint32(high16adjusted(v)))
		o2 = AOP_IRR(c.opload(p.As), uint32(p.To.Reg), REGTMP, uint32(v))
		o3 = LOP_RRR(OP_EXTSB, uint32(p.To.Reg), uint32(p.To.Reg), 0)

	case 40: 
		o1 = uint32(c.regoff(&p.From))

	case 41: 
		o1 = AOP_RRR(c.opirr(p.As), uint32(p.From.Reg), uint32(p.To.Reg), 0) | (uint32(c.regoff(p.GetFrom3()))&0x7F)<<11

	case 42: 
		o1 = AOP_RRR(c.opirr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), 0) | (uint32(c.regoff(p.GetFrom3()))&0x7F)<<11

	case 43: 
		
		
		
		
		

		
		
		
		
		if p.To.Type == obj.TYPE_NONE {
			o1 = AOP_RRR(c.oprrr(p.As), 0, uint32(p.From.Index), uint32(p.From.Reg))
		} else {
			th := c.regoff(&p.To)
			o1 = AOP_RRR(c.oprrr(p.As), uint32(th), uint32(p.From.Index), uint32(p.From.Reg))
		}

	case 44: 
		o1 = AOP_RRR(c.opstorex(p.As), uint32(p.From.Reg), uint32(p.To.Index), uint32(p.To.Reg))

	case 45: 
		switch p.As {
		
		
		
		
		case ALBAR, ALHAR, ALWAR, ALDAR:
			if p.From3Type() != obj.TYPE_NONE {
				eh := int(c.regoff(p.GetFrom3()))
				if eh > 1 {
					c.ctxt.Diag("illegal EH field\n%v", p)
				}
				o1 = AOP_RRRI(c.oploadx(p.As), uint32(p.To.Reg), uint32(p.From.Index), uint32(p.From.Reg), uint32(eh))
			} else {
				o1 = AOP_RRR(c.oploadx(p.As), uint32(p.To.Reg), uint32(p.From.Index), uint32(p.From.Reg))
			}
		default:
			o1 = AOP_RRR(c.oploadx(p.As), uint32(p.To.Reg), uint32(p.From.Index), uint32(p.From.Reg))
		}
	case 46: 
		o1 = c.oprrr(p.As)

	case 47: 
		r := int(p.From.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), 0)

	case 48: 
		r := int(p.From.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = LOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(r), 0)

	case 49: 
		if p.From.Type != obj.TYPE_REG { 
			v := c.regoff(&p.From) & 1
			o1 = AOP_RRR(c.oprrr(p.As), 0, 0, uint32(p.To.Reg)) | uint32(v)<<21
		} else {
			o1 = AOP_RRR(c.oprrr(p.As), 0, 0, uint32(p.From.Reg))
		}

	case 50: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		v := c.oprrr(p.As)
		t := v & (1<<10 | 1) 
		o1 = AOP_RRR(v&^t, REGTMP, uint32(r), uint32(p.From.Reg))
		o2 = AOP_RRR(OP_MULLW, REGTMP, REGTMP, uint32(p.From.Reg))
		o3 = AOP_RRR(OP_SUBF|t, uint32(p.To.Reg), REGTMP, uint32(r))
		if p.As == AREMU {
			o4 = o3

			
			o3 = OP_RLW(OP_RLDIC, REGTMP, REGTMP, 0, 0, 0) | 1<<5
		}

	case 51: 
		r := int(p.Reg)

		if r == 0 {
			r = int(p.To.Reg)
		}
		v := c.oprrr(p.As)
		t := v & (1<<10 | 1) 
		o1 = AOP_RRR(v&^t, REGTMP, uint32(r), uint32(p.From.Reg))
		o2 = AOP_RRR(OP_MULLD, REGTMP, REGTMP, uint32(p.From.Reg))
		o3 = AOP_RRR(OP_SUBF|t, uint32(p.To.Reg), REGTMP, uint32(r))
		

		

	case 52: 
		v := c.regoff(&p.From) & 31

		o1 = AOP_RRR(c.oprrr(p.As), uint32(v), 0, 0)

	case 53: 
		o1 = AOP_RRR(OP_MFFS, uint32(p.To.Reg), 0, 0)

	case 54: 
		if oclass(&p.From) == C_REG {
			if p.As == AMOVD {
				o1 = AOP_RRR(OP_MTMSRD, uint32(p.From.Reg), 0, 0)
			} else {
				o1 = AOP_RRR(OP_MTMSR, uint32(p.From.Reg), 0, 0)
			}
		} else {
			o1 = AOP_RRR(OP_MFMSR, uint32(p.To.Reg), 0, 0)
		}

	case 55: 
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), 0, uint32(p.From.Reg))

	case 56: 
		v := c.regoff(&p.From)

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = AOP_RRR(c.opirr(p.As), uint32(r), uint32(p.To.Reg), uint32(v)&31)
		if (p.As == ASRAD || p.As == ASRADCC) && (v&0x20 != 0) {
			o1 |= 1 << 1 
		}

	case 57: 
		v := c.regoff(&p.From)

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}

		
		if v < 0 {
			v = 0
		} else if v > 32 {
			v = 32
		}
		var mask [2]uint8
		switch p.As {
		case AROTLW:
			mask[0], mask[1] = 0, 31
		case ASRW, ASRWCC:
			mask[0], mask[1] = uint8(v), 31
			v = 32 - v
		default:
			mask[0], mask[1] = 0, uint8(31-v)
		}
		o1 = OP_RLW(OP_RLWINM, uint32(p.To.Reg), uint32(r), uint32(v), uint32(mask[0]), uint32(mask[1]))
		if p.As == ASLWCC || p.As == ASRWCC {
			o1 |= 1 
		}

	case 58: 
		v := c.regoff(&p.From)

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 = LOP_IRR(c.opirr(p.As), uint32(p.To.Reg), uint32(r), uint32(v))

	case 59: 
		v := c.regoff(&p.From)

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		switch p.As {
		case AOR:
			o1 = LOP_IRR(c.opirr(AORIS), uint32(p.To.Reg), uint32(r), uint32(v)>>16) 
		case AXOR:
			o1 = LOP_IRR(c.opirr(AXORIS), uint32(p.To.Reg), uint32(r), uint32(v)>>16)
		case AANDCC:
			o1 = LOP_IRR(c.opirr(AANDISCC), uint32(p.To.Reg), uint32(r), uint32(v)>>16)
		default:
			o1 = LOP_IRR(c.opirr(p.As), uint32(p.To.Reg), uint32(r), uint32(v))
		}

	case 60: 
		r := int(c.regoff(&p.From) & 31)

		o1 = AOP_RRR(c.oprrr(p.As), uint32(r), uint32(p.Reg), uint32(p.To.Reg))

	case 61: 
		r := int(c.regoff(&p.From) & 31)

		v := c.regoff(&p.To)
		o1 = AOP_IRR(c.opirr(p.As), uint32(r), uint32(p.Reg), uint32(v))

	case 62: 
		v := c.regoff(&p.From)
		switch p.As {
		case ACLRLSLWI:
			b := c.regoff(p.GetFrom3())
			
			
			
			if v < 0 || v > 32 || b > 32 {
				c.ctxt.Diag("Invalid n or b for CLRLSLWI: %x %x\n%v", v, b)
			}
			o1 = OP_RLW(OP_RLWINM, uint32(p.To.Reg), uint32(p.Reg), uint32(v), uint32(b-v), uint32(31-v))
		default:
			var mask [2]uint8
			c.maskgen(p, mask[:], uint32(c.regoff(p.GetFrom3())))
			o1 = AOP_RRR(c.opirr(p.As), uint32(p.Reg), uint32(p.To.Reg), uint32(v))
			o1 |= (uint32(mask[0])&31)<<6 | (uint32(mask[1])&31)<<1
		}

	case 63: 
		v := c.regoff(&p.From)
		switch p.As {
		case ACLRLSLWI:
			b := c.regoff(p.GetFrom3())
			if v > b || b > 32 {
				
				
				c.ctxt.Diag("Invalid n or b for CLRLSLWI: %x %x\n%v", v, b)
			}
			
			
			
			o1 = OP_RLW(OP_RLWINM, uint32(p.To.Reg), uint32(p.Reg), uint32(v), uint32(b-v), uint32(31-v))
		default:
			var mask [2]uint8
			c.maskgen(p, mask[:], uint32(c.regoff(p.GetFrom3())))
			o1 = AOP_RRR(c.opirr(p.As), uint32(p.Reg), uint32(p.To.Reg), uint32(v))
			o1 |= (uint32(mask[0])&31)<<6 | (uint32(mask[1])&31)<<1
		}

	case 64: 
		var v int32
		if p.From3Type() != obj.TYPE_NONE {
			v = c.regoff(p.GetFrom3()) & 255
		} else {
			v = 255
		}
		o1 = OP_MTFSF | uint32(v)<<17 | uint32(p.From.Reg)<<11

	case 65: 
		if p.To.Reg == 0 {
			c.ctxt.Diag("must specify FPSCR(n)\n%v", p)
		}
		o1 = OP_MTFSFI | (uint32(p.To.Reg)&15)<<23 | (uint32(c.regoff(&p.From))&31)<<12

	case 66: 
		var r int
		var v int32
		if REG_R0 <= p.From.Reg && p.From.Reg <= REG_R31 {
			r = int(p.From.Reg)
			v = int32(p.To.Reg)
			if REG_DCR0 <= v && v <= REG_DCR0+1023 {
				o1 = OPVCC(31, 451, 0, 0) 
			} else {
				o1 = OPVCC(31, 467, 0, 0) 
			}
		} else {
			r = int(p.To.Reg)
			v = int32(p.From.Reg)
			if REG_DCR0 <= v && v <= REG_DCR0+1023 {
				o1 = OPVCC(31, 323, 0, 0) 
			} else {
				o1 = OPVCC(31, 339, 0, 0) 
			}
		}

		o1 = AOP_RRR(o1, uint32(r), 0, 0) | (uint32(v)&0x1f)<<16 | ((uint32(v)>>5)&0x1f)<<11

	case 67: 
		if p.From.Type != obj.TYPE_REG || p.From.Reg < REG_CR0 || REG_CR7 < p.From.Reg || p.To.Type != obj.TYPE_REG || p.To.Reg < REG_CR0 || REG_CR7 < p.To.Reg {
			c.ctxt.Diag("illegal CR field number\n%v", p)
		}
		o1 = AOP_RRR(OP_MCRF, ((uint32(p.To.Reg) & 7) << 2), ((uint32(p.From.Reg) & 7) << 2), 0)

	case 68: 
		if p.From.Type == obj.TYPE_REG && REG_CR0 <= p.From.Reg && p.From.Reg <= REG_CR7 {
			v := int32(1 << uint(7-(p.To.Reg&7)))                                 
			o1 = AOP_RRR(OP_MFCR, uint32(p.To.Reg), 0, 0) | 1<<20 | uint32(v)<<12 
		} else {
			o1 = AOP_RRR(OP_MFCR, uint32(p.To.Reg), 0, 0) 
		}

	case 69: 
		var v int32
		if p.From3Type() != obj.TYPE_NONE {
			if p.To.Reg != 0 {
				c.ctxt.Diag("can't use both mask and CR(n)\n%v", p)
			}
			v = c.regoff(p.GetFrom3()) & 0xff
		} else {
			if p.To.Reg == 0 {
				v = 0xff 
			} else {
				v = 1 << uint(7-(p.To.Reg&7)) 
			}
		}

		o1 = AOP_RRR(OP_MTCRF, uint32(p.From.Reg), 0, 0) | uint32(v)<<12

	case 70: 
		var r int
		if p.Reg == 0 {
			r = 0
		} else {
			r = (int(p.Reg) & 7) << 2
		}
		o1 = AOP_RRR(c.oprrr(p.As), uint32(r), uint32(p.From.Reg), uint32(p.To.Reg))

	case 71: 
		var r int
		if p.Reg == 0 {
			r = 0
		} else {
			r = (int(p.Reg) & 7) << 2
		}
		o1 = AOP_RRR(c.opirr(p.As), uint32(r), uint32(p.From.Reg), 0) | uint32(c.regoff(&p.To))&0xffff

	case 72: 
		o1 = AOP_RRR(c.oprrr(p.As), uint32(p.From.Reg), 0, uint32(p.To.Reg))

	case 73: 
		if p.From.Type != obj.TYPE_REG || p.From.Reg != REG_FPSCR || p.To.Type != obj.TYPE_REG || p.To.Reg < REG_CR0 || REG_CR7 < p.To.Reg {
			c.ctxt.Diag("illegal FPSCR/CR field number\n%v", p)
		}
		o1 = AOP_RRR(OP_MCRFS, ((uint32(p.To.Reg) & 7) << 2), ((0 & 7) << 2), 0)

	case 77: 
		if p.From.Type == obj.TYPE_CONST {
			if p.From.Offset > BIG || p.From.Offset < -BIG {
				c.ctxt.Diag("illegal syscall, sysnum too large: %v", p)
			}
			o1 = AOP_IRR(OP_ADDI, REGZERO, REGZERO, uint32(p.From.Offset))
		} else if p.From.Type == obj.TYPE_REG {
			o1 = LOP_RRR(OP_OR, REGZERO, uint32(p.From.Reg), uint32(p.From.Reg))
		} else {
			c.ctxt.Diag("illegal syscall: %v", p)
			o1 = 0x7fe00008 
		}

		o2 = c.oprrr(p.As)
		o3 = AOP_RRR(c.oprrr(AXOR), REGZERO, REGZERO, REGZERO) 

	case 78: 
		o1 = 0 

	
	case 74:
		v := c.vregoff(&p.To)
		
		inst := c.opstore(p.As)
		if c.opform(inst) == DS_FORM && v&0x3 != 0 {
			log.Fatalf("invalid offset for DS form load/store %v", p)
		}
		o1, o2 = c.symbolAccess(p.To.Sym, v, p.From.Reg, inst)

	

	case 75:
		v := c.vregoff(&p.From)
		
		inst := c.opload(p.As)
		if c.opform(inst) == DS_FORM && v&0x3 != 0 {
			log.Fatalf("invalid offset for DS form load/store %v", p)
		}
		o1, o2 = c.symbolAccess(p.From.Sym, v, p.To.Reg, inst)

	

	case 76:
		v := c.vregoff(&p.From)
		
		inst := c.opload(p.As)
		if c.opform(inst) == DS_FORM && v&0x3 != 0 {
			log.Fatalf("invalid offset for DS form load/store %v", p)
		}
		o1, o2 = c.symbolAccess(p.From.Sym, v, p.To.Reg, inst)
		o3 = LOP_RRR(OP_EXTSB, uint32(p.To.Reg), uint32(p.To.Reg), 0)

		

	case 79:
		if p.From.Offset != 0 {
			c.ctxt.Diag("invalid offset against tls var %v", p)
		}
		o1 = AOP_IRR(OP_ADDI, uint32(p.To.Reg), REGZERO, 0)
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc)
		rel.Siz = 4
		rel.Sym = p.From.Sym
		rel.Type = objabi.R_POWER_TLS_LE

	case 80:
		if p.From.Offset != 0 {
			c.ctxt.Diag("invalid offset against tls var %v", p)
		}
		o1 = AOP_IRR(OP_ADDIS, uint32(p.To.Reg), REG_R2, 0)
		o2 = AOP_IRR(c.opload(AMOVD), uint32(p.To.Reg), uint32(p.To.Reg), 0)
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Type = objabi.R_POWER_TLS_IE

	case 81:
		v := c.vregoff(&p.To)
		if v != 0 {
			c.ctxt.Diag("invalid offset against GOT slot %v", p)
		}

		o1 = AOP_IRR(OP_ADDIS, uint32(p.To.Reg), REG_R2, 0)
		o2 = AOP_IRR(c.opload(AMOVD), uint32(p.To.Reg), uint32(p.To.Reg), 0)
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Type = objabi.R_ADDRPOWER_GOT
	case 82: 
		if p.From.Type == obj.TYPE_REG {
			
			
			
			o1 = AOP_RRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg))
		} else if p.From3Type() == obj.TYPE_CONST {
			
			
			six := int(c.regoff(&p.From))
			st := int(c.regoff(p.GetFrom3()))
			o1 = AOP_IIRR(c.opiirr(p.As), uint32(p.To.Reg), uint32(p.Reg), uint32(st), uint32(six))
		} else if p.From3Type() == obj.TYPE_NONE && p.Reg != 0 {
			
			
			uim := int(c.regoff(&p.From))
			o1 = AOP_VIRR(c.opirr(p.As), uint32(p.To.Reg), uint32(p.Reg), uint32(uim))
		} else {
			
			
			sim := int(c.regoff(&p.From))
			o1 = AOP_IR(c.opirr(p.As), uint32(p.To.Reg), uint32(sim))
		}

	case 83: 
		if p.From.Type == obj.TYPE_REG {
			
			
			o1 = AOP_RRRR(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), uint32(p.GetFrom3().Reg))
		} else if p.From.Type == obj.TYPE_CONST {
			
			
			shb := int(c.regoff(&p.From))
			o1 = AOP_IRRR(c.opirrr(p.As), uint32(p.To.Reg), uint32(p.Reg), uint32(p.GetFrom3().Reg), uint32(shb))
		}

	case 84: 
		bc := c.vregoff(&p.From)

		
		o1 = AOP_ISEL(OP_ISEL, uint32(p.To.Reg), uint32(p.Reg), uint32(p.GetFrom3().Reg), uint32(bc))

	case 85: 
		
		
		o1 = AOP_RR(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg))

	case 86: 
		
		
		o1 = AOP_XX1(c.opstorex(p.As), uint32(p.From.Reg), uint32(p.To.Index), uint32(p.To.Reg))

	case 87: 
		
		
		o1 = AOP_XX1(c.oploadx(p.As), uint32(p.To.Reg), uint32(p.From.Index), uint32(p.From.Reg))

	case 88: 
		
		
		
		xt := int32(p.To.Reg)
		xs := int32(p.From.Reg)
		
		if REG_V0 <= xt && xt <= REG_V31 {
			
			xt = xt + 64
			o1 = AOP_XX1(c.oprrr(p.As), uint32(xt), uint32(p.From.Reg), uint32(p.Reg))
		} else if REG_F0 <= xt && xt <= REG_F31 {
			
			xt = xt + 64
			o1 = AOP_XX1(c.oprrr(p.As), uint32(xt), uint32(p.From.Reg), uint32(p.Reg))
		} else if REG_VS0 <= xt && xt <= REG_VS63 {
			o1 = AOP_XX1(c.oprrr(p.As), uint32(xt), uint32(p.From.Reg), uint32(p.Reg))
		} else if REG_V0 <= xs && xs <= REG_V31 {
			
			xs = xs + 64
			o1 = AOP_XX1(c.oprrr(p.As), uint32(xs), uint32(p.To.Reg), uint32(p.Reg))
		} else if REG_F0 <= xs && xs <= REG_F31 {
			xs = xs + 64
			o1 = AOP_XX1(c.oprrr(p.As), uint32(xs), uint32(p.To.Reg), uint32(p.Reg))
		} else if REG_VS0 <= xs && xs <= REG_VS63 {
			o1 = AOP_XX1(c.oprrr(p.As), uint32(xs), uint32(p.To.Reg), uint32(p.Reg))
		}

	case 89: 
		
		
		uim := int(c.regoff(p.GetFrom3()))
		o1 = AOP_XX2(c.oprrr(p.As), uint32(p.To.Reg), uint32(uim), uint32(p.From.Reg))

	case 90: 
		if p.From3Type() == obj.TYPE_NONE {
			
			
			o1 = AOP_XX3(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg))
		} else if p.From3Type() == obj.TYPE_CONST {
			
			
			dm := int(c.regoff(p.GetFrom3()))
			o1 = AOP_XX3I(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), uint32(dm))
		}

	case 91: 
		
		
		o1 = AOP_XX4(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), uint32(p.GetFrom3().Reg))

	case 92: 
		if p.To.Type == obj.TYPE_CONST {
			
			xf := int32(p.From.Reg)
			if REG_F0 <= xf && xf <= REG_F31 {
				
				bf := int(c.regoff(&p.To)) << 2
				o1 = AOP_RRR(c.opirr(p.As), uint32(bf), uint32(p.From.Reg), uint32(p.Reg))
			} else {
				
				l := int(c.regoff(&p.To))
				o1 = AOP_RRR(c.opirr(p.As), uint32(l), uint32(p.From.Reg), uint32(p.Reg))
			}
		} else if p.From3Type() == obj.TYPE_CONST {
			
			
			l := int(c.regoff(p.GetFrom3()))
			o1 = AOP_RRR(c.opirr(p.As), uint32(l), uint32(p.To.Reg), uint32(p.From.Reg))
		} else if p.To.Type == obj.TYPE_REG {
			cr := int32(p.To.Reg)
			if REG_CR0 <= cr && cr <= REG_CR7 {
				
				
				bf := (int(p.To.Reg) & 7) << 2
				o1 = AOP_RRR(c.opirr(p.As), uint32(bf), uint32(p.From.Reg), uint32(p.Reg))
			} else if p.From.Type == obj.TYPE_CONST {
				
				
				l := int(c.regoff(&p.From))
				o1 = AOP_RRR(c.opirr(p.As), uint32(p.To.Reg), uint32(l), uint32(p.Reg))
			} else {
				switch p.As {
				case ACOPY, APASTECC:
					o1 = AOP_RRR(c.opirr(p.As), uint32(1), uint32(p.From.Reg), uint32(p.To.Reg))
				default:
					
					
					o1 = AOP_RRR(c.oprrr(p.As), uint32(p.From.Reg), uint32(p.To.Reg), uint32(p.Reg))
				}
			}
		}

	case 93: 
		if p.To.Type == obj.TYPE_CONST {
			
			
			bf := int(c.regoff(&p.To)) << 2
			o1 = AOP_RR(c.opirr(p.As), uint32(bf), uint32(p.From.Reg))
		} else if p.Reg == 0 {
			
			
			o1 = AOP_RRR(c.oprrr(p.As), uint32(p.From.Reg), uint32(p.To.Reg), uint32(p.Reg))
		}

	case 94: 
		
		
		cy := int(c.regoff(p.GetFrom3()))
		o1 = AOP_Z23I(c.oprrr(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), uint32(cy))

	case 95: 
		
		v := c.vregoff(&p.From)
		if v != 0 {
			c.ctxt.Diag("invalid offset against TOC slot %v", p)
		}

		inst := c.opload(p.As)
		if c.opform(inst) != DS_FORM {
			c.ctxt.Diag("invalid form for a TOC access in %v", p)
		}

		o1 = AOP_IRR(OP_ADDIS, uint32(p.To.Reg), REG_R2, 0)
		o2 = AOP_IRR(inst, uint32(p.To.Reg), uint32(p.To.Reg), 0)
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Type = objabi.R_ADDRPOWER_TOCREL_DS

	case 96: 
		
		
		dq := int16(c.regoff(&p.From))
		if (dq & 15) != 0 {
			c.ctxt.Diag("invalid offset for DQ form load/store %v", dq)
		}
		o1 = AOP_DQ(c.opload(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(dq))

	case 97: 
		
		
		dq := int16(c.regoff(&p.To))
		if (dq & 15) != 0 {
			c.ctxt.Diag("invalid offset for DQ form load/store %v", dq)
		}
		o1 = AOP_DQ(c.opstore(p.As), uint32(p.From.Reg), uint32(p.To.Reg), uint32(dq))
	case 98: 
		
		o1 = AOP_XX1(c.opload(p.As), uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg))
	case 99: 
		
		o1 = AOP_XX1(c.opstore(p.As), uint32(p.From.Reg), uint32(p.Reg), uint32(p.To.Reg))
	case 100: 
		if p.From.Type == obj.TYPE_CONST {
			
			uim := int(c.regoff(&p.From))
			
			
			o1 = AOP_XX1(c.oprrr(p.As), uint32(p.To.Reg), uint32(0), uint32(uim))
		} else {
			c.ctxt.Diag("invalid ops for %v", p.As)
		}
	case 101:
		o1 = AOP_XX2(c.oprrr(p.As), uint32(p.To.Reg), uint32(0), uint32(p.From.Reg))
	}

	out[0] = o1
	out[1] = o2
	out[2] = o3
	out[3] = o4
	out[4] = o5
}

func (c *ctxt9) vregoff(a *obj.Addr) int64 {
	c.instoffset = 0
	if a != nil {
		c.aclass(a)
	}
	return c.instoffset
}

func (c *ctxt9) regoff(a *obj.Addr) int32 {
	return int32(c.vregoff(a))
}

func (c *ctxt9) oprrr(a obj.As) uint32 {
	switch a {
	case AADD:
		return OPVCC(31, 266, 0, 0)
	case AADDCC:
		return OPVCC(31, 266, 0, 1)
	case AADDV:
		return OPVCC(31, 266, 1, 0)
	case AADDVCC:
		return OPVCC(31, 266, 1, 1)
	case AADDC:
		return OPVCC(31, 10, 0, 0)
	case AADDCCC:
		return OPVCC(31, 10, 0, 1)
	case AADDCV:
		return OPVCC(31, 10, 1, 0)
	case AADDCVCC:
		return OPVCC(31, 10, 1, 1)
	case AADDE:
		return OPVCC(31, 138, 0, 0)
	case AADDECC:
		return OPVCC(31, 138, 0, 1)
	case AADDEV:
		return OPVCC(31, 138, 1, 0)
	case AADDEVCC:
		return OPVCC(31, 138, 1, 1)
	case AADDME:
		return OPVCC(31, 234, 0, 0)
	case AADDMECC:
		return OPVCC(31, 234, 0, 1)
	case AADDMEV:
		return OPVCC(31, 234, 1, 0)
	case AADDMEVCC:
		return OPVCC(31, 234, 1, 1)
	case AADDZE:
		return OPVCC(31, 202, 0, 0)
	case AADDZECC:
		return OPVCC(31, 202, 0, 1)
	case AADDZEV:
		return OPVCC(31, 202, 1, 0)
	case AADDZEVCC:
		return OPVCC(31, 202, 1, 1)
	case AADDEX:
		return OPVCC(31, 170, 0, 0) 

	case AAND:
		return OPVCC(31, 28, 0, 0)
	case AANDCC:
		return OPVCC(31, 28, 0, 1)
	case AANDN:
		return OPVCC(31, 60, 0, 0)
	case AANDNCC:
		return OPVCC(31, 60, 0, 1)

	case ACMP:
		return OPVCC(31, 0, 0, 0) | 1<<21 
	case ACMPU:
		return OPVCC(31, 32, 0, 0) | 1<<21
	case ACMPW:
		return OPVCC(31, 0, 0, 0) 
	case ACMPWU:
		return OPVCC(31, 32, 0, 0)
	case ACMPB:
		return OPVCC(31, 508, 0, 0) 
	case ACMPEQB:
		return OPVCC(31, 224, 0, 0) 

	case ACNTLZW:
		return OPVCC(31, 26, 0, 0)
	case ACNTLZWCC:
		return OPVCC(31, 26, 0, 1)
	case ACNTLZD:
		return OPVCC(31, 58, 0, 0)
	case ACNTLZDCC:
		return OPVCC(31, 58, 0, 1)

	case ACRAND:
		return OPVCC(19, 257, 0, 0)
	case ACRANDN:
		return OPVCC(19, 129, 0, 0)
	case ACREQV:
		return OPVCC(19, 289, 0, 0)
	case ACRNAND:
		return OPVCC(19, 225, 0, 0)
	case ACRNOR:
		return OPVCC(19, 33, 0, 0)
	case ACROR:
		return OPVCC(19, 449, 0, 0)
	case ACRORN:
		return OPVCC(19, 417, 0, 0)
	case ACRXOR:
		return OPVCC(19, 193, 0, 0)

	case ADCBF:
		return OPVCC(31, 86, 0, 0)
	case ADCBI:
		return OPVCC(31, 470, 0, 0)
	case ADCBST:
		return OPVCC(31, 54, 0, 0)
	case ADCBT:
		return OPVCC(31, 278, 0, 0)
	case ADCBTST:
		return OPVCC(31, 246, 0, 0)
	case ADCBZ:
		return OPVCC(31, 1014, 0, 0)

	case AMODUD:
		return OPVCC(31, 265, 0, 0) 
	case AMODUW:
		return OPVCC(31, 267, 0, 0) 
	case AMODSD:
		return OPVCC(31, 777, 0, 0) 
	case AMODSW:
		return OPVCC(31, 779, 0, 0) 

	case ADIVW, AREM:
		return OPVCC(31, 491, 0, 0)

	case ADIVWCC:
		return OPVCC(31, 491, 0, 1)

	case ADIVWV:
		return OPVCC(31, 491, 1, 0)

	case ADIVWVCC:
		return OPVCC(31, 491, 1, 1)

	case ADIVWU, AREMU:
		return OPVCC(31, 459, 0, 0)

	case ADIVWUCC:
		return OPVCC(31, 459, 0, 1)

	case ADIVWUV:
		return OPVCC(31, 459, 1, 0)

	case ADIVWUVCC:
		return OPVCC(31, 459, 1, 1)

	case ADIVD, AREMD:
		return OPVCC(31, 489, 0, 0)

	case ADIVDCC:
		return OPVCC(31, 489, 0, 1)

	case ADIVDE:
		return OPVCC(31, 425, 0, 0)

	case ADIVDECC:
		return OPVCC(31, 425, 0, 1)

	case ADIVDEU:
		return OPVCC(31, 393, 0, 0)

	case ADIVDEUCC:
		return OPVCC(31, 393, 0, 1)

	case ADIVDV:
		return OPVCC(31, 489, 1, 0)

	case ADIVDVCC:
		return OPVCC(31, 489, 1, 1)

	case ADIVDU, AREMDU:
		return OPVCC(31, 457, 0, 0)

	case ADIVDUCC:
		return OPVCC(31, 457, 0, 1)

	case ADIVDUV:
		return OPVCC(31, 457, 1, 0)

	case ADIVDUVCC:
		return OPVCC(31, 457, 1, 1)

	case AEIEIO:
		return OPVCC(31, 854, 0, 0)

	case AEQV:
		return OPVCC(31, 284, 0, 0)
	case AEQVCC:
		return OPVCC(31, 284, 0, 1)

	case AEXTSB:
		return OPVCC(31, 954, 0, 0)
	case AEXTSBCC:
		return OPVCC(31, 954, 0, 1)
	case AEXTSH:
		return OPVCC(31, 922, 0, 0)
	case AEXTSHCC:
		return OPVCC(31, 922, 0, 1)
	case AEXTSW:
		return OPVCC(31, 986, 0, 0)
	case AEXTSWCC:
		return OPVCC(31, 986, 0, 1)

	case AFABS:
		return OPVCC(63, 264, 0, 0)
	case AFABSCC:
		return OPVCC(63, 264, 0, 1)
	case AFADD:
		return OPVCC(63, 21, 0, 0)
	case AFADDCC:
		return OPVCC(63, 21, 0, 1)
	case AFADDS:
		return OPVCC(59, 21, 0, 0)
	case AFADDSCC:
		return OPVCC(59, 21, 0, 1)
	case AFCMPO:
		return OPVCC(63, 32, 0, 0)
	case AFCMPU:
		return OPVCC(63, 0, 0, 0)
	case AFCFID:
		return OPVCC(63, 846, 0, 0)
	case AFCFIDCC:
		return OPVCC(63, 846, 0, 1)
	case AFCFIDU:
		return OPVCC(63, 974, 0, 0)
	case AFCFIDUCC:
		return OPVCC(63, 974, 0, 1)
	case AFCFIDS:
		return OPVCC(59, 846, 0, 0)
	case AFCFIDSCC:
		return OPVCC(59, 846, 0, 1)
	case AFCTIW:
		return OPVCC(63, 14, 0, 0)
	case AFCTIWCC:
		return OPVCC(63, 14, 0, 1)
	case AFCTIWZ:
		return OPVCC(63, 15, 0, 0)
	case AFCTIWZCC:
		return OPVCC(63, 15, 0, 1)
	case AFCTID:
		return OPVCC(63, 814, 0, 0)
	case AFCTIDCC:
		return OPVCC(63, 814, 0, 1)
	case AFCTIDZ:
		return OPVCC(63, 815, 0, 0)
	case AFCTIDZCC:
		return OPVCC(63, 815, 0, 1)
	case AFDIV:
		return OPVCC(63, 18, 0, 0)
	case AFDIVCC:
		return OPVCC(63, 18, 0, 1)
	case AFDIVS:
		return OPVCC(59, 18, 0, 0)
	case AFDIVSCC:
		return OPVCC(59, 18, 0, 1)
	case AFMADD:
		return OPVCC(63, 29, 0, 0)
	case AFMADDCC:
		return OPVCC(63, 29, 0, 1)
	case AFMADDS:
		return OPVCC(59, 29, 0, 0)
	case AFMADDSCC:
		return OPVCC(59, 29, 0, 1)

	case AFMOVS, AFMOVD:
		return OPVCC(63, 72, 0, 0) 
	case AFMOVDCC:
		return OPVCC(63, 72, 0, 1)
	case AFMSUB:
		return OPVCC(63, 28, 0, 0)
	case AFMSUBCC:
		return OPVCC(63, 28, 0, 1)
	case AFMSUBS:
		return OPVCC(59, 28, 0, 0)
	case AFMSUBSCC:
		return OPVCC(59, 28, 0, 1)
	case AFMUL:
		return OPVCC(63, 25, 0, 0)
	case AFMULCC:
		return OPVCC(63, 25, 0, 1)
	case AFMULS:
		return OPVCC(59, 25, 0, 0)
	case AFMULSCC:
		return OPVCC(59, 25, 0, 1)
	case AFNABS:
		return OPVCC(63, 136, 0, 0)
	case AFNABSCC:
		return OPVCC(63, 136, 0, 1)
	case AFNEG:
		return OPVCC(63, 40, 0, 0)
	case AFNEGCC:
		return OPVCC(63, 40, 0, 1)
	case AFNMADD:
		return OPVCC(63, 31, 0, 0)
	case AFNMADDCC:
		return OPVCC(63, 31, 0, 1)
	case AFNMADDS:
		return OPVCC(59, 31, 0, 0)
	case AFNMADDSCC:
		return OPVCC(59, 31, 0, 1)
	case AFNMSUB:
		return OPVCC(63, 30, 0, 0)
	case AFNMSUBCC:
		return OPVCC(63, 30, 0, 1)
	case AFNMSUBS:
		return OPVCC(59, 30, 0, 0)
	case AFNMSUBSCC:
		return OPVCC(59, 30, 0, 1)
	case AFCPSGN:
		return OPVCC(63, 8, 0, 0)
	case AFCPSGNCC:
		return OPVCC(63, 8, 0, 1)
	case AFRES:
		return OPVCC(59, 24, 0, 0)
	case AFRESCC:
		return OPVCC(59, 24, 0, 1)
	case AFRIM:
		return OPVCC(63, 488, 0, 0)
	case AFRIMCC:
		return OPVCC(63, 488, 0, 1)
	case AFRIP:
		return OPVCC(63, 456, 0, 0)
	case AFRIPCC:
		return OPVCC(63, 456, 0, 1)
	case AFRIZ:
		return OPVCC(63, 424, 0, 0)
	case AFRIZCC:
		return OPVCC(63, 424, 0, 1)
	case AFRIN:
		return OPVCC(63, 392, 0, 0)
	case AFRINCC:
		return OPVCC(63, 392, 0, 1)
	case AFRSP:
		return OPVCC(63, 12, 0, 0)
	case AFRSPCC:
		return OPVCC(63, 12, 0, 1)
	case AFRSQRTE:
		return OPVCC(63, 26, 0, 0)
	case AFRSQRTECC:
		return OPVCC(63, 26, 0, 1)
	case AFSEL:
		return OPVCC(63, 23, 0, 0)
	case AFSELCC:
		return OPVCC(63, 23, 0, 1)
	case AFSQRT:
		return OPVCC(63, 22, 0, 0)
	case AFSQRTCC:
		return OPVCC(63, 22, 0, 1)
	case AFSQRTS:
		return OPVCC(59, 22, 0, 0)
	case AFSQRTSCC:
		return OPVCC(59, 22, 0, 1)
	case AFSUB:
		return OPVCC(63, 20, 0, 0)
	case AFSUBCC:
		return OPVCC(63, 20, 0, 1)
	case AFSUBS:
		return OPVCC(59, 20, 0, 0)
	case AFSUBSCC:
		return OPVCC(59, 20, 0, 1)

	case AICBI:
		return OPVCC(31, 982, 0, 0)
	case AISYNC:
		return OPVCC(19, 150, 0, 0)

	case AMTFSB0:
		return OPVCC(63, 70, 0, 0)
	case AMTFSB0CC:
		return OPVCC(63, 70, 0, 1)
	case AMTFSB1:
		return OPVCC(63, 38, 0, 0)
	case AMTFSB1CC:
		return OPVCC(63, 38, 0, 1)

	case AMULHW:
		return OPVCC(31, 75, 0, 0)
	case AMULHWCC:
		return OPVCC(31, 75, 0, 1)
	case AMULHWU:
		return OPVCC(31, 11, 0, 0)
	case AMULHWUCC:
		return OPVCC(31, 11, 0, 1)
	case AMULLW:
		return OPVCC(31, 235, 0, 0)
	case AMULLWCC:
		return OPVCC(31, 235, 0, 1)
	case AMULLWV:
		return OPVCC(31, 235, 1, 0)
	case AMULLWVCC:
		return OPVCC(31, 235, 1, 1)

	case AMULHD:
		return OPVCC(31, 73, 0, 0)
	case AMULHDCC:
		return OPVCC(31, 73, 0, 1)
	case AMULHDU:
		return OPVCC(31, 9, 0, 0)
	case AMULHDUCC:
		return OPVCC(31, 9, 0, 1)
	case AMULLD:
		return OPVCC(31, 233, 0, 0)
	case AMULLDCC:
		return OPVCC(31, 233, 0, 1)
	case AMULLDV:
		return OPVCC(31, 233, 1, 0)
	case AMULLDVCC:
		return OPVCC(31, 233, 1, 1)

	case ANAND:
		return OPVCC(31, 476, 0, 0)
	case ANANDCC:
		return OPVCC(31, 476, 0, 1)
	case ANEG:
		return OPVCC(31, 104, 0, 0)
	case ANEGCC:
		return OPVCC(31, 104, 0, 1)
	case ANEGV:
		return OPVCC(31, 104, 1, 0)
	case ANEGVCC:
		return OPVCC(31, 104, 1, 1)
	case ANOR:
		return OPVCC(31, 124, 0, 0)
	case ANORCC:
		return OPVCC(31, 124, 0, 1)
	case AOR:
		return OPVCC(31, 444, 0, 0)
	case AORCC:
		return OPVCC(31, 444, 0, 1)
	case AORN:
		return OPVCC(31, 412, 0, 0)
	case AORNCC:
		return OPVCC(31, 412, 0, 1)

	case APOPCNTD:
		return OPVCC(31, 506, 0, 0) 
	case APOPCNTW:
		return OPVCC(31, 378, 0, 0) 
	case APOPCNTB:
		return OPVCC(31, 122, 0, 0) 
	case ACNTTZW:
		return OPVCC(31, 538, 0, 0) 
	case ACNTTZWCC:
		return OPVCC(31, 538, 0, 1) 
	case ACNTTZD:
		return OPVCC(31, 570, 0, 0) 
	case ACNTTZDCC:
		return OPVCC(31, 570, 0, 1) 

	case ARFI:
		return OPVCC(19, 50, 0, 0)
	case ARFCI:
		return OPVCC(19, 51, 0, 0)
	case ARFID:
		return OPVCC(19, 18, 0, 0)
	case AHRFID:
		return OPVCC(19, 274, 0, 0)

	case ARLWMI:
		return OPVCC(20, 0, 0, 0)
	case ARLWMICC:
		return OPVCC(20, 0, 0, 1)
	case ARLWNM:
		return OPVCC(23, 0, 0, 0)
	case ARLWNMCC:
		return OPVCC(23, 0, 0, 1)

	case ARLDCL:
		return OPVCC(30, 8, 0, 0)
	case ARLDCLCC:
		return OPVCC(30, 0, 0, 1)

	case ARLDCR:
		return OPVCC(30, 9, 0, 0)
	case ARLDCRCC:
		return OPVCC(30, 9, 0, 1)

	case ARLDICL:
		return OPVCC(30, 0, 0, 0)
	case ARLDICLCC:
		return OPVCC(30, 0, 0, 1)
	case ARLDICR:
		return OPVCC(30, 0, 0, 0) | 2<<1 
	case ARLDICRCC:
		return OPVCC(30, 0, 0, 1) | 2<<1 

	case ARLDIC:
		return OPVCC(30, 0, 0, 0) | 4<<1 
	case ARLDICCC:
		return OPVCC(30, 0, 0, 1) | 4<<1 

	case ASYSCALL:
		return OPVCC(17, 1, 0, 0)

	case ASLW:
		return OPVCC(31, 24, 0, 0)
	case ASLWCC:
		return OPVCC(31, 24, 0, 1)
	case ASLD:
		return OPVCC(31, 27, 0, 0)
	case ASLDCC:
		return OPVCC(31, 27, 0, 1)

	case ASRAW:
		return OPVCC(31, 792, 0, 0)
	case ASRAWCC:
		return OPVCC(31, 792, 0, 1)
	case ASRAD:
		return OPVCC(31, 794, 0, 0)
	case ASRADCC:
		return OPVCC(31, 794, 0, 1)

	case ASRW:
		return OPVCC(31, 536, 0, 0)
	case ASRWCC:
		return OPVCC(31, 536, 0, 1)
	case ASRD:
		return OPVCC(31, 539, 0, 0)
	case ASRDCC:
		return OPVCC(31, 539, 0, 1)

	case ASUB:
		return OPVCC(31, 40, 0, 0)
	case ASUBCC:
		return OPVCC(31, 40, 0, 1)
	case ASUBV:
		return OPVCC(31, 40, 1, 0)
	case ASUBVCC:
		return OPVCC(31, 40, 1, 1)
	case ASUBC:
		return OPVCC(31, 8, 0, 0)
	case ASUBCCC:
		return OPVCC(31, 8, 0, 1)
	case ASUBCV:
		return OPVCC(31, 8, 1, 0)
	case ASUBCVCC:
		return OPVCC(31, 8, 1, 1)
	case ASUBE:
		return OPVCC(31, 136, 0, 0)
	case ASUBECC:
		return OPVCC(31, 136, 0, 1)
	case ASUBEV:
		return OPVCC(31, 136, 1, 0)
	case ASUBEVCC:
		return OPVCC(31, 136, 1, 1)
	case ASUBME:
		return OPVCC(31, 232, 0, 0)
	case ASUBMECC:
		return OPVCC(31, 232, 0, 1)
	case ASUBMEV:
		return OPVCC(31, 232, 1, 0)
	case ASUBMEVCC:
		return OPVCC(31, 232, 1, 1)
	case ASUBZE:
		return OPVCC(31, 200, 0, 0)
	case ASUBZECC:
		return OPVCC(31, 200, 0, 1)
	case ASUBZEV:
		return OPVCC(31, 200, 1, 0)
	case ASUBZEVCC:
		return OPVCC(31, 200, 1, 1)

	case ASYNC:
		return OPVCC(31, 598, 0, 0)
	case ALWSYNC:
		return OPVCC(31, 598, 0, 0) | 1<<21

	case APTESYNC:
		return OPVCC(31, 598, 0, 0) | 2<<21

	case ATLBIE:
		return OPVCC(31, 306, 0, 0)
	case ATLBIEL:
		return OPVCC(31, 274, 0, 0)
	case ATLBSYNC:
		return OPVCC(31, 566, 0, 0)
	case ASLBIA:
		return OPVCC(31, 498, 0, 0)
	case ASLBIE:
		return OPVCC(31, 434, 0, 0)
	case ASLBMFEE:
		return OPVCC(31, 915, 0, 0)
	case ASLBMFEV:
		return OPVCC(31, 851, 0, 0)
	case ASLBMTE:
		return OPVCC(31, 402, 0, 0)

	case ATW:
		return OPVCC(31, 4, 0, 0)
	case ATD:
		return OPVCC(31, 68, 0, 0)

	
	
	
	case AVAND:
		return OPVX(4, 1028, 0, 0) 
	case AVANDC:
		return OPVX(4, 1092, 0, 0) 
	case AVNAND:
		return OPVX(4, 1412, 0, 0) 

	case AVOR:
		return OPVX(4, 1156, 0, 0) 
	case AVORC:
		return OPVX(4, 1348, 0, 0) 
	case AVNOR:
		return OPVX(4, 1284, 0, 0) 
	case AVXOR:
		return OPVX(4, 1220, 0, 0) 
	case AVEQV:
		return OPVX(4, 1668, 0, 0) 

	case AVADDUBM:
		return OPVX(4, 0, 0, 0) 
	case AVADDUHM:
		return OPVX(4, 64, 0, 0) 
	case AVADDUWM:
		return OPVX(4, 128, 0, 0) 
	case AVADDUDM:
		return OPVX(4, 192, 0, 0) 
	case AVADDUQM:
		return OPVX(4, 256, 0, 0) 

	case AVADDCUQ:
		return OPVX(4, 320, 0, 0) 
	case AVADDCUW:
		return OPVX(4, 384, 0, 0) 

	case AVADDUBS:
		return OPVX(4, 512, 0, 0) 
	case AVADDUHS:
		return OPVX(4, 576, 0, 0) 
	case AVADDUWS:
		return OPVX(4, 640, 0, 0) 

	case AVADDSBS:
		return OPVX(4, 768, 0, 0) 
	case AVADDSHS:
		return OPVX(4, 832, 0, 0) 
	case AVADDSWS:
		return OPVX(4, 896, 0, 0) 

	case AVADDEUQM:
		return OPVX(4, 60, 0, 0) 
	case AVADDECUQ:
		return OPVX(4, 61, 0, 0) 

	case AVMULESB:
		return OPVX(4, 776, 0, 0) 
	case AVMULOSB:
		return OPVX(4, 264, 0, 0) 
	case AVMULEUB:
		return OPVX(4, 520, 0, 0) 
	case AVMULOUB:
		return OPVX(4, 8, 0, 0) 
	case AVMULESH:
		return OPVX(4, 840, 0, 0) 
	case AVMULOSH:
		return OPVX(4, 328, 0, 0) 
	case AVMULEUH:
		return OPVX(4, 584, 0, 0) 
	case AVMULOUH:
		return OPVX(4, 72, 0, 0) 
	case AVMULESW:
		return OPVX(4, 904, 0, 0) 
	case AVMULOSW:
		return OPVX(4, 392, 0, 0) 
	case AVMULEUW:
		return OPVX(4, 648, 0, 0) 
	case AVMULOUW:
		return OPVX(4, 136, 0, 0) 
	case AVMULUWM:
		return OPVX(4, 137, 0, 0) 

	case AVPMSUMB:
		return OPVX(4, 1032, 0, 0) 
	case AVPMSUMH:
		return OPVX(4, 1096, 0, 0) 
	case AVPMSUMW:
		return OPVX(4, 1160, 0, 0) 
	case AVPMSUMD:
		return OPVX(4, 1224, 0, 0) 

	case AVMSUMUDM:
		return OPVX(4, 35, 0, 0) 

	case AVSUBUBM:
		return OPVX(4, 1024, 0, 0) 
	case AVSUBUHM:
		return OPVX(4, 1088, 0, 0) 
	case AVSUBUWM:
		return OPVX(4, 1152, 0, 0) 
	case AVSUBUDM:
		return OPVX(4, 1216, 0, 0) 
	case AVSUBUQM:
		return OPVX(4, 1280, 0, 0) 

	case AVSUBCUQ:
		return OPVX(4, 1344, 0, 0) 
	case AVSUBCUW:
		return OPVX(4, 1408, 0, 0) 

	case AVSUBUBS:
		return OPVX(4, 1536, 0, 0) 
	case AVSUBUHS:
		return OPVX(4, 1600, 0, 0) 
	case AVSUBUWS:
		return OPVX(4, 1664, 0, 0) 

	case AVSUBSBS:
		return OPVX(4, 1792, 0, 0) 
	case AVSUBSHS:
		return OPVX(4, 1856, 0, 0) 
	case AVSUBSWS:
		return OPVX(4, 1920, 0, 0) 

	case AVSUBEUQM:
		return OPVX(4, 62, 0, 0) 
	case AVSUBECUQ:
		return OPVX(4, 63, 0, 0) 

	case AVRLB:
		return OPVX(4, 4, 0, 0) 
	case AVRLH:
		return OPVX(4, 68, 0, 0) 
	case AVRLW:
		return OPVX(4, 132, 0, 0) 
	case AVRLD:
		return OPVX(4, 196, 0, 0) 

	case AVMRGOW:
		return OPVX(4, 1676, 0, 0) 
	case AVMRGEW:
		return OPVX(4, 1932, 0, 0) 

	case AVSLB:
		return OPVX(4, 260, 0, 0) 
	case AVSLH:
		return OPVX(4, 324, 0, 0) 
	case AVSLW:
		return OPVX(4, 388, 0, 0) 
	case AVSL:
		return OPVX(4, 452, 0, 0) 
	case AVSLO:
		return OPVX(4, 1036, 0, 0) 
	case AVSRB:
		return OPVX(4, 516, 0, 0) 
	case AVSRH:
		return OPVX(4, 580, 0, 0) 
	case AVSRW:
		return OPVX(4, 644, 0, 0) 
	case AVSR:
		return OPVX(4, 708, 0, 0) 
	case AVSRO:
		return OPVX(4, 1100, 0, 0) 
	case AVSLD:
		return OPVX(4, 1476, 0, 0) 
	case AVSRD:
		return OPVX(4, 1732, 0, 0) 

	case AVSRAB:
		return OPVX(4, 772, 0, 0) 
	case AVSRAH:
		return OPVX(4, 836, 0, 0) 
	case AVSRAW:
		return OPVX(4, 900, 0, 0) 
	case AVSRAD:
		return OPVX(4, 964, 0, 0) 

	case AVBPERMQ:
		return OPVC(4, 1356, 0, 0) 
	case AVBPERMD:
		return OPVC(4, 1484, 0, 0) 

	case AVCLZB:
		return OPVX(4, 1794, 0, 0) 
	case AVCLZH:
		return OPVX(4, 1858, 0, 0) 
	case AVCLZW:
		return OPVX(4, 1922, 0, 0) 
	case AVCLZD:
		return OPVX(4, 1986, 0, 0) 

	case AVPOPCNTB:
		return OPVX(4, 1795, 0, 0) 
	case AVPOPCNTH:
		return OPVX(4, 1859, 0, 0) 
	case AVPOPCNTW:
		return OPVX(4, 1923, 0, 0) 
	case AVPOPCNTD:
		return OPVX(4, 1987, 0, 0) 

	case AVCMPEQUB:
		return OPVC(4, 6, 0, 0) 
	case AVCMPEQUBCC:
		return OPVC(4, 6, 0, 1) 
	case AVCMPEQUH:
		return OPVC(4, 70, 0, 0) 
	case AVCMPEQUHCC:
		return OPVC(4, 70, 0, 1) 
	case AVCMPEQUW:
		return OPVC(4, 134, 0, 0) 
	case AVCMPEQUWCC:
		return OPVC(4, 134, 0, 1) 
	case AVCMPEQUD:
		return OPVC(4, 199, 0, 0) 
	case AVCMPEQUDCC:
		return OPVC(4, 199, 0, 1) 

	case AVCMPGTUB:
		return OPVC(4, 518, 0, 0) 
	case AVCMPGTUBCC:
		return OPVC(4, 518, 0, 1) 
	case AVCMPGTUH:
		return OPVC(4, 582, 0, 0) 
	case AVCMPGTUHCC:
		return OPVC(4, 582, 0, 1) 
	case AVCMPGTUW:
		return OPVC(4, 646, 0, 0) 
	case AVCMPGTUWCC:
		return OPVC(4, 646, 0, 1) 
	case AVCMPGTUD:
		return OPVC(4, 711, 0, 0) 
	case AVCMPGTUDCC:
		return OPVC(4, 711, 0, 1) 
	case AVCMPGTSB:
		return OPVC(4, 774, 0, 0) 
	case AVCMPGTSBCC:
		return OPVC(4, 774, 0, 1) 
	case AVCMPGTSH:
		return OPVC(4, 838, 0, 0) 
	case AVCMPGTSHCC:
		return OPVC(4, 838, 0, 1) 
	case AVCMPGTSW:
		return OPVC(4, 902, 0, 0) 
	case AVCMPGTSWCC:
		return OPVC(4, 902, 0, 1) 
	case AVCMPGTSD:
		return OPVC(4, 967, 0, 0) 
	case AVCMPGTSDCC:
		return OPVC(4, 967, 0, 1) 

	case AVCMPNEZB:
		return OPVC(4, 263, 0, 0) 
	case AVCMPNEZBCC:
		return OPVC(4, 263, 0, 1) 
	case AVCMPNEB:
		return OPVC(4, 7, 0, 0) 
	case AVCMPNEBCC:
		return OPVC(4, 7, 0, 1) 
	case AVCMPNEH:
		return OPVC(4, 71, 0, 0) 
	case AVCMPNEHCC:
		return OPVC(4, 71, 0, 1) 
	case AVCMPNEW:
		return OPVC(4, 135, 0, 0) 
	case AVCMPNEWCC:
		return OPVC(4, 135, 0, 1) 

	case AVPERM:
		return OPVX(4, 43, 0, 0) 
	case AVPERMXOR:
		return OPVX(4, 45, 0, 0) 
	case AVPERMR:
		return OPVX(4, 59, 0, 0) 

	case AVSEL:
		return OPVX(4, 42, 0, 0) 

	case AVCIPHER:
		return OPVX(4, 1288, 0, 0) 
	case AVCIPHERLAST:
		return OPVX(4, 1289, 0, 0) 
	case AVNCIPHER:
		return OPVX(4, 1352, 0, 0) 
	case AVNCIPHERLAST:
		return OPVX(4, 1353, 0, 0) 
	case AVSBOX:
		return OPVX(4, 1480, 0, 0) 
	

	
	
	case AMFVSRD, AMFVRD, AMFFPRD:
		return OPVXX1(31, 51, 0) 
	case AMFVSRWZ:
		return OPVXX1(31, 115, 0) 
	case AMFVSRLD:
		return OPVXX1(31, 307, 0) 

	case AMTVSRD, AMTFPRD, AMTVRD:
		return OPVXX1(31, 179, 0) 
	case AMTVSRWA:
		return OPVXX1(31, 211, 0) 
	case AMTVSRWZ:
		return OPVXX1(31, 243, 0) 
	case AMTVSRDD:
		return OPVXX1(31, 435, 0) 
	case AMTVSRWS:
		return OPVXX1(31, 403, 0) 

	case AXXLAND:
		return OPVXX3(60, 130, 0) 
	case AXXLANDC:
		return OPVXX3(60, 138, 0) 
	case AXXLEQV:
		return OPVXX3(60, 186, 0) 
	case AXXLNAND:
		return OPVXX3(60, 178, 0) 

	case AXXLORC:
		return OPVXX3(60, 170, 0) 
	case AXXLNOR:
		return OPVXX3(60, 162, 0) 
	case AXXLOR, AXXLORQ:
		return OPVXX3(60, 146, 0) 
	case AXXLXOR:
		return OPVXX3(60, 154, 0) 

	case AXXSEL:
		return OPVXX4(60, 3, 0) 

	case AXXMRGHW:
		return OPVXX3(60, 18, 0) 
	case AXXMRGLW:
		return OPVXX3(60, 50, 0) 

	case AXXSPLTW:
		return OPVXX2(60, 164, 0) 

	case AXXSPLTIB:
		return OPVCC(60, 360, 0, 0) 

	case AXXPERM:
		return OPVXX3(60, 26, 0) 
	case AXXPERMDI:
		return OPVXX3(60, 10, 0) 

	case AXXSLDWI:
		return OPVXX3(60, 2, 0) 

	case AXXBRQ:
		return OPVXX2VA(60, 475, 31) 
	case AXXBRD:
		return OPVXX2VA(60, 475, 23) 
	case AXXBRW:
		return OPVXX2VA(60, 475, 15) 
	case AXXBRH:
		return OPVXX2VA(60, 475, 7) 

	case AXSCVDPSP:
		return OPVXX2(60, 265, 0) 
	case AXSCVSPDP:
		return OPVXX2(60, 329, 0) 
	case AXSCVDPSPN:
		return OPVXX2(60, 267, 0) 
	case AXSCVSPDPN:
		return OPVXX2(60, 331, 0) 

	case AXVCVDPSP:
		return OPVXX2(60, 393, 0) 
	case AXVCVSPDP:
		return OPVXX2(60, 457, 0) 

	case AXSCVDPSXDS:
		return OPVXX2(60, 344, 0) 
	case AXSCVDPSXWS:
		return OPVXX2(60, 88, 0) 
	case AXSCVDPUXDS:
		return OPVXX2(60, 328, 0) 
	case AXSCVDPUXWS:
		return OPVXX2(60, 72, 0) 

	case AXSCVSXDDP:
		return OPVXX2(60, 376, 0) 
	case AXSCVUXDDP:
		return OPVXX2(60, 360, 0) 
	case AXSCVSXDSP:
		return OPVXX2(60, 312, 0) 
	case AXSCVUXDSP:
		return OPVXX2(60, 296, 0) 

	case AXVCVDPSXDS:
		return OPVXX2(60, 472, 0) 
	case AXVCVDPSXWS:
		return OPVXX2(60, 216, 0) 
	case AXVCVDPUXDS:
		return OPVXX2(60, 456, 0) 
	case AXVCVDPUXWS:
		return OPVXX2(60, 200, 0) 
	case AXVCVSPSXDS:
		return OPVXX2(60, 408, 0) 
	case AXVCVSPSXWS:
		return OPVXX2(60, 152, 0) 
	case AXVCVSPUXDS:
		return OPVXX2(60, 392, 0) 
	case AXVCVSPUXWS:
		return OPVXX2(60, 136, 0) 

	case AXVCVSXDDP:
		return OPVXX2(60, 504, 0) 
	case AXVCVSXWDP:
		return OPVXX2(60, 248, 0) 
	case AXVCVUXDDP:
		return OPVXX2(60, 488, 0) 
	case AXVCVUXWDP:
		return OPVXX2(60, 232, 0) 
	case AXVCVSXDSP:
		return OPVXX2(60, 440, 0) 
	case AXVCVSXWSP:
		return OPVXX2(60, 184, 0) 
	case AXVCVUXDSP:
		return OPVXX2(60, 424, 0) 
	case AXVCVUXWSP:
		return OPVXX2(60, 168, 0) 
	

	case AMADDHD:
		return OPVX(4, 48, 0, 0) 
	case AMADDHDU:
		return OPVX(4, 49, 0, 0) 
	case AMADDLD:
		return OPVX(4, 51, 0, 0) 

	case AXOR:
		return OPVCC(31, 316, 0, 0)
	case AXORCC:
		return OPVCC(31, 316, 0, 1)
	}

	c.ctxt.Diag("bad r/r, r/r/r or r/r/r/r opcode %v", a)
	return 0
}

func (c *ctxt9) opirrr(a obj.As) uint32 {
	switch a {
	
	
	
	case AVSLDOI:
		return OPVX(4, 44, 0, 0) 
	}

	c.ctxt.Diag("bad i/r/r/r opcode %v", a)
	return 0
}

func (c *ctxt9) opiirr(a obj.As) uint32 {
	switch a {
	
	
	case AVSHASIGMAW:
		return OPVX(4, 1666, 0, 0) 
	case AVSHASIGMAD:
		return OPVX(4, 1730, 0, 0) 
	}

	c.ctxt.Diag("bad i/i/r/r opcode %v", a)
	return 0
}

func (c *ctxt9) opirr(a obj.As) uint32 {
	switch a {
	case AADD:
		return OPVCC(14, 0, 0, 0)
	case AADDC:
		return OPVCC(12, 0, 0, 0)
	case AADDCCC:
		return OPVCC(13, 0, 0, 0)
	case AADDIS:
		return OPVCC(15, 0, 0, 0) 

	case AANDCC:
		return OPVCC(28, 0, 0, 0)
	case AANDISCC:
		return OPVCC(29, 0, 0, 0) 

	case ABR:
		return OPVCC(18, 0, 0, 0)
	case ABL:
		return OPVCC(18, 0, 0, 0) | 1
	case obj.ADUFFZERO:
		return OPVCC(18, 0, 0, 0) | 1
	case obj.ADUFFCOPY:
		return OPVCC(18, 0, 0, 0) | 1
	case ABC:
		return OPVCC(16, 0, 0, 0)
	case ABCL:
		return OPVCC(16, 0, 0, 0) | 1

	case ABEQ:
		return AOP_RRR(16<<26, 12, 2, 0)
	case ABGE:
		return AOP_RRR(16<<26, 4, 0, 0)
	case ABGT:
		return AOP_RRR(16<<26, 12, 1, 0)
	case ABLE:
		return AOP_RRR(16<<26, 4, 1, 0)
	case ABLT:
		return AOP_RRR(16<<26, 12, 0, 0)
	case ABNE:
		return AOP_RRR(16<<26, 4, 2, 0)
	case ABVC:
		return AOP_RRR(16<<26, 4, 3, 0) 
	case ABVS:
		return AOP_RRR(16<<26, 12, 3, 0) 

	case ACMP:
		return OPVCC(11, 0, 0, 0) | 1<<21 
	case ACMPU:
		return OPVCC(10, 0, 0, 0) | 1<<21
	case ACMPW:
		return OPVCC(11, 0, 0, 0) 
	case ACMPWU:
		return OPVCC(10, 0, 0, 0)
	case ACMPEQB:
		return OPVCC(31, 224, 0, 0) 

	case ALSW:
		return OPVCC(31, 597, 0, 0)

	case ACOPY:
		return OPVCC(31, 774, 0, 0) 
	case APASTECC:
		return OPVCC(31, 902, 0, 1) 
	case ADARN:
		return OPVCC(31, 755, 0, 0) 

	case AMULLW:
		return OPVCC(7, 0, 0, 0)

	case AOR:
		return OPVCC(24, 0, 0, 0)
	case AORIS:
		return OPVCC(25, 0, 0, 0) 

	case ARLWMI:
		return OPVCC(20, 0, 0, 0) 
	case ARLWMICC:
		return OPVCC(20, 0, 0, 1)
	case ARLDMI:
		return OPVCC(30, 0, 0, 0) | 3<<2 
	case ARLDMICC:
		return OPVCC(30, 0, 0, 1) | 3<<2
	case ARLDIMI:
		return OPVCC(30, 0, 0, 0) | 3<<2 
	case ARLDIMICC:
		return OPVCC(30, 0, 0, 1) | 3<<2
	case ARLWNM:
		return OPVCC(21, 0, 0, 0) 
	case ARLWNMCC:
		return OPVCC(21, 0, 0, 1)

	case ARLDCL:
		return OPVCC(30, 0, 0, 0) 
	case ARLDCLCC:
		return OPVCC(30, 0, 0, 1)
	case ARLDCR:
		return OPVCC(30, 1, 0, 0) 
	case ARLDCRCC:
		return OPVCC(30, 1, 0, 1)
	case ARLDC:
		return OPVCC(30, 0, 0, 0) | 2<<2
	case ARLDCCC:
		return OPVCC(30, 0, 0, 1) | 2<<2

	case ASRAW:
		return OPVCC(31, 824, 0, 0)
	case ASRAWCC:
		return OPVCC(31, 824, 0, 1)
	case ASRAD:
		return OPVCC(31, (413 << 1), 0, 0)
	case ASRADCC:
		return OPVCC(31, (413 << 1), 0, 1)

	case ASTSW:
		return OPVCC(31, 725, 0, 0)

	case ASUBC:
		return OPVCC(8, 0, 0, 0)

	case ATW:
		return OPVCC(3, 0, 0, 0)
	case ATD:
		return OPVCC(2, 0, 0, 0)

	
	
	
	case AVSPLTB:
		return OPVX(4, 524, 0, 0) 
	case AVSPLTH:
		return OPVX(4, 588, 0, 0) 
	case AVSPLTW:
		return OPVX(4, 652, 0, 0) 

	case AVSPLTISB:
		return OPVX(4, 780, 0, 0) 
	case AVSPLTISH:
		return OPVX(4, 844, 0, 0) 
	case AVSPLTISW:
		return OPVX(4, 908, 0, 0) 
	

	case AFTDIV:
		return OPVCC(63, 128, 0, 0) 
	case AFTSQRT:
		return OPVCC(63, 160, 0, 0) 

	case AXOR:
		return OPVCC(26, 0, 0, 0) 
	case AXORIS:
		return OPVCC(27, 0, 0, 0) 
	}

	c.ctxt.Diag("bad opcode i/r or i/r/r %v", a)
	return 0
}


func (c *ctxt9) opload(a obj.As) uint32 {
	switch a {
	case AMOVD:
		return OPVCC(58, 0, 0, 0) 
	case AMOVDU:
		return OPVCC(58, 0, 0, 1) 
	case AMOVWZ:
		return OPVCC(32, 0, 0, 0) 
	case AMOVWZU:
		return OPVCC(33, 0, 0, 0) 
	case AMOVW:
		return OPVCC(58, 0, 0, 0) | 1<<1 
	case ALXV:
		return OPDQ(61, 1, 0) 
	case ALXVL:
		return OPVXX1(31, 269, 0) 
	case ALXVLL:
		return OPVXX1(31, 301, 0) 
	case ALXVX:
		return OPVXX1(31, 268, 0) 

		
	case AMOVB, AMOVBZ:
		return OPVCC(34, 0, 0, 0)
		

	case AMOVBU, AMOVBZU:
		return OPVCC(35, 0, 0, 0)
	case AFMOVD:
		return OPVCC(50, 0, 0, 0)
	case AFMOVDU:
		return OPVCC(51, 0, 0, 0)
	case AFMOVS:
		return OPVCC(48, 0, 0, 0)
	case AFMOVSU:
		return OPVCC(49, 0, 0, 0)
	case AMOVH:
		return OPVCC(42, 0, 0, 0)
	case AMOVHU:
		return OPVCC(43, 0, 0, 0)
	case AMOVHZ:
		return OPVCC(40, 0, 0, 0)
	case AMOVHZU:
		return OPVCC(41, 0, 0, 0)
	case AMOVMW:
		return OPVCC(46, 0, 0, 0) 
	}

	c.ctxt.Diag("bad load opcode %v", a)
	return 0
}


func (c *ctxt9) oploadx(a obj.As) uint32 {
	switch a {
	case AMOVWZ:
		return OPVCC(31, 23, 0, 0) 
	case AMOVWZU:
		return OPVCC(31, 55, 0, 0) 
	case AMOVW:
		return OPVCC(31, 341, 0, 0) 
	case AMOVWU:
		return OPVCC(31, 373, 0, 0) 

	case AMOVB, AMOVBZ:
		return OPVCC(31, 87, 0, 0) 

	case AMOVBU, AMOVBZU:
		return OPVCC(31, 119, 0, 0) 
	case AFMOVD:
		return OPVCC(31, 599, 0, 0) 
	case AFMOVDU:
		return OPVCC(31, 631, 0, 0) 
	case AFMOVS:
		return OPVCC(31, 535, 0, 0) 
	case AFMOVSU:
		return OPVCC(31, 567, 0, 0) 
	case AFMOVSX:
		return OPVCC(31, 855, 0, 0) 
	case AFMOVSZ:
		return OPVCC(31, 887, 0, 0) 
	case AMOVH:
		return OPVCC(31, 343, 0, 0) 
	case AMOVHU:
		return OPVCC(31, 375, 0, 0) 
	case AMOVHBR:
		return OPVCC(31, 790, 0, 0) 
	case AMOVWBR:
		return OPVCC(31, 534, 0, 0) 
	case AMOVDBR:
		return OPVCC(31, 532, 0, 0) 
	case AMOVHZ:
		return OPVCC(31, 279, 0, 0) 
	case AMOVHZU:
		return OPVCC(31, 311, 0, 0) 
	case AECIWX:
		return OPVCC(31, 310, 0, 0) 
	case ALBAR:
		return OPVCC(31, 52, 0, 0) 
	case ALHAR:
		return OPVCC(31, 116, 0, 0) 
	case ALWAR:
		return OPVCC(31, 20, 0, 0) 
	case ALDAR:
		return OPVCC(31, 84, 0, 0) 
	case ALSW:
		return OPVCC(31, 533, 0, 0) 
	case AMOVD:
		return OPVCC(31, 21, 0, 0) 
	case AMOVDU:
		return OPVCC(31, 53, 0, 0) 
	case ALDMX:
		return OPVCC(31, 309, 0, 0) 

	
	case ALVEBX:
		return OPVCC(31, 7, 0, 0) 
	case ALVEHX:
		return OPVCC(31, 39, 0, 0) 
	case ALVEWX:
		return OPVCC(31, 71, 0, 0) 
	case ALVX:
		return OPVCC(31, 103, 0, 0) 
	case ALVXL:
		return OPVCC(31, 359, 0, 0) 
	case ALVSL:
		return OPVCC(31, 6, 0, 0) 
	case ALVSR:
		return OPVCC(31, 38, 0, 0) 
		

	
	case ALXVX:
		return OPVXX1(31, 268, 0) 
	case ALXVD2X:
		return OPVXX1(31, 844, 0) 
	case ALXVW4X:
		return OPVXX1(31, 780, 0) 
	case ALXVH8X:
		return OPVXX1(31, 812, 0) 
	case ALXVB16X:
		return OPVXX1(31, 876, 0) 
	case ALXVDSX:
		return OPVXX1(31, 332, 0) 
	case ALXSDX:
		return OPVXX1(31, 588, 0) 
	case ALXSIWAX:
		return OPVXX1(31, 76, 0) 
	case ALXSIWZX:
		return OPVXX1(31, 12, 0) 
	}

	c.ctxt.Diag("bad loadx opcode %v", a)
	return 0
}


func (c *ctxt9) opstore(a obj.As) uint32 {
	switch a {
	case AMOVB, AMOVBZ:
		return OPVCC(38, 0, 0, 0) 

	case AMOVBU, AMOVBZU:
		return OPVCC(39, 0, 0, 0) 
	case AFMOVD:
		return OPVCC(54, 0, 0, 0) 
	case AFMOVDU:
		return OPVCC(55, 0, 0, 0) 
	case AFMOVS:
		return OPVCC(52, 0, 0, 0) 
	case AFMOVSU:
		return OPVCC(53, 0, 0, 0) 

	case AMOVHZ, AMOVH:
		return OPVCC(44, 0, 0, 0) 

	case AMOVHZU, AMOVHU:
		return OPVCC(45, 0, 0, 0) 
	case AMOVMW:
		return OPVCC(47, 0, 0, 0) 
	case ASTSW:
		return OPVCC(31, 725, 0, 0) 

	case AMOVWZ, AMOVW:
		return OPVCC(36, 0, 0, 0) 

	case AMOVWZU, AMOVWU:
		return OPVCC(37, 0, 0, 0) 
	case AMOVD:
		return OPVCC(62, 0, 0, 0) 
	case AMOVDU:
		return OPVCC(62, 0, 0, 1) 
	case ASTXV:
		return OPDQ(61, 5, 0) 
	case ASTXVL:
		return OPVXX1(31, 397, 0) 
	case ASTXVLL:
		return OPVXX1(31, 429, 0) 
	case ASTXVX:
		return OPVXX1(31, 396, 0) 

	}

	c.ctxt.Diag("unknown store opcode %v", a)
	return 0
}


func (c *ctxt9) opstorex(a obj.As) uint32 {
	switch a {
	case AMOVB, AMOVBZ:
		return OPVCC(31, 215, 0, 0) 

	case AMOVBU, AMOVBZU:
		return OPVCC(31, 247, 0, 0) 
	case AFMOVD:
		return OPVCC(31, 727, 0, 0) 
	case AFMOVDU:
		return OPVCC(31, 759, 0, 0) 
	case AFMOVS:
		return OPVCC(31, 663, 0, 0) 
	case AFMOVSU:
		return OPVCC(31, 695, 0, 0) 
	case AFMOVSX:
		return OPVCC(31, 983, 0, 0) 

	case AMOVHZ, AMOVH:
		return OPVCC(31, 407, 0, 0) 
	case AMOVHBR:
		return OPVCC(31, 918, 0, 0) 

	case AMOVHZU, AMOVHU:
		return OPVCC(31, 439, 0, 0) 

	case AMOVWZ, AMOVW:
		return OPVCC(31, 151, 0, 0) 

	case AMOVWZU, AMOVWU:
		return OPVCC(31, 183, 0, 0) 
	case ASTSW:
		return OPVCC(31, 661, 0, 0) 
	case AMOVWBR:
		return OPVCC(31, 662, 0, 0) 
	case AMOVDBR:
		return OPVCC(31, 660, 0, 0) 
	case ASTBCCC:
		return OPVCC(31, 694, 0, 1) 
	case ASTHCCC:
		return OPVCC(31, 726, 0, 1) 
	case ASTWCCC:
		return OPVCC(31, 150, 0, 1) 
	case ASTDCCC:
		return OPVCC(31, 214, 0, 1) 
	case AECOWX:
		return OPVCC(31, 438, 0, 0) 
	case AMOVD:
		return OPVCC(31, 149, 0, 0) 
	case AMOVDU:
		return OPVCC(31, 181, 0, 0) 

	
	case ASTVEBX:
		return OPVCC(31, 135, 0, 0) 
	case ASTVEHX:
		return OPVCC(31, 167, 0, 0) 
	case ASTVEWX:
		return OPVCC(31, 199, 0, 0) 
	case ASTVX:
		return OPVCC(31, 231, 0, 0) 
	case ASTVXL:
		return OPVCC(31, 487, 0, 0) 
		

	
	case ASTXVX:
		return OPVXX1(31, 396, 0) 
	case ASTXVD2X:
		return OPVXX1(31, 972, 0) 
	case ASTXVW4X:
		return OPVXX1(31, 908, 0) 
	case ASTXVH8X:
		return OPVXX1(31, 940, 0) 
	case ASTXVB16X:
		return OPVXX1(31, 1004, 0) 

	case ASTXSDX:
		return OPVXX1(31, 716, 0) 

	case ASTXSIWX:
		return OPVXX1(31, 140, 0) 

		

	}

	c.ctxt.Diag("unknown storex opcode %v", a)
	return 0
}
