//go:build !amd64
// +build !amd64



package encoder

func init() {
	ForceUseVM()
}
