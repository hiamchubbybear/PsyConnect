




package cpuinfo


func HasBMI1() bool {
	return hasBMI1
}


func HasBMI2() bool {
	return hasBMI2
}



func DisableBMI2() func() {
	old := hasBMI2
	hasBMI2 = false
	return func() {
		hasBMI2 = old
	}
}


func HasBMI() bool {
	return HasBMI1() && HasBMI2()
}

var hasBMI1 bool
var hasBMI2 bool
