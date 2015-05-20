package store

import (
	elastigo "github.com/mattbaird/elastigo/lib"
	. "github.com/tj/go-debug"
	"net/url"
	"os"
)

var debug = Debug("store")

type Storer interface {
	Store(Storable)
	IsStored(string) bool
}

type Storable interface {
	Id() string
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

func (es *ElasticSearch) Store(item Storable) {
	_, err := es.Index("craigjr", "post", item.Id(), nil, item)
	if err != nil {
		es.Errors++
		debug("ElasticSearch error %v", err)
	}
}

func (es *ElasticSearch) IsStored(id string) bool {
	idQuery := "_id:" + id
	query := map[string]interface{}{"q": idQuery}
	result, err := es.SearchUri("craigjr", "post", query)
	if err == nil {
		return result.Hits.Total > 0
	} else {
		es.Errors++
		return false
	}
}
