


package loader

import (
    `github.com/bytedance/sonic/loader/internal/rt`
)












func (self Loader) LoadOne(text []byte, funcName string, frameSize int, argSize int, argPtrs []bool, localPtrs []bool) Function {
    size := uint32(len(text))

    fn := Func{
        Name: funcName,
        TextSize: size,
        ArgsSize: int32(argSize),
    }

    
    fn.Pcsp = &Pcdata{
        {PC: size, Val: int32(frameSize)},
    }

    if self.NoPreempt {
        fn.PcUnsafePoint = &Pcdata{
            {PC: size, Val: PCDATA_UnsafePointUnsafe},
        }
    } else {
        fn.PcUnsafePoint = &Pcdata{
            {PC: size, Val: PCDATA_UnsafePointSafe},
        }
    }

    
    fn.PcStackMapIndex = &Pcdata{
        {PC: size, Val: 0},
    }

    if argPtrs != nil {
        args := rt.StackMapBuilder{}
        for _, b := range argPtrs {
            args.AddField(b)
        }
        fn.ArgsPointerMaps = args.Build()
    }
    
    if localPtrs != nil {
        locals := rt.StackMapBuilder{}
        for _, b := range localPtrs {
            locals.AddField(b)
        }
        fn.LocalsPointerMaps = locals.Build()
    }

    out := Load(text, []Func{fn}, self.Name + funcName, []string{self.File})
    return out[0]
}




func Load(text []byte, funcs []Func, modulename string, filenames []string) (out []Function) {
    ids := make([]string, len(funcs))
    for i, f := range funcs {
        ids[i] = f.Name
    }
    
    mod := makeModuledata(modulename, filenames, &funcs, text)

    
    moduledataverify1(mod)
    registerModule(mod)

    
    
    out = make([]Function, len(funcs))
    for i, s := range ids {
        for _, f := range funcs {
            if f.Name == s {
                m := uintptr(mod.text + uintptr(f.EntryOff))
                out[i] = Function(&m)
            }
        }
    }
    return 
}
