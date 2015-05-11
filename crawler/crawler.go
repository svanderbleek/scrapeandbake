package crawler

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
	urls chan int
}

func NewCrawler(pool Pooler) *Crawler {
	return &Crawler{
		Pooler: pool,
	}
}

func (crawler *Crawler) Crawl(pages PaginationIterator) {
	for url, done := pages.Next(); !done; url, done = pages.Next() {
		body := crawler.MustGet(url)
		for _, item := range pages.Items(body) {
			crawler.UrlStream() <- item
		}
	}
	crawler.Return(crawler)
}
