package proxy

import (
	"log"
)

type proxySink chan *Proxy
type proxySource func() *Result

type Result struct {
	Proxies []*Proxy
	Error   error
}

func Load(sink proxySink) {
	go tryLoadProxies(sink, kingProxy)
	go tryLoadProxies(sink, proxIsRight)
	go tryLoadProxies(sink, inCloak)
}

func tryLoadProxies(sink proxySink, source proxySource) {
	result := source()
	if result.Error == nil {
		log.Printf("Proxies loaded from %v", source)
		loadProxies(sink, result.Proxies)
	} else {
		log.Printf("Proxies failed to load from %v", source)
		log.Println(result.Error.Error())
	}
}

func loadProxies(sink proxySink, proxies []*Proxy) {
	for _, proxy := range proxies {
		sink <- proxy
		log.Printf("Proxy loaded %v", proxy)
	}
}
