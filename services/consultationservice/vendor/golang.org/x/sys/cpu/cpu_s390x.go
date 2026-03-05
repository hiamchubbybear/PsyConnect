



package cpu

const cacheLineSize = 256

func initOptions() {
	options = []option{
		{Name: "zarch", Feature: &S390X.HasZARCH, Required: true},
		{Name: "stfle", Feature: &S390X.HasSTFLE, Required: true},
		{Name: "ldisp", Feature: &S390X.HasLDISP, Required: true},
		{Name: "eimm", Feature: &S390X.HasEIMM, Required: true},
		{Name: "dfp", Feature: &S390X.HasDFP},
		{Name: "etf3eh", Feature: &S390X.HasETF3EH},
		{Name: "msa", Feature: &S390X.HasMSA},
		{Name: "aes", Feature: &S390X.HasAES},
		{Name: "aescbc", Feature: &S390X.HasAESCBC},
		{Name: "aesctr", Feature: &S390X.HasAESCTR},
		{Name: "aesgcm", Feature: &S390X.HasAESGCM},
		{Name: "ghash", Feature: &S390X.HasGHASH},
		{Name: "sha1", Feature: &S390X.HasSHA1},
		{Name: "sha256", Feature: &S390X.HasSHA256},
		{Name: "sha3", Feature: &S390X.HasSHA3},
		{Name: "sha512", Feature: &S390X.HasSHA512},
		{Name: "vx", Feature: &S390X.HasVX},
		{Name: "vxe", Feature: &S390X.HasVXE},
	}
}



func bitIsSet(bits []uint64, index uint) bool {
	return bits[index/64]&((1<<63)>>(index%64)) != 0
}


type facility uint8

const (
	
	zarch  facility = 1  
	stflef facility = 7  
	ldisp  facility = 18 
	eimm   facility = 21 

	
	dfp    facility = 42 
	etf3eh facility = 30 

	
	msa  facility = 17  
	msa3 facility = 76  
	msa4 facility = 77  
	msa5 facility = 57  
	msa8 facility = 146 
	msa9 facility = 155 

	
	vx   facility = 129 
	vxe  facility = 135 
	vxe2 facility = 148 
)




type facilityList struct {
	bits [4]uint64
}


func (s *facilityList) Has(fs ...facility) bool {
	if len(fs) == 0 {
		panic("no facility bits provided")
	}
	for _, f := range fs {
		if !bitIsSet(s.bits[:], uint(f)) {
			return false
		}
	}
	return true
}


type function uint8

const (
	
	aes128 function = 18 
	aes192 function = 19 
	aes256 function = 20 

	
	sha1     function = 1  
	sha256   function = 2  
	sha512   function = 3  
	sha3_224 function = 32 
	sha3_256 function = 33 
	sha3_384 function = 34 
	sha3_512 function = 35 
	shake128 function = 36 
	shake256 function = 37 

	
	ghash function = 65 
)




type queryResult struct {
	bits [2]uint64
}


func (q *queryResult) Has(fns ...function) bool {
	if len(fns) == 0 {
		panic("no function codes provided")
	}
	for _, f := range fns {
		if !bitIsSet(q.bits[:], uint(f)) {
			return false
		}
	}
	return true
}

func doinit() {
	initS390Xbase()

	
	
	if !haveAsmFunctions() {
		return
	}

	
	if S390X.HasMSA {
		aes := []function{aes128, aes192, aes256}

		
		km, kmc := kmQuery(), kmcQuery()
		S390X.HasAES = km.Has(aes...)
		S390X.HasAESCBC = kmc.Has(aes...)
		if S390X.HasSTFLE {
			facilities := stfle()
			if facilities.Has(msa4) {
				kmctr := kmctrQuery()
				S390X.HasAESCTR = kmctr.Has(aes...)
			}
			if facilities.Has(msa8) {
				kma := kmaQuery()
				S390X.HasAESGCM = kma.Has(aes...)
			}
		}

		
		kimd := kimdQuery() 
		klmd := klmdQuery() 
		S390X.HasSHA1 = kimd.Has(sha1) && klmd.Has(sha1)
		S390X.HasSHA256 = kimd.Has(sha256) && klmd.Has(sha256)
		S390X.HasSHA512 = kimd.Has(sha512) && klmd.Has(sha512)
		S390X.HasGHASH = kimd.Has(ghash) 
		sha3 := []function{
			sha3_224, sha3_256, sha3_384, sha3_512,
			shake128, shake256,
		}
		S390X.HasSHA3 = kimd.Has(sha3...) && klmd.Has(sha3...)
	}
}
