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
	crawlers := crawler.NewPool(proxy)
	craigslist.CitiesAsPosts(crawlers.Crawl)
	visitors := visitor.NewPool(proxy)
	visitors.Visit(crawlers.Urls)
	drainPosts(visitors.Posts)
}

func drainPosts(posts <-chan *craigslist.Post) {
	for {
		post := <-posts
		fmt.Println(post)
	}
}
