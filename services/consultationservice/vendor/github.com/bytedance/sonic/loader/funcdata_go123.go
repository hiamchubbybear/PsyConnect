//go:build go1.23 && !go1.25
// +build go1.23,!go1.25



package loader

import (
    `unsafe`
    `github.com/bytedance/sonic/loader/internal/rt`
)

const (
    _Magic uint32 = 0xFFFFFFF1
)

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
    covctrs, ecovctrs     uintptr
    end, gcdata, gcbss    uintptr
    types, etypes         uintptr
    rodata                uintptr
    gofunc                uintptr 

    textsectmap []textSection 
    typelinks   []int32 
    itablinks   []*rt.GoItab

    ptab []ptabEntry

    pluginpath string
    pkghashes  []modulehash

    
	
	inittasks []unsafe.Pointer

    modulename   string
    modulehashes []modulehash

    hasmain uint8 
    bad bool 

    gcdatamask, gcbssmask bitVector

    typemap map[int32]*rt.GoType 

    next *moduledata
}

type _func struct {
    entryOff uint32 
    nameOff  int32  

    args        int32  
    deferreturn uint32 

    pcsp      uint32 
    pcfile    uint32
    pcln      uint32
    npcdata   uint32
    cuOffset  uint32 
    startLine int32  
    funcID    uint8 
    flag      uint8
    _         [1]byte 
    nfuncdata uint8   
    
    
    
    

    
    
    
    
    
    
    
    

    
    
    
    
    
    
    
    
}
