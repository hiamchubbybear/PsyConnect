


package attribute 

import (
	"cmp"
	"encoding/json"
	"reflect"
	"slices"
	"sort"
)

type (
	
	
	
	
	
	
	
	
	
	
	
	
	
	Set struct {
		equivalent Distinct
	}

	
	
	
	
	
	Distinct struct {
		iface interface{}
	}

	
	
	
	
	
	Sortable []KeyValue
)

var (
	
	keyValueType = reflect.TypeOf(KeyValue{})

	
	emptySet = &Set{
		equivalent: Distinct{
			iface: [0]KeyValue{},
		},
	}
)




func EmptySet() *Set {
	return emptySet
}


func (d Distinct) reflectValue() reflect.Value {
	return reflect.ValueOf(d.iface)
}


func (d Distinct) Valid() bool {
	return d.iface != nil
}


func (l *Set) Len() int {
	if l == nil || !l.equivalent.Valid() {
		return 0
	}
	return l.equivalent.reflectValue().Len()
}


func (l *Set) Get(idx int) (KeyValue, bool) {
	if l == nil || !l.equivalent.Valid() {
		return KeyValue{}, false
	}
	value := l.equivalent.reflectValue()

	if idx >= 0 && idx < value.Len() {
		
		
		return value.Index(idx).Interface().(KeyValue), true
	}

	return KeyValue{}, false
}


func (l *Set) Value(k Key) (Value, bool) {
	if l == nil || !l.equivalent.Valid() {
		return Value{}, false
	}
	rValue := l.equivalent.reflectValue()
	vlen := rValue.Len()

	idx := sort.Search(vlen, func(idx int) bool {
		return rValue.Index(idx).Interface().(KeyValue).Key >= k
	})
	if idx >= vlen {
		return Value{}, false
	}
	keyValue := rValue.Index(idx).Interface().(KeyValue)
	if k == keyValue.Key {
		return keyValue.Value, true
	}
	return Value{}, false
}


func (l *Set) HasValue(k Key) bool {
	if l == nil {
		return false
	}
	_, ok := l.Value(k)
	return ok
}


func (l *Set) Iter() Iterator {
	return Iterator{
		storage: l,
		idx:     -1,
	}
}



func (l *Set) ToSlice() []KeyValue {
	iter := l.Iter()
	return iter.ToSlice()
}





func (l *Set) Equivalent() Distinct {
	if l == nil || !l.equivalent.Valid() {
		return emptySet.equivalent
	}
	return l.equivalent
}


func (l *Set) Equals(o *Set) bool {
	return l.Equivalent() == o.Equivalent()
}


func (l *Set) Encoded(encoder Encoder) string {
	if l == nil || encoder == nil {
		return ""
	}

	return encoder.Encode(l.Iter())
}

func empty() Set {
	return Set{
		equivalent: emptySet.equivalent,
	}
}






func NewSet(kvs ...KeyValue) Set {
	s, _ := NewSetWithFiltered(kvs, nil)
	return s
}







func NewSetWithSortable(kvs []KeyValue, _ *Sortable) Set {
	s, _ := NewSetWithFiltered(kvs, nil)
	return s
}






func NewSetWithFiltered(kvs []KeyValue, filter Filter) (Set, []KeyValue) {
	
	if len(kvs) == 0 {
		return empty(), nil
	}

	
	
	slices.SortStableFunc(kvs, func(a, b KeyValue) int {
		return cmp.Compare(a.Key, b.Key)
	})

	position := len(kvs) - 1
	offset := position - 1

	
	
	
	
	
	
	for ; offset >= 0; offset-- {
		if kvs[offset].Key == kvs[position].Key {
			continue
		}
		position--
		kvs[offset], kvs[position] = kvs[position], kvs[offset]
	}
	kvs = kvs[position:]

	if filter != nil {
		if div := filteredToFront(kvs, filter); div != 0 {
			return Set{equivalent: computeDistinct(kvs[div:])}, kvs[:div]
		}
	}
	return Set{equivalent: computeDistinct(kvs)}, nil
}


























func NewSetWithSortableFiltered(kvs []KeyValue, _ *Sortable, filter Filter) (Set, []KeyValue) {
	return NewSetWithFiltered(kvs, filter)
}





func filteredToFront(slice []KeyValue, keep Filter) int {
	n := len(slice)
	j := n
	for i := n - 1; i >= 0; i-- {
		if keep(slice[i]) {
			j--
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	return j
}



func (l *Set) Filter(re Filter) (Set, []KeyValue) {
	if re == nil {
		return *l, nil
	}

	
	n := l.Len()
	first := n - 1
	for ; first >= 0; first-- {
		kv, _ := l.Get(first)
		if !re(kv) {
			break
		}
	}

	
	if first < 0 {
		return *l, nil
	}

	
	
	
	
	slice := l.ToSlice()

	
	if first == 0 {
		
		
		return Set{equivalent: computeDistinct(slice[1:])}, slice[:1]
	}

	
	kv := slice[first]
	copy(slice[1:first+1], slice[:first])
	slice[0] = kv

	
	div := filteredToFront(slice[1:first+1], re) + 1
	return Set{equivalent: computeDistinct(slice[div:])}, slice[:div]
}




func computeDistinct(kvs []KeyValue) Distinct {
	iface := computeDistinctFixed(kvs)
	if iface == nil {
		iface = computeDistinctReflect(kvs)
	}
	return Distinct{
		iface: iface,
	}
}



func computeDistinctFixed(kvs []KeyValue) interface{} {
	switch len(kvs) {
	case 1:
		return [1]KeyValue(kvs)
	case 2:
		return [2]KeyValue(kvs)
	case 3:
		return [3]KeyValue(kvs)
	case 4:
		return [4]KeyValue(kvs)
	case 5:
		return [5]KeyValue(kvs)
	case 6:
		return [6]KeyValue(kvs)
	case 7:
		return [7]KeyValue(kvs)
	case 8:
		return [8]KeyValue(kvs)
	case 9:
		return [9]KeyValue(kvs)
	case 10:
		return [10]KeyValue(kvs)
	default:
		return nil
	}
}



func computeDistinctReflect(kvs []KeyValue) interface{} {
	at := reflect.New(reflect.ArrayOf(len(kvs), keyValueType)).Elem()
	for i, keyValue := range kvs {
		*(at.Index(i).Addr().Interface().(*KeyValue)) = keyValue
	}
	return at.Interface()
}


func (l *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.equivalent.iface)
}


func (l Set) MarshalLog() interface{} {
	kvs := make(map[string]string)
	for _, kv := range l.ToSlice() {
		kvs[string(kv.Key)] = kv.Value.Emit()
	}
	return kvs
}


func (l *Sortable) Len() int {
	return len(*l)
}


func (l *Sortable) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}


func (l *Sortable) Less(i, j int) bool {
	return (*l)[i].Key < (*l)[j].Key
}
