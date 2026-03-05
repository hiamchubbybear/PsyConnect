



















package funcr

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
)


func New(fn func(prefix, args string), opts Options) logr.Logger {
	return logr.New(newSink(fn, NewFormatter(opts)))
}



func NewJSON(fn func(obj string), opts Options) logr.Logger {
	fnWrapper := func(_, obj string) {
		fn(obj)
	}
	return logr.New(newSink(fnWrapper, NewFormatterJSON(opts)))
}





type Underlier interface {
	GetUnderlying() func(prefix, args string)
}

func newSink(fn func(prefix, args string), formatter Formatter) logr.LogSink {
	l := &fnlogger{
		Formatter: formatter,
		write:     fn,
	}
	
	l.Formatter.AddCallDepth(1)
	return l
}


type Options struct {
	
	
	LogCaller MessageClass

	
	
	LogCallerFunc bool

	
	
	LogTimestamp bool

	
	
	
	TimestampFormat string

	
	
	
	LogInfoLevel *string

	
	
	
	Verbosity int

	
	
	
	
	
	
	
	
	
	
	
	RenderBuiltinsHook func(kvList []any) []any

	
	
	
	RenderValuesHook func(kvList []any) []any

	
	
	
	RenderArgsHook func(kvList []any) []any

	
	
	
	
	
	
	MaxLogDepth int
}


type MessageClass int

const (
	
	None MessageClass = iota
	
	All
	
	Info
	
	Error
)



type fnlogger struct {
	Formatter
	write func(prefix, args string)
}

func (l fnlogger) WithName(name string) logr.LogSink {
	l.Formatter.AddName(name)
	return &l
}

func (l fnlogger) WithValues(kvList ...any) logr.LogSink {
	l.Formatter.AddValues(kvList)
	return &l
}

func (l fnlogger) WithCallDepth(depth int) logr.LogSink {
	l.Formatter.AddCallDepth(depth)
	return &l
}

func (l fnlogger) Info(level int, msg string, kvList ...any) {
	prefix, args := l.FormatInfo(level, msg, kvList)
	l.write(prefix, args)
}

func (l fnlogger) Error(err error, msg string, kvList ...any) {
	prefix, args := l.FormatError(err, msg, kvList)
	l.write(prefix, args)
}

func (l fnlogger) GetUnderlying() func(prefix, args string) {
	return l.write
}


var _ logr.LogSink = &fnlogger{}
var _ logr.CallDepthLogSink = &fnlogger{}
var _ Underlier = &fnlogger{}


func NewFormatter(opts Options) Formatter {
	return newFormatter(opts, outputKeyValue)
}


func NewFormatterJSON(opts Options) Formatter {
	return newFormatter(opts, outputJSON)
}


const defaultTimestampFormat = "2006-01-02 15:04:05.000000"
const defaultMaxLogDepth = 16

func newFormatter(opts Options, outfmt outputFormat) Formatter {
	if opts.TimestampFormat == "" {
		opts.TimestampFormat = defaultTimestampFormat
	}
	if opts.MaxLogDepth == 0 {
		opts.MaxLogDepth = defaultMaxLogDepth
	}
	if opts.LogInfoLevel == nil {
		opts.LogInfoLevel = new(string)
		*opts.LogInfoLevel = "level"
	}
	f := Formatter{
		outputFormat: outfmt,
		prefix:       "",
		values:       nil,
		depth:        0,
		opts:         &opts,
	}
	return f
}




type Formatter struct {
	outputFormat outputFormat
	prefix       string
	values       []any
	valuesStr    string
	depth        int
	opts         *Options
	groupName    string 
	groups       []groupDef
}


type outputFormat int

const (
	
	outputKeyValue outputFormat = iota
	
	outputJSON
)



type groupDef struct {
	name   string
	values string
}


type PseudoStruct []any


func (f Formatter) render(builtins, args []any) string {
	
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	if f.outputFormat == outputJSON {
		buf.WriteByte('{') 
	}

	
	vals := builtins
	if hook := f.opts.RenderBuiltinsHook; hook != nil {
		vals = hook(f.sanitize(vals))
	}
	f.flatten(buf, vals, false) 
	continuing := len(builtins) > 0

	
	argsStr := func() string {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))

		vals = args
		if hook := f.opts.RenderArgsHook; hook != nil {
			vals = hook(f.sanitize(vals))
		}
		f.flatten(buf, vals, true) 

		return buf.String()
	}()

	
	bodyStr := f.renderGroup(f.groupName, f.valuesStr, argsStr)
	for i := len(f.groups) - 1; i >= 0; i-- {
		grp := &f.groups[i]
		if grp.values == "" && bodyStr == "" {
			
			continue
		}
		bodyStr = f.renderGroup(grp.name, grp.values, bodyStr)
	}

	if bodyStr != "" {
		if continuing {
			buf.WriteByte(f.comma())
		}
		buf.WriteString(bodyStr)
	}

	if f.outputFormat == outputJSON {
		buf.WriteByte('}') 
	}

	return buf.String()
}







func (f Formatter) renderGroup(name string, values string, args string) string {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	needClosingBrace := false
	if name != "" && (values != "" || args != "") {
		buf.WriteString(f.quoted(name, true)) 
		buf.WriteByte(f.colon())
		buf.WriteByte('{')
		needClosingBrace = true
	}

	continuing := false
	if values != "" {
		buf.WriteString(values)
		continuing = true
	}

	if args != "" {
		if continuing {
			buf.WriteByte(f.comma())
		}
		buf.WriteString(args)
	}

	if needClosingBrace {
		buf.WriteByte('}')
	}

	return buf.String()
}








func (f Formatter) flatten(buf *bytes.Buffer, kvList []any, escapeKeys bool) []any {
	
	
	if len(kvList)%2 != 0 {
		kvList = append(kvList, noValue)
	}
	copied := false
	for i := 0; i < len(kvList); i += 2 {
		k, ok := kvList[i].(string)
		if !ok {
			if !copied {
				newList := make([]any, len(kvList))
				copy(newList, kvList)
				kvList = newList
				copied = true
			}
			k = f.nonStringKey(kvList[i])
			kvList[i] = k
		}
		v := kvList[i+1]

		if i > 0 {
			if f.outputFormat == outputJSON {
				buf.WriteByte(f.comma())
			} else {
				
				
				buf.WriteByte(' ')
			}
		}

		buf.WriteString(f.quoted(k, escapeKeys))
		buf.WriteByte(f.colon())
		buf.WriteString(f.pretty(v))
	}
	return kvList
}

func (f Formatter) quoted(str string, escape bool) string {
	if escape {
		return prettyString(str)
	}
	
	return `"` + str + `"`
}

func (f Formatter) comma() byte {
	if f.outputFormat == outputJSON {
		return ','
	}
	return ' '
}

func (f Formatter) colon() byte {
	if f.outputFormat == outputJSON {
		return ':'
	}
	return '='
}

func (f Formatter) pretty(value any) string {
	return f.prettyWithFlags(value, 0, 0)
}

const (
	flagRawStruct = 0x1 
)


func (f Formatter) prettyWithFlags(value any, flags uint32, depth int) string {
	if depth > f.opts.MaxLogDepth {
		return `"<max-log-depth-exceeded>"`
	}

	
	if v, ok := value.(logr.Marshaler); ok {
		
		
		value = invokeMarshaler(v)
	}

	
	switch v := value.(type) {
	case fmt.Stringer:
		value = invokeStringer(v)
	case error:
		value = invokeError(v)
	}

	
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case string:
		return prettyString(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case complex64:
		return `"` + strconv.FormatComplex(complex128(v), 'f', -1, 64) + `"`
	case complex128:
		return `"` + strconv.FormatComplex(v, 'f', -1, 128) + `"`
	case PseudoStruct:
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		v = f.sanitize(v)
		if flags&flagRawStruct == 0 {
			buf.WriteByte('{')
		}
		for i := 0; i < len(v); i += 2 {
			if i > 0 {
				buf.WriteByte(f.comma())
			}
			k, _ := v[i].(string) 
			
			buf.WriteString(prettyString(k))
			buf.WriteByte(f.colon())
			buf.WriteString(f.prettyWithFlags(v[i+1], 0, depth+1))
		}
		if flags&flagRawStruct == 0 {
			buf.WriteByte('}')
		}
		return buf.String()
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	t := reflect.TypeOf(value)
	if t == nil {
		return "null"
	}
	v := reflect.ValueOf(value)
	switch t.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return prettyString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(int64(v.Int()), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(uint64(v.Uint()), 10)
	case reflect.Float32:
		return strconv.FormatFloat(float64(v.Float()), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64:
		return `"` + strconv.FormatComplex(complex128(v.Complex()), 'f', -1, 64) + `"`
	case reflect.Complex128:
		return `"` + strconv.FormatComplex(v.Complex(), 'f', -1, 128) + `"`
	case reflect.Struct:
		if flags&flagRawStruct == 0 {
			buf.WriteByte('{')
		}
		printComma := false 
		for i := 0; i < t.NumField(); i++ {
			fld := t.Field(i)
			if fld.PkgPath != "" {
				
				continue
			}
			if !v.Field(i).CanInterface() {
				
				continue
			}
			name := ""
			omitempty := false
			if tag, found := fld.Tag.Lookup("json"); found {
				if tag == "-" {
					continue
				}
				if comma := strings.Index(tag, ","); comma != -1 {
					if n := tag[:comma]; n != "" {
						name = n
					}
					rest := tag[comma:]
					if strings.Contains(rest, ",omitempty,") || strings.HasSuffix(rest, ",omitempty") {
						omitempty = true
					}
				} else {
					name = tag
				}
			}
			if omitempty && isEmpty(v.Field(i)) {
				continue
			}
			if printComma {
				buf.WriteByte(f.comma())
			}
			printComma = true 
			if fld.Anonymous && fld.Type.Kind() == reflect.Struct && name == "" {
				buf.WriteString(f.prettyWithFlags(v.Field(i).Interface(), flags|flagRawStruct, depth+1))
				continue
			}
			if name == "" {
				name = fld.Name
			}
			
			buf.WriteString(f.quoted(name, false))
			buf.WriteByte(f.colon())
			buf.WriteString(f.prettyWithFlags(v.Field(i).Interface(), 0, depth+1))
		}
		if flags&flagRawStruct == 0 {
			buf.WriteByte('}')
		}
		return buf.String()
	case reflect.Slice, reflect.Array:
		
		
		
		if f.outputFormat == outputJSON {
			if rm, ok := value.(json.RawMessage); ok {
				
				if len(rm) > 0 {
					buf.Write(rm)
				} else {
					buf.WriteString("null")
				}
				return buf.String()
			}
		}
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(f.comma())
			}
			e := v.Index(i)
			buf.WriteString(f.prettyWithFlags(e.Interface(), 0, depth+1))
		}
		buf.WriteByte(']')
		return buf.String()
	case reflect.Map:
		buf.WriteByte('{')
		
		it := v.MapRange()
		i := 0
		for it.Next() {
			if i > 0 {
				buf.WriteByte(f.comma())
			}
			
			keystr := ""
			if m, ok := it.Key().Interface().(encoding.TextMarshaler); ok {
				txt, err := m.MarshalText()
				if err != nil {
					keystr = fmt.Sprintf("<error-MarshalText: %s>", err.Error())
				} else {
					keystr = string(txt)
				}
				keystr = prettyString(keystr)
			} else {
				
				keystr = f.prettyWithFlags(it.Key().Interface(), 0, depth+1)
				if t.Key().Kind() != reflect.String {
					
					
					keystr = prettyString(keystr)
				}
			}
			buf.WriteString(keystr)
			buf.WriteByte(f.colon())
			buf.WriteString(f.prettyWithFlags(it.Value().Interface(), 0, depth+1))
			i++
		}
		buf.WriteByte('}')
		return buf.String()
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "null"
		}
		return f.prettyWithFlags(v.Elem().Interface(), 0, depth)
	}
	return fmt.Sprintf(`"<unhandled-%s>"`, t.Kind().String())
}

func prettyString(s string) string {
	
	if needsEscape(s) {
		return strconv.Quote(s)
	}
	b := bytes.NewBuffer(make([]byte, 0, 1024))
	b.WriteByte('"')
	b.WriteString(s)
	b.WriteByte('"')
	return b.String()
}



func needsEscape(s string) bool {
	for _, r := range s {
		if !strconv.IsPrint(r) || r == '\\' || r == '"' {
			return true
		}
	}
	return false
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func invokeMarshaler(m logr.Marshaler) (ret any) {
	defer func() {
		if r := recover(); r != nil {
			ret = fmt.Sprintf("<panic: %s>", r)
		}
	}()
	return m.MarshalLog()
}

func invokeStringer(s fmt.Stringer) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			ret = fmt.Sprintf("<panic: %s>", r)
		}
	}()
	return s.String()
}

func invokeError(e error) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			ret = fmt.Sprintf("<panic: %s>", r)
		}
	}()
	return e.Error()
}







type Caller struct {
	
	File string `json:"file"`
	
	Line int `json:"line"`
	
	
	Func string `json:"function,omitempty"`
}

func (f Formatter) caller() Caller {
	
	pc, file, line, ok := runtime.Caller(f.depth + 2)
	if !ok {
		return Caller{"<unknown>", 0, ""}
	}
	fn := ""
	if f.opts.LogCallerFunc {
		if fp := runtime.FuncForPC(pc); fp != nil {
			fn = fp.Name()
		}
	}

	return Caller{filepath.Base(file), line, fn}
}

const noValue = "<no-value>"

func (f Formatter) nonStringKey(v any) string {
	return fmt.Sprintf("<non-string-key: %s>", f.snippet(v))
}


func (f Formatter) snippet(v any) string {
	const snipLen = 16

	snip := f.pretty(v)
	if len(snip) > snipLen {
		snip = snip[:snipLen]
	}
	return snip
}




func (f Formatter) sanitize(kvList []any) []any {
	if len(kvList)%2 != 0 {
		kvList = append(kvList, noValue)
	}
	for i := 0; i < len(kvList); i += 2 {
		_, ok := kvList[i].(string)
		if !ok {
			kvList[i] = f.nonStringKey(kvList[i])
		}
	}
	return kvList
}




func (f *Formatter) startGroup(name string) {
	
	if name == "" {
		return
	}

	n := len(f.groups)
	f.groups = append(f.groups[:n:n], groupDef{f.groupName, f.valuesStr})

	
	f.groupName = name
	f.valuesStr = ""
	f.values = nil
}




func (f *Formatter) Init(info logr.RuntimeInfo) {
	f.depth += info.CallDepth
}


func (f Formatter) Enabled(level int) bool {
	return level <= f.opts.Verbosity
}



func (f Formatter) GetDepth() int {
	return f.depth
}




func (f Formatter) FormatInfo(level int, msg string, kvList []any) (prefix, argsStr string) {
	args := make([]any, 0, 64) 
	prefix = f.prefix
	if f.outputFormat == outputJSON {
		args = append(args, "logger", prefix)
		prefix = ""
	}
	if f.opts.LogTimestamp {
		args = append(args, "ts", time.Now().Format(f.opts.TimestampFormat))
	}
	if policy := f.opts.LogCaller; policy == All || policy == Info {
		args = append(args, "caller", f.caller())
	}
	if key := *f.opts.LogInfoLevel; key != "" {
		args = append(args, key, level)
	}
	args = append(args, "msg", msg)
	return prefix, f.render(args, kvList)
}




func (f Formatter) FormatError(err error, msg string, kvList []any) (prefix, argsStr string) {
	args := make([]any, 0, 64) 
	prefix = f.prefix
	if f.outputFormat == outputJSON {
		args = append(args, "logger", prefix)
		prefix = ""
	}
	if f.opts.LogTimestamp {
		args = append(args, "ts", time.Now().Format(f.opts.TimestampFormat))
	}
	if policy := f.opts.LogCaller; policy == All || policy == Error {
		args = append(args, "caller", f.caller())
	}
	args = append(args, "msg", msg)
	var loggableErr any
	if err != nil {
		loggableErr = err.Error()
	}
	args = append(args, "error", loggableErr)
	return prefix, f.render(args, kvList)
}




func (f *Formatter) AddName(name string) {
	if len(f.prefix) > 0 {
		f.prefix += "/"
	}
	f.prefix += name
}



func (f *Formatter) AddValues(kvList []any) {
	
	n := len(f.values)
	f.values = append(f.values[:n:n], kvList...)

	vals := f.values
	if hook := f.opts.RenderValuesHook; hook != nil {
		vals = hook(f.sanitize(vals))
	}

	
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	f.flatten(buf, vals, true) 
	f.valuesStr = buf.String()
}



func (f *Formatter) AddCallDepth(depth int) {
	f.depth += depth
}
