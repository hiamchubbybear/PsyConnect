// +build amd64,go1.17,!go1.25 arm64,go1.20,!go1.25



package encoder

import (
    `github.com/bytedance/sonic/internal/encoder`
)


const EnableFallback = false


type Encoder = encoder.Encoder


type StreamEncoder = encoder.StreamEncoder


type Options = encoder.Options

const (
    
    
    
    SortMapKeys Options = encoder.SortMapKeys

    
    
    
    EscapeHTML Options = encoder.EscapeHTML

    
    
    CompactMarshaler Options = encoder.CompactMarshaler

    
    
    NoQuoteTextMarshaler Options = encoder.NoQuoteTextMarshaler

    
    
    NoNullSliceOrMap Options = encoder.NoNullSliceOrMap

    
    
    ValidateString Options = encoder.ValidateString

    
    
    NoValidateJSONMarshaler Options = encoder.NoValidateJSONMarshaler

    
    NoEncoderNewline Options = encoder.NoEncoderNewline

    
    CompatibleWithStd Options = encoder.CompatibleWithStd

    
    EncodeNullForInfOrNan Options = encoder.EncodeNullForInfOrNan
)


var (
    
    Encode = encoder.Encode

    
    EncodeIndented = encoder.EncodeIndented

    
    
    
    EncodeInto = encoder.EncodeInto

    
    
    
    
    
    
    HTMLEscape = encoder.HTMLEscape

    
    
    
    
    
    Pretouch = encoder.Pretouch

    
    Quote = encoder.Quote

    
    
    
    
    
    Valid = encoder.Valid

    
    
    
    NewStreamEncoder = encoder.NewStreamEncoder
)
