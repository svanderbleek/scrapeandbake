package visitor

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
	postBody := visitor.MustGet(url)
	post := craigslist.NewPost(url, postBody)
	if post.IsContactInfo() {
		contactInfoBody := visitor.MustGet(post.ContactInfoUrl())
		post.ContactInfo = contactInfoBody
	}
	visitor.Return(post, visitor)
}
