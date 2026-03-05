



package objabi






type FuncID uint8

const (
	FuncID_normal FuncID = iota 
	FuncID_runtime_main
	FuncID_goexit
	FuncID_jmpdefer
	FuncID_mcall
	FuncID_morestack
	FuncID_mstart
	FuncID_rt0_go
	FuncID_asmcgocall
	FuncID_sigpanic
	FuncID_runfinq
	FuncID_gcBgMarkWorker
	FuncID_systemstack_switch
	FuncID_systemstack
	FuncID_cgocallback_gofunc
	FuncID_gogo
	FuncID_externalthreadhandler
	FuncID_debugCallV1
	FuncID_gopanic
	FuncID_panicwrap
	FuncID_handleAsyncEvent
	FuncID_asyncPreempt
	FuncID_wrapper 
)



func GetFuncID(name string, isWrapper bool) FuncID {
	if isWrapper {
		return FuncID_wrapper
	}
	switch name {
	case "runtime.main":
		return FuncID_runtime_main
	case "runtime.goexit":
		return FuncID_goexit
	case "runtime.jmpdefer":
		return FuncID_jmpdefer
	case "runtime.mcall":
		return FuncID_mcall
	case "runtime.morestack":
		return FuncID_morestack
	case "runtime.mstart":
		return FuncID_mstart
	case "runtime.rt0_go":
		return FuncID_rt0_go
	case "runtime.asmcgocall":
		return FuncID_asmcgocall
	case "runtime.sigpanic":
		return FuncID_sigpanic
	case "runtime.runfinq":
		return FuncID_runfinq
	case "runtime.gcBgMarkWorker":
		return FuncID_gcBgMarkWorker
	case "runtime.systemstack_switch":
		return FuncID_systemstack_switch
	case "runtime.systemstack":
		return FuncID_systemstack
	case "runtime.cgocallback_gofunc":
		return FuncID_cgocallback_gofunc
	case "runtime.gogo":
		return FuncID_gogo
	case "runtime.externalthreadhandler":
		return FuncID_externalthreadhandler
	case "runtime.debugCallV1":
		return FuncID_debugCallV1
	case "runtime.gopanic":
		return FuncID_gopanic
	case "runtime.panicwrap":
		return FuncID_panicwrap
	case "runtime.handleAsyncEvent":
		return FuncID_handleAsyncEvent
	case "runtime.asyncPreempt":
		return FuncID_asyncPreempt
	case "runtime.deferreturn":
		
		return FuncID_wrapper
	case "runtime.runOpenDeferFrame":
		
		return FuncID_wrapper
	case "runtime.reflectcallSave":
		
		return FuncID_wrapper
	}
	return FuncID_normal
}
