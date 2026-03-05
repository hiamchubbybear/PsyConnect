



package objabi

import (
	"os"
	"path/filepath"
	"strings"
)




func WorkingDir() string {
	var path string
	path, _ = os.Getwd()
	if path == "" {
		path = "/???"
	}
	return filepath.ToSlash(path)
}










func AbsFile(dir, file, rewrites string) string {
	abs := file
	if dir != "" && !filepath.IsAbs(file) {
		abs = filepath.Join(dir, file)
	}

	start := 0
	for i := 0; i <= len(rewrites); i++ {
		if i == len(rewrites) || rewrites[i] == ';' {
			if new, ok := applyRewrite(abs, rewrites[start:i]); ok {
				abs = new
				goto Rewritten
			}
			start = i + 1
		}
	}
	if hasPathPrefix(abs, GOROOT) {
		abs = "$GOROOT" + abs[len(GOROOT):]
	}

Rewritten:
	if abs == "" {
		abs = "??"
	}
	return abs
}




func applyRewrite(path, rewrite string) (string, bool) {
	prefix, replace := rewrite, ""
	if j := strings.LastIndex(rewrite, "=>"); j >= 0 {
		prefix, replace = rewrite[:j], rewrite[j+len("=>"):]
	}

	if prefix == "" || !hasPathPrefix(path, prefix) {
		return path, false
	}
	if len(path) == len(prefix) {
		return replace, true
	}
	if replace == "" {
		return path[len(prefix)+1:], true
	}
	return replace + path[len(prefix):], true
}








func hasPathPrefix(s string, t string) bool {
	if len(t) > len(s) {
		return false
	}
	var i int
	for i = 0; i < len(t); i++ {
		cs := int(s[i])
		ct := int(t[i])
		if 'A' <= cs && cs <= 'Z' {
			cs += 'a' - 'A'
		}
		if 'A' <= ct && ct <= 'Z' {
			ct += 'a' - 'A'
		}
		if cs == '\\' {
			cs = '/'
		}
		if ct == '\\' {
			ct = '/'
		}
		if cs != ct {
			return false
		}
	}
	return i >= len(s) || s[i] == '/' || s[i] == '\\'
}
