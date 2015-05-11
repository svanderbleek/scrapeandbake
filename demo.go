package main

import (
	"fmt"
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/crawler"
	"github.com/rentapplication/craigjr/proxy"
	"github.com/rentapplication/craigjr/visitor"
)

func main() {
	proxy := proxy.NewList()
	proxy.LoadDefault()
	posts := make(chan crawler.PaginationIterator)
	go craigslist.StreamCities(posts)
	crawlers := crawler.NewPool(proxy)
	go crawlers.Crawl(posts)
	visitors := visitor.NewPool(proxy)
	go visitors.Visit(crawlers.Urls)
	drainPosts(visitors.Posts)
}

func drainPosts(posts <-chan *craigslist.Post) {
	for {
		post := <-posts
		fmt.Println(post)
	}
}
