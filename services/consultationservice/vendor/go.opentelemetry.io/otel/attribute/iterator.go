


package attribute 



type Iterator struct {
	storage *Set
	idx     int
}




type MergeIterator struct {
	one     oneIterator
	two     oneIterator
	current KeyValue
}

type oneIterator struct {
	iter Iterator
	done bool
	attr KeyValue
}



func (i *Iterator) Next() bool {
	i.idx++
	return i.idx < i.Len()
}





func (i *Iterator) Label() KeyValue {
	return i.Attribute()
}



func (i *Iterator) Attribute() KeyValue {
	kv, _ := i.storage.Get(i.idx)
	return kv
}





func (i *Iterator) IndexedLabel() (int, KeyValue) {
	return i.idx, i.Attribute()
}



func (i *Iterator) IndexedAttribute() (int, KeyValue) {
	return i.idx, i.Attribute()
}


func (i *Iterator) Len() int {
	return i.storage.Len()
}




func (i *Iterator) ToSlice() []KeyValue {
	l := i.Len()
	if l == 0 {
		return nil
	}
	i.idx = -1
	slice := make([]KeyValue, 0, l)
	for i.Next() {
		slice = append(slice, i.Attribute())
	}
	return slice
}



func NewMergeIterator(s1, s2 *Set) MergeIterator {
	mi := MergeIterator{
		one: makeOne(s1.Iter()),
		two: makeOne(s2.Iter()),
	}
	return mi
}

func makeOne(iter Iterator) oneIterator {
	oi := oneIterator{
		iter: iter,
	}
	oi.advance()
	return oi
}

func (oi *oneIterator) advance() {
	if oi.done = !oi.iter.Next(); !oi.done {
		oi.attr = oi.iter.Attribute()
	}
}


func (m *MergeIterator) Next() bool {
	if m.one.done && m.two.done {
		return false
	}
	if m.one.done {
		m.current = m.two.attr
		m.two.advance()
		return true
	}
	if m.two.done {
		m.current = m.one.attr
		m.one.advance()
		return true
	}
	if m.one.attr.Key == m.two.attr.Key {
		m.current = m.one.attr 
		m.one.advance()
		m.two.advance()
		return true
	}
	if m.one.attr.Key < m.two.attr.Key {
		m.current = m.one.attr
		m.one.advance()
		return true
	}
	m.current = m.two.attr
	m.two.advance()
	return true
}




func (m *MergeIterator) Label() KeyValue {
	return m.current
}


func (m *MergeIterator) Attribute() KeyValue {
	return m.current
}
