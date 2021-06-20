package main

import (
	"fmt"
	"github.com/89z/mech"
	"log"
	"net/http"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		wg.Done()
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, wg)
	}
	return
}

type URLFetcher struct {

}

func (f URLFetcher) Fetch(url string) (body string, urls []string, err error) {
	r, err := http.Get(url)
	if err != nil {
		log.Fatalln("Could not get URL: ", url)
		return "", nil, err
	}
	defer r.Body.Close()

	doc, err := mech.Parse(r.Body)
	if err != nil {
		log.Fatalln("Could not parse body of URL: ", url)
		return "", nil, err
	}
	a := doc.ByTag("a")
	var parsedURLs []string
	for a.Scan() {
		href := a.Attr("href")
		fmt.Println(href)
		parsedURLs = append(parsedURLs, href)
	}
	return "", parsedURLs, nil
}

var fetcher = URLFetcher{}


func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("https://golang.org/", 4, fetcher, &wg)
	wg.Wait()
}