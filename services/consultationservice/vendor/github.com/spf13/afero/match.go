













package afero

import (
	"path/filepath"
	"sort"
	"strings"
)












func Glob(fs Fs, pattern string) (matches []string, err error) {
	if !hasMeta(pattern) {
		
		if _, err = lstatIfPossible(fs, pattern); err != nil {
			return nil, nil
		}
		return []string{pattern}, nil
	}

	dir, file := filepath.Split(pattern)
	switch dir {
	case "":
		dir = "."
	case string(filepath.Separator):
	
	default:
		dir = dir[0 : len(dir)-1] 
	}

	if !hasMeta(dir) {
		return glob(fs, dir, file, nil)
	}

	var m []string
	m, err = Glob(fs, dir)
	if err != nil {
		return
	}
	for _, d := range m {
		matches, err = glob(fs, d, file, matches)
		if err != nil {
			return
		}
	}
	return
}





func glob(fs Fs, dir, pattern string, matches []string) (m []string, e error) {
	m = matches
	fi, err := fs.Stat(dir)
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return
	}
	d, err := fs.Open(dir)
	if err != nil {
		return
	}
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	sort.Strings(names)

	for _, n := range names {
		matched, err := filepath.Match(pattern, n)
		if err != nil {
			return m, err
		}
		if matched {
			m = append(m, filepath.Join(dir, n))
		}
	}
	return
}



func hasMeta(path string) bool {
	
	return strings.ContainsAny(path, "*?[")
}
