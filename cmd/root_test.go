package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const tpl = `<!DOCTYPE HTML>
	<html>
	<head>
	<title>Testing Kelsier</title>
	</head>
	<body>
	<header>
	<h1 class="title">Testing</h1>
	</header>
	<a href="https://www.google.com">Google</a>
	<a href="www.bing.com">Bing</a>
	<a href="https://www.github.com">GitHub</a>
	<a href="#title">Title</a>
	<a href="/blog">Blog</a>
	</body>
	</html>`

func TestFetchLinks(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, tpl)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &http.Client{}
	links, err := fetchLinks(server.URL, client)
	if err != nil {
		t.Fatalf("while fetching links: %v", err)
	}

	tt := []struct {
		name string
		link string
	}{
		{"Google", "https://www.google.com"},
		{"Bing", "https://www.bing.com/"},
		{"GitHub", "https://www.github.com"},
		{"ID", server.URL + "#title"},
		{"Page", server.URL + "/blog"},
	}

	fmt.Println(links)
	for _, tc := range tt {
		t.Run(tc.link, func(t *testing.T) {
			if ok := links[tc.link]; !ok {
				t.Fatalf("expected to found %q", tc.link)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	baseURL := "https://www.golang.com/"
	tt := []struct {
		name        string
		url         string
		expectedURL string
	}{
		{"Root", "/", baseURL},
		{"HTTPS", "www.golang.com", "https://www.golang.com/"},
		{"ID", "#about", baseURL + "#about"},
		{"Page", "/blog/", baseURL + "blog/"},
		{"Mail", "mailto:me@somewhere.com", "me@somewhere.com"},
		{"Empty", "", baseURL},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			parsedURL := parseURL(baseURL, tc.url)
			if parsedURL != tc.expectedURL {
				t.Fatalf("expected %q. got=%q", tc.expectedURL, parsedURL)
			}
		})
	}
}

func BenchmarkFetchLinks(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, tpl)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &http.Client{}
	for i := 0; i < b.N; i++ {
		fetchLinks(server.URL, client)
	}
}
