



package objabi



const (
	STACKSYSTEM = 0
	StackSystem = STACKSYSTEM
	StackBig    = 4096
	StackSmall  = 128
)

const (
	StackPreempt = -1314 
)


var StackGuard = 928*stackGuardMultiplier() + StackSystem
var StackLimit = StackGuard - StackSystem - StackSmall




func stackGuardMultiplier() int {
	
	if GOOS == "aix" {
		return 2
	}
	return 1
}
