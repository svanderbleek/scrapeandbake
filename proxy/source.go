package proxy

import (
	"fmt"
)

type ProxyStream chan *Proxy

type ProxySource interface {
	fmt.Stringer
	Result() *Result
}

type Result struct {
	Proxies []*Proxy
	Error   error
	Source  fmt.Stringer
}

func Fetch(stream ProxyStream, source ProxySource) {
	result := source.Result()
	result.Source = source
	stream.Load(result)
}

func (stream ProxyStream) Load(result *Result) {
	if result.Error == nil {
		for _, proxy := range result.Proxies {
			stream <- proxy
			debug("Proxy loaded %v from %v", proxy, result.Source)
		}
	} else {
		debug("Proxies failed to load from %v", result.Source)
	}
}
