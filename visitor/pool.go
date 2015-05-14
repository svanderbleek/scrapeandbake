package visitor

import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/proxy"
	. "github.com/tj/go-debug"
)

var debug = Debug("pool")

type Pool struct {
	proxy.Proxier
	Posts    chan *craigslist.Post
	Count    int
	Visited  int
	visitors chan *Visitor
}

func NewPool(proxier proxy.Proxier) *Pool {
	return &Pool{
		Proxier:  proxier,
		Posts:    make(chan *craigslist.Post),
		visitors: make(chan *Visitor),
	}
}

func (pool *Pool) Return(post *craigslist.Post, visitor *Visitor) {
	pool.Posts <- post
	pool.visitors <- visitor
}

func (pool *Pool) Visit(urls <-chan string) {
	for url := range urls {
		pool.Visited++
		visitor := pool.Visitor(url)
		go visitor.Visit(url)
	}
}

func (pool *Pool) Visitor(url string) *Visitor {
	var visitor *Visitor
	select {
	case visitor = <-pool.visitors:
	default:
		pool.Count++
		debug("Visitor count is %v", pool.Count)
		visitor = NewVisitor(pool)
	}
	return visitor
}
