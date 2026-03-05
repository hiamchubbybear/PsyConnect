



package sys






func RaceDetectorSupported(goos, goarch string) bool {
	switch goos {
	case "linux":
		return goarch == "amd64" || goarch == "ppc64le" || goarch == "arm64"
	case "darwin", "freebsd", "netbsd", "windows":
		return goarch == "amd64"
	default:
		return false
	}
}



func MSanSupported(goos, goarch string) bool {
	switch goos {
	case "linux":
		return goarch == "amd64" || goarch == "arm64"
	default:
		return false
	}
}


func MustLinkExternal(goos, goarch string) bool {
	switch goos {
	case "android":
		if goarch != "arm64" {
			return true
		}
	case "darwin":
		if goarch == "arm64" {
			return true
		}
	}
	return false
}



func BuildModeSupported(compiler, buildmode, goos, goarch string) bool {
	if compiler == "gccgo" {
		return true
	}

	platform := goos + "/" + goarch

	switch buildmode {
	case "archive":
		return true

	case "c-archive":
		
		
		return platform != "linux/ppc64"

	case "c-shared":
		switch platform {
		case "linux/amd64", "linux/arm", "linux/arm64", "linux/386", "linux/ppc64le", "linux/s390x",
			"android/amd64", "android/arm", "android/arm64", "android/386",
			"freebsd/amd64",
			"darwin/amd64",
			"windows/amd64", "windows/386":
			return true
		}
		return false

	case "default":
		return true

	case "exe":
		return true

	case "pie":
		switch platform {
		case "linux/386", "linux/amd64", "linux/arm", "linux/arm64", "linux/ppc64le", "linux/s390x",
			"android/amd64", "android/arm", "android/arm64", "android/386",
			"freebsd/amd64",
			"darwin/amd64",
			"aix/ppc64",
			"windows/386", "windows/amd64", "windows/arm":
			return true
		}
		return false

	case "shared":
		switch platform {
		case "linux/386", "linux/amd64", "linux/arm", "linux/arm64", "linux/ppc64le", "linux/s390x":
			return true
		}
		return false

	case "plugin":
		switch platform {
		case "linux/amd64", "linux/arm", "linux/arm64", "linux/386", "linux/s390x", "linux/ppc64le",
			"android/amd64", "android/arm", "android/arm64", "android/386",
			"darwin/amd64",
			"freebsd/amd64":
			return true
		}
		return false

	default:
		return false
	}
}
