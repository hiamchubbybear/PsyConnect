

package option

var (
    
    DefaultDecoderBufferSize  uint = 4 * 1024

    
    DefaultEncoderBufferSize  uint = 4 * 1024

    
    DefaultAstBufferSize  uint = 4 * 1024

    
    
    LimitBufferSize uint = 1024 * 1024
)


type CompileOptions struct {
    
    MaxInlineDepth int

    
    RecursiveDepth int
}

var (
    
    
    
    DefaultMaxInlineDepth = 3

    
    
    DefaultRecursiveDepth = 1
)


func DefaultCompileOptions() CompileOptions {
    return CompileOptions{
        RecursiveDepth: DefaultRecursiveDepth,
        MaxInlineDepth: DefaultMaxInlineDepth,
    }
}


type CompileOption func(o *CompileOptions)








func WithCompileRecursiveDepth(loop int) CompileOption {
    return func(o *CompileOptions) {
            if loop < 0 {
                panic("loop must be >= 0")
            }
            o.RecursiveDepth = loop
        }
}





func WithCompileMaxInlineDepth(depth int) CompileOption {
    return func(o *CompileOptions) {
            if depth <= 0 {
                panic("depth must be > 0")
            }
            o.MaxInlineDepth = depth
        }
}
