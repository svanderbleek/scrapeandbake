package proxy

const (
	MAX_PROXIES = 200
)

type Response struct {
	Url  string
	Body []byte
	Err  error
}

type List struct {
	Proxies chan *Proxy
}

func NewList() *List {
	list := &List{
		Proxies: make(chan *Proxy, MAX_PROXIES),
	}
	go Load(list.Proxies)
	return list
}

func (l *List) Get(url string) *Response {
	proxy := l.borrowProxy()
	body, err := proxy.getBody(url)
	l.returnProxy(proxy)
	return &Response{
		Url:  url,
		Body: body,
		Err:  err,
	}
}

func (l *List) borrowProxy() *Proxy {
	return <-l.Proxies
}

func (l *List) returnProxy(proxy *Proxy) {
	l.Proxies <- proxy
}
