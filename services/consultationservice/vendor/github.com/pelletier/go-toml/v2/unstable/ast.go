package unstable

import (
	"fmt"
	"unsafe"

	"github.com/pelletier/go-toml/v2/internal/danger"
)












type Iterator struct {
	started bool
	node    *Node
}



func (c *Iterator) Next() bool {
	if !c.started {
		c.started = true
	} else if c.node.Valid() {
		c.node = c.node.Next()
	}
	return c.node.Valid()
}



func (c *Iterator) IsLast() bool {
	return c.node.next == 0
}


func (c *Iterator) Node() *Node {
	return c.node
}
















type Node struct {
	Kind Kind
	Raw  Range  
	Data []byte 

	
	
	
	next  int 
	child int 
}


type Range struct {
	Offset uint32
	Length uint32
}


func (n *Node) Next() *Node {
	if n.next == 0 {
		return nil
	}
	ptr := unsafe.Pointer(n)
	size := unsafe.Sizeof(Node{})
	return (*Node)(danger.Stride(ptr, size, n.next))
}




func (n *Node) Child() *Node {
	if n.child == 0 {
		return nil
	}
	ptr := unsafe.Pointer(n)
	size := unsafe.Sizeof(Node{})
	return (*Node)(danger.Stride(ptr, size, n.child))
}


func (n *Node) Valid() bool {
	return n != nil
}




func (n *Node) Key() Iterator {
	switch n.Kind {
	case KeyValue:
		value := n.Child()
		if !value.Valid() {
			panic(fmt.Errorf("KeyValue should have at least two children"))
		}
		return Iterator{node: value.Next()}
	case Table, ArrayTable:
		return Iterator{node: n.Child()}
	default:
		panic(fmt.Errorf("Key() is not supported on a %s", n.Kind))
	}
}




func (n *Node) Value() *Node {
	return n.Child()
}


func (n *Node) Children() Iterator {
	return Iterator{node: n.Child()}
}
