package visitor

import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/proxy"
	. "github.com/tj/go-debug"
	"time"
)

var debug = Debug("pool")

const (
	MAX_WAIT  = 30 * time.Second
	INIT_POOL = 40
)

type Pool struct {
	proxy.Proxier
	Posts    chan *craigslist.Post
	Count    int
	Visited  int
	visitors chan *Visitor
}

func NewPool(proxier proxy.Proxier) *Pool {
	pool := &Pool{
		Proxier:  proxier,
		Posts:    make(chan *craigslist.Post),
		visitors: make(chan *Visitor, INIT_POOL),
	}
	for i := 0; i < INIT_POOL; i++ {
		pool.visitors <- NewVisitor(pool)
		pool.Count++
	}
	return pool
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
		pool.Visited++
		debug("Visited count is %v", pool.Visited)
	case <-time.After(MAX_WAIT):
		pool.Count++
		debug("Visitor count is %v", pool.Count)
		visitor = NewVisitor(pool)
	}
	return visitor
}
