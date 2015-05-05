package proxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Proxy struct {
	*http.Client
}

func NewProxy(proxyUrl string) *Proxy {
	proxy := proxy(proxyUrl)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	return &Proxy{client}
}

func proxy(proxyUrl string) func(*http.Request) (*url.URL, error) {
	url, _ := url.Parse(proxyUrl)
	return http.ProxyURL(url)
}

func (p *Proxy) getBody(url string) ([]byte, error) {
	var body []byte
	response, err := p.Get(url)
	if err == nil {
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
	}
	return body, err
}
