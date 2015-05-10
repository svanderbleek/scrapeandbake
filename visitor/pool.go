package visitor

import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/proxy"
)

type Pool struct {
	Posts    chan *craigslist.Post
	visitors chan *Visitor
	proxy    proxy.Proxier
}

func NewPool(proxier proxy.Proxier) *Pool {
	return &Pool{
		Posts:    make(chan *craigslist.Post),
		proxy:    proxier,
		visitors: make(chan *Visitor),
	}
}

func (pool *Pool) Visit(urls <-chan string) {
	go pool.visitUrls(urls)
}

func (pool *Pool) visitUrls(urls <-chan string) {
	url := <-urls
	visitor := pool.Visitor(url)
	go visitor.Visit(url)
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
