









package writeconcern 

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const majority = "majority"




var ErrInconsistent = errors.New("a write concern cannot have both w=0 and j=true")




var ErrEmptyWriteConcern = errors.New("a write concern must have at least one field set")




var ErrNegativeW = errors.New("write concern `w` field cannot be a negative number")




var ErrNegativeWTimeout = errors.New("write concern `wtimeout` field cannot be negative")







type WriteConcern struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	W interface{}

	
	
	
	
	
	
	Journal *bool

	
	
	
	
	
	
	
	
	
	WTimeout time.Duration
}






func Unacknowledged() *WriteConcern {
	return &WriteConcern{W: 0}
}







func W1() *WriteConcern {
	return &WriteConcern{W: 1}
}









func Journaled() *WriteConcern {
	journal := true
	return &WriteConcern{Journal: &journal}
}












func Majority() *WriteConcern {
	return &WriteConcern{W: majority}
}







func Custom(tag string) *WriteConcern {
	return &WriteConcern{W: tag}
}















type Option func(concern *WriteConcern)















func New(options ...Option) *WriteConcern {
	concern := &WriteConcern{}

	for _, option := range options {
		option(concern)
	}

	return concern
}
















func W(w int) Option {
	return func(concern *WriteConcern) {
		concern.W = w
	}
}





func WMajority() Option {
	return func(concern *WriteConcern) {
		concern.W = majority
	}
}





func WTagSet(tag string) Option {
	return func(concern *WriteConcern) {
		concern.W = tag
	}
}
















func J(j bool) Option {
	return func(concern *WriteConcern) {
		
		
		
		
		if j {
			concern.Journal = &j
		}
	}
}



















func WTimeout(d time.Duration) Option {
	return func(concern *WriteConcern) {
		concern.WTimeout = d
	}
}





func (wc *WriteConcern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if wc == nil {
		return 0, nil, ErrEmptyWriteConcern
	}

	var elems []byte
	if wc.W != nil {
		
		
		switch w := wc.W.(type) {
		case int:
			if w < 0 {
				return 0, nil, ErrNegativeW
			}

			
			
			if wc.Journal != nil && *wc.Journal && w == 0 {
				return 0, nil, ErrInconsistent
			}

			elems = bsoncore.AppendInt32Element(elems, "w", int32(w))
		case string:
			elems = bsoncore.AppendStringElement(elems, "w", w)
		default:
			return 0,
				nil,
				fmt.Errorf("WriteConcern.W must be a string or int, but is a %T", wc.W)
		}
	}

	if wc.Journal != nil {
		elems = bsoncore.AppendBooleanElement(elems, "j", *wc.Journal)
	}

	if wc.WTimeout < 0 {
		return 0, nil, ErrNegativeWTimeout
	}

	if wc.WTimeout != 0 {
		elems = bsoncore.AppendInt64Element(elems, "wtimeout", int64(wc.WTimeout/time.Millisecond))
	}

	if len(elems) == 0 {
		return 0, nil, ErrEmptyWriteConcern
	}
	return bson.TypeEmbeddedDocument, bsoncore.BuildDocument(nil, elems), nil
}





func AcknowledgedValue(rawv bson.RawValue) bool {
	doc, ok := bsoncore.Value{Type: rawv.Type, Data: rawv.Value}.DocumentOK()
	if !ok {
		return false
	}

	val, err := doc.LookupErr("w")
	if err != nil {
		
		return true
	}

	i32, ok := val.Int32OK()
	if !ok {
		return false
	}
	return i32 != 0
}


func (wc *WriteConcern) Acknowledged() bool {
	
	
	return wc == nil || wc.W != 0 || (wc.Journal != nil && *wc.Journal)
}


func (wc *WriteConcern) IsValid() bool {
	if wc == nil {
		return true
	}

	switch w := wc.W.(type) {
	case int:
		
		
		return w >= 0 && (w > 0 || wc.Journal == nil || !*wc.Journal)
	case string, nil:
		
		return true
	default:
		
		return false
	}
}




func (wc *WriteConcern) GetW() interface{} {
	return wc.W
}




func (wc *WriteConcern) GetJ() bool {
	
	
	
	return wc.Journal != nil && *wc.Journal
}




func (wc *WriteConcern) GetWTimeout() time.Duration {
	return wc.WTimeout
}















func (wc *WriteConcern) WithOptions(options ...Option) *WriteConcern {
	if wc == nil {
		return New(options...)
	}
	newWC := &WriteConcern{}
	*newWC = *wc

	for _, option := range options {
		option(newWC)
	}

	return newWC
}




func AckWrite(wc *WriteConcern) bool {
	return wc == nil || wc.Acknowledged()
}
