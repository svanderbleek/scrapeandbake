package proxy

import (
	"fmt"
	"github.com/rentapplication/craigjr/web"
	"os"
)

const (
	proxIsRightList = "https://theproxisright.com/api/proxy/get?onlyActive=true&onlySupportsCraigslist=true&format=json&apiKey=%v"
)

type proxIsRightInfo struct {
	Host string `json:"host"`
}

type proxIsRightResult struct {
	Proxies []proxIsRightInfo `json:"list"`
	Error   error
}

func proxIsRight() *Result {
	result := &Result{}
	query := proxIsRightFetch()
	if query.Error == nil {
		result.Proxies = proxIsRightBuildProxies(query.Result.(*proxIsRightResult))
	} else {
		result.Error = query.Error
	}
	return result
}

func proxIsRightFetch() *web.Query {
	query := &web.Query{
		Url:    proxIsRightUrl(),
		Result: &proxIsRightResult{},
	}
	query.Fetch()
	return query
}

func proxIsRightUrl() string {
	return fmt.Sprintf(proxIsRightList, os.Getenv("PROX_IS_RIGHT_API_KEY"))
}

func proxIsRightBuildProxies(result *proxIsRightResult) []*Proxy {
	var proxies []*Proxy
	for _, proxy := range result.Proxies {
		proxies = append(proxies, NewProxy(proxy.Host))
	}
	return proxies
}
