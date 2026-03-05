



//go:build gccgo

package cpu



func haveAsmFunctions() bool { return false }




func stfle() facilityList     { panic("not implemented for gccgo") }
func kmQuery() queryResult    { panic("not implemented for gccgo") }
func kmcQuery() queryResult   { panic("not implemented for gccgo") }
func kmctrQuery() queryResult { panic("not implemented for gccgo") }
func kmaQuery() queryResult   { panic("not implemented for gccgo") }
func kimdQuery() queryResult  { panic("not implemented for gccgo") }
func klmdQuery() queryResult  { panic("not implemented for gccgo") }
