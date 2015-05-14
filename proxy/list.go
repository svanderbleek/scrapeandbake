package proxy

import (
	"container/heap"
	. "github.com/tj/go-debug"
	"sync"
	"time"
)

var debug = Debug("proxy")

type Proxier interface {
	MustGet(url string) string
}

type List struct {
	Proxies []*Proxy
	sync.Mutex
	in ProxyStream
}

func NewList() *List {
	list := &List{
		in: make(ProxyStream),
	}
	heap.Init(list)
	return list
}

// Asynchronously load proxies from sources
func (list *List) Load(sources ...ProxySource) {
	for _, source := range sources {
		go Fetch(list.in, source)
	}
}

func (list *List) LoadDefault() {
	list.Load(KingProxy{}, InCloak{})
}

// Borrow proxies until successful response
func (list *List) MustGet(url string) string {
	response := &Response{Error: EmptyResponseError{}}
	for response.Error != nil {
		response = list.Get(url)
	}
	return response.Body
}

type Response struct {
	Url   string
	Body  string
	Error error
}

type EmptyResponseError struct{}

func (ere EmptyResponseError) Error() string {
	return "Response is empty"
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
	proxy := list.Borrow()
	debug("Get url %v with Proxy %v", url, proxy)
	body, err := proxy.getBody(url)
	debug("Got url %v with Proxy %v and err %v", url, proxy, err)
	if err == nil {
		err = list.discardProxyIfBlocked(proxy, body)
	} else {
		list.Return(proxy)
	}
	return body, err
}

func (list *List) discardProxyIfBlocked(proxy *Proxy, body string) error {
	if proxy.isBlocked(body) {
		debug("Proxy %v is blocked, discarding", proxy)
		return ProxyBlockedError{}
	} else {
		list.Return(proxy)
		return nil
	}
}

func (list *List) Borrow() *Proxy {
	proxy := list.remove()
	if proxy == nil {
		debug("Waiting for proxy")
		proxy = list.waitForProxy()
	}
	debug("Borrowing proxy %v", proxy)
	return proxy
}

func (list *List) remove() *Proxy {
	var proxy *Proxy
	list.Lock()
	if list.Len() > 0 {
		proxy = heap.Pop(list).(*Proxy)
	}
	list.Unlock()
	return proxy
}

const PROXY_WAIT_TIME = 5 * time.Second

func (list *List) waitForProxy() *Proxy {
	var proxy *Proxy
	for proxy == nil {
		select {
		case proxy = <-list.in:
		case <-time.After(PROXY_WAIT_TIME):
			proxy = list.remove()
		}
	}
	return proxy
}

const MAX_PROXY_ERRORS = 2

func (list *List) Return(proxy *Proxy) {
	if proxy.Errors < MAX_PROXY_ERRORS {
		debug("Returning proxy %v", proxy)
		list.add(proxy)
	} else {
		debug("Failing proxy   %v", proxy)
	}
}

func (list *List) add(proxy *Proxy) {
	list.Lock()
	heap.Push(list, proxy)
	list.Unlock()
}

// Heap interface{} implementation

func (list *List) Pop() interface{} {
	index := list.Len() - 1
	popped := list.Proxies[index]
	list.Proxies = list.Proxies[:index]
	return popped
}

func (list *List) Push(proxy interface{}) {
	list.Proxies = append(list.Proxies, proxy.(*Proxy))
}

// Sort interface{} implementation

func (list *List) Len() int {
	return len(list.Proxies)
}

func (list *List) Less(in, jn int) bool {
	i := list.Proxies[in]
	j := list.Proxies[jn]
	return j.Successes <= i.Successes
}

func (list *List) Swap(i, j int) {
	list.Proxies[i], list.Proxies[j] = list.Proxies[j], list.Proxies[i]
}
