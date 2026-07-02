package main

import (
	"fmt"

	"github.com/thedeepak12/arachne/internal/config"
	"github.com/thedeepak12/arachne/internal/crawler"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Arachne - Web Crawler\n")
	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("Workers: %d\n", cfg.Workers)
	fmt.Printf("Depth: %d\n", cfg.Depth)
	fmt.Printf("Timeout: %d\n", cfg.Timeout)

	if cfg.URL == "" {
		fmt.Println("Please provide a URL with -url flag.")
		return
	}

	c := crawler.New(cfg.Timeout)
	c.Crawl(cfg.URL)

	fmt.Println("Crawling complete")
}
