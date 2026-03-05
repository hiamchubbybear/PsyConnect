



//go:build race

package impl





type RaceDetectHookData struct {
	shadowPresence *[]byte
}



func (data *RaceDetectHookData) raceDetectHookAlloc(size presenceSize) {
	sp := make([]byte, size)
	atomicStoreShadowPresence(&data.shadowPresence, &sp)
}

func (p presence) raceDetectHookPresent(num uint32) {
	data := p.toRaceDetectData()
	if data == nil {
		return
	}
	sp := atomicLoadShadowPresence(&data.shadowPresence)
	if sp != nil {
		_ = (*sp)[num]
	}
}

func (p presence) raceDetectHookSetPresent(num uint32, size presenceSize) {
	data := p.toRaceDetectData()
	if data == nil {
		return
	}
	sp := atomicLoadShadowPresence(&data.shadowPresence)
	if sp == nil {
		data.raceDetectHookAlloc(size)
		sp = atomicLoadShadowPresence(&data.shadowPresence)
	}
	(*sp)[num] = 1
}

func (p presence) raceDetectHookClearPresent(num uint32) {
	data := p.toRaceDetectData()
	if data == nil {
		return
	}
	sp := atomicLoadShadowPresence(&data.shadowPresence)
	if sp != nil {
		(*sp)[num] = 0

	}
}



func (p presence) raceDetectHookAllocAndCopy(q presence) {
	sData := q.toRaceDetectData()
	dData := p.toRaceDetectData()
	if sData == nil {
		return
	}
	srcSp := atomicLoadShadowPresence(&sData.shadowPresence)
	if srcSp == nil {
		atomicStoreShadowPresence(&dData.shadowPresence, nil)
		return
	}
	n := len(*srcSp)
	dSlice := make([]byte, n)
	atomicStoreShadowPresence(&dData.shadowPresence, &dSlice)
	for i := 0; i < n; i++ {
		dSlice[i] = (*srcSp)[i]
	}
}





func raceDetectHookPresent(field *uint32, num uint32) {
	data := findPointerToRaceDetectData(field, num)
	if data == nil {
		return
	}
	sp := atomicLoadShadowPresence(&data.shadowPresence)
	if sp != nil {
		_ = (*sp)[num]
	}
}





func raceDetectHookSetPresent(field *uint32, num uint32, size presenceSize) {
	data := findPointerToRaceDetectData(field, num)
	if data == nil {
		return
	}
	sp := atomicLoadShadowPresence(&data.shadowPresence)
	if sp == nil {
		data.raceDetectHookAlloc(size)
		sp = atomicLoadShadowPresence(&data.shadowPresence)
	}
	(*sp)[num] = 1
}





func raceDetectHookClearPresent(field *uint32, num uint32) {
	data := findPointerToRaceDetectData(field, num)
	if data == nil {
		return
	}
	sp := atomicLoadShadowPresence(&data.shadowPresence)
	if sp != nil {
		(*sp)[num] = 0
	}
}
