package crawler

type PaginationIterator interface {
	Next() bool
	Items() []string
}

type Queue interface {
	Push(string)
}

func Crawl(pages PaginationIterator, queue Queue) {
	for pages.Next() {
		for _, item := range pages.Items() {
			queue.Push(item)
		}
	}
}
