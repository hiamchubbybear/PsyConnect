


package trace 

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	maxListMembers = 32

	listDelimiters  = ","
	memberDelimiter = "="

	errInvalidKey    errorConst = "invalid tracestate key"
	errInvalidValue  errorConst = "invalid tracestate value"
	errInvalidMember errorConst = "invalid tracestate list-member"
	errMemberNumber  errorConst = "too many list-members in tracestate"
	errDuplicate     errorConst = "duplicate list-member in tracestate"
)

type member struct {
	Key   string
	Value string
}



func checkValueChar(v byte) bool {
	return v >= '\x20' && v <= '\x7e' && v != '\x2c' && v != '\x3d'
}


func checkValueLast(v byte) bool {
	return v >= '\x21' && v <= '\x7e' && v != '\x2c' && v != '\x3d'
}








func checkValue(val string) bool {
	n := len(val)
	if n == 0 || n > 256 {
		return false
	}
	for i := 0; i < n-1; i++ {
		if !checkValueChar(val[i]) {
			return false
		}
	}
	return checkValueLast(val[n-1])
}

func checkKeyRemain(key string) bool {
	
	for _, v := range key {
		if isAlphaNum(byte(v)) {
			continue
		}
		switch v {
		case '_', '-', '*', '/':
			continue
		}
		return false
	}
	return true
}







func checkKeyPart(key string, n int) bool {
	if len(key) == 0 {
		return false
	}
	first := key[0] 
	ret := len(key[1:]) <= n
	ret = ret && first >= 'a' && first <= 'z'
	return ret && checkKeyRemain(key[1:])
}

func isAlphaNum(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	return c >= '0' && c <= '9'
}






func checkKeyTenant(key string, n int) bool {
	if len(key) == 0 {
		return false
	}
	return isAlphaNum(key[0]) && len(key[1:]) <= n && checkKeyRemain(key[1:])
}











func checkKey(key string) bool {
	tenant, system, ok := strings.Cut(key, "@")
	if !ok {
		return checkKeyPart(key, 255)
	}
	return checkKeyTenant(tenant, 240) && checkKeyPart(system, 13)
}

func newMember(key, value string) (member, error) {
	if !checkKey(key) {
		return member{}, errInvalidKey
	}
	if !checkValue(value) {
		return member{}, errInvalidValue
	}
	return member{Key: key, Value: value}, nil
}

func parseMember(m string) (member, error) {
	key, val, ok := strings.Cut(m, memberDelimiter)
	if !ok {
		return member{}, fmt.Errorf("%w: %s", errInvalidMember, m)
	}
	key = strings.TrimLeft(key, " \t")
	val = strings.TrimRight(val, " \t")
	result, e := newMember(key, val)
	if e != nil {
		return member{}, fmt.Errorf("%w: %s", errInvalidMember, m)
	}
	return result, nil
}



func (m member) String() string {
	return m.Key + "=" + m.Value
}












type TraceState struct { 
	
	list []member
}

var _ json.Marshaler = TraceState{}




func ParseTraceState(ts string) (TraceState, error) {
	if ts == "" {
		return TraceState{}, nil
	}

	wrapErr := func(err error) error {
		return fmt.Errorf("failed to parse tracestate: %w", err)
	}

	var members []member
	found := make(map[string]struct{})
	for ts != "" {
		var memberStr string
		memberStr, ts, _ = strings.Cut(ts, listDelimiters)
		if len(memberStr) == 0 {
			continue
		}

		m, err := parseMember(memberStr)
		if err != nil {
			return TraceState{}, wrapErr(err)
		}

		if _, ok := found[m.Key]; ok {
			return TraceState{}, wrapErr(errDuplicate)
		}
		found[m.Key] = struct{}{}

		members = append(members, m)
		if n := len(members); n > maxListMembers {
			return TraceState{}, wrapErr(errMemberNumber)
		}
	}

	return TraceState{list: members}, nil
}


func (ts TraceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.String())
}




func (ts TraceState) String() string {
	if len(ts.list) == 0 {
		return ""
	}
	var n int
	n += len(ts.list)     
	n += len(ts.list) - 1 
	for _, mem := range ts.list {
		n += len(mem.Key)
		n += len(mem.Value)
	}

	var sb strings.Builder
	sb.Grow(n)
	_, _ = sb.WriteString(ts.list[0].Key)
	_ = sb.WriteByte('=')
	_, _ = sb.WriteString(ts.list[0].Value)
	for i := 1; i < len(ts.list); i++ {
		_ = sb.WriteByte(listDelimiters[0])
		_, _ = sb.WriteString(ts.list[i].Key)
		_ = sb.WriteByte('=')
		_, _ = sb.WriteString(ts.list[i].Value)
	}
	return sb.String()
}



func (ts TraceState) Get(key string) string {
	for _, member := range ts.list {
		if member.Key == key {
			return member.Value
		}
	}

	return ""
}



func (ts TraceState) Walk(f func(key, value string) bool) {
	for _, m := range ts.list {
		if !f(m.Key, m.Value) {
			break
		}
	}
}













func (ts TraceState) Insert(key, value string) (TraceState, error) {
	m, err := newMember(key, value)
	if err != nil {
		return ts, err
	}
	n := len(ts.list)
	found := n
	for i := range ts.list {
		if ts.list[i].Key == key {
			found = i
		}
	}
	cTS := TraceState{}
	if found == n && n < maxListMembers {
		cTS.list = make([]member, n+1)
	} else {
		cTS.list = make([]member, n)
	}
	cTS.list[0] = m
	
	copy(cTS.list[1:], ts.list[0:found])
	if found < n {
		copy(cTS.list[1+found:], ts.list[found+1:])
	}
	return cTS, nil
}



func (ts TraceState) Delete(key string) TraceState {
	members := make([]member, ts.Len())
	copy(members, ts.list)
	for i, member := range ts.list {
		if member.Key == key {
			members = append(members[:i], members[i+1:]...)
			
			break
		}
	}
	return TraceState{list: members}
}


func (ts TraceState) Len() int {
	return len(ts.list)
}
