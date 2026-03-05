package jsoniter


func (iter *Iterator) ReadArray() (ret bool) {
	c := iter.nextToken()
	switch c {
	case 'n':
		iter.skipThreeBytes('u', 'l', 'l')
		return false 
	case '[':
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			return true
		}
		return false
	case ']':
		return false
	case ',':
		return true
	default:
		iter.ReportError("ReadArray", "expect [ or , or ] or n, but found "+string([]byte{c}))
		return
	}
}


func (iter *Iterator) ReadArrayCB(callback func(*Iterator) bool) (ret bool) {
	c := iter.nextToken()
	if c == '[' {
		if !iter.incrementDepth() {
			return false
		}
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			if !callback(iter) {
				iter.decrementDepth()
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				if !callback(iter) {
					iter.decrementDepth()
					return false
				}
				c = iter.nextToken()
			}
			if c != ']' {
				iter.ReportError("ReadArrayCB", "expect ] in the end, but found "+string([]byte{c}))
				iter.decrementDepth()
				return false
			}
			return iter.decrementDepth()
		}
		return iter.decrementDepth()
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true 
	}
	iter.ReportError("ReadArrayCB", "expect [ or n, but found "+string([]byte{c}))
	return false
}
