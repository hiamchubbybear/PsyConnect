

package loader

import (
	`reflect`
	`unsafe`

	`github.com/bytedance/sonic/loader/internal/abi`
	`github.com/bytedance/sonic/loader/internal/rt`
)

var _C_Redzone = []bool{false, false, false, false}


type CFunc struct {
	
	Name     string

	
	EntryOff uint32

	
	TextSize uint32

	
	MaxStack uintptr

	
	Pcsp     [][2]uint32
}


type GoC struct {
	
	CName     string

	
	
	CEntry   *uintptr

	
	
	
	
	
	
	GoFunc   interface{} 
}


func WrapGoC(text []byte, natives []CFunc, stubs []GoC, modulename string, filename string) {
	funcs := make([]Func, len(natives))
	
	
	for i, f := range natives {
		fn := Func{
			Flag: FuncFlag_ASM,
			EntryOff: f.EntryOff,
			TextSize: f.TextSize,
			Name: f.Name,
		}
		if len(f.Pcsp) != 0 {
			fn.Pcsp = (*Pcdata)(unsafe.Pointer(&natives[i].Pcsp))
		}
		
		fn.PcUnsafePoint = &Pcdata{
			{PC: f.TextSize, Val: PCDATA_UnsafePointUnsafe},
		}
		
		fn.Pcfile = &Pcdata{
			{PC: f.TextSize, Val: 0},
		}
		
		fn.Pcline = &Pcdata{
			{PC: f.TextSize, Val: 1},
		}
		
		fn.PcStackMapIndex = &Pcdata{
			{PC: f.TextSize, Val: 0},
		}
		sm := rt.StackMapBuilder{}
		sm.AddField(false)
		fn.ArgsPointerMaps = sm.Build()
		fn.LocalsPointerMaps = sm.Build()
		funcs[i] = fn
	}
	rets := Load(text, funcs, modulename, []string{filename})

	
	native_entry := **(**uintptr)(unsafe.Pointer(&rets[0]))
	

	wraps := make([]Func, 0, len(stubs))
	wrapIds := make([]int, 0, len(stubs))
	code := make([]byte, 0, len(wraps))
	entryOff := uint32(0)

	
	for i := range stubs {
		for j := range natives {
			if stubs[i].CName != natives[j].Name {
				continue
			}
			
			
			pc := uintptr(native_entry + uintptr(natives[j].EntryOff))
			if stubs[i].CEntry != nil {
				*stubs[i].CEntry = pc
			}

			
			if stubs[i].GoFunc == nil {
				continue
			}

			
			layout := abi.NewFunctionLayout(reflect.TypeOf(stubs[i].GoFunc).Elem())
			frame := abi.NewFrame(&layout, _C_Redzone, true) 
			tcode := abi.CallC(pc, frame, natives[j].MaxStack)
			code = append(code, tcode...)
			size := uint32(len(tcode))
		
			fn := Func{
				Flag: FuncFlag_ASM,
				ArgsSize: int32(layout.ArgSize()),
				EntryOff: entryOff,
				TextSize: size,
				Name: stubs[i].CName + "_go",
			}

			
			fn.Pcsp = &Pcdata{
				{PC: uint32(frame.StackCheckTextSize()), Val: 0},
				{PC: size - uint32(frame.GrowStackTextSize()), Val: int32(frame.Size())},
				{PC: size, Val: 0},
			}
			
			fn.Pcfile = &Pcdata{
				{PC: size, Val: 0},
			}
			
			fn.Pcline = &Pcdata{
				{PC: size, Val: 1},
			}
			
			fn.PcUnsafePoint = &Pcdata{
				{PC: size, Val: PCDATA_UnsafePointUnsafe},
			}

			
			fn.PcStackMapIndex = &Pcdata{
				{PC: size, Val: 0},
			}
			fn.ArgsPointerMaps = frame.ArgPtrs()
			fn.LocalsPointerMaps = frame.LocalPtrs()

			entryOff += size
			wraps = append(wraps, fn)
			wrapIds = append(wrapIds, i)
		}
	}
	gofuncs := Load(code, wraps, modulename+"/go", []string{filename+".go"})

	
	for i := range gofuncs {
		idx := wrapIds[i]
		w := rt.UnpackEface(stubs[idx].GoFunc)
		*(*Function)(w.Value) = gofuncs[i]
	}
}
