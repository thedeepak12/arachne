package parser

import (
	"golang.org/x/net/html"
	"net/url"
	"path"
	"strings"
)

func Parse(htmlContent, baseURL string) []string {
	links := make([]string, 0)

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return links
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					normalized := normalizeURL(attr.Val, baseURL)
					if normalized != "" {
						links = append(links, normalized)
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return links
}

func normalizeURL(rawURL, baseURL string) string {
	if strings.HasPrefix(rawURL, "mailto:") || strings.HasPrefix(rawURL, "javascript:") || strings.HasPrefix(rawURL, "tel:") || strings.HasPrefix(rawURL, "data:") {
		return ""
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	if !parsed.IsAbs() {
		base, err := url.Parse(baseURL)
		if err != nil {
			return ""
		}

		parsed = base.ResolveReference(parsed)
	}

	parsed.Fragment = ""
	parsed.Path = path.Clean(parsed.Path)

	return parsed.String()
}
