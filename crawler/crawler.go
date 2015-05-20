package crawler

import "github.com/rentapplication/craigjr/store"

type PaginationIterator interface {
	Next() (string, bool)
	Items(string) []string
}

type Pooler interface {
	MustGet(string) string
	Return(*Crawler)
	UrlStream() chan string
}

type Crawler struct {
	Pooler
	store.Storer
	urls chan int
}

func NewCrawler(pool Pooler) *Crawler {
	return &Crawler{
		Pooler: pool,
		Storer: store.New(),
	}
}

func (crawler *Crawler) Crawl(pages PaginationIterator) {
	for url, done := pages.Next(); !done; url, done = pages.Next() {
		crawler.streamItems(pages, url)
	}
	crawler.Return(crawler)
}

func (crawler *Crawler) streamItems(pages PaginationIterator, url string) {
	body := crawler.MustGet(url)
	for _, item := range pages.Items(body) {
		if !crawler.IsStored(item) {
			crawler.UrlStream() <- item
		}
	}
}
