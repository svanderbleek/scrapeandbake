package scraper

type UrlQueue struct {
	queue chan string
}

func NewUrlQueue() *UrlQueue {
	return &UrlQueue{queue: make(chan string)}
}

func (uq *UrlQueue) Push(url string) {
	uq.queue <- url
}

func (uq *UrlQueue) Pop() string {
	return <-uq.queue
}
