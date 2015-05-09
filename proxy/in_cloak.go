package proxy

import (
	"github.com/rentapplication/craigjr/web"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type InCloak struct{}

func (ic InCloak) String() string {
	return "InCloak.com"
}

func (ic InCloak) Result() *Result {
	result := &Result{}
	query := ic.Query()
	if query.Error == nil {
		result.Proxies = ic.Proxies(query.Result.(string))
	} else {
		result.Error = query.Error
	}
	return result
}

func (ic InCloak) Query() *web.Query {
	err := ic.Login()
	if err == nil {
		return ic.ListQuery()
	} else {
		return &web.Query{Error: err}
	}
}

const inCloakLogin = "http://incloak.com/login"

type InCloakLoginError struct{}

func (icle InCloakLoginError) Error() string {
	return "InCloak Login Failed"
}

func (ic InCloak) Login() error {
	queryString := url.Values{"c": {os.Getenv("IN_CLOAK_API_KEY")}}
	response, err := http.PostForm(inCloakLogin, queryString)
	if err == nil && response.StatusCode >= 400 {
		return InCloakLoginError{}
	}
	return err
}

const inCloakList = "http://incloak.com/api/proxylist.txt?out=plain&lang=en"

func (ic InCloak) ListQuery() *web.Query {
	query := &web.Query{
		Url: inCloakList,
	}
	query.Fetch()
	return query
}

func (ic InCloak) Proxies(result string) []*Proxy {
	lines := strings.Split(result, "\r\n")
	var proxies []*Proxy
	for _, proxy := range lines {
		proxies = append(proxies, ic.Proxy(proxy))
	}
	return proxies
}

func (ic InCloak) Proxy(proxy string) *Proxy {
	return NewProxy("http://" + proxy)
}
