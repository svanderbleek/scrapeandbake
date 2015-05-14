package crawler

import (
	"github.com/rentapplication/craigjr/proxy"
	. "github.com/tj/go-debug"
	"time"
)

var debug = Debug("pool")

const (
	MAX_POOL = 20
	MAX_WAIT = 30 * time.Second
)

type Pool struct {
	proxy.Proxier
	Urls     chan string
	Count    int
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
	case <-time.After(MAX_WAIT):
		if pool.Count < MAX_POOL {
			pool.Count++
			debug("Crawler count is %v", pool.Count)
			crawler = NewCrawler(pool)
		} else {
			crawler = pool.Crawler()
		}
	}
	return crawler
}

func (pool *Pool) Return(crawler *Crawler) {
	pool.crawlers <- crawler
}

func (pool *Pool) UrlStream() chan string {
	return pool.Urls
}
