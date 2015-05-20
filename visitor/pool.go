package visitor

import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/proxy"
	. "github.com/tj/go-debug"
	"time"
)

var debug = Debug("pool")

const (
	INIT_POOL = 100
	MAX_WAIT  = 60 * time.Second
)

type Pool struct {
	proxy.Proxier
	Posts    chan *craigslist.Post
	Count    int
	visitors chan *Visitor
}

func NewPool(proxier proxy.Proxier) *Pool {
	pool := &Pool{
		Proxier:  proxier,
		Posts:    make(chan *craigslist.Post),
		visitors: make(chan *Visitor, INIT_POOL),
	}
	pool.init()
	return pool
}

func (pool *Pool) init() {
	for i := 0; i < INIT_POOL; i++ {
		pool.visitors <- NewVisitor(pool)
		pool.Count++
	}
}

func (pool *Pool) Return(post *craigslist.Post, visitor *Visitor) {
	pool.Posts <- post
	pool.visitors <- visitor
}

func (pool *Pool) Visit(urls <-chan string) {
	for url := range urls {
		visitor := pool.Visitor(url)
		go visitor.Visit(url)
	}
}

func (pool *Pool) Visitor(url string) *Visitor {
	var visitor *Visitor
	select {
	case visitor = <-pool.visitors:
	case <-time.After(MAX_WAIT):
		pool.Count++
		debug("Visitor count is %v", pool.Count)
		visitor = NewVisitor(pool)
	}
	return visitor
}
