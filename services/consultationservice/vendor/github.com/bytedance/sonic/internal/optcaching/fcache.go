

package caching

import (
	"strings"
	"unicode"
	"unsafe"

	"github.com/bytedance/sonic/internal/envs"
	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/resolver"
	"github.com/bytedance/sonic/internal/rt"
)

const _AlignSize =  32
const _PaddingSize =  32

type FieldLookup interface {
	Set(fields []resolver.FieldMeta)
	Get(name string, caseSensitive bool) int
}

func isAscii(s string) bool {
	for i :=0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func NewFieldLookup(fields []resolver.FieldMeta) FieldLookup {
	var f FieldLookup
	isAsc := true
	n := len(fields)

	
	for _, f := range fields {
		if !isAscii(f.Name) {
			isAsc = false
			break
		}
	}

	if n <= 8 {
		f =  NewSmallFieldMap(n)
	} else if envs.UseFastMap && n <= 128 && isAsc {
		f =   NewNormalFieldMap(n)
	} else {
		f =   NewFallbackFieldMap(n)
	}

	f.Set(fields)
	return f
}


type SmallFieldMap struct {
	keys []string
	lowerKeys []string
}

func NewSmallFieldMap (hint int) *SmallFieldMap {
	return &SmallFieldMap{
		keys: make([]string, hint, hint),
		lowerKeys: make([]string, hint, hint),
	}
}

func (self *SmallFieldMap) Set(fields []resolver.FieldMeta) {
	if len(fields) > 8 {
		panic("small field map should use in small struct")
	}

	for i, f := range fields {
		self.keys[i] = f.Name
		self.lowerKeys[i] = strings.ToLower(f.Name)
	}
}

func (self *SmallFieldMap) Get(name string, caseSensitive bool) int {
	for i, k := range self.keys {
		if len(k) == len(name) && k == name {
			return i
		}
	}
	if caseSensitive {
		return -1
	}
	name = strings.ToLower(name)
	for i, k := range self.lowerKeys {
		if len(k) == len(name) && k == name {
			return i
		}
	}
	return -1
}












type NormalFieldMap struct {
	keys  			[]byte
	longKeys		[]keyEntry
	
	lowOffset	    int
}

type keyEntry struct {
	key 		string
	lowerKey	string
	index		uint
}

func NewNormalFieldMap(n int) *NormalFieldMap {
	return &NormalFieldMap{
	}
}

const _HdrSlot = 33
const _HdrSize = _HdrSlot * 5


func (self *NormalFieldMap) Get(name string, caseSensitive bool) int {
	
	if len(name) <= 32 {
		_ = native.LookupSmallKey
		lowOffset := self.lowOffset
		if caseSensitive {
			lowOffset = -1
		}
		return native.LookupSmallKey(&name, &self.keys, lowOffset);
	}
	return self.getLongKey(name, caseSensitive)
}

func (self *NormalFieldMap) getLongKey(name string, caseSensitive bool) int {
	for _, k := range self.longKeys {
		if len(k.key) != len(name) {
			continue;
		}
		if k.key == name {
			return int(k.index)
		}
	}

	if caseSensitive {
		return -1
	}

	lower := strings.ToLower(name)
	for _, k := range self.longKeys {
		if len(k.key) != len(name) {
			continue;
		}

		if k.lowerKey == lower {
			return int(k.index)
		}
	}
	return -1
}

func (self *NormalFieldMap) Getdouble(name string) int {
	if len(name) > 32 {
		for _, k := range self.longKeys {
			if len(k.key) != len(name) {
				continue;
			}
			if k.key == name {
				return int(k.index)
			}
		}
		return self.getCaseInsensitive(name)
	}

	
	cnt := int(self.keys[5 * len(name)])
	if cnt == 0 {
		return -1
	}
	p := ((*rt.GoSlice)(unsafe.Pointer(&self.keys))).Ptr
	offset := int(*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(5 * len(name) + 1)))) + _HdrSize
	for i := 0; i < cnt; i++ {
		key := rt.Mem2Str(self.keys[offset: offset + len(name)])
		if key == name {
			return int(self.keys[offset + len(name)])
		}
		offset += len(name) + 1
	}

	return self.getCaseInsensitive(name)
}

func (self *NormalFieldMap) getCaseInsensitive(name string) int {
	lower := strings.ToLower(name)
	if len(name) > 32 {
		for _, k := range self.longKeys {
			if len(k.key) != len(name) {
				continue;
			}

			if k.lowerKey == lower {
				return int(k.index)
			}
		}
		return -1
	}

	cnt := int(self.keys[5 * len(name)])
	p := ((*rt.GoSlice)(unsafe.Pointer(&self.keys))).Ptr
	offset := int(*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(5 * len(name) + 1)))) + self.lowOffset
	for i := 0; i < cnt; i++ {
		key := rt.Mem2Str(self.keys[offset: offset + len(name)])
		if key == lower {
			return int(self.keys[offset + len(name)])
		}
		offset += len(name) + 1
	}

	return -1
}

type keysInfo struct {
	counts int
	lenSum int
	offset int
	cur    int
}

func (self *NormalFieldMap) Set(fields []resolver.FieldMeta) {
	if len(fields) <=8 || len(fields) > 128 {
		panic("normal field map should use in small struct")
	}

	
	var keyLenSum [_HdrSlot]keysInfo

	for i := 0; i < _HdrSlot; i++ {
		keyLenSum[i].offset = 0
		keyLenSum[i].counts = 0
		keyLenSum[i].lenSum = 0
		keyLenSum[i].cur = 0
	}

	kvLen := 0
	for _, f := range(fields) {
		len := len(f.Name)
		if len <= 32 {
			kvLen += len + 1 
			keyLenSum[len].counts++
			keyLenSum[len].lenSum += len + 1
		}

	}

	
	self.keys = make([]byte, _HdrSize + 2 * kvLen, _HdrSize + 2 * kvLen + _PaddingSize)
	self.lowOffset = _HdrSize + kvLen

	
	self.keys[0] = byte(keyLenSum[0].counts)
	
	i := 1
	p := ((*rt.GoSlice)(unsafe.Pointer(&self.keys))).Ptr
	for i < _HdrSlot {
		keyLenSum[i].offset = keyLenSum[i-1].offset + keyLenSum[i-1].lenSum
		self.keys[i * 5] = byte(keyLenSum[i].counts)
		
		*(*int32)(unsafe.Pointer(uintptr(p) + uintptr(i * 5 + 1))) = int32(keyLenSum[i].offset)
		i += 1

	}

	
	for i, f := range(fields) {
		len := len(f.Name)
		if len <= 32 {
			offset := keyLenSum[len].offset +  keyLenSum[len].cur
			copy(self.keys[_HdrSize + offset: ], f.Name)
			copy(self.keys[self.lowOffset + offset: ], strings.ToLower(f.Name))
			self.keys[_HdrSize + offset + len] = byte(i)
			self.keys[self.lowOffset + offset + len] = byte(i)
			keyLenSum[len].cur += len + 1

		} else {
			self.longKeys = append(self.longKeys, keyEntry{f.Name, strings.ToLower(f.Name), uint(i)})
		}
	}

}


type FallbackFieldMap struct {
	oders  []string
	inner  map[string]int
	backup map[string]int
}
 
 func NewFallbackFieldMap(n int) *FallbackFieldMap {
	 return &FallbackFieldMap{
		 oders:  make([]string, n, n),
		 inner:  make(map[string]int, n*2),
		 backup: make(map[string]int, n*2),
	 }
 }
 
 func (self *FallbackFieldMap) Get(name string, caseSensitive bool) int {
	 if i, ok := self.inner[name]; ok {
		 return i
	 } else if !caseSensitive {
		 return self.getCaseInsensitive(name)
	 } else {
		return -1
	 }
 }
 
 func (self *FallbackFieldMap) Set(fields []resolver.FieldMeta) {

	for i, f := range(fields) {
		name := f.Name
		self.oders[i] = name
		self.inner[name] = i
	
		
		key := strings.ToLower(name)
		if v, ok := self.backup[key]; !ok || i < v {
			self.backup[key] = i
		}
	}
 }
 
 func (self *FallbackFieldMap) getCaseInsensitive(name string) int {
	 if i, ok := self.backup[strings.ToLower(name)]; ok {
		 return i
	 } else {
		 return -1
	 }
 }
 