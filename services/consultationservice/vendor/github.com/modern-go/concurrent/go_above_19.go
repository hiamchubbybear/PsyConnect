

package concurrent

import "sync"


type Map struct {
	sync.Map
}


func NewMap() *Map {
	return &Map{}
}
