package proxy

type ProxyStream chan *Proxy
type ProxyArray []*Proxy

type ProxySource interface {
	Fetch(ProxyStream)
}

func streamArray(stream ProxyStream, array ProxyArray) {
	for _, proxy := range array {
		stream <- proxy
	}
}
