package crawler

import (
	"fmt"
	"github.com/thedeepak12/arachne/internal/fetcher"
	"github.com/thedeepak12/arachne/internal/frontier"
	"github.com/thedeepak12/arachne/internal/parser"
)

type Crawler struct {
	fetcher *fetcher.Fetcher
	queue   *frontier.Queue
	visited *frontier.Visited
	maxDepth int
}

func New(timeout int, maxDepth int) *Crawler {
	return &Crawler{
		fetcher: fetcher.New(timeout),
		queue:   frontier.NewQueue(),
		visited: frontier.NewVisited(),
		maxDepth: maxDepth,
	}
}

func (c *Crawler) Crawl(seedURL string) {
	seedTask := &frontier.Task{
		URL: seedURL,
		Depth: 0,
	}

	c.queue.Push(seedTask)
	c.visited.Add(seedURL)

	for !c.queue.IsEmpty() {
		task := c.queue.Pop()

		fmt.Printf("Crawling: %s (depth: %d)\n", task.URL, task.Depth)

		body, err := c.fetcher.Fetch(task.URL)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", task.URL, err)
			continue
		}

		links := parser.Parse(body, task.URL)
		fmt.Printf("Found %d links\n", len(links))

		if task.Depth < c.maxDepth {
			for _, link := range links {
				if !c.visited.Contains(link) {
					c.visited.Add(link)
					newTask := &frontier.Task{
						URL: link,
						Depth: task.Depth + 1,
					}
					c.queue.Push(newTask)
				}
			}
		}
	}
}
