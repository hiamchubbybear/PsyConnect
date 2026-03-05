

package loader

import (
    `encoding`
    `encoding/binary`
    `fmt`
    `reflect`
    `strings`
    `sync`
    `unsafe`
)

const (
    _MinLC uint8 = 1
    _PtrSize uint8 = 8
)

const (
    _N_FUNCDATA = 8
    _INVALID_FUNCDATA_OFFSET = ^uint32(0)
    _FUNC_SIZE = unsafe.Sizeof(_func{})
    
    _MINFUNC = 16 
    _BUCKETSIZE    = 256 * _MINFUNC
    _SUBBUCKETS    = 16
    _SUB_BUCKETSIZE = _BUCKETSIZE / _SUBBUCKETS
)


const (
	FuncFlag_TOPFRAME = 1 << iota
	FuncFlag_SPWRITE
	FuncFlag_ASM
)




const (
    _FUNCDATA_ArgsPointerMaps    = 0
    _FUNCDATA_LocalsPointerMaps  = 1
    _FUNCDATA_StackObjects       = 2
    _FUNCDATA_InlTree            = 3
    _FUNCDATA_OpenCodedDeferInfo = 4
    _FUNCDATA_ArgInfo            = 5
    _FUNCDATA_ArgLiveInfo        = 6
    _FUNCDATA_WrapInfo           = 7

    
    
    
    
    ArgsSizeUnknown = -0x80000000
)


var moduleCache = struct {
    m map[*moduledata][]byte
    sync.Mutex
}{
    m: make(map[*moduledata][]byte),
}


type Func struct {
    ID          uint8  
    Flag        uint8  
    ArgsSize    int32  
    EntryOff    uint32 
    TextSize    uint32 
    DeferReturn uint32 
    FileIndex   uint32 
    Name        string 

    
    Pcsp            *Pcdata 
    Pcfile          *Pcdata 
    Pcline          *Pcdata 
    PcUnsafePoint   *Pcdata 
    PcStackMapIndex *Pcdata 
    PcInlTreeIndex  *Pcdata 
    PcArgLiveIndex  *Pcdata 
    
    
    ArgsPointerMaps    encoding.BinaryMarshaler 
    LocalsPointerMaps  encoding.BinaryMarshaler 
    StackObjects       encoding.BinaryMarshaler
    InlTree            encoding.BinaryMarshaler
    OpenCodedDeferInfo encoding.BinaryMarshaler
    ArgInfo            encoding.BinaryMarshaler
    ArgLiveInfo        encoding.BinaryMarshaler
    WrapInfo           encoding.BinaryMarshaler
}

func getOffsetOf(data interface{}, field string) uintptr {
    t := reflect.TypeOf(data)
    fv, ok := t.FieldByName(field)
    if !ok {
        panic(fmt.Sprintf("field %s not found in struct %s", field, t.Name()))
    }
    return fv.Offset
}

func rnd(v int64, r int64) int64 {
    if r <= 0 {
        return v
    }
    v += r - 1
    c := v % r
    if c < 0 {
        c += r
    }
    v -= c
    return v
}

var (
    byteOrder binary.ByteOrder = binary.LittleEndian
)

func funcNameParts(name string) (string, string, string) {
    i := strings.IndexByte(name, '[')
    if i < 0 {
        return name, "", ""
    }
    
    j := len(name) - 1
    for j > i && name[j] != ']' {
        j--
    }
    if j <= i {
        return name, "", ""
    }
    return name[:i], "[...]", name[j+1:]
}






func makeFuncnameTab(funcs []Func) (tab []byte, offs []int32) {
    offs = make([]int32, len(funcs))
    offset := 1
    tab = []byte{0}

    for i, f := range funcs {
        offs[i] = int32(offset)

        a, b, c := funcNameParts(f.Name)
        tab = append(tab, a...)
        tab = append(tab, b...)
        tab = append(tab, c...)
        tab = append(tab, 0)
        offset += len(a) + len(b) + len(c) + 1
    }

    return
}












func makeFilenametab(cus []compilationUnit) (cutab []uint32, filetab []byte, cuOffsets []uint32) {
    cuOffsets = make([]uint32, len(cus))
    cuOffset := 0
    fileOffset := 0

    for i, cu := range cus {
        cuOffsets[i] = uint32(cuOffset)

        for _, name := range cu.fileNames {
            cutab = append(cutab, uint32(fileOffset))

            fileOffset += len(name) + 1
            filetab = append(filetab, name...)
            filetab = append(filetab, 0)
        }

        cuOffset += len(cu.fileNames)
    }

    return
}

func writeFuncdata(out *[]byte, funcs []Func) (fstart int, funcdataOffs [][]uint32) {
    fstart = len(*out)
    *out = append(*out, byte(0))
    offs := uint32(1)

    funcdataOffs = make([][]uint32, len(funcs))
    for i, f := range funcs {

        var writer = func(fd encoding.BinaryMarshaler) {
            var ab []byte
            var err error
            if fd != nil {
                ab, err = fd.MarshalBinary()
                if err != nil {
                    panic(err)
                }
                funcdataOffs[i] = append(funcdataOffs[i], offs)
            } else {
                ab = []byte{0}
                funcdataOffs[i] = append(funcdataOffs[i], _INVALID_FUNCDATA_OFFSET)
            }
            *out = append(*out, ab...)
            offs += uint32(len(ab))
        }

        writer(f.ArgsPointerMaps)
        writer(f.LocalsPointerMaps)
        writer(f.StackObjects)
        writer(f.InlTree)
        writer(f.OpenCodedDeferInfo)
        writer(f.ArgInfo)
        writer(f.ArgLiveInfo)
        writer(f.WrapInfo)
    }
    return 
}
