


package codes 

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	
	Unset Code = 0

	
	
	
	
	
	Error Code = 1

	
	
	
	
	
	
	Ok Code = 2

	maxCode = 3
)


type Code uint32

var codeToStr = map[Code]string{
	Unset: "Unset",
	Error: "Error",
	Ok:    "Ok",
}

var strToCode = map[string]Code{
	`"Unset"`: Unset,
	`"Error"`: Error,
	`"Ok"`:    Ok,
}


func (c Code) String() string {
	return codeToStr[c]
}





func (c *Code) UnmarshalJSON(b []byte) error {
	
	
	
	if string(b) == "null" {
		return nil
	}
	if c == nil {
		return errors.New("nil receiver passed to UnmarshalJSON")
	}

	var x interface{}
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	switch x.(type) {
	case string:
		if jc, ok := strToCode[string(b)]; ok {
			*c = jc
			return nil
		}
		return fmt.Errorf("invalid code: %q", string(b))
	case float64:
		if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
			if ci >= maxCode {
				return fmt.Errorf("invalid code: %q", ci)
			}

			*c = Code(ci) 
			return nil
		}
		return fmt.Errorf("invalid code: %q", string(b))
	default:
		return fmt.Errorf("invalid code: %q", string(b))
	}
}


func (c *Code) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	str, ok := codeToStr[*c]
	if !ok {
		return nil, fmt.Errorf("invalid code: %d", *c)
	}
	return []byte(fmt.Sprintf("%q", str)), nil
}
