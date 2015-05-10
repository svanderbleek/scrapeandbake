package visitor

import (
	"github.com/rentapplication/craigjr/craigslist"
)

type Visitor struct {
	pool *Pool
}

func NewVisitor(pool *Pool) *Visitor {
	return &Visitor{pool}
}

func (visitor *Visitor) Visit(url string) {
	body := visitor.pool.proxy.MustGet(url)
	post := craigslist.NewPost(body)
	visitor.pool.Posts <- post
	visitor.pool.visitors <- visitor
}
