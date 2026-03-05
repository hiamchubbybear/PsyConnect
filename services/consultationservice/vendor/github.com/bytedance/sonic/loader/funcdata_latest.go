
// +build go1.18,!go1.25



package loader

import (
    `os`
    `sort`
    `unsafe`

    `github.com/bytedance/sonic/loader/internal/rt`
)

type funcTab struct {
    entry   uint32
    funcoff uint32
}

type pcHeader struct {
    magic          uint32  
    pad1, pad2     uint8   
    minLC          uint8   
    ptrSize        uint8   
    nfunc          int     
    nfiles         uint    
    textStart      uintptr 
    funcnameOffset uintptr 
    cuOffset       uintptr 
    filetabOffset  uintptr 
    pctabOffset    uintptr 
    pclnOffset     uintptr 
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

func makeFtab(funcs []_func, maxpc uint32) (ftab []funcTab, pclntabSize int64, startLocations []uint32) {
    
    
    
    pclntabSize = int64(len(funcs)*2*int(_PtrSize) + int(_PtrSize))
    startLocations = make([]uint32, len(funcs))
    for i, f := range funcs {
        pclntabSize = rnd(pclntabSize, int64(_PtrSize))
        
        startLocations[i] = uint32(pclntabSize)
        pclntabSize += int64(uint8(_FUNC_SIZE)+f.nfuncdata*4+uint8(f.npcdata)*4)
    }

    ftab = make([]funcTab, 0, len(funcs)+1)

    
    for i, f := range funcs {
        ftab = append(ftab, funcTab{uint32(f.entryOff), uint32(startLocations[i])})
    }

    
    ftab = append(ftab, funcTab{maxpc, 0})
    return
}


func makePclntable(size int64, startLocations []uint32, funcs []_func, maxpc uint32, pcdataOffs [][]uint32, funcdataOffs [][]uint32) (pclntab []byte) {
    
    
    
    pclntab = make([]byte, size, size)

    
    offs := 0
    for i, f := range funcs {
        byteOrder.PutUint32(pclntab[offs:offs+4], uint32(f.entryOff))
        byteOrder.PutUint32(pclntab[offs+4:offs+8], uint32(startLocations[i]))
        offs += 8
    }
    
    byteOrder.PutUint32(pclntab[offs:offs+4], maxpc)

    
    for i := range funcs {
        off := startLocations[i]

        
        fb := rt.BytesFrom(unsafe.Pointer(&funcs[i]), int(_FUNC_SIZE), int(_FUNC_SIZE))
        copy(pclntab[off:off+uint32(_FUNC_SIZE)], fb)
        off += uint32(_FUNC_SIZE)

        
        for j := 3; j < len(pcdataOffs[i]); j++ {
            byteOrder.PutUint32(pclntab[off:off+4], uint32(pcdataOffs[i][j]))
            off += 4
        }

        
        for _, funcdata := range funcdataOffs[i] {
            byteOrder.PutUint32(pclntab[off:off+4], uint32(funcdata))
            off += 4
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

        
        var pc = min + uint32((i+1)*_BUCKETSIZE)
        for ; e < len(ftab)-1 && ftab[e+1].entry <= pc; e++ {}

        for j := 0; j<_SUBBUCKETS && (i*_SUBBUCKETS+j)<int(n); j++ {
            
            fb._SUBBUCKETS[j] = byte(uint32(s) - fb.idx)

            
            pc = min + uint32(i*_BUCKETSIZE) + uint32((j+1)*_SUB_BUCKETSIZE)
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
    pctab, pcdataOffs, _funcs := makePctab(funcs, cuOff, nameOffs)
    mod.pctab = pctab

    
    
    
    cache := make([]byte, 0, len(funcs)*int(_PtrSize)) 
    fstart, funcdataOffs := writeFuncdata(&cache, funcs)

    
    ftab, pclntSize, startLocations := makeFtab(_funcs, uint32(len(text)))
    mod.ftab = ftab

    
    ffstart := writeFindfunctab(&cache, ftab)

    
    moduleCache.Lock()
    moduleCache.m[mod] = cache
    moduleCache.Unlock()
    mod.gofunc = uintptr(unsafe.Pointer(&cache[fstart]))
    mod.findfunctab = uintptr(unsafe.Pointer(&cache[ffstart]))

    
    pclntab := makePclntable(pclntSize, startLocations, _funcs, uint32(len(text)), pcdataOffs, funcdataOffs)
    mod.pclntable = pclntab

    
    mod.pcHeader = &pcHeader {
        magic   : _Magic,
        minLC   : _MinLC,
        ptrSize : _PtrSize,
        nfunc   : len(funcs),
        nfiles: uint(len(cu)),
        textStart: mod.text,
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



func makePctab(funcs []Func, cuOffset uint32, nameOffset []int32) (pctab []byte, pcdataOffs [][]uint32, _funcs []_func) {
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
        
        _f.entryOff = f.EntryOff
        _f.nameOff = nameOffset[i]
        _f.args = f.ArgsSize
        _f.deferreturn = f.DeferReturn
        
        _f.npcdata = uint32(_N_PCDATA)
        _f.cuOffset = cuOffset
        _f.funcID = f.ID
        _f.flag = f.Flag
        _f.nfuncdata = uint8(_N_FUNCDATA)
    }

    return
}

func registerFunction(name string, pc uintptr, textSize uintptr, fp int, args int, size uintptr, argptrs uintptr, localptrs uintptr) {} 
