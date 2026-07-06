package crawler

import (
	"sync"
	"time"

	"github.com/thedeepak12/arachne/internal/fetcher"
	"github.com/thedeepak12/arachne/internal/frontier"
)

type Crawler struct {
	fetcher        *fetcher.Fetcher
	queue          *frontier.Queue
	visited        *frontier.Visited
	maxDepth       int
	numWorkers     int
	totalLimit     int
	tasksProcessed int
}

func New(timeout int, maxDepth int, numWorkers int) *Crawler {
	return &Crawler{
		fetcher:        fetcher.New(timeout),
		queue:          frontier.NewQueue(),
		visited:        frontier.NewVisited(),
		maxDepth:       maxDepth,
		numWorkers:     numWorkers,
		totalLimit:     1000,
		tasksProcessed: 0,
	}
}

func (c *Crawler) Crawl(seedURL string) {
	taskChan := make(chan *frontier.Task, 100)

	workers := make([]*Worker, c.numWorkers)
	for i := 0; i < c.numWorkers; i++ {
		workers[i] = NewWorker(i, c.fetcher, c.queue, c.visited, c.maxDepth)
	}

	seedTask := &frontier.Task{URL: seedURL, Depth: 0}
	c.queue.Push(seedTask)
	c.visited.Add(seedURL)

	var activeWorkers int
	var mu sync.Mutex

	go func() {
		emptyCount := 0
		for {
			task := c.queue.Pop()
			if task != nil {
				emptyCount = 0
				c.tasksProcessed++
				if c.tasksProcessed >= c.totalLimit {
					break
				}
				mu.Lock()
				activeWorkers++
				mu.Unlock()
				taskChan <- task
			} else {
				emptyCount++
				mu.Lock()
				noActive := activeWorkers == 0
				mu.Unlock()
				if emptyCount > 10 && noActive {
					break
				}
				time.Sleep(50 * time.Millisecond)
			}
		}
		close(taskChan)
	}()

	var wg sync.WaitGroup
	for i := 0; i < c.numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range taskChan {
				workers[workerID].processTask(task)
				mu.Lock()
				activeWorkers--
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
}
