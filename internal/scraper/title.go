package scraper

import (
	"context"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// FetchTitle fetches the <title> tag from a URL.
// Returns URL as fallback on any error.
func FetchTitle(ctx context.Context, url string) string {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return url
	}
	req.Header.Set("User-Agent", "LinkVault/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return url
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return url
	}

	title := extractTitle(doc)
	if title == "" {
		return url
	}
	return strings.TrimSpace(title)
}

func extractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil {
			return n.FirstChild.Data
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if t := extractTitle(c); t != "" {
			return t
		}
	}
	return ""
}