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
}

func tryLoadProxies(sink proxySink, source proxySource) {
	result := source()
	if result.Error == nil {
		loadProxies(sink, result.Proxies)
		log.Printf("%v Proxies loaded", len(result.Proxies))
	} else {
		log.Printf("Proxies failed to load")
		log.Println(result.Error.Error())
	}
}

func loadProxies(sink proxySink, proxies []*Proxy) {
	for _, proxy := range proxies {
		sink <- proxy
	}
}
