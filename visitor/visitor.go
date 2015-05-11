package visitor

// TODO Wrong
import (
	"github.com/rentapplication/craigjr/craigslist"
)

type Pooler interface {
	MustGet(string) string
	Return(*craigslist.Post, *Visitor)
}

type Visitor struct {
	Pooler
}

func NewVisitor(pool Pooler) *Visitor {
	return &Visitor{pool}
}

func (visitor *Visitor) Visit(url string) {
	body := visitor.MustGet(url)
	post := craigslist.NewPost(body)
	visitor.Return(post, visitor)
}
