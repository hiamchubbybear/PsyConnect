

package encoder

import (
	"encoding/json"
	"io"

	"github.com/bytedance/sonic/internal/encoder/vars"
)


type StreamEncoder struct {
    w io.Writer
    Encoder
}




func NewStreamEncoder(w io.Writer) *StreamEncoder {
    return &StreamEncoder{w: w}
}


func (enc *StreamEncoder) Encode(val interface{}) (err error) {
    out := vars.NewBytes()

    
    err = EncodeInto(out, val, enc.Opts)
    if err != nil {
        goto free_bytes
    }

    if enc.indent != "" || enc.prefix != "" {
        
        buf := vars.NewBuffer()
        err = json.Indent(buf, *out, enc.prefix, enc.indent)
        if err != nil {
            vars.FreeBuffer(buf)
            goto free_bytes
        }

        
        if enc.Opts & NoEncoderNewline == 0 {
            buf.WriteByte('\n')
        }

        
        _, err = io.Copy(enc.w, buf)
        if err != nil {
            vars.FreeBuffer(buf)
            goto free_bytes
        }

    } else {
        
        var n int
        buf := *out
        for len(buf) > 0 {
            n, err = enc.w.Write(buf)
            buf = buf[n:]
            if err != nil {
                goto free_bytes
            }
        }

        
        if enc.Opts & NoEncoderNewline == 0 {
            enc.w.Write([]byte{'\n'})
        }
    }

free_bytes:
    vars.FreeBytes(out)
    return err
}
