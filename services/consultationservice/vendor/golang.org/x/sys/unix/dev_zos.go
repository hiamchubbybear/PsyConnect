



//go:build zos && s390x






package unix


func Major(dev uint64) uint32 {
	return uint32((dev >> 16) & 0x0000FFFF)
}


func Minor(dev uint64) uint32 {
	return uint32(dev & 0x0000FFFF)
}



func Mkdev(major, minor uint32) uint64 {
	return (uint64(major) << 16) | uint64(minor)
}
