

package loader

import (
    `encoding/binary`
)

const (
    _N_PCDATA   = 4

    _PCDATA_UnsafePoint   = 0
    _PCDATA_StackMapIndex = 1
    _PCDATA_InlTreeIndex  = 2
    _PCDATA_ArgLiveIndex  = 3

    _PCDATA_INVALID_OFFSET = 0
)

const (
    
    PCDATA_UnsafePointSafe   = -1 
    PCDATA_UnsafePointUnsafe = -2 

    
    
    
    
    
    PCDATA_Restart1 = -3
    PCDATA_Restart2 = -4

    
    
    PCDATA_RestartAtEntry = -5

    _PCDATA_START_VAL = -1
)

var emptyByte byte



type Pcvalue struct {
    PC  uint32 
    Val int32  
}




type Pcdata []Pcvalue


func (self Pcdata) MarshalBinary() (data []byte, err error) {
    
    sv := int32(_PCDATA_START_VAL)
    sp := uint32(0)
    buf := make([]byte, binary.MaxVarintLen32)
    for _, v := range self {
        if v.PC < sp {
            panic("PC must be in ascending order!")
        }
        dp := uint64(v.PC - sp)
        dv := int64(v.Val - sv)
        if dv == 0 || dp == 0 {
            continue
        }
        n := binary.PutVarint(buf, dv)
        data = append(data, buf[:n]...)
        n2 := binary.PutUvarint(buf, dp)
        data = append(data, buf[:n2]...)
        sp = v.PC
        sv = v.Val
    }
    
    data = append(data, 0)
    return
}
