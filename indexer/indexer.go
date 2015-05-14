package indexer

import (
	elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/rentapplication/craigjr/craigslist"
	"log"
	"net/url"
	"os"
)

type Indexer interface {
	Index(*craigslist.Post)
	IsIndexed(*craigslist.Post) bool
}

type ElasticSearch struct {
	connection *elastigo.Conn
}

func New() Indexer {
	host, _ := url.Parse(os.Getenv("ELASTIC_SEARCH_URL"))
	connection := elastigo.NewConn()
	configure(connection, host)
	return &ElasticSearch{connection}
}

func configure(c *elastigo.Conn, u *url.URL) {
	c.Protocol = u.Scheme
	c.Username = u.User.Username()
	c.Password, _ = u.User.Password()
	c.SetHosts([]string{u.Host})
}

func (es *ElasticSearch) Index(post *craigslist.Post) {
	response, err := es.connection.Index("craigjr", "post", "", nil, post)
	log.Printf("Response %v Error %v", response, err)
}

func (es *ElasticSearch) IsIndexed(post *craigslist.Post) bool {
	return false
}
