package frontier

import (
	"sync"
)

type Queue struct {
	urls []string
	mu   sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		urls: make([]string, 0),
	}
}

func (q *Queue) Push(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.urls = append(q.urls, url)
}

func (q *Queue) Pop() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.urls) == 0 {
		return ""
	}

	url := q.urls[0]
	q.urls = q.urls[1:]
	return url
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.urls) == 0
}
