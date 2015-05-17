package main

import (
	"fmt"
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/crawler"
	"github.com/rentapplication/craigjr/proxy"
	"github.com/rentapplication/craigjr/store"
	//"github.com/rentapplication/craigjr/visitor"
)

func main() {
	proxy := proxy.NewList()
	proxy.LoadDefault()

	posts := make(chan crawler.PaginationIterator)
	go craigslist.StreamCities(posts)

	crawlers := crawler.NewPool(proxy)
	go crawlers.Crawl(posts)

	drainUrls(crawlers.Urls)

	//visitors := visitor.NewPool(proxy)
	//go visitors.Visit(crawlers.Urls)

	//indexPosts(visitors.Posts)
}

func drainUrls(urls chan string) {
	count := 0
	for url := range urls {
		count++
		fmt.Println(url)
		fmt.Println(count)
	}
}

func indexPosts(posts <-chan *craigslist.Post) {
	index := store.New()
	for {
		post := <-posts
		index.Store(post)
	}
}
