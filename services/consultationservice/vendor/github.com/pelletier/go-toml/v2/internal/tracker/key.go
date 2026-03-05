package tracker

import "github.com/pelletier/go-toml/v2/unstable"



type KeyTracker struct {
	k []string
}


func (t *KeyTracker) UpdateTable(node *unstable.Node) {
	t.reset()
	t.Push(node)
}


func (t *KeyTracker) UpdateArrayTable(node *unstable.Node) {
	t.reset()
	t.Push(node)
}


func (t *KeyTracker) Push(node *unstable.Node) {
	it := node.Key()
	for it.Next() {
		t.k = append(t.k, string(it.Node().Data))
	}
}


func (t *KeyTracker) Pop(node *unstable.Node) {
	it := node.Key()
	for it.Next() {
		t.k = t.k[:len(t.k)-1]
	}
}


func (t *KeyTracker) Key() []string {
	k := make([]string, len(t.k))
	copy(k, t.k)
	return k
}

func (t *KeyTracker) reset() {
	t.k = t.k[:0]
}
