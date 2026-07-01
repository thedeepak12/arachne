package main

import (
	"fmt"
	"github.com/thedeepak12/arachne/internal/config"
	"github.com/thedeepak12/arachne/internal/fetcher"
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

	f := fetcher.New(cfg.Timeout)
	body, err := f.Fetch(cfg.URL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Successfully fetched %d bytes\n", len(body))
}
