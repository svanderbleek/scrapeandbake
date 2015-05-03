package proxy

import (
	"net/http"
	"net/url"
)

type Proxy struct {
	*http.Client
}

func New(proxyUrl string) Proxy {
	proxy := proxy(proxyUrl)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	return Proxy{client}
}

func proxy(proxyUrl string) func(*http.Request) (*url.URL, error) {
	url, _ := url.Parse(proxyUrl)
	return http.ProxyURL(url)
}
