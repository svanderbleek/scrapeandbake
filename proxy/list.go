package proxy

import (
	"container/heap"
	"log"
	"sync"
)

type List struct {
	Proxies []*Proxy
	state   *sync.Mutex
	in      ProxyStream
}

var defaultList = NewList()

func LoadDefault() {
	defaultList.Load(KingProxy{}, InCloak{}, ProxIsRight{})
}

func NewList() *List {
	list := &List{
		in:    make(ProxyStream),
		state: &sync.Mutex{},
	}
	heap.Init(list)
	return list
}

// Asynchronous Proxy loader
func (list *List) Load(sources ...ProxySource) {
	for _, source := range sources {
		go Fetch(list.in, source)
	}
}

// Borrows Proxies until successful Response
func MustGet(url string) *Response {
	response := &Response{Error: EmptyResponseError{}}
	for response.Error != nil {
		response = defaultList.Get(url)
	}
	return response
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
	if proxy.isBlocked(body) {
		log.Printf("Proxy %v is blocked, discarding", proxy)
		return ProxyBlockedError{}
	} else {
		list.Return(proxy)
		return nil
	}
}

func (list *List) Borrow() *Proxy {
	var proxy *Proxy
	if list.Len() > 0 { // TODO race condition need to wrap in lock
		proxy = list.lockedRemove()
		log.Printf("Borrowing proxy %v", proxy)
	} else {
		log.Printf("Waiting for proxy")
		proxy = <-list.in
	}
	return proxy
}

func (list *List) lockedRemove() *Proxy {
	list.state.Lock()
	proxy := heap.Pop(list).(*Proxy)
	list.state.Unlock()
	return proxy
}

const MAX_PROXY_ERRORS = 2

func (list *List) Return(proxy *Proxy) {
	if proxy.Errors < MAX_PROXY_ERRORS {
		log.Printf("Returning proxy %v", proxy)
		list.lockedAdd(proxy)
	} else {
		log.Printf("Failing proxy   %v", proxy)
	}
}

func (list *List) lockedAdd(proxy *Proxy) {
	list.state.Lock()
	heap.Push(list, proxy)
	list.state.Unlock()
}

// Heap interface{} Implementation
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

func (list *List) Less(in, jn int) bool {
	i := list.Proxies[in]
	j := list.Proxies[jn]
	return j.Successes <= i.Successes
}

func (list *List) Swap(i, j int) {
	list.Proxies[i], list.Proxies[j] = list.Proxies[j], list.Proxies[i]
}
