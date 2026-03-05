
package gotenv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	
	linePattern = `\A\s*(?:export\s+)?([\w\.]+)(?:\s*=\s*|:\s+?)('(?:\'|[^'])*'|"(?:\"|[^"])*"|[^#\n]+)?\s*(?:\s*\#.*)?\z`

	
	variablePattern = `(\\)?(\$)(\{?([A-Z0-9_]+)?\}?)`
)


var (
	bomUTF8    = []byte("\xEF\xBB\xBF")
	bomUTF16LE = []byte("\xFF\xFE")
	bomUTF16BE = []byte("\xFE\xFF")
)


type Env map[string]string




func Load(filenames ...string) error {
	return loadenv(false, filenames...)
}


func OverLoad(filenames ...string) error {
	return loadenv(true, filenames...)
}


func Must(fn func(filenames ...string) error, filenames ...string) {
	if err := fn(filenames...); err != nil {
		panic(err.Error())
	}
}


func Apply(r io.Reader) error {
	return parset(r, false)
}


func OverApply(r io.Reader) error {
	return parset(r, true)
}

func loadenv(override bool, filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		err = parset(f, override)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}


func parset(r io.Reader, override bool) error {
	env, err := strictParse(r, override)
	if err != nil {
		return err
	}

	for key, val := range env {
		setenv(key, val, override)
	}

	return nil
}

func setenv(key, val string, override bool) {
	if override {
		os.Setenv(key, val)
	} else {
		if _, present := os.LookupEnv(key); !present {
			os.Setenv(key, val)
		}
	}
}




func Parse(r io.Reader) Env {
	env, _ := strictParse(r, false)
	return env
}




func StrictParse(r io.Reader) (Env, error) {
	return strictParse(r, false)
}




func Read(filename string) (Env, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return strictParse(f, false)
}




func Unmarshal(str string) (Env, error) {
	return strictParse(strings.NewReader(str), false)
}



func Marshal(env Env) (string, error) {
	lines := make([]string, 0, len(env))
	for k, v := range env {
		if d, err := strconv.Atoi(v); err == nil {
			lines = append(lines, fmt.Sprintf(`%s=%d`, k, d))
		} else {
			lines = append(lines, fmt.Sprintf(`%s=%q`, k, v))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}


func Write(env Env, filename string) error {
	content, err := Marshal(env)
	if err != nil {
		return err
	}
	
	if err := os.MkdirAll(filepath.Dir(filename), 0o775); err != nil {
		return err
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}

	return file.Sync()
}



func splitLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, bufio.ErrFinalToken
	}

	idx := bytes.IndexAny(data, "\r\n")
	switch {
	case atEOF && idx < 0:
		return len(data), data, bufio.ErrFinalToken

	case idx < 0:
		return 0, nil, nil
	}

	
	eol := idx + 1
	
	if len(data) > eol && data[eol-1] == '\r' && data[eol] == '\n' {
		eol++
	}

	return eol, data[:idx], nil
}

func strictParse(r io.Reader, override bool) (Env, error) {
	env := make(Env)

	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)

	
	bomByteBuffer := make([]byte, 3)
	_, err := tee.Read(bomByteBuffer)
	if err != nil && err != io.EOF {
		return env, err
	}

	z := io.MultiReader(buf, r)

	
	var scanner *bufio.Scanner

	if bytes.HasPrefix(bomByteBuffer, bomUTF8) {
		scanner = bufio.NewScanner(transform.NewReader(z, unicode.UTF8BOM.NewDecoder()))
	} else if bytes.HasPrefix(bomByteBuffer, bomUTF16LE) {
		scanner = bufio.NewScanner(transform.NewReader(z, unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder()))
	} else if bytes.HasPrefix(bomByteBuffer, bomUTF16BE) {
		scanner = bufio.NewScanner(transform.NewReader(z, unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewDecoder()))
	} else {
		scanner = bufio.NewScanner(z)
	}

	scanner.Split(splitLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return env, err
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		quote := ""
		
		idx := strings.Index(line, "=")
		if idx == -1 {
			idx = strings.Index(line, ":")
		}
		
		if idx > 0 && idx < len(line)-1 {
			val := strings.TrimSpace(line[idx+1:])
			if val[0] == '"' || val[0] == '\'' {
				quote = val[:1]
				
				idx = strings.LastIndex(strings.TrimSpace(val[1:]), quote)
				if idx >= 0 && val[idx] != '\\' {
					quote = ""
				}
			}
		}
		
		for quote != "" && scanner.Scan() {
			l := scanner.Text()
			line += "\n" + l
			idx := strings.LastIndex(l, quote)
			if idx > 0 && l[idx-1] == '\\' {
				
				continue
			}
			if idx >= 0 {
				
				quote = ""
			}
		}

		if quote != "" {
			return env, fmt.Errorf("missing quotes")
		}

		err := parseLine(line, env, override)
		if err != nil {
			return env, err
		}
	}

	return env, scanner.Err()
}

var (
	lineRgx     = regexp.MustCompile(linePattern)
	unescapeRgx = regexp.MustCompile(`\\([^$])`)
	varRgx      = regexp.MustCompile(variablePattern)
)

func parseLine(s string, env Env, override bool) error {
	rm := lineRgx.FindStringSubmatch(s)

	if len(rm) == 0 {
		return checkFormat(s, env)
	}

	key := strings.TrimSpace(rm[1])
	val := strings.TrimSpace(rm[2])

	var hsq, hdq bool

	
	if l := len(val); l >= 2 {
		l -= 1
		
		hdq = val[0] == '"' && val[l] == '"'
		
		hsq = val[0] == '\'' && val[l] == '\''

		
		if hsq || hdq {
			val = val[1:l]
		}
	}

	if hdq {
		val = strings.ReplaceAll(val, `\n`, "\n")
		val = strings.ReplaceAll(val, `\r`, "\r")

		
		val = unescapeRgx.ReplaceAllString(val, "$1")
	}

	if !hsq {
		fv := func(s string) string {
			return varReplacement(s, hsq, env, override)
		}
		val = varRgx.ReplaceAllStringFunc(val, fv)
	}

	env[key] = val
	return nil
}

func parseExport(st string, env Env) error {
	if strings.HasPrefix(st, "export") {
		vs := strings.SplitN(st, " ", 2)

		if len(vs) > 1 {
			if _, ok := env[vs[1]]; !ok {
				return fmt.Errorf("line `%s` has an unset variable", st)
			}
		}
	}

	return nil
}

var varNameRgx = regexp.MustCompile(`(\$)(\{?([A-Z0-9_]+)\}?)`)

func varReplacement(s string, hsq bool, env Env, override bool) string {
	if s == "" {
		return s
	}

	if s[0] == '\\' {
		
		return s[1:]
	}

	if hsq {
		return s
	}

	mn := varNameRgx.FindStringSubmatch(s)

	if len(mn) == 0 {
		return s
	}

	v := mn[3]

	if replace, ok := os.LookupEnv(v); ok && !override {
		return replace
	}

	if replace, ok := env[v]; ok {
		return replace
	}

	return os.Getenv(v)
}

func checkFormat(s string, env Env) error {
	st := strings.TrimSpace(s)

	if st == "" || st[0] == '#' {
		return nil
	}

	if err := parseExport(st, env); err != nil {
		return err
	}

	return fmt.Errorf("line `%s` doesn't match format", s)
}
