package crawler

import (
	"fmt"
	"github.com/thedeepak12/arachne/internal/fetcher"
	"github.com/thedeepak12/arachne/internal/frontier"
	"github.com/thedeepak12/arachne/internal/parser"
)

type Worker struct {
	id       int
	fetcher  *fetcher.Fetcher
	queue    *frontier.Queue
	visited  *frontier.Visited
	maxDepth int
}

func NewWorker(id int, fetcher *fetcher.Fetcher, queue *frontier.Queue, visited *frontier.Visited, maxDepth int) *Worker {
	return &Worker{
		id:       id,
		fetcher:  fetcher,
		queue:    queue,
		visited:  visited,
		maxDepth: maxDepth,
	}
}

func (w *Worker) Run(taskChan chan *frontier.Task) {
	for task := range taskChan {
		w.processTask(task)
	}
}

func (w *Worker) processTask(task *frontier.Task) {
	fmt.Printf("[worker-%d] Crawling: %s (depth: %d)\n", w.id, task.URL, task.Depth)

	body, err := w.fetcher.Fetch(task.URL)
	if err != nil {
		fmt.Printf("[worker-%d] Error fetching %s: %v\n", w.id, task.URL, err)
		return
	}

	links := parser.Parse(body, task.URL)
	fmt.Printf("[worker-%d] Found %d links\n", w.id, len(links))

	if task.Depth < w.maxDepth {
		for _, link := range links {
			if !w.visited.Contains(link) {
				w.visited.Add(link)
				w.queue.Push(&frontier.Task{
					URL:   link,
					Depth: task.Depth + 1,
				})
			}
		}
	}
}
