package proxy

import (
	"github.com/rentapplication/craigjr/web"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	inCloakList  = "http://incloak.com/api/proxylist.txt?out=plain&lang=en"
	inCloakLogin = "http://incloak.com/login"
)

type InCloakLoginError struct {
}

func (icle InCloakLoginError) Error() string {
	return "InCloak Login Failed"
}

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
	err := inCloakLoginPost()
	if err == nil {
		return inCloakQuery()
	} else {
		return &web.Query{Error: err}
	}
}

func inCloakLoginPost() error {
	response, err := http.PostForm(inCloakLogin, url.Values{"c": {os.Getenv("IN_CLOAK_API_KEY")}})
	if err == nil && response.StatusCode >= 400 {
		return InCloakLoginError{}
	}
	return err
}

func inCloakQuery() *web.Query {
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
