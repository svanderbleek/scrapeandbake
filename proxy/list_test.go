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

func randomWork(length int) {
	duration := time.Duration(rand.Intn(length))
	time.Sleep(duration * time.Millisecond)
}

func TestAsyncAccess(t *testing.T) {
	go list.Load(&FakeProxySource{})
	runProxyConsumers(100)
	close(list.in)
	assertOrder(t, list)
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
