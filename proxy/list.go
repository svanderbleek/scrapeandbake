package proxy

import "net/http"

type List struct {
	Proxies []Proxy
}

type Request func(proxy string) (*http.Response, error)

func (list *List) Proxy(request Request) Proxy {
	return Proxy{}
}
