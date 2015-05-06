package proxy

import (
	"fmt"
	"github.com/rentapplication/craigjr/web"
	"os"
)

const (
	kingProxyList = "http://kingproxies.com/api/v1/proxies.json?supports=craigslist&limit=1000&key=%v"
	kingProxyHost = "http://%v:%v"
)

type kingProxyInfo struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type kingProxyResult struct {
	Data struct {
		Proxies []kingProxyInfo `json:"proxies"`
	} `json:"data"`
	Error error
}

func kingProxy() *Result {
	result := &Result{}
	query := kingProxyFetch()
	if query.Error == nil {
		result.Proxies = kingProxyBuildProxies(query.Result.(*kingProxyResult))
	} else {
		result.Error = query.Error
	}
	return result
}

func kingProxyFetch() *web.Query {
	query := &web.Query{
		Url:    kingProxyUrl(),
		Result: &kingProxyResult{},
	}
	query.Fetch()
	return query
}

func kingProxyUrl() string {
	return fmt.Sprintf(kingProxyList, os.Getenv("KING_PROXY_API_KEY"))
}

func kingProxyBuildProxies(result *kingProxyResult) []*Proxy {
	var proxies []*Proxy
	for _, proxy := range result.Data.Proxies {
		proxies = append(proxies, kingProxyBuildProxy(proxy))
	}
	return proxies
}

func kingProxyBuildProxy(proxy kingProxyInfo) *Proxy {
	host := fmt.Sprintf(kingProxyHost, proxy.Ip, proxy.Port)
	return NewProxy(host)
}
