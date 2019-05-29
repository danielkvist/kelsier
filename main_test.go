package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalize(t *testing.T) {
	tc := []struct {
		rawURL   string
		base     string
		expected string
	}{
		{
			rawURL:   "https://mysite.com",
			base:     "",
			expected: "https://mysite.com",
		},
		{
			rawURL:   "a.es",
			base:     "",
			expected: "https://a.es",
		},
		{
			rawURL:   "/",
			base:     "https://site.com",
			expected: "https://site.com/",
		},
		{
			rawURL:   "www.abc.com",
			base:     "",
			expected: "https://abc.com",
		},
		{
			rawURL:   "#id",
			base:     "https://qwerty.cat",
			expected: "https://qwerty.cat/#id",
		},
		{
			rawURL:   "mailto:me@me.com",
			base:     "",
			expected: "me@me.com",
		},
		{
			rawURL:   "/a/b/c/#id",
			base:     "https://site.com",
			expected: "https://site.com/a/b/c/#id",
		},
	}

	for _, tt := range tc {
		t.Run(tt.rawURL, func(t *testing.T) {
			normalized := normalize(tt.base, tt.rawURL)
			if tt.expected != normalized {
				t.Fatalf("expected URL %q to be normalized as %q. got=%q", tt.rawURL, tt.expected, normalized)
			}
		})
	}
}

func TestFetchLinks(t *testing.T) {
	template := `
		<a href="#"></a>
		<section>
			<a href="/"></a>
			<a href="mysite.com"></a>
		</section>
		<footer>
			<div>
				<a href="www.me.cat"></a>
			<div>
		</footer>
	`
	handler := testingHandler(template)
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &http.Client{}
	expectedLinks := []string{
		server.URL + "/#",
		server.URL + "/",
		"https://mysite.com",
		"https://me.cat",
	}

	links, err := fetchLinks(c, server.URL)
	if err != nil {
		t.Fatalf("while extracting links for test: %v", err)
	}

	if len(expectedLinks) != len(links) {
		t.Fatalf("expected %v of links. got=%v", len(expectedLinks), len(links))
	}

	for i, l := range links {
		if l != expectedLinks[i] {
			t.Errorf("expected link %q. got=%q", expectedLinks[i], l)
		}
	}
}

func TestLinksOut(t *testing.T) {
	testLinks := map[string]struct{}{
		"https://me.com":    struct{}{},
		"https://abc.cat":   struct{}{},
		"https://qwerty.es": struct{}{},
		"https://as.ru":     struct{}{},
	}

	var testLinksKeys []string
	for k := range testLinks {
		testLinksKeys = append(testLinksKeys, k)
	}

	var links []string
	linksChan := linksOut("", testLinksKeys)
	for l := range linksChan {
		links = append(links, l)
	}

	for _, l := range links {
		if _, ok := testLinks[l]; !ok {
			t.Errorf("expecting to found link %q on fan out result.", l)
		}
	}
}

func TestCheck(t *testing.T) {
	handler := testingHandler("<h1>Hello, World!</h1>")
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	c := &http.Client{}
	linksChan := make(chan string)
	go func() { linksChan <- server.URL }()

	expectedStatus := status(http.StatusOK, server.URL)
	in := check(c, linksChan)
	result := <-in
	if result != expectedStatus {
		t.Errorf("expected result to be %q. got=%q", expectedStatus, result)
	}
}

func testingHandler(tpl string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, tpl)
	}
}
