package crawler

import (
	"testing"
)

type FakePaginationIterator struct {
	page int
	done chan int
}

const PAGE_LIMIT = 10

func (fpi *FakePaginationIterator) Next() bool {
	fpi.page++
	if fpi.page < PAGE_LIMIT {
		return true
	} else {
		return false
	}
}

func (fpi *FakePaginationIterator) Items() []string {
	return []string{"test"}
}

func TestCrawlers(t *testing.T) {
	urls := make(chan string, 1000)
	var crawlers []*Crawler
	var iterators []*FakePaginationIterator
	for i := 0; i < 100; i++ {
		pages := &FakePaginationIterator{}
		crawler := New(urls)
		crawlers = append(crawlers, crawler)
		iterators = append(iterators, pages)
		go crawler.Crawl(pages)
	}
	for _, crawler := range crawlers {
		crawler.WaitDone()
	}
	for _, iterator := range iterators {
		got := iterator.page
		if got != PAGE_LIMIT {
			t.Errorf("Got %v Wanted %v", got, PAGE_LIMIT)
		}
	}
}
