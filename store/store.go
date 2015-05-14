package store

import (
	elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/rentapplication/craigjr/craigslist"
	"net/url"
	"os"
)

type Storer interface {
	Store(*craigslist.Post)
	IsStored(*craigslist.Post) bool
}

type ElasticSearch struct {
	*elastigo.Conn
	Errors int
}

func New() Storer {
	connection := elastigo.NewConn()
	es := &ElasticSearch{Conn: connection}
	host, _ := url.Parse(os.Getenv("ELASTIC_SEARCH_URL"))
	es.configure(host)
	return es
}

func (es *ElasticSearch) configure(u *url.URL) {
	es.Protocol = u.Scheme
	es.Username = u.User.Username()
	es.Password, _ = u.User.Password()
	es.SetHosts([]string{u.Host})
}

func (es *ElasticSearch) Store(post *craigslist.Post) {
	_, err := es.Index("craigjr", "post", "", nil, post)
	if err != nil {
		es.Errors++
	}
}

func (es *ElasticSearch) IsStored(post *craigslist.Post) bool {
	idQuery := "_id:" + post.Url
	query := map[string]interface{}{"q": idQuery}
	result, err := es.SearchUri("craigjr", "post", query)
	if err == nil {
		return result.Hits.Total > 0
	} else {
		es.Errors++
		return false
	}
}
