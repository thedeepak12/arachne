package crawler

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Stats struct {
	startTime       time.Time
	urlsCrawled     atomic.Int64
	successfulFetch atomic.Int64
	failedFetch     atomic.Int64
	linksFound      atomic.Int64
	totalBytes      atomic.Int64
}

func NewStats() *Stats {
	return &Stats{
		startTime: time.Now(),
	}
}

func (s *Stats) Start() {
	s.startTime = time.Now()
}

func (s *Stats) IncrementURLsCrawled() {
	s.urlsCrawled.Add(1)
}

func (s *Stats) IncrementSuccessfulFetch() {
	s.successfulFetch.Add(1)
}

func (s *Stats) IncrementFailedFetch() {
	s.failedFetch.Add(1)
}

func (s *Stats) AddLinksFound(count int) {
	s.linksFound.Add(int64(count))
}

func (s *Stats) AddBytes(count int64) {
	s.totalBytes.Add(count)
}

func (s *Stats) GetDuration() time.Duration {
	return time.Since(s.startTime)
}

func (s *Stats) String() string {
	duration := s.GetDuration()
	return fmt.Sprintf("\n=== Crawl Statistics ===\n"+
		"Duration: %v\n"+
		"URLs Crawled: %d\n"+
		"Successful Fetches: %d\n"+
		"Failed Fetches: %d\n"+
		"Links Found: %d\n"+
		"Total Bytes: %s\n"+
		"========================",
		duration.Round(time.Millisecond),
		s.urlsCrawled.Load(),
		s.successfulFetch.Load(),
		s.failedFetch.Load(),
		s.linksFound.Load(),
		formatBytes(s.totalBytes.Load()))
}

func formatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	exp := 0
	for n >= unit && exp < 4 {
		n /= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n), "KMGT"[exp-1])
}
