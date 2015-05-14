package main

import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/crawler"
	"github.com/rentapplication/craigjr/indexer"
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

	indexPosts(visitors.Posts)
}

func indexPosts(posts <-chan *craigslist.Post) {
	index := indexer.New()
	for {
		post := <-posts
		index.Index(post)
	}
}
