



package cpu

import "strconv"






func parseRelease(rel string) (major, minor, patch int, ok bool) {
	
	for i := 0; i < len(rel); i++ {
		if rel[i] == '-' || rel[i] == '+' {
			rel = rel[:i]
			break
		}
	}

	next := func() (int, bool) {
		for i := 0; i < len(rel); i++ {
			if rel[i] == '.' {
				ver, err := strconv.Atoi(rel[:i])
				rel = rel[i+1:]
				return ver, err == nil
			}
		}
		ver, err := strconv.Atoi(rel)
		rel = ""
		return ver, err == nil
	}
	if major, ok = next(); !ok || rel == "" {
		return
	}
	if minor, ok = next(); !ok || rel == "" {
		return
	}
	patch, ok = next()
	return
}
