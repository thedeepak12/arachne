package frontier

import (
	"sync"
)

type Visited struct {
	urls map[string]bool
	mu   sync.Mutex
}

func NewVisited() *Visited {
	return &Visited{
		urls: make(map[string]bool),
	}
}

func (v *Visited) Add(url string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.urls[url] = true
}

func (v *Visited) Contains(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.urls[url]
}
