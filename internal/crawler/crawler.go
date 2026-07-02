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
}

func New(timeout int) *Crawler {
	return &Crawler{
		fetcher: fetcher.New(timeout),
		queue:   frontier.NewQueue(),
		visited: frontier.NewVisited(),
	}
}

func (c *Crawler) Crawl(seedURL string) {
	c.queue.Push(seedURL)
	c.visited.Add(seedURL)

	count := 0
	for !c.queue.IsEmpty() {
		url := c.queue.Pop()

		fmt.Printf("Crawling: %s\n", url)

		body, err := c.fetcher.Fetch(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
			continue
		}

		links := parser.Parse(body, url)
		fmt.Printf("Found %d links\n", len(links))

		for _, link := range links {
			if !c.visited.Contains(link) {
				c.visited.Add(link)
				c.queue.Push(link)
			}
		}

		count++
		if count >= 100 {
			fmt.Println("Reached test limit of 100 pages")
			break
		}
	}
}
