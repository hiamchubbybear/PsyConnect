







package primitive 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)


type Binary struct {
	Subtype byte
	Data    []byte
}


func (bp Binary) Equal(bp2 Binary) bool {
	if bp.Subtype != bp2.Subtype {
		return false
	}
	return bytes.Equal(bp.Data, bp2.Data)
}


func (bp Binary) IsZero() bool {
	return bp.Subtype == 0 && len(bp.Data) == 0
}


type Undefined struct{}


type DateTime int64

var _ json.Marshaler = DateTime(0)
var _ json.Unmarshaler = (*DateTime)(nil)


func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time().UTC())
}


func (d *DateTime) UnmarshalJSON(data []byte) error {
	
	
	
	if string(data) == "null" {
		return nil
	}

	var tempTime time.Time
	if err := json.Unmarshal(data, &tempTime); err != nil {
		return err
	}

	*d = NewDateTimeFromTime(tempTime)
	return nil
}


func (d DateTime) Time() time.Time {
	return time.Unix(int64(d)/1000, int64(d)%1000*1000000)
}


func NewDateTimeFromTime(t time.Time) DateTime {
	return DateTime(t.Unix()*1e3 + int64(t.Nanosecond())/1e6)
}


type Null struct{}


type Regex struct {
	Pattern string
	Options string
}

func (rp Regex) String() string {
	return fmt.Sprintf(`{"pattern": "%s", "options": "%s"}`, rp.Pattern, rp.Options)
}


func (rp Regex) Equal(rp2 Regex) bool {
	return rp.Pattern == rp2.Pattern && rp.Options == rp2.Options
}


func (rp Regex) IsZero() bool {
	return rp.Pattern == "" && rp.Options == ""
}


type DBPointer struct {
	DB      string
	Pointer ObjectID
}

func (d DBPointer) String() string {
	return fmt.Sprintf(`{"db": "%s", "pointer": "%s"}`, d.DB, d.Pointer)
}


func (d DBPointer) Equal(d2 DBPointer) bool {
	return d == d2
}


func (d DBPointer) IsZero() bool {
	return d.DB == "" && d.Pointer.IsZero()
}


type JavaScript string


type Symbol string


type CodeWithScope struct {
	Code  JavaScript
	Scope interface{}
}

func (cws CodeWithScope) String() string {
	return fmt.Sprintf(`{"code": "%s", "scope": %v}`, cws.Code, cws.Scope)
}


type Timestamp struct {
	T uint32
	I uint32
}


func (tp Timestamp) After(tp2 Timestamp) bool {
	return tp.T > tp2.T || (tp.T == tp2.T && tp.I > tp2.I)
}


func (tp Timestamp) Before(tp2 Timestamp) bool {
	return tp.T < tp2.T || (tp.T == tp2.T && tp.I < tp2.I)
}


func (tp Timestamp) Equal(tp2 Timestamp) bool {
	return tp.T == tp2.T && tp.I == tp2.I
}


func (tp Timestamp) IsZero() bool {
	return tp.T == 0 && tp.I == 0
}



func (tp Timestamp) Compare(tp2 Timestamp) int {
	switch {
	case tp.Equal(tp2):
		return 0
	case tp.Before(tp2):
		return -1
	default:
		return +1
	}
}





func CompareTimestamp(tp, tp2 Timestamp) int {
	return tp.Compare(tp2)
}


type MinKey struct{}


type MaxKey struct{}







type D []E





func (d D) Map() M {
	m := make(M, len(d))
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}


type E struct {
	Key   string
	Value interface{}
}








type M map[string]interface{}






type A []interface{}
