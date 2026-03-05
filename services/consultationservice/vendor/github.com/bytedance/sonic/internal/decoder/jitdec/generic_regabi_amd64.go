// +build go1.17,!go1.25



package jitdec

import (
    `encoding/json`
    `fmt`
    `reflect`

    `github.com/bytedance/sonic/internal/jit`
    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/twitchyliquid64/golang-asm/obj`
    `github.com/bytedance/sonic/internal/rt`
)



const (
    _VD_args   = 8      
    _VD_fargs  = 64     
    _VD_saves  = 48     
    _VD_locals = 96     
)

const (
    _VD_offs = _VD_fargs + _VD_saves + _VD_locals
    _VD_size = _VD_offs + 8     
)

var (
    _VAR_ss = _VAR_ss_Vt
    _VAR_df = jit.Ptr(_SP, _VD_fargs + _VD_saves)
)

var (
    _VAR_ss_Vt = jit.Ptr(_SP, _VD_fargs + _VD_saves + 8)
    _VAR_ss_Dv = jit.Ptr(_SP, _VD_fargs + _VD_saves + 16)
    _VAR_ss_Iv = jit.Ptr(_SP, _VD_fargs + _VD_saves + 24)
    _VAR_ss_Ep = jit.Ptr(_SP, _VD_fargs + _VD_saves + 32)
    _VAR_ss_Db = jit.Ptr(_SP, _VD_fargs + _VD_saves + 40)
    _VAR_ss_Dc = jit.Ptr(_SP, _VD_fargs + _VD_saves + 48)
)

var (
    _VAR_R9 = jit.Ptr(_SP, _VD_fargs + _VD_saves + 56)
)
type _ValueDecoder struct {
    jit.BaseAssembler
}

var (
    _VAR_cs_LR = jit.Ptr(_SP, _VD_fargs + _VD_saves + 64)
    _VAR_cs_p = jit.Ptr(_SP, _VD_fargs + _VD_saves + 72)
    _VAR_cs_n = jit.Ptr(_SP, _VD_fargs + _VD_saves + 80)
    _VAR_cs_d = jit.Ptr(_SP, _VD_fargs + _VD_saves + 88)
)

func (self *_ValueDecoder) build() uintptr {
    self.Init(self.compile)
    return *(*uintptr)(self.Load("decode_value", _VD_size, _VD_args, argPtrs_generic, localPtrs_generic))
}



func (self *_ValueDecoder) save(r ...obj.Addr) {
    for i, v := range r {
        if i > _VD_saves / 8 - 1 {
            panic("too many registers to save")
        } else {
            self.Emit("MOVQ", v, jit.Ptr(_SP, _VD_fargs + int64(i) * 8))
        }
    }
}

func (self *_ValueDecoder) load(r ...obj.Addr) {
    for i, v := range r {
        if i > _VD_saves / 8 - 1 {
            panic("too many registers to load")
        } else {
            self.Emit("MOVQ", jit.Ptr(_SP, _VD_fargs + int64(i) * 8), v)
        }
    }
}

func (self *_ValueDecoder) call(fn obj.Addr) {
    self.Emit("MOVQ", fn, _R9)  
    self.Rjmp("CALL", _R9)      
}

func (self *_ValueDecoder) call_go(fn obj.Addr) {
    self.save(_REG_go...)   
    self.call(fn)           
    self.load(_REG_go...)   
}

func (self *_ValueDecoder) callc(fn obj.Addr) {
    self.save(_IP)  
    self.call(fn)
    self.load(_IP)  
}

func (self *_ValueDecoder) call_c(fn obj.Addr) {
    self.Emit("XCHGQ", _IC, _BX)
    self.callc(fn)
    self.Emit("XCHGQ", _IC, _BX)
}



const (
    _S_val = iota + 1
    _S_arr
    _S_arr_0
    _S_obj
    _S_obj_0
    _S_obj_delim
    _S_obj_sep
)

const (
    _S_omask_key = (1 << _S_obj_0) | (1 << _S_obj_sep)
    _S_omask_end = (1 << _S_obj_0) | (1 << _S_obj)
    _S_vmask = (1 << _S_val) | (1 << _S_arr_0)
)

const (
    _A_init_len = 1
    _A_init_cap = 16
)

const (
    _ST_Sp = 0
    _ST_Vt = _PtrBytes
    _ST_Vp = _PtrBytes * (types.MAX_RECURSE + 1)
)

var (
    _V_true  = jit.Imm(int64(pbool(true)))
    _V_false = jit.Imm(int64(pbool(false)))
    _F_value = jit.Imm(int64(native.S_value))
)

var (
    _V_max     = jit.Imm(int64(types.V_MAX))
    _E_eof     = jit.Imm(int64(types.ERR_EOF))
    _E_invalid = jit.Imm(int64(types.ERR_INVALID_CHAR))
    _E_recurse = jit.Imm(int64(types.ERR_RECURSE_EXCEED_MAX))
)

var (
    _F_convTslice    = jit.Func(rt.ConvTslice)
    _F_convTstring   = jit.Func(rt.ConvTstring)
    _F_invalid_vtype = jit.Func(invalid_vtype)
)

var (
    _T_map     = jit.Type(reflect.TypeOf((map[string]interface{})(nil)))
    _T_bool    = jit.Type(reflect.TypeOf(false))
    _T_int64   = jit.Type(reflect.TypeOf(int64(0)))
    _T_eface   = jit.Type(reflect.TypeOf((*interface{})(nil)).Elem())
    _T_slice   = jit.Type(reflect.TypeOf(([]interface{})(nil)))
    _T_string  = jit.Type(reflect.TypeOf(""))
    _T_number  = jit.Type(reflect.TypeOf(json.Number("")))
    _T_float64 = jit.Type(reflect.TypeOf(float64(0)))
)

var _R_tab = map[int]string {
    '[': "_decode_V_ARRAY",
    '{': "_decode_V_OBJECT",
    ':': "_decode_V_KEY_SEP",
    ',': "_decode_V_ELEM_SEP",
    ']': "_decode_V_ARRAY_END",
    '}': "_decode_V_OBJECT_END",
}

func (self *_ValueDecoder) compile() {
    self.Emit("SUBQ", jit.Imm(_VD_size), _SP)       
    self.Emit("MOVQ", _BP, jit.Ptr(_SP, _VD_offs))  
    self.Emit("LEAQ", jit.Ptr(_SP, _VD_offs), _BP)  

    
    self.Emit("XORL", _CX, _CX)                                 
    self.Emit("MOVQ", _DF, _VAR_df)                             
    
    self.Emit("MOVQ", jit.Imm(_MaxDigitNums), _VAR_ss_Dc)       
    self.Emit("LEAQ", jit.Ptr(_ST, _DbufOffset), _AX)           
    self.Emit("MOVQ", _AX, _VAR_ss_Db)                          
    
    self.Emit("ADDQ", jit.Imm(_FsmOffset), _ST)                 
    self.Emit("MOVQ", _CX, jit.Ptr(_ST, _ST_Sp))                
    self.WriteRecNotAX(0, _VP, jit.Ptr(_ST, _ST_Vp), false)                
    self.Emit("MOVQ", jit.Imm(_S_val), jit.Ptr(_ST, _ST_Vt))    
    self.Sjmp("JMP" , "_next")                                  

    
    self.Link("_set_value")                                 
    self.Emit("MOVL" , jit.Imm(_S_vmask), _DX)              
    self.Emit("MOVQ" , jit.Ptr(_ST, _ST_Sp), _CX)           
    self.Emit("MOVQ" , jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)   
    self.Emit("BTQ"  , _AX, _DX)                            
    self.Sjmp("JNC"  , "_vtype_error")                      
    self.Emit("XORL" , _SI, _SI)                            
    self.Emit("SUBQ" , jit.Imm(1), jit.Ptr(_ST, _ST_Sp))    
    self.Emit("XCHGQ", jit.Sib(_ST, _CX, 8, _ST_Vp), _SI)   
    self.Emit("MOVQ" , _R8, jit.Ptr(_SI, 0))                
    self.WriteRecNotAX(1, _R9, jit.Ptr(_SI, 8), false)           

    
    self.Link("_next")                              
    self.Emit("MOVQ" , jit.Ptr(_ST, _ST_Sp), _AX)   
    self.Emit("TESTQ", _AX, _AX)                    
    self.Sjmp("JS"   , "_return")                   

    
    self.Emit("CMPQ"   , _IC, _IL)                      
    self.Sjmp("JAE"    , "_decode_V_EOF")               
    self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  
    self.Emit("MOVQ"   , jit.Imm(_BM_space), _DX)       
    self.Emit("CMPQ"   , _AX, jit.Imm(' '))             
    self.Sjmp("JA"     , "_decode_fast")                
    self.Emit("BTQ"    , _AX, _DX)                      
    self.Sjmp("JNC"    , "_decode_fast")                
    self.Emit("ADDQ"   , jit.Imm(1), _IC)               

    
    for i := 0; i < 3; i++ {
        self.Emit("CMPQ"   , _IC, _IL)                      
        self.Sjmp("JAE"    , "_decode_V_EOF")               
        self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  
        self.Emit("CMPQ"   , _AX, jit.Imm(' '))             
        self.Sjmp("JA"     , "_decode_fast")                
        self.Emit("BTQ"    , _AX, _DX)                      
        self.Sjmp("JNC"    , "_decode_fast")                
        self.Emit("ADDQ"   , jit.Imm(1), _IC)               
    }

    
    self.Emit("CMPQ"   , _IC, _IL)                      
    self.Sjmp("JAE"    , "_decode_V_EOF")               
    self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  

    
    self.Link("_decode_fast")                           
    self.Byte(0x48, 0x8d, 0x3d)                         
    self.Sref("_decode_tab", 4)                         
    self.Emit("MOVLQSX", jit.Sib(_DI, _AX, 4, 0), _AX)  
    self.Emit("TESTQ"  , _AX, _AX)                      
    self.Sjmp("JZ"     , "_decode_native")              
    self.Emit("ADDQ"   , jit.Imm(1), _IC)               
    self.Emit("ADDQ"   , _DI, _AX)                      
    self.Rjmp("JMP"    , _AX)                           

    
    self.Link("_decode_native")         
    self.Emit("MOVQ", _IP, _DI)         
    self.Emit("MOVQ", _IL, _SI)         
    self.Emit("MOVQ", _IC, _DX)         
    self.Emit("LEAQ", _VAR_ss, _CX)     
    self.Emit("MOVQ", _VAR_df, _R8)     
    self.Emit("BTSQ", jit.Imm(_F_allow_control), _R8)  
    self.callc(_F_value)                
    self.Emit("MOVQ", _AX, _IC)         

    
    self.Emit("MOVQ" , _VAR_ss_Vt, _AX)     
    self.Emit("TESTQ", _AX, _AX)            
    self.Sjmp("JS"   , "_parsing_error")       
    self.Sjmp("JZ"   , "_invalid_vtype")    
    self.Emit("CMPQ" , _AX, _V_max)         
    self.Sjmp("JA"   , "_invalid_vtype")    

    
    self.Byte(0x48, 0x8d, 0x3d)                             
    self.Sref("_switch_table", 4)                           
    self.Emit("MOVLQSX", jit.Sib(_DI, _AX, 4, -4), _AX)     
    self.Emit("ADDQ"   , _DI, _AX)                          
    self.Rjmp("JMP"    , _AX)                               

    
    self.Link("_decode_V_EOF")          
    self.Emit("MOVL", _E_eof, _EP)      
    self.Sjmp("JMP" , "_error")         

    
    self.Link("_decode_V_NULL")                 
    self.Emit("XORL", _R8, _R8)                 
    self.Emit("XORL", _R9, _R9)                 
    self.Emit("LEAQ", jit.Ptr(_IC, -4), _DI)    
    self.Sjmp("JMP" , "_set_value")             

    
    self.Link("_decode_V_TRUE")                 
    self.Emit("MOVQ", _T_bool, _R8)             
    
    self.Emit("MOVQ", _V_true, _R9)             
    self.Emit("LEAQ", jit.Ptr(_IC, -4), _DI)    
    self.Sjmp("JMP" , "_set_value")             

    
    self.Link("_decode_V_FALSE")                
    self.Emit("MOVQ", _T_bool, _R8)             
    self.Emit("MOVQ", _V_false, _R9)            
    self.Emit("LEAQ", jit.Ptr(_IC, -5), _DI)    
    self.Sjmp("JMP" , "_set_value")             

    
    self.Link("_decode_V_ARRAY")                            
    self.Emit("MOVL", jit.Imm(_S_vmask), _DX)               
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)    
    self.Emit("BTQ" , _AX, _DX)                             
    self.Sjmp("JNC" , "_invalid_char")                      

    
    self.Emit("MOVQ", _T_eface, _AX)                            
    self.Emit("MOVQ", jit.Imm(_A_init_len), _BX)                
    self.Emit("MOVQ", jit.Imm(_A_init_cap), _CX)                
    self.call_go(_F_makeslice)                                  

    
    self.Emit("MOVQ", jit.Imm(_A_init_len), _BX)                
    self.Emit("MOVQ", jit.Imm(_A_init_cap), _CX)                
    self.call_go(_F_convTslice)                                 
    self.Emit("MOVQ", _AX, _R8)                                 

    
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)                        
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vp), _SI)                
    self.Emit("MOVQ", jit.Imm(_S_arr), jit.Sib(_ST, _CX, 8, _ST_Vt))    
    self.Emit("MOVQ", _T_slice, _AX)                                    
    self.Emit("MOVQ", _AX, jit.Ptr(_SI, 0))                             
    self.WriteRecNotAX(2, _R8, jit.Ptr(_SI, 8), false)                  

    
    self.Emit("ADDQ", jit.Imm(1), _CX)                                  
    self.Emit("CMPQ", _CX, jit.Imm(types.MAX_RECURSE))                  
    self.Sjmp("JAE"  , "_stack_overflow")                                
    self.Emit("MOVQ", jit.Ptr(_R8, 0), _AX)                             
    self.Emit("MOVQ", _CX, jit.Ptr(_ST, _ST_Sp))                        
    self.WritePtrAX(3, jit.Sib(_ST, _CX, 8, _ST_Vp), false)             
    self.Emit("MOVQ", jit.Imm(_S_arr_0), jit.Sib(_ST, _CX, 8, _ST_Vt))  
    self.Sjmp("JMP" , "_next")                                          

    
    self.Link("_decode_V_OBJECT")                                       
    self.Emit("MOVL", jit.Imm(_S_vmask), _DX)                           
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)                        
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)                
    self.Emit("BTQ" , _AX, _DX)                                         
    self.Sjmp("JNC" , "_invalid_char")                                  
    self.call_go(_F_makemap_small)                                      
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)                        
    self.Emit("MOVQ", jit.Imm(_S_obj_0), jit.Sib(_ST, _CX, 8, _ST_Vt))    
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vp), _SI)                
    self.Emit("MOVQ", _T_map, _DX)                                      
    self.Emit("MOVQ", _DX, jit.Ptr(_SI, 0))                             
    self.WritePtrAX(4, jit.Ptr(_SI, 8), false)                          
    self.Sjmp("JMP" , "_next")                                          

    
    self.Link("_decode_V_STRING")       
    self.Emit("MOVQ", _VAR_ss_Iv, _CX)  
    self.Emit("MOVQ", _IC, _AX)         
    self.Emit("SUBQ", _CX, _AX)         

    
    self.Emit("CMPQ", _VAR_ss_Ep, jit.Imm(-1))          
    self.Sjmp("JNE" , "_unquote")                       
    self.Emit("SUBQ", jit.Imm(1), _AX)                  
    self.Emit("LEAQ", jit.Sib(_IP, _CX, 1, 0), _R8)     
    self.Byte(0x48, 0x8d, 0x3d)                         
    self.Sref("_copy_string_end", 4)
    self.Emit("BTQ", jit.Imm(_F_copy_string), _VAR_df)
    self.Sjmp("JC", "copy_string")
    self.Link("_copy_string_end")                                 
    self.Emit("XORL", _DX, _DX)   

    
    self.Link("_noescape")                                  
    self.Emit("MOVL", jit.Imm(_S_omask_key), _DI)               
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vt), _SI)    
    self.Emit("BTQ" , _SI, _DI)                             
    self.Sjmp("JC"  , "_object_key")                        

    
    self.Emit("TESTQ", _DX, _DX)                
    self.Sjmp("JNZ"  , "_packed_str")           
    self.Emit("MOVQ" , _AX, _BX)                
    self.Emit("MOVQ" , _R8, _AX)                
    self.call_go(_F_convTstring)                
    self.Emit("MOVQ" , _AX, _R9)                

    
    self.Link("_packed_str")            
    self.Emit("MOVQ", _T_string, _R8)   
    self.Emit("MOVQ", _VAR_ss_Iv, _DI)  
    self.Emit("SUBQ", jit.Imm(1), _DI)  
    self.Sjmp("JMP" , "_set_value")     

    
    self.Link("_object_key")
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vp), _SI)    
    self.Emit("MOVQ", jit.Ptr(_SI, 8), _SI)                 

    
    self.Emit("ADDQ", jit.Imm(1), _CX)                                      
    self.Emit("CMPQ", _CX, jit.Imm(types.MAX_RECURSE))                      
    self.Sjmp("JAE"  , "_stack_overflow")                                    
    self.Emit("MOVQ", _CX, jit.Ptr(_ST, _ST_Sp))                            
    self.Emit("MOVQ", jit.Imm(_S_obj_delim), jit.Sib(_ST, _CX, 8, _ST_Vt))  

    
    self.Emit("MOVQ", _AX, _DI)                         
    self.Emit("MOVQ", _T_map, _AX)                      
    self.Emit("MOVQ", _SI, _BX)                         
    self.Emit("MOVQ", _R8, _CX)                         
    self.call_go(_F_mapassign_faststr)                  

    
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)                 
    self.WritePtrAX(6, jit.Sib(_ST, _CX, 8, _ST_Vp), false)    
    self.Sjmp("JMP" , "_next")                                   

    
    self.Link("_unquote")                               
    self.Emit("ADDQ", jit.Imm(15), _AX)                 
    self.Emit("MOVQ", _T_byte, _BX)                     
    self.Emit("MOVB", jit.Imm(0), _CX)                  
    self.call_go(_F_mallocgc)                           
    self.Emit("MOVQ", _AX, _R9)                         

    
    self.Emit("MOVQ" , _VAR_ss_Iv, _CX)                         
    self.Emit("LEAQ" , jit.Sib(_IP, _CX, 1, 0), _DI)            
    self.Emit("NEGQ" , _CX)                                     
    self.Emit("LEAQ" , jit.Sib(_IC, _CX, 1, -1), _SI)           
    self.Emit("LEAQ" , jit.Ptr(_R9, 16), _DX)                   
    self.Emit("LEAQ" , _VAR_ss_Ep, _CX)                         
    self.Emit("XORL" , _R8, _R8)                                
    self.Emit("BTQ"  , jit.Imm(_F_disable_urc), _VAR_df)        
    self.Emit("SETCC", _R8)                                     
    self.Emit("SHLQ" , jit.Imm(types.B_UNICODE_REPLACE), _R8)   

    
    self.Emit("MOVQ", _R9, _VAR_R9)             
    self.call_c(_F_unquote)                     
    self.Emit("MOVQ", _VAR_R9, _R9)             

    
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , "_unquote_error")        
    self.Emit("MOVL" , jit.Imm(1), _DX)         
    self.Emit("LEAQ" , jit.Ptr(_R9, 16), _R8)   
    self.Emit("MOVQ" , _R8, jit.Ptr(_R9, 0))    
    self.Emit("MOVQ" , _AX, jit.Ptr(_R9, 8))    
    self.Sjmp("JMP"  , "_noescape")             

    
    self.Link("_decode_V_DOUBLE")                           
    self.Emit("BTQ"  , jit.Imm(_F_use_number), _VAR_df)     
    self.Sjmp("JC"   , "_use_number")                       
    self.Emit("MOVSD", _VAR_ss_Dv, _X0)                     
    self.Sjmp("JMP"  , "_use_float64")                      

    
    self.Link("_decode_V_INTEGER")                          
    self.Emit("BTQ"     , jit.Imm(_F_use_number), _VAR_df)  
    self.Sjmp("JC"      , "_use_number")                    
    self.Emit("BTQ"     , jit.Imm(_F_use_int64), _VAR_df)   
    self.Sjmp("JC"      , "_use_int64")                     
    
    self.Emit("MOVSD", _VAR_ss_Dv, _X0)                  

    
    self.Link("_use_float64")                   
    self.Emit("MOVQ" , _X0, _AX)                
    self.call_go(_F_convT64)                    
    self.Emit("MOVQ" , _T_float64, _R8)         
    self.Emit("MOVQ" , _AX, _R9)                
    self.Emit("MOVQ" , _VAR_ss_Ep, _DI)         
    self.Sjmp("JMP"  , "_set_value")            

    
    self.Link("_use_number")                            
    self.Emit("MOVQ", _VAR_ss_Ep, _AX)                  
    self.Emit("LEAQ", jit.Sib(_IP, _AX, 1, 0), _SI)     
    self.Emit("MOVQ", _IC, _CX)                         
    self.Emit("SUBQ", _AX, _CX)                         
    self.Emit("MOVQ", _SI, _AX)                         
    self.Emit("MOVQ", _CX, _BX)                         
    self.call_go(_F_convTstring)                        
    self.Emit("MOVQ", _T_number, _R8)                   
    self.Emit("MOVQ", _AX, _R9)                         
    self.Emit("MOVQ", _VAR_ss_Ep, _DI)                  
    self.Sjmp("JMP" , "_set_value")                     

    
    self.Link("_use_int64")                     
    self.Emit("MOVQ", _VAR_ss_Iv, _AX)          
    self.call_go(_F_convT64)                    
    self.Emit("MOVQ", _T_int64, _R8)            
    self.Emit("MOVQ", _AX, _R9)                 
    self.Emit("MOVQ", _VAR_ss_Ep, _DI)          
    self.Sjmp("JMP" , "_set_value")             

    
    self.Link("_decode_V_KEY_SEP")                                          
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)                            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)                    
    self.Emit("CMPQ", _AX, jit.Imm(_S_obj_delim))                           
    self.Sjmp("JNE" , "_invalid_char")                                      
    self.Emit("MOVQ", jit.Imm(_S_val), jit.Sib(_ST, _CX, 8, _ST_Vt))        
    self.Emit("MOVQ", jit.Imm(_S_obj), jit.Sib(_ST, _CX, 8, _ST_Vt - 8))    
    self.Sjmp("JMP" , "_next")                                              

    
    self.Link("_decode_V_ELEM_SEP")                          
    self.Emit("MOVQ" , jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ" , jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)    
    self.Emit("CMPQ" , _AX, jit.Imm(_S_arr))      
    self.Sjmp("JE"   , "_array_sep")                         
    self.Emit("CMPQ" , _AX, jit.Imm(_S_obj))                 
    self.Sjmp("JNE"  , "_invalid_char")                      
    self.Emit("MOVQ" , jit.Imm(_S_obj_sep), jit.Sib(_ST, _CX, 8, _ST_Vt))
    self.Sjmp("JMP"  , "_next")                              

    
    self.Link("_array_sep")
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vp), _SI)    
    self.Emit("MOVQ", jit.Ptr(_SI, 8), _SI)                 
    self.Emit("MOVQ", jit.Ptr(_SI, 8), _DX)                 
    self.Emit("CMPQ", _DX, jit.Ptr(_SI, 16))                
    self.Sjmp("JAE" , "_array_more")                        

    
    self.Link("_array_append")                                          
    self.Emit("ADDQ", jit.Imm(1), jit.Ptr(_SI, 8))                      
    self.Emit("MOVQ", jit.Ptr(_SI, 0), _SI)                             
    self.Emit("ADDQ", jit.Imm(1), _CX)                                  
    self.Emit("CMPQ", _CX, jit.Imm(types.MAX_RECURSE))                  
    self.Sjmp("JAE"  , "_stack_overflow")                                
    self.Emit("SHLQ", jit.Imm(1), _DX)                                  
    self.Emit("LEAQ", jit.Sib(_SI, _DX, 8, 0), _SI)                     
    self.Emit("MOVQ", _CX, jit.Ptr(_ST, _ST_Sp))                        
    self.WriteRecNotAX(7 , _SI, jit.Sib(_ST, _CX, 8, _ST_Vp), false)           
    self.Emit("MOVQ", jit.Imm(_S_val), jit.Sib(_ST, _CX, 8, _ST_Vt))    
    self.Sjmp("JMP" , "_next")                                          

    
    self.Link("_decode_V_ARRAY_END")                        
    self.Emit("XORL", _DX, _DX)                             
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)    
    self.Emit("CMPQ", _AX, jit.Imm(_S_arr_0))               
    self.Sjmp("JE"  , "_first_item")                        
    self.Emit("CMPQ", _AX, jit.Imm(_S_arr))                 
    self.Sjmp("JNE" , "_invalid_char")                      
    self.Emit("SUBQ", jit.Imm(1), jit.Ptr(_ST, _ST_Sp))     
    self.Emit("MOVQ", _DX, jit.Sib(_ST, _CX, 8, _ST_Vp))    
    self.Sjmp("JMP" , "_next")                              

    
    self.Link("_first_item")                                    
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)                
    self.Emit("SUBQ", jit.Imm(2), jit.Ptr(_ST, _ST_Sp))         
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vp - 8), _SI)    
    self.Emit("MOVQ", jit.Ptr(_SI, 8), _SI)                     
    self.Emit("MOVQ", _DX, jit.Sib(_ST, _CX, 8, _ST_Vp - 8))    
    self.Emit("MOVQ", _DX, jit.Sib(_ST, _CX, 8, _ST_Vp))        
    self.Emit("MOVQ", _DX, jit.Ptr(_SI, 8))                     
    self.Sjmp("JMP" , "_next")                                  

    
    self.Link("_decode_V_OBJECT_END")                       
    self.Emit("MOVL", jit.Imm(_S_omask_end), _DI)           
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vt), _AX)    
    self.Emit("BTQ" , _AX, _DI)                    
    self.Sjmp("JNC" , "_invalid_char")                      
    self.Emit("XORL", _AX, _AX)                             
    self.Emit("SUBQ", jit.Imm(1), jit.Ptr(_ST, _ST_Sp))     
    self.Emit("MOVQ", _AX, jit.Sib(_ST, _CX, 8, _ST_Vp))    
    self.Sjmp("JMP" , "_next")                              

    
    self.Link("_return")                            
    self.Emit("XORL", _EP, _EP)                     
    self.Emit("MOVQ", _EP, jit.Ptr(_ST, _ST_Vp))    
    self.Link("_epilogue")                          
    self.Emit("SUBQ", jit.Imm(_FsmOffset), _ST)     
    self.Emit("MOVQ", jit.Ptr(_SP, _VD_offs), _BP)  
    self.Emit("ADDQ", jit.Imm(_VD_size), _SP)       
    self.Emit("RET")                                

    
    self.Link("_array_more")                    
    self.Emit("MOVQ" , _T_eface, _AX)           
    self.Emit("MOVQ" , jit.Ptr(_SI, 0), _BX)    
    self.Emit("MOVQ" , jit.Ptr(_SI, 8), _CX)    
    self.Emit("MOVQ" , jit.Ptr(_SI, 16), _DI)   
    self.Emit("MOVQ" , _DI, _SI)                
    self.Emit("SHLQ" , jit.Imm(1), _SI)         
    self.call_go(_F_growslice)                  
    self.Emit("MOVQ" , _AX, _DI)                
    self.Emit("MOVQ" , _BX, _DX)                
    self.Emit("MOVQ" , _CX, _AX)                
             
    
    self.Emit("MOVQ", jit.Ptr(_ST, _ST_Sp), _CX)            
    self.Emit("MOVQ", jit.Sib(_ST, _CX, 8, _ST_Vp), _SI)    
    self.Emit("MOVQ", jit.Ptr(_SI, 8), _SI)                 
    self.Emit("MOVQ", _DX, jit.Ptr(_SI, 8))                 
    self.Emit("MOVQ", _AX, jit.Ptr(_SI, 16))                
    self.WriteRecNotAX(8 , _DI, jit.Ptr(_SI, 0), false)                 
    self.Sjmp("JMP" , "_array_append")                      

    
    self.Link("copy_string")  
    self.Emit("MOVQ", _R8, _VAR_cs_p)
    self.Emit("MOVQ", _AX, _VAR_cs_n)
    self.Emit("MOVQ", _DI, _VAR_cs_LR)
    self.Emit("MOVQ", _AX, _BX)
    self.Emit("MOVQ", _AX, _CX)
    self.Emit("MOVQ", _T_byte, _AX)
    self.call_go(_F_makeslice)                              
    self.Emit("MOVQ", _AX, _VAR_cs_d)                    
    self.Emit("MOVQ", _VAR_cs_p, _BX)
    self.Emit("MOVQ", _VAR_cs_n, _CX)
    self.call_go(_F_memmove)
    self.Emit("MOVQ", _VAR_cs_d, _R8)
    self.Emit("MOVQ", _VAR_cs_n, _AX)
    self.Emit("MOVQ", _VAR_cs_LR, _DI)
    self.Rjmp("JMP", _DI)

    
    self.Link("_stack_overflow")
    self.Emit("MOVL" , _E_recurse, _EP)         
    self.Sjmp("JMP"  , "_error")                
    self.Link("_vtype_error")                   
    self.Emit("MOVQ" , _DI, _IC)                
    self.Emit("MOVL" , _E_invalid, _EP)         
    self.Sjmp("JMP"  , "_error")                
    self.Link("_invalid_char")                  
    self.Emit("SUBQ" , jit.Imm(1), _IC)         
    self.Emit("MOVL" , _E_invalid, _EP)         
    self.Sjmp("JMP"  , "_error")                
    self.Link("_unquote_error")                 
    self.Emit("MOVQ" , _VAR_ss_Iv, _IC)         
    self.Emit("SUBQ" , jit.Imm(1), _IC)         
    self.Link("_parsing_error")                 
    self.Emit("NEGQ" , _AX)                     
    self.Emit("MOVQ" , _AX, _EP)                
    self.Link("_error")                         
    self.Emit("PXOR" , _X0, _X0)                
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 0))    
    self.Sjmp("JMP"  , "_epilogue")             

    
    self.Link("_invalid_vtype")
    self.call_go(_F_invalid_vtype)                 
    self.Emit("UD2")                            

    
    self.Link("_switch_table")              
    self.Sref("_decode_V_EOF", 0)           
    self.Sref("_decode_V_NULL", -4)         
    self.Sref("_decode_V_TRUE", -8)         
    self.Sref("_decode_V_FALSE", -12)       
    self.Sref("_decode_V_ARRAY", -16)       
    self.Sref("_decode_V_OBJECT", -20)      
    self.Sref("_decode_V_STRING", -24)      
    self.Sref("_decode_V_DOUBLE", -28)      
    self.Sref("_decode_V_INTEGER", -32)     
    self.Sref("_decode_V_KEY_SEP", -36)     
    self.Sref("_decode_V_ELEM_SEP", -40)    
    self.Sref("_decode_V_ARRAY_END", -44)   
    self.Sref("_decode_V_OBJECT_END", -48)  

    
    self.Link("_decode_tab")        
    self.Sref("_decode_V_EOF", 0)   

    
    for i := 1; i < 256; i++ {
        if to, ok := _R_tab[i]; ok {
            self.Sref(to, -int64(i) * 4)
        } else {
            self.Byte(0x00, 0x00, 0x00, 0x00)
        }
    }
}



var (
    _subr_decode_value = new(_ValueDecoder).build()
)

//go:nosplit
func invalid_vtype(vt types.ValueType) {
    rt.Throw(fmt.Sprintf("invalid value type: %d", vt))
}
