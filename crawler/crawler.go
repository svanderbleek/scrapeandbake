package crawler

import "fmt"

type PaginationIterator interface {
	Next() (string, bool)
	Items(string) []string
}

type Crawler struct {
	pool *Pool
}

func NewCrawler(pool *Pool) *Crawler {
	return &Crawler{pool}
}

func (crawler *Crawler) Crawl(pages PaginationIterator) {
	fmt.Println("????????")
	for url, done := pages.Next(); !done; url, done = pages.Next() {
		body := crawler.pool.proxy.MustGet(url)
		for _, item := range pages.Items(body) {
			crawler.pool.Urls <- item
		}
	}
	crawler.pool.crawlers <- crawler
}
