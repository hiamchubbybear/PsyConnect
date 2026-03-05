



//go:build !race

package impl


type RaceDetectHookData struct{}


func (presence) raceDetectHookPresent(num uint32)                       {}
func (presence) raceDetectHookSetPresent(num uint32, size presenceSize) {}
func (presence) raceDetectHookClearPresent(num uint32)                  {}
func (presence) raceDetectHookAllocAndCopy(src presence)                {}





func raceDetectHookPresent(field *uint32, num uint32) {}





func raceDetectHookSetPresent(field *uint32, num uint32, size presenceSize) {}





func raceDetectHookClearPresent(field *uint32, num uint32) {}
