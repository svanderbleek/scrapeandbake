package main

import (
	"fmt"
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/crawler"
	"github.com/rentapplication/craigjr/proxy"
)

func main() {
	proxy.LoadDefault()
	urls := make(chan string)
	c1 := crawler.New(urls)
	c2 := crawler.New(urls)
	go c1.Crawl(&craigslist.Posts{City: "atlanta"})
	go c2.Crawl(&craigslist.Posts{City: "atlanta"})
	go stupidUrlConsumer(urls)
	c1.WaitDone()
	c2.WaitDone()
	close(urls)
}

func stupidUrlConsumer(urls crawler.UrlStream) {
	defer func() { recover() }()
	<-urls
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
