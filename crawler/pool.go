package crawler

import (
	"github.com/rentapplication/craigjr/proxy"
	. "github.com/tj/go-debug"
	"time"
)

var debug = Debug("pool")

const (
	INIT_POOL = 10
	MAX_WAIT  = 30 * time.Second
)

type Pool struct {
	proxy.Proxier
	Urls     chan string
	Count    int
	Crawled  int
	crawlers chan *Crawler
	proxy    proxy.Proxier
}

func NewPool(proxier proxy.Proxier) *Pool {
	pool := &Pool{
		Urls:     make(chan string),
		Proxier:  proxier,
		crawlers: make(chan *Crawler, INIT_POOL),
	}
	for i := 0; i < INIT_POOL; i++ {
		pool.crawlers <- NewCrawler(pool)
		pool.Count++
	}
	return pool
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
		pool.Crawled++
		debug("Crawled count is %v", pool.Crawled)
	case <-time.After(MAX_WAIT):
		pool.Count++
		debug("Crawler count is %v", pool.Count)
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
