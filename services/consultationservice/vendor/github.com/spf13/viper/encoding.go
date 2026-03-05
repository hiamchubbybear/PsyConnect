package viper

import (
	"errors"
	"strings"
	"sync"

	"github.com/spf13/viper/internal/encoding/dotenv"
	"github.com/spf13/viper/internal/encoding/json"
	"github.com/spf13/viper/internal/encoding/toml"
	"github.com/spf13/viper/internal/encoding/yaml"
)



type Encoder interface {
	Encode(v map[string]any) ([]byte, error)
}



type Decoder interface {
	Decode(b []byte, v map[string]any) error
}


type Codec interface {
	Encoder
	Decoder
}








type EncoderRegistry interface {
	Encoder(format string) (Encoder, error)
}






type DecoderRegistry interface {
	Decoder(format string) (Decoder, error)
}


type CodecRegistry interface {
	EncoderRegistry
	DecoderRegistry
}


func WithEncoderRegistry(r EncoderRegistry) Option {
	return optionFunc(func(v *Viper) {
		if r == nil {
			return
		}

		v.encoderRegistry = r
	})
}


func WithDecoderRegistry(r DecoderRegistry) Option {
	return optionFunc(func(v *Viper) {
		if r == nil {
			return
		}

		v.decoderRegistry = r
	})
}


func WithCodecRegistry(r CodecRegistry) Option {
	return optionFunc(func(v *Viper) {
		if r == nil {
			return
		}

		v.encoderRegistry = r
		v.decoderRegistry = r
	})
}


type DefaultCodecRegistry struct {
	codecs map[string]Codec

	mu   sync.RWMutex
	once sync.Once
}


func NewCodecRegistry() *DefaultCodecRegistry {
	r := &DefaultCodecRegistry{}

	r.init()

	return r
}

func (r *DefaultCodecRegistry) init() {
	r.once.Do(func() {
		r.codecs = map[string]Codec{}
	})
}




func (r *DefaultCodecRegistry) RegisterCodec(format string, codec Codec) error {
	r.init()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.codecs[strings.ToLower(format)] = codec

	return nil
}




func (r *DefaultCodecRegistry) Encoder(format string) (Encoder, error) {
	encoder, ok := r.codec(format)
	if !ok {
		return nil, errors.New("encoder not found for this format")
	}

	return encoder, nil
}




func (r *DefaultCodecRegistry) Decoder(format string) (Decoder, error) {
	decoder, ok := r.codec(format)
	if !ok {
		return nil, errors.New("decoder not found for this format")
	}

	return decoder, nil
}

func (r *DefaultCodecRegistry) codec(format string) (Codec, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	format = strings.ToLower(format)

	if r.codecs != nil {
		codec, ok := r.codecs[format]
		if ok {
			return codec, true
		}
	}

	switch format {
	case "yaml", "yml":
		return yaml.Codec{}, true

	case "json":
		return json.Codec{}, true

	case "toml":
		return toml.Codec{}, true

	case "dotenv", "env":
		return &dotenv.Codec{}, true
	}

	return nil, false
}
