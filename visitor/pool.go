package visitor

import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/proxy"
)

type Pool struct {
	proxy.Proxier
	Posts    chan *craigslist.Post
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
		visitor := pool.Visitor(url)
		go visitor.Visit(url)
	}
}

func (pool *Pool) Visitor(url string) *Visitor {
	var visitor *Visitor
	select {
	case visitor = <-pool.visitors:
	default:
		visitor = NewVisitor(pool)
	}
	return visitor
}
