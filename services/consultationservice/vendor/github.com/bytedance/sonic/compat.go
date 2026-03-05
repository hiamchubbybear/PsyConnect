// +build !amd64,!arm64 go1.25 !go1.17 arm64,!go1.20



package sonic

import (
    `bytes`
    `encoding/json`
    `io`
    `reflect`

    `github.com/bytedance/sonic/option`
)

const apiKind = UseStdJSON

type frozenConfig struct {
    Config
}


func (cfg Config) Froze() API {
    api := &frozenConfig{Config: cfg}
    return api
}

func (cfg frozenConfig) marshalOptions(val interface{}, prefix, indent string) ([]byte, error) {
    w := bytes.NewBuffer([]byte{})
    enc := json.NewEncoder(w)
    enc.SetEscapeHTML(cfg.EscapeHTML)
    enc.SetIndent(prefix, indent)
    err := enc.Encode(val)
	out := w.Bytes()

	
	
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return out, err
}


func (cfg frozenConfig) Marshal(val interface{}) ([]byte, error) {
    if !cfg.EscapeHTML {
        return cfg.marshalOptions(val, "", "")
    }
    return json.Marshal(val)
}


func (cfg frozenConfig) MarshalToString(val interface{}) (string, error) {
    out, err := cfg.Marshal(val)
    return string(out), err
}


func (cfg frozenConfig) MarshalIndent(val interface{}, prefix, indent string) ([]byte, error) {
    if !cfg.EscapeHTML {
        return cfg.marshalOptions(val, prefix, indent)
    }
    return json.MarshalIndent(val, prefix, indent)
}


func (cfg frozenConfig) UnmarshalFromString(buf string, val interface{}) error {
    r := bytes.NewBufferString(buf)
    dec := json.NewDecoder(r)
    if cfg.UseNumber {
        dec.UseNumber()
    }
    if cfg.DisallowUnknownFields {
        dec.DisallowUnknownFields()
    }
    return dec.Decode(val)
}


func (cfg frozenConfig) Unmarshal(buf []byte, val interface{}) error {
    return cfg.UnmarshalFromString(string(buf), val)
}


func (cfg frozenConfig) NewEncoder(writer io.Writer) Encoder {
    enc := json.NewEncoder(writer)
    if !cfg.EscapeHTML {
        enc.SetEscapeHTML(cfg.EscapeHTML)
    }
    return enc
}


func (cfg frozenConfig) NewDecoder(reader io.Reader) Decoder {
    dec := json.NewDecoder(reader)
    if cfg.UseNumber {
        dec.UseNumber()
    }
    if cfg.DisallowUnknownFields {
        dec.DisallowUnknownFields()
    }
    return dec
}


func (cfg frozenConfig) Valid(data []byte) bool {
    return json.Valid(data)
}







func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
    return nil
}

