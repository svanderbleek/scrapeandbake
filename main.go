package main

import (
	"./craigslist"
	"./crawler"
	"./scraper"
	"fmt"
)

func main() {
	city := craigslist.Cities[0]
	posts := &craigslist.Posts{City: city}
	queue := scraper.NewUrlQueue()
	go crawler.Crawl(posts, queue)
	for {
		fmt.Println(queue.Pop())
	}
}
