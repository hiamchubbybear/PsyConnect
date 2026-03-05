package urn

import (
	"encoding/json"
	"fmt"
	"strings"
)

const errInvalidURN = "invalid URN: %s"








type URN struct {
	prefix     string 
	ID         string 
	SS         string 
	norm       string 
	kind       Kind
	scim       *SCIM
	rComponent string 
	qComponent string 
	fComponent string 
	rStart     bool   
	qStart     bool   
	tolower    []int
}




func (u *URN) Normalize() *URN {
	return &URN{
		prefix: "urn",
		ID:     strings.ToLower(u.ID),
		SS:     u.norm,
		
		
		
	}
}


func (u *URN) Equal(x *URN) bool {
	if x == nil {
		return false
	}
	nu := u.Normalize()
	nx := x.Normalize()

	return nu.prefix == nx.prefix && nu.ID == nx.ID && nu.SS == nx.SS
}







func (u *URN) String() string {
	var res string
	if u.ID != "" && u.SS != "" {
		if u.prefix == "" {
			res += "urn"
		}
		res += u.prefix + ":" + u.ID + ":" + u.SS
		if u.rComponent != "" {
			res += "?+" + u.rComponent
		}
		if u.qComponent != "" {
			res += "?=" + u.qComponent
		}
		if u.fComponent != "" {
			res += "#" + u.fComponent
		}
	}

	return res
}


func Parse(u []byte, options ...Option) (*URN, bool) {
	urn, err := NewMachine(options...).Parse(u)
	if err != nil {
		return nil, false
	}

	return urn, true
}


func (u URN) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}


func (u *URN) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	if value, ok := Parse([]byte(str)); !ok {
		return fmt.Errorf(errInvalidURN, str)
	} else {
		*u = *value
	}

	return nil
}

func (u *URN) IsSCIM() bool {
	return u.kind == RFC7643
}

func (u *URN) SCIM() *SCIM {
	if u.kind != RFC7643 {
		return nil
	}

	return u.scim
}

func (u *URN) RFC() Kind {
	return u.kind
}

func (u *URN) FComponent() string {
	return u.fComponent
}

func (u *URN) QComponent() string {
	return u.qComponent
}

func (u *URN) RComponent() string {
	return u.rComponent
}
