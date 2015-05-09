package proxy

import "log"

type ProxyStream chan *Proxy

type ProxySource interface {
	Result() *Result
}

type Result struct {
	Proxies []*Proxy
	Error   error
	Source  interface{}
}

func Fetch(stream ProxyStream, source ProxySource) {
	result := source.Result()
	stream.Load(result)
}

func (stream ProxyStream) Load(result *Result) {
	if result.Error == nil {
		for _, proxy := range result.Proxies {
			stream <- proxy
			log.Printf("Proxy loaded %v", proxy)
		}
		log.Printf("Proxies loaded from %v", result.Source)
	} else {
		log.Printf("Proxies failed to load from %v", result.Source)
	}
}
