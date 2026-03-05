//go:build (amd64 && go1.17 && !go1.25) || (arm64 && go1.20 && !go1.25)
// +build amd64,go1.17,!go1.25 arm64,go1.20,!go1.25




package decoder

import (
    `github.com/bytedance/sonic/internal/decoder/api`
)


type Decoder = api.Decoder


type SyntaxError = api.SyntaxError


type MismatchTypeError = api.MismatchTypeError


type Options = api.Options

const (
    OptionUseInt64         Options = api.OptionUseInt64
    OptionUseNumber        Options = api.OptionUseNumber
    OptionUseUnicodeErrors Options = api.OptionUseUnicodeErrors
    OptionDisableUnknown   Options = api.OptionDisableUnknown
    OptionCopyString       Options = api.OptionCopyString
    OptionValidateString   Options = api.OptionValidateString
    OptionNoValidateJSON   Options = api.OptionNoValidateJSON
    OptionCaseSensitive    Options = api.OptionCaseSensitive
)


type StreamDecoder = api.StreamDecoder

var (
    
    NewDecoder = api.NewDecoder

    
    
    
    NewStreamDecoder = api.NewStreamDecoder

    
    
    
    
    
    Pretouch = api.Pretouch
    
    
    
    Skip = api.Skip
)
