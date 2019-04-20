package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/danielkvist/fitchner"
)

// Run parses the URL argument an executes the
// makeRequest function asynchronously.
func Run() {
	url := flag.String("url", "https://www.google.com", "URL")
	flag.Parse()

	var wg sync.WaitGroup
	client := &http.Client{}

	wg.Add(1)
	go makeRequest(*url, client, &wg)
	wg.Wait()
}

func makeRequest(url string, c *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("while creating a new request for %q: %v\n", url, err)
		return
	}

	status := checkStatus(c, req)
	printResult(os.Stdout, status, url)

	if status != http.StatusOK {
		return
	}

	b, err := fitchner.Fetch(c, req)
	if err != nil {
		log.Printf("while fetching from %q: %v\n", url, err)
		return
	}

	nodes, err := fitchner.Filter(b, "", "href", "")
	if err != nil {
		log.Printf("while filtering body of %q: %v", url, err)
		return
	}

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				wg.Add(1)
				parsedURL := parseURL(url, attr.Val)
				go func() {
					defer wg.Done()
					makeRequest(parsedURL, c, wg)
				}()
				break
			}
		}
	}
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
	if strings.HasPrefix(url, "/") {
		return baseURL + url
	}

	if strings.HasPrefix(url, "#") {
		return baseURL + url
	}

	if strings.HasPrefix(url, "mailto:") {
		return strings.TrimPrefix(url, "mailto:")
	}

	return url
}

func printResult(w io.Writer, status int, url string) {
	fmt.Fprintf(w, "%v - %s\n", status, url)
}
