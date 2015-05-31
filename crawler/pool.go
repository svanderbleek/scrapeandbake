package crawler

import (
	"github.com/rentapplication/craigjr/proxy"
	. "github.com/tj/go-debug"
)

var debug = Debug("pool")

const (
	INIT_POOL = 100
)

type Pool struct {
	proxy.Proxier
	Urls     chan string
	Count    int
	crawlers chan *Crawler
	proxy    proxy.Proxier
}

func NewPool(proxier proxy.Proxier) *Pool {
	pool := &Pool{
		Urls:     make(chan string),
		Proxier:  proxier,
		crawlers: make(chan *Crawler, INIT_POOL),
	}
	pool.init()
	return pool
}

func (pool *Pool) init() {
	for i := 0; i < INIT_POOL; i++ {
		pool.crawlers <- NewCrawler(pool)
		pool.Count++
	}
}

func (pool *Pool) Crawl(posts chan PaginationIterator) {
	for posts := range posts {
		crawler := pool.Crawler()
		go crawler.Crawl(posts)
	}
}

func (pool *Pool) Crawler() *Crawler {
	return <-pool.crawlers
}

func (pool *Pool) Return(crawler *Crawler) {
	pool.crawlers <- crawler
}

func (pool *Pool) UrlStream() chan string {
	return pool.Urls
}
