package proxy

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestAsyncAccess(t *testing.T) {
	list := NewList()
	runProxySource(list)
	runProxyConsumers(list, 100)
	assertOrder(t, list)
}

type FakeProxySource struct{}

func (fps FakeProxySource) Result() *Result {
	randomWork(20)
	var proxies []*Proxy
	for i := 0; i < 1000; i++ {
		proxies = append(proxies, &Proxy{})
	}
	return &Result{Proxies: proxies}
}

func randomWork(length int) {
	duration := time.Duration(rand.Intn(length))
	time.Sleep(duration * time.Millisecond)
}

func runProxySource(list *List) {
	go list.Load(FakeProxySource{})
}

func runProxyConsumers(list *List, consumers int) {
	var wg sync.WaitGroup
	wg.Add(consumers)
	for ; consumers > 0; consumers-- {
		go func() {
			defer wg.Done()
			useProxyAttempts(list, 10)
		}()
	}
	wg.Wait()
}

func useProxyAttempts(list *List, attempts int) {
	for ; attempts > 0; attempts-- {
		proxy := list.Borrow()
		randomWork(10)
		randomSuccessOrError(proxy)
		list.Return(proxy)
	}
}

const PERCENT_FAILURE = 0.30

func randomSuccessOrError(proxy *Proxy) {
	if rand.Float32() > PERCENT_FAILURE {
		proxy.Successes++
	} else {
		proxy.Errors++
	}
}

func assertOrder(t *testing.T, list *List) {
	i := list.Borrow()
	for list.Len() > 0 {
		j := list.Borrow()
		log.Printf("Comparing proxy %v with %v", i, j)
		if !(i.Successes >= j.Successes && i.Errors < j.Errors+MAX_PROXY_ERRORS) {
			t.Errorf("Invalid Ordering of %v before %v", i, j)
		}
		i = j
	}
}
