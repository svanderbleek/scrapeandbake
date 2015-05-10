package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type Proxy struct {
	*http.Client
	Url       string
	Errors    int
	Successes int
}

func (p *Proxy) String() string {
	return fmt.Sprintf("Proxy %v Errors %v Successes %v", p.Url, p.Errors, p.Successes)
}

func NewProxy(url string) *Proxy {
	proxy := proxyUrl(url)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	return &Proxy{Client: client, Url: url}
}

func proxyUrl(host string) func(*http.Request) (*url.URL, error) {
	proxy, _ := url.Parse(host)
	return http.ProxyURL(proxy)
}

func (proxy *Proxy) getBody(url string) (string, error) {
	body, err := proxy.getBodyStringWithError(url)
	if err != nil {
		proxy.Errors++
	} else {
		proxy.Successes++
	}
	return body, err
}

func (proxy *Proxy) getBodyStringWithError(url string) (string, error) {
	var body []byte
	response, err := proxy.Get(url)
	if err == nil {
		body, err = proxy.readBody(response)
	}
	return string(body), err
}

type BadResponseError int

func (bre BadResponseError) Error() string {
	return fmt.Sprintf("Response code: %v", int(bre))
}

func (proxy *Proxy) readBody(response *http.Response) ([]byte, error) {
	status := response.StatusCode
	if status == 200 {
		defer response.Body.Close()
		return ioutil.ReadAll(response.Body)
	} else {
		return nil, BadResponseError(status)
	}
}

type ProxyBlockedError struct{}

func (pbe ProxyBlockedError) Error() string {
	return "Proxy blocked"
}

var blockedMessage = regexp.MustCompile(`This IP has been automatically blocked`)

func (proxy *Proxy) isBlocked(body string) bool {
	return blockedMessage.MatchString(body)
}
