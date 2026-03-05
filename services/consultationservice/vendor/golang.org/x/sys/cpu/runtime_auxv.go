



package cpu



var getAuxvFn func() []uintptr

func getAuxv() []uintptr {
	if getAuxvFn == nil {
		return nil
	}
	return getAuxvFn()
}
