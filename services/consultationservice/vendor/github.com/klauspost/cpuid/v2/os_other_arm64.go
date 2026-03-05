

//go:build arm64 && !linux && !darwin
// +build arm64,!linux,!darwin

package cpuid

import "runtime"

func detectOS(c *CPUInfo) bool {
	c.PhysicalCores = runtime.NumCPU()
	
	c.ThreadsPerCore = 1
	c.LogicalCores = c.PhysicalCores
	return false
}
