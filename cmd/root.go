package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/danielkvist/fitchner"
)

// Run receives an URL from when it extracts all the links.
// It parses and check each one asynchronously printing the results.
func Run() {
	url := flag.String("url", "https://www.google.com", "URL")
	flag.Parse()

	parsedURL := parseURL("", *url)
	links, err := fetchLinks(parsedURL)
	if err != nil {
		log.Fatalf("while searching for links on %q: %v\n", parsedURL, err)
	}

	var wg sync.WaitGroup
	client := &http.Client{}
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodGet, link, nil)
			if err != nil {
				log.Printf("while creating a new request for %q: %v\n", link, err)
				return
			}

			fmt.Printf("%v - %s\n", checkStatus(client, req), link)
		}(link)
	}
	wg.Wait()
}

func fetchLinks(url string) ([]string, error) {
	optReq := fitchner.WithSimpleGetRequest(url)
	f, err := fitchner.NewFetcher(optReq)
	if err != nil {
		return nil, fmt.Errorf("while creating a new Fetcher to make request: %v", err)
	}

	b, err := f.Do()
	if err != nil {
		return nil, fmt.Errorf("while fetching from %q: %v", url, err)
	}

	data := bytes.NewReader(b)
	links, err := fitchner.Links(data)
	if err != nil {
		return nil, fmt.Errorf("while searching for links on response for %q: %v", url, err)
	}

	parsedLinks := make([]string, 0, len(links))
	for _, l := range links {
		parsedLinks = append(parsedLinks, parseURL(url, l))
	}

	return parsedLinks, nil
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
}
