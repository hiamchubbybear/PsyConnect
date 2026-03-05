package json

import (
	"io"

	"github.com/goccy/go-json/internal/decoder"
	"github.com/goccy/go-json/internal/encoder"
)

type EncodeOption = encoder.Option
type EncodeOptionFunc func(*EncodeOption)


func UnorderedMap() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag |= encoder.UnorderedMapOption
	}
}


func DisableHTMLEscape() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag &= ^encoder.HTMLEscapeOption
	}
}





func DisableNormalizeUTF8() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag &= ^encoder.NormalizeUTF8Option
	}
}


func Debug() EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag |= encoder.DebugOption
	}
}


func DebugWith(w io.Writer) EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.DebugOut = w
	}
}


func DebugDOT(w io.WriteCloser) EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.DebugDOTOut = w
	}
}


func Colorize(scheme *ColorScheme) EncodeOptionFunc {
	return func(opt *EncodeOption) {
		opt.Flag |= encoder.ColorizeOption
		opt.ColorScheme = scheme
	}
}

type DecodeOption = decoder.Option
type DecodeOptionFunc func(*DecodeOption)







func DecodeFieldPriorityFirstWin() DecodeOptionFunc {
	return func(opt *DecodeOption) {
		opt.Flags |= decoder.FirstWinOption
	}
}
