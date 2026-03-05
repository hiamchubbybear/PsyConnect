

package base64x

import (
    `encoding/base64`

    "github.com/cloudwego/base64x/internal/native"
)






type Encoding int

const (
    _MODE_URL  = 1 << 0
    _MODE_RAW  = 1 << 1
    _MODE_AVX2 = 1 << 2
    _MODE_JSON = 1 << 3
)



const StdEncoding Encoding = 0



const URLEncoding Encoding = _MODE_URL





const RawStdEncoding Encoding = _MODE_RAW





const RawURLEncoding Encoding = _MODE_RAW | _MODE_URL


const JSONStdEncoding Encoding = _MODE_JSON;

var (
    archFlags = 0
)












func (self Encoding) Encode(out []byte, src []byte) {
    if len(src) != 0 {
        if buf := out[:0:len(out)]; self.EncodedLen(len(src)) <= len(out) {
            self.EncodeUnsafe(&buf, src)
        } else {
            panic("encoder output buffer is too small")
        }
    }
}





func (self Encoding) EncodeUnsafe(out *[]byte, src []byte) {
    native.B64Encode(out, &src, int(self) | archFlags)
}


func (self Encoding) EncodeToString(src []byte) string {
    nbs := len(src)
    ret := make([]byte, 0, self.EncodedLen(nbs))

    
    self.EncodeUnsafe(&ret, src)
    return mem2str(ret)
}



func (self Encoding) EncodedLen(n int) int {
    if (self & _MODE_RAW) == 0 {
        return (n + 2) / 3 * 4
    } else {
        return (n * 8 + 5) / 6
    }
}












func (self Encoding) Decode(out []byte, src []byte) (int, error) {
    if len(src) == 0 {
        return 0, nil
    } else if buf := out[:0:len(out)]; self.DecodedLen(len(src)) <= len(out) {
        return self.DecodeUnsafe(&buf, src)
    } else {
        panic("decoder output buffer is too small")
    }
}





func (self Encoding) DecodeUnsafe(out *[]byte, src []byte) (int, error) {
    if n := native.B64Decode(out, mem2addr(src), len(src), int(self) | archFlags); n >= 0 {
        return n, nil
    } else {
        return 0, base64.CorruptInputError(-n - 1)
    }
}


func (self Encoding) DecodeString(s string) ([]byte, error) {
    src := str2mem(s)
    ret := make([]byte, 0, self.DecodedLen(len(s)))

    
    if _, err := self.DecodeUnsafe(&ret, src); err != nil {
        return nil, err
    } else {
        return ret, nil
    }
}



func (self Encoding) DecodedLen(n int) int {
    if (self & _MODE_RAW) == 0 {
        return n / 4 * 3
    } else {
        return n * 6 / 8
    }
}
