









package viper

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/spf13/cast"
)


type ConfigParseError struct {
	err error
}


func (pe ConfigParseError) Error() string {
	return fmt.Sprintf("While parsing config: %s", pe.err.Error())
}


func (pe ConfigParseError) Unwrap() error {
	return pe.err
}



func toCaseInsensitiveValue(value any) any {
	switch v := value.(type) {
	case map[any]any:
		value = copyAndInsensitiviseMap(cast.ToStringMap(v))
	case map[string]any:
		value = copyAndInsensitiviseMap(v)
	}

	return value
}



func copyAndInsensitiviseMap(m map[string]any) map[string]any {
	nm := make(map[string]any)

	for key, val := range m {
		lkey := strings.ToLower(key)
		switch v := val.(type) {
		case map[any]any:
			nm[lkey] = copyAndInsensitiviseMap(cast.ToStringMap(v))
		case map[string]any:
			nm[lkey] = copyAndInsensitiviseMap(v)
		default:
			nm[lkey] = v
		}
	}

	return nm
}

func insensitiviseVal(val any) any {
	switch v := val.(type) {
	case map[any]any:
		
		val = cast.ToStringMap(val)
		insensitiviseMap(val.(map[string]any))
	case map[string]any:
		
		insensitiviseMap(v)
	case []any:
		
		insensitiveArray(v)
	}
	return val
}

func insensitiviseMap(m map[string]any) {
	for key, val := range m {
		val = insensitiviseVal(val)
		lower := strings.ToLower(key)
		if key != lower {
			
			delete(m, key)
		}
		
		m[lower] = val
	}
}

func insensitiveArray(a []any) {
	for i, val := range a {
		a[i] = insensitiviseVal(val)
	}
}

func absPathify(logger *slog.Logger, inPath string) string {
	logger.Info("trying to resolve absolute path", "path", inPath)

	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		inPath = userHomeDir() + inPath[5:]
	}

	inPath = os.ExpandEnv(inPath)

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	logger.Error(fmt.Errorf("could not discover absolute path: %w", err).Error())

	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func safeMul(a, b uint) uint {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		return 0
	}
	return c
}


func parseSizeInBytes(sizeStr string) uint {
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			if lastChar > 1 {
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size := cast.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

	return safeMul(uint(size), multiplier)
}








func deepSearch(m map[string]any, path []string) map[string]any {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			
			
			m3 := make(map[string]any)
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]any)
		if !ok {
			
			
			m3 = make(map[string]any)
			m[k] = m3
		}
		
		m = m3
	}
	return m
}
