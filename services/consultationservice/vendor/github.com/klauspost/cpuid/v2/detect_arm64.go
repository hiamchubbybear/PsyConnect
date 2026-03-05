

//go:build arm64 && !gccgo && !noasm && !appengine
// +build arm64,!gccgo,!noasm,!appengine

package cpuid

import "runtime"

func getMidr() (midr uint64)
func getProcFeatures() (procFeatures uint64)
func getInstAttributes() (instAttrReg0, instAttrReg1 uint64)
func getVectorLength() (vl, pl uint64)

func initCPU() {
	cpuid = func(uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
	cpuidex = func(x, y uint32) (a, b, c, d uint32) { return 0, 0, 0, 0 }
	xgetbv = func(uint32) (a, b uint32) { return 0, 0 }
	rdtscpAsm = func() (a, b, c, d uint32) { return 0, 0, 0, 0 }
}

func addInfo(c *CPUInfo, safe bool) {
	
	c.CacheLine = 64
	detectOS(c)

	
	if safe && !c.Has(ARMCPUID) && runtime.GOOS != "freebsd" {
		return
	}
	midr := getMidr()

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	switch (midr >> 24) & 0xff {
	case 0xC0:
		c.VendorString = "Ampere Computing"
		c.VendorID = Ampere
	case 0x41:
		c.VendorString = "Arm Limited"
		c.VendorID = ARM
	case 0x42:
		c.VendorString = "Broadcom Corporation"
		c.VendorID = Broadcom
	case 0x43:
		c.VendorString = "Cavium Inc"
		c.VendorID = Cavium
	case 0x44:
		c.VendorString = "Digital Equipment Corporation"
		c.VendorID = DEC
	case 0x46:
		c.VendorString = "Fujitsu Ltd"
		c.VendorID = Fujitsu
	case 0x49:
		c.VendorString = "Infineon Technologies AG"
		c.VendorID = Infineon
	case 0x4D:
		c.VendorString = "Motorola or Freescale Semiconductor Inc"
		c.VendorID = Motorola
	case 0x4E:
		c.VendorString = "NVIDIA Corporation"
		c.VendorID = NVIDIA
	case 0x50:
		c.VendorString = "Applied Micro Circuits Corporation"
		c.VendorID = AMCC
	case 0x51:
		c.VendorString = "Qualcomm Inc"
		c.VendorID = Qualcomm
	case 0x56:
		c.VendorString = "Marvell International Ltd"
		c.VendorID = Marvell
	case 0x69:
		c.VendorString = "Intel Corporation"
		c.VendorID = Intel
	}

	
	
	
	
	
	
	
	
	
	
	
	
	
	c.Family = int(midr>>16) & 0xff

	
	
	
	
	
	
	c.Model = int(midr) & 0xffff

	procFeatures := getProcFeatures()

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	var f flagSet
	
	
	
	f.setIf(procFeatures&(0xf<<32) != 0, SVE)
	if procFeatures&(0xf<<20) != 15<<20 {
		f.set(ASIMD)
		
		
		f.setIf(procFeatures&(0xf<<20) == 1<<20, FPHP, ASIMDHP)
	}
	f.setIf(procFeatures&(0xf<<16) != 0, FP)

	instAttrReg0, instAttrReg1 := getInstAttributes()

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	f.setIf(instAttrReg0&(0xf<<60) != 0, RNDR)
	f.setIf(instAttrReg0&(0xf<<56) != 0, TLB)
	f.setIf(instAttrReg0&(0xf<<52) != 0, TS)
	f.setIf(instAttrReg0&(0xf<<48) != 0, FHM)
	f.setIf(instAttrReg0&(0xf<<44) != 0, ASIMDDP)
	f.setIf(instAttrReg0&(0xf<<40) != 0, SM4)
	f.setIf(instAttrReg0&(0xf<<36) != 0, SM3)
	f.setIf(instAttrReg0&(0xf<<32) != 0, SHA3)
	f.setIf(instAttrReg0&(0xf<<28) != 0, ASIMDRDM)
	f.setIf(instAttrReg0&(0xf<<20) != 0, ATOMICS)
	f.setIf(instAttrReg0&(0xf<<16) != 0, CRC32)
	f.setIf(instAttrReg0&(0xf<<12) != 0, SHA2)
	
	
	f.setIf(instAttrReg0&(0xf<<12) == 2<<12, SHA512)
	f.setIf(instAttrReg0&(0xf<<8) != 0, SHA1)
	f.setIf(instAttrReg0&(0xf<<4) != 0, AESARM)
	
	
	f.setIf(instAttrReg0&(0xf<<4) == 2<<4, PMULL)

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	
	
	
	f.setIf(instAttrReg1&(0xf<<28) != 24, GPA)
	f.setIf(instAttrReg1&(0xf<<20) != 0, LRCPC)
	f.setIf(instAttrReg1&(0xf<<16) != 0, FCMA)
	f.setIf(instAttrReg1&(0xf<<12) != 0, JSCVT)
	
	
	
	
	
	
	f.setIf(instAttrReg1&(0xf<<0) != 0, DCPOP)

	
	c.featureSet.or(f)
}
