


//go:build go1.7 && !go1.12 && (safe || codec.safe || appengine)
// +build go1.7
// +build !go1.12
// +build safe codec.safe appengine

package codec

import "reflect"

type mapIter struct {
	m      reflect.Value
	keys   []reflect.Value
	j      int
	values bool
}

func (t *mapIter) Next() (r bool) {
	t.j++
	return t.j < len(t.keys)
}

func (t *mapIter) Key() reflect.Value {
	return t.keys[t.j]
}

func (t *mapIter) Value() (r reflect.Value) {
	if t.values {
		return t.m.MapIndex(t.keys[t.j])
	}
	return
}

func (t *mapIter) Done() {}

func mapRange(t *mapIter, m, k, v reflect.Value, values bool) {
	*t = mapIter{
		m:      m,
		keys:   m.MapKeys(),
		values: values,
		j:      -1,
	}
}
