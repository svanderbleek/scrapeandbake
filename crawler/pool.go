package crawler

import (
	"github.com/rentapplication/craigjr/proxy"
)

type Pool struct {
	Urls chan string
	proxy.Proxier
	crawlers chan *Crawler
	proxy    proxy.Proxier
}

func NewPool(proxier proxy.Proxier) *Pool {
	return &Pool{
		Urls:     make(chan string),
		Proxier:  proxier,
		crawlers: make(chan *Crawler),
	}
}

func (pool *Pool) Crawl(posts chan PaginationIterator) {
	for posts := range posts {
		crawler := pool.Crawler()
		go crawler.Crawl(posts)
	}
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

func (pool *Pool) Return(crawler *Crawler) {
	pool.crawlers <- crawler
}

func (pool *Pool) UrlStream() chan string {
	return pool.Urls
}
