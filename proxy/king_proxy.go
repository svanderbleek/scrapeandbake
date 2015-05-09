package proxy

import (
	"fmt"
	"github.com/rentapplication/craigjr/web"
	"os"
)

type KingProxy struct{}

func (kp KingProxy) String() string {
	return "KingProxy.com"
}

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

func (kp KingProxy) Result() *Result {
	result := &Result{}
	query := kp.Query()
	if query.Error == nil {
		queryResult := query.Result.(*kingProxyResult)
		result.Proxies = kp.Proxies(queryResult)
	} else {
		result.Error = query.Error
	}
	return result
}

func (kp KingProxy) Query() *web.Query {
	query := &web.Query{
		Url:    kingProxyUrl(),
		Result: &kingProxyResult{},
	}
	query.FetchJson()
	return query
}

const kingProxyList = "http://kingproxies.com/api/v1/proxies.json?supports=craigslist&limit=1000&key=%v"

func kingProxyUrl() string {
	return fmt.Sprintf(kingProxyList, os.Getenv("KING_PROXY_API_KEY"))
}

func (kp KingProxy) Proxies(result *kingProxyResult) []*Proxy {
	var proxies []*Proxy
	for _, proxy := range result.Data.Proxies {
		proxies = append(proxies, kp.Proxy(proxy))
	}
	return proxies
}

const kingProxyHost = "http://%v:%v"

func (kp KingProxy) Proxy(proxy kingProxyInfo) *Proxy {
	host := fmt.Sprintf(kingProxyHost, proxy.Ip, proxy.Port)
	return NewProxy(host)
}
