

package sonic

import (
    `io`

    `github.com/bytedance/sonic/ast`
    `github.com/bytedance/sonic/internal/rt`
)

const (
    
	UseStdJSON = iota
    
	UseSonicJSON
)


const APIKind = apiKind


type Config struct {
    
    
    
    EscapeHTML                    bool

    
    
    
    SortMapKeys                   bool

    
    
    CompactMarshaler              bool

    
    
    NoQuoteTextMarshaler          bool

    
    
    NoNullSliceOrMap              bool

    
    
    UseInt64                      bool

    
    
    UseNumber                     bool

    
    
    UseUnicodeErrors              bool

    
    
    
    DisallowUnknownFields         bool

    
    CopyString                    bool

    
    
    ValidateString                bool

    
    
    NoValidateJSONMarshaler       bool

    
    
    NoValidateJSONSkip bool
    
    
    NoEncoderNewline bool

    
    EncodeNullForInfOrNan bool
}
 
var (
    
    ConfigDefault = Config{}.Froze()
 
    
    ConfigStd = Config{
        EscapeHTML : true,
        SortMapKeys: true,
        CompactMarshaler: true,
        CopyString : true,
        ValidateString : true,
    }.Froze()
 
    
    ConfigFastest = Config{
        NoQuoteTextMarshaler: true,
        NoValidateJSONMarshaler: true,
        NoValidateJSONSkip: true,
    }.Froze()
)
 
 



type API interface {
    
    MarshalToString(v interface{}) (string, error)
    
    Marshal(v interface{}) ([]byte, error)
    
    MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
    
    UnmarshalFromString(str string, v interface{}) error
    
    Unmarshal(data []byte, v interface{}) error
    
    NewEncoder(writer io.Writer) Encoder
    
    NewDecoder(reader io.Reader) Decoder
    
    Valid(data []byte) bool
}


type Encoder interface {
    
    Encode(val interface{}) error
    
    
    
    SetEscapeHTML(on bool)
    
    
    
    SetIndent(prefix, indent string)
}


type Decoder interface {
    
    Decode(val interface{}) error
    
    
    Buffered() io.Reader
    
    
    DisallowUnknownFields()
    
    More() bool
    
    UseNumber()
}


func Marshal(val interface{}) ([]byte, error) {
    return ConfigDefault.Marshal(val)
}




func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
    return ConfigDefault.MarshalIndent(v, prefix, indent)
}


func MarshalString(val interface{}) (string, error) {
    return ConfigDefault.MarshalToString(val)
}




func Unmarshal(buf []byte, val interface{}) error {
    return ConfigDefault.Unmarshal(buf, val)
}


func UnmarshalString(buf string, val interface{}) error {
    return ConfigDefault.UnmarshalFromString(buf, val)
}












func Get(src []byte, path ...interface{}) (ast.Node, error) {
    return GetCopyFromString(rt.Mem2Str(src), path...)
}



func GetWithOptions(src []byte, opts ast.SearchOptions, path ...interface{}) (ast.Node, error) {
    s := ast.NewSearcher(rt.Mem2Str(src))
    s.SearchOptions = opts
    return s.GetByPath(path...)
}






func GetFromString(src string, path ...interface{}) (ast.Node, error) {
    return ast.NewSearcher(src).GetByPath(path...)
}


func GetCopyFromString(src string, path ...interface{}) (ast.Node, error) {
    return ast.NewSearcher(src).GetByPathCopy(path...)
}


func Valid(data []byte) bool {
    return ConfigDefault.Valid(data)
}


func ValidString(data string) bool {
    return ConfigDefault.Valid(rt.Str2Mem(data))
}
