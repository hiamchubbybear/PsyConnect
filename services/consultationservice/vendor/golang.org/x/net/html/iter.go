



//go:build go1.23

package html

import "iter"




func (n *Node) Ancestors() iter.Seq[*Node] {
	_ = n.Parent 

	return func(yield func(*Node) bool) {
		for p := n.Parent; p != nil && yield(p); p = p.Parent {
		}
	}
}





func (n *Node) ChildNodes() iter.Seq[*Node] {
	_ = n.FirstChild 

	return func(yield func(*Node) bool) {
		for c := n.FirstChild; c != nil && yield(c); c = c.NextSibling {
		}
	}

}





func (n *Node) Descendants() iter.Seq[*Node] {
	_ = n.FirstChild 

	return func(yield func(*Node) bool) {
		n.descendants(yield)
	}
}

func (n *Node) descendants(yield func(*Node) bool) bool {
	for c := range n.ChildNodes() {
		if !yield(c) || !c.descendants(yield) {
			return false
		}
	}
	return true
}
