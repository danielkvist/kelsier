package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/danielkvist/fitchner"
)

// Run receives an URL via the command line and extracts all
// the links he finds in it to then check them asynchronously.
func Run() {
	url := flag.String("url", "https://www.google.com", "URL")
	flag.Parse()

	client := &http.Client{}

	parsedURL := parseURL("", *url)
	links, err := fetchLinks(parsedURL, client)
	if err != nil {
		log.Fatalf("while fetching links on %q: %v\n", parsedURL, err)
	}

	var wg sync.WaitGroup
	for link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			makeRequest(link, client)
		}(link)
	}
	wg.Wait()
}

func fetchLinks(url string, c *http.Client) (map[string]bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request for %q: %v", url, err)
	}

	b, err := fitchner.Fetch(c, req)
	if err != nil {
		return nil, fmt.Errorf("while fetching from %q: %v", url, err)
	}

	nodes, err := fitchner.Filter(b, "", "href", "")
	if err != nil {
		return nil, fmt.Errorf("while filtering body of %q: %v", url, err)
	}

	links := map[string]bool{}

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				parsedURL := parseURL(url, attr.Val)
				links[parsedURL] = true
				break
			}
		}
	}

	return links, nil
}

func makeRequest(url string, c *http.Client) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("while creating a new request for %q: %v\n", url, err)
		return
	}

	status := checkStatus(c, req)
	printResult(status, url)
}

func checkStatus(c *http.Client, req *http.Request) int {
	resp, err := c.Do(req)
	if err != nil {
		return http.StatusBadRequest
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func parseURL(baseURL, url string) string {
	switch {
	case len(url) <= 1:
		return baseURL
	case strings.HasPrefix(url, "/"):
		return strings.TrimSuffix(baseURL, "/") + url
	case strings.HasPrefix(url, "#"):
		return baseURL + url
	case strings.HasPrefix(url, "www"):
		return "https://" + url + "/"
	case strings.HasPrefix(url, "mailto:"):
		return strings.TrimPrefix(url, "mailto:")
	default:
		return url
	}
}

func printResult(status int, url string) {
	fmt.Printf("%v - %s\n", status, url)
}
