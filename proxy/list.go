package proxy

import (
	"container/heap"
	"log"
	"regexp"
)

const MAX_PROXY_ERRORS = 2

var (
	blockedMessage    = regexp.MustCompile(`This IP has been automatically blocked`)
	proxyListInstance = NewList()
)

func MustGet(url string) *Response {
	response := &Response{Error: EmptyResponseError{}}
	for response.Error != nil {
		response = proxyListInstance.Get(url)
	}
	return response
}

type ProxyBlockedError struct {
}

func (pbe ProxyBlockedError) Error() string {
	return "Proxy blocked"
}

type EmptyResponseError struct {
}

func (ere EmptyResponseError) Error() string {
	return "Response is empty"
}

type Response struct {
	Url   string
	Body  string
	Error error
}

type List struct {
	Proxies []*Proxy
	load    chan *Proxy
}

func NewList() *List {
	list := &List{
		load: make(chan *Proxy),
	}
	go Load(list.load)
	return list
}

func (list *List) Get(url string) *Response {
	body, err := list.getBodyWithProxy(url)
	return &Response{
		Url:   url,
		Body:  body,
		Error: err,
	}
}

func (list *List) getBodyWithProxy(url string) (string, error) {
	log.Printf("Proxies status: %v", list.Proxies)
	proxy := list.Borrow()
	log.Printf("Get url %v with Proxy %v", url, proxy)
	body, err := proxy.getBody(url)
	log.Printf("Got url %v with Proxy %v and err %v", url, proxy, err)
	if err == nil {
		err = list.discardProxyIfBlocked(proxy, body)
	} else {
		list.Return(proxy)
	}
	return body, err
}

func (list *List) discardProxyIfBlocked(proxy *Proxy, body string) error {
	if !proxyIsBlocked(body) {
		list.Return(proxy)
		return nil
	} else {
		log.Printf("Proxy %v is blocked, discarding", proxy)
		return ProxyBlockedError{}
	}
}

func proxyIsBlocked(body string) bool {
	return blockedMessage.MatchString(body)
}

func (list *List) Borrow() *Proxy {
	var proxy *Proxy
	if list.Len() > 0 {
		proxy = heap.Pop(list).(*Proxy)
	} else {
		log.Printf("Waiting for proxy")
		proxy = <-list.load
	}
	log.Printf("Borrowing proxy %v", proxy)
	return proxy
}

func (list *List) Return(proxy *Proxy) {
	if proxy.Errors < MAX_PROXY_ERRORS {
		log.Printf("Returning proxy %v", proxy)
		heap.Push(list, proxy)
	}
}

func (list *List) Pop() interface{} {
	index := list.Len() - 1
	popped := list.Proxies[index]
	list.Proxies = list.Proxies[:index]
	return popped
}

func (list *List) Push(proxy interface{}) {
	list.Proxies = append(list.Proxies, proxy.(*Proxy))
}

func (list *List) Len() int {
	return len(list.Proxies)
}

func (list *List) Less(i, j int) bool {
	return list.Proxies[i].Errors < list.Proxies[j].Errors
}

func (list *List) Swap(i, j int) {
	list.Proxies[i], list.Proxies[j] = list.Proxies[j], list.Proxies[i]
}
