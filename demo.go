package main

import (
	"fmt"
	"github.com/rentapplication/craigjr/proxy"
)

func main() {
	proxyList := proxy.NewList()
	responses := make(chan string)

	for i := 0; i < 4; i++ {
		go func() {
			for {
				response := proxyList.Get("http://atlanta.craigslist.org/search/apa?s=100")
				if response.Error == nil {
					responses <- "ok"
				} else {
					responses <- response.Error.Error()
				}
			}
		}()
	}

	for {
		select {
		case r := <-responses:
			fmt.Println(r)
		}
	}
}
