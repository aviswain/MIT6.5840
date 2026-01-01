package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

//
// Serial crawler
//

func SerialCrawl(url string, fetcher Fetcher, fetched map[string]string) {
	if _, exists := fetched[url]; exists {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fetched[url] = ""
		return
	}
	fetched[url] = body
	for _, u := range urls {
		SerialCrawl(u, fetcher, fetched)
	}
}

//
// Concurrent crawler with shared state and Mutex
//

type protectedFetchedURLs struct {
	mtx sync.Mutex
	fetchedURLs map[string]*string
}

func MutexConcurrentCrawl(url string, fetcher Fetcher, fetched *protectedFetchedURLs) {
	fetched.mtx.Lock()
	_, exists := fetched.fetchedURLs[url]
	if exists {
		fetched.mtx.Unlock()
		return
	}
	fetched.fetchedURLs[url] = nil
	fetched.mtx.Unlock()

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fetched.mtx.Lock()
		emptyStrCpy := ""
		fetched.fetchedURLs[url] = &emptyStrCpy
		fetched.mtx.Unlock()
		return
	}
	
	fetched.mtx.Lock()
	bodyCpy := body
	fetched.fetchedURLs[url] = &bodyCpy
	fetched.mtx.Unlock()

	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			MutexConcurrentCrawl(u, fetcher, fetched)
		}(u)
	}
	done.Wait()
}

//
// Concurrent crawler with channels
//

type fetchResult struct {
	url	 string
	body string
	urls []string
	err  error
}

func worker(url string, ch chan fetchResult, fetcher Fetcher) {
	body, urls, err := fetcher.Fetch(url)
	ch <- fetchResult{url: url, body: body, urls: urls, err: err}
}

func coordinator(ch chan fetchResult, fetcher Fetcher) {
	pending := 1
	fetched := make(map[string]string)
	for result := range ch {
		if result.err != nil {
			fetched[result.url] = ""
		} else {
			fetched[result.url] = result.body
		}
		for _, u := range result.urls {
			_, exists := fetched[u]
			if !exists {
				fetched[u] = ""
				pending++
				go worker(u, ch, fetcher)
			}
		}
		pending -= 1
		if pending == 0 {
			close(ch)
			return
		}
	}
}

func ChannelConcurrentCrawl(url string, fetcher Fetcher) {
	ch := make(chan fetchResult, 1)
	go worker(url, ch, fetcher)
	coordinator(ch, fetcher)
}

func main() {

	fmt.Println("=== SERIAL CRAWL ===")
	SerialCrawl("https://golang.org/", fetcher, make(map[string]string))

	concurrentFetchedURLs := &protectedFetchedURLs{
		fetchedURLs: make(map[string]*string),
	}
	fmt.Println("=== CONCURRENT MUTEX ===")
	MutexConcurrentCrawl("https://golang.org/", fetcher, concurrentFetchedURLs)

	fmt.Println("=== CONCURRENT CHANNEL ===")
	ChannelConcurrentCrawl("https://golang.org/", fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("Found: %s %s\n", url, res.body)
		return res.body, res.urls, nil
	}
	fmt.Printf("missing: %s\n", url)
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
