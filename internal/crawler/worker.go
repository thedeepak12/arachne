package crawler

import (
	"context"
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
	stats    *Stats
}

func NewWorker(id int, fetcher *fetcher.Fetcher, queue *frontier.Queue, visited *frontier.Visited, maxDepth int, stats *Stats) *Worker {
	return &Worker{
		id:       id,
		fetcher:  fetcher,
		queue:    queue,
		visited:  visited,
		maxDepth: maxDepth,
		stats:    stats,
	}
}

func (w *Worker) Run(ctx context.Context, taskChan chan *frontier.Task) {
	for task := range taskChan {
		w.processTask(ctx, task)
	}
}

func (w *Worker) processTask(ctx context.Context, task *frontier.Task) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	w.stats.IncrementURLsCrawled()
	fmt.Printf("[worker-%d] Crawling: %s (depth: %d)\n", w.id, task.URL, task.Depth)

	body, err := w.fetcher.Fetch(task.URL)
	if err != nil {
		w.stats.IncrementFailedFetch()
		fmt.Printf("[worker-%d] Error fetching %s: %v\n", w.id, task.URL, err)
		return
	}

	w.stats.IncrementSuccessfulFetch()
	w.stats.AddBytes(int64(len(body)))

	select {
	case <-ctx.Done():
		return
	default:
	}

	links := parser.Parse(body, task.URL)
	w.stats.AddLinksFound(len(links))
	fmt.Printf("[worker-%d] Found %d links\n", w.id, len(links))

	if task.Depth < w.maxDepth {
		for _, link := range links {
			select {
			case <-ctx.Done():
				return
			default:
				if !w.visited.Contains(link) {
					w.visited.Add(link)
					w.queue.Push(&frontier.Task{URL: link, Depth: task.Depth + 1})
				}
			}
		}
	}
}
