

package jsoniter

import (
	"fmt"
	"io"
)

func (iter *Iterator) skipNumber() {
	if !iter.trySkipNumber() {
		iter.unreadByte()
		if iter.Error != nil && iter.Error != io.EOF {
			return
		}
		iter.ReadFloat64()
		if iter.Error != nil && iter.Error != io.EOF {
			iter.Error = nil
			iter.ReadBigFloat()
		}
	}
}

func (iter *Iterator) trySkipNumber() bool {
	dotFound := false
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if dotFound {
				iter.ReportError("validateNumber", `more than one dot found in number`)
				return true 
			}
			if i+1 == iter.tail {
				return false
			}
			c = iter.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				iter.ReportError("validateNumber", `missing digit after dot`)
				return true 
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if iter.head == i {
					return false 
				}
				iter.head = i
				return true 
			}
			return false 
		}
	}
	return false
}

func (iter *Iterator) skipString() {
	if !iter.trySkipString() {
		iter.unreadByte()
		iter.ReadString()
	}
}

func (iter *Iterator) trySkipString() bool {
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		if c == '"' {
			iter.head = i + 1
			return true 
		} else if c == '\\' {
			return false
		} else if c < ' ' {
			iter.ReportError("trySkipString",
				fmt.Sprintf(`invalid control character found: %d`, c))
			return true 
		}
	}
	return false
}

func (iter *Iterator) skipObject() {
	iter.unreadByte()
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		iter.Skip()
		return true
	})
}

func (iter *Iterator) skipArray() {
	iter.unreadByte()
	iter.ReadArrayCB(func(iter *Iterator) bool {
		iter.Skip()
		return true
	})
}
