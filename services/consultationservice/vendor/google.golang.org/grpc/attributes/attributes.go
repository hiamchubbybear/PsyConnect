








package attributes

import (
	"fmt"
	"strings"
)







type Attributes struct {
	m map[any]any
}


func New(key, value any) *Attributes {
	return &Attributes{m: map[any]any{key: value}}
}





func (a *Attributes) WithValue(key, value any) *Attributes {
	if a == nil {
		return New(key, value)
	}
	n := &Attributes{m: make(map[any]any, len(a.m)+1)}
	for k, v := range a.m {
		n.m[k] = v
	}
	n.m[key] = value
	return n
}



func (a *Attributes) Value(key any) any {
	if a == nil {
		return nil
	}
	return a.m[key]
}







func (a *Attributes) Equal(o *Attributes) bool {
	if a == nil && o == nil {
		return true
	}
	if a == nil || o == nil {
		return false
	}
	if len(a.m) != len(o.m) {
		return false
	}
	for k, v := range a.m {
		ov, ok := o.m[k]
		if !ok {
			
			return false
		}
		if eq, ok := v.(interface{ Equal(o any) bool }); ok {
			if !eq.Equal(ov) {
				return false
			}
		} else if v != ov {
			
			return false
		}
	}
	return true
}



func (a *Attributes) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for k, v := range a.m {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q: %q ", str(k), str(v)))
		first = false
	}
	sb.WriteString("}")
	return sb.String()
}

func str(x any) (s string) {
	if v, ok := x.(fmt.Stringer); ok {
		return fmt.Sprint(v)
	} else if v, ok := x.(string); ok {
		return v
	}
	return fmt.Sprintf("<%p>", x)
}







func (a *Attributes) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}
