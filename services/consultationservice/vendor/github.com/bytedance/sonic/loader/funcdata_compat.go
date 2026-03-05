//go:build !go1.17 || go1.25
// +build !go1.17 go1.25



package loader

import (
   `os`
   `unsafe`
   `sort`

   `github.com/bytedance/sonic/loader/internal/rt`
)

const (
    _Magic uint32 = 0xfffffffa
)

type pcHeader struct {
    magic          uint32  
    pad1, pad2     uint8   
    minLC          uint8   
    ptrSize        uint8   
    nfunc          int     
    nfiles         uint    
    funcnameOffset uintptr 
    cuOffset       uintptr 
    filetabOffset  uintptr 
    pctabOffset    uintptr 
    pclnOffset     uintptr 
}

type moduledata struct {
    pcHeader     *pcHeader
    funcnametab  []byte
    cutab        []uint32
    filetab      []byte
    pctab        []byte
    pclntable    []byte
    ftab         []funcTab
    findfunctab  uintptr
    minpc, maxpc uintptr 

    text, etext           uintptr 
    noptrdata, enoptrdata uintptr
    data, edata           uintptr
    bss, ebss             uintptr
    noptrbss, enoptrbss   uintptr
    end, gcdata, gcbss    uintptr
    types, etypes         uintptr
    
    textsectmap []textSection 
    typelinks   []int32 
    itablinks   []*rt.GoItab

    ptab []ptabEntry

    pluginpath string
    pkghashes  []modulehash

    modulename   string
    modulehashes []modulehash

    hasmain uint8 

    gcdatamask, gcbssmask bitVector

    typemap map[int32]*rt.GoType 

    bad bool 

    next *moduledata
}

type _func struct {
    entry    uintptr 
    nameOff  int32  

    args        int32  
    deferreturn uint32 

    pcsp      uint32 
    pcfile    uint32
    pcln      uint32
    npcdata   uint32
    cuOffset  uint32 
    funcID    uint8  
    _         [2]byte 
    nfuncdata uint8   
    
    
    
    

    
    
    
    
    
    
    
    

    
    
    
    
    
    
    
    
}

type funcTab struct {
    entry   uintptr
    funcoff uintptr
}

type bitVector struct {
    n        int32 
    bytedata *uint8
}

type ptabEntry struct {
    name int32
    typ  int32
}

type textSection struct {
    vaddr    uintptr 
    end      uintptr 
    baseaddr uintptr 
}

type modulehash struct {
    modulename   string
    linktimehash string
    runtimehash  *string
}









type findfuncbucket struct {
    idx        uint32
    _SUBBUCKETS [16]byte
}


type compilationUnit struct {
    fileNames []string
}

func makeFtab(funcs []_func, maxpc uintptr) (ftab []funcTab, pclntabSize int64, startLocations []uint32) {
    
    
    
    pclntabSize = int64(len(funcs)*2*int(_PtrSize) + int(_PtrSize))
    startLocations = make([]uint32, len(funcs))
    for i, f := range funcs {
        pclntabSize = rnd(pclntabSize, int64(_PtrSize))
        
        startLocations[i] = uint32(pclntabSize)
        pclntabSize += int64(uint8(_FUNC_SIZE) + f.nfuncdata*_PtrSize + uint8(f.npcdata)*4)
    }
    ftab = make([]funcTab, 0, len(funcs)+1)

    
    for i, f := range funcs {
        ftab = append(ftab, funcTab{uintptr(f.entry), uintptr(startLocations[i])})
    }

    
    ftab = append(ftab, funcTab{maxpc, 0})
    return
}


func makePclntable(size int64, startLocations []uint32, funcs []_func, maxpc uintptr, pcdataOffs [][]uint32, funcdataAddr uintptr, funcdataOffs [][]uint32) (pclntab []byte) {
    pclntab = make([]byte, size, size)

    
    offs := 0
    for i, f := range funcs {
        byteOrder.PutUint64(pclntab[offs:offs+8], uint64(f.entry))
        byteOrder.PutUint64(pclntab[offs+8:offs+16], uint64(startLocations[i]))
        offs += 16
    }
    
    byteOrder.PutUint64(pclntab[offs:offs+8], uint64(maxpc))
    offs += 8

    
    for i, f := range funcs {
        off := startLocations[i]

        
        byteOrder.PutUint64(pclntab[off:off+8], uint64(f.entry))
        off += 8
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.nameOff))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.args))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.deferreturn))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.pcsp))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.pcfile))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.pcln))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.npcdata))
        off += 4
        byteOrder.PutUint32(pclntab[off:off+4], uint32(f.cuOffset))
        off += 4
        pclntab[off] = f.funcID
        
        off += 3
        pclntab[off] = f.nfuncdata
        off += 1

        
        for j := 3; j < len(pcdataOffs[i]); j++ {
            byteOrder.PutUint32(pclntab[off:off+4], uint32(pcdataOffs[i][j]))
            off += 4
        }

        off = uint32(rnd(int64(off), int64(_PtrSize)))

        
        for _, funcdata := range funcdataOffs[i] {
            if funcdata == _INVALID_FUNCDATA_OFFSET {
                byteOrder.PutUint64(pclntab[off:off+8], 0)
            } else {
                byteOrder.PutUint64(pclntab[off:off+8], uint64(funcdataAddr)+uint64(funcdata))
            }
            off += 8
        }
    }

    return
}








func writeFindfunctab(out *[]byte, ftab []funcTab) (start int) {
    start = len(*out)

    max := ftab[len(ftab)-1].entry
    min := ftab[0].entry
    nbuckets := (max - min + _BUCKETSIZE - 1) / _BUCKETSIZE
    n := (max - min + _SUB_BUCKETSIZE - 1) / _SUB_BUCKETSIZE

    tab := make([]findfuncbucket, 0, nbuckets)
    var s, e = 0, 0
    for i := 0; i<int(nbuckets); i++ {
        
        var fb = findfuncbucket{idx: uint32(s)}

        
        var pc = min + uintptr((i+1)*_BUCKETSIZE)
        for ; e < len(ftab)-1 && ftab[e+1].entry <= pc; e++ {}
        
        for j := 0; j<_SUBBUCKETS && (i*_SUBBUCKETS+j)<int(n); j++ {
            
            fb._SUBBUCKETS[j] = byte(uint32(s) - fb.idx)
            
            
            pc = min + uintptr(i*_BUCKETSIZE) + uintptr((j+1)*_SUB_BUCKETSIZE)
            for ; s < len(ftab)-1 && ftab[s+1].entry <= pc; s++ {}            
        }

        s = e
        tab = append(tab, fb)
    }

    
    if len(tab) > 0 {
        size := int(unsafe.Sizeof(findfuncbucket{}))*len(tab)
        *out = append(*out, rt.BytesFrom(unsafe.Pointer(&tab[0]), size, size)...)
    }
    return 
}

func makeModuledata(name string, filenames []string, funcsp *[]Func, text []byte) (mod *moduledata) {
    mod = new(moduledata)
    mod.modulename = name

    
    funcs := *funcsp
    sort.Slice(funcs, func(i, j int) bool {
        return funcs[i].EntryOff < funcs[j].EntryOff
    })
    *funcsp = funcs

    
    cu := make([]string, 0, len(filenames))
    cu = append(cu, filenames...)
    cutab, filetab, cuOffs := makeFilenametab([]compilationUnit{{cu}})
    mod.cutab = cutab
    mod.filetab = filetab

    
    funcnametab, nameOffs := makeFuncnameTab(funcs)
    mod.funcnametab = funcnametab

    
    p := os.Getpagesize()
    size := int(rnd(int64(len(text)), int64(p)))
    addr := mmap(size)
    
    s := rt.BytesFrom(unsafe.Pointer(addr), len(text), size)
    copy(s, text)
    
    mprotect(addr, size)

    
    mod.text = addr
    mod.etext = addr + uintptr(size)
    mod.minpc = addr
    mod.maxpc = addr + uintptr(len(text))

    
    
    cuOff := cuOffs[0]
    pctab, pcdataOffs, _funcs := makePctab(funcs, addr, cuOff, nameOffs)
    mod.pctab = pctab

    
    
    
    cache := make([]byte, 0, len(funcs)*int(_PtrSize)) 
    fstart, funcdataOffs := writeFuncdata(&cache, funcs)

    
    ftab, pclntSize, startLocations := makeFtab(_funcs, mod.maxpc)
    mod.ftab = ftab

    
    ffstart := writeFindfunctab(&cache, ftab)

    
    moduleCache.Lock()
    moduleCache.m[mod] = cache
    moduleCache.Unlock()
    mod.findfunctab = uintptr(rt.IndexByte(cache, ffstart))
    funcdataAddr := uintptr(rt.IndexByte(cache, fstart))

    
    pclntab := makePclntable(pclntSize, startLocations, _funcs, mod.maxpc, pcdataOffs, funcdataAddr, funcdataOffs)
    mod.pclntable = pclntab

    
    mod.pcHeader = &pcHeader {
        magic   : _Magic,
        minLC   : _MinLC,
        ptrSize : _PtrSize,
        nfunc   : len(funcs),
        nfiles: uint(len(cu)),
        funcnameOffset: getOffsetOf(moduledata{}, "funcnametab"),
        cuOffset: getOffsetOf(moduledata{}, "cutab"),
        filetabOffset: getOffsetOf(moduledata{}, "filetab"),
        pctabOffset: getOffsetOf(moduledata{}, "pctab"),
        pclnOffset: getOffsetOf(moduledata{}, "pclntable"),
    }

    
    mod.gcdata = uintptr(unsafe.Pointer(&emptyByte))
    mod.gcbss = uintptr(unsafe.Pointer(&emptyByte))

    return
}



func makePctab(funcs []Func, addr uintptr, cuOffset uint32, nameOffset []int32) (pctab []byte, pcdataOffs [][]uint32, _funcs []_func) {
    _funcs = make([]_func, len(funcs))

    
    
    
    pctab = make([]byte, 1, 12*len(funcs)+1)
    pcdataOffs = make([][]uint32, len(funcs))

    for i, f := range funcs {
        _f := &_funcs[i]

        var writer = func(pc *Pcdata) {
            var ab []byte
            var err error
            if pc != nil {
                ab, err = pc.MarshalBinary()
                if err != nil {
                    panic(err)
                }
                pcdataOffs[i] = append(pcdataOffs[i], uint32(len(pctab)))
            } else {
                ab = []byte{0}
                pcdataOffs[i] = append(pcdataOffs[i], _PCDATA_INVALID_OFFSET)
            }
            pctab = append(pctab, ab...)
        }

        if f.Pcsp != nil {
            _f.pcsp = uint32(len(pctab))
        }
        writer(f.Pcsp)
        if f.Pcfile != nil {
            _f.pcfile = uint32(len(pctab))
        }
        writer(f.Pcfile)
        if f.Pcline != nil {
            _f.pcln = uint32(len(pctab))
        }
        writer(f.Pcline)
        writer(f.PcUnsafePoint)
        writer(f.PcStackMapIndex)
        writer(f.PcInlTreeIndex)
        writer(f.PcArgLiveIndex)
        
        _f.entry = addr + uintptr(f.EntryOff)
        _f.nameOff = nameOffset[i]
        _f.args = f.ArgsSize
        _f.deferreturn = f.DeferReturn
        
        _f.npcdata = uint32(_N_PCDATA)
        _f.cuOffset = cuOffset
        _f.funcID = f.ID
        _f.nfuncdata = uint8(_N_FUNCDATA)
    }

    return
}

func registerFunction(name string, pc uintptr, textSize uintptr, fp int, args int, size uintptr, argptrs uintptr, localptrs uintptr) {} 
