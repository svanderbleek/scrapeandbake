package proxy

import (
	"fmt"
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

func (pis ProxIsRight) Query() *web.Query {
	query := &web.Query{
		Url:    pis.ListUrl(),
		Result: &proxIsRightResult{},
	}
	query.FetchJson()
	return query
}

const proxIsRightList = "https://theproxisright.com/api/proxy/get?onlyActive=true&onlySupportsCraigslist=true&format=json&apiKey=%v"

func (pis ProxIsRight) ListUrl() string {
	return fmt.Sprintf(proxIsRightList, os.Getenv("PROX_IS_RIGHT_API_KEY"))
}

func (pis ProxIsRight) Proxies(result *proxIsRightResult) []*Proxy {
	var proxies []*Proxy
	for _, proxy := range result.Proxies {
		proxies = append(proxies, NewProxy(proxy.Host))
	}
	return proxies
}
