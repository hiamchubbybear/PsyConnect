package rand

import (
	"math/rand"
	"sync"
)


func Int() int { return pseudo.Int() }



func Intn(n int) int { return pseudo.Intn(n) }



func Int63n(n int64) int64 { return pseudo.Int63n(n) }


func Perm(n int) []int { return pseudo.Perm(n) }




func Seed(n int64) { pseudo.Seed(n) }

var pseudo = rand.New(&source{src: rand.NewSource(1)})

type source struct {
	src rand.Source
	mu  sync.Mutex
}

func (s *source) Int63() int64 {
	s.mu.Lock()
	n := s.src.Int63()
	s.mu.Unlock()
	return n
}

func (s *source) Seed(seed int64) {
	s.mu.Lock()
	s.src.Seed(seed)
	s.mu.Unlock()
}




func Shuffle(n int, swap func(i, j int)) { pseudo.Shuffle(n, swap) }
