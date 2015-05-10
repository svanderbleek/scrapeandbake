package crawler

import (
	"fmt"
	"github.com/rentapplication/craigjr/proxy"
)

type Pool struct {
	Urls     chan string
	crawlers chan *Crawler
	proxy    proxy.Proxier
}

func NewPool(proxier proxy.Proxier) *Pool {
	return &Pool{
		Urls:  make(chan string),
		proxy: proxier,
	}
}

func (pool *Pool) Crawl(pages PaginationIterator) {
	crawler := pool.Crawler()
	fmt.Println(crawler)
	go crawler.Crawl(pages)
}

func (pool *Pool) Crawler() *Crawler {
	var crawler *Crawler
	select {
	case crawler = <-pool.crawlers:
	default:
		crawler = NewCrawler(pool)
	}
	return crawler
}
