//go:build go1.17 && !go1.25
// +build go1.17,!go1.25



package jitdec

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"github.com/bytedance/sonic/internal/caching"
	"github.com/bytedance/sonic/internal/jit"
	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/twitchyliquid64/golang-asm/obj"
)





const (
    _FP_args   = 72     
    _FP_fargs  = 80     
    _FP_saves  = 48     
    _FP_locals = 144    
)

const (
    _FP_offs = _FP_fargs + _FP_saves + _FP_locals
    _FP_size = _FP_offs + 8     
    _FP_base = _FP_size + 8     
)

const (
    _IM_null = 0x6c6c756e   
    _IM_true = 0x65757274   
    _IM_alse = 0x65736c61   
)

const (
    _BM_space = (1 << ' ') | (1 << '\t') | (1 << '\r') | (1 << '\n')
)

const (
    _MODE_JSON = 1 << 3 
)

const (
    _LB_error           = "_error"
    _LB_im_error        = "_im_error"
    _LB_eof_error       = "_eof_error"
    _LB_type_error      = "_type_error"
    _LB_field_error     = "_field_error"
    _LB_range_error     = "_range_error"
    _LB_stack_error     = "_stack_error"
    _LB_base64_error    = "_base64_error"
    _LB_unquote_error   = "_unquote_error"
    _LB_parsing_error   = "_parsing_error"
    _LB_parsing_error_v = "_parsing_error_v"
    _LB_mismatch_error   = "_mismatch_error"
)

const (
    _LB_char_0_error  = "_char_0_error"
    _LB_char_1_error  = "_char_1_error"
    _LB_char_2_error  = "_char_2_error"
    _LB_char_3_error  = "_char_3_error"
    _LB_char_4_error  = "_char_4_error"
    _LB_char_m2_error = "_char_m2_error"
    _LB_char_m3_error = "_char_m3_error"
)

const (
    _LB_skip_one = "_skip_one"
    _LB_skip_key_value = "_skip_key_value"
)

var (
    _AX = jit.Reg("AX")
    _BX = jit.Reg("BX")
    _CX = jit.Reg("CX")
    _DX = jit.Reg("DX")
    _DI = jit.Reg("DI")
    _SI = jit.Reg("SI")
    _BP = jit.Reg("BP")
    _SP = jit.Reg("SP")
    _R8 = jit.Reg("R8")
    _R9 = jit.Reg("R9")
    _X0 = jit.Reg("X0")
    _X1 = jit.Reg("X1")
    _X15 = jit.Reg("X15")
)

var (
    _IP = jit.Reg("R10")  
    _IC = jit.Reg("R11")  
    _IL = jit.Reg("R12")
    _ST = jit.Reg("R13")
    _VP = jit.Reg("R15")
)

var (
    _DF = jit.Reg("AX")    
    _ET = jit.Reg("AX")
    _EP = jit.Reg("BX")
)



var (
    _ARG_s  = _ARG_sp
    _ARG_sp = jit.Ptr(_SP, _FP_base + 0)
    _ARG_sl = jit.Ptr(_SP, _FP_base + 8)
    _ARG_ic = jit.Ptr(_SP, _FP_base + 16)
    _ARG_vp = jit.Ptr(_SP, _FP_base + 24)
    _ARG_sb = jit.Ptr(_SP, _FP_base + 32)
    _ARG_fv = jit.Ptr(_SP, _FP_base + 40)
)

var (
    _ARG_sv   = _ARG_sv_p
    _ARG_sv_p = jit.Ptr(_SP, _FP_base + 48)
    _ARG_sv_n = jit.Ptr(_SP, _FP_base + 56)
    _ARG_vk   = jit.Ptr(_SP, _FP_base + 64)
)

var (
    _VAR_st = _VAR_st_Vt
    _VAR_sr = jit.Ptr(_SP, _FP_fargs + _FP_saves)
)

var (
    _VAR_st_Vt = jit.Ptr(_SP, _FP_fargs + _FP_saves + 0)
    _VAR_st_Dv = jit.Ptr(_SP, _FP_fargs + _FP_saves + 8)
    _VAR_st_Iv = jit.Ptr(_SP, _FP_fargs + _FP_saves + 16)
    _VAR_st_Ep = jit.Ptr(_SP, _FP_fargs + _FP_saves + 24)
    _VAR_st_Db = jit.Ptr(_SP, _FP_fargs + _FP_saves + 32)
    _VAR_st_Dc = jit.Ptr(_SP, _FP_fargs + _FP_saves + 40)
)

var (
    _VAR_ss_AX = jit.Ptr(_SP, _FP_fargs + _FP_saves + 48)
    _VAR_ss_CX = jit.Ptr(_SP, _FP_fargs + _FP_saves + 56)
    _VAR_ss_SI = jit.Ptr(_SP, _FP_fargs + _FP_saves + 64)
    _VAR_ss_R8 = jit.Ptr(_SP, _FP_fargs + _FP_saves + 72)
    _VAR_ss_R9 = jit.Ptr(_SP, _FP_fargs + _FP_saves + 80)
)

var (
    _VAR_bs_p = jit.Ptr(_SP, _FP_fargs + _FP_saves + 88)
    _VAR_bs_n = jit.Ptr(_SP, _FP_fargs + _FP_saves + 96)
    _VAR_bs_LR = jit.Ptr(_SP, _FP_fargs + _FP_saves + 104)
)

var _VAR_fl = jit.Ptr(_SP, _FP_fargs + _FP_saves + 112)

var (
    _VAR_et = jit.Ptr(_SP, _FP_fargs + _FP_saves + 120) 
    _VAR_pc = jit.Ptr(_SP, _FP_fargs + _FP_saves + 128) 
    _VAR_ic = jit.Ptr(_SP, _FP_fargs + _FP_saves + 136) 
)

type _Assembler struct {
    jit.BaseAssembler
    p _Program
    name string
}

func newAssembler(p _Program) *_Assembler {
    return new(_Assembler).Init(p)
}



func (self *_Assembler) Load() _Decoder {
    return ptodec(self.BaseAssembler.Load("decode_"+self.name, _FP_size, _FP_args, argPtrs, localPtrs))
}

func (self *_Assembler) Init(p _Program) *_Assembler {
    self.p = p
    self.BaseAssembler.Init(self.compile)
    return self
}

func (self *_Assembler) compile() {
    self.prologue()
    self.instrs()
    self.epilogue()
    self.copy_string()
    self.escape_string()
    self.escape_string_twice()
    self.skip_one()
    self.skip_key_value()
    self.type_error()
    self.mismatch_error()
    self.field_error()
    self.range_error()
    self.stack_error()
    self.base64_error()
    self.parsing_error()
}



var _OpFuncTab = [256]func(*_Assembler, *_Instr) {
    _OP_any              : (*_Assembler)._asm_OP_any,
    _OP_dyn              : (*_Assembler)._asm_OP_dyn,
    _OP_str              : (*_Assembler)._asm_OP_str,
    _OP_bin              : (*_Assembler)._asm_OP_bin,
    _OP_bool             : (*_Assembler)._asm_OP_bool,
    _OP_num              : (*_Assembler)._asm_OP_num,
    _OP_i8               : (*_Assembler)._asm_OP_i8,
    _OP_i16              : (*_Assembler)._asm_OP_i16,
    _OP_i32              : (*_Assembler)._asm_OP_i32,
    _OP_i64              : (*_Assembler)._asm_OP_i64,
    _OP_u8               : (*_Assembler)._asm_OP_u8,
    _OP_u16              : (*_Assembler)._asm_OP_u16,
    _OP_u32              : (*_Assembler)._asm_OP_u32,
    _OP_u64              : (*_Assembler)._asm_OP_u64,
    _OP_f32              : (*_Assembler)._asm_OP_f32,
    _OP_f64              : (*_Assembler)._asm_OP_f64,
    _OP_unquote          : (*_Assembler)._asm_OP_unquote,
    _OP_nil_1            : (*_Assembler)._asm_OP_nil_1,
    _OP_nil_2            : (*_Assembler)._asm_OP_nil_2,
    _OP_nil_3            : (*_Assembler)._asm_OP_nil_3,
    _OP_empty_bytes      : (*_Assembler)._asm_OP_empty_bytes,
    _OP_deref            : (*_Assembler)._asm_OP_deref,
    _OP_index            : (*_Assembler)._asm_OP_index,
    _OP_is_null          : (*_Assembler)._asm_OP_is_null,
    _OP_is_null_quote    : (*_Assembler)._asm_OP_is_null_quote,
    _OP_map_init         : (*_Assembler)._asm_OP_map_init,
    _OP_map_key_i8       : (*_Assembler)._asm_OP_map_key_i8,
    _OP_map_key_i16      : (*_Assembler)._asm_OP_map_key_i16,
    _OP_map_key_i32      : (*_Assembler)._asm_OP_map_key_i32,
    _OP_map_key_i64      : (*_Assembler)._asm_OP_map_key_i64,
    _OP_map_key_u8       : (*_Assembler)._asm_OP_map_key_u8,
    _OP_map_key_u16      : (*_Assembler)._asm_OP_map_key_u16,
    _OP_map_key_u32      : (*_Assembler)._asm_OP_map_key_u32,
    _OP_map_key_u64      : (*_Assembler)._asm_OP_map_key_u64,
    _OP_map_key_f32      : (*_Assembler)._asm_OP_map_key_f32,
    _OP_map_key_f64      : (*_Assembler)._asm_OP_map_key_f64,
    _OP_map_key_str      : (*_Assembler)._asm_OP_map_key_str,
    _OP_map_key_utext    : (*_Assembler)._asm_OP_map_key_utext,
    _OP_map_key_utext_p  : (*_Assembler)._asm_OP_map_key_utext_p,
    _OP_array_skip       : (*_Assembler)._asm_OP_array_skip,
    _OP_array_clear      : (*_Assembler)._asm_OP_array_clear,
    _OP_array_clear_p    : (*_Assembler)._asm_OP_array_clear_p,
    _OP_slice_init       : (*_Assembler)._asm_OP_slice_init,
    _OP_slice_append     : (*_Assembler)._asm_OP_slice_append,
    _OP_object_next      : (*_Assembler)._asm_OP_object_next,
    _OP_struct_field     : (*_Assembler)._asm_OP_struct_field,
    _OP_unmarshal        : (*_Assembler)._asm_OP_unmarshal,
    _OP_unmarshal_p      : (*_Assembler)._asm_OP_unmarshal_p,
    _OP_unmarshal_text   : (*_Assembler)._asm_OP_unmarshal_text,
    _OP_unmarshal_text_p : (*_Assembler)._asm_OP_unmarshal_text_p,
    _OP_lspace           : (*_Assembler)._asm_OP_lspace,
    _OP_match_char       : (*_Assembler)._asm_OP_match_char,
    _OP_check_char       : (*_Assembler)._asm_OP_check_char,
    _OP_load             : (*_Assembler)._asm_OP_load,
    _OP_save             : (*_Assembler)._asm_OP_save,
    _OP_drop             : (*_Assembler)._asm_OP_drop,
    _OP_drop_2           : (*_Assembler)._asm_OP_drop_2,
    _OP_recurse          : (*_Assembler)._asm_OP_recurse,
    _OP_goto             : (*_Assembler)._asm_OP_goto,
    _OP_switch           : (*_Assembler)._asm_OP_switch,
    _OP_check_char_0     : (*_Assembler)._asm_OP_check_char_0,
    _OP_dismatch_err     : (*_Assembler)._asm_OP_dismatch_err,
    _OP_go_skip          : (*_Assembler)._asm_OP_go_skip,
    _OP_skip_emtpy       : (*_Assembler)._asm_OP_skip_empty,
    _OP_add              : (*_Assembler)._asm_OP_add,
    _OP_check_empty      : (*_Assembler)._asm_OP_check_empty,
    _OP_unsupported      : (*_Assembler)._asm_OP_unsupported,
    _OP_debug            : (*_Assembler)._asm_OP_debug,
}

func (self *_Assembler) _asm_OP_debug(_ *_Instr) {
    self.Byte(0xcc)
}

func (self *_Assembler) instr(v *_Instr) {
    if fn := _OpFuncTab[v.op()]; fn != nil {
        fn(self, v)
    } else {
        panic(fmt.Sprintf("invalid opcode: %d", v.op()))
    }
}

func (self *_Assembler) instrs() {
    for i, v := range self.p {
        self.Mark(i)
        self.instr(&v)
        self.debug_instr(i, &v)
    }
}

func (self *_Assembler) epilogue() {
    self.Mark(len(self.p))
    self.Emit("XORL", _EP, _EP)                     
    self.Emit("MOVQ", _VAR_et, _ET)                 
    self.Emit("TESTQ", _ET, _ET)                    
    self.Sjmp("JNZ", _LB_mismatch_error)            
    self.Link(_LB_error)                            
    self.Emit("MOVQ", _EP, _CX)                     
    self.Emit("MOVQ", _ET, _BX)                     
    self.Emit("MOVQ", _IC, _AX)                     
    self.Emit("MOVQ", jit.Imm(0), _ARG_sp)          
    self.Emit("MOVQ", jit.Imm(0), _ARG_vp)          
    self.Emit("MOVQ", jit.Imm(0), _ARG_sv_p)        
    self.Emit("MOVQ", jit.Imm(0), _ARG_vk)          
    self.Emit("MOVQ", jit.Ptr(_SP, _FP_offs), _BP)  
    self.Emit("ADDQ", jit.Imm(_FP_size), _SP)       
    self.Emit("RET")                                
}

func (self *_Assembler) prologue() {
    self.Emit("SUBQ", jit.Imm(_FP_size), _SP)       
    self.Emit("MOVQ", _BP, jit.Ptr(_SP, _FP_offs))  
    self.Emit("LEAQ", jit.Ptr(_SP, _FP_offs), _BP)  
    self.Emit("MOVQ", _AX, _ARG_sp)                 
    self.Emit("MOVQ", _AX, _IP)                     
    self.Emit("MOVQ", _BX, _ARG_sl)                 
    self.Emit("MOVQ", _BX, _IL)                     
    self.Emit("MOVQ", _CX, _ARG_ic)                 
    self.Emit("MOVQ", _CX, _IC)                     
    self.Emit("MOVQ", _DI, _ARG_vp)                 
    self.Emit("MOVQ", _DI, _VP)                     
    self.Emit("MOVQ", _SI, _ARG_sb)                 
    self.Emit("MOVQ", _SI, _ST)                     
    self.Emit("MOVQ", _R8, _ARG_fv)                 
    self.Emit("MOVQ", jit.Imm(0), _ARG_sv_p)        
    self.Emit("MOVQ", jit.Imm(0), _ARG_sv_n)        
    self.Emit("MOVQ", jit.Imm(0), _ARG_vk)          
    self.Emit("MOVQ", jit.Imm(0), _VAR_et)          
    
    self.Emit("MOVQ", jit.Imm(_MaxDigitNums), _VAR_st_Dc)    
    self.Emit("LEAQ", jit.Ptr(_ST, _DbufOffset), _AX)        
    self.Emit("MOVQ", _AX, _VAR_st_Db)                       
}



var (
    _REG_go = []obj.Addr { _ST, _VP, _IP, _IL, _IC }
    _REG_rt = []obj.Addr { _ST, _VP, _IP, _IL, _IC }
)

func (self *_Assembler) save(r ...obj.Addr) {
    for i, v := range r {
        if i > _FP_saves / 8 - 1 {
            panic("too many registers to save")
        } else {
            self.Emit("MOVQ", v, jit.Ptr(_SP, _FP_fargs + int64(i) * 8))
        }
    }
}

func (self *_Assembler) load(r ...obj.Addr) {
    for i, v := range r {
        if i > _FP_saves / 8 - 1 {
            panic("too many registers to load")
        } else {
            self.Emit("MOVQ", jit.Ptr(_SP, _FP_fargs + int64(i) * 8), v)
        }
    }
}

func (self *_Assembler) call(fn obj.Addr) {
    self.Emit("MOVQ", fn, _R9)  
    self.Rjmp("CALL", _R9)      
}

func (self *_Assembler) call_go(fn obj.Addr) {
    self.save(_REG_go...)   
    self.call(fn)
    self.load(_REG_go...)   
}

func (self *_Assembler) callc(fn obj.Addr) {
    self.save(_IP)
    self.call(fn)
    self.Emit("XORPS", _X15, _X15)
    self.load(_IP)
}

func (self *_Assembler) call_c(fn obj.Addr) {
    self.Emit("XCHGQ", _IC, _BX)
    self.callc(fn)
    self.Emit("XCHGQ", _IC, _BX)
}

func (self *_Assembler) call_sf(fn obj.Addr) {
    self.Emit("LEAQ", _ARG_s, _DI)                      
    self.Emit("MOVQ", _IC, _ARG_ic)                     
    self.Emit("LEAQ", _ARG_ic, _SI)                     
    self.Emit("LEAQ", jit.Ptr(_ST, _FsmOffset), _DX)    
    self.Emit("MOVQ", _ARG_fv, _CX)
    self.callc(fn)
    self.Emit("MOVQ", _ARG_ic, _IC)                     
}

func (self *_Assembler) call_vf(fn obj.Addr) {
    self.Emit("LEAQ", _ARG_s, _DI)      
    self.Emit("MOVQ", _IC, _ARG_ic)     
    self.Emit("LEAQ", _ARG_ic, _SI)     
    self.Emit("LEAQ", _VAR_st, _DX)     
    self.callc(fn)
    self.Emit("MOVQ", _ARG_ic, _IC)     
}



var (
    _F_convT64        = jit.Func(rt.ConvT64)
    _F_error_wrap     = jit.Func(error_wrap)
    _F_error_type     = jit.Func(error_type)
    _F_error_field    = jit.Func(error_field)
    _F_error_value    = jit.Func(error_value)
    _F_error_mismatch = jit.Func(error_mismatch)
)

var (
    _I_int8    , _T_int8    = rtype(reflect.TypeOf(int8(0)))
    _I_int16   , _T_int16   = rtype(reflect.TypeOf(int16(0)))
    _I_int32   , _T_int32   = rtype(reflect.TypeOf(int32(0)))
    _I_uint8   , _T_uint8   = rtype(reflect.TypeOf(uint8(0)))
    _I_uint16  , _T_uint16  = rtype(reflect.TypeOf(uint16(0)))
    _I_uint32  , _T_uint32  = rtype(reflect.TypeOf(uint32(0)))
    _I_float32 , _T_float32 = rtype(reflect.TypeOf(float32(0)))
)

var (
    _T_error                    = rt.UnpackType(errorType)
    _I_base64_CorruptInputError = jit.Itab(_T_error, base64CorruptInputError)
)

var (
    _V_stackOverflow              = jit.Imm(int64(uintptr(unsafe.Pointer(&stackOverflow))))
    _I_json_UnsupportedValueError = jit.Itab(_T_error, reflect.TypeOf(new(json.UnsupportedValueError)))
    _I_json_MismatchTypeError     = jit.Itab(_T_error, reflect.TypeOf(new(MismatchTypeError)))
    _I_json_MismatchQuotedError   = jit.Itab(_T_error, reflect.TypeOf(new(MismatchQuotedError)))
)

func (self *_Assembler) type_error() {
    self.Link(_LB_type_error)                   
    self.call_go(_F_error_type)                 
    self.Sjmp("JMP" , _LB_error)                
}

func (self *_Assembler) mismatch_error() {
    self.Link(_LB_mismatch_error)                     
    self.Emit("MOVQ", _VAR_et, _ET)                   
    self.Emit("MOVQ", _I_json_MismatchTypeError, _CX) 
    self.Emit("CMPQ", _ET, _CX)                       
    self.Emit("MOVQ", jit.Ptr(_ST, _EpOffset), _EP)   
    self.Sjmp("JE"  , _LB_error)                      
    self.Emit("MOVQ", _ARG_sp, _AX)
    self.Emit("MOVQ", _ARG_sl, _BX)
    self.Emit("MOVQ", _VAR_ic, _CX)
    self.Emit("MOVQ", _VAR_et, _DI)
    self.call_go(_F_error_mismatch)             
    self.Sjmp("JMP" , _LB_error)                
}

func (self *_Assembler) field_error() {
    self.Link(_LB_field_error)                  
    self.Emit("MOVQ", _ARG_sv_p, _AX)           
    self.Emit("MOVQ", _ARG_sv_n, _BX)           
    self.call_go(_F_error_field)                
    self.Sjmp("JMP" , _LB_error)                
}

func (self *_Assembler) range_error() {
    self.Link(_LB_range_error)                  
    self.Emit("MOVQ", _ET, _CX)                 
    self.slice_from(_VAR_st_Ep, 0)              
    self.Emit("MOVQ", _DI, _AX)                 
    self.Emit("MOVQ", _EP, _DI)                 
    self.Emit("MOVQ", _SI, _BX)                 
    self.call_go(_F_error_value)                
    self.Sjmp("JMP" , _LB_error)                
}

func (self *_Assembler) stack_error() {
    self.Link(_LB_stack_error)                              
    self.Emit("MOVQ", _V_stackOverflow, _EP)                
    self.Emit("MOVQ", _I_json_UnsupportedValueError, _ET)   
    self.Sjmp("JMP" , _LB_error)                            
}

func (self *_Assembler) base64_error() {
    self.Link(_LB_base64_error)
    self.Emit("NEGQ", _AX)                                  
    self.Emit("SUBQ", jit.Imm(1), _AX)                      
    self.call_go(_F_convT64)                                
    self.Emit("MOVQ", _AX, _EP)                             
    self.Emit("MOVQ", _I_base64_CorruptInputError, _ET)     
    self.Sjmp("JMP" , _LB_error)                            
}

func (self *_Assembler) parsing_error() {
    self.Link(_LB_eof_error)                                            
    self.Emit("MOVQ" , _IL, _IC)                                        
    self.Emit("MOVL" , jit.Imm(int64(types.ERR_EOF)), _EP)              
    self.Sjmp("JMP"  , _LB_parsing_error)                               
    self.Link(_LB_unquote_error)                                        
    self.Emit("SUBQ" , _VAR_sr, _SI)                                    
    self.Emit("SUBQ" , _SI, _IC)                                        
    self.Link(_LB_parsing_error_v)                                      
    self.Emit("MOVQ" , _AX, _EP)                                        
    self.Emit("NEGQ" , _EP)                                             
    self.Sjmp("JMP"  , _LB_parsing_error)                               
    self.Link(_LB_char_m3_error)                                        
    self.Emit("SUBQ" , jit.Imm(1), _IC)                                 
    self.Link(_LB_char_m2_error)                                        
    self.Emit("SUBQ" , jit.Imm(2), _IC)                                 
    self.Sjmp("JMP"  , _LB_char_0_error)                                
    self.Link(_LB_im_error)                                             
    self.Emit("CMPB" , _CX, jit.Sib(_IP, _IC, 1, 0))                    
    self.Sjmp("JNE"  , _LB_char_0_error)                                
    self.Emit("SHRL" , jit.Imm(8), _CX)                                 
    self.Emit("CMPB" , _CX, jit.Sib(_IP, _IC, 1, 1))                    
    self.Sjmp("JNE"  , _LB_char_1_error)                                
    self.Emit("SHRL" , jit.Imm(8), _CX)                                 
    self.Emit("CMPB" , _CX, jit.Sib(_IP, _IC, 1, 2))                    
    self.Sjmp("JNE"  , _LB_char_2_error)                                
    self.Sjmp("JMP"  , _LB_char_3_error)                                
    self.Link(_LB_char_4_error)                                         
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 
    self.Link(_LB_char_3_error)                                         
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 
    self.Link(_LB_char_2_error)                                         
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 
    self.Link(_LB_char_1_error)                                         
    self.Emit("ADDQ" , jit.Imm(1), _IC)                                 
    self.Link(_LB_char_0_error)                                         
    self.Emit("MOVL" , jit.Imm(int64(types.ERR_INVALID_CHAR)), _EP)     
    self.Link(_LB_parsing_error)                                        
    self.Emit("MOVQ" , _EP, _DI)                                        
    self.Emit("MOVQ",  _ARG_sp, _AX)                                     
    self.Emit("MOVQ",  _ARG_sl, _BX)                                     
    self.Emit("MOVQ" , _IC, _CX)                                        
    self.call_go(_F_error_wrap)                                         
    self.Sjmp("JMP"  , _LB_error)                                       
}

func (self *_Assembler) _asm_OP_dismatch_err(p *_Instr) {
    self.Emit("MOVQ", _IC, _VAR_ic)      
    self.Emit("MOVQ", jit.Type(p.vt()), _ET)     
    self.Emit("MOVQ", _ET, _VAR_et)
}

func (self *_Assembler) _asm_OP_go_skip(p *_Instr) {
    self.Byte(0x4c, 0x8d, 0x0d)         
    self.Xref(p.vi(), 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)            
}

var _F_IndexByte = jit.Func(strings.IndexByte)

func (self *_Assembler) _asm_OP_skip_empty(p *_Instr) {
    self.call_sf(_F_skip_one)                   
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_parsing_error_v)     
    self.Emit("BTQ", jit.Imm(_F_disable_unknown), _ARG_fv) 
    self.Xjmp("JNC", p.vi())
    self.Emit("LEAQ", jit.Sib(_IC, _AX, 1, 0), _BX)
    self.Emit("MOVQ", _BX, _ARG_sv_n)
    self.Emit("LEAQ", jit.Sib(_IP, _AX, 1, 0), _AX)
    self.Emit("MOVQ", _AX, _ARG_sv_p)
    self.Emit("MOVQ", jit.Imm(':'), _CX)
    self.call_go(_F_IndexByte)
    self.Emit("TESTQ", _AX, _AX)
    
    self.Sjmp("JNS", _LB_field_error)
}

func (self *_Assembler) skip_one() {
    self.Link(_LB_skip_one)                     
    self.Emit("MOVQ", _VAR_ic, _IC)             
    self.call_sf(_F_skip_one)                   
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_parsing_error_v)     
    self.Emit("MOVQ" , _VAR_pc, _R9)            
    self.Rjmp("JMP"  , _R9)                     
}

func (self *_Assembler) skip_key_value() {
    self.Link(_LB_skip_key_value)               
    
    self.Emit("MOVQ", _VAR_ic, _IC)             
    self.call_sf(_F_skip_one)                   
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_parsing_error_v)     
    
    self.lspace("_global_1")
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(':'))
    self.Sjmp("JNE"  , _LB_parsing_error_v)     
    self.Emit("ADDQ", jit.Imm(1), _IC)          
    self.lspace("_global_2")
    
    self.call_sf(_F_skip_one)                   
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_parsing_error_v)     
    
    self.Emit("MOVQ" , _VAR_pc, _R9)            
    self.Rjmp("JMP"  , _R9)                     
}




var (
    _T_byte     = jit.Type(byteType)
    _F_mallocgc = jit.Func(rt.Mallocgc)
)

func (self *_Assembler) malloc_AX(nb obj.Addr, ret obj.Addr) {
    self.Emit("MOVQ", nb, _AX)                  
    self.Emit("MOVQ", _T_byte, _BX)             
    self.Emit("XORL", _CX, _CX)                 
    self.call_go(_F_mallocgc)                   
    self.Emit("MOVQ", _AX, ret)                 
}

func (self *_Assembler) valloc(vt reflect.Type, ret obj.Addr) {
    self.Emit("MOVQ", jit.Imm(int64(vt.Size())), _AX)   
    self.Emit("MOVQ", jit.Type(vt), _BX)                
    self.Emit("MOVB", jit.Imm(1), _CX)                  
    self.call_go(_F_mallocgc)                           
    self.Emit("MOVQ", _AX, ret)                         
}

func (self *_Assembler) valloc_AX(vt reflect.Type) {
    self.Emit("MOVQ", jit.Imm(int64(vt.Size())), _AX)   
    self.Emit("MOVQ", jit.Type(vt), _BX)                
    self.Emit("MOVB", jit.Imm(1), _CX)                  
    self.call_go(_F_mallocgc)                           
}

func (self *_Assembler) vfollow(vt reflect.Type) {
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _AX)    
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JNZ"  , "_end_{n}")              
    self.valloc_AX(vt)                          
    self.WritePtrAX(1, jit.Ptr(_VP, 0), true)   
    self.Link("_end_{n}")                       
    self.Emit("MOVQ" , _AX, _VP)                
}



var (
    _F_vstring   = jit.Imm(int64(native.S_vstring))
    _F_vnumber   = jit.Imm(int64(native.S_vnumber))
    _F_vsigned   = jit.Imm(int64(native.S_vsigned))
    _F_vunsigned = jit.Imm(int64(native.S_vunsigned))
)

func (self *_Assembler) check_err(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ" , _VAR_st_Vt, _AX)         
    self.Emit("TESTQ", _AX, _AX)                
    
    if vt != nil {
        self.Sjmp("JNS" , "_check_err_{n}")        
        self.Emit("MOVQ", jit.Type(vt), _ET)         
        self.Emit("MOVQ", _ET, _VAR_et)
        if pin2 != -1 {
            self.Emit("SUBQ", jit.Imm(1), _BX)
            self.Emit("MOVQ", _BX, _VAR_ic)
            self.Byte(0x4c  , 0x8d, 0x0d)         
            self.Xref(pin2, 4)
            self.Emit("MOVQ", _R9, _VAR_pc)
            self.Sjmp("JMP" , _LB_skip_key_value)
        } else {
            self.Emit("MOVQ", _BX, _VAR_ic)
            self.Byte(0x4c  , 0x8d, 0x0d)         
            self.Sref(pin, 4)
            self.Emit("MOVQ", _R9, _VAR_pc)
            self.Sjmp("JMP" , _LB_skip_one)
        }
        self.Link("_check_err_{n}")
    } else {
        self.Sjmp("JS"   , _LB_parsing_error_v)     
    }
}

func (self *_Assembler) check_eof(d int64) {
    if d == 1 {
        self.Emit("CMPQ", _IC, _IL)         
        self.Sjmp("JAE" , _LB_eof_error)    
    } else {
        self.Emit("LEAQ", jit.Ptr(_IC, d), _AX)     
        self.Emit("CMPQ", _AX, _IL)                 
        self.Sjmp("JA"  , _LB_eof_error)            
    }
}


func (self *_Assembler) parse_string() {
    self.Emit("MOVQ", _ARG_fv, _CX)
    self.call_vf(_F_vstring)
    self.check_err(nil, "", -1)
}

func (self *_Assembler) parse_number(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ", _IC, _BX)       
    self.call_vf(_F_vnumber)
    self.check_err(vt, pin, pin2)
}

func (self *_Assembler) parse_signed(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ", _IC, _BX)       
    self.call_vf(_F_vsigned)
    self.check_err(vt, pin, pin2)
}

func (self *_Assembler) parse_unsigned(vt reflect.Type, pin string, pin2 int) {
    self.Emit("MOVQ", _IC, _BX)       
    self.call_vf(_F_vunsigned)
    self.check_err(vt, pin, pin2)
}


func (self *_Assembler) copy_string() {
    self.Link("_copy_string")
    self.Emit("MOVQ", _DI, _VAR_bs_p)
    self.Emit("MOVQ", _SI, _VAR_bs_n)
    self.Emit("MOVQ", _R9, _VAR_bs_LR)
    self.malloc_AX(_SI, _ARG_sv_p)                              
    self.Emit("MOVQ", _VAR_bs_p, _BX)
    self.Emit("MOVQ", _VAR_bs_n, _CX)
    self.call_go(_F_memmove)
    self.Emit("MOVQ", _ARG_sv_p, _DI)
    self.Emit("MOVQ", _VAR_bs_n, _SI)
    self.Emit("MOVQ", _VAR_bs_LR, _R9)
    self.Rjmp("JMP", _R9)
}


func (self *_Assembler) escape_string() {
    self.Link("_escape_string")
    self.Emit("MOVQ" , _DI, _VAR_bs_p)
    self.Emit("MOVQ" , _SI, _VAR_bs_n)
    self.Emit("MOVQ" , _R9, _VAR_bs_LR)
    self.malloc_AX(_SI, _DX)                                    
    self.Emit("MOVQ" , _DX, _ARG_sv_p)
    self.Emit("MOVQ" , _VAR_bs_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_n, _SI)                                  
    self.Emit("LEAQ" , _VAR_sr, _CX)                            
    self.Emit("XORL" , _R8, _R8)                                
    self.Emit("BTQ"  , jit.Imm(_F_disable_urc), _ARG_fv)        
    self.Emit("SETCC", _R8)                                     
    self.Emit("SHLQ" , jit.Imm(types.B_UNICODE_REPLACE), _R8)   
    self.call_c(_F_unquote)                                       
    self.Emit("MOVQ" , _VAR_bs_n, _SI)                                  
    self.Emit("ADDQ" , jit.Imm(1), _SI)                         
    self.Emit("TESTQ", _AX, _AX)                                
    self.Sjmp("JS"   , _LB_unquote_error)                       
    self.Emit("MOVQ" , _AX, _SI)
    self.Emit("MOVQ" , _ARG_sv_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_LR, _R9)
    self.Rjmp("JMP", _R9)
}

func (self *_Assembler) escape_string_twice() {
    self.Link("_escape_string_twice")
    self.Emit("MOVQ" , _DI, _VAR_bs_p)
    self.Emit("MOVQ" , _SI, _VAR_bs_n)
    self.Emit("MOVQ" , _R9, _VAR_bs_LR)
    self.malloc_AX(_SI, _DX)                                        
    self.Emit("MOVQ" , _DX, _ARG_sv_p)
    self.Emit("MOVQ" , _VAR_bs_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_n, _SI)        
    self.Emit("LEAQ" , _VAR_sr, _CX)                                
    self.Emit("MOVL" , jit.Imm(types.F_DOUBLE_UNQUOTE), _R8)        
    self.Emit("BTQ"  , jit.Imm(_F_disable_urc), _ARG_fv)            
    self.Emit("XORL" , _AX, _AX)                                    
    self.Emit("SETCC", _AX)                                         
    self.Emit("SHLQ" , jit.Imm(types.B_UNICODE_REPLACE), _AX)       
    self.Emit("ORQ"  , _AX, _R8)                                    
    self.call_c(_F_unquote)                                         
    self.Emit("MOVQ" , _VAR_bs_n, _SI)                              
    self.Emit("ADDQ" , jit.Imm(3), _SI)                             
    self.Emit("TESTQ", _AX, _AX)                                    
    self.Sjmp("JS"   , _LB_unquote_error)                           
    self.Emit("MOVQ" , _AX, _SI)
    self.Emit("MOVQ" , _ARG_sv_p, _DI)
    self.Emit("MOVQ" , _VAR_bs_LR, _R9)
    self.Rjmp("JMP", _R9)
}



var (
    _V_max_f32 = jit.Imm(int64(uintptr(unsafe.Pointer(_Vp_max_f32))))
    _V_min_f32 = jit.Imm(int64(uintptr(unsafe.Pointer(_Vp_min_f32))))
)

var (
    _Vp_max_f32 = new(float32)
    _Vp_min_f32 = new(float32)
)

func init() {
    *_Vp_max_f32 = math.MaxFloat32
    *_Vp_min_f32 = -math.MaxFloat32
}

func (self *_Assembler) range_single_X0() {
    self.Emit("CVTSD2SS", _VAR_st_Dv, _X0)              
    self.Emit("MOVQ"    , _V_max_f32, _CX)              
    self.Emit("MOVQ"    , jit.Gitab(_I_float32), _ET)   
    self.Emit("MOVQ"    , jit.Gtype(_T_float32), _EP)   
    self.Emit("UCOMISS" , jit.Ptr(_CX, 0), _X0)         
    self.Sjmp("JA"      , _LB_range_error)              
    self.Emit("MOVQ"    , _V_min_f32, _CX)              
    self.Emit("UCOMISS" , jit.Ptr(_CX, 0), _X0)         
    self.Sjmp("JB"      , _LB_range_error)              
}

func (self *_Assembler) range_signed_CX(i *rt.GoItab, t *rt.GoType, a int64, b int64) {
    self.Emit("MOVQ", _VAR_st_Iv, _CX)      
    self.Emit("MOVQ", jit.Gitab(i), _ET)    
    self.Emit("MOVQ", jit.Gtype(t), _EP)    
    self.Emit("CMPQ", _CX, jit.Imm(a))      
    self.Sjmp("JL"  , _LB_range_error)      
    self.Emit("CMPQ", _CX, jit.Imm(b))      
    self.Sjmp("JG"  , _LB_range_error)      
}

func (self *_Assembler) range_unsigned_CX(i *rt.GoItab, t *rt.GoType, v uint64) {
    self.Emit("MOVQ" , _VAR_st_Iv, _CX)         
    self.Emit("MOVQ" , jit.Gitab(i), _ET)       
    self.Emit("MOVQ" , jit.Gtype(t), _EP)       
    self.Emit("TESTQ", _CX, _CX)                
    self.Sjmp("JS"   , _LB_range_error)         
    self.Emit("CMPQ" , _CX, jit.Imm(int64(v)))  
    self.Sjmp("JA"   , _LB_range_error)         
}



var (
    _F_unquote = jit.Imm(int64(native.S_unquote))
)

func (self *_Assembler) slice_from(p obj.Addr, d int64) {
    self.Emit("MOVQ", p, _SI)   
    self.slice_from_r(_SI, d)   
}

func (self *_Assembler) slice_from_r(p obj.Addr, d int64) {
    self.Emit("LEAQ", jit.Sib(_IP, p, 1, 0), _DI)   
    self.Emit("NEGQ", p)                            
    self.Emit("LEAQ", jit.Sib(_IC, p, 1, d), _SI)   
}

func (self *_Assembler) unquote_once(p obj.Addr, n obj.Addr, stack bool, copy bool) {
    self.slice_from(_VAR_st_Iv, -1)                             
    self.Emit("CMPQ", _VAR_st_Ep, jit.Imm(-1))                 
    self.Sjmp("JE"  , "_noescape_{n}")                         
    self.Byte(0x4c, 0x8d, 0x0d)         
    self.Sref("_unquote_once_write_{n}", 4)
    self.Sjmp("JMP" , "_escape_string")
    self.Link("_noescape_{n}")
    if copy {
        self.Emit("BTQ" , jit.Imm(_F_copy_string), _ARG_fv)    
        self.Sjmp("JNC", "_unquote_once_write_{n}")
        self.Byte(0x4c, 0x8d, 0x0d)         
        self.Sref("_unquote_once_write_{n}", 4)
        self.Sjmp("JMP", "_copy_string")
    }
    self.Link("_unquote_once_write_{n}")
    self.Emit("MOVQ", _SI, n)                                  
    if stack {
        self.Emit("MOVQ", _DI, p) 
    } else {
        self.WriteRecNotAX(10, _DI, p, false, false)
    }
}

func (self *_Assembler) unquote_twice(p obj.Addr, n obj.Addr, stack bool) {
    self.Emit("CMPQ" , _VAR_st_Ep, jit.Imm(-1))                     
    self.Sjmp("JE"   , _LB_eof_error)                               
    self.Emit("CMPB" , jit.Sib(_IP, _IC, 1, -3), jit.Imm('\\'))     
    self.Sjmp("JNE"  , _LB_char_m3_error)                           
    self.Emit("CMPB" , jit.Sib(_IP, _IC, 1, -2), jit.Imm('"'))      
    self.Sjmp("JNE"  , _LB_char_m2_error)                           
    self.slice_from(_VAR_st_Iv, -3)                                 
    self.Emit("MOVQ" , _SI, _AX)                                    
    self.Emit("ADDQ" , _VAR_st_Iv, _AX)                             
    self.Emit("CMPQ" , _VAR_st_Ep, _AX)                             
    self.Sjmp("JE"   , "_noescape_{n}")                             
    self.Byte(0x4c, 0x8d, 0x0d)         
    self.Sref("_unquote_twice_write_{n}", 4)
    self.Sjmp("JMP" , "_escape_string_twice")
    self.Link("_noescape_{n}")                                      
    self.Emit("BTQ"  , jit.Imm(_F_copy_string), _ARG_fv)    
    self.Sjmp("JNC", "_unquote_twice_write_{n}") 
    self.Byte(0x4c, 0x8d, 0x0d)         
    self.Sref("_unquote_twice_write_{n}", 4)
    self.Sjmp("JMP", "_copy_string")
    self.Link("_unquote_twice_write_{n}")
    self.Emit("MOVQ" , _SI, n)                                      
    if stack {
        self.Emit("MOVQ", _DI, p) 
    } else {
        self.WriteRecNotAX(12, _DI, p, false, false)
    }
    self.Link("_unquote_twice_end_{n}")
}



var (
    _F_memclrHasPointers    = jit.Func(rt.MemclrHasPointers)
    _F_memclrNoHeapPointers = jit.Func(rt.MemclrNoHeapPointers)
)

func (self *_Assembler) mem_clear_fn(ptrfree bool) {
    if !ptrfree {
        self.call_go(_F_memclrHasPointers)
    } else {
        self.call_go(_F_memclrNoHeapPointers)
    }
}

func (self *_Assembler) mem_clear_rem(size int64, ptrfree bool) {
    self.Emit("MOVQ", jit.Imm(size), _BX)               
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             
    self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 0), _AX)     
    self.Emit("SUBQ", _VP, _AX)                         
    self.Emit("ADDQ", _AX, _BX)                         
    self.Emit("MOVQ", _VP, _AX)                         
    self.mem_clear_fn(ptrfree)                          
}



var (
    _F_mapassign           = jit.Func(rt.Mapassign)
    _F_mapassign_fast32    = jit.Func(rt.Mapassign_fast32)
    _F_mapassign_faststr   = jit.Func(rt.Mapassign_faststr)
    _F_mapassign_fast64ptr = jit.Func(rt.Mapassign_fast64ptr)
)

var (
    _F_decodeJsonUnmarshaler obj.Addr
    _F_decodeJsonUnmarshalerQuoted obj.Addr
    _F_decodeTextUnmarshaler obj.Addr
)

func init() {
    _F_decodeJsonUnmarshaler = jit.Func(decodeJsonUnmarshaler)
    _F_decodeJsonUnmarshalerQuoted = jit.Func(decodeJsonUnmarshalerQuoted)
    _F_decodeTextUnmarshaler = jit.Func(decodeTextUnmarshaler)
}

func (self *_Assembler) mapaccess_ptr(t reflect.Type) {
    if rt.MapType(rt.UnpackType(t)).IndirectElem() {
        self.vfollow(t.Elem())
    }
}

func (self *_Assembler) mapassign_std(t reflect.Type, v obj.Addr) {
    self.Emit("LEAQ", v, _AX)               
    self.mapassign_call_from_AX(t, _F_mapassign)    
}

func (self *_Assembler) mapassign_str_fast(t reflect.Type, p obj.Addr, n obj.Addr) {
    self.Emit("MOVQ", jit.Type(t), _AX)         
    self.Emit("MOVQ", _VP, _BX)                 
    self.Emit("MOVQ", p, _CX)                   
    self.Emit("MOVQ", n, _DI)                   
    self.call_go(_F_mapassign_faststr)          
    self.Emit("MOVQ", _AX, _VP)                 
    self.mapaccess_ptr(t)
}

func (self *_Assembler) mapassign_call_from_AX(t reflect.Type, fn obj.Addr) {
    self.Emit("MOVQ", _AX, _CX)
    self.Emit("MOVQ", jit.Type(t), _AX)         
    self.Emit("MOVQ", _VP, _BX)                 
    self.call_go(fn)                            
    self.Emit("MOVQ", _AX, _VP)                 
}

func (self *_Assembler) mapassign_fastx(t reflect.Type, fn obj.Addr) {
    self.mapassign_call_from_AX(t, fn)
    self.mapaccess_ptr(t)
}

func (self *_Assembler) mapassign_utext(t reflect.Type, addressable bool) {
    pv := false
    vk := t.Key()
    tk := t.Key()

    
    if vk.Kind() == reflect.Ptr {
        pv = true
        vk = vk.Elem()
    }

    
    if addressable {
        pv = false
        tk = reflect.PtrTo(tk)
    }

    
    self.valloc(vk, _BX)                        
    
    self.Emit("MOVQ" , _BX, _ARG_vk)
    self.Emit("MOVQ" , jit.Type(tk), _AX)       
    self.Emit("MOVQ" , _ARG_sv_p, _CX)          
    self.Emit("MOVQ" , _ARG_sv_n, _DI)          
    self.call_go(_F_decodeTextUnmarshaler)      
    self.Emit("TESTQ", _ET, _ET)                
    self.Sjmp("JNZ"  , _LB_error)               
    self.Emit("MOVQ" , _ARG_vk, _AX)            
    self.Emit("MOVQ", jit.Imm(0), _ARG_vk)

    
    if !pv {
        self.mapassign_call_from_AX(t, _F_mapassign)
    } else {
        self.mapassign_fastx(t, _F_mapassign_fast64ptr)
    }
}



var (
    _F_skip_one = jit.Imm(int64(native.S_skip_one))
    _F_skip_array  = jit.Imm(int64(native.S_skip_array))
    _F_skip_number = jit.Imm(int64(native.S_skip_number))
)

func (self *_Assembler) unmarshal_json(t reflect.Type, deref bool, f obj.Addr) {
    self.call_sf(_F_skip_one)                                   
    self.Emit("TESTQ", _AX, _AX)                                
    self.Sjmp("JS"   , _LB_parsing_error_v)                     
    self.Emit("MOVQ", _IC, _VAR_ic)                             
    self.slice_from_r(_AX, 0)                                   
    self.Emit("MOVQ" , _DI, _ARG_sv_p)                          
    self.Emit("MOVQ" , _SI, _ARG_sv_n)                          
    self.unmarshal_func(t, f, deref)     
}

func (self *_Assembler) unmarshal_text(t reflect.Type, deref bool) {
    self.parse_string()                                         
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)        
    self.unmarshal_func(t, _F_decodeTextUnmarshaler, deref)     
}

func (self *_Assembler) unmarshal_func(t reflect.Type, fn obj.Addr, deref bool) {
    pt := t
    vk := t.Kind()

    
    if deref && vk == reflect.Ptr {
        self.Emit("MOVQ" , _VP, _BX)                
        self.Emit("MOVQ" , jit.Ptr(_BX, 0), _BX)    
        self.Emit("TESTQ", _BX, _BX)                
        self.Sjmp("JNZ"  , "_deref_{n}")            
        self.valloc(t.Elem(), _BX)                  
        self.WriteRecNotAX(3, _BX, jit.Ptr(_VP, 0), false, false)    
        self.Link("_deref_{n}")                     
    } else {
        
        self.Emit("MOVQ", _VP, _BX)                 
    }

    
    self.Emit("MOVQ", jit.Type(pt), _AX)        

    
    self.Emit("MOVQ" , _ARG_sv_p, _CX)          
    self.Emit("MOVQ" , _ARG_sv_n, _DI)          
    self.call_go(fn)                            
    self.Emit("TESTQ", _ET, _ET)                
    if fn == _F_decodeJsonUnmarshalerQuoted {
        self.Sjmp("JZ"  , "_unmarshal_func_end_{n}")            
        self.Emit("MOVQ", _I_json_MismatchQuotedError, _CX)     
        self.Emit("CMPQ", _ET, _CX)            
        self.Sjmp("JNE" , _LB_error)           
        self.Emit("MOVQ", jit.Type(t), _CX)    
        self.Emit("MOVQ", _CX, _VAR_et)        
        self.Emit("MOVQ", _VAR_ic, _IC)        
        self.Emit("XORL", _ET, _ET)            
        self.Link("_unmarshal_func_end_{n}")
    } else {
        self.Sjmp("JNE" , _LB_error)           
    }
}



var (
    _F_decodeTypedPointer obj.Addr
)

func init() {
    _F_decodeTypedPointer = jit.Func(decodeTypedPointer)
}

func (self *_Assembler) decode_dynamic(vt obj.Addr, vp obj.Addr) {
    self.Emit("MOVQ" , vp, _SI)    
    self.Emit("MOVQ" , vt, _DI)    
    self.Emit("MOVQ", _ARG_sp, _AX)            
    self.Emit("MOVQ", _ARG_sl, _BX)            
    self.Emit("MOVQ" , _IC, _CX)                
    self.Emit("MOVQ" , _ST, _R8)                
    self.Emit("MOVQ" , _ARG_fv, _R9)            
    self.save(_REG_rt...)
    self.Emit("MOVQ", _F_decodeTypedPointer, _IL)  
    self.Rjmp("CALL", _IL)      
    self.load(_REG_rt...)
    self.Emit("MOVQ" , _AX, _IC)                
    self.Emit("MOVQ" , _BX, _ET)                
    self.Emit("MOVQ" , _CX, _EP)                
    self.Emit("TESTQ", _ET, _ET)                
    self.Sjmp("JE", "_decode_dynamic_end_{n}")  
    self.Emit("MOVQ", _I_json_MismatchTypeError, _CX) 
    self.Emit("CMPQ", _ET, _CX)                 
    self.Sjmp("JNE",  _LB_error)                
    self.Emit("MOVQ", _ET, _VAR_et)             
    self.WriteRecNotAX(14, _EP, jit.Ptr(_ST, _EpOffset), false, false) 
    self.Link("_decode_dynamic_end_{n}")
}



var (
    _F_memequal         = jit.Func(rt.MemEqual)
    _F_memmove          = jit.Func(rt.Memmove)
    _F_growslice        = jit.Func(rt.GrowSlice)
    _F_makeslice        = jit.Func(rt.MakeSliceStd)
    _F_makemap_small    = jit.Func(rt.MakemapSmall)
    _F_mapassign_fast64 = jit.Func(rt.Mapassign_fast64)
)

var (
    _F_lspace  = jit.Imm(int64(native.S_lspace))
    _F_strhash = jit.Imm(int64(caching.S_strhash))
)

var (
    _F_b64decode   = jit.Imm(int64(rt.SubrB64Decode))
    _F_decodeValue = jit.Imm(int64(_subr_decode_value))
)

var (
    _F_FieldMap_GetCaseInsensitive obj.Addr
    _Empty_Slice = []byte{}
    _Zero_Base = int64(uintptr(((*rt.GoSlice)(unsafe.Pointer(&_Empty_Slice))).Ptr))
)

const (
    _MODE_AVX2 = 1 << 2
)

const (
    _Fe_ID   = int64(unsafe.Offsetof(caching.FieldEntry{}.ID))
    _Fe_Name = int64(unsafe.Offsetof(caching.FieldEntry{}.Name))
    _Fe_Hash = int64(unsafe.Offsetof(caching.FieldEntry{}.Hash))
)

const (
    _Vk_Ptr       = int64(reflect.Ptr)
    _Gt_KindFlags = int64(unsafe.Offsetof(rt.GoType{}.KindFlags))
)

func init() {
    _F_FieldMap_GetCaseInsensitive = jit.Func((*caching.FieldMap).GetCaseInsensitive)
}

func (self *_Assembler) _asm_OP_any(_ *_Instr) {
    self.Emit("MOVQ"   , jit.Ptr(_VP, 8), _CX)              
    self.Emit("TESTQ"  , _CX, _CX)                          
    self.Sjmp("JZ"     , "_decode_{n}")                     
    self.Emit("CMPQ"   , _CX, _VP)                          
    self.Sjmp("JE"     , "_decode_{n}")                     
    self.Emit("MOVQ"   , jit.Ptr(_VP, 0), _AX)              
    self.Emit("MOVBLZX", jit.Ptr(_AX, _Gt_KindFlags), _DX)  
    self.Emit("ANDL"   , jit.Imm(rt.F_kind_mask), _DX)      
    self.Emit("CMPL"   , _DX, jit.Imm(_Vk_Ptr))             
    self.Sjmp("JNE"    , "_decode_{n}")                     
    self.Emit("LEAQ"   , jit.Ptr(_VP, 8), _DI)              
    self.decode_dynamic(_AX, _DI)                           
    self.Sjmp("JMP"    , "_decode_end_{n}")                 
    self.Link("_decode_{n}")                                
    self.Emit("MOVQ"   , _ARG_fv, _DF)                      
    self.Emit("MOVQ"   , _ST, jit.Ptr(_SP, 0))              
    self.call(_F_decodeValue)                               
    self.Emit("MOVQ"   , jit.Imm(0), jit.Ptr(_SP, 0))              
    self.Emit("TESTQ"  , _EP, _EP)                          
    self.Sjmp("JNZ"    , _LB_parsing_error)                 
    self.Link("_decode_end_{n}")                            
}

func (self *_Assembler) _asm_OP_dyn(p *_Instr) {
    self.Emit("MOVQ"   , jit.Type(p.vt()), _ET)             
    self.Emit("CMPQ"   , jit.Ptr(_VP, 8), jit.Imm(0))       
    self.Sjmp("JNE"     , "_decode_dyn_non_nil_{n}")                    

    
    self.Emit("MOVQ", _IC, _VAR_ic)
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)       
    self.Sref("_decode_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)

    self.Link("_decode_dyn_non_nil_{n}")                    
    self.Emit("MOVQ"   , jit.Ptr(_VP, 0), _CX)              
    self.Emit("MOVQ"   , jit.Ptr(_CX, 8), _CX)              
    self.Emit("MOVBLZX", jit.Ptr(_CX, _Gt_KindFlags), _DX)  
    self.Emit("ANDL"   , jit.Imm(rt.F_kind_mask), _DX)      
    self.Emit("CMPL"   , _DX, jit.Imm(_Vk_Ptr))             
    self.Sjmp("JE"    , "_decode_dyn_ptr_{n}")              

    self.Emit("MOVQ", _IC, _VAR_ic)
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)       
    self.Sref("_decode_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)

    self.Link("_decode_dyn_ptr_{n}")                        
    self.Emit("LEAQ"   , jit.Ptr(_VP, 8), _DI)              
    self.decode_dynamic(_CX, _DI)                           
    self.Link("_decode_end_{n}")                            
}

func (self *_Assembler) _asm_OP_unsupported(p *_Instr) {
    self.Emit("MOVQ", jit.Type(p.vt()), _ET)               
    self.Sjmp("JMP" , _LB_type_error)                      
}

func (self *_Assembler) _asm_OP_str(_ *_Instr) {
    self.parse_string()                                     
    self.unquote_once(jit.Ptr(_VP, 0), jit.Ptr(_VP, 8), false, true)     
}

func (self *_Assembler) _asm_OP_bin(_ *_Instr) {
    self.parse_string()                                 
    self.slice_from(_VAR_st_Iv, -1)                     
    self.Emit("MOVQ" , _DI, jit.Ptr(_VP, 0))            
    self.Emit("MOVQ" , _SI, jit.Ptr(_VP, 8))            
    self.Emit("SHRQ" , jit.Imm(2), _SI)                 
    self.Emit("LEAQ" , jit.Sib(_SI, _SI, 2, 0), _SI)    
    self.Emit("MOVQ" , _SI, jit.Ptr(_VP, 16))           
    self.malloc_AX(_SI, _SI)                               

    
    self.Emit("MOVL", jit.Imm(_MODE_JSON), _CX)          

    
    self.Emit("XORL" , _DX, _DX)                
    self.Emit("MOVQ" , _VP, _DI)                

    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _R8)    
    self.WriteRecNotAX(4, _SI, jit.Ptr(_VP, 0), true, false)    
    self.Emit("MOVQ" , _R8, _SI)

    self.Emit("XCHGQ", _DX, jit.Ptr(_VP, 8))    
    self.call_c(_F_b64decode)                     
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_base64_error)        
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 8))    
}

func (self *_Assembler) _asm_OP_bool(_ *_Instr) {
    self.Emit("LEAQ", jit.Ptr(_IC, 4), _AX)                     
    self.Emit("CMPQ", _AX, _IL)                                 
    self.Sjmp("JA"  , _LB_eof_error)                            
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('f'))    
    self.Sjmp("JE"  , "_false_{n}")                             
    self.Emit("MOVL", jit.Imm(_IM_true), _CX)                   
    self.Emit("CMPL", _CX, jit.Sib(_IP, _IC, 1, 0))             
    self.Sjmp("JE" , "_bool_true_{n}")          
    
    self.Emit("MOVQ", _IC, _VAR_ic)           
    self.Emit("MOVQ", _T_bool, _ET)     
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)         
    self.Sref("_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one) 

    self.Link("_bool_true_{n}")
    self.Emit("MOVQ", _AX, _IC)                                 
    self.Emit("MOVB", jit.Imm(1), jit.Ptr(_VP, 0))              
    self.Sjmp("JMP" , "_end_{n}")                               
    self.Link("_false_{n}")                                     
    self.Emit("ADDQ", jit.Imm(1), _AX)                          
    self.Emit("ADDQ", jit.Imm(1), _IC)                          
    self.Emit("CMPQ", _AX, _IL)                                 
    self.Sjmp("JA"  , _LB_eof_error)                            
    self.Emit("MOVL", jit.Imm(_IM_alse), _CX)                   
    self.Emit("CMPL", _CX, jit.Sib(_IP, _IC, 1, 0))             
    self.Sjmp("JNE" , _LB_im_error)                             
    self.Emit("MOVQ", _AX, _IC)                                 
    self.Emit("XORL", _AX, _AX)                                 
    self.Emit("MOVB", _AX, jit.Ptr(_VP, 0))                     
    self.Link("_end_{n}")                                       
}

func (self *_Assembler) _asm_OP_num(_ *_Instr) {
    self.Emit("MOVQ", jit.Imm(0), _VAR_fl)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('"'))
    self.Emit("MOVQ", _IC, _BX)
    self.Sjmp("JNE", "_skip_number_{n}")
    self.Emit("MOVQ", jit.Imm(1), _VAR_fl)
    self.Emit("ADDQ", jit.Imm(1), _IC)
    self.Link("_skip_number_{n}")

    
    self.Emit("LEAQ", _ARG_s, _DI)                      
    self.Emit("MOVQ", _IC, _ARG_ic)                     
    self.Emit("LEAQ", _ARG_ic, _SI)                     
    self.callc(_F_skip_number)                          
    self.Emit("MOVQ", _ARG_ic, _IC)                     
    self.Emit("TESTQ", _AX, _AX)                        
    self.Sjmp("JNS"   , "_num_next_{n}")

    
    self.Emit("MOVQ", _BX, _VAR_ic)           
    self.Emit("MOVQ", _T_number, _ET)     
    self.Emit("MOVQ", _ET, _VAR_et)
    self.Byte(0x4c, 0x8d, 0x0d)       
    self.Sref("_num_end_{n}", 4)
    self.Emit("MOVQ", _R9, _VAR_pc)
    self.Sjmp("JMP"  , _LB_skip_one)

    
    self.Link("_num_next_{n}")
    self.slice_from_r(_AX, 0)
    self.Emit("BTQ", jit.Imm(_F_copy_string), _ARG_fv)
    self.Sjmp("JNC", "_num_write_{n}")
    self.Byte(0x4c, 0x8d, 0x0d)         
    self.Sref("_num_write_{n}", 4)
    self.Sjmp("JMP", "_copy_string")
    self.Link("_num_write_{n}")
    self.Emit("MOVQ", _SI, jit.Ptr(_VP, 8))     
    self.WriteRecNotAX(13, _DI, jit.Ptr(_VP, 0), false, false)
    self.Emit("CMPQ", _VAR_fl, jit.Imm(1))
    self.Sjmp("JNE", "_num_end_{n}")
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('"'))
    self.Sjmp("JNE", _LB_char_0_error)
    self.Emit("ADDQ", jit.Imm(1), _IC)
    self.Link("_num_end_{n}")
}

func (self *_Assembler) _asm_OP_i8(_ *_Instr) {
    var pin = "_i8_end_{n}"
    self.parse_signed(int8Type, pin, -1)                                                 
    self.range_signed_CX(_I_int8, _T_int8, math.MinInt8, math.MaxInt8)     
    self.Emit("MOVB", _CX, jit.Ptr(_VP, 0))                             
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_i16(_ *_Instr) {
    var pin = "_i16_end_{n}"
    self.parse_signed(int16Type, pin, -1)                                                     
    self.range_signed_CX(_I_int16, _T_int16, math.MinInt16, math.MaxInt16)     
    self.Emit("MOVW", _CX, jit.Ptr(_VP, 0))                                 
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_i32(_ *_Instr) {
    var pin = "_i32_end_{n}"
    self.parse_signed(int32Type, pin, -1)                                                     
    self.range_signed_CX(_I_int32, _T_int32, math.MinInt32, math.MaxInt32)     
    self.Emit("MOVL", _CX, jit.Ptr(_VP, 0))                                 
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_i64(_ *_Instr) {
    var pin = "_i64_end_{n}"
    self.parse_signed(int64Type, pin, -1)                         
    self.Emit("MOVQ", _VAR_st_Iv, _AX)          
    self.Emit("MOVQ", _AX, jit.Ptr(_VP, 0))     
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u8(_ *_Instr) {
    var pin = "_u8_end_{n}"
    self.parse_unsigned(uint8Type, pin, -1)                                   
    self.range_unsigned_CX(_I_uint8, _T_uint8, math.MaxUint8)  
    self.Emit("MOVB", _CX, jit.Ptr(_VP, 0))                 
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u16(_ *_Instr) {
    var pin = "_u16_end_{n}"
    self.parse_unsigned(uint16Type, pin, -1)                                       
    self.range_unsigned_CX(_I_uint16, _T_uint16, math.MaxUint16)   
    self.Emit("MOVW", _CX, jit.Ptr(_VP, 0))                     
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u32(_ *_Instr) {
    var pin = "_u32_end_{n}"
    self.parse_unsigned(uint32Type, pin, -1)                                       
    self.range_unsigned_CX(_I_uint32, _T_uint32, math.MaxUint32)   
    self.Emit("MOVL", _CX, jit.Ptr(_VP, 0))                     
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_u64(_ *_Instr) {
    var pin = "_u64_end_{n}"
    self.parse_unsigned(uint64Type, pin, -1)                       
    self.Emit("MOVQ", _VAR_st_Iv, _AX)          
    self.Emit("MOVQ", _AX, jit.Ptr(_VP, 0))     
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_f32(_ *_Instr) {
    var pin = "_f32_end_{n}"
    self.parse_number(float32Type, pin, -1)                         
    self.range_single_X0()                         
    self.Emit("MOVSS", _X0, jit.Ptr(_VP, 0))    
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_f64(_ *_Instr) {
    var pin = "_f64_end_{n}"
    self.parse_number(float64Type, pin, -1)                         
    self.Emit("MOVSD", _VAR_st_Dv, _X0)         
    self.Emit("MOVSD", _X0, jit.Ptr(_VP, 0))    
    self.Link(pin)
}

func (self *_Assembler) _asm_OP_unquote(_ *_Instr) {
    self.check_eof(2)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm('\\'))   
    self.Sjmp("JNE" , _LB_char_0_error)                         
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 1), jit.Imm('"'))    
    self.Sjmp("JNE" , _LB_char_1_error)                         
    self.Emit("ADDQ", jit.Imm(2), _IC)                          
    self.parse_string()                                         
    self.unquote_twice(jit.Ptr(_VP, 0), jit.Ptr(_VP, 8), false)        
}

func (self *_Assembler) _asm_OP_nil_1(_ *_Instr) {
    self.Emit("XORL", _AX, _AX)                 
    self.Emit("MOVQ", _AX, jit.Ptr(_VP, 0))     
}

func (self *_Assembler) _asm_OP_nil_2(_ *_Instr) {
    self.Emit("PXOR" , _X0, _X0)                
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 0))    
}

func (self *_Assembler) _asm_OP_nil_3(_ *_Instr) {
    self.Emit("XORL" , _AX, _AX)                
    self.Emit("PXOR" , _X0, _X0)                
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 0))    
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 16))   
}

var (
    bytes []byte = make([]byte, 0)
    zerobytes = (*rt.GoSlice)(unsafe.Pointer(&bytes)).Ptr
    _ZERO_PTR = jit.Imm(int64(uintptr(zerobytes)))
)

func (self *_Assembler) _asm_OP_empty_bytes(_ *_Instr) {
    self.Emit("MOVQ", _ZERO_PTR, _AX)
    self.Emit("PXOR" , _X0, _X0)
    self.Emit("MOVQ", _AX,  jit.Ptr(_VP, 0))
    self.Emit("MOVOU", _X0, jit.Ptr(_VP, 8))
}

func (self *_Assembler) _asm_OP_deref(p *_Instr) {
    self.vfollow(p.vt())
}

func (self *_Assembler) _asm_OP_index(p *_Instr) {
    self.Emit("MOVQ", jit.Imm(p.i64()), _AX)    
    self.Emit("ADDQ", _AX, _VP)                 
}

func (self *_Assembler) _asm_OP_is_null(p *_Instr) {
    self.Emit("LEAQ"   , jit.Ptr(_IC, 4), _AX)                          
    self.Emit("CMPQ"   , _AX, _IL)                                      
    self.Sjmp("JA"     , "_not_null_{n}")                               
    self.Emit("CMPL"   , jit.Sib(_IP, _IC, 1, 0), jit.Imm(_IM_null))    
    self.Emit("CMOVQEQ", _AX, _IC)                                      
    self.Xjmp("JE"     , p.vi())                                        
    self.Link("_not_null_{n}")                                          
}

func (self *_Assembler) _asm_OP_is_null_quote(p *_Instr) {
    self.Emit("LEAQ"   , jit.Ptr(_IC, 5), _AX)                          
    self.Emit("CMPQ"   , _AX, _IL)                                      
    self.Sjmp("JA"     , "_not_null_quote_{n}")                         
    self.Emit("CMPL"   , jit.Sib(_IP, _IC, 1, 0), jit.Imm(_IM_null))    
    self.Sjmp("JNE"    , "_not_null_quote_{n}")                         
    self.Emit("CMPB"   , jit.Sib(_IP, _IC, 1, 4), jit.Imm('"'))         
    self.Emit("CMOVQEQ", _AX, _IC)                                      
    self.Xjmp("JE"     , p.vi())                                        
    self.Link("_not_null_quote_{n}")                                    
}

func (self *_Assembler) _asm_OP_map_init(_ *_Instr) {
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _AX)    
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JNZ"  , "_end_{n}")              
    self.call_go(_F_makemap_small)              
    self.WritePtrAX(6, jit.Ptr(_VP, 0), false)    
    self.Link("_end_{n}")                       
    self.Emit("MOVQ" , _AX, _VP)                
}

func (self *_Assembler) _asm_OP_map_key_i8(p *_Instr) {
    self.parse_signed(int8Type, "", p.vi())                                                 
    self.range_signed_CX(_I_int8, _T_int8, math.MinInt8, math.MaxInt8)     
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                              
}

func (self *_Assembler) _asm_OP_map_key_i16(p *_Instr) {
    self.parse_signed(int16Type, "", p.vi())                                                     
    self.range_signed_CX(_I_int16, _T_int16, math.MinInt16, math.MaxInt16)     
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                                  
}

func (self *_Assembler) _asm_OP_map_key_i32(p *_Instr) {
    self.parse_signed(int32Type, "", p.vi())                                                     
    self.range_signed_CX(_I_int32, _T_int32, math.MinInt32, math.MaxInt32)     
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)                                  
    } else {
        self.Emit("MOVQ", _CX, _AX)                                         
        self.mapassign_fastx(vt, _F_mapassign_fast32)                       
    }
}

func (self *_Assembler) _asm_OP_map_key_i64(p *_Instr) {
    self.parse_signed(int64Type, "", p.vi())                                 
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)              
    } else {
        self.Emit("MOVQ", _VAR_st_Iv, _AX)              
        self.mapassign_fastx(vt, _F_mapassign_fast64)   
    }
}

func (self *_Assembler) _asm_OP_map_key_u8(p *_Instr) {
    self.parse_unsigned(uint8Type, "", p.vi())                                   
    self.range_unsigned_CX(_I_uint8, _T_uint8, math.MaxUint8)  
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                    
}

func (self *_Assembler) _asm_OP_map_key_u16(p *_Instr) {
    self.parse_unsigned(uint16Type, "", p.vi())                                       
    self.range_unsigned_CX(_I_uint16, _T_uint16, math.MaxUint16)   
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Iv)                      
}

func (self *_Assembler) _asm_OP_map_key_u32(p *_Instr) {
    self.parse_unsigned(uint32Type, "", p.vi())                                       
    self.range_unsigned_CX(_I_uint32, _T_uint32, math.MaxUint32)   
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)                      
    } else {
        self.Emit("MOVQ", _CX, _AX)                             
        self.mapassign_fastx(vt, _F_mapassign_fast32)           
    }
}

func (self *_Assembler) _asm_OP_map_key_u64(p *_Instr) {
    self.parse_unsigned(uint64Type, "", p.vi())                                       
    self.match_char('"')
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.mapassign_std(vt, _VAR_st_Iv)                      
    } else {
        self.Emit("MOVQ", _VAR_st_Iv, _AX)                      
        self.mapassign_fastx(vt, _F_mapassign_fast64)           
    }
}

func (self *_Assembler) _asm_OP_map_key_f32(p *_Instr) {
    self.parse_number(float32Type, "", p.vi())                     
    self.range_single_X0()                     
    self.Emit("MOVSS", _X0, _VAR_st_Dv)     
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Dv)  
}

func (self *_Assembler) _asm_OP_map_key_f64(p *_Instr) {
    self.parse_number(float64Type, "", p.vi())                     
    self.match_char('"')
    self.mapassign_std(p.vt(), _VAR_st_Dv)  
}

func (self *_Assembler) _asm_OP_map_key_str(p *_Instr) {
    self.parse_string()                          
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)      
    if vt := p.vt(); !rt.IsMapfast(vt) {
        self.valloc(vt.Key(), _DI)
        self.Emit("MOVOU", _ARG_sv, _X0)
        self.Emit("MOVOU", _X0, jit.Ptr(_DI, 0))
        self.mapassign_std(vt, jit.Ptr(_DI, 0))        
    } else {
        self.mapassign_str_fast(vt, _ARG_sv_p, _ARG_sv_n)    
    }
}

func (self *_Assembler) _asm_OP_map_key_utext(p *_Instr) {
    self.parse_string()                         
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)     
    self.mapassign_utext(p.vt(), false)         
}

func (self *_Assembler) _asm_OP_map_key_utext_p(p *_Instr) {
    self.parse_string()                         
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, true)     
    self.mapassign_utext(p.vt(), true)          
}

func (self *_Assembler) _asm_OP_array_skip(_ *_Instr) {
    self.call_sf(_F_skip_array)                 
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_parsing_error_v)     
}

func (self *_Assembler) _asm_OP_array_clear(p *_Instr) {
    self.mem_clear_rem(p.i64(), true)
}

func (self *_Assembler) _asm_OP_array_clear_p(p *_Instr) {
    self.mem_clear_rem(p.i64(), false)
}

func (self *_Assembler) _asm_OP_slice_init(p *_Instr) {
    self.Emit("XORL" , _AX, _AX)                    
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 8))        
    self.Emit("MOVQ" , jit.Ptr(_VP, 16), _BX)       
    self.Emit("TESTQ", _BX, _BX)                    
    self.Sjmp("JNZ"  , "_done_{n}")                 
    self.Emit("MOVQ" , jit.Imm(_MinSlice), _CX)     
    self.Emit("MOVQ" , _CX, jit.Ptr(_VP, 16))       
    self.Emit("MOVQ" , jit.Type(p.vt()), _AX)       
    self.call_go(_F_makeslice)                      
    self.WritePtrAX(7, jit.Ptr(_VP, 0), false)      
    self.Emit("XORL" , _AX, _AX)                    
    self.Emit("MOVQ" , _AX, jit.Ptr(_VP, 8))        
    self.Link("_done_{n}")                          
}

func (self *_Assembler) _asm_OP_check_empty(p *_Instr) {
    rbracket := p.vb()
    if rbracket == ']' {
        self.check_eof(1)
        self.Emit("LEAQ", jit.Ptr(_IC, 1), _AX)                              
        self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(rbracket))) 
        self.Sjmp("JNE" , "_not_empty_array_{n}")                            
        self.Emit("MOVQ", _AX, _IC)                                          
        self.Emit("MOVQ", jit.Imm(_Zero_Base), _AX)
        self.WritePtrAX(9, jit.Ptr(_VP, 0), false)
        self.Emit("PXOR", _X0, _X0)                                          
        self.Emit("MOVOU", _X0, jit.Ptr(_VP, 8))                             
        self.Xjmp("JMP" , p.vi())                                            
        self.Link("_not_empty_array_{n}")
    } else {
        panic("only implement check empty array here!")
    }
}

func (self *_Assembler) _asm_OP_slice_append(p *_Instr) {
    self.Emit("MOVQ" , jit.Ptr(_VP, 8), _AX)            
    self.Emit("CMPQ" , _AX, jit.Ptr(_VP, 16))           
    self.Sjmp("JB"   , "_index_{n}")                    
    self.Emit("MOVQ" , _AX, _SI)                        
    self.Emit("SHLQ" , jit.Imm(1), _SI)                 
    self.Emit("MOVQ" , jit.Type(p.vt()), _AX)           
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _BX)            
    self.Emit("MOVQ" , jit.Ptr(_VP, 8), _CX)            
    self.Emit("MOVQ" , jit.Ptr(_VP, 16), _DI)           
    self.call_go(_F_growslice)                          
    self.WritePtrAX(8, jit.Ptr(_VP, 0), false)          
    self.Emit("MOVQ" , _BX, jit.Ptr(_VP, 8))            
    self.Emit("MOVQ" , _CX, jit.Ptr(_VP, 16))           

    
    
    if rt.UnpackType(p.vt()).PtrData == 0 {
        self.Emit("MOVQ" , _CX, _DI)                        
        self.Emit("SUBQ" , _BX, _DI)                        
    
        self.Emit("ADDQ" , jit.Imm(1), jit.Ptr(_VP, 8))     
        self.Emit("MOVQ" , _AX, _VP)                        
        self.Emit("MOVQ" , jit.Imm(int64(p.vlen())), _CX)   
        self.Emit("MOVQ" , _BX, _AX)                        
        self.From("MULQ" , _CX)                             
        self.Emit("ADDQ" , _AX, _VP)                        

        self.Emit("MOVQ" , _DI, _AX)                        
        self.From("MULQ" , _CX)                             
        self.Emit("MOVQ" , _AX, _BX)                        
        self.Emit("MOVQ" , _VP, _AX)                        
        self.mem_clear_fn(true)                             
        self.Sjmp("JMP", "_append_slice_end_{n}")
    }

    self.Emit("MOVQ" , _BX, _AX)                        
    self.Link("_index_{n}")                             
    self.Emit("ADDQ" , jit.Imm(1), jit.Ptr(_VP, 8))     
    self.Emit("MOVQ" , jit.Ptr(_VP, 0), _VP)            
    self.Emit("MOVQ" , jit.Imm(int64(p.vlen())), _CX)   
    self.From("MULQ" , _CX)                             
    self.Emit("ADDQ" , _AX, _VP)                        
    self.Link("_append_slice_end_{n}")
}

func (self *_Assembler) _asm_OP_object_next(_ *_Instr) {
    self.call_sf(_F_skip_one)                   
    self.Emit("TESTQ", _AX, _AX)                
    self.Sjmp("JS"   , _LB_parsing_error_v)     
}

func (self *_Assembler) _asm_OP_struct_field(p *_Instr) {
    assert_eq(caching.FieldEntrySize, 32, "invalid field entry size")
    self.Emit("MOVQ" , jit.Imm(-1), _AX)                        
    self.Emit("MOVQ" , _AX, _VAR_sr)                            
    self.parse_string()                                         
    self.unquote_once(_ARG_sv_p, _ARG_sv_n, true, false)                     
    self.Emit("LEAQ" , _ARG_sv, _AX)                            
    self.Emit("XORL" , _BX, _BX)                                
    self.call_go(_F_strhash)                                    
    self.Emit("MOVQ" , _AX, _R9)                                
    self.Emit("MOVQ" , jit.Imm(freezeFields(p.vf())), _CX)      
    self.Emit("MOVQ" , jit.Ptr(_CX, caching.FieldMap_b), _SI)   
    self.Emit("MOVQ" , jit.Ptr(_CX, caching.FieldMap_N), _CX)   
    self.Emit("TESTQ", _CX, _CX)                                
    self.Sjmp("JZ"   , "_try_lowercase_{n}")                    
    self.Link("_loop_{n}")                                      
    self.Emit("XORL" , _DX, _DX)                                
    self.From("DIVQ" , _CX)                                     
    self.Emit("LEAQ" , jit.Ptr(_DX, 1), _AX)                    
    self.Emit("SHLQ" , jit.Imm(5), _DX)                         
    self.Emit("LEAQ" , jit.Sib(_SI, _DX, 1, 0), _DI)            
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_Hash), _R8)             
    self.Emit("TESTQ", _R8, _R8)                                
    self.Sjmp("JZ"   , "_try_lowercase_{n}")                    
    self.Emit("CMPQ" , _R8, _R9)                                
    self.Sjmp("JNE"  , "_loop_{n}")                             
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_Name + 8), _DX)         
    self.Emit("CMPQ" , _DX, _ARG_sv_n)                          
    self.Sjmp("JNE"  , "_loop_{n}")                             
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_ID), _R8)               
    self.Emit("MOVQ" , _AX, _VAR_ss_AX)                         
    self.Emit("MOVQ" , _CX, _VAR_ss_CX)                         
    self.Emit("MOVQ" , _SI, _VAR_ss_SI)                         
    self.Emit("MOVQ" , _R8, _VAR_ss_R8)                         
    self.Emit("MOVQ" , _R9, _VAR_ss_R9)                         
    self.Emit("MOVQ" , _ARG_sv_p, _AX)                          
    self.Emit("MOVQ" , jit.Ptr(_DI, _Fe_Name), _CX)             
    self.Emit("MOVQ" , _CX, _BX)                                
    self.Emit("MOVQ" , _DX, _CX)                                
    self.call_go(_F_memequal)                                   
    self.Emit("MOVB" , _AX, _DX)                                
    self.Emit("MOVQ" , _VAR_ss_AX, _AX)                         
    self.Emit("MOVQ" , _VAR_ss_CX, _CX)                         
    self.Emit("MOVQ" , _VAR_ss_SI, _SI)                         
    self.Emit("MOVQ" , _VAR_ss_R9, _R9)                         
    self.Emit("TESTB", _DX, _DX)                                
    self.Sjmp("JZ"   , "_loop_{n}")                             
    self.Emit("MOVQ" , _VAR_ss_R8, _R8)                         
    self.Emit("MOVQ" , _R8, _VAR_sr)                            
    self.Sjmp("JMP"  , "_end_{n}")                              
    self.Link("_try_lowercase_{n}")                             
    self.Emit("BTQ"  , jit.Imm(_F_case_sensitive), _ARG_fv)     
    self.Sjmp("JC"   , "_unknown_{n}")                         
    self.Emit("MOVQ" , jit.Imm(referenceFields(p.vf())), _AX)   
    self.Emit("MOVQ", _ARG_sv_p, _BX)                            
    self.Emit("MOVQ", _ARG_sv_n, _CX)                            
    self.call_go(_F_FieldMap_GetCaseInsensitive)                
    self.Emit("MOVQ" , _AX, _VAR_sr)                            
    self.Emit("TESTQ", _AX, _AX)                                
    self.Sjmp("JNS"  , "_end_{n}")                              
    self.Link("_unknown_{n}")
    
    self.Emit("MOVQ" , jit.Imm(-1), _AX)                        
    self.Emit("MOVQ" , _AX, _VAR_sr)                            
    self.Emit("BTQ"  , jit.Imm(_F_disable_unknown), _ARG_fv)    
    self.Sjmp("JC"   , _LB_field_error)                         
    self.Link("_end_{n}")                                       
}

func (self *_Assembler) _asm_OP_unmarshal(p *_Instr) {
    if iv := p.i64(); iv != 0 {
        self.unmarshal_json(p.vt(), true, _F_decodeJsonUnmarshalerQuoted)
    } else {
        self.unmarshal_json(p.vt(), true, _F_decodeJsonUnmarshaler)
    }
}

func (self *_Assembler) _asm_OP_unmarshal_p(p *_Instr) {
    if iv := p.i64(); iv != 0 {
        self.unmarshal_json(p.vt(), false, _F_decodeJsonUnmarshalerQuoted)
    } else {
        self.unmarshal_json(p.vt(), false, _F_decodeJsonUnmarshaler)
    }
}

func (self *_Assembler) _asm_OP_unmarshal_text(p *_Instr) {
    self.unmarshal_text(p.vt(), true)
}

func (self *_Assembler) _asm_OP_unmarshal_text_p(p *_Instr) {
    self.unmarshal_text(p.vt(), false)
}

func (self *_Assembler) _asm_OP_lspace(_ *_Instr) {
    self.lspace("_{n}")
}

func (self *_Assembler) lspace(subfix string) {
    var label = "_lspace" + subfix
    self.Emit("CMPQ"   , _IC, _IL)                      
    self.Sjmp("JAE"    , _LB_eof_error)                 
    self.Emit("MOVQ"   , jit.Imm(_BM_space), _DX)       
    self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  
    self.Emit("CMPQ"   , _AX, jit.Imm(' '))             
    self.Sjmp("JA"     , label)                
    self.Emit("BTQ"    , _AX, _DX)                      
    self.Sjmp("JNC"    , label)                

    
    for i := 0; i < 3; i++ {
        self.Emit("ADDQ"   , jit.Imm(1), _IC)               
        self.Emit("CMPQ"   , _IC, _IL)                      
        self.Sjmp("JAE"    , _LB_eof_error)                 
        self.Emit("MOVBQZX", jit.Sib(_IP, _IC, 1, 0), _AX)  
        self.Emit("CMPQ"   , _AX, jit.Imm(' '))             
        self.Sjmp("JA"     , label)                
        self.Emit("BTQ"    , _AX, _DX)                      
        self.Sjmp("JNC"    , label)                
    }

    
    self.Emit("MOVQ"   , _IP, _DI)                      
    self.Emit("MOVQ"   , _IL, _SI)                      
    self.Emit("MOVQ"   , _IC, _DX)                      
    self.callc(_F_lspace)                                
    self.Emit("TESTQ"  , _AX, _AX)                      
    self.Sjmp("JS"     , _LB_parsing_error_v)           
    self.Emit("CMPQ"   , _AX, _IL)                      
    self.Sjmp("JAE"    , _LB_eof_error)                 
    self.Emit("MOVQ"   , _AX, _IC)                      
    self.Link(label)                           
}

func (self *_Assembler) _asm_OP_match_char(p *_Instr) {
    self.match_char(p.vb())
}

func (self *_Assembler) match_char(char byte) {
    self.check_eof(1)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(char)))  
    self.Sjmp("JNE" , _LB_char_0_error)                                 
    self.Emit("ADDQ", jit.Imm(1), _IC)                                  
}

func (self *_Assembler) _asm_OP_check_char(p *_Instr) {
    self.check_eof(1)
    self.Emit("LEAQ"   , jit.Ptr(_IC, 1), _AX)                              
    self.Emit("CMPB"   , jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(p.vb())))   
    self.Emit("CMOVQEQ", _AX, _IC)                                          
    self.Xjmp("JE"     , p.vi())                                            
}

func (self *_Assembler) _asm_OP_check_char_0(p *_Instr) {
    self.check_eof(1)
    self.Emit("CMPB", jit.Sib(_IP, _IC, 1, 0), jit.Imm(int64(p.vb())))   
    self.Xjmp("JE"  , p.vi())                                            
}

func (self *_Assembler) _asm_OP_add(p *_Instr) {
    self.Emit("ADDQ", jit.Imm(int64(p.vi())), _IC)  
}

func (self *_Assembler) _asm_OP_load(_ *_Instr) {
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             
    self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 0), _VP)     
}

func (self *_Assembler) _asm_OP_save(_ *_Instr) {
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _CX)             
    self.Emit("CMPQ", _CX, jit.Imm(_MaxStackBytes))     
    self.Sjmp("JAE"  , _LB_stack_error)                  
    self.WriteRecNotAX(0 , _VP, jit.Sib(_ST, _CX, 1, 8), false, false) 
    self.Emit("ADDQ", jit.Imm(8), _CX)                  
    self.Emit("MOVQ", _CX, jit.Ptr(_ST, 0))             
}

func (self *_Assembler) _asm_OP_drop(_ *_Instr) {
    self.Emit("MOVQ", jit.Ptr(_ST, 0), _AX)             
    self.Emit("SUBQ", jit.Imm(8), _AX)                  
    self.Emit("MOVQ", jit.Sib(_ST, _AX, 1, 8), _VP)     
    self.Emit("MOVQ", _AX, jit.Ptr(_ST, 0))             
    self.Emit("XORL", _BX, _BX)                         
    self.Emit("MOVQ", _BX, jit.Sib(_ST, _AX, 1, 8))     
}

func (self *_Assembler) _asm_OP_drop_2(_ *_Instr) {
    self.Emit("MOVQ" , jit.Ptr(_ST, 0), _AX)            
    self.Emit("SUBQ" , jit.Imm(16), _AX)                
    self.Emit("MOVQ" , jit.Sib(_ST, _AX, 1, 8), _VP)    
    self.Emit("MOVQ" , _AX, jit.Ptr(_ST, 0))            
    self.Emit("PXOR" , _X0, _X0)                        
    self.Emit("MOVOU", _X0, jit.Sib(_ST, _AX, 1, 8))    
}

func (self *_Assembler) _asm_OP_recurse(p *_Instr) {
    self.Emit("MOVQ", jit.Type(p.vt()), _AX)    
    self.decode_dynamic(_AX, _VP)               
}

func (self *_Assembler) _asm_OP_goto(p *_Instr) {
    self.Xjmp("JMP", p.vi())
}

func (self *_Assembler) _asm_OP_switch(p *_Instr) {
    self.Emit("MOVQ", _VAR_sr, _AX)             
    self.Emit("CMPQ", _AX, jit.Imm(p.i64()))    
    self.Sjmp("JAE" , "_default_{n}")           

    
    self.Byte(0x48, 0x8d, 0x3d)                         
    self.Sref("_switch_table_{n}", 4)                   
    self.Emit("MOVLQSX", jit.Sib(_DI, _AX, 4, 0), _AX)  
    self.Emit("ADDQ"   , _DI, _AX)                      
    self.Rjmp("JMP"    , _AX)                           
    self.Link("_switch_table_{n}")                      

    
    for i, v := range p.vs() {
        self.Xref(v, int64(-i) * 4)
    }

    
    self.Link("_default_{n}")
    self.NOP()
}

func (self *_Assembler) print_gc(i int, p1 *_Instr, p2 *_Instr) {
    self.Emit("MOVQ", jit.Imm(int64(p2.op())),  _CX)
    self.Emit("MOVQ", jit.Imm(int64(p1.op())),  _BX) 
    self.Emit("MOVQ", jit.Imm(int64(i)),  _AX)       
    self.call_go(_F_println)
}
