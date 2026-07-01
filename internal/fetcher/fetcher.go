package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	client *http.Client
}

func New(timeout int) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (f *Fetcher) Fetch(url string) (string, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body from %s: %w", url, err)
	}

	return string(body), nil
}
