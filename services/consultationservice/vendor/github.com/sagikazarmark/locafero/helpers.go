package locafero

import "fmt"




func NameWithExtensions(baseName string, extensions ...string) []string {
	var names []string

	if baseName == "" {
		return names
	}

	for _, ext := range extensions {
		if ext == "" {
			continue
		}

		names = append(names, fmt.Sprintf("%s.%s", baseName, ext))
	}

	return names
}





func NameWithOptionalExtensions(baseName string, extensions ...string) []string {
	var names []string

	if baseName == "" {
		return names
	}

	names = NameWithExtensions(baseName, extensions...)
	names = append(names, baseName)

	return names
}
