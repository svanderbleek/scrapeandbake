package main

import (
	"./proxy"
	"fmt"
	"sync"
)

func main() {
	proxyList := proxy.NewList()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			response := proxyList.Get("http://atlanta.craigslist.org/search/apa?s=100")
			fmt.Println(string(response.Body))
		}()
	}
	wg.Wait()
}
