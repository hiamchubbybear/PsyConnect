




























package s390x

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/objabi"
	"fmt"
	"log"
	"math"
	"sort"
)




type ctxtz struct {
	ctxt       *obj.Link
	newprog    obj.ProgAlloc
	cursym     *obj.LSym
	autosize   int32
	instoffset int64
	pc         int64
}


const (
	funcAlign = 16
)

type Optab struct {
	as obj.As 
	i  uint8  
	a1 uint8  
	a2 uint8  
	a3 uint8  
	a4 uint8  
	a5 uint8  
	a6 uint8  
}

var optab = []Optab{
	
	{i: 0, as: obj.ATEXT, a1: C_ADDR, a6: C_TEXTSIZE},
	{i: 0, as: obj.ATEXT, a1: C_ADDR, a3: C_LCON, a6: C_TEXTSIZE},
	{i: 0, as: obj.APCDATA, a1: C_LCON, a6: C_LCON},
	{i: 0, as: obj.AFUNCDATA, a1: C_SCON, a6: C_ADDR},
	{i: 0, as: obj.ANOP},
	{i: 0, as: obj.ANOP, a1: C_SAUTO},

	
	{i: 1, as: AMOVD, a1: C_REG, a6: C_REG},
	{i: 1, as: AMOVB, a1: C_REG, a6: C_REG},
	{i: 1, as: AMOVBZ, a1: C_REG, a6: C_REG},
	{i: 1, as: AMOVW, a1: C_REG, a6: C_REG},
	{i: 1, as: AMOVWZ, a1: C_REG, a6: C_REG},
	{i: 1, as: AFMOVD, a1: C_FREG, a6: C_FREG},
	{i: 1, as: AMOVDBR, a1: C_REG, a6: C_REG},

	
	{i: 26, as: AMOVD, a1: C_LACON, a6: C_REG},
	{i: 26, as: AMOVW, a1: C_LACON, a6: C_REG},
	{i: 26, as: AMOVWZ, a1: C_LACON, a6: C_REG},
	{i: 3, as: AMOVD, a1: C_DCON, a6: C_REG},
	{i: 3, as: AMOVW, a1: C_DCON, a6: C_REG},
	{i: 3, as: AMOVWZ, a1: C_DCON, a6: C_REG},
	{i: 3, as: AMOVB, a1: C_DCON, a6: C_REG},
	{i: 3, as: AMOVBZ, a1: C_DCON, a6: C_REG},

	
	{i: 72, as: AMOVD, a1: C_SCON, a6: C_LAUTO},
	{i: 72, as: AMOVD, a1: C_ADDCON, a6: C_LAUTO},
	{i: 72, as: AMOVW, a1: C_SCON, a6: C_LAUTO},
	{i: 72, as: AMOVW, a1: C_ADDCON, a6: C_LAUTO},
	{i: 72, as: AMOVWZ, a1: C_SCON, a6: C_LAUTO},
	{i: 72, as: AMOVWZ, a1: C_ADDCON, a6: C_LAUTO},
	{i: 72, as: AMOVB, a1: C_SCON, a6: C_LAUTO},
	{i: 72, as: AMOVB, a1: C_ADDCON, a6: C_LAUTO},
	{i: 72, as: AMOVBZ, a1: C_SCON, a6: C_LAUTO},
	{i: 72, as: AMOVBZ, a1: C_ADDCON, a6: C_LAUTO},
	{i: 72, as: AMOVD, a1: C_SCON, a6: C_LOREG},
	{i: 72, as: AMOVD, a1: C_ADDCON, a6: C_LOREG},
	{i: 72, as: AMOVW, a1: C_SCON, a6: C_LOREG},
	{i: 72, as: AMOVW, a1: C_ADDCON, a6: C_LOREG},
	{i: 72, as: AMOVWZ, a1: C_SCON, a6: C_LOREG},
	{i: 72, as: AMOVWZ, a1: C_ADDCON, a6: C_LOREG},
	{i: 72, as: AMOVB, a1: C_SCON, a6: C_LOREG},
	{i: 72, as: AMOVB, a1: C_ADDCON, a6: C_LOREG},
	{i: 72, as: AMOVBZ, a1: C_SCON, a6: C_LOREG},
	{i: 72, as: AMOVBZ, a1: C_ADDCON, a6: C_LOREG},

	
	{i: 35, as: AMOVD, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVW, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVWZ, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVBZ, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVB, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVDBR, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVHBR, a1: C_REG, a6: C_LAUTO},
	{i: 35, as: AMOVD, a1: C_REG, a6: C_LOREG},
	{i: 35, as: AMOVW, a1: C_REG, a6: C_LOREG},
	{i: 35, as: AMOVWZ, a1: C_REG, a6: C_LOREG},
	{i: 35, as: AMOVBZ, a1: C_REG, a6: C_LOREG},
	{i: 35, as: AMOVB, a1: C_REG, a6: C_LOREG},
	{i: 35, as: AMOVDBR, a1: C_REG, a6: C_LOREG},
	{i: 35, as: AMOVHBR, a1: C_REG, a6: C_LOREG},
	{i: 74, as: AMOVD, a1: C_REG, a6: C_ADDR},
	{i: 74, as: AMOVW, a1: C_REG, a6: C_ADDR},
	{i: 74, as: AMOVWZ, a1: C_REG, a6: C_ADDR},
	{i: 74, as: AMOVBZ, a1: C_REG, a6: C_ADDR},
	{i: 74, as: AMOVB, a1: C_REG, a6: C_ADDR},

	
	{i: 36, as: AMOVD, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVW, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVWZ, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVBZ, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVB, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVDBR, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVHBR, a1: C_LAUTO, a6: C_REG},
	{i: 36, as: AMOVD, a1: C_LOREG, a6: C_REG},
	{i: 36, as: AMOVW, a1: C_LOREG, a6: C_REG},
	{i: 36, as: AMOVWZ, a1: C_LOREG, a6: C_REG},
	{i: 36, as: AMOVBZ, a1: C_LOREG, a6: C_REG},
	{i: 36, as: AMOVB, a1: C_LOREG, a6: C_REG},
	{i: 36, as: AMOVDBR, a1: C_LOREG, a6: C_REG},
	{i: 36, as: AMOVHBR, a1: C_LOREG, a6: C_REG},
	{i: 75, as: AMOVD, a1: C_ADDR, a6: C_REG},
	{i: 75, as: AMOVW, a1: C_ADDR, a6: C_REG},
	{i: 75, as: AMOVWZ, a1: C_ADDR, a6: C_REG},
	{i: 75, as: AMOVBZ, a1: C_ADDR, a6: C_REG},
	{i: 75, as: AMOVB, a1: C_ADDR, a6: C_REG},

	
	{i: 99, as: ALAAG, a1: C_REG, a2: C_REG, a6: C_LOREG},

	
	{i: 2, as: AADD, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 2, as: AADD, a1: C_REG, a6: C_REG},
	{i: 22, as: AADD, a1: C_LCON, a2: C_REG, a6: C_REG},
	{i: 22, as: AADD, a1: C_LCON, a6: C_REG},
	{i: 12, as: AADD, a1: C_LOREG, a6: C_REG},
	{i: 12, as: AADD, a1: C_LAUTO, a6: C_REG},
	{i: 21, as: ASUB, a1: C_LCON, a2: C_REG, a6: C_REG},
	{i: 21, as: ASUB, a1: C_LCON, a6: C_REG},
	{i: 12, as: ASUB, a1: C_LOREG, a6: C_REG},
	{i: 12, as: ASUB, a1: C_LAUTO, a6: C_REG},
	{i: 4, as: AMULHD, a1: C_REG, a6: C_REG},
	{i: 4, as: AMULHD, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 62, as: AMLGR, a1: C_REG, a6: C_REG},
	{i: 2, as: ADIVW, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 2, as: ADIVW, a1: C_REG, a6: C_REG},
	{i: 10, as: ASUB, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 10, as: ASUB, a1: C_REG, a6: C_REG},
	{i: 47, as: ANEG, a1: C_REG, a6: C_REG},
	{i: 47, as: ANEG, a6: C_REG},

	
	{i: 6, as: AAND, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 6, as: AAND, a1: C_REG, a6: C_REG},
	{i: 23, as: AAND, a1: C_LCON, a6: C_REG},
	{i: 12, as: AAND, a1: C_LOREG, a6: C_REG},
	{i: 12, as: AAND, a1: C_LAUTO, a6: C_REG},
	{i: 6, as: AANDW, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 6, as: AANDW, a1: C_REG, a6: C_REG},
	{i: 24, as: AANDW, a1: C_LCON, a6: C_REG},
	{i: 12, as: AANDW, a1: C_LOREG, a6: C_REG},
	{i: 12, as: AANDW, a1: C_LAUTO, a6: C_REG},
	{i: 7, as: ASLD, a1: C_REG, a6: C_REG},
	{i: 7, as: ASLD, a1: C_REG, a2: C_REG, a6: C_REG},
	{i: 7, as: ASLD, a1: C_SCON, a2: C_REG, a6: C_REG},
	{i: 7, as: ASLD, a1: C_SCON, a6: C_REG},
	{i: 13, as: ARNSBG, a1: C_SCON, a3: C_SCON, a4: C_SCON, a5: C_REG, a6: C_REG},

	
	{i: 79, as: ACSG, a1: C_REG, a2: C_REG, a6: C_SOREG},

	
	{i: 32, as: AFADD, a1: C_FREG, a6: C_FREG},
	{i: 33, as: AFABS, a1: C_FREG, a6: C_FREG},
	{i: 33, as: AFABS, a6: C_FREG},
	{i: 34, as: AFMADD, a1: C_FREG, a2: C_FREG, a6: C_FREG},
	{i: 32, as: AFMUL, a1: C_FREG, a6: C_FREG},
	{i: 36, as: AFMOVD, a1: C_LAUTO, a6: C_FREG},
	{i: 36, as: AFMOVD, a1: C_LOREG, a6: C_FREG},
	{i: 75, as: AFMOVD, a1: C_ADDR, a6: C_FREG},
	{i: 35, as: AFMOVD, a1: C_FREG, a6: C_LAUTO},
	{i: 35, as: AFMOVD, a1: C_FREG, a6: C_LOREG},
	{i: 74, as: AFMOVD, a1: C_FREG, a6: C_ADDR},
	{i: 67, as: AFMOVD, a1: C_ZCON, a6: C_FREG},
	{i: 81, as: ALDGR, a1: C_REG, a6: C_FREG},
	{i: 81, as: ALGDR, a1: C_FREG, a6: C_REG},
	{i: 82, as: ACEFBRA, a1: C_REG, a6: C_FREG},
	{i: 83, as: ACFEBRA, a1: C_FREG, a6: C_REG},
	{i: 48, as: AFIEBR, a1: C_SCON, a2: C_FREG, a6: C_FREG},
	{i: 49, as: ACPSDR, a1: C_FREG, a2: C_FREG, a6: C_FREG},
	{i: 50, as: ALTDBR, a1: C_FREG, a6: C_FREG},
	{i: 51, as: ATCDB, a1: C_FREG, a6: C_SCON},

	
	{i: 19, as: AMOVD, a1: C_SYMADDR, a6: C_REG},
	{i: 93, as: AMOVD, a1: C_GOTADDR, a6: C_REG},
	{i: 94, as: AMOVD, a1: C_TLS_LE, a6: C_REG},
	{i: 95, as: AMOVD, a1: C_TLS_IE, a6: C_REG},

	
	{i: 5, as: ASYSCALL},
	{i: 77, as: ASYSCALL, a1: C_SCON},

	
	{i: 16, as: ABEQ, a6: C_SBRA},
	{i: 16, as: ABRC, a1: C_SCON, a6: C_SBRA},
	{i: 11, as: ABR, a6: C_LBRA},
	{i: 16, as: ABC, a1: C_SCON, a2: C_REG, a6: C_LBRA},
	{i: 18, as: ABR, a6: C_REG},
	{i: 18, as: ABR, a1: C_REG, a6: C_REG},
	{i: 15, as: ABR, a6: C_ZOREG},
	{i: 15, as: ABC, a6: C_ZOREG},

	
	{i: 89, as: ACGRJ, a1: C_SCON, a2: C_REG, a3: C_REG, a6: C_SBRA},
	{i: 89, as: ACMPBEQ, a1: C_REG, a2: C_REG, a6: C_SBRA},
	{i: 89, as: ACLGRJ, a1: C_SCON, a2: C_REG, a3: C_REG, a6: C_SBRA},
	{i: 89, as: ACMPUBEQ, a1: C_REG, a2: C_REG, a6: C_SBRA},
	{i: 90, as: ACGIJ, a1: C_SCON, a2: C_REG, a3: C_ADDCON, a6: C_SBRA},
	{i: 90, as: ACGIJ, a1: C_SCON, a2: C_REG, a3: C_SCON, a6: C_SBRA},
	{i: 90, as: ACMPBEQ, a1: C_REG, a3: C_ADDCON, a6: C_SBRA},
	{i: 90, as: ACMPBEQ, a1: C_REG, a3: C_SCON, a6: C_SBRA},
	{i: 90, as: ACLGIJ, a1: C_SCON, a2: C_REG, a3: C_ADDCON, a6: C_SBRA},
	{i: 90, as: ACMPUBEQ, a1: C_REG, a3: C_ANDCON, a6: C_SBRA},

	
	{i: 41, as: ABRCT, a1: C_REG, a6: C_SBRA},
	{i: 41, as: ABRCTG, a1: C_REG, a6: C_SBRA},

	
	{i: 17, as: AMOVDEQ, a1: C_REG, a6: C_REG},

	
	{i: 25, as: ALOCGR, a1: C_SCON, a2: C_REG, a6: C_REG},

	
	{i: 8, as: AFLOGR, a1: C_REG, a6: C_REG},

	
	{i: 9, as: APOPCNT, a1: C_REG, a6: C_REG},

	
	{i: 70, as: ACMP, a1: C_REG, a6: C_REG},
	{i: 71, as: ACMP, a1: C_REG, a6: C_LCON},
	{i: 70, as: ACMPU, a1: C_REG, a6: C_REG},
	{i: 71, as: ACMPU, a1: C_REG, a6: C_LCON},
	{i: 70, as: AFCMPO, a1: C_FREG, a6: C_FREG},
	{i: 70, as: AFCMPO, a1: C_FREG, a2: C_REG, a6: C_FREG},

	
	{i: 91, as: ATMHH, a1: C_REG, a6: C_ANDCON},

	
	{i: 92, as: AIPM, a1: C_REG},

	
	{i: 76, as: ASPM, a1: C_REG},

	
	{i: 68, as: AMOVW, a1: C_AREG, a6: C_REG},
	{i: 68, as: AMOVWZ, a1: C_AREG, a6: C_REG},
	{i: 69, as: AMOVW, a1: C_REG, a6: C_AREG},
	{i: 69, as: AMOVWZ, a1: C_REG, a6: C_AREG},

	
	{i: 96, as: ACLEAR, a1: C_LCON, a6: C_LOREG},
	{i: 96, as: ACLEAR, a1: C_LCON, a6: C_LAUTO},

	
	{i: 97, as: ASTMG, a1: C_REG, a2: C_REG, a6: C_LOREG},
	{i: 97, as: ASTMG, a1: C_REG, a2: C_REG, a6: C_LAUTO},
	{i: 98, as: ALMG, a1: C_LOREG, a2: C_REG, a6: C_REG},
	{i: 98, as: ALMG, a1: C_LAUTO, a2: C_REG, a6: C_REG},

	
	{i: 40, as: ABYTE, a1: C_SCON},
	{i: 40, as: AWORD, a1: C_LCON},
	{i: 31, as: ADWORD, a1: C_LCON},
	{i: 31, as: ADWORD, a1: C_DCON},

	
	{i: 80, as: ASYNC},

	
	{i: 88, as: ASTCK, a6: C_SAUTO},
	{i: 88, as: ASTCK, a6: C_SOREG},

	
	{i: 84, as: AMVC, a1: C_SCON, a3: C_LOREG, a6: C_LOREG},
	{i: 84, as: AMVC, a1: C_SCON, a3: C_LOREG, a6: C_LAUTO},
	{i: 84, as: AMVC, a1: C_SCON, a3: C_LAUTO, a6: C_LAUTO},

	
	{i: 85, as: ALARL, a1: C_LCON, a6: C_REG},
	{i: 85, as: ALARL, a1: C_SYMADDR, a6: C_REG},
	{i: 86, as: ALA, a1: C_SOREG, a6: C_REG},
	{i: 86, as: ALA, a1: C_SAUTO, a6: C_REG},
	{i: 87, as: AEXRL, a1: C_SYMADDR, a6: C_REG},

	
	{i: 78, as: obj.AUNDEF},

	
	{i: 66, as: ANOPH},

	

	
	{i: 100, as: AVST, a1: C_VREG, a6: C_SOREG},
	{i: 100, as: AVST, a1: C_VREG, a6: C_SAUTO},
	{i: 100, as: AVSTEG, a1: C_SCON, a2: C_VREG, a6: C_SOREG},
	{i: 100, as: AVSTEG, a1: C_SCON, a2: C_VREG, a6: C_SAUTO},

	
	{i: 101, as: AVL, a1: C_SOREG, a6: C_VREG},
	{i: 101, as: AVL, a1: C_SAUTO, a6: C_VREG},
	{i: 101, as: AVLEG, a1: C_SCON, a3: C_SOREG, a6: C_VREG},
	{i: 101, as: AVLEG, a1: C_SCON, a3: C_SAUTO, a6: C_VREG},

	
	{i: 102, as: AVSCEG, a1: C_SCON, a2: C_VREG, a6: C_SOREG},
	{i: 102, as: AVSCEG, a1: C_SCON, a2: C_VREG, a6: C_SAUTO},

	
	{i: 103, as: AVGEG, a1: C_SCON, a3: C_SOREG, a6: C_VREG},
	{i: 103, as: AVGEG, a1: C_SCON, a3: C_SAUTO, a6: C_VREG},

	
	{i: 104, as: AVESLG, a1: C_SCON, a2: C_VREG, a6: C_VREG},
	{i: 104, as: AVESLG, a1: C_REG, a2: C_VREG, a6: C_VREG},
	{i: 104, as: AVESLG, a1: C_SCON, a6: C_VREG},
	{i: 104, as: AVESLG, a1: C_REG, a6: C_VREG},
	{i: 104, as: AVLGVG, a1: C_SCON, a2: C_VREG, a6: C_REG},
	{i: 104, as: AVLGVG, a1: C_REG, a2: C_VREG, a6: C_REG},
	{i: 104, as: AVLVGG, a1: C_SCON, a2: C_REG, a6: C_VREG},
	{i: 104, as: AVLVGG, a1: C_REG, a2: C_REG, a6: C_VREG},

	
	{i: 105, as: AVSTM, a1: C_VREG, a2: C_VREG, a6: C_SOREG},
	{i: 105, as: AVSTM, a1: C_VREG, a2: C_VREG, a6: C_SAUTO},

	
	{i: 106, as: AVLM, a1: C_SOREG, a2: C_VREG, a6: C_VREG},
	{i: 106, as: AVLM, a1: C_SAUTO, a2: C_VREG, a6: C_VREG},

	
	{i: 107, as: AVSTL, a1: C_REG, a2: C_VREG, a6: C_SOREG},
	{i: 107, as: AVSTL, a1: C_REG, a2: C_VREG, a6: C_SAUTO},

	
	{i: 108, as: AVLL, a1: C_REG, a3: C_SOREG, a6: C_VREG},
	{i: 108, as: AVLL, a1: C_REG, a3: C_SAUTO, a6: C_VREG},

	
	{i: 109, as: AVGBM, a1: C_ANDCON, a6: C_VREG},
	{i: 109, as: AVZERO, a6: C_VREG},
	{i: 109, as: AVREPIG, a1: C_ADDCON, a6: C_VREG},
	{i: 109, as: AVREPIG, a1: C_SCON, a6: C_VREG},
	{i: 109, as: AVLEIG, a1: C_SCON, a3: C_ADDCON, a6: C_VREG},
	{i: 109, as: AVLEIG, a1: C_SCON, a3: C_SCON, a6: C_VREG},

	
	{i: 110, as: AVGMG, a1: C_SCON, a3: C_SCON, a6: C_VREG},

	
	{i: 111, as: AVREPG, a1: C_UCON, a2: C_VREG, a6: C_VREG},

	
	
	{i: 112, as: AVERIMG, a1: C_SCON, a2: C_VREG, a3: C_VREG, a6: C_VREG},
	{i: 112, as: AVSLDB, a1: C_SCON, a2: C_VREG, a3: C_VREG, a6: C_VREG},

	
	{i: 113, as: AVFTCIDB, a1: C_SCON, a2: C_VREG, a6: C_VREG},

	
	{i: 114, as: AVLR, a1: C_VREG, a6: C_VREG},

	
	{i: 115, as: AVECG, a1: C_VREG, a6: C_VREG},

	
	{i: 117, as: AVCEQG, a1: C_VREG, a2: C_VREG, a6: C_VREG},
	{i: 117, as: AVFAEF, a1: C_VREG, a2: C_VREG, a6: C_VREG},
	{i: 117, as: AVPKSG, a1: C_VREG, a2: C_VREG, a6: C_VREG},

	
	{i: 118, as: AVAQ, a1: C_VREG, a2: C_VREG, a6: C_VREG},
	{i: 118, as: AVAQ, a1: C_VREG, a6: C_VREG},
	{i: 118, as: AVNOT, a1: C_VREG, a6: C_VREG},
	{i: 123, as: AVPDI, a1: C_SCON, a2: C_VREG, a3: C_VREG, a6: C_VREG},

	
	{i: 119, as: AVERLLVG, a1: C_VREG, a2: C_VREG, a6: C_VREG},
	{i: 119, as: AVERLLVG, a1: C_VREG, a6: C_VREG},

	
	{i: 120, as: AVACQ, a1: C_VREG, a2: C_VREG, a3: C_VREG, a6: C_VREG},

	
	{i: 121, as: AVSEL, a1: C_VREG, a2: C_VREG, a3: C_VREG, a6: C_VREG},

	
	{i: 122, as: AVLVGP, a1: C_REG, a2: C_REG, a6: C_VREG},
}

var oprange [ALAST & obj.AMask][]Optab

var xcmp [C_NCLASS][C_NCLASS]bool

func spanz(ctxt *obj.Link, cursym *obj.LSym, newprog obj.ProgAlloc) {
	if ctxt.Retpoline {
		ctxt.Diag("-spectre=ret not supported on s390x")
		ctxt.Retpoline = false 
	}

	p := cursym.Func.Text
	if p == nil || p.Link == nil { 
		return
	}

	if oprange[AORW&obj.AMask] == nil {
		ctxt.Diag("s390x ops not initialized, call s390x.buildop first")
	}

	c := ctxtz{ctxt: ctxt, newprog: newprog, cursym: cursym, autosize: int32(p.To.Offset)}

	buffer := make([]byte, 0)
	changed := true
	loop := 0
	for changed {
		if loop > 100 {
			c.ctxt.Diag("stuck in spanz loop")
			break
		}
		changed = false
		buffer = buffer[:0]
		c.cursym.R = make([]obj.Reloc, 0)
		for p := c.cursym.Func.Text; p != nil; p = p.Link {
			pc := int64(len(buffer))
			if pc != p.Pc {
				changed = true
			}
			p.Pc = pc
			c.pc = p.Pc
			c.asmout(p, &buffer)
			if pc == int64(len(buffer)) {
				switch p.As {
				case obj.ANOP, obj.AFUNCDATA, obj.APCDATA, obj.ATEXT:
					
				default:
					c.ctxt.Diag("zero-width instruction\n%v", p)
				}
			}
		}
		loop++
	}

	c.cursym.Size = int64(len(buffer))
	if c.cursym.Size%funcAlign != 0 {
		c.cursym.Size += funcAlign - (c.cursym.Size % funcAlign)
	}
	c.cursym.Grow(c.cursym.Size)
	copy(c.cursym.P, buffer)

	
	
	
	
	obj.MarkUnsafePoints(c.ctxt, c.cursym.Func.Text, c.newprog, c.isUnsafePoint, nil)
}


func (c *ctxtz) isUnsafePoint(p *obj.Prog) bool {
	if p.From.Reg == REGTMP || p.To.Reg == REGTMP || p.Reg == REGTMP {
		return true
	}
	for _, a := range p.RestArgs {
		if a.Reg == REGTMP {
			return true
		}
	}
	return p.Mark&USETMP != 0
}

func isint32(v int64) bool {
	return int64(int32(v)) == v
}

func isuint32(v uint64) bool {
	return uint64(uint32(v)) == v
}

func (c *ctxtz) aclass(a *obj.Addr) int {
	switch a.Type {
	case obj.TYPE_NONE:
		return C_NONE

	case obj.TYPE_REG:
		if REG_R0 <= a.Reg && a.Reg <= REG_R15 {
			return C_REG
		}
		if REG_F0 <= a.Reg && a.Reg <= REG_F15 {
			return C_FREG
		}
		if REG_AR0 <= a.Reg && a.Reg <= REG_AR15 {
			return C_AREG
		}
		if REG_V0 <= a.Reg && a.Reg <= REG_V31 {
			return C_VREG
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
			if a.Sym.Type == objabi.STLSBSS {
				if c.ctxt.Flag_shared {
					return C_TLS_IE 
				}
				return C_TLS_LE 
			}
			return C_ADDR

		case obj.NAME_GOTREF:
			return C_GOTADDR

		case obj.NAME_AUTO:
			if a.Reg == REGSP {
				
				
				a.Reg = obj.REG_NONE
			}
			c.instoffset = int64(c.autosize) + a.Offset
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SAUTO
			}
			return C_LAUTO

		case obj.NAME_PARAM:
			if a.Reg == REGSP {
				
				
				a.Reg = obj.REG_NONE
			}
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
		if f64, ok := a.Val.(float64); ok && math.Float64bits(f64) == 0 {
			return C_ZCON
		}
		c.ctxt.Diag("cannot handle the floating point constant %v", a.Val)

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

			return C_SYMADDR

		case obj.NAME_AUTO:
			if a.Reg == REGSP {
				
				
				a.Reg = obj.REG_NONE
			}
			c.instoffset = int64(c.autosize) + a.Offset
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SACON
			}
			return C_LACON

		case obj.NAME_PARAM:
			if a.Reg == REGSP {
				
				
				a.Reg = obj.REG_NONE
			}
			c.instoffset = int64(c.autosize) + a.Offset + c.ctxt.FixedFrameSize()
			if c.instoffset >= -BIG && c.instoffset < BIG {
				return C_SACON
			}
			return C_LACON

		default:
			return C_GOK
		}

		if c.instoffset == 0 {
			return C_ZCON
		}
		if c.instoffset >= 0 {
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
		return C_SBRA
	}

	return C_GOK
}

func (c *ctxtz) oplook(p *obj.Prog) *Optab {
	
	if p.Optab != 0 {
		return &optab[p.Optab-1]
	}
	if len(p.RestArgs) > 3 {
		c.ctxt.Diag("too many RestArgs: got %v, maximum is 3\n", len(p.RestArgs))
		return nil
	}

	
	p.From.Class = int8(c.aclass(&p.From) + 1)
	p.To.Class = int8(c.aclass(&p.To) + 1)
	for i := range p.RestArgs {
		p.RestArgs[i].Class = int8(c.aclass(&p.RestArgs[i]) + 1)
	}

	
	args := [...]int8{
		p.From.Class - 1,
		C_NONE, 
		C_NONE, 
		C_NONE, 
		C_NONE, 
		p.To.Class - 1,
	}
	
	switch {
	case REG_R0 <= p.Reg && p.Reg <= REG_R15:
		args[1] = C_REG
	case REG_V0 <= p.Reg && p.Reg <= REG_V31:
		args[1] = C_VREG
	case REG_F0 <= p.Reg && p.Reg <= REG_F15:
		args[1] = C_FREG
	case REG_AR0 <= p.Reg && p.Reg <= REG_AR15:
		args[1] = C_AREG
	}
	
	for i, a := range p.RestArgs {
		args[2+i] = a.Class - 1
	}

	
	ops := oprange[p.As&obj.AMask]
	cmp := [len(args)]*[C_NCLASS]bool{}
	for i := range cmp {
		cmp[i] = &xcmp[args[i]]
	}
	for i := range ops {
		op := &ops[i]
		if cmp[0][op.a1] && cmp[1][op.a2] &&
			cmp[2][op.a3] && cmp[3][op.a4] &&
			cmp[4][op.a5] && cmp[5][op.a6] {
			p.Optab = uint16(cap(optab) - cap(ops) + i + 1)
			return op
		}
	}

	
	s := ""
	for _, a := range args {
		s += fmt.Sprintf(" %v", DRconv(int(a)))
	}
	c.ctxt.Diag("illegal combination %v%v\n", p.As, s)
	c.ctxt.Diag("prog: %v\n", p)
	return nil
}

func cmp(a int, b int) bool {
	if a == b {
		return true
	}
	switch a {
	case C_DCON:
		if b == C_LCON {
			return true
		}
		fallthrough
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

	case C_UCON:
		if b == C_ZCON || b == C_SCON {
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

	case C_LAUTO:
		if b == C_SAUTO {
			return true
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
func opset(a, b obj.As) {
	oprange[a&obj.AMask] = oprange[b&obj.AMask]
}

func buildop(ctxt *obj.Link) {
	if oprange[AORW&obj.AMask] != nil {
		
		
		
		return
	}

	for i := 0; i < C_NCLASS; i++ {
		for n := 0; n < C_NCLASS; n++ {
			if cmp(n, i) {
				xcmp[i][n] = true
			}
		}
	}
	sort.Sort(ocmp(optab))
	for i := 0; i < len(optab); i++ {
		r := optab[i].as
		start := i
		for ; i+1 < len(optab); i++ {
			if optab[i+1].as != r {
				break
			}
		}
		oprange[r&obj.AMask] = optab[start : i+1]

		
		
		switch r {
		case AADD:
			opset(AADDC, r)
			opset(AADDW, r)
			opset(AADDE, r)
			opset(AMULLD, r)
			opset(AMULLW, r)
		case ADIVW:
			opset(ADIVD, r)
			opset(ADIVDU, r)
			opset(ADIVWU, r)
			opset(AMODD, r)
			opset(AMODDU, r)
			opset(AMODW, r)
			opset(AMODWU, r)
		case AMULHD:
			opset(AMULHDU, r)
		case AMOVBZ:
			opset(AMOVH, r)
			opset(AMOVHZ, r)
		case ALA:
			opset(ALAY, r)
		case AMVC:
			opset(AMVCIN, r)
			opset(ACLC, r)
			opset(AXC, r)
			opset(AOC, r)
			opset(ANC, r)
		case ASTCK:
			opset(ASTCKC, r)
			opset(ASTCKE, r)
			opset(ASTCKF, r)
		case ALAAG:
			opset(ALAA, r)
			opset(ALAAL, r)
			opset(ALAALG, r)
			opset(ALAN, r)
			opset(ALANG, r)
			opset(ALAX, r)
			opset(ALAXG, r)
			opset(ALAO, r)
			opset(ALAOG, r)
		case ASTMG:
			opset(ASTMY, r)
		case ALMG:
			opset(ALMY, r)
		case ABEQ:
			opset(ABGE, r)
			opset(ABGT, r)
			opset(ABLE, r)
			opset(ABLT, r)
			opset(ABNE, r)
			opset(ABVC, r)
			opset(ABVS, r)
			opset(ABLEU, r)
			opset(ABLTU, r)
		case ABR:
			opset(ABL, r)
		case ABC:
			opset(ABCL, r)
		case AFABS:
			opset(AFNABS, r)
			opset(ALPDFR, r)
			opset(ALNDFR, r)
			opset(AFNEG, r)
			opset(AFNEGS, r)
			opset(ALEDBR, r)
			opset(ALDEBR, r)
			opset(AFSQRT, r)
			opset(AFSQRTS, r)
		case AFADD:
			opset(AFADDS, r)
			opset(AFDIV, r)
			opset(AFDIVS, r)
			opset(AFSUB, r)
			opset(AFSUBS, r)
		case AFMADD:
			opset(AFMADDS, r)
			opset(AFMSUB, r)
			opset(AFMSUBS, r)
		case AFMUL:
			opset(AFMULS, r)
		case AFCMPO:
			opset(AFCMPU, r)
			opset(ACEBR, r)
		case AAND:
			opset(AOR, r)
			opset(AXOR, r)
		case AANDW:
			opset(AORW, r)
			opset(AXORW, r)
		case ASLD:
			opset(ASRD, r)
			opset(ASLW, r)
			opset(ASRW, r)
			opset(ASRAD, r)
			opset(ASRAW, r)
			opset(ARLL, r)
			opset(ARLLG, r)
		case ARNSBG:
			opset(ARXSBG, r)
			opset(AROSBG, r)
			opset(ARNSBGT, r)
			opset(ARXSBGT, r)
			opset(AROSBGT, r)
			opset(ARISBG, r)
			opset(ARISBGN, r)
			opset(ARISBGZ, r)
			opset(ARISBGNZ, r)
			opset(ARISBHG, r)
			opset(ARISBLG, r)
			opset(ARISBHGZ, r)
			opset(ARISBLGZ, r)
		case ACSG:
			opset(ACS, r)
		case ASUB:
			opset(ASUBC, r)
			opset(ASUBE, r)
			opset(ASUBW, r)
		case ANEG:
			opset(ANEGW, r)
		case AFMOVD:
			opset(AFMOVS, r)
		case AMOVDBR:
			opset(AMOVWBR, r)
		case ACMP:
			opset(ACMPW, r)
		case ACMPU:
			opset(ACMPWU, r)
		case ATMHH:
			opset(ATMHL, r)
			opset(ATMLH, r)
			opset(ATMLL, r)
		case ACEFBRA:
			opset(ACDFBRA, r)
			opset(ACEGBRA, r)
			opset(ACDGBRA, r)
			opset(ACELFBR, r)
			opset(ACDLFBR, r)
			opset(ACELGBR, r)
			opset(ACDLGBR, r)
		case ACFEBRA:
			opset(ACFDBRA, r)
			opset(ACGEBRA, r)
			opset(ACGDBRA, r)
			opset(ACLFEBR, r)
			opset(ACLFDBR, r)
			opset(ACLGEBR, r)
			opset(ACLGDBR, r)
		case AFIEBR:
			opset(AFIDBR, r)
		case ACMPBEQ:
			opset(ACMPBGE, r)
			opset(ACMPBGT, r)
			opset(ACMPBLE, r)
			opset(ACMPBLT, r)
			opset(ACMPBNE, r)
		case ACMPUBEQ:
			opset(ACMPUBGE, r)
			opset(ACMPUBGT, r)
			opset(ACMPUBLE, r)
			opset(ACMPUBLT, r)
			opset(ACMPUBNE, r)
		case ACGRJ:
			opset(ACRJ, r)
		case ACLGRJ:
			opset(ACLRJ, r)
		case ACGIJ:
			opset(ACIJ, r)
		case ACLGIJ:
			opset(ACLIJ, r)
		case AMOVDEQ:
			opset(AMOVDGE, r)
			opset(AMOVDGT, r)
			opset(AMOVDLE, r)
			opset(AMOVDLT, r)
			opset(AMOVDNE, r)
		case ALOCGR:
			opset(ALOCR, r)
		case ALTDBR:
			opset(ALTEBR, r)
		case ATCDB:
			opset(ATCEB, r)
		case AVL:
			opset(AVLLEZB, r)
			opset(AVLLEZH, r)
			opset(AVLLEZF, r)
			opset(AVLLEZG, r)
			opset(AVLREPB, r)
			opset(AVLREPH, r)
			opset(AVLREPF, r)
			opset(AVLREPG, r)
		case AVLEG:
			opset(AVLBB, r)
			opset(AVLEB, r)
			opset(AVLEH, r)
			opset(AVLEF, r)
			opset(AVLEG, r)
			opset(AVLREP, r)
		case AVSTEG:
			opset(AVSTEB, r)
			opset(AVSTEH, r)
			opset(AVSTEF, r)
		case AVSCEG:
			opset(AVSCEF, r)
		case AVGEG:
			opset(AVGEF, r)
		case AVESLG:
			opset(AVESLB, r)
			opset(AVESLH, r)
			opset(AVESLF, r)
			opset(AVERLLB, r)
			opset(AVERLLH, r)
			opset(AVERLLF, r)
			opset(AVERLLG, r)
			opset(AVESRAB, r)
			opset(AVESRAH, r)
			opset(AVESRAF, r)
			opset(AVESRAG, r)
			opset(AVESRLB, r)
			opset(AVESRLH, r)
			opset(AVESRLF, r)
			opset(AVESRLG, r)
		case AVLGVG:
			opset(AVLGVB, r)
			opset(AVLGVH, r)
			opset(AVLGVF, r)
		case AVLVGG:
			opset(AVLVGB, r)
			opset(AVLVGH, r)
			opset(AVLVGF, r)
		case AVZERO:
			opset(AVONE, r)
		case AVREPIG:
			opset(AVREPIB, r)
			opset(AVREPIH, r)
			opset(AVREPIF, r)
		case AVLEIG:
			opset(AVLEIB, r)
			opset(AVLEIH, r)
			opset(AVLEIF, r)
		case AVGMG:
			opset(AVGMB, r)
			opset(AVGMH, r)
			opset(AVGMF, r)
		case AVREPG:
			opset(AVREPB, r)
			opset(AVREPH, r)
			opset(AVREPF, r)
		case AVERIMG:
			opset(AVERIMB, r)
			opset(AVERIMH, r)
			opset(AVERIMF, r)
		case AVFTCIDB:
			opset(AWFTCIDB, r)
		case AVLR:
			opset(AVUPHB, r)
			opset(AVUPHH, r)
			opset(AVUPHF, r)
			opset(AVUPLHB, r)
			opset(AVUPLHH, r)
			opset(AVUPLHF, r)
			opset(AVUPLB, r)
			opset(AVUPLHW, r)
			opset(AVUPLF, r)
			opset(AVUPLLB, r)
			opset(AVUPLLH, r)
			opset(AVUPLLF, r)
			opset(AVCLZB, r)
			opset(AVCLZH, r)
			opset(AVCLZF, r)
			opset(AVCLZG, r)
			opset(AVCTZB, r)
			opset(AVCTZH, r)
			opset(AVCTZF, r)
			opset(AVCTZG, r)
			opset(AVLDEB, r)
			opset(AWLDEB, r)
			opset(AVFLCDB, r)
			opset(AWFLCDB, r)
			opset(AVFLNDB, r)
			opset(AWFLNDB, r)
			opset(AVFLPDB, r)
			opset(AWFLPDB, r)
			opset(AVFSQDB, r)
			opset(AWFSQDB, r)
			opset(AVISTRB, r)
			opset(AVISTRH, r)
			opset(AVISTRF, r)
			opset(AVISTRBS, r)
			opset(AVISTRHS, r)
			opset(AVISTRFS, r)
			opset(AVLCB, r)
			opset(AVLCH, r)
			opset(AVLCF, r)
			opset(AVLCG, r)
			opset(AVLPB, r)
			opset(AVLPH, r)
			opset(AVLPF, r)
			opset(AVLPG, r)
			opset(AVPOPCT, r)
			opset(AVSEGB, r)
			opset(AVSEGH, r)
			opset(AVSEGF, r)
		case AVECG:
			opset(AVECB, r)
			opset(AVECH, r)
			opset(AVECF, r)
			opset(AVECLB, r)
			opset(AVECLH, r)
			opset(AVECLF, r)
			opset(AVECLG, r)
			opset(AWFCDB, r)
			opset(AWFKDB, r)
		case AVCEQG:
			opset(AVCEQB, r)
			opset(AVCEQH, r)
			opset(AVCEQF, r)
			opset(AVCEQBS, r)
			opset(AVCEQHS, r)
			opset(AVCEQFS, r)
			opset(AVCEQGS, r)
			opset(AVCHB, r)
			opset(AVCHH, r)
			opset(AVCHF, r)
			opset(AVCHG, r)
			opset(AVCHBS, r)
			opset(AVCHHS, r)
			opset(AVCHFS, r)
			opset(AVCHGS, r)
			opset(AVCHLB, r)
			opset(AVCHLH, r)
			opset(AVCHLF, r)
			opset(AVCHLG, r)
			opset(AVCHLBS, r)
			opset(AVCHLHS, r)
			opset(AVCHLFS, r)
			opset(AVCHLGS, r)
		case AVFAEF:
			opset(AVFAEB, r)
			opset(AVFAEH, r)
			opset(AVFAEBS, r)
			opset(AVFAEHS, r)
			opset(AVFAEFS, r)
			opset(AVFAEZB, r)
			opset(AVFAEZH, r)
			opset(AVFAEZF, r)
			opset(AVFAEZBS, r)
			opset(AVFAEZHS, r)
			opset(AVFAEZFS, r)
			opset(AVFEEB, r)
			opset(AVFEEH, r)
			opset(AVFEEF, r)
			opset(AVFEEBS, r)
			opset(AVFEEHS, r)
			opset(AVFEEFS, r)
			opset(AVFEEZB, r)
			opset(AVFEEZH, r)
			opset(AVFEEZF, r)
			opset(AVFEEZBS, r)
			opset(AVFEEZHS, r)
			opset(AVFEEZFS, r)
			opset(AVFENEB, r)
			opset(AVFENEH, r)
			opset(AVFENEF, r)
			opset(AVFENEBS, r)
			opset(AVFENEHS, r)
			opset(AVFENEFS, r)
			opset(AVFENEZB, r)
			opset(AVFENEZH, r)
			opset(AVFENEZF, r)
			opset(AVFENEZBS, r)
			opset(AVFENEZHS, r)
			opset(AVFENEZFS, r)
		case AVPKSG:
			opset(AVPKSH, r)
			opset(AVPKSF, r)
			opset(AVPKSHS, r)
			opset(AVPKSFS, r)
			opset(AVPKSGS, r)
			opset(AVPKLSH, r)
			opset(AVPKLSF, r)
			opset(AVPKLSG, r)
			opset(AVPKLSHS, r)
			opset(AVPKLSFS, r)
			opset(AVPKLSGS, r)
		case AVAQ:
			opset(AVAB, r)
			opset(AVAH, r)
			opset(AVAF, r)
			opset(AVAG, r)
			opset(AVACCB, r)
			opset(AVACCH, r)
			opset(AVACCF, r)
			opset(AVACCG, r)
			opset(AVACCQ, r)
			opset(AVN, r)
			opset(AVNC, r)
			opset(AVAVGB, r)
			opset(AVAVGH, r)
			opset(AVAVGF, r)
			opset(AVAVGG, r)
			opset(AVAVGLB, r)
			opset(AVAVGLH, r)
			opset(AVAVGLF, r)
			opset(AVAVGLG, r)
			opset(AVCKSM, r)
			opset(AVX, r)
			opset(AVFADB, r)
			opset(AWFADB, r)
			opset(AVFCEDB, r)
			opset(AVFCEDBS, r)
			opset(AWFCEDB, r)
			opset(AWFCEDBS, r)
			opset(AVFCHDB, r)
			opset(AVFCHDBS, r)
			opset(AWFCHDB, r)
			opset(AWFCHDBS, r)
			opset(AVFCHEDB, r)
			opset(AVFCHEDBS, r)
			opset(AWFCHEDB, r)
			opset(AWFCHEDBS, r)
			opset(AVFMDB, r)
			opset(AWFMDB, r)
			opset(AVGFMB, r)
			opset(AVGFMH, r)
			opset(AVGFMF, r)
			opset(AVGFMG, r)
			opset(AVMXB, r)
			opset(AVMXH, r)
			opset(AVMXF, r)
			opset(AVMXG, r)
			opset(AVMXLB, r)
			opset(AVMXLH, r)
			opset(AVMXLF, r)
			opset(AVMXLG, r)
			opset(AVMNB, r)
			opset(AVMNH, r)
			opset(AVMNF, r)
			opset(AVMNG, r)
			opset(AVMNLB, r)
			opset(AVMNLH, r)
			opset(AVMNLF, r)
			opset(AVMNLG, r)
			opset(AVMRHB, r)
			opset(AVMRHH, r)
			opset(AVMRHF, r)
			opset(AVMRHG, r)
			opset(AVMRLB, r)
			opset(AVMRLH, r)
			opset(AVMRLF, r)
			opset(AVMRLG, r)
			opset(AVMEB, r)
			opset(AVMEH, r)
			opset(AVMEF, r)
			opset(AVMLEB, r)
			opset(AVMLEH, r)
			opset(AVMLEF, r)
			opset(AVMOB, r)
			opset(AVMOH, r)
			opset(AVMOF, r)
			opset(AVMLOB, r)
			opset(AVMLOH, r)
			opset(AVMLOF, r)
			opset(AVMHB, r)
			opset(AVMHH, r)
			opset(AVMHF, r)
			opset(AVMLHB, r)
			opset(AVMLHH, r)
			opset(AVMLHF, r)
			opset(AVMLH, r)
			opset(AVMLHW, r)
			opset(AVMLF, r)
			opset(AVNO, r)
			opset(AVO, r)
			opset(AVPKH, r)
			opset(AVPKF, r)
			opset(AVPKG, r)
			opset(AVSUMGH, r)
			opset(AVSUMGF, r)
			opset(AVSUMQF, r)
			opset(AVSUMQG, r)
			opset(AVSUMB, r)
			opset(AVSUMH, r)
		case AVERLLVG:
			opset(AVERLLVB, r)
			opset(AVERLLVH, r)
			opset(AVERLLVF, r)
			opset(AVESLVB, r)
			opset(AVESLVH, r)
			opset(AVESLVF, r)
			opset(AVESLVG, r)
			opset(AVESRAVB, r)
			opset(AVESRAVH, r)
			opset(AVESRAVF, r)
			opset(AVESRAVG, r)
			opset(AVESRLVB, r)
			opset(AVESRLVH, r)
			opset(AVESRLVF, r)
			opset(AVESRLVG, r)
			opset(AVFDDB, r)
			opset(AWFDDB, r)
			opset(AVFSDB, r)
			opset(AWFSDB, r)
			opset(AVSL, r)
			opset(AVSLB, r)
			opset(AVSRA, r)
			opset(AVSRAB, r)
			opset(AVSRL, r)
			opset(AVSRLB, r)
			opset(AVSB, r)
			opset(AVSH, r)
			opset(AVSF, r)
			opset(AVSG, r)
			opset(AVSQ, r)
			opset(AVSCBIB, r)
			opset(AVSCBIH, r)
			opset(AVSCBIF, r)
			opset(AVSCBIG, r)
			opset(AVSCBIQ, r)
		case AVACQ:
			opset(AVACCCQ, r)
			opset(AVGFMAB, r)
			opset(AVGFMAH, r)
			opset(AVGFMAF, r)
			opset(AVGFMAG, r)
			opset(AVMALB, r)
			opset(AVMALHW, r)
			opset(AVMALF, r)
			opset(AVMAHB, r)
			opset(AVMAHH, r)
			opset(AVMAHF, r)
			opset(AVMALHB, r)
			opset(AVMALHH, r)
			opset(AVMALHF, r)
			opset(AVMAEB, r)
			opset(AVMAEH, r)
			opset(AVMAEF, r)
			opset(AVMALEB, r)
			opset(AVMALEH, r)
			opset(AVMALEF, r)
			opset(AVMAOB, r)
			opset(AVMAOH, r)
			opset(AVMAOF, r)
			opset(AVMALOB, r)
			opset(AVMALOH, r)
			opset(AVMALOF, r)
			opset(AVSTRCB, r)
			opset(AVSTRCH, r)
			opset(AVSTRCF, r)
			opset(AVSTRCBS, r)
			opset(AVSTRCHS, r)
			opset(AVSTRCFS, r)
			opset(AVSTRCZB, r)
			opset(AVSTRCZH, r)
			opset(AVSTRCZF, r)
			opset(AVSTRCZBS, r)
			opset(AVSTRCZHS, r)
			opset(AVSTRCZFS, r)
			opset(AVSBCBIQ, r)
			opset(AVSBIQ, r)
			opset(AVMSLG, r)
			opset(AVMSLEG, r)
			opset(AVMSLOG, r)
			opset(AVMSLEOG, r)
		case AVSEL:
			opset(AVFMADB, r)
			opset(AWFMADB, r)
			opset(AVFMSDB, r)
			opset(AWFMSDB, r)
			opset(AVPERM, r)
		}
	}
}

const (
	op_A       uint32 = 0x5A00 
	op_AD      uint32 = 0x6A00 
	op_ADB     uint32 = 0xED1A 
	op_ADBR    uint32 = 0xB31A 
	op_ADR     uint32 = 0x2A00 
	op_ADTR    uint32 = 0xB3D2 
	op_ADTRA   uint32 = 0xB3D2 
	op_AE      uint32 = 0x7A00 
	op_AEB     uint32 = 0xED0A 
	op_AEBR    uint32 = 0xB30A 
	op_AER     uint32 = 0x3A00 
	op_AFI     uint32 = 0xC209 
	op_AG      uint32 = 0xE308 
	op_AGF     uint32 = 0xE318 
	op_AGFI    uint32 = 0xC208 
	op_AGFR    uint32 = 0xB918 
	op_AGHI    uint32 = 0xA70B 
	op_AGHIK   uint32 = 0xECD9 
	op_AGR     uint32 = 0xB908 
	op_AGRK    uint32 = 0xB9E8 
	op_AGSI    uint32 = 0xEB7A 
	op_AH      uint32 = 0x4A00 
	op_AHHHR   uint32 = 0xB9C8 
	op_AHHLR   uint32 = 0xB9D8 
	op_AHI     uint32 = 0xA70A 
	op_AHIK    uint32 = 0xECD8 
	op_AHY     uint32 = 0xE37A 
	op_AIH     uint32 = 0xCC08 
	op_AL      uint32 = 0x5E00 
	op_ALC     uint32 = 0xE398 
	op_ALCG    uint32 = 0xE388 
	op_ALCGR   uint32 = 0xB988 
	op_ALCR    uint32 = 0xB998 
	op_ALFI    uint32 = 0xC20B 
	op_ALG     uint32 = 0xE30A 
	op_ALGF    uint32 = 0xE31A 
	op_ALGFI   uint32 = 0xC20A 
	op_ALGFR   uint32 = 0xB91A 
	op_ALGHSIK uint32 = 0xECDB 
	op_ALGR    uint32 = 0xB90A 
	op_ALGRK   uint32 = 0xB9EA 
	op_ALGSI   uint32 = 0xEB7E 
	op_ALHHHR  uint32 = 0xB9CA 
	op_ALHHLR  uint32 = 0xB9DA 
	op_ALHSIK  uint32 = 0xECDA 
	op_ALR     uint32 = 0x1E00 
	op_ALRK    uint32 = 0xB9FA 
	op_ALSI    uint32 = 0xEB6E 
	op_ALSIH   uint32 = 0xCC0A 
	op_ALSIHN  uint32 = 0xCC0B 
	op_ALY     uint32 = 0xE35E 
	op_AP      uint32 = 0xFA00 
	op_AR      uint32 = 0x1A00 
	op_ARK     uint32 = 0xB9F8 
	op_ASI     uint32 = 0xEB6A 
	op_AU      uint32 = 0x7E00 
	op_AUR     uint32 = 0x3E00 
	op_AW      uint32 = 0x6E00 
	op_AWR     uint32 = 0x2E00 
	op_AXBR    uint32 = 0xB34A 
	op_AXR     uint32 = 0x3600 
	op_AXTR    uint32 = 0xB3DA 
	op_AXTRA   uint32 = 0xB3DA 
	op_AY      uint32 = 0xE35A 
	op_BAKR    uint32 = 0xB240 
	op_BAL     uint32 = 0x4500 
	op_BALR    uint32 = 0x0500 
	op_BAS     uint32 = 0x4D00 
	op_BASR    uint32 = 0x0D00 
	op_BASSM   uint32 = 0x0C00 
	op_BC      uint32 = 0x4700 
	op_BCR     uint32 = 0x0700 
	op_BCT     uint32 = 0x4600 
	op_BCTG    uint32 = 0xE346 
	op_BCTGR   uint32 = 0xB946 
	op_BCTR    uint32 = 0x0600 
	op_BPP     uint32 = 0xC700 
	op_BPRP    uint32 = 0xC500 
	op_BRAS    uint32 = 0xA705 
	op_BRASL   uint32 = 0xC005 
	op_BRC     uint32 = 0xA704 
	op_BRCL    uint32 = 0xC004 
	op_BRCT    uint32 = 0xA706 
	op_BRCTG   uint32 = 0xA707 
	op_BRCTH   uint32 = 0xCC06 
	op_BRXH    uint32 = 0x8400 
	op_BRXHG   uint32 = 0xEC44 
	op_BRXLE   uint32 = 0x8500 
	op_BRXLG   uint32 = 0xEC45 
	op_BSA     uint32 = 0xB25A 
	op_BSG     uint32 = 0xB258 
	op_BSM     uint32 = 0x0B00 
	op_BXH     uint32 = 0x8600 
	op_BXHG    uint32 = 0xEB44 
	op_BXLE    uint32 = 0x8700 
	op_BXLEG   uint32 = 0xEB45 
	op_C       uint32 = 0x5900 
	op_CD      uint32 = 0x6900 
	op_CDB     uint32 = 0xED19 
	op_CDBR    uint32 = 0xB319 
	op_CDFBR   uint32 = 0xB395 
	op_CDFBRA  uint32 = 0xB395 
	op_CDFR    uint32 = 0xB3B5 
	op_CDFTR   uint32 = 0xB951 
	op_CDGBR   uint32 = 0xB3A5 
	op_CDGBRA  uint32 = 0xB3A5 
	op_CDGR    uint32 = 0xB3C5 
	op_CDGTR   uint32 = 0xB3F1 
	op_CDGTRA  uint32 = 0xB3F1 
	op_CDLFBR  uint32 = 0xB391 
	op_CDLFTR  uint32 = 0xB953 
	op_CDLGBR  uint32 = 0xB3A1 
	op_CDLGTR  uint32 = 0xB952 
	op_CDR     uint32 = 0x2900 
	op_CDS     uint32 = 0xBB00 
	op_CDSG    uint32 = 0xEB3E 
	op_CDSTR   uint32 = 0xB3F3 
	op_CDSY    uint32 = 0xEB31 
	op_CDTR    uint32 = 0xB3E4 
	op_CDUTR   uint32 = 0xB3F2 
	op_CDZT    uint32 = 0xEDAA 
	op_CE      uint32 = 0x7900 
	op_CEB     uint32 = 0xED09 
	op_CEBR    uint32 = 0xB309 
	op_CEDTR   uint32 = 0xB3F4 
	op_CEFBR   uint32 = 0xB394 
	op_CEFBRA  uint32 = 0xB394 
	op_CEFR    uint32 = 0xB3B4 
	op_CEGBR   uint32 = 0xB3A4 
	op_CEGBRA  uint32 = 0xB3A4 
	op_CEGR    uint32 = 0xB3C4 
	op_CELFBR  uint32 = 0xB390 
	op_CELGBR  uint32 = 0xB3A0 
	op_CER     uint32 = 0x3900 
	op_CEXTR   uint32 = 0xB3FC 
	op_CFC     uint32 = 0xB21A 
	op_CFDBR   uint32 = 0xB399 
	op_CFDBRA  uint32 = 0xB399 
	op_CFDR    uint32 = 0xB3B9 
	op_CFDTR   uint32 = 0xB941 
	op_CFEBR   uint32 = 0xB398 
	op_CFEBRA  uint32 = 0xB398 
	op_CFER    uint32 = 0xB3B8 
	op_CFI     uint32 = 0xC20D 
	op_CFXBR   uint32 = 0xB39A 
	op_CFXBRA  uint32 = 0xB39A 
	op_CFXR    uint32 = 0xB3BA 
	op_CFXTR   uint32 = 0xB949 
	op_CG      uint32 = 0xE320 
	op_CGDBR   uint32 = 0xB3A9 
	op_CGDBRA  uint32 = 0xB3A9 
	op_CGDR    uint32 = 0xB3C9 
	op_CGDTR   uint32 = 0xB3E1 
	op_CGDTRA  uint32 = 0xB3E1 
	op_CGEBR   uint32 = 0xB3A8 
	op_CGEBRA  uint32 = 0xB3A8 
	op_CGER    uint32 = 0xB3C8 
	op_CGF     uint32 = 0xE330 
	op_CGFI    uint32 = 0xC20C 
	op_CGFR    uint32 = 0xB930 
	op_CGFRL   uint32 = 0xC60C 
	op_CGH     uint32 = 0xE334 
	op_CGHI    uint32 = 0xA70F 
	op_CGHRL   uint32 = 0xC604 
	op_CGHSI   uint32 = 0xE558 
	op_CGIB    uint32 = 0xECFC 
	op_CGIJ    uint32 = 0xEC7C 
	op_CGIT    uint32 = 0xEC70 
	op_CGR     uint32 = 0xB920 
	op_CGRB    uint32 = 0xECE4 
	op_CGRJ    uint32 = 0xEC64 
	op_CGRL    uint32 = 0xC608 
	op_CGRT    uint32 = 0xB960 
	op_CGXBR   uint32 = 0xB3AA 
	op_CGXBRA  uint32 = 0xB3AA 
	op_CGXR    uint32 = 0xB3CA 
	op_CGXTR   uint32 = 0xB3E9 
	op_CGXTRA  uint32 = 0xB3E9 
	op_CH      uint32 = 0x4900 
	op_CHF     uint32 = 0xE3CD 
	op_CHHR    uint32 = 0xB9CD 
	op_CHHSI   uint32 = 0xE554 
	op_CHI     uint32 = 0xA70E 
	op_CHLR    uint32 = 0xB9DD 
	op_CHRL    uint32 = 0xC605 
	op_CHSI    uint32 = 0xE55C 
	op_CHY     uint32 = 0xE379 
	op_CIB     uint32 = 0xECFE 
	op_CIH     uint32 = 0xCC0D 
	op_CIJ     uint32 = 0xEC7E 
	op_CIT     uint32 = 0xEC72 
	op_CKSM    uint32 = 0xB241 
	op_CL      uint32 = 0x5500 
	op_CLC     uint32 = 0xD500 
	op_CLCL    uint32 = 0x0F00 
	op_CLCLE   uint32 = 0xA900 
	op_CLCLU   uint32 = 0xEB8F 
	op_CLFDBR  uint32 = 0xB39D 
	op_CLFDTR  uint32 = 0xB943 
	op_CLFEBR  uint32 = 0xB39C 
	op_CLFHSI  uint32 = 0xE55D 
	op_CLFI    uint32 = 0xC20F 
	op_CLFIT   uint32 = 0xEC73 
	op_CLFXBR  uint32 = 0xB39E 
	op_CLFXTR  uint32 = 0xB94B 
	op_CLG     uint32 = 0xE321 
	op_CLGDBR  uint32 = 0xB3AD 
	op_CLGDTR  uint32 = 0xB942 
	op_CLGEBR  uint32 = 0xB3AC 
	op_CLGF    uint32 = 0xE331 
	op_CLGFI   uint32 = 0xC20E 
	op_CLGFR   uint32 = 0xB931 
	op_CLGFRL  uint32 = 0xC60E 
	op_CLGHRL  uint32 = 0xC606 
	op_CLGHSI  uint32 = 0xE559 
	op_CLGIB   uint32 = 0xECFD 
	op_CLGIJ   uint32 = 0xEC7D 
	op_CLGIT   uint32 = 0xEC71 
	op_CLGR    uint32 = 0xB921 
	op_CLGRB   uint32 = 0xECE5 
	op_CLGRJ   uint32 = 0xEC65 
	op_CLGRL   uint32 = 0xC60A 
	op_CLGRT   uint32 = 0xB961 
	op_CLGT    uint32 = 0xEB2B 
	op_CLGXBR  uint32 = 0xB3AE 
	op_CLGXTR  uint32 = 0xB94A 
	op_CLHF    uint32 = 0xE3CF 
	op_CLHHR   uint32 = 0xB9CF 
	op_CLHHSI  uint32 = 0xE555 
	op_CLHLR   uint32 = 0xB9DF 
	op_CLHRL   uint32 = 0xC607 
	op_CLI     uint32 = 0x9500 
	op_CLIB    uint32 = 0xECFF 
	op_CLIH    uint32 = 0xCC0F 
	op_CLIJ    uint32 = 0xEC7F 
	op_CLIY    uint32 = 0xEB55 
	op_CLM     uint32 = 0xBD00 
	op_CLMH    uint32 = 0xEB20 
	op_CLMY    uint32 = 0xEB21 
	op_CLR     uint32 = 0x1500 
	op_CLRB    uint32 = 0xECF7 
	op_CLRJ    uint32 = 0xEC77 
	op_CLRL    uint32 = 0xC60F 
	op_CLRT    uint32 = 0xB973 
	op_CLST    uint32 = 0xB25D 
	op_CLT     uint32 = 0xEB23 
	op_CLY     uint32 = 0xE355 
	op_CMPSC   uint32 = 0xB263 
	op_CP      uint32 = 0xF900 
	op_CPSDR   uint32 = 0xB372 
	op_CPYA    uint32 = 0xB24D 
	op_CR      uint32 = 0x1900 
	op_CRB     uint32 = 0xECF6 
	op_CRDTE   uint32 = 0xB98F 
	op_CRJ     uint32 = 0xEC76 
	op_CRL     uint32 = 0xC60D 
	op_CRT     uint32 = 0xB972 
	op_CS      uint32 = 0xBA00 
	op_CSCH    uint32 = 0xB230 
	op_CSDTR   uint32 = 0xB3E3 
	op_CSG     uint32 = 0xEB30 
	op_CSP     uint32 = 0xB250 
	op_CSPG    uint32 = 0xB98A 
	op_CSST    uint32 = 0xC802 
	op_CSXTR   uint32 = 0xB3EB 
	op_CSY     uint32 = 0xEB14 
	op_CU12    uint32 = 0xB2A7 
	op_CU14    uint32 = 0xB9B0 
	op_CU21    uint32 = 0xB2A6 
	op_CU24    uint32 = 0xB9B1 
	op_CU41    uint32 = 0xB9B2 
	op_CU42    uint32 = 0xB9B3 
	op_CUDTR   uint32 = 0xB3E2 
	op_CUSE    uint32 = 0xB257 
	op_CUTFU   uint32 = 0xB2A7 
	op_CUUTF   uint32 = 0xB2A6 
	op_CUXTR   uint32 = 0xB3EA 
	op_CVB     uint32 = 0x4F00 
	op_CVBG    uint32 = 0xE30E 
	op_CVBY    uint32 = 0xE306 
	op_CVD     uint32 = 0x4E00 
	op_CVDG    uint32 = 0xE32E 
	op_CVDY    uint32 = 0xE326 
	op_CXBR    uint32 = 0xB349 
	op_CXFBR   uint32 = 0xB396 
	op_CXFBRA  uint32 = 0xB396 
	op_CXFR    uint32 = 0xB3B6 
	op_CXFTR   uint32 = 0xB959 
	op_CXGBR   uint32 = 0xB3A6 
	op_CXGBRA  uint32 = 0xB3A6 
	op_CXGR    uint32 = 0xB3C6 
	op_CXGTR   uint32 = 0xB3F9 
	op_CXGTRA  uint32 = 0xB3F9 
	op_CXLFBR  uint32 = 0xB392 
	op_CXLFTR  uint32 = 0xB95B 
	op_CXLGBR  uint32 = 0xB3A2 
	op_CXLGTR  uint32 = 0xB95A 
	op_CXR     uint32 = 0xB369 
	op_CXSTR   uint32 = 0xB3FB 
	op_CXTR    uint32 = 0xB3EC 
	op_CXUTR   uint32 = 0xB3FA 
	op_CXZT    uint32 = 0xEDAB 
	op_CY      uint32 = 0xE359 
	op_CZDT    uint32 = 0xEDA8 
	op_CZXT    uint32 = 0xEDA9 
	op_D       uint32 = 0x5D00 
	op_DD      uint32 = 0x6D00 
	op_DDB     uint32 = 0xED1D 
	op_DDBR    uint32 = 0xB31D 
	op_DDR     uint32 = 0x2D00 
	op_DDTR    uint32 = 0xB3D1 
	op_DDTRA   uint32 = 0xB3D1 
	op_DE      uint32 = 0x7D00 
	op_DEB     uint32 = 0xED0D 
	op_DEBR    uint32 = 0xB30D 
	op_DER     uint32 = 0x3D00 
	op_DIDBR   uint32 = 0xB35B 
	op_DIEBR   uint32 = 0xB353 
	op_DL      uint32 = 0xE397 
	op_DLG     uint32 = 0xE387 
	op_DLGR    uint32 = 0xB987 
	op_DLR     uint32 = 0xB997 
	op_DP      uint32 = 0xFD00 
	op_DR      uint32 = 0x1D00 
	op_DSG     uint32 = 0xE30D 
	op_DSGF    uint32 = 0xE31D 
	op_DSGFR   uint32 = 0xB91D 
	op_DSGR    uint32 = 0xB90D 
	op_DXBR    uint32 = 0xB34D 
	op_DXR     uint32 = 0xB22D 
	op_DXTR    uint32 = 0xB3D9 
	op_DXTRA   uint32 = 0xB3D9 
	op_EAR     uint32 = 0xB24F 
	op_ECAG    uint32 = 0xEB4C 
	op_ECTG    uint32 = 0xC801 
	op_ED      uint32 = 0xDE00 
	op_EDMK    uint32 = 0xDF00 
	op_EEDTR   uint32 = 0xB3E5 
	op_EEXTR   uint32 = 0xB3ED 
	op_EFPC    uint32 = 0xB38C 
	op_EPAIR   uint32 = 0xB99A 
	op_EPAR    uint32 = 0xB226 
	op_EPSW    uint32 = 0xB98D 
	op_EREG    uint32 = 0xB249 
	op_EREGG   uint32 = 0xB90E 
	op_ESAIR   uint32 = 0xB99B 
	op_ESAR    uint32 = 0xB227 
	op_ESDTR   uint32 = 0xB3E7 
	op_ESEA    uint32 = 0xB99D 
	op_ESTA    uint32 = 0xB24A 
	op_ESXTR   uint32 = 0xB3EF 
	op_ETND    uint32 = 0xB2EC 
	op_EX      uint32 = 0x4400 
	op_EXRL    uint32 = 0xC600 
	op_FIDBR   uint32 = 0xB35F 
	op_FIDBRA  uint32 = 0xB35F 
	op_FIDR    uint32 = 0xB37F 
	op_FIDTR   uint32 = 0xB3D7 
	op_FIEBR   uint32 = 0xB357 
	op_FIEBRA  uint32 = 0xB357 
	op_FIER    uint32 = 0xB377 
	op_FIXBR   uint32 = 0xB347 
	op_FIXBRA  uint32 = 0xB347 
	op_FIXR    uint32 = 0xB367 
	op_FIXTR   uint32 = 0xB3DF 
	op_FLOGR   uint32 = 0xB983 
	op_HDR     uint32 = 0x2400 
	op_HER     uint32 = 0x3400 
	op_HSCH    uint32 = 0xB231 
	op_IAC     uint32 = 0xB224 
	op_IC      uint32 = 0x4300 
	op_ICM     uint32 = 0xBF00 
	op_ICMH    uint32 = 0xEB80 
	op_ICMY    uint32 = 0xEB81 
	op_ICY     uint32 = 0xE373 
	op_IDTE    uint32 = 0xB98E 
	op_IEDTR   uint32 = 0xB3F6 
	op_IEXTR   uint32 = 0xB3FE 
	op_IIHF    uint32 = 0xC008 
	op_IIHH    uint32 = 0xA500 
	op_IIHL    uint32 = 0xA501 
	op_IILF    uint32 = 0xC009 
	op_IILH    uint32 = 0xA502 
	op_IILL    uint32 = 0xA503 
	op_IPK     uint32 = 0xB20B 
	op_IPM     uint32 = 0xB222 
	op_IPTE    uint32 = 0xB221 
	op_ISKE    uint32 = 0xB229 
	op_IVSK    uint32 = 0xB223 
	op_KDB     uint32 = 0xED18 
	op_KDBR    uint32 = 0xB318 
	op_KDTR    uint32 = 0xB3E0 
	op_KEB     uint32 = 0xED08 
	op_KEBR    uint32 = 0xB308 
	op_KIMD    uint32 = 0xB93E 
	op_KLMD    uint32 = 0xB93F 
	op_KM      uint32 = 0xB92E 
	op_KMAC    uint32 = 0xB91E 
	op_KMC     uint32 = 0xB92F 
	op_KMCTR   uint32 = 0xB92D 
	op_KMF     uint32 = 0xB92A 
	op_KMO     uint32 = 0xB92B 
	op_KXBR    uint32 = 0xB348 
	op_KXTR    uint32 = 0xB3E8 
	op_L       uint32 = 0x5800 
	op_LA      uint32 = 0x4100 
	op_LAA     uint32 = 0xEBF8 
	op_LAAG    uint32 = 0xEBE8 
	op_LAAL    uint32 = 0xEBFA 
	op_LAALG   uint32 = 0xEBEA 
	op_LAE     uint32 = 0x5100 
	op_LAEY    uint32 = 0xE375 
	op_LAM     uint32 = 0x9A00 
	op_LAMY    uint32 = 0xEB9A 
	op_LAN     uint32 = 0xEBF4 
	op_LANG    uint32 = 0xEBE4 
	op_LAO     uint32 = 0xEBF6 
	op_LAOG    uint32 = 0xEBE6 
	op_LARL    uint32 = 0xC000 
	op_LASP    uint32 = 0xE500 
	op_LAT     uint32 = 0xE39F 
	op_LAX     uint32 = 0xEBF7 
	op_LAXG    uint32 = 0xEBE7 
	op_LAY     uint32 = 0xE371 
	op_LB      uint32 = 0xE376 
	op_LBH     uint32 = 0xE3C0 
	op_LBR     uint32 = 0xB926 
	op_LCDBR   uint32 = 0xB313 
	op_LCDFR   uint32 = 0xB373 
	op_LCDR    uint32 = 0x2300 
	op_LCEBR   uint32 = 0xB303 
	op_LCER    uint32 = 0x3300 
	op_LCGFR   uint32 = 0xB913 
	op_LCGR    uint32 = 0xB903 
	op_LCR     uint32 = 0x1300 
	op_LCTL    uint32 = 0xB700 
	op_LCTLG   uint32 = 0xEB2F 
	op_LCXBR   uint32 = 0xB343 
	op_LCXR    uint32 = 0xB363 
	op_LD      uint32 = 0x6800 
	op_LDE     uint32 = 0xED24 
	op_LDEB    uint32 = 0xED04 
	op_LDEBR   uint32 = 0xB304 
	op_LDER    uint32 = 0xB324 
	op_LDETR   uint32 = 0xB3D4 
	op_LDGR    uint32 = 0xB3C1 
	op_LDR     uint32 = 0x2800 
	op_LDXBR   uint32 = 0xB345 
	op_LDXBRA  uint32 = 0xB345 
	op_LDXR    uint32 = 0x2500 
	op_LDXTR   uint32 = 0xB3DD 
	op_LDY     uint32 = 0xED65 
	op_LE      uint32 = 0x7800 
	op_LEDBR   uint32 = 0xB344 
	op_LEDBRA  uint32 = 0xB344 
	op_LEDR    uint32 = 0x3500 
	op_LEDTR   uint32 = 0xB3D5 
	op_LER     uint32 = 0x3800 
	op_LEXBR   uint32 = 0xB346 
	op_LEXBRA  uint32 = 0xB346 
	op_LEXR    uint32 = 0xB366 
	op_LEY     uint32 = 0xED64 
	op_LFAS    uint32 = 0xB2BD 
	op_LFH     uint32 = 0xE3CA 
	op_LFHAT   uint32 = 0xE3C8 
	op_LFPC    uint32 = 0xB29D 
	op_LG      uint32 = 0xE304 
	op_LGAT    uint32 = 0xE385 
	op_LGB     uint32 = 0xE377 
	op_LGBR    uint32 = 0xB906 
	op_LGDR    uint32 = 0xB3CD 
	op_LGF     uint32 = 0xE314 
	op_LGFI    uint32 = 0xC001 
	op_LGFR    uint32 = 0xB914 
	op_LGFRL   uint32 = 0xC40C 
	op_LGH     uint32 = 0xE315 
	op_LGHI    uint32 = 0xA709 
	op_LGHR    uint32 = 0xB907 
	op_LGHRL   uint32 = 0xC404 
	op_LGR     uint32 = 0xB904 
	op_LGRL    uint32 = 0xC408 
	op_LH      uint32 = 0x4800 
	op_LHH     uint32 = 0xE3C4 
	op_LHI     uint32 = 0xA708 
	op_LHR     uint32 = 0xB927 
	op_LHRL    uint32 = 0xC405 
	op_LHY     uint32 = 0xE378 
	op_LLC     uint32 = 0xE394 
	op_LLCH    uint32 = 0xE3C2 
	op_LLCR    uint32 = 0xB994 
	op_LLGC    uint32 = 0xE390 
	op_LLGCR   uint32 = 0xB984 
	op_LLGF    uint32 = 0xE316 
	op_LLGFAT  uint32 = 0xE39D 
	op_LLGFR   uint32 = 0xB916 
	op_LLGFRL  uint32 = 0xC40E 
	op_LLGH    uint32 = 0xE391 
	op_LLGHR   uint32 = 0xB985 
	op_LLGHRL  uint32 = 0xC406 
	op_LLGT    uint32 = 0xE317 
	op_LLGTAT  uint32 = 0xE39C 
	op_LLGTR   uint32 = 0xB917 
	op_LLH     uint32 = 0xE395 
	op_LLHH    uint32 = 0xE3C6 
	op_LLHR    uint32 = 0xB995 
	op_LLHRL   uint32 = 0xC402 
	op_LLIHF   uint32 = 0xC00E 
	op_LLIHH   uint32 = 0xA50C 
	op_LLIHL   uint32 = 0xA50D 
	op_LLILF   uint32 = 0xC00F 
	op_LLILH   uint32 = 0xA50E 
	op_LLILL   uint32 = 0xA50F 
	op_LM      uint32 = 0x9800 
	op_LMD     uint32 = 0xEF00 
	op_LMG     uint32 = 0xEB04 
	op_LMH     uint32 = 0xEB96 
	op_LMY     uint32 = 0xEB98 
	op_LNDBR   uint32 = 0xB311 
	op_LNDFR   uint32 = 0xB371 
	op_LNDR    uint32 = 0x2100 
	op_LNEBR   uint32 = 0xB301 
	op_LNER    uint32 = 0x3100 
	op_LNGFR   uint32 = 0xB911 
	op_LNGR    uint32 = 0xB901 
	op_LNR     uint32 = 0x1100 
	op_LNXBR   uint32 = 0xB341 
	op_LNXR    uint32 = 0xB361 
	op_LOC     uint32 = 0xEBF2 
	op_LOCG    uint32 = 0xEBE2 
	op_LOCGR   uint32 = 0xB9E2 
	op_LOCR    uint32 = 0xB9F2 
	op_LPD     uint32 = 0xC804 
	op_LPDBR   uint32 = 0xB310 
	op_LPDFR   uint32 = 0xB370 
	op_LPDG    uint32 = 0xC805 
	op_LPDR    uint32 = 0x2000 
	op_LPEBR   uint32 = 0xB300 
	op_LPER    uint32 = 0x3000 
	op_LPGFR   uint32 = 0xB910 
	op_LPGR    uint32 = 0xB900 
	op_LPQ     uint32 = 0xE38F 
	op_LPR     uint32 = 0x1000 
	op_LPSW    uint32 = 0x8200 
	op_LPSWE   uint32 = 0xB2B2 
	op_LPTEA   uint32 = 0xB9AA 
	op_LPXBR   uint32 = 0xB340 
	op_LPXR    uint32 = 0xB360 
	op_LR      uint32 = 0x1800 
	op_LRA     uint32 = 0xB100 
	op_LRAG    uint32 = 0xE303 
	op_LRAY    uint32 = 0xE313 
	op_LRDR    uint32 = 0x2500 
	op_LRER    uint32 = 0x3500 
	op_LRL     uint32 = 0xC40D 
	op_LRV     uint32 = 0xE31E 
	op_LRVG    uint32 = 0xE30F 
	op_LRVGR   uint32 = 0xB90F 
	op_LRVH    uint32 = 0xE31F 
	op_LRVR    uint32 = 0xB91F 
	op_LT      uint32 = 0xE312 
	op_LTDBR   uint32 = 0xB312 
	op_LTDR    uint32 = 0x2200 
	op_LTDTR   uint32 = 0xB3D6 
	op_LTEBR   uint32 = 0xB302 
	op_LTER    uint32 = 0x3200 
	op_LTG     uint32 = 0xE302 
	op_LTGF    uint32 = 0xE332 
	op_LTGFR   uint32 = 0xB912 
	op_LTGR    uint32 = 0xB902 
	op_LTR     uint32 = 0x1200 
	op_LTXBR   uint32 = 0xB342 
	op_LTXR    uint32 = 0xB362 
	op_LTXTR   uint32 = 0xB3DE 
	op_LURA    uint32 = 0xB24B 
	op_LURAG   uint32 = 0xB905 
	op_LXD     uint32 = 0xED25 
	op_LXDB    uint32 = 0xED05 
	op_LXDBR   uint32 = 0xB305 
	op_LXDR    uint32 = 0xB325 
	op_LXDTR   uint32 = 0xB3DC 
	op_LXE     uint32 = 0xED26 
	op_LXEB    uint32 = 0xED06 
	op_LXEBR   uint32 = 0xB306 
	op_LXER    uint32 = 0xB326 
	op_LXR     uint32 = 0xB365 
	op_LY      uint32 = 0xE358 
	op_LZDR    uint32 = 0xB375 
	op_LZER    uint32 = 0xB374 
	op_LZXR    uint32 = 0xB376 
	op_M       uint32 = 0x5C00 
	op_MAD     uint32 = 0xED3E 
	op_MADB    uint32 = 0xED1E 
	op_MADBR   uint32 = 0xB31E 
	op_MADR    uint32 = 0xB33E 
	op_MAE     uint32 = 0xED2E 
	op_MAEB    uint32 = 0xED0E 
	op_MAEBR   uint32 = 0xB30E 
	op_MAER    uint32 = 0xB32E 
	op_MAY     uint32 = 0xED3A 
	op_MAYH    uint32 = 0xED3C 
	op_MAYHR   uint32 = 0xB33C 
	op_MAYL    uint32 = 0xED38 
	op_MAYLR   uint32 = 0xB338 
	op_MAYR    uint32 = 0xB33A 
	op_MC      uint32 = 0xAF00 
	op_MD      uint32 = 0x6C00 
	op_MDB     uint32 = 0xED1C 
	op_MDBR    uint32 = 0xB31C 
	op_MDE     uint32 = 0x7C00 
	op_MDEB    uint32 = 0xED0C 
	op_MDEBR   uint32 = 0xB30C 
	op_MDER    uint32 = 0x3C00 
	op_MDR     uint32 = 0x2C00 
	op_MDTR    uint32 = 0xB3D0 
	op_MDTRA   uint32 = 0xB3D0 
	op_ME      uint32 = 0x7C00 
	op_MEE     uint32 = 0xED37 
	op_MEEB    uint32 = 0xED17 
	op_MEEBR   uint32 = 0xB317 
	op_MEER    uint32 = 0xB337 
	op_MER     uint32 = 0x3C00 
	op_MFY     uint32 = 0xE35C 
	op_MGHI    uint32 = 0xA70D 
	op_MH      uint32 = 0x4C00 
	op_MHI     uint32 = 0xA70C 
	op_MHY     uint32 = 0xE37C 
	op_ML      uint32 = 0xE396 
	op_MLG     uint32 = 0xE386 
	op_MLGR    uint32 = 0xB986 
	op_MLR     uint32 = 0xB996 
	op_MP      uint32 = 0xFC00 
	op_MR      uint32 = 0x1C00 
	op_MS      uint32 = 0x7100 
	op_MSCH    uint32 = 0xB232 
	op_MSD     uint32 = 0xED3F 
	op_MSDB    uint32 = 0xED1F 
	op_MSDBR   uint32 = 0xB31F 
	op_MSDR    uint32 = 0xB33F 
	op_MSE     uint32 = 0xED2F 
	op_MSEB    uint32 = 0xED0F 
	op_MSEBR   uint32 = 0xB30F 
	op_MSER    uint32 = 0xB32F 
	op_MSFI    uint32 = 0xC201 
	op_MSG     uint32 = 0xE30C 
	op_MSGF    uint32 = 0xE31C 
	op_MSGFI   uint32 = 0xC200 
	op_MSGFR   uint32 = 0xB91C 
	op_MSGR    uint32 = 0xB90C 
	op_MSR     uint32 = 0xB252 
	op_MSTA    uint32 = 0xB247 
	op_MSY     uint32 = 0xE351 
	op_MVC     uint32 = 0xD200 
	op_MVCDK   uint32 = 0xE50F 
	op_MVCIN   uint32 = 0xE800 
	op_MVCK    uint32 = 0xD900 
	op_MVCL    uint32 = 0x0E00 
	op_MVCLE   uint32 = 0xA800 
	op_MVCLU   uint32 = 0xEB8E 
	op_MVCOS   uint32 = 0xC800 
	op_MVCP    uint32 = 0xDA00 
	op_MVCS    uint32 = 0xDB00 
	op_MVCSK   uint32 = 0xE50E 
	op_MVGHI   uint32 = 0xE548 
	op_MVHHI   uint32 = 0xE544 
	op_MVHI    uint32 = 0xE54C 
	op_MVI     uint32 = 0x9200 
	op_MVIY    uint32 = 0xEB52 
	op_MVN     uint32 = 0xD100 
	op_MVO     uint32 = 0xF100 
	op_MVPG    uint32 = 0xB254 
	op_MVST    uint32 = 0xB255 
	op_MVZ     uint32 = 0xD300 
	op_MXBR    uint32 = 0xB34C 
	op_MXD     uint32 = 0x6700 
	op_MXDB    uint32 = 0xED07 
	op_MXDBR   uint32 = 0xB307 
	op_MXDR    uint32 = 0x2700 
	op_MXR     uint32 = 0x2600 
	op_MXTR    uint32 = 0xB3D8 
	op_MXTRA   uint32 = 0xB3D8 
	op_MY      uint32 = 0xED3B 
	op_MYH     uint32 = 0xED3D 
	op_MYHR    uint32 = 0xB33D 
	op_MYL     uint32 = 0xED39 
	op_MYLR    uint32 = 0xB339 
	op_MYR     uint32 = 0xB33B 
	op_N       uint32 = 0x5400 
	op_NC      uint32 = 0xD400 
	op_NG      uint32 = 0xE380 
	op_NGR     uint32 = 0xB980 
	op_NGRK    uint32 = 0xB9E4 
	op_NI      uint32 = 0x9400 
	op_NIAI    uint32 = 0xB2FA 
	op_NIHF    uint32 = 0xC00A 
	op_NIHH    uint32 = 0xA504 
	op_NIHL    uint32 = 0xA505 
	op_NILF    uint32 = 0xC00B 
	op_NILH    uint32 = 0xA506 
	op_NILL    uint32 = 0xA507 
	op_NIY     uint32 = 0xEB54 
	op_NR      uint32 = 0x1400 
	op_NRK     uint32 = 0xB9F4 
	op_NTSTG   uint32 = 0xE325 
	op_NY      uint32 = 0xE354 
	op_O       uint32 = 0x5600 
	op_OC      uint32 = 0xD600 
	op_OG      uint32 = 0xE381 
	op_OGR     uint32 = 0xB981 
	op_OGRK    uint32 = 0xB9E6 
	op_OI      uint32 = 0x9600 
	op_OIHF    uint32 = 0xC00C 
	op_OIHH    uint32 = 0xA508 
	op_OIHL    uint32 = 0xA509 
	op_OILF    uint32 = 0xC00D 
	op_OILH    uint32 = 0xA50A 
	op_OILL    uint32 = 0xA50B 
	op_OIY     uint32 = 0xEB56 
	op_OR      uint32 = 0x1600 
	op_ORK     uint32 = 0xB9F6 
	op_OY      uint32 = 0xE356 
	op_PACK    uint32 = 0xF200 
	op_PALB    uint32 = 0xB248 
	op_PC      uint32 = 0xB218 
	op_PCC     uint32 = 0xB92C 
	op_PCKMO   uint32 = 0xB928 
	op_PFD     uint32 = 0xE336 
	op_PFDRL   uint32 = 0xC602 
	op_PFMF    uint32 = 0xB9AF 
	op_PFPO    uint32 = 0x010A 
	op_PGIN    uint32 = 0xB22E 
	op_PGOUT   uint32 = 0xB22F 
	op_PKA     uint32 = 0xE900 
	op_PKU     uint32 = 0xE100 
	op_PLO     uint32 = 0xEE00 
	op_POPCNT  uint32 = 0xB9E1 
	op_PPA     uint32 = 0xB2E8 
	op_PR      uint32 = 0x0101 
	op_PT      uint32 = 0xB228 
	op_PTF     uint32 = 0xB9A2 
	op_PTFF    uint32 = 0x0104 
	op_PTI     uint32 = 0xB99E 
	op_PTLB    uint32 = 0xB20D 
	op_QADTR   uint32 = 0xB3F5 
	op_QAXTR   uint32 = 0xB3FD 
	op_RCHP    uint32 = 0xB23B 
	op_RISBG   uint32 = 0xEC55 
	op_RISBGN  uint32 = 0xEC59 
	op_RISBHG  uint32 = 0xEC5D 
	op_RISBLG  uint32 = 0xEC51 
	op_RLL     uint32 = 0xEB1D 
	op_RLLG    uint32 = 0xEB1C 
	op_RNSBG   uint32 = 0xEC54 
	op_ROSBG   uint32 = 0xEC56 
	op_RP      uint32 = 0xB277 
	op_RRBE    uint32 = 0xB22A 
	op_RRBM    uint32 = 0xB9AE 
	op_RRDTR   uint32 = 0xB3F7 
	op_RRXTR   uint32 = 0xB3FF 
	op_RSCH    uint32 = 0xB238 
	op_RXSBG   uint32 = 0xEC57 
	op_S       uint32 = 0x5B00 
	op_SAC     uint32 = 0xB219 
	op_SACF    uint32 = 0xB279 
	op_SAL     uint32 = 0xB237 
	op_SAM24   uint32 = 0x010C 
	op_SAM31   uint32 = 0x010D 
	op_SAM64   uint32 = 0x010E 
	op_SAR     uint32 = 0xB24E 
	op_SCHM    uint32 = 0xB23C 
	op_SCK     uint32 = 0xB204 
	op_SCKC    uint32 = 0xB206 
	op_SCKPF   uint32 = 0x0107 
	op_SD      uint32 = 0x6B00 
	op_SDB     uint32 = 0xED1B 
	op_SDBR    uint32 = 0xB31B 
	op_SDR     uint32 = 0x2B00 
	op_SDTR    uint32 = 0xB3D3 
	op_SDTRA   uint32 = 0xB3D3 
	op_SE      uint32 = 0x7B00 
	op_SEB     uint32 = 0xED0B 
	op_SEBR    uint32 = 0xB30B 
	op_SER     uint32 = 0x3B00 
	op_SFASR   uint32 = 0xB385 
	op_SFPC    uint32 = 0xB384 
	op_SG      uint32 = 0xE309 
	op_SGF     uint32 = 0xE319 
	op_SGFR    uint32 = 0xB919 
	op_SGR     uint32 = 0xB909 
	op_SGRK    uint32 = 0xB9E9 
	op_SH      uint32 = 0x4B00 
	op_SHHHR   uint32 = 0xB9C9 
	op_SHHLR   uint32 = 0xB9D9 
	op_SHY     uint32 = 0xE37B 
	op_SIGP    uint32 = 0xAE00 
	op_SL      uint32 = 0x5F00 
	op_SLA     uint32 = 0x8B00 
	op_SLAG    uint32 = 0xEB0B 
	op_SLAK    uint32 = 0xEBDD 
	op_SLB     uint32 = 0xE399 
	op_SLBG    uint32 = 0xE389 
	op_SLBGR   uint32 = 0xB989 
	op_SLBR    uint32 = 0xB999 
	op_SLDA    uint32 = 0x8F00 
	op_SLDL    uint32 = 0x8D00 
	op_SLDT    uint32 = 0xED40 
	op_SLFI    uint32 = 0xC205 
	op_SLG     uint32 = 0xE30B 
	op_SLGF    uint32 = 0xE31B 
	op_SLGFI   uint32 = 0xC204 
	op_SLGFR   uint32 = 0xB91B 
	op_SLGR    uint32 = 0xB90B 
	op_SLGRK   uint32 = 0xB9EB 
	op_SLHHHR  uint32 = 0xB9CB 
	op_SLHHLR  uint32 = 0xB9DB 
	op_SLL     uint32 = 0x8900 
	op_SLLG    uint32 = 0xEB0D 
	op_SLLK    uint32 = 0xEBDF 
	op_SLR     uint32 = 0x1F00 
	op_SLRK    uint32 = 0xB9FB 
	op_SLXT    uint32 = 0xED48 
	op_SLY     uint32 = 0xE35F 
	op_SP      uint32 = 0xFB00 
	op_SPKA    uint32 = 0xB20A 
	op_SPM     uint32 = 0x0400 
	op_SPT     uint32 = 0xB208 
	op_SPX     uint32 = 0xB210 
	op_SQD     uint32 = 0xED35 
	op_SQDB    uint32 = 0xED15 
	op_SQDBR   uint32 = 0xB315 
	op_SQDR    uint32 = 0xB244 
	op_SQE     uint32 = 0xED34 
	op_SQEB    uint32 = 0xED14 
	op_SQEBR   uint32 = 0xB314 
	op_SQER    uint32 = 0xB245 
	op_SQXBR   uint32 = 0xB316 
	op_SQXR    uint32 = 0xB336 
	op_SR      uint32 = 0x1B00 
	op_SRA     uint32 = 0x8A00 
	op_SRAG    uint32 = 0xEB0A 
	op_SRAK    uint32 = 0xEBDC 
	op_SRDA    uint32 = 0x8E00 
	op_SRDL    uint32 = 0x8C00 
	op_SRDT    uint32 = 0xED41 
	op_SRK     uint32 = 0xB9F9 
	op_SRL     uint32 = 0x8800 
	op_SRLG    uint32 = 0xEB0C 
	op_SRLK    uint32 = 0xEBDE 
	op_SRNM    uint32 = 0xB299 
	op_SRNMB   uint32 = 0xB2B8 
	op_SRNMT   uint32 = 0xB2B9 
	op_SRP     uint32 = 0xF000 
	op_SRST    uint32 = 0xB25E 
	op_SRSTU   uint32 = 0xB9BE 
	op_SRXT    uint32 = 0xED49 
	op_SSAIR   uint32 = 0xB99F 
	op_SSAR    uint32 = 0xB225 
	op_SSCH    uint32 = 0xB233 
	op_SSKE    uint32 = 0xB22B 
	op_SSM     uint32 = 0x8000 
	op_ST      uint32 = 0x5000 
	op_STAM    uint32 = 0x9B00 
	op_STAMY   uint32 = 0xEB9B 
	op_STAP    uint32 = 0xB212 
	op_STC     uint32 = 0x4200 
	op_STCH    uint32 = 0xE3C3 
	op_STCK    uint32 = 0xB205 
	op_STCKC   uint32 = 0xB207 
	op_STCKE   uint32 = 0xB278 
	op_STCKF   uint32 = 0xB27C 
	op_STCM    uint32 = 0xBE00 
	op_STCMH   uint32 = 0xEB2C 
	op_STCMY   uint32 = 0xEB2D 
	op_STCPS   uint32 = 0xB23A 
	op_STCRW   uint32 = 0xB239 
	op_STCTG   uint32 = 0xEB25 
	op_STCTL   uint32 = 0xB600 
	op_STCY    uint32 = 0xE372 
	op_STD     uint32 = 0x6000 
	op_STDY    uint32 = 0xED67 
	op_STE     uint32 = 0x7000 
	op_STEY    uint32 = 0xED66 
	op_STFH    uint32 = 0xE3CB 
	op_STFL    uint32 = 0xB2B1 
	op_STFLE   uint32 = 0xB2B0 
	op_STFPC   uint32 = 0xB29C 
	op_STG     uint32 = 0xE324 
	op_STGRL   uint32 = 0xC40B 
	op_STH     uint32 = 0x4000 
	op_STHH    uint32 = 0xE3C7 
	op_STHRL   uint32 = 0xC407 
	op_STHY    uint32 = 0xE370 
	op_STIDP   uint32 = 0xB202 
	op_STM     uint32 = 0x9000 
	op_STMG    uint32 = 0xEB24 
	op_STMH    uint32 = 0xEB26 
	op_STMY    uint32 = 0xEB90 
	op_STNSM   uint32 = 0xAC00 
	op_STOC    uint32 = 0xEBF3 
	op_STOCG   uint32 = 0xEBE3 
	op_STOSM   uint32 = 0xAD00 
	op_STPQ    uint32 = 0xE38E 
	op_STPT    uint32 = 0xB209 
	op_STPX    uint32 = 0xB211 
	op_STRAG   uint32 = 0xE502 
	op_STRL    uint32 = 0xC40F 
	op_STRV    uint32 = 0xE33E 
	op_STRVG   uint32 = 0xE32F 
	op_STRVH   uint32 = 0xE33F 
	op_STSCH   uint32 = 0xB234 
	op_STSI    uint32 = 0xB27D 
	op_STURA   uint32 = 0xB246 
	op_STURG   uint32 = 0xB925 
	op_STY     uint32 = 0xE350 
	op_SU      uint32 = 0x7F00 
	op_SUR     uint32 = 0x3F00 
	op_SVC     uint32 = 0x0A00 
	op_SW      uint32 = 0x6F00 
	op_SWR     uint32 = 0x2F00 
	op_SXBR    uint32 = 0xB34B 
	op_SXR     uint32 = 0x3700 
	op_SXTR    uint32 = 0xB3DB 
	op_SXTRA   uint32 = 0xB3DB 
	op_SY      uint32 = 0xE35B 
	op_TABORT  uint32 = 0xB2FC 
	op_TAM     uint32 = 0x010B 
	op_TAR     uint32 = 0xB24C 
	op_TB      uint32 = 0xB22C 
	op_TBDR    uint32 = 0xB351 
	op_TBEDR   uint32 = 0xB350 
	op_TBEGIN  uint32 = 0xE560 
	op_TBEGINC uint32 = 0xE561 
	op_TCDB    uint32 = 0xED11 
	op_TCEB    uint32 = 0xED10 
	op_TCXB    uint32 = 0xED12 
	op_TDCDT   uint32 = 0xED54 
	op_TDCET   uint32 = 0xED50 
	op_TDCXT   uint32 = 0xED58 
	op_TDGDT   uint32 = 0xED55 
	op_TDGET   uint32 = 0xED51 
	op_TDGXT   uint32 = 0xED59 
	op_TEND    uint32 = 0xB2F8 
	op_THDER   uint32 = 0xB358 
	op_THDR    uint32 = 0xB359 
	op_TM      uint32 = 0x9100 
	op_TMH     uint32 = 0xA700 
	op_TMHH    uint32 = 0xA702 
	op_TMHL    uint32 = 0xA703 
	op_TML     uint32 = 0xA701 
	op_TMLH    uint32 = 0xA700 
	op_TMLL    uint32 = 0xA701 
	op_TMY     uint32 = 0xEB51 
	op_TP      uint32 = 0xEBC0 
	op_TPI     uint32 = 0xB236 
	op_TPROT   uint32 = 0xE501 
	op_TR      uint32 = 0xDC00 
	op_TRACE   uint32 = 0x9900 
	op_TRACG   uint32 = 0xEB0F 
	op_TRAP2   uint32 = 0x01FF 
	op_TRAP4   uint32 = 0xB2FF 
	op_TRE     uint32 = 0xB2A5 
	op_TROO    uint32 = 0xB993 
	op_TROT    uint32 = 0xB992 
	op_TRT     uint32 = 0xDD00 
	op_TRTE    uint32 = 0xB9BF 
	op_TRTO    uint32 = 0xB991 
	op_TRTR    uint32 = 0xD000 
	op_TRTRE   uint32 = 0xB9BD 
	op_TRTT    uint32 = 0xB990 
	op_TS      uint32 = 0x9300 
	op_TSCH    uint32 = 0xB235 
	op_UNPK    uint32 = 0xF300 
	op_UNPKA   uint32 = 0xEA00 
	op_UNPKU   uint32 = 0xE200 
	op_UPT     uint32 = 0x0102 
	op_X       uint32 = 0x5700 
	op_XC      uint32 = 0xD700 
	op_XG      uint32 = 0xE382 
	op_XGR     uint32 = 0xB982 
	op_XGRK    uint32 = 0xB9E7 
	op_XI      uint32 = 0x9700 
	op_XIHF    uint32 = 0xC006 
	op_XILF    uint32 = 0xC007 
	op_XIY     uint32 = 0xEB57 
	op_XR      uint32 = 0x1700 
	op_XRK     uint32 = 0xB9F7 
	op_XSCH    uint32 = 0xB276 
	op_XY      uint32 = 0xE357 
	op_ZAP     uint32 = 0xF800 

	
	op_CXPT   uint32 = 0xEDAF 
	op_CDPT   uint32 = 0xEDAE 
	op_CPXT   uint32 = 0xEDAD 
	op_CPDT   uint32 = 0xEDAC 
	op_LZRF   uint32 = 0xE33B 
	op_LZRG   uint32 = 0xE32A 
	op_LCCB   uint32 = 0xE727 
	op_LOCHHI uint32 = 0xEC4E 
	op_LOCHI  uint32 = 0xEC42 
	op_LOCGHI uint32 = 0xEC46 
	op_LOCFH  uint32 = 0xEBE0 
	op_LOCFHR uint32 = 0xB9E0 
	op_LLZRGF uint32 = 0xE33A 
	op_STOCFH uint32 = 0xEBE1 
	op_VA     uint32 = 0xE7F3 
	op_VACC   uint32 = 0xE7F1 
	op_VAC    uint32 = 0xE7BB 
	op_VACCC  uint32 = 0xE7B9 
	op_VN     uint32 = 0xE768 
	op_VNC    uint32 = 0xE769 
	op_VAVG   uint32 = 0xE7F2 
	op_VAVGL  uint32 = 0xE7F0 
	op_VCKSM  uint32 = 0xE766 
	op_VCEQ   uint32 = 0xE7F8 
	op_VCH    uint32 = 0xE7FB 
	op_VCHL   uint32 = 0xE7F9 
	op_VCLZ   uint32 = 0xE753 
	op_VCTZ   uint32 = 0xE752 
	op_VEC    uint32 = 0xE7DB 
	op_VECL   uint32 = 0xE7D9 
	op_VERIM  uint32 = 0xE772 
	op_VERLL  uint32 = 0xE733 
	op_VERLLV uint32 = 0xE773 
	op_VESLV  uint32 = 0xE770 
	op_VESL   uint32 = 0xE730 
	op_VESRA  uint32 = 0xE73A 
	op_VESRAV uint32 = 0xE77A 
	op_VESRL  uint32 = 0xE738 
	op_VESRLV uint32 = 0xE778 
	op_VX     uint32 = 0xE76D 
	op_VFAE   uint32 = 0xE782 
	op_VFEE   uint32 = 0xE780 
	op_VFENE  uint32 = 0xE781 
	op_VFA    uint32 = 0xE7E3 
	op_WFK    uint32 = 0xE7CA 
	op_VFCE   uint32 = 0xE7E8 
	op_VFCH   uint32 = 0xE7EB 
	op_VFCHE  uint32 = 0xE7EA 
	op_WFC    uint32 = 0xE7CB 
	op_VCDG   uint32 = 0xE7C3 
	op_VCDLG  uint32 = 0xE7C1 
	op_VCGD   uint32 = 0xE7C2 
	op_VCLGD  uint32 = 0xE7C0 
	op_VFD    uint32 = 0xE7E5 
	op_VLDE   uint32 = 0xE7C4 
	op_VLED   uint32 = 0xE7C5 
	op_VFM    uint32 = 0xE7E7 
	op_VFMA   uint32 = 0xE78F 
	op_VFMS   uint32 = 0xE78E 
	op_VFPSO  uint32 = 0xE7CC 
	op_VFSQ   uint32 = 0xE7CE 
	op_VFS    uint32 = 0xE7E2 
	op_VFTCI  uint32 = 0xE74A 
	op_VGFM   uint32 = 0xE7B4 
	op_VGFMA  uint32 = 0xE7BC 
	op_VGEF   uint32 = 0xE713 
	op_VGEG   uint32 = 0xE712 
	op_VGBM   uint32 = 0xE744 
	op_VGM    uint32 = 0xE746 
	op_VISTR  uint32 = 0xE75C 
	op_VL     uint32 = 0xE706 
	op_VLR    uint32 = 0xE756 
	op_VLREP  uint32 = 0xE705 
	op_VLC    uint32 = 0xE7DE 
	op_VLEH   uint32 = 0xE701 
	op_VLEF   uint32 = 0xE703 
	op_VLEG   uint32 = 0xE702 
	op_VLEB   uint32 = 0xE700 
	op_VLEIH  uint32 = 0xE741 
	op_VLEIF  uint32 = 0xE743 
	op_VLEIG  uint32 = 0xE742 
	op_VLEIB  uint32 = 0xE740 
	op_VFI    uint32 = 0xE7C7 
	op_VLGV   uint32 = 0xE721 
	op_VLLEZ  uint32 = 0xE704 
	op_VLM    uint32 = 0xE736 
	op_VLP    uint32 = 0xE7DF 
	op_VLBB   uint32 = 0xE707 
	op_VLVG   uint32 = 0xE722 
	op_VLVGP  uint32 = 0xE762 
	op_VLL    uint32 = 0xE737 
	op_VMX    uint32 = 0xE7FF 
	op_VMXL   uint32 = 0xE7FD 
	op_VMRH   uint32 = 0xE761 
	op_VMRL   uint32 = 0xE760 
	op_VMN    uint32 = 0xE7FE 
	op_VMNL   uint32 = 0xE7FC 
	op_VMAE   uint32 = 0xE7AE 
	op_VMAH   uint32 = 0xE7AB 
	op_VMALE  uint32 = 0xE7AC 
	op_VMALH  uint32 = 0xE7A9 
	op_VMALO  uint32 = 0xE7AD 
	op_VMAL   uint32 = 0xE7AA 
	op_VMAO   uint32 = 0xE7AF 
	op_VME    uint32 = 0xE7A6 
	op_VMH    uint32 = 0xE7A3 
	op_VMLE   uint32 = 0xE7A4 
	op_VMLH   uint32 = 0xE7A1 
	op_VMLO   uint32 = 0xE7A5 
	op_VML    uint32 = 0xE7A2 
	op_VMO    uint32 = 0xE7A7 
	op_VNO    uint32 = 0xE76B 
	op_VO     uint32 = 0xE76A 
	op_VPK    uint32 = 0xE794 
	op_VPKLS  uint32 = 0xE795 
	op_VPKS   uint32 = 0xE797 
	op_VPERM  uint32 = 0xE78C 
	op_VPDI   uint32 = 0xE784 
	op_VPOPCT uint32 = 0xE750 
	op_VREP   uint32 = 0xE74D 
	op_VREPI  uint32 = 0xE745 
	op_VSCEF  uint32 = 0xE71B 
	op_VSCEG  uint32 = 0xE71A 
	op_VSEL   uint32 = 0xE78D 
	op_VSL    uint32 = 0xE774 
	op_VSLB   uint32 = 0xE775 
	op_VSLDB  uint32 = 0xE777 
	op_VSRA   uint32 = 0xE77E 
	op_VSRAB  uint32 = 0xE77F 
	op_VSRL   uint32 = 0xE77C 
	op_VSRLB  uint32 = 0xE77D 
	op_VSEG   uint32 = 0xE75F 
	op_VST    uint32 = 0xE70E 
	op_VSTEH  uint32 = 0xE709 
	op_VSTEF  uint32 = 0xE70B 
	op_VSTEG  uint32 = 0xE70A 
	op_VSTEB  uint32 = 0xE708 
	op_VSTM   uint32 = 0xE73E 
	op_VSTL   uint32 = 0xE73F 
	op_VSTRC  uint32 = 0xE78A 
	op_VS     uint32 = 0xE7F7 
	op_VSCBI  uint32 = 0xE7F5 
	op_VSBCBI uint32 = 0xE7BD 
	op_VSBI   uint32 = 0xE7BF 
	op_VSUMG  uint32 = 0xE765 
	op_VSUMQ  uint32 = 0xE767 
	op_VSUM   uint32 = 0xE764 
	op_VTM    uint32 = 0xE7D8 
	op_VUPH   uint32 = 0xE7D7 
	op_VUPLH  uint32 = 0xE7D5 
	op_VUPLL  uint32 = 0xE7D4 
	op_VUPL   uint32 = 0xE7D6 
	op_VMSL   uint32 = 0xE7B8 
)

func oclass(a *obj.Addr) int {
	return int(a.Class) - 1
}



func (c *ctxtz) addrilreloc(sym *obj.LSym, add int64) *obj.Reloc {
	if sym == nil {
		c.ctxt.Diag("require symbol to apply relocation")
	}
	offset := int64(2) 
	rel := obj.Addrel(c.cursym)
	rel.Off = int32(c.pc + offset)
	rel.Siz = 4
	rel.Sym = sym
	rel.Add = add + offset + int64(rel.Siz)
	rel.Type = objabi.R_PCRELDBL
	return rel
}

func (c *ctxtz) addrilrelocoffset(sym *obj.LSym, add, offset int64) *obj.Reloc {
	if sym == nil {
		c.ctxt.Diag("require symbol to apply relocation")
	}
	offset += int64(2) 
	rel := obj.Addrel(c.cursym)
	rel.Off = int32(c.pc + offset)
	rel.Siz = 4
	rel.Sym = sym
	rel.Add = add + offset + int64(rel.Siz)
	rel.Type = objabi.R_PCRELDBL
	return rel
}



func (c *ctxtz) addcallreloc(sym *obj.LSym, add int64) *obj.Reloc {
	if sym == nil {
		c.ctxt.Diag("require symbol to apply relocation")
	}
	offset := int64(2) 
	rel := obj.Addrel(c.cursym)
	rel.Off = int32(c.pc + offset)
	rel.Siz = 4
	rel.Sym = sym
	rel.Add = add + offset + int64(rel.Siz)
	rel.Type = objabi.R_CALL
	return rel
}

func (c *ctxtz) branchMask(p *obj.Prog) CCMask {
	switch p.As {
	case ABRC, ALOCR, ALOCGR,
		ACRJ, ACGRJ, ACIJ, ACGIJ,
		ACLRJ, ACLGRJ, ACLIJ, ACLGIJ:
		return CCMask(p.From.Offset)
	case ABEQ, ACMPBEQ, ACMPUBEQ, AMOVDEQ:
		return Equal
	case ABGE, ACMPBGE, ACMPUBGE, AMOVDGE:
		return GreaterOrEqual
	case ABGT, ACMPBGT, ACMPUBGT, AMOVDGT:
		return Greater
	case ABLE, ACMPBLE, ACMPUBLE, AMOVDLE:
		return LessOrEqual
	case ABLT, ACMPBLT, ACMPUBLT, AMOVDLT:
		return Less
	case ABNE, ACMPBNE, ACMPUBNE, AMOVDNE:
		return NotEqual
	case ABLEU: 
		return NotGreater
	case ABLTU: 
		return LessOrUnordered
	case ABVC:
		return Never 
	case ABVS:
		return Unordered
	}
	c.ctxt.Diag("unknown conditional branch %v", p.As)
	return Always
}

func regtmp(p *obj.Prog) uint32 {
	p.Mark |= USETMP
	return REGTMP
}

func (c *ctxtz) asmout(p *obj.Prog, asm *[]byte) {
	o := c.oplook(p)

	if o == nil {
		return
	}

	
	

	switch o.i {
	default:
		c.ctxt.Diag("unknown index %d", o.i)

	case 0: 
		break

	case 1: 
		switch p.As {
		default:
			c.ctxt.Diag("unhandled operation: %v", p.As)
		case AMOVD:
			zRRE(op_LGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		
		case AMOVW:
			zRRE(op_LGFR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		case AMOVH:
			zRRE(op_LGHR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		case AMOVB:
			zRRE(op_LGBR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		
		case AMOVWZ:
			zRRE(op_LLGFR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		case AMOVHZ:
			zRRE(op_LLGHR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		case AMOVBZ:
			zRRE(op_LLGCR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		
		case AMOVDBR:
			zRRE(op_LRVGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		case AMOVWBR:
			zRRE(op_LRVR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		
		case AFMOVD, AFMOVS:
			zRR(op_LDR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		}

	case 2: 
		r := p.Reg
		if r == 0 {
			r = p.To.Reg
		}

		var opcode uint32

		switch p.As {
		default:
			c.ctxt.Diag("invalid opcode")
		case AADD:
			opcode = op_AGRK
		case AADDC:
			opcode = op_ALGRK
		case AADDE:
			opcode = op_ALCGR
		case AADDW:
			opcode = op_ARK
		case AMULLW:
			opcode = op_MSGFR
		case AMULLD:
			opcode = op_MSGR
		case ADIVW, AMODW:
			opcode = op_DSGFR
		case ADIVWU, AMODWU:
			opcode = op_DLR
		case ADIVD, AMODD:
			opcode = op_DSGR
		case ADIVDU, AMODDU:
			opcode = op_DLGR
		}

		switch p.As {
		default:

		case AADD, AADDC, AADDW:
			if p.As == AADDW && r == p.To.Reg {
				zRR(op_AR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else {
				zRRF(opcode, uint32(p.From.Reg), 0, uint32(p.To.Reg), uint32(r), asm)
			}

		case AADDE, AMULLW, AMULLD:
			if r == p.To.Reg {
				zRRE(opcode, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else if p.From.Reg == p.To.Reg {
				zRRE(opcode, uint32(p.To.Reg), uint32(r), asm)
			} else {
				zRRE(op_LGR, uint32(p.To.Reg), uint32(r), asm)
				zRRE(opcode, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			}

		case ADIVW, ADIVWU, ADIVD, ADIVDU:
			if p.As == ADIVWU || p.As == ADIVDU {
				zRI(op_LGHI, regtmp(p), 0, asm)
			}
			zRRE(op_LGR, REGTMP2, uint32(r), asm)
			zRRE(opcode, regtmp(p), uint32(p.From.Reg), asm)
			zRRE(op_LGR, uint32(p.To.Reg), REGTMP2, asm)

		case AMODW, AMODWU, AMODD, AMODDU:
			if p.As == AMODWU || p.As == AMODDU {
				zRI(op_LGHI, regtmp(p), 0, asm)
			}
			zRRE(op_LGR, REGTMP2, uint32(r), asm)
			zRRE(opcode, regtmp(p), uint32(p.From.Reg), asm)
			zRRE(op_LGR, uint32(p.To.Reg), regtmp(p), asm)

		}

	case 3: 
		v := c.vregoff(&p.From)
		switch p.As {
		case AMOVBZ:
			v = int64(uint8(v))
		case AMOVHZ:
			v = int64(uint16(v))
		case AMOVWZ:
			v = int64(uint32(v))
		case AMOVB:
			v = int64(int8(v))
		case AMOVH:
			v = int64(int16(v))
		case AMOVW:
			v = int64(int32(v))
		}
		if int64(int16(v)) == v {
			zRI(op_LGHI, uint32(p.To.Reg), uint32(v), asm)
		} else if v&0xffff0000 == v {
			zRI(op_LLILH, uint32(p.To.Reg), uint32(v>>16), asm)
		} else if v&0xffff00000000 == v {
			zRI(op_LLIHL, uint32(p.To.Reg), uint32(v>>32), asm)
		} else if uint64(v)&0xffff000000000000 == uint64(v) {
			zRI(op_LLIHH, uint32(p.To.Reg), uint32(v>>48), asm)
		} else if int64(int32(v)) == v {
			zRIL(_a, op_LGFI, uint32(p.To.Reg), uint32(v), asm)
		} else if int64(uint32(v)) == v {
			zRIL(_a, op_LLILF, uint32(p.To.Reg), uint32(v), asm)
		} else if uint64(v)&0xffffffff00000000 == uint64(v) {
			zRIL(_a, op_LLIHF, uint32(p.To.Reg), uint32(v>>32), asm)
		} else {
			zRIL(_a, op_LLILF, uint32(p.To.Reg), uint32(v), asm)
			zRIL(_a, op_IIHF, uint32(p.To.Reg), uint32(v>>32), asm)
		}

	case 4: 
		r := p.Reg
		if r == 0 {
			r = p.To.Reg
		}
		zRRE(op_LGR, REGTMP2, uint32(r), asm)
		zRRE(op_MLGR, regtmp(p), uint32(p.From.Reg), asm)
		switch p.As {
		case AMULHDU:
			
			zRRE(op_LGR, uint32(p.To.Reg), regtmp(p), asm)
		case AMULHD:
			
			
			zRSY(op_SRAG, REGTMP2, uint32(p.From.Reg), 0, 63, asm)
			zRRE(op_NGR, REGTMP2, uint32(r), asm)
			zRRE(op_SGR, regtmp(p), REGTMP2, asm)
			zRSY(op_SRAG, REGTMP2, uint32(r), 0, 63, asm)
			zRRE(op_NGR, REGTMP2, uint32(p.From.Reg), asm)
			zRRF(op_SGRK, REGTMP2, 0, uint32(p.To.Reg), regtmp(p), asm)
		}

	case 5: 
		zI(op_SVC, 0, asm)

	case 6: 
		var oprr, oprre, oprrf uint32
		switch p.As {
		case AAND:
			oprre = op_NGR
			oprrf = op_NGRK
		case AANDW:
			oprr = op_NR
			oprrf = op_NRK
		case AOR:
			oprre = op_OGR
			oprrf = op_OGRK
		case AORW:
			oprr = op_OR
			oprrf = op_ORK
		case AXOR:
			oprre = op_XGR
			oprrf = op_XGRK
		case AXORW:
			oprr = op_XR
			oprrf = op_XRK
		}
		if p.Reg == 0 {
			if oprr != 0 {
				zRR(oprr, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else {
				zRRE(oprre, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			}
		} else {
			zRRF(oprrf, uint32(p.Reg), 0, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		}

	case 7: 
		d2 := c.vregoff(&p.From)
		b2 := p.From.Reg
		r3 := p.Reg
		if r3 == 0 {
			r3 = p.To.Reg
		}
		r1 := p.To.Reg
		var opcode uint32
		switch p.As {
		default:
		case ASLD:
			opcode = op_SLLG
		case ASRD:
			opcode = op_SRLG
		case ASLW:
			opcode = op_SLLK
		case ASRW:
			opcode = op_SRLK
		case ARLL:
			opcode = op_RLL
		case ARLLG:
			opcode = op_RLLG
		case ASRAW:
			opcode = op_SRAK
		case ASRAD:
			opcode = op_SRAG
		}
		zRSY(opcode, uint32(r1), uint32(r3), uint32(b2), uint32(d2), asm)

	case 8: 
		if p.To.Reg&1 != 0 {
			c.ctxt.Diag("target must be an even-numbered register")
		}
		
		zRRE(op_FLOGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 9: 
		zRRE(op_POPCNT, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 10: 
		r := int(p.Reg)

		switch p.As {
		default:
		case ASUB:
			if r == 0 {
				zRRE(op_SGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else {
				zRRF(op_SGRK, uint32(p.From.Reg), 0, uint32(p.To.Reg), uint32(r), asm)
			}
		case ASUBC:
			if r == 0 {
				zRRE(op_SLGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else {
				zRRF(op_SLGRK, uint32(p.From.Reg), 0, uint32(p.To.Reg), uint32(r), asm)
			}
		case ASUBE:
			if r == 0 {
				r = int(p.To.Reg)
			}
			if r == int(p.To.Reg) {
				zRRE(op_SLBGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else if p.From.Reg == p.To.Reg {
				zRRE(op_LGR, regtmp(p), uint32(p.From.Reg), asm)
				zRRE(op_LGR, uint32(p.To.Reg), uint32(r), asm)
				zRRE(op_SLBGR, uint32(p.To.Reg), regtmp(p), asm)
			} else {
				zRRE(op_LGR, uint32(p.To.Reg), uint32(r), asm)
				zRRE(op_SLBGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			}
		case ASUBW:
			if r == 0 {
				zRR(op_SR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
			} else {
				zRRF(op_SRK, uint32(p.From.Reg), 0, uint32(p.To.Reg), uint32(r), asm)
			}
		}

	case 11: 
		v := int32(0)

		if p.To.Target() != nil {
			v = int32((p.To.Target().Pc - p.Pc) >> 1)
		}

		if p.As == ABR && p.To.Sym == nil && int32(int16(v)) == v {
			zRI(op_BRC, 0xF, uint32(v), asm)
		} else {
			if p.As == ABL {
				zRIL(_b, op_BRASL, uint32(REG_LR), uint32(v), asm)
			} else {
				zRIL(_c, op_BRCL, 0xF, uint32(v), asm)
			}
			if p.To.Sym != nil {
				c.addcallreloc(p.To.Sym, p.To.Offset)
			}
		}

	case 12:
		r1 := p.To.Reg
		d2 := c.vregoff(&p.From)
		b2 := p.From.Reg
		if b2 == 0 {
			b2 = REGSP
		}
		x2 := p.From.Index
		if -DISP20/2 > d2 || d2 >= DISP20/2 {
			zRIL(_a, op_LGFI, regtmp(p), uint32(d2), asm)
			if x2 != 0 {
				zRX(op_LA, regtmp(p), regtmp(p), uint32(x2), 0, asm)
			}
			x2 = int16(regtmp(p))
			d2 = 0
		}
		var opx, opxy uint32
		switch p.As {
		case AADD:
			opxy = op_AG
		case AADDC:
			opxy = op_ALG
		case AADDE:
			opxy = op_ALCG
		case AADDW:
			opx = op_A
			opxy = op_AY
		case AMULLW:
			opx = op_MS
			opxy = op_MSY
		case AMULLD:
			opxy = op_MSG
		case ASUB:
			opxy = op_SG
		case ASUBC:
			opxy = op_SLG
		case ASUBE:
			opxy = op_SLBG
		case ASUBW:
			opx = op_S
			opxy = op_SY
		case AAND:
			opxy = op_NG
		case AANDW:
			opx = op_N
			opxy = op_NY
		case AOR:
			opxy = op_OG
		case AORW:
			opx = op_O
			opxy = op_OY
		case AXOR:
			opxy = op_XG
		case AXORW:
			opx = op_X
			opxy = op_XY
		}
		if opx != 0 && 0 <= d2 && d2 < DISP12 {
			zRX(opx, uint32(r1), uint32(x2), uint32(b2), uint32(d2), asm)
		} else {
			zRXY(opxy, uint32(r1), uint32(x2), uint32(b2), uint32(d2), asm)
		}

	case 13: 
		r1 := p.To.Reg
		r2 := p.RestArgs[2].Reg
		i3 := uint8(p.From.Offset)        
		i4 := uint8(p.RestArgs[0].Offset) 
		i5 := uint8(p.RestArgs[1].Offset) 
		switch p.As {
		case ARNSBGT, ARXSBGT, AROSBGT:
			i3 |= 0x80 
		case ARISBGZ, ARISBGNZ, ARISBHGZ, ARISBLGZ:
			i4 |= 0x80 
		}
		var opcode uint32
		switch p.As {
		case ARNSBG, ARNSBGT:
			opcode = op_RNSBG
		case ARXSBG, ARXSBGT:
			opcode = op_RXSBG
		case AROSBG, AROSBGT:
			opcode = op_ROSBG
		case ARISBG, ARISBGZ:
			opcode = op_RISBG
		case ARISBGN, ARISBGNZ:
			opcode = op_RISBGN
		case ARISBHG, ARISBHGZ:
			opcode = op_RISBHG
		case ARISBLG, ARISBLGZ:
			opcode = op_RISBLG
		}
		zRIE(_f, uint32(opcode), uint32(r1), uint32(r2), 0, uint32(i3), uint32(i4), 0, uint32(i5), asm)

	case 15: 
		r := p.To.Reg
		if p.As == ABCL || p.As == ABL {
			zRR(op_BASR, uint32(REG_LR), uint32(r), asm)
		} else {
			zRR(op_BCR, uint32(Always), uint32(r), asm)
		}

	case 16: 
		v := int32(0)
		if p.To.Target() != nil {
			v = int32((p.To.Target().Pc - p.Pc) >> 1)
		}
		mask := uint32(c.branchMask(p))
		if p.To.Sym == nil && int32(int16(v)) == v {
			zRI(op_BRC, mask, uint32(v), asm)
		} else {
			zRIL(_c, op_BRCL, mask, uint32(v), asm)
		}
		if p.To.Sym != nil {
			c.addrilreloc(p.To.Sym, p.To.Offset)
		}

	case 17: 
		m3 := uint32(c.branchMask(p))
		zRRF(op_LOCGR, m3, 0, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 18: 
		if p.As == ABL {
			zRR(op_BASR, uint32(REG_LR), uint32(p.To.Reg), asm)
		} else {
			zRR(op_BCR, uint32(Always), uint32(p.To.Reg), asm)
		}

	case 19: 
		d := c.vregoff(&p.From)
		zRIL(_b, op_LARL, uint32(p.To.Reg), 0, asm)
		if d&1 != 0 {
			zRX(op_LA, uint32(p.To.Reg), uint32(p.To.Reg), 0, 1, asm)
			d -= 1
		}
		c.addrilreloc(p.From.Sym, d)

	case 21: 
		v := c.vregoff(&p.From)
		r := p.Reg
		if r == 0 {
			r = p.To.Reg
		}
		switch p.As {
		case ASUB:
			zRIL(_a, op_LGFI, uint32(regtmp(p)), uint32(v), asm)
			zRRF(op_SLGRK, uint32(regtmp(p)), 0, uint32(p.To.Reg), uint32(r), asm)
		case ASUBC:
			if r != p.To.Reg {
				zRRE(op_LGR, uint32(p.To.Reg), uint32(r), asm)
			}
			zRIL(_a, op_SLGFI, uint32(p.To.Reg), uint32(v), asm)
		case ASUBW:
			if r != p.To.Reg {
				zRR(op_LR, uint32(p.To.Reg), uint32(r), asm)
			}
			zRIL(_a, op_SLFI, uint32(p.To.Reg), uint32(v), asm)
		}

	case 22: 
		v := c.vregoff(&p.From)
		r := p.Reg
		if r == 0 {
			r = p.To.Reg
		}
		var opri, opril, oprie uint32
		switch p.As {
		case AADD:
			opri = op_AGHI
			opril = op_AGFI
			oprie = op_AGHIK
		case AADDC:
			opril = op_ALGFI
			oprie = op_ALGHSIK
		case AADDW:
			opri = op_AHI
			opril = op_AFI
			oprie = op_AHIK
		case AMULLW:
			opri = op_MHI
			opril = op_MSFI
		case AMULLD:
			opri = op_MGHI
			opril = op_MSGFI
		}
		if r != p.To.Reg && (oprie == 0 || int64(int16(v)) != v) {
			switch p.As {
			case AADD, AADDC, AMULLD:
				zRRE(op_LGR, uint32(p.To.Reg), uint32(r), asm)
			case AADDW, AMULLW:
				zRR(op_LR, uint32(p.To.Reg), uint32(r), asm)
			}
			r = p.To.Reg
		}
		if opri != 0 && r == p.To.Reg && int64(int16(v)) == v {
			zRI(opri, uint32(p.To.Reg), uint32(v), asm)
		} else if oprie != 0 && int64(int16(v)) == v {
			zRIE(_d, oprie, uint32(p.To.Reg), uint32(r), uint32(v), 0, 0, 0, 0, asm)
		} else {
			zRIL(_a, opril, uint32(p.To.Reg), uint32(v), asm)
		}

	case 23: 
		
		v := c.vregoff(&p.From)
		switch p.As {
		default:
			c.ctxt.Diag("%v is not supported", p)
		case AAND:
			if v >= 0 { 
				zRIL(_a, op_LGFI, regtmp(p), uint32(v), asm)
				zRRE(op_NGR, uint32(p.To.Reg), regtmp(p), asm)
			} else if int64(int16(v)) == v {
				zRI(op_NILL, uint32(p.To.Reg), uint32(v), asm)
			} else { 
				zRIL(_a, op_NILF, uint32(p.To.Reg), uint32(v), asm)
			}
		case AOR:
			if int64(uint32(v)) != v { 
				zRIL(_a, op_LGFI, regtmp(p), uint32(v), asm)
				zRRE(op_OGR, uint32(p.To.Reg), regtmp(p), asm)
			} else if int64(uint16(v)) == v {
				zRI(op_OILL, uint32(p.To.Reg), uint32(v), asm)
			} else {
				zRIL(_a, op_OILF, uint32(p.To.Reg), uint32(v), asm)
			}
		case AXOR:
			if int64(uint32(v)) != v { 
				zRIL(_a, op_LGFI, regtmp(p), uint32(v), asm)
				zRRE(op_XGR, uint32(p.To.Reg), regtmp(p), asm)
			} else {
				zRIL(_a, op_XILF, uint32(p.To.Reg), uint32(v), asm)
			}
		}

	case 24: 
		v := c.vregoff(&p.From)
		switch p.As {
		case AANDW:
			if uint32(v&0xffff0000) == 0xffff0000 {
				zRI(op_NILL, uint32(p.To.Reg), uint32(v), asm)
			} else if uint32(v&0x0000ffff) == 0x0000ffff {
				zRI(op_NILH, uint32(p.To.Reg), uint32(v)>>16, asm)
			} else {
				zRIL(_a, op_NILF, uint32(p.To.Reg), uint32(v), asm)
			}
		case AORW:
			if uint32(v&0xffff0000) == 0 {
				zRI(op_OILL, uint32(p.To.Reg), uint32(v), asm)
			} else if uint32(v&0x0000ffff) == 0 {
				zRI(op_OILH, uint32(p.To.Reg), uint32(v)>>16, asm)
			} else {
				zRIL(_a, op_OILF, uint32(p.To.Reg), uint32(v), asm)
			}
		case AXORW:
			zRIL(_a, op_XILF, uint32(p.To.Reg), uint32(v), asm)
		}

	case 25: 
		m3 := uint32(c.branchMask(p))
		var opcode uint32
		switch p.As {
		case ALOCR:
			opcode = op_LOCR
		case ALOCGR:
			opcode = op_LOCGR
		}
		zRRF(opcode, m3, 0, uint32(p.To.Reg), uint32(p.Reg), asm)

	case 26: 
		v := c.regoff(&p.From)
		r := p.From.Reg
		if r == 0 {
			r = REGSP
		}
		i := p.From.Index
		if v >= 0 && v < DISP12 {
			zRX(op_LA, uint32(p.To.Reg), uint32(r), uint32(i), uint32(v), asm)
		} else if v >= -DISP20/2 && v < DISP20/2 {
			zRXY(op_LAY, uint32(p.To.Reg), uint32(r), uint32(i), uint32(v), asm)
		} else {
			zRIL(_a, op_LGFI, regtmp(p), uint32(v), asm)
			zRX(op_LA, uint32(p.To.Reg), uint32(r), regtmp(p), uint32(i), asm)
		}

	case 31: 
		wd := uint64(c.vregoff(&p.From))
		*asm = append(*asm,
			uint8(wd>>56),
			uint8(wd>>48),
			uint8(wd>>40),
			uint8(wd>>32),
			uint8(wd>>24),
			uint8(wd>>16),
			uint8(wd>>8),
			uint8(wd))

	case 32: 
		var opcode uint32
		switch p.As {
		default:
			c.ctxt.Diag("invalid opcode")
		case AFADD:
			opcode = op_ADBR
		case AFADDS:
			opcode = op_AEBR
		case AFDIV:
			opcode = op_DDBR
		case AFDIVS:
			opcode = op_DEBR
		case AFMUL:
			opcode = op_MDBR
		case AFMULS:
			opcode = op_MEEBR
		case AFSUB:
			opcode = op_SDBR
		case AFSUBS:
			opcode = op_SEBR
		}
		zRRE(opcode, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 33: 
		r := p.From.Reg
		if oclass(&p.From) == C_NONE {
			r = p.To.Reg
		}
		var opcode uint32
		switch p.As {
		default:
		case AFABS:
			opcode = op_LPDBR
		case AFNABS:
			opcode = op_LNDBR
		case ALPDFR:
			opcode = op_LPDFR
		case ALNDFR:
			opcode = op_LNDFR
		case AFNEG:
			opcode = op_LCDFR
		case AFNEGS:
			opcode = op_LCEBR
		case ALEDBR:
			opcode = op_LEDBR
		case ALDEBR:
			opcode = op_LDEBR
		case AFSQRT:
			opcode = op_SQDBR
		case AFSQRTS:
			opcode = op_SQEBR
		}
		zRRE(opcode, uint32(p.To.Reg), uint32(r), asm)

	case 34: 
		var opcode uint32
		switch p.As {
		default:
			c.ctxt.Diag("invalid opcode")
		case AFMADD:
			opcode = op_MADBR
		case AFMADDS:
			opcode = op_MAEBR
		case AFMSUB:
			opcode = op_MSDBR
		case AFMSUBS:
			opcode = op_MSEBR
		}
		zRRD(opcode, uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), asm)

	case 35: 
		d2 := c.regoff(&p.To)
		b2 := p.To.Reg
		if b2 == 0 {
			b2 = REGSP
		}
		x2 := p.To.Index
		if d2 < -DISP20/2 || d2 >= DISP20/2 {
			zRIL(_a, op_LGFI, regtmp(p), uint32(d2), asm)
			if x2 != 0 {
				zRX(op_LA, regtmp(p), regtmp(p), uint32(x2), 0, asm)
			}
			x2 = int16(regtmp(p))
			d2 = 0
		}
		
		if op, ok := c.zopstore12(p.As); ok && isU12(d2) {
			zRX(op, uint32(p.From.Reg), uint32(x2), uint32(b2), uint32(d2), asm)
		} else {
			zRXY(c.zopstore(p.As), uint32(p.From.Reg), uint32(x2), uint32(b2), uint32(d2), asm)
		}

	case 36: 
		d2 := c.regoff(&p.From)
		b2 := p.From.Reg
		if b2 == 0 {
			b2 = REGSP
		}
		x2 := p.From.Index
		if d2 < -DISP20/2 || d2 >= DISP20/2 {
			zRIL(_a, op_LGFI, regtmp(p), uint32(d2), asm)
			if x2 != 0 {
				zRX(op_LA, regtmp(p), regtmp(p), uint32(x2), 0, asm)
			}
			x2 = int16(regtmp(p))
			d2 = 0
		}
		
		if op, ok := c.zopload12(p.As); ok && isU12(d2) {
			zRX(op, uint32(p.To.Reg), uint32(x2), uint32(b2), uint32(d2), asm)
		} else {
			zRXY(c.zopload(p.As), uint32(p.To.Reg), uint32(x2), uint32(b2), uint32(d2), asm)
		}

	case 40: 
		wd := uint32(c.regoff(&p.From))
		if p.As == AWORD { 
			*asm = append(*asm, uint8(wd>>24), uint8(wd>>16), uint8(wd>>8), uint8(wd))
		} else { 
			*asm = append(*asm, uint8(wd))
		}

	case 41: 
		r1 := p.From.Reg
		ri2 := (p.To.Target().Pc - p.Pc) >> 1
		if int64(int16(ri2)) != ri2 {
			c.ctxt.Diag("branch target too far away")
		}
		var opcode uint32
		switch p.As {
		case ABRCT:
			opcode = op_BRCT
		case ABRCTG:
			opcode = op_BRCTG
		}
		zRI(opcode, uint32(r1), uint32(ri2), asm)

	case 47: 
		r := p.From.Reg
		if r == 0 {
			r = p.To.Reg
		}
		switch p.As {
		case ANEG:
			zRRE(op_LCGR, uint32(p.To.Reg), uint32(r), asm)
		case ANEGW:
			zRRE(op_LCGFR, uint32(p.To.Reg), uint32(r), asm)
		}

	case 48: 
		m3 := c.vregoff(&p.From)
		if 0 > m3 || m3 > 7 {
			c.ctxt.Diag("mask (%v) must be in the range [0, 7]", m3)
		}
		var opcode uint32
		switch p.As {
		case AFIEBR:
			opcode = op_FIEBR
		case AFIDBR:
			opcode = op_FIDBR
		}
		zRRF(opcode, uint32(m3), 0, uint32(p.To.Reg), uint32(p.Reg), asm)

	case 49: 
		zRRF(op_CPSDR, uint32(p.From.Reg), 0, uint32(p.To.Reg), uint32(p.Reg), asm)

	case 50: 
		var opcode uint32
		switch p.As {
		case ALTEBR:
			opcode = op_LTEBR
		case ALTDBR:
			opcode = op_LTDBR
		}
		zRRE(opcode, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 51: 
		var opcode uint32
		switch p.As {
		case ATCEB:
			opcode = op_TCEB
		case ATCDB:
			opcode = op_TCDB
		}
		d2 := c.regoff(&p.To)
		zRXE(opcode, uint32(p.From.Reg), 0, 0, uint32(d2), 0, asm)

	case 62: 
		zRRE(op_MLGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 66:
		zRR(op_BCR, uint32(Never), 0, asm)

	case 67: 
		var opcode uint32
		switch p.As {
		case AFMOVS:
			opcode = op_LZER
		case AFMOVD:
			opcode = op_LZDR
		}
		zRRE(opcode, uint32(p.To.Reg), 0, asm)

	case 68: 
		zRRE(op_EAR, uint32(p.To.Reg), uint32(p.From.Reg-REG_AR0), asm)

	case 69: 
		zRRE(op_SAR, uint32(p.To.Reg-REG_AR0), uint32(p.From.Reg), asm)

	case 70: 
		if p.As == ACMPW || p.As == ACMPWU {
			zRR(c.zoprr(p.As), uint32(p.From.Reg), uint32(p.To.Reg), asm)
		} else {
			zRRE(c.zoprre(p.As), uint32(p.From.Reg), uint32(p.To.Reg), asm)
		}

	case 71: 
		v := c.vregoff(&p.To)
		switch p.As {
		case ACMP, ACMPW:
			if int64(int32(v)) != v {
				c.ctxt.Diag("%v overflows an int32", v)
			}
		case ACMPU, ACMPWU:
			if int64(uint32(v)) != v {
				c.ctxt.Diag("%v overflows a uint32", v)
			}
		}
		if p.As == ACMP && int64(int16(v)) == v {
			zRI(op_CGHI, uint32(p.From.Reg), uint32(v), asm)
		} else if p.As == ACMPW && int64(int16(v)) == v {
			zRI(op_CHI, uint32(p.From.Reg), uint32(v), asm)
		} else {
			zRIL(_a, c.zopril(p.As), uint32(p.From.Reg), uint32(v), asm)
		}

	case 72: 
		v := c.regoff(&p.From)
		d := c.regoff(&p.To)
		r := p.To.Reg
		if p.To.Index != 0 {
			c.ctxt.Diag("cannot use index register")
		}
		if r == 0 {
			r = REGSP
		}
		var opcode uint32
		switch p.As {
		case AMOVD:
			opcode = op_MVGHI
		case AMOVW, AMOVWZ:
			opcode = op_MVHI
		case AMOVH, AMOVHZ:
			opcode = op_MVHHI
		case AMOVB, AMOVBZ:
			opcode = op_MVI
		}
		if d < 0 || d >= DISP12 {
			if r == int16(regtmp(p)) {
				c.ctxt.Diag("displacement must be in range [0, 4096) to use %v", r)
			}
			if d >= -DISP20/2 && d < DISP20/2 {
				if opcode == op_MVI {
					opcode = op_MVIY
				} else {
					zRXY(op_LAY, uint32(regtmp(p)), 0, uint32(r), uint32(d), asm)
					r = int16(regtmp(p))
					d = 0
				}
			} else {
				zRIL(_a, op_LGFI, regtmp(p), uint32(d), asm)
				zRX(op_LA, regtmp(p), regtmp(p), uint32(r), 0, asm)
				r = int16(regtmp(p))
				d = 0
			}
		}
		switch opcode {
		case op_MVI:
			zSI(opcode, uint32(v), uint32(r), uint32(d), asm)
		case op_MVIY:
			zSIY(opcode, uint32(v), uint32(r), uint32(d), asm)
		default:
			zSIL(opcode, uint32(r), uint32(d), uint32(v), asm)
		}

	case 74: 
		i2 := c.regoff(&p.To)
		switch p.As {
		case AMOVD:
			zRIL(_b, op_STGRL, uint32(p.From.Reg), 0, asm)
		case AMOVW, AMOVWZ: 
			zRIL(_b, op_STRL, uint32(p.From.Reg), 0, asm)
		case AMOVH, AMOVHZ: 
			zRIL(_b, op_STHRL, uint32(p.From.Reg), 0, asm)
		case AMOVB, AMOVBZ: 
			zRIL(_b, op_LARL, regtmp(p), 0, asm)
			adj := uint32(0) 
			if i2&1 != 0 {
				i2 -= 1
				adj = 1
			}
			zRX(op_STC, uint32(p.From.Reg), 0, regtmp(p), adj, asm)
		case AFMOVD:
			zRIL(_b, op_LARL, regtmp(p), 0, asm)
			zRX(op_STD, uint32(p.From.Reg), 0, regtmp(p), 0, asm)
		case AFMOVS:
			zRIL(_b, op_LARL, regtmp(p), 0, asm)
			zRX(op_STE, uint32(p.From.Reg), 0, regtmp(p), 0, asm)
		}
		c.addrilreloc(p.To.Sym, int64(i2))

	case 75: 
		i2 := c.regoff(&p.From)
		switch p.As {
		case AMOVD:
			if i2&1 != 0 {
				zRIL(_b, op_LARL, regtmp(p), 0, asm)
				zRXY(op_LG, uint32(p.To.Reg), regtmp(p), 0, 1, asm)
				i2 -= 1
			} else {
				zRIL(_b, op_LGRL, uint32(p.To.Reg), 0, asm)
			}
		case AMOVW:
			zRIL(_b, op_LGFRL, uint32(p.To.Reg), 0, asm)
		case AMOVWZ:
			zRIL(_b, op_LLGFRL, uint32(p.To.Reg), 0, asm)
		case AMOVH:
			zRIL(_b, op_LGHRL, uint32(p.To.Reg), 0, asm)
		case AMOVHZ:
			zRIL(_b, op_LLGHRL, uint32(p.To.Reg), 0, asm)
		case AMOVB, AMOVBZ:
			zRIL(_b, op_LARL, regtmp(p), 0, asm)
			adj := uint32(0) 
			if i2&1 != 0 {
				i2 -= 1
				adj = 1
			}
			switch p.As {
			case AMOVB:
				zRXY(op_LGB, uint32(p.To.Reg), 0, regtmp(p), adj, asm)
			case AMOVBZ:
				zRXY(op_LLGC, uint32(p.To.Reg), 0, regtmp(p), adj, asm)
			}
		case AFMOVD:
			zRIL(_a, op_LARL, regtmp(p), 0, asm)
			zRX(op_LD, uint32(p.To.Reg), 0, regtmp(p), 0, asm)
		case AFMOVS:
			zRIL(_a, op_LARL, regtmp(p), 0, asm)
			zRX(op_LE, uint32(p.To.Reg), 0, regtmp(p), 0, asm)
		}
		c.addrilreloc(p.From.Sym, int64(i2))

	case 76: 
		zRR(op_SPM, uint32(p.From.Reg), 0, asm)

	case 77: 
		if p.From.Offset > 255 || p.From.Offset < 1 {
			c.ctxt.Diag("illegal system call; system call number out of range: %v", p)
			zE(op_TRAP2, asm) 
		} else {
			zI(op_SVC, uint32(p.From.Offset), asm)
		}

	case 78: 
		
		
		*asm = append(*asm, 0, 0, 0, 0)

	case 79: 
		v := c.regoff(&p.To)
		if v < 0 {
			v = 0
		}
		if p.As == ACS {
			zRS(op_CS, uint32(p.From.Reg), uint32(p.Reg), uint32(p.To.Reg), uint32(v), asm)
		} else if p.As == ACSG {
			zRSY(op_CSG, uint32(p.From.Reg), uint32(p.Reg), uint32(p.To.Reg), uint32(v), asm)
		}

	case 80: 
		zRR(op_BCR, uint32(NotEqual), 0, asm)

	case 81: 
		switch p.As {
		case ALDGR:
			zRRE(op_LDGR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		case ALGDR:
			zRRE(op_LGDR, uint32(p.To.Reg), uint32(p.From.Reg), asm)
		}

	case 82: 
		var opcode uint32
		switch p.As {
		default:
			log.Fatalf("unexpected opcode %v", p.As)
		case ACEFBRA:
			opcode = op_CEFBRA
		case ACDFBRA:
			opcode = op_CDFBRA
		case ACEGBRA:
			opcode = op_CEGBRA
		case ACDGBRA:
			opcode = op_CDGBRA
		case ACELFBR:
			opcode = op_CELFBR
		case ACDLFBR:
			opcode = op_CDLFBR
		case ACELGBR:
			opcode = op_CELGBR
		case ACDLGBR:
			opcode = op_CDLGBR
		}
		
		
		
		
		zRRF(opcode, 0, 0, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 83: 
		var opcode uint32
		switch p.As {
		default:
			log.Fatalf("unexpected opcode %v", p.As)
		case ACFEBRA:
			opcode = op_CFEBRA
		case ACFDBRA:
			opcode = op_CFDBRA
		case ACGEBRA:
			opcode = op_CGEBRA
		case ACGDBRA:
			opcode = op_CGDBRA
		case ACLFEBR:
			opcode = op_CLFEBR
		case ACLFDBR:
			opcode = op_CLFDBR
		case ACLGEBR:
			opcode = op_CLGEBR
		case ACLGDBR:
			opcode = op_CLGDBR
		}
		
		
		zRRF(opcode, 5, 0, uint32(p.To.Reg), uint32(p.From.Reg), asm)

	case 84: 
		l := c.regoff(&p.From)
		if l < 1 || l > 256 {
			c.ctxt.Diag("number of bytes (%v) not in range [1,256]", l)
		}
		if p.GetFrom3().Index != 0 || p.To.Index != 0 {
			c.ctxt.Diag("cannot use index reg")
		}
		b1 := p.To.Reg
		b2 := p.GetFrom3().Reg
		if b1 == 0 {
			b1 = REGSP
		}
		if b2 == 0 {
			b2 = REGSP
		}
		d1 := c.regoff(&p.To)
		d2 := c.regoff(p.GetFrom3())
		if d1 < 0 || d1 >= DISP12 {
			if b2 == int16(regtmp(p)) {
				c.ctxt.Diag("regtmp(p) conflict")
			}
			if b1 != int16(regtmp(p)) {
				zRRE(op_LGR, regtmp(p), uint32(b1), asm)
			}
			zRIL(_a, op_AGFI, regtmp(p), uint32(d1), asm)
			if d1 == d2 && b1 == b2 {
				d2 = 0
				b2 = int16(regtmp(p))
			}
			d1 = 0
			b1 = int16(regtmp(p))
		}
		if d2 < 0 || d2 >= DISP12 {
			if b1 == REGTMP2 {
				c.ctxt.Diag("REGTMP2 conflict")
			}
			if b2 != REGTMP2 {
				zRRE(op_LGR, REGTMP2, uint32(b2), asm)
			}
			zRIL(_a, op_AGFI, REGTMP2, uint32(d2), asm)
			d2 = 0
			b2 = REGTMP2
		}
		var opcode uint32
		switch p.As {
		default:
			c.ctxt.Diag("unexpected opcode %v", p.As)
		case AMVC:
			opcode = op_MVC
		case AMVCIN:
			opcode = op_MVCIN
		case ACLC:
			opcode = op_CLC
			
			b1, b2 = b2, b1
			d1, d2 = d2, d1
		case AXC:
			opcode = op_XC
		case AOC:
			opcode = op_OC
		case ANC:
			opcode = op_NC
		}
		zSS(_a, opcode, uint32(l-1), 0, uint32(b1), uint32(d1), uint32(b2), uint32(d2), asm)

	case 85: 
		v := c.regoff(&p.From)
		if p.From.Sym == nil {
			if (v & 1) != 0 {
				c.ctxt.Diag("cannot use LARL with odd offset: %v", v)
			}
		} else {
			c.addrilreloc(p.From.Sym, int64(v))
			v = 0
		}
		zRIL(_b, op_LARL, uint32(p.To.Reg), uint32(v>>1), asm)

	case 86: 
		d := c.vregoff(&p.From)
		x := p.From.Index
		b := p.From.Reg
		if b == 0 {
			b = REGSP
		}
		switch p.As {
		case ALA:
			zRX(op_LA, uint32(p.To.Reg), uint32(x), uint32(b), uint32(d), asm)
		case ALAY:
			zRXY(op_LAY, uint32(p.To.Reg), uint32(x), uint32(b), uint32(d), asm)
		}

	case 87: 
		v := c.vregoff(&p.From)
		if p.From.Sym == nil {
			if v&1 != 0 {
				c.ctxt.Diag("cannot use EXRL with odd offset: %v", v)
			}
		} else {
			c.addrilreloc(p.From.Sym, v)
			v = 0
		}
		zRIL(_b, op_EXRL, uint32(p.To.Reg), uint32(v>>1), asm)

	case 88: 
		var opcode uint32
		switch p.As {
		case ASTCK:
			opcode = op_STCK
		case ASTCKC:
			opcode = op_STCKC
		case ASTCKE:
			opcode = op_STCKE
		case ASTCKF:
			opcode = op_STCKF
		}
		v := c.vregoff(&p.To)
		r := p.To.Reg
		if r == 0 {
			r = REGSP
		}
		zS(opcode, uint32(r), uint32(v), asm)

	case 89: 
		var v int32
		if p.To.Target() != nil {
			v = int32((p.To.Target().Pc - p.Pc) >> 1)
		}

		
		r1, r2 := p.From.Reg, p.Reg
		if p.From.Type == obj.TYPE_CONST {
			r1, r2 = p.Reg, p.RestArgs[0].Reg
		}
		m3 := uint32(c.branchMask(p))

		var opcode uint32
		switch p.As {
		case ACRJ:
			
			opcode = op_CRJ
		case ACGRJ, ACMPBEQ, ACMPBGE, ACMPBGT, ACMPBLE, ACMPBLT, ACMPBNE:
			
			opcode = op_CGRJ
		case ACLRJ:
			
			opcode = op_CLRJ
		case ACLGRJ, ACMPUBEQ, ACMPUBGE, ACMPUBGT, ACMPUBLE, ACMPUBLT, ACMPUBNE:
			
			opcode = op_CLGRJ
		}

		if int32(int16(v)) != v {
			
			
			
			
			
			
			
			
			m3 ^= 0xe 
			zRIE(_b, opcode, uint32(r1), uint32(r2), uint32(sizeRIE+sizeRIL)/2, 0, 0, m3, 0, asm)
			zRIL(_c, op_BRCL, uint32(Always), uint32(v-sizeRIE/2), asm)
		} else {
			zRIE(_b, opcode, uint32(r1), uint32(r2), uint32(v), 0, 0, m3, 0, asm)
		}

	case 90: 
		var v int32
		if p.To.Target() != nil {
			v = int32((p.To.Target().Pc - p.Pc) >> 1)
		}

		
		r1, i2 := p.From.Reg, p.RestArgs[0].Offset
		if p.From.Type == obj.TYPE_CONST {
			r1 = p.Reg
		}
		m3 := uint32(c.branchMask(p))

		var opcode uint32
		switch p.As {
		case ACIJ:
			opcode = op_CIJ
		case ACGIJ, ACMPBEQ, ACMPBGE, ACMPBGT, ACMPBLE, ACMPBLT, ACMPBNE:
			opcode = op_CGIJ
		case ACLIJ:
			opcode = op_CLIJ
		case ACLGIJ, ACMPUBEQ, ACMPUBGE, ACMPUBGT, ACMPUBLE, ACMPUBLT, ACMPUBNE:
			opcode = op_CLGIJ
		}
		if int32(int16(v)) != v {
			
			
			
			
			
			
			
			
			m3 ^= 0xe 
			zRIE(_c, opcode, uint32(r1), m3, uint32(sizeRIE+sizeRIL)/2, 0, 0, 0, uint32(i2), asm)
			zRIL(_c, op_BRCL, uint32(Always), uint32(v-sizeRIE/2), asm)
		} else {
			zRIE(_c, opcode, uint32(r1), m3, uint32(v), 0, 0, 0, uint32(i2), asm)
		}

	case 91: 
		var opcode uint32
		switch p.As {
		case ATMHH:
			opcode = op_TMHH
		case ATMHL:
			opcode = op_TMHL
		case ATMLH:
			opcode = op_TMLH
		case ATMLL:
			opcode = op_TMLL
		}
		zRI(opcode, uint32(p.From.Reg), uint32(c.vregoff(&p.To)), asm)

	case 92: 
		zRRE(op_IPM, uint32(p.From.Reg), 0, asm)

	case 93: 
		v := c.vregoff(&p.To)
		if v != 0 {
			c.ctxt.Diag("invalid offset against GOT slot %v", p)
		}
		zRIL(_b, op_LGRL, uint32(p.To.Reg), 0, asm)
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc + 2)
		rel.Siz = 4
		rel.Sym = p.From.Sym
		rel.Type = objabi.R_GOTPCREL
		rel.Add = 2 + int64(rel.Siz)

	case 94: 
		zRIL(_b, op_LARL, regtmp(p), (sizeRIL+sizeRXY+sizeRI)>>1, asm)
		zRXY(op_LG, uint32(p.To.Reg), regtmp(p), 0, 0, asm)
		zRI(op_BRC, 0xF, (sizeRI+8)>>1, asm)
		*asm = append(*asm, 0, 0, 0, 0, 0, 0, 0, 0)
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc + sizeRIL + sizeRXY + sizeRI)
		rel.Siz = 8
		rel.Sym = p.From.Sym
		rel.Type = objabi.R_TLS_LE
		rel.Add = 0

	case 95: 
		
		
		
		
		
		
		
		
		

		
		zRIL(_b, op_LARL, regtmp(p), 0, asm)
		ieent := obj.Addrel(c.cursym)
		ieent.Off = int32(c.pc + 2)
		ieent.Siz = 4
		ieent.Sym = p.From.Sym
		ieent.Type = objabi.R_TLS_IE
		ieent.Add = 2 + int64(ieent.Siz)

		
		zRXY(op_LGF, uint32(p.To.Reg), regtmp(p), 0, 0, asm)
		
		

	case 96: 
		length := c.vregoff(&p.From)
		offset := c.vregoff(&p.To)
		reg := p.To.Reg
		if reg == 0 {
			reg = REGSP
		}
		if length <= 0 {
			c.ctxt.Diag("cannot CLEAR %d bytes, must be greater than 0", length)
		}
		for length > 0 {
			if offset < 0 || offset >= DISP12 {
				if offset >= -DISP20/2 && offset < DISP20/2 {
					zRXY(op_LAY, regtmp(p), uint32(reg), 0, uint32(offset), asm)
				} else {
					if reg != int16(regtmp(p)) {
						zRRE(op_LGR, regtmp(p), uint32(reg), asm)
					}
					zRIL(_a, op_AGFI, regtmp(p), uint32(offset), asm)
				}
				reg = int16(regtmp(p))
				offset = 0
			}
			size := length
			if size > 256 {
				size = 256
			}

			switch size {
			case 1:
				zSI(op_MVI, 0, uint32(reg), uint32(offset), asm)
			case 2:
				zSIL(op_MVHHI, uint32(reg), uint32(offset), 0, asm)
			case 4:
				zSIL(op_MVHI, uint32(reg), uint32(offset), 0, asm)
			case 8:
				zSIL(op_MVGHI, uint32(reg), uint32(offset), 0, asm)
			default:
				zSS(_a, op_XC, uint32(size-1), 0, uint32(reg), uint32(offset), uint32(reg), uint32(offset), asm)
			}

			length -= size
			offset += size
		}

	case 97: 
		rstart := p.From.Reg
		rend := p.Reg
		offset := c.regoff(&p.To)
		reg := p.To.Reg
		if reg == 0 {
			reg = REGSP
		}
		if offset < -DISP20/2 || offset >= DISP20/2 {
			if reg != int16(regtmp(p)) {
				zRRE(op_LGR, regtmp(p), uint32(reg), asm)
			}
			zRIL(_a, op_AGFI, regtmp(p), uint32(offset), asm)
			reg = int16(regtmp(p))
			offset = 0
		}
		switch p.As {
		case ASTMY:
			if offset >= 0 && offset < DISP12 {
				zRS(op_STM, uint32(rstart), uint32(rend), uint32(reg), uint32(offset), asm)
			} else {
				zRSY(op_STMY, uint32(rstart), uint32(rend), uint32(reg), uint32(offset), asm)
			}
		case ASTMG:
			zRSY(op_STMG, uint32(rstart), uint32(rend), uint32(reg), uint32(offset), asm)
		}

	case 98: 
		rstart := p.Reg
		rend := p.To.Reg
		offset := c.regoff(&p.From)
		reg := p.From.Reg
		if reg == 0 {
			reg = REGSP
		}
		if offset < -DISP20/2 || offset >= DISP20/2 {
			if reg != int16(regtmp(p)) {
				zRRE(op_LGR, regtmp(p), uint32(reg), asm)
			}
			zRIL(_a, op_AGFI, regtmp(p), uint32(offset), asm)
			reg = int16(regtmp(p))
			offset = 0
		}
		switch p.As {
		case ALMY:
			if offset >= 0 && offset < DISP12 {
				zRS(op_LM, uint32(rstart), uint32(rend), uint32(reg), uint32(offset), asm)
			} else {
				zRSY(op_LMY, uint32(rstart), uint32(rend), uint32(reg), uint32(offset), asm)
			}
		case ALMG:
			zRSY(op_LMG, uint32(rstart), uint32(rend), uint32(reg), uint32(offset), asm)
		}

	case 99: 
		if p.To.Index != 0 {
			c.ctxt.Diag("cannot use indexed address")
		}
		offset := c.regoff(&p.To)
		if offset < -DISP20/2 || offset >= DISP20/2 {
			c.ctxt.Diag("%v does not fit into 20-bit signed integer", offset)
		}
		var opcode uint32
		switch p.As {
		case ALAA:
			opcode = op_LAA
		case ALAAG:
			opcode = op_LAAG
		case ALAAL:
			opcode = op_LAAL
		case ALAALG:
			opcode = op_LAALG
		case ALAN:
			opcode = op_LAN
		case ALANG:
			opcode = op_LANG
		case ALAX:
			opcode = op_LAX
		case ALAXG:
			opcode = op_LAXG
		case ALAO:
			opcode = op_LAO
		case ALAOG:
			opcode = op_LAOG
		}
		zRSY(opcode, uint32(p.Reg), uint32(p.From.Reg), uint32(p.To.Reg), uint32(offset), asm)

	case 100: 
		op, m3, _ := vop(p.As)
		v1 := p.From.Reg
		if p.Reg != 0 {
			m3 = uint32(c.vregoff(&p.From))
			v1 = p.Reg
		}
		b2 := p.To.Reg
		if b2 == 0 {
			b2 = REGSP
		}
		d2 := uint32(c.vregoff(&p.To))
		zVRX(op, uint32(v1), uint32(p.To.Index), uint32(b2), d2, m3, asm)

	case 101: 
		op, m3, _ := vop(p.As)
		src := &p.From
		if p.GetFrom3() != nil {
			m3 = uint32(c.vregoff(&p.From))
			src = p.GetFrom3()
		}
		b2 := src.Reg
		if b2 == 0 {
			b2 = REGSP
		}
		d2 := uint32(c.vregoff(src))
		zVRX(op, uint32(p.To.Reg), uint32(src.Index), uint32(b2), d2, m3, asm)

	case 102: 
		op, _, _ := vop(p.As)
		m3 := uint32(c.vregoff(&p.From))
		b2 := p.To.Reg
		if b2 == 0 {
			b2 = REGSP
		}
		d2 := uint32(c.vregoff(&p.To))
		zVRV(op, uint32(p.Reg), uint32(p.To.Index), uint32(b2), d2, m3, asm)

	case 103: 
		op, _, _ := vop(p.As)
		m3 := uint32(c.vregoff(&p.From))
		b2 := p.GetFrom3().Reg
		if b2 == 0 {
			b2 = REGSP
		}
		d2 := uint32(c.vregoff(p.GetFrom3()))
		zVRV(op, uint32(p.To.Reg), uint32(p.GetFrom3().Index), uint32(b2), d2, m3, asm)

	case 104: 
		op, m4, _ := vop(p.As)
		fr := p.Reg
		if fr == 0 {
			fr = p.To.Reg
		}
		bits := uint32(c.vregoff(&p.From))
		zVRS(op, uint32(p.To.Reg), uint32(fr), uint32(p.From.Reg), bits, m4, asm)

	case 105: 
		op, _, _ := vop(p.As)
		offset := uint32(c.vregoff(&p.To))
		reg := p.To.Reg
		if reg == 0 {
			reg = REGSP
		}
		zVRS(op, uint32(p.From.Reg), uint32(p.Reg), uint32(reg), offset, 0, asm)

	case 106: 
		op, _, _ := vop(p.As)
		offset := uint32(c.vregoff(&p.From))
		reg := p.From.Reg
		if reg == 0 {
			reg = REGSP
		}
		zVRS(op, uint32(p.Reg), uint32(p.To.Reg), uint32(reg), offset, 0, asm)

	case 107: 
		op, _, _ := vop(p.As)
		offset := uint32(c.vregoff(&p.To))
		reg := p.To.Reg
		if reg == 0 {
			reg = REGSP
		}
		zVRS(op, uint32(p.Reg), uint32(p.From.Reg), uint32(reg), offset, 0, asm)

	case 108: 
		op, _, _ := vop(p.As)
		offset := uint32(c.vregoff(p.GetFrom3()))
		reg := p.GetFrom3().Reg
		if reg == 0 {
			reg = REGSP
		}
		zVRS(op, uint32(p.To.Reg), uint32(p.From.Reg), uint32(reg), offset, 0, asm)

	case 109: 
		op, m3, _ := vop(p.As)
		i2 := uint32(c.vregoff(&p.From))
		if p.GetFrom3() != nil {
			m3 = uint32(c.vregoff(&p.From))
			i2 = uint32(c.vregoff(p.GetFrom3()))
		}
		switch p.As {
		case AVZERO:
			i2 = 0
		case AVONE:
			i2 = 0xffff
		}
		zVRIa(op, uint32(p.To.Reg), i2, m3, asm)

	case 110:
		op, m4, _ := vop(p.As)
		i2 := uint32(c.vregoff(&p.From))
		i3 := uint32(c.vregoff(p.GetFrom3()))
		zVRIb(op, uint32(p.To.Reg), i2, i3, m4, asm)

	case 111:
		op, m4, _ := vop(p.As)
		i2 := uint32(c.vregoff(&p.From))
		zVRIc(op, uint32(p.To.Reg), uint32(p.Reg), i2, m4, asm)

	case 112:
		op, m5, _ := vop(p.As)
		i4 := uint32(c.vregoff(&p.From))
		zVRId(op, uint32(p.To.Reg), uint32(p.Reg), uint32(p.GetFrom3().Reg), i4, m5, asm)

	case 113:
		op, m4, _ := vop(p.As)
		m5 := singleElementMask(p.As)
		i3 := uint32(c.vregoff(&p.From))
		zVRIe(op, uint32(p.To.Reg), uint32(p.Reg), i3, m5, m4, asm)

	case 114: 
		op, m3, m5 := vop(p.As)
		m4 := singleElementMask(p.As)
		zVRRa(op, uint32(p.To.Reg), uint32(p.From.Reg), m5, m4, m3, asm)

	case 115: 
		op, m3, m5 := vop(p.As)
		m4 := singleElementMask(p.As)
		zVRRa(op, uint32(p.From.Reg), uint32(p.To.Reg), m5, m4, m3, asm)

	case 117: 
		op, m4, m5 := vop(p.As)
		zVRRb(op, uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), m5, m4, asm)

	case 118: 
		op, m4, m6 := vop(p.As)
		m5 := singleElementMask(p.As)
		v3 := p.Reg
		if v3 == 0 {
			v3 = p.To.Reg
		}
		zVRRc(op, uint32(p.To.Reg), uint32(p.From.Reg), uint32(v3), m6, m5, m4, asm)

	case 119: 
		op, m4, m6 := vop(p.As)
		m5 := singleElementMask(p.As)
		v2 := p.Reg
		if v2 == 0 {
			v2 = p.To.Reg
		}
		zVRRc(op, uint32(p.To.Reg), uint32(v2), uint32(p.From.Reg), m6, m5, m4, asm)

	case 120: 
		op, m6, _ := vop(p.As)
		m5 := singleElementMask(p.As)
		v1 := uint32(p.To.Reg)
		v2 := uint32(p.From.Reg)
		v3 := uint32(p.Reg)
		v4 := uint32(p.GetFrom3().Reg)
		zVRRd(op, v1, v2, v3, m6, m5, v4, asm)

	case 121: 
		op, m6, _ := vop(p.As)
		m5 := singleElementMask(p.As)
		v1 := uint32(p.To.Reg)
		v2 := uint32(p.From.Reg)
		v3 := uint32(p.Reg)
		v4 := uint32(p.GetFrom3().Reg)
		zVRRe(op, v1, v2, v3, m6, m5, v4, asm)

	case 122: 
		op, _, _ := vop(p.As)
		zVRRf(op, uint32(p.To.Reg), uint32(p.From.Reg), uint32(p.Reg), asm)

	case 123: 
		op, _, _ := vop(p.As)
		m4 := c.regoff(&p.From)
		zVRRc(op, uint32(p.To.Reg), uint32(p.Reg), uint32(p.GetFrom3().Reg), 0, 0, uint32(m4), asm)
	}
}

func (c *ctxtz) vregoff(a *obj.Addr) int64 {
	c.instoffset = 0
	if a != nil {
		c.aclass(a)
	}
	return c.instoffset
}

func (c *ctxtz) regoff(a *obj.Addr) int32 {
	return int32(c.vregoff(a))
}


func isU12(displacement int32) bool {
	return displacement >= 0 && displacement < DISP12
}


func (c *ctxtz) zopload12(a obj.As) (uint32, bool) {
	switch a {
	case AFMOVD:
		return op_LD, true
	case AFMOVS:
		return op_LE, true
	}
	return 0, false
}


func (c *ctxtz) zopload(a obj.As) uint32 {
	switch a {
	
	case AMOVD:
		return op_LG
	case AMOVW:
		return op_LGF
	case AMOVWZ:
		return op_LLGF
	case AMOVH:
		return op_LGH
	case AMOVHZ:
		return op_LLGH
	case AMOVB:
		return op_LGB
	case AMOVBZ:
		return op_LLGC

	
	case AFMOVD:
		return op_LDY
	case AFMOVS:
		return op_LEY

	
	case AMOVDBR:
		return op_LRVG
	case AMOVWBR:
		return op_LRV
	case AMOVHBR:
		return op_LRVH
	}

	c.ctxt.Diag("unknown store opcode %v", a)
	return 0
}


func (c *ctxtz) zopstore12(a obj.As) (uint32, bool) {
	switch a {
	case AFMOVD:
		return op_STD, true
	case AFMOVS:
		return op_STE, true
	case AMOVW, AMOVWZ:
		return op_ST, true
	case AMOVH, AMOVHZ:
		return op_STH, true
	case AMOVB, AMOVBZ:
		return op_STC, true
	}
	return 0, false
}


func (c *ctxtz) zopstore(a obj.As) uint32 {
	switch a {
	
	case AMOVD:
		return op_STG
	case AMOVW, AMOVWZ:
		return op_STY
	case AMOVH, AMOVHZ:
		return op_STHY
	case AMOVB, AMOVBZ:
		return op_STCY

	
	case AFMOVD:
		return op_STDY
	case AFMOVS:
		return op_STEY

	
	case AMOVDBR:
		return op_STRVG
	case AMOVWBR:
		return op_STRV
	case AMOVHBR:
		return op_STRVH
	}

	c.ctxt.Diag("unknown store opcode %v", a)
	return 0
}


func (c *ctxtz) zoprre(a obj.As) uint32 {
	switch a {
	case ACMP:
		return op_CGR
	case ACMPU:
		return op_CLGR
	case AFCMPO: 
		return op_KDBR
	case AFCMPU: 
		return op_CDBR
	case ACEBR:
		return op_CEBR
	}
	c.ctxt.Diag("unknown rre opcode %v", a)
	return 0
}


func (c *ctxtz) zoprr(a obj.As) uint32 {
	switch a {
	case ACMPW:
		return op_CR
	case ACMPWU:
		return op_CLR
	}
	c.ctxt.Diag("unknown rr opcode %v", a)
	return 0
}


func (c *ctxtz) zopril(a obj.As) uint32 {
	switch a {
	case ACMP:
		return op_CGFI
	case ACMPU:
		return op_CLGFI
	case ACMPW:
		return op_CFI
	case ACMPWU:
		return op_CLFI
	}
	c.ctxt.Diag("unknown ril opcode %v", a)
	return 0
}


const (
	sizeE    = 2
	sizeI    = 2
	sizeIE   = 4
	sizeMII  = 6
	sizeRI   = 4
	sizeRI1  = 4
	sizeRI2  = 4
	sizeRI3  = 4
	sizeRIE  = 6
	sizeRIE1 = 6
	sizeRIE2 = 6
	sizeRIE3 = 6
	sizeRIE4 = 6
	sizeRIE5 = 6
	sizeRIE6 = 6
	sizeRIL  = 6
	sizeRIL1 = 6
	sizeRIL2 = 6
	sizeRIL3 = 6
	sizeRIS  = 6
	sizeRR   = 2
	sizeRRD  = 4
	sizeRRE  = 4
	sizeRRF  = 4
	sizeRRF1 = 4
	sizeRRF2 = 4
	sizeRRF3 = 4
	sizeRRF4 = 4
	sizeRRF5 = 4
	sizeRRR  = 2
	sizeRRS  = 6
	sizeRS   = 4
	sizeRS1  = 4
	sizeRS2  = 4
	sizeRSI  = 4
	sizeRSL  = 6
	sizeRSY  = 6
	sizeRSY1 = 6
	sizeRSY2 = 6
	sizeRX   = 4
	sizeRX1  = 4
	sizeRX2  = 4
	sizeRXE  = 6
	sizeRXF  = 6
	sizeRXY  = 6
	sizeRXY1 = 6
	sizeRXY2 = 6
	sizeS    = 4
	sizeSI   = 4
	sizeSIL  = 6
	sizeSIY  = 6
	sizeSMI  = 6
	sizeSS   = 6
	sizeSS1  = 6
	sizeSS2  = 6
	sizeSS3  = 6
	sizeSS4  = 6
	sizeSS5  = 6
	sizeSS6  = 6
	sizeSSE  = 6
	sizeSSF  = 6
)


type form int

const (
	_a form = iota
	_b
	_c
	_d
	_e
	_f
)

func zE(op uint32, asm *[]byte) {
	*asm = append(*asm, uint8(op>>8), uint8(op))
}

func zI(op, i1 uint32, asm *[]byte) {
	*asm = append(*asm, uint8(op>>8), uint8(i1))
}

func zMII(op, m1, ri2, ri3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(m1)<<4)|uint8((ri2>>8)&0x0F),
		uint8(ri2),
		uint8(ri3>>16),
		uint8(ri3>>8),
		uint8(ri3))
}

func zRI(op, r1_m1, i2_ri2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1_m1)<<4)|(uint8(op)&0x0F),
		uint8(i2_ri2>>8),
		uint8(i2_ri2))
}












func zRIE(f form, op, r1, r2_m3_r3, i2_ri4_ri2, i3, i4, m3, i2_i5 uint32, asm *[]byte) {
	*asm = append(*asm, uint8(op>>8), uint8(r1)<<4|uint8(r2_m3_r3&0x0F))

	switch f {
	default:
		*asm = append(*asm, uint8(i2_ri4_ri2>>8), uint8(i2_ri4_ri2))
	case _f:
		*asm = append(*asm, uint8(i3), uint8(i4))
	}

	switch f {
	case _a, _b:
		*asm = append(*asm, uint8(m3)<<4)
	default:
		*asm = append(*asm, uint8(i2_i5))
	}

	*asm = append(*asm, uint8(op))
}

func zRIL(f form, op, r1_m1, i2_ri2 uint32, asm *[]byte) {
	if f == _a || f == _b {
		r1_m1 = r1_m1 - obj.RBaseS390X 
	}
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1_m1)<<4)|(uint8(op)&0x0F),
		uint8(i2_ri2>>24),
		uint8(i2_ri2>>16),
		uint8(i2_ri2>>8),
		uint8(i2_ri2))
}

func zRIS(op, r1, m3, b4, d4, i2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1)<<4)|uint8(m3&0x0F),
		(uint8(b4)<<4)|(uint8(d4>>8)&0x0F),
		uint8(d4),
		uint8(i2),
		uint8(op))
}

func zRR(op, r1, r2 uint32, asm *[]byte) {
	*asm = append(*asm, uint8(op>>8), (uint8(r1)<<4)|uint8(r2&0x0F))
}

func zRRD(op, r1, r3, r2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(op),
		uint8(r1)<<4,
		(uint8(r3)<<4)|uint8(r2&0x0F))
}

func zRRE(op, r1, r2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(op),
		0,
		(uint8(r1)<<4)|uint8(r2&0x0F))
}

func zRRF(op, r3_m3, m4, r1, r2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(op),
		(uint8(r3_m3)<<4)|uint8(m4&0x0F),
		(uint8(r1)<<4)|uint8(r2&0x0F))
}

func zRRS(op, r1, r2, b4, d4, m3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1)<<4)|uint8(r2&0x0F),
		(uint8(b4)<<4)|uint8((d4>>8)&0x0F),
		uint8(d4),
		uint8(m3)<<4,
		uint8(op))
}

func zRS(op, r1, r3_m3, b2, d2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1)<<4)|uint8(r3_m3&0x0F),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2))
}

func zRSI(op, r1, r3, ri2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1)<<4)|uint8(r3&0x0F),
		uint8(ri2>>8),
		uint8(ri2))
}

func zRSL(op, l1, b2, d2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(l1),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2),
		uint8(op))
}

func zRSY(op, r1, r3_m3, b2, d2 uint32, asm *[]byte) {
	dl2 := uint16(d2) & 0x0FFF
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1)<<4)|uint8(r3_m3&0x0F),
		(uint8(b2)<<4)|(uint8(dl2>>8)&0x0F),
		uint8(dl2),
		uint8(d2>>12),
		uint8(op))
}

func zRX(op, r1_m1, x2, b2, d2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1_m1)<<4)|uint8(x2&0x0F),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2))
}

func zRXE(op, r1, x2, b2, d2, m3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1)<<4)|uint8(x2&0x0F),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2),
		uint8(m3)<<4,
		uint8(op))
}

func zRXF(op, r3, x2, b2, d2, m1 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r3)<<4)|uint8(x2&0x0F),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2),
		uint8(m1)<<4,
		uint8(op))
}

func zRXY(op, r1_m1, x2, b2, d2 uint32, asm *[]byte) {
	dl2 := uint16(d2) & 0x0FFF
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r1_m1)<<4)|uint8(x2&0x0F),
		(uint8(b2)<<4)|(uint8(dl2>>8)&0x0F),
		uint8(dl2),
		uint8(d2>>12),
		uint8(op))
}

func zS(op, b2, d2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(op),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2))
}

func zSI(op, i2, b1, d1 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(i2),
		(uint8(b1)<<4)|uint8((d1>>8)&0x0F),
		uint8(d1))
}

func zSIL(op, b1, d1, i2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(op),
		(uint8(b1)<<4)|uint8((d1>>8)&0x0F),
		uint8(d1),
		uint8(i2>>8),
		uint8(i2))
}

func zSIY(op, i2, b1, d1 uint32, asm *[]byte) {
	dl1 := uint16(d1) & 0x0FFF
	*asm = append(*asm,
		uint8(op>>8),
		uint8(i2),
		(uint8(b1)<<4)|(uint8(dl1>>8)&0x0F),
		uint8(dl1),
		uint8(d1>>12),
		uint8(op))
}

func zSMI(op, m1, b3, d3, ri2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(m1)<<4,
		(uint8(b3)<<4)|uint8((d3>>8)&0x0F),
		uint8(d3),
		uint8(ri2>>8),
		uint8(ri2))
}











func zSS(f form, op, l1_r1, l2_i3_r3, b1_b2, d1_d2, b2_b4, d2_d4 uint32, asm *[]byte) {
	*asm = append(*asm, uint8(op>>8))

	switch f {
	case _a:
		*asm = append(*asm, uint8(l1_r1))
	case _b, _c, _d, _e:
		*asm = append(*asm, (uint8(l1_r1)<<4)|uint8(l2_i3_r3&0x0F))
	case _f:
		*asm = append(*asm, uint8(l2_i3_r3))
	}

	*asm = append(*asm,
		(uint8(b1_b2)<<4)|uint8((d1_d2>>8)&0x0F),
		uint8(d1_d2),
		(uint8(b2_b4)<<4)|uint8((d2_d4>>8)&0x0F),
		uint8(d2_d4))
}

func zSSE(op, b1, d1, b2, d2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(op),
		(uint8(b1)<<4)|uint8((d1>>8)&0x0F),
		uint8(d1),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2))
}

func zSSF(op, r3, b1, d1, b2, d2 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(r3)<<4)|(uint8(op)&0x0F),
		(uint8(b1)<<4)|uint8((d1>>8)&0x0F),
		uint8(d1),
		(uint8(b2)<<4)|uint8((d2>>8)&0x0F),
		uint8(d2))
}

func rxb(va, vb, vc, vd uint32) uint8 {
	mask := uint8(0)
	if va >= REG_V16 && va <= REG_V31 {
		mask |= 0x8
	}
	if vb >= REG_V16 && vb <= REG_V31 {
		mask |= 0x4
	}
	if vc >= REG_V16 && vc <= REG_V31 {
		mask |= 0x2
	}
	if vd >= REG_V16 && vd <= REG_V31 {
		mask |= 0x1
	}
	return mask
}

func zVRX(op, v1, x2, b2, d2, m3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(x2)&0xf),
		(uint8(b2)<<4)|(uint8(d2>>8)&0xf),
		uint8(d2),
		(uint8(m3)<<4)|rxb(v1, 0, 0, 0),
		uint8(op))
}

func zVRV(op, v1, v2, b2, d2, m3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		(uint8(b2)<<4)|(uint8(d2>>8)&0xf),
		uint8(d2),
		(uint8(m3)<<4)|rxb(v1, v2, 0, 0),
		uint8(op))
}

func zVRS(op, v1, v3_r3, b2, d2, m4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v3_r3)&0xf),
		(uint8(b2)<<4)|(uint8(d2>>8)&0xf),
		uint8(d2),
		(uint8(m4)<<4)|rxb(v1, v3_r3, 0, 0),
		uint8(op))
}

func zVRRa(op, v1, v2, m5, m4, m3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		0,
		(uint8(m5)<<4)|(uint8(m4)&0xf),
		(uint8(m3)<<4)|rxb(v1, v2, 0, 0),
		uint8(op))
}

func zVRRb(op, v1, v2, v3, m5, m4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		uint8(v3)<<4,
		uint8(m5)<<4,
		(uint8(m4)<<4)|rxb(v1, v2, v3, 0),
		uint8(op))
}

func zVRRc(op, v1, v2, v3, m6, m5, m4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		uint8(v3)<<4,
		(uint8(m6)<<4)|(uint8(m5)&0xf),
		(uint8(m4)<<4)|rxb(v1, v2, v3, 0),
		uint8(op))
}

func zVRRd(op, v1, v2, v3, m5, m6, v4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		(uint8(v3)<<4)|(uint8(m5)&0xf),
		uint8(m6)<<4,
		(uint8(v4)<<4)|rxb(v1, v2, v3, v4),
		uint8(op))
}

func zVRRe(op, v1, v2, v3, m6, m5, v4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		(uint8(v3)<<4)|(uint8(m6)&0xf),
		uint8(m5),
		(uint8(v4)<<4)|rxb(v1, v2, v3, v4),
		uint8(op))
}

func zVRRf(op, v1, r2, r3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(r2)&0xf),
		uint8(r3)<<4,
		0,
		rxb(v1, 0, 0, 0),
		uint8(op))
}

func zVRIa(op, v1, i2, m3 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(v1)<<4,
		uint8(i2>>8),
		uint8(i2),
		(uint8(m3)<<4)|rxb(v1, 0, 0, 0),
		uint8(op))
}

func zVRIb(op, v1, i2, i3, m4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		uint8(v1)<<4,
		uint8(i2),
		uint8(i3),
		(uint8(m4)<<4)|rxb(v1, 0, 0, 0),
		uint8(op))
}

func zVRIc(op, v1, v3, i2, m4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v3)&0xf),
		uint8(i2>>8),
		uint8(i2),
		(uint8(m4)<<4)|rxb(v1, v3, 0, 0),
		uint8(op))
}

func zVRId(op, v1, v2, v3, i4, m5 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		uint8(v3)<<4,
		uint8(i4),
		(uint8(m5)<<4)|rxb(v1, v2, v3, 0),
		uint8(op))
}

func zVRIe(op, v1, v2, i3, m5, m4 uint32, asm *[]byte) {
	*asm = append(*asm,
		uint8(op>>8),
		(uint8(v1)<<4)|(uint8(v2)&0xf),
		uint8(i3>>4),
		(uint8(i3)<<4)|(uint8(m5)&0xf),
		(uint8(m4)<<4)|rxb(v1, v2, 0, 0),
		uint8(op))
}
