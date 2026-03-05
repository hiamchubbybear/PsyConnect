package json

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/goccy/go-json/internal/decoder"
	"github.com/goccy/go-json/internal/errors"
	"github.com/goccy/go-json/internal/runtime"
)

type Decoder struct {
	s *decoder.Stream
}

const (
	nul = '\000'
)

type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

func unmarshal(data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	src := make([]byte, len(data)+1) 
	copy(src, data)

	header := (*emptyInterface)(unsafe.Pointer(&v))

	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
		return err
	}
	dec, err := decoder.CompileToGetDecoder(header.typ)
	if err != nil {
		return err
	}
	ctx := decoder.TakeRuntimeContext()
	ctx.Buf = src
	ctx.Option.Flags = 0
	for _, optFunc := range optFuncs {
		optFunc(ctx.Option)
	}
	cursor, err := dec.Decode(ctx, 0, 0, header.ptr)
	if err != nil {
		decoder.ReleaseRuntimeContext(ctx)
		return err
	}
	decoder.ReleaseRuntimeContext(ctx)
	return validateEndBuf(src, cursor)
}

func unmarshalContext(ctx context.Context, data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	src := make([]byte, len(data)+1) 
	copy(src, data)

	header := (*emptyInterface)(unsafe.Pointer(&v))

	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
		return err
	}
	dec, err := decoder.CompileToGetDecoder(header.typ)
	if err != nil {
		return err
	}
	rctx := decoder.TakeRuntimeContext()
	rctx.Buf = src
	rctx.Option.Flags = 0
	rctx.Option.Flags |= decoder.ContextOption
	rctx.Option.Context = ctx
	for _, optFunc := range optFuncs {
		optFunc(rctx.Option)
	}
	cursor, err := dec.Decode(rctx, 0, 0, header.ptr)
	if err != nil {
		decoder.ReleaseRuntimeContext(rctx)
		return err
	}
	decoder.ReleaseRuntimeContext(rctx)
	return validateEndBuf(src, cursor)
}

var (
	pathDecoder = decoder.NewPathDecoder()
)

func extractFromPath(path *Path, data []byte, optFuncs ...DecodeOptionFunc) ([][]byte, error) {
	if path.path.RootSelectorOnly {
		return [][]byte{data}, nil
	}
	src := make([]byte, len(data)+1) 
	copy(src, data)

	ctx := decoder.TakeRuntimeContext()
	ctx.Buf = src
	ctx.Option.Flags = 0
	ctx.Option.Flags |= decoder.PathOption
	ctx.Option.Path = path.path
	for _, optFunc := range optFuncs {
		optFunc(ctx.Option)
	}
	paths, cursor, err := pathDecoder.DecodePath(ctx, 0, 0)
	if err != nil {
		decoder.ReleaseRuntimeContext(ctx)
		return nil, err
	}
	decoder.ReleaseRuntimeContext(ctx)
	if err := validateEndBuf(src, cursor); err != nil {
		return nil, err
	}
	return paths, nil
}

func unmarshalNoEscape(data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	src := make([]byte, len(data)+1) 
	copy(src, data)

	header := (*emptyInterface)(unsafe.Pointer(&v))

	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
		return err
	}
	dec, err := decoder.CompileToGetDecoder(header.typ)
	if err != nil {
		return err
	}

	ctx := decoder.TakeRuntimeContext()
	ctx.Buf = src
	ctx.Option.Flags = 0
	for _, optFunc := range optFuncs {
		optFunc(ctx.Option)
	}
	cursor, err := dec.Decode(ctx, 0, 0, noescape(header.ptr))
	if err != nil {
		decoder.ReleaseRuntimeContext(ctx)
		return err
	}
	decoder.ReleaseRuntimeContext(ctx)
	return validateEndBuf(src, cursor)
}

func validateEndBuf(src []byte, cursor int64) error {
	for {
		switch src[cursor] {
		case ' ', '\t', '\n', '\r':
			cursor++
			continue
		case nul:
			return nil
		}
		return errors.ErrSyntax(
			fmt.Sprintf("invalid character '%c' after top-level value", src[cursor]),
			cursor+1,
		)
	}
}


//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func validateType(typ *runtime.Type, p uintptr) error {
	if typ == nil || typ.Kind() != reflect.Ptr || p == 0 {
		return &InvalidUnmarshalError{Type: runtime.RType2Type(typ)}
	}
	return nil
}





func NewDecoder(r io.Reader) *Decoder {
	s := decoder.NewStream(r)
	return &Decoder{
		s: s,
	}
}



func (d *Decoder) Buffered() io.Reader {
	return d.s.Buffered()
}






func (d *Decoder) Decode(v interface{}) error {
	return d.DecodeWithOption(v)
}



func (d *Decoder) DecodeContext(ctx context.Context, v interface{}) error {
	d.s.Option.Flags |= decoder.ContextOption
	d.s.Option.Context = ctx
	return d.DecodeWithOption(v)
}

func (d *Decoder) DecodeWithOption(v interface{}, optFuncs ...DecodeOptionFunc) error {
	header := (*emptyInterface)(unsafe.Pointer(&v))
	typ := header.typ
	ptr := uintptr(header.ptr)
	typeptr := uintptr(unsafe.Pointer(typ))
	
	copiedType := *(**runtime.Type)(unsafe.Pointer(&typeptr))

	if err := validateType(copiedType, ptr); err != nil {
		return err
	}

	dec, err := decoder.CompileToGetDecoder(typ)
	if err != nil {
		return err
	}
	if err := d.s.PrepareForDecode(); err != nil {
		return err
	}
	s := d.s
	for _, optFunc := range optFuncs {
		optFunc(s.Option)
	}
	if err := dec.DecodeStream(s, 0, header.ptr); err != nil {
		return err
	}
	s.Reset()
	return nil
}

func (d *Decoder) More() bool {
	return d.s.More()
}

func (d *Decoder) Token() (Token, error) {
	return d.s.Token()
}




func (d *Decoder) DisallowUnknownFields() {
	d.s.DisallowUnknownFields = true
}

func (d *Decoder) InputOffset() int64 {
	return d.s.TotalOffset()
}



func (d *Decoder) UseNumber() {
	d.s.UseNumber = true
}
