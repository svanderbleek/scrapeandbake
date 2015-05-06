package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const PROXY_TIMEOUT = 20 * time.Second

type BadResponseError int

func (bre BadResponseError) Error() string {
	return fmt.Sprintf("Response code: %v", int(bre))
}

type Proxy struct {
	*http.Client
	Url    string
	Errors int
}

func (p *Proxy) String() string {
	return fmt.Sprintf("Proxy %v Errors %v", p.Url, p.Errors)
}

func NewProxy(proxyUrl string) *Proxy {
	proxy := proxy(proxyUrl)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	return &Proxy{client, proxyUrl, 0}
}

func proxy(proxyUrl string) func(*http.Request) (*url.URL, error) {
	url, _ := url.Parse(proxyUrl)
	return http.ProxyURL(url)
}

func (proxy *Proxy) getBody(url string) (string, error) {
	body, err := proxy.getBodyStringWithError(url)
	if err != nil {
		proxy.Errors++
	}
	return body, err
}

func (proxy *Proxy) getBodyStringWithError(url string) (string, error) {
	var body []byte
	response, err := proxy.Get(url)
	if err == nil {
		body, err = readBodyAndStatus(response)
	}
	return string(body), err
}

func readBodyAndStatus(response *http.Response) ([]byte, error) {
	status := response.StatusCode
	if status == 200 {
		defer response.Body.Close()
		return ioutil.ReadAll(response.Body)
	} else {
		return nil, BadResponseError(status)
	}
}
