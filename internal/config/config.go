package config

import (
	"flag"
)

type Config struct {
	URL     string
	Workers int
	Depth   int
	Timeout int
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.URL, "url", "", "URL to crawl")
	flag.IntVar(&cfg.Workers, "workers", 5, "Number of workers")
	flag.IntVar(&cfg.Depth, "depth", 1, "Depth to crawl")
	flag.IntVar(&cfg.Timeout, "timeout", 10, "Timeout in seconds")
	flag.Parse()

	return cfg
}
