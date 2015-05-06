package proxy

import (
	"github.com/rentapplication/craigjr/web"
	"strings"
)

const (
	inCloakList = "http://incloak.com/api/proxylist.txt?out=plain&lang=en"
)

func inCloak() *Result {
	result := &Result{}
	query := inCloakFetch()
	if query.Error == nil {
		result.Proxies = inCloakBuildProxies(query.Result.(string))
	} else {
		result.Error = query.Error
	}
	return result
}

func inCloakFetch() *web.Query {
	query := &web.Query{
		Url: inCloakList,
	}
	query.Fetch()
	return query
}

func inCloakBuildProxies(result string) []*Proxy {
	lines := strings.Split(result, "\r\n")
	var proxies []*Proxy
	for _, proxy := range lines {
		proxies = append(proxies, inCloakBuildProxy(proxy))
	}
	return proxies
}

func inCloakBuildProxy(proxy string) *Proxy {
	return NewProxy("http://" + proxy)
}
