// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 243.

// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/html"
)

// extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
//
// copied from ch5/links with some change
func extract(url string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Cancel = done

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

var depthflag = flag.Int("depth", 0, "for example: -depth=3, default value 0 means only passed urls will be traversed")

func crawl(url string) []string {
	fmt.Println(url)
	list, err := extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

type Link struct {
	url   string
	depth int
}

var done = make(chan struct{})

func canceled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

//!+
func main() {
	flag.Parse()

	worklist := make(chan []Link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	{
		links := []Link{}
		for _, arg := range os.Args[flag.NFlag()+1:] {
			links = append(links, Link{url: arg, depth: 0})
		}

		go func() {
			worklist <- links
		}()
	}

	wg := new(sync.WaitGroup)

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			for {
				select {
				case <-done:
					wg.Done()
					return
				case link := <-unseenLinks:
					foundLinks := []Link{}

					foundURLs := crawl(link.url)
					for _, url := range foundURLs {
						foundLinks = append(foundLinks, Link{url: url, depth: link.depth + 1})
					}

					go func() {
						select {
						case <-done:
							fmt.Printf("function with id %d is done\n", id)
							return
						case worklist <- foundLinks:
						}
					}()
				}

			}
		}(i)
	}

	go func() {
		wg.Wait()
		fmt.Println("workgroup ready")
		close(worklist)
		fmt.Println("closed worklist")
	}()

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if link.depth > *depthflag {
				continue
			}
			if !seen[link.url] {
				seen[link.url] = true
				unseenLinks <- link
			}
		}
	}
	fmt.Println("herer")
}

//!-
