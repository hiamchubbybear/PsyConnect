//go:build (amd64 && go1.17 && !go1.25) || (arm64 && go1.20 && !go1.25)
// +build amd64,go1.17,!go1.25 arm64,go1.20,!go1.25



//go:generate make
package sonic

import (
    `io`
    `reflect`

    `github.com/bytedance/sonic/decoder`
    `github.com/bytedance/sonic/encoder`
    `github.com/bytedance/sonic/option`
    `github.com/bytedance/sonic/internal/rt`
)

const apiKind = UseSonicJSON

type frozenConfig struct {
    Config
    encoderOpts encoder.Options
    decoderOpts decoder.Options
}


func (cfg Config) Froze() API {
    api := &frozenConfig{Config: cfg}

    
    if cfg.EscapeHTML {
        api.encoderOpts |= encoder.EscapeHTML
    }
    if cfg.SortMapKeys {
        api.encoderOpts |= encoder.SortMapKeys
    }
    if cfg.CompactMarshaler {
        api.encoderOpts |= encoder.CompactMarshaler
    }
    if cfg.NoQuoteTextMarshaler {
        api.encoderOpts |= encoder.NoQuoteTextMarshaler
    }
    if cfg.NoNullSliceOrMap {
        api.encoderOpts |= encoder.NoNullSliceOrMap
    }
    if cfg.ValidateString {
        api.encoderOpts |= encoder.ValidateString
    }
    if cfg.NoValidateJSONMarshaler {
        api.encoderOpts |= encoder.NoValidateJSONMarshaler
    }
    if cfg.NoEncoderNewline {
        api.encoderOpts |= encoder.NoEncoderNewline
    }
    if cfg.EncodeNullForInfOrNan {
        api.encoderOpts |= encoder.EncodeNullForInfOrNan
    }

    
    if cfg.NoValidateJSONSkip {
        api.decoderOpts |= decoder.OptionNoValidateJSON
    }
    if cfg.UseInt64 {
        api.decoderOpts |= decoder.OptionUseInt64
    }
    if cfg.UseNumber {
        api.decoderOpts |= decoder.OptionUseNumber
    }
    if cfg.DisallowUnknownFields {
        api.decoderOpts |= decoder.OptionDisableUnknown
    }
    if cfg.CopyString {
        api.decoderOpts |= decoder.OptionCopyString
    }
    if cfg.ValidateString {
        api.decoderOpts |= decoder.OptionValidateString
    }
    return api
}


func (cfg frozenConfig) Marshal(val interface{}) ([]byte, error) {
    return encoder.Encode(val, cfg.encoderOpts)
}


func (cfg frozenConfig) MarshalToString(val interface{}) (string, error) {
    buf, err := encoder.Encode(val, cfg.encoderOpts)
    return rt.Mem2Str(buf), err
}


func (cfg frozenConfig) MarshalIndent(val interface{}, prefix, indent string) ([]byte, error) {
    return encoder.EncodeIndented(val, prefix, indent, cfg.encoderOpts)
}


func (cfg frozenConfig) UnmarshalFromString(buf string, val interface{}) error {
    dec := decoder.NewDecoder(buf)
    dec.SetOptions(cfg.decoderOpts)
    err := dec.Decode(val)

    
    if err != nil {
        return err
    }

    return dec.CheckTrailings()
}


func (cfg frozenConfig) Unmarshal(buf []byte, val interface{}) error {
    return cfg.UnmarshalFromString(string(buf), val)
}


func (cfg frozenConfig) NewEncoder(writer io.Writer) Encoder {
    enc := encoder.NewStreamEncoder(writer)
    enc.Opts = cfg.encoderOpts
    return enc
}


func (cfg frozenConfig) NewDecoder(reader io.Reader) Decoder {
    dec := decoder.NewStreamDecoder(reader)
    dec.SetOptions(cfg.decoderOpts)
    return dec
}


func (cfg frozenConfig) Valid(data []byte) bool {
    ok, _ := encoder.Valid(data)
    return ok
}






func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
    if err := encoder.Pretouch(vt, opts...); err != nil {
        return err
    } 
    if err := decoder.Pretouch(vt, opts...); err != nil {
        return err
    }
    
    if vt.Kind() == reflect.Ptr {
        vt = vt.Elem()
    } else {
        vt = reflect.PtrTo(vt)
    }
    if err := encoder.Pretouch(vt, opts...); err != nil {
        return err
    } 
    if err := decoder.Pretouch(vt, opts...); err != nil {
        return err
    }
    return nil
}
