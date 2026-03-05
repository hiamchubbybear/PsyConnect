

package vars

import (
	"os"
	"unsafe"
)

const (
	MaxStack = 4096 
	StackSize = unsafe.Sizeof(Stack{})
	StateSize  = int64(unsafe.Sizeof(State{}))
	StackLimit = MaxStack * StateSize
)

const (
	MAX_ILBUF  = 100000 
	MAX_FIELDS = 50     
)

var (
	DebugSyncGC   = os.Getenv("SONIC_SYNC_GC") != ""
	DebugAsyncGC  = os.Getenv("SONIC_NO_ASYNC_GC") == ""
	DebugCheckPtr = os.Getenv("SONIC_CHECK_POINTER") != ""
)

var UseVM = os.Getenv("SONIC_ENCODER_USE_VM") != ""
