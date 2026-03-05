





package pragma

import "sync"


type NoUnkeyedLiterals struct{}






type DoNotImplement interface{ ProtoInternal(DoNotImplement) }


type DoNotCompare [0]func()






type DoNotCopy [0]sync.Mutex
