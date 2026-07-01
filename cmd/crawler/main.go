package main

import (
	"fmt"
	"github.com/thedeepak12/arachne/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Arachne - Web Crawler\n")
	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("Workers: %d\n", cfg.Workers)
	fmt.Printf("Depth: %d\n", cfg.Depth)
	fmt.Printf("Timeout: %d\n", cfg.Timeout)
}
