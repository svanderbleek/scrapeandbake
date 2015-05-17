package proxy

import (
	"github.com/rentapplication/craigjr/web"
	"os"
)

type ProxIsRight struct{}

func (pis ProxIsRight) String() string {
	return "TheProxIsRight.com"
}

func (pis ProxIsRight) Result() *Result {
	result := &Result{}
	query := pis.Query()
	if query.Error == nil {
		queryResult := query.Result.(*proxIsRightResult)
		result.Proxies = pis.Proxies(queryResult)
	} else {
		result.Error = query.Error
	}
	return result
}

type proxIsRightInfo struct {
	Host string `json:"host"`
}

type proxIsRightResult struct {
	Proxies []proxIsRightInfo `json:"list"`
	Error   error
}

const proxIsRightList = "https://proxies.p.mashape.com/api/proxy/get?onlyActive=true&onlySupportsCraigslist=true"

func (pis ProxIsRight) Query() *web.Query {
	query := &web.Query{
		Url:    proxIsRightList,
		Header: &web.Header{"X-Mashape-Key", os.Getenv("MASHAPE_API_KEY")},
		Result: &proxIsRightResult{},
	}
	query.FetchJson()
	return query
}

func (pis ProxIsRight) Proxies(result *proxIsRightResult) []*Proxy {
	var proxies []*Proxy
	for _, proxy := range result.Proxies {
		proxies = append(proxies, NewProxy("http://"+proxy.Host))
	}
	return proxies
}
