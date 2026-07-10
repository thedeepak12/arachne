package crawler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	stats          *Stats
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
		stats:          NewStats(),
	}
}

func (c *Crawler) Crawl(seedURL string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		fmt.Println("\n\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	defer func() {
		fmt.Println(c.stats.String())
	}()

	taskChan := make(chan *frontier.Task, 100)

	workers := make([]*Worker, c.numWorkers)
	for i := 0; i < c.numWorkers; i++ {
		workers[i] = NewWorker(i, c.fetcher, c.queue, c.visited, c.maxDepth, c.stats)
	}

	seedTask := &frontier.Task{URL: seedURL, Depth: 0}
	c.queue.Push(seedTask)
	c.visited.Add(seedURL)

	var wg sync.WaitGroup
	var tasksInFlight sync.WaitGroup

	go func() {
		emptyCount := 0
		for {
			select {
			case <-ctx.Done():
				close(taskChan)
				return
			default:
				task := c.queue.Pop()
				if task != nil {
					emptyCount = 0
					c.tasksProcessed++
					if c.tasksProcessed >= c.totalLimit {
						close(taskChan)
						return
					}
					tasksInFlight.Add(1)
					taskChan <- task
				} else {
					tasksInFlight.Wait()
					if c.queue.IsEmpty() {
						emptyCount++
						if emptyCount >= 50 {
							close(taskChan)
							return
						}
						time.Sleep(100 * time.Millisecond)
					} else {
						emptyCount = 0
					}
				}
			}
		}
	}()

	for i := 0; i < c.numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskChan:
					if !ok {
						return
					}
					workers[workerID].processTask(ctx, task)
					tasksInFlight.Done()
				}
			}
		}(i)
	}

	wg.Wait()
}
