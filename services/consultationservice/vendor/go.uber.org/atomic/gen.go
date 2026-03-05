



















package atomic

//go:generate bin/gen-atomicint -name=Int32 -wrapped=int32 -file=int32.go
//go:generate bin/gen-atomicint -name=Int64 -wrapped=int64 -file=int64.go
//go:generate bin/gen-atomicint -name=Uint32 -wrapped=uint32 -unsigned -file=uint32.go
//go:generate bin/gen-atomicint -name=Uint64 -wrapped=uint64 -unsigned -file=uint64.go
//go:generate bin/gen-atomicint -name=Uintptr -wrapped=uintptr -unsigned -file=uintptr.go
