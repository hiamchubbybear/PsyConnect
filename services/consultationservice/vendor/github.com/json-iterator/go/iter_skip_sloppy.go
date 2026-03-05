

package jsoniter



func (iter *Iterator) skipNumber() {
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case ' ', '\n', '\r', '\t', ',', '}', ']':
				iter.head = i
				return
			}
		}
		if !iter.loadMore() {
			return
		}
	}
}

func (iter *Iterator) skipArray() {
	level := 1
	if !iter.incrementDepth() {
		return
	}
	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': 
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 
			case '[': 
				level++
				if !iter.incrementDepth() {
					return
				}
			case ']': 
				level--
				if !iter.decrementDepth() {
					return
				}

				
				if level == 0 {
					iter.head = i + 1
					return
				}
			}
		}
		if !iter.loadMore() {
			iter.ReportError("skipObject", "incomplete array")
			return
		}
	}
}

func (iter *Iterator) skipObject() {
	level := 1
	if !iter.incrementDepth() {
		return
	}

	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': 
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 
			case '{': 
				level++
				if !iter.incrementDepth() {
					return
				}
			case '}': 
				level--
				if !iter.decrementDepth() {
					return
				}

				
				if level == 0 {
					iter.head = i + 1
					return
				}
			}
		}
		if !iter.loadMore() {
			iter.ReportError("skipObject", "incomplete object")
			return
		}
	}
}

func (iter *Iterator) skipString() {
	for {
		end, escaped := iter.findStringEnd()
		if end == -1 {
			if !iter.loadMore() {
				iter.ReportError("skipString", "incomplete string")
				return
			}
			if escaped {
				iter.head = 1 
			}
		} else {
			iter.head = end
			return
		}
	}
}




func (iter *Iterator) findStringEnd() (int, bool) {
	escaped := false
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		if c == '"' {
			if !escaped {
				return i + 1, false
			}
			j := i - 1
			for {
				if j < iter.head || iter.buf[j] != '\\' {
					
					
					return i + 1, true
				}
				j--
				if j < iter.head || iter.buf[j] != '\\' {
					
					
					break
				}
				j--
			}
		} else if c == '\\' {
			escaped = true
		}
	}
	j := iter.tail - 1
	for {
		if j < iter.head || iter.buf[j] != '\\' {
			
			
			return -1, false 
		}
		j--
		if j < iter.head || iter.buf[j] != '\\' {
			
			
			break
		}
		j--

	}
	return -1, true 
}
