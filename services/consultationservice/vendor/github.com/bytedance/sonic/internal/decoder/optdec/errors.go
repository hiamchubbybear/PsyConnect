

 package optdec

 import (
	 "encoding/json"
	 "errors"
	 "reflect"
	 "strconv"
 
	 "github.com/bytedance/sonic/internal/rt"
 )

 
 
 var stackOverflow = &json.UnsupportedValueError{
	 Str:   "Value nesting too deep",
	 Value: reflect.ValueOf("..."),
 }
 
 func error_type(vt *rt.GoType) error {
	 return &json.UnmarshalTypeError{Type: vt.Pack()}
 }
 
 func error_mismatch(node Node, ctx *context, typ reflect.Type) error {
	 return MismatchTypeError{
		 Pos:  node.Position(),
		 Src:  ctx.Parser.Json,
		 Type: typ,
	 }
 }
 
 func newUnmatched(pos int, vt *rt.GoType) error {
	 return MismatchTypeError{
		Pos:  pos,
		Src:  "",
		Type: vt.Pack(),
	 }
 }

 func error_field(name string) error {
	 return errors.New("json: unknown field " + strconv.Quote(name))
 }
 
 func error_value(value string, vtype reflect.Type) error {
	 return &json.UnmarshalTypeError{
		 Type:  vtype,
		 Value: value,
	 }
 }
 
 func error_syntax(pos int, src string, msg string) error {
	 return SyntaxError{
		 Pos: pos,
		 Src: src,
		 Msg: msg,
	 }
 }

 func error_unsuppoted(typ *rt.GoType) error {
	return &json.UnsupportedTypeError{
		Type: typ.Pack(),
	}
}
 