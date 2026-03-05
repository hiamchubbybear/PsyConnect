


package codec

import (
	"math"
	"strconv"
)











func parseFloat32(b []byte) (f float32, err error) {
	return parseFloat32_custom(b)
}

func parseFloat64(b []byte) (f float64, err error) {
	return parseFloat64_custom(b)
}

func parseFloat32_strconv(b []byte) (f float32, err error) {
	f64, err := strconv.ParseFloat(stringView(b), 32)
	f = float32(f64)
	return
}

func parseFloat64_strconv(b []byte) (f float64, err error) {
	return strconv.ParseFloat(stringView(b), 64)
}



































const (
	fMaxMultiplierForExactPow10_64 = 1e15
	fMaxMultiplierForExactPow10_32 = 1e7

	fUint64Cutoff = (1<<64-1)/10 + 1
	

	fBase = 10
)

const (
	thousand    = 1000
	million     = thousand * thousand
	billion     = thousand * million
	trillion    = thousand * billion
	quadrillion = thousand * trillion
	quintillion = thousand * quadrillion
)


var uint64pow10 = [...]uint64{
	1, 10, 100,
	1 * thousand, 10 * thousand, 100 * thousand,
	1 * million, 10 * million, 100 * million,
	1 * billion, 10 * billion, 100 * billion,
	1 * trillion, 10 * trillion, 100 * trillion,
	1 * quadrillion, 10 * quadrillion, 100 * quadrillion,
	1 * quintillion, 10 * quintillion,
}
var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}
var float32pow10 = [...]float32{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10,
}

type floatinfo struct {
	mantbits uint8

	
	
	

	exactPow10 int8 

	exactInts int8 

	

	mantCutoffIsUint64Cutoff bool

	mantCutoff uint64
}

var fi32 = floatinfo{23, 10, 7, false, 1<<23 - 1}
var fi64 = floatinfo{52, 22, 15, false, 1<<52 - 1}

var fi64u = floatinfo{0, 19, 0, true, fUint64Cutoff}

func noFrac64(fbits uint64) bool {
	if fbits == 0 {
		return true
	}

	exp := uint64(fbits>>52)&0x7FF - 1023 
	
	return exp < 52 && fbits<<(12+exp) == 0 
}

func noFrac32(fbits uint32) bool {
	if fbits == 0 {
		return true
	}

	exp := uint32(fbits>>23)&0xFF - 127 
	
	return exp < 23 && fbits<<(9+exp) == 0 
}

func strconvParseErr(b []byte, fn string) error {
	return &strconv.NumError{
		Func: fn,
		Err:  strconv.ErrSyntax,
		Num:  string(b),
	}
}

func parseFloat32_reader(r readFloatResult) (f float32, fail bool) {
	f = float32(r.mantissa)
	if r.exp == 0 {
	} else if r.exp < 0 { 
		f /= float32pow10[uint8(-r.exp)]
	} else { 
		if r.exp > fi32.exactPow10 {
			f *= float32pow10[r.exp-fi32.exactPow10]
			if f > fMaxMultiplierForExactPow10_32 { 
				fail = true
				return 
			}
			f *= float32pow10[fi32.exactPow10]
		} else {
			f *= float32pow10[uint8(r.exp)]
		}
	}
	if r.neg {
		f = -f
	}
	return
}

func parseFloat32_custom(b []byte) (f float32, err error) {
	r := readFloat(b, fi32)
	if r.bad {
		return 0, strconvParseErr(b, "ParseFloat")
	}
	if r.ok {
		f, r.bad = parseFloat32_reader(r)
		if !r.bad {
			return
		}
	}
	return parseFloat32_strconv(b)
}

func parseFloat64_reader(r readFloatResult) (f float64, fail bool) {
	f = float64(r.mantissa)
	if r.exp == 0 {
	} else if r.exp < 0 { 
		f /= float64pow10[-uint8(r.exp)]
	} else { 
		if r.exp > fi64.exactPow10 {
			f *= float64pow10[r.exp-fi64.exactPow10]
			if f > fMaxMultiplierForExactPow10_64 { 
				fail = true
				return
			}
			f *= float64pow10[fi64.exactPow10]
		} else {
			f *= float64pow10[uint8(r.exp)]
		}
	}
	if r.neg {
		f = -f
	}
	return
}

func parseFloat64_custom(b []byte) (f float64, err error) {
	r := readFloat(b, fi64)
	if r.bad {
		return 0, strconvParseErr(b, "ParseFloat")
	}
	if r.ok {
		f, r.bad = parseFloat64_reader(r)
		if !r.bad {
			return
		}
	}
	return parseFloat64_strconv(b)
}

func parseUint64_simple(b []byte) (n uint64, ok bool) {
	var i int
	var n1 uint64
	var c uint8
LOOP:
	if i < len(b) {
		c = b[i]
		
		
		
		if n >= fUint64Cutoff || c < '0' || c > '9' {
			return
		} else if c == '0' {
			n *= fBase
		} else {
			n1 = n
			n = n*fBase + uint64(c-'0')
			if n < n1 {
				return
			}
		}
		i++
		goto LOOP
	}
	ok = true
	return
}

func parseUint64_reader(r readFloatResult) (f uint64, fail bool) {
	f = r.mantissa
	if r.exp == 0 {
	} else if r.exp < 0 { 
		if f%uint64pow10[uint8(-r.exp)] != 0 {
			fail = true
		} else {
			f /= uint64pow10[uint8(-r.exp)]
		}
	} else { 
		f *= uint64pow10[uint8(r.exp)]
	}
	return
}

func parseInteger_bytes(b []byte) (u uint64, neg, ok bool) {
	if len(b) == 0 {
		ok = true
		return
	}
	if b[0] == '-' {
		if len(b) == 1 {
			return
		}
		neg = true
		b = b[1:]
	}

	u, ok = parseUint64_simple(b)
	if ok {
		return
	}

	r := readFloat(b, fi64u)
	if r.ok {
		var fail bool
		u, fail = parseUint64_reader(r)
		if fail {
			f, err := parseFloat64(b)
			if err != nil {
				return
			}
			if !noFrac64(math.Float64bits(f)) {
				return
			}
			u = uint64(f)
		}
		ok = true
		return
	}
	return
}



func parseNumber(b []byte, z *fauxUnion, preferSignedInt bool) (err error) {
	var ok, neg bool
	var f uint64

	if len(b) == 0 {
		return
	}

	if b[0] == '-' {
		neg = true
		f, ok = parseUint64_simple(b[1:])
	} else {
		f, ok = parseUint64_simple(b)
	}

	if ok {
		if neg {
			z.v = valueTypeInt
			if chkOvf.Uint2Int(f, neg) {
				return strconvParseErr(b, "ParseInt")
			}
			z.i = -int64(f)
		} else if preferSignedInt {
			z.v = valueTypeInt
			if chkOvf.Uint2Int(f, neg) {
				return strconvParseErr(b, "ParseInt")
			}
			z.i = int64(f)
		} else {
			z.v = valueTypeUint
			z.u = f
		}
		return
	}

	z.v = valueTypeFloat
	z.f, err = parseFloat64_custom(b)
	return
}

type readFloatResult struct {
	mantissa uint64
	exp      int8
	neg      bool
	trunc    bool
	bad      bool 
	hardexp  bool 
	ok       bool
	
	
	
}

func readFloat(s []byte, y floatinfo) (r readFloatResult) {
	var i uint 
	var slen = uint(len(s))
	if slen == 0 {
		
		
		r.ok = true
		return
	}

	if s[0] == '-' {
		r.neg = true
		i++
	}

	
	

	var nd, ndMant, dp int8
	var sawdot, sawexp bool
	var xu uint64

LOOP:
	for ; i < slen; i++ {
		switch s[i] {
		case '.':
			if sawdot {
				r.bad = true
				return
			}
			sawdot = true
			dp = nd
		case 'e', 'E':
			sawexp = true
			break LOOP
		case '0':
			if nd == 0 {
				dp--
				continue LOOP
			}
			nd++
			if r.mantissa < y.mantCutoff {
				r.mantissa *= fBase
				ndMant++
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			nd++
			if y.mantCutoffIsUint64Cutoff && r.mantissa < fUint64Cutoff {
				r.mantissa *= fBase
				xu = r.mantissa + uint64(s[i]-'0')
				if xu < r.mantissa {
					r.trunc = true
					return
				}
				r.mantissa = xu
			} else if r.mantissa < y.mantCutoff {
				
				r.mantissa = r.mantissa*fBase + uint64(s[i]-'0')
			} else {
				r.trunc = true
				return
			}
			ndMant++
		default:
			r.bad = true
			return
		}
	}

	if !sawdot {
		dp = nd
	}

	if sawexp {
		i++
		if i < slen {
			var eneg bool
			if s[i] == '+' {
				i++
			} else if s[i] == '-' {
				i++
				eneg = true
			}
			if i < slen {
				
				
				if i+2 < slen {
					r.hardexp = true
					return
				}
				var e int8
				if s[i] < '0' || s[i] > '9' { 
					r.bad = true
					return
				}
				e = int8(s[i] - '0')
				i++
				if i < slen {
					if s[i] < '0' || s[i] > '9' { 
						r.bad = true
						return
					}
					e = e*fBase + int8(s[i]-'0') 
					i++
				}
				if eneg {
					dp -= e
				} else {
					dp += e
				}
			}
		}
	}

	if r.mantissa != 0 {
		r.exp = dp - ndMant
		
		if r.exp < -y.exactPow10 ||
			r.exp > y.exactInts+y.exactPow10 ||
			(y.mantbits != 0 && r.mantissa>>y.mantbits != 0) {
			r.hardexp = true
			return
		}
	}

	r.ok = true
	return
}
