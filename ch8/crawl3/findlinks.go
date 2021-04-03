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
	"os"

	"gopl.io/ch5/links"
)

var depthflag = flag.Int("depth", 0, "for example: -depth=3, default value 0 means only passed urls will be traversed")

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

type Link struct {
	url   string
	depth int
}

//!+
func main() {
	flag.Parse()

	worklist := make(chan []Link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	{
		links := []Link{}
		for _, arg := range os.Args[flag.NFlag()+1:] {
			links = append(links, Link{url: arg, depth: 0})
		}

		go func() {
			worklist <- links
		}()
	}

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := []Link{}

				foundURLs := crawl(link.url)
				for _, url := range foundURLs {
					foundLinks = append(foundLinks, Link{url: url, depth: link.depth + 1})
				}

				go func() { worklist <- foundLinks }()
			}
		}()
	}

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
}

//!-
