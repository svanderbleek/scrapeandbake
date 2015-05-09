package crawler

type PaginationIterator interface {
	Next() bool
	Items() []string
}

type UrlStream chan string

type Crawler struct {
	urls UrlStream
	done chan int
}

func New(urls UrlStream) *Crawler {
	return &Crawler{
		done: make(chan int, 1),
		urls: urls,
	}
}

func (crawler *Crawler) Crawl(pages PaginationIterator) {
	for pages.Next() {
		for _, item := range pages.Items() {
			crawler.urls <- item
		}
	}
	crawler.Done()
}

func (crawler *Crawler) Done() {
	crawler.done <- 1
}

func (crawler *Crawler) WaitDone() {
	<-crawler.done
}
