

package rt

import (
    `fmt`
    `strings`
    `unsafe`

)

type Bitmap struct {
    N int
    B []byte
}

func (self *Bitmap) grow() {
    if self.N >= len(self.B) * 8 {
        self.B = append(self.B, 0)
    }
}

func (self *Bitmap) mark(i int, bv int) {
    if bv != 0 {
        self.B[i / 8] |= 1 << (i % 8)
    } else {
        self.B[i / 8] &^= 1 << (i % 8)
    }
}

func (self *Bitmap) Set(i int, bv int) {
    if i >= self.N {
        panic("bitmap: invalid bit position")
    } else {
        self.mark(i, bv)
    }
}

func (self *Bitmap) Append(bv int) {
    self.grow()
    self.mark(self.N, bv)
    self.N++
}

func (self *Bitmap) AppendMany(n int, bv int) {
    for i := 0; i < n; i++ {
        self.Append(bv)
    }
}

func (b *Bitmap) String() string {
    var buf strings.Builder
    for _, byteVal := range b.B {
        fmt.Fprintf(&buf, "%08b ", byteVal)
    }
    return fmt.Sprintf("Bits: %s(total %d bits, %d bytes)", buf.String(), b.N, len(b.B))
}

type BitVec struct {
    N uintptr
    B unsafe.Pointer
}

func (self BitVec) Bit(i uintptr) byte {
    return (*(*byte)(unsafe.Pointer(uintptr(self.B) + i / 8)) >> (i % 8)) & 1
}

func (self BitVec) String() string {
    var i uintptr
    var v []string

    
    for i = 0; i < self.N; i++ {
        v = append(v, fmt.Sprintf("%d", self.Bit(i)))
    }

    
    return fmt.Sprintf(
        "BitVec { %s }",
        strings.Join(v, ", "),
    )
}



type StackMap struct {
    
    N int32
    
    L int32
    
    B [1]byte
}

func (self *StackMap) Get(i int) BitVec {
    return BitVec {
        N: uintptr(self.L),
        B: unsafe.Pointer(uintptr(unsafe.Pointer(&self.B)) + uintptr(i * self.BitmapBytes())),
    }
}

func (self *StackMap) String() string {
    sb := strings.Builder{}
    sb.WriteString("StackMap {")

    
    for i := 0; i < int(self.N); i++ {
        sb.WriteRune('\n')
        sb.WriteString("    " + self.Get(i).String())
    }

    
    sb.WriteString("\n}")
    return sb.String()
}

func (self *StackMap) BitmapLen() int {
    return int(self.L)
}

func (self *StackMap) BitmapBytes() int {
    return int(self.L + 7) >> 3
}

func (self *StackMap) BitmapNums() int {
    return int(self.N)
}

func (self *StackMap) StackMapHeaderSize() int {
    return int(unsafe.Sizeof(self.L)) + int(unsafe.Sizeof(self.N)) 
}

func (self *StackMap) MarshalBinary() ([]byte, error) {
    size := self.BinaryLen()
    return BytesFrom(unsafe.Pointer(self), size, size), nil
}

func (self *StackMap) BinaryLen() int {
    return self.StackMapHeaderSize() + self.BitmapBytes() * self.BitmapNums()
}

var (
    byteType = UnpackEface(byte(0)).Type
)

const (
    _StackMapSize = unsafe.Sizeof(StackMap{})
)

//go:linkname mallocgc runtime.mallocgc

func mallocgc(nb uintptr, vt *GoType, zero bool) unsafe.Pointer

type StackMapBuilder struct {
    b Bitmap
}


//go:nocheckptr
func (self *StackMapBuilder) Build() (p *StackMap) {
    nb := len(self.b.B)
    allocatedSize := _StackMapSize + uintptr(nb) - 1
    bm := mallocgc(allocatedSize, byteType, false)

    
    p = (*StackMap)(bm)
    p.N, p.L = 1, int32(self.b.N)
    copy(BytesFrom(unsafe.Pointer(&p.B), nb, nb), self.b.B)

    
    if allocatedSize < uintptr(p.BinaryLen()) {
        panic(fmt.Sprintf("stackmap allocation too small: allocated %d, required %d", allocatedSize, p.BinaryLen()))
    }
    return
}

func (self *StackMapBuilder) AddField(ptr bool) {
    if ptr {
        self.b.Append(1)
    } else {
        self.b.Append(0)
    }
}

func (self *StackMapBuilder) AddFields(n int, ptr bool) {
    if ptr {
        self.b.AppendMany(n, 1)
    } else {
        self.b.AppendMany(n, 0)
    }
}