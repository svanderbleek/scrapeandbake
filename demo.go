package main

import (
	"fmt"
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/crawler"
	"github.com/rentapplication/craigjr/proxy"
)

func main() {
	atlanta := &craigslist.Posts{City: "atlanta"}
	urls := crawler.NewUrlQueue()
	go crawler.Crawl(atlanta, urls)
	for {
		url := urls.Pop()
		fmt.Println(url)
	}
}

func proxyDemo() {
	proxyList := proxy.NewList()
	responses := make(chan string)

	for i := 0; i < 4; i++ {
		go func() {
			for {
				response := proxyList.Get("http://atlanta.craigslist.org/search/apa?s=100")
				if response.Error == nil {
					responses <- "ok"
				} else {
					responses <- response.Error.Error()
				}
			}
		}()
	}

	for {
		select {
		case r := <-responses:
			fmt.Println(r)
		}
	}
}
