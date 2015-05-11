package proxy

import (
	"fmt"
	"log"
)

type ProxyStream chan *Proxy

type ProxySource interface {
	Result() *Result
	fmt.Stringer
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
			log.Printf("Proxy loaded %v from %v", proxy, result.Source)
		}
	} else {
		log.Printf("Proxies failed to load from %v", result.Source)
	}
}
