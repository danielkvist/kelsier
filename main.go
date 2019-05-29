package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/danielkvist/fitchner"
)

func main() {
	urls := os.Args[1:]
	if len(urls) < 1 {
		fmt.Printf("no URLs provided.")
		os.Exit(1)
	}

	c := &http.Client{}
	chans := []<-chan string{}
	for _, url := range urls {
		normalized := normalize("", url)
		links, err := fetchLinks(c, normalized)
		if err != nil {
			fmt.Printf("while extracting links from %q: %v", url, err)
			continue
		}

		ch := linksOut(url, links)
		chans = append(chans, check(c, ch))
	}

	for result := range merge(chans...) {
		fmt.Printf(result)
	}
}

func normalize(base, url string) string {
	base = strings.TrimSuffix(base, "/")

	switch {
	case strings.HasPrefix(url, "/"):
		return base + url
	case strings.HasPrefix(url, "#"):
		return base + "/" + url
	case strings.HasPrefix(url, "www."):
		url = strings.TrimPrefix(url, "www.")
		return "https://" + url
	case strings.HasPrefix(url, "mailto:"):
		return strings.TrimPrefix(url, "mailto:")
	case !strings.HasPrefix(url, "http"):
		return "https://" + url
	default:
		return url
	}
}

func fetchLinks(c *http.Client, url string) ([]string, error) {
	url = normalize("", url)
	f, err := fitchner.NewFetcher(fitchner.WithClient(c), fitchner.WithSimpleGetRequest(url))
	if err != nil {
		return nil, fmt.Errorf("while creating a new Fetcher to extract the links from %q: %v", url, err)
	}

	d, err := f.Do()
	if err != nil {
		return nil, fmt.Errorf("while making a request to %q: %v", url, err)
	}

	data := bytes.NewReader(d)
	links, err := fitchner.Links(data)
	if err != nil {
		return nil, fmt.Errorf("while extracting the links of the response body of %q: %v", url, err)
	}

	for i, link := range links {
		links[i] = normalize(url, link)
	}

	return links, nil
}

func linksOut(base string, links []string) <-chan string {
	out := make(chan string, len(links))
	go func() {
		for _, l := range links {
			func(link string) {
				out <- normalize(base, l)
			}(l)
		}
		close(out)
	}()

	return out
}

func check(c *http.Client, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for l := range in {
			func(link string) {
				req, err := http.NewRequest(http.MethodGet, link, nil)
				if err != nil {
					out <- status(0, link)
					return
				}

				resp, err := c.Do(req)
				if err != nil {
					out <- status(http.StatusBadRequest, link)
					return
				}

				out <- status(resp.StatusCode, link)
			}(l)
		}

		close(out)
	}()

	return out
}

func merge(chans ...<-chan string) <-chan string {
	out := make(chan string)
	go func() {
		var wg sync.WaitGroup

		wg.Add(len(chans))
		for _, c := range chans {
			go func(c <-chan string) {
				defer wg.Done()

				for msg := range c {
					out <- msg
				}
			}(c)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func status(s int, url string) string {
	return fmt.Sprintf("%v - %s\n", s, url)
}
