


package attribute 

import (
	"bytes"
	"sync"
	"sync/atomic"
)

type (
	
	
	
	
	Encoder interface {
		
		
		Encode(iterator Iterator) string

		
		
		ID() EncoderID
	}

	
	
	EncoderID struct {
		value uint64
	}

	
	
	
	
	defaultAttrEncoder struct {
		
		
		
		pool sync.Pool 
	}
)






const escapeChar = '\\'

var (
	_ Encoder = &defaultAttrEncoder{}

	
	encoderIDCounter uint64

	defaultEncoderOnce     sync.Once
	defaultEncoderID       = NewEncoderID()
	defaultEncoderInstance *defaultAttrEncoder
)




func NewEncoderID() EncoderID {
	return EncoderID{value: atomic.AddUint64(&encoderIDCounter, 1)}
}








func DefaultEncoder() Encoder {
	defaultEncoderOnce.Do(func() {
		defaultEncoderInstance = &defaultAttrEncoder{
			pool: sync.Pool{
				New: func() interface{} {
					return &bytes.Buffer{}
				},
			},
		}
	})
	return defaultEncoderInstance
}


func (d *defaultAttrEncoder) Encode(iter Iterator) string {
	buf := d.pool.Get().(*bytes.Buffer)
	defer d.pool.Put(buf)
	buf.Reset()

	for iter.Next() {
		i, keyValue := iter.IndexedAttribute()
		if i > 0 {
			_, _ = buf.WriteRune(',')
		}
		copyAndEscape(buf, string(keyValue.Key))

		_, _ = buf.WriteRune('=')

		if keyValue.Value.Type() == STRING {
			copyAndEscape(buf, keyValue.Value.AsString())
		} else {
			_, _ = buf.WriteString(keyValue.Value.Emit())
		}
	}
	return buf.String()
}


func (*defaultAttrEncoder) ID() EncoderID {
	return defaultEncoderID
}



func copyAndEscape(buf *bytes.Buffer, val string) {
	for _, ch := range val {
		switch ch {
		case '=', ',', escapeChar:
			_, _ = buf.WriteRune(escapeChar)
		}
		_, _ = buf.WriteRune(ch)
	}
}



func (id EncoderID) Valid() bool {
	return id.value != 0
}
