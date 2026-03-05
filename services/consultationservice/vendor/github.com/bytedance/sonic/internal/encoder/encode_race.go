//go:build race
// +build race



package encoder

import (
    `encoding/json`

    `github.com/bytedance/sonic/internal/rt`
)


func helpDetectDataRace(val interface{}) {
    var out []byte
    defer func() {
        if v := recover(); v != nil {
            
            println("panic when encoding on: ", truncate(out))
            panic(v)
        }
    }()
    out, _ = json.Marshal(val)
}

func encodeIntoCheckRace(buf *[]byte, val interface{}, opts Options) error {
	err := encodeInto(buf, val, opts)
    
    helpDetectDataRace(val)
    return err
}

func truncate(json []byte) string {
    if len(json) <= 256 {
        return rt.Mem2Str(json)
    } else {
        return rt.Mem2Str(json[len(json)-256:])
    }
}
