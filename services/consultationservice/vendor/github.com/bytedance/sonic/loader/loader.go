

package loader

import (
    `unsafe`
)


type Function unsafe.Pointer


type Options struct {
    
    NoPreempt bool
}


type Loader struct {
    Name string 
    File string 
    Options 
}
