package proxy

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var (
	list = NewList()
)

type FakeProxySource struct {
	stop chan int
}

func (fps *FakeProxySource) Fetch(stream ProxyStream) {
	defer func() {
		recover()
	}()
	for {
		stream <- &Proxy{}
		randomWork(10)
	}
}

func runProxyConsumers(consumers int) {
	var wg sync.WaitGroup
	wg.Add(consumers)
	for ; consumers > 0; consumers-- {
		go func() {
			defer wg.Done()
			useProxyAttempts(10)
		}()
	}
	wg.Wait()
}

func useProxyAttempts(attempts int) {
	for ; attempts > 0; attempts-- {
		proxy := list.Borrow()
		randomWork(10)
		randomSuccessOrError(proxy)
		list.Return(proxy)
	}
}

func randomWork(length int) {
	duration := time.Duration(rand.Intn(length))
	time.Sleep(duration * time.Millisecond)
}

func randomSuccessOrError(proxy *Proxy) {
	if rand.Intn(2) > 0 {
		proxy.Successes++
	} else {
		proxy.Errors++
	}
}

func TestAsyncAccess(t *testing.T) {
	go list.Load(&FakeProxySource{})
	runProxyConsumers(100)
	close(list.in)
	proxy := list.Borrow()
	for list.Len() > 0 {
		next := list.Borrow()
		log.Printf("Comparing proxy %v with %v", proxy, next)
		if proxy.Errors > next.Errors || (proxy.Errors == next.Errors && proxy.Successes < next.Successes) {
			t.Errorf("Invalid Ordering of %v before %v", proxy, next)
		}
		proxy = next
	}
}
